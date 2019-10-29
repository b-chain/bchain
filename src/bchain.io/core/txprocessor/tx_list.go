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
// @File: tx_list.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package txprocessor

import (
	"container/heap"
	"math"
	"math/big"
	"sort"

	"bchain.io/core/transaction"

	"bchain.io/common/types"
)

// nonceHeap is a heap.Interface implementation over 64bit unsigned integers for
// retrieving sorted transaction.Transactions from the possibly gapped future queue.
type nonceHeap []uint64

func (h nonceHeap) Len() int           { return len(h) }
func (h nonceHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h nonceHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *nonceHeap) Push(x interface{}) {
	*h = append(*h, x.(uint64))
}

func (h *nonceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// txSortedMap is a nonce->transaction.Transaction hash map with a heap based index to allow
// iterating over the contents in a nonce-incrementing way.
type txSortedMap struct {
	items map[uint64]*transaction.Transaction // Hash map storing the transaction.Transaction data
	index *nonceHeap                    // Heap of nonces of all the stored transaction.Transactions (non-strict mode)
	cache transaction.Transactions            // Cache of the transaction.Transactions already sorted
}

// newTxSortedMap creates a new nonce-sorted transaction.Transaction map.
func newTxSortedMap() *txSortedMap {
	return &txSortedMap{
		items: make(map[uint64]*transaction.Transaction),
		index: new(nonceHeap),
	}
}

// Get retrieves the current transaction.Transactions associated with the given nonce.
func (m *txSortedMap) Get(nonce uint64) *transaction.Transaction {
	return m.items[nonce]
}

// Put inserts a new transaction.Transaction into the map, also updating the map's nonce
// index. If a transaction.Transaction already exists with the same nonce, it's overwritten.
func (m *txSortedMap) Put(tx *transaction.Transaction) {
	nonce := tx.Nonce()
	if m.items[nonce] == nil {
		heap.Push(m.index, nonce)
	}
	m.items[nonce], m.cache = tx, nil
}

// Forward removes all transaction.Transactions from the map with a nonce lower than the
// provided threshold. Every removed transaction.Transaction is returned for any post-removal
// maintenance.
func (m *txSortedMap) Forward(threshold uint64) transaction.Transactions {
	var removed transaction.Transactions

	// Pop off heap items until the threshold is reached
	for m.index.Len() > 0 && (*m.index)[0] < threshold {
		nonce := heap.Pop(m.index).(uint64)
		removed = append(removed, m.items[nonce])
		delete(m.items, nonce)
	}
	// If we had a cached order, shift the front
	if m.cache != nil {
		m.cache = m.cache[len(removed):]
	}
	return removed
}

// Filter iterates over the list of transaction.Transactions and removes all of them for which
// the specified function evaluates to true.
func (m *txSortedMap) Filter(filter func(*transaction.Transaction) bool) transaction.Transactions {
	var removed transaction.Transactions

	// Collect all the transaction.Transactions to filter out

	for nonce, tx := range m.items {
		if filter(tx) {
			removed = append(removed, tx)
			delete(m.items, nonce)
		}
	}
	// If transaction.Transactions were removed, the heap and cache are ruined
	if len(removed) > 0 {
		*m.index = make([]uint64, 0, len(m.items))
		for nonce := range m.items {
			*m.index = append(*m.index, nonce)
		}
		heap.Init(m.index)

		m.cache = nil
	}
	return removed
}

// Cap places a hard limit on the number of items, returning all transaction.Transactions
// exceeding that limit.
func (m *txSortedMap) Cap(threshold int) transaction.Transactions {
	// Short circuit if the number of items is under the limit
	if len(m.items) <= threshold {
		return nil
	}
	// Otherwise gather and drop the highest nonce'd transaction.Transactions
	var drops transaction.Transactions

	sort.Sort(*m.index)
	for size := len(m.items); size > threshold; size-- {
		drops = append(drops, m.items[(*m.index)[size-1]])
		delete(m.items, (*m.index)[size-1])
	}
	*m.index = (*m.index)[:threshold]
	heap.Init(m.index)

	// If we had a cache, shift the back
	if m.cache != nil {
		m.cache = m.cache[:len(m.cache)-len(drops)]
	}
	return drops
}

// Remove deletes a transaction.Transaction from the maintained map, returning whether the
// transaction.Transaction was found.
func (m *txSortedMap) Remove(nonce uint64) bool {
	// Short circuit if no transaction.Transaction is present
	_, ok := m.items[nonce]
	if !ok {
		return false
	}
	// Otherwise delete the transaction.Transaction and fix the heap index
	for i := 0; i < m.index.Len(); i++ {
		if (*m.index)[i] == nonce {
			heap.Remove(m.index, i)
			break
		}
	}
	delete(m.items, nonce)
	m.cache = nil

	return true
}

// Ready retrieves a sequentially increasing list of transaction.Transactions starting at the
// provided nonce that is ready for processing. The returned transaction.Transactions will be
// removed from the list.
//
// Note, all transaction.Transactions with nonces lower than start will also be returned to
// prevent getting into and invalid state. This is not something that should ever
// happen but better to be self correcting than failing!
func (m *txSortedMap) Ready(start uint64) transaction.Transactions {
	// Short circuit if no transaction.Transactions are available
	if m.index.Len() == 0 || (*m.index)[0] > start {
		return nil
	}
	// Otherwise start accumulating incremental transaction.Transactions
	var ready transaction.Transactions
	for next := (*m.index)[0]; m.index.Len() > 0 && (*m.index)[0] == next; next++ {
		ready = append(ready, m.items[next])
		delete(m.items, next)
		heap.Pop(m.index)
	}
	m.cache = nil

	return ready
}

// Len returns the length of the transaction.Transaction map.
func (m *txSortedMap) Len() int {
	return len(m.items)
}

// Flatten creates a nonce-sorted slice of transaction.Transactions based on the loosely
// sorted internal representation. The result of the sorting is cached in case
// it's requested again before any modifications are made to the contents.
func (m *txSortedMap) Flatten() transaction.Transactions {
	// If the sorting was not cached yet, create and cache it
	if m.cache == nil {
		m.cache = make(transaction.Transactions, 0, len(m.items))
		for _, tx := range m.items {
			m.cache = append(m.cache, tx)
		}
		sort.Sort(transaction.TxByNonce(m.cache))
	}
	// Copy the cache to prevent accidental modifications
	txs := make(transaction.Transactions, len(m.cache))
	copy(txs, m.cache)
	return txs
}

// txList is a "list" of transaction.Transactions belonging to an account, sorted by account
// nonce. The same type can be used both for storing contiguous transaction.Transactions for
// the executable/pending queue; and for storing gapped transaction.Transactions for the non-
// executable/future queue, with minor behavioral changes.
type txList struct {
	strict bool         // Whether nonces are strictly continuous or not
	txs    *txSortedMap // Heap indexed sorted hash map of the transaction.Transactions

	priority *big.Int //  (reset only if exceeds balance)
}

// newTxList create a new transaction.Transaction list for maintaining nonce-indexable fast,
// gapped, sortable transaction.Transaction lists.
func newTxList(strict bool) *txList {
	return &txList{
		strict:  strict,
		txs:     newTxSortedMap(),
		priority: new(big.Int),
	}
}

// Overlaps returns whether the transaction.Transaction specified has the same nonce as one
// already contained within the list.
func (l *txList) Overlaps(tx *transaction.Transaction) bool {
	return l.txs.Get(tx.Nonce()) != nil
}

// Add tries to insert a new transaction.Transaction into the list, returning whether the
// transaction.Transaction was accepted, and if yes, any previous transaction.Transaction it replaced.
//

func (l *txList) Add(tx *transaction.Transaction, rev uint64) (bool, *transaction.Transaction) {
	// If there's an older better transaction.Transaction, abort
	//already has
	old := l.txs.Get(tx.Nonce())
	if old != nil {
		//check priority
		if old.Priority().Int64() >= tx.Priority().Int64() {
			//Add false , because the old transaction's priority is higher than new
			return false , nil
		}
	}
	// Otherwise overwrite the old transaction.Transaction with the current one
	//has not yet
	l.txs.Put(tx)

	if priority := tx.Priority(); l.priority.Cmp(priority) < 0 {
		l.priority = priority
	}

	return true, old
}

// Forward removes all transaction.Transactions from the list with a nonce lower than the
// provided threshold. Every removed transaction.Transaction is returned for any post-removal
// maintenance.
func (l *txList) Forward(threshold uint64) transaction.Transactions {
	return l.txs.Forward(threshold)
}

// Filter removes all transaction.Transactions from the list with a cost  higher
// than the provided thresholds. Every removed transaction.Transaction is returned for any
// post-removal maintenance. Strict-mode invalidated transaction.Transactions are also
// returned.
//
// This method uses the cached costcap  to quickly decide if there's even
// a point in calculating all the costs or if the balance covers all.
func (l *txList) Filter(priority *big.Int, reserve uint64) (transaction.Transactions, transaction.Transactions) {
	// If all transaction.Transactions are below the threshold, short circuit
	if l.priority.Cmp(priority) <= 0  {
		return nil, nil
	}
	l.priority = new(big.Int).Set(priority) // Lower the caps to the thresholds


	// Filter out all the transaction.Transactions above the account's funds

	removed := l.txs.Filter(func(tx *transaction.Transaction) bool {
		return tx.Priority().Cmp(priority) > 0  })

	// If the list was strict, filter anything above the lowest nonce
	var invalids transaction.Transactions

	if l.strict && len(removed) > 0 {
		lowest := uint64(math.MaxUint64)
		for _, tx := range removed {
			if nonce := tx.Nonce(); lowest > nonce {
				lowest = nonce
			}
		}
		invalids = l.txs.Filter(func(tx *transaction.Transaction) bool { return tx.Nonce() > lowest })
	}
	return removed, invalids
}

// Cap places a hard limit on the number of items, returning all transaction.Transactions
// exceeding that limit.
func (l *txList) Cap(threshold int) transaction.Transactions {
	return l.txs.Cap(threshold)
}

// Remove deletes a transaction.Transaction from the maintained list, returning whether the
// transaction.Transaction was found, and also returning any transaction.Transaction invalidated due to
// the deletion (strict mode only).
func (l *txList) Remove(tx *transaction.Transaction) (bool, transaction.Transactions) {
	// Remove the transaction.Transaction from the set
	nonce := tx.Nonce()
	if removed := l.txs.Remove(nonce); !removed {
		return false, nil
	}
	// In strict mode, filter out non-executable transaction.Transactions
	if l.strict {
		return true, l.txs.Filter(func(tx *transaction.Transaction) bool { return tx.Nonce() > nonce })
	}
	return true, nil
}

// Ready retrieves a sequentially increasing list of transaction.Transactions starting at the
// provided nonce that is ready for processing. The returned transaction.Transactions will be
// removed from the list.
//
// Note, all transaction.Transactions with nonces lower than start will also be returned to
// prevent getting into and invalid state. This is not something that should ever
// happen but better to be self correcting than failing!
func (l *txList) Ready(start uint64) transaction.Transactions {
	return l.txs.Ready(start)
}

// Len returns the length of the transaction.Transaction list.
func (l *txList) Len() int {
	return l.txs.Len()
}

// Empty returns whether the list of transaction.Transactions is empty or not.
func (l *txList) Empty() bool {
	return l.Len() == 0
}

// Flatten creates a nonce-sorted slice of transaction.Transactions based on the loosely
// sorted internal representation. The result of the sorting is cached in case
// it's requested again before any modifications are made to the contents.
func (l *txList) Flatten() transaction.Transactions {
	return l.txs.Flatten()
}



type priorityHeap	[]*transaction.Transaction

func (h priorityHeap) Len() int           { return len(h) }

func (h priorityHeap) Less(i, j int) bool { return h[i].Priority().Cmp(h[j].Priority()) < 0 }
func (h priorityHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *priorityHeap)Push(x interface{}){
	*h = append(*h , x.(*transaction.Transaction))
}

func (h *priorityHeap)Pop()interface{}{
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type txPriorityList struct{
	all		*map[types.Hash]*transaction.Transaction
	items 	*priorityHeap
	stales	int
}

func newTxPriorityList(all *map[types.Hash]*transaction.Transaction)*txPriorityList{
	return &txPriorityList{
		all:	all,
		items:	new(priorityHeap),
	}
}


func (l *txPriorityList) Put(tx *transaction.Transaction) {
	heap.Push(l.items, tx)
}


func (l *txPriorityList) Removed() {
	// Bump the stale counter, but exit if still too low (< 25%)
	l.stales++
	if l.stales <= len(*l.items)/4 {
		return
	}

	reheap := make(priorityHeap, 0, len(*l.all))

	l.stales, l.items = 0, &reheap
	for _, tx := range *l.all {
		*l.items = append(*l.items, tx)
	}
	heap.Init(l.items)
}


func (l *txPriorityList) Cap(threshold *big.Int, local *accountSet) transaction.Transactions {
	drop := make(transaction.Transactions, 0, 128)
	save := make(transaction.Transactions, 0, 64)

	for len(*l.items) > 0 {

		tx := heap.Pop(l.items).(*transaction.Transaction)
		if _, ok := (*l.all)[tx.Hash()]; !ok {
			l.stales--
			continue
		}

		if tx.Priority().Cmp(threshold) >= 0 {
			save = append(save, tx)
			break
		}

		if local.containsTx(tx) {
			save = append(save, tx)
		} else {
			drop = append(drop, tx)
		}
	}

	for _, tx := range save {
		heap.Push(l.items, tx)
	}
	return drop
}


func (l *txPriorityList) Underpriority(tx *transaction.Transaction, local *accountSet ,threshold *big.Int) bool {

	if local.containsTx(tx) {
		return false
	}

	for len(*l.items) > 0 {
		head := []*transaction.Transaction(*l.items)[0]

		if _, ok := (*l.all)[head.Hash()]; !ok {
			l.stales--
			heap.Pop(l.items)
			continue
		}
		break
	}

	//if len(*l.items) == 0 {
	//	logger.Error("Pricing query for empty pool") // This cannot happen, print to catch programming errors
	//	return false
	//}
	cheapest := threshold

	return cheapest.Cmp(tx.Priority()) > 0
}


func (l *txPriorityList) Discard(count int, local *accountSet) transaction.Transactions {
	drop := make(transaction.Transactions, 0, count) // Remote underpriced transactions to drop
	save := make(transaction.Transactions, 0, 64)    // Local underpriced transactions to keep

	for len(*l.items) > 0 && count > 0 {

		tx := heap.Pop(l.items).(*transaction.Transaction)
		if _, ok := (*l.all)[tx.Hash()]; !ok {
			l.stales--
			continue
		}

		if local.containsTx(tx) {
			save = append(save, tx)
		} else {
			drop = append(drop, tx)
			count--
		}
	}
	for _, tx := range save {
		heap.Push(l.items, tx)
	}
	return drop
}

