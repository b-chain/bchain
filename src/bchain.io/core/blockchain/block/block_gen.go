package block

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
	"bchain.io/common/types"
	"bchain.io/core/transaction"
)

// DecodeMsg implements msgp.Decodable
func (z *Block) DecodeMsg(dc *msgp.Reader) (err error) {
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
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.H = nil
			} else {
				if z.H == nil {
					z.H = new(Header)
				}
				err = z.H.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "B":
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
				case "Transactions":
					var zb0003 uint32
					zb0003, err = dc.ReadArrayHeader()
					if err != nil {
						return
					}
					if cap(z.B.Transactions) >= int(zb0003) {
						z.B.Transactions = (z.B.Transactions)[:zb0003]
					} else {
						z.B.Transactions = make([]*transaction.Transaction, zb0003)
					}
					for za0001 := range z.B.Transactions {
						if dc.IsNil() {
							err = dc.ReadNil()
							if err != nil {
								return
							}
							z.B.Transactions[za0001] = nil
						} else {
							if z.B.Transactions[za0001] == nil {
								z.B.Transactions[za0001] = new(transaction.Transaction)
							}
							err = z.B.Transactions[za0001].DecodeMsg(dc)
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
func (z *Block) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "H"
	err = en.Append(0x82, 0xa1, 0x48)
	if err != nil {
		return
	}
	if z.H == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.H.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "B"
	// map header, size 1
	// write "Transactions"
	err = en.Append(0xa1, 0x42, 0x81, 0xac, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.B.Transactions)))
	if err != nil {
		return
	}
	for za0001 := range z.B.Transactions {
		if z.B.Transactions[za0001] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z.B.Transactions[za0001].EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Block) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "H"
	o = append(o, 0x82, 0xa1, 0x48)
	if z.H == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.H.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "B"
	// map header, size 1
	// string "Transactions"
	o = append(o, 0xa1, 0x42, 0x81, 0xac, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.B.Transactions)))
	for za0001 := range z.B.Transactions {
		if z.B.Transactions[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.B.Transactions[za0001].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Block) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.H = nil
			} else {
				if z.H == nil {
					z.H = new(Header)
				}
				bts, err = z.H.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "B":
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
				case "Transactions":
					var zb0003 uint32
					zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
					if err != nil {
						return
					}
					if cap(z.B.Transactions) >= int(zb0003) {
						z.B.Transactions = (z.B.Transactions)[:zb0003]
					} else {
						z.B.Transactions = make([]*transaction.Transaction, zb0003)
					}
					for za0001 := range z.B.Transactions {
						if msgp.IsNil(bts) {
							bts, err = msgp.ReadNilBytes(bts)
							if err != nil {
								return
							}
							z.B.Transactions[za0001] = nil
						} else {
							if z.B.Transactions[za0001] == nil {
								z.B.Transactions[za0001] = new(transaction.Transaction)
							}
							bts, err = z.B.Transactions[za0001].UnmarshalMsg(bts)
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
func (z *Block) Msgsize() (s int) {
	s = 1 + 2
	if z.H == nil {
		s += msgp.NilSize
	} else {
		s += z.H.Msgsize()
	}
	s += 2 + 1 + 13 + msgp.ArrayHeaderSize
	for za0001 := range z.B.Transactions {
		if z.B.Transactions[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += z.B.Transactions[za0001].Msgsize()
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Body) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Transactions":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Transactions) >= int(zb0002) {
				z.Transactions = (z.Transactions)[:zb0002]
			} else {
				z.Transactions = make([]*transaction.Transaction, zb0002)
			}
			for za0001 := range z.Transactions {
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						return
					}
					z.Transactions[za0001] = nil
				} else {
					if z.Transactions[za0001] == nil {
						z.Transactions[za0001] = new(transaction.Transaction)
					}
					err = z.Transactions[za0001].DecodeMsg(dc)
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
func (z *Body) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Transactions"
	err = en.Append(0x81, 0xac, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Transactions)))
	if err != nil {
		return
	}
	for za0001 := range z.Transactions {
		if z.Transactions[za0001] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z.Transactions[za0001].EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Body) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Transactions"
	o = append(o, 0x81, 0xac, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Transactions)))
	for za0001 := range z.Transactions {
		if z.Transactions[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.Transactions[za0001].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Body) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Transactions":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Transactions) >= int(zb0002) {
				z.Transactions = (z.Transactions)[:zb0002]
			} else {
				z.Transactions = make([]*transaction.Transaction, zb0002)
			}
			for za0001 := range z.Transactions {
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					z.Transactions[za0001] = nil
				} else {
					if z.Transactions[za0001] == nil {
						z.Transactions[za0001] = new(transaction.Transaction)
					}
					bts, err = z.Transactions[za0001].UnmarshalMsg(bts)
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
func (z *Body) Msgsize() (s int) {
	s = 1 + 13 + msgp.ArrayHeaderSize
	for za0001 := range z.Transactions {
		if z.Transactions[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += z.Transactions[za0001].Msgsize()
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ConsensusData) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Id":
			z.Id, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Para":
			z.Para, err = dc.ReadBytes(z.Para)
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
func (z *ConsensusData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "Id"
	err = en.Append(0x82, 0xa2, 0x49, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.Id)
	if err != nil {
		return
	}
	// write "Para"
	err = en.Append(0xa4, 0x50, 0x61, 0x72, 0x61)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Para)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ConsensusData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Id"
	o = append(o, 0x82, 0xa2, 0x49, 0x64)
	o = msgp.AppendString(o, z.Id)
	// string "Para"
	o = append(o, 0xa4, 0x50, 0x61, 0x72, 0x61)
	o = msgp.AppendBytes(o, z.Para)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ConsensusData) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Id":
			z.Id, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Para":
			z.Para, bts, err = msgp.ReadBytesBytes(bts, z.Para)
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
func (z *ConsensusData) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.Id) + 5 + msgp.BytesPrefixSize + len(z.Para)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Header) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "ParentHash":
			err = z.ParentHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "StateRootHash":
			err = z.StateRootHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "TxRootHash":
			err = z.TxRootHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "ReceiptRootHash":
			err = z.ReceiptRootHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Bloom":
			err = z.Bloom.DecodeMsg(dc)
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
		case "Time":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Time = nil
			} else {
				if z.Time == nil {
					z.Time = new(types.BigInt)
				}
				err = z.Time.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Cdata":
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
				case "Id":
					z.Cdata.Id, err = dc.ReadString()
					if err != nil {
						return
					}
				case "Para":
					z.Cdata.Para, err = dc.ReadBytes(z.Cdata.Para)
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
		case "Extra":
			z.Extra, err = dc.ReadBytes(z.Extra)
			if err != nil {
				return
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
func (z *Header) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 12
	// write "ParentHash"
	err = en.Append(0x8c, 0xaa, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.ParentHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "StateRootHash"
	err = en.Append(0xad, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.StateRootHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "TxRootHash"
	err = en.Append(0xaa, 0x54, 0x78, 0x52, 0x6f, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.TxRootHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "ReceiptRootHash"
	err = en.Append(0xaf, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x6f, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.ReceiptRootHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Bloom"
	err = en.Append(0xa5, 0x42, 0x6c, 0x6f, 0x6f, 0x6d)
	if err != nil {
		return
	}
	err = z.Bloom.EncodeMsg(en)
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
	// write "Time"
	err = en.Append(0xa4, 0x54, 0x69, 0x6d, 0x65)
	if err != nil {
		return
	}
	if z.Time == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Time.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Cdata"
	// map header, size 2
	// write "Id"
	err = en.Append(0xa5, 0x43, 0x64, 0x61, 0x74, 0x61, 0x82, 0xa2, 0x49, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.Cdata.Id)
	if err != nil {
		return
	}
	// write "Para"
	err = en.Append(0xa4, 0x50, 0x61, 0x72, 0x61)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Cdata.Para)
	if err != nil {
		return
	}
	// write "Extra"
	err = en.Append(0xa5, 0x45, 0x78, 0x74, 0x72, 0x61)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Extra)
	if err != nil {
		return
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
func (z *Header) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 12
	// string "ParentHash"
	o = append(o, 0x8c, 0xaa, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68)
	o, err = z.ParentHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "StateRootHash"
	o = append(o, 0xad, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68)
	o, err = z.StateRootHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "TxRootHash"
	o = append(o, 0xaa, 0x54, 0x78, 0x52, 0x6f, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68)
	o, err = z.TxRootHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "ReceiptRootHash"
	o = append(o, 0xaf, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x6f, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68)
	o, err = z.ReceiptRootHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Bloom"
	o = append(o, 0xa5, 0x42, 0x6c, 0x6f, 0x6f, 0x6d)
	o, err = z.Bloom.MarshalMsg(o)
	if err != nil {
		return
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
	// string "Time"
	o = append(o, 0xa4, 0x54, 0x69, 0x6d, 0x65)
	if z.Time == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Time.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Cdata"
	// map header, size 2
	// string "Id"
	o = append(o, 0xa5, 0x43, 0x64, 0x61, 0x74, 0x61, 0x82, 0xa2, 0x49, 0x64)
	o = msgp.AppendString(o, z.Cdata.Id)
	// string "Para"
	o = append(o, 0xa4, 0x50, 0x61, 0x72, 0x61)
	o = msgp.AppendBytes(o, z.Cdata.Para)
	// string "Extra"
	o = append(o, 0xa5, 0x45, 0x78, 0x74, 0x72, 0x61)
	o = msgp.AppendBytes(o, z.Extra)
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
func (z *Header) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ParentHash":
			bts, err = z.ParentHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "StateRootHash":
			bts, err = z.StateRootHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "TxRootHash":
			bts, err = z.TxRootHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "ReceiptRootHash":
			bts, err = z.ReceiptRootHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Bloom":
			bts, err = z.Bloom.UnmarshalMsg(bts)
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
		case "Time":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Time = nil
			} else {
				if z.Time == nil {
					z.Time = new(types.BigInt)
				}
				bts, err = z.Time.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Cdata":
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
				case "Id":
					z.Cdata.Id, bts, err = msgp.ReadStringBytes(bts)
					if err != nil {
						return
					}
				case "Para":
					z.Cdata.Para, bts, err = msgp.ReadBytesBytes(bts, z.Cdata.Para)
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
		case "Extra":
			z.Extra, bts, err = msgp.ReadBytesBytes(bts, z.Extra)
			if err != nil {
				return
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
func (z *Header) Msgsize() (s int) {
	s = 1 + 11 + z.ParentHash.Msgsize() + 14 + z.StateRootHash.Msgsize() + 11 + z.TxRootHash.Msgsize() + 16 + z.ReceiptRootHash.Msgsize() + 6 + z.Bloom.Msgsize() + 7
	if z.Number == nil {
		s += msgp.NilSize
	} else {
		s += z.Number.Msgsize()
	}
	s += 5
	if z.Time == nil {
		s += msgp.NilSize
	} else {
		s += z.Time.Msgsize()
	}
	s += 6 + 1 + 3 + msgp.StringPrefixSize + len(z.Cdata.Id) + 5 + msgp.BytesPrefixSize + len(z.Cdata.Para) + 6 + msgp.BytesPrefixSize + len(z.Extra) + 2
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

// DecodeMsg implements msgp.Decodable
func (z *HeaderNoSig) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "ParentHash":
			err = z.ParentHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "StateRootHash":
			err = z.StateRootHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "TxRootHash":
			err = z.TxRootHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "ReceiptRootHash":
			err = z.ReceiptRootHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Bloom":
			err = z.Bloom.DecodeMsg(dc)
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
		case "Time":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Time = nil
			} else {
				if z.Time == nil {
					z.Time = new(types.BigInt)
				}
				err = z.Time.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Cdata":
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
				case "Id":
					z.Cdata.Id, err = dc.ReadString()
					if err != nil {
						return
					}
				case "Para":
					z.Cdata.Para, err = dc.ReadBytes(z.Cdata.Para)
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
		case "Extra":
			z.Extra, err = dc.ReadBytes(z.Extra)
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
func (z *HeaderNoSig) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 9
	// write "ParentHash"
	err = en.Append(0x89, 0xaa, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.ParentHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "StateRootHash"
	err = en.Append(0xad, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.StateRootHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "TxRootHash"
	err = en.Append(0xaa, 0x54, 0x78, 0x52, 0x6f, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.TxRootHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "ReceiptRootHash"
	err = en.Append(0xaf, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x6f, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.ReceiptRootHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Bloom"
	err = en.Append(0xa5, 0x42, 0x6c, 0x6f, 0x6f, 0x6d)
	if err != nil {
		return
	}
	err = z.Bloom.EncodeMsg(en)
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
	// write "Time"
	err = en.Append(0xa4, 0x54, 0x69, 0x6d, 0x65)
	if err != nil {
		return
	}
	if z.Time == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Time.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Cdata"
	// map header, size 2
	// write "Id"
	err = en.Append(0xa5, 0x43, 0x64, 0x61, 0x74, 0x61, 0x82, 0xa2, 0x49, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.Cdata.Id)
	if err != nil {
		return
	}
	// write "Para"
	err = en.Append(0xa4, 0x50, 0x61, 0x72, 0x61)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Cdata.Para)
	if err != nil {
		return
	}
	// write "Extra"
	err = en.Append(0xa5, 0x45, 0x78, 0x74, 0x72, 0x61)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Extra)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *HeaderNoSig) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 9
	// string "ParentHash"
	o = append(o, 0x89, 0xaa, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68)
	o, err = z.ParentHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "StateRootHash"
	o = append(o, 0xad, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68)
	o, err = z.StateRootHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "TxRootHash"
	o = append(o, 0xaa, 0x54, 0x78, 0x52, 0x6f, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68)
	o, err = z.TxRootHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "ReceiptRootHash"
	o = append(o, 0xaf, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x6f, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68)
	o, err = z.ReceiptRootHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Bloom"
	o = append(o, 0xa5, 0x42, 0x6c, 0x6f, 0x6f, 0x6d)
	o, err = z.Bloom.MarshalMsg(o)
	if err != nil {
		return
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
	// string "Time"
	o = append(o, 0xa4, 0x54, 0x69, 0x6d, 0x65)
	if z.Time == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Time.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Cdata"
	// map header, size 2
	// string "Id"
	o = append(o, 0xa5, 0x43, 0x64, 0x61, 0x74, 0x61, 0x82, 0xa2, 0x49, 0x64)
	o = msgp.AppendString(o, z.Cdata.Id)
	// string "Para"
	o = append(o, 0xa4, 0x50, 0x61, 0x72, 0x61)
	o = msgp.AppendBytes(o, z.Cdata.Para)
	// string "Extra"
	o = append(o, 0xa5, 0x45, 0x78, 0x74, 0x72, 0x61)
	o = msgp.AppendBytes(o, z.Extra)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *HeaderNoSig) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ParentHash":
			bts, err = z.ParentHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "StateRootHash":
			bts, err = z.StateRootHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "TxRootHash":
			bts, err = z.TxRootHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "ReceiptRootHash":
			bts, err = z.ReceiptRootHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Bloom":
			bts, err = z.Bloom.UnmarshalMsg(bts)
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
		case "Time":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Time = nil
			} else {
				if z.Time == nil {
					z.Time = new(types.BigInt)
				}
				bts, err = z.Time.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Cdata":
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
				case "Id":
					z.Cdata.Id, bts, err = msgp.ReadStringBytes(bts)
					if err != nil {
						return
					}
				case "Para":
					z.Cdata.Para, bts, err = msgp.ReadBytesBytes(bts, z.Cdata.Para)
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
		case "Extra":
			z.Extra, bts, err = msgp.ReadBytesBytes(bts, z.Extra)
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
func (z *HeaderNoSig) Msgsize() (s int) {
	s = 1 + 11 + z.ParentHash.Msgsize() + 14 + z.StateRootHash.Msgsize() + 11 + z.TxRootHash.Msgsize() + 16 + z.ReceiptRootHash.Msgsize() + 6 + z.Bloom.Msgsize() + 7
	if z.Number == nil {
		s += msgp.NilSize
	} else {
		s += z.Number.Msgsize()
	}
	s += 5
	if z.Time == nil {
		s += msgp.NilSize
	} else {
		s += z.Time.Msgsize()
	}
	s += 6 + 1 + 3 + msgp.StringPrefixSize + len(z.Cdata.Id) + 5 + msgp.BytesPrefixSize + len(z.Cdata.Para) + 6 + msgp.BytesPrefixSize + len(z.Extra)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Headers) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Headers":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Headers) >= int(zb0002) {
				z.Headers = (z.Headers)[:zb0002]
			} else {
				z.Headers = make([]*Header, zb0002)
			}
			for za0001 := range z.Headers {
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						return
					}
					z.Headers[za0001] = nil
				} else {
					if z.Headers[za0001] == nil {
						z.Headers[za0001] = new(Header)
					}
					err = z.Headers[za0001].DecodeMsg(dc)
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
func (z *Headers) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Headers"
	err = en.Append(0x81, 0xa7, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Headers)))
	if err != nil {
		return
	}
	for za0001 := range z.Headers {
		if z.Headers[za0001] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z.Headers[za0001].EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Headers) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Headers"
	o = append(o, 0x81, 0xa7, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Headers)))
	for za0001 := range z.Headers {
		if z.Headers[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.Headers[za0001].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Headers) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Headers":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Headers) >= int(zb0002) {
				z.Headers = (z.Headers)[:zb0002]
			} else {
				z.Headers = make([]*Header, zb0002)
			}
			for za0001 := range z.Headers {
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					z.Headers[za0001] = nil
				} else {
					if z.Headers[za0001] == nil {
						z.Headers[za0001] = new(Header)
					}
					bts, err = z.Headers[za0001].UnmarshalMsg(bts)
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
func (z *Headers) Msgsize() (s int) {
	s = 1 + 8 + msgp.ArrayHeaderSize
	for za0001 := range z.Headers {
		if z.Headers[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += z.Headers[za0001].Msgsize()
		}
	}
	return
}
