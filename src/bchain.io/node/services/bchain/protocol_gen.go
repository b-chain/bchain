package bchain

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
)

// DecodeMsg implements msgp.Decodable
func (z *BlockBodiesData) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Bodys":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Bodys) >= int(zb0002) {
				z.Bodys = (z.Bodys)[:zb0002]
			} else {
				z.Bodys = make([]*block.Body, zb0002)
			}
			for za0001 := range z.Bodys {
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						return
					}
					z.Bodys[za0001] = nil
				} else {
					if z.Bodys[za0001] == nil {
						z.Bodys[za0001] = new(block.Body)
					}
					err = z.Bodys[za0001].DecodeMsg(dc)
					if err != nil {
						return
					}
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
func (z *BlockBodiesData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Bodys"
	err = en.Append(0x81, 0xa5, 0x42, 0x6f, 0x64, 0x79, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Bodys)))
	if err != nil {
		return
	}
	for za0001 := range z.Bodys {
		if z.Bodys[za0001] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z.Bodys[za0001].EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *BlockBodiesData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Bodys"
	o = append(o, 0x81, 0xa5, 0x42, 0x6f, 0x64, 0x79, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Bodys)))
	for za0001 := range z.Bodys {
		if z.Bodys[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.Bodys[za0001].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BlockBodiesData) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Bodys":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Bodys) >= int(zb0002) {
				z.Bodys = (z.Bodys)[:zb0002]
			} else {
				z.Bodys = make([]*block.Body, zb0002)
			}
			for za0001 := range z.Bodys {
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					z.Bodys[za0001] = nil
				} else {
					if z.Bodys[za0001] == nil {
						z.Bodys[za0001] = new(block.Body)
					}
					bts, err = z.Bodys[za0001].UnmarshalMsg(bts)
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
func (z *BlockBodiesData) Msgsize() (s int) {
	s = 1 + 6 + msgp.ArrayHeaderSize
	for za0001 := range z.Bodys {
		if z.Bodys[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += z.Bodys[za0001].Msgsize()
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *BlockCertificateData) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Certificates":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Certificates) >= int(zb0002) {
				z.Certificates = (z.Certificates)[:zb0002]
			} else {
				z.Certificates = make([][]byte, zb0002)
			}
			for za0001 := range z.Certificates {
				z.Certificates[za0001], err = dc.ReadBytes(z.Certificates[za0001])
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
func (z *BlockCertificateData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Certificates"
	err = en.Append(0x81, 0xac, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Certificates)))
	if err != nil {
		return
	}
	for za0001 := range z.Certificates {
		err = en.WriteBytes(z.Certificates[za0001])
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *BlockCertificateData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Certificates"
	o = append(o, 0x81, 0xac, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Certificates)))
	for za0001 := range z.Certificates {
		o = msgp.AppendBytes(o, z.Certificates[za0001])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BlockCertificateData) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Certificates":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Certificates) >= int(zb0002) {
				z.Certificates = (z.Certificates)[:zb0002]
			} else {
				z.Certificates = make([][]byte, zb0002)
			}
			for za0001 := range z.Certificates {
				z.Certificates[za0001], bts, err = msgp.ReadBytesBytes(bts, z.Certificates[za0001])
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
func (z *BlockCertificateData) Msgsize() (s int) {
	s = 1 + 13 + msgp.ArrayHeaderSize
	for za0001 := range z.Certificates {
		s += msgp.BytesPrefixSize + len(z.Certificates[za0001])
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *GetBlockHeadersData) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Origin":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			for zb0002 > 0 {
				zb0002--
				field, err = dc.ReadMapKeyPtr()
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				case "Hash":
					err = z.Origin.Hash.DecodeMsg(dc)
					if err != nil {
						return
					}
				case "Number":
					z.Origin.Number, err = dc.ReadUint64()
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
		case "Amount":
			z.Amount, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Skip":
			z.Skip, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Reverse":
			z.Reverse, err = dc.ReadBool()
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
func (z *GetBlockHeadersData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "Origin"
	// map header, size 2
	// write "Hash"
	err = en.Append(0x84, 0xa6, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x82, 0xa4, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.Origin.Hash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Number"
	err = en.Append(0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Origin.Number)
	if err != nil {
		return
	}
	// write "Amount"
	err = en.Append(0xa6, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Amount)
	if err != nil {
		return
	}
	// write "Skip"
	err = en.Append(0xa4, 0x53, 0x6b, 0x69, 0x70)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Skip)
	if err != nil {
		return
	}
	// write "Reverse"
	err = en.Append(0xa7, 0x52, 0x65, 0x76, 0x65, 0x72, 0x73, 0x65)
	if err != nil {
		return
	}
	err = en.WriteBool(z.Reverse)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *GetBlockHeadersData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Origin"
	// map header, size 2
	// string "Hash"
	o = append(o, 0x84, 0xa6, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x82, 0xa4, 0x48, 0x61, 0x73, 0x68)
	o, err = z.Origin.Hash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Number"
	o = append(o, 0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	o = msgp.AppendUint64(o, z.Origin.Number)
	// string "Amount"
	o = append(o, 0xa6, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74)
	o = msgp.AppendUint64(o, z.Amount)
	// string "Skip"
	o = append(o, 0xa4, 0x53, 0x6b, 0x69, 0x70)
	o = msgp.AppendUint64(o, z.Skip)
	// string "Reverse"
	o = append(o, 0xa7, 0x52, 0x65, 0x76, 0x65, 0x72, 0x73, 0x65)
	o = msgp.AppendBool(o, z.Reverse)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *GetBlockHeadersData) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Origin":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			for zb0002 > 0 {
				zb0002--
				field, bts, err = msgp.ReadMapKeyZC(bts)
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				case "Hash":
					bts, err = z.Origin.Hash.UnmarshalMsg(bts)
					if err != nil {
						return
					}
				case "Number":
					z.Origin.Number, bts, err = msgp.ReadUint64Bytes(bts)
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
		case "Amount":
			z.Amount, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Skip":
			z.Skip, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Reverse":
			z.Reverse, bts, err = msgp.ReadBoolBytes(bts)
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
func (z *GetBlockHeadersData) Msgsize() (s int) {
	s = 1 + 7 + 1 + 5 + z.Origin.Hash.Msgsize() + 7 + msgp.Uint64Size + 7 + msgp.Uint64Size + 5 + msgp.Uint64Size + 8 + msgp.BoolSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *HashOrNumber) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Hash":
			err = z.Hash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Number":
			z.Number, err = dc.ReadUint64()
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
func (z *HashOrNumber) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "Hash"
	err = en.Append(0x82, 0xa4, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.Hash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Number"
	err = en.Append(0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Number)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *HashOrNumber) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Hash"
	o = append(o, 0x82, 0xa4, 0x48, 0x61, 0x73, 0x68)
	o, err = z.Hash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Number"
	o = append(o, 0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	o = msgp.AppendUint64(o, z.Number)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *HashOrNumber) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Hash":
			bts, err = z.Hash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Number":
			z.Number, bts, err = msgp.ReadUint64Bytes(bts)
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
func (z *HashOrNumber) Msgsize() (s int) {
	s = 1 + 5 + z.Hash.Msgsize() + 7 + msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *NewBlockData) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Block":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Block = nil
			} else {
				if z.Block == nil {
					z.Block = new(block.Block)
				}
				err = z.Block.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Number":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Number = nil
			} else {
				if z.Number == nil {
					z.Number = new(types.BigInt)
				}
				err = z.Number.DecodeMsg(dc)
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
func (z *NewBlockData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "Block"
	err = en.Append(0x82, 0xa5, 0x42, 0x6c, 0x6f, 0x63, 0x6b)
	if err != nil {
		return
	}
	if z.Block == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Block.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Number"
	err = en.Append(0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	if err != nil {
		return
	}
	if z.Number == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Number.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *NewBlockData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Block"
	o = append(o, 0x82, 0xa5, 0x42, 0x6c, 0x6f, 0x63, 0x6b)
	if z.Block == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Block.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Number"
	o = append(o, 0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	if z.Number == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Number.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *NewBlockData) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Block":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Block = nil
			} else {
				if z.Block == nil {
					z.Block = new(block.Block)
				}
				bts, err = z.Block.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Number":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Number = nil
			} else {
				if z.Number == nil {
					z.Number = new(types.BigInt)
				}
				bts, err = z.Number.UnmarshalMsg(bts)
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
func (z *NewBlockData) Msgsize() (s int) {
	s = 1 + 6
	if z.Block == nil {
		s += msgp.NilSize
	} else {
		s += z.Block.Msgsize()
	}
	s += 7
	if z.Number == nil {
		s += msgp.NilSize
	} else {
		s += z.Number.Msgsize()
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *NewBlockHashesData) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(NewBlockHashesData, zb0002)
	}
	for zb0001 := range *z {
		var field []byte
		_ = field
		var zb0003 uint32
		zb0003, err = dc.ReadMapHeader()
		if err != nil {
			return
		}
		for zb0003 > 0 {
			zb0003--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			switch msgp.UnsafeString(field) {
			case "Hash":
				err = (*z)[zb0001].Hash.DecodeMsg(dc)
				if err != nil {
					return
				}
			case "Number":
				(*z)[zb0001].Number, err = dc.ReadUint64()
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
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z NewBlockHashesData) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zb0004 := range z {
		// map header, size 2
		// write "Hash"
		err = en.Append(0x82, 0xa4, 0x48, 0x61, 0x73, 0x68)
		if err != nil {
			return
		}
		err = z[zb0004].Hash.EncodeMsg(en)
		if err != nil {
			return
		}
		// write "Number"
		err = en.Append(0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
		if err != nil {
			return
		}
		err = en.WriteUint64(z[zb0004].Number)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z NewBlockHashesData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zb0004 := range z {
		// map header, size 2
		// string "Hash"
		o = append(o, 0x82, 0xa4, 0x48, 0x61, 0x73, 0x68)
		o, err = z[zb0004].Hash.MarshalMsg(o)
		if err != nil {
			return
		}
		// string "Number"
		o = append(o, 0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
		o = msgp.AppendUint64(o, z[zb0004].Number)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *NewBlockHashesData) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(NewBlockHashesData, zb0002)
	}
	for zb0001 := range *z {
		var field []byte
		_ = field
		var zb0003 uint32
		zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
		for zb0003 > 0 {
			zb0003--
			field, bts, err = msgp.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			switch msgp.UnsafeString(field) {
			case "Hash":
				bts, err = (*z)[zb0001].Hash.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			case "Number":
				(*z)[zb0001].Number, bts, err = msgp.ReadUint64Bytes(bts)
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
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z NewBlockHashesData) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0004 := range z {
		s += 1 + 5 + z[zb0004].Hash.Msgsize() + 7 + msgp.Uint64Size
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *NodeData) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Nodes":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Nodes) >= int(zb0002) {
				z.Nodes = (z.Nodes)[:zb0002]
			} else {
				z.Nodes = make([][]byte, zb0002)
			}
			for za0001 := range z.Nodes {
				z.Nodes[za0001], err = dc.ReadBytes(z.Nodes[za0001])
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
func (z *NodeData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Nodes"
	err = en.Append(0x81, 0xa5, 0x4e, 0x6f, 0x64, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Nodes)))
	if err != nil {
		return
	}
	for za0001 := range z.Nodes {
		err = en.WriteBytes(z.Nodes[za0001])
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *NodeData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Nodes"
	o = append(o, 0x81, 0xa5, 0x4e, 0x6f, 0x64, 0x65, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Nodes)))
	for za0001 := range z.Nodes {
		o = msgp.AppendBytes(o, z.Nodes[za0001])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *NodeData) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Nodes":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Nodes) >= int(zb0002) {
				z.Nodes = (z.Nodes)[:zb0002]
			} else {
				z.Nodes = make([][]byte, zb0002)
			}
			for za0001 := range z.Nodes {
				z.Nodes[za0001], bts, err = msgp.ReadBytesBytes(bts, z.Nodes[za0001])
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
func (z *NodeData) Msgsize() (s int) {
	s = 1 + 6 + msgp.ArrayHeaderSize
	for za0001 := range z.Nodes {
		s += msgp.BytesPrefixSize + len(z.Nodes[za0001])
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *StatusData) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "ProtocolVersion":
			z.ProtocolVersion, err = dc.ReadUint32()
			if err != nil {
				return
			}
		case "NetworkId":
			z.NetworkId, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Number":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Number = nil
			} else {
				if z.Number == nil {
					z.Number = new(types.BigInt)
				}
				err = z.Number.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "CurrentBlock":
			err = z.CurrentBlock.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "GenesisBlock":
			err = z.GenesisBlock.DecodeMsg(dc)
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
func (z *StatusData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "ProtocolVersion"
	err = en.Append(0x85, 0xaf, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.ProtocolVersion)
	if err != nil {
		return
	}
	// write "NetworkId"
	err = en.Append(0xa9, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x64)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.NetworkId)
	if err != nil {
		return
	}
	// write "Number"
	err = en.Append(0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	if err != nil {
		return
	}
	if z.Number == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Number.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "CurrentBlock"
	err = en.Append(0xac, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b)
	if err != nil {
		return
	}
	err = z.CurrentBlock.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "GenesisBlock"
	err = en.Append(0xac, 0x47, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x42, 0x6c, 0x6f, 0x63, 0x6b)
	if err != nil {
		return
	}
	err = z.GenesisBlock.EncodeMsg(en)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *StatusData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "ProtocolVersion"
	o = append(o, 0x85, 0xaf, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint32(o, z.ProtocolVersion)
	// string "NetworkId"
	o = append(o, 0xa9, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x64)
	o = msgp.AppendUint64(o, z.NetworkId)
	// string "Number"
	o = append(o, 0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	if z.Number == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Number.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "CurrentBlock"
	o = append(o, 0xac, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b)
	o, err = z.CurrentBlock.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "GenesisBlock"
	o = append(o, 0xac, 0x47, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x42, 0x6c, 0x6f, 0x63, 0x6b)
	o, err = z.GenesisBlock.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *StatusData) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ProtocolVersion":
			z.ProtocolVersion, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				return
			}
		case "NetworkId":
			z.NetworkId, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Number":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Number = nil
			} else {
				if z.Number == nil {
					z.Number = new(types.BigInt)
				}
				bts, err = z.Number.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "CurrentBlock":
			bts, err = z.CurrentBlock.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "GenesisBlock":
			bts, err = z.GenesisBlock.UnmarshalMsg(bts)
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
func (z *StatusData) Msgsize() (s int) {
	s = 1 + 16 + msgp.Uint32Size + 10 + msgp.Uint64Size + 7
	if z.Number == nil {
		s += msgp.NilSize
	} else {
		s += z.Number.Msgsize()
	}
	s += 13 + z.CurrentBlock.Msgsize() + 13 + z.GenesisBlock.Msgsize()
	return
}
