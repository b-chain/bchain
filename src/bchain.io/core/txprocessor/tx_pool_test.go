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
// @File: tx_pool_test.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package txprocessor

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"math/rand"
	"bchain.io/common/types"
	"bchain.io/core"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
	"bchain.io/params"
	"bchain.io/utils/crypto"
	"bchain.io/utils/database"
	"bchain.io/utils/event"
	"testing"
	"time"
)

// Tests that transactions can be added to strict lists and list contents and
// nonce boundaries are correctly maintained.
var mSigner = transaction.NewMSigner(big.NewInt(1))

var testTxPoolConfig TxPoolConfig

func init() {
	testTxPoolConfig = DefaultTxPoolConfig
	testTxPoolConfig.Journal = ""
}

type testBlockChain struct {
	statedb       *state.StateDB
	chainHeadFeed *event.Feed
}

func (bc *testBlockChain) CurrentBlock() *block.Block {
	return block.NewBlock(&block.Header{}, nil, nil)
}

func (bc *testBlockChain) GetBlock(hash types.Hash, num uint64) *block.Block {
	return bc.CurrentBlock()
}

func (bc *testBlockChain) StateAt(hash types.Hash) (*state.StateDB, error) {
	return bc.statedb, nil
}

func (bc *testBlockChain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return bc.chainHeadFeed.Subscribe(ch)
}

func randomActions() []transaction.Action {
	address := types.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed")
	t := rand.Intn(10)
	r := []transaction.Action{}

	for i := 0; i < t; i++ {
		action := transaction.Action{}
		action.Contract = address

		action.Params = []byte{}
		action.Params = append(action.Params, byte(t))
		r = append(r, action)

	}
	return r
}

func xtransaction(nonce uint64, key *ecdsa.PrivateKey, pool *TxPool) *transaction.Transaction {

	return AsignedTransaction(nonce, key, pool)
}

func AsignedTransaction(nonce uint64, key *ecdsa.PrivateKey, pool *TxPool) *transaction.Transaction {
	var actions transaction.Actions
	for _, a := range randomActions() {
		actions = append(actions, &a)
	}
	tx, _ := transaction.SignTx(transaction.NewTransaction(nonce, actions), mSigner, key)
	pool.inter.SetPriorityForTransaction(tx)
	return tx
}

func newxtransaction(nonce uint64, key *ecdsa.PrivateKey, pool *TxPool) *transaction.Transaction {
	var actions transaction.Actions
	for _, a := range randomActions() {
		actions = append(actions, &a)
	}
	tx, _ := transaction.SignTx(transaction.NewTransaction(nonce, actions), mSigner, key)
	pool.inter.SetPriorityForTransaction(tx)
	return tx
}

var (
	TestChainConfig = &params.ChainConfig{big.NewInt(1)}
)

func setupTxPool() (*TxPool, *ecdsa.PrivateKey) {
	db, _ := database.OpenMemDB()
	statedb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	blockchain := &testBlockChain{statedb, new(event.Feed)}

	key, _ := crypto.GenerateKey()
	pool := NewTxPool(testTxPoolConfig, TestChainConfig, blockchain)

	return pool, key

}

func validateTxPoolInternals(pool *TxPool) error {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	pending, queued := pool.stats()

	if total := len(pool.all); total != pending+queued {
		return fmt.Errorf("total transaction count %d != %d pending + %d queued", total, pending, queued)
	}

	for addr, txs := range pool.pending {
		var last uint64
		for nonce := range txs.txs.items {
			if last < nonce {
				last = nonce
			}
		}

		if nonce := pool.pendingState.GetNonce(addr); nonce != last+1 {
			return fmt.Errorf("pending nonce mismatch:have %v , want %v", nonce, last+1)
		}
	}
	return nil
}

func validateEvents(events chan core.TxPreEvent, count int) error {
	for i := 0; i < count; i++ {
		select {
		case <-events:
		case <-time.After(time.Second):
			return fmt.Errorf("event #%d not fired", i)

		}
	}

	select {
	case tx := <-events:
		return fmt.Errorf("more than %d events fired: %v", count, tx.Tx)
	case <-time.After(50 * time.Millisecond):

	}
	return nil
}

func deriveSender(tx *transaction.Transaction) (types.Address, error) {
	return transaction.Sender(mSigner, tx)
}

type testChain struct {
	*testBlockChain
	address types.Address
	trigger *bool
}

func (c *testChain) State() (*state.StateDB, error) {
	stdb := c.statedb
	if *c.trigger {
		db, _ := database.OpenMemDB()
		c.statedb, _ = state.New(types.Hash{}, state.NewDatabase(db))
		c.statedb.SetNonce(c.address, 2)
		//c.statedb.SetBalance(c.address, new(big.Int).SetUint64(100000))
		*c.trigger = false
	}

	return stdb, nil
}

func TestStateChangeDuringTransactionPoolReset(t *testing.T) {
	t.Parallel()

	var (
		db, _      = database.OpenMemDB()
		key, _     = crypto.GenerateKey()
		address    = crypto.PubkeyToAddress(key.PublicKey)
		statedb, _ = state.New(types.Hash{}, state.NewDatabase(db))
		trigger    = false
	)

	//statedb.SetBalance(address, new(big.Int).SetUint64(100000000))
	blockchain := &testChain{&testBlockChain{statedb, new(event.Feed)}, address, &trigger}

	pool := NewTxPool(testTxPoolConfig, TestChainConfig, blockchain)
	defer pool.Stop()
	tx0 := xtransaction(0, key, pool)
	tx1 := xtransaction(1, key, pool)

	nonce := pool.State().GetNonce(address)
	if nonce != 0 {
		t.Fatalf("Invalid nonce , want 0 , got %d", nonce)
	}
	//set priority

	fmt.Println("Before AddRemotes")
	pool.AddRemotes(transaction.Transactions{tx0, tx1})
	fmt.Println("After AddRemotes")
	fmt.Println("Pending Len :=", len(pool.pending[address].txs.items))
	return
	nonce = pool.State().GetNonce(address)
	fmt.Println("GetNonce=", nonce)
	if nonce != 2 {
		t.Fatalf("Invalid nonce, want 2, got %d", nonce)
	}
	fmt.Println("xxxxxxxxxxxxxxxxxxxx")
	trigger = true
	pool.lockedReset(nil, nil)
	pendingTx, err := pool.Pending()
	if err != nil {
		t.Fatalf("Could not fetch pending transactions: %v", err)
	}
	fmt.Println("..................pending Len := ", len(pendingTx))
	for addr, txs := range pendingTx {
		t.Logf("%0x: %d\n", addr, len(txs))
		fmt.Println(".....")
	}

	nonce = pool.State().GetNonce(address)
	if nonce != 2 {
		t.Fatalf("Invalid nonce, want 2, got %d", nonce)
	}
	fmt.Println("nonce last get=", nonce)

}

func TestInvalidTransactions(t *testing.T) {
	t.Parallel()

	pool, key := setupTxPool()
	defer pool.Stop()

	tx := xtransaction(0, key, pool)

	from, _ := deriveSender(tx)
	//fmt.Println("Before AddBalance:", pool.currentState.GetBalance(from).Int64())
	//pool.currentState.AddBalance(from, big.NewInt(1))
	////pool.currentState.AddBalance(from,big.NewInt(1000))
	//fmt.Println("After AddBalance:", pool.currentState.GetBalance(from).Int64())

	if err := pool.AddRemote(tx); err != nil {
		fmt.Println("test addremote Err:", err)
	}

	pool.currentState.SetNonce(from, 1)
	//pool.currentState.AddBalance(from, big.NewInt(111111111))
	tx = xtransaction(0, key, pool)

	if err := pool.AddRemote(tx); err != nil {
		fmt.Println("test addremote2 Err:", err)
	}
}

func TestTransactionQueue(t *testing.T) {
	t.Parallel()

	pool, key := setupTxPool()
	defer pool.Stop()

	tx := xtransaction(0, key, pool)

	from, _ := deriveSender(tx)

	//pool.currentState.AddBalance(from, big.NewInt(1000))
	pool.lockedReset(nil, nil)
	pool.enqueueTx(tx.Hash(), tx)

	pool.promoteExecutables([]types.Address{from})

	if len(pool.pending) != 1 {
		fmt.Println("Excepted valid txs to be 1 is ", len(pool.pending))
	}

	tx = xtransaction(1, key, pool)

	from, _ = deriveSender(tx)

	pool.currentState.SetNonce(from, 2)

	_, err := pool.enqueueTx(tx.Hash(), tx)
	if err != nil {
		fmt.Println("EnqueueTx:", err)
	}

	pool.promoteExecutables([]types.Address{from})

	if _, ok := pool.pending[from].txs.items[tx.Nonce()]; ok {
		fmt.Println("expected transaction to be in tx pool")
	}

	if len(pool.queue) > 0 {
		fmt.Println("expected transaction queue to be empty.is ", len(pool.queue))
	}

	pool, key = setupTxPool()
	defer pool.Stop()

	tx1 := xtransaction(0, key, pool)
	tx2 := xtransaction(10, key, pool)
	tx3 := xtransaction(11, key, pool)

	from, _ = deriveSender(tx1)

	//pool.currentState.AddBalance(from, big.NewInt(1000))
	pool.lockedReset(nil, nil)

	fmt.Println("QueueLen before = ", len(pool.queue))
	pool.enqueueTx(tx1.Hash(), tx1)
	pool.enqueueTx(tx2.Hash(), tx2)
	pool.enqueueTx(tx3.Hash(), tx3)
	fmt.Println("QueueLen after = ", len(pool.queue[from].txs.items))

	pool.promoteExecutables([]types.Address{from})

	fmt.Println("pending len=", len(pool.pending[from].txs.items))
	fmt.Println("QueueLen = ", len(pool.queue[from].txs.items))

}

func TestTransactionNegativeValue(t *testing.T) {
	t.Parallel()

	pool, key := setupTxPool()
	defer pool.Stop()

	var actions transaction.Actions
	for _, a := range randomActions() {
		actions = append(actions, &a)
	}
	tx, _ := transaction.SignTx(transaction.NewTransaction(0, actions), mSigner, key)
	pool.inter.SetPriorityForTransaction(tx)
	from, _ := deriveSender(tx)
	_ = from
	//pool.currentState.AddBalance(from, big.NewInt(1))

	if err := pool.AddRemote(tx); err != nil {
		fmt.Println("Get A err:", err)
	}

}

func TestTransactionChainFork(t *testing.T) {
	t.Parallel()

	pool, key := setupTxPool()
	defer pool.Stop()

	addr := crypto.PubkeyToAddress(key.PublicKey)
	resetState := func() {
		db, _ := database.OpenMemDB()
		statedb, _ := state.New(types.Hash{}, state.NewDatabase(db))
		_ = addr
		//statedb.AddBalance(addr, big.NewInt(100000000000000))

		pool.chain = &testBlockChain{statedb, new(event.Feed)}
		pool.lockedReset(nil, nil)
	}
	resetState()

	tx := xtransaction(0, key, pool)
	if _, err := pool.add(tx, false); err != nil {
		fmt.Println("didn't expect error:", err)
	}

	fmt.Println("pending Len:", len(pool.pending))
	fmt.Println("queue Len:", len(pool.queue))
	pool.removeTx(tx.Hash())

	//reset the pool's internal state
	resetState()

	fmt.Println("pending Len:", len(pool.pending))
	fmt.Println("queue Len:", len(pool.queue))

	if _, err := pool.add(tx, false); err != nil {
		fmt.Println("didn't expect error:", err)
	}
	fmt.Println("pending Len:", len(pool.pending))
	fmt.Println("queue Len:", len(pool.queue))

}

func TestTransactionDoubleNonce(t *testing.T) {
	t.Parallel()

	pool, key := setupTxPool()
	defer pool.Stop()

	addr := crypto.PubkeyToAddress(key.PublicKey)
	resetState := func() {

		db, _ := database.OpenMemDB()

		statedb, _ := state.New(types.Hash{}, state.NewDatabase(db))

		//statedb.AddBalance(addr, big.NewInt(100000000000000))

		pool.chain = &testBlockChain{statedb, new(event.Feed)}
		pool.lockedReset(nil, nil)
	}

	resetState()
	fmt.Println("Notice:If All Transaction's nonce are same,the txpool just exist one ")
	var actions transaction.Actions
	for _, a := range randomActions() {
		actions = append(actions, &a)
	}
	tx1, _ := transaction.SignTx(transaction.NewTransaction(0, actions), mSigner, key)
	pool.inter.SetPriorityForTransaction(tx1)
	tx2, _ := transaction.SignTx(transaction.NewTransaction(0, actions), mSigner, key)
	pool.inter.SetPriorityForTransaction(tx2)
	tx3, _ := transaction.SignTx(transaction.NewTransaction(0, actions), mSigner, key)
	pool.inter.SetPriorityForTransaction(tx3)

	fmt.Println("tx1..1:", tx1.Priority().Int64(), "  tx..2:", tx2.Priority().Int64(), "   tx3..:", tx3.Priority().Int64())
	fmt.Println("The first one must insert Ok.....")
	if replace, err := pool.add(tx1, false); err != nil || replace {
		fmt.Printf("Error:first transaction\n ")
	}

	if replace, err := pool.add(tx2, false); err != nil || !replace {
		fmt.Println("tx2..:", tx2.Priority().Int64(), " Err:", err.Error())
	}

	pool.promoteExecutables([]types.Address{addr})

	if tx := pool.pending[addr].txs.items[0]; tx.Hash() != tx2.Hash() {
		fmt.Printf("transaction mismatch:have %x , wat %x\n", tx.Hash(), tx2.Hash())
	}

	if tx := pool.pending[addr].txs.items[0]; tx.Hash() != tx2.Hash() {
		fmt.Printf("transaction mismatch:have %x , wat %x\n", tx.Hash(), tx2.Hash())
	}

	pool.add(tx3, false)

	pool.promoteExecutables([]types.Address{addr})

	if pool.pending[addr].Len() != 1 {
		t.Errorf("Error:expected %d pending trasactions , got\n", pool.pending[addr].Len())
	}

	if tx := pool.pending[addr].txs.items[0]; tx.Hash() != tx2.Hash() {
		fmt.Printf("transaction mismatch :have %x ,want %x\n", tx.Hash(), tx2.Hash())
	}

	if len(pool.all) != 1 {
		t.Errorf("Error:expected 1 total transactions , got %d\n", len(pool.all))
	}

	fmt.Println("Test ok.....")

}

func TestTransactionMissingNonce(t *testing.T) {
	t.Parallel()

	pool, key := setupTxPool()
	defer pool.Stop()

	addr := crypto.PubkeyToAddress(key.PublicKey)
	//pool.currentState.AddBalance(addr, big.NewInt(100000000000000))

	tx := xtransaction(1, key, pool)
	if _, err := pool.add(tx, false); err != nil {
		fmt.Printf("didn't expect error:%v", err)
	}

	if len(pool.pending) != 0 {
		fmt.Errorf("expected 0 pending transactions,got %d", len(pool.pending))
	}

	if pool.queue[addr].Len() != 1 {
		fmt.Printf("expected 1 queued transaction , got %d\n", pool.queue[addr].Len())
	}

	if len(pool.all) != 1 {
		fmt.Printf("expected 1 total transactions,got %d\n", len(pool.all))
	}

	fmt.Println("test ok....")
}

func TestTransactionDroppng(t *testing.T) {
	t.Parallel()

	pool, key := setupTxPool()
	defer pool.Stop()

	account, _ := deriveSender(xtransaction(0, key, pool))
	//pool.currentState.AddBalance(account, big.NewInt(1000))

	var (
		tx0  = newxtransaction(0, key, pool)
		tx1  = newxtransaction(1, key, pool)
		tx2  = newxtransaction(2, key, pool)
		tx10 = newxtransaction(10, key, pool)
		tx11 = newxtransaction(11, key, pool)
		tx12 = newxtransaction(12, key, pool)
	)

	pool.promoteTx(account, tx0.Hash(), tx0)
	pool.promoteTx(account, tx1.Hash(), tx1)
	pool.promoteTx(account, tx2.Hash(), tx2)

	pool.enqueueTx(tx10.Hash(), tx10)
	pool.enqueueTx(tx11.Hash(), tx11)
	pool.enqueueTx(tx12.Hash(), tx12)
	fmt.Println("Do something 1")
	fmt.Println("Do something 2")
	fmt.Println("Do something 3")
	fmt.Println("Do something 4")
	if pool.pending[account].Len() != 3 {
		t.Errorf("pending transaction mismatch:have %d ,want %d", pool.pending[account].Len(), 3)
	}

	if pool.queue[account].Len() != 3 {
		t.Errorf("queue tranaction mismatch :have %d ,want %d", pool.queue[account].Len(), 3)
	}

	if len(pool.all) != 6 {
		t.Errorf("total transaction mismatch :have %d , want %d", pool.pending[account].Len(), 3)
	}

	pool.lockedReset(nil, nil)

	if pool.pending[account].Len() != 3 {
		t.Errorf("pending transaction mismatch :have %d ,want %d", pool.pending[account].Len(), 3)
	}

	if pool.queue[account].Len() != 3 {
		t.Errorf("queued tranaction mismatch :have %d , want %d", pool.queue[account].Len(), 3)
	}

	if len(pool.all) != 6 {
		t.Errorf("total transaction mismatch :have %d ,want %d", len(pool.all), 6)
	}
	//fmt.Println("1now Balance:", pool.currentState.GetBalance(account))
	fmt.Println("1len pending:", len(pool.pending[account].txs.items))
	fmt.Println("1len queue:", len(pool.queue[account].txs.items))

	//pool.currentState.AddBalance(account, big.NewInt(-500))
	pool.lockedReset(nil, nil)

	//fmt.Println("2now Balance:", pool.currentState.GetBalance(account))
	fmt.Println("2len pending:", len(pool.pending[account].txs.items))
	fmt.Println("2len queue:", len(pool.queue[account].txs.items))
	if _, ok := pool.pending[account].txs.items[tx0.Nonce()]; !ok {
		t.Errorf("funded pending transaction missing: %v", tx0)
	}
	if _, ok := pool.pending[account].txs.items[tx1.Nonce()]; !ok {
		t.Errorf("funded pending transaction missing: %v", tx1)
	}
	fmt.Println("len all=", len(pool.all))
	fmt.Println("len pending=", len(pool.pending[account].txs.items))
	if _, ok := pool.pending[account].txs.items[tx2.Nonce()]; !ok {
		t.Errorf("out-of-fund pending transaction present: %v", tx2)
	}

	if _, ok := pool.queue[account].txs.items[tx10.Nonce()]; !ok {
		t.Errorf("funded queued transaction missing: %v", tx10)
	}
	if _, ok := pool.queue[account].txs.items[tx11.Nonce()]; !ok {
		t.Errorf("funded queued transaction missing: %v", tx11)
	}
	if _, ok := pool.queue[account].txs.items[tx12.Nonce()]; !ok {
		t.Errorf("out-of-fund queued transaction present: %v", tx12)
	}

	if len(pool.all) != 6 {
		t.Errorf("total transaction mismatch: have %d, want %d", len(pool.all), 4)
	}

	fmt.Println("test ok....")
}

func TestPendingState(t *testing.T) {
	t.Parallel()

	pool, key := setupTxPool()
	defer pool.Stop()

	account, _ := deriveSender(xtransaction(0, key, pool))
	//pool.currentState.AddBalance(account, big.NewInt(1000))

	var (
		tx0 = newxtransaction(0, key, pool)
		tx1 = newxtransaction(1, key, pool)
		tx2 = newxtransaction(2, key, pool)
	)
	pool.enqueueTx(tx1.Hash(), tx1)
	pool.promoteExecutables([]types.Address{account})
	fmt.Println("pendingState:", pool.pendingState.GetNonce(account))

	pool.enqueueTx(tx0.Hash(), tx0)
	pool.promoteExecutables([]types.Address{account})
	fmt.Println("pendingState:", pool.pendingState.GetNonce(account))

	pool.enqueueTx(tx2.Hash(), tx2)
	pool.promoteExecutables([]types.Address{account})
	fmt.Println("pendingState:", pool.pendingState.GetNonce(account))

}

func TestTransactionPostponing(t *testing.T) {
	t.Parallel()

	pool, key := setupTxPool()
	defer pool.Stop()

	account, _ := deriveSender(newxtransaction(0, key, pool))
	//pool.currentState.AddBalance(account, big.NewInt(1000))

	txns := []*transaction.Transaction{}

	for i := 0; i < 100; i++ {
		var tx *transaction.Transaction
		if i%2 == 0 {
			tx = newxtransaction(uint64(i), key, pool)
		} else {
			tx = newxtransaction(uint64(i), key, pool)
		}
		pool.promoteTx(account, tx.Hash(), tx)
		txns = append(txns, tx)
	}

	if pool.pending[account].Len() != len(txns) {
		t.Errorf("pending transaction mismatch :have %d , want %d", pool.pending[account].Len(), len(txns))
	}

	if len(pool.queue) != 0 {
		t.Errorf("queued transaction mismatch :have %d， want %d", pool.queue[account].Len(), 0)
	}

	if len(pool.all) != len(txns) {
		t.Errorf("total transaction mismatch: have %d, want %d", len(pool.all), len(txns))
	}

	pool.lockedReset(nil, nil)

	if pool.pending[account].Len() != len(txns) {
		t.Errorf("pending transaction mismatch: have %d, want %d", pool.pending[account].Len(), len(txns))
	}
	if len(pool.queue) != 0 {
		t.Errorf("queued transaction mismatch: have %d, want %d", pool.queue[account].Len(), 0)
	}
	if len(pool.all) != len(txns) {
		t.Errorf("total transaction mismatch: have %d, want %d", len(pool.all), len(txns))
	}

	// Reduce the balance of the account, and check that transactions are reorganised
	//pool.currentState.AddBalance(account, big.NewInt(-750))
	pool.lockedReset(nil, nil)

	if _, ok := pool.pending[account].txs.items[txns[0].Nonce()]; !ok {
		t.Errorf("tx %d: valid and funded transaction missing from pending pool: %v", 0, txns[0])
	}

	if _, ok := pool.queue[account]; ok {
		t.Errorf("tx %d: valid and funded transaction present in future queue: %v", 0, txns[0])
	}

}

func TestTransactionGapFilling(t *testing.T) {

	t.Parallel()
	pool, key := setupTxPool()
	defer pool.Stop()

	account, _ := deriveSender(newxtransaction(0, key, pool))
	_ = account
	//pool.currentState.AddBalance(account, big.NewInt(1000000))

	events := make(chan core.TxPreEvent, 70)
	sub := pool.txFeed.Subscribe(events)
	defer sub.Unsubscribe()

	if err := pool.AddRemote(newxtransaction(0, key, pool)); err != nil {
		t.Fatalf("failed to add pending transaction:%v", err)
	}

	if err := pool.AddRemote(newxtransaction(2, key, pool)); err != nil {
		t.Fatalf("failed to add queued transaction:%v", err)
	}

	pending, queued := pool.stats()

	if pending != 1 {
		t.Fatalf("pending transactions mismatched :have %d , want %d", pending, 1)
	}

	if queued != 1 {
		t.Fatalf("queued transactions mismatched :have %d , want %d", queued, 1)
	}

	if err := validateEvents(events, 1); err != nil {
		t.Fatalf("original event firing failed :%v", err)
	}

	if err := validateTxPoolInternals(pool); err != nil {
		t.Fatalf("pool internal state corrupted :%v", err)
	}

	if err := pool.AddRemote(newxtransaction(1, key, pool)); err != nil {
		t.Fatalf("failed to add gapped transaction:%v", err)
	}

	pending, queued = pool.stats()

	if pending != 3 {
		t.Fatalf("pending transactions mismatched :have %d , want %d", pending, 3)
	}

	if queued != 0 {
		t.Fatalf("queued transactions mismatched :have %d , want %d", queued, 0)
	}

	if err := validateEvents(events, 2); err != nil {
		t.Fatalf("gap-filling event firing failed:%v", err)
	}

	if err := validateTxPoolInternals(pool); err != nil {
		t.Fatalf("pool internal state corrupted :%v", err)
	}
}

func TestTransactionQueueAccountLimiting(t *testing.T) {
	t.Parallel()

	pool, key := setupTxPool()
	defer pool.Stop()

	account, _ := deriveSender(newxtransaction(0, key, pool))
	//pool.currentState.AddBalance(account, big.NewInt(1000000))

	for i := uint64(1); i <= testTxPoolConfig.AccountQueue+5; i++ {
		if err := pool.AddRemote(newxtransaction(i, key, pool)); err != nil {
			t.Fatalf("tx %d : failed to add transaction :%v", i, err)
		}

		if len(pool.pending) != 0 {
			t.Errorf("tx %d:pending pool size mismatch :have %d , want %d", i, len(pool.pending), 0)
		}

		if i <= testTxPoolConfig.AccountQueue {
			if pool.queue[account].Len() != int(i) {
				t.Errorf("tx %d: queue size mismatch :have %d , want %d", i, pool.queue[account].Len(), i)
			}
		} else {
			if pool.queue[account].Len() != int(testTxPoolConfig.AccountQueue) {
				t.Errorf("tx %d:queue limit mismatch :have %d ,want %d", i, pool.queue[account].Len(), testTxPoolConfig.AccountQueue)
			}
		}

	}
	if len(pool.all) != int(testTxPoolConfig.AccountQueue) {
		t.Errorf("total transaction mismatch :have %d , want %d", len(pool.all), testTxPoolConfig.AccountQueue)
	}
}

func TestBefore(t *testing.T) {
	t.Parallel()

	var (
		db, _      = database.OpenMemDB()
		key, _     = crypto.GenerateKey()
		address    = crypto.PubkeyToAddress(key.PublicKey)
		statedb, _ = state.New(types.Hash{}, state.NewDatabase(db))
		trigger    = false
	)

	//statedb.SetBalance(address, new(big.Int).SetUint64(100000000))
	blockchain := &testChain{&testBlockChain{statedb, new(event.Feed)}, address, &trigger}
	pool := NewTxPool(testTxPoolConfig, TestChainConfig, blockchain)
	defer pool.Stop()
	tx0 := xtransaction(0, key, pool)
	tx1 := xtransaction(1, key, pool)

	nonce := pool.State().GetNonce(address)
	if nonce != 0 {
		t.Fatalf("Invalid nonce , want 0 , got %d", nonce)
	}

	fmt.Println("Before AddRemotes")
	pool.AddRemotes(transaction.Transactions{tx0, tx1})
	fmt.Println("After AddRemotes")
	fmt.Println("Pending Len :=", len(pool.pending[address].txs.items))
}

func TestTransactionPendingLimiting(t *testing.T) {
	t.Parallel()

	db, _ := database.OpenMemDB()
	statedb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	blockchain := &testBlockChain{statedb, new(event.Feed)}
	config := testTxPoolConfig
	config.GlobalQueue = config.AccountQueue*3 - 1
	config.GlobalSlots = 350
	fmt.Println("config.AccountSlots:", config.AccountSlots)
	fmt.Println("config.GlobalSlots:", config.GlobalSlots)

	pool := NewTxPool(config, TestChainConfig, blockchain)
	defer pool.Stop()

	keys := make([]*ecdsa.PrivateKey, 5)
	address := make([]types.Address, 5)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
		address[i] = crypto.PubkeyToAddress(keys[i].PublicKey)
		//pool.currentState.AddBalance(address[i], big.NewInt(1000000))

	}

	local := keys[len(keys)-1]
	localAddress := address[len(keys)-1]

	_ = local
	_ = localAddress

	txCnts := make([]int, 0)
	txCnts = append(txCnts, 15)
	txCnts = append(txCnts, 80)
	txCnts = append(txCnts, 200)
	txCnts = append(txCnts, 230)

	nonces := make(map[types.Address]uint64)

	txs := make(transaction.Transactions, 0, 3*config.GlobalQueue)

	for i := 0; i < len(txCnts); i++ {
		addr := address[i]
		key := keys[i]

		for j := 0; j < txCnts[i]; j++ {
			txs = append(txs, newxtransaction(nonces[addr], key, pool))
			nonces[addr]++
		}

	}

	for _, cnt := range nonces {
		fmt.Println("Sig Address Txs Src Len:", cnt)
	}

	fmt.Println("Before------all Txs Length:", len(txs), " GlobalQueue:", config.GlobalQueue)
	pool.AddRemotes(txs)
	_, nowQueueLen := pool.stats()
	fmt.Println("After-------all Queue Length:", nowQueueLen, "GlobalQueue:", config.GlobalQueue)

	for _, list := range pool.pending {
		fmt.Println("Sig Address Txs Len:", len(list.txs.items))
	}

}

func TestTransactionQueueGlobalLimiting(t *testing.T) {
	testTransactionQueueGlobalLimiting(t, false)
}

func testTransactionQueueGlobalLimiting(t *testing.T, nolocals bool) {
	t.Parallel()

	db, _ := database.OpenMemDB()
	statedb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	blockchain := &testBlockChain{statedb, new(event.Feed)}
	config := testTxPoolConfig
	config.NoLocals = nolocals
	config.GlobalQueue = config.AccountQueue*3 - 1
	fmt.Println("config.AccountQueue:", config.AccountQueue)
	fmt.Println("config.GlobalQueue:", config.GlobalQueue)

	pool := NewTxPool(config, TestChainConfig, blockchain)
	defer pool.Stop()

	keys := make([]*ecdsa.PrivateKey, 5)
	address := make([]types.Address, 5)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
		address[i] = crypto.PubkeyToAddress(keys[i].PublicKey)
		//pool.currentState.AddBalance(address[i], big.NewInt(1000000))

	}

	local := keys[len(keys)-1]
	localAddress := address[len(keys)-1]

	_ = local
	_ = localAddress

	nonces := make(map[types.Address]uint64)

	txs := make(transaction.Transactions, 0, 3*config.GlobalQueue)

	for len(txs) < cap(txs) {
		i := rand.Intn(len(keys) - 1)

		addr := address[i]
		key := keys[i]

		txs = append(txs, newxtransaction(nonces[addr]+1, key, pool))
		nonces[addr]++
	}

	for _, cnt := range nonces {
		fmt.Println("Sig Address Txs Src Len:", cnt)
	}
	fmt.Println("Before------all Txs Length:", len(txs), " GlobalQueue:", config.GlobalQueue)
	pool.AddRemotes(txs)
	_, nowQueueLen := pool.stats()
	fmt.Println("After-------all Queue Length:", nowQueueLen, "GlobalQueue:", config.GlobalQueue)

	for _, list := range pool.queue {
		fmt.Println("Sig Address Txs Len:", len(list.txs.items))
	}

}
