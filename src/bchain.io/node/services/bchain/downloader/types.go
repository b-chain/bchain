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
// @File: types.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package downloader

import (
	"fmt"

	"bchain.io/core/blockchain/block"
	"bchain.io/core/transaction"
)

// peerDropFn is a callback type for dropping a peer detected as malicious.
type peerDropFn func(id string)

type certificateCheckFn func(header *block.Header, data []byte) error

// dataPack is a data message returned by a peer for some query.
type dataPack interface {
	PeerId() string
	Items() int
	Stats() string
}

// headerPack is a batch of block headers returned by a peer.
type headerPack struct {
	peerId  string
	headers []*block.Header
}

func (p *headerPack) PeerId() string { return p.peerId }
func (p *headerPack) Items() int     { return len(p.headers) }
func (p *headerPack) Stats() string  { return fmt.Sprintf("%d", len(p.headers)) }

// bodyPack is a batch of block bodies returned by a peer.
type bodyPack struct {
	peerId       string
	transactions [][]*transaction.Transaction
}

func (p *bodyPack) PeerId() string { return p.peerId }
func (p *bodyPack) Items() int {
	return len(p.transactions)
}
func (p *bodyPack) Stats() string { return fmt.Sprintf("%d", len(p.transactions)) }

// certificatePack is a batch of block certificates returned by a peer.
type certificatePack struct {
	peerId       string
	certificates [][]byte
}

func (p *certificatePack) PeerId() string { return p.peerId }
func (p *certificatePack) Items() int { return len(p.certificates) }
func (p *certificatePack) Stats() string { return fmt.Sprintf("%d", len(p.certificates)) }

// receiptPack is a batch of receipts returned by a peer.
type receiptPack struct {
	peerId   string
	receipts [][]*transaction.Receipt
}

func (p *receiptPack) PeerId() string { return p.peerId }
func (p *receiptPack) Items() int     { return len(p.receipts) }
func (p *receiptPack) Stats() string  { return fmt.Sprintf("%d", len(p.receipts)) }

// statePack is a batch of states returned by a peer.
type statePack struct {
	peerId string
	states [][]byte
}

func (p *statePack) PeerId() string { return p.peerId }
func (p *statePack) Items() int     { return len(p.states) }
func (p *statePack) Stats() string  { return fmt.Sprintf("%d", len(p.states)) }
