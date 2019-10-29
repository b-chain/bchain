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
// @File: filter_system_test.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package filters

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"bchain.io/common/types"
	"bchain.io/communication/rpc"
	"bchain.io/consensus"
	"bchain.io/core"
	"bchain.io/core/blockchain"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/blockchain/chainmaker"
	"bchain.io/core/genesis"
	"bchain.io/core/transaction"
	"bchain.io/params"
	"bchain.io/utils/bloom"
	"bchain.io/utils/database"
	"bchain.io/utils/event"
)

type testBackend struct {
	mux        *event.TypeMux
	db         database.IDatabase
	sections   uint64
	txFeed     *event.Feed
	rmLogsFeed *event.Feed
	logsFeed   *event.Feed
	chainFeed  *event.Feed
}

func (b *testBackend) ChainDb() database.IDatabase {
	return b.db
}

func (b *testBackend) EventMux() *event.TypeMux {
	return b.mux
}

func (b *testBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*block.Header, error) {
	var hash types.Hash
	var num uint64
	if blockNr == rpc.LatestBlockNumber {
		hash = blockchain.GetHeadBlockHash(b.db)
		num = blockchain.GetBlockNumber(b.db, hash)
	} else {
		num = uint64(blockNr)
		hash = blockchain.GetCanonicalHash(b.db, num)
	}
	return blockchain.GetHeader(b.db, hash, num), nil
}

func (b *testBackend) GetReceipts(ctx context.Context, blockHash types.Hash) (transaction.Receipts, error) {
	num := blockchain.GetBlockNumber(b.db, blockHash)
	return blockchain.GetBlockReceipts(b.db, blockHash, num), nil
}

func (b *testBackend) SubscribeTxPreEvent(ch chan<- core.TxPreEvent) event.Subscription {
	return b.txFeed.Subscribe(ch)
}

func (b *testBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.rmLogsFeed.Subscribe(ch)
}

func (b *testBackend) SubscribeLogsEvent(ch chan<- []*transaction.Log) event.Subscription {
	return b.logsFeed.Subscribe(ch)
}

func (b *testBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.chainFeed.Subscribe(ch)
}

func (b *testBackend) BloomStatus() (uint64, uint64) {
	return params.BloomBitsBlocks, b.sections
}

func (b *testBackend) ServiceFilter(ctx context.Context, session *bloom.MatcherSession) {
	requests := make(chan chan *bloom.Retrieval)

	go session.Multiplex(16, 0, requests)
	go func() {
		for {
			// Wait for a service request or a shutdown
			select {
			case <-ctx.Done():
				return

			case request := <-requests:
				task := <-request

				task.Bitsets = make([][]byte, len(task.Sections))
				for i, section := range task.Sections {
					if rand.Int()%4 != 0 { // Handle occasional missing deliveries
						head := blockchain.GetCanonicalHash(b.db, (section+1)*params.BloomBitsBlocks-1)
						task.Bitsets[i], _ = blockchain.GetBloomBits(b.db, task.Bit, section, head)
					}
				}
				request <- task
			}
		}
	}()
}

// TestBlockSubscription tests if a block subscription returns block hashes for posted chain events.
// It creates multiple subscriptions:
// - one at the start and should receive all posted chain events and a second (blockHashes)
// - one that is created after a cutoff moment and uninstalled after a second cutoff moment (blockHashes[cutoff1:cutoff2])
// - one that is created after the second cutoff moment (blockHashes[cutoff2:])
func TestBlockSubscription(t *testing.T) {
	t.Parallel()

	var (
		mux         = new(event.TypeMux)
		db, _       = database.OpenMemDB()
		txFeed      = new(event.Feed)
		rmLogsFeed  = new(event.Feed)
		logsFeed    = new(event.Feed)
		chainFeed   = new(event.Feed)
		backend     = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed}
		api         = NewPublicFilterAPI(backend, false)
		genesis     = new(genesis.Genesis).MustCommit(db)
		chain, _    = chainmaker.GenerateChain(params.TestChainConfig, genesis, &consensus.Engine_empty{}, db, 10, func(i int, gen *chainmaker.BlockGen) {})
		chainEvents = []core.ChainEvent{}
	)

	for _, blk := range chain {
		chainEvents = append(chainEvents, core.ChainEvent{Hash: blk.Hash(), Block: blk})
	}

	chan0 := make(chan *block.Header)
	sub0 := api.events.SubscribeNewHeads(chan0)
	chan1 := make(chan *block.Header)
	sub1 := api.events.SubscribeNewHeads(chan1)

	go func() { // simulate client
		i1, i2 := 0, 0
		for i1 != len(chainEvents) || i2 != len(chainEvents) {
			select {
			case header := <-chan0:
				if chainEvents[i1].Hash != header.Hash() {
					t.Errorf("sub0 received invalid hash on index %d, want %x, got %x", i1, chainEvents[i1].Hash, header.Hash())
				}
				i1++
			case header := <-chan1:
				if chainEvents[i2].Hash != header.Hash() {
					t.Errorf("sub1 received invalid hash on index %d, want %x, got %x", i2, chainEvents[i2].Hash, header.Hash())
				}
				i2++
			}
		}

		sub0.Unsubscribe()
		sub1.Unsubscribe()
	}()

	time.Sleep(1 * time.Second)
	for _, e := range chainEvents {
		chainFeed.Send(e)
	}

	<-sub0.Err()
	<-sub1.Err()
}

// TestPendingTxFilter tests whether pending tx filters retrieve all pending transactions that are posted to the event mux.
func TestPendingTxFilter(t *testing.T) {
	t.Parallel()

	var (
		mux        = new(event.TypeMux)
		db, _      = database.OpenMemDB()
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed}
		api        = NewPublicFilterAPI(backend, false)

		actions      = transaction.Actions{&transaction.Action{Contract: types.Address{}, Params: []byte{}}}
		transactions = []*transaction.Transaction{
			transaction.NewTransaction(0, actions),
			transaction.NewTransaction(1, actions),
			transaction.NewTransaction(2, actions),
			transaction.NewTransaction(3, actions),
			transaction.NewTransaction(4, actions),
		}

		hashes []types.Hash
	)

	fid0 := api.NewPendingTransactionFilter()

	time.Sleep(1 * time.Second)
	for _, tx := range transactions {
		ev := core.TxPreEvent{Tx: tx}
		txFeed.Send(ev)
	}

	timeout := time.Now().Add(1 * time.Second)
	for {
		results, err := api.GetFilterChanges(fid0)
		if err != nil {
			t.Fatalf("Unable to retrieve logs: %v", err)
		}

		h := results.([]types.Hash)
		hashes = append(hashes, h...)
		if len(hashes) >= len(transactions) {
			break
		}
		// check timeout
		if time.Now().After(timeout) {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	if len(hashes) != len(transactions) {
		t.Errorf("invalid number of transactions, want %d transactions(s), got %d", len(transactions), len(hashes))
		return
	}
	for i := range hashes {
		if hashes[i] != transactions[i].Hash() {
			t.Errorf("hashes[%d] invalid, want %x, got %x", i, transactions[i].Hash(), hashes[i])
		}
	}
}

// TestLogFilterCreation test whether a given filter criteria makes sense.
// If not it must return an error.
func TestLogFilterCreation(t *testing.T) {
	var (
		mux        = new(event.TypeMux)
		db, _      = database.OpenMemDB()
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed}
		api        = NewPublicFilterAPI(backend, false)

		testCases = []struct {
			crit    FilterCriteria
			success bool
		}{
			// defaults
			{FilterCriteria{}, true},
			// valid block number range
			{FilterCriteria{FromBlock: big.NewInt(1), ToBlock: big.NewInt(2)}, true},
			// "produced" block range to pending
			{FilterCriteria{FromBlock: big.NewInt(1), ToBlock: big.NewInt(rpc.LatestBlockNumber.Int64())}, true},
			// new produced and pending blocks
			{FilterCriteria{FromBlock: big.NewInt(rpc.LatestBlockNumber.Int64()), ToBlock: big.NewInt(rpc.PendingBlockNumber.Int64())}, true},
			// from block "higher" than to block
			{FilterCriteria{FromBlock: big.NewInt(2), ToBlock: big.NewInt(1)}, false},
			// from block "higher" than to block
			{FilterCriteria{FromBlock: big.NewInt(rpc.LatestBlockNumber.Int64()), ToBlock: big.NewInt(100)}, false},
			// from block "higher" than to block
			{FilterCriteria{FromBlock: big.NewInt(rpc.PendingBlockNumber.Int64()), ToBlock: big.NewInt(100)}, false},
			// from block "higher" than to block
			{FilterCriteria{FromBlock: big.NewInt(rpc.PendingBlockNumber.Int64()), ToBlock: big.NewInt(rpc.LatestBlockNumber.Int64())}, false},
		}
	)

	for i, test := range testCases {
		_, err := api.NewFilter(test.crit)
		if test.success && err != nil {
			t.Errorf("expected filter creation for case %d to success, got %v", i, err)
		}
		if !test.success && err == nil {
			t.Errorf("expected testcase %d to fail with an error", i)
		}
	}
}

// TestInvalidLogFilterCreation tests whether invalid filter log criteria results in an error
// when the filter is created.
func TestInvalidLogFilterCreation(t *testing.T) {
	t.Parallel()

	var (
		mux        = new(event.TypeMux)
		db, _      = database.OpenMemDB()
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed}
		api        = NewPublicFilterAPI(backend, false)
	)

	// different situations where log filter creation should fail.
	// Reason: fromBlock > toBlock
	testCases := []FilterCriteria{
		0: {FromBlock: big.NewInt(rpc.PendingBlockNumber.Int64()), ToBlock: big.NewInt(rpc.LatestBlockNumber.Int64())},
		1: {FromBlock: big.NewInt(rpc.PendingBlockNumber.Int64()), ToBlock: big.NewInt(100)},
		2: {FromBlock: big.NewInt(rpc.LatestBlockNumber.Int64()), ToBlock: big.NewInt(100)},
	}

	for i, test := range testCases {
		if _, err := api.NewFilter(test); err == nil {
			t.Errorf("Expected NewFilter for case #%d to fail", i)
		}
	}
}

// TestLogFilter tests whether log filters match the correct logs that are posted to the event feed.
func TestLogFilter(t *testing.T) {
	t.Parallel()

	var (
		mux        = new(event.TypeMux)
		db, _      = database.OpenMemDB()
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed}
		api        = NewPublicFilterAPI(backend, false)

		firstAddr      = types.HexToAddress("0x1111111111111111111111111111111111111111")
		secondAddr     = types.HexToAddress("0x2222222222222222222222222222222222222222")
		thirdAddress   = types.HexToAddress("0x3333333333333333333333333333333333333333")
		notUsedAddress = types.HexToAddress("0x9999999999999999999999999999999999999999")
		firstTopic     = types.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
		secondTopic    = types.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222")
		notUsedTopic   = types.HexToHash("0x9999999999999999999999999999999999999999999999999999999999999999")

		// posted twice, once as vm.Logs and once as core.PendingLogsEvent
		allLogs = []*transaction.Log{
			{Address: firstAddr},
			{Address: firstAddr, Topics: []types.Hash{firstTopic}, BlockNumber: 1},
			{Address: secondAddr, Topics: []types.Hash{firstTopic}, BlockNumber: 1},
			{Address: thirdAddress, Topics: []types.Hash{secondTopic}, BlockNumber: 2},
			{Address: thirdAddress, Topics: []types.Hash{secondTopic}, BlockNumber: 3},
		}

		expectedCase7  = []*transaction.Log{allLogs[3], allLogs[4], allLogs[0], allLogs[1], allLogs[2], allLogs[3], allLogs[4]}
		expectedCase11 = []*transaction.Log{allLogs[1], allLogs[2], allLogs[1], allLogs[2]}

		testCases = []struct {
			crit     FilterCriteria
			expected []*transaction.Log
			id       rpc.ID
		}{
			// match all
			0: {FilterCriteria{}, allLogs, ""},
			// match none due to no matching addresses
			1: {FilterCriteria{Addresses: []types.Address{{}, notUsedAddress}, Topics: [][]types.Hash{nil}}, []*transaction.Log{}, ""},
			// match logs based on addresses, ignore topics
			2: {FilterCriteria{Addresses: []types.Address{firstAddr}}, allLogs[:2], ""},
			// match none due to no matching topics (match with address)
			3: {FilterCriteria{Addresses: []types.Address{secondAddr}, Topics: [][]types.Hash{{notUsedTopic}}}, []*transaction.Log{}, ""},
			// match logs based on addresses and topics
			4: {FilterCriteria{Addresses: []types.Address{thirdAddress}, Topics: [][]types.Hash{{firstTopic, secondTopic}}}, allLogs[3:5], ""},
			// match logs based on multiple addresses and "or" topics
			5: {FilterCriteria{Addresses: []types.Address{secondAddr, thirdAddress}, Topics: [][]types.Hash{{firstTopic, secondTopic}}}, allLogs[2:5], ""},
			// logs in the pending block
			6: {FilterCriteria{Addresses: []types.Address{firstAddr}, FromBlock: big.NewInt(rpc.PendingBlockNumber.Int64()), ToBlock: big.NewInt(rpc.PendingBlockNumber.Int64())}, allLogs[:2], ""},
			// produced logs with block num >= 2 or pending logs
			7: {FilterCriteria{FromBlock: big.NewInt(2), ToBlock: big.NewInt(rpc.PendingBlockNumber.Int64())}, expectedCase7, ""},
			// all "produced" logs with block num >= 2
			8: {FilterCriteria{FromBlock: big.NewInt(2), ToBlock: big.NewInt(rpc.LatestBlockNumber.Int64())}, allLogs[3:], ""},
			// all "produced" logs
			9: {FilterCriteria{ToBlock: big.NewInt(rpc.LatestBlockNumber.Int64())}, allLogs, ""},
			// all "produced" logs with 1>= block num <=2 and topic secondTopic
			10: {FilterCriteria{FromBlock: big.NewInt(1), ToBlock: big.NewInt(2), Topics: [][]types.Hash{{secondTopic}}}, allLogs[3:4], ""},
			// all "produced" and pending logs with topic firstTopic
			11: {FilterCriteria{FromBlock: big.NewInt(rpc.LatestBlockNumber.Int64()), ToBlock: big.NewInt(rpc.PendingBlockNumber.Int64()), Topics: [][]types.Hash{{firstTopic}}}, expectedCase11, ""},
			// match all logs due to wildcard topic
			12: {FilterCriteria{Topics: [][]types.Hash{nil}}, allLogs[1:], ""},
		}
	)

	// create all filters
	for i := range testCases {
		testCases[i].id, _ = api.NewFilter(testCases[i].crit)
	}

	// raise events
	time.Sleep(1 * time.Second)
	if nsend := logsFeed.Send(allLogs); nsend == 0 {
		t.Fatal("Shoud have at least one subscription")
	}
	if err := mux.Post(core.PendingLogsEvent{Logs: allLogs}); err != nil {
		t.Fatal(err)
	}

	for i, tt := range testCases {
		var fetched []*transaction.Log
		timeout := time.Now().Add(1 * time.Second)
		for { // fetch all expected logs
			results, err := api.GetFilterChanges(tt.id)
			if err != nil {
				t.Fatalf("Unable to fetch logs: %v", err)
			}

			fetched = append(fetched, results.([]*transaction.Log)...)
			if len(fetched) >= len(tt.expected) {
				break
			}
			// check timeout
			if time.Now().After(timeout) {
				break
			}

			time.Sleep(100 * time.Millisecond)
		}

		if len(fetched) != len(tt.expected) {
			t.Errorf("invalid number of logs for case %d, want %d log(s), got %d", i, len(tt.expected), len(fetched))
			return
		}

		for l := range fetched {
			if fetched[l].Removed {
				t.Errorf("expected log not to be removed for log %d in case %d", l, i)
			}
			if !reflect.DeepEqual(fetched[l], tt.expected[l]) {
				t.Errorf("invalid log on index %d for case %d", l, i)
			}
		}
	}
}

// TestPendingLogsSubscription tests if a subscription receives the correct pending logs that are posted to the event feed.
func TestPendingLogsSubscription(t *testing.T) {
	t.Parallel()

	var (
		mux        = new(event.TypeMux)
		db, _      = database.OpenMemDB()
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed}
		api        = NewPublicFilterAPI(backend, false)

		firstAddr      = types.HexToAddress("0x1111111111111111111111111111111111111111")
		secondAddr     = types.HexToAddress("0x2222222222222222222222222222222222222222")
		thirdAddress   = types.HexToAddress("0x3333333333333333333333333333333333333333")
		notUsedAddress = types.HexToAddress("0x9999999999999999999999999999999999999999")
		firstTopic     = types.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
		secondTopic    = types.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222")
		thirdTopic     = types.HexToHash("0x3333333333333333333333333333333333333333333333333333333333333333")
		fourthTopic    = types.HexToHash("0x4444444444444444444444444444444444444444444444444444444444444444")
		notUsedTopic   = types.HexToHash("0x9999999999999999999999999999999999999999999999999999999999999999")

		allLogs = []core.PendingLogsEvent{
			{Logs: []*transaction.Log{{Address: firstAddr, Topics: []types.Hash{}, BlockNumber: 0}}},
			{Logs: []*transaction.Log{{Address: firstAddr, Topics: []types.Hash{firstTopic}, BlockNumber: 1}}},
			{Logs: []*transaction.Log{{Address: secondAddr, Topics: []types.Hash{firstTopic}, BlockNumber: 2}}},
			{Logs: []*transaction.Log{{Address: thirdAddress, Topics: []types.Hash{secondTopic}, BlockNumber: 3}}},
			{Logs: []*transaction.Log{{Address: thirdAddress, Topics: []types.Hash{secondTopic}, BlockNumber: 4}}},
			{Logs: []*transaction.Log{
				{Address: thirdAddress, Topics: []types.Hash{firstTopic}, BlockNumber: 5},
				{Address: thirdAddress, Topics: []types.Hash{thirdTopic}, BlockNumber: 5},
				{Address: thirdAddress, Topics: []types.Hash{fourthTopic}, BlockNumber: 5},
				{Address: firstAddr, Topics: []types.Hash{firstTopic}, BlockNumber: 5},
			}},
		}

		convertLogs = func(pl []core.PendingLogsEvent) []*transaction.Log {
			var logs []*transaction.Log
			for _, l := range pl {
				logs = append(logs, l.Logs...)
			}
			return logs
		}

		testCases = []struct {
			crit     FilterCriteria
			expected []*transaction.Log
			c        chan []*transaction.Log
			sub      *Subscription
		}{
			// match all
			{FilterCriteria{}, convertLogs(allLogs), nil, nil},
			// match none due to no matching addresses
			{FilterCriteria{Addresses: []types.Address{{}, notUsedAddress}, Topics: [][]types.Hash{nil}}, []*transaction.Log{}, nil, nil},
			// match logs based on addresses, ignore topics
			{FilterCriteria{Addresses: []types.Address{firstAddr}}, append(convertLogs(allLogs[:2]), allLogs[5].Logs[3]), nil, nil},
			// match none due to no matching topics (match with address)
			{FilterCriteria{Addresses: []types.Address{secondAddr}, Topics: [][]types.Hash{{notUsedTopic}}}, []*transaction.Log{}, nil, nil},
			// match logs based on addresses and topics
			{FilterCriteria{Addresses: []types.Address{thirdAddress}, Topics: [][]types.Hash{{firstTopic, secondTopic}}}, append(convertLogs(allLogs[3:5]), allLogs[5].Logs[0]), nil, nil},
			// match logs based on multiple addresses and "or" topics
			{FilterCriteria{Addresses: []types.Address{secondAddr, thirdAddress}, Topics: [][]types.Hash{{firstTopic, secondTopic}}}, append(convertLogs(allLogs[2:5]), allLogs[5].Logs[0]), nil, nil},
			// block numbers are ignored for filters created with New***Filter, these return all logs that match the given criteria when the state changes
			{FilterCriteria{Addresses: []types.Address{firstAddr}, FromBlock: big.NewInt(2), ToBlock: big.NewInt(3)}, append(convertLogs(allLogs[:2]), allLogs[5].Logs[3]), nil, nil},
			// multiple pending logs, should match only 2 topics from the logs in block 5
			{FilterCriteria{Addresses: []types.Address{thirdAddress}, Topics: [][]types.Hash{{firstTopic, fourthTopic}}}, []*transaction.Log{allLogs[5].Logs[0], allLogs[5].Logs[2]}, nil, nil},
		}
	)

	// create all subscriptions, this ensures all subscriptions are created before the events are posted.
	// on slow machines this could otherwise lead to missing events when the subscription is created after
	// (some) events are posted.
	for i := range testCases {
		testCases[i].c = make(chan []*transaction.Log)
		testCases[i].sub, _ = api.events.SubscribeLogs(testCases[i].crit, testCases[i].c)
	}

	for n, test := range testCases {
		i := n
		tt := test
		go func() {
			var fetched []*transaction.Log
		fetchLoop:
			for {
				logs := <-tt.c
				fetched = append(fetched, logs...)
				if len(fetched) >= len(tt.expected) {
					break fetchLoop
				}
			}

			if len(fetched) != len(tt.expected) {
				panic(fmt.Sprintf("invalid number of logs for case %d, want %d log(s), got %d", i, len(tt.expected), len(fetched)))
			}

			for l := range fetched {
				if fetched[l].Removed {
					panic(fmt.Sprintf("expected log not to be removed for log %d in case %d", l, i))
				}
				if !reflect.DeepEqual(fetched[l], tt.expected[l]) {
					panic(fmt.Sprintf("invalid log on index %d for case %d", l, i))
				}
			}
		}()
	}

	// raise events
	time.Sleep(1 * time.Second)
	// allLogs are type of core.PendingLogsEvent
	for _, l := range allLogs {
		if err := mux.Post(l); err != nil {
			t.Fatal(err)
		}
	}
}
