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
// @File: empty_engine.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package consensus

import (
	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
)

type Engine_empty struct {
}

func (empty *Engine_empty) Author(chain ChainReader, header *block.Header) (types.Address, error) {
	return header.Producer, nil
}

func (empty *Engine_empty) VerifyHeader(chain ChainReader, header *block.Header, seal bool) error {
	return nil
}

func (empty *Engine_empty) VerifyHeaders(chain ChainReader, headers []*block.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort, results := make(chan struct{}), make(chan error, len(headers))
	for i := 0; i < len(headers); i++ {
		results <- nil
	}
	return abort, results
}

func (empty *Engine_empty) VerifySeal(chain ChainReader, header *block.Header) error {
	return nil
}

func (empty *Engine_empty) Prepare(chain ChainReader, header *block.Header) error {
	return nil
}

func (empty *Engine_empty) Finalize(chain ChainReader, header *block.Header, state *state.StateDB, txs []*transaction.Transaction, receipts []*transaction.Receipt, sign bool) (*block.Block, error) {
	//reward := big.NewInt(5e+18)
	//state.AddBalance(header.BlockProducer, reward)
	header.StateRootHash = state.IntermediateRoot()
	return block.NewBlock(header, txs, receipts), nil
}

func (empty *Engine_empty) Seal(chain ChainReader, block *block.Block, stop <-chan struct{}) (*block.Block, error) {
	header := block.Header()
	return block.WithSeal(header), nil
}

func (empty *Engine_empty) Incentive(producer types.Address, state *state.StateDB, header *block.Header) (*transaction.Transaction, error) {
	return nil, nil
}