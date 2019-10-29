package discv5

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Findnode) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Target":
			err = z.Target.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Expiration":
			z.Expiration, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Rest":
			z.Rest, err = dc.ReadBytes(z.Rest)
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
func (z *Findnode) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Target"
	err = en.Append(0x83, 0xa6, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74)
	if err != nil {
		return
	}
	err = z.Target.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Expiration"
	err = en.Append(0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Expiration)
	if err != nil {
		return
	}
	// write "Rest"
	err = en.Append(0xa4, 0x52, 0x65, 0x73, 0x74)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Rest)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Findnode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Target"
	o = append(o, 0x83, 0xa6, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74)
	o, err = z.Target.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Expiration"
	o = append(o, 0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint64(o, z.Expiration)
	// string "Rest"
	o = append(o, 0xa4, 0x52, 0x65, 0x73, 0x74)
	o = msgp.AppendBytes(o, z.Rest)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Findnode) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Target":
			bts, err = z.Target.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Expiration":
			z.Expiration, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Rest":
			z.Rest, bts, err = msgp.ReadBytesBytes(bts, z.Rest)
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
func (z *Findnode) Msgsize() (s int) {
	s = 1 + 7 + z.Target.Msgsize() + 11 + msgp.Uint64Size + 5 + msgp.BytesPrefixSize + len(z.Rest)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FindnodeHash) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Target":
			err = z.Target.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Expiration":
			z.Expiration, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Rest":
			z.Rest, err = dc.ReadBytes(z.Rest)
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
func (z *FindnodeHash) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Target"
	err = en.Append(0x83, 0xa6, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74)
	if err != nil {
		return
	}
	err = z.Target.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Expiration"
	err = en.Append(0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Expiration)
	if err != nil {
		return
	}
	// write "Rest"
	err = en.Append(0xa4, 0x52, 0x65, 0x73, 0x74)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Rest)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FindnodeHash) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Target"
	o = append(o, 0x83, 0xa6, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74)
	o, err = z.Target.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Expiration"
	o = append(o, 0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint64(o, z.Expiration)
	// string "Rest"
	o = append(o, 0xa4, 0x52, 0x65, 0x73, 0x74)
	o = msgp.AppendBytes(o, z.Rest)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FindnodeHash) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Target":
			bts, err = z.Target.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Expiration":
			z.Expiration, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Rest":
			z.Rest, bts, err = msgp.ReadBytesBytes(bts, z.Rest)
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
func (z *FindnodeHash) Msgsize() (s int) {
	s = 1 + 7 + z.Target.Msgsize() + 11 + msgp.Uint64Size + 5 + msgp.BytesPrefixSize + len(z.Rest)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Neighbors) DecodeMsg(dc *msgp.Reader) (err error) {
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
				z.Nodes = make([]RpcNode, zb0002)
			}
			for za0001 := range z.Nodes {
				err = z.Nodes[za0001].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Expiration":
			z.Expiration, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Rest":
			z.Rest, err = dc.ReadBytes(z.Rest)
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
func (z *Neighbors) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Nodes"
	err = en.Append(0x83, 0xa5, 0x4e, 0x6f, 0x64, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Nodes)))
	if err != nil {
		return
	}
	for za0001 := range z.Nodes {
		err = z.Nodes[za0001].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Expiration"
	err = en.Append(0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Expiration)
	if err != nil {
		return
	}
	// write "Rest"
	err = en.Append(0xa4, 0x52, 0x65, 0x73, 0x74)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Rest)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Neighbors) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Nodes"
	o = append(o, 0x83, 0xa5, 0x4e, 0x6f, 0x64, 0x65, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Nodes)))
	for za0001 := range z.Nodes {
		o, err = z.Nodes[za0001].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Expiration"
	o = append(o, 0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint64(o, z.Expiration)
	// string "Rest"
	o = append(o, 0xa4, 0x52, 0x65, 0x73, 0x74)
	o = msgp.AppendBytes(o, z.Rest)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Neighbors) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
				z.Nodes = make([]RpcNode, zb0002)
			}
			for za0001 := range z.Nodes {
				bts, err = z.Nodes[za0001].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Expiration":
			z.Expiration, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Rest":
			z.Rest, bts, err = msgp.ReadBytesBytes(bts, z.Rest)
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
func (z *Neighbors) Msgsize() (s int) {
	s = 1 + 6 + msgp.ArrayHeaderSize
	for za0001 := range z.Nodes {
		s += z.Nodes[za0001].Msgsize()
	}
	s += 11 + msgp.Uint64Size + 5 + msgp.BytesPrefixSize + len(z.Rest)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Ping) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Version":
			z.Version, err = dc.ReadUint()
			if err != nil {
				return
			}
		case "From":
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
				case "IP":
					err = z.From.IP.DecodeMsg(dc)
					if err != nil {
						return
					}
				case "UDP":
					z.From.UDP, err = dc.ReadUint16()
					if err != nil {
						return
					}
				case "TCP":
					z.From.TCP, err = dc.ReadUint16()
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
		case "To":
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
				case "IP":
					err = z.To.IP.DecodeMsg(dc)
					if err != nil {
						return
					}
				case "UDP":
					z.To.UDP, err = dc.ReadUint16()
					if err != nil {
						return
					}
				case "TCP":
					z.To.TCP, err = dc.ReadUint16()
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
		case "Expiration":
			z.Expiration, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Topics":
			var zb0004 uint32
			zb0004, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Topics) >= int(zb0004) {
				z.Topics = (z.Topics)[:zb0004]
			} else {
				z.Topics = make([]Topic, zb0004)
			}
			for za0001 := range z.Topics {
				err = z.Topics[za0001].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Rest":
			z.Rest, err = dc.ReadBytes(z.Rest)
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
func (z *Ping) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 6
	// write "Version"
	err = en.Append(0x86, 0xa7, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint(z.Version)
	if err != nil {
		return
	}
	// write "From"
	// map header, size 3
	// write "IP"
	err = en.Append(0xa4, 0x46, 0x72, 0x6f, 0x6d, 0x83, 0xa2, 0x49, 0x50)
	if err != nil {
		return
	}
	err = z.From.IP.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "UDP"
	err = en.Append(0xa3, 0x55, 0x44, 0x50)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.From.UDP)
	if err != nil {
		return
	}
	// write "TCP"
	err = en.Append(0xa3, 0x54, 0x43, 0x50)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.From.TCP)
	if err != nil {
		return
	}
	// write "To"
	// map header, size 3
	// write "IP"
	err = en.Append(0xa2, 0x54, 0x6f, 0x83, 0xa2, 0x49, 0x50)
	if err != nil {
		return
	}
	err = z.To.IP.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "UDP"
	err = en.Append(0xa3, 0x55, 0x44, 0x50)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.To.UDP)
	if err != nil {
		return
	}
	// write "TCP"
	err = en.Append(0xa3, 0x54, 0x43, 0x50)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.To.TCP)
	if err != nil {
		return
	}
	// write "Expiration"
	err = en.Append(0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Expiration)
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
	// write "Rest"
	err = en.Append(0xa4, 0x52, 0x65, 0x73, 0x74)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Rest)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Ping) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 6
	// string "Version"
	o = append(o, 0x86, 0xa7, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint(o, z.Version)
	// string "From"
	// map header, size 3
	// string "IP"
	o = append(o, 0xa4, 0x46, 0x72, 0x6f, 0x6d, 0x83, 0xa2, 0x49, 0x50)
	o, err = z.From.IP.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "UDP"
	o = append(o, 0xa3, 0x55, 0x44, 0x50)
	o = msgp.AppendUint16(o, z.From.UDP)
	// string "TCP"
	o = append(o, 0xa3, 0x54, 0x43, 0x50)
	o = msgp.AppendUint16(o, z.From.TCP)
	// string "To"
	// map header, size 3
	// string "IP"
	o = append(o, 0xa2, 0x54, 0x6f, 0x83, 0xa2, 0x49, 0x50)
	o, err = z.To.IP.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "UDP"
	o = append(o, 0xa3, 0x55, 0x44, 0x50)
	o = msgp.AppendUint16(o, z.To.UDP)
	// string "TCP"
	o = append(o, 0xa3, 0x54, 0x43, 0x50)
	o = msgp.AppendUint16(o, z.To.TCP)
	// string "Expiration"
	o = append(o, 0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint64(o, z.Expiration)
	// string "Topics"
	o = append(o, 0xa6, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Topics)))
	for za0001 := range z.Topics {
		o, err = z.Topics[za0001].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Rest"
	o = append(o, 0xa4, 0x52, 0x65, 0x73, 0x74)
	o = msgp.AppendBytes(o, z.Rest)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Ping) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Version":
			z.Version, bts, err = msgp.ReadUintBytes(bts)
			if err != nil {
				return
			}
		case "From":
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
				case "IP":
					bts, err = z.From.IP.UnmarshalMsg(bts)
					if err != nil {
						return
					}
				case "UDP":
					z.From.UDP, bts, err = msgp.ReadUint16Bytes(bts)
					if err != nil {
						return
					}
				case "TCP":
					z.From.TCP, bts, err = msgp.ReadUint16Bytes(bts)
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
		case "To":
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
				case "IP":
					bts, err = z.To.IP.UnmarshalMsg(bts)
					if err != nil {
						return
					}
				case "UDP":
					z.To.UDP, bts, err = msgp.ReadUint16Bytes(bts)
					if err != nil {
						return
					}
				case "TCP":
					z.To.TCP, bts, err = msgp.ReadUint16Bytes(bts)
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
		case "Expiration":
			z.Expiration, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Topics":
			var zb0004 uint32
			zb0004, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Topics) >= int(zb0004) {
				z.Topics = (z.Topics)[:zb0004]
			} else {
				z.Topics = make([]Topic, zb0004)
			}
			for za0001 := range z.Topics {
				bts, err = z.Topics[za0001].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Rest":
			z.Rest, bts, err = msgp.ReadBytesBytes(bts, z.Rest)
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
func (z *Ping) Msgsize() (s int) {
	s = 1 + 8 + msgp.UintSize + 5 + 1 + 3 + z.From.IP.Msgsize() + 4 + msgp.Uint16Size + 4 + msgp.Uint16Size + 3 + 1 + 3 + z.To.IP.Msgsize() + 4 + msgp.Uint16Size + 4 + msgp.Uint16Size + 11 + msgp.Uint64Size + 7 + msgp.ArrayHeaderSize
	for za0001 := range z.Topics {
		s += z.Topics[za0001].Msgsize()
	}
	s += 5 + msgp.BytesPrefixSize + len(z.Rest)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Pong) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "To":
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
				case "IP":
					err = z.To.IP.DecodeMsg(dc)
					if err != nil {
						return
					}
				case "UDP":
					z.To.UDP, err = dc.ReadUint16()
					if err != nil {
						return
					}
				case "TCP":
					z.To.TCP, err = dc.ReadUint16()
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
		case "ReplyTok":
			z.ReplyTok, err = dc.ReadBytes(z.ReplyTok)
			if err != nil {
				return
			}
		case "Expiration":
			z.Expiration, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "TopicHash":
			err = z.TopicHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "TicketSerial":
			z.TicketSerial, err = dc.ReadUint32()
			if err != nil {
				return
			}
		case "WaitPeriods":
			var zb0003 uint32
			zb0003, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.WaitPeriods) >= int(zb0003) {
				z.WaitPeriods = (z.WaitPeriods)[:zb0003]
			} else {
				z.WaitPeriods = make([]uint32, zb0003)
			}
			for za0001 := range z.WaitPeriods {
				z.WaitPeriods[za0001], err = dc.ReadUint32()
				if err != nil {
					return
				}
			}
		case "Rest":
			z.Rest, err = dc.ReadBytes(z.Rest)
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
func (z *Pong) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 7
	// write "To"
	// map header, size 3
	// write "IP"
	err = en.Append(0x87, 0xa2, 0x54, 0x6f, 0x83, 0xa2, 0x49, 0x50)
	if err != nil {
		return
	}
	err = z.To.IP.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "UDP"
	err = en.Append(0xa3, 0x55, 0x44, 0x50)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.To.UDP)
	if err != nil {
		return
	}
	// write "TCP"
	err = en.Append(0xa3, 0x54, 0x43, 0x50)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.To.TCP)
	if err != nil {
		return
	}
	// write "ReplyTok"
	err = en.Append(0xa8, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x6b)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.ReplyTok)
	if err != nil {
		return
	}
	// write "Expiration"
	err = en.Append(0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Expiration)
	if err != nil {
		return
	}
	// write "TopicHash"
	err = en.Append(0xa9, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.TopicHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "TicketSerial"
	err = en.Append(0xac, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.TicketSerial)
	if err != nil {
		return
	}
	// write "WaitPeriods"
	err = en.Append(0xab, 0x57, 0x61, 0x69, 0x74, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.WaitPeriods)))
	if err != nil {
		return
	}
	for za0001 := range z.WaitPeriods {
		err = en.WriteUint32(z.WaitPeriods[za0001])
		if err != nil {
			return
		}
	}
	// write "Rest"
	err = en.Append(0xa4, 0x52, 0x65, 0x73, 0x74)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Rest)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Pong) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 7
	// string "To"
	// map header, size 3
	// string "IP"
	o = append(o, 0x87, 0xa2, 0x54, 0x6f, 0x83, 0xa2, 0x49, 0x50)
	o, err = z.To.IP.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "UDP"
	o = append(o, 0xa3, 0x55, 0x44, 0x50)
	o = msgp.AppendUint16(o, z.To.UDP)
	// string "TCP"
	o = append(o, 0xa3, 0x54, 0x43, 0x50)
	o = msgp.AppendUint16(o, z.To.TCP)
	// string "ReplyTok"
	o = append(o, 0xa8, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x6b)
	o = msgp.AppendBytes(o, z.ReplyTok)
	// string "Expiration"
	o = append(o, 0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint64(o, z.Expiration)
	// string "TopicHash"
	o = append(o, 0xa9, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x48, 0x61, 0x73, 0x68)
	o, err = z.TopicHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "TicketSerial"
	o = append(o, 0xac, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c)
	o = msgp.AppendUint32(o, z.TicketSerial)
	// string "WaitPeriods"
	o = append(o, 0xab, 0x57, 0x61, 0x69, 0x74, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.WaitPeriods)))
	for za0001 := range z.WaitPeriods {
		o = msgp.AppendUint32(o, z.WaitPeriods[za0001])
	}
	// string "Rest"
	o = append(o, 0xa4, 0x52, 0x65, 0x73, 0x74)
	o = msgp.AppendBytes(o, z.Rest)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Pong) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "To":
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
				case "IP":
					bts, err = z.To.IP.UnmarshalMsg(bts)
					if err != nil {
						return
					}
				case "UDP":
					z.To.UDP, bts, err = msgp.ReadUint16Bytes(bts)
					if err != nil {
						return
					}
				case "TCP":
					z.To.TCP, bts, err = msgp.ReadUint16Bytes(bts)
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
		case "ReplyTok":
			z.ReplyTok, bts, err = msgp.ReadBytesBytes(bts, z.ReplyTok)
			if err != nil {
				return
			}
		case "Expiration":
			z.Expiration, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "TopicHash":
			bts, err = z.TopicHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "TicketSerial":
			z.TicketSerial, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				return
			}
		case "WaitPeriods":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.WaitPeriods) >= int(zb0003) {
				z.WaitPeriods = (z.WaitPeriods)[:zb0003]
			} else {
				z.WaitPeriods = make([]uint32, zb0003)
			}
			for za0001 := range z.WaitPeriods {
				z.WaitPeriods[za0001], bts, err = msgp.ReadUint32Bytes(bts)
				if err != nil {
					return
				}
			}
		case "Rest":
			z.Rest, bts, err = msgp.ReadBytesBytes(bts, z.Rest)
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
func (z *Pong) Msgsize() (s int) {
	s = 1 + 3 + 1 + 3 + z.To.IP.Msgsize() + 4 + msgp.Uint16Size + 4 + msgp.Uint16Size + 9 + msgp.BytesPrefixSize + len(z.ReplyTok) + 11 + msgp.Uint64Size + 10 + z.TopicHash.Msgsize() + 13 + msgp.Uint32Size + 12 + msgp.ArrayHeaderSize + (len(z.WaitPeriods) * (msgp.Uint32Size)) + 5 + msgp.BytesPrefixSize + len(z.Rest)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *RpcEndpoint) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "IP":
			err = z.IP.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "UDP":
			z.UDP, err = dc.ReadUint16()
			if err != nil {
				return
			}
		case "TCP":
			z.TCP, err = dc.ReadUint16()
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
func (z *RpcEndpoint) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "IP"
	err = en.Append(0x83, 0xa2, 0x49, 0x50)
	if err != nil {
		return
	}
	err = z.IP.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "UDP"
	err = en.Append(0xa3, 0x55, 0x44, 0x50)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.UDP)
	if err != nil {
		return
	}
	// write "TCP"
	err = en.Append(0xa3, 0x54, 0x43, 0x50)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.TCP)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *RpcEndpoint) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "IP"
	o = append(o, 0x83, 0xa2, 0x49, 0x50)
	o, err = z.IP.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "UDP"
	o = append(o, 0xa3, 0x55, 0x44, 0x50)
	o = msgp.AppendUint16(o, z.UDP)
	// string "TCP"
	o = append(o, 0xa3, 0x54, 0x43, 0x50)
	o = msgp.AppendUint16(o, z.TCP)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *RpcEndpoint) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "IP":
			bts, err = z.IP.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "UDP":
			z.UDP, bts, err = msgp.ReadUint16Bytes(bts)
			if err != nil {
				return
			}
		case "TCP":
			z.TCP, bts, err = msgp.ReadUint16Bytes(bts)
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
func (z *RpcEndpoint) Msgsize() (s int) {
	s = 1 + 3 + z.IP.Msgsize() + 4 + msgp.Uint16Size + 4 + msgp.Uint16Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *RpcNode) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "IP":
			err = z.IP.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "UDP":
			z.UDP, err = dc.ReadUint16()
			if err != nil {
				return
			}
		case "TCP":
			z.TCP, err = dc.ReadUint16()
			if err != nil {
				return
			}
		case "ID":
			err = z.ID.DecodeMsg(dc)
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
func (z *RpcNode) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "IP"
	err = en.Append(0x84, 0xa2, 0x49, 0x50)
	if err != nil {
		return
	}
	err = z.IP.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "UDP"
	err = en.Append(0xa3, 0x55, 0x44, 0x50)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.UDP)
	if err != nil {
		return
	}
	// write "TCP"
	err = en.Append(0xa3, 0x54, 0x43, 0x50)
	if err != nil {
		return
	}
	err = en.WriteUint16(z.TCP)
	if err != nil {
		return
	}
	// write "ID"
	err = en.Append(0xa2, 0x49, 0x44)
	if err != nil {
		return
	}
	err = z.ID.EncodeMsg(en)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *RpcNode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "IP"
	o = append(o, 0x84, 0xa2, 0x49, 0x50)
	o, err = z.IP.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "UDP"
	o = append(o, 0xa3, 0x55, 0x44, 0x50)
	o = msgp.AppendUint16(o, z.UDP)
	// string "TCP"
	o = append(o, 0xa3, 0x54, 0x43, 0x50)
	o = msgp.AppendUint16(o, z.TCP)
	// string "ID"
	o = append(o, 0xa2, 0x49, 0x44)
	o, err = z.ID.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *RpcNode) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "IP":
			bts, err = z.IP.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "UDP":
			z.UDP, bts, err = msgp.ReadUint16Bytes(bts)
			if err != nil {
				return
			}
		case "TCP":
			z.TCP, bts, err = msgp.ReadUint16Bytes(bts)
			if err != nil {
				return
			}
		case "ID":
			bts, err = z.ID.UnmarshalMsg(bts)
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
func (z *RpcNode) Msgsize() (s int) {
	s = 1 + 3 + z.IP.Msgsize() + 4 + msgp.Uint16Size + 4 + msgp.Uint16Size + 3 + z.ID.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TopicNodes) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Echo":
			err = z.Echo.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Nodes":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Nodes) >= int(zb0002) {
				z.Nodes = (z.Nodes)[:zb0002]
			} else {
				z.Nodes = make([]RpcNode, zb0002)
			}
			for za0001 := range z.Nodes {
				err = z.Nodes[za0001].DecodeMsg(dc)
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
func (z *TopicNodes) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "Echo"
	err = en.Append(0x82, 0xa4, 0x45, 0x63, 0x68, 0x6f)
	if err != nil {
		return
	}
	err = z.Echo.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Nodes"
	err = en.Append(0xa5, 0x4e, 0x6f, 0x64, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Nodes)))
	if err != nil {
		return
	}
	for za0001 := range z.Nodes {
		err = z.Nodes[za0001].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TopicNodes) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Echo"
	o = append(o, 0x82, 0xa4, 0x45, 0x63, 0x68, 0x6f)
	o, err = z.Echo.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Nodes"
	o = append(o, 0xa5, 0x4e, 0x6f, 0x64, 0x65, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Nodes)))
	for za0001 := range z.Nodes {
		o, err = z.Nodes[za0001].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TopicNodes) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Echo":
			bts, err = z.Echo.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Nodes":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Nodes) >= int(zb0002) {
				z.Nodes = (z.Nodes)[:zb0002]
			} else {
				z.Nodes = make([]RpcNode, zb0002)
			}
			for za0001 := range z.Nodes {
				bts, err = z.Nodes[za0001].UnmarshalMsg(bts)
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
func (z *TopicNodes) Msgsize() (s int) {
	s = 1 + 5 + z.Echo.Msgsize() + 6 + msgp.ArrayHeaderSize
	for za0001 := range z.Nodes {
		s += z.Nodes[za0001].Msgsize()
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TopicQuery) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Topic":
			err = z.Topic.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Expiration":
			z.Expiration, err = dc.ReadUint64()
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
func (z *TopicQuery) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "Topic"
	err = en.Append(0x82, 0xa5, 0x54, 0x6f, 0x70, 0x69, 0x63)
	if err != nil {
		return
	}
	err = z.Topic.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Expiration"
	err = en.Append(0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Expiration)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TopicQuery) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Topic"
	o = append(o, 0x82, 0xa5, 0x54, 0x6f, 0x70, 0x69, 0x63)
	o, err = z.Topic.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Expiration"
	o = append(o, 0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint64(o, z.Expiration)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TopicQuery) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Topic":
			bts, err = z.Topic.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Expiration":
			z.Expiration, bts, err = msgp.ReadUint64Bytes(bts)
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
func (z *TopicQuery) Msgsize() (s int) {
	s = 1 + 6 + z.Topic.Msgsize() + 11 + msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TopicRegister) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Topics":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Topics) >= int(zb0002) {
				z.Topics = (z.Topics)[:zb0002]
			} else {
				z.Topics = make([]Topic, zb0002)
			}
			for za0001 := range z.Topics {
				err = z.Topics[za0001].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Idx":
			z.Idx, err = dc.ReadUint()
			if err != nil {
				return
			}
		case "Pong":
			z.Pong, err = dc.ReadBytes(z.Pong)
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
func (z *TopicRegister) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Topics"
	err = en.Append(0x83, 0xa6, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73)
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
	// write "Idx"
	err = en.Append(0xa3, 0x49, 0x64, 0x78)
	if err != nil {
		return
	}
	err = en.WriteUint(z.Idx)
	if err != nil {
		return
	}
	// write "Pong"
	err = en.Append(0xa4, 0x50, 0x6f, 0x6e, 0x67)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Pong)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TopicRegister) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Topics"
	o = append(o, 0x83, 0xa6, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Topics)))
	for za0001 := range z.Topics {
		o, err = z.Topics[za0001].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Idx"
	o = append(o, 0xa3, 0x49, 0x64, 0x78)
	o = msgp.AppendUint(o, z.Idx)
	// string "Pong"
	o = append(o, 0xa4, 0x50, 0x6f, 0x6e, 0x67)
	o = msgp.AppendBytes(o, z.Pong)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TopicRegister) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Topics":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Topics) >= int(zb0002) {
				z.Topics = (z.Topics)[:zb0002]
			} else {
				z.Topics = make([]Topic, zb0002)
			}
			for za0001 := range z.Topics {
				bts, err = z.Topics[za0001].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Idx":
			z.Idx, bts, err = msgp.ReadUintBytes(bts)
			if err != nil {
				return
			}
		case "Pong":
			z.Pong, bts, err = msgp.ReadBytesBytes(bts, z.Pong)
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
func (z *TopicRegister) Msgsize() (s int) {
	s = 1 + 7 + msgp.ArrayHeaderSize
	for za0001 := range z.Topics {
		s += z.Topics[za0001].Msgsize()
	}
	s += 4 + msgp.UintSize + 5 + msgp.BytesPrefixSize + len(z.Pong)
	return
}
