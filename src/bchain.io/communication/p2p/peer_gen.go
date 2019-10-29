package p2p

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ProtoHandshake) DecodeMsg(dc *msgp.Reader) (err error) {
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
			z.Version, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Caps":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Caps) >= int(zb0002) {
				z.Caps = (z.Caps)[:zb0002]
			} else {
				z.Caps = make([]Cap, zb0002)
			}
			for za0001 := range z.Caps {
				err = z.Caps[za0001].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "ListenPort":
			z.ListenPort, err = dc.ReadUint64()
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
func (z *ProtoHandshake) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "Version"
	err = en.Append(0x85, 0xa7, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Version)
	if err != nil {
		return
	}
	// write "Name"
	err = en.Append(0xa4, 0x4e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "Caps"
	err = en.Append(0xa4, 0x43, 0x61, 0x70, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Caps)))
	if err != nil {
		return
	}
	for za0001 := range z.Caps {
		err = z.Caps[za0001].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "ListenPort"
	err = en.Append(0xaa, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x50, 0x6f, 0x72, 0x74)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.ListenPort)
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
func (z *ProtoHandshake) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "Version"
	o = append(o, 0x85, 0xa7, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint64(o, z.Version)
	// string "Name"
	o = append(o, 0xa4, 0x4e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "Caps"
	o = append(o, 0xa4, 0x43, 0x61, 0x70, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Caps)))
	for za0001 := range z.Caps {
		o, err = z.Caps[za0001].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "ListenPort"
	o = append(o, 0xaa, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x50, 0x6f, 0x72, 0x74)
	o = msgp.AppendUint64(o, z.ListenPort)
	// string "ID"
	o = append(o, 0xa2, 0x49, 0x44)
	o, err = z.ID.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ProtoHandshake) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
			z.Version, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Caps":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Caps) >= int(zb0002) {
				z.Caps = (z.Caps)[:zb0002]
			} else {
				z.Caps = make([]Cap, zb0002)
			}
			for za0001 := range z.Caps {
				bts, err = z.Caps[za0001].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "ListenPort":
			z.ListenPort, bts, err = msgp.ReadUint64Bytes(bts)
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
func (z *ProtoHandshake) Msgsize() (s int) {
	s = 1 + 8 + msgp.Uint64Size + 5 + msgp.StringPrefixSize + len(z.Name) + 5 + msgp.ArrayHeaderSize
	for za0001 := range z.Caps {
		s += z.Caps[za0001].Msgsize()
	}
	s += 11 + msgp.Uint64Size + 3 + z.ID.Msgsize()
	return
}
