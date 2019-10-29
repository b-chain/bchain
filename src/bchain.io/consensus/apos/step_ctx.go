package apos

import (
	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
)

//stepCtx contains all functions the stepObj will use
type stepCtx struct {
	getStep   func() int // get the number of step in the round
	getRound  func() uint64
	stopStep  func() // stop the step
	stopRound func() // stop all the step in the round, and end the round

	//getCredential func() signature
	//getEphemeralSig func(signed []byte) signature
	esig                  func(pEphemeralSign *EphemeralSign) error
	sendInner             func(pack dataPack) error
	propagateMsg          func(dataPack) error
	getCredential         func() *CredentialSign
	setRound              func(*Round)
	makeEmptyBlockForTest func(cs *CredentialSign) *block.Block
	getEmptyBlockHash     func() types.Hash
	getEphemeralSig       func(signed []byte) Signature
	getProducerNewBlock   func(data *block.ConsensusData, timeLimit int64) *block.Block
	//getPrivKey

	//gilad
	writeRet                 func(data *VoteData) //x
	verifyBlock              func(b *block.Block) bool
	getCredentialByStep      func(step uint64) *CredentialSign
	startVoteTimer           func(delay int)
	makeBlockConsensusData   func(bp *BlockProposal) *block.ConsensusData
	setBpResult              func(hash types.Hash)
	setReductionResult       func(hash types.Hash)
	setBbaResult             func(hash types.Hash)
	setFinalResult           func(hash types.Hash)
	getCommonCoinMinHashRslt func(step int) int
	getWeight                func(r uint64, addr types.Address) (int64, int64)
}
