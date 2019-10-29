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
// @File: chain_makers.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package chainmaker

import (
	"fmt"
	"math/big"
	"bchain.io/common"
	"bchain.io/common/types"
	"bchain.io/consensus"
	"bchain.io/core/actioncontext"
	"bchain.io/core/blockchain"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/genesis"
	"bchain.io/core/state"
	"bchain.io/core/stateprocessor"
	"bchain.io/core/transaction"
	"bchain.io/params"
	"bchain.io/utils/database"
)

// So we can deterministically seed different blockchains
var (
	canonicalSeed = 1
	forkSeed      = 2
)
var defaultChainConfig = params.TestChainConfig

// BlockGen creates blocks for testing.
// See GenerateChain for a detailed explanation.
type BlockGen struct {
	i           int
	parent      *block.Block
	chain       []*block.Block
	chainReader consensus.ChainReader
	header      *block.Header
	statedb     *state.StateDB
	db          database.IDatabase

	txs      []*transaction.Transaction
	receipts []*transaction.Receipt

	config *params.ChainConfig
	engine consensus.Engine
}

// SetCoinbase sets the coinbase of the generated block.
// It can be called at most once.
func (b *BlockGen) SetCoinbase(addr types.Address) {
	b.header.Producer = addr
}

// SetExtra sets the extra data field of the generated block.
func (b *BlockGen) SetConsensusData(data []byte) {
	b.header.Cdata.Para = data
}

// AddTx adds a transaction to the generated block. If no coinbase has
// been set, the block's coinbase is set to the zero address.
//
// AddTx panics if the transaction cannot be executed. In addition to
// the protocol-imposed limitations , there are some
// further limitations on the content of transactions that can be
// added. Notably, contract code relying on the BLOCKHASH instruction
// will panic during execution.
func (b *BlockGen) AddTx(tx *transaction.Transaction) {

	tmpDb, _ := database.OpenMemDB()
	blkCtx := actioncontext.NewBlockContext(b.statedb, b.db, tmpDb, &b.header.Number.IntVal, b.header.Producer)
	b.statedb.Prepare(tx.Hash(), types.Hash{}, len(b.txs))
	receipt, err := stateprocessor.ApplyTransaction(b.config, b.header, tx, blkCtx)
	if err != nil {
		panic(err)
	}
	b.txs = append(b.txs, tx)
	b.receipts = append(b.receipts, receipt)
}

// Number returns the block number of the block being generated.
func (b *BlockGen) Number() *big.Int {
	return new(big.Int).Set(&b.header.Number.IntVal)
}

// AddUncheckedReceipt forcefully adds a receipts to the block without a
// backing transaction.
//
// AddUncheckedReceipt will cause consensus failures when used during real
// chain processing. This is best used in conjunction with raw block insertion.
func (b *BlockGen) AddUncheckedReceipt(receipt *transaction.Receipt) {
	b.receipts = append(b.receipts, receipt)
}

// TxNonce returns the next valid transaction nonce for the
// account at addr. It panics if the account does not exist.
func (b *BlockGen) TxNonce(addr types.Address) uint64 {
	if !b.statedb.Exist(addr) {
		panic("account does not exist")
	}
	return b.statedb.GetNonce(addr)
}

// PrevBlock returns a previously generated block by number. It panics if
// num is greater or equal to the number of the block being generated.
// For index -1, PrevBlock returns the parent block given to GenerateChain.
func (b *BlockGen) PrevBlock(index int) *block.Block {
	if index >= b.i {
		panic("block index out of range")
	}
	if index == -1 {
		return b.parent
	}
	return b.chain[index]
}

// OffsetTime modifies the time instance of a block, implicitly changing its
// associated difficulty. It's useful to test scenarios where forking is not
// tied to chain length directly.
func (b *BlockGen) OffsetTime(seconds int64) {
	b.header.Time.IntVal.Add(&b.header.Time.IntVal, new(big.Int).SetInt64(seconds))
	if b.header.Time.IntVal.Cmp(&b.parent.Header().Time.IntVal) <= 0 {
		panic("block time out of range")
	}
	//b.header.Difficulty = b.engine.CalcDifficulty(b.chainReader, b.header.Time.Uint64(), b.parent.Header())
}

// GenerateChain creates a chain of n blocks. The first block's
// parent will be the provided parent. db is used to store
// intermediate states and should contain the parent's state trie.
//
// The generator function is called with a new block generator for
// every block. Any transactions.If gen is nil, the blocks will be empty
// and their coinbase will be the zero address.
//
// Blocks created by GenerateChain do not contain valid proof of work
// values. Inserting them into BlockChain requires use of FakePow or
// a similar non-validating proof of work implementation.
func GenerateChain(config *params.ChainConfig, parent *block.Block, engine consensus.Engine, db database.IDatabase, n int, gen func(int, *BlockGen)) ([]*block.Block, []transaction.Receipts) {
	if config == nil {
		config = defaultChainConfig
	}
	blocks, receipts := make(block.Blocks, n), make([]transaction.Receipts, n)
	genblock := func(i int, parent *block.Block, statedb *state.StateDB) (*block.Block, transaction.Receipts) {
		// TODO(karalabe): This is needed for clique, which depends on multiple blocks.
		// It's nonetheless ugly to spin up a blockchain here. Get rid of this somehow.
		blockchain, _ := blockchain.NewBlockChain(db, config, engine)
		defer blockchain.Stop()

		b := &BlockGen{i: i, parent: parent, chain: blocks, chainReader: blockchain, statedb: statedb, config: config, engine: engine, db: db}
		b.header = makeHeader(b.chainReader, parent, statedb, b.engine)

		// Mutate the state and block according to any hard-fork specs

		// Execute any user modifications to the block and finalize it
		if gen != nil {
			gen(i, b)
		}

		if b.engine != nil {
			//b.header.StateHash = statedb.IntermediateRoot()
			block, _ := b.engine.Finalize(b.chainReader, b.header, statedb, b.txs, b.receipts, true)
			//b.header.StateHash = statedb.IntermediateRoot()
			//block.B_header.StateHash = statedb.IntermediateRoot(true)
			// Write state changes to db
			_, _, err := statedb.CommitTo(db, false)
			if err != nil {
				panic(fmt.Sprintf("state write error: %v", err))
			}
			return block, b.receipts
		}
		return nil, nil
	}
	for i := 0; i < n; i++ {
		statedb, err := state.New(parent.Root(), state.NewDatabase(db))
		if err != nil {
			panic(err)
		}
		block, receipt := genblock(i, parent, statedb)
		blocks[i] = block
		receipts[i] = receipt
		parent = block
	}
	return blocks, receipts
}

func makeHeader(chain consensus.ChainReader, parent *block.Block, state *state.StateDB, engine consensus.Engine) *block.Header {
	var time *big.Int
	if parent.Time() == nil {
		time = big.NewInt(10)
	} else {
		time = new(big.Int).Add(parent.Time(), big.NewInt(10)) // block time is fixed at 10 seconds
	}

	return &block.Header{
		StateRootHash: state.IntermediateRoot(),
		ParentHash:    parent.Hash(),
		Producer:      parent.Producer(),
		Number:        types.NewBigInt(*new(big.Int).Add(parent.Number(), common.Big1)),
		Time:          types.NewBigInt(*time),
	}
}

// newCanonical creates a chain database, and injects a deterministic canonical
// chain. Depending on the full flag, if creates either a full block chain or a
// header only chain.
func newCanonical(engine consensus.Engine, n int, full bool) (database.IDatabase, *blockchain.BlockChain, error) {
	// Initialize a fresh chain with only a genesis block
	gspec := new(genesis.Genesis)
	db, _ := database.OpenMemDB()
	genesis := gspec.MustCommit(db)

	blockchain, _ := blockchain.NewBlockChain(db, defaultChainConfig, engine)
	// Create and inject the requested chain
	if n == 0 {
		return db, blockchain, nil
	}
	if full {
		// Full block-chain requested
		blocks := makeBlockChain(genesis, n, engine, db, canonicalSeed)
		_, err := blockchain.InsertChain(blocks)
		return db, blockchain, err
	}
	// Header-only chain requested
	headers := makeHeaderChain(genesis.Header(), n, engine, db, canonicalSeed)
	_, err := blockchain.InsertHeaderChain(headers, 1)
	return db, blockchain, err
}

// makeHeaderChain creates a deterministic chain of headers rooted at parent.
func makeHeaderChain(parent *block.Header, n int, engine consensus.Engine, db database.IDatabase, seed int) []*block.Header {
	blocks := makeBlockChain(block.NewBlockWithHeader(parent), n, engine, db, seed)
	headers := make([]*block.Header, len(blocks))
	for i, blk := range blocks {
		headers[i] = blk.Header()
	}
	return headers
}

// makeBlockChain creates a deterministic chain of blocks rooted at parent.
func makeBlockChain(parent *block.Block, n int, engine consensus.Engine, db database.IDatabase, seed int) []*block.Block {
	blocks, _ := GenerateChain(nil, parent, engine, db, n, func(i int, b *BlockGen) {
		b.SetCoinbase(types.Address{0: byte(seed), 19: byte(i)})
	})
	return blocks
}
