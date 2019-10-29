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
// @File: bloomindexer.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package chainindexer

import (
	"time"

	"bchain.io/utils/database"
	"bchain.io/utils/bloom"
	"bchain.io/common/types"
	"bchain.io/common/bitutil"
	"bchain.io/core/blockchain"
	"bchain.io/core/blockchain/block"
)

const (
	// bloomConfirms is the number of confirmation blocks before a bloom section is
	// considered probably final and its rotated bits are calculated.
	bloomConfirms = 256

	// bloomThrottling is the time to wait between processing two consecutive index
	// sections. It's useful during chain upgrades to prevent disk overload.
	bloomThrottling = 100 * time.Millisecond
)

// BloomIndexer implements a core.ChainIndexer, building up a rotated bloom bits index
// for the bchain header bloom filters, permitting blazing fast filtering.
type BloomIndexer struct {
	size uint64             // section size to generate bloombits for

	db  database.IDatabase  // database instance to write index data and metadata into
	gen *bloom.Generator    // generator to rotate the bloom bits crating the bloom index

	section uint64          // Section is the section number being processed currently
	head    types.Hash      // Head is the hash of the last header processed
}

// NewBloomIndexer returns a chain indexer that generates bloom bits data for the
// canonical chain for fast logs filtering.
func NewBloomIndexer(db database.IDatabase, size uint64) *ChainIndexer {
	backend := &BloomIndexer{
		db:   db,
		size: size,
	}
	table := database.NewTable(db, string(blockchain.BloomBitsIndexPrefix))

	return NewChainIndexer(db, table, backend, size, bloomConfirms, bloomThrottling, "bloombits")
}

// Reset implements core.ChainIndexerBackend, starting a new bloombits index
// section.
func (b *BloomIndexer) Reset(section uint64, lastSectionHead types.Hash) error {
	gen, err := bloom.NewGenerator(uint(b.size))
	b.gen, b.section, b.head = gen, section, types.Hash{}
	return err
}

// Process implements core.ChainIndexerBackend, adding a new header's bloom into
// the index.
func (b *BloomIndexer) Process(header *block.Header) {
	b.gen.AddBloom(uint(header.Number.IntVal.Uint64()-b.section*b.size), header.Bloom)
	b.head = header.Hash()
}

// Commit implements core.ChainIndexerBackend, finalizing the bloom section and
// writing it out into the database.
func (b *BloomIndexer) Commit() error {
	batch := b.db.NewBatch()

	for i := 0; i < types.BloomBitLength; i++ {
		bits, err := b.gen.Bitset(uint(i))
		if err != nil {
			return err
		}
		blockchain.WriteBloomBits(batch, uint(i), b.section, b.head, bitutil.CompressBytes(bits))
	}
	return batch.Write()
}
