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
// @File: udp.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package discv5

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net"
	"time"

	"bchain.io/common/types"
	"bchain.io/utils/crypto"
	"bchain.io/communication/p2p/nat"
	"bchain.io/communication/p2p/netutil"

	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

const Version = 4

// Errors
var (
	errPacketTooSmall   = errors.New("too small")
	errBadHash          = errors.New("bad hash")
	errExpired          = errors.New("expired")
	errUnsolicitedReply = errors.New("unsolicited reply")
	errUnknownNode      = errors.New("unknown node")
	errTimeout          = errors.New("RPC timeout")
	errClockWarp        = errors.New("reply deadline too far in the future")
	errClosed           = errors.New("socket closed")
)

// Timeouts
const (
	respTimeout = 500 * time.Millisecond
	sendTimeout = 500 * time.Millisecond
	expiration  = 20 * time.Second

	ntpFailureThreshold = 32               // Continuous timeouts after which to check NTP
	ntpWarningCooldown  = 10 * time.Minute // Minimum amount of time to pass before repeating NTP warning
	driftThreshold      = 10 * time.Second // Allowed clock drift before warning user
)

// RPC request structures
type (
	Ping struct {
		Version    uint
		From, To   RpcEndpoint
		Expiration uint64

		// v5
		Topics []Topic

		// Ignore additional fields (for forward compatibility).
		Rest []byte `msg:"-"`
	}

	// pong is the reply to Ping.
	Pong struct {
		// This field should mirror the UDP envelope address
		// of the Ping packet, which provides a way to discover the
		// the external address (after NAT).
		To RpcEndpoint

		ReplyTok   []byte // This contains the hash of the Ping packet.
		Expiration uint64 // Absolute timestamp at which the packet becomes invalid.

		// v5
		TopicHash    types.Hash `msg:",extension"`
		TicketSerial uint32
		WaitPeriods  []uint32

		// Ignore additional fields (for forward compatibility).
		Rest []byte `msg:"-"`
	}

	// Findnode is a query for nodes close to the given target.
	Findnode struct {
		Target     NodeID `msg:",extension"` // doesn't need to be an actual public key
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []byte `msg:"-"`
	}

	// FindnodeHash is a query for nodes close to the given target.
	FindnodeHash struct {
		Target     types.Hash `msg:",extension"`
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []byte `msg:"-"`
	}

	// reply to Findnode
	Neighbors struct {
		Nodes      []RpcNode
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []byte `msg:"-"`
	}

	TopicRegister struct {
		Topics []Topic
		Idx    uint
		Pong   []byte
	}

	TopicQuery struct {
		Topic      Topic
		Expiration uint64
	}

	// reply to TopicQuery
	TopicNodes struct {
		Echo  types.Hash `msg:",extension"`
		Nodes []RpcNode
	}

	RpcNode struct {
		IP  types.IP `msg:",extension"` // len 4 for IPv4 or 16 for IPv6
		UDP uint16   // for discovery protocol
		TCP uint16   // for Msgpx protocol
		ID  NodeID
	}

	RpcEndpoint struct {
		IP  types.IP `msg:",extension"` // len 4 for IPv4 or 16 for IPv6
		UDP uint16   // for discovery protocol
		TCP uint16   // for Msgpx protocol
	}
)

const (
	macSize  = 256 / 8
	sigSize  = 520 / 8
	headSize = macSize + sigSize // space of packet frame data
)

// Neighbors replies are sent across multiple packets to
// stay below the 1280 byte limit. We compute the maximum number
// of entries by stuffing a packet until it grows too large.
var maxNeighbors = func() int {
	p := Neighbors{Expiration: ^uint64(0)}
	maxSizeNode := RpcNode{IP: *types.NewIP(net.IP(net.IPv4(0, 0, 0, 0))), UDP: ^uint16(0), TCP: ^uint16(0)}
	for n := 0; ; n++ {
		p.Nodes = append(p.Nodes, maxSizeNode)
		buf := bytes.Buffer{}
		err := msgp.Encode(&buf, &p)
		if err != nil {
			// If this ever happens, it will be caught by the unit tests.
			panic("cannot encode: " + err.Error())
		}
		if headSize+buf.Len()+1 >= 1280 {
			return n
		}
	}
}()

var maxTopicNodes = func() int {
	p := TopicNodes{}
	maxSizeNode := RpcNode{IP: *types.NewIP(net.IP(net.IPv4(0, 0, 0, 0))), UDP: ^uint16(0), TCP: ^uint16(0)}
	for n := 0; ; n++ {
		p.Nodes = append(p.Nodes, maxSizeNode)
		buf := bytes.Buffer{}
		err := msgp.Encode(&buf, &p)
		if err != nil {
			// If this ever happens, it will be caught by the unit tests.
			panic("cannot encode: " + err.Error())
		}
		if headSize+buf.Len()+1 >= 1280 {
			return n
		}
	}
}()

func makeEndpoint(addr *net.UDPAddr, tcpPort uint16) RpcEndpoint {
	ip := addr.IP.To4()
	if ip == nil {
		ip = addr.IP.To16()
	}
	return RpcEndpoint{IP: *types.NewIP(ip), UDP: uint16(addr.Port), TCP: tcpPort}
}

func (e1 RpcEndpoint) equal(e2 RpcEndpoint) bool {
	return e1.UDP == e2.UDP && e1.TCP == e2.TCP && e1.IP.Get().Equal(e2.IP.Get())
}

func nodeFromRPC(sender *net.UDPAddr, rn RpcNode) (*Node, error) {
	if err := netutil.CheckRelayIP(sender.IP, rn.IP.Get()); err != nil {
		return nil, err
	}
	n := NewNode(rn.ID, rn.IP.Get(), rn.UDP, rn.TCP)
	err := n.validateComplete()
	return n, err
}

func nodeToRPC(n *Node) RpcNode {
	return RpcNode{ID: n.ID, IP: n.IP, UDP: n.UDP, TCP: n.TCP}
}

type ingressPacket struct {
	remoteID   NodeID
	remoteAddr *net.UDPAddr
	ev         nodeEvent
	hash       []byte
	data       interface{} // one of the RPC structs
	rawData    []byte
}

type conn interface {
	ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	WriteToUDP(b []byte, addr *net.UDPAddr) (n int, err error)
	Close() error
	LocalAddr() net.Addr
}

// udp implements the RPC protocol.
type udp struct {
	conn        conn
	priv        *ecdsa.PrivateKey
	ourEndpoint RpcEndpoint
	nat         nat.Interface
	net         *Network
}

// ListenUDP returns a new table that listens for UDP packets on laddr.
func ListenUDP(priv *ecdsa.PrivateKey, laddr string, natm nat.Interface, nodeDBPath string, netrestrict *netutil.Netlist) (*Network, error) {
	transport, err := listenUDP(priv, laddr)
	if err != nil {
		return nil, err
	}
	net, err := newNetwork(transport, priv.PublicKey, natm, nodeDBPath, netrestrict)
	if err != nil {
		return nil, err
	}
	transport.net = net
	go transport.readLoop()
	return net, nil
}

func listenUDP(priv *ecdsa.PrivateKey, laddr string) (*udp, error) {
	addr, err := net.ResolveUDPAddr("udp", laddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return &udp{conn: conn, priv: priv, ourEndpoint: makeEndpoint(addr, uint16(addr.Port))}, nil
}

func (t *udp) localAddr() *net.UDPAddr {
	return t.conn.LocalAddr().(*net.UDPAddr)
}

func (t *udp) Close() {
	t.conn.Close()
}

func (t *udp) send(remote *Node, ptype nodeEvent, data interface{}) (hash []byte) {
	hash, _ = t.sendPacket(remote.ID, remote.addr(), byte(ptype), data)
	return hash
}

func (t *udp) sendPing(remote *Node, toaddr *net.UDPAddr, topics []Topic) (hash []byte) {
	hash, _ = t.sendPacket(remote.ID, toaddr, byte(pingPacket), Ping{
		Version:    Version,
		From:       t.ourEndpoint,
		To:         makeEndpoint(toaddr, uint16(toaddr.Port)), // TODO: maybe use known TCP port from DB
		Expiration: uint64(time.Now().Add(expiration).Unix()),
		Topics:     topics,
	})
	return hash
}

func (t *udp) sendFindnode(remote *Node, target NodeID) {
	t.sendPacket(remote.ID, remote.addr(), byte(findnodePacket), Findnode{
		Target:     target,
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	})
}

func (t *udp) sendNeighbours(remote *Node, results []*Node) {
	// Send Neighbors in chunks with at most maxNeighbors per packet
	// to stay below the 1280 byte limit.
	p := Neighbors{Expiration: uint64(time.Now().Add(expiration).Unix())}
	for i, result := range results {
		p.Nodes = append(p.Nodes, nodeToRPC(result))
		if len(p.Nodes) == maxNeighbors || i == len(results)-1 {
			t.sendPacket(remote.ID, remote.addr(), byte(neighborsPacket), p)
			p.Nodes = p.Nodes[:0]
		}
	}
}

func (t *udp) sendFindnodeHash(remote *Node, target types.Hash) {
	t.sendPacket(remote.ID, remote.addr(), byte(findnodeHashPacket), FindnodeHash{
		Target:     target,
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	})
}

func (t *udp) sendTopicRegister(remote *Node, topics []Topic, idx int, pong []byte) {
	t.sendPacket(remote.ID, remote.addr(), byte(topicRegisterPacket), TopicRegister{
		Topics: topics,
		Idx:    uint(idx),
		Pong:   pong,
	})
}

func (t *udp) sendTopicNodes(remote *Node, queryHash types.Hash, nodes []*Node) {
	p := TopicNodes{Echo: queryHash}
	if len(nodes) == 0 {
		t.sendPacket(remote.ID, remote.addr(), byte(topicNodesPacket), p)
		return
	}
	for i, result := range nodes {
		if netutil.CheckRelayIP(remote.IP.Get(), result.IP.Get()) != nil {
			continue
		}
		p.Nodes = append(p.Nodes, nodeToRPC(result))
		if len(p.Nodes) == maxTopicNodes || i == len(nodes)-1 {
			t.sendPacket(remote.ID, remote.addr(), byte(topicNodesPacket), p)
			p.Nodes = p.Nodes[:0]
		}
	}
}

func (t *udp) sendPacket(toid NodeID, toaddr *net.UDPAddr, ptype byte, req interface{}) (hash []byte, err error) {
	//fmt.Println("sendPacket", nodeEvent(ptype), toaddr.String(), toid.String())
	packet, hash, err := encodePacket(t.priv, ptype, req)
	if err != nil {
		//fmt.Println(err)
		return hash, err
	}
	logger.Trace(fmt.Sprintf(">>> %v to %x@%v", nodeEvent(ptype), toid[:8], toaddr))
	if _, err = t.conn.WriteToUDP(packet, toaddr); err != nil {
		logger.Trace(fmt.Sprint("UDP send failed:", err))
	}
	//fmt.Println(err)
	return hash, err
}

// zeroed padding space for encodePacket.
var headSpace = make([]byte, headSize)

func encodePacket(priv *ecdsa.PrivateKey, ptype byte, req interface{}) (p, hash []byte, err error) {
	reqVal, ok := req.(msgp.Encodable)
	if !ok {
		logger.Error("Can't convert to msgp.Encodable")
		return nil, nil, errors.New("Can't convert to msgp.Encodable.")
	}

	b := new(bytes.Buffer)
	b.Write(headSpace)
	b.WriteByte(ptype)

	if err := msgp.Encode(b, reqVal); err != nil {
		logger.Error("Can't encode discv5 packet.", "err", err)
		return nil, nil, err
	}
	packet := b.Bytes()
	sig, err := crypto.Sign(crypto.Keccak256(packet[headSize:]), priv)
	if err != nil {
		logger.Error("Can't sign discv5 packet.", "err", err)
		return nil, nil, err
	}
	copy(packet[macSize:], sig)
	// add the hash to the front. Note: this doesn't protect the
	// packet in any way.
	hash = crypto.Keccak256(packet[macSize:])
	copy(packet, hash)
	return packet, hash, nil
}

// readLoop runs in its own goroutine. it injects ingress UDP packets
// into the network loop.
func (t *udp) readLoop() {
	defer t.conn.Close()
	// Discovery packets are defined to be no larger than 1280 bytes.
	// Packets larger than this size will be cut at the end and treated
	// as invalid because their hash won't match.
	buf := make([]byte, 1280)
	for {
		nbytes, from, err := t.conn.ReadFromUDP(buf)
		if netutil.IsTemporaryError(err) {
			// Ignore temporary read errors.
			logger.Debug(fmt.Sprintf("Temporary read error: %v", err))
			continue
		} else if err != nil {
			// Shut down the loop for permament errors.
			logger.Debug(fmt.Sprintf("Read error: %v", err))
			return
		}
		t.handlePacket(from, buf[:nbytes])
	}
}

func (t *udp) handlePacket(from *net.UDPAddr, buf []byte) error {
	pkt := ingressPacket{remoteAddr: from}
	if err := decodePacket(buf, &pkt); err != nil {
		logger.Debug(fmt.Sprintf("Bad packet from %v: %v", from, err))
		//fmt.Println("bad packet", err)
		return err
	}
	t.net.reqReadPacket(pkt)
	return nil
}

func decodePacket(buffer []byte, pkt *ingressPacket) error {
	if len(buffer) < headSize+1 {
		return errPacketTooSmall
	}
	buf := make([]byte, len(buffer))
	copy(buf, buffer)
	hash, sig, sigdata := buf[:macSize], buf[macSize:headSize], buf[headSize:]
	shouldhash := crypto.Keccak256(buf[macSize:])
	if !bytes.Equal(hash, shouldhash) {
		return errBadHash
	}
	fromID, err := recoverNodeID(crypto.Keccak256(buf[headSize:]), sig)
	if err != nil {
		return err
	}
	pkt.rawData = buf
	pkt.hash = hash
	pkt.remoteID = fromID
	switch pkt.ev = nodeEvent(sigdata[0]); pkt.ev {
	case pingPacket:
		pkt.data = new(Ping)
	case pongPacket:
		pkt.data = new(Pong)
	case findnodePacket:
		pkt.data = new(Findnode)
	case neighborsPacket:
		pkt.data = new(Neighbors)
	case findnodeHashPacket:
		pkt.data = new(FindnodeHash)
	case topicRegisterPacket:
		pkt.data = new(TopicRegister)
	case topicQueryPacket:
		pkt.data = new(TopicQuery)
	case topicNodesPacket:
		pkt.data = new(TopicNodes)
	default:
		return fmt.Errorf("unknown packet type: %d", sigdata[0])
	}

	reqVal, ok := pkt.data.(msgp.Decodable)
	if !ok {
		logger.Error("Can't convert to msgp.Decodable")
		return fmt.Errorf("Can't convert to msgp.Decodable")
	}

	byteBuf := bytes.NewBuffer(sigdata[1:])
	err = msgp.Decode(byteBuf, reqVal)
	if err == nil {
		pkt.data = reqVal
	}

	return err
}
