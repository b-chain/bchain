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
// @File: queue.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

// Contains the block download scheduler to collect download tasks and schedule
// them in an ordered, and throttled way.

package downloader

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"bchain.io/common"
	"bchain.io/common/types"
	"bchain.io/utils/metrics"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/transaction"
)

var blockCacheLimit = 8192 // Maximum number of blocks to cache before throttling the download

var (
	errNoFetchesPending = errors.New("no fetches pending")
	errStaleDelivery    = errors.New("stale delivery")
)

// fetchRequest is a currently running data retrieval operation.
type fetchRequest struct {
	Peer    *peerConnection     // Peer to which the request was sent
	From    uint64              // [bchain/62] Requested chain element index (used for skeleton fills only)
	Hashes  map[types.Hash]int //  [bchain/61] Requested hashes with their insertion index (priority)
	Headers []*block.Header     // [bchain/62] Requested headers, sorted by request order
	Time    time.Time           // Time when the request was made
}

// fetchResult is a struct collecting partial results from data fetchers until
// all outstanding pieces complete and the result as a whole can be processed.
type fetchResult struct {
	Pending int // Number of data fetches still pending

	Header       *block.Header
	Transactions transaction.Transactions
	Receipts     transaction.Receipts
	Certificate  []byte
}

// queue represents hashes that are either need fetching or are being fetched
type queue struct {
	mode          SyncMode // Synchronisation mode to decide on the block parts to schedule for fetching
	fastSyncPivot uint64   // Block number where the fast sync pivots into archive synchronisation mode

	headerHead types.Hash // [bchain/62] Hash of the last queued header to verify order

	// Headers are "special", they download in batches, supported by a skeleton chain
	headerTaskPool  map[uint64]*block.Header       // [bchain/62] Pending header retrieval tasks, mapping starting indexes to skeleton headers
	headerTaskQueue *prque.Prque                   // [bchain/62] Priority queue of the skeleton indexes to fetch the filling headers for
	headerPeerMiss  map[string]map[uint64]struct{} // [bchain/62] Set of per-peer header batches known to be unavailable
	headerPendPool  map[string]*fetchRequest       // [bchain/62] Currently pending header retrieval operations
	headerResults   []*block.Header                // [bchain/62] Result cache accumulating the completed headers
	headerProced    int                            // [bchain/62] Number of headers already processed from the results
	headerOffset    uint64                         // [bchain/62] Number of the first header in the result cache
	headerContCh    chan bool                      // [bchain/62] Channel to notify when header download finishes

	// All data retrievals below are based on an already assembles header chain
	blockTaskPool  map[types.Hash]*block.Header // [bchain/62] Pending block (body) retrieval tasks, mapping hashes to headers
	blockTaskQueue *prque.Prque                  // [bchain/62] Priority queue of the headers to fetch the blocks (bodies) for
	blockPendPool  map[string]*fetchRequest      // [bchain/62] Currently pending block (body) retrieval operations
	blockDonePool  map[types.Hash]struct{}      // [bchain/62] Set of the completed block (body) fetches

	certificateTaskPool  map[types.Hash]*block.Header  // [bchain/62] Pending block (certificate) retrieval tasks, mapping hashes to headers
	certificateTaskQueue *prque.Prque                   // [bchain/62] Priority queue of the headers to fetch the blocks (certificate) for
	certificatePendPool  map[string]*fetchRequest       // [bchain/62] Currently pending block (certificate) retrieval operations
	certificateDonePool  map[types.Hash]struct{}       // [bchain/62] Set of the completed block (certificate) fetches

	receiptTaskPool  map[types.Hash]*block.Header // [bchain/63] Pending receipt retrieval tasks, mapping hashes to headers
	receiptTaskQueue *prque.Prque                  // [bchain/63] Priority queue of the headers to fetch the receipts for
	receiptPendPool  map[string]*fetchRequest      // [bchain/63] Currently pending receipt retrieval operations
	receiptDonePool  map[types.Hash]struct{}      // [bchain/63] Set of the completed receipt fetches

	resultCache  []*fetchResult // Downloaded but not yet delivered fetch results
	resultOffset uint64         // Offset of the first cached fetch result in the block chain

	lock   *sync.Mutex
	active *sync.Cond
	closed bool
	certificateCheck certificateCheckFn
}

// newQueue creates a new download queue for scheduling block retrieval.
func newQueue(certificateCheck certificateCheckFn) *queue {
	lock := new(sync.Mutex)
	return &queue{
		headerPendPool:   make(map[string]*fetchRequest),
		headerContCh:     make(chan bool),
		blockTaskPool:    make(map[types.Hash]*block.Header),
		blockTaskQueue:   prque.New(),
		blockPendPool:    make(map[string]*fetchRequest),
		blockDonePool:    make(map[types.Hash]struct{}),
		certificateTaskPool:    make(map[types.Hash]*block.Header),
		certificateTaskQueue:   prque.New(),
		certificatePendPool:    make(map[string]*fetchRequest),
		certificateDonePool:    make(map[types.Hash]struct{}),
		receiptTaskPool:  make(map[types.Hash]*block.Header),
		receiptTaskQueue: prque.New(),
		receiptPendPool:  make(map[string]*fetchRequest),
		receiptDonePool:  make(map[types.Hash]struct{}),
		resultCache:      make([]*fetchResult, blockCacheLimit),
		active:           sync.NewCond(lock),
		lock:             lock,
		certificateCheck: certificateCheck,
	}
}

// Reset clears out the queue contents.
func (q *queue) Reset() {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.closed = false
	q.mode = FullSync
	q.fastSyncPivot = 0

	q.headerHead = types.Hash{}

	q.headerPendPool = make(map[string]*fetchRequest)

	q.blockTaskPool = make(map[types.Hash]*block.Header)
	q.blockTaskQueue.Reset()
	q.blockPendPool = make(map[string]*fetchRequest)
	q.blockDonePool = make(map[types.Hash]struct{})

	q.certificateTaskPool = make(map[types.Hash]*block.Header)
	q.certificateTaskQueue.Reset()
	q.certificatePendPool = make(map[string]*fetchRequest)
	q.certificateDonePool = make(map[types.Hash]struct{})

	q.receiptTaskPool = make(map[types.Hash]*block.Header)
	q.receiptTaskQueue.Reset()
	q.receiptPendPool = make(map[string]*fetchRequest)
	q.receiptDonePool = make(map[types.Hash]struct{})

	q.resultCache = make([]*fetchResult, blockCacheLimit)
	q.resultOffset = 0
}

// Close marks the end of the sync, unblocking WaitResults.
// It may be called even if the queue is already closed.
func (q *queue) Close() {
	q.lock.Lock()
	q.closed = true
	q.lock.Unlock()
	q.active.Broadcast()
}

// PendingHeaders retrieves the number of header requests pending for retrieval.
func (q *queue) PendingHeaders() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.headerTaskQueue.Size()
}

// PendingBlocks retrieves the number of block (body) requests pending for retrieval.
func (q *queue) PendingBlocks() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.blockTaskQueue.Size()
}

// PendingCertificates retrieves the number of block (certificate) requests pending for retrieval.
func (q *queue) PendingCertificates() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.certificateTaskQueue.Size()
}

// PendingReceipts retrieves the number of block receipts pending for retrieval.
func (q *queue) PendingReceipts() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.receiptTaskQueue.Size()
}

// InFlightHeaders retrieves whether there are header fetch requests currently
// in flight.
func (q *queue) InFlightHeaders() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.headerPendPool) > 0
}

// InFlightBlocks retrieves whether there are block fetch requests currently in
// flight.
func (q *queue) InFlightBlocks() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.blockPendPool) > 0
}

// InFlightCertificates retrieves whether there are block fetch requests currently in
// flight.
func (q *queue) InFlightCertificates() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.certificatePendPool) > 0
}

// InFlightReceipts retrieves whether there are receipt fetch requests currently
// in flight.
func (q *queue) InFlightReceipts() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.receiptPendPool) > 0
}

// Idle returns if the queue is fully idle or has some data still inside.
func (q *queue) Idle() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	queued := q.blockTaskQueue.Size() + q.receiptTaskQueue.Size() + q.certificateTaskQueue.Size()
	pending := len(q.blockPendPool) + len(q.receiptPendPool) + len(q.certificatePendPool)
	cached := len(q.blockDonePool) + len(q.receiptDonePool) + len(q.certificateDonePool)

	return (queued + pending + cached) == 0
}

// FastSyncPivot retrieves the currently used fast sync pivot point.
func (q *queue) FastSyncPivot() uint64 {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.fastSyncPivot
}

// ShouldThrottleBlocks checks if the download should be throttled (active block (body)
// fetches exceed block cache).
func (q *queue) ShouldThrottleBlocks() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Calculate the currently in-flight block (body) requests
	pending := 0
	for _, request := range q.blockPendPool {
		pending += len(request.Hashes) + len(request.Headers)
	}
	// Throttle if more blocks (bodies) are in-flight than free space in the cache
	return pending >= len(q.resultCache)-len(q.blockDonePool)
}

// ShouldThrottleCertificates checks if the download should be throttled (active block (certificate)
// fetches exceed block cache).
func (q *queue) ShouldThrottleCertificates() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Calculate the currently in-flight block (certificate) requests
	pending := 0
	for _, request := range q.certificatePendPool {
		pending += len(request.Hashes) + len(request.Headers)
	}
	// Throttle if more blocks (certificates) are in-flight than free space in the cache
	return pending >= len(q.resultCache)-len(q.certificateDonePool)
}

// ShouldThrottleReceipts checks if the download should be throttled (active receipt
// fetches exceed block cache).
func (q *queue) ShouldThrottleReceipts() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Calculate the currently in-flight receipt requests
	pending := 0
	for _, request := range q.receiptPendPool {
		pending += len(request.Headers)
	}
	// Throttle if more receipts are in-flight than free space in the cache
	return pending >= len(q.resultCache)-len(q.receiptDonePool)
}

// ScheduleSkeleton adds a batch of header retrieval tasks to the queue to fill
// up an already retrieved header skeleton.
func (q *queue) ScheduleSkeleton(from uint64, skeleton []*block.Header) {
	q.lock.Lock()
	defer q.lock.Unlock()

	// No skeleton retrieval can be in progress, fail hard if so (huge implementation bug)
	if q.headerResults != nil {
		panic("skeleton assembly already in progress")
	}
	// Shedule all the header retrieval tasks for the skeleton assembly
	q.headerTaskPool = make(map[uint64]*block.Header)
	q.headerTaskQueue = prque.New()
	q.headerPeerMiss = make(map[string]map[uint64]struct{}) // Reset availability to correct invalid chains
	q.headerResults = make([]*block.Header, len(skeleton)*MaxHeaderFetch)
	q.headerProced = 0
	q.headerOffset = from
	q.headerContCh = make(chan bool, 1)

	for i, header := range skeleton {
		index := from + uint64(i*MaxHeaderFetch)

		q.headerTaskPool[index] = header
		q.headerTaskQueue.Push(index, -float32(index))
	}
}

// RetrieveHeaders retrieves the header chain assemble based on the scheduled
// skeleton.
func (q *queue) RetrieveHeaders() ([]*block.Header, int) {
	q.lock.Lock()
	defer q.lock.Unlock()

	headers, proced := q.headerResults, q.headerProced
	q.headerResults, q.headerProced = nil, 0

	return headers, proced
}

// Schedule adds a set of headers for the download queue for scheduling, returning
// the new headers encountered.
func (q *queue) Schedule(headers []*block.Header, from uint64) []*block.Header {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Insert all the headers prioritised by the contained block number
	inserts := make([]*block.Header, 0, len(headers))
	for _, header := range headers {
		// Make sure chain order is honoured and preserved throughout
		hash := header.Hash()
		if header.Number == nil || header.Number.IntVal.Uint64() != from {
			logger.Warn("Header broke chain ordering", "number", header.Number.IntVal.String(), "hash", hash.String(), "expected", from)
			break
		}
		if q.headerHead != (types.Hash{}) && q.headerHead != header.ParentHash {
			logger.Warn("Header broke chain ancestry", "number", header.Number.IntVal.String(), "hash", hash.String())
			break
		}
		// Make sure no duplicate requests are executed
		if _, ok := q.blockTaskPool[hash]; ok {
			logger.Warn("Header  already scheduled for block fetch", "number", header.Number.IntVal.String(), "hash", hash.String())
			continue
		}
		if _, ok := q.receiptTaskPool[hash]; ok {
			logger.Warn("Header already scheduled for receipt fetch", "number", header.Number.IntVal.String(), "hash", hash.String())
			continue
		}
		if _, ok := q.certificateTaskPool[hash]; ok {
			logger.Warn("Header already scheduled for certificate fetch", "number", header.Number.IntVal.String(), "hash", hash.String())
			continue
		}
		// Queue the header for content retrieval
		q.blockTaskPool[hash] = header
		q.blockTaskQueue.Push(header, -float32(header.Number.IntVal.Uint64()))

		q.certificateTaskPool[hash] = header
		q.certificateTaskQueue.Push(header, -float32(header.Number.IntVal.Uint64()))

		if q.mode == FastSync && header.Number.IntVal.Uint64() <= q.fastSyncPivot {
			// Fast phase of the fast sync, retrieve receipts too
			q.receiptTaskPool[hash] = header
			q.receiptTaskQueue.Push(header, -float32(header.Number.IntVal.Uint64()))
		}
		inserts = append(inserts, header)
		q.headerHead = hash
		from++
	}
	return inserts
}

// WaitResults retrieves and permanently removes a batch of fetch
// results from the cache. the result slice will be empty if the queue
// has been closed.
func (q *queue) WaitResults() []*fetchResult {
	q.lock.Lock()
	defer q.lock.Unlock()

	nproc := q.countProcessableItems()
	for nproc == 0 && !q.closed {
		q.active.Wait()
		nproc = q.countProcessableItems()
	}
	results := make([]*fetchResult, nproc)
	copy(results, q.resultCache[:nproc])
	if len(results) > 0 {
		// Mark results as done before dropping them from the cache.
		for _, result := range results {
			hash := result.Header.Hash()
			delete(q.blockDonePool, hash)
			delete(q.certificateDonePool, hash)
			delete(q.receiptDonePool, hash)
		}
		// Delete the results from the cache and clear the tail.
		copy(q.resultCache, q.resultCache[nproc:])
		for i := len(q.resultCache) - nproc; i < len(q.resultCache); i++ {
			q.resultCache[i] = nil
		}
		// Advance the expected block number of the first cache entry.
		q.resultOffset += uint64(nproc)
	}
	return results
}

// countProcessableItems counts the processable items.
func (q *queue) countProcessableItems() int {
	for i, result := range q.resultCache {
		// Don't process incomplete or unavailable items.
		if result == nil || result.Pending > 0 {
			return i
		}
		// Stop before processing the pivot block to ensure that
		// resultCache has space for fsHeaderForceVerify items. Not
		// doing this could leave us unable to download the required
		// amount of headers.
		if q.mode == FastSync && result.Header.Number.IntVal.Uint64() == q.fastSyncPivot {
			for j := 0; j < fsHeaderForceVerify; j++ {
				if i+j+1 >= len(q.resultCache) || q.resultCache[i+j+1] == nil {
					return i
				}
			}
		}
	}
	return len(q.resultCache)
}

// ReserveHeaders reserves a set of headers for the given peer, skipping any
// previously failed batches.
func (q *queue) ReserveHeaders(p *peerConnection, count int) *fetchRequest {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Short circuit if the peer's already downloading something (sanity check to
	// not corrupt state)
	if _, ok := q.headerPendPool[p.id]; ok {
		return nil
	}
	// Retrieve a batch of hashes, skipping previously failed ones
	send, skip := uint64(0), []uint64{}
	for send == 0 && !q.headerTaskQueue.Empty() {
		from, _ := q.headerTaskQueue.Pop()
		if q.headerPeerMiss[p.id] != nil {
			if _, ok := q.headerPeerMiss[p.id][from.(uint64)]; ok {
				skip = append(skip, from.(uint64))
				continue
			}
		}
		send = from.(uint64)
	}
	// Merge all the skipped batches back
	for _, from := range skip {
		q.headerTaskQueue.Push(from, -float32(from))
	}
	// Assemble and return the block download request
	if send == 0 {
		return nil
	}
	request := &fetchRequest{
		Peer: p,
		From: send,
		Time: time.Now(),
	}
	q.headerPendPool[p.id] = request
	return request
}

// ReserveBodies reserves a set of body fetches for the given peer, skipping any
// previously failed downloads. Beside the next batch of needed fetches, it also
// returns a flag whether empty blocks were queued requiring processing.
func (q *queue) ReserveBodies(p *peerConnection, count int) (*fetchRequest, bool, error) {
	isNoop := func(header *block.Header) bool {
		return header.TxRootHash == block.EmptyRootHash
	}
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.reserveHeaders(p, count, q.blockTaskPool, q.blockTaskQueue, q.blockPendPool, q.blockDonePool, isNoop)
}

// ReserveCertificate reserves a set of certificate fetches for the given peer, skipping any
// previously failed downloads. Beside the next batch of needed fetches, it also
// returns a flag whether empty blocks were queued requiring processing.
func (q *queue) ReserveCertificate(p *peerConnection, count int) (*fetchRequest, bool, error) {
	isNoop := func(header *block.Header) bool {
		return false
	}
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.reserveHeaders(p, count, q.certificateTaskPool, q.certificateTaskQueue, q.certificatePendPool, q.certificateDonePool, isNoop)
}

// ReserveReceipts reserves a set of receipt fetches for the given peer, skipping
// any previously failed downloads. Beside the next batch of needed fetches, it
// also returns a flag whether empty receipts were queued requiring importing.
func (q *queue) ReserveReceipts(p *peerConnection, count int) (*fetchRequest, bool, error) {
	isNoop := func(header *block.Header) bool {
		return header.ReceiptRootHash == block.EmptyRootHash
	}
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.reserveHeaders(p, count, q.receiptTaskPool, q.receiptTaskQueue, q.receiptPendPool, q.receiptDonePool, isNoop)
}

// reserveHeaders reserves a set of data download operations for a given peer,
// skipping any previously failed ones. This method is a generic version used
// by the individual special reservation functions.
//
// Note, this method expects the queue lock to be already held for writing. The
// reason the lock is not obtained in here is because the parameters already need
// to access the queue, so they already need a lock anyway.
func (q *queue) reserveHeaders(p *peerConnection, count int, taskPool map[types.Hash]*block.Header, taskQueue *prque.Prque,
	pendPool map[string]*fetchRequest, donePool map[types.Hash]struct{}, isNoop func(*block.Header) bool) (*fetchRequest, bool, error) {
	// Short circuit if the pool has been depleted, or if the peer's already
	// downloading something (sanity check not to corrupt state)
	if taskQueue.Empty() {
		return nil, false, nil
	}
	if _, ok := pendPool[p.id]; ok {
		return nil, false, nil
	}
	// Calculate an upper limit on the items we might fetch (i.e. throttling)
	space := len(q.resultCache) - len(donePool)
	for _, request := range pendPool {
		space -= len(request.Headers)
	}
	// Retrieve a batch of tasks, skipping previously failed ones
	send := make([]*block.Header, 0, count)
	skip := make([]*block.Header, 0)

	progress := false
	for proc := 0; proc < space && len(send) < count && !taskQueue.Empty(); proc++ {
		header := taskQueue.PopItem().(*block.Header)

		// If we're the first to request this task, initialise the result container
		index := int(header.Number.IntVal.Int64() - int64(q.resultOffset))
		if index >= len(q.resultCache) || index < 0 {
			common.Report("index allocation went beyond available resultCache space")
			return nil, false, errInvalidChain
		}
		if q.resultCache[index] == nil {
			components := 2
			if q.mode == FastSync && header.Number.IntVal.Uint64() <= q.fastSyncPivot {
				components = 3
			}
			q.resultCache[index] = &fetchResult{
				Pending: components,
				Header:  header,
			}
		}
		// If this fetch task is a noop, skip this fetch operation
		if isNoop(header) {
			donePool[header.Hash()] = struct{}{}
			delete(taskPool, header.Hash())

			space, proc = space-1, proc-1
			q.resultCache[index].Pending--
			progress = true
			continue
		}
		// Otherwise unless the peer is known not to have the data, add to the retrieve list
		if p.Lacks(header.Hash()) {
			skip = append(skip, header)
		} else {
			send = append(send, header)
		}
	}
	// Merge all the skipped headers back
	for _, header := range skip {
		taskQueue.Push(header, -float32(header.Number.IntVal.Uint64()))
	}
	if progress {
		// Wake WaitResults, resultCache was modified
		q.active.Signal()
	}
	// Assemble and return the block download request
	if len(send) == 0 {
		return nil, progress, nil
	}
	request := &fetchRequest{
		Peer:    p,
		Headers: send,
		Time:    time.Now(),
	}
	pendPool[p.id] = request

	return request, progress, nil
}

// CancelHeaders aborts a fetch request, returning all pending skeleton indexes to the queue.
func (q *queue) CancelHeaders(request *fetchRequest) {
	q.cancel(request, q.headerTaskQueue, q.headerPendPool)
}

// CancelBodies aborts a body fetch request, returning all pending headers to the
// task queue.
func (q *queue) CancelBodies(request *fetchRequest) {
	q.cancel(request, q.blockTaskQueue, q.blockPendPool)
}

// CancelCertificate aborts a certificate fetch request, returning all pending headers to the
// task queue.
func (q *queue) CancelCertificate(request *fetchRequest) {
	q.cancel(request, q.certificateTaskQueue, q.certificatePendPool)
}

// CancelReceipts aborts a body fetch request, returning all pending headers to
// the task queue.
func (q *queue) CancelReceipts(request *fetchRequest) {
	q.cancel(request, q.receiptTaskQueue, q.receiptPendPool)
}

// Cancel aborts a fetch request, returning all pending hashes to the task queue.
func (q *queue) cancel(request *fetchRequest, taskQueue *prque.Prque, pendPool map[string]*fetchRequest) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if request.From > 0 {
		taskQueue.Push(request.From, -float32(request.From))
	}
	for hash, index := range request.Hashes {
		taskQueue.Push(hash, float32(index))
	}
	for _, header := range request.Headers {
		taskQueue.Push(header, -float32(header.Number.IntVal.Uint64()))
	}
	delete(pendPool, request.Peer.id)
}

// Revoke cancels all pending requests belonging to a given peer. This method is
// meant to be called during a peer drop to quickly reassign owned data fetches
// to remaining nodes.
func (q *queue) Revoke(peerId string) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if request, ok := q.blockPendPool[peerId]; ok {
		for _, header := range request.Headers {
			q.blockTaskQueue.Push(header, -float32(header.Number.IntVal.Uint64()))
		}
		delete(q.blockPendPool, peerId)
	}
	if request, ok := q.certificatePendPool[peerId]; ok {
		for _, header := range request.Headers {
			q.certificateTaskQueue.Push(header, -float32(header.Number.IntVal.Uint64()))
		}
		delete(q.certificatePendPool, peerId)
	}
	if request, ok := q.receiptPendPool[peerId]; ok {
		for _, header := range request.Headers {
			q.receiptTaskQueue.Push(header, -float32(header.Number.IntVal.Uint64()))
		}
		delete(q.receiptPendPool, peerId)
	}
}

// ExpireHeaders checks for in flight requests that exceeded a timeout allowance,
// canceling them and returning the responsible peers for penalisation.
func (q *queue) ExpireHeaders(timeout time.Duration) map[string]int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.expire(timeout, q.headerPendPool, q.headerTaskQueue, headerTimeoutMeter)
}

// ExpireBodies checks for in flight block body requests that exceeded a timeout
// allowance, canceling them and returning the responsible peers for penalisation.
func (q *queue) ExpireBodies(timeout time.Duration) map[string]int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.expire(timeout, q.blockPendPool, q.blockTaskQueue, bodyTimeoutMeter)
}

// ExpireCertificates checks for in flight block certificate requests that exceeded a timeout
// allowance, canceling them and returning the responsible peers for penalisation.
func (q *queue) ExpireCertificates(timeout time.Duration) map[string]int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.expire(timeout, q.certificatePendPool, q.certificateTaskQueue, certificateTimeoutMeter)
}

// ExpireReceipts checks for in flight receipt requests that exceeded a timeout
// allowance, canceling them and returning the responsible peers for penalisation.
func (q *queue) ExpireReceipts(timeout time.Duration) map[string]int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.expire(timeout, q.receiptPendPool, q.receiptTaskQueue, receiptTimeoutMeter)
}

// expire is the generic check that move expired tasks from a pending pool back
// into a task pool, returning all entities caught with expired tasks.
//
// Note, this method expects the queue lock to be already held. The
// reason the lock is not obtained in here is because the parameters already need
// to access the queue, so they already need a lock anyway.
func (q *queue) expire(timeout time.Duration, pendPool map[string]*fetchRequest, taskQueue *prque.Prque, timeoutMeter metrics.Meter) map[string]int {
	// Iterate over the expired requests and return each to the queue
	expiries := make(map[string]int)
	for id, request := range pendPool {
		if time.Since(request.Time) > timeout {
			// Update the metrics with the timeout
			timeoutMeter.Mark(1)

			// Return any non satisfied requests to the pool
			if request.From > 0 {
				taskQueue.Push(request.From, -float32(request.From))
			}
			for hash, index := range request.Hashes {
				taskQueue.Push(hash, float32(index))
			}
			for _, header := range request.Headers {
				taskQueue.Push(header, -float32(header.Number.IntVal.Uint64()))
			}
			// Add the peer to the expiry report along the the number of failed requests
			expirations := len(request.Hashes)
			if expirations < len(request.Headers) {
				expirations = len(request.Headers)
			}
			expiries[id] = expirations
		}
	}
	// Remove the expired requests from the pending pool
	for id := range expiries {
		delete(pendPool, id)
	}
	return expiries
}

// DeliverHeaders injects a header retrieval response into the header results
// cache. This method either accepts all headers it received, or none of them
// if they do not map correctly to the skeleton.
//
// If the headers are accepted, the method makes an attempt to deliver the set
// of ready headers to the processor to keep the pipeline full. However it will
// not block to prevent stalling other pending deliveries.
func (q *queue) DeliverHeaders(id string, headers []*block.Header, headerProcCh chan []*block.Header) (int, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Short circuit if the data was never requested
	request := q.headerPendPool[id]
	if request == nil {
		return 0, errNoFetchesPending
	}
	headerReqTimer.UpdateSince(request.Time)
	delete(q.headerPendPool, id)

	// Ensure headers can be mapped onto the skeleton chain
	target := q.headerTaskPool[request.From].Hash()

	accepted := len(headers) == MaxHeaderFetch
	if accepted {
		if headers[0].Number.IntVal.Uint64() != request.From {
			logger.Trace("First header broke chain ordering", "peer", id, "number", headers[0].Number.IntVal.String(), "hash", headers[0].Hash().String(), request.From)
			accepted = false
		} else if headers[len(headers)-1].Hash() != target {
			logger.Trace("Last header broke skeleton structure ", "peer", id, "number", headers[len(headers)-1].Number.IntVal.String(), "hash", headers[len(headers)-1].Hash().String(), "expected", target)
			accepted = false
		}
	}
	if accepted {
		for i, header := range headers[1:] {
			hash := header.Hash()
			if want := request.From + 1 + uint64(i); header.Number.IntVal.Uint64() != want {
				logger.Warn("Header broke chain ordering", "peer", id, "number", header.Number.IntVal.String(), "hash", hash.String(), "expected", want)
				accepted = false
				break
			}
			if headers[i].Hash() != header.ParentHash {
				logger.Warn("Header broke chain ancestry", "peer", id, "number", header.Number.IntVal.String(), "hash", hash.String())
				accepted = false
				break
			}
		}
	}
	// If the batch of headers wasn't accepted, mark as unavailable
	if !accepted {
		logger.Trace("Skeleton filling not accepted", "peer", id, "from", request.From)

		miss := q.headerPeerMiss[id]
		if miss == nil {
			q.headerPeerMiss[id] = make(map[uint64]struct{})
			miss = q.headerPeerMiss[id]
		}
		miss[request.From] = struct{}{}

		q.headerTaskQueue.Push(request.From, -float32(request.From))
		return 0, errors.New("delivery not accepted")
	}
	// Clean up a successful fetch and try to deliver any sub-results
	copy(q.headerResults[request.From-q.headerOffset:], headers)
	delete(q.headerTaskPool, request.From)

	ready := 0
	for q.headerProced+ready < len(q.headerResults) && q.headerResults[q.headerProced+ready] != nil {
		ready += MaxHeaderFetch
	}
	if ready > 0 {
		// Headers are ready for delivery, gather them and push forward (non blocking)
		process := make([]*block.Header, ready)
		copy(process, q.headerResults[q.headerProced:q.headerProced+ready])

		select {
		case headerProcCh <- process:
			logger.Trace("Pre-scheduled new headers", "peer", id, "count", len(process), "from", process[0].Number.IntVal.String())
			q.headerProced += len(process)
		default:
		}
	}
	// Check for termination and return
	if len(q.headerTaskPool) == 0 {
		q.headerContCh <- false
	}
	return len(headers), nil
}

// DeliverBodies injects a block body retrieval response into the results queue.
// The method returns the number of blocks bodies accepted from the delivery and
// also wakes any threads waiting for data delivery.
func (q *queue) DeliverBodies(id string, txLists [][]*transaction.Transaction) (int, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	reconstruct := func(header *block.Header, index int, result *fetchResult) error {
		if block.DeriveSha(transaction.Transactions(txLists[index])) != header.TxRootHash {
			return errInvalidBody
		}
		result.Transactions = txLists[index]
		return nil
	}
	return q.deliver(id, q.blockTaskPool, q.blockTaskQueue, q.blockPendPool, q.blockDonePool, bodyReqTimer, len(txLists), reconstruct)
}

// DeliverCertificates injects a block certificate retrieval response into the results queue.
// The method returns the number of blocks certificates accepted from the delivery and
// also wakes any threads waiting for data delivery.
func (q *queue) DeliverCertificates(id string, certificates [][]byte) (int, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	reconstruct := func(header *block.Header, index int, result *fetchResult) error {
		err := q.certificateCheck(header, certificates[index])
		if err != nil {
			return err
		}
		//if cert[0].ParentHash != header.ParentHash {
		//	return errInvalidCert
		//}
		result.Certificate = certificates[index]
		return nil
	}
	return q.deliver(id, q.certificateTaskPool, q.certificateTaskQueue, q.certificatePendPool, q.certificateDonePool, certificateReqTimer, len(certificates), reconstruct)
}

// DeliverReceipts injects a receipt retrieval response into the results queue.
// The method returns the number of transaction receipts accepted from the delivery
// and also wakes any threads waiting for data delivery.
func (q *queue) DeliverReceipts(id string, receiptList [][]*transaction.Receipt) (int, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	reconstruct := func(header *block.Header, index int, result *fetchResult) error {

		if block.DeriveSha(transaction.Receipts(receiptList[index])) != header.ReceiptRootHash {
			return errInvalidReceipt
		}
		result.Receipts = receiptList[index]
		return nil
	}
	return q.deliver(id, q.receiptTaskPool, q.receiptTaskQueue, q.receiptPendPool, q.receiptDonePool, receiptReqTimer, len(receiptList), reconstruct)
}

// deliver injects a data retrieval response into the results queue.
//
// Note, this method expects the queue lock to be already held for writing. The
// reason the lock is not obtained in here is because the parameters already need
// to access the queue, so they already need a lock anyway.
func (q *queue) deliver(id string, taskPool map[types.Hash]*block.Header, taskQueue *prque.Prque,
	pendPool map[string]*fetchRequest, donePool map[types.Hash]struct{}, reqTimer metrics.Timer,
	results int, reconstruct func(header *block.Header, index int, result *fetchResult) error) (int, error) {

	// Short circuit if the data was never requested
	request := pendPool[id]
	if request == nil {
		return 0, errNoFetchesPending
	}
	reqTimer.UpdateSince(request.Time)
	delete(pendPool, id)

	// If no data items were retrieved, mark them as unavailable for the origin peer
	if results == 0 {
		for _, header := range request.Headers {
			request.Peer.MarkLacking(header.Hash())
		}
	}
	// Assemble each of the results with their headers and retrieved data parts
	var (
		accepted int
		failure  error
		useful   bool
	)
	for i, header := range request.Headers {
		// Short circuit assembly if no more fetch results are found
		if i >= results {
			break
		}
		// Reconstruct the next result if contents match up
		index := int(header.Number.IntVal.Int64() - int64(q.resultOffset))
		if index >= len(q.resultCache) || index < 0 || q.resultCache[index] == nil {
			failure = errInvalidChain
			break
		}
		if err := reconstruct(header, i, q.resultCache[index]); err != nil {
			failure = err
			break
		}
		donePool[header.Hash()] = struct{}{}
		q.resultCache[index].Pending--
		useful = true
		accepted++

		// Clean up a successful fetch
		request.Headers[i] = nil
		delete(taskPool, header.Hash())
	}
	// Return all failed or missing fetches to the queue
	for _, header := range request.Headers {
		if header != nil {
			taskQueue.Push(header, -float32(header.Number.IntVal.Uint64()))
		}
	}
	// Wake up WaitResults
	if accepted > 0 {
		q.active.Signal()
	}
	// If none of the data was good, it's a stale delivery
	switch {
	case failure == nil || failure == errInvalidChain:
		return accepted, failure
	case useful:
		return accepted, fmt.Errorf("partial failure: %v", failure)
	default:
		return accepted, errStaleDelivery
	}
}

// Prepare configures the result cache to allow accepting and caching inbound
// fetch results.
func (q *queue) Prepare(offset uint64, mode SyncMode, pivot uint64, head *block.Header) {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Prepare the queue for sync results
	if q.resultOffset < offset {
		q.resultOffset = offset
	}
	q.fastSyncPivot = pivot
	q.mode = mode
}
