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
// @File: common.go
// @Date: 2018/06/14 14:14:14
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"bytes"
	"errors"
	"github.com/tinylib/msgp/msgp"
	"math/big"
	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
	"bchain.io/params"
)

// priority queue Item
type pqItem struct {
	value    interface{}
	priority *big.Int
}

//priority Queue
type priorityQueue []*pqItem

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority.Cmp(pq[j].priority) > 0
}
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*pqItem)
	*pq = append(*pq, item)
}

//pop the highest priority item
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

//todo future need use VRF function
//restore the seed by round
func restoreSeed(round uint64) (types.Hash, []byte, error) {
	sigByte := gCommonTools.GetQrSignature(round)
	sd := SeedData{}
	sd.Signature.init()
	err := sd.Signature.get(sigByte)
	if err != nil {
		return types.Hash{}, nil, err
	}
	sd.Round = round
	return sd.Hash(), sigByte, nil
}

func generateSeedBySeedHeader(seedHeader *block.Header) (types.Hash, []byte, error) {
	sigByte := seedHeader.Cdata.Para
	sd := SeedData{}
	sd.Signature.init()
	err := sd.Signature.get(sigByte)
	if err != nil {
		return types.Hash{}, nil, err
	}
	sd.Round = seedHeader.Number.IntVal.Uint64()
	return sd.Hash(), sigByte, nil
}

func makeEmptyBlockConsensusData(round uint64) *block.ConsensusData {
	bcd := &block.ConsensusData{}
	bcd.Id = ConsensusDataId

	sd := SeedData{}
	sd.init()
	sd.Round = round
	sd.sign(params.RewordPrikey)

	bcd.Para = sd.toBytes()
	return bcd
}

func makeBlockConsensusData(bp *BlockProposal, ct CommonTools) *block.ConsensusData {
	bcd := &block.ConsensusData{}
	bcd.Id = ConsensusDataId

	sd := &SeedData{}
	sd.init()
	sd.Round = bp.Credential.Round
	ct.SeedSig(sd)

	bcd.Para = sd.toBytes()
	return bcd
}

func senderFromBlock(header *block.Header, parent *block.Header) (types.Address, error) {
	sd := SeedData{}
	sd.init()
	err := sd.Signature.get(header.Cdata.Para)
	if err != nil {
		return types.Address{}, err
	}
	sd.Round = header.Number.IntVal.Uint64()
	return sd.sender(parent)

}

func CheckCertificate(header *block.Header, data []byte) error {
	var blockCertificate BlockCertificate
	byteBuf := bytes.NewBuffer(data)
	err := msgp.Decode(byteBuf, &blockCertificate)
	if err != nil {
		logger.Error("blockCertificate.Decode err", "err", err)
		return err
	}
	if len(blockCertificate) == 0 {
		return errors.New("BlockCertificate len is zero")
	}
	var step uint64
	hash := types.Hash{}
	for i, cs := range blockCertificate {
		if i == 0 {
			step = cs.Step
			if step >= uint64(Config().maxStep) || step == 0 {
				return ErrCertStep
			}
			hash = cs.ParentHash
		} else {
			if step != cs.Step {
				return ErrCertSteps
			}
			if hash != cs.ParentHash {
				return ErrCertPHash
			}
		}
	}

	if hash != header.ParentHash {
		return errors.New("retrieved block certificate is invalid, parent hash is not right")
	}

	logger.Debug("checkCertificate ok. step:", step, "certificates len:", len(blockCertificate))
	return nil
}

func getThreshold(step int) int64 {
	if step == STEP_FINAL {
		return Config().tFinalThreshold
	} else if step == STEP_BP {
		return Config().tProposer
	} else {
		return Config().tStepThreshold
	}
}

func GetCredentialVotes(cs *CredentialSign) float64 {
	return getCredentialVotes(cs, nil)
}

func getCredentialVotes(cs *CredentialSign, seedBlock *block.Header) float64 {
	sigHash := cs.Signature.Hash()
	var tao int64
	if cs.Step == STEP_BP {
		tao = Config().tProposer
	} else if cs.Step == StepFinal {
		tao = Config().tFinal
	} else {
		tao = Config().tStep
	}
	sender, err := cs.senderBySeedBlock(seedBlock)
	if err != nil {
		return 0
	}
	addr := sender
	w, W := gCommonTools.GetWeight(cs.Round, addr)

	cs.votes = Config().sp.getSortitionPriorityByHash(sigHash, w, tao, W)
	return cs.votes
}
