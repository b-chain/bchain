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
// @File: log.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package transaction

import (
	"fmt"
	"bchain.io/common/types"
)

//go:generate msgp
//go:generate gencodec -type Log -field-override logMarshaling -out gen_log_json.go

// Log represents a contract log event. These events are generated by the LOG opcode and
// stored/indexed by the node.
type Log struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address types.Address `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	Topics []types.Hash `json:"topics" gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	Data [][]byte `json:"data" gencodec:"required"`

	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	BlockNumber uint64 `json:"blockNumber"`
	// hash of the transaction
	TxHash types.Hash `json:"transactionHash" gencodec:"required"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex" gencodec:"required"`
	// hash of the block in which the transaction was included
	BlockHash types.Hash `json:"blockHash"`
	// index of the log in the receipt
	Index uint `json:"logIndex" gencodec:"required"`

	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool `json:"removed"  msg:"-"`
}

type logMarshaling struct {
	Data        []types.BytesForJson
	BlockNumber types.Uint64ForJson
	TxIndex     types.UintForJson
	Index       types.UintForJson
}

func (l *Log) String() string {
	return fmt.Sprintf(`log: %x %x %x %x %d %x %d`, l.Address, l.Topics, l.Data, l.TxHash, l.TxIndex, l.BlockHash, l.Index)
}

// LogProtocol is the consensus encoding of a Log.
type LogProtocol struct {
	Address types.Address
	// list of topics provided by the contract.
	Topics []types.Hash
	// supplied by the contract, usually ABI-encoded
	Data [][]byte
}
