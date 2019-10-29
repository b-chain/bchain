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
// @File: handler.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package bchain

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"time"
	"bchain.io/core"
	"bchain.io/params"
	"bchain.io/node/services/bchain/downloader"
	"bchain.io/node/services/bchain/fetcher"
	"bchain.io/communication/p2p"
	"bchain.io/utils/event"
	"bchain.io/consensus"
	"bchain.io/communication/p2p/discover"
	"bchain.io/common"
	"bchain.io/core/blockchain"
	"bchain.io/utils/database"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/transaction"
	"bchain.io/common/types"
	"bchain.io/consensus/apos"
)

const (
	softResponseLimit = 2 * 1024 * 1024 // Target maximum size of returned blocks, headers or node data.
	estHeaderMsgpSize  = 500             // Approximate size of an MSGP encoded block header

	// txChanSize is the size of channel listening to TxPreEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096

	aposChanSize = 4096
)

var (
	daoChallengeTimeout = 15 * time.Second // Time allowance for a node to reply to the DAO handshake challenge
)

// errIncompatibleConfig is returned if the requested protocols and configs are
// not compatible (low protocol version restrictions and high requirements).
var errIncompatibleConfig = errors.New("incompatible configuration")


type ProtocolManager struct {
	networkId uint64

	fastSync  uint32 // Flag whether fast sync is enabled (gets disabled if we already have blocks)
	acceptTxs uint32 // Flag whether we're considered synchronised (enables transaction processing)

	txpool      txPool
	blockchain  *blockchain.BlockChain
	chaindb     database.IDatabase
	chainconfig *params.ChainConfig
	maxPeers    int

	downloader *downloader.Downloader
	fetcher    *fetcher.Fetcher
	peers      *peerSet

	SubProtocols []p2p.Protocol

	eventMux      *event.TypeMux
	txCh          chan core.TxPreEvent
	txSub         event.Subscription

	csCh          chan apos.CsEvent
	csSub         event.Subscription
	bpCh          chan apos.BpEvent
	bpSub         event.Subscription
	baCh          chan apos.BaEvent
	baSub         event.Subscription

	apos          *apos.Apos

	producedBlockSub *event.TypeMuxSubscription

	// channels for fetcher, syncer, txsyncLoop
	newPeerCh   chan *peer
	txsyncCh    chan *txsync
	quitSync    chan struct{}
	noMorePeers chan struct{}

	// wait group is used for graceful shutdowns during downloading
	// and processing
	wg sync.WaitGroup
}

// NewProtocolManager returns a new bchain sub protocol manager. The bchain sub protocol manages peers capable
// with the bchain network.
func NewProtocolManager(config *params.ChainConfig, mode downloader.SyncMode, networkId uint64, mux *event.TypeMux, txpool txPool, engine consensus.Engine, blockchain *blockchain.BlockChain, chaindb database.IDatabase, aPos *apos.Apos) (*ProtocolManager, error) {
	// Create the protocol manager with the base fields
	manager := &ProtocolManager{
		networkId:   networkId,
		eventMux:    mux,
		txpool:      txpool,
		blockchain:  blockchain,
		chaindb:     chaindb,
		chainconfig: config,
		peers:       newPeerSet(),
		newPeerCh:   make(chan *peer),
		noMorePeers: make(chan struct{}),
		txsyncCh:    make(chan *txsync),
		quitSync:    make(chan struct{}),
		apos    :    aPos,
	}
	// Figure out whether to allow fast sync or not
	if mode == downloader.FastSync && blockchain.CurrentBlock().NumberU64() > 0 {
		logger.Warn("Blockchain not empty, fast sync disabled")
		mode = downloader.FullSync
	}
	if mode == downloader.FastSync {
		manager.fastSync = uint32(1)
	}
	// Initiate a sub-protocol for every implemented version we can handle
	manager.SubProtocols = make([]p2p.Protocol, 0, len(ProtocolVersions))
	for i, version := range ProtocolVersions {
		// Skip protocol version if incompatible with the mode of operation
		if mode == downloader.FastSync && version < bchain63 {
			continue
		}
		// Compatible; initialise the sub-protocol
		version := version // Closure for the run
		manager.SubProtocols = append(manager.SubProtocols, p2p.Protocol{
			Name:    ProtocolName,
			Version: version,
			Length:  ProtocolLengths[i],
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				peer := manager.newPeer(int(version), p, rw)
				select {
				case manager.newPeerCh <- peer:
					manager.wg.Add(1)
					defer manager.wg.Done()
					return manager.handle(peer)
				case <-manager.quitSync:
					return p2p.DiscQuitting
				}
			},
			NodeInfo: func() interface{} {
				return manager.NodeInfo()
			},
			PeerInfo: func(id discover.NodeID) interface{} {
				if p := manager.peers.Peer(fmt.Sprintf("%x", id[:8])); p != nil {
					return p.Info()
				}
				return nil
			},
		})
	}
	if len(manager.SubProtocols) == 0 {
		return nil, errIncompatibleConfig
	}
	// Construct the different synchronisation mechanisms
	manager.downloader = downloader.New(mode, chaindb, manager.eventMux, blockchain, nil, manager.removePeer, apos.CheckCertificate)

	validator := func(header *block.Header) error {
		return engine.VerifyHeader(blockchain, header, true)
	}
	heighter := func() uint64 {
		return blockchain.CurrentBlock().NumberU64()
	}
	inserter := func(blocks block.Blocks) (int, error) {
		// If fast sync is running, deny importing weird blocks
		if atomic.LoadUint32(&manager.fastSync) == 1 {
			logger.Warn("Discarded bad propagated block", "number", blocks[0].Number(), "hash", blocks[0].Hash())
			return 0, nil
		}
		atomic.StoreUint32(&manager.acceptTxs, 1) // Mark initial sync done on any fetcher import
		return manager.blockchain.InsertChain(blocks)
	}
	manager.fetcher = fetcher.New(blockchain.GetBlockByHash, validator, manager.BroadcastBlock, heighter, inserter, manager.removePeer)

	return manager, nil
}

func (pm *ProtocolManager) removePeer(id string) {
	// Short circuit if the peer was already removed
	peer := pm.peers.Peer(id)
	if peer == nil {
		return
	}
	logger.Debug("Removing bchain peer", "peer", id)

	// Unregister the peer from the downloader and bchain peer set
	pm.downloader.UnregisterPeer(id)
	if err := pm.peers.Unregister(id); err != nil {
		logger.Error("Peer removal failed", "peer", id, "err", err)
	}
	// Hard disconnect at the networking layer
	if peer != nil {
		peer.Peer.Disconnect(p2p.DiscUselessPeer)
	}
}

func (pm *ProtocolManager) Start(maxPeers int) {
	pm.maxPeers = maxPeers

	pm.csCh = make(chan apos.CsEvent, aposChanSize)
	pm.csSub = apos.MsgTransfer().SubscribeCsEvent(pm.csCh)
	go pm.csBroadcastLoop()

	pm.bpCh = make(chan apos.BpEvent, aposChanSize)
	pm.bpSub = apos.MsgTransfer().SubscribeBpEvent(pm.bpCh)
	go pm.bpBroadcastLoop()

	pm.baCh = make(chan apos.BaEvent, aposChanSize)
	pm.baSub = apos.MsgTransfer().SubscribeBaEvent(pm.baCh)
	go pm.baBroadcastLoop()

	// broadcast transactions
	pm.txCh = make(chan core.TxPreEvent, txChanSize)
	pm.txSub = pm.txpool.SubscribeTxPreEvent(pm.txCh)
	go pm.txBroadcastLoop()

	// broadcast produced blocks
	pm.producedBlockSub = pm.eventMux.Subscribe(core.NewProducedBlockEvent{})
	go pm.producedBroadcastLoop()

	// start sync handlers
	go pm.syncer()
	go pm.txsyncLoop()
}

func (pm *ProtocolManager) Stop() {
	logger.Info("Stopping bchain protocol")

	pm.txSub.Unsubscribe()         // quits txBroadcastLoop
	pm.producedBlockSub.Unsubscribe() // quits blockBroadcastLoop

	pm.csSub.Unsubscribe()
	pm.bpSub.Unsubscribe()
	pm.baSub.Unsubscribe()

	// Quit the sync loop.
	// After this send has completed, no new peers will be accepted.
	pm.noMorePeers <- struct{}{}

	// Quit fetcher, txsyncLoop.
	close(pm.quitSync)

	// Disconnect existing sessions.
	// This also closes the gate for any new registrations on the peer set.
	// sessions which are already established but not added to pm.peers yet
	// will exit when they try to register.
	pm.peers.Close()

	// Wait for all peer handler goroutines and the loops to come down.
	pm.wg.Wait()

	logger.Info("bchain protocol stopped")
}

func (pm *ProtocolManager) newPeer(pv int, p *p2p.Peer, rw p2p.MsgReadWriter) *peer {
	return newPeer(pv, p, newMeteredMsgWriter(rw))
}

// handle is the callback invoked to manage the life cycle of an bchain peer. When
// this function terminates, the peer is disconnected.
func (pm *ProtocolManager) handle(p *peer) error {
	if pm.peers.Len() >= pm.maxPeers {
		return p2p.DiscTooManyPeers
	}
	logger.Debug("bchain peer connected", "name", p.Name())

	// Execute the bchain handshake
	number, head, genesis := pm.blockchain.Status()
	if err := p.Handshake(pm.networkId, number, head, genesis); err != nil {
		logger.Error("bchain handshake failed", "err", err)
		return err
	}
	logger.Debugf("Handshake info, block number(self: %d, peer: %d)" , number.Int64(), p.number.Int64())

	if rw, ok := p.rw.(*meteredMsgReadWriter); ok {
		rw.Init(p.version)
	}
	// Register the peer locally
	if err := pm.peers.Register(p); err != nil {
		logger.Error("bchain peer registration failed", "err", err)
		return err
	}
	defer pm.removePeer(p.id)

	// Register the peer in the downloader. If the downloader considers it banned, we disconnect
	if err := pm.downloader.RegisterPeer(p.id, p.version, p); err != nil {
		return err
	}
	// Propagate existing transactions. new transactions appearing
	// after this will be sent via broadcasts.
	pm.syncTransactions(p)

	// main loop. handle incoming messages.
	for {
		if err := pm.handleMsg(p); err != nil {
			logger.Info("bchain message handling failed", "err", err)
			return err
		}
	}
}

// handleMsg is invoked whenever an inbound message is received from a remote
// peer. The remote connection is torn down upon returning any error.
func (pm *ProtocolManager) handleMsg(p *peer) error {
	// Read the next message from the remote peer, and ensure it's fully consumed
	msg, err := p.rw.ReadMsg()
	if err != nil {
		return err
	}
	logger.Debugf("recive a msg: code = %v ", msg.Code)

	if msg.Size > ProtocolMaxMsgSize {
		return errResp(ErrMsgTooLarge, "%v > %v", msg.Size, ProtocolMaxMsgSize)
	}
	defer msg.Discard()

	// Handle the message depending on its contents
	switch {
	case msg.Code == StatusMsg:
		fmt.Println("msgCode == StatusMsg")
		// Status messages should never arrive after the handshake
		return errResp(ErrExtraStatusMsg, "uncontrolled status message")

	// Block header query, collect the requested headers and reply
	case msg.Code == GetBlockHeadersMsg:
		// Decode the complex header query
		var query GetBlockHeadersData
		if err := msg.Decode(&query); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		hashMode := query.Origin.Hash != (types.Hash{})

		// Gather headers until the fetch or network limits is reached
		var (
			bytes   common.StorageSize
			headers block.Headers
			unknown bool
		)
		for !unknown && len(headers.Headers) < int(query.Amount) && bytes < softResponseLimit && len(headers.Headers) < downloader.MaxHeaderFetch {
			// Retrieve the next header satisfying the query
			var origin *block.Header
			if hashMode {
				origin = pm.blockchain.GetHeaderByHash(query.Origin.Hash)
			} else {
				origin = pm.blockchain.GetHeaderByNumber(query.Origin.Number)
			}
			if origin == nil {
				break
			}
			number := origin.Number.IntVal.Uint64()
			headers.Headers = append(headers.Headers, origin)
			bytes += estHeaderMsgpSize

			// Advance to the next header of the query
			switch {
			case query.Origin.Hash != (types.Hash{}) && query.Reverse:
				// Hash based traversal towards the genesis block
				for i := 0; i < int(query.Skip)+1; i++ {
					if header := pm.blockchain.GetHeader(query.Origin.Hash, number); header != nil {
						query.Origin.Hash = header.ParentHash
						number--
					} else {
						unknown = true
						break
					}
				}
			case query.Origin.Hash != (types.Hash{}) && !query.Reverse:
				// Hash based traversal towards the leaf block
				var (
					current = origin.Number.IntVal.Uint64()
					next    = current + query.Skip + 1
				)
				if next <= current {
					infos, _ := json.MarshalIndent(p.Peer.Info(), "", "  ")
					logger.Warn("GetBlockHeaders skip overflow attack", "current", current, "skip", query.Skip, "next", next, "attacker", infos)
					unknown = true
				} else {
					if header := pm.blockchain.GetHeaderByNumber(next); header != nil {
						if pm.blockchain.GetBlockHashesFromHash(header.Hash(), query.Skip+1)[query.Skip] == query.Origin.Hash {
							query.Origin.Hash = header.Hash()
						} else {
							unknown = true
						}
					} else {
						unknown = true
					}
				}
			case query.Reverse:
				// Number based traversal towards the genesis block
				if query.Origin.Number >= query.Skip+1 {
					query.Origin.Number -= query.Skip + 1
				} else {
					unknown = true
				}

			case !query.Reverse:
				// Number based traversal towards the leaf block
				query.Origin.Number += query.Skip + 1
			}
		}
		return p.SendBlockHeaders(&headers)

	case msg.Code == BlockHeadersMsg:
		// A batch of headers arrived to one of our previous requests
		var headers block.Headers
		if err := msg.Decode(&headers); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}

		// Filter out any explicitly requested headers, deliver the rest to the downloader
		filter := len(headers.Headers) == 1
		if filter {
			// Irrelevant of the fork checks, send the header to the fetcher just in case
			headers.Headers = pm.fetcher.FilterHeaders(p.id, headers.Headers, time.Now())
		}
		if len(headers.Headers) > 0 || !filter {
			err := pm.downloader.DeliverHeaders(p.id, headers.Headers)
			if err != nil {
				logger.Debug("Failed to deliver headers", "err", err)
			}
		}

	case msg.Code == GetBlockBodiesMsg:
		// Decode the retrieval message
		var hashs types.Hashs
		if err := msg.Decode(&hashs); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Gather blocks until the fetch or network limits is reached
		var (
			bodies []*block.Body
		)
		for i:=0; i< len(hashs.Hashs)  && len(bodies) < downloader.MaxBlockFetch;i++ {
			hash := hashs.Hashs[i]
			// Retrieve the requested block body, stopping if enough was found
			if body := pm.blockchain.GetBody(*hash); body != nil {
				bodies = append(bodies, body)
			}
		}
		return p.SendBlockBodies(bodies)

	case msg.Code == BlockBodiesMsg:
		// A batch of block bodies arrived to one of our previous requests
		var request BlockBodiesData
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver them all to the downloader for queuing
		trasactions := make([][]*transaction.Transaction, len(request.Bodys))

		for i, body := range request.Bodys {
			trasactions[i] = body.Transactions
		}
		// Filter out any explicitly requested bodies, deliver the rest to the downloader
		filter := len(trasactions) > 0
		if filter {
			trasactions = pm.fetcher.FilterBodies(p.id, trasactions, time.Now())
		}
		if len(trasactions) > 0  || !filter {
			err := pm.downloader.DeliverBodies(p.id, trasactions)
			if err != nil {
				logger.Debug("Failed to deliver bodies", "err", err)
			}
		}
	case msg.Code == GetBlockCertificateMsg:
		// Decode the retrieval message
		var hashs types.Hashs
		if err := msg.Decode(&hashs); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Gather blocks until the fetch or network limits is reached
		var (
			certificates [][]byte
		)
		for i:=0; i< len(hashs.Hashs)  && len(certificates) < downloader.MaxBlockFetch;i++ {
			hash := hashs.Hashs[i]
			// Retrieve the requested block body, stopping if enough was found
			if certificate := pm.blockchain.GetConsensusData(*hash); certificate != nil {
				certificates = append(certificates, certificate)
			}
		}
		return p.SendBlockCertificates(certificates)
	case msg.Code == BlockCertificateMsg:
		// A batch of block bodies arrived to one of our previous requests
		var request BlockCertificateData
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}

		if err := pm.downloader.DeliverCertificates(p.id, request.Certificates); err != nil {
			logger.Debug("Failed to deliver certificate", "err", err)
		}

	case p.version >= bchain63 && msg.Code == GetNodeDataMsg:
		// Decode the retrieval message
		var hashs types.Hashs
		if err := msg.Decode(&hashs); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Gather state data until the fetch or network limits is reached
		var (
			bytes int
			data  NodeData
		)
		for i:=0; i<len(hashs.Hashs) && bytes < softResponseLimit && len(data.Nodes) < downloader.MaxStateFetch; i++ {
			hash := hashs.Hashs[i]
			// Retrieve the requested state entry, stopping if enough was found
			if entry, err := pm.chaindb.Get(hash.Bytes()); err == nil {
				data.Nodes = append(data.Nodes, entry)
				bytes += len(entry)
			}
		}
		return p.SendNodeData(&data)

	case p.version >= bchain63 && msg.Code == NodeDataMsg:
		// A batch of node state data arrived to one of our previous requests
		var data NodeData
		if err := msg.Decode(&data); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver all to the downloader
		if err := pm.downloader.DeliverNodeData(p.id, data.Nodes); err != nil {
			logger.Debug("Failed to deliver node state data", "err", err)
		}

	case p.version >= bchain63 && msg.Code == GetReceiptsMsg:
		// Decode the retrieval message
		var hashs types.Hashs
		if err := msg.Decode(&hashs); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Gather state data until the fetch or network limits is reached
		var (
			receipts transaction.Receipts_s
		)
		for i:=0;i< len(hashs.Hashs)  && len(receipts.Receipts_s) < downloader.MaxReceiptFetch;i++ {
			// Retrieve the hash of the next block
			hash := hashs.Hashs[i]
			// Retrieve the requested block's receipts, skipping if unknown to us
			results := blockchain.GetBlockReceipts(pm.chaindb, *hash, blockchain.GetBlockNumber(pm.chaindb, *hash))
			if results == nil {
				if header := pm.blockchain.GetHeaderByHash(*hash); header == nil || header.ReceiptRootHash != block.EmptyRootHash {
					continue
				}
			}
			receiptPs := []*transaction.ReceiptProtocol{}
			for _, receipt := range results {
				logPs := []*transaction.LogProtocol{}
				for _, log := range receipt.Logs {
					logP := &transaction.LogProtocol{log.Address,log.Topics,log.Data}
					logPs = append(logPs, logP)
				}
				loreceiptP := &transaction.ReceiptProtocol{receipt.Status,receipt.Bloom, logPs}
				receiptPs = append(receiptPs, loreceiptP)
			}
			// If known, encode and queue for response packet
			receipts.Receipts_s = append(receipts.Receipts_s, receiptPs)

		}
		return p.SendReceipts(&receipts)

	case p.version >= bchain63 && msg.Code == ReceiptsMsg:
		// A batch of receipts arrived to one of our previous requests
		var receipts transaction.Receipts_s
		if err := msg.Decode(&receipts); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}

		var receipts_s [][]*transaction.Receipt
		for _, rs := range receipts.Receipts_s {
			receipt_s := []*transaction.Receipt{}
			for _, r := range rs {
				logs := []*transaction.Log{}
				for _, l := range r.Logs {
					log:= &transaction.Log{Address:l.Address, Topics:l.Topics, Data:l.Data}
					logs = append(logs,log)
				}
				receipt := &transaction.Receipt{
					Status: r.Status,
					Bloom: r.Bloom,
					Logs: logs,
				}
				receipt_s = append(receipt_s, receipt)
			}

			receipts_s = append(receipts_s, receipt_s)
		}
		// Deliver all to the downloader
		if err := pm.downloader.DeliverReceipts(p.id, receipts_s); err != nil {
			logger.Debug("Failed to deliver receipts", "err", err)
		}

	case msg.Code == NewBlockHashesMsg:
		var announces NewBlockHashesData
		if err := msg.Decode(&announces); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		// Mark the hashes as present at the remote node
		for _, block := range announces {
			p.MarkBlock(block.Hash)
		}
		// Schedule all the unknown hashes for retrieval
		//unknown := make(NewBlockHashesData, 0, len(announces))
		//for _, block := range announces {
		//	if !pm.blockchain.HasBlock(block.Hash, block.Number) {
		//		unknown = append(unknown, block)
		//	}
		//}
		//for _, block := range unknown {
		//	pm.fetcher.Notify(p.id, block.Hash, block.Number, time.Now(), p.RequestOneHeader, p.RequestBodies)
		//}
		//todo, the safety will be discussed and tested
		if len(announces) == 1 {
			var (
				trueHead = announces[0].Hash
				trueNumber = big.NewInt(int64(announces[0].Number))
			)
			if _, number := p.Head(); trueNumber.Cmp(number) > 0 {
				p.SetHead(trueHead, trueNumber)
				currentBlock := pm.blockchain.CurrentBlock()
				if trueNumber.Cmp(currentBlock.Number()) > 0 {
					logger.Info("peer's block number is higher, try synchronise")
					go pm.synchronise(p)
				}
			}
		}

	case msg.Code == NewBlockMsg:
		// Retrieve and decode the propagated block
		var request NewBlockData
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		request.Block.ReceivedAt = msg.ReceivedAt
		request.Block.ReceivedFrom = p

		// Mark the peer as owning the block and schedule it for import
		p.MarkBlock(request.Block.Hash())
		pm.fetcher.Enqueue(p.id, request.Block)

		// Assuming the block is importable by the peer, but possibly not yet done so,
		// calculate the head hash and TD that the peer truly must have.
		var (
			trueHead = request.Block.ParentHash()
			//trueNumber   = request.Number.IntVal
			trueNumber = new(big.Int).Sub(&request.Number.IntVal, big.NewInt(1))
		)

		if _, number := p.Head(); trueNumber.Cmp(number) > 0 {
			p.SetHead(trueHead, trueNumber)

			// Schedule a sync if above ours. Note, this will not fire a sync for a gap of
			// a singe block (as the true TD is below the propagated block), however this
			// scenario should easily be covered by the fetcher.
			currentBlock := pm.blockchain.CurrentBlock()
			if trueNumber.Cmp(currentBlock.Number()) > 0 {
				go pm.synchronise(p)
			}
		}

	case msg.Code == TxMsg:
		// Transactions arrived, make sure we have a valid and fresh chain to handle them
		if atomic.LoadUint32(&pm.acceptTxs) == 0 {
			break
		}
		// Transactions can be processed, parse all of them and deliver to the pool
		var txs transaction.Transactions
		if err := msg.Decode(&txs); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		for i, tx := range txs {
			// Validate and mark the remote transaction
			if tx == nil {
				return errResp(ErrDecode, "transaction %d is nil", i)
			}
			p.MarkTransaction(tx.Hash())
		}
		pm.txpool.AddRemotes(txs)

	case msg.Code == CsMsg:
		var cs apos.CredentialSign
		if err := msg.Decode(&cs); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}

		p.MarkCredential(cs.CsHash())

		msgCs := apos.NewMsgCredential(&cs)
		msgCs.Send()
		msgCs.Close()
		// use bba step 1 msg for recovery trigger sync block
		if cs.Step == 1 && pm.apos.InRecovering() {
			var (
				trueHead = cs.ParentHash
				trueNumber =  big.NewInt(int64(cs.Round - 1))
			)
			if _, number := p.Head(); trueNumber.Cmp(number) >= 0 {
				p.SetHead(trueHead, trueNumber)
				currentBlock := pm.blockchain.CurrentBlock()
				if trueNumber.Cmp(currentBlock.Number()) > 0 {
					logger.Info("peer's block number(based bba step 1 credential) is higher, try synchronise")
					go pm.synchronise(p)
				} else if  trueNumber.Cmp(currentBlock.Number()) == 0 {
					if pm.aposRecoveryForking(trueHead, currentBlock, p) {
						logger.Info("network forking (based bba step 1 credential), try synchronise")
						go pm.synchronise(p)
					}
				}
			}
		}
	case msg.Code == BpMsg:
		var bp apos.BlockProposal
		if err := msg.Decode(&bp); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}

		p.MarkBlockProposal(bp.Block.Hash())

		msgBp := apos.NewMsgBlockProposal(&bp)
		msgBp.Send()
		msgBp.Close()
	case msg.Code == BaMsg:
		var ba apos.ByzantineAgreementStar
		if err := msg.Decode(&ba); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}

		p.MarkByzantineAgreementStar(ba.BaHash())

		msgBa := apos.NewMsgByzantineAgreementStar(&ba)
		msgBa.Send()
		msgBa.Close()

	default:
		return errResp(ErrInvalidMsgCode, "%v", msg.Code)
	}
	return nil
}

func (pm *ProtocolManager) aposRecoveryForking(trueHead types.Hash, currentBlock *block.Block, p *peer) bool {
	if trueHead != currentBlock.Hash() {
		lastBlock := pm.blockchain.GetBlock(currentBlock.ParentHash(), currentBlock.NumberU64()-1)
		if lastBlock != nil {
			if currentBlock.Time().Cmp(lastBlock.Time()) == 0 {
				// empty block
				return true
			} else {
				a := new(big.Int).SetBytes(trueHead.Bytes())
				b := new(big.Int).SetBytes(currentBlock.Hash().Bytes())
				hashCmpRet := a.Cmp(b)
				if hashCmpRet > 0 {
					return true
				}
			}
		}
	}
	return false
}

// BroadcastBlock will either propagate a block to a subset of it's peers, or
// will only announce it's availability (depending what's requested).
func (pm *ProtocolManager) BroadcastBlock(block *block.Block, propagate bool) {
	hash := block.Hash()
	peers := pm.peers.PeersWithoutBlock(hash)

	// If propagation is requested, send to a subset of the peer
	if propagate {
		// Calculate the TD of the block (it's not imported yet, so block.Td is not valid)
		var number *big.Int
		if parent := pm.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1); parent != nil {
			number = block.Number()
		} else {
			logger.Errorf("Propagating dangling block. number %d, hash %x", block.Number(), hash)
			return
		}
		// Send the block to a subset of our peers
		transfer := peers[:int(math.Sqrt(float64(len(peers))))]
		for _, peer := range transfer {
			peer.SendNewBlock(block, number)
		}
		logger.Tracef("Propagated block hash %x, recipiens %d, duration %d", hash, len(transfer), int64(time.Since(block.ReceivedAt)))
		return
	}
	// Otherwise if the block is indeed in out own chain, announce it
	if pm.blockchain.HasBlock(hash, block.NumberU64()) {
		for _, peer := range peers {
			peer.SendNewBlockHashes([]types.Hash{hash}, []uint64{block.NumberU64()})
		}
		logger.Tracef("Announced block hash %x, recipients %d, duration %d", hash, len(peers), int64(time.Since(block.ReceivedAt)))
	}
}

// BroadcastTx will propagate a transaction to all peers which are not known to
// already have the given transaction.
func (pm *ProtocolManager) BroadcastTx(hash types.Hash, tx *transaction.Transaction) {
	// Broadcast transaction to a batch of peers not knowing about it
	peers := pm.peers.PeersWithoutTx(hash)
	//FIXME include this again: peers = peers[:int(math.Sqrt(float64(len(peers))))]
	for _, peer := range peers {
		err := peer.SendTransactions(transaction.Transactions{tx})
		if err != nil {
			logger.Error("p2p send tx fail", err)
		}
	}
	logger.Tracef("Broadcast transaction. hash %x, recipients %v", hash, len(peers))
}


func (pm *ProtocolManager) BroadcastCs(hash types.Hash, cs *apos.CredentialSign) {
	// Broadcast Credential to a batch of peers not knowing about it
	peers := pm.peers.PeersWithoutCss(hash)

	for _, peer := range peers {
		err := peer.SendCredential(cs)
		if err != nil {
			logger.Error("p2p send cs fail", err)
		}
	}
	logger.Tracef("Broadcast credential. hash %x, recipients %v", hash, len(peers))
}

func (pm *ProtocolManager) BroadcastBp(hash types.Hash, bp *apos.BlockProposal) {
	// Broadcast BlockProposal to a batch of peers not knowing about it
	peers := pm.peers.PeersWithoutBps(hash)

	for _, peer := range peers {
		err := peer.SendBlockProposal(bp)
		if err != nil {
			logger.Error("p2p send bp fail", err)
		}
	}
	logger.Tracef("Broadcast BlockProposal. hash %x, recipients %v", hash, len(peers))
}

func (pm *ProtocolManager) BroadcastBa(hash types.Hash, ba *apos.ByzantineAgreementStar) {
	// Broadcast BinaryByzantineAgreement to a batch of peers not knowing about it
	peers := pm.peers.PeersWithoutBas(hash)

	for _, peer := range peers {
		err := peer.SendByzantineAgreementStar(ba)
		if err != nil {
			logger.Error("p2p send bba fail", err)
		}
	}
	logger.Tracef("Broadcast ByzantineAgreementStar. hash %x, recipients %v", hash, len(peers))
}

// Produced broadcast loop
func (self *ProtocolManager) producedBroadcastLoop() {
	// automatically stops if unsubscribe
	for obj := range self.producedBlockSub.Chan() {
		switch ev := obj.Data.(type) {
		case core.NewProducedBlockEvent:
			self.BroadcastBlock(ev.Block, true)  // First propagate block to peers
			self.BroadcastBlock(ev.Block, false) // Only then announce to the rest
		}
	}
}

func (self *ProtocolManager) txBroadcastLoop() {
	for {
		select {
		case event := <-self.txCh:
			self.BroadcastTx(event.Tx.Hash(), event.Tx)

		// Err() channel will be closed when unsubscribing.
		case <-self.txSub.Err():
			return
		}
	}
}

func (self *ProtocolManager) csBroadcastLoop() {
	for {
		select {
		case event := <-self.csCh:
			self.BroadcastCs(event.Cs.CsHash(), event.Cs)

		// Err() channel will be closed when unsubscribing.
		case <-self.csSub.Err():
			return
		}
	}
}

func (self *ProtocolManager) bpBroadcastLoop() {
	for {
		select {
		case event := <-self.bpCh:
			self.BroadcastBp(event.Bp.Block.Hash(), event.Bp)

		// Err() channel will be closed when unsubscribing.
		case <-self.bpSub.Err():
			return
		}
	}
}

func (self *ProtocolManager) baBroadcastLoop() {
	for {
		select {
		case event := <-self.baCh:
			self.BroadcastBa(event.Ba.BaHash(), event.Ba)

			// Err() channel will be closed when unsubscribing.
		case <-self.baSub.Err():
			return
		}
	}
}

// NodeInfo represents a short summary of the bchain sub-protocol metadata
// known about the host peer.
type NodeInfo struct {
	Network    uint64              `json:"network"`    // bchain network ID
	Number     *big.Int            `json:"difficulty"`     // height of the host's blockchain
	Genesis    types.Hash          `json:"genesis"`    // SHA3 hash of the host's genesis block
	Config     *params.ChainConfig `json:"config"`     // Chain configuration for the fork rules
	Head       types.Hash          `json:"head"`       // SHA3 hash of the host's best owned block
}

// NodeInfo retrieves some protocol metadata about the running host node.
func (self *ProtocolManager) NodeInfo() *NodeInfo {
	currentBlock := self.blockchain.CurrentBlock()
	return &NodeInfo{
		Network:    self.networkId,
		Number:     currentBlock.Number(),
		Genesis:    self.blockchain.Genesis().Hash(),
		Config:     self.blockchain.Config(),
		Head:       currentBlock.Hash(),
	}
}
