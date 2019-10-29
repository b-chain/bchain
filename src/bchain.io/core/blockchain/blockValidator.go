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
// @File: blockValidator.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package blockchain

import (
	"fmt"
	"bchain.io/core/state"
	"errors"
	"bchain.io/core/transaction"
	"bchain.io/consensus"
	"bchain.io/core/blockchain/block"
	"bchain.io/core"
	"bchain.io/params"
)

// BlockValidator implements Validator.
type BlockValidator struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for validating
}

// NewBlockValidator returns a new block validator which is safe for re-use
func NewBlockValidator(config *params.ChainConfig, blockchain *BlockChain, engine consensus.Engine) *BlockValidator {
	validator := &BlockValidator{
		config: config,
		engine: engine,
		bc:     blockchain,
	}
	return validator
}

// ValidateBody verifies the the block
// header's transaction root. The headers are assumed to be already
// validated at this point.
func (v *BlockValidator) ValidateBody(blk *block.Block) error {
	// Check whether the block's known, and if not, that it's linkable
	if v.bc.HasBlockAndState(blk.Hash()) {
		return core.ErrKnownBlock
	}
	if !v.bc.HasBlockAndState(blk.ParentHash()) {
		if !v.bc.HasBlock(blk.ParentHash(), blk.NumberU64()-1) {
			return errors.New("unknown ancestor")
		}
		return errors.New("pruned ancestor")
	}
	// Header validity is known at this point
	header := blk.Header()

	if hash := block.DeriveSha(blk.Transactions()); hash != header.TxRootHash {
		return fmt.Errorf("transaction root hash mismatch: have %x, want %x", hash.String(), header.TxRootHash.String())
	}
	return nil
}

// ValidateState validates the various changes that happen after a state
// transition, such as the receipt roots and the state root
// itself. ValidateState returns a database batch if the validation was a success
// otherwise nil and an error is returned.
func (v *BlockValidator) ValidateState(blk, parent *block.Block, statedb *state.StateDB, receipts transaction.Receipts) error {
	header := blk.Header()

	// Validate the received block's bloom with the one derived from the generated receipts.
	// For valid blocks this should always validate to true.
	rbloom := transaction.CreateBloom(receipts)
	if rbloom != header.Bloom {
		return fmt.Errorf("invalid bloom (remote: %x  local: %x)", header.Bloom, rbloom)
	}
	// Tre receipt Trie's root (R = (Tr [[H1, R1], ... [Hn, R1]]))
	receiptSha := block.DeriveSha(receipts)
	if receiptSha != header.ReceiptRootHash {
		return fmt.Errorf("invalid receipt root hash (remote: %v local: %v)", header.ReceiptRootHash.String(), receiptSha.String())
	}
	// Validate the state root against the received state root and throw
	// an error if they don't match.
	if root := statedb.IntermediateRoot(); header.StateRootHash != root {
		return fmt.Errorf("invalid merkle root (remote: %v local: %v)", header.StateRootHash.String(), root.String())
	}
	return nil
}