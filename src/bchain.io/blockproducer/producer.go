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
// @File: producer.go
// @Date: 2018/05/08 17:23:08
////////////////////////////////////////////////////////////////////////////////

package blockproducer

import (
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"crypto/ecdsa"
	"fmt"
	"bchain.io/common"
	"bchain.io/common/types"
	"bchain.io/consensus"
	"bchain.io/core"
	"bchain.io/core/actioncontext"
	"bchain.io/core/blockchain"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/state"
	"bchain.io/core/stateprocessor"
	"bchain.io/core/transaction"
	"bchain.io/params"
	"bchain.io/utils/database"
	"bchain.io/utils/event"
)

const (
	resultQueueSize     = 10
	producingLogAtDepth = 5

	// txChanSize is the size of channel listening to TxPreEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096
	// chainHeadChanSize is the size of channel listening to ChainHeadEvent.
	chainHeadChanSize = 10
	// chainSideChanSize is the size of channel listening to ChainSideEvent.
	chainSideChanSize = 10

	txsSizeLimit = common.StorageSize(10*1000*1000 - 2000)
)

// Agent can register themself with the worker
type Agent interface {
	Work() chan<- *Work
	SetReturnCh(chan<- *Result)
	Stop()
	Start()
	GetHashRate() int64
}

// Work is the workers current environment and holds
// all of the current state information
type Work struct {
	config        *params.ChainConfig
	signer        transaction.Signer
	state         *state.StateDB // apply state changes here
	db            database.IDatabase
	stateRootHash types.Hash
	tcount        int          // tx count in cycle
	Block         *block.Block // the new block
	header        *block.Header
	txs           []*transaction.Transaction
	failTxHashs   []types.Hash
	receipts      []*transaction.Receipt
	bchain          Backend
	createdAt     time.Time
	txTimeLimit   int64
	maxBlkSize    common.StorageSize
}

type Result struct {
	Work  *Work
	Block *block.Block
}

type blockRequest struct {
	data      *block.ConsensusData
	timeLimit int64
}

// worker is the main object which takes care of applying messages to the new state
type producer struct {
	config *params.ChainConfig
	engine consensus.Engine

	mu sync.Mutex

	// update loop
	mux *event.TypeMux
	wg  sync.WaitGroup

	agents map[Agent]struct{}
	recv   chan *Result

	bchain    Backend
	chain   *blockchain.BlockChain
	proc    blockchain.Validator
	chainDb database.IDatabase

	coinbase types.Address
	priKey   *ecdsa.PrivateKey
	//extra    []byte
	maxBlkSize uint64

	currentMu sync.Mutex
	current   *Work
	// atomic status counters
	producing          int32
	atWork             int32

	createRequestChan  chan blockRequest
	createResponseChan chan *block.Block
}

func newProducer(config *params.ChainConfig, engine consensus.Engine, bchain Backend, mux *event.TypeMux, maxBlkSize uint64) *producer {
	producer := &producer{
		config:             config,
		engine:             engine,
		bchain:             bchain,
		mux:                mux,
		chainDb:            bchain.ChainDb(),
		recv:               make(chan *Result, resultQueueSize),
		chain:              bchain.BlockChain(),
		proc:               bchain.BlockChain().Validator(),
		coinbase:           types.Address{},
		maxBlkSize:         maxBlkSize,
		agents:             make(map[Agent]struct{}),
		createRequestChan:  make(chan blockRequest, 128),
		createResponseChan: make(chan *block.Block, 128),
	}
	// Subscribe TxPreEvent for tx pool

	// Subscribe events for blockchain

	go producer.wait()

	//producer.commitNewWork()

	return producer
}

func (self *producer) setPrikey(pri *ecdsa.PrivateKey) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.priKey = pri
}

func (self *producer) setCoinbase(addr types.Address) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.coinbase = addr
}

func (self *producer) setExtra(extra []byte) {
	self.mu.Lock()
	defer self.mu.Unlock()
}

func (self *producer) pending() (*block.Block, *state.StateDB) {
	self.currentMu.Lock()
	defer self.currentMu.Unlock()

	if atomic.LoadInt32(&self.producing) == 0 {
		return block.NewBlock(
			self.current.header,
			self.current.txs,
			self.current.receipts,
		), self.current.state.Copy()
	}
	return self.current.Block, self.current.state.Copy()
}

func (self *producer) pendingBlock() *block.Block {
	self.currentMu.Lock()
	defer self.currentMu.Unlock()

	if atomic.LoadInt32(&self.producing) == 0 {
		return block.NewBlock(
			self.current.header,
			self.current.txs,
			self.current.receipts,
		)
	}
	return self.current.Block
}

func (self *producer) start() {
	self.mu.Lock()
	defer self.mu.Unlock()

	atomic.StoreInt32(&self.producing, 1)

	// spin up agents
	for agent := range self.agents {
		agent.Start()
	}
	go self.DealRequest()
}

func (self *producer) stop() {
	self.wg.Wait()

	self.mu.Lock()
	defer self.mu.Unlock()
	if atomic.LoadInt32(&self.producing) == 1 {
		for agent := range self.agents {
			agent.Stop()
		}
	}
	atomic.StoreInt32(&self.producing, 0)
	atomic.StoreInt32(&self.atWork, 0)
}

func (self *producer) register(agent Agent) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.agents[agent] = struct{}{}
	agent.SetReturnCh(self.recv)
}

func (self *producer) unregister(agent Agent) {
	self.mu.Lock()
	defer self.mu.Unlock()
	delete(self.agents, agent)
	agent.Stop()
}

func (this *producer) DealRequest() {
	for {
		select {
		case ask := <-this.createRequestChan:
			this.commitNewWork(ask.data, ask.timeLimit)
		}

	}
}

func (this *producer) ProduceNewBlock(data *block.ConsensusData, timeLimit int64) *block.Block {
	br := blockRequest{data: data, timeLimit: timeLimit}
	this.createRequestChan <- br
	timer := time.Tick(time.Duration(timeLimit+1) * time.Second)
	select {
	case <-timer:
		return nil
	case newBlock := <-this.createResponseChan:
		return newBlock
	}

	return nil
}

func (self *producer) wait() {
	for {
		//mustCommitNewWork := true

		for result := range self.recv {
			atomic.AddInt32(&self.atWork, -1)

			if result == nil {
				continue
			}
			block := result.Block
			work := result.Work

			// Update the block hash in all logs since it is now available and not when the
			// receipt/log of individual transactions were created.

			for _, r := range work.receipts {
				for _, l := range r.Logs {
					l.BlockHash = block.Hash()
				}
			}
			for _, log := range work.state.Logs() {
				log.BlockHash = block.Hash()
			}

			// Remove fail txs in tx pools
			self.bchain.TxPool().RemoveTxs(work.failTxHashs)

			self.createResponseChan <- block

		}
	}
}

// push sends a new work task to currently live blockproducer agents.
func (self *producer) push(work *Work) {
	if atomic.LoadInt32(&self.producing) != 1 {
		return
	}
	for agent := range self.agents {

		atomic.AddInt32(&self.atWork, 1)
		if ch := agent.Work(); ch != nil {
			ch <- work
		}
	}
}

// makeCurrent creates a new environment for the current cycle.
func (self *producer) makeCurrent(parent *block.Block, header *block.Header, timeLimit int64) error {
	state, err := self.chain.StateAt(parent.Root())
	if err != nil {
		return err
	}
	work := &Work{
		config:        self.config,
		signer:        transaction.NewMSigner(self.config.ChainId),
		state:         state,
		db:            self.chainDb,
		stateRootHash: parent.Root(),
		header:        header,
		createdAt:     time.Now(),
		txTimeLimit:   timeLimit,
		maxBlkSize:    common.StorageSize(self.maxBlkSize),
	}

	// Keep track of transactions which return errors so they can be removed
	work.tcount = 0
	self.current = work
	return nil
}

//now,one request , one commitNewWork
func (self *producer) commitNewWork(data *block.ConsensusData, timeLimit int64) {
	//time.Sleep(10*time.Second)
	self.mu.Lock()
	defer self.mu.Unlock()
	self.currentMu.Lock()
	defer self.currentMu.Unlock()

	tstart := time.Now()
	parent := self.chain.CurrentBlock()

	tstamp := tstart.Unix()
	if parent.Time().Cmp(new(big.Int).SetInt64(tstamp)) >= 0 {
		tstamp = parent.Time().Int64() + 1
	}

	// this will ensure we're not going off too far in the future
	if now := time.Now().Unix(); tstamp > now {
		wait := time.Duration(tstamp-now) * time.Second
		logger.Info("Producing too far in the future", "wait", common.PrettyDuration(wait))
		time.Sleep(wait)
	}

	num := parent.Number()

	header := &block.Header{
		ParentHash: parent.Hash(),
		Number:     &types.BigInt{*num.Add(num, common.Big1)},
		Time:       &types.BigInt{*big.NewInt(tstamp)},
		Cdata:      *data,
	}
	// Only set the coinbase if we are producing (avoid spurious block rewards)
	if atomic.LoadInt32(&self.producing) == 1 {
		header.Producer = self.coinbase
	}

	if err := self.engine.Prepare(self.chain, header); err != nil {
		logger.Error("Failed to prepare header for producing", "err", err)
		return
	}

	// Could potentially happen if starting to produce block in an odd state.
	err := self.makeCurrent(parent, header, timeLimit)
	if err != nil {
		logger.Error("Failed to create producing context", "err", err)
		return
	}
	// Create the current work task and check any fork transitions needed
	work := self.current

	//add test tx
	//self.addTestTransactions()
	//self.DoTestTransactionsQuery()
	//self.addTestTransactionsPledge()
	//self.addTestErrorTransactions()	//add a error transaction
	pending, err := self.bchain.TxPool().Pending()
	if err != nil {
		logger.Error("Failed to fetch pending transactions", "err", err)
		return
	}

	//txs := transaction.NewTransactionsForProducing(self.current.signer, pending)
	//actions := transaction.Actions{}
	//action := transaction.Action{types.Address{}, balancetransfer.MakeActionParamsReword(header.BlockProducer)}
	//actions = append(actions, action)

	//sysNonce := self.bchain.TxPool().State().GetNonce(params.Address)

	//tx := transaction.NewTransaction(sysNonce, actions)

	//txReword, err := transaction.SignTx(tx, self.current.signer, params.RewordPrikey)
	//if err != nil {
	//	logger.Error("Failed to make reword transaction", "err", err)
	//	return
	//}
	//txReword.Priority = big.NewInt(10)
	fmt.Println("-----currentBlock num:", self.bchain.BlockChain().CurrentBlockNum())

	incentiveTx, err := self.engine.Incentive(self.coinbase, work.state, header)
	if err != nil {
		logger.Error("Failed to fetch incentive transaction", "err", err)
		return
	}

	//txBadPublicPrivateTx := self.makeBadPublicPrivateKeyTx()
	//txsReward = append(txsReward , txBadPublicPrivateTx)

	txs := transaction.NewTransactionsByPriorityAndNonce(self.current.signer, pending, incentiveTx)
	fmt.Println("!!!!!!!!!!Current All Txs:", txs.Length())
	logger.Info(">>>>>Producer will commit transactions.......")
	work.commitTransactions(self.mux, txs, self.chain, self.coinbase)

	blkTx := work.txs
	if len(work.txs) > 0 {
		// remove incentive transaction in block body
		if work.txs[0] == incentiveTx {
			blkTx = work.txs[1:]
		}
	}

	// Create the new block to seal with the consensus engine
	if work.Block, err = self.engine.Finalize(self.chain, header, work.state, blkTx, work.receipts, true); err != nil {
		logger.Error("Failed to finalize block for sealing", "err", err)
		return
	}

	work.bchain = self.bchain
	self.push(work)
}

func (env *Work) commitTransactions(mux *event.TypeMux, txs *transaction.TransactionsByPriorityAndNonce, bc *blockchain.BlockChain, coinbase types.Address) {
	var coalescedLogs []*transaction.Log
	tmpDb, _ := database.OpenMemDB()
	blkCtx := actioncontext.NewBlockContext(env.state, env.db, tmpDb, &env.header.Number.IntVal, coinbase)
	logger.Info("CommitTransactions TimeLimit:", env.txTimeLimit)
	limitDuration := time.Duration(env.txTimeLimit) * time.Second
	start := time.Now()
	txsSize := common.StorageSize(0)
	hasExecOnce := bool(false)

	for {
		// Retrieve the next transaction and abort if all done
		tx := txs.Peek()
		if tx == nil {
			break
		}

		cust := time.Since(start)

		// time limit
		if cust+time.Duration(len(tx.Actions()))*time.Second >= limitDuration {
			logger.Info("end of commitTransactions by limit time duration", cust)
			break
		}

		// txs size limit
		txsSize += tx.Size()
		if hasExecOnce && txsSize > env.maxBlkSize {
			logger.Info("end of commitTransactions by limit txs size", txsSize)
			break
		}

		// Error may be ignored here. The error has already been checked
		// during transaction acceptance is the transaction pool.
		//
		// We use the eip155 signer regardless of the current hf.
		from, _ := transaction.Sender(env.signer, tx)
		// Check whether the tx is replay protected. If we're not in the EIP155 hf
		// phase, start ignoring the sender until we do.
		if false {
			if tx.Protected() {
				logger.Tracef("Ignoring reply protected transaction hash:%x\n", tx.Hash())

				txs.Pop()
				continue
			}
		}

		// Start executing the transaction
		env.state.Prepare(tx.Hash(), types.Hash{}, env.tcount)

		err, logs := env.commitTransaction(tx, blkCtx)
		switch err {
		case core.ErrNonceTooLow:
			// New head notification data race between the transaction pool and blockproducer, shift
			logger.Error("Skipping transaction with low nonce", "sender", from.HexLower(), "nonce", tx.Nonce(),"hash", tx.Hash().String())
			txs.Shift()
			env.failTxHashs = append(env.failTxHashs, tx.Hash())

		case core.ErrNonceTooHigh:
			// Reorg notification data race between the transaction pool and blockproducer, skip account =
			logger.Info("Skipping account with hight nonce", "sender", from.HexLower(), "nonce", tx.Nonce(), "hash", tx.Hash().String())
			txs.Pop()

		case nil:
			// Everything ok, collect the logs and shift in the next transaction from the same account
			coalescedLogs = append(coalescedLogs, logs...)
			env.tcount++
			txs.Shift()

		default:
			// Strange error, discard the transaction and get the next in line (note, the
			// nonce-too-high clause will prevent us from executing in vain).
			logger.Error("Transaction failed, account skipped", "hash", "sender", from.HexLower(), tx.Hash().String(), "err", err, "nonce", tx.Nonce())
			txs.Shift()
			env.failTxHashs = append(env.failTxHashs, tx.Hash())
		}
		hasExecOnce = true
	}

	if len(coalescedLogs) > 0 || env.tcount > 0 {
		// make a copy, the state caches the logs and these logs get "upgraded" from pending to produced block
		// logs by filling in the block hash when the block was produced block by the local blockproducer. This can
		// cause a race condition if a log was "upgraded" before the PendingLogsEvent is processed.
		cpy := make([]*transaction.Log, len(coalescedLogs))
		for i, l := range coalescedLogs {
			cpy[i] = new(transaction.Log)
			*cpy[i] = *l
		}
		go func(logs []*transaction.Log, tcount int) {
			if len(logs) > 0 {
				mux.Post(core.PendingLogsEvent{Logs: logs})
			}
			if tcount > 0 {
				mux.Post(core.PendingStateEvent{})
			}
		}(cpy, env.tcount)
	}
}

func (env *Work) commitTransaction(tx *transaction.Transaction, blkCtx *actioncontext.BlockContext) (error, []*transaction.Log) {
	snap := env.state.Snapshot()
	//                                ApplyTransaction(this.config,&coinbase,this.state ,header,tx)
	receipt, err := stateprocessor.ApplyTransaction(env.config, env.header, tx, blkCtx)
	if err != nil {
		env.state.RevertToSnapshot(snap)
		return err, nil
	}
	env.txs = append(env.txs, tx)
	env.receipts = append(env.receipts, receipt)

	return nil, receipt.Logs
}
