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
// @File: recovery.go
// @Date: 2018/08/09 14:59:09
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"errors"
	"github.com/hashicorp/golang-lru"
	"math/big"
	"bchain.io/common/types"
	"reflect"
	"sync"
	"time"
)

var (
	ErrDuplicateMsg          = errors.New("recovery duplicate credential message")
	ErrDuplicatePropagateMsg = errors.New("recovery duplicate Propagate credential message")
)

const (
	nonCsCacheLimit = 10240
)

type recoveryRound struct {
	round       uint64
	apos        *Apos
	credentials map[int]*CredentialSign
	cretLock    sync.RWMutex
	msgs        map[types.Address]*peerMsgs
	roundOverCh chan interface{}
	stopCh      chan int
	isStop      bool
	parentHash  types.Hash
	countVote   *countVote
	lock        sync.RWMutex
	nonCsCache  *lru.Cache // Cache for the bp and bba* message
}

func newRecoveryRound(round int, parentHash types.Hash, apos *Apos, roundOverCh chan interface{}) *recoveryRound {
	r := new(recoveryRound)
	r.init(round, apos, roundOverCh)
	r.parentHash = parentHash
	r.isStop = false
	return r
}

func (this *recoveryRound) init(round int, apos *Apos, roundOverCh chan interface{}) {
	this.round = uint64(round)
	this.apos = apos
	this.roundOverCh = roundOverCh
	this.stopCh = make(chan int, 1)
	this.credentials = make(map[int]*CredentialSign)
	this.msgs = make(map[types.Address]*peerMsgs)
	nonCsCache, _ := lru.New(nonCsCacheLimit)
	this.nonCsCache = nonCsCache
}

func (this *recoveryRound) stop(val int) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.isStop == false {
		this.stopCh <- val
		this.isStop = true
	}
}

func (this *recoveryRound) run() {
	// make verifiers Credential
	this.generateCredentials()

	// broadcast Credentials
	this.broadcastCredentials(false)

	this.startStepObjs()

	this.commonProcess()
	this.roundOverCh <- 1 //inform the caller,the mission complete
}

// Generate valid Credentials in current round
func (this *recoveryRound) generateCredentials() {
	sp := Config().sp
	csTime := types.NewBigInt(*big.NewInt(time.Now().Unix()))
	addr := this.apos.commonTools.GetCoinBase()
	w, W := this.apos.commonTools.GetWeight(this.round, addr)

	logger.Info(COLOR_FRONT_YELLOW, "CoinBase:", addr.Hex(), "w:", w, " W:", W, COLOR_SHORT_RESET)
	for i := 1; i < int(Config().maxStep); i++ {
		credential := this.apos.makeCredential(i, sp, csTime, w, W)
		isVerifier := this.apos.judgeVerifier(credential, i)
		if isVerifier {
			this.credentials[i] = credential
			this.apos.commonTools.CreateTmpPriKey(i)
		}
	}
	for i := STEP_BP; i < STEP_IDLE; i++ {
		credential := this.apos.makeCredential(i, sp, csTime, w, W)
		isVerifier := this.apos.judgeVerifier(credential, i)
		if isVerifier {
			this.credentials[i] = credential
			this.apos.commonTools.CreateTmpPriKey(i)
		}
	}
}

func (this *recoveryRound) broadcastCredentials(refreshTime bool) {
	this.cretLock.Lock()
	defer this.cretLock.Unlock()
	newTime := types.NewBigInt(*big.NewInt(time.Now().Unix()))
	for i, credential := range this.credentials {
		if refreshTime {
			credential.Time = newTime
			this.apos.outMsger.PropagateMsg(credential)
		} else {
			this.apos.outMsger.SendInner(credential)
		}
		logger.Info("send credential recovery round", this.round, "step", i, "time", credential.Time.IntVal.String())
	}
}

func (this *recoveryRound) startStepObjs() {
	this.countVote = newCountVote(this.recoveryOk, nil, types.Hash{})
	if this.countVote == nil {
		logger.Error("this.countVote == nil...........")
	}

	go this.countVote.run()
}

func (this *recoveryRound) recoveryOk(step int, hash types.Hash) {
	logger.Info("recovery round", this.round, "step", step, "success! exit")
	go this.stop(0)
}

func (this *recoveryRound) commonProcess() {
	timer := time.NewTicker(time.Second * time.Duration(Config().delayRecovery))
	defer timer.Stop()
	for {
		select {
		// receive message
		case outData := <-this.apos.outMsger.GetDataMsg():
			switch v := outData.(type) {
			case *CredentialSign:
				this.receiveMsgCs(v)
			case *BlockProposal:
				logger.Info("recovery buffer loop back message, type bp", v.Credential.Round)
				this.bufferNormalMsg(v.Block.Hash(), v)
			case *ByzantineAgreementStar:
				logger.Info("recovery buffer loop back message, type bba*", v.Credential.Round, v.Credential.Step)
				this.bufferNormalMsg(v.BaHash(), v)
			default:
				logger.Debug("recovery ignore message, type ", reflect.TypeOf(v))
			}
		case val := <-this.stopCh:
			logger.Info(COLOR_FRONT_RED, "recoveryRound run exit", val, this.round, COLOR_SHORT_RESET)
			this.broadCastStop()
			if val == 0 {
				this.sendNormalMsg()
			}
			//this.nonCsCache.Purge()
			return
		case <-timer.C:
			go this.broadcastCredentials(true)
		}
	}
}

// bp and bba message should buffer, when recovery round finish, normal round should process
func (this *recoveryRound) bufferNormalMsg(hash types.Hash, msg interface{}) {
	if !this.nonCsCache.Contains(hash) {
		this.nonCsCache.Add(hash, msg)
	}
}

func (this *recoveryRound) sendNormalMsg() {
	for _, hash := range this.nonCsCache.Keys() {
		if msg, exist := this.nonCsCache.Peek(hash); exist {
			this.apos.outMsger.SendInner(msg)
		}
	}
}

func (this *recoveryRound) receiveMsgCs(msg *CredentialSign) {
	logger.Debug("Receive message CredentialSign [r:s]:", msg.Round, msg.Step)
	if msg.Round != this.round {
		logger.Warn("verify fail, Credential msg is not in current round, want:", this.round, " but:", msg.Round)
		return
	}

	if msg.ParentHash != this.parentHash {
		logger.Warn("verify fail, Credential msg is not in current block chain", msg.ParentHash.String(), this.parentHash.String())
		return
	}

	//prevent DDOS, p2p global propagate time is within 4 second, peer's time too far, should ignore
	nowTime := big.NewInt(time.Now().Unix())
	diffTime := new(big.Int).Sub(&msg.Time.IntVal, nowTime)
	absDiff := diffTime.Abs(diffTime).Uint64()
	if absDiff > 15 {
		logger.Warn("Credential time is not sync with local time", msg.Time.IntVal.String())
		return
	}

	//duplicate message check
	if err := this.filterMsgCs(msg); err != nil {
		logger.Debug("filter Credential fail", err)
		if err == ErrDuplicatePropagateMsg {
			this.apos.outMsger.PropagateMsg(msg)
		}
		return
	}

	//Propagate message via p2p
	this.apos.outMsger.PropagateMsg(msg)

	ba := &ByzantineAgreementStar{
		Hash:       types.Hash{},
		Credential: msg,
	}
	this.countVote.sendMsg(ba)
}

//duplicate message check
func (this *recoveryRound) filterMsgCs(msg *CredentialSign) error {
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
				//prevent DDOS, time must increase with seconds
				if msg.Time.IntVal.Uint64()-mCs.Time.IntVal.Uint64() >= 5 {
					mCs.Time = msg.Time
					return ErrDuplicatePropagateMsg
				} else {
					return ErrDuplicateMsg
				}
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

func (this *recoveryRound) broadCastStop() {
	logger.Debug(COLOR_FRONT_RED, "In BroadCastStop...", COLOR_SHORT_RESET)
	this.countVote.stop()
}
