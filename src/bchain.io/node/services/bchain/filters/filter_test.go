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
// @File: filter_test.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package filters

import (
	"context"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"bchain.io/common/types"
	"bchain.io/consensus"
	"bchain.io/core/blockchain"
	"bchain.io/core/blockchain/chainmaker"
	"bchain.io/core/genesis"
	"bchain.io/core/transaction"
	"bchain.io/params"
	"bchain.io/utils/bloom"
	"bchain.io/utils/crypto"
	"bchain.io/utils/database"
	"bchain.io/utils/event"
)

func makeReceipt(addr types.Address) *transaction.Receipt {
	receipt := transaction.NewReceipt(false)
	receipt.Logs = []*transaction.Log{
		{Address: addr},
	}

	topics := []bloom.BloomByte{}
	for _, log := range receipt.Logs {
		topics = append(topics, log.Address)
		for _, topic := range log.Topics {
			topics = append(topics, topic)
		}
	}
	receipt.Bloom = bloom.CreateBloom(topics)
	return receipt
}

func BenchmarkFilters(b *testing.B) {
	dir, err := ioutil.TempDir("", "filtertest")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	var (
		db, _      = database.OpenLDB(dir, 0, 0)
		mux        = new(event.TypeMux)
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed}
		key1, _    = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1      = crypto.PubkeyToAddress(key1.PublicKey)
		addr2      = types.BytesToAddress([]byte("jeff"))
		addr3      = types.BytesToAddress([]byte("bchain"))
		addr4      = types.BytesToAddress([]byte("random addresses please"))
	)
	defer db.Close()

	genesis := genesis.GenesisBlockForTesting(db, addr1, big.NewInt(1000000))
	chain, receipts := chainmaker.GenerateChain(params.TestChainConfig, genesis, &consensus.Engine_empty{}, db, 100010, func(i int, gen *chainmaker.BlockGen) {
		switch i {
		case 2403:
			receipt := makeReceipt(addr1)
			gen.AddUncheckedReceipt(receipt)
		case 1034:
			receipt := makeReceipt(addr2)
			gen.AddUncheckedReceipt(receipt)
		case 34:
			receipt := makeReceipt(addr3)
			gen.AddUncheckedReceipt(receipt)
		case 99999:
			receipt := makeReceipt(addr4)
			gen.AddUncheckedReceipt(receipt)

		}
	})
	for i, block := range chain {
		blockchain.WriteBlock(db, block)
		if err := blockchain.WriteCanonicalHash(db, block.Hash(), block.NumberU64()); err != nil {
			b.Fatalf("failed to insert block number: %v", err)
		}
		if err := blockchain.WriteHeadBlockHash(db, block.Hash()); err != nil {
			b.Fatalf("failed to insert block number: %v", err)
		}
		if err := blockchain.WriteBlockReceipts(db, block.Hash(), block.NumberU64(), receipts[i]); err != nil {
			b.Fatal("error writing block receipts:", err)
		}
	}
	b.ResetTimer()

	filter := New(backend, 0, -1, []types.Address{addr1, addr2, addr3, addr4}, nil)

	for i := 0; i < b.N; i++ {
		logs, _ := filter.Logs(context.Background())
		if len(logs) != 4 {
			b.Fatal("expected 4 logs, got", len(logs))
		}
	}
}

func TestFilters(t *testing.T) {
	dir, err := ioutil.TempDir("", "filtertest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	var (
		db, _      = database.OpenLDB(dir, 0, 0)
		mux        = new(event.TypeMux)
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed}
		key1, _    = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr       = crypto.PubkeyToAddress(key1.PublicKey)

		hash1 = types.BytesToHash([]byte("topic1"))
		hash2 = types.BytesToHash([]byte("topic2"))
		hash3 = types.BytesToHash([]byte("topic3"))
		hash4 = types.BytesToHash([]byte("topic4"))
	)
	defer db.Close()

	genesis := genesis.GenesisBlockForTesting(db, addr, big.NewInt(1000000))
	chain, receipts := chainmaker.GenerateChain(params.TestChainConfig, genesis, &consensus.Engine_empty{}, db, 1000, func(i int, gen *chainmaker.BlockGen) {
		switch i {
		case 1:
			receipt := transaction.NewReceipt(false)
			receipt.Logs = []*transaction.Log{
				{
					Address: addr,
					Topics:  []types.Hash{hash1},
				},
			}
			gen.AddUncheckedReceipt(receipt)
		case 2:
			receipt := transaction.NewReceipt(false)
			receipt.Logs = []*transaction.Log{
				{
					Address: addr,
					Topics:  []types.Hash{hash2},
				},
			}
			gen.AddUncheckedReceipt(receipt)
		case 998:
			receipt := transaction.NewReceipt(false)
			receipt.Logs = []*transaction.Log{
				{
					Address: addr,
					Topics:  []types.Hash{hash3},
				},
			}
			gen.AddUncheckedReceipt(receipt)
		case 999:
			receipt := transaction.NewReceipt(false)
			receipt.Logs = []*transaction.Log{
				{
					Address: addr,
					Topics:  []types.Hash{hash4},
				},
			}
			gen.AddUncheckedReceipt(receipt)
		}
	})
	for i, block := range chain {
		blockchain.WriteBlock(db, block)
		if err := blockchain.WriteCanonicalHash(db, block.Hash(), block.NumberU64()); err != nil {
			t.Fatalf("failed to insert block number: %v", err)
		}
		if err := blockchain.WriteHeadBlockHash(db, block.Hash()); err != nil {
			t.Fatalf("failed to insert block number: %v", err)
		}
		if err := blockchain.WriteBlockReceipts(db, block.Hash(), block.NumberU64(), receipts[i]); err != nil {
			t.Fatal("error writing block receipts:", err)
		}
	}

	filter := New(backend, 0, -1, []types.Address{addr}, [][]types.Hash{{hash1, hash2, hash3, hash4}})

	logs, _ := filter.Logs(context.Background())
	if len(logs) != 4 {
		t.Error("expected 4 log, got", len(logs))
	}

	filter = New(backend, 900, 999, []types.Address{addr}, [][]types.Hash{{hash3}})
	logs, _ = filter.Logs(context.Background())
	if len(logs) != 1 {
		t.Error("expected 1 log, got", len(logs))
	}
	if len(logs) > 0 && logs[0].Topics[0] != hash3 {
		t.Errorf("expected log[0].Topics[0] to be %x, got %x", hash3, logs[0].Topics[0])
	}

	filter = New(backend, 990, -1, []types.Address{addr}, [][]types.Hash{{hash3}})
	logs, _ = filter.Logs(context.Background())
	if len(logs) != 1 {
		t.Error("expected 1 log, got", len(logs))
	}
	if len(logs) > 0 && logs[0].Topics[0] != hash3 {
		t.Errorf("expected log[0].Topics[0] to be %x, got %x", hash3, logs[0].Topics[0])
	}

	filter = New(backend, 1, 10, nil, [][]types.Hash{{hash1, hash2}})

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 2 {
		t.Error("expected 2 log, got", len(logs))
	}

	failHash := types.BytesToHash([]byte("fail"))
	filter = New(backend, 0, -1, nil, [][]types.Hash{{failHash}})

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 0 {
		t.Error("expected 0 log, got", len(logs))
	}

	failAddr := types.BytesToAddress([]byte("failmenow"))
	filter = New(backend, 0, -1, []types.Address{failAddr}, nil)

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 0 {
		t.Error("expected 0 log, got", len(logs))
	}

	filter = New(backend, 0, -1, nil, [][]types.Hash{{failHash}, {hash1}})

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 0 {
		t.Error("expected 0 log, got", len(logs))
	}
}
