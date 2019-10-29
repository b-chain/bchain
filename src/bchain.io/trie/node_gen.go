package trie

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *FullNode) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Children":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if zb0002 != uint32(17) {
				err = msgp.ArrayError{Wanted: uint32(17), Got: zb0002}
				return
			}
			for za0001 := range z.Children {
				{
					var zb0003 interface{}
					zb0003, err = dc.ReadIntf()
					if err != nil {
						return
					}
					z.Children[za0001], err = toNode(zb0003)
				}
				if err != nil {
					return
				}
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
func (z *FullNode) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Children"
	err = en.Append(0x81, 0xa8, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(17))
	if err != nil {
		return
	}
	for za0001 := range z.Children {
		var zb0001 interface{}
		zb0001, err = fromNode(z.Children[za0001])
		if err != nil {
			return
		}
		err = en.WriteIntf(zb0001)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FullNode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Children"
	o = append(o, 0x81, 0xa8, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e)
	o = msgp.AppendArrayHeader(o, uint32(17))
	for za0001 := range z.Children {
		var zb0001 interface{}
		zb0001, err = fromNode(z.Children[za0001])
		if err != nil {
			return
		}
		o, err = msgp.AppendIntf(o, zb0001)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FullNode) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Children":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if zb0002 != uint32(17) {
				err = msgp.ArrayError{Wanted: uint32(17), Got: zb0002}
				return
			}
			for za0001 := range z.Children {
				{
					var zb0003 interface{}
					zb0003, bts, err = msgp.ReadIntfBytes(bts)
					if err != nil {
						return
					}
					z.Children[za0001], err = toNode(zb0003)
					if err != nil {
						return
					}
				}
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
func (z *FullNode) Msgsize() (s int) {
	s = 1 + 9 + msgp.ArrayHeaderSize
	for za0001 := range z.Children {
		var zb0001 interface{}
		_ = z.Children[za0001]
		s += msgp.GuessSize(zb0001)
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *HashNode) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 []byte
		zb0001, err = dc.ReadBytes([]byte((*z)))
		if err != nil {
			return
		}
		(*z) = HashNode(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z HashNode) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteBytes([]byte(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z HashNode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendBytes(o, []byte(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *HashNode) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 []byte
		zb0001, bts, err = msgp.ReadBytesBytes(bts, []byte((*z)))
		if err != nil {
			return
		}
		(*z) = HashNode(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z HashNode) Msgsize() (s int) {
	s = msgp.BytesPrefixSize + len([]byte(z))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *NodeIntf) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 interface{}
		zb0001, err = dc.ReadIntf()
		if err != nil {
			return
		}
		(*z), err = toNode(zb0001)
	}
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z NodeIntf) EncodeMsg(en *msgp.Writer) (err error) {
	var zb0001 interface{}
	zb0001, err = fromNode(z)
	if err != nil {
		return
	}
	err = en.WriteIntf(zb0001)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z NodeIntf) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	var zb0001 interface{}
	zb0001, err = fromNode(z)
	if err != nil {
		return
	}
	o, err = msgp.AppendIntf(o, zb0001)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *NodeIntf) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 interface{}
		zb0001, bts, err = msgp.ReadIntfBytes(bts)
		if err != nil {
			return
		}
		(*z), err = toNode(zb0001)
		if err != nil {
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z NodeIntf) Msgsize() (s int) {
	var zb0001 interface{}
	_ = z
	s += msgp.GuessSize(zb0001)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ShortNode) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Key":
			z.Key, err = dc.ReadBytes(z.Key)
			if err != nil {
				return
			}
		case "Val":
			{
				var zb0002 interface{}
				zb0002, err = dc.ReadIntf()
				if err != nil {
					return
				}
				z.Val, err = toNode(zb0002)
			}
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
func (z *ShortNode) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "Key"
	err = en.Append(0x82, 0xa3, 0x4b, 0x65, 0x79)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Key)
	if err != nil {
		return
	}
	// write "Val"
	err = en.Append(0xa3, 0x56, 0x61, 0x6c)
	if err != nil {
		return
	}
	var zb0001 interface{}
	zb0001, err = fromNode(z.Val)
	if err != nil {
		return
	}
	err = en.WriteIntf(zb0001)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ShortNode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Key"
	o = append(o, 0x82, 0xa3, 0x4b, 0x65, 0x79)
	o = msgp.AppendBytes(o, z.Key)
	// string "Val"
	o = append(o, 0xa3, 0x56, 0x61, 0x6c)
	var zb0001 interface{}
	zb0001, err = fromNode(z.Val)
	if err != nil {
		return
	}
	o, err = msgp.AppendIntf(o, zb0001)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ShortNode) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Key":
			z.Key, bts, err = msgp.ReadBytesBytes(bts, z.Key)
			if err != nil {
				return
			}
		case "Val":
			{
				var zb0002 interface{}
				zb0002, bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
				z.Val, err = toNode(zb0002)
				if err != nil {
					return
				}
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
func (z *ShortNode) Msgsize() (s int) {
	s = 1 + 4 + msgp.BytesPrefixSize + len(z.Key) + 4
	var zb0001 interface{}
	_ = z.Val
	s += msgp.GuessSize(zb0001)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ValueNode) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 []byte
		zb0001, err = dc.ReadBytes([]byte((*z)))
		if err != nil {
			return
		}
		(*z) = ValueNode(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z ValueNode) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteBytes([]byte(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z ValueNode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendBytes(o, []byte(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ValueNode) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 []byte
		zb0001, bts, err = msgp.ReadBytesBytes(bts, []byte((*z)))
		if err != nil {
			return
		}
		(*z) = ValueNode(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ValueNode) Msgsize() (s int) {
	s = msgp.BytesPrefixSize + len([]byte(z))
	return
}
