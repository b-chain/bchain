package para_paser

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Arg) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "type":
			z.Type, err = dc.ReadString()
			if err != nil {
				return
			}
		case "val":
			z.Data, err = dc.ReadBytes(z.Data)
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
func (z *Arg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "type"
	err = en.Append(0x82, 0xa4, 0x74, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Type)
	if err != nil {
		return
	}
	// write "val"
	err = en.Append(0xa3, 0x76, 0x61, 0x6c)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Data)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Arg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "type"
	o = append(o, 0x82, 0xa4, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendString(o, z.Type)
	// string "val"
	o = append(o, 0xa3, 0x76, 0x61, 0x6c)
	o = msgp.AppendBytes(o, z.Data)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Arg) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "type":
			z.Type, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "val":
			z.Data, bts, err = msgp.ReadBytesBytes(bts, z.Data)
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
func (z *Arg) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Type) + 4 + msgp.BytesPrefixSize + len(z.Data)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *WasmPara) DecodeMsg(dc *msgp.Reader) (err error) {
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
				z.Args = make([]Arg, zb0002)
			}
			for za0001 := range z.Args {
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
					case "type":
						z.Args[za0001].Type, err = dc.ReadString()
						if err != nil {
							return
						}
					case "val":
						z.Args[za0001].Data, err = dc.ReadBytes(z.Args[za0001].Data)
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
func (z *WasmPara) EncodeMsg(en *msgp.Writer) (err error) {
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
		// map header, size 2
		// write "type"
		err = en.Append(0x82, 0xa4, 0x74, 0x79, 0x70, 0x65)
		if err != nil {
			return
		}
		err = en.WriteString(z.Args[za0001].Type)
		if err != nil {
			return
		}
		// write "val"
		err = en.Append(0xa3, 0x76, 0x61, 0x6c)
		if err != nil {
			return
		}
		err = en.WriteBytes(z.Args[za0001].Data)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *WasmPara) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "funcName"
	o = append(o, 0x82, 0xa8, 0x66, 0x75, 0x6e, 0x63, 0x4e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.FuncName)
	// string "args"
	o = append(o, 0xa4, 0x61, 0x72, 0x67, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Args)))
	for za0001 := range z.Args {
		// map header, size 2
		// string "type"
		o = append(o, 0x82, 0xa4, 0x74, 0x79, 0x70, 0x65)
		o = msgp.AppendString(o, z.Args[za0001].Type)
		// string "val"
		o = append(o, 0xa3, 0x76, 0x61, 0x6c)
		o = msgp.AppendBytes(o, z.Args[za0001].Data)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *WasmPara) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
				z.Args = make([]Arg, zb0002)
			}
			for za0001 := range z.Args {
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
					case "type":
						z.Args[za0001].Type, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					case "val":
						z.Args[za0001].Data, bts, err = msgp.ReadBytesBytes(bts, z.Args[za0001].Data)
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
func (z *WasmPara) Msgsize() (s int) {
	s = 1 + 9 + msgp.StringPrefixSize + len(z.FuncName) + 5 + msgp.ArrayHeaderSize
	for za0001 := range z.Args {
		s += 1 + 5 + msgp.StringPrefixSize + len(z.Args[za0001].Type) + 4 + msgp.BytesPrefixSize + len(z.Args[za0001].Data)
	}
	return
}
