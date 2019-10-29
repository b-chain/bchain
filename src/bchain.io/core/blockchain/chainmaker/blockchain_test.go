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
// @File: blockchain_test.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package chainmaker

import (
	"fmt"
	"math/big"
	"bchain.io/common/types"
	"bchain.io/consensus"
	"bchain.io/core"
	"bchain.io/core/blockchain"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/genesis"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
	"bchain.io/utils/crypto"
	"bchain.io/utils/database"
	"sync"
	"testing"
	"time"
)

// newTestBlockChain creates a blockchain without validation.
func newTestBlockChain(fake bool) *blockchain.BlockChain {
	db, _ := database.OpenMemDB()
	gspec := &genesis.Genesis{
		Config: defaultChainConfig,
	}
	gspec.MustCommit(db)

	var engine consensus.Engine
	key, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f292")
	engine = consensus.NewBasicEngine(key)
	if fake {
		engine = &consensus.Engine_empty{}
	}

	blockchain, err := blockchain.NewBlockChain(db, gspec.Config, engine)
	if err != nil {
		panic(err)
	}
	blockchain.SetValidator(bproc{})
	return blockchain
}

// Test fork of length N starting from block i
func testFork(t *testing.T, blockchain *blockchain.BlockChain, i, n int, full bool, comparator func(td1, td2 *big.Int)) {
	// Copy old chain up to #i into a new db
	engine := &consensus.Engine_empty{}
	db, blockchain2, err := newCanonical(engine, i, full)
	if err != nil {
		t.Fatal("could not make new canonical in testFork", err)
	}
	defer blockchain2.Stop()

	// Assert the chains have the same header/block at #i
	var hash1, hash2 types.Hash
	if full {
		hash1 = blockchain.GetBlockByNumber(uint64(i)).Hash()
		hash2 = blockchain2.GetBlockByNumber(uint64(i)).Hash()
	} else {
		hash1 = blockchain.GetHeaderByNumber(uint64(i)).Hash()
		hash2 = blockchain2.GetHeaderByNumber(uint64(i)).Hash()
	}
	if hash1 != hash2 {
		t.Errorf("chain content mismatch at %d: have hash %v, want hash %v", i, hash2, hash1)
	}
	// Extend the newly created chain
	var (
		blockChainB  []*block.Block
		headerChainB []*block.Header
	)
	if full {
		blockChainB = makeBlockChain(blockchain2.CurrentBlock(), n, engine, db, forkSeed)
		if _, err := blockchain2.InsertChain(blockChainB); err != nil {
			t.Fatalf("failed to insert forking chain: %v", err)
		}
	} else {
		headerChainB = makeHeaderChain(blockchain2.CurrentHeader(), n, engine, db, forkSeed)
		if _, err := blockchain2.InsertHeaderChain(headerChainB, 1); err != nil {
			t.Fatalf("failed to insert forking chain: %v", err)
		}
	}
	// Sanity check that the forked chain can be imported into the original
	var numPre, numPost *big.Int

	if full {
		numPre = blockchain.CurrentBlock().Number()
		if err := testBlockChainImport(blockChainB, blockchain); err != nil {
			t.Fatalf("failed to import forked block chain: %v", err)
		}
		numPost = blockChainB[len(blockChainB)-1].Number()
	} else {
		numPre = &blockchain.CurrentHeader().Number.IntVal
		if err := testHeaderChainImport(headerChainB, blockchain); err != nil {
			t.Fatalf("failed to import forked header chain: %v", err)
		}
		numPost = &headerChainB[len(headerChainB)-1].Number.IntVal
	}
	// Compare the total difficulties of the chains
	comparator(numPre, numPost)
}

func printChain(bc *blockchain.BlockChain) {
	for i := bc.CurrentBlock().Number().Uint64(); i > 0; i-- {
		b := bc.GetBlockByNumber(uint64(i))
		fmt.Printf("\t%x %v\n", b.Hash(), b.Number())
	}
}

// testBlockChainImport tries to process a chain of blocks, writing them into
// the database if successful.
func testBlockChainImport(chain block.Blocks, bc *blockchain.BlockChain) error {
	for _, block := range chain {
		// Try and process the block
		err := bc.Engine().VerifyHeader(bc, block.Header(), true)
		if err == nil {
			err = bc.Validate().ValidateBody(block)
		}
		if err != nil {
			if err == core.ErrKnownBlock {
				continue
			}
			return err
		}
		statedb, err := state.New(bc.GetBlockByHash(block.ParentHash()).Root(), bc.StateCache())
		if err != nil {
			return err
		}
		receipts, _, err := bc.Processor().Process(block, statedb, bc.GetDb(), bc.Config())
		if err != nil {
			bc.ReportBlock(block, receipts, err)
			return err
		}
		err = bc.Validate().ValidateState(block, bc.GetBlockByHash(block.ParentHash()), statedb, receipts)
		if err != nil {
			bc.ReportBlock(block, receipts, err)
			return err
		}
		bc.MuLock()
		blockchain.WriteBlock(bc.GetDb(), block)
		statedb.CommitTo(bc.GetDb(), false)
		bc.MuUnLock()
	}
	return nil
}

// testHeaderChainImport tries to process a chain of header, writing them into
// the database if successful.
func testHeaderChainImport(chain []*block.Header, bc *blockchain.BlockChain) error {
	for _, header := range chain {
		// Try and validate the header
		if err := bc.Engine().VerifyHeader(bc, header, false); err != nil {
			return err
		}
		// Manually insert the header into the database, but don't reorganise (allows subsequent testing)
		bc.MuLock()
		blockchain.WriteHeader(bc.GetDb(), header)
		bc.MuUnLock()
	}
	return nil
}

func insertChain(done chan bool, blockchain *blockchain.BlockChain, chain block.Blocks, t *testing.T) {
	_, err := blockchain.InsertChain(chain)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	done <- true
}

func TestLastBlock(t *testing.T) {
	bchain := newTestBlockChain(false)
	defer bchain.Stop()
	key, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f292")
	block := makeBlockChain(bchain.CurrentBlock(), 1, consensus.NewBasicEngine(key), bchain.GetDb(), 0)[0]
	bchain.Test_insert(block)
	if block.Hash() != blockchain.GetHeadBlockHash(bchain.GetDb()) {
		t.Errorf("Write/Get HeadBlockHash failed")
	}
}

// Tests that given a starting canonical chain of a given size, it can be extended
// with various length chains.
func TestExtendCanonicalHeaders(t *testing.T) { testExtendCanonical(t, false) }
func TestExtendCanonicalBlocks(t *testing.T)  { testExtendCanonical(t, true) }

func testExtendCanonical(t *testing.T, full bool) {
	length := 5

	// Make first chain starting from genesis
	_, processor, err := newCanonical(&consensus.Engine_empty{}, length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the difficulty comparator
	better := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) <= 0 {
			t.Errorf("total difficulty mismatch: have %v, expected more than %v", td2, td1)
		}
	}
	// Start fork from current height
	testFork(t, processor, length, 1, full, better)
	testFork(t, processor, length, 2, full, better)
	testFork(t, processor, length, 5, full, better)
	testFork(t, processor, length, 10, full, better)
}

// Tests that given a starting canonical chain of a given size, creating shorter
// forks do not take canonical ownership.
func TestShorterForkHeaders(t *testing.T) { testShorterFork(t, false) }
func TestShorterForkBlocks(t *testing.T)  { testShorterFork(t, true) }

func testShorterFork(t *testing.T, full bool) {
	length := 10

	// Make first chain starting from genesis
	_, processor, err := newCanonical(&consensus.Engine_empty{}, length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the difficulty comparator
	worse := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) >= 0 {
			t.Errorf("total difficulty mismatch: have %v, expected less than %v", td2, td1)
		}
	}
	// Sum of numbers must be less than `length` for this to be a shorter fork
	testFork(t, processor, 0, 3, full, worse)
	testFork(t, processor, 0, 7, full, worse)
	testFork(t, processor, 1, 1, full, worse)
	testFork(t, processor, 1, 7, full, worse)
	testFork(t, processor, 5, 3, full, worse)
	testFork(t, processor, 5, 4, full, worse)
}

// Tests that given a starting canonical chain of a given size, creating longer
// forks do take canonical ownership.
func TestLongerForkHeaders(t *testing.T) { testLongerFork(t, false) }
func TestLongerForkBlocks(t *testing.T)  { testLongerFork(t, true) }

func testLongerFork(t *testing.T, full bool) {
	length := 10

	// Make first chain starting from genesis
	_, processor, err := newCanonical(&consensus.Engine_empty{}, length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the difficulty comparator
	better := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) <= 0 {
			t.Errorf("total difficulty mismatch: have %v, expected more than %v", td2, td1)
		}
	}
	// Sum of numbers must be greater than `length` for this to be a longer fork
	testFork(t, processor, 0, 11, full, better)
	testFork(t, processor, 0, 15, full, better)
	testFork(t, processor, 1, 10, full, better)
	testFork(t, processor, 1, 12, full, better)
	testFork(t, processor, 5, 6, full, better)
	testFork(t, processor, 5, 8, full, better)
}

// Tests that given a starting canonical chain of a given size, creating equal
// forks do take canonical ownership.
func TestEqualForkHeaders(t *testing.T) { testEqualFork(t, false) }
func TestEqualForkBlocks(t *testing.T)  { testEqualFork(t, true) }

func testEqualFork(t *testing.T, full bool) {
	length := 10

	// Make first chain starting from genesis
	_, processor, err := newCanonical(&consensus.Engine_empty{}, length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the difficulty comparator
	equal := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) != 0 {
			t.Errorf("total difficulty mismatch: have %v, want %v", td2, td1)
		}
	}
	// Sum of numbers must be equal to `length` for this to be an equal fork
	testFork(t, processor, 0, 10, full, equal)
	testFork(t, processor, 1, 9, full, equal)
	testFork(t, processor, 2, 8, full, equal)
	testFork(t, processor, 5, 5, full, equal)
	testFork(t, processor, 6, 4, full, equal)
	testFork(t, processor, 9, 1, full, equal)
}

// Tests that chains missing links do not get accepted by the processor.
//func TestBrokenHeaderChain(t *testing.T) { testBrokenChain(t, false) }
func TestBrokenBlockChain(t *testing.T) { testBrokenChain(t, true) }

func testBrokenChain(t *testing.T, full bool) {
	// Make chain starting from genesis
	db, blockchain, err := newCanonical(&consensus.Engine_empty{}, 10, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer blockchain.Stop()

	// Create a forked chain, and try to insert with a missing link
	if full {
		chain := makeBlockChain(blockchain.CurrentBlock(), 5, &consensus.Engine_empty{}, db, forkSeed)[1:]
		if err := testBlockChainImport(chain, blockchain); err == nil {
			t.Errorf("broken block chain not reported")
		}
	} else {
		chain := makeHeaderChain(blockchain.CurrentHeader(), 5, &consensus.Engine_empty{}, db, forkSeed)[1:]
		if err := testHeaderChainImport(chain, blockchain); err == nil {
			t.Errorf("broken header chain not reported")
		}
	}
}

type bproc struct{}

func (bproc) ValidateBody(*block.Block) error { return nil }
func (bproc) ValidateState(block, parent *block.Block, state *state.StateDB, receipts transaction.Receipts) error {
	return nil
}
func (bproc) Process(block *block.Block, statedb *state.StateDB) (transaction.Receipts, []*transaction.Log, uint64, error) {
	return nil, nil, 0, nil
}

func makeHeaderChainWithDiff(genesis *block.Block, n int, seed byte) []*block.Header {
	blocks := makeBlockChainWithDiff(genesis, n, seed)
	headers := make([]*block.Header, len(blocks))
	for i, blk := range blocks {
		headers[i] = blk.Header()
	}
	return headers
}

func makeBlockChainWithDiff(genesis *block.Block, n int, seed byte) []*block.Block {
	var chain []*block.Block
	singner := block.NewBlockSigner(defaultChainConfig.ChainId)

	var (
		key, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f292")
	)

	for i := 0; i < n; i++ {
		header := &block.Header{
			Producer:        types.Address{seed},
			Number:          types.NewBigInt(*big.NewInt(int64(i + 1))),
			TxRootHash:      block.EmptyRootHash,
			ReceiptRootHash: block.EmptyRootHash,
			Time:            types.NewBigInt(*big.NewInt(int64(i) + 1)),
		}
		if i == 0 {
			header.ParentHash = genesis.Hash()
		} else {
			header.ParentHash = chain[i-1].Hash()
		}
		signHeaer, _ := block.SignHeader(header, singner, key)
		//fmt.Println(signHeaer)
		block := block.NewBlockWithHeader(signHeaer)
		chain = append(chain, block)
	}
	return chain
}

// Tests that reorganising a long difficult chain after a short easy one
// overwrites the canonical numbers and links in the database.
func TestReorgLongHeaders(t *testing.T) { testReorgLong(t, false) }
func TestReorgLongBlocks(t *testing.T)  { testReorgLong(t, true) }

func testReorgLong(t *testing.T, full bool) {
	testReorg(t, 3, 4, 4, full)
}

// Tests that reorganising a short difficult chain after a long easy one
// overwrites the canonical numbers and links in the database.
func TestReorgShortHeaders(t *testing.T) { testReorgShort(t, false) }
func TestReorgShortBlocks(t *testing.T)  { testReorgShort(t, true) }

func testReorgShort(t *testing.T, full bool) {
	testReorg(t, 5, 3, 5, full)
}

func testReorg(t *testing.T, first, second int, num uint64, full bool) {
	bc := newTestBlockChain(false)
	defer bc.Stop()

	// Insert an easy and a difficult chain afterwards
	if full {
		bc.InsertChain(makeBlockChainWithDiff(bc.GetGenesis(), first, 11))
		bc.InsertChain(makeBlockChainWithDiff(bc.GetGenesis(), second, 22))
	} else {
		bc.InsertHeaderChain(makeHeaderChainWithDiff(bc.GetGenesis(), first, 11), 1)
		bc.InsertHeaderChain(makeHeaderChainWithDiff(bc.GetGenesis(), second, 22), 1)
	}
	// Check that the chain is valid number and link wise
	if full {
		prev := bc.CurrentBlock()
		for block := bc.GetBlockByNumber(bc.CurrentBlock().NumberU64() - 1); block.NumberU64() != 0; prev, block = block, bc.GetBlockByNumber(block.NumberU64()-1) {
			if prev.ParentHash() != block.Hash() {
				t.Errorf("parent block hash mismatch: have %x, want %x", prev.ParentHash(), block.Hash())
			}
		}
	} else {
		prev := bc.CurrentHeader()
		for header := bc.GetHeaderByNumber(bc.CurrentHeader().Number.IntVal.Uint64() - 1); header.Number.IntVal.Uint64() != 0; prev, header = header, bc.GetHeaderByNumber(header.Number.IntVal.Uint64()-1) {
			if prev.ParentHash != header.Hash() {
				t.Errorf("parent header hash mismatch: have %x, want %x", prev.ParentHash, header.Hash())
			}
		}
	}
	// Make sure the chain total difficulty is the correct one
	want := num
	if full {
		if have := bc.CurrentBlock().Number().Uint64(); want != have {
			t.Errorf("total difficulty mismatch: have %v, want %v", have, want)
		}
	} else {
		if have := bc.CurrentHeader().Number.IntVal.Uint64(); want != have {
			t.Errorf("total difficulty mismatch: have %v, want %v", have, want)
		}
	}
}

/*
// Tests chain insertions in the face of one entity containing an invalid nonce.
func TestHeadersInsertNonceError(t *testing.T) { testInsertNonceError(t, false) }
func TestBlocksInsertNonceError(t *testing.T)  { testInsertNonceError(t, true) }

func testInsertNonceError(t *testing.T, full bool) {
	for i := 1; i < 25 && !t.Failed(); i++ {
		// Create a pristine chain and database
		db, blockchain, err := newCanonical(&consensus.Engine_empty{}, 0, full)
		if err != nil {
			t.Fatalf("failed to create pristine chain: %v", err)
		}
		defer blockchain.Stop()

		// Create and insert a chain with a failing nonce
		var (
			failAt  int
			failRes int
			failNum uint64
		)
		if full {
			blocks := makeBlockChain(blockchain.CurrentBlock(), i, &consensus.Engine_empty{}, db, 0)

			failAt = rand.Int() % len(blocks)
			failNum = blocks[failAt].NumberU64()

			failRes, err = blockchain.InsertChain(blocks)
		} else {
			headers := makeHeaderChain(blockchain.CurrentHeader(), i, &consensus.Engine_empty{}, db, 0)

			failAt = rand.Int() % len(headers)
			failNum = headers[failAt].Number.IntVal.Uint64()

			failRes, err = blockchain.InsertHeaderChain(headers, 1)
		}
		// Check that the returned error indicates the failure.
		if failRes != failAt {
			t.Errorf("test %d: failure index mismatch: have %d, want %d", i, failRes, failAt)
		}
		// Check that all no blocks after the failing block have been inserted.
		for j := 0; j < i-failAt; j++ {
			if full {
				if block := blockchain.GetBlockByNumber(failNum + uint64(j)); block != nil {
					t.Errorf("test %d: invalid block in chain: %v", i, block)
				}
			} else {
				if header := blockchain.GetHeaderByNumber(failNum + uint64(j)); header != nil {
					t.Errorf("test %d: invalid header in chain: %v", i, header)
				}
			}
		}
	}
}
*/
// Tests that fast importing a block chain produces the same chain data as the
// classical full block processing.
func TestFastVsFullChains(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb, _ = database.OpenMemDB()
		key, _   = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address  = crypto.PubkeyToAddress(key.PublicKey)
		//funds    = big.NewInt(1000000000)
		gspec = &genesis.Genesis{
			Config: defaultChainConfig,
			Alloc:  genesis.GenesisAlloc{address: {}},
		}
		genesis = gspec.MustCommit(gendb)
		signer  = transaction.NewMSigner(gspec.Config.ChainId)
		actions = transaction.Actions{}
	)
	actions = append(actions, &transaction.Action{types.Address{0x00}, []byte{1, 2, 3, 4, 5}})
	blocks, receipts := GenerateChain(gspec.Config, genesis, &consensus.Engine_empty{}, gendb, 1024, func(i int, blockGen *BlockGen) {
		blockGen.SetCoinbase(types.Address{0x00})

		// If the block number is multiple of 3, send a few bonus transactions to the blockproducer
		if i%3 == 2 {
			for j := 0; j < i%4+1; j++ {
				tx, err := transaction.SignTx(transaction.NewTransaction(blockGen.TxNonce(address), actions), signer, key)
				if err != nil {
					panic(err)
				}
				blockGen.AddTx(tx)
			}
		}
	})
	// Import the chain as an archive node for the comparison baseline
	archiveDb, _ := database.OpenMemDB()
	gspec.MustCommit(archiveDb)
	archive, _ := blockchain.NewBlockChain(archiveDb, gspec.Config, &consensus.Engine_empty{})
	defer archive.Stop()

	if n, err := archive.InsertChain(blocks); err != nil {
		t.Fatalf("failed to process block %d: %v", n, err)
	}
	// Fast import the chain as a non-archive node to test
	fastDb, _ := database.OpenMemDB()
	gspec.MustCommit(fastDb)
	fast, _ := blockchain.NewBlockChain(fastDb, gspec.Config, &consensus.Engine_empty{})
	defer fast.Stop()

	headers := make([]*block.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := fast.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := fast.InsertReceiptChain(blocks, receipts); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	// Iterate over all chain data components, and cross reference
	for i := 0; i < len(blocks); i++ {
		num, hash := blocks[i].NumberU64(), blocks[i].Hash()

		if fheader, aheader := fast.GetHeaderByHash(hash), archive.GetHeaderByHash(hash); fheader.Hash() != aheader.Hash() {
			t.Errorf("block #%d [%x]: header mismatch: have %v, want %v", num, hash, fheader, aheader)
		}
		if fblock, ablock := fast.GetBlockByHash(hash), archive.GetBlockByHash(hash); fblock.Hash() != ablock.Hash() {
			t.Errorf("block #%d [%x]: block mismatch: have %v, want %v", num, hash, fblock, ablock)
		} else if block.DeriveSha(fblock.Transactions()) != block.DeriveSha(ablock.Transactions()) {
			t.Errorf("block #%d [%x]: transactions mismatch: have %v, want %v", num, hash, fblock.Transactions(), ablock.Transactions())
		}
		if freceipts, areceipts := blockchain.GetBlockReceipts(fastDb, hash, blockchain.GetBlockNumber(fastDb, hash)), blockchain.GetBlockReceipts(archiveDb, hash, blockchain.GetBlockNumber(archiveDb, hash)); block.DeriveSha(freceipts) != block.DeriveSha(areceipts) {
			t.Errorf("block #%d [%x]: receipts mismatch: have %v, want %v", num, hash, freceipts, areceipts)
		}
	}
	// Check that the canonical chains are the same between the databases
	for i := 0; i < len(blocks)+1; i++ {
		if fhash, ahash := blockchain.GetCanonicalHash(fastDb, uint64(i)), blockchain.GetCanonicalHash(archiveDb, uint64(i)); fhash != ahash {
			t.Errorf("block #%d: canonical hash mismatch: have %v, want %v", i, fhash, ahash)
		}
	}
}

// Tests that various import methods move the chain head pointers to the correct
// positions.
func TestLightVsFastVsFullChainHeads(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb, _ = database.OpenMemDB()
		key, _   = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address  = crypto.PubkeyToAddress(key.PublicKey)
		//funds    = big.NewInt(1000000000)
		gspec   = &genesis.Genesis{Config: defaultChainConfig, Alloc: genesis.GenesisAlloc{address: {}}}
		genesis = gspec.MustCommit(gendb)
	)
	height := uint64(1024)
	blocks, receipts := GenerateChain(gspec.Config, genesis, &consensus.Engine_empty{}, gendb, int(height), nil)

	// Configure a subchain to roll back
	remove := []types.Hash{}
	for _, block := range blocks[height/2:] {
		remove = append(remove, block.Hash())
	}
	// Create a small assertion method to check the three heads
	assert := func(t *testing.T, kind string, chain *blockchain.BlockChain, header uint64, fast uint64, block uint64) {
		if num := chain.CurrentBlock().NumberU64(); num != block {
			t.Errorf("%s head block mismatch: have #%v, want #%v", kind, num, block)
		}
		if num := chain.CurrentFastBlock().NumberU64(); num != fast {
			t.Errorf("%s head fast-block mismatch: have #%v, want #%v", kind, num, fast)
		}
		if num := chain.CurrentHeader().Number.IntVal.Uint64(); num != header {
			t.Errorf("%s head header mismatch: have #%v, want #%v", kind, num, header)
		}
	}
	// Import the chain as an archive node and ensure all pointers are updated
	archiveDb, _ := database.OpenMemDB()
	gspec.MustCommit(archiveDb)

	archive, _ := blockchain.NewBlockChain(archiveDb, gspec.Config, &consensus.Engine_empty{})
	if n, err := archive.InsertChain(blocks); err != nil {
		t.Fatalf("failed to process block %d: %v", n, err)
	}
	defer archive.Stop()

	assert(t, "archive", archive, height, height, height)
	archive.Rollback(remove)
	assert(t, "archive", archive, height/2, height/2, height/2)

	// Import the chain as a non-archive node and ensure all pointers are updated
	fastDb, _ := database.OpenMemDB()
	gspec.MustCommit(fastDb)
	fast, _ := blockchain.NewBlockChain(fastDb, gspec.Config, &consensus.Engine_empty{})
	defer fast.Stop()

	headers := make([]*block.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := fast.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := fast.InsertReceiptChain(blocks, receipts); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	assert(t, "fast", fast, height, height, 0)
	fast.Rollback(remove)
	assert(t, "fast", fast, height/2, height/2, 0)

	// Import the chain as a light node and ensure all pointers are updated
	lightDb, _ := database.OpenMemDB()
	gspec.MustCommit(lightDb)

	light, _ := blockchain.NewBlockChain(lightDb, gspec.Config, &consensus.Engine_empty{})
	if n, err := light.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	defer light.Stop()

	assert(t, "light", light, height, 0, 0)
	light.Rollback(remove)
	assert(t, "light", light, height/2, 0, 0)
}

// Tests that chain reorganisations handle transaction removals and reinsertions.
func TestChainTxReorgs(t *testing.T) {
	var (
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key2, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		key3, _ = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		addr2   = crypto.PubkeyToAddress(key2.PublicKey)
		addr3   = crypto.PubkeyToAddress(key3.PublicKey)
		db, _   = database.OpenMemDB()
		gspec   = &genesis.Genesis{
			Config: defaultChainConfig,
			Alloc: genesis.GenesisAlloc{
				addr1: {},
				addr2: {},
				addr3: {},
			},
		}
		genesis = gspec.MustCommit(db)
		signer  = transaction.NewMSigner(gspec.Config.ChainId)
		actions = transaction.Actions{}
	)
	actions = append(actions, &transaction.Action{types.Address{0x00}, []byte{1, 2, 3, 4, 5}})

	// Create two transactions shared between the chains:
	//  - postponed: transaction included at a later block in the forked chain
	//  - swapped: transaction included at the same block number in the forked chain
	postponed, _ := transaction.SignTx(transaction.NewTransaction(0, actions), signer, key1)
	//todo set nonce 0 for fake test. future need change to 1
	swapped, _ := transaction.SignTx(transaction.NewTransaction(0, actions), signer, key1)

	// Create two transactions that will be dropped by the forked chain:
	//  - pastDrop: transaction dropped retroactively from a past block
	//  - freshDrop: transaction dropped exactly at the block where the reorg is detected
	var pastDrop, freshDrop *transaction.Transaction

	// Create three transactions that will be added in the forked chain:
	//  - pastAdd:   transaction added before the reorganization is detected
	//  - freshAdd:  transaction added at the exact block the reorg is detected
	//  - futureAdd: transaction added after the reorg has already finished
	var pastAdd, freshAdd, futureAdd *transaction.Transaction

	chain, _ := GenerateChain(gspec.Config, genesis, &consensus.Engine_empty{}, db, 3, func(i int, gen *BlockGen) {
		switch i {
		case 0:
			pastDrop, _ = transaction.SignTx(transaction.NewTransaction(gen.TxNonce(addr2), actions), signer, key2)

			gen.AddTx(pastDrop)  // This transaction will be dropped in the fork from below the split point
			gen.AddTx(postponed) // This transaction will be postponed till block #3 in the fork

		case 2:
			freshDrop, _ = transaction.SignTx(transaction.NewTransaction(gen.TxNonce(addr2), actions), signer, key2)

			gen.AddTx(freshDrop) // This transaction will be dropped in the fork from exactly at the split point
			gen.AddTx(swapped)   // This transaction will be swapped out at the exact height

			gen.OffsetTime(9) // Lower the block difficulty to simulate a weaker chain
		}
	})
	// Import the chain. This runs all block validation rules.
	bc, _ := blockchain.NewBlockChain(db, gspec.Config, &consensus.Engine_empty{})
	if i, err := bc.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert original chain[%d]: %v", i, err)
	}
	defer bc.Stop()

	// overwrite the old chain
	chain, _ = GenerateChain(gspec.Config, genesis, &consensus.Engine_empty{}, db, 5, func(i int, gen *BlockGen) {
		switch i {
		case 0:
			pastAdd, _ = transaction.SignTx(transaction.NewTransaction(gen.TxNonce(addr3), actions), signer, key3)
			gen.AddTx(pastAdd) // This transaction needs to be injected during reorg

		case 2:
			gen.AddTx(postponed) // This transaction was postponed from block #1 in the original chain
			gen.AddTx(swapped)   // This transaction was swapped from the exact current spot in the original chain

			freshAdd, _ = transaction.SignTx(transaction.NewTransaction(gen.TxNonce(addr3), actions), signer, key3)
			gen.AddTx(freshAdd) // This transaction will be added exactly at reorg time

		case 3:
			futureAdd, _ = transaction.SignTx(transaction.NewTransaction(gen.TxNonce(addr3), actions), signer, key3)
			gen.AddTx(futureAdd) // This transaction will be added after a full reorg
		}
	})
	if _, err := bc.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert forked chain: %v", err)
	}

	// removed tx
	for i, tx := range (transaction.Transactions{pastDrop, freshDrop}) {
		if txn, _, _, _ := blockchain.GetTransaction(db, tx.Hash()); txn != nil {
			t.Errorf("drop %d: tx %v found while shouldn't have been", i, txn)
		}
		if rcpt, _, _, _ := blockchain.GetReceipt(db, tx.Hash()); rcpt != nil {
			t.Errorf("drop %d: receipt %v found while shouldn't have been", i, rcpt)
		}
	}
	// added tx
	for i, tx := range (transaction.Transactions{pastAdd, freshAdd, futureAdd}) {
		if txn, _, _, _ := blockchain.GetTransaction(db, tx.Hash()); txn == nil {
			t.Errorf("add %d: expected tx to be found", i)
		}
		if rcpt, _, _, _ := blockchain.GetReceipt(db, tx.Hash()); rcpt == nil {
			t.Errorf("add %d: expected receipt to be found", i)
		}
	}
	// shared tx
	for i, tx := range (transaction.Transactions{postponed, swapped}) {
		if txn, _, _, _ := blockchain.GetTransaction(db, tx.Hash()); txn == nil {
			t.Errorf("share %d: expected tx to be found", i)
		}
		if rcpt, _, _, _ := blockchain.GetReceipt(db, tx.Hash()); rcpt == nil {
			t.Errorf("share %d: expected receipt to be found", i)
		}
	}
}

func TestLogReorgs(t *testing.T) {

	var (
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		db, _   = database.OpenMemDB()
		// this code generates a log
		// todo code    = util.Hex2Bytes("60606040525b7f24ec1d3ff24c2f6ff210738839dbc339cd45a5294d85c79361016243157aae7b60405180905060405180910390a15b600a8060416000396000f360606040526008565b00")
		gspec   = &genesis.Genesis{Config: defaultChainConfig, Alloc: genesis.GenesisAlloc{addr1: {}}}
		genesis = gspec.MustCommit(db)
		signer  = transaction.NewMSigner(gspec.Config.ChainId)
		actions = transaction.Actions{}
	)
	actions = append(actions, &transaction.Action{types.Address{0x00}, []byte{1, 2, 3, 4, 5}})
	blockchain, _ := blockchain.NewBlockChain(db, gspec.Config, &consensus.Engine_empty{})
	defer blockchain.Stop()

	rmLogsCh := make(chan core.RemovedLogsEvent)
	blockchain.SubscribeRemovedLogsEvent(rmLogsCh)
	chain, _ := GenerateChain(defaultChainConfig, genesis, &consensus.Engine_empty{}, db, 2, func(i int, gen *BlockGen) {
		if i == 1 {
			tx, err := transaction.SignTx(transaction.NewTransaction(gen.TxNonce(addr1), actions), signer, key1)
			// todo : tx, err := transaction.SignTx(transaction.NewContractCreation(gen.TxNonce(addr1), big.NewInt(100), 1000000, big.NewInt(100), code), signer, key1)
			if err != nil {
				t.Fatalf("failed to create tx: %v", err)
			}
			gen.AddTx(tx)
		}
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	chain, _ = GenerateChain(defaultChainConfig, genesis, &consensus.Engine_empty{}, db, 3, func(i int, gen *BlockGen) {})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert forked chain: %v", err)
	}

	timeout := time.NewTimer(1 * time.Second)
	select {
	case ev := <-rmLogsCh:
		if len(ev.Logs) == 0 {
			t.Error("expected logs")
		}
	case <-timeout.C:
		// todo t.Fatal("Timeout. There is no RemovedLogsEvent has been sent.")
	}
}

func TestReorgSideEvent(t *testing.T) {
	var (
		db, _   = database.OpenMemDB()
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		gspec   = &genesis.Genesis{
			Config: defaultChainConfig,
			Alloc:  genesis.GenesisAlloc{addr1: {}},
		}
		genesis = gspec.MustCommit(db)
		signer  = transaction.NewMSigner(gspec.Config.ChainId)
		actions = transaction.Actions{}
	)
	actions = append(actions, &transaction.Action{types.Address{0x00}, []byte{1, 2, 3, 4, 5}})
	blockchain, _ := blockchain.NewBlockChain(db, gspec.Config, &consensus.Engine_empty{})
	defer blockchain.Stop()

	chain, _ := GenerateChain(gspec.Config, genesis, &consensus.Engine_empty{}, db, 3, func(i int, gen *BlockGen) {})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	replacementBlocks, _ := GenerateChain(gspec.Config, genesis, &consensus.Engine_empty{}, db, 4, func(i int, gen *BlockGen) {
		tx, err := transaction.SignTx(transaction.NewTransaction(gen.TxNonce(addr1), actions), signer, key1)
		//tx, err := transaction.SignTx(transaction.NewContractCreation(gen.TxNonce(addr1), big.NewInt(100), 1000000, big.NewInt(100), nil), signer, key1)
		if i == 2 {
			gen.OffsetTime(-9)
		}
		if err != nil {
			t.Fatalf("failed to create tx: %v", err)
		}
		gen.AddTx(tx)
	})
	chainSideCh := make(chan core.ChainSideEvent, 64)
	blockchain.SubscribeChainSideEvent(chainSideCh)
	if _, err := blockchain.InsertChain(replacementBlocks); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	// first two block of the secondary chain are for a brief moment considered
	// side chains because up to that point the first one is considered the
	// heavier chain.
	expectedSideHashes := map[types.Hash]bool{
		replacementBlocks[0].Hash(): true,
		replacementBlocks[1].Hash(): true,
		replacementBlocks[2].Hash(): true,
		chain[0].Hash():             true,
		chain[1].Hash():             true,
		chain[2].Hash():             true,
	}

	i := 0

	const timeoutDura = 10 * time.Second
	timeout := time.NewTimer(timeoutDura)
done:
	for {
		select {
		case ev := <-chainSideCh:
			block := ev.Block
			if _, ok := expectedSideHashes[block.Hash()]; !ok {
				t.Errorf("%d: didn't expect %x to be in side chain", i, block.Hash())
			}
			i++

			if i == len(expectedSideHashes) {
				timeout.Stop()

				break done
			}
			timeout.Reset(timeoutDura)

		case <-timeout.C:
			t.Fatal("Timeout. Possibly not all blocks were triggered for sideevent")
		}
	}

	// make sure no more events are fired
	select {
	case e := <-chainSideCh:
		t.Errorf("unexpected event fired: %v", e)
	case <-time.After(250 * time.Millisecond):
	}

}

// Tests if the canonical block can be fetched from the database during chain insertion.
func TestCanonicalBlockRetrieval(t *testing.T) {
	bc := newTestBlockChain(true)
	defer bc.Stop()

	chain, _ := GenerateChain(defaultChainConfig, bc.GetGenesis(), &consensus.Engine_empty{}, bc.GetDb(), 10, func(i int, gen *BlockGen) {})

	var pend sync.WaitGroup
	pend.Add(len(chain))

	for i := range chain {
		go func(block *block.Block) {
			defer pend.Done()

			// try to retrieve a block by its canonical hash and see if the block data can be retrieved.
			for {
				ch := blockchain.GetCanonicalHash(bc.GetDb(), block.NumberU64())
				if ch == (types.Hash{}) {
					continue // busy wait for canonical hash to be written
				}
				if ch != block.Hash() {
					t.Fatalf("unknown canonical hash, want %s, got %s", block.Hash().Hex(), ch.Hex())
				}
				fb := blockchain.GetBlock(bc.GetDb(), ch, block.NumberU64())
				if fb == nil {
					t.Fatalf("unable to retrieve block %d for canonical hash: %s", block.NumberU64(), ch.Hex())
				}
				if fb.Hash() != block.Hash() {
					t.Fatalf("invalid block hash for block %d, want %s, got %s", block.NumberU64(), block.Hash().Hex(), fb.Hash().Hex())
				}
				return
			}
		}(chain[i])

		if _, err := bc.InsertChain(block.Blocks{chain[i]}); err != nil {
			t.Fatalf("failed to insert block %d: %v", i, err)
		}
	}
	pend.Wait()
}
