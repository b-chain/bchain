package types

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Hash) DecodeMsg(dc *msgp.Reader) (err error) {
	err = dc.ReadExactBytes((z)[:])
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Hash) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteBytes((z)[:])
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Hash) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendBytes(o, (z)[:])
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Hash) UnmarshalMsg(bts []byte) (o []byte, err error) {
	bts, err = msgp.ReadExactBytes(bts, (z)[:])
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Hash) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize + (HashLength * (msgp.ByteSize))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Hashs) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Hashs":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Hashs) >= int(zb0002) {
				z.Hashs = (z.Hashs)[:zb0002]
			} else {
				z.Hashs = make([]*Hash, zb0002)
			}
			for za0001 := range z.Hashs {
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						return
					}
					z.Hashs[za0001] = nil
				} else {
					if z.Hashs[za0001] == nil {
						z.Hashs[za0001] = new(Hash)
					}
					err = dc.ReadExactBytes((*z.Hashs[za0001])[:])
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
func (z *Hashs) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Hashs"
	err = en.Append(0x81, 0xa5, 0x48, 0x61, 0x73, 0x68, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Hashs)))
	if err != nil {
		return
	}
	for za0001 := range z.Hashs {
		if z.Hashs[za0001] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteBytes((*z.Hashs[za0001])[:])
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Hashs) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Hashs"
	o = append(o, 0x81, 0xa5, 0x48, 0x61, 0x73, 0x68, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Hashs)))
	for za0001 := range z.Hashs {
		if z.Hashs[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendBytes(o, (*z.Hashs[za0001])[:])
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Hashs) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Hashs":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Hashs) >= int(zb0002) {
				z.Hashs = (z.Hashs)[:zb0002]
			} else {
				z.Hashs = make([]*Hash, zb0002)
			}
			for za0001 := range z.Hashs {
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					z.Hashs[za0001] = nil
				} else {
					if z.Hashs[za0001] == nil {
						z.Hashs[za0001] = new(Hash)
					}
					bts, err = msgp.ReadExactBytes(bts, (*z.Hashs[za0001])[:])
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
func (z *Hashs) Msgsize() (s int) {
	s = 1 + 6 + msgp.ArrayHeaderSize
	for za0001 := range z.Hashs {
		if z.Hashs[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += msgp.ArrayHeaderSize + (HashLength * (msgp.ByteSize))
		}
	}
	return
}
