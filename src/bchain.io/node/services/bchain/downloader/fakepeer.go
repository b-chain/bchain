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
// @File: fakepeer.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package downloader

import (
	"math/big"

	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/transaction"
	"bchain.io/utils/database"
	"bchain.io/core/blockchain"
)

// FakePeer is a mock downloader peer that operates on a local database instance
// instead of being an actual live node. It's useful for testing and to implement
// sync commands from an xisting local database.
type FakePeer struct {
	id string
	db database.IDatabase

	hc *blockchain.HeaderChain
	dl *Downloader
}

// NewFakePeer creates a new mock downloader peer with the given data sources.
func NewFakePeer(id string, db database.IDatabase, hc *blockchain.HeaderChain, dl *Downloader) *FakePeer {
	return &FakePeer{id: id, db: db, hc: hc, dl: dl}
}

// Head implements downloader.Peer, returning the current head hash and number
// of the best known header.
func (p *FakePeer) Head() (types.Hash, *big.Int) {
	header := p.hc.CurrentHeader()
	return header.Hash(), &header.Number.IntVal
}

// RequestHeadersByHash implements downloader.Peer, returning a batch of headers
// defined by the origin hash and the associaed query parameters.
func (p *FakePeer) RequestHeadersByHash(hash types.Hash, amount int, skip int, reverse bool) error {
	var (
		headers []*block.Header
		unknown bool
	)
	for !unknown && len(headers) < amount {
		origin := p.hc.GetHeaderByHash(hash)
		if origin == nil {
			break
		}
		number := origin.Number.IntVal.Uint64()
		headers = append(headers, origin)
		if reverse {
			for i := 0; i <= skip; i++ {
				if header := p.hc.GetHeader(hash, number); header != nil {
					hash = header.ParentHash
					number--
				} else {
					unknown = true
					break
				}
			}
		} else {
			var (
				current = origin.Number.IntVal.Uint64()
				next    = current + uint64(skip) + 1
			)
			if header := p.hc.GetHeaderByNumber(next); header != nil {
				if p.hc.GetBlockHashesFromHash(header.Hash(), uint64(skip+1))[skip] == hash {
					hash = header.Hash()
				} else {
					unknown = true
				}
			} else {
				unknown = true
			}
		}
	}
	p.dl.DeliverHeaders(p.id, headers)
	return nil
}

// RequestHeadersByNumber implements downloader.Peer, returning a batch of headers
// defined by the origin number and the associaed query parameters.
func (p *FakePeer) RequestHeadersByNumber(number uint64, amount int, skip int, reverse bool) error {
	var (
		headers []*block.Header
		unknown bool
	)
	for !unknown && len(headers) < amount {
		origin := p.hc.GetHeaderByNumber(number)
		if origin == nil {
			break
		}
		if reverse {
			if number >= uint64(skip+1) {
				number -= uint64(skip + 1)
			} else {
				unknown = true
			}
		} else {
			number += uint64(skip + 1)
		}
		headers = append(headers, origin)
	}
	p.dl.DeliverHeaders(p.id, headers)
	return nil
}

// RequestBodies implements downloader.Peer, returning a batch of block bodies
// corresponding to the specified block hashes.
func (p *FakePeer) RequestBodies(hashes []types.Hash) error {
	var (
		txs    [][]*transaction.Transaction
	)
	for _, hash := range hashes {
		block := blockchain.GetBlock(p.db, hash, p.hc.GetBlockNumber(hash))

		txs = append(txs, block.Transactions())

	}
	p.dl.DeliverBodies(p.id, txs)
	return nil
}

// RequestReceipts implements downloader.Peer, returning a batch of transaction
// receipts corresponding to the specified block hashes.
func (p *FakePeer) RequestReceipts(hashes []types.Hash) error {
	var receipts [][]*transaction.Receipt
	for _, hash := range hashes {
		receipts = append(receipts, blockchain.GetBlockReceipts(p.db, hash, p.hc.GetBlockNumber(hash)))
	}
	p.dl.DeliverReceipts(p.id, receipts)
	return nil
}

// RequestNodeData implements downloader.Peer, returning a batch of state trie
// nodes corresponding to the specified trie hashes.
func (p *FakePeer) RequestNodeData(hashes []types.Hash) error {
	var data [][]byte
	for _, hash := range hashes {
		if entry, err := p.db.Get(hash.Bytes()); err == nil {
			data = append(data, entry)
		}
	}
	p.dl.DeliverNodeData(p.id, data)
	return nil
}
