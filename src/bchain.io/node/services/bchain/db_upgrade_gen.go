package bchain

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Meta) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "BlockHash":
			err = z.BlockHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "BlockIndex":
			z.BlockIndex, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Index":
			z.Index, err = dc.ReadUint64()
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
func (z *Meta) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "BlockHash"
	err = en.Append(0x83, 0xa9, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.BlockHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "BlockIndex"
	err = en.Append(0xaa, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x64, 0x65, 0x78)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.BlockIndex)
	if err != nil {
		return
	}
	// write "Index"
	err = en.Append(0xa5, 0x49, 0x6e, 0x64, 0x65, 0x78)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Index)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Meta) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "BlockHash"
	o = append(o, 0x83, 0xa9, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68)
	o, err = z.BlockHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "BlockIndex"
	o = append(o, 0xaa, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x64, 0x65, 0x78)
	o = msgp.AppendUint64(o, z.BlockIndex)
	// string "Index"
	o = append(o, 0xa5, 0x49, 0x6e, 0x64, 0x65, 0x78)
	o = msgp.AppendUint64(o, z.Index)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Meta) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "BlockHash":
			bts, err = z.BlockHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "BlockIndex":
			z.BlockIndex, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Index":
			z.Index, bts, err = msgp.ReadUint64Bytes(bts)
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
func (z *Meta) Msgsize() (s int) {
	s = 1 + 10 + z.BlockHash.Msgsize() + 11 + msgp.Uint64Size + 6 + msgp.Uint64Size
	return
}
