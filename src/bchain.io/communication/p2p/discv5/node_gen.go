package discv5

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Node) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "IP":
			err = z.IP.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "UDP":
			z.UDP, err = dc.ReadUint16()
			if err != nil {
				return
			}
		case "TCP":
			z.TCP, err = dc.ReadUint16()
			if err != nil {
				return
			}
		case "ID":
			err = dc.ReadExactBytes((z.ID)[:])
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Node) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "IP"
	err = en.Append(0x84, 0xa2, 0x49, 0x50)
	if err != nil {
		return
	}
	err = z.IP.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "UDP"
	err = en.Append(0xa3, 0x55, 0x44, 0x50)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.UDP)
	if err != nil {
		return
	}
	// write "TCP"
	err = en.Append(0xa3, 0x54, 0x43, 0x50)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.TCP)
	if err != nil {
		return
	}
	// write "ID"
	err = en.Append(0xa2, 0x49, 0x44)
	if err != nil {
		return
	}
	err = en.WriteBytes((z.ID)[:])
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Node) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "IP"
	o = append(o, 0x84, 0xa2, 0x49, 0x50)
	o, err = z.IP.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "UDP"
	o = append(o, 0xa3, 0x55, 0x44, 0x50)
	o = msgp.AppendUint16(o, z.UDP)
	// string "TCP"
	o = append(o, 0xa3, 0x54, 0x43, 0x50)
	o = msgp.AppendUint16(o, z.TCP)
	// string "ID"
	o = append(o, 0xa2, 0x49, 0x44)
	o = msgp.AppendBytes(o, (z.ID)[:])
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Node) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "IP":
			bts, err = z.IP.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "UDP":
			z.UDP, bts, err = msgp.ReadUint16Bytes(bts)
			if err != nil {
				return
			}
		case "TCP":
			z.TCP, bts, err = msgp.ReadUint16Bytes(bts)
			if err != nil {
				return
			}
		case "ID":
			bts, err = msgp.ReadExactBytes(bts, (z.ID)[:])
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Node) Msgsize() (s int) {
	s = 1 + 3 + z.IP.Msgsize() + 4 + msgp.Uint16Size + 4 + msgp.Uint16Size + 3 + msgp.ArrayHeaderSize + (NodeIDBytes * (msgp.ByteSize))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *NodeID) DecodeMsg(dc *msgp.Reader) (err error) {
	err = dc.ReadExactBytes((z)[:])
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *NodeID) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteBytes((z)[:])
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *NodeID) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendBytes(o, (z)[:])
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *NodeID) UnmarshalMsg(bts []byte) (o []byte, err error) {
	bts, err = msgp.ReadExactBytes(bts, (z)[:])
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *NodeID) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize + (NodeIDBytes * (msgp.ByteSize))
	return
}
