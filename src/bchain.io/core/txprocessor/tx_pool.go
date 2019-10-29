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
// @File: tx_pool.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package txprocessor

import (
	"errors"
	"fmt"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
	"math"
	"math/big"
	"bchain.io/common/types"
	"bchain.io/core"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
	"bchain.io/params"
	"bchain.io/utils/event"
	"bchain.io/utils/metrics"
	"sort"
	"sync"
	"time"
)

const (
	// chainHeadChanSize is the size of channel listening to ChainHeadEvent.
	chainHeadChanSize = 10
	// rmTxChanSize is the size of channel listening to RemovedTransactionEvent.
	rmTxChanSize = 10

	//todo:test priority
	testPriorityValue = 10
)

var (
	// ErrInvalidSender is returned if the transaction contains an invalid signature.
	ErrInvalidSender = errors.New("invalid sender")

	// ErrNonceTooLow is returned if the nonce of a transaction is lower than the
	// one present in the local chain.
	ErrNonceTooLow = errors.New("nonce too low")

	//ErrWrongTransactionAmount is returned if a transaction's amount is nil
	ErrWrongTransactionAmount = errors.New("transaction's Amount is wrong")

	//ErrUnderpriority is returned if a transaction's priority is below the minimum
	//configured for the transaction pool
	ErrUnderPriority = errors.New("transaction underpriority")

	ErrReplaceUnderpriority = errors.New("replacement transaction underpriority")

	// ErrInsufficientFunds is returned if the total cost of executing a transaction
	// is higher than the balance of the user's account.
	ErrInsufficientFunds = errors.New("insufficient funds for value")

	// ErrNegativeValue is a sanity error to ensure noone is able to specify a
	// transaction with a negative value.
	ErrNegativeValue = errors.New("negative value")

	// ErrOversizedData is returned if the input data of a transaction is greater
	// than some meaningful limit a user might use. This is not a consensus error
	// making the transaction invalid, rather a DOS protection.
	ErrOversizedData = errors.New("oversized data")
)

var (
	evictionInterval    = time.Minute     // Time interval to check for evictable transactions
	statsReportInterval = 8 * time.Second // Time interval to report transaction pool stats
)

var (
	// Metrics for the pending pool

	pendingDiscardCounter   = metrics.NewRegisteredCounter("txpool/pending/discard", nil)
	pendingReplaceCounter   = metrics.NewRegisteredCounter("txpool/pending/replace", nil)
	pendingRateLimitCounter = metrics.NewRegisteredCounter("txpool/pending/ratelimit", nil) // Dropped due to rate limiting
	pendingNofundsCounter   = metrics.NewRegisteredCounter("txpool/pending/nofunds", nil)   // Dropped due to out-of-funds

	// Metrics for the queued pool
	queuedDiscardCounter   = metrics.NewRegisteredCounter("txpool/queued/discard", nil)
	queuedReplaceCounter   = metrics.NewRegisteredCounter("txpool/queued/replace", nil)
	queuedRateLimitCounter = metrics.NewRegisteredCounter("txpool/queued/ratelimit", nil) // Dropped due to rate limiting
	queuedNofundsCounter   = metrics.NewRegisteredCounter("txpool/queued/nofunds", nil)   // Dropped due to out-of-funds

	// General tx metrics
	invalidTxCounter = metrics.NewRegisteredCounter("txpool/invalid", nil)
)

// TxStatus is the current status of a transaction as seen by the pool.
type TxStatus uint

const (
	TxStatusUnknown TxStatus = iota
	TxStatusQueued
	TxStatusPending
	TxStatusIncluded
)

// blockChain provides the state of blockchain  to do
// some pre checks in tx pool and event subscribers.
type blockChain interface {
	CurrentBlock() *block.Block
	GetBlock(hash types.Hash, number uint64) *block.Block
	StateAt(root types.Hash) (*state.StateDB, error)

	SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription
}

// TxPoolConfig are the configuration parameters of the transaction pool.
type TxPoolConfig struct {
	NoLocals  bool          // Whether local transaction handling should be disabled
	Journal   string        // Journal of local transactions to survive node restarts
	Rejournal time.Duration // Time interval to regenerate the local transaction journal

	AccountSlots uint64 // Minimum number of executable transaction slots guaranteed per account
	GlobalSlots  uint64 // Maximum number of executable transaction slots for all accounts
	AccountQueue uint64 // Maximum number of non-executable transaction slots permitted per account
	GlobalQueue  uint64 // Maximum number of non-executable transaction slots for all accounts
	Priority 	 uint64
	Lifetime time.Duration // Maximum amount of time non-executable transaction are queued
}


// DefaultTxPoolConfig contains the default configurations for the transaction
// pool.
var DefaultTxPoolConfig = TxPoolConfig{
	Journal:   "transactions.msgp",
	Rejournal: time.Hour,

	NoLocals : true,
	AccountSlots: 16,
	GlobalSlots:  65536,
	AccountQueue: 640,
	GlobalQueue:  10240,
	Priority: 	0,
	Lifetime: 30 * time.Minute,
}

// sanitize checks the provided user configurations and changes anything that's
// unreasonable or unworkable.
func (config *TxPoolConfig) sanitize() TxPoolConfig {
	conf := *config
	if conf.Rejournal < time.Second {
		logger.Warn("Sanitizing invalid txpool journal time", "provided", conf.Rejournal, "updated", time.Second)
		conf.Rejournal = time.Second
	}

	return conf
}

// TxPool contains all currently known transactions. Transactions
// enter the pool when they are received from the network or submitted
// locally. They exit the pool when they are included in the blockchain.
//
// The pool separates processable transactions (which can be applied to the
// current state) and future transactions. Transactions move between those
// two states over time as they are received and processed.
type TxPool struct {
	config       TxPoolConfig
	chainconfig  *params.ChainConfig
	chain        blockChain
	txFeed       event.Feed
	scope        event.SubscriptionScope
	chainHeadCh  chan core.ChainHeadEvent
	chainHeadSub event.Subscription
	signer       transaction.Signer
	mu           sync.RWMutex

	currentState *state.StateDB      // Current state in the blockchain head
	pendingState *state.ManagedState // Pending state tracking virtual nonces

	priorityThreshold *big.Int		//txpool priority threshold
	inter    Interpreter

	locals  *accountSet // Set of local transaction to exempt from eviction rules
	journal *txJournal  // Journal of local transaction to back up to disk

	pending   map[types.Address]*txList               // All currently processable transactions
	queue     map[types.Address]*txList               // Queued but non-processable transactions
	beats     map[types.Address]time.Time             // Last heartbeat from each known account
	all       map[types.Hash]*transaction.Transaction // All transactions to allow lookups
	priorited *txPriorityList

	wg sync.WaitGroup // for shutdown sync

}

// NewTxPool creates a new transaction pool to gather, sort and filter inbound
// transactions from the network.
func NewTxPool(config TxPoolConfig, chainconfig *params.ChainConfig, chain blockChain) *TxPool {

	config = (&config).sanitize()
	// Create the transaction pool with its initial settings
	pool := &TxPool{
		config:      config,
		chainconfig: chainconfig,
		chain:       chain,
		signer:      transaction.NewMSigner(chainconfig.ChainId),
		pending:     make(map[types.Address]*txList),
		queue:       make(map[types.Address]*txList),
		beats:       make(map[types.Address]time.Time),
		all:         make(map[types.Hash]*transaction.Transaction),
		chainHeadCh: make(chan core.ChainHeadEvent, chainHeadChanSize),
		priorityThreshold: 	big.NewInt(0),		//
	}
	//set test interpreter for test
	pool.inter = new(testInterpreter)
	pool.locals = newAccountSet(pool.signer)
	pool.priorited = newTxPriorityList(&pool.all)
	pool.reset(nil, chain.CurrentBlock().Header())

	// If local transactions and journaling is enabled, load from disk
	if !config.NoLocals && config.Journal != "" {
		pool.journal = newTxJournal(config.Journal)

		if err := pool.journal.load(pool.AddLocal); err != nil {
			logger.Warn("Failed to load transaction journal", "err", err)
		}
		if err := pool.journal.rotate(pool.local()); err != nil {
			logger.Warn("Failed to rotate transaction journal", "err", err)
		}
	}
	// Subscribe events from blockchain
	pool.chainHeadSub = pool.chain.SubscribeChainHeadEvent(pool.chainHeadCh)

	// Start the event loop and return
	pool.wg.Add(1)

	logger.Info("config.Priority:" , config.Priority)
	//just for test
	pool.SetPriorityThreshold(big.NewInt(int64(config.Priority)))

	go pool.loop()

	return pool
}

// loop is the transaction pool's main event loop, waiting for and reacting to
// outside blockchain events as well as for various reporting and transaction
// eviction events.
func (pool *TxPool) loop() {
	defer pool.wg.Done()

	// Start the stats reporting and transaction eviction tickers
	var prevPending, prevQueued int

	report := time.NewTicker(statsReportInterval)
	defer report.Stop()

	evict := time.NewTicker(evictionInterval)
	defer evict.Stop()

	journal := time.NewTicker(pool.config.Rejournal)
	defer journal.Stop()

	// Track the previous head headers for transaction reorgs
	head := pool.chain.CurrentBlock()

	// Keep waiting for and reacting to the various events
	for {
		select {
		// Handle ChainHeadEvent
		case ev := <-pool.chainHeadCh:
			if ev.Block != nil {
				pool.mu.Lock()
				pool.reset(head.Header(), ev.Block.Header())
				head = ev.Block

				pool.mu.Unlock()
			}
		// Be unsubscribed due to system stopped
		case <-pool.chainHeadSub.Err():
			return

		// Handle stats reporting ticks
		case <-report.C:
			pool.mu.RLock()
			pending, queued := pool.stats()
			stales := pool.priorited.stales
			pool.mu.RUnlock()
			logger.Info("report Pending:" , pending , "   report Queued:" , queued)
			if pending != prevPending || queued != prevQueued {
				logger.Debug("Transaction pool status report", "executable", pending, "queued", queued, "  stales:", stales)
				prevPending, prevQueued = pending, queued
			}

		// Handle inactive account transaction eviction
		case <-evict.C:
			pool.mu.Lock()
			for addr := range pool.queue {
				// Any old enough should be removed
				if time.Since(pool.beats[addr]) > pool.config.Lifetime {
					for _, tx := range pool.queue[addr].Flatten() {
						logger.Warn("delete from queue pool addr", addr.HexLower(), "txHash", tx.Hash().String())
						pool.removeTx(tx.Hash())
					}
				}
			}
			for addr := range pool.pending {
				// Any old enough should be removed
				if time.Since(pool.beats[addr]) > pool.config.Lifetime {
					for _, tx := range pool.pending[addr].Flatten() {
						logger.Warn("delete from pending pool addr", addr.HexLower(), "txHash", tx.Hash().String())
						pool.removeTx(tx.Hash())
					}
				}
			}
			pool.mu.Unlock()

		// Handle local transaction journal rotation
		case <-journal.C:
			if pool.journal != nil {
				pool.mu.Lock()
				if err := pool.journal.rotate(pool.local()); err != nil {
					logger.Warn("Failed to rotate local tx journal", "err", err)
				}
				pool.mu.Unlock()
			}
		}
	}
}

// lockedReset is a wrapper around reset to allow calling it in a thread safe
// manner. This method is only ever used in the tester!
func (pool *TxPool) lockedReset(oldHead, newHead *block.Header) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pool.reset(oldHead, newHead)
}

// reset retrieves the current state of the blockchain and ensures the content
// of the transaction pool is valid with regard to the chain state.
func (pool *TxPool) reset(oldHead, newHead *block.Header) {
	// If we're reorging an old state, reinject all dropped transactions
	var reinject transaction.Transactions

	if oldHead != nil && oldHead.Hash() != newHead.ParentHash {
		// If the reorg is too deep, avoid doing it (will happen during fast sync)
		oldNum := oldHead.Number.IntVal.Uint64()
		newNum := newHead.Number.IntVal.Uint64()

		if depth := uint64(math.Abs(float64(oldNum) - float64(newNum))); depth > 64 {
			logger.Debug("Skipping deep transaction reorg", "depth", depth)
		} else {
			// Reorg seems shallow enough to pull in all transactions into memory
			var discarded, included transaction.Transactions

			var (
				rem = pool.chain.GetBlock(oldHead.Hash(), oldHead.Number.IntVal.Uint64())
				add = pool.chain.GetBlock(newHead.Hash(), newHead.Number.IntVal.Uint64())
			)
			for rem.NumberU64() > add.NumberU64() {
				discarded = append(discarded, rem.Transactions()...)
				if rem = pool.chain.GetBlock(rem.ParentHash(), rem.NumberU64()-1); rem == nil {
					logger.Error("Unrooted old chain seen by tx pool", "block", oldHead.Number, "hash", oldHead.Hash())
					return
				}
			}
			for add.NumberU64() > rem.NumberU64() {
				included = append(included, add.Transactions()...)
				if add = pool.chain.GetBlock(add.ParentHash(), add.NumberU64()-1); add == nil {
					logger.Error("Unrooted new chain seen by tx pool", "block", newHead.Number, "hash", newHead.Hash())
					return
				}
			}
			for rem.Hash() != add.Hash() {
				discarded = append(discarded, rem.Transactions()...)
				if rem = pool.chain.GetBlock(rem.ParentHash(), rem.NumberU64()-1); rem == nil {
					logger.Error("Unrooted old chain seen by tx pool", "block", oldHead.Number, "hash", oldHead.Hash())
					return
				}
				included = append(included, add.Transactions()...)
				if add = pool.chain.GetBlock(add.ParentHash(), add.NumberU64()-1); add == nil {
					logger.Error("Unrooted new chain seen by tx pool", "block", newHead.Number, "hash", newHead.Hash())
					return
				}
			}
			reinject = transaction.TxDifference(discarded, included)
		}
	}
	// Initialize the internal state to the current head
	if newHead == nil {
		newHead = pool.chain.CurrentBlock().Header() // Special case during testing
	}
	statedb, err := pool.chain.StateAt(newHead.StateRootHash)
	if err != nil {
		logger.Error("Failed to reset txpool state", "err", err)
		return
	}
	pool.currentState = statedb
	pool.pendingState = state.ManageState(statedb)

	// Inject any transactions discarded due to reorgs
	logger.Debug("Reinjecting stale transactions", "count", len(reinject))
	pool.addTxsLocked(reinject, false)

	// validate the pool of pending transactions, this will remove
	// any transactions that have been included in the block or
	// have been invalidated because of another transaction
	pool.demoteUnexecutables()

	// Update all accounts to the latest known pending nonce
	for addr, list := range pool.pending {
		txs := list.Flatten() // Heavy but will be cached and is needed by the blockproducer anyway
		pool.pendingState.SetNonce(addr, txs[len(txs)-1].Nonce()+1)
	}
	// Check the queue and move transactions over to the pending if possible
	// or remove those that have become invalid
	pool.promoteExecutables(nil)
}

// Stop terminates the transaction pool.
func (pool *TxPool) Stop() {
	// Unsubscribe all subscriptions registered from txpool
	pool.scope.Close()

	// Unsubscribe subscriptions registered from blockchain
	pool.chainHeadSub.Unsubscribe()
	pool.wg.Wait()

	if pool.journal != nil {
		pool.journal.close()
	}
	logger.Info("Transaction pool stopped")
}

// SubscribeTxPreEvent registers a subscription of TxPreEvent and
// starts sending event to the given channel.
func (pool *TxPool) SubscribeTxPreEvent(ch chan<- core.TxPreEvent) event.Subscription {
	return pool.scope.Track(pool.txFeed.Subscribe(ch))
}

//setPriority updates the minimum priority required by the transaction pool for a new transaction,
//and drops all transactions below this threshold
func (pool *TxPool) PriorityThreshold() *big.Int {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return new(big.Int).Set(pool.priorityThreshold)
}

func (pool *TxPool) SetPriorityThreshold(priority *big.Int) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pool.priorityThreshold = priority

	for _, tx := range pool.priorited.Cap(priority, pool.locals) {
		pool.removeTx(tx.Hash())
	}
	logger.Info("Transaction pool priority threshold updated", "priority:", priority.Int64())
}


// State returns the virtual managed state of the transaction pool.
func (pool *TxPool) State() *state.ManagedState {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.pendingState
}

// Stats retrieves the current pool stats, namely the number of pending and the
// number of queued (non-executable) transactions.
func (pool *TxPool) Stats() (int, int) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.stats()
}

// stats retrieves the current pool stats, namely the number of pending and the
// number of queued (non-executable) transactions.
func (pool *TxPool) stats() (int, int) {
	pending := 0
	for _, list := range pool.pending {
		pending += list.Len()
	}
	queued := 0
	for _, list := range pool.queue {
		queued += list.Len()
	}
	return pending, queued
}

// Content retrieves the data content of the transaction pool, returning all the
// pending as well as queued transactions, grouped by account and sorted by nonce.
func (pool *TxPool) Content() (map[types.Address]transaction.Transactions, map[types.Address]transaction.Transactions) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pending := make(map[types.Address]transaction.Transactions)
	for addr, list := range pool.pending {
		pending[addr] = list.Flatten()
	}
	queued := make(map[types.Address]transaction.Transactions)
	for addr, list := range pool.queue {
		queued[addr] = list.Flatten()
	}
	return pending, queued
}

// Pending retrieves all currently processable transactions, groupped by origin
// account and sorted by nonce. The returned transaction set is a copy and can be
// freely modified by calling code.
func (pool *TxPool) Pending() (map[types.Address]transaction.Transactions, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pending := make(map[types.Address]transaction.Transactions)
	for addr, list := range pool.pending {
		pending[addr] = list.Flatten()
	}
	return pending, nil
}

// local retrieves all currently known local transactions, groupped by origin
// account and sorted by nonce. The returned transaction set is a copy and can be
// freely modified by calling code.
func (pool *TxPool) local() map[types.Address]transaction.Transactions {
	txs := make(map[types.Address]transaction.Transactions)
	for addr := range pool.locals.accounts {
		if pending := pool.pending[addr]; pending != nil {
			txs[addr] = append(txs[addr], pending.Flatten()...)
		}
		if queued := pool.queue[addr]; queued != nil {
			txs[addr] = append(txs[addr], queued.Flatten()...)
		}
	}
	return txs
}

//validateTx return nil,we should check the tx validation at interpreter
func (pool *TxPool) validateTx(tx *transaction.Transaction) error {
	// Heuristic limit, reject transactions over 32KB to prevent DOS attacks
	if tx.Size() > 32*1024 {
		return ErrOversizedData
	}

	// Make sure the transaction is signed properly
	from, err := transaction.Sender(pool.signer, tx)
	if err != nil {
		logger.Error("Why invalidSender :",err)
		return ErrInvalidSender
	}

	// Ensure the transaction adheres to nonce ordering
	if pool.currentState.GetNonce(from) > tx.Nonce() {
		logger.Errorf("Account :%x , stateNonce:%d   tx.Nonce:%d" , from , pool.currentState.GetNonce(from) , tx.Nonce())
		return ErrNonceTooLow
	}

	//always check the priority
	if pool.priorited.Underpriority(tx, pool.locals, pool.priorityThreshold) {
		logger.Info("Discard by Priority:   tx.Priority:" , tx.Priority().Uint64() , "  But Miner need more than:" , pool.priorityThreshold.Uint64())
		return ErrUnderPriority
	}

	return nil
}

// add validates a transaction and inserts it into the non-executable queue for
// later pending promotion and execution. If the transaction is a replacement for
// an already pending or queued one, it overwrites the previous and returns this
// so outer code doesn't uselessly call promote.
//
// If a newly added transaction is marked as local, its sending account will be
// whitelisted
func (pool *TxPool) add(tx *transaction.Transaction, local bool) (bool, error) {
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	//get Priority
	//todo:we should get priority by tx.actions,maybe the first action is the value of priofity
	//todo:like : interpreter.getPriorityByActions(xxxx)
	//todo:here we set a const value
	tx.ParsePriority()
	//tx.SetPriority(big.NewInt(int64(testPriorityValue)))


	// If the transaction is already known, discard it
	hash := tx.Hash()
	if pool.all[hash] != nil {
		//logger.Tracef("Discarding already known transaction hash:0x%x",  hash)
		return false, fmt.Errorf("known transaction")
	}
	// If the transaction fails basic validation, discard it


	if err := pool.validateTx(tx); err != nil {
		logger.Trace("Discarding invalid transaction hash:0x%x , err:%s",  hash, err.Error())
		invalidTxCounter.Inc(1)
		return false, err
	}


	if uint64(len(pool.all)) >= pool.config.GlobalSlots+pool.config.GlobalQueue {
		//do not add more transactions
		//if pool.priorited.Underpriority(tx, pool.locals, pool.priorityThreshold) {
		//	fmt.Println("add 3....return")
		//	return false, ErrUnderPriority
		//}

		//New transaction is better than our worse ones , make room for it
		drop := pool.priorited.Discard(len(pool.all)-int(pool.config.GlobalSlots+pool.config.GlobalQueue-1), pool.locals)
		for _, tx := range drop {
			pool.removeTx(tx.Hash())
		}

		return false, fmt.Errorf("pool.all more than config.GlobalQueue")
	}

	// If the transaction is replacing an already pending one, do directly
	from, _ := transaction.Sender(pool.signer, tx) // already validated
	if list := pool.pending[from]; list != nil && list.Overlaps(tx) {

		inserted, old := list.Add(tx, 0)
		if !inserted {
			//have same nonce ,but priority is lower than before
			return false, ErrReplaceUnderpriority
		}

		// if old != nil,the tx has been here before
		if old != nil {
			//delete the old transaction
			delete(pool.all, old.Hash())
			pool.priorited.Removed()
		} else {
			//old == nil,mean here is no the transaction in the pool before
			pool.all[tx.Hash()] = tx
			pool.journalTx(from, tx)

			logger.Trace("Pooled new executable transaction hash:0x%x , from:0x%x", hash, from)

			// We've directly injected a replacement transaction, notify subsystems
			logger.Debugf("!!!!!!!!!!!add  From:%x  Nonce:%d", from, tx.Data.H.Nonce)
			go pool.txFeed.Send(core.TxPreEvent{tx})

		}

		return true, nil
	}
	// New transaction isn't replacing a pending one, push into queue
	replace, err := pool.enqueueTx(hash, tx)
	if err != nil {
		return false, err
	}
	// Mark local addresses and journal local transactions
	if local {
		pool.locals.add(from)
	}
	pool.journalTx(from, tx)

	logger.Tracef("Pooled new future transaction hash:0x%x  from:0x%x", hash, from)
	if replace {

		logger.Debug("add-replace a old one")
	} else {
		//fmt.Println("add--a new one")
	}

	return replace, nil
}

// enqueueTx inserts a new transaction into the non-executable transaction queue.
//
// Note, this method assumes the pool lock is held!
//result bool,meaning replace a old one ,result error,meaning insert right or not
func (pool *TxPool) enqueueTx(hash types.Hash, tx *transaction.Transaction) (bool, error) {
	// Try to insert the transaction into the future queue
	from, _ := transaction.Sender(pool.signer, tx) // already validated
	if pool.queue[from] == nil {
		pool.queue[from] = newTxList(false)
	}
	inserted, old := pool.queue[from].Add(tx, 0)
	if !inserted {
		//new transaction's priority is lower than the old one
		return false, ErrReplaceUnderpriority
	}

	if old != nil {
		//old != nil,it's mean we insert new transaction in the queue,
		//we should delete the old one

		delete(pool.all, old.Hash())
		pool.priorited.Removed()
	}
	//notice , if no the same tx before ,we should not return a true boolean,
	//should false,because not replace
	//old == nil,no same tx before

	pool.all[hash] = tx
	pool.priorited.Put(tx)
	return old != nil, nil
}

// journalTx adds the specified transaction to the local disk journal if it is
// deemed to have been sent from a local account.
func (pool *TxPool) journalTx(from types.Address, tx *transaction.Transaction) {
	// Only journal if it's enabled and the transaction is local
	if pool.journal == nil || !pool.locals.contains(from) {
		return
	}
	if err := pool.journal.insert(tx); err != nil {
		logger.Warn("Failed to journal local transaction", "err", err)
	}
}

// promoteTx adds a transaction to the pending (processable) list of transactions.
//
// Note, this method assumes the pool lock is held!
func (pool *TxPool) promoteTx(addr types.Address, hash types.Hash, tx *transaction.Transaction) {
	// Try to insert the transaction into the pending queue
	if pool.pending[addr] == nil {
		pool.pending[addr] = newTxList(true)
	}
	list := pool.pending[addr]

	inserted, old := list.Add(tx, 0)
	if !inserted {
		// An older transaction was better, discard this
		delete(pool.all, hash)
		pool.priorited.Removed()
		pendingDiscardCounter.Inc(1)
		return
	}
	// Otherwise discard any previous transaction and mark this
	if old != nil {
		delete(pool.all, old.Hash())
		pool.priorited.Removed()
	}
	// Failsafe to work around direct pending inserts (tests)
	if pool.all[hash] == nil {
		pool.all[hash] = tx
		pool.priorited.Put(tx)
	}
	// Set the potentially new pending nonce and notify any subsystems of the new tx
	pool.beats[addr] = time.Now()
	pool.pendingState.SetNonce(addr, tx.Nonce()+1)
	logger.Debugf("!!!!!!!!!!!promoteTx From:%x  Nonce:%d", addr, tx.Nonce())
	go pool.txFeed.Send(core.TxPreEvent{tx})
}

// AddLocal enqueues a single transaction into the pool if it is valid, marking
// the sender as a local one in the mean time
func (pool *TxPool) AddLocal(tx *transaction.Transaction) error {
	return pool.addTx(tx, !pool.config.NoLocals)
}

// AddRemote enqueues a single transaction into the pool if it is valid.
func (pool *TxPool) AddRemote(tx *transaction.Transaction) error {
	return pool.addTx(tx, false)
}

// AddLocals enqueues a batch of transactions into the pool if they are valid
func (pool *TxPool) AddLocals(txs []*transaction.Transaction) []error {
	return pool.addTxs(txs, !pool.config.NoLocals)
}

// AddRemotes enqueues a batch of transactions into the pool if they are valid.
func (pool *TxPool) AddRemotes(txs []*transaction.Transaction) []error {
	return pool.addTxs(txs, false)
}

// addTx enqueues a single transaction into the pool if it is valid.
func (pool *TxPool) addTx(tx *transaction.Transaction, local bool) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Try to inject the transaction and update any state
	replace, err := pool.add(tx, local)
	if err != nil {
		return err
	}
	// If we added a new transaction, run promotion checks and return
	if !replace {
		from, _ := transaction.Sender(pool.signer, tx) // already validated
		pool.promoteExecutables([]types.Address{from})
	}
	return nil
}

// addTxs attempts to queue a batch of transactions if they are valid.
func (pool *TxPool) addTxs(txs []*transaction.Transaction, local bool) []error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.addTxsLocked(txs, local)
}

// addTxsLocked attempts to queue a batch of transactions if they are valid,
// whilst assuming the transaction pool lock is already held.
func (pool *TxPool) addTxsLocked(txs []*transaction.Transaction, local bool) []error {
	// Add the batch of transaction, tracking the accepted ones
	dirty := make(map[types.Address]struct{})
	errs := make([]error, len(txs))

	for i, tx := range txs {

		var replace bool
		if replace, errs[i] = pool.add(tx, local); errs[i] == nil {

			if !replace {
				from, _ := transaction.Sender(pool.signer, tx) // already validated
				dirty[from] = struct{}{}
			}
		} else {
			if errs[i].Error() != "known transaction" {
				logger.Errorf("errs[%d]=%v", i, errs[i])
			}
		}
	}

	// Only reprocess the internal state if something was actually added
	if len(dirty) > 0 {
		addrs := make([]types.Address, 0, len(dirty))
		for addr := range dirty {
			addrs = append(addrs, addr)
		}
		_, queLen := pool.stats()

		logger.Debugf("[addTxsLocked] Get:%d", queLen)
		pool.promoteExecutables(addrs)
	}

	return errs
}

// Status returns the status (unknown/pending/queued) of a batch of transactions
// identified by their hashes.
func (pool *TxPool) Status(hashes []types.Hash) []TxStatus {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	status := make([]TxStatus, len(hashes))
	for i, hash := range hashes {
		if tx := pool.all[hash]; tx != nil {
			from, _ := transaction.Sender(pool.signer, tx) // already validated
			if pool.pending[from] != nil && pool.pending[from].txs.items[tx.Nonce()] != nil {
				status[i] = TxStatusPending
			} else {
				status[i] = TxStatusQueued
			}
		}
	}
	return status
}

// Get returns a transaction if it is contained in the pool
// and nil otherwise.
func (pool *TxPool) Get(hash types.Hash) *transaction.Transaction {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.all[hash]
}

func (pool *TxPool) RemoveTxs(hashs []types.Hash) {
	if len(hashs) == 0 {
		return
	}
	go func() {
		pool.mu.Lock()
		defer pool.mu.Unlock()
		for _, txhash := range hashs {
			pool.removeTx(txhash)
		}
	}()
}

// removeTx removes a single transaction from the queue, moving all subsequent
// transactions back to the future queue.
func (pool *TxPool) removeTx(hash types.Hash) {
	// Fetch the transaction we wish to delete
	tx, ok := pool.all[hash]
	if !ok {
		return
	}
	addr, _ := transaction.Sender(pool.signer, tx) // already validated during insertion

	// Remove it from the list of known transactions
	delete(pool.all, hash)
	pool.priorited.Removed()
	// Remove the transaction from the pending lists and reset the account nonce
	if pending := pool.pending[addr]; pending != nil {
		if removed, invalids := pending.Remove(tx); removed {
			// If no more transactions are left, remove the list
			if pending.Empty() {
				delete(pool.pending, addr)
				delete(pool.beats, addr)
			} else {
				// Otherwise postpone any invalidated transactions
				for _, tx := range invalids {
					pool.enqueueTx(tx.Hash(), tx)
				}
			}
			// Update the account nonce if needed
			if nonce := tx.Nonce(); pool.pendingState.GetNonce(addr) > nonce {
				pool.pendingState.SetNonce(addr, nonce)
			}
			return
		}
	}
	// Transaction is in the future queue
	if future := pool.queue[addr]; future != nil {
		future.Remove(tx)
		if future.Empty() {
			delete(pool.queue, addr)
		}
	}
}

// promoteExecutables moves transactions that have become processable from the
// future queue to the set of pending transactions. During this process, all
// invalidated transactions (low nonce, low balance) are deleted.
func (pool *TxPool) promoteExecutables(accounts []types.Address) {
	// Gather all the accounts potentially needing updates
	if accounts == nil {
		accounts = make([]types.Address, 0, len(pool.queue))
		for addr := range pool.queue {
			accounts = append(accounts, addr)
		}
	}
	// Iterate over all accounts and promote any executable transactions
	for _, addr := range accounts {
		list := pool.queue[addr]
		if list == nil {
			continue // Just in case someone calls with a non existing account
		}
		//fmt.Println("[promoteExecutables]List Len Before:Forward:" , len(list.txs.items))
		// Drop all transactions that are deemed too old (low nonce)
		for _, tx := range list.Forward(pool.currentState.GetNonce(addr)) {
			hash := tx.Hash()

			logger.Tracef("Removed old queued transaction hash:0x%x", hash)
			delete(pool.all, hash)
			pool.priorited.Removed()
		}
		//fmt.Println("[promoteExecutables]List Len Before:Filter:" , len(list.txs.items))
		// Drop all transactions that are too costly (low balance )

		//drops, _ := list.Filter(pool.currentState.GetBalance(addr), 0)
		//for _, tx := range drops {
		//	hash := tx.Hash()
		//	logger.Tracef("Removed unpayable queued transaction hash:0x%x", hash)
		//	delete(pool.all, hash)
		//	pool.priorited.Removed()
		//	queuedNofundsCounter.Inc(1)
		//}

		//fmt.Println("[promoteExecutables]List Len Before:Ready:" , len(list.txs.items))
		// Gather all executable transactions and promote them
		for _, tx := range list.Ready(pool.pendingState.GetNonce(addr)) {
			hash := tx.Hash()
			logger.Trace("Promoting queued transaction hash:", hash.String())

			pool.promoteTx(addr, hash, tx)
		}
		// Drop all transactions over the allowed limit
		//fmt.Println("[promoteExecutables]List Len Before:Cap:" , len(list.txs.items))
		if !pool.locals.contains(addr) {
			for _, tx := range list.Cap(int(pool.config.AccountQueue)) {
				hash := tx.Hash()
				delete(pool.all, hash)
				pool.priorited.Removed()
				queuedRateLimitCounter.Inc(1)
				logger.Tracef("Removed cap-exceeding queued transaction hash:0x%x", hash)
			}
		}
		// Delete the entire queue entry if it became empty.
		if list.Empty() {
			delete(pool.queue, addr)
		}
		//fmt.Println("[promoteExecutables]List Len AT Range End:" , len(list.txs.items))
	}

	// If the pending limit is overflown, start equalizing allowances
	pending := uint64(0)
	for _, list := range pool.pending {
		pending += uint64(list.Len())
	}
	//fmt.Println("[promoteExecutables]Pending+++ :" , pending)
	if pending > pool.config.GlobalSlots {
		pendingBeforeCap := pending
		// Assemble a spam order to penalize large transactors first
		spammers := prque.New()
		for addr, list := range pool.pending {
			// Only evict transactions from high rollers
			if !pool.locals.contains(addr) && uint64(list.Len()) > pool.config.AccountSlots {
				spammers.Push(addr, float32(list.Len()))
			}
		}
		// Gradually drop transactions from offenders
		offenders := []types.Address{}
		for pending > pool.config.GlobalSlots && !spammers.Empty() {
			// Retrieve the next offender if not local address
			offender, _ := spammers.Pop()
			offenders = append(offenders, offender.(types.Address))

			// Equalize balances until all the same or below threshold
			if len(offenders) > 1 {
				// Calculate the equalization threshold for all current offenders
				threshold := pool.pending[offender.(types.Address)].Len()

				// Iteratively reduce all offenders until below limit or threshold reached
				for pending > pool.config.GlobalSlots && pool.pending[offenders[len(offenders)-2]].Len() > threshold {
					for i := 0; i < len(offenders)-1; i++ {
						list := pool.pending[offenders[i]]
						for _, tx := range list.Cap(list.Len() - 1) {
							// Drop the transaction from the global pools too
							hash := tx.Hash()
							delete(pool.all, hash)
							pool.priorited.Removed()
							// Update the account nonce to the dropped transaction
							if nonce := tx.Nonce(); pool.pendingState.GetNonce(offenders[i]) > nonce {
								pool.pendingState.SetNonce(offenders[i], nonce)
							}
							logger.Trace("Removed fairness-exceeding pending transaction hash:0x%x", hash)
						}
						pending--
					}
				}
			}
		}
		// If still above threshold, reduce to limit or min allowance
		if pending > pool.config.GlobalSlots && len(offenders) > 0 {
			for pending > pool.config.GlobalSlots && uint64(pool.pending[offenders[len(offenders)-1]].Len()) > pool.config.AccountSlots {
				for _, addr := range offenders {
					list := pool.pending[addr]
					for _, tx := range list.Cap(list.Len() - 1) {
						// Drop the transaction from the global pools too
						hash := tx.Hash()
						delete(pool.all, hash)
						pool.priorited.Removed()
						// Update the account nonce to the dropped transaction
						if nonce := tx.Nonce(); pool.pendingState.GetNonce(addr) > nonce {
							pool.pendingState.SetNonce(addr, nonce)
						}
						logger.Trace("Removed fairness-exceeding pending transaction hash:0x%x", hash)
					}
					pending--
				}
			}
		}
		pendingRateLimitCounter.Inc(int64(pendingBeforeCap - pending))
	}
	// If we've queued more transactions than the hard limit, drop oldest ones

	queued := uint64(0)
	for _, list := range pool.queue {
		queued += uint64(list.Len())
	}

	//fmt.Println("[promoteExecutables]Queued+++:" , queued , "  GlobalQueue:" , pool.config.GlobalQueue)
	if queued > pool.config.GlobalQueue {
		// Sort all accounts with queued transactions by heartbeat
		addresses := make(addresssByHeartbeat, 0, len(pool.queue))
		for addr := range pool.queue {
			if !pool.locals.contains(addr) { // don't drop locals
				logger.Info("[promoteExecutables]Beats :", pool.beats[addr].String())
				addresses = append(addresses, addressByHeartbeat{addr, pool.beats[addr]})
			}
		}
		sort.Sort(addresses)

		// Drop transactions until the total is below the limit or only locals remain
		for drop := queued - pool.config.GlobalQueue; drop > 0 && len(addresses) > 0; {
			logger.Info("[promoteExecutables]drop:", drop)
			addr := addresses[len(addresses)-1]
			list := pool.queue[addr.address]

			addresses = addresses[:len(addresses)-1]

			// Drop all transactions if they are less than the overflow
			logger.Info("[promoteExecutables] Will Drop size:", list.Len())
			if size := uint64(list.Len()); size <= drop {
				for _, tx := range list.Flatten() {
					pool.removeTx(tx.Hash())
				}
				drop -= size
				queuedRateLimitCounter.Inc(int64(size))
				continue
			}
			// Otherwise drop only last few transactions
			txs := list.Flatten()
			for i := len(txs) - 1; i >= 0 && drop > 0; i-- {
				pool.removeTx(txs[i].Hash())
				drop--
				queuedRateLimitCounter.Inc(1)
			}
		}
	}
}

// demoteUnexecutables removes invalid and processed transactions from the pools
// executable/pending queue and any subsequent transactions that become unexecutable
// are moved back into the future queue.
func (pool *TxPool) demoteUnexecutables() {
	// Iterate over all accounts and demote any non-executable transactions
	for addr, list := range pool.pending {
		nonce := pool.currentState.GetNonce(addr)

		// Drop all transactions that are deemed too old (low nonce)
		for _, tx := range list.Forward(nonce) {
			hash := tx.Hash()
			logger.Tracef("Removed old pending transaction hash:0x%x", hash)
			delete(pool.all, hash)
			pool.priorited.Removed()
		}
		// Drop all transactions that are too costly (low balance ), and queue any invalids back for later
		//drops, invalids := list.Filter(pool.currentState.GetBalance(addr), 0)
		//for _, tx := range drops {
		//	hash := tx.Hash()
		//	logger.Tracef("Removed unpayable pending transaction hash:0x%x", hash)
		//	delete(pool.all, hash)
		//	pool.priorited.Removed()
		//	pendingNofundsCounter.Inc(1)
		//}
		//for _, tx := range invalids {
		//	hash := tx.Hash()
		//	logger.Tracef("Demoting pending transaction  hash:0x%x", hash)
		//	pool.enqueueTx(hash, tx)
		//}

		// If there's a gap in front, warn (should never happen) and postpone all transactions
		if list.Len() > 0 && list.txs.Get(nonce) == nil {
			for _, tx := range list.Cap(0) {
				hash := tx.Hash()
				logger.Errorf("Demoting invalidated transaction hash:0x%x", hash)
				pool.enqueueTx(hash, tx)
			}
		}
		// Delete the entire queue entry if it became empty.
		if list.Empty() {
			delete(pool.pending, addr)
			delete(pool.beats, addr)
		}
	}
}

// addressByHeartbeat is an account address tagged with its last activity timestamp.
type addressByHeartbeat struct {
	address   types.Address
	heartbeat time.Time
}

type addresssByHeartbeat []addressByHeartbeat

func (a addresssByHeartbeat) Len() int           { return len(a) }
func (a addresssByHeartbeat) Less(i, j int) bool { return a[i].heartbeat.Before(a[j].heartbeat) }
func (a addresssByHeartbeat) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// accountSet is simply a set of addresses to check for existence, and a signer
// capable of deriving addresses from transactions.
type accountSet struct {
	accounts map[types.Address]struct{}
	signer   transaction.Signer
}

// newAccountSet creates a new address set with an associated signer for sender
// derivations.
func newAccountSet(signer transaction.Signer) *accountSet {
	return &accountSet{
		accounts: make(map[types.Address]struct{}),
		signer:   signer,
	}
}

// contains checks if a given address is contained within the set.
func (as *accountSet) contains(addr types.Address) bool {
	_, exist := as.accounts[addr]
	return exist
}

// containsTx checks if the sender of a given tx is within the set. If the sender
// cannot be derived, this method returns false.
func (as *accountSet) containsTx(tx *transaction.Transaction) bool {
	if addr, err := transaction.Sender(as.signer, tx); err == nil {
		return as.contains(addr)
	}
	return false
}

// add inserts a new address into the set to track.
func (as *accountSet) add(addr types.Address) {
	as.accounts[addr] = struct{}{}
}
