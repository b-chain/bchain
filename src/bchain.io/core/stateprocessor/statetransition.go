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
// @File: statetransition.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package stateprocessor

import (
	"errors"
	"bchain.io/common/types"
	"bchain.io/core"
	"bchain.io/core/actioncontext"
	"bchain.io/core/interpreter"
	"bchain.io/core/transaction"
)

/*
The State Transitioning Model

A state transition is a change made when a transaction is applied to the current world state
The state transitioning model does all all the necessary work to work out a valid new state root.

1) Nonce handling
3) Create a new state object if the recipient is \0*32
4) Value transfer
== If contract creation ==
  4a) Attempt to run transaction data
  4b) If valid, use result as code for the new state object
== end ==
5) Run Script section
6) Derive new state root
*/
type StateTransition struct {
	nonce       uint64
	bCheckNonce bool
	sender      types.Address
	actions     transaction.Actions
	blkCtx      *actioncontext.BlockContext
}

// Message represents a message sent to a contract.
type Message interface {
	From() types.Address
	Actions() []transaction.Action
	Nonce() uint64
	CheckNonce() bool
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(tx *transaction.Transaction, sender types.Address, blkCtx *actioncontext.BlockContext) *StateTransition {
	newActions := transaction.Actions{}
	newActions = append(newActions, tx.Data.Acts...)
	return &StateTransition{
		nonce:       tx.Nonce(),
		bCheckNonce: tx.CheckNonce(),
		sender:      sender,
		actions:     newActions,
		blkCtx:      blkCtx,
	}
}

/*
// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
func ApplyMessage(statedb *state.StateDB, msg Message, coinBase types.Address, header *block.Header) ([]byte, bool, error) {
	return NewStateTransition(statedb, msg, coinBase, header).TransitionDb()
}
*/
func (st *StateTransition) from() types.Address {
	f := st.sender
	if !st.blkCtx.GetState().Exist(f) {
		st.blkCtx.GetState().CreateAccount(f)
	}
	return f
}

func (st *StateTransition) preCheck() error {
	sender := st.sender

	// Make sure this transaction's nonce is correct
	if st.bCheckNonce {
		nonce := st.blkCtx.GetState().GetNonce(sender)
		if nonce < st.nonce {
			return core.ErrNonceTooHigh
		} else if nonce > st.nonce {
			return core.ErrNonceTooLow
		}
	}
	return nil
}

// TransitionDb will transition the state by applying the current message and
// returning the result. It returns an error if it
// failed. An error indicates a consensus issue.
func (st *StateTransition) TransitionDb() (ret []byte, contracts []types.Address, failed bool, err error) {
	if err = st.preCheck(); err != nil {
		return
	}

	sender := st.from() // err checked in preCheck
	// Snapshot !!!!!!!!!!!!!!!!!
	snapshot := st.blkCtx.GetState().Snapshot()

	logger.Debugf("Just process actions transaction.")
	st.blkCtx.GetState().SetNonce(sender, st.blkCtx.GetState().GetNonce(sender)+1)

	contracts = make([]types.Address, 0)
	for _, act := range st.actions {
		ctx := actioncontext.NewContext(sender, act, st.blkCtx)
		if ctx == nil {
			logger.Warn("new context return nil by contract ", act.Contract.Hex())
			err = errors.New("new context fail")
			break
		}
		err = ctx.Exec(interpreter.Singleton())
		if err != nil {
			break
		}
		contracts = append(contracts, ctx.GetContracts()...)
	}

	if err != nil {
		//clean the slice
		contracts = append([]types.Address{})
		st.blkCtx.GetState().RevertToSnapshot(snapshot)
		return nil, contracts, true, err
	}
	return ret, contracts, false, err
}
