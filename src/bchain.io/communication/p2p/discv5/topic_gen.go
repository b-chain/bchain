package discv5

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Topic) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 string
		zb0001, err = dc.ReadString()
		if err != nil {
			return
		}
		(*z) = Topic(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Topic) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteString(string(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Topic) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendString(o, string(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Topic) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 string
		zb0001, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			return
		}
		(*z) = Topic(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Topic) Msgsize() (s int) {
	s = msgp.StringPrefixSize + len(string(z))
	return
}
