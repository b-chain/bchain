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
// @File: apos_engine.go
// @Date: 2018/08/07 13:59:07
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"bchain.io/common"
	"bchain.io/common/types"
	"bchain.io/consensus"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
	"bchain.io/params"
	"runtime"
	"encoding/json"
	"bchain.io/core/interpreter/wasmre/para_paser"
)

type EngineApos struct {
	//key for sign header
	prv    *ecdsa.PrivateKey
	signer transaction.Signer
}

var (
	ErrBlockTime = errors.New("timestamp less than or equal parent's timestamp")
	ErrSignature = errors.New("signature is not right")
	ErrAposID    = errors.New("consensus Id is not apos")
	ErrAposData  = errors.New("apos consensus fail")
	ErrAposSign  = errors.New("apos consensus signature is not equal to header")
	ErrCertSteps = errors.New("apos consensus certificate step is not same")
	ErrCertPHash = errors.New("apos consensus certificate parent hash is not same")
	ErrCertStep  = errors.New("apos consensus certificate step is not bba step")
	ErrCertVotes = errors.New("apos consensus certificate votes is not enough")
)

func NewAposEngine(prv *ecdsa.PrivateKey) *EngineApos {
	return &EngineApos{
		prv:    prv,
		signer: transaction.NewMSigner(Config().chainId),
	}
}

func (this *EngineApos) SetKey(prv *ecdsa.PrivateKey) {
	this.prv = prv
}

func (this *EngineApos) Author(chain consensus.ChainReader, header *block.Header) (types.Address, error) {
	singner := block.NewBlockSigner(chain.Config().ChainId)
	return singner.Sender(header)
}

func (this *EngineApos) VerifyHeader(chain consensus.ChainReader, header *block.Header, seal bool) error {
	//if the header is known, verify success
	number := header.Number.IntVal.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	// Verify that the block number is parent's +1
	if diff := new(big.Int).Sub(&header.Number.IntVal, &parent.Number.IntVal); diff.Cmp(common.Big1) != 0 {
		return consensus.ErrInvalidNumber
	}

	//verify time
	cmpResult := header.Time.IntVal.Cmp(&parent.Time.IntVal)
	if cmpResult < 0 {
		return ErrBlockTime
	}

	//verify ConsensusData
	if header.Cdata.Id != ConsensusDataId {
		return ErrAposID
	}
	//apos
	senderApos, err := senderFromBlock(header, parent)
	if err != nil {
		return ErrAposData
	}

	if seal {
		if err := this.VerifySeal(chain, header); err != nil {
			return err
		}
	}

	//verify signature
	singner := block.NewBlockSigner(chain.Config().ChainId)
	sender, err := singner.Sender(header)
	if err != nil {
		return ErrSignature
	}

	if sender != senderApos {
		return ErrAposSign
	}

	if cmpResult == 0 {
		//for apos, sender == params.Address means that empty block
		if sender != params.Address {
			return ErrBlockTime
		}
	}
	header.Producer = sender
	return nil
}

func (this *EngineApos) verifyHeader(chain consensus.ChainReader, header, parent *block.Header, seal bool, seedBlock *block.Header) error {
	//if the header is known, verify success
	number := header.Number.IntVal.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}

	// Verify that the block number is parent's +1
	if diff := new(big.Int).Sub(&header.Number.IntVal, &parent.Number.IntVal); diff.Cmp(common.Big1) != 0 {
		return consensus.ErrInvalidNumber
	}

	//verify time
	cmpResult := header.Time.IntVal.Cmp(&parent.Time.IntVal)
	if cmpResult < 0 {
		return ErrBlockTime
	}

	//verify ConsensusData
	if header.Cdata.Id != ConsensusDataId {
		return ErrAposID
	}
	//apos
	senderApos, err := senderFromBlock(header, parent)
	if err != nil {
		return ErrAposData
	}
	if seal {
		if err := this.verifySeal(chain, header, seedBlock); err != nil {
			return err
		}
	}

	//verify signature
	singner := block.NewBlockSigner(chain.Config().ChainId)
	sender, err := singner.Sender(header)
	if err != nil {
		return ErrSignature
	}

	if sender != senderApos {
		return ErrAposSign
	}

	if cmpResult == 0 {
		//for apos, sender == params.Address means that empty block
		if sender != params.Address {
			return ErrBlockTime
		}
	}
	header.Producer = sender
	return nil
}

func (this *EngineApos) verifyHeaderWorker(chain consensus.ChainReader, headers []*block.Header, seals []bool, index int) error {
	var parent *block.Header
	var seedBlock *block.Header
	if index == 0 {
		parent = chain.GetHeader(headers[0].ParentHash, headers[0].Number.IntVal.Uint64()-1)
	} else if headers[index-1].Hash() == headers[index].ParentHash {
		parent = headers[index-1]
	}
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	if chain.GetHeader(headers[index].Hash(), headers[index].Number.IntVal.Uint64()) != nil {
		return nil // known block
	}
	numStart := headers[0].Number.IntVal.Uint64()
	number := headers[index].Number.IntVal.Uint64()
	R := uint64(Config().R)
	seedRound := R
	if number < R {
		seedRound = 0
	} else {
		seedRound = number - 1 - (number % R)
	}
	if seedRound >= numStart {
		seedBlock = headers[seedRound-numStart]
	} else {
		seedBlock = nil
	}
	return this.verifyHeader(chain, headers[index], parent, seals[index], seedBlock)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications.
func (this *EngineApos) VerifyHeaders(chain consensus.ChainReader, headers []*block.Header, seals []bool) (chan<- struct{}, <-chan error) {
	workers := runtime.GOMAXPROCS(0)
	if len(headers) < workers {
		workers = len(headers)
	}
	// Create a task channel and spawn the verifiers
	var (
		inputs = make(chan int)
		done   = make(chan int, workers)
		errors = make([]error, len(headers))
		abort  = make(chan struct{})
	)
	for i := 0; i < workers; i++ {
		go func() {
			for index := range inputs {
				errors[index] = this.verifyHeaderWorker(chain, headers, seals, index)
				done <- index
			}
		}()
	}

	errorsOut := make(chan error, len(headers))
	go func() {
		defer close(inputs)
		var (
			in, out = 0, 0
			checked = make([]bool, len(headers))
			inputs  = inputs
		)
		for {
			select {
			case inputs <- in:
				if in++; in == len(headers) {
					// Reached end of headers. Stop sending to workers.
					inputs = nil
				}
			case index := <-done:
				for checked[index] = true; checked[out]; out++ {
					errorsOut <- errors[out]
					if out == len(headers)-1 {
						return
					}
				}
			case <-abort:
				return
			}
		}
	}()
	return abort, errorsOut
}

// verify block certificate
func (this *EngineApos) VerifySeal(chain consensus.ChainReader, header *block.Header) error {
	//for apos, header validate the certificate
	blockCertificate := gCommonTools.GetBlockCertificate(header.Hash())
	var sumVotes float64
	var step uint64
	for i, cs := range blockCertificate {
		if i == 0 {
			step = cs.Step
			if step >= uint64(Config().maxStep) || step == 0 {
				return ErrCertStep
			}
		} else {
			if step != cs.Step {
				return ErrCertSteps
			}
		}
		votes := getCredentialVotes(cs, nil)
		sumVotes += votes
	}
	if sumVotes > float64(getThreshold(int(step))) {
		logger.Debug("VerifySeal ok. s votes:", step, sumVotes)
		return nil
	} else {
		return ErrCertVotes
	}
}

// verify block certificate
func (this *EngineApos) verifySeal(chain consensus.ChainReader, header *block.Header, seedBlock *block.Header) error {
	//for apos, header validate the certificate
	blockCertificate := gCommonTools.GetBlockCertificate(header.Hash())
	var sumVotes float64
	var step uint64
	for i, cs := range blockCertificate {
		if i == 0 {
			step = cs.Step
			if step >= uint64(Config().maxStep) || step == 0 {
				return ErrCertStep
			}
		} else {
			if step != cs.Step {
				return ErrCertSteps
			}
		}
		votes := getCredentialVotes(cs, seedBlock)
		sumVotes += votes
	}
	if sumVotes > float64(getThreshold(int(step))) {
		logger.Debug("VerifySeal ok. s votes:", step, sumVotes)
		return nil
	} else {
		return ErrCertVotes
	}
}

func (this *EngineApos) Prepare(chain consensus.ChainReader, header *block.Header) error {
	return nil
}

//todo this need interpreter process
//interpreter need change state
func (this *EngineApos) Finalize(chain consensus.ChainReader, header *block.Header, state *state.StateDB, txs []*transaction.Transaction, receipts []*transaction.Receipt, sign bool) (*block.Block, error) {
	//reward := big.NewInt(5e+18)
	//state.AddBalance(header.BlockProducer, reward)
	header.StateRootHash = state.IntermediateRoot()

	//sign header
	if sign {
		if this.prv == nil {
			return nil, errors.New("No key found fo sign header")
		}
		blk := block.NewBlock(header, txs, receipts)

		err := block.SignHeaderInner(blk.H, block.NewBlockSigner(chain.Config().ChainId), this.prv)
		if err != nil {
			return nil, err
		}
		return blk, nil
	} else {
		return block.NewBlock(header, txs, receipts), nil
	}

}

func (this *EngineApos) Seal(chain consensus.ChainReader, block *block.Block, stop <-chan struct{}) (*block.Block, error) {
	header := block.Header()
	return block.WithSeal(header), nil
}

func (this *EngineApos) Incentive(producer types.Address, state *state.StateDB, header *block.Header) (*transaction.Transaction, error) {
	if header.Producer == params.Address && block.EmptyRootHash == header.TxRootHash {
		// apos empty block
		logger.Info("Incentive, empty block, no incentive transaction")
		return nil, nil
	}
	return this.makeRewardTransaction(producer, state)
}

func (this *EngineApos) makeRewardTransaction(producer types.Address, state *state.StateDB) (*transaction.Transaction, error) {
	actions := transaction.Actions{}
	wp := &para_paser.WasmPara {
		FuncName: "reword",
		Args:     []para_paser.Arg{},
	}
	para, _ := json.Marshal(wp)

	action := transaction.Action{types.HexToAddress("0xb78f12Cb3924607A8BC6a66799e159E3459097e9"), para}
	actions = append(actions, &action)

	//make tx
	//nc := this.backend.TxPool().State().GetNonce(params.Address)
	nc := state.GetNonce(params.Address)
	tx := transaction.NewTransaction(nc, actions)

	txSign, err := transaction.SignTx(tx, this.signer, params.RewordPrikey)
	if err != nil {
		logger.Errorf("sign Tx with private key err :", err.Error())
		return nil, err
	}
	return txSign, nil
}
