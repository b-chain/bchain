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
// @File: interfaces.go
// @Date: 2018/05/07 11:22:07
////////////////////////////////////////////////////////////////////////////////

package common

import (
	"context"
	"math/big"
	"bchain.io/common/types"
)

type Subscription interface {
	// Unsubscribe cancels the sending of events to the data channel
	// and closes the error channel.
	Unsubscribe()
	// Err returns the subscription error channel. The error channel receives
	// a value if there is an issue with the subscription (e.g. the network connection
	// delivering the events has been closed). Only one value will ever be sent.
	// The error channel is closed by Unsubscribe.
	Err() <-chan error
}

// ChainStateReader wraps access to the state trie of the canonical blockchain. Note that
// implementations of the interface may be unable to return state values for old blocks.
// In many cases, using CallContract can be preferable to reading raw contract storage.
type ChainStateReader interface {
	BalanceAt(ctx context.Context, account types.Address, blockNumber *big.Int) (*big.Int, error)
	StorageAt(ctx context.Context, account types.Address, key types.Hash, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account types.Address, blockNumber *big.Int) ([]byte, error)
	NonceAt(ctx context.Context, account types.Address, blockNumber *big.Int) (uint64, error)
}

// SyncProgress gives progress indications when the node is synchronising with
// the Bchain network.
type SyncProgress struct {
	StartingBlock uint64 // Block number where sync began
	CurrentBlock  uint64 // Current block number where sync is at
	HighestBlock  uint64 // Highest alleged block number in the chain
	PulledStates  uint64 // Number of state trie entries already downloaded
	KnownStates   uint64 // Total number of state trie entries known about
}
