package nativere

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *DumpData) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Amount":
			z.Amount, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "PledgeData":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.PledgeData = nil
			} else {
				if z.PledgeData == nil {
					z.PledgeData = new(PledgeData)
				}
				err = z.PledgeData.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "PoolData":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.PoolData = nil
			} else {
				if z.PoolData == nil {
					z.PoolData = new(PoolData)
				}
				err = z.PoolData.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "ProducerData":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.ProducerData = nil
			} else {
				if z.ProducerData == nil {
					z.ProducerData = new(ProducerData)
				}
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
					case "ChildHead":
						err = z.ProducerData.ChildHead.DecodeMsg(dc)
						if err != nil {
							return
						}
					case "Amount":
						z.ProducerData.Amount, err = dc.ReadUint64()
						if err != nil {
							return
						}
					case "ProducerCert":
						z.ProducerData.ProducerCert, err = dc.ReadUint64()
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
func (z *DumpData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "Amount"
	err = en.Append(0x84, 0xa6, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Amount)
	if err != nil {
		return
	}
	// write "PledgeData"
	err = en.Append(0xaa, 0x50, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61)
	if err != nil {
		return
	}
	if z.PledgeData == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.PledgeData.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "PoolData"
	err = en.Append(0xa8, 0x50, 0x6f, 0x6f, 0x6c, 0x44, 0x61, 0x74, 0x61)
	if err != nil {
		return
	}
	if z.PoolData == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.PoolData.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "ProducerData"
	err = en.Append(0xac, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61)
	if err != nil {
		return
	}
	if z.ProducerData == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		// map header, size 3
		// write "ChildHead"
		err = en.Append(0x83, 0xa9, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x48, 0x65, 0x61, 0x64)
		if err != nil {
			return
		}
		err = z.ProducerData.ChildHead.EncodeMsg(en)
		if err != nil {
			return
		}
		// write "Amount"
		err = en.Append(0xa6, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74)
		if err != nil {
			return
		}
		err = en.WriteUint64(z.ProducerData.Amount)
		if err != nil {
			return
		}
		// write "ProducerCert"
		err = en.Append(0xac, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x43, 0x65, 0x72, 0x74)
		if err != nil {
			return
		}
		err = en.WriteUint64(z.ProducerData.ProducerCert)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *DumpData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Amount"
	o = append(o, 0x84, 0xa6, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74)
	o = msgp.AppendUint64(o, z.Amount)
	// string "PledgeData"
	o = append(o, 0xaa, 0x50, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61)
	if z.PledgeData == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.PledgeData.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "PoolData"
	o = append(o, 0xa8, 0x50, 0x6f, 0x6f, 0x6c, 0x44, 0x61, 0x74, 0x61)
	if z.PoolData == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.PoolData.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "ProducerData"
	o = append(o, 0xac, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61)
	if z.ProducerData == nil {
		o = msgp.AppendNil(o)
	} else {
		// map header, size 3
		// string "ChildHead"
		o = append(o, 0x83, 0xa9, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x48, 0x65, 0x61, 0x64)
		o, err = z.ProducerData.ChildHead.MarshalMsg(o)
		if err != nil {
			return
		}
		// string "Amount"
		o = append(o, 0xa6, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74)
		o = msgp.AppendUint64(o, z.ProducerData.Amount)
		// string "ProducerCert"
		o = append(o, 0xac, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x43, 0x65, 0x72, 0x74)
		o = msgp.AppendUint64(o, z.ProducerData.ProducerCert)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DumpData) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Amount":
			z.Amount, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "PledgeData":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.PledgeData = nil
			} else {
				if z.PledgeData == nil {
					z.PledgeData = new(PledgeData)
				}
				bts, err = z.PledgeData.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "PoolData":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.PoolData = nil
			} else {
				if z.PoolData == nil {
					z.PoolData = new(PoolData)
				}
				bts, err = z.PoolData.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "ProducerData":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.ProducerData = nil
			} else {
				if z.ProducerData == nil {
					z.ProducerData = new(ProducerData)
				}
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
					case "ChildHead":
						bts, err = z.ProducerData.ChildHead.UnmarshalMsg(bts)
						if err != nil {
							return
						}
					case "Amount":
						z.ProducerData.Amount, bts, err = msgp.ReadUint64Bytes(bts)
						if err != nil {
							return
						}
					case "ProducerCert":
						z.ProducerData.ProducerCert, bts, err = msgp.ReadUint64Bytes(bts)
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
func (z *DumpData) Msgsize() (s int) {
	s = 1 + 7 + msgp.Uint64Size + 11
	if z.PledgeData == nil {
		s += msgp.NilSize
	} else {
		s += z.PledgeData.Msgsize()
	}
	s += 9
	if z.PoolData == nil {
		s += msgp.NilSize
	} else {
		s += z.PoolData.Msgsize()
	}
	s += 13
	if z.ProducerData == nil {
		s += msgp.NilSize
	} else {
		s += 1 + 10 + z.ProducerData.ChildHead.Msgsize() + 7 + msgp.Uint64Size + 13 + msgp.Uint64Size
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *NativePara) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "funcName":
			z.FuncName, err = dc.ReadString()
			if err != nil {
				return
			}
		case "args":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Args) >= int(zb0002) {
				z.Args = (z.Args)[:zb0002]
			} else {
				z.Args = make([]string, zb0002)
			}
			for za0001 := range z.Args {
				z.Args[za0001], err = dc.ReadString()
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
func (z *NativePara) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "funcName"
	err = en.Append(0x82, 0xa8, 0x66, 0x75, 0x6e, 0x63, 0x4e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.FuncName)
	if err != nil {
		return
	}
	// write "args"
	err = en.Append(0xa4, 0x61, 0x72, 0x67, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Args)))
	if err != nil {
		return
	}
	for za0001 := range z.Args {
		err = en.WriteString(z.Args[za0001])
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *NativePara) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "funcName"
	o = append(o, 0x82, 0xa8, 0x66, 0x75, 0x6e, 0x63, 0x4e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.FuncName)
	// string "args"
	o = append(o, 0xa4, 0x61, 0x72, 0x67, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Args)))
	for za0001 := range z.Args {
		o = msgp.AppendString(o, z.Args[za0001])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *NativePara) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "funcName":
			z.FuncName, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "args":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Args) >= int(zb0002) {
				z.Args = (z.Args)[:zb0002]
			} else {
				z.Args = make([]string, zb0002)
			}
			for za0001 := range z.Args {
				z.Args[za0001], bts, err = msgp.ReadStringBytes(bts)
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
func (z *NativePara) Msgsize() (s int) {
	s = 1 + 9 + msgp.StringPrefixSize + len(z.FuncName) + 5 + msgp.ArrayHeaderSize
	for za0001 := range z.Args {
		s += msgp.StringPrefixSize + len(z.Args[za0001])
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PledgeData) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Parent":
			err = z.Parent.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Prev":
			err = z.Prev.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Next":
			err = z.Next.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Amount":
			z.Amount, err = dc.ReadUint64()
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
func (z *PledgeData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "Parent"
	err = en.Append(0x84, 0xa6, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74)
	if err != nil {
		return
	}
	err = z.Parent.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Prev"
	err = en.Append(0xa4, 0x50, 0x72, 0x65, 0x76)
	if err != nil {
		return
	}
	err = z.Prev.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Next"
	err = en.Append(0xa4, 0x4e, 0x65, 0x78, 0x74)
	if err != nil {
		return
	}
	err = z.Next.EncodeMsg(en)
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
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *PledgeData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Parent"
	o = append(o, 0x84, 0xa6, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74)
	o, err = z.Parent.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Prev"
	o = append(o, 0xa4, 0x50, 0x72, 0x65, 0x76)
	o, err = z.Prev.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Next"
	o = append(o, 0xa4, 0x4e, 0x65, 0x78, 0x74)
	o, err = z.Next.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Amount"
	o = append(o, 0xa6, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74)
	o = msgp.AppendUint64(o, z.Amount)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PledgeData) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Parent":
			bts, err = z.Parent.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Prev":
			bts, err = z.Prev.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Next":
			bts, err = z.Next.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Amount":
			z.Amount, bts, err = msgp.ReadUint64Bytes(bts)
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
func (z *PledgeData) Msgsize() (s int) {
	s = 1 + 7 + z.Parent.Msgsize() + 5 + z.Prev.Msgsize() + 5 + z.Next.Msgsize() + 7 + msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PoolData) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Producer":
			err = z.Producer.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Prev":
			err = z.Prev.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Next":
			err = z.Next.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Amount":
			z.Amount, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "ChildHead":
			err = z.ChildHead.DecodeMsg(dc)
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
func (z *PoolData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "Producer"
	err = en.Append(0x85, 0xa8, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72)
	if err != nil {
		return
	}
	err = z.Producer.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Prev"
	err = en.Append(0xa4, 0x50, 0x72, 0x65, 0x76)
	if err != nil {
		return
	}
	err = z.Prev.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Next"
	err = en.Append(0xa4, 0x4e, 0x65, 0x78, 0x74)
	if err != nil {
		return
	}
	err = z.Next.EncodeMsg(en)
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
	// write "ChildHead"
	err = en.Append(0xa9, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x48, 0x65, 0x61, 0x64)
	if err != nil {
		return
	}
	err = z.ChildHead.EncodeMsg(en)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *PoolData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "Producer"
	o = append(o, 0x85, 0xa8, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72)
	o, err = z.Producer.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Prev"
	o = append(o, 0xa4, 0x50, 0x72, 0x65, 0x76)
	o, err = z.Prev.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Next"
	o = append(o, 0xa4, 0x4e, 0x65, 0x78, 0x74)
	o, err = z.Next.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Amount"
	o = append(o, 0xa6, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74)
	o = msgp.AppendUint64(o, z.Amount)
	// string "ChildHead"
	o = append(o, 0xa9, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x48, 0x65, 0x61, 0x64)
	o, err = z.ChildHead.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PoolData) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Producer":
			bts, err = z.Producer.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Prev":
			bts, err = z.Prev.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Next":
			bts, err = z.Next.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Amount":
			z.Amount, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "ChildHead":
			bts, err = z.ChildHead.UnmarshalMsg(bts)
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
func (z *PoolData) Msgsize() (s int) {
	s = 1 + 9 + z.Producer.Msgsize() + 5 + z.Prev.Msgsize() + 5 + z.Next.Msgsize() + 7 + msgp.Uint64Size + 10 + z.ChildHead.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ProducerData) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "ChildHead":
			err = z.ChildHead.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Amount":
			z.Amount, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "ProducerCert":
			z.ProducerCert, err = dc.ReadUint64()
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
func (z *ProducerData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "ChildHead"
	err = en.Append(0x83, 0xa9, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x48, 0x65, 0x61, 0x64)
	if err != nil {
		return
	}
	err = z.ChildHead.EncodeMsg(en)
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
	// write "ProducerCert"
	err = en.Append(0xac, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x43, 0x65, 0x72, 0x74)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.ProducerCert)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ProducerData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "ChildHead"
	o = append(o, 0x83, 0xa9, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x48, 0x65, 0x61, 0x64)
	o, err = z.ChildHead.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Amount"
	o = append(o, 0xa6, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74)
	o = msgp.AppendUint64(o, z.Amount)
	// string "ProducerCert"
	o = append(o, 0xac, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x43, 0x65, 0x72, 0x74)
	o = msgp.AppendUint64(o, z.ProducerCert)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ProducerData) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ChildHead":
			bts, err = z.ChildHead.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Amount":
			z.Amount, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "ProducerCert":
			z.ProducerCert, bts, err = msgp.ReadUint64Bytes(bts)
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
func (z *ProducerData) Msgsize() (s int) {
	s = 1 + 10 + z.ChildHead.Msgsize() + 7 + msgp.Uint64Size + 13 + msgp.Uint64Size
	return
}
