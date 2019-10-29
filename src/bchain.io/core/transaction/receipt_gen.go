package transaction

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
	"bchain.io/common/types"
)

// DecodeMsg implements msgp.Decodable
func (z *Receipt) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Status":
			z.Status, err = dc.ReadUint()
			if err != nil {
				return
			}
		case "Bloom":
			err = z.Bloom.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Logs":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Logs) >= int(zb0002) {
				z.Logs = (z.Logs)[:zb0002]
			} else {
				z.Logs = make([]*Log, zb0002)
			}
			for za0001 := range z.Logs {
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						return
					}
					z.Logs[za0001] = nil
				} else {
					if z.Logs[za0001] == nil {
						z.Logs[za0001] = new(Log)
					}
					err = z.Logs[za0001].DecodeMsg(dc)
					if err != nil {
						return
					}
				}
			}
		case "TxHash":
			err = z.TxHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "ContractAddress":
			var zb0003 uint32
			zb0003, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.ContractAddress) >= int(zb0003) {
				z.ContractAddress = (z.ContractAddress)[:zb0003]
			} else {
				z.ContractAddress = make([]types.Address, zb0003)
			}
			for za0002 := range z.ContractAddress {
				err = z.ContractAddress[za0002].DecodeMsg(dc)
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
func (z *Receipt) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "Status"
	err = en.Append(0x85, 0xa6, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint(z.Status)
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
	// write "Logs"
	err = en.Append(0xa4, 0x4c, 0x6f, 0x67, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Logs)))
	if err != nil {
		return
	}
	for za0001 := range z.Logs {
		if z.Logs[za0001] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z.Logs[za0001].EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	// write "TxHash"
	err = en.Append(0xa6, 0x54, 0x78, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.TxHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "ContractAddress"
	err = en.Append(0xaf, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.ContractAddress)))
	if err != nil {
		return
	}
	for za0002 := range z.ContractAddress {
		err = z.ContractAddress[za0002].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Receipt) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "Status"
	o = append(o, 0x85, 0xa6, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73)
	o = msgp.AppendUint(o, z.Status)
	// string "Bloom"
	o = append(o, 0xa5, 0x42, 0x6c, 0x6f, 0x6f, 0x6d)
	o, err = z.Bloom.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Logs"
	o = append(o, 0xa4, 0x4c, 0x6f, 0x67, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Logs)))
	for za0001 := range z.Logs {
		if z.Logs[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.Logs[za0001].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	// string "TxHash"
	o = append(o, 0xa6, 0x54, 0x78, 0x48, 0x61, 0x73, 0x68)
	o, err = z.TxHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "ContractAddress"
	o = append(o, 0xaf, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.ContractAddress)))
	for za0002 := range z.ContractAddress {
		o, err = z.ContractAddress[za0002].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Receipt) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Status":
			z.Status, bts, err = msgp.ReadUintBytes(bts)
			if err != nil {
				return
			}
		case "Bloom":
			bts, err = z.Bloom.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Logs":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Logs) >= int(zb0002) {
				z.Logs = (z.Logs)[:zb0002]
			} else {
				z.Logs = make([]*Log, zb0002)
			}
			for za0001 := range z.Logs {
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					z.Logs[za0001] = nil
				} else {
					if z.Logs[za0001] == nil {
						z.Logs[za0001] = new(Log)
					}
					bts, err = z.Logs[za0001].UnmarshalMsg(bts)
					if err != nil {
						return
					}
				}
			}
		case "TxHash":
			bts, err = z.TxHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "ContractAddress":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.ContractAddress) >= int(zb0003) {
				z.ContractAddress = (z.ContractAddress)[:zb0003]
			} else {
				z.ContractAddress = make([]types.Address, zb0003)
			}
			for za0002 := range z.ContractAddress {
				bts, err = z.ContractAddress[za0002].UnmarshalMsg(bts)
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
func (z *Receipt) Msgsize() (s int) {
	s = 1 + 7 + msgp.UintSize + 6 + z.Bloom.Msgsize() + 5 + msgp.ArrayHeaderSize
	for za0001 := range z.Logs {
		if z.Logs[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += z.Logs[za0001].Msgsize()
		}
	}
	s += 7 + z.TxHash.Msgsize() + 16 + msgp.ArrayHeaderSize
	for za0002 := range z.ContractAddress {
		s += z.ContractAddress[za0002].Msgsize()
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ReceiptProtocol) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Status":
			z.Status, err = dc.ReadUint()
			if err != nil {
				return
			}
		case "Bloom":
			err = z.Bloom.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Logs":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Logs) >= int(zb0002) {
				z.Logs = (z.Logs)[:zb0002]
			} else {
				z.Logs = make([]*LogProtocol, zb0002)
			}
			for za0001 := range z.Logs {
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						return
					}
					z.Logs[za0001] = nil
				} else {
					if z.Logs[za0001] == nil {
						z.Logs[za0001] = new(LogProtocol)
					}
					err = z.Logs[za0001].DecodeMsg(dc)
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
func (z *ReceiptProtocol) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Status"
	err = en.Append(0x83, 0xa6, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint(z.Status)
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
	// write "Logs"
	err = en.Append(0xa4, 0x4c, 0x6f, 0x67, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Logs)))
	if err != nil {
		return
	}
	for za0001 := range z.Logs {
		if z.Logs[za0001] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z.Logs[za0001].EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ReceiptProtocol) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Status"
	o = append(o, 0x83, 0xa6, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73)
	o = msgp.AppendUint(o, z.Status)
	// string "Bloom"
	o = append(o, 0xa5, 0x42, 0x6c, 0x6f, 0x6f, 0x6d)
	o, err = z.Bloom.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Logs"
	o = append(o, 0xa4, 0x4c, 0x6f, 0x67, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Logs)))
	for za0001 := range z.Logs {
		if z.Logs[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.Logs[za0001].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ReceiptProtocol) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Status":
			z.Status, bts, err = msgp.ReadUintBytes(bts)
			if err != nil {
				return
			}
		case "Bloom":
			bts, err = z.Bloom.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Logs":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Logs) >= int(zb0002) {
				z.Logs = (z.Logs)[:zb0002]
			} else {
				z.Logs = make([]*LogProtocol, zb0002)
			}
			for za0001 := range z.Logs {
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					z.Logs[za0001] = nil
				} else {
					if z.Logs[za0001] == nil {
						z.Logs[za0001] = new(LogProtocol)
					}
					bts, err = z.Logs[za0001].UnmarshalMsg(bts)
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
func (z *ReceiptProtocol) Msgsize() (s int) {
	s = 1 + 7 + msgp.UintSize + 6 + z.Bloom.Msgsize() + 5 + msgp.ArrayHeaderSize
	for za0001 := range z.Logs {
		if z.Logs[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += z.Logs[za0001].Msgsize()
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ReceiptProtocols) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(ReceiptProtocols, zb0002)
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
				(*z)[zb0001] = new(ReceiptProtocol)
			}
			err = (*z)[zb0001].DecodeMsg(dc)
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z ReceiptProtocols) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zb0003 := range z {
		if z[zb0003] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z[zb0003].EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z ReceiptProtocols) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zb0003 := range z {
		if z[zb0003] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z[zb0003].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ReceiptProtocols) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(ReceiptProtocols, zb0002)
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
				(*z)[zb0001] = new(ReceiptProtocol)
			}
			bts, err = (*z)[zb0001].UnmarshalMsg(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ReceiptProtocols) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0003 := range z {
		if z[zb0003] == nil {
			s += msgp.NilSize
		} else {
			s += z[zb0003].Msgsize()
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Receipts) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(Receipts, zb0002)
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
				(*z)[zb0001] = new(Receipt)
			}
			err = (*z)[zb0001].DecodeMsg(dc)
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Receipts) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zb0003 := range z {
		if z[zb0003] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z[zb0003].EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Receipts) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zb0003 := range z {
		if z[zb0003] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z[zb0003].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Receipts) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(Receipts, zb0002)
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
				(*z)[zb0001] = new(Receipt)
			}
			bts, err = (*z)[zb0001].UnmarshalMsg(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Receipts) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0003 := range z {
		if z[zb0003] == nil {
			s += msgp.NilSize
		} else {
			s += z[zb0003].Msgsize()
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Receipts_s) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Receipts_s":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Receipts_s) >= int(zb0002) {
				z.Receipts_s = (z.Receipts_s)[:zb0002]
			} else {
				z.Receipts_s = make([]ReceiptProtocols, zb0002)
			}
			for za0001 := range z.Receipts_s {
				var zb0003 uint32
				zb0003, err = dc.ReadArrayHeader()
				if err != nil {
					return
				}
				if cap(z.Receipts_s[za0001]) >= int(zb0003) {
					z.Receipts_s[za0001] = (z.Receipts_s[za0001])[:zb0003]
				} else {
					z.Receipts_s[za0001] = make(ReceiptProtocols, zb0003)
				}
				for za0002 := range z.Receipts_s[za0001] {
					if dc.IsNil() {
						err = dc.ReadNil()
						if err != nil {
							return
						}
						z.Receipts_s[za0001][za0002] = nil
					} else {
						if z.Receipts_s[za0001][za0002] == nil {
							z.Receipts_s[za0001][za0002] = new(ReceiptProtocol)
						}
						err = z.Receipts_s[za0001][za0002].DecodeMsg(dc)
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
func (z *Receipts_s) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Receipts_s"
	err = en.Append(0x81, 0xaa, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x73, 0x5f, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Receipts_s)))
	if err != nil {
		return
	}
	for za0001 := range z.Receipts_s {
		err = en.WriteArrayHeader(uint32(len(z.Receipts_s[za0001])))
		if err != nil {
			return
		}
		for za0002 := range z.Receipts_s[za0001] {
			if z.Receipts_s[za0001][za0002] == nil {
				err = en.WriteNil()
				if err != nil {
					return
				}
			} else {
				err = z.Receipts_s[za0001][za0002].EncodeMsg(en)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Receipts_s) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Receipts_s"
	o = append(o, 0x81, 0xaa, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x73, 0x5f, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Receipts_s)))
	for za0001 := range z.Receipts_s {
		o = msgp.AppendArrayHeader(o, uint32(len(z.Receipts_s[za0001])))
		for za0002 := range z.Receipts_s[za0001] {
			if z.Receipts_s[za0001][za0002] == nil {
				o = msgp.AppendNil(o)
			} else {
				o, err = z.Receipts_s[za0001][za0002].MarshalMsg(o)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Receipts_s) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Receipts_s":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Receipts_s) >= int(zb0002) {
				z.Receipts_s = (z.Receipts_s)[:zb0002]
			} else {
				z.Receipts_s = make([]ReceiptProtocols, zb0002)
			}
			for za0001 := range z.Receipts_s {
				var zb0003 uint32
				zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
				if err != nil {
					return
				}
				if cap(z.Receipts_s[za0001]) >= int(zb0003) {
					z.Receipts_s[za0001] = (z.Receipts_s[za0001])[:zb0003]
				} else {
					z.Receipts_s[za0001] = make(ReceiptProtocols, zb0003)
				}
				for za0002 := range z.Receipts_s[za0001] {
					if msgp.IsNil(bts) {
						bts, err = msgp.ReadNilBytes(bts)
						if err != nil {
							return
						}
						z.Receipts_s[za0001][za0002] = nil
					} else {
						if z.Receipts_s[za0001][za0002] == nil {
							z.Receipts_s[za0001][za0002] = new(ReceiptProtocol)
						}
						bts, err = z.Receipts_s[za0001][za0002].UnmarshalMsg(bts)
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
func (z *Receipts_s) Msgsize() (s int) {
	s = 1 + 11 + msgp.ArrayHeaderSize
	for za0001 := range z.Receipts_s {
		s += msgp.ArrayHeaderSize
		for za0002 := range z.Receipts_s[za0001] {
			if z.Receipts_s[za0001][za0002] == nil {
				s += msgp.NilSize
			} else {
				s += z.Receipts_s[za0001][za0002].Msgsize()
			}
		}
	}
	return
}
