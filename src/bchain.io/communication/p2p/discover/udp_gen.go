package discover

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
			err = dc.ReadExtension(&z.Target)
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
func (z Findnode) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "Target"
	err = en.Append(0x82, 0xa6, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74)
	if err != nil {
		return
	}
	err = en.WriteExtension(&z.Target)
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
func (z Findnode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Target"
	o = append(o, 0x82, 0xa6, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74)
	o, err = msgp.AppendExtension(o, &z.Target)
	if err != nil {
		return
	}
	// string "Expiration"
	o = append(o, 0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint64(o, z.Expiration)
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
			bts, err = msgp.ReadExtensionBytes(bts, &z.Target)
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
func (z Findnode) Msgsize() (s int) {
	s = 1 + 7 + msgp.ExtensionPrefixSize + z.Target.Len() + 11 + msgp.Uint64Size
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
	// map header, size 2
	// write "Nodes"
	err = en.Append(0x82, 0xa5, 0x4e, 0x6f, 0x64, 0x65, 0x73)
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
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Neighbors) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Nodes"
	o = append(o, 0x82, 0xa5, 0x4e, 0x6f, 0x64, 0x65, 0x73)
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
	s += 11 + msgp.Uint64Size
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
					err = dc.ReadExtension(&z.From.IP)
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
					err = dc.ReadExtension(&z.To.IP)
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
	// map header, size 4
	// write "Version"
	err = en.Append(0x84, 0xa7, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
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
	err = en.WriteExtension(&z.From.IP)
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
	err = en.WriteExtension(&z.To.IP)
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
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Ping) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Version"
	o = append(o, 0x84, 0xa7, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint(o, z.Version)
	// string "From"
	// map header, size 3
	// string "IP"
	o = append(o, 0xa4, 0x46, 0x72, 0x6f, 0x6d, 0x83, 0xa2, 0x49, 0x50)
	o, err = msgp.AppendExtension(o, &z.From.IP)
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
	o, err = msgp.AppendExtension(o, &z.To.IP)
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
					bts, err = msgp.ReadExtensionBytes(bts, &z.From.IP)
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
					bts, err = msgp.ReadExtensionBytes(bts, &z.To.IP)
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
	s = 1 + 8 + msgp.UintSize + 5 + 1 + 3 + msgp.ExtensionPrefixSize + z.From.IP.Len() + 4 + msgp.Uint16Size + 4 + msgp.Uint16Size + 3 + 1 + 3 + msgp.ExtensionPrefixSize + z.To.IP.Len() + 4 + msgp.Uint16Size + 4 + msgp.Uint16Size + 11 + msgp.Uint64Size
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
					err = dc.ReadExtension(&z.To.IP)
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
	// map header, size 3
	// write "To"
	// map header, size 3
	// write "IP"
	err = en.Append(0x83, 0xa2, 0x54, 0x6f, 0x83, 0xa2, 0x49, 0x50)
	if err != nil {
		return
	}
	err = en.WriteExtension(&z.To.IP)
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
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Pong) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "To"
	// map header, size 3
	// string "IP"
	o = append(o, 0x83, 0xa2, 0x54, 0x6f, 0x83, 0xa2, 0x49, 0x50)
	o, err = msgp.AppendExtension(o, &z.To.IP)
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
					bts, err = msgp.ReadExtensionBytes(bts, &z.To.IP)
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
	s = 1 + 3 + 1 + 3 + msgp.ExtensionPrefixSize + z.To.IP.Len() + 4 + msgp.Uint16Size + 4 + msgp.Uint16Size + 9 + msgp.BytesPrefixSize + len(z.ReplyTok) + 11 + msgp.Uint64Size
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
			err = dc.ReadExtension(&z.IP)
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
func (z RpcEndpoint) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "IP"
	err = en.Append(0x83, 0xa2, 0x49, 0x50)
	if err != nil {
		return
	}
	err = en.WriteExtension(&z.IP)
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
func (z RpcEndpoint) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "IP"
	o = append(o, 0x83, 0xa2, 0x49, 0x50)
	o, err = msgp.AppendExtension(o, &z.IP)
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
			bts, err = msgp.ReadExtensionBytes(bts, &z.IP)
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
func (z RpcEndpoint) Msgsize() (s int) {
	s = 1 + 3 + msgp.ExtensionPrefixSize + z.IP.Len() + 4 + msgp.Uint16Size + 4 + msgp.Uint16Size
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
			err = dc.ReadExtension(&z.IP)
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
			err = dc.ReadExtension(&z.ID)
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
	err = en.WriteExtension(&z.IP)
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
	err = en.WriteExtension(&z.ID)
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
	o, err = msgp.AppendExtension(o, &z.IP)
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
	o, err = msgp.AppendExtension(o, &z.ID)
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
			bts, err = msgp.ReadExtensionBytes(bts, &z.IP)
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
			bts, err = msgp.ReadExtensionBytes(bts, &z.ID)
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
	s = 1 + 3 + msgp.ExtensionPrefixSize + z.IP.Len() + 4 + msgp.Uint16Size + 4 + msgp.Uint16Size + 3 + msgp.ExtensionPrefixSize + z.ID.Len()
	return
}
