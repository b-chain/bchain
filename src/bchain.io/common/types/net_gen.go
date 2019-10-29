package types

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *IP) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "ip":
			{
				var zb0002 []byte
				zb0002, err = dc.ReadBytes([]byte(z.Ip))
				if err != nil {
					return
				}
				z.Ip = fromBytes(zb0002)
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
func (z *IP) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "ip"
	err = en.Append(0x81, 0xa2, 0x69, 0x70)
	if err != nil {
		return
	}
	err = en.WriteBytes(toBytes(z.Ip))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *IP) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "ip"
	o = append(o, 0x81, 0xa2, 0x69, 0x70)
	o = msgp.AppendBytes(o, toBytes(z.Ip))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *IP) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ip":
			{
				var zb0002 []byte
				zb0002, bts, err = msgp.ReadBytesBytes(bts, toBytes(z.Ip))
				if err != nil {
					return
				}
				z.Ip = fromBytes(zb0002)
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
func (z *IP) Msgsize() (s int) {
	s = 1 + 3 + msgp.BytesPrefixSize + len(toBytes(z.Ip))
	return
}
