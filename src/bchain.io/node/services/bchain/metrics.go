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
// @File: metrics.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package bchain

import (
	"bchain.io/communication/p2p"
	"bchain.io/utils/metrics"
)

var (
	propTxnInPacketsMeter     = metrics.NewRegisteredMeter("bchain/prop/txns/in/packets",nil)
	propTxnInTrafficMeter     = metrics.NewRegisteredMeter("bchain/prop/txns/in/traffic",nil)
	propTxnOutPacketsMeter    = metrics.NewRegisteredMeter("bchain/prop/txns/out/packets",nil)
	propTxnOutTrafficMeter    = metrics.NewRegisteredMeter("bchain/prop/txns/out/traffic",nil)
	propHashInPacketsMeter    = metrics.NewRegisteredMeter("bchain/prop/hashes/in/packets",nil)
	propHashInTrafficMeter    = metrics.NewRegisteredMeter("bchain/prop/hashes/in/traffic",nil)
	propHashOutPacketsMeter   = metrics.NewRegisteredMeter("bchain/prop/hashes/out/packets",nil)
	propHashOutTrafficMeter   = metrics.NewRegisteredMeter("bchain/prop/hashes/out/traffic",nil)
	propBlockInPacketsMeter   = metrics.NewRegisteredMeter("bchain/prop/blocks/in/packets",nil)
	propBlockInTrafficMeter   = metrics.NewRegisteredMeter("bchain/prop/blocks/in/traffic",nil)
	propBlockOutPacketsMeter  = metrics.NewRegisteredMeter("bchain/prop/blocks/out/packets",nil)
	propBlockOutTrafficMeter  = metrics.NewRegisteredMeter("bchain/prop/blocks/out/traffic",nil)
	reqHeaderInPacketsMeter   = metrics.NewRegisteredMeter("bchain/req/headers/in/packets",nil)
	reqHeaderInTrafficMeter   = metrics.NewRegisteredMeter("bchain/req/headers/in/traffic",nil)
	reqHeaderOutPacketsMeter  = metrics.NewRegisteredMeter("bchain/req/headers/out/packets",nil)
	reqHeaderOutTrafficMeter  = metrics.NewRegisteredMeter("bchain/req/headers/out/traffic",nil)
	reqBodyInPacketsMeter     = metrics.NewRegisteredMeter("bchain/req/bodies/in/packets",nil)
	reqBodyInTrafficMeter     = metrics.NewRegisteredMeter("bchain/req/bodies/in/traffic",nil)
	reqBodyOutPacketsMeter    = metrics.NewRegisteredMeter("bchain/req/bodies/out/packets",nil)
	reqBodyOutTrafficMeter    = metrics.NewRegisteredMeter("bchain/req/bodies/out/traffic",nil)
	reqStateInPacketsMeter    = metrics.NewRegisteredMeter("bchain/req/states/in/packets",nil)
	reqStateInTrafficMeter    = metrics.NewRegisteredMeter("bchain/req/states/in/traffic",nil)
	reqStateOutPacketsMeter   = metrics.NewRegisteredMeter("bchain/req/states/out/packets",nil)
	reqStateOutTrafficMeter   = metrics.NewRegisteredMeter("bchain/req/states/out/traffic",nil)
	reqReceiptInPacketsMeter  = metrics.NewRegisteredMeter("bchain/req/receipts/in/packets",nil)
	reqReceiptInTrafficMeter  = metrics.NewRegisteredMeter("bchain/req/receipts/in/traffic",nil)
	reqReceiptOutPacketsMeter = metrics.NewRegisteredMeter("bchain/req/receipts/out/packets",nil)
	reqReceiptOutTrafficMeter = metrics.NewRegisteredMeter("bchain/req/receipts/out/traffic",nil)
	miscInPacketsMeter        = metrics.NewRegisteredMeter("bchain/misc/in/packets",nil)
	miscInTrafficMeter        = metrics.NewRegisteredMeter("bchain/misc/in/traffic",nil)
	miscOutPacketsMeter       = metrics.NewRegisteredMeter("bchain/misc/out/packets",nil)
	miscOutTrafficMeter       = metrics.NewRegisteredMeter("bchain/misc/out/traffic",nil)
)

// meteredMsgReadWriter is a wrapper around a p2p.MsgReadWriter, capable of
// accumulating the above defined metrics based on the data stream contents.
type meteredMsgReadWriter struct {
	p2p.MsgReadWriter     // Wrapped message stream to meter
	version           int // Protocol version to select correct meters
}

// newMeteredMsgWriter wraps a p2p MsgReadWriter with metering support. If the
// metrics system is disabled, this function returns the original object.
func newMeteredMsgWriter(rw p2p.MsgReadWriter) p2p.MsgReadWriter {
	if !metrics.Enabled {
		return rw
	}
	return &meteredMsgReadWriter{MsgReadWriter: rw}
}

// Init sets the protocol version used by the stream to know which meters to
// increment in case of overlapping message ids between protocol versions.
func (rw *meteredMsgReadWriter) Init(version int) {
	rw.version = version
}

func (rw *meteredMsgReadWriter) ReadMsg() (p2p.Msg, error) {
	// Read the message and short circuit in case of an error
	msg, err := rw.MsgReadWriter.ReadMsg()
	if err != nil {
		return msg, err
	}
	// Account for the data traffic
	packets, traffic := miscInPacketsMeter, miscInTrafficMeter
	switch {
	case msg.Code == BlockHeadersMsg:
		packets, traffic = reqHeaderInPacketsMeter, reqHeaderInTrafficMeter
	case msg.Code == BlockBodiesMsg:
		packets, traffic = reqBodyInPacketsMeter, reqBodyInTrafficMeter

	case msg.Code == NodeDataMsg:
		packets, traffic = reqStateInPacketsMeter, reqStateInTrafficMeter
	case msg.Code == ReceiptsMsg:
		packets, traffic = reqReceiptInPacketsMeter, reqReceiptInTrafficMeter

	case msg.Code == NewBlockHashesMsg:
		packets, traffic = propHashInPacketsMeter, propHashInTrafficMeter
	case msg.Code == NewBlockMsg:
		packets, traffic = propBlockInPacketsMeter, propBlockInTrafficMeter
	case msg.Code == TxMsg:
		packets, traffic = propTxnInPacketsMeter, propTxnInTrafficMeter
	}
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

	return msg, err
}

func (rw *meteredMsgReadWriter) WriteMsg(msg p2p.Msg) error {
	// Account for the data traffic
	packets, traffic := miscOutPacketsMeter, miscOutTrafficMeter
	switch {
	case msg.Code == BlockHeadersMsg:
		packets, traffic = reqHeaderOutPacketsMeter, reqHeaderOutTrafficMeter
	case msg.Code == BlockBodiesMsg:
		packets, traffic = reqBodyOutPacketsMeter, reqBodyOutTrafficMeter

	case msg.Code == NodeDataMsg:
		packets, traffic = reqStateOutPacketsMeter, reqStateOutTrafficMeter
	case msg.Code == ReceiptsMsg:
		packets, traffic = reqReceiptOutPacketsMeter, reqReceiptOutTrafficMeter

	case msg.Code == NewBlockHashesMsg:
		packets, traffic = propHashOutPacketsMeter, propHashOutTrafficMeter
	case msg.Code == NewBlockMsg:
		packets, traffic = propBlockOutPacketsMeter, propBlockOutTrafficMeter
	case msg.Code == TxMsg:
		packets, traffic = propTxnOutPacketsMeter, propTxnOutTrafficMeter
	}
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

	// Send the packet to the p2p layer
	return rw.MsgReadWriter.WriteMsg(msg)
}
