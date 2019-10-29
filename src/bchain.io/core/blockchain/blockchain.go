////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The bchain-go Authors.
//
// The bchain-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// @File: blockchain.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package blockchain

import (
	"errors"
	"fmt"
	"github.com/hashicorp/golang-lru"
	"github.com/tinylib/msgp/msgp"
	"io"
	"math/big"
	"bchain.io/common"
	"bchain.io/common/mclock"
	"bchain.io/common/types"
	"bchain.io/consensus"
	"bchain.io/core"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/state"
	"bchain.io/core/stateprocessor"
	"bchain.io/core/transaction"
	"bchain.io/params"
	"bchain.io/trie"
	"bchain.io/utils/database"
	"bchain.io/utils/event"
	"bchain.io/utils/metrics"
	"sync"
	"sync/atomic"
	"time"
)

var (
	blockInsertTimer = metrics.NewRegisteredTimer("chain/inserts", nil)

	ErrNoGenesis = errors.New("Genesis not found in chain")
)

const (
	bodyCacheLimit           = 256
	blockCacheLimit          = 256
	maxFutureBlocks          = 256
	maxTimeFutureBlocks      = 30
	badBlockLimit            = 10
	triesInMemory            = 128
	blockConsensusCacheLimit = 256

	// BlockChainVersion ensures that an incompatible database forces a resync from scratch.
	BlockChainVersion = 3
)

// CacheConfig contains the configuration values for the trie caching/pruning
// that's resident in a blockchain.
type CacheConfig struct {
	Disabled      bool          // Whether to disable trie write caching (archive node)
	TrieNodeLimit int           // Memory limit (MB) at which to flush the current in-memory trie to disk
	TrieTimeLimit time.Duration // Time limit after which to flush the current in-memory trie to disk
}

// WriteStatus status of write
type WriteStatus byte

const (
	NonStatTy WriteStatus = iota
	CanonStatTy
	SideStatTy
)

type Validator interface {
	// ValidateBody validates the given block's content.
	ValidateBody(block *block.Block) error

	// ValidateState validates the given statedb and optionally the receipts
	ValidateState(block, parent *block.Block, state *state.StateDB, receipts transaction.Receipts) error
}

// Processor is an interface for processing blocks using a given initial state.
//
// Process takes the block to be processed and the statedb upon which the
// initial state is based. It should return the receipts generated
// and return an error if any of the internal rules failed.
type Processor interface {
	Process(block *block.Block, statedb *state.StateDB, db database.IDatabaseGetter, config *params.ChainConfig) (transaction.Receipts, []*transaction.Log, error)
}

// BlockChain represents the canonical chain given a database with a genesis
// block. The Blockchain manages chain imports, reverts, chain reorganisations.
//
// Importing blocks in to the block chain happens according to the set of rules
// defined by the two stage Validator. Processing of blocks is done using the
// Processor which processes the included transaction. The validation of the state
// is done in the second part of the Validator. Failing results in aborting of
// the import.
//
// The BlockChain also helps in returning blocks from **any** chain included
// in the database as well as blocks that represents the canonical chain. It's
// important to note that GetBlock can return any block and does not need to be
// included in the canonical one where as GetBlockByNumber always represents the
// canonical chain.
type BlockChain struct {
	config *params.ChainConfig // chain & network configuration

	hc            *HeaderChain
	chainDb       database.IDatabase
	rmLogsFeed    event.Feed
	chainFeed     event.Feed
	chainSideFeed event.Feed
	chainHeadFeed event.Feed
	logsFeed      event.Feed
	scope         event.SubscriptionScope
	genesisBlock  *block.Block

	mu      sync.RWMutex // global mutex for locking chain operations
	chainmu sync.RWMutex // blockchain insertion lock
	procmu  sync.RWMutex // block processor lock

	checkpoint       int          // checkpoint counts towards the new checkpoint
	currentBlock     *block.Block // Current head of the block chain
	currentFastBlock *block.Block // Current head of the fast-sync chain (may be above the block chain!)

	stateCache        state.Database // State database to reuse between imports (contains state cache)
	bodyCache         *lru.Cache     // Cache for the most recent block bodies
	bodyMsgpCache     *lru.Cache     // Cache for the most recent block bodies in Msgp encoded format
	blockCache        *lru.Cache     // Cache for the most recent entire blocks
	futureBlocks      *lru.Cache     // future blocks are blocks added for later processing
	blkConsensusCache *lru.Cache     // Cache for block consensus data cash
	totalTxs          *lru.Cache     // Cache for block total txs

	quit    chan struct{} // blockchain quit channel
	running int32         // running must be called atomically
	// procInterrupt must be atomically called
	procInterrupt int32          // interrupt signaler for block processing
	wg            sync.WaitGroup // chain processing wait group for shutting down

	engine    consensus.Engine
	processor Processor // block processor interface
	validator Validator // block and state validator interface
	//vmConfig  vm.Config //todo for future vm

	badBlocks *lru.Cache // Bad block cache
}

// NewBlockChain returns a fully initialised block chain using information
// available in the database. It initialises the default bchain Validator and
// Processor.
func NewBlockChain(chainDb database.IDatabase, config *params.ChainConfig, engine consensus.Engine) (*BlockChain, error) {
	bodyCache, _ := lru.New(bodyCacheLimit)
	bodyMsgpCache, _ := lru.New(bodyCacheLimit)
	blockCache, _ := lru.New(blockCacheLimit)
	futureBlocks, _ := lru.New(maxFutureBlocks)
	badBlocks, _ := lru.New(badBlockLimit)
	blkConsensusData, _ := lru.New(blockConsensusCacheLimit)
	totalTxs, _ := lru.New(blockCacheLimit)

	bc := &BlockChain{
		config:            config,
		chainDb:           chainDb,
		stateCache:        state.NewDatabase(chainDb),
		quit:              make(chan struct{}),
		bodyCache:         bodyCache,
		bodyMsgpCache:     bodyMsgpCache,
		blockCache:        blockCache,
		futureBlocks:      futureBlocks,
		engine:            engine,
		badBlocks:         badBlocks,
		blkConsensusCache: blkConsensusData,
		totalTxs:          totalTxs,
	}
	bc.SetValidator(NewBlockValidator(config, bc, engine))
	bc.SetProcessor(stateprocessor.NewStateProcessor(config, bc, engine, chainDb))

	var err error
	bc.hc, err = NewHeaderChain(chainDb, config, engine, bc.getProcInterrupt)
	if err != nil {
		return nil, err
	}
	bc.genesisBlock = bc.GetBlockByNumber(0)
	if bc.genesisBlock == nil {
		return nil, ErrNoGenesis
	}
	if err := bc.loadLastState(); err != nil {
		return nil, err
	}

	// Take ownership of this particular state
	go bc.update()
	return bc, nil
}

func (bc *BlockChain) getProcInterrupt() bool {
	return atomic.LoadInt32(&bc.procInterrupt) == 1
}

// loadLastState loads the last known chain state from the database. This method
// assumes that the chain manager mutex is held.
func (bc *BlockChain) loadLastState() error {
	// Restore the last known head block
	headHash := GetHeadBlockHash(bc.chainDb)
	if headHash == (types.Hash{}) {
		// Corrupt or empty database, init from scratch
		logger.Warn("Empty database, resetting chain")
		return bc.Reset()
	}
	// Make sure the entire head block is available
	currentBlock := bc.GetBlockByHash(headHash)
	if currentBlock == nil {
		// Corrupt or empty database, init from scratch
		logger.Warn("Head block missing, resetting chain", "hash", headHash.String())
		return bc.Reset()
	}
	// Make sure the state associated with the block is available
	if _, err := state.New(currentBlock.Root(), bc.stateCache); err != nil {
		// Dangling block without a state associated, init from scratch
		logger.Warn("Head state missing, repairing chain", "number", currentBlock.Number().String(), "hash", currentBlock.Hash().String())
		return bc.Reset()
	}
	// Everything seems to be fine, set as the head block
	bc.currentBlock = currentBlock

	// Restore the last known head header
	currentHeader := bc.currentBlock.Header()
	if head := GetHeadHeaderHash(bc.chainDb); head != (types.Hash{}) {
		if header := bc.GetHeaderByHash(head); header != nil {
			currentHeader = header
		}
	}
	bc.hc.SetCurrentHeader(currentHeader)

	// Restore the last known head fast block
	bc.currentFastBlock = bc.currentBlock
	if head := GetHeadFastBlockHash(bc.chainDb); head != (types.Hash{}) {
		if block := bc.GetBlockByHash(head); block != nil {
			bc.currentFastBlock = block
		}
	}

	// Issue a status log for the user

	logger.Info("Loaded most recent local header", "number", currentHeader.Number.IntVal.String(), "hash", currentHeader.Hash().String())
	logger.Info("Loaded most recent local full block", "number", bc.currentBlock.Number().String(), "hash", bc.currentBlock.Hash().String())
	logger.Info("Loaded most recent local fast block", "number", bc.currentFastBlock.Number().String(), "hash", bc.currentFastBlock.Hash().String())

	return nil
}

// SetHead rewinds the local chain to a new head. In the case of headers, everything
// above the new head will be deleted and the new one set. In the case of blocks
// though, the head may be further rewound if block bodies are missing (non-archive
// nodes after a fast sync).
func (bc *BlockChain) SetHead(head uint64) error {
	logger.Warn("Rewinding blockchain", "target", head)

	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Rewind the header chain, deleting all block bodies until then
	delFn := func(hash types.Hash, num uint64) {
		DeleteBody(bc.chainDb, hash, num)
	}
	bc.hc.SetHead(head, delFn)
	currentHeader := bc.hc.CurrentHeader()

	// Clear out any stale content from the caches
	bc.bodyCache.Purge()
	bc.bodyMsgpCache.Purge()
	bc.blockCache.Purge()
	bc.futureBlocks.Purge()

	// Rewind the block chain, ensuring we don't end up with a stateless head block
	if bc.currentBlock != nil && currentHeader.Number.IntVal.Uint64() < bc.currentBlock.NumberU64() {
		bc.currentBlock = bc.GetBlock(currentHeader.Hash(), currentHeader.Number.IntVal.Uint64())
	}
	if bc.currentBlock != nil {
		if _, err := state.New(bc.currentBlock.Root(), bc.stateCache); err != nil {
			// Rewound state missing, rolled back to before pivot, reset to genesis
			bc.currentBlock = nil
		}
	}
	// Rewind the fast block in a simpleton way to the target head
	if bc.currentFastBlock != nil && currentHeader.Number.IntVal.Uint64() < bc.currentFastBlock.NumberU64() {
		bc.currentFastBlock = bc.GetBlock(currentHeader.Hash(), currentHeader.Number.IntVal.Uint64())
	}
	// If either blocks reached nil, reset to the genesis state
	if bc.currentBlock == nil {
		bc.currentBlock = bc.genesisBlock
	}
	if bc.currentFastBlock == nil {
		bc.currentFastBlock = bc.genesisBlock
	}
	if err := WriteHeadBlockHash(bc.chainDb, bc.currentBlock.Hash()); err != nil {
		logger.Critical("Failed to reset head full block", "err", err)
	}
	if err := WriteHeadFastBlockHash(bc.chainDb, bc.currentFastBlock.Hash()); err != nil {
		logger.Critical("Failed to reset head fast block", "err", err)
	}
	return bc.loadLastState()
}

// FastSyncCommitHead sets the current head block to the one defined by the hash
// irrelevant what the chain contents were prior.
func (bc *BlockChain) FastSyncCommitHead(hash types.Hash) error {
	// Make sure that both the block as well at its state trie exists
	block := bc.GetBlockByHash(hash)
	if block == nil {
		return fmt.Errorf("non existent block [%x…]", hash[:4])
	}
	if _, err := trie.NewSecure(block.Root(), bc.chainDb, 0); err != nil {
		return err
	}
	// If all checks out, manually set the head block
	bc.mu.Lock()
	bc.currentBlock = block
	bc.mu.Unlock()

	logger.Info("Committed new head block", "number", block.Number().String(), "hash", hash.String())
	return nil
}

// LastBlockHash return the hash of the HEAD block.
func (bc *BlockChain) LastBlockHash() types.Hash {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.currentBlock.Hash()
}

func (bc *BlockChain) GetNowBlockHash() types.Hash {
	return bc.LastBlockHash()
}

// CurrentBlock retrieves the current head block of the canonical chain. The
// block is retrieved from the blockchain's internal cache.
func (bc *BlockChain) CurrentBlock() *block.Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.currentBlock
}

func (bc *BlockChain) CurrentBlockNum() uint64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.currentBlock.Header().Number.IntVal.Uint64()
}

// CurrentFastBlock retrieves the current fast-sync head block of the canonical
// chain. The block is retrieved from the blockchain's internal cache.
func (bc *BlockChain) CurrentFastBlock() *block.Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.currentFastBlock
}

// Status returns status information about the current chain such as the HEAD Number,
// the HEAD hash and the hash of the genesis block.
func (bc *BlockChain) Status() (numnber *big.Int, currentBlock types.Hash, genesisBlock types.Hash) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.currentBlock.Number(), bc.currentBlock.Hash(), bc.genesisBlock.Hash()
}

// SetProcessor sets the processor required for making state modifications.
func (bc *BlockChain) SetProcessor(processor Processor) {
	bc.procmu.Lock()
	defer bc.procmu.Unlock()
	bc.processor = processor
}

// SetValidator sets the validator which is used to validate incoming blocks.
func (bc *BlockChain) SetValidator(validator Validator) {
	bc.procmu.Lock()
	defer bc.procmu.Unlock()
	bc.validator = validator
}

// Validator returns the current validator.
func (bc *BlockChain) Validator() Validator {
	bc.procmu.RLock()
	defer bc.procmu.RUnlock()
	return bc.validator
}

// Processor returns the current processor.
func (bc *BlockChain) Processor() Processor {
	bc.procmu.RLock()
	defer bc.procmu.RUnlock()
	return bc.processor
}

// State returns a new mutable state based on the current HEAD block.
func (bc *BlockChain) State() (*state.StateDB, error) {
	return bc.StateAt(bc.CurrentBlock().Root())
}

// StateAt returns a new mutable state based on a particular point in time.
func (bc *BlockChain) StateAt(root types.Hash) (*state.StateDB, error) {
	return state.New(root, bc.stateCache)
}

// Reset purges the entire blockchain, restoring it to its genesis state.
func (bc *BlockChain) Reset() error {
	return bc.ResetWithGenesisBlock(bc.genesisBlock)
}

// ResetWithGenesisBlock purges the entire blockchain, restoring it to the
// specified genesis state.
func (bc *BlockChain) ResetWithGenesisBlock(genesis *block.Block) error {
	// Dump the entire block chain and purge the caches
	if err := bc.SetHead(0); err != nil {
		return err
	}
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if err := WriteBlock(bc.chainDb, genesis); err != nil {
		logger.Critical("Failed to write genesis block", "err", err)
	}
	bc.genesisBlock = genesis
	bc.insert(bc.genesisBlock)
	bc.currentBlock = bc.genesisBlock
	bc.hc.SetGenesis(bc.genesisBlock.Header())
	bc.hc.SetCurrentHeader(bc.genesisBlock.Header())
	bc.currentFastBlock = bc.genesisBlock

	return nil
}

// Export writes the active chain to the given writer.
func (bc *BlockChain) Export(w io.Writer) error {
	return bc.ExportN(w, uint64(0), bc.currentBlock.NumberU64())
}

// ExportN writes a subset of the active chain to the given writer.
func (bc *BlockChain) ExportN(w io.Writer, first uint64, last uint64) error {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if first > last {
		return fmt.Errorf("export failed: first (%d) is greater than last (%d)", first, last)
	}
	logger.Info("Exporting batch of blocks", "count", last-first+1)

	for nr := first; nr <= last; nr++ {
		blk := bc.GetBlockByNumber(nr)
		if blk == nil {
			return fmt.Errorf("export failed on #%d: not found", nr)
		}
		if err := msgp.Encode(w, blk); err != nil {
			return err
		}
	}

	return nil
}

// insert injects a new head block into the current block chain. This method
// assumes that the block is indeed a true head. It will also reset the head
// header and the head fast sync block to this very same block if they are older
// or if they are on a different side chain.
//
// Note, this function assumes that the `mu` mutex is held!
func (bc *BlockChain) insert(block *block.Block) {
	// If the block is on a side chain or an unknown one, force other heads onto it too
	updateHeads := GetCanonicalHash(bc.chainDb, block.NumberU64()) != block.Hash()

	// Add the block to the canonical chain number scheme and mark as the head
	if err := WriteCanonicalHash(bc.chainDb, block.Hash(), block.NumberU64()); err != nil {
		logger.Critical("Failed to insert block number", "err", err)
	}
	if err := WriteHeadBlockHash(bc.chainDb, block.Hash()); err != nil {
		logger.Critical("Failed to insert head block hash", "err", err)
	}
	bc.currentBlock = block

	// If the block is better than our head or is on a different chain, force update heads
	if updateHeads {
		bc.hc.SetCurrentHeader(block.Header())

		if err := WriteHeadFastBlockHash(bc.chainDb, block.Hash()); err != nil {
			logger.Critical("Failed to insert head fast block hash", "err", err)
		}
		bc.currentFastBlock = block
	}
}

// Genesis retrieves the chain's genesis block.
func (bc *BlockChain) Genesis() *block.Block {
	return bc.genesisBlock
}

// GetBody retrieves a block body from the database by
// hash, caching it if found.
func (bc *BlockChain) GetBody(hash types.Hash) *block.Body {
	// Short circuit if the body's already in the cache, retrieve otherwise
	if cached, ok := bc.bodyCache.Get(hash); ok {
		body := cached.(*block.Body)
		return body
	}
	body := GetBody(bc.chainDb, hash, bc.hc.GetBlockNumber(hash))
	if body == nil {
		return nil
	}
	// Cache the found body for next time and return
	bc.bodyCache.Add(hash, body)
	return body
}

func (bc *BlockChain) GetBodyMsgp(hash types.Hash) []byte {
	// Short circuit if the body's already in the cache, retrieve otherwise
	if cached, ok := bc.bodyMsgpCache.Get(hash); ok {
		return cached.([]byte)
	}
	body := GetBodyMsgp(bc.chainDb, hash, bc.hc.GetBlockNumber(hash))
	if len(body) == 0 {
		return nil
	}
	// Cache the found body for next time and return
	bc.bodyMsgpCache.Add(hash, body)
	return body
}

// HasBlock checks if a block is fully present in the database or not.
func (bc *BlockChain) HasBlock(hash types.Hash, number uint64) bool {
	if bc.blockCache.Contains(hash) {
		return true
	}
	ok, _ := bc.chainDb.Has(blockBodyKey(hash, number))
	return ok
}

// HasState checks if state trie is fully present in the database or not.
func (bc *BlockChain) HasState(hash types.Hash) bool {
	_, err := bc.stateCache.OpenTrie(hash)
	return err == nil
}

// HasBlockAndState checks if a block and associated state trie is fully present
// in the database or not, caching it if present.
func (bc *BlockChain) HasBlockAndState(hash types.Hash) bool {
	// Check first that the block itself is known
	block := bc.GetBlockByHash(hash)
	if block == nil {
		return false
	}
	return bc.HasState(block.Root())
}

// GetBlock retrieves a block from the database by hash and number,
// caching it if found.
func (bc *BlockChain) GetBlock(hash types.Hash, number uint64) *block.Block {
	// Short circuit if the block's already in the cache, retrieve otherwise
	if blk, ok := bc.blockCache.Get(hash); ok {
		return blk.(*block.Block)
	}
	blk := GetBlock(bc.chainDb, hash, number)
	if blk == nil {
		return nil
	}
	// Cache the found block for next time and return
	bc.blockCache.Add(blk.Hash(), blk)
	return blk
}

// GetBlockByHash retrieves a block from the database by hash, caching it if found.
func (bc *BlockChain) GetBlockByHash(hash types.Hash) *block.Block {
	return bc.GetBlock(hash, bc.hc.GetBlockNumber(hash))
}

// GetBlockByNumber retrieves a block from the database by number, caching it
// (associated with its hash) if found.
func (bc *BlockChain) GetBlockByNumber(number uint64) *block.Block {
	hash := GetCanonicalHash(bc.chainDb, number)
	if hash == (types.Hash{}) {
		return nil
	}
	return bc.GetBlock(hash, number)
}

func (bc *BlockChain) VerifyNextRoundBlock(block *block.Block) bool {

	currentBlock := bc.CurrentBlock()
	//check hash
	if currentBlock.Header().Hash() != block.H.ParentHash {
		logger.Error(" parent hash not match", "self", currentBlock.Header().Hash().String(), "peer", block.H.ParentHash.String())
		return false
	}
	//check num
	if currentBlock.H.Number.IntVal.Int64()+1 != block.H.Number.IntVal.Int64() {
		logger.Error("currentBlock.B_header.Number.IntVal.Int64() + 1 != block.B_header.Number.IntVal.Int64()")
		return false
	}

	//verify header
	var err error
	if err = bc.engine.VerifyHeader(bc, block.Header(), false); err != nil {
		logger.Error("err = bc.engine.VerifyHeader : ", err.Error())
		return false
	}
	//verify body
	if err = bc.Validator().ValidateBody(block); err != nil {
		logger.Error("err = bc.Validator().ValidateBody(block) :", err.Error())
		return false
	}

	return true

}

// GetReceiptsByHash retrieves the receipts for all transactions in a given block.
func (bc *BlockChain) GetReceiptsByHash(hash types.Hash) transaction.Receipts {
	return GetBlockReceipts(bc.chainDb, hash, GetBlockNumber(bc.chainDb, hash))
}

// GetBlocksFromHash returns the block corresponding to hash and up to n-1 ancestors.
func (bc *BlockChain) GetBlocksFromHash(hash types.Hash, n int) (blocks []*block.Block) {
	number := bc.hc.GetBlockNumber(hash)
	for i := 0; i < n; i++ {
		block := bc.GetBlock(hash, number)
		if block == nil {
			break
		}
		blocks = append(blocks, block)
		hash = block.ParentHash()
		number--
	}
	return
}

// Stop stops the blockchain service. If any imports are currently in progress
// it will abort them using the procInterrupt.
func (bc *BlockChain) Stop() {
	if !atomic.CompareAndSwapInt32(&bc.running, 0, 1) {
		return
	}
	// Unsubscribe all subscriptions registered from blockchain
	bc.scope.Close()
	close(bc.quit)
	atomic.StoreInt32(&bc.procInterrupt, 1)

	bc.wg.Wait()
	//logger.Info("Blockchain manager stopped")
}

func (bc *BlockChain) procFutureBlocks() {
	blocks := make([]*block.Block, 0, bc.futureBlocks.Len())
	for _, hash := range bc.futureBlocks.Keys() {
		if blk, exist := bc.futureBlocks.Peek(hash); exist {
			blocks = append(blocks, blk.(*block.Block))
		}
	}
	if len(blocks) > 0 {
		block.BlockBy(block.Number).Sort(blocks)

		// Insert one by one as chain insertion needs contiguous ancestry between blocks
		for i := range blocks {
			bc.InsertChain(blocks[i : i+1])
		}
	}
}

// Rollback is designed to remove a chain of links from the database that aren't
// certain enough to be valid.
func (bc *BlockChain) Rollback(chain []types.Hash) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	for i := len(chain) - 1; i >= 0; i-- {
		hash := chain[i]

		currentHeader := bc.hc.CurrentHeader()
		if currentHeader.Hash() == hash {
			bc.hc.SetCurrentHeader(bc.GetHeader(currentHeader.ParentHash, currentHeader.Number.IntVal.Uint64()-1))
		}
		if bc.currentFastBlock.Hash() == hash {
			bc.currentFastBlock = bc.GetBlock(bc.currentFastBlock.ParentHash(), bc.currentFastBlock.NumberU64()-1)
			WriteHeadFastBlockHash(bc.chainDb, bc.currentFastBlock.Hash())
		}
		if bc.currentBlock.Hash() == hash {
			bc.currentBlock = bc.GetBlock(bc.currentBlock.ParentHash(), bc.currentBlock.NumberU64()-1)
			WriteHeadBlockHash(bc.chainDb, bc.currentBlock.Hash())
		}
	}
}

// SetReceiptsData computes all the non-consensus fields of the receipts
func SetReceiptsData(config *params.ChainConfig, block *block.Block, receipts transaction.Receipts) {
	//signer := transaction.MakeSigner(config, block.Number())

	transactions, logIndex := block.Transactions(), uint(0)

	// todo: need modify when len of receipts is not equal len of transactions
	// todo eg: for j := 1; j < len(receipts); j++ {
	for j := 0; j < len(receipts); j++ {
		// The transaction hash can be retrieved from the transaction itself
		receipts[j].TxHash = transactions[j].Hash()

		// todo The contract address can be derived from the transaction itself
		//if len(transactions[j].Data.Acts)==2 && transactions[j].Data.Acts[1].Contract == nil {
		//	// Deriving the signer is expensive, only do if it's actually needed
		//	from, _ := transaction.Sender(signer, transactions[j])
		//	receipts[j].ContractAddress = crypto.CreateAddress(from, transactions[j].Nonce())
		//}

		// The derived log fields can simply be set from the block and transaction
		for k := 0; k < len(receipts[j].Logs); k++ {
			receipts[j].Logs[k].BlockNumber = block.NumberU64()
			receipts[j].Logs[k].BlockHash = block.Hash()
			receipts[j].Logs[k].TxHash = receipts[j].TxHash
			receipts[j].Logs[k].TxIndex = uint(j)
			receipts[j].Logs[k].Index = logIndex
			logIndex++
		}
	}
}

const IdealBatchSize = 100 * 1024

// InsertReceiptChain attempts to complete an already existing header chain with
// transaction and receipt data.
func (bc *BlockChain) InsertReceiptChain(blockChain block.Blocks, receiptChain []transaction.Receipts) (int, error) {
	bc.wg.Add(1)
	defer bc.wg.Done()

	// Do a sanity check that the provided chain is actually ordered and linked
	for i := 1; i < len(blockChain); i++ {
		if blockChain[i].NumberU64() != blockChain[i-1].NumberU64()+1 || blockChain[i].ParentHash() != blockChain[i-1].Hash() {
			logger.Error("Non contiguous receipt insert", "number", blockChain[i].Number().String(), "hash", blockChain[i].Hash().String(), "parent", blockChain[i].ParentHash().String(),
				"prevnumber", blockChain[i-1].Number().String(), "prevhash", blockChain[i-1].Hash().String())
			return 0, fmt.Errorf("non contiguous insert: item %d is #%d [%x…], item %d is #%d [%x…] (parent [%x…])", i-1, blockChain[i-1].NumberU64(),
				blockChain[i-1].Hash().Bytes()[:4], i, blockChain[i].NumberU64(), blockChain[i].Hash().Bytes()[:4], blockChain[i].ParentHash().Bytes()[:4])
		}
	}

	var (
		stats = struct{ processed, ignored int32 }{}
		start = time.Now()
		bytes = 0
		batch = bc.chainDb.NewBatch()
	)
	for i, block := range blockChain {
		receipts := receiptChain[i]
		// Short circuit insertion if shutting down or processing failed
		if atomic.LoadInt32(&bc.procInterrupt) == 1 {
			return 0, nil
		}
		// Short circuit if the owner header is unknown
		if !bc.HasHeader(block.Hash(), block.NumberU64()) {
			return i, fmt.Errorf("containing header #%d [%x…] unknown", block.Number(), block.Hash().Bytes()[:4])
		}
		// Skip if the entire data is already known
		if bc.HasBlock(block.Hash(), block.NumberU64()) {
			stats.ignored++
			continue
		}
		// Compute all the non-consensus fields of the receipts
		SetReceiptsData(bc.config, block, receipts)
		// Write all the data out into the database
		if err := WriteBody(batch, block.Hash(), block.NumberU64(), block.Body()); err != nil {
			return i, fmt.Errorf("failed to write block body: %v", err)
		}
		if err := WriteBlockReceipts(batch, block.Hash(), block.NumberU64(), receipts); err != nil {
			return i, fmt.Errorf("failed to write block receipts: %v", err)
		}
		if err := WriteTxLookupEntries(batch, block); err != nil {
			return i, fmt.Errorf("failed to write lookup metadata: %v", err)
		}
		stats.processed++

		if batch.ValueSize() >= IdealBatchSize {
			if err := batch.Write(); err != nil {
				return 0, err
			}
			bytes += batch.ValueSize()
			batch.Reset()
		}
	}
	if batch.ValueSize() > 0 {
		bytes += batch.ValueSize()
		if err := batch.Write(); err != nil {
			return 0, err
		}
	}

	// Update the head fast sync block if better
	bc.mu.Lock()
	head := blockChain[len(blockChain)-1]
	if head.NumberU64() > bc.currentFastBlock.Number().Uint64() {
		if err := WriteHeadFastBlockHash(bc.chainDb, head.Hash()); err != nil {
			logger.Critical("Failed to update head fast block hash", "err", err)
		}
		bc.currentFastBlock = head
	}
	bc.mu.Unlock()

	logger.Info("Imported new block receipts",
		"count", stats.processed,
		"elapsed", common.PrettyDuration(time.Since(start)),
		"bytes", bytes,
		"number", head.Number().String(),
		"hash", head.Hash().String(),
		"ignored", stats.ignored)
	return 0, nil
}

func (bc *BlockChain) WriteConsensusData(key types.Hash, value []byte) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	batch := bc.chainDb.NewBatch()

	dbKey := append(consensusPrefix, key.Bytes()...)
	if err := batch.Put(dbKey, value); err != nil {
		logger.Critical("WriteConsensusData, failed to store value", "err", err)
		return err
	}

	if err := batch.Write(); err != nil {
		logger.Critical("WriteConsensusData, batch write fail", "err", err)
		return err
	}
	// Cache
	bc.blkConsensusCache.Add(key, value)
	logger.Debug("WriteConsensusData OK! block hash", key.String())
	return nil
}

func (bc *BlockChain) GetConsensusData(key types.Hash) []byte {
	if data, ok := bc.blkConsensusCache.Get(key); ok {
		return data.([]byte)
	}
	dbKey := append(consensusPrefix, key.Bytes()...)
	data, _ := bc.chainDb.Get(dbKey)
	if len(data) == 0 {
		return nil
	}

	// Cache
	bc.blkConsensusCache.Add(key, data)
	return data
}

func (bc *BlockChain) GetExtra(key []byte) []byte {
	return GetExtra(bc.chainDb, key)
}

func (bc *BlockChain) GetBlockStat(hash types.Hash) *BlockStat {
	if data, ok := bc.totalTxs.Get(hash); ok {
		return data.(*BlockStat)
	}
	stat := GetBlockStat(bc.chainDb, hash)
	if stat == nil {
		return nil
	}
	// Cache
	bc.totalTxs.Add(hash, stat)
	return stat
}

func (bc *BlockChain) writeBlockStat(hash types.Hash, stat *BlockStat) error {
	if err := WriteBlockStat(bc.chainDb, hash, stat); err != nil {
		return err
	}
	bc.totalTxs.Add(hash, stat)
	return nil
}

// WriteBlock writes the block to the chain.
func (bc *BlockChain) WriteBlockAndState(block *block.Block, receipts []*transaction.Receipt, state *state.StateDB) (status WriteStatus, err error) {
	bc.wg.Add(1)
	defer bc.wg.Done()

	// Make sure no inconsistent state is leaked during insertion
	bc.mu.Lock()
	defer bc.mu.Unlock()

	bcCurrentBlockNumber := bc.currentBlock.NumberU64()

	blockNumber := block.Number().Uint64()

	pstat := bc.GetBlockStat(block.ParentHash())
	if pstat == nil {
		return NonStatTy, consensus.ErrUnknownAncestor
	}
	txlen := big.NewInt(int64(len(block.Transactions())))

	// Write other block data using a batch.
	batch := bc.chainDb.NewBatch()
	if err := WriteBlock(batch, block); err != nil {
		return NonStatTy, err
	}

	_, si, err := state.CommitTo(batch, true)
	if err != nil {
		return NonStatTy, err
	}
	if err := WriteBlockReceipts(batch, block.Hash(), block.NumberU64(), receipts); err != nil {
		return NonStatTy, err
	}

	// Write hash preimages
	if err := WritePreimages(bc.chainDb, block.NumberU64(), state.Preimages()); err != nil {
		return NonStatTy, err
	}

	//if new bigger block number coming, need reorg the chain
	reorg := blockNumber > bcCurrentBlockNumber

	// for apos, rollback Canon chain
	// empty block time is equal to parent time. choose normal block
	// if two block time is same, choose bigger hash block
	if (blockNumber == bcCurrentBlockNumber) {
		cmpRet := bc.currentBlock.Time().Cmp(block.Time())
		if cmpRet < 0 {
			reorg = true
		} else {
			a := new(big.Int).SetBytes(block.Hash().Bytes())
			b := new(big.Int).SetBytes(bc.currentBlock.Hash().Bytes())
			hashCmpRet := a.Cmp(b)
			if hashCmpRet > 0 {
				reorg = true
			}
		}
	}

	if reorg {
		// Reorganise the chain if the parent is not the head block
		if block.ParentHash() != bc.currentBlock.Hash() {
			if err := bc.reorg(bc.currentBlock, block); err != nil {
				return NonStatTy, err
			}
		}
		// Write the positional metadata for transaction and receipt lookups
		if err := WriteTxLookupEntries(batch, block); err != nil {
			return NonStatTy, err
		}
		status = CanonStatTy
	} else {
		status = SideStatTy
	}

	stat := &BlockStat{
		Ttxs:        types.NewBigInt(*new(big.Int).Add(txlen, &pstat.Ttxs.IntVal)),
		TsoContract: types.NewBigInt(*new(big.Int).Add(big.NewInt(int64(si.TnewsoContract)), &pstat.TsoContract.IntVal)),
		TsoNormal:   types.NewBigInt(*new(big.Int).Add(big.NewInt(int64(si.TnewsoNormal)), &pstat.TsoNormal.IntVal)),
		TstateNum:   types.NewBigInt(*new(big.Int).Add(big.NewInt(int64(si.TnewState)), &pstat.TstateNum.IntVal)),
	}
	if err := bc.writeBlockStat(block.Hash(), stat); err != nil {
		return NonStatTy, err
	}

	if err := batch.Write(); err != nil {
		return NonStatTy, err
	}

	// Set new head.
	if status == CanonStatTy {
		bc.insert(block)
	}
	bc.futureBlocks.Remove(block.Hash())
	return status, nil
}

// InsertChain attempts to insert the given batch of blocks in to the canonical
// chain or, otherwise, create a fork. If an error is returned it will return
// the index number of the failing block as well an error describing what went
// wrong.
//
// After insertion is done, all accumulated events will be fired.
func (bc *BlockChain) InsertChain(chain block.Blocks) (int, error) {
	n, events, logs, err := bc.insertChain(chain)
	bc.PostChainEvents(events, logs)
	return n, err
}

// insertChain will execute the actual chain insertion and event aggregation. The
// only reason this method exists as a separate one is to make locking cleaner
// with deferred statements.
func (bc *BlockChain) insertChain(chain block.Blocks) (int, []interface{}, []*transaction.Log, error) {
	logger.Info("insertChain in ....")
	// Do a sanity check that the provided chain is actually ordered and linked
	for i := 1; i < len(chain); i++ {
		if chain[i].NumberU64() != chain[i-1].NumberU64()+1 || chain[i].ParentHash() != chain[i-1].Hash() {
			// Chain broke ancestry, log a messge (programming error) and skip insertion
			logger.Error("Non contiguous block insert", "number", chain[i].Number().String(), "hash", chain[i].Hash().String(),
				"parent", chain[i].ParentHash().String(), "prevnumber", chain[i-1].Number().String(), "prevhash", chain[i-1].Hash().String())

			return 0, nil, nil, fmt.Errorf("non contiguous insert: item %d is #%d [%x…], item %d is #%d [%x…] (parent [%x…])", i-1, chain[i-1].NumberU64(),
				chain[i-1].Hash().Bytes()[:4], i, chain[i].NumberU64(), chain[i].Hash().Bytes()[:4], chain[i].ParentHash().Bytes()[:4])
		}
	}
	// Pre-checks passed, start the full block imports
	bc.wg.Add(1)
	defer bc.wg.Done()

	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	// A queued approach to delivering events. This is generally
	// faster than direct delivery and requires much less mutex
	// acquiring.
	var (
		stats         = insertStats{startTime: mclock.Now()}
		events        = make([]interface{}, 0, len(chain))
		lastCanon     *block.Block
		coalescedLogs []*transaction.Log
	)
	// Start the parallel header verifier
	headers := make([]*block.Header, len(chain))
	seals := make([]bool, len(chain))

	for i, blk := range chain {
		headers[i] = blk.Header()
		seals[i] = true
	}
	abort, results := bc.engine.VerifyHeaders(bc, headers, seals)
	logger.Info("VerifyHeaders return channel ....")
	defer close(abort)

	// Iterate over the blocks and insert when the verifier permits
	for i, blk := range chain {
		// If the chain is terminating, stop processing blocks
		if atomic.LoadInt32(&bc.procInterrupt) == 1 {
			logger.Debug("Premature abort during blocks processing")
			break
		}

		// Wait for the block's verification to complete
		bstart := time.Now()

		err := <-results
		if err == nil {
			err = bc.Validator().ValidateBody(blk)
		}
		if err != nil {
			if err == core.ErrKnownBlock {
				stats.ignored++
				continue
			}

			if err == consensus.ErrFutureBlock {
				// Allow up to MaxFuture second in the future blocks. If this limit
				// is exceeded the chain is discarded and processed at a later time
				// if given.
				max := big.NewInt(time.Now().Unix() + maxTimeFutureBlocks)
				if blk.Time().Cmp(max) > 0 {
					return i, events, coalescedLogs, fmt.Errorf("future block: %v > %v", blk.Time(), max)
				}
				bc.futureBlocks.Add(blk.Hash(), blk)
				stats.queued++
				continue
			}

			if err == consensus.ErrUnknownAncestor && bc.futureBlocks.Contains(blk.ParentHash()) {
				bc.futureBlocks.Add(blk.Hash(), blk)
				stats.queued++
				continue
			}

			bc.reportBlock(blk, nil, err)
			return i, events, coalescedLogs, err
		}
		// Create a new statedb using the parent block and report an
		// error if it fails.
		var parent *block.Block
		if i == 0 {
			parent = bc.GetBlock(blk.ParentHash(), blk.NumberU64()-1)
		} else {
			parent = chain[i-1]
		}
		state, err := state.New(parent.Root(), bc.stateCache)
		if err != nil {
			return i, events, coalescedLogs, err
		}
		// Process block using the parent state as reference point.
		logger.Info(">>>InsertChain will process transactions.........")
		receipts, logs, err := bc.processor.Process(blk, state, bc.chainDb, bc.Config())
		if err != nil {
			bc.reportBlock(blk, receipts, err)
			return i, events, coalescedLogs, err
		}
		// Validate the state using the default validator
		err = bc.Validator().ValidateState(blk, parent, state, receipts)
		if err != nil {
			bc.reportBlock(blk, receipts, err)
			return i, events, coalescedLogs, err
		}

		// Write the block to the chain and get the status.
		status, err := bc.WriteBlockAndState(blk, receipts, state)
		if err != nil {
			return i, events, coalescedLogs, err
		}

		switch status {
		case CanonStatTy:
			//logger.Debug("Inserted new block", "number", blk.Number().String(), "hash", blk.Hash().String(),
			//	"txs", len(blk.Transactions()), "elapsed", common.PrettyDuration(time.Since(bstart)))
			logger.Infof("\033[31m Inserted new block number:%d \033[0m  hash:%s  len(txs):%d elapsed:%d", blk.Number().Int64(), blk.Hash().String(), len(blk.Transactions()), common.PrettyDuration(time.Since(bstart)))
			coalescedLogs = append(coalescedLogs, logs...)
			blockInsertTimer.UpdateSince(bstart)
			events = append(events, core.ChainEvent{blk, blk.Hash(), logs})
			lastCanon = blk

		case SideStatTy:
			logger.Debug("Inserted forked block", "number", blk.Number().String(), "hash", blk.Hash().String(), "elapsed",
				common.PrettyDuration(time.Since(bstart)), "txs", len(blk.Transactions()))

			blockInsertTimer.UpdateSince(bstart)
			events = append(events, core.ChainSideEvent{blk})
		}
		stats.processed++
		stats.report(chain, i)
	}
	// Append a single chain head event if we've progressed the chain
	if lastCanon != nil && bc.CurrentBlock().Hash() == lastCanon.Hash() {
		events = append(events, core.ChainHeadEvent{lastCanon})
	}
	return 0, events, coalescedLogs, nil
}

// insertStats tracks and reports on block insertion.
type insertStats struct {
	queued, processed, ignored int
	lastIndex                  int
	startTime                  mclock.AbsTime
}

// statsReportLimit is the time limit during import after which we always print
// out progress. This avoids the user wondering what's going on.
const statsReportLimit = 8 * time.Second

// report prints statistics if some number of blocks have been processed
// or more than a few seconds have passed since the last message.
func (st *insertStats) report(chain []*block.Block, index int) {
	// Fetch the timings for the batch
	var (
		now     = mclock.Now()
		elapsed = time.Duration(now) - time.Duration(st.startTime)
	)
	// If we're at the last block of the batch or report period reached, log
	if index == len(chain)-1 || elapsed >= statsReportLimit {
		var (
			end = chain[index]
			txs = countTransactions(chain[st.lastIndex : index+1])
		)
		context := []interface{}{
			"blocks", st.processed,
			"txs", txs,
			"elapsed", common.PrettyDuration(elapsed),
			"number", end.Number().String(),
			"hash", end.Hash().String(),
		}
		if st.queued > 0 {
			context = append(context, []interface{}{"queued", st.queued}...)
		}
		if st.ignored > 0 {
			context = append(context, []interface{}{"ignored", st.ignored}...)
		}
		logger.Info("Imported new chain segment")
		logger.Info(context...)

		*st = insertStats{startTime: now, lastIndex: index + 1}
	}
}

func countTransactions(chain []*block.Block) (c int) {
	for _, b := range chain {
		c += len(b.Transactions())
	}
	return c
}

// reorgs takes two blocks, an old chain and a new chain and will reconstruct the blocks and inserts them
// to be part of the new canonical chain and accumulates potential missing transactions and post an
// event about them
func (bc *BlockChain) reorg(oldBlock, newBlock *block.Block) error {
	var (
		newChain    block.Blocks
		oldChain    block.Blocks
		commonBlock *block.Block
		deletedTxs  transaction.Transactions
		deletedLogs []*transaction.Log
		// collectLogs collects the logs that were generated during the
		// processing of the block that corresponds with the given hash.
		// These logs are later announced as deleted.
		collectLogs = func(h types.Hash) {
			// Coalesce logs and set 'Removed'.
			receipts := GetBlockReceipts(bc.chainDb, h, bc.hc.GetBlockNumber(h))
			for _, receipt := range receipts {
				for _, log := range receipt.Logs {
					del := *log
					del.Removed = true
					deletedLogs = append(deletedLogs, &del)
				}
			}
		}
	)

	// reduce new chain and append new chain blocks for inserting later on
	for ; newBlock != nil && newBlock.NumberU64() != oldBlock.NumberU64(); newBlock = bc.GetBlock(newBlock.ParentHash(), newBlock.NumberU64()-1) {
		newChain = append(newChain, newBlock)
	}
	if oldBlock == nil {
		return fmt.Errorf("Invalid old chain")
	}
	if newBlock == nil {
		return fmt.Errorf("Invalid new chain")
	}

	for {
		if oldBlock.Hash() == newBlock.Hash() {
			commonBlock = oldBlock
			break
		}

		oldChain = append(oldChain, oldBlock)
		newChain = append(newChain, newBlock)
		deletedTxs = append(deletedTxs, oldBlock.Transactions()...)
		collectLogs(oldBlock.Hash())

		oldBlock, newBlock = bc.GetBlock(oldBlock.ParentHash(), oldBlock.NumberU64()-1), bc.GetBlock(newBlock.ParentHash(), newBlock.NumberU64()-1)
		if oldBlock == nil {
			return fmt.Errorf("Invalid old chain")
		}
		if newBlock == nil {
			return fmt.Errorf("Invalid new chain")
		}
	}
	// Ensure the user sees large reorgs
	if len(oldChain) > 0 && len(newChain) > 0 {
		logFn := logger.Debug
		if len(oldChain) > 63 {
			logFn = logger.Warn
		}
		logFn("Chain split detected", "number", commonBlock.Number(), "hash", commonBlock.Hash(),
			"drop", len(oldChain), "dropfrom", oldChain[0].Hash(), "add", len(newChain), "addfrom", newChain[0].Hash())
	} else {
		logger.Error("Impossible reorg, please file an issue", "oldnum", oldBlock.Number().String(), "oldhash", oldBlock.Hash().String(), "newnum", newBlock.Number().String(), "newhash", newBlock.Hash().String())
	}
	// Insert the new chain, taking care of the proper incremental order
	var addedTxs transaction.Transactions
	for i := len(newChain) - 1; i >= 0; i-- {
		// insert the block in the canonical way, re-writing history
		bc.insert(newChain[i])
		// write lookup entries for hash based transaction/receipt searches
		if err := WriteTxLookupEntries(bc.chainDb, newChain[i]); err != nil {
			return err
		}
		addedTxs = append(addedTxs, newChain[i].Transactions()...)
	}
	// calculate the difference between deleted and added transactions
	diff := transaction.TxDifference(deletedTxs, addedTxs)
	// When transactions get deleted from the database that means the
	// receipts that were created in the fork must also be deleted
	for _, tx := range diff {
		DeleteTxLookupEntry(bc.chainDb, tx.Hash())
		DeleteTxAddrNonceLookupEntry(bc.chainDb, tx)
	}
	if len(deletedLogs) > 0 {
		go bc.rmLogsFeed.Send(core.RemovedLogsEvent{deletedLogs})
	}
	if len(oldChain) > 0 {
		go func() {
			for _, block := range oldChain {
				bc.chainSideFeed.Send(core.ChainSideEvent{Block: block})
			}
		}()
	}

	return nil
}

// PostChainEvents iterates over the events generated by a chain insertion and
// posts them into the event feed.
// TODO: Should not expose PostChainEvents. The chain events should be posted in WriteBlock.
func (bc *BlockChain) PostChainEvents(events []interface{}, logs []*transaction.Log) {
	// post event logs for further processing
	if logs != nil {
		bc.logsFeed.Send(logs)
	}
	for _, event := range events {
		switch ev := event.(type) {
		case core.ChainEvent:
			bc.chainFeed.Send(ev)

		case core.ChainHeadEvent:
			bc.chainHeadFeed.Send(ev)

		case core.ChainSideEvent:
			bc.chainSideFeed.Send(ev)
		}
	}
}

func (bc *BlockChain) update() {
	futureTimer := time.NewTicker(5 * time.Second)
	defer futureTimer.Stop()
	for {
		select {
		case <-futureTimer.C:
			bc.procFutureBlocks()
		case <-bc.quit:
			return
		}
	}
}

// BadBlockArgs represents the entries in the list returned when bad blocks are queried.
type BadBlockArgs struct {
	Hash   types.Hash    `msg:"hash"`
	Header *block.Header `msg:"header"`
}

// BadBlocks returns a list of the last 'bad blocks' that the client has seen on the network
func (bc *BlockChain) BadBlocks() ([]BadBlockArgs, error) {
	headers := make([]BadBlockArgs, 0, bc.badBlocks.Len())
	for _, hash := range bc.badBlocks.Keys() {
		if hdr, exist := bc.badBlocks.Peek(hash); exist {
			header := hdr.(*block.Header)
			headers = append(headers, BadBlockArgs{header.Hash(), header})
		}
	}
	return headers, nil
}

// addBadBlock adds a bad block to the bad-block LRU cache
func (bc *BlockChain) addBadBlock(block *block.Block) {
	bc.badBlocks.Add(block.Header().Hash(), block.Header())
}

// reportBlock logs a bad block error.
func (bc *BlockChain) reportBlock(block *block.Block, receipts transaction.Receipts, err error) {
	bc.addBadBlock(block)

	var receiptString string
	for _, receipt := range receipts {
		receiptString += fmt.Sprintf("\t%v\n", receipt)
	}
	logger.Error(fmt.Sprintf(`
########## BAD BLOCK #########
Chain config: %v

Number: %v
Hash: 0x%x
%v

Error: %v
##############################
`, bc.config, block.Number().String(), block.Hash(), receiptString, err))
}

// InsertHeaderChain attempts to insert the given header chain in to the local
// chain, possibly creating a reorg. If an error is returned, it will return the
// index number of the failing header as well an error describing what went wrong.
//
// The verify parameter can be used to fine tune whether nonce verification
// should be done or not. The reason behind the optional check is because some
// of the header retrieval mechanisms already need to verify nonces, as well as
// because nonces can be verified sparsely, not needing to check each.
func (bc *BlockChain) InsertHeaderChain(chain []*block.Header, checkFreq int) (int, error) {
	start := time.Now()
	if i, err := bc.hc.ValidateHeaderChain(chain, checkFreq); err != nil {
		return i, err
	}

	// Make sure only one thread manipulates the chain at once
	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	bc.wg.Add(1)
	defer bc.wg.Done()

	whFunc := func(header *block.Header) error {
		bc.mu.Lock()
		defer bc.mu.Unlock()

		_, err := bc.hc.WriteHeader(header)
		return err
	}

	return bc.hc.InsertHeaderChain(chain, whFunc, start)
}

// writeHeader writes a header into the local chain, given that its parent is
// already known. If the total difficulty of the newly inserted header becomes
// greater than the current known TD, the canonical chain is re-routed.
//
// Note: This method is not concurrent-safe with inserting blocks simultaneously
// into the chain, as side effects caused by reorganisations cannot be emulated
// without the real blocks. Hence, writing headers directly should only be done
// in two scenarios: pure-header mode of operation (light clients), or properly
// separated header/block phases (non-archive clients).
func (bc *BlockChain) writeHeader(header *block.Header) error {
	bc.wg.Add(1)
	defer bc.wg.Done()

	bc.mu.Lock()
	defer bc.mu.Unlock()

	_, err := bc.hc.WriteHeader(header)
	return err
}

// CurrentHeader retrieves the current head header of the canonical chain. The
// header is retrieved from the HeaderChain's internal cache.
func (bc *BlockChain) CurrentHeader() *block.Header {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.hc.CurrentHeader()
}

// GetHeader retrieves a block header from the database by hash and number,
// caching it if found.
func (bc *BlockChain) GetHeader(hash types.Hash, number uint64) *block.Header {
	return bc.hc.GetHeader(hash, number)
}

// GetHeaderByHash retrieves a block header from the database by hash, caching it if
// found.
func (bc *BlockChain) GetHeaderByHash(hash types.Hash) *block.Header {
	return bc.hc.GetHeaderByHash(hash)
}

// HasHeader checks if a block header is present in the database or not, caching
// it if present.
func (bc *BlockChain) HasHeader(hash types.Hash, number uint64) bool {
	return bc.hc.HasHeader(hash, number)
}

// GetBlockHashesFromHash retrieves a number of block hashes starting at a given
// hash, fetching towards the genesis block.
func (bc *BlockChain) GetBlockHashesFromHash(hash types.Hash, max uint64) []types.Hash {
	return bc.hc.GetBlockHashesFromHash(hash, max)
}

// GetHeaderByNumber retrieves a block header from the database by number,
// caching it (associated with its hash) if found.
func (bc *BlockChain) GetHeaderByNumber(number uint64) *block.Header {
	return bc.hc.GetHeaderByNumber(number)
}

// Config retrieves the blockchain's chain configuration.
func (bc *BlockChain) Config() *params.ChainConfig { return bc.config }

// Engine retrieves the blockchain's consensus engine.
func (bc *BlockChain) Engine() consensus.Engine { return bc.engine }

func (bc *BlockChain) Validate() Validator { return bc.validator }

func (bc *BlockChain) StateCache() state.Database { return bc.stateCache }

func (bc *BlockChain) ReportBlock(block *block.Block, receipts transaction.Receipts, err error) {
	bc.reportBlock(block, receipts, err)
}

func (bc *BlockChain) MuLock()   { bc.mu.Lock() }
func (bc *BlockChain) MuUnLock() { bc.mu.Unlock() }

func (bc *BlockChain) GetDb() database.IDatabase { return bc.chainDb }

func (bc *BlockChain) Test_insert(block *block.Block) { bc.insert(block) }

func (bc *BlockChain) GetGenesis() *block.Block { return bc.genesisBlock }

// SubscribeRemovedLogsEvent registers a subscription of RemovedLogsEvent.
func (bc *BlockChain) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return bc.scope.Track(bc.rmLogsFeed.Subscribe(ch))
}

// SubscribeChainEvent registers a subscription of ChainEvent.
func (bc *BlockChain) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return bc.scope.Track(bc.chainFeed.Subscribe(ch))
}

// SubscribeChainHeadEvent registers a subscription of ChainHeadEvent.
func (bc *BlockChain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return bc.scope.Track(bc.chainHeadFeed.Subscribe(ch))
}

// SubscribeChainSideEvent registers a subscription of ChainSideEvent.
func (bc *BlockChain) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return bc.scope.Track(bc.chainSideFeed.Subscribe(ch))
}

// SubscribeLogsEvent registers a subscription of []*types.Log.
func (bc *BlockChain) SubscribeLogsEvent(ch chan<- []*transaction.Log) event.Subscription {
	return bc.scope.Track(bc.logsFeed.Subscribe(ch))
}
