package transaction

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
	"bchain.io/common/types"
)

// DecodeMsg implements msgp.Decodable
func (z *Log) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Address":
			err = z.Address.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Topics":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Topics) >= int(zb0002) {
				z.Topics = (z.Topics)[:zb0002]
			} else {
				z.Topics = make([]types.Hash, zb0002)
			}
			for za0001 := range z.Topics {
				err = z.Topics[za0001].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Data":
			var zb0003 uint32
			zb0003, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Data) >= int(zb0003) {
				z.Data = (z.Data)[:zb0003]
			} else {
				z.Data = make([][]byte, zb0003)
			}
			for za0002 := range z.Data {
				z.Data[za0002], err = dc.ReadBytes(z.Data[za0002])
				if err != nil {
					return
				}
			}
		case "BlockNumber":
			z.BlockNumber, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "TxHash":
			err = z.TxHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "TxIndex":
			z.TxIndex, err = dc.ReadUint()
			if err != nil {
				return
			}
		case "BlockHash":
			err = z.BlockHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Index":
			z.Index, err = dc.ReadUint()
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
func (z *Log) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 8
	// write "Address"
	err = en.Append(0x88, 0xa7, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73)
	if err != nil {
		return
	}
	err = z.Address.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Topics"
	err = en.Append(0xa6, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Topics)))
	if err != nil {
		return
	}
	for za0001 := range z.Topics {
		err = z.Topics[za0001].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Data"
	err = en.Append(0xa4, 0x44, 0x61, 0x74, 0x61)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Data)))
	if err != nil {
		return
	}
	for za0002 := range z.Data {
		err = en.WriteBytes(z.Data[za0002])
		if err != nil {
			return
		}
	}
	// write "BlockNumber"
	err = en.Append(0xab, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.BlockNumber)
	if err != nil {
		return
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
	// write "TxIndex"
	err = en.Append(0xa7, 0x54, 0x78, 0x49, 0x6e, 0x64, 0x65, 0x78)
	if err != nil {
		return
	}
	err = en.WriteUint(z.TxIndex)
	if err != nil {
		return
	}
	// write "BlockHash"
	err = en.Append(0xa9, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.BlockHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Index"
	err = en.Append(0xa5, 0x49, 0x6e, 0x64, 0x65, 0x78)
	if err != nil {
		return
	}
	err = en.WriteUint(z.Index)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Log) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 8
	// string "Address"
	o = append(o, 0x88, 0xa7, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73)
	o, err = z.Address.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Topics"
	o = append(o, 0xa6, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Topics)))
	for za0001 := range z.Topics {
		o, err = z.Topics[za0001].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Data"
	o = append(o, 0xa4, 0x44, 0x61, 0x74, 0x61)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Data)))
	for za0002 := range z.Data {
		o = msgp.AppendBytes(o, z.Data[za0002])
	}
	// string "BlockNumber"
	o = append(o, 0xab, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	o = msgp.AppendUint64(o, z.BlockNumber)
	// string "TxHash"
	o = append(o, 0xa6, 0x54, 0x78, 0x48, 0x61, 0x73, 0x68)
	o, err = z.TxHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "TxIndex"
	o = append(o, 0xa7, 0x54, 0x78, 0x49, 0x6e, 0x64, 0x65, 0x78)
	o = msgp.AppendUint(o, z.TxIndex)
	// string "BlockHash"
	o = append(o, 0xa9, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68)
	o, err = z.BlockHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Index"
	o = append(o, 0xa5, 0x49, 0x6e, 0x64, 0x65, 0x78)
	o = msgp.AppendUint(o, z.Index)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Log) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Address":
			bts, err = z.Address.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Topics":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Topics) >= int(zb0002) {
				z.Topics = (z.Topics)[:zb0002]
			} else {
				z.Topics = make([]types.Hash, zb0002)
			}
			for za0001 := range z.Topics {
				bts, err = z.Topics[za0001].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Data":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Data) >= int(zb0003) {
				z.Data = (z.Data)[:zb0003]
			} else {
				z.Data = make([][]byte, zb0003)
			}
			for za0002 := range z.Data {
				z.Data[za0002], bts, err = msgp.ReadBytesBytes(bts, z.Data[za0002])
				if err != nil {
					return
				}
			}
		case "BlockNumber":
			z.BlockNumber, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "TxHash":
			bts, err = z.TxHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "TxIndex":
			z.TxIndex, bts, err = msgp.ReadUintBytes(bts)
			if err != nil {
				return
			}
		case "BlockHash":
			bts, err = z.BlockHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Index":
			z.Index, bts, err = msgp.ReadUintBytes(bts)
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
func (z *Log) Msgsize() (s int) {
	s = 1 + 8 + z.Address.Msgsize() + 7 + msgp.ArrayHeaderSize
	for za0001 := range z.Topics {
		s += z.Topics[za0001].Msgsize()
	}
	s += 5 + msgp.ArrayHeaderSize
	for za0002 := range z.Data {
		s += msgp.BytesPrefixSize + len(z.Data[za0002])
	}
	s += 12 + msgp.Uint64Size + 7 + z.TxHash.Msgsize() + 8 + msgp.UintSize + 10 + z.BlockHash.Msgsize() + 6 + msgp.UintSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *LogProtocol) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Address":
			err = z.Address.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Topics":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Topics) >= int(zb0002) {
				z.Topics = (z.Topics)[:zb0002]
			} else {
				z.Topics = make([]types.Hash, zb0002)
			}
			for za0001 := range z.Topics {
				err = z.Topics[za0001].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Data":
			var zb0003 uint32
			zb0003, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Data) >= int(zb0003) {
				z.Data = (z.Data)[:zb0003]
			} else {
				z.Data = make([][]byte, zb0003)
			}
			for za0002 := range z.Data {
				z.Data[za0002], err = dc.ReadBytes(z.Data[za0002])
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
func (z *LogProtocol) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Address"
	err = en.Append(0x83, 0xa7, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73)
	if err != nil {
		return
	}
	err = z.Address.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Topics"
	err = en.Append(0xa6, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Topics)))
	if err != nil {
		return
	}
	for za0001 := range z.Topics {
		err = z.Topics[za0001].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Data"
	err = en.Append(0xa4, 0x44, 0x61, 0x74, 0x61)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Data)))
	if err != nil {
		return
	}
	for za0002 := range z.Data {
		err = en.WriteBytes(z.Data[za0002])
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *LogProtocol) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Address"
	o = append(o, 0x83, 0xa7, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73)
	o, err = z.Address.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Topics"
	o = append(o, 0xa6, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Topics)))
	for za0001 := range z.Topics {
		o, err = z.Topics[za0001].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Data"
	o = append(o, 0xa4, 0x44, 0x61, 0x74, 0x61)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Data)))
	for za0002 := range z.Data {
		o = msgp.AppendBytes(o, z.Data[za0002])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *LogProtocol) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Address":
			bts, err = z.Address.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Topics":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Topics) >= int(zb0002) {
				z.Topics = (z.Topics)[:zb0002]
			} else {
				z.Topics = make([]types.Hash, zb0002)
			}
			for za0001 := range z.Topics {
				bts, err = z.Topics[za0001].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Data":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Data) >= int(zb0003) {
				z.Data = (z.Data)[:zb0003]
			} else {
				z.Data = make([][]byte, zb0003)
			}
			for za0002 := range z.Data {
				z.Data[za0002], bts, err = msgp.ReadBytesBytes(bts, z.Data[za0002])
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
func (z *LogProtocol) Msgsize() (s int) {
	s = 1 + 8 + z.Address.Msgsize() + 7 + msgp.ArrayHeaderSize
	for za0001 := range z.Topics {
		s += z.Topics[za0001].Msgsize()
	}
	s += 5 + msgp.ArrayHeaderSize
	for za0002 := range z.Data {
		s += msgp.BytesPrefixSize + len(z.Data[za0002])
	}
	return
}
