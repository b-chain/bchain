package blockchain

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	msgp "github.com/tinylib/msgp/msgp"
	"bchain.io/common/types"
)

// DecodeMsg implements msgp.Decodable
func (z *BlockStat) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Ttxs":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Ttxs = nil
			} else {
				if z.Ttxs == nil {
					z.Ttxs = new(types.BigInt)
				}
				err = z.Ttxs.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "TsoNormal":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.TsoNormal = nil
			} else {
				if z.TsoNormal == nil {
					z.TsoNormal = new(types.BigInt)
				}
				err = z.TsoNormal.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "TsoContract":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.TsoContract = nil
			} else {
				if z.TsoContract == nil {
					z.TsoContract = new(types.BigInt)
				}
				err = z.TsoContract.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "TstateNum":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.TstateNum = nil
			} else {
				if z.TstateNum == nil {
					z.TstateNum = new(types.BigInt)
				}
				err = z.TstateNum.DecodeMsg(dc)
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
func (z *BlockStat) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "Ttxs"
	err = en.Append(0x84, 0xa4, 0x54, 0x74, 0x78, 0x73)
	if err != nil {
		return
	}
	if z.Ttxs == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Ttxs.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "TsoNormal"
	err = en.Append(0xa9, 0x54, 0x73, 0x6f, 0x4e, 0x6f, 0x72, 0x6d, 0x61, 0x6c)
	if err != nil {
		return
	}
	if z.TsoNormal == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.TsoNormal.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "TsoContract"
	err = en.Append(0xab, 0x54, 0x73, 0x6f, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74)
	if err != nil {
		return
	}
	if z.TsoContract == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.TsoContract.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "TstateNum"
	err = en.Append(0xa9, 0x54, 0x73, 0x74, 0x61, 0x74, 0x65, 0x4e, 0x75, 0x6d)
	if err != nil {
		return
	}
	if z.TstateNum == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.TstateNum.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *BlockStat) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Ttxs"
	o = append(o, 0x84, 0xa4, 0x54, 0x74, 0x78, 0x73)
	if z.Ttxs == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Ttxs.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "TsoNormal"
	o = append(o, 0xa9, 0x54, 0x73, 0x6f, 0x4e, 0x6f, 0x72, 0x6d, 0x61, 0x6c)
	if z.TsoNormal == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.TsoNormal.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "TsoContract"
	o = append(o, 0xab, 0x54, 0x73, 0x6f, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74)
	if z.TsoContract == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.TsoContract.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "TstateNum"
	o = append(o, 0xa9, 0x54, 0x73, 0x74, 0x61, 0x74, 0x65, 0x4e, 0x75, 0x6d)
	if z.TstateNum == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.TstateNum.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BlockStat) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Ttxs":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Ttxs = nil
			} else {
				if z.Ttxs == nil {
					z.Ttxs = new(types.BigInt)
				}
				bts, err = z.Ttxs.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "TsoNormal":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.TsoNormal = nil
			} else {
				if z.TsoNormal == nil {
					z.TsoNormal = new(types.BigInt)
				}
				bts, err = z.TsoNormal.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "TsoContract":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.TsoContract = nil
			} else {
				if z.TsoContract == nil {
					z.TsoContract = new(types.BigInt)
				}
				bts, err = z.TsoContract.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "TstateNum":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.TstateNum = nil
			} else {
				if z.TstateNum == nil {
					z.TstateNum = new(types.BigInt)
				}
				bts, err = z.TstateNum.UnmarshalMsg(bts)
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
func (z *BlockStat) Msgsize() (s int) {
	s = 1 + 5
	if z.Ttxs == nil {
		s += msgp.NilSize
	} else {
		s += z.Ttxs.Msgsize()
	}
	s += 10
	if z.TsoNormal == nil {
		s += msgp.NilSize
	} else {
		s += z.TsoNormal.Msgsize()
	}
	s += 12
	if z.TsoContract == nil {
		s += msgp.NilSize
	} else {
		s += z.TsoContract.Msgsize()
	}
	s += 10
	if z.TstateNum == nil {
		s += msgp.NilSize
	} else {
		s += z.TstateNum.Msgsize()
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TxLookupEntry) DecodeMsg(dc *msgp.Reader) (err error) {
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
func (z *TxLookupEntry) EncodeMsg(en *msgp.Writer) (err error) {
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
func (z *TxLookupEntry) MarshalMsg(b []byte) (o []byte, err error) {
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
func (z *TxLookupEntry) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
func (z *TxLookupEntry) Msgsize() (s int) {
	s = 1 + 10 + z.BlockHash.Msgsize() + 11 + msgp.Uint64Size + 6 + msgp.Uint64Size
	return
}
