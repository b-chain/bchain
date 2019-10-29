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
// @File: handler_test.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package bchain

import (
	"math"
	"math/big"
	"math/rand"
	"testing"

	"bchain.io/node/services/bchain/downloader"
	"bchain.io/communication/p2p"
	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
	"bchain.io/utils/crypto"
	"bchain.io/core/transaction"
	"bchain.io/core/blockchain/chainmaker"
	"bchain.io/core/blockchain"
	"bchain.io/utils/database"
	"bchain.io/core/state"
)

// Tests that protocol versions and modes of operations are matched up properly.
func TestProtocolCompatibility(t *testing.T) {
	// Define the compatibility chart
	tests := []struct {
		version    uint
		mode       downloader.SyncMode
		compatible bool
	}{
		{61, downloader.FullSync, true}, {62, downloader.FullSync, true}, {63, downloader.FullSync, true},
		{61, downloader.FastSync, false}, {62, downloader.FastSync, false}, {63, downloader.FastSync, true},
	}
	// Make sure anything we screw up is restored
	backup := ProtocolVersions
	defer func() { ProtocolVersions = backup }()

	// Try all available compatibility configs and check for errors
	for i, tt := range tests {
		ProtocolVersions = []uint{tt.version}

		pm, err := newTestProtocolManager(tt.mode, 0, nil, nil)
		if pm != nil {
			defer pm.Stop()
		}
		if (err == nil && !tt.compatible) || (err != nil && tt.compatible) {
			t.Errorf("test %d: compatibility mismatch: have error %v, want compatibility %v", i, err, tt.compatible)
		}
	}
}

// Tests that block headers can be retrieved from a remote chain based on user queries.
func TestGetBlockHeaders62(t *testing.T) { testGetBlockHeaders(t, 62) }
func TestGetBlockHeaders63(t *testing.T) { testGetBlockHeaders(t, 63) }

func testGetBlockHeaders(t *testing.T, protocol int) {
	pm := newTestProtocolManagerMust(t, downloader.FullSync, downloader.MaxHashFetch+15, nil, nil)
	peer, _ := newTestPeer("peer", protocol, pm, true)
	defer peer.close()

	// Create a "random" unknown hash for testing
	var unknown types.Hash
	for i := range unknown {
		unknown[i] = byte(i)
	}
	// Create a batch of tests for various scenarios
	limit := uint64(downloader.MaxHeaderFetch)
	tests := []struct {
		query  *GetBlockHeadersData // The query to execute for header retrieval
		expect []types.Hash        // The hashes of the block whose headers are expected
	}{
		// A single random block should be retrievable by hash and number too
		{
			&GetBlockHeadersData{Origin: HashOrNumber{Hash: pm.blockchain.GetBlockByNumber(limit / 2).Hash()}, Amount: 1},
			[]types.Hash{pm.blockchain.GetBlockByNumber(limit / 2).Hash()},
		}, {
			&GetBlockHeadersData{Origin: HashOrNumber{Number: limit / 2}, Amount: 1},
			[]types.Hash{pm.blockchain.GetBlockByNumber(limit / 2).Hash()},
		},
		// Multiple headers should be retrievable in both directions
		{
			&GetBlockHeadersData{Origin: HashOrNumber{Number: limit / 2}, Amount: 3},
			[]types.Hash{
				pm.blockchain.GetBlockByNumber(limit / 2).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 + 1).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 + 2).Hash(),
			},
		}, {
			&GetBlockHeadersData{Origin: HashOrNumber{Number: limit / 2}, Amount: 3, Reverse: true},
			[]types.Hash{
				pm.blockchain.GetBlockByNumber(limit / 2).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 - 1).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 - 2).Hash(),
			},
		},
		// Multiple headers with skip lists should be retrievable
		{
			&GetBlockHeadersData{Origin: HashOrNumber{Number: limit / 2}, Skip: 3, Amount: 3},
			[]types.Hash{
				pm.blockchain.GetBlockByNumber(limit / 2).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 + 4).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 + 8).Hash(),
			},
		}, {
			&GetBlockHeadersData{Origin: HashOrNumber{Number: limit / 2}, Skip: 3, Amount: 3, Reverse: true},
			[]types.Hash{
				pm.blockchain.GetBlockByNumber(limit / 2).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 - 4).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 - 8).Hash(),
			},
		},
		// The chain endpoints should be retrievable
		{
			&GetBlockHeadersData{Origin: HashOrNumber{Number: 0}, Amount: 1},
			[]types.Hash{pm.blockchain.GetBlockByNumber(0).Hash()},
		}, {
			&GetBlockHeadersData{Origin: HashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64()}, Amount: 1},
			[]types.Hash{pm.blockchain.CurrentBlock().Hash()},
		},
		// Ensure protocol limits are honored
		{
			&GetBlockHeadersData{Origin: HashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64() - 1}, Amount: limit + 10, Reverse: true},
			pm.blockchain.GetBlockHashesFromHash(pm.blockchain.CurrentBlock().Hash(), limit),
		},
		// Check that requesting more than available is handled gracefully
		{
			&GetBlockHeadersData{Origin: HashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64() - 4}, Skip: 3, Amount: 3},
			[]types.Hash{
				pm.blockchain.GetBlockByNumber(pm.blockchain.CurrentBlock().NumberU64() - 4).Hash(),
				pm.blockchain.GetBlockByNumber(pm.blockchain.CurrentBlock().NumberU64()).Hash(),
			},
		}, {
			&GetBlockHeadersData{Origin: HashOrNumber{Number: 4}, Skip: 3, Amount: 3, Reverse: true},
			[]types.Hash{
				pm.blockchain.GetBlockByNumber(4).Hash(),
				pm.blockchain.GetBlockByNumber(0).Hash(),
			},
		},
		// Check that requesting more than available is handled gracefully, even if mid skip
		{
			&GetBlockHeadersData{Origin: HashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64() - 4}, Skip: 2, Amount: 3},
			[]types.Hash{
				pm.blockchain.GetBlockByNumber(pm.blockchain.CurrentBlock().NumberU64() - 4).Hash(),
				pm.blockchain.GetBlockByNumber(pm.blockchain.CurrentBlock().NumberU64() - 1).Hash(),
			},
		}, {
			&GetBlockHeadersData{Origin: HashOrNumber{Number: 4}, Skip: 2, Amount: 3, Reverse: true},
			[]types.Hash{
				pm.blockchain.GetBlockByNumber(4).Hash(),
				pm.blockchain.GetBlockByNumber(1).Hash(),
			},
		},
		// Check a corner case where requesting more can iterate past the endpoints
		{
			&GetBlockHeadersData{Origin: HashOrNumber{Number: 2}, Amount: 5, Reverse: true},
			[]types.Hash{
				pm.blockchain.GetBlockByNumber(2).Hash(),
				pm.blockchain.GetBlockByNumber(1).Hash(),
				pm.blockchain.GetBlockByNumber(0).Hash(),
			},
		},
		// Check a corner case where skipping overflow loops back into the chain start
		{
			&GetBlockHeadersData{Origin: HashOrNumber{Hash: pm.blockchain.GetBlockByNumber(3).Hash()}, Amount: 2, Reverse: false, Skip: math.MaxUint64 - 1},
			[]types.Hash{
				pm.blockchain.GetBlockByNumber(3).Hash(),
			},
		},
		// Check a corner case where skipping overflow loops back to the same header
		{
			&GetBlockHeadersData{Origin: HashOrNumber{Hash: pm.blockchain.GetBlockByNumber(1).Hash()}, Amount: 2, Reverse: false, Skip: math.MaxUint64},
			[]types.Hash{
				pm.blockchain.GetBlockByNumber(1).Hash(),
			},
		},
		// Check that non existing headers aren't returned
		{
			&GetBlockHeadersData{Origin: HashOrNumber{Hash: unknown}, Amount: 1},
			[]types.Hash{},
		}, {
			&GetBlockHeadersData{Origin: HashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64() + 1}, Amount: 1},
			[]types.Hash{},
		},
	}
	// Run each of the tests and verify the results against the chain
	for i, tt := range tests {
		// Collect the headers to expect in the response
		var headers block.Headers
		//headers := block.Headers{}
		for _, hash := range tt.expect {
			headers.Headers = append(headers.Headers, pm.blockchain.GetBlockByHash(hash).Header())
		}
		// Send the hash request and verify the response
		p2p.Send(peer.app, 0x03, tt.query)
		if err := p2p.ExpectMsg(peer.app, 0x04, &headers); err != nil {
			t.Errorf("test %d: headers mismatch: %v", i, err)
		}
		// If the test used number origins, repeat with hashes as the too
		if tt.query.Origin.Hash == (types.Hash{}) {
			if origin := pm.blockchain.GetBlockByNumber(tt.query.Origin.Number); origin != nil {
				tt.query.Origin.Hash, tt.query.Origin.Number = origin.Hash(), 0

				p2p.Send(peer.app, 0x03, tt.query)
				if err := p2p.ExpectMsg(peer.app, 0x04, &headers); err != nil {
					t.Errorf("test %d: headers mismatch: %v", i, err)
				}
			}
		}
	}
}

// Tests that block contents can be retrieved from a remote chain based on their hashes.
func TestGetBlockBodies62(t *testing.T) { testGetBlockBodies(t, 62) }
func TestGetBlockBodies63(t *testing.T) { testGetBlockBodies(t, 63) }

func testGetBlockBodies(t *testing.T, protocol int) {
	pm := newTestProtocolManagerMust(t, downloader.FullSync, downloader.MaxBlockFetch+15, nil, nil)
	peer, _ := newTestPeer("peer", protocol, pm, true)
	defer peer.close()

	// Create a batch of tests for various scenarios
	limit := downloader.MaxBlockFetch
	tests := []struct {
		random    int           // Number of blocks to fetch randomly from the chain
		explicit  []types.Hash // Explicitly requested blocks
		available []bool        // Availability of explicitly requested blocks
		expected  int           // Total number of existing blocks to expect
	}{
		{1, nil, nil, 1},                                                         // A single random block should be retrievable
		{10, nil, nil, 10},                                                       // Multiple random blocks should be retrievable
		{limit, nil, nil, limit},                                                 // The maximum possible blocks should be retrievable
		{limit + 1, nil, nil, limit},                                             // No more than the possible block count should be returned
		{0, []types.Hash{pm.blockchain.Genesis().Hash()}, []bool{true}, 1},      // The genesis block should be retrievable
		{0, []types.Hash{pm.blockchain.CurrentBlock().Hash()}, []bool{true}, 1}, // The chains head block should be retrievable
		{0, []types.Hash{{}}, []bool{false}, 0},                                 // A non existent block should not be returned

		// Existing and non-existing blocks interleaved should not cause problems
		{0, []types.Hash{
			{},
			pm.blockchain.GetBlockByNumber(1).Hash(),
			{},
			pm.blockchain.GetBlockByNumber(10).Hash(),
			{},
			pm.blockchain.GetBlockByNumber(100).Hash(),
			{},
		}, []bool{false, true, false, true, false, true, false}, 3},
	}
	// Run each of the tests and verify the results against the chain
	for i, tt := range tests {
		// Collect the hashes to request, and the response to expect
		//hashes, seen := []types.Hashs{}, make(map[int64]bool)
		seen := make(map[int64]bool)
		var hashes types.Hashs
		bodies := BlockBodiesData{}

		for j := 0; j < tt.random; j++ {
			for {
				num := rand.Int63n(int64(pm.blockchain.CurrentBlock().NumberU64()))
				if !seen[num] {
					seen[num] = true

					blk := pm.blockchain.GetBlockByNumber(uint64(num))
					blkHash := blk.Hash()
					rblkHash := &blkHash
					hashes.Hashs = append(hashes.Hashs, rblkHash)
					if len(bodies.Bodys) < tt.expected {
						bodies.Bodys = append(bodies.Bodys, &block.Body{Transactions: blk.Transactions()})
					}
					break
				}
			}
		}
		for j, hash := range tt.explicit {
			rhash := hash
			hashes.Hashs = append(hashes.Hashs, &rhash)
			if tt.available[j] && len(bodies.Bodys) < tt.expected {
				blk := pm.blockchain.GetBlockByHash(hash)
				bodies.Bodys = append(bodies.Bodys, &block.Body{Transactions: blk.Transactions()})
			}
		}
		// Send the hash request and verify the response
		p2p.Send(peer.app, 0x05, &hashes)
		if err := p2p.ExpectMsg(peer.app, 0x06, &bodies); err != nil {
			t.Errorf("test %d: bodies mismatch: %v", i, err)
		}
	}
}

// Tests that the node state database can be retrieved based on hashes.
func TestGetNodeData63(t *testing.T) { testGetNodeData(t, 63) }

func testGetNodeData(t *testing.T, protocol int) {
	// Define three accounts to simulate transactions with
	acc1Key, _ := crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
	acc2Key, _ := crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	acc1Addr := crypto.PubkeyToAddress(acc1Key.PublicKey)
	acc2Addr := crypto.PubkeyToAddress(acc2Key.PublicKey)

	signer := transaction.MakeSigner(testChainConfig, big.NewInt(0))
	// Create a chain generator with some simple transactions (blatantly stolen from @fjl/chain_markets_test)
	generator := func(i int, block *chainmaker.BlockGen) {
		switch i {
		case 0:
			// In block 1, the test bank sends account #1 some bchain coin .
			tx, _ := transaction.SignTx(transaction.NewTransaction(block.TxNonce(testBank), acc1Addr, big.NewInt(10000),0, nil, nil), signer, testBankKey)
			block.AddTx(tx)
		case 1:
			// In block 2, the test bank sends some more bchain coin to account #1.
			// acc1Addr passes it on to account #2.
			tx1, _ := transaction.SignTx(transaction.NewTransaction(block.TxNonce(testBank), acc1Addr, big.NewInt(1000), 0, nil, nil), signer, testBankKey)
			tx2, _ := transaction.SignTx(transaction.NewTransaction(block.TxNonce(acc1Addr), acc2Addr, big.NewInt(1000), 0, nil, nil), signer, acc1Key)
			block.AddTx(tx1)
			block.AddTx(tx2)
		case 2:
			// Block 3 is empty but was produced by account #2.
			block.SetCoinbase(acc2Addr)
			block.SetConsensusData([]byte("yeehaw"))


		}
	}
	// Assemble the test environment
	pm := newTestProtocolManagerMust(t, downloader.FullSync, 4, generator, nil)
	peer, _ := newTestPeer("peer", protocol, pm, true)
	defer peer.close()

	// Fetch for now the entire chain db
	hashes := types.Hashs{}
	for _, key := range pm.chaindb.(*database.MemDatabase).Keys() {
		if len(key) == len(types.Hash{}) {
			hash := types.BytesToHash(key)
			rHash := &hash
			hashes.Hashs = append(hashes.Hashs, rHash)
		}
	}
	p2p.Send(peer.app, 0x0d, &hashes)
	msg, err := peer.app.ReadMsg()
	if err != nil {
		t.Fatalf("failed to read node data response: %v", err)
	}
	if msg.Code != 0x0e {
		t.Fatalf("response packet code mismatch: have %x, want %x", msg.Code, 0x0c)
	}
	var data NodeData
	if err := msg.Decode(&data); err != nil {
		t.Fatalf("failed to decode response node data: %v", err)
	}
	// Verify that all hashes correspond to the requested data, and reconstruct a state tree
	for i, want := range hashes.Hashs {
		if hash := crypto.Keccak256Hash(data.Nodes[i]); hash != *want {
			t.Errorf("data hash mismatch: have %x, want %x", hash, want)
		}
	}
	statedb, _ := database.OpenMemDB()
	for i := 0; i < len(data.Nodes); i++ {
		statedb.Put(hashes.Hashs[i].Bytes(), data.Nodes[i])
	}
	accounts := []types.Address{testBank, acc1Addr, acc2Addr}
	for i := uint64(0); i <= pm.blockchain.CurrentBlock().NumberU64(); i++ {
		trie, _ := state.New(pm.blockchain.GetBlockByNumber(i).Root(), state.NewDatabase(statedb))

		for j, acc := range accounts {
			state, _ := pm.blockchain.State()
			bw := state.GetBalance(acc)
			bh := trie.GetBalance(acc)

			if (bw != nil && bh == nil) || (bw == nil && bh != nil) {
				t.Errorf("test %d, account %d: balance mismatch: have %v, want %v", i, j, bh, bw)
			}
			if bw != nil && bh != nil && bw.Cmp(bw) != 0 {
				t.Errorf("test %d, account %d: balance mismatch: have %v, want %v", i, j, bh, bw)
			}
		}
	}
}

// Tests that the transaction receipts can be retrieved based on hashes.
func TestGetReceipt63(t *testing.T) { testGetReceipt(t, 63) }

func testGetReceipt(t *testing.T, protocol int) {
	// Define three accounts to simulate transactions with
	acc1Key, _ := crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
	acc2Key, _ := crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	acc1Addr := crypto.PubkeyToAddress(acc1Key.PublicKey)
	acc2Addr := crypto.PubkeyToAddress(acc2Key.PublicKey)

	signer := transaction.MakeSigner(testChainConfig, big.NewInt(0))
	// Create a chain generator with some simple transactions (blatantly stolen from @fjl/chain_markets_test)
	generator := func(i int, block *chainmaker.BlockGen) {
		switch i {
		case 0:
			// In block 1, the test bank sends account #1 some bchain coin.
			tx, _ := transaction.SignTx(transaction.NewTransaction(block.TxNonce(testBank), acc1Addr, big.NewInt(10000), 0, nil, nil), signer, testBankKey)
			block.AddTx(tx)
		case 1:
			// In block 2, the test bank sends some more bchain coin to account #1.
			// acc1Addr passes it on to account #2.
			tx1, _ := transaction.SignTx(transaction.NewTransaction(block.TxNonce(testBank), acc1Addr, big.NewInt(1000), 0, nil, nil), signer, testBankKey)
			tx2, _ := transaction.SignTx(transaction.NewTransaction(block.TxNonce(acc1Addr), acc2Addr, big.NewInt(1000), 0, nil, nil), signer, acc1Key)
			block.AddTx(tx1)
			block.AddTx(tx2)
		case 2:
			// Block 3 is empty but was produced by account #2.
			block.SetCoinbase(acc2Addr)
			block.SetConsensusData([]byte("yeehaw"))

		}
	}
	// Assemble the test environment
	pm := newTestProtocolManagerMust(t, downloader.FullSync, 4, generator, nil)
	peer, _ := newTestPeer("peer", protocol, pm, true)
	defer peer.close()

	// Collect the hashes to request, and the response to expect
	receipts :=  transaction.Receipts_s{}
	var hashes types.Hashs
	for i := uint64(0); i <= pm.blockchain.CurrentBlock().NumberU64(); i++ {
		block := pm.blockchain.GetBlockByNumber(i)
		blkHash := block.Hash()
		rblkHash := & blkHash
		hashes.Hashs = append(hashes.Hashs, rblkHash)
		receipts.Receipts_s = append(receipts.Receipts_s, blockchain.GetBlockReceipts(pm.chaindb, block.Hash(), block.NumberU64()))
	}
	// Send the hash request and verify the response
	p2p.Send(peer.app, 0x0f, &hashes)
	if err := p2p.ExpectMsg(peer.app, 0x10, &receipts); err != nil {
		t.Errorf("receipts mismatch: %v", err)
	}
}

