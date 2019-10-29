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
// @File: apos.go
// @Date: 2018/06/15 11:35:15
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
	"bchain.io/node/services/bchain/downloader"
	"bchain.io/utils/event"
	"reflect"
	"sync"
	"sync/atomic"
	"math/big"
	"time"
)

var (
	ConsensusDataId = "apos"
)

const (
	ROUND_OK = iota
	ROUND_HANGFOREVER
	ROUND_NO_BLOCK
	ROUND_STOP
	ROUND_INSERT_ERR
	ROUND_WRITE_CERTIFICATE_ERR
)

type peerMsgs struct {
	msgCs  map[int]*CredentialSign
	msgBas map[int]*ByzantineAgreementStar

	//0 :default honesty peer. 1: malicious peer
	honesty uint //judge by baMsg filter , if one node send different hash at same step
}

type mainStepOutput struct {
	bp        types.Hash
	reduction types.Hash
	bba       types.Hash
	final     types.Hash
	mu        sync.Mutex
}

func (this *mainStepOutput) setBpResult(bp types.Hash) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.bp = bp
}
func (this *mainStepOutput) setReductionResult(reduction types.Hash) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.reduction = reduction
}
func (this *mainStepOutput) setBbaResult(bba types.Hash) bool {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.bba = bba
	nullHash := types.Hash{}

	if this.final != nullHash {
		return true
	}
	return false
}
func (this *mainStepOutput) setFinalResult(final types.Hash) bool {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.final = final
	nullHash := types.Hash{}
	if this.bba != nullHash {
		return true
	}
	return false
}

//round context
type Round struct {
	round       uint64
	apos        *Apos
	credentials map[int]*CredentialSign
	lock        sync.RWMutex
	emptyBlock  *block.Block
	msgs        map[types.Address]*peerMsgs
	resultCh    chan *block.Block //the result of Round,a block was made,should insert into blockchain
	roundOverCh chan int
	roundHangCh chan interface{}
	bpObj       *BpObj
	voteObj     *VoteObj
	//version 1.1
	mainStepRlt mainStepOutput
	parentHash  types.Hash
	countVote   *countVote
	stopCh      chan interface{}
}

//gilad tools
func (this *Round) startVoteTimer(delay int) {
	this.countVote.startTimer(delay)
}

func (this *Round) makeBlockConsensusData(bp *BlockProposal) *block.ConsensusData {
	return makeBlockConsensusData(bp, this.apos.commonTools)
}

func (this *Round) getCredentialByStep(step uint64) *CredentialSign {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if c, ok := this.credentials[int(step)]; ok {
		return c
	}
	return nil
}

func (this *Round) verifyBlock(b *block.Block) bool {
	lastHash := this.apos.commonTools.GetNowBlockHash()

	//here we just compare the parent hash is right or not
	if lastHash.Equal(&b.H.ParentHash) {
		return true
	}
	return false
}

func newRound(round int, parentHash types.Hash, apos *Apos, roundOverCh chan int) *Round {
	r := new(Round)
	r.init(round, apos, roundOverCh)
	r.parentHash = parentHash
	return r
}

func (this *Round) getEmptyBlockHash() types.Hash {
	return this.emptyBlock.Hash()
}

func (this *Round) init(round int, apos *Apos, roundOverCh chan int) {
	this.round = uint64(round)
	this.apos = apos
	this.roundOverCh = roundOverCh

	// this.maxLeaderNum = this.apos.algoParam.maxLeaderNum
	this.credentials = make(map[int]*CredentialSign)
	emptyBlock := this.apos.commonTools.MakeEmptyBlock(makeEmptyBlockConsensusData(this.round))
	this.emptyBlock = emptyBlock

	this.roundHangCh = make(chan interface{}, 1)
	this.resultCh = make(chan *block.Block, 1)
	this.stopCh = make(chan interface{}, 1)

	this.msgs = make(map[types.Address]*peerMsgs)
}

func (this *Round) hangForever() {
	this.roundHangCh <- 1
}

func (this *Round) stop() {
	this.stopCh <- 1
}

func (this *Round) setBpResult(hash types.Hash) {

	logger.Info("round", this.round, "setBpResult", hash.String())

	this.mainStepRlt.setBpResult(hash)
}
func (this *Round) setReductionResult(hash types.Hash) {
	logger.Info("round", this.round, "setReductionResult", hash.String())
	this.mainStepRlt.setReductionResult(hash)
}

func (this *Round) setBbaResult(hash types.Hash) {
	logger.Info("round", this.round, "setBbaResult", hash.String())
	complete := this.mainStepRlt.setBbaResult(hash)
	if complete {
		if hash == this.mainStepRlt.final {
			logger.Info("Final consensus!!!")
		} else {
			logger.Info("Tentative consensus!!!")
		}
		consensusBlock := this.bpObj.getExistBlock(hash)
		//not get the exist block

		if hash == this.getEmptyBlockHash() {
			consensusBlock = this.emptyBlock
		}
		if consensusBlock != nil {
			this.resultCh <- consensusBlock
		} else {
			logger.Error(COLOR_FRONT_RED, "SetFinalResult Get a nil block", this.mainStepRlt.bba.Hex(), COLOR_SHORT_RESET)
			//need download this block based on hash, here round just exit, recovery round will download
			blk := this.getBufferBlock(hash)
			this.resultCh <- blk
		}
	}
}

func (this *Round) setFinalResult(hash types.Hash) {
	logger.Info("round", this.round, "setFinalResult", hash.String())
	complete := this.mainStepRlt.setFinalResult(hash)
	if complete {
		if hash == this.mainStepRlt.bba {
			logger.Info("Final consensus!!!")
		} else {
			logger.Info("Tentative consensus!!!")
		}
		consensusBlock := this.bpObj.getExistBlock(this.mainStepRlt.bba)

		if this.mainStepRlt.bba == this.getEmptyBlockHash() {
			consensusBlock = this.emptyBlock
		}
		if consensusBlock != nil {
			this.resultCh <- consensusBlock
		} else {
			logger.Error(COLOR_FRONT_RED, "SetFinalResult Get a nil block", this.mainStepRlt.bba.Hex(), COLOR_SHORT_RESET)
			//need download this block based on hash, here round just exit, recovery round will download
			blk := this.getBufferBlock(hash)
			this.resultCh <- blk
		}
	} else {
		logger.Debug(COLOR_FRONT_RED, "setFinalResult,but not complete......", COLOR_SHORT_RESET)
	}
}

//inform stepObj to stop running
func (this *Round) broadCastStop() {
	logger.Debug(COLOR_FRONT_RED, "In BroadCastStop...", COLOR_SHORT_RESET)
	this.bpObj.stop()
	this.voteObj.stop()
	this.countVote.stop()
}

// Generate valid Credentials in current round
func (this *Round) generateCredentials() {
	sp := Config().sp
	csTime := types.NewBigInt(*big.NewInt(time.Now().Unix()))
	addr := this.apos.commonTools.GetCoinBase()
	w, W := this.apos.commonTools.GetWeight(this.round, addr)

	logger.Info(COLOR_FRONT_YELLOW, "CoinBase:", addr.Hex(), "w:", w, " W:", W, COLOR_SHORT_RESET)
	for i := 1; i < int(Config().maxStep); i++ {
		credential := this.apos.makeCredential(i, sp, csTime, w, W)
		isVerifier := this.apos.judgeVerifier(credential, i)
		//logger.Info("GenerateCredential step:",i,"  isVerifier:",isVerfier)
		if isVerifier {
			logger.Info("GenerateCredential step:", i, "  votes:", credential.votes)
			this.credentials[i] = credential
			this.apos.commonTools.CreateTmpPriKey(i)
		}
	}

	for i := STEP_BP; i < STEP_IDLE; i++ {
		credential := this.apos.makeCredential(i, sp, csTime, w, W)
		isVerifier := this.apos.judgeVerifier(credential, i)
		//logger.Info("GenerateCredential step:",i,"  isVerifier:",isVerfier)
		if isVerifier {
			logger.Info("GenerateCredential step:", i, " votes:", credential.votes)
			this.credentials[i] = credential
			this.apos.commonTools.CreateTmpPriKey(i)
		}
	}
}

func (this *Round) broadcastCredentials() {
	for i, credential := range this.credentials {
		_ = i
		logger.Info("SendCredential round", this.round, "step", i, "votes", credential.votes)
		this.apos.outMsger.SendInner(credential)
	}
}

func (this *Round) startStepObjs() {
	stepCtx := &stepCtx{}

	stepCtx.setBpResult = this.setBpResult
	stepCtx.setReductionResult = this.setReductionResult
	stepCtx.setBbaResult = this.setBbaResult
	stepCtx.setFinalResult = this.setFinalResult

	//ctx for new step obj
	//stepCtx.verifyBlock = this.verifyBlock
	stepCtx.verifyBlock = this.apos.commonTools.VerifyNextRoundBlock
	stepCtx.getCredentialByStep = this.getCredentialByStep
	stepCtx.startVoteTimer = this.startVoteTimer
	stepCtx.getProducerNewBlock = this.apos.commonTools.GetProducerNewBlock
	stepCtx.makeBlockConsensusData = this.makeBlockConsensusData

	roundRt := this.round
	stepCtx.getRound = func() uint64 {
		return roundRt
	}

	stepCtx.esig = this.apos.commonTools.Esig
	stepCtx.sendInner = this.apos.outMsger.SendInner
	stepCtx.propagateMsg = this.apos.outMsger.PropagateMsg
	stepCtx.getEmptyBlockHash = this.getEmptyBlockHash
	stepCtx.getWeight = this.apos.commonTools.GetWeight

	this.bpObj = makeBpObj(stepCtx)
	this.voteObj = makeVoteObj(stepCtx)

	sendVoteData := func(step int, hash types.Hash) {
		this.voteObj.SendVoteData(this.round, uint64(step), hash)
	}
	this.countVote = newCountVote(sendVoteData, this.hangForever, this.emptyBlock.Hash())

	//here set getCommonCoinMinHashRslt,because voteObj has not run yet,no error
	stepCtx.getCommonCoinMinHashRslt = this.countVote.getCommonCoinMinHashRslt

	go this.bpObj.run()
	go this.voteObj.run()
	go this.countVote.run()

}

//duplicate message check
func (this *Round) filterMsgCs(msg *CredentialSign) error {
	address, err := msg.sender()
	if err != nil {
		return err
	}
	step := msg.Step

	if peermsgs, ok := this.msgs[address]; ok {
		if peermsgs.honesty == 1 {
			return errors.New("not honesty peer")
		}
		if mCs, ok := peermsgs.msgCs[int(step)]; ok {
			if mCs.Step == step {
				return errors.New("duplicate message Credential")
			}
		} else {
			peermsgs.msgCs[int(step)] = msg
		}
	} else {
		ps := &peerMsgs{
			msgBas:  make(map[int]*ByzantineAgreementStar),
			msgCs:   make(map[int]*CredentialSign),
			honesty: 0,
		}
		ps.msgCs[int(step)] = msg
		this.msgs[address] = ps
	}
	return nil
}

// process the Credential message
func (this *Round) receiveMsgCs(msg *CredentialSign) {
	//logger.Debug("Receive message CredentialSign [r:s]:", msg.Round, msg.Step)
	if msg.Round != this.round {
		logger.Debug("verify fail, Credential msg is not in current round:", msg.Round, "want:", this.round)
		return
	}

	//duplicate message check
	if err := this.filterMsgCs(msg); err != nil {
		logger.Debug("filter Credential fail", err)
		return
	}
	//Propagate message via p2p
	this.apos.outMsger.PropagateMsg(msg)
}

func (this *Round) receiveMsgBp(msg *BlockProposal) {
	//verify msg
	logger.Debug("Receive message BlockProposal [r:s:hash:votes]:", msg.Credential.Round, msg.Credential.Step, msg.Block.Hash().String(), msg.Credential.votes)
	if msg.Credential.Round != this.round {
		logger.Warn("verify fail, BlockProposal msg is not in current round", msg.Credential.Round, this.round)
		return
	}
	this.bpObj.sendMsg(msg)
	// for BP Propagate process will in stepObj
}

//filter duplicate msg and same step with different hash
func (this *Round) filterMsgBa(msg *ByzantineAgreementStar) error {
	address, err := msg.Credential.sender()
	if err != nil {
		return err
	}
	step := msg.Credential.Step

	if peerMsgBas, ok := this.msgs[address]; ok {
		if peerMsgBas.honesty == 1 {
			return errors.New("not honesty peer")
		}
		if peerba, ok := peerMsgBas.msgBas[int(step)]; ok {
			if peerba.Hash == msg.Hash {
				return errors.New("duplicate message ByzantineAgreementStar")
			} else {
				peerMsgBas.honesty = 1
				return errors.New("receive different hash in BA message, it must a malicious peer")
			}
		} else {
			peerMsgBas.msgBas[int(step)] = msg
		}
	} else {
		ps := &peerMsgs{
			msgBas:  make(map[int]*ByzantineAgreementStar),
			msgCs:   make(map[int]*CredentialSign),
			honesty: 0,
		}
		ps.msgBas[int(step)] = msg
		this.msgs[address] = ps
	}
	return nil
}
func (this *Round) receiveMsgBaStar(msg *ByzantineAgreementStar) {
	//verify msg
	if msg.Credential.Round != this.round {
		logger.Debug("verify fail, ba msg is not in current round", this.round, "message round", msg.Credential.Round)
		return
	}
	if msg.Credential.ParentHash != this.parentHash {
		logger.Warn("verify fail, ba msg is not in current block chain", msg.Credential.ParentHash.String(), this.parentHash.String())
		return
	}
	if err := this.filterMsgBa(msg); err != nil {
		logger.Debug("filter ba message fail:", err)
		return
	}

	this.countVote.sendMsg(msg)
	//Propagate message via p2p
	this.apos.outMsger.PropagateMsg(msg)
}

// get channel's buffer block as possible
func (this *Round) getBufferBlock(hash types.Hash) *block.Block{
	timer := time.NewTimer(300 * time.Millisecond)
	defer timer.Stop()
	for timeout := false; !timeout; {
		select {
		case outData := <-this.apos.outMsger.GetDataMsg():
			switch v := outData.(type) {
			case *BlockProposal:
				if hash == v.Block.Hash() {
					if ok := this.apos.commonTools.VerifyNextRoundBlock(v.Block); ok {
						logger.Info("getBufferBp OK ")
						return v.Block
					}
				}
			default:
				logger.Info("getBufferBp ignore message type ", reflect.TypeOf(v))
			}
		case <-timer.C:
			timeout = true
			break
		}
	}
	if blk := this.bpObj.getExistBlock(hash); blk != nil {
		logger.Info("getExistBlock OK ")
		return blk
	}

	logger.Info("getBufferBp still return nil.")
	return nil
}


func (this *Round) commonProcess() int {
	defer this.broadCastStop()
	for {
		select {
		// receive message
		case outData := <-this.apos.outMsger.GetDataMsg():
			switch v := outData.(type) {
			case *CredentialSign:
				this.receiveMsgCs(v)
			case *BlockProposal:
				this.receiveMsgBp(v)
			case *ByzantineAgreementStar:
				this.receiveMsgBaStar(v)
			default:
				logger.Warn("invalid message type ", reflect.TypeOf(v))
			}
		//Round complete , insert the block into blockchain
		case consensusBlock := <-this.resultCh:
			//fmt.Println("CommonProcess end block:", consensusBlock)
			if consensusBlock == nil {
				logger.Error("consensusBlock is nil, need download this block")
				return ROUND_NO_BLOCK
			}
			certificate := this.countVote.getBlockCertificate()
			for i, cs := range certificate {
				fmt.Println("block certificate:", i, cs)
			}
			if err := this.apos.commonTools.WriteBlockCertificate(consensusBlock, certificate); err != nil {
				logger.Error("write block certificate fail", err)
				return ROUND_WRITE_CERTIFICATE_ERR
			}

			bs := block.Blocks{}
			bs = append(bs, consensusBlock)
			logger.Info("InsertChain start")
			_, err := this.apos.commonTools.InsertChain(bs)
			if err != nil {
				logger.Info("InsertChain error")
				return ROUND_INSERT_ERR
			}
			fmt.Println("InsertOneBlock    ErrStatus:", err)

			logger.Info("round exit ")
			return ROUND_OK
		case <-this.roundHangCh:
			logger.Warn("commonProcess: hang for ever")
			return ROUND_HANGFOREVER
		case <-this.stopCh:
			logger.Warn("commonProcess: stop")
			return ROUND_STOP
		}
	}
}

func (this *Round) run() {
	// make verifiers Credential
	this.generateCredentials()

	// broadcast Credentials
	this.broadcastCredentials()

	this.startStepObjs()

	ret := this.commonProcess()

	this.roundOverCh <- ret //inform the caller,the mission complete
}

type Apos struct {
	commonTools CommonTools
	outMsger    OutMsger

	roundCtx       *Round
	roundOverCh    chan int
	recoveryRound  *recoveryRound
	recoveryOverCh chan interface{}

	lock sync.RWMutex
	mux  *event.TypeMux

	canStart    int32 // can start indicates whether we can start the apos run operation
	shouldStart int32 // should start indicates whether we should start after sync
	running     int32
	recovering  int32
	runLock     sync.RWMutex
}

//Create Apos
func NewApos(bcHandler BlockChainHandler, producerHandler BlockProducerHandler, mux *event.TypeMux) *Apos {
	logger.Info("NewApos....................")
	a := new(Apos)
	cmTools := newAposTools(Config().chainId, bcHandler, producerHandler)
	a.commonTools = cmTools
	gCommonTools = cmTools
	a.roundOverCh = make(chan int, 1)
	a.recoveryOverCh = make(chan interface{}, 1)
	a.outMsger = MsgTransfer()

	a.mux = mux
	a.canStart = 1

	go a.monitorDownloader()
	return a
}

func (this *Apos) monitorDownloader() {
	events := this.mux.Subscribe(downloader.StartEvent{}, downloader.DoneEvent{}, downloader.FailedEvent{})
	for ev := range events.Chan() {
		switch ev.Data.(type) {
		case downloader.StartEvent:
			logger.Info(COLOR_FRONT_PINK, "downloader.Star	tEvent", COLOR_SHORT_RESET)

			atomic.StoreInt32(&this.canStart, 0)

			if atomic.LoadInt32(&this.running) > 0 {
				this.stopAllRound()
				atomic.StoreInt32(&this.shouldStart, 1)
			}
			//this.stopAllRound()
		case downloader.DoneEvent, downloader.FailedEvent:
			logger.Info(COLOR_FRONT_PINK, "downloader.DoneEvent", COLOR_SHORT_RESET)
			shouldStart := atomic.LoadInt32(&this.shouldStart) == 1

			atomic.StoreInt32(&this.canStart, 1)
			atomic.StoreInt32(&this.shouldStart, 0)
			if shouldStart {
				logger.Info(COLOR_FRONT_PINK, "downloader.DoneEvent shouldStart", COLOR_SHORT_RESET)
				go this.Start()
			}
		}
	}
}

func (this *Apos) SetPriKey(priKey *ecdsa.PrivateKey) {
	this.commonTools.SetPriKey(priKey)
}

func (this *Apos) SetCoinBase(coinBase types.Address) {
	this.commonTools.SetCoinBase(coinBase)
}

func (this *Apos) Start() {
	logger.Info(COLOR_FRONT_GREEN, "apos Start....", COLOR_SHORT_RESET)
	atomic.StoreInt32(&this.shouldStart, 1)
	if 0 == atomic.LoadInt32(&this.canStart) {
		logger.Info(COLOR_FRONT_GREEN, "Network syncing, will start apos afterwards", COLOR_SHORT_RESET)
		return
	}
	atomic.StoreInt32(&this.running, 1)
	this.outMsger.setAposRunning(true)
	go this.run()
}

//this is the main loop of Apos
func (this *Apos) run() {
	this.runLock.Lock()
	defer this.runLock.Unlock()
	for {
		if 0 != atomic.LoadInt32(&this.canStart) {
			this.runRecovery()
		}

		if 0 != atomic.LoadInt32(&this.canStart) {
			this.runNormal()
		}

		if 0 == atomic.LoadInt32(&this.canStart) {
			logger.Debug(COLOR_FRONT_PINK, "Apos Run exit", COLOR_SHORT_RESET)
			return
		}
	}
}

func (this *Apos) runNormal() int {
	logger.Info("apos run round:", this.commonTools.GetNextRound())
	this.runNormalRound()
	defer this.setNormalRoundNil() //for GC
	for {
		select {
		case roundRet := <-this.roundOverCh:
			logger.Info("apos new round running...............round return", roundRet)
			if roundRet != ROUND_OK {
				return roundRet
			}
			this.runNormalRound()
		}
	}
}

func (this *Apos) runRecovery() {
	logger.Info("apos recovery is running.....")
	atomic.StoreInt32(&this.recovering, 1)
	defer atomic.StoreInt32(&this.recovering, 0)
	r := this.runRecoveryRound()
	defer this.setRecoveryRoundNil() //for GC
	for {
		select {
		case <-this.recoveryOverCh:
			logger.Info("apos recovery finish!!!", r)
			return
		}
	}
}

func (this *Apos) runNormalRound() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.roundCtx = newRound(this.commonTools.GetNextRound(), this.commonTools.GetNowBlockHash(), this, this.roundOverCh)
	go this.roundCtx.run()
}

func (this *Apos) setNormalRoundNil() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.roundCtx = nil
}

func (this *Apos) runRecoveryRound() int{
	this.lock.Lock()
	defer this.lock.Unlock()
	r := this.commonTools.GetNextRound()
	parentHash := this.commonTools.GetNowBlockHash()
	logger.Info("apos run recovery round:", r)
	this.recoveryRound = newRecoveryRound(r, parentHash, this, this.recoveryOverCh)
	go this.recoveryRound.run()
	return r
}

func (this *Apos) setRecoveryRoundNil() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.recoveryRound = nil
}

func (this *Apos) stopAllRound() {
	this.lock.Lock()
	defer this.lock.Unlock()

	atomic.StoreInt32(&this.running, 0)

	if this.roundCtx != nil {
		logger.Info("stopAllRound, stop normal round", this.roundCtx.round)
		go this.roundCtx.stop()
	}
	if this.recoveryRound != nil {
		logger.Info("stopAllRound, stop recovery round", this.recoveryRound.round)
		go this.recoveryRound.stop(1)
	}
}

func (this *Apos) InRecovering() bool {
	return atomic.LoadInt32(&this.recovering) > 0
}

//Create The Credential
func (this *Apos) makeCredential(s int, sp SortitionPriority, csTime *types.BigInt, w, W int64) *CredentialSign {
	r := this.commonTools.GetNextRound()
	c := new(CredentialSign)
	c.Signature.init()
	c.Round = uint64(r)
	c.Step = uint64(s)
	c.ParentHash = this.commonTools.GetNowBlockHash()
	c.Time = csTime

	err := this.commonTools.Sig(c)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	var tao int64
	if s == STEP_BP {
		tao = Config().tProposer
	} else if s == StepFinal {
		tao = Config().tFinal
	} else {
		tao = Config().tStep
	}

	c.votes = sp.getSortitionPriorityByHash(c.Signature.Hash(), w, tao, W)
	logger.Debug(COLOR_FRONT_GREEN, "***Credential Votes Show:  Round:", c.Round, " Step:", c.Step, "  Votes:", c.votes, COLOR_SHORT_RESET)

	return c
}

func (this *Apos) judgeVerifier(cs *CredentialSign, setp int) bool {
	if cs.votes > 0 {
		return true
	} else {
		return false
	}
}
