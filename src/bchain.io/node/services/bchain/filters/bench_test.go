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
// @File: bench_test.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package filters

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"bchain.io/common/bitutil"
	"bchain.io/common/types"
	"bchain.io/core/blockchain"
	"bchain.io/core/blockchain/block"
	"bchain.io/node"
	"bchain.io/utils/bloom"
	"bchain.io/utils/database"
	"bchain.io/utils/event"
)

func BenchmarkBloomBits512(b *testing.B) {
	benchmarkBloomBits(b, 512)
}

func BenchmarkBloomBits1k(b *testing.B) {
	benchmarkBloomBits(b, 1024)
}

func BenchmarkBloomBits2k(b *testing.B) {
	benchmarkBloomBits(b, 2048)
}

func BenchmarkBloomBits4k(b *testing.B) {
	benchmarkBloomBits(b, 4096)
}

func BenchmarkBloomBits8k(b *testing.B) {
	benchmarkBloomBits(b, 8192)
}

func BenchmarkBloomBits16k(b *testing.B) {
	benchmarkBloomBits(b, 16384)
}

func BenchmarkBloomBits32k(b *testing.B) {
	benchmarkBloomBits(b, 32768)
}

const benchFilterCnt = 2000

func benchmarkBloomBits(b *testing.B, sectionSize uint64) {
	benchDataDir := node.DefaultDataDir() + "/bchaind/chaindata"
	fmt.Println("Running bloombits benchmark   section size:", sectionSize)

	db, err := database.OpenLDB(benchDataDir, 128, 1024)
	if err != nil {
		b.Fatalf("error opening database at %v: %v", benchDataDir, err)
	}
	head := blockchain.GetHeadBlockHash(db)
	if head == (types.Hash{}) {
		b.Fatalf("chain data not found at %v", benchDataDir)
	}

	clearBloomBits(db)
	fmt.Println("Generating bloombits data...")
	headNum := blockchain.GetBlockNumber(db, head)
	if headNum < sectionSize+512 {
		b.Fatalf("not enough blocks for running a benchmark")
	}

	start := time.Now()
	cnt := (headNum - 512) / sectionSize
	var dataSize, compSize uint64
	for sectionIdx := uint64(0); sectionIdx < cnt; sectionIdx++ {
		bc, err := bloom.NewGenerator(uint(sectionSize))
		if err != nil {
			b.Fatalf("failed to create generator: %v", err)
		}
		var header *block.Header
		for i := sectionIdx * sectionSize; i < (sectionIdx+1)*sectionSize; i++ {
			hash := blockchain.GetCanonicalHash(db, i)
			header = blockchain.GetHeader(db, hash, i)
			if header == nil {
				b.Fatalf("Error creating bloomBits data")
			}
			bc.AddBloom(uint(i-sectionIdx*sectionSize), header.Bloom)
		}
		sectionHead := blockchain.GetCanonicalHash(db, (sectionIdx+1)*sectionSize-1)
		for i := 0; i < types.BloomBitLength; i++ {
			data, err := bc.Bitset(uint(i))
			if err != nil {
				b.Fatalf("failed to retrieve bitset: %v", err)
			}
			comp := bitutil.CompressBytes(data)
			dataSize += uint64(len(data))
			compSize += uint64(len(comp))
			blockchain.WriteBloomBits(db, uint(i), sectionIdx, sectionHead, comp)
		}
		//if sectionIdx%50 == 0 {
		//	fmt.Println(" section", sectionIdx, "/", cnt)
		//}
	}

	d := time.Since(start)
	fmt.Println("Finished generating bloombits data")
	fmt.Println(" ", d, "total  ", d/time.Duration(cnt*sectionSize), "per block")
	fmt.Println(" data size:", dataSize, "  compressed size:", compSize, "  compression ratio:", float64(compSize)/float64(dataSize))

	fmt.Println("Running filter benchmarks...")
	start = time.Now()
	mux := new(event.TypeMux)
	var backend *testBackend

	for i := 0; i < benchFilterCnt; i++ {
		if i%20 == 0 {
			db.Close()
			db, _ = database.OpenLDB(benchDataDir, 128, 1024)
			backend = &testBackend{mux, db, cnt, new(event.Feed), new(event.Feed), new(event.Feed), new(event.Feed)}
		}
		var addr types.Address
		addr[0] = byte(i)
		addr[1] = byte(i / 256)
		filter := New(backend, 0, int64(cnt*sectionSize-1), []types.Address{addr}, nil)
		if _, err := filter.Logs(context.Background()); err != nil {
			b.Error("filter.Find error:", err)
		}
	}
	d = time.Since(start)
	fmt.Println("Finished running filter benchmarks")
	fmt.Println(" ", d, "total  ", d/time.Duration(benchFilterCnt), "per address", d*time.Duration(1000000)/time.Duration(benchFilterCnt*cnt*sectionSize), "per million blocks")
	db.Close()
}

func forEachKey(db database.IDatabase, startPrefix, endPrefix []byte, fn func(key []byte)) {
	it := db.(*database.LDatabase).NewIterator()
	it.Seek(startPrefix)
	for it.Valid() {
		key := it.Key()
		cmpLen := len(key)
		if len(endPrefix) < cmpLen {
			cmpLen = len(endPrefix)
		}
		if bytes.Compare(key[:cmpLen], endPrefix) == 1 {
			break
		}
		fn(types.CopyBytes(key))
		it.Next()
	}
	it.Release()
}

var bloomBitsPrefix = []byte("bloomBits-")

func clearBloomBits(db database.IDatabase) {
	fmt.Println("Clearing bloombits data...")
	forEachKey(db, bloomBitsPrefix, bloomBitsPrefix, func(key []byte) {
		db.Delete(key)
	})
}

func BenchmarkNoBloomBits(b *testing.B) {
	benchDataDir := node.DefaultDataDir() + "/bchaind/chaindata"
	fmt.Println("Running benchmark without bloombits")
	db, err := database.OpenLDB(benchDataDir, 128, 1024)
	if err != nil {
		b.Fatalf("error opening database at %v: %v", benchDataDir, err)
	}
	head := blockchain.GetHeadBlockHash(db)
	if head == (types.Hash{}) {
		b.Fatalf("chain data not found at %v", benchDataDir)
	}
	headNum := blockchain.GetBlockNumber(db, head)

	clearBloomBits(db)

	fmt.Println("Running filter benchmarks...")
	start := time.Now()
	mux := new(event.TypeMux)
	backend := &testBackend{mux, db, 0, new(event.Feed), new(event.Feed), new(event.Feed), new(event.Feed)}
	filter := New(backend, 0, int64(headNum), []types.Address{{}}, nil)
	filter.Logs(context.Background())
	d := time.Since(start)
	fmt.Println("Finished running filter benchmarks")
	fmt.Println(" ", d, "total  ", d*time.Duration(1000000)/time.Duration(headNum+1), "per million blocks")
	db.Close()
}
