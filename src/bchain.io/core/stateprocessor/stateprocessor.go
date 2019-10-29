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
// @File: stateprocessor.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package stateprocessor

import (
	"bchain.io/consensus"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
	"bchain.io/params"
	"bchain.io/utils/bloom"
	"bchain.io/utils/database"
	"bchain.io/core/actioncontext"
)

type IChainForState interface {
	consensus.ChainReader
}

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	cs     IChainForState      // chain interface for state processor
	engine consensus.Engine    // Consensus engine used for block rewards
	db     database.IDatabase
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(config *params.ChainConfig, cs IChainForState, engine consensus.Engine, db database.IDatabase) *StateProcessor {
	return &StateProcessor{
		cs:     cs,
		engine: engine,
		config: config,
		db:     db,
	}
}

// Process processes the state changes according to the Bchain rules by running
// the transaction messages using the statedb and applying any rewards to
// the processor (coinbase).
//
// Process returns the receipts and logs accumulated during the process.
// If any of the transactions failed  it will return an error.
func (p *StateProcessor) Process(blk *block.Block, statedb *state.StateDB, db database.IDatabaseGetter, config *params.ChainConfig) (transaction.Receipts, []*transaction.Log, error) {
	var (
		receipts transaction.Receipts
		header   = blk.Header()
		allLogs  []*transaction.Log
	)

	singner := block.NewBlockSigner(config.ChainId)
	coinbase, err := singner.Sender(header)
	if err != nil {
		logger.Error("Process: block signature is not right", err)
		return nil, nil, err
	}

	logger.Trace("Process: coinbase", coinbase.Hex())

	tmpDb, _ := database.OpenMemDB()
	blkCtx := actioncontext.NewBlockContext(statedb, p.db, tmpDb, &header.Number.IntVal,coinbase)

	txs := []*transaction.Transaction{}
	incentiveTx, err := p.engine.Incentive(coinbase, statedb, header)
	if err != nil {
		logger.Error("Failed to fetch incentive transaction", "err", err)
		return nil, nil, err
	}
	if incentiveTx != nil {
		txs = append(txs, incentiveTx)
	}
	txs = append(txs, blk.Transactions()...)
	// Iterate over and process the individual transactions
	for i, tx := range txs {
		statedb.Prepare(tx.Hash(), blk.Hash(), i)
		receipt, err := ApplyTransaction(p.config, header, tx, blkCtx)
		if err != nil {
			logger.Errorf("ApplyTransacton Wrong.....:", err.Error())

			return nil, nil, err
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}

	// TODO: need to be compeleted, now skip this step
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	if p.engine != nil {
		p.engine.Finalize(p.cs, header, statedb, blk.Transactions(), receipts, false)
	}

	return receipts, allLogs, nil
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, header *block.Header, tx *transaction.Transaction, blkCtx *actioncontext.BlockContext) (*transaction.Receipt, error) {
	signer := transaction.MakeSigner(config, &header.Number.IntVal)
	sender, err := transaction.Sender(signer, tx)
	if err != nil {
		return nil, err
	}
	// Apply the transaction to the current state (included in the env)
	st := NewStateTransition(tx, sender, blkCtx)
	_, contracts, failed, err := st.TransitionDb()
	if err != nil {
		return nil, err
	}
	// Update the state with pending changes
	blkCtx.GetState().Finalise(true)

	// Create a new receipt for the transaction, storing the intermediate root  by the tx
	// based on the mip phase, we're passing wether the root touch-delete accounts.
	receipt := transaction.NewReceipt(failed)
	receipt.TxHash = tx.Hash()
	/*
		// if the transaction created a contract, store the creation address in the receipt.
		if len(tx.Data.Acts) == 2 && len(tx.Data.Acts[1].Contract) == 0 {
			receipt.ContractAddress = crypto.CreateAddress(msg(), tx.Nonce())
		}
	*/
	// if the transaction created contracts, store the creation addresses in the receipt.
	receipt.ContractAddress = append(receipt.ContractAddress, contracts...)
	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = blkCtx.GetState().GetLogs(tx.Hash())

	topics := make([]bloom.BloomByte, 0)
	for _, log := range receipt.Logs {
		topics = append(topics, log.Address)
		for _, topic := range log.Topics {
			topics = append(topics, topic)
		}
	}
	receipt.Bloom = bloom.CreateBloom(topics)

	return receipt, err
}
