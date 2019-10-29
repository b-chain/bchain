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
// @File: protocol.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package bchain

import (
	"bchain.io/core/transaction"
	"bchain.io/core"
	"bchain.io/utils/event"
	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
)

// Constants to match up protocol versions and messages
const (
	bchain63 = 63
)

// Official short name of the protocol used during capability negotiation.
var ProtocolName = "bchain"

// Supported versions of the bchain protocol (first is primary).
var ProtocolVersions = []uint{bchain63}

// Number of implemented message corresponding to different protocol versions.
var ProtocolLengths = []uint64{38}

const ProtocolMaxMsgSize = 10 * 1024 * 1024 // Maximum cap on the size of a protocol message

// bchain protocol message codes
const (
	StatusMsg          = 0x00
	NewBlockHashesMsg  = 0x01
	TxMsg              = 0x02
	GetBlockHeadersMsg = 0x03
	BlockHeadersMsg    = 0x04
	GetBlockBodiesMsg  = 0x05
	BlockBodiesMsg     = 0x06
	NewBlockMsg        = 0x07

	GetNodeDataMsg = 0x0d
	NodeDataMsg    = 0x0e
	GetReceiptsMsg = 0x0f
	ReceiptsMsg    = 0x10

	CsMsg          = 0x20
	BpMsg          = 0x21
	BaMsg          = 0x22

	GetBlockCertificateMsg  = 0x23
	BlockCertificateMsg     = 0x24
)

type errCode int

const (
	ErrMsgTooLarge = iota
	ErrDecode
	ErrInvalidMsgCode
	ErrProtocolVersionMismatch
	ErrNetworkIdMismatch
	ErrGenesisBlockMismatch
	ErrNoStatusMsg
	ErrExtraStatusMsg
	ErrSuspendedPeer
)

func (e errCode) String() string {
	return errorToString[int(e)]
}

// XXX change once legacy code is out
var errorToString = map[int]string{
	ErrMsgTooLarge:             "Message too long",
	ErrDecode:                  "Invalid message",
	ErrInvalidMsgCode:          "Invalid message code",
	ErrProtocolVersionMismatch: "Protocol version mismatch",
	ErrNetworkIdMismatch:       "NetworkId mismatch",
	ErrGenesisBlockMismatch:    "Genesis block mismatch",
	ErrNoStatusMsg:             "No status message",
	ErrExtraStatusMsg:          "Extra status message",
	ErrSuspendedPeer:           "Suspended peer",
}

type txPool interface {
	// AddRemotes should add the given transactions to the pool.
	AddRemotes([]*transaction.Transaction) []error

	// Pending should return pending transactions.
	// The slice should be modifiable by the caller.
	Pending() (map[types.Address]transaction.Transactions, error)

	// SubscribeTxPreEvent should return an event subscription of
	// TxPreEvent and send events to the given channel.
	SubscribeTxPreEvent(chan<- core.TxPreEvent) event.Subscription
}

//go:generate msgp
// statusData is the network packet for the status message.
type StatusData struct {
	ProtocolVersion uint32
	NetworkId       uint64
	Number          *types.BigInt     //current chain height
	CurrentBlock    types.Hash
	GenesisBlock    types.Hash
}

// newBlockHashesData is the network packet for the block announcements.
type NewBlockHashesData []struct {
	Hash   types.Hash // Hash of one particular block being announced
	Number uint64      // Number of one particular block being announced
}

// getBlockHeadersData represents a block header query.
type GetBlockHeadersData struct {
	Origin  HashOrNumber // Block from which to retrieve headers
	Amount  uint64       // Maximum number of headers to retrieve
	Skip    uint64       // Blocks to skip between consecutive headers
	Reverse bool         // Query direction (false = rising towards latest, true = falling towards genesis)
}

// hashOrNumber is a combined field for specifying an origin block.
type HashOrNumber struct {
	Hash   types.Hash // Block hash from which to retrieve headers (excludes Number)
	Number uint64      // Block hash from which to retrieve headers (excludes Hash)
}


// newBlockData is the network packet for the block propagation message.
type NewBlockData struct {
	Block     *block.Block
	Number    *types.BigInt     //block number
}

// blockBodiesData is the network packet for block content distribution.
type BlockBodiesData struct {
	Bodys []*block.Body
}

// blockCertificateData is the network packet for block apos certificates.
type BlockCertificateData struct {
	Certificates [][]byte
}

type NodeData struct {
	Nodes [][]byte
}
