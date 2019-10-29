package types

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *BigInt) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "bigint":
			{
				var zb0002 interface{}
				zb0002, err = dc.ReadIntf()
				if err != nil {
					return
				}
				z.IntVal = bigFromBytes(zb0002)
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
func (z BigInt) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "bigint"
	err = en.Append(0x81, 0xa6, 0x62, 0x69, 0x67, 0x69, 0x6e, 0x74)
	if err != nil {
		return
	}
	err = en.WriteIntf(bigToBytes(z.IntVal))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z BigInt) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "bigint"
	o = append(o, 0x81, 0xa6, 0x62, 0x69, 0x67, 0x69, 0x6e, 0x74)
	o, err = msgp.AppendIntf(o, bigToBytes(z.IntVal))
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BigInt) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "bigint":
			{
				var zb0002 interface{}
				zb0002, bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
				z.IntVal = bigFromBytes(zb0002)
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
func (z BigInt) Msgsize() (s int) {
	s = 1 + 7 + msgp.GuessSize(bigToBytes(z.IntVal))
	return
}
