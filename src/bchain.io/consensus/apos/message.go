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
// @File: message.go
// @Date: 2018/06/22 14:40:22
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"errors"
	"fmt"
	"bchain.io/common"
	"bchain.io/common/types"
	"bchain.io/consensus/message"
	"bchain.io/core"
	"bchain.io/core/blockchain/block"
	"bchain.io/utils/event"
	"reflect"
	"sync"
	"time"
)

const (
	STEP_BP = iota + 0xffff
	STEP_REDUCTION_1
	STEP_REDUCTION_2
	STEP_FINAL
	STEP_IDLE
)

//go:generate msgp
/*
bufferMsg:the msg come from other peer with higher block number than currentNumber,
the msg should not be dealed now or discarded , how to do that,so ,we should block the msg
until the currentNum >= msg.Round - 1,the function bufferMsg will return ,and the msg will be dealed by
main task.
*/
func bufferMsg(peerNumber uint64, chainSub event.Subscription, chainChan chan core.ChainEvent) {
	currentNumber := gCommonTools.GetNowBlockNum()
	if peerNumber > currentNumber {
		logger.Debug("need buffer msg", peerNumber, currentNumber)
		timer := time.NewTimer(60 * time.Second)
		defer timer.Stop()
		//future msg, need buffer
		chainSub = gCommonTools.SubscribeChainEvent(chainChan)
		defer chainSub.Unsubscribe()
		for {
			select {
			case event := <-chainChan:
				if peerNumber <= event.Block.NumberU64() {
					return
				}
				// Err() channel will be closed when unsubscribing.
			case <-chainSub.Err():
				return
			case <-timer.C:
				logger.Debug("msg wait too long, ignore", peerNumber, currentNumber)
				return
			}
		}
	}
}

func (cs *CredentialSign) validate() (types.Address, error) {

	//1. validate parentHash
	parentBlock := gCommonTools.GetBlockByHash(cs.ParentHash)
	if parentBlock == nil {
		return types.Address{}, errors.New(fmt.Sprintf("verify CredentialSig fail: Round %d can't get block form hash %s", cs.Round, cs.ParentHash.Hex()))
	}

	//2. round check
	if parentBlock.H.Number.IntVal.Uint64()+1 != cs.Round {
		return types.Address{}, errors.New(fmt.Sprintf("verify CredentialSig fail: Round %d is not equal block number", cs.Round))

	}

	//3.signature check and sender check
	sender, err := cs.sender()
	if err != nil {
		return types.Address{}, errors.New(fmt.Sprintf("verify CredentialSig fail: %s", err))
	}

	//4.sortition check
	//todo 2. validate right
	sigHash := cs.Signature.Hash()

	var tao int64
	if cs.Step == STEP_BP {
		tao = Config().tProposer
	} else if cs.Step == StepFinal {
		tao = Config().tFinal
	} else {
		tao = Config().tStep
	}

	addr := sender
	w, W := gCommonTools.GetWeight(cs.Round, addr)

	cs.votes = Config().sp.getSortitionPriorityByHash(sigHash, w, tao, W)
	logger.Debug("message", "r s vote", cs.Round, cs.Step, cs.votes, tao)
	if cs.votes <= 0 {
		return types.Address{}, errors.New(fmt.Sprintf("verify votes fail: Round %d peer %v no ritht %f to verify", cs.Round, sender.Hash(), cs.votes))
	}

	return sender, nil
}

type msgCredentialSig struct {
	cs *CredentialSign
	*message.MsgPriv
	chainChan chan core.ChainEvent
	chainSub  event.Subscription
}

func NewMsgCredential(c *CredentialSign) *msgCredentialSig {
	msgCs := &msgCredentialSig{
		cs:        c,
		MsgPriv:   message.NewMsgPriv(),
		chainChan: make(chan core.ChainEvent, 10),
	}
	message.Msgcore().Handle(msgCs)
	return msgCs
}

func (c *msgCredentialSig) DataHandle(data interface{}) {
	mt := MsgTransfer()
	if mt.aposRunning == false {
		return
	}
	logger.Debug("msgCredentialSig data handle")
	bufferMsg(c.cs.Round-1, c.chainSub, c.chainChan)
	if _, err := c.cs.validate(); err != nil {
		logger.Info("message CredentialSig validate error:", err)
		return
	}
	mt.Send2Apos(c.cs)
}

func (c *msgCredentialSig) StopHandle() {
	logger.Debug("msgCredentialSig stop ...")
}

// step1 (Block Proposal) message
// m(r,1) = (Br, esig(H(Br)), σr1)
type BlockProposal struct {
	Block      *block.Block
	Esig       *EphemeralSign
	Credential *CredentialSign
}

func newBlockProposal() *BlockProposal {
	b := new(BlockProposal)
	b.Esig = new(EphemeralSign)
	return b
}

func (bp *BlockProposal) validate() error {
	//verify step
	if bp.Credential.Step != STEP_BP {
		return errors.New(fmt.Sprintf("Block Proposal step is not 1: %d", bp.Credential.Step))
	}

	//verify Credential
	cretSender, err := bp.Credential.validate()
	if err != nil {
		return err
	}

	//verify ephemeral signature,check  sender is equal CretSender or not
	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = bp.Credential.Step
	bp.Esig.val = bp.Block.Hash().Bytes()
	sender, err := bp.Esig.sender()
	if err != nil {
		return errors.New(fmt.Sprintf("BP verify ephemeral signature fail: %s", err))
	}
	if cretSender != sender {
		logger.Debug("Block Proposal Ephemeral signature address is not equal to Credential signature address", sender.Hex(), cretSender.Hex())
		return errors.New("sender's address between Credential and Ephemeral is not equal")
	}

	//todo block validate

	return nil
}

type msgBlockProposal struct {
	bp *BlockProposal
	*message.MsgPriv
	chainChan chan core.ChainEvent
	chainSub  event.Subscription
}

// new a message
func NewMsgBlockProposal(bp *BlockProposal) *msgBlockProposal {
	msgBp := &msgBlockProposal{
		bp:        bp,
		MsgPriv:   message.NewMsgPriv(),
		chainChan: make(chan core.ChainEvent, 10),
	}
	message.Msgcore().Handle(msgBp)
	return msgBp
}

func (bp *msgBlockProposal) DataHandle(data interface{}) {
	mt := MsgTransfer()
	if mt.aposRunning == false {
		return
	}
	logger.Debug("msgBlockProposal data handle")
	bufferMsg(bp.bp.Credential.Round-1, bp.chainSub, bp.chainChan)
	if err := bp.bp.validate(); err != nil {
		logger.Info("message BlockProposal validate error:", err)
		return
	}
	mt.Send2Apos(bp.bp)
}

func (bp *msgBlockProposal) StopHandle() {
	logger.Debug("msgBlockProposal stop ...")
}

type ByzantineAgreementStar struct {
	Hash       types.Hash     //voted block's hash.
	Esig       *EphemeralSign //the signature of somebody's ephemeral secret key
	Credential *CredentialSign
}

func newByzantineAgreementStar() *ByzantineAgreementStar {
	b := new(ByzantineAgreementStar)
	b.Esig = new(EphemeralSign)
	return b
}

func (ba *ByzantineAgreementStar) validate() error {
	//verify step
	if ba.Credential.Step < 1 || uint(ba.Credential.Step) >= Config().maxStep {
		if ba.Credential.Step != STEP_REDUCTION_1 && ba.Credential.Step != STEP_REDUCTION_2 && ba.Credential.Step != STEP_FINAL {
			return errors.New(fmt.Sprintf("Byzantine Agreement Star step is not right: %d", ba.Credential.Step))
		}
	}
	//verify Credential
	cretSender, err := ba.Credential.validate()
	if err != nil {
		return err
	}

	//verify ephemeral signature
	ba.Esig.round = ba.Credential.Round
	ba.Esig.step = ba.Credential.Step
	ba.Esig.val = ba.Hash.Bytes()
	sender, err := ba.Esig.sender()
	if err != nil {
		return errors.New(fmt.Sprintf("BA* verify ephemeral signature fail: %s", err))
	}

	if cretSender != sender {
		logger.Debug("BA* Ephemeral hash signature address is not equal to Credential signature address", sender.Hex(), cretSender.Hex())
		return errors.New("sender's address between Credential and Hash Ephemeral is not equal")
	}

	return nil
}

func (ba *ByzantineAgreementStar) BaHash() types.Hash {
	hash, err := common.MsgpHash(ba)
	if err != nil {
		return types.Hash{}
	}
	return hash
}

type msgByzantineAgreementStar struct {
	ba *ByzantineAgreementStar
	*message.MsgPriv
	chainChan chan core.ChainEvent
	chainSub  event.Subscription
}

func NewMsgByzantineAgreementStar(ba *ByzantineAgreementStar) *msgByzantineAgreementStar {
	msgBba := &msgByzantineAgreementStar{
		ba:        ba,
		MsgPriv:   message.NewMsgPriv(),
		chainChan: make(chan core.ChainEvent, 10),
	}
	message.Msgcore().Handle(msgBba)
	return msgBba
}

func (ba *msgByzantineAgreementStar) DataHandle(data interface{}) {
	mt := MsgTransfer()
	if mt.aposRunning == false {
		return
	}
	logger.Debug("msgByzantineAgreementStar data handle", ba.ba.Credential.Round, ba.ba.Credential.Step)
	bufferMsg(ba.ba.Credential.Round-1, ba.chainSub, ba.chainChan)
	if err := ba.ba.validate(); err != nil {
		logger.Info("message ByzantineAgreementStar validate error:", err)
		return
	}
	mt.Send2Apos(ba.ba)
}

func (bba *msgByzantineAgreementStar) StopHandle() {
	logger.Debug("msgByzantineAgreementStar stop ...")
}

//message transfer between msg and Apos
type msgTransfer struct {
	receiveChan chan dataPack //receive message from BBa, Gc, Bp and etc.
	sendChan    chan dataPack

	aposRunning bool

	csFeed event.Feed
	bpFeed event.Feed
	baFeed event.Feed
	scope  event.SubscriptionScope
}

// about MsgTransfer singleton
var (
	msgTransferInstance *msgTransfer
	msgTransferOnce     sync.Once
)

// get the MsgTransfer singleton
func MsgTransfer() *msgTransfer {
	msgTransferOnce.Do(func() {
		msgTransferInstance = &msgTransfer{
			receiveChan: make(chan dataPack, 10),
			sendChan:    make(chan dataPack, 10),
			aposRunning: false,
		}
	})
	return msgTransferInstance
}

func (mt *msgTransfer) setAposRunning(v bool) {
	mt.aposRunning = v
}

func (mt *msgTransfer) GetDataMsg() <-chan dataPack {
	return mt.receiveChan
}

func (mt *msgTransfer) sendInner(data dataPack) {
	mt.receiveChan <- data
}

func (mt *msgTransfer) SendInner(data dataPack) error {
	//todo here need to validate process??
	logger.Debug("SendInner type:", reflect.TypeOf(data))
	go mt.sendInner(data)

	return nil
}

func (mt *msgTransfer) PropagateMsg(data dataPack) error {
	//logger.Debug("msgTransfer PropagateMsg in, data type:", reflect.TypeOf(data))
	switch v := data.(type) {
	case *CredentialSign:
		go mt.csFeed.Send(CsEvent{v})
	case *BlockProposal:
		go mt.bpFeed.Send(BpEvent{v})
	case *ByzantineAgreementStar:
		go mt.baFeed.Send(BaEvent{v})
	default:
		logger.Warn("in PropagateMsg invalid message type ", reflect.TypeOf(v))
	}
	return nil
}

func (mt *msgTransfer) Send2Apos(data dataPack) {
	mt.receiveChan <- data
}

//called by protocol manager,when we want send data to other peer,should call it before send(PropagateMsg)
func (mt *msgTransfer) SubscribeCsEvent(ch chan<- CsEvent) event.Subscription {
	return mt.scope.Track(mt.csFeed.Subscribe(ch))
}
func (mt *msgTransfer) SubscribeBpEvent(ch chan<- BpEvent) event.Subscription {
	return mt.scope.Track(mt.bpFeed.Subscribe(ch))
}
func (mt *msgTransfer) SubscribeBaEvent(ch chan<- BaEvent) event.Subscription {
	return mt.scope.Track(mt.baFeed.Subscribe(ch))
}

type CsEvent struct{ Cs *CredentialSign }
type BpEvent struct{ Bp *BlockProposal }
type BaEvent struct{ Ba *ByzantineAgreementStar }
