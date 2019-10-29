package types

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Bloom) DecodeMsg(dc *msgp.Reader) (err error) {
	err = dc.ReadExactBytes((z)[:])
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Bloom) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteBytes((z)[:])
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Bloom) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendBytes(o, (z)[:])
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Bloom) UnmarshalMsg(bts []byte) (o []byte, err error) {
	bts, err = msgp.ReadExactBytes(bts, (z)[:])
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Bloom) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize + (BloomByteLength * (msgp.ByteSize))
	return
}
