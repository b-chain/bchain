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
// @File: db_operate_test.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package blockchain

import (
	"bytes"
	"math/big"
	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/transaction"
	"bchain.io/utils/crypto/sha3"
	"bchain.io/utils/database"
	"testing"
)

// Tests block header storage and retrieval operations.
func TestHeaderStorage(t *testing.T) {
	db, _ := database.OpenMemDB()

	// Create a test header to move around the database and make sure it's really new
	header := &block.Header{Number: types.NewBigInt(*big.NewInt(42)), Cdata: block.ConsensusData{"test", []byte("test header")}}
	if entry := GetHeader(db, header.Hash(), header.Number.IntVal.Uint64()); entry != nil {
		t.Fatalf("Non existent header returned: %v", entry)
	}
	// Write and verify the header in the database
	if err := WriteHeader(db, header); err != nil {
		t.Fatalf("Failed to write header into database: %v", err)
	}
	if entry := GetHeader(db, header.Hash(), header.Number.IntVal.Uint64()); entry == nil {
		t.Fatalf("Stored header not found")
	} else if entry.Hash() != header.Hash() {
		t.Fatalf("Retrieved header mismatch: have %v, want %v", entry, header)
	}
	if entry := GetHeaderMsgp(db, header.Hash(), header.Number.IntVal.Uint64()); entry == nil {
		t.Fatalf("Stored header MSGP not found")
	} else {
		hasher := sha3.NewKeccak256()
		hasher.Write(entry)

		if hash := types.BytesToHash(hasher.Sum(nil)); hash != header.Hash() {
			t.Fatalf("Retrieved MSGP header mismatch: have %v, want %v", entry, header)
		}
	}
	// Delete the header and verify the execution
	DeleteHeader(db, header.Hash(), header.Number.IntVal.Uint64())
	if entry := GetHeader(db, header.Hash(), header.Number.IntVal.Uint64()); entry != nil {
		t.Fatalf("Deleted header returned: %v", entry)
	}
}
func sh3Hash(x interface{}) (h types.Hash) {
	h3 := sha3.NewKeccak256()
	h3.Write(x.([]byte))
	h3.Sum(h[:0])
	return h
}

// Tests block body storage and retrieval operations.
func TestBodyStorage(t *testing.T) {
	db, _ := database.OpenMemDB()

	// Create a test body to move around the database and make sure it's really new
	body := &block.Body{}

	encData, err := body.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("body enc err: %v", err)
	}
	hash := sh3Hash(encData)

	if entry := GetBody(db, hash, 0); entry != nil {
		t.Fatalf("Non existent body returned: %v", entry)
	}
	// Write and verify the body in the database
	if err := WriteBody(db, hash, 0, body); err != nil {
		t.Fatalf("Failed to write body into database: %v", err)
	}
	if entry := GetBody(db, hash, 0); entry == nil {
		t.Fatalf("Stored body not found")
	} else if block.DeriveSha(transaction.Transactions(entry.Transactions)) != block.DeriveSha(transaction.Transactions(body.Transactions)) {
		t.Fatalf("Retrieved body mismatch: have %v, want %v", entry, body)
	}
	if entry := GetBodyMsgp(db, hash, 0); entry == nil {
		t.Fatalf("Stored body MSGP not found")
	} else {
		hasher := sha3.NewKeccak256()
		hasher.Write(entry)

		if calc := types.BytesToHash(hasher.Sum(nil)); calc != hash {
			t.Fatalf("Retrieved MSGP body mismatch: have %v, want %v", entry, body)
		}
	}
	// Delete the body and verify the execution
	DeleteBody(db, hash, 0)
	if entry := GetBody(db, hash, 0); entry != nil {
		t.Fatalf("Deleted body returned: %v", entry)
	}
}

// Tests block storage and retrieval operations.
func TestBlockStorage(t *testing.T) {
	db, _ := database.OpenMemDB()

	// Create a test block to move around the database and make sure it's really new
	blk := block.NewBlockWithHeader(&block.Header{
		Cdata:           block.ConsensusData{"test", []byte("test block")},
		TxRootHash:      block.EmptyRootHash,
		ReceiptRootHash: block.EmptyRootHash,
	})
	if entry := GetBlock(db, blk.Hash(), blk.NumberU64()); entry != nil {
		t.Fatalf("Non existent block returned: %v", entry)
	}
	if entry := GetHeader(db, blk.Hash(), blk.NumberU64()); entry != nil {
		t.Fatalf("Non existent header returned: %v", entry)
	}
	if entry := GetBody(db, blk.Hash(), blk.NumberU64()); entry != nil {
		t.Fatalf("Non existent body returned: %v", entry)
	}
	// Write and verify the block in the database
	if err := WriteBlock(db, blk); err != nil {
		t.Fatalf("Failed to write block into database: %v", err)
	}
	if entry := GetBlock(db, blk.Hash(), blk.NumberU64()); entry == nil {
		t.Fatalf("Stored block not found")
	} else if entry.Hash() != blk.Hash() {
		t.Fatalf("Retrieved block mismatch: have %v, want %v", entry, blk)
	}
	if entry := GetHeader(db, blk.Hash(), blk.NumberU64()); entry == nil {
		t.Fatalf("Stored header not found")
	} else if entry.Hash() != blk.Header().Hash() {
		t.Fatalf("Retrieved header mismatch: have %v, want %v", entry, blk.Header())
	}
	if entry := GetBody(db, blk.Hash(), blk.NumberU64()); entry == nil {
		t.Fatalf("Stored body not found")
	} else if block.DeriveSha(transaction.Transactions(entry.Transactions)) != block.DeriveSha(blk.Transactions()) {
		t.Fatalf("Retrieved body mismatch: have %v, want %v", entry, blk.Body())
	}
	// Delete the block and verify the execution
	DeleteBlock(db, blk.Hash(), blk.NumberU64())
	if entry := GetBlock(db, blk.Hash(), blk.NumberU64()); entry != nil {
		t.Fatalf("Deleted block returned: %v", entry)
	}
	if entry := GetHeader(db, blk.Hash(), blk.NumberU64()); entry != nil {
		t.Fatalf("Deleted header returned: %v", entry)
	}
	if entry := GetBody(db, blk.Hash(), blk.NumberU64()); entry != nil {
		t.Fatalf("Deleted body returned: %v", entry)
	}
}

// Tests that partial block contents don't get reassembled into full blocks.
func TestPartialBlockStorage(t *testing.T) {
	db, _ := database.OpenMemDB()
	block := block.NewBlockWithHeader(&block.Header{
		Cdata:           block.ConsensusData{"test", []byte("test block")},
		TxRootHash:      block.EmptyRootHash,
		ReceiptRootHash: block.EmptyRootHash,
	})
	// Store a header and check that it's not recognized as a block
	if err := WriteHeader(db, block.Header()); err != nil {
		t.Fatalf("Failed to write header into database: %v", err)
	}
	if entry := GetBlock(db, block.Hash(), block.NumberU64()); entry != nil {
		t.Fatalf("Non existent block returned: %v", entry)
	}
	DeleteHeader(db, block.Hash(), block.NumberU64())

	// Store a body and check that it's not recognized as a block
	if err := WriteBody(db, block.Hash(), block.NumberU64(), block.Body()); err != nil {
		t.Fatalf("Failed to write body into database: %v", err)
	}
	if entry := GetBlock(db, block.Hash(), block.NumberU64()); entry != nil {
		t.Fatalf("Non existent block returned: %v", entry)
	}
	DeleteBody(db, block.Hash(), block.NumberU64())

	// Store a header and a body separately and check reassembly
	if err := WriteHeader(db, block.Header()); err != nil {
		t.Fatalf("Failed to write header into database: %v", err)
	}
	if err := WriteBody(db, block.Hash(), block.NumberU64(), block.Body()); err != nil {
		t.Fatalf("Failed to write body into database: %v", err)
	}
	if entry := GetBlock(db, block.Hash(), block.NumberU64()); entry == nil {
		t.Fatalf("Stored block not found")
	} else if entry.Hash() != block.Hash() {
		t.Fatalf("Retrieved block mismatch: have %v, want %v", entry, block)
	}
}

// Tests that canonical numbers can be mapped to hashes and retrieved.
func TestCanonicalMappingStorage(t *testing.T) {
	db, _ := database.OpenMemDB()

	// Create a test canonical number and assinged hash to move around
	hash, number := types.Hash{0: 0xff}, uint64(314)
	if entry := GetCanonicalHash(db, number); entry != (types.Hash{}) {
		t.Fatalf("Non existent canonical mapping returned: %v", entry)
	}
	// Write and verify the TD in the database
	if err := WriteCanonicalHash(db, hash, number); err != nil {
		t.Fatalf("Failed to write canonical mapping into database: %v", err)
	}
	if entry := GetCanonicalHash(db, number); entry == (types.Hash{}) {
		t.Fatalf("Stored canonical mapping not found")
	} else if entry != hash {
		t.Fatalf("Retrieved canonical mapping mismatch: have %v, want %v", entry, hash)
	}
	// Delete the TD and verify the execution
	DeleteCanonicalHash(db, number)
	if entry := GetCanonicalHash(db, number); entry != (types.Hash{}) {
		t.Fatalf("Deleted canonical mapping returned: %v", entry)
	}
}

// Tests that head headers and head blocks can be assigned, individually.
func TestHeadStorage(t *testing.T) {
	db, _ := database.OpenMemDB()

	blockHead := block.NewBlockWithHeader(&block.Header{Cdata: block.ConsensusData{"test", []byte("test block header")}})
	blockFull := block.NewBlockWithHeader(&block.Header{Cdata: block.ConsensusData{"test", []byte("test block full")}})
	blockFast := block.NewBlockWithHeader(&block.Header{Cdata: block.ConsensusData{"test", []byte("test block fast")}})

	// Check that no head entries are in a pristine database
	if entry := GetHeadHeaderHash(db); entry != (types.Hash{}) {
		t.Fatalf("Non head header entry returned: %v", entry)
	}
	if entry := GetHeadBlockHash(db); entry != (types.Hash{}) {
		t.Fatalf("Non head block entry returned: %v", entry)
	}
	if entry := GetHeadFastBlockHash(db); entry != (types.Hash{}) {
		t.Fatalf("Non fast head block entry returned: %v", entry)
	}
	// Assign separate entries for the head header and block
	if err := WriteHeadHeaderHash(db, blockHead.Hash()); err != nil {
		t.Fatalf("Failed to write head header hash: %v", err)
	}
	if err := WriteHeadBlockHash(db, blockFull.Hash()); err != nil {
		t.Fatalf("Failed to write head block hash: %v", err)
	}
	if err := WriteHeadFastBlockHash(db, blockFast.Hash()); err != nil {
		t.Fatalf("Failed to write fast head block hash: %v", err)
	}
	// Check that both heads are present, and different (i.e. two heads maintained)
	if entry := GetHeadHeaderHash(db); entry != blockHead.Hash() {
		t.Fatalf("Head header hash mismatch: have %v, want %v", entry, blockHead.Hash())
	}
	if entry := GetHeadBlockHash(db); entry != blockFull.Hash() {
		t.Fatalf("Head block hash mismatch: have %v, want %v", entry, blockFull.Hash())
	}
	if entry := GetHeadFastBlockHash(db); entry != blockFast.Hash() {
		t.Fatalf("Fast head block hash mismatch: have %v, want %v", entry, blockFast.Hash())
	}
}

// Tests that positional lookup metadata can be stored and retrieved.
func TestLookupStorage(t *testing.T) {
	db, _ := database.OpenMemDB()

	tx1 := transaction.NewTransaction(1, transaction.Actions{&transaction.Action{Contract: types.BytesToAddress([]byte{0x11}), Params: []byte{0x11, 0x11, 0x11}}})
	tx2 := transaction.NewTransaction(2, transaction.Actions{&transaction.Action{Contract: types.BytesToAddress([]byte{0x22}), Params: []byte{0x22, 0x22, 0x22}}})
	tx3 := transaction.NewTransaction(3, transaction.Actions{&transaction.Action{Contract: types.BytesToAddress([]byte{0x33}), Params: []byte{0x33, 0x33, 0x33}}})
	txs := []*transaction.Transaction{tx1, tx2, tx3}

	block := block.NewBlock(&block.Header{Number: types.NewBigInt(*big.NewInt(314))}, txs, nil)

	// Check that no transactions entries are in a pristine database
	for i, tx := range txs {
		if txn, _, _, _ := GetTransaction(db, tx.Hash()); txn != nil {
			t.Fatalf("tx #%d [%x]: non existent transaction returned: %v", i, tx.Hash(), txn)
		}
	}
	// Insert all the transactions into the database, and verify contents
	if err := WriteBlock(db, block); err != nil {
		t.Fatalf("failed to write block contents: %v", err)
	}
	if err := WriteTxLookupEntries(db, block); err != nil {
		t.Fatalf("failed to write transactions: %v", err)
	}
	for i, tx := range txs {
		if txn, hash, number, index := GetTransaction(db, tx.Hash()); txn == nil {
			t.Fatalf("tx #%d [%x]: transaction not found", i, tx.Hash())
		} else {
			if hash != block.Hash() || number != block.NumberU64() || index != uint64(i) {
				t.Fatalf("tx #%d [%x]: positional metadata mismatch: have %x/%d/%d, want %x/%v/%v", i, tx.Hash(), hash, number, index, block.Hash(), block.NumberU64(), i)
			}
			if tx.String() != txn.String() {
				t.Fatalf("tx #%d [%x]: transaction mismatch: have %v, want %v", i, tx.Hash(), txn, tx)
			}
		}
	}
	// Delete the transactions and check purge
	for i, tx := range txs {
		DeleteTxLookupEntry(db, tx.Hash())
		if txn, _, _, _ := GetTransaction(db, tx.Hash()); txn != nil {
			t.Fatalf("tx #%d [%x]: deleted transaction returned: %v", i, tx.Hash(), txn)
		}
	}
}

// Tests that receipts associated with a single block can be stored and retrieved.
func TestBlockReceiptStorage(t *testing.T) {
	db, _ := database.OpenMemDB()

	receipt1 := &transaction.Receipt{
		Status: transaction.ReceiptStatusFailed,
		Logs: []*transaction.Log{
			{Address: types.BytesToAddress([]byte{0x11})},
			{Address: types.BytesToAddress([]byte{0x01, 0x11})},
		},
		TxHash: types.BytesToHash([]byte{0x11, 0x11}),
		ContractAddress: []types.Address{
			types.BytesToAddress([]byte{0x01, 0x11, 0x11}),
		},
	}
	receipt2 := &transaction.Receipt{
		//PostState: types.Hash{2}.Bytes(),
		Logs: []*transaction.Log{
			{Address: types.BytesToAddress([]byte{0x22})},
			{Address: types.BytesToAddress([]byte{0x02, 0x22})},
		},
		TxHash: types.BytesToHash([]byte{0x22, 0x22}),
		ContractAddress: []types.Address{
			types.BytesToAddress([]byte{0x02, 0x22, 0x22}),
		},
	}
	receipts := []*transaction.Receipt{receipt1, receipt2}

	// Check that no receipt entries are in a pristine database
	hash := types.BytesToHash([]byte{0x03, 0x14})
	if rs := GetBlockReceipts(db, hash, 0); len(rs) != 0 {
		t.Fatalf("non existent receipts returned: %v", rs)
	}
	// Insert the receipt slice into the database and check presence
	if err := WriteBlockReceipts(db, hash, 0, receipts); err != nil {
		t.Fatalf("failed to write block receipts: %v", err)
	}
	if rs := GetBlockReceipts(db, hash, 0); len(rs) == 0 {
		t.Fatalf("no receipts returned")
	} else {
		for i := 0; i < len(receipts); i++ {
			//var bufHava bytes.Buffer
			//var bufWant bytes.Buffer
			msgpHave, _ := rs[i].MarshalMsg(nil)
			msgpWant, _ := receipts[i].MarshalMsg(nil)

			if !bytes.Equal(msgpHave, msgpWant) {
				t.Fatalf("receipt #%d: receipt mismatch: have %v, want %v", i, rs[i], receipts[i])
			}
		}
	}
	// Delete the receipt slice and check purge
	DeleteBlockReceipts(db, hash, 0)
	if rs := GetBlockReceipts(db, hash, 0); len(rs) != 0 {
		t.Fatalf("deleted receipts returned: %v", rs)
	}
}
