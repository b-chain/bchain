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
// @File: count_vote.go
// @Date: 2018/07/19 13:34:19
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"errors"
	"math/big"
	"bchain.io/common"
	"bchain.io/common/types"
	"sync"
	"time"
)

type targetVotes struct {
	total  float64
	detail []*CredentialSign
}

type stepVotes struct {
	svLock sync.RWMutex
	counts map[types.Hash]*targetVotes
	//flag for vote result
	isFinish bool
	value    types.Hash
}

func (this *stepVotes) getTargetVotes(h types.Hash) *targetVotes {
	this.svLock.RLock()
	defer this.svLock.RUnlock()

	if tv, ok := this.counts[h]; ok {
		return tv
	}
	return nil
}

func (this *stepVotes) setTargetVotes(h types.Hash, tv *targetVotes) {
	this.svLock.Lock()
	defer this.svLock.Unlock()

	this.counts[h] = tv
}

func newStepVotes() *stepVotes {
	sv := new(stepVotes)
	sv.counts = make(map[types.Hash]*targetVotes)
	return sv
}

type countVote struct {
	cvLock            sync.RWMutex
	voteRecord        map[int]*stepVotes //map[step]*stepVotes
	CommonCoinMinHash map[int]types.Hash
	msgCh             chan *ByzantineAgreementStar
	stopCh            chan interface{}
	timer             *time.Timer
	timerStep         uint
	emptyBlock        types.Hash
	bbaFinish         bool
	bbaFinishStep     int
	sendVoteResult    func(s int, hash types.Hash)
	hangForeverFn     func()
}

func (cv *countVote) getCommonCoinMinHash(step int) (types.Hash, error) {
	cv.cvLock.RLock()
	defer cv.cvLock.RUnlock()

	if h, ok := cv.CommonCoinMinHash[step]; ok {
		return h, nil
	}
	return types.Hash{}, errors.New("none exist")
}

func (cv *countVote) setCommonCoinMinHash(step int, h types.Hash) {
	cv.cvLock.Lock()
	defer cv.cvLock.Unlock()

	//need check exist?
	cv.CommonCoinMinHash[step] = h
}
func (cv *countVote) getStepVotes(step int) *stepVotes {
	cv.cvLock.RLock()
	defer cv.cvLock.RUnlock()

	if sv, ok := cv.voteRecord[step]; ok {
		return sv
	}
	return nil
}

func (cv *countVote) setStepVotes(step int, sv *stepVotes) {
	cv.cvLock.Lock()
	defer cv.cvLock.Unlock()
	//need check exist?
	cv.voteRecord[step] = sv
}

func newCountVote(sendVoteResult func(s int, hash types.Hash), hangForeverFn func(), emptyBlock types.Hash) *countVote {
	cv := new(countVote)
	cv.init()
	cv.sendVoteResult = sendVoteResult
	cv.hangForeverFn = hangForeverFn
	cv.emptyBlock = emptyBlock
	return cv
}

func (cv *countVote) init() {
	cv.voteRecord = make(map[int]*stepVotes)
	cv.msgCh = make(chan *ByzantineAgreementStar, 1)
	cv.stopCh = make(chan interface{}, 1)
	cv.CommonCoinMinHash = make(map[int]types.Hash)
	cv.timerStep = STEP_IDLE
	cv.timer = time.NewTimer(time.Duration(2*Config().delayBlock) * time.Second)
}

//this function should be called by BP handle
func (cv *countVote) startTimer(delay int) {
	delayDuration := time.Second * time.Duration(delay)
	cv.timer.Reset(delayDuration)
	cv.timerStep = uint(cv.getNextTimerStep(STEP_BP))
}

func (cv *countVote) run() {
	for {
		select {
		// receive message
		case voteMsg := <-cv.msgCh:
			//add votes and check the votes whether complete
			step, hash, complete := cv.processMsg(voteMsg)
			logger.Debug(COLOR_FRONT_RED, "CountVote Get Msg:", "step:", step, "  hash:", hash.Hex(), COLOR_SHORT_RESET)
			if complete {
				logger.Info(COLOR_FRONT_RED, "CountVote Complete Step:", step, "hash", hash.Hex(), COLOR_SHORT_RESET)
				cv.countSuccess(step, hash)
			}
		//timeout message
		case <-cv.timer.C:
			logger.Info(COLOR_FRONT_RED, "countVote timeout, step", cv.timerStep, COLOR_SHORT_RESET)
			cv.timeoutHandle()
		case <-cv.stopCh:
			logger.Debug(COLOR_FRONT_RED, "countVote run exit", cv.timerStep, COLOR_SHORT_RESET)
			cv.timer.Stop()
			return
		}
	}
}

/*
normal scene:
timeoutStep = nowStep skip to next logic step


special scene:
nowStep = stepFinal
timeouStep = StepIdle

when countVote.bbaFinish == true,
skip to stepFinal

when next logic Step finish,
continue step searching

*/
func (cv *countVote) getNextTimerStep(step int) int {
	timeoutStep := step
	for {
		switch {
		case timeoutStep == STEP_BP:
			timeoutStep = STEP_REDUCTION_1
		case timeoutStep == STEP_REDUCTION_1:
			timeoutStep = STEP_REDUCTION_2
		case timeoutStep == STEP_REDUCTION_2:
			timeoutStep = 1
		case timeoutStep < int(Config().maxStep):
			timeoutStep++
			if timeoutStep == int(Config().maxStep) {
				timeoutStep = STEP_IDLE
			}
		case timeoutStep == STEP_FINAL:
			timeoutStep = STEP_IDLE
		default:
			//ignore
			timeoutStep = STEP_IDLE
		}
		//if the bba is finished , turn to STEP_FINAL
		if timeoutStep < int(Config().maxStep) && cv.bbaFinish {
			timeoutStep = STEP_FINAL
		}
		sv := cv.getStepVotes(timeoutStep)
		if sv != nil && sv.isFinish == true {
			continue
		}
		//if future step is finished before current step , continue find a timeoutStep
		//if sv, ok := cv.voteRecord[timeoutStep]; ok {
		//	if sv.isFinish == true {
		//		continue
		//	}
		//}
		break
	}
	return timeoutStep
}

func (cv *countVote) timeoutHandle() {
	timeoutStep := int(cv.timerStep)

	if timeoutStep == STEP_IDLE {
		logger.Debug(COLOR_FRONT_RED, "timeoutHandle timeoutStep == STEP_IDLE , Step", timeoutStep, COLOR_SHORT_RESET)
		//ignore
		return
	}

	logger.Debug(COLOR_FRONT_RED, "timeoutHandle  Step:", timeoutStep, COLOR_SHORT_RESET)
	//fill results in voteRecord
	sv := cv.getStepVotes(timeoutStep)
	if sv == nil {
		svNew := newStepVotes()
		svNew.isFinish = true
		svNew.value = TimeOut
		cv.setStepVotes(timeoutStep, svNew)
	} else {
		sv.isFinish = true
		sv.value = TimeOut
	}

	resetTimer := true
	nextTimoutStep := cv.getNextTimerStep(timeoutStep)
	cv.timerStep = uint(nextTimoutStep)
	if nextTimoutStep == STEP_IDLE {
		resetTimer = false
	}

	if resetTimer {
		delay := time.Second * time.Duration(Config().delayStep)
		cv.timer.Reset(delay)
	}

	cv.sendVoteResult(timeoutStep, TimeOut)

	if timeoutStep == int(Config().maxStep-1) && nextTimoutStep == STEP_IDLE {
		cv.hangForever()
	}
}

func (cv *countVote) countSuccess(step int, hash types.Hash) {
	//send result
	cv.sendVoteResult(step, hash)

	resetTimer := false
	nextTimoutStep := STEP_IDLE
	if int(cv.timerStep) == step {
		//reset timer,this step has been completed
		resetTimer = true
		nextTimoutStep = cv.getNextTimerStep(step)
	}

	if step < int(Config().maxStep) {
		bbaIdex := step % 3
		if bbaIdex == 1 && hash != cv.emptyBlock {
			//bba complete: block hash
			cv.bbaFinish = true
			cv.bbaFinishStep = step
			if cv.timerStep < (Config().maxStep) {
				nextTimoutStep = STEP_FINAL
				resetTimer = true
			}
		} else if bbaIdex == 2 && hash == cv.emptyBlock {
			//bba complete: empty block hash
			cv.bbaFinish = true
			cv.bbaFinishStep = step
			if cv.timerStep < (Config().maxStep) {
				nextTimoutStep = STEP_FINAL
				resetTimer = true
			}
		}
	}

	//bba last step success, and timer in current step, bba is not finished
	//mean that bba will hang for ever
	if step == int(Config().maxStep-1) && int(cv.timerStep) == step && cv.bbaFinish == false {
		logger.Warn("bba last step can not make bba finished, hangForever...")
		cv.hangForever()
		resetTimer = false
	}

	if resetTimer {
		cv.timerStep = uint(nextTimoutStep)
		delay := time.Second * time.Duration(Config().delayStep)
		if nextTimoutStep != STEP_IDLE {
			cv.timer.Reset(delay)
		}
	}
}

func (cv *countVote) addVotes(ba *ByzantineAgreementStar) (types.Hash, float64) {
	hash := ba.Hash
	step := int(ba.Credential.Step)
	votes := ba.Credential.votes

	sv := cv.getStepVotes(step)
	if sv == nil {
		svNew := newStepVotes()
		targetVotes := &targetVotes{votes, make([]*CredentialSign, 0)}
		targetVotes.detail = append(targetVotes.detail, ba.Credential)
		svNew.setTargetVotes(hash, targetVotes)
		cv.setStepVotes(step, svNew)
		//cv.voteRecord[step] = svNew
		return hash, votes
	} else {
		hashVote := sv.getTargetVotes(hash)
		if hashVote != nil {
			sumVote := hashVote.total + votes
			hashVote.total = sumVote
			hashVote.detail = append(hashVote.detail, ba.Credential)
			return hash, sumVote
		} else {
			targetVotes := &targetVotes{votes, make([]*CredentialSign, 0)}
			targetVotes.detail = append(targetVotes.detail, ba.Credential)
			sv.setTargetVotes(hash, targetVotes)
			return hash, votes
		}
	}
}

func (cv *countVote) getCommonCoinMinHashRslt(step int) int {

	minhash, err := cv.getCommonCoinMinHash(step)

	if err != nil {
		return 1
	}
	lastByte := uint(minhash[31])
	return int(lastByte % 2)
}
func (cv *countVote) calCommonCoinMinHash(step int, hash types.Hash, votes int) {
	//get currentMinhash
	var minhash types.Hash
	minhash, err := cv.getCommonCoinMinHash(step)
	if err != nil {
		minhash := TimeOut
		cv.setCommonCoinMinHash(step, minhash)
	}

	cnt := votes //because c is float64,but here we need int
	for j := 1; j <= cnt; j++ {
		cch := &CommonCoinMinHash{hash, j}
		hashCch, err := common.MsgpHash(cch)
		if err != nil {
			continue
		}

		if new(big.Int).SetBytes(hashCch[:]).Cmp(new(big.Int).SetBytes(minhash[:])) < 0 {
			minhash = hashCch
			cv.setCommonCoinMinHash(step, minhash)
		}
	}
}

func (cv *countVote) processMsg(ba *ByzantineAgreementStar) (int, types.Hash, bool) {
	logger.Trace("processMsg", ba.Credential.Step, ba.Hash.String())
	step := int(ba.Credential.Step)
	//check bba is finished or not
	if step < int(Config().maxStep) && cv.bbaFinish {
		logger.Info("all bba finished. step ", step, "will ignore")
		return step, types.Hash{}, false
	}
	//check stepVotes is finished or not
	sv := cv.getStepVotes(step)
	if sv != nil {
		//check this step whether is finish
		if sv.isFinish {
			logger.Info("step", step, "is finished, ignore vote")
			return step, types.Hash{}, false
		}
	}

	hash, votes := cv.addVotes(ba)
	//checked again after addVotes
	sv = cv.getStepVotes(step)
	//sv = cv.voteRecord[step]
	logger.Info(COLOR_FRONT_PINK, "ProcessMsg step", step, "getThreshold:", getThreshold(step), "Now Votes:", votes, COLOR_SHORT_RESET)

	if votes > float64(getThreshold(step)) {
		sv.isFinish = true
		sv.value = hash
		return step, hash, true
	}
	if step < int(Config().maxStep) && step%3 == 0 {
		cv.calCommonCoinMinHash(step, ba.Credential.Signature.Hash(), int(ba.Credential.votes)+1)
	}

	return step, hash, false
}

func (cv *countVote) getBlockCertificate() BlockCertificate {
	if cv.bbaFinish == false {
		logger.Warn("getBlockCertificate fail, bba is not stop")
		return nil
	}
	sv := cv.getStepVotes(cv.bbaFinishStep)
	hash := sv.value

	tv := sv.getTargetVotes(hash)
	if tv == nil {
		return nil
	}
	return tv.detail
}

func (cv *countVote) hangForever() {
	logger.Warn("bba last step timeout or not leading bba finished, hangForever...")
	cv.hangForeverFn()
}

func (cv *countVote) sendMsg(ba *ByzantineAgreementStar) {
	cv.msgCh <- ba
}

func (cv *countVote) stop() {
	cv.stopCh <- 1
}
