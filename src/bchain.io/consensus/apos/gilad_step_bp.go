package apos

import (
	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
	"sort"
	"sync"
	"time"
)

type BpWithPriority struct {
	j  float64 //the priofity
	bp *BlockProposal
}

type BpWithPriorityHeap []*BpWithPriority

func (h BpWithPriorityHeap) Len() int           { return len(h) }
func (h BpWithPriorityHeap) Less(i, j int) bool {

	if h[i].j > h[j].j {
		return true
	} else if h[i].j < h[j].j {
		return false
	}
	bigI := h[i].bp.Credential.sigHashHashBig()
	bigJ := h[j].bp.Credential.sigHashHashBig()

	ret := bigI.Cmp(bigJ)
	if ret > 0 {
		return true
	}
	return false

}
func (h BpWithPriorityHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *BpWithPriorityHeap) Push(x interface{}) {
	*h = append(*h, x.(*BpWithPriority))
}

func (h *BpWithPriorityHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type BpObj struct {
	lock        sync.RWMutex
	BpHeap      BpWithPriorityHeap
	totalPri    float64
	existMap    map[types.Hash]*BlockProposal //todo:here should be change to *bp
	msgChan     chan *BlockProposal
	exit        chan interface{}
	ctx         *stepCtx
	priorityBp  *BlockProposal
	nothingTodo bool
}

func makeBpObj(ctx *stepCtx) *BpObj {
	s := new(BpObj)
	s.ctx = ctx
	s.BpHeap = make(BpWithPriorityHeap, 0)
	s.msgChan = make(chan *BlockProposal)
	s.existMap = make(map[types.Hash]*BlockProposal)
	s.exit = make(chan interface{}, 2)
	s.totalPri = 0
	logger.Debug(COLOR_FRONT_PINK, "makeBpObj", COLOR_SHORT_RESET)
	return s
}

func (this *BpObj) isExistBlock(blockHash types.Hash) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if _, ok := this.existMap[blockHash]; ok {
		return true
	}
	return false
}

func (this *BpObj) addExistBlock(bp *BlockProposal) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.existMap[bp.Block.Hash()] = bp
}

func (this *BpObj) getExistBlock(blockHash types.Hash) *block.Block {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if bp, ok := this.existMap[blockHash]; ok {
		return bp.Block
	}
	return nil
}
func (this *BpObj) CommitteeVote(data *VoteData) {
	//get credential from Round
	cret := this.ctx.getCredentialByStep(uint64(data.Step))
	if cret == nil {
		return
	}
	if cret.votes <= 0 {
		return
	}
	//todo :need pack ba msg

	msgBa := newByzantineAgreementStar()
	//hash
	msgBa.Hash = data.Value
	//Credential
	msgBa.Credential = cret

	//Esig
	msgBa.Esig.round = msgBa.Credential.Round
	msgBa.Esig.step = msgBa.Credential.Step
	msgBa.Esig.val = make([]byte, 0)
	msgBa.Esig.val = append(msgBa.Esig.val, msgBa.Hash.Bytes()...)

	err := this.ctx.esig(msgBa.Esig)
	if err != nil {
		logger.Error("CommitteeVote Esig Err:", err.Error())
		return
	}

	this.ctx.sendInner(msgBa)

}

func (this *BpObj) makeBlock() {
	bp := newBlockProposal()
	bp.Credential = this.ctx.getCredentialByStep(StepBp)
	if nil == bp.Credential {
		logger.Warn("makeBlock getCredentialByStep--->nil")
		return
	}

	bcd := this.ctx.makeBlockConsensusData(bp)

	bp.Block = this.ctx.getProducerNewBlock(bcd ,int64(Config().delayBlock - 2) )

	if bp.Block == nil {
		logger.Error("makeBlock getProducerNewBlock return nil")
		return
	}

	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = StepBp
	bp.Esig.val = make([]byte, 0)
	h := bp.Block.Hash()
	bp.Esig.val = append(bp.Esig.val, h[:]...)
	err := this.ctx.esig(bp.Esig)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	//delay here ?

	this.ctx.sendInner(bp)

	logger.Debug(COLOR_FRONT_PINK, "***[A]Out M1 CreHash:", bp.Credential.Signature.Hash().String(), " BlockHash", bp.Block.H.Hash().String(), COLOR_SHORT_RESET)
}
func (this *BpObj) stop() {
	logger.Debug(COLOR_FRONT_PINK, "Call BpObj Exit....:", COLOR_SHORT_RESET)
	this.exit <- 1
}
func (this *BpObj) run() {
	rd := this.ctx.getRound()
	//make block and send out
	go this.makeBlock()
	tProposer := int(Config().tProposer)
	timer := time.NewTicker(time.Duration(Config().delayBlock) * time.Second)
	defer timer.Stop()
	logger.Debug(COLOR_FRONT_PINK, "#########In BpObj-", rd, COLOR_SHORT_RESET)
	defer func() {
		logger.Debug(COLOR_FRONT_PINK, "#######Out BpObj", rd, COLOR_SHORT_RESET)
	}()
	for {
		select {
		case <-this.exit:
			logger.Debug(COLOR_FRONT_PINK, "Round:", rd, "   BpObj Exit....:return", COLOR_SHORT_RESET)
			return
		case <-timer.C:
			if this.nothingTodo {
				continue
			}
			value := types.Hash{}
			if this.BpHeap.Len() == 0 {
				//output empty hash
				value = this.ctx.getEmptyBlockHash()
			} else {
				sort.Sort(&this.BpHeap)
				value = this.BpHeap[0].bp.Block.Hash()
				//make reduction input data
			}
			vd := new(VoteData)
			vd.Round = this.ctx.getRound()
			vd.Step = StepReduction1
			vd.Value = value
			logger.Debug(COLOR_FRONT_PINK, "BpObj timeOut dataOutput hash:", vd.Value.Hex(), COLOR_SHORT_RESET)
			this.CommitteeVote(vd)

			this.ctx.setBpResult(value)
			this.ctx.startVoteTimer(int(Config().delayStep))
			this.nothingTodo = true
		case bp := <-this.msgChan:

			//logic do
			//verify the block
			if !this.ctx.verifyBlock(bp.Block) {
				logger.Debug(COLOR_FRONT_PINK, "!this.ctx.verifyBlock(bp.Block) Wrong:hash:", bp.Block.Hash().Hex(), COLOR_SHORT_RESET)
				continue
			}
			//check is exist a same block
			if this.isExistBlock(bp.Block.Hash()) {
				logger.Debug(COLOR_FRONT_PINK, "this.isExistBlock(bp.Block.Hash()) Wrong:hash:", bp.Block.Hash().Hex(), COLOR_SHORT_RESET)
				continue
			}
			this.addExistBlock(bp)

			if this.nothingTodo {
				continue
			}

			//check the node has the right to produce a block
			pri := bp.Credential.votes
			bpp := new(BpWithPriority)
			bpp.j = pri
			bpp.bp = bp

			this.BpHeap.Push(bpp)

			this.totalPri += pri
			if this.priorityBp == nil {
				this.ctx.propagateMsg(bp)
			} else if pri > this.priorityBp.Credential.votes {
				this.priorityBp = bp
				this.ctx.propagateMsg(bp)
			} else if pri == this.priorityBp.Credential.votes {
				if bp.Credential.sigHashHashBig().Cmp(this.priorityBp.Credential.sigHashHashBig()) > 0 {
					this.priorityBp = bp
					this.ctx.propagateMsg(bp)
				}
			}

			if this.totalPri >= float64(tProposer) {
				sort.Sort(&this.BpHeap)
				//get the bigger one
				x := this.BpHeap[0]
				_ = x

				vd := new(VoteData)
				vd.Round = x.bp.Credential.Round
				vd.Step = StepReduction1
				vd.Value = x.bp.Block.Hash()
				logger.Debug(COLOR_FRONT_PINK, "BpObj >tProposer dataOutput hash:", vd.Value.Hex(), COLOR_SHORT_RESET)
				this.CommitteeVote(vd)
				this.ctx.startVoteTimer(int(Config().delayStep))
				this.ctx.setBpResult(x.bp.Block.Hash())
				//todo:inform the reduction

				this.nothingTodo = true
			}
		}
	}
}

func (this *BpObj) sendMsg(bp *BlockProposal) {
	this.msgChan <- bp
}
