package aposgensis

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *WeightInfo) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "AddrStr":
			err = z.AddrStr.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Wt":
			z.Wt, err = dc.ReadUint64()
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
func (z *WeightInfo) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "AddrStr"
	err = en.Append(0x82, 0xa7, 0x41, 0x64, 0x64, 0x72, 0x53, 0x74, 0x72)
	if err != nil {
		return
	}
	err = z.AddrStr.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Wt"
	err = en.Append(0xa2, 0x57, 0x74)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Wt)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *WeightInfo) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "AddrStr"
	o = append(o, 0x82, 0xa7, 0x41, 0x64, 0x64, 0x72, 0x53, 0x74, 0x72)
	o, err = z.AddrStr.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Wt"
	o = append(o, 0xa2, 0x57, 0x74)
	o = msgp.AppendUint64(o, z.Wt)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *WeightInfo) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "AddrStr":
			bts, err = z.AddrStr.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Wt":
			z.Wt, bts, err = msgp.ReadUint64Bytes(bts)
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
func (z *WeightInfo) Msgsize() (s int) {
	s = 1 + 8 + z.AddrStr.Msgsize() + 3 + msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *WeightInfos) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(WeightInfos, zb0002)
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
			case "AddrStr":
				err = (*z)[zb0001].AddrStr.DecodeMsg(dc)
				if err != nil {
					return
				}
			case "Wt":
				(*z)[zb0001].Wt, err = dc.ReadUint64()
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
func (z WeightInfos) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zb0004 := range z {
		// map header, size 2
		// write "AddrStr"
		err = en.Append(0x82, 0xa7, 0x41, 0x64, 0x64, 0x72, 0x53, 0x74, 0x72)
		if err != nil {
			return
		}
		err = z[zb0004].AddrStr.EncodeMsg(en)
		if err != nil {
			return
		}
		// write "Wt"
		err = en.Append(0xa2, 0x57, 0x74)
		if err != nil {
			return
		}
		err = en.WriteUint64(z[zb0004].Wt)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z WeightInfos) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zb0004 := range z {
		// map header, size 2
		// string "AddrStr"
		o = append(o, 0x82, 0xa7, 0x41, 0x64, 0x64, 0x72, 0x53, 0x74, 0x72)
		o, err = z[zb0004].AddrStr.MarshalMsg(o)
		if err != nil {
			return
		}
		// string "Wt"
		o = append(o, 0xa2, 0x57, 0x74)
		o = msgp.AppendUint64(o, z[zb0004].Wt)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *WeightInfos) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(WeightInfos, zb0002)
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
			case "AddrStr":
				bts, err = (*z)[zb0001].AddrStr.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			case "Wt":
				(*z)[zb0001].Wt, bts, err = msgp.ReadUint64Bytes(bts)
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
func (z WeightInfos) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0004 := range z {
		s += 1 + 8 + z[zb0004].AddrStr.Msgsize() + 3 + msgp.Uint64Size
	}
	return
}
