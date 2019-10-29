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
// @File: peer.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package bchain

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"
	"gopkg.in/fatih/set.v0"
	"bchain.io/communication/p2p"
	"bchain.io/common/types"
	"bchain.io/core/transaction"
	"bchain.io/core/blockchain/block"
	"bchain.io/consensus/apos"
)

var (
	errClosed            = errors.New("peer set is closed")
	errAlreadyRegistered = errors.New("peer is already registered")
	errNotRegistered     = errors.New("peer is not registered")
)

const (
	maxKnownTxs      = 32768 // Maximum transactions hashes to keep in the known list (prevent DOS)
	maxKnownBlocks   = 1024  // Maximum block hashes to keep in the known list (prevent DOS)
	maxKnownApos     = 32768      // Maximum apos message hashes to keep in the known list (prevent DOS)
	handshakeTimeout = 5 * time.Second
)

func errResp(code errCode, format string, v ...interface{}) error {
	return fmt.Errorf("%v - %v", code, fmt.Sprintf(format, v...))
}

// PeerInfo represents a short summary of the bchain sub-protocol metadata known
// about a connected peer.
type PeerInfo struct {
	Version    int      `json:"version"`    // bchain protocol version negotiated
	Number     *big.Int `json:"number"`     // height of the peer's blockchain
	Head       string   `json:"head"`       // SHA3 hash of the peer's best owned block
}

type peer struct {
	id string

	*p2p.Peer
	rw p2p.MsgReadWriter

	version  int         // Protocol version negotiated
	forkDrop *time.Timer // Timed connection dropper if forks aren't validated in time

	head types.Hash
	number   *big.Int    // block chain height
	lock sync.RWMutex

	knownTxs    *set.Set // Set of transaction hashes known to be known by this peer
	knownBlocks *set.Set // Set of block hashes known to be known by this peer

	knownCss    *set.Set
	knownBps    *set.Set
	knownBas    *set.Set
}

func newPeer(version int, p *p2p.Peer, rw p2p.MsgReadWriter) *peer {
	id := p.ID()

	return &peer{
		Peer:        p,
		rw:          rw,
		version:     version,
		id:          fmt.Sprintf("%x", id[:8]),
		knownTxs:    set.New(),
		knownBlocks: set.New(),
		knownCss:    set.New(),
		knownBps:    set.New(),
		knownBas:    set.New(),
	}
}

// Info gathers and returns a collection of metadata known about a peer.
func (p *peer) Info() *PeerInfo {
	hash, number := p.Head()

	return &PeerInfo{
		Version:    p.version,
		Number:     number,
		Head:       hash.Hex(),
	}
}

// Head retrieves a copy of the current head hash and block number of the
// peer.
func (p *peer) Head() (hash types.Hash, number *big.Int) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	copy(hash[:], p.head[:])
	return hash, new(big.Int).Set(p.number)
}

// SetHead updates the head hash and total difficulty of the peer.
func (p *peer) SetHead(hash types.Hash, number *big.Int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	copy(p.head[:], hash[:])
	p.number.Set(number)
}

// MarkBlock marks a block as known for the peer, ensuring that the block will
// never be propagated to this particular peer.
func (p *peer) MarkBlock(hash types.Hash) {
	// If we reached the memory allowance, drop a previously known block hash
	for p.knownBlocks.Size() >= maxKnownBlocks {
		p.knownBlocks.Pop()
	}
	p.knownBlocks.Add(hash)
}

// MarkTransaction marks a transaction as known for the peer, ensuring that it
// will never be propagated to this particular peer.
func (p *peer) MarkTransaction(hash types.Hash) {
	// If we reached the memory allowance, drop a previously known transaction hash
	for p.knownTxs.Size() >= maxKnownTxs {
		p.knownTxs.Pop()
	}
	p.knownTxs.Add(hash)
}

func (p *peer) MarkCredential(hash types.Hash) {
	// If we reached the memory allowance, drop a previously known  hash
	for p.knownCss.Size() >= maxKnownApos {
		p.knownCss.Pop()
	}
	p.knownCss.Add(hash)
}

func (p *peer) MarkBlockProposal(hash types.Hash) {
	// If we reached the memory allowance, drop a previously known  hash
	for p.knownBps.Size() >= maxKnownApos {
		p.knownBps.Pop()
	}
	p.knownBps.Add(hash)
}

func (p *peer) MarkByzantineAgreementStar(hash types.Hash) {
	// If we reached the memory allowance, drop a previously known  hash
	for p.knownBas.Size() >= maxKnownApos {
		p.knownBas.Pop()
	}
	p.knownBas.Add(hash)
}

func (p *peer) SendCredential(cs *apos.CredentialSign) error {
	p.knownCss.Add(cs.CsHash())
	return p2p.Send(p.rw, CsMsg, cs)
}
func (p *peer) SendBlockProposal(bp *apos.BlockProposal) error {
	//p.knownBps.Add(bp.Block.Hash())
	return p2p.Send(p.rw, BpMsg, bp)
}
func (p *peer) SendByzantineAgreementStar(ba *apos.ByzantineAgreementStar) error {
	//p.knownBas.Add(ba.BaHash())
	return p2p.Send(p.rw, BaMsg, ba)
}

// SendTransactions sends transactions to the peer and includes the hashes
// in its transaction hash set for future reference.
func (p *peer) SendTransactions(txs transaction.Transactions) error {
	for _, tx := range txs {
		p.knownTxs.Add(tx.Hash())
	}
	return p2p.Send(p.rw, TxMsg, txs)
}

// SendNewBlockHashes announces the availability of a number of blocks through
// a hash notification.
func (p *peer) SendNewBlockHashes(hashes []types.Hash, numbers []uint64) error {
	for _, hash := range hashes {
		p.knownBlocks.Add(hash)
	}
	request := make(NewBlockHashesData, len(hashes))
	for i := 0; i < len(hashes); i++ {
		request[i].Hash = hashes[i]
		request[i].Number = numbers[i]
	}
	return p2p.Send(p.rw, NewBlockHashesMsg, request)
}

// SendNewBlock propagates an entire block to a remote peer.
func (p *peer) SendNewBlock(block *block.Block, number *big.Int) error {
	p.knownBlocks.Add(block.Hash())
	Number := types.NewBigInt(*number)
	return p2p.Send(p.rw, NewBlockMsg, &NewBlockData{block, Number})
}

// SendBlockHeaders sends a batch of block headers to the remote peer.
func (p *peer) SendBlockHeaders(headers *block.Headers) error {
	return p2p.Send(p.rw, BlockHeadersMsg, headers)
}

// SendBlockBodies sends a batch of block contents to the remote peer.
func (p *peer) SendBlockBodies(bodies []*block.Body) error {
	return p2p.Send(p.rw, BlockBodiesMsg, &BlockBodiesData{Bodys: bodies})
}

// SendBlockCertificates sends a batch of block certificates to the remote peer.
func (p *peer) SendBlockCertificates(certificates [][]byte) error {
	return p2p.Send(p.rw, BlockCertificateMsg, &BlockCertificateData{Certificates: certificates})
}

// SendBlockBodiesMSGP sends a batch of block contents to the remote peer from
// an already MSGP encoded format.
// todo
//func (p *peer) SendBlockBodiesMSGP(bodies [][]byte) error {
//	return p2p.SendEx(p.rw, BlockBodiesMsg, bodies)
//}


// SendNodeData sends a batch of arbitrary internal data, corresponding to the
// hashes requested.
func (p *peer) SendNodeData(data *NodeData) error {
	return p2p.Send(p.rw, NodeDataMsg, data)
}

// SendReceiptsMSGP sends a batch of transaction receipts, corresponding to the
// ones requested from an already MSPG encoded format.
func (p *peer) SendReceiptsMSGP(receipts [][]byte) error {
	return p2p.Send(p.rw, ReceiptsMsg, receipts)
}

func (p *peer) SendReceipts(receipts *transaction.Receipts_s) error {
	return p2p.Send(p.rw, ReceiptsMsg, receipts)
}

// RequestOneHeader is a wrapper around the header query functions to fetch a
// single header. It is used solely by the fetcher.
func (p *peer) RequestOneHeader(hash types.Hash) error {
	logger.Debug("Fetching single header", "hash", hash.String())
	return p2p.Send(p.rw, GetBlockHeadersMsg, &GetBlockHeadersData{Origin: HashOrNumber{Hash: hash}, Amount: uint64(1), Skip: uint64(0), Reverse: false})
}

// RequestHeadersByHash fetches a batch of blocks' headers corresponding to the
// specified header query, based on the hash of an origin block.
func (p *peer) RequestHeadersByHash(origin types.Hash, amount int, skip int, reverse bool) error {
	logger.Debug("Fetching batch of headers", "count", amount, "fromhash", origin, "skip", skip, "reverse", reverse)
	return p2p.Send(p.rw, GetBlockHeadersMsg, &GetBlockHeadersData{Origin: HashOrNumber{Hash: origin}, Amount: uint64(amount), Skip: uint64(skip), Reverse: reverse})
}

// RequestHeadersByNumber fetches a batch of blocks' headers corresponding to the
// specified header query, based on the number of an origin block.
func (p *peer) RequestHeadersByNumber(origin uint64, amount int, skip int, reverse bool) error {
	logger.Debug("Fetching batch of headers", "count", amount, "fromnum", origin, "skip", skip, "reverse", reverse)
	return p2p.Send(p.rw, GetBlockHeadersMsg, &GetBlockHeadersData{Origin: HashOrNumber{Number: origin}, Amount: uint64(amount), Skip: uint64(skip), Reverse: reverse})
}

// RequestBodies fetches a batch of blocks' bodies corresponding to the hashes
// specified.
func (p *peer) RequestBodies(hashes []types.Hash) error {
	logger.Debug("Fetching batch of block bodies", "count", len(hashes))
	var hashs_send types.Hashs
	for _, hash := range hashes{
		sHash := hash
		hashs_send.Hashs = append(hashs_send.Hashs, &sHash)
	}

	return p2p.Send(p.rw, GetBlockBodiesMsg, &hashs_send)
}

// RequestBodies fetches a batch of blocks' bodies corresponding to the hashes
// specified.
func (p *peer) RequestCertificates(hashes []types.Hash) error {
	logger.Debug("Fetching batch of block certificates", "count", len(hashes))
	var hashs_send types.Hashs
	for _, hash := range hashes{
		sHash := hash
		hashs_send.Hashs = append(hashs_send.Hashs, &sHash)
	}

	return p2p.Send(p.rw, GetBlockCertificateMsg, &hashs_send)
}

// RequestNodeData fetches a batch of arbitrary data from a node's known state
// data, corresponding to the specified hashes.
func (p *peer) RequestNodeData(hashes []types.Hash) error {
	logger.Debug("Fetching batch of state data", "count", len(hashes))
	return p2p.Send(p.rw, GetNodeDataMsg, hashes)
}

// RequestReceipts fetches a batch of transaction receipts from a remote node.
func (p *peer) RequestReceipts(hashes []types.Hash) error {
	logger.Debug("Fetching batch of receipts", "count", len(hashes))
	var hashs_send types.Hashs
	for _, hash := range hashes{
		sHash := hash
		hashs_send.Hashs = append(hashs_send.Hashs, &sHash)
	}
	return p2p.Send(p.rw, GetReceiptsMsg, &hashs_send)
}

// Handshake executes the bchain protocol handshake, negotiating version number,
// network IDs, difficulties, head and genesis blocks.
func (p *peer) Handshake(network uint64, number *big.Int, head types.Hash, genesis types.Hash) error {
	// Send out own handshake in a new thread
	errc := make(chan error, 2)
	var status StatusData // safe to read after two values have been received from errc

	go func() {
		errc <- p2p.Send(p.rw, StatusMsg, &StatusData{
			ProtocolVersion: uint32(p.version),
			NetworkId:       network,
			Number:          types.NewBigInt(*number),
			CurrentBlock:    head,
			GenesisBlock:    genesis,
		})
	}()
	go func() {
		errc <- p.readStatus(network, &status, genesis)
	}()
	timeout := time.NewTimer(handshakeTimeout)
	defer timeout.Stop()
	for i := 0; i < 2; i++ {
		select {
		case err := <-errc:
			if err != nil {
				return err
			}
		case <-timeout.C:
			return p2p.DiscReadTimeout
		}
	}
	p.number, p.head = &status.Number.IntVal, status.CurrentBlock
	return nil
}

func (p *peer) readStatus(network uint64, status *StatusData, genesis types.Hash) (err error) {
	msg, err := p.rw.ReadMsg()
	if err != nil {
		logger.Errorf("readStatus fail %v", err)
		return err
	}
	if msg.Code != StatusMsg {
		return errResp(ErrNoStatusMsg, "first msg has code %x (!= %x)", msg.Code, StatusMsg)
	}
	if msg.Size > ProtocolMaxMsgSize {
		return errResp(ErrMsgTooLarge, "%v > %v", msg.Size, ProtocolMaxMsgSize)
	}
	// Decode the handshake and make sure everything matches
	if err := msg.Decode(status); err != nil {
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	if status.GenesisBlock != genesis {
		return errResp(ErrGenesisBlockMismatch, "%x (!= %x)", status.GenesisBlock[:8], genesis[:8])
	}
	if status.NetworkId != network {
		return errResp(ErrNetworkIdMismatch, "%d (!= %d)", status.NetworkId, network)
	}
	if int(status.ProtocolVersion) != p.version {
		return errResp(ErrProtocolVersionMismatch, "%d (!= %d)", status.ProtocolVersion, p.version)
	}
	return nil
}

// String implements fmt.Stringer.
func (p *peer) String() string {
	return fmt.Sprintf("Peer %s [%s]", p.id,
		fmt.Sprintf("bchain/%2d", p.version),
	)
}

// peerSet represents the collection of active peers currently participating in
// the bchain sub-protocol.
type peerSet struct {
	peers  map[string]*peer
	lock   sync.RWMutex
	closed bool
}

// newPeerSet creates a new peer set to track the active participants.
func newPeerSet() *peerSet {
	return &peerSet{
		peers: make(map[string]*peer),
	}
}

// Register injects a new peer into the working set, or returns an error if the
// peer is already known.
func (ps *peerSet) Register(p *peer) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if ps.closed {
		return errClosed
	}
	if _, ok := ps.peers[p.id]; ok {
		return errAlreadyRegistered
	}
	ps.peers[p.id] = p
	return nil
}

// Unregister removes a remote peer from the active set, disabling any further
// actions to/from that particular entity.
func (ps *peerSet) Unregister(id string) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if _, ok := ps.peers[id]; !ok {
		return errNotRegistered
	}
	delete(ps.peers, id)
	return nil
}

// Peer retrieves the registered peer with the given id.
func (ps *peerSet) Peer(id string) *peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return ps.peers[id]
}

// Len returns if the current number of peers in the set.
func (ps *peerSet) Len() int {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return len(ps.peers)
}

// PeersWithoutBlock retrieves a list of peers that do not have a given block in
// their set of known hashes.
func (ps *peerSet) PeersWithoutBlock(hash types.Hash) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownBlocks.Has(hash) {
			list = append(list, p)
		}
	}
	return list
}

// PeersWithoutTx retrieves a list of peers that do not have a given transaction
// in their set of known hashes.
func (ps *peerSet) PeersWithoutTx(hash types.Hash) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownTxs.Has(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) PeersWithoutCss(hash types.Hash) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownCss.Has(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) PeersWithoutBps(hash types.Hash) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownBps.Has(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) PeersWithoutBas(hash types.Hash) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownBas.Has(hash) {
			list = append(list, p)
		}
	}
	return list
}

// BestPeer retrieves the known peer with the currently highest chain.
func (ps *peerSet) BestPeer() *peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	var (
		bestPeer *peer
		bestHeight   *big.Int
	)
	for _, p := range ps.peers {
		if _, number := p.Head(); bestPeer == nil || number.Cmp(bestHeight) > 0 {
			bestPeer, bestHeight = p, number
		}
	}
	return bestPeer
}

// Close disconnects all peers.
// No new peers can be registered after Close has returned.
func (ps *peerSet) Close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, p := range ps.peers {
		p.Disconnect(p2p.DiscQuitting)
	}
	ps.closed = true
}
