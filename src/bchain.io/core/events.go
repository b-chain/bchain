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
// @File: events.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package core

import (
	"bchain.io/core/blockchain/block"
	"bchain.io/core/transaction"
	"bchain.io/common/types"
)

// TxPreEvent is posted when a transaction enters the transaction pool.
type TxPreEvent struct{ Tx *transaction.Transaction}

// PendingLogsEvent is posted pre producing and notifies of pending logs.
type PendingLogsEvent struct {
	Logs []*transaction.Log
}

// PendingStateEvent is posted pre producing and notifies of pending state changes.
type PendingStateEvent struct{}

// NewProducedBlockEvent is posted when a block has been imported.
type NewProducedBlockEvent struct{ Block *block.Block }

// RemovedTransactionEvent is posted when a reorg happens
type RemovedTransactionEvent struct{ Txs transaction.Transactions }

// RemovedLogsEvent is posted when a reorg happens
type RemovedLogsEvent struct{ Logs []*transaction.Log }

type ChainEvent struct {
	Block *block.Block
	Hash  types.Hash
	Logs  []*transaction.Log
}

type ChainSideEvent struct {
	Block *block.Block
}

type ChainHeadEvent struct{ Block *block.Block }
