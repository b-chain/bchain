package transaction

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
	"bchain.io/common/types"
)

// DecodeMsg implements msgp.Decodable
func (z *Action) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Contract":
			err = z.Contract.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Params":
			z.Params, err = dc.ReadBytes(z.Params)
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
func (z *Action) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "Contract"
	err = en.Append(0x82, 0xa8, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74)
	if err != nil {
		return
	}
	err = z.Contract.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Params"
	err = en.Append(0xa6, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Params)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Action) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Contract"
	o = append(o, 0x82, 0xa8, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74)
	o, err = z.Contract.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Params"
	o = append(o, 0xa6, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73)
	o = msgp.AppendBytes(o, z.Params)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Action) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Contract":
			bts, err = z.Contract.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Params":
			z.Params, bts, err = msgp.ReadBytesBytes(bts, z.Params)
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
func (z *Action) Msgsize() (s int) {
	s = 1 + 9 + z.Contract.Msgsize() + 7 + msgp.BytesPrefixSize + len(z.Params)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Actions) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(Actions, zb0002)
	}
	for zb0001 := range *z {
		if dc.IsNil() {
			err = dc.ReadNil()
			if err != nil {
				return
			}
			(*z)[zb0001] = nil
		} else {
			if (*z)[zb0001] == nil {
				(*z)[zb0001] = new(Action)
			}
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
				case "Contract":
					err = (*z)[zb0001].Contract.DecodeMsg(dc)
					if err != nil {
						return
					}
				case "Params":
					(*z)[zb0001].Params, err = dc.ReadBytes((*z)[zb0001].Params)
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
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Actions) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zb0004 := range z {
		if z[zb0004] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			// map header, size 2
			// write "Contract"
			err = en.Append(0x82, 0xa8, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74)
			if err != nil {
				return
			}
			err = z[zb0004].Contract.EncodeMsg(en)
			if err != nil {
				return
			}
			// write "Params"
			err = en.Append(0xa6, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73)
			if err != nil {
				return
			}
			err = en.WriteBytes(z[zb0004].Params)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Actions) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zb0004 := range z {
		if z[zb0004] == nil {
			o = msgp.AppendNil(o)
		} else {
			// map header, size 2
			// string "Contract"
			o = append(o, 0x82, 0xa8, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74)
			o, err = z[zb0004].Contract.MarshalMsg(o)
			if err != nil {
				return
			}
			// string "Params"
			o = append(o, 0xa6, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73)
			o = msgp.AppendBytes(o, z[zb0004].Params)
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Actions) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(Actions, zb0002)
	}
	for zb0001 := range *z {
		if msgp.IsNil(bts) {
			bts, err = msgp.ReadNilBytes(bts)
			if err != nil {
				return
			}
			(*z)[zb0001] = nil
		} else {
			if (*z)[zb0001] == nil {
				(*z)[zb0001] = new(Action)
			}
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
				case "Contract":
					bts, err = (*z)[zb0001].Contract.UnmarshalMsg(bts)
					if err != nil {
						return
					}
				case "Params":
					(*z)[zb0001].Params, bts, err = msgp.ReadBytesBytes(bts, (*z)[zb0001].Params)
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
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Actions) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0004 := range z {
		if z[zb0004] == nil {
			s += msgp.NilSize
		} else {
			s += 1 + 9 + z[zb0004].Contract.Msgsize() + 7 + msgp.BytesPrefixSize + len(z[zb0004].Params)
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Transaction) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Data":
			err = z.Data.DecodeMsg(dc)
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
func (z *Transaction) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Data"
	err = en.Append(0x81, 0xa4, 0x44, 0x61, 0x74, 0x61)
	if err != nil {
		return
	}
	err = z.Data.EncodeMsg(en)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Transaction) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Data"
	o = append(o, 0x81, 0xa4, 0x44, 0x61, 0x74, 0x61)
	o, err = z.Data.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Transaction) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Data":
			bts, err = z.Data.UnmarshalMsg(bts)
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
func (z *Transaction) Msgsize() (s int) {
	s = 1 + 5 + z.Data.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Transactions) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(Transactions, zb0002)
	}
	for zb0001 := range *z {
		if dc.IsNil() {
			err = dc.ReadNil()
			if err != nil {
				return
			}
			(*z)[zb0001] = nil
		} else {
			if (*z)[zb0001] == nil {
				(*z)[zb0001] = new(Transaction)
			}
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
				case "Data":
					err = (*z)[zb0001].Data.DecodeMsg(dc)
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
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Transactions) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zb0004 := range z {
		if z[zb0004] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			// map header, size 1
			// write "Data"
			err = en.Append(0x81, 0xa4, 0x44, 0x61, 0x74, 0x61)
			if err != nil {
				return
			}
			err = z[zb0004].Data.EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Transactions) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zb0004 := range z {
		if z[zb0004] == nil {
			o = msgp.AppendNil(o)
		} else {
			// map header, size 1
			// string "Data"
			o = append(o, 0x81, 0xa4, 0x44, 0x61, 0x74, 0x61)
			o, err = z[zb0004].Data.MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Transactions) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(Transactions, zb0002)
	}
	for zb0001 := range *z {
		if msgp.IsNil(bts) {
			bts, err = msgp.ReadNilBytes(bts)
			if err != nil {
				return
			}
			(*z)[zb0001] = nil
		} else {
			if (*z)[zb0001] == nil {
				(*z)[zb0001] = new(Transaction)
			}
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
				case "Data":
					bts, err = (*z)[zb0001].Data.UnmarshalMsg(bts)
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
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Transactions) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0004 := range z {
		if z[zb0004] == nil {
			s += msgp.NilSize
		} else {
			s += 1 + 5 + z[zb0004].Data.Msgsize()
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TxByNonce) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(TxByNonce, zb0002)
	}
	for zb0001 := range *z {
		if dc.IsNil() {
			err = dc.ReadNil()
			if err != nil {
				return
			}
			(*z)[zb0001] = nil
		} else {
			if (*z)[zb0001] == nil {
				(*z)[zb0001] = new(Transaction)
			}
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
				case "Data":
					err = (*z)[zb0001].Data.DecodeMsg(dc)
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
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z TxByNonce) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zb0004 := range z {
		if z[zb0004] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			// map header, size 1
			// write "Data"
			err = en.Append(0x81, 0xa4, 0x44, 0x61, 0x74, 0x61)
			if err != nil {
				return
			}
			err = z[zb0004].Data.EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z TxByNonce) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zb0004 := range z {
		if z[zb0004] == nil {
			o = msgp.AppendNil(o)
		} else {
			// map header, size 1
			// string "Data"
			o = append(o, 0x81, 0xa4, 0x44, 0x61, 0x74, 0x61)
			o, err = z[zb0004].Data.MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TxByNonce) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(TxByNonce, zb0002)
	}
	for zb0001 := range *z {
		if msgp.IsNil(bts) {
			bts, err = msgp.ReadNilBytes(bts)
			if err != nil {
				return
			}
			(*z)[zb0001] = nil
		} else {
			if (*z)[zb0001] == nil {
				(*z)[zb0001] = new(Transaction)
			}
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
				case "Data":
					bts, err = (*z)[zb0001].Data.UnmarshalMsg(bts)
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
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z TxByNonce) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0004 := range z {
		if z[zb0004] == nil {
			s += msgp.NilSize
		} else {
			s += 1 + 5 + z[zb0004].Data.Msgsize()
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TxByPriority) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(TxByPriority, zb0002)
	}
	for zb0001 := range *z {
		if dc.IsNil() {
			err = dc.ReadNil()
			if err != nil {
				return
			}
			(*z)[zb0001] = nil
		} else {
			if (*z)[zb0001] == nil {
				(*z)[zb0001] = new(Transaction)
			}
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
				case "Data":
					err = (*z)[zb0001].Data.DecodeMsg(dc)
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
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z TxByPriority) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zb0004 := range z {
		if z[zb0004] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			// map header, size 1
			// write "Data"
			err = en.Append(0x81, 0xa4, 0x44, 0x61, 0x74, 0x61)
			if err != nil {
				return
			}
			err = z[zb0004].Data.EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z TxByPriority) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zb0004 := range z {
		if z[zb0004] == nil {
			o = msgp.AppendNil(o)
		} else {
			// map header, size 1
			// string "Data"
			o = append(o, 0x81, 0xa4, 0x44, 0x61, 0x74, 0x61)
			o, err = z[zb0004].Data.MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TxByPriority) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(TxByPriority, zb0002)
	}
	for zb0001 := range *z {
		if msgp.IsNil(bts) {
			bts, err = msgp.ReadNilBytes(bts)
			if err != nil {
				return
			}
			(*z)[zb0001] = nil
		} else {
			if (*z)[zb0001] == nil {
				(*z)[zb0001] = new(Transaction)
			}
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
				case "Data":
					bts, err = (*z)[zb0001].Data.UnmarshalMsg(bts)
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
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z TxByPriority) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0004 := range z {
		if z[zb0004] == nil {
			s += msgp.NilSize
		} else {
			s += 1 + 5 + z[zb0004].Data.Msgsize()
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TxHeader) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Nonce":
			z.Nonce, err = dc.ReadUint64()
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
func (z TxHeader) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Nonce"
	err = en.Append(0x81, 0xa5, 0x4e, 0x6f, 0x6e, 0x63, 0x65)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Nonce)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z TxHeader) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Nonce"
	o = append(o, 0x81, 0xa5, 0x4e, 0x6f, 0x6e, 0x63, 0x65)
	o = msgp.AppendUint64(o, z.Nonce)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TxHeader) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Nonce":
			z.Nonce, bts, err = msgp.ReadUint64Bytes(bts)
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
func (z TxHeader) Msgsize() (s int) {
	s = 1 + 6 + msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Txdata) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "H":
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
				case "Nonce":
					z.H.Nonce, err = dc.ReadUint64()
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
		case "Acts":
			var zb0003 uint32
			zb0003, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Acts) >= int(zb0003) {
				z.Acts = (z.Acts)[:zb0003]
			} else {
				z.Acts = make(Actions, zb0003)
			}
			for za0001 := range z.Acts {
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						return
					}
					z.Acts[za0001] = nil
				} else {
					if z.Acts[za0001] == nil {
						z.Acts[za0001] = new(Action)
					}
					var zb0004 uint32
					zb0004, err = dc.ReadMapHeader()
					if err != nil {
						return
					}
					for zb0004 > 0 {
						zb0004--
						field, err = dc.ReadMapKeyPtr()
						if err != nil {
							return
						}
						switch msgp.UnsafeString(field) {
						case "Contract":
							err = z.Acts[za0001].Contract.DecodeMsg(dc)
							if err != nil {
								return
							}
						case "Params":
							z.Acts[za0001].Params, err = dc.ReadBytes(z.Acts[za0001].Params)
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
			}
		case "V":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.V = nil
			} else {
				if z.V == nil {
					z.V = new(types.BigInt)
				}
				err = z.V.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "R":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.R = nil
			} else {
				if z.R == nil {
					z.R = new(types.BigInt)
				}
				err = z.R.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "S":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.S = nil
			} else {
				if z.S == nil {
					z.S = new(types.BigInt)
				}
				err = z.S.DecodeMsg(dc)
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
func (z *Txdata) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "H"
	// map header, size 1
	// write "Nonce"
	err = en.Append(0x85, 0xa1, 0x48, 0x81, 0xa5, 0x4e, 0x6f, 0x6e, 0x63, 0x65)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.H.Nonce)
	if err != nil {
		return
	}
	// write "Acts"
	err = en.Append(0xa4, 0x41, 0x63, 0x74, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Acts)))
	if err != nil {
		return
	}
	for za0001 := range z.Acts {
		if z.Acts[za0001] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			// map header, size 2
			// write "Contract"
			err = en.Append(0x82, 0xa8, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74)
			if err != nil {
				return
			}
			err = z.Acts[za0001].Contract.EncodeMsg(en)
			if err != nil {
				return
			}
			// write "Params"
			err = en.Append(0xa6, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73)
			if err != nil {
				return
			}
			err = en.WriteBytes(z.Acts[za0001].Params)
			if err != nil {
				return
			}
		}
	}
	// write "V"
	err = en.Append(0xa1, 0x56)
	if err != nil {
		return
	}
	if z.V == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.V.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "R"
	err = en.Append(0xa1, 0x52)
	if err != nil {
		return
	}
	if z.R == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.R.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "S"
	err = en.Append(0xa1, 0x53)
	if err != nil {
		return
	}
	if z.S == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.S.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Txdata) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "H"
	// map header, size 1
	// string "Nonce"
	o = append(o, 0x85, 0xa1, 0x48, 0x81, 0xa5, 0x4e, 0x6f, 0x6e, 0x63, 0x65)
	o = msgp.AppendUint64(o, z.H.Nonce)
	// string "Acts"
	o = append(o, 0xa4, 0x41, 0x63, 0x74, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Acts)))
	for za0001 := range z.Acts {
		if z.Acts[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			// map header, size 2
			// string "Contract"
			o = append(o, 0x82, 0xa8, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74)
			o, err = z.Acts[za0001].Contract.MarshalMsg(o)
			if err != nil {
				return
			}
			// string "Params"
			o = append(o, 0xa6, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73)
			o = msgp.AppendBytes(o, z.Acts[za0001].Params)
		}
	}
	// string "V"
	o = append(o, 0xa1, 0x56)
	if z.V == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.V.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "R"
	o = append(o, 0xa1, 0x52)
	if z.R == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.R.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "S"
	o = append(o, 0xa1, 0x53)
	if z.S == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.S.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Txdata) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "H":
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
				case "Nonce":
					z.H.Nonce, bts, err = msgp.ReadUint64Bytes(bts)
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
		case "Acts":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Acts) >= int(zb0003) {
				z.Acts = (z.Acts)[:zb0003]
			} else {
				z.Acts = make(Actions, zb0003)
			}
			for za0001 := range z.Acts {
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					z.Acts[za0001] = nil
				} else {
					if z.Acts[za0001] == nil {
						z.Acts[za0001] = new(Action)
					}
					var zb0004 uint32
					zb0004, bts, err = msgp.ReadMapHeaderBytes(bts)
					if err != nil {
						return
					}
					for zb0004 > 0 {
						zb0004--
						field, bts, err = msgp.ReadMapKeyZC(bts)
						if err != nil {
							return
						}
						switch msgp.UnsafeString(field) {
						case "Contract":
							bts, err = z.Acts[za0001].Contract.UnmarshalMsg(bts)
							if err != nil {
								return
							}
						case "Params":
							z.Acts[za0001].Params, bts, err = msgp.ReadBytesBytes(bts, z.Acts[za0001].Params)
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
			}
		case "V":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.V = nil
			} else {
				if z.V == nil {
					z.V = new(types.BigInt)
				}
				bts, err = z.V.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "R":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.R = nil
			} else {
				if z.R == nil {
					z.R = new(types.BigInt)
				}
				bts, err = z.R.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "S":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.S = nil
			} else {
				if z.S == nil {
					z.S = new(types.BigInt)
				}
				bts, err = z.S.UnmarshalMsg(bts)
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
func (z *Txdata) Msgsize() (s int) {
	s = 1 + 2 + 1 + 6 + msgp.Uint64Size + 5 + msgp.ArrayHeaderSize
	for za0001 := range z.Acts {
		if z.Acts[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += 1 + 9 + z.Acts[za0001].Contract.Msgsize() + 7 + msgp.BytesPrefixSize + len(z.Acts[za0001].Params)
		}
	}
	s += 2
	if z.V == nil {
		s += msgp.NilSize
	} else {
		s += z.V.Msgsize()
	}
	s += 2
	if z.R == nil {
		s += msgp.NilSize
	} else {
		s += z.R.Msgsize()
	}
	s += 2
	if z.S == nil {
		s += msgp.NilSize
	} else {
		s += z.S.Msgsize()
	}
	return
}
