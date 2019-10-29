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
// @File: consensus.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package consensus

import (
	"errors"
	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
	"bchain.io/params"
)

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header
type ChainReader interface {
	// Config retrieves the blockchain's chain configuration.
	Config() *params.ChainConfig

	// CurrentHeader retrieves the current header from the local chain.
	CurrentHeader() *block.Header

	// GetHeader retrieves a block header from the database by hash and number.
	GetHeader(hash types.Hash, number uint64) *block.Header

	// GetHeaderByNumber retrieves a block header from the database by number.
	GetHeaderByNumber(number uint64) *block.Header

	// GetHeaderByHash retrieves a block header from the database by its hash.
	GetHeaderByHash(hash types.Hash) *block.Header

	// GetBlock retrieves a block from the database by hash and number.
	GetBlock(hash types.Hash, number uint64) *block.Block
}

// Engine is an algorithm agnostic consensus engine.
type Engine interface {
	// Author retrieves the Bchain address of the account that minted the given
	// block, which may be different from the header's coinbase if a consensus
	// engine is based on signatures.
	Author(chain ChainReader, header *block.Header) (types.Address, error)

	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// given engine. Verifying the seal may be done optionally here, or explicitly
	// via the VerifySeal method.
	VerifyHeader(chain ChainReader, header *block.Header, seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	VerifyHeaders(chain ChainReader, headers []*block.Header, seals []bool) (chan<- struct{}, <-chan error)

	// VerifySeal checks whether the crypto seal on a header is valid according to
	// the consensus rules of the given engine.
	VerifySeal(chain ChainReader, header *block.Header) error

	// Prepare initializes the consensus fields of a block header according to the
	// rules of a particular engine. The changes are executed inline.
	Prepare(chain ChainReader, header *block.Header) error

	// Finalize runs any post-transaction state modifications (e.g. block rewards)
	// and assembles the final block.
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	Finalize(chain ChainReader, header *block.Header, state *state.StateDB, txs []*transaction.Transaction, receipts []*transaction.Receipt, sign bool) (*block.Block, error)

	// Seal generates a new block for the given input block with the local blockproducer's
	// seal place on top.
	Seal(chain ChainReader, block *block.Block, stop <-chan struct{}) (*block.Block, error)

	// APIs returns the RPC APIs this consensus engine provides.
	//APIs(chain ChainReader) []rpc.API

	// Incentive generates consensus transaction which only affect block state
	Incentive(producer types.Address, state *state.StateDB, header *block.Header) (*transaction.Transaction, error)
}

var (
	// ErrUnknownAncestor is returned when validating a block requires an ancestor
	// that is unknown.
	ErrUnknownAncestor = errors.New("unknown ancestor")

	// ErrFutureBlock is returned when a block's timestamp is in the future according
	// to the current node.
	ErrFutureBlock = errors.New("block in the future")

	// ErrInvalidNumber is returned if a block's number doesn't equal it's parent's
	// plus one.
	ErrInvalidNumber = errors.New("invalid block number")
)
