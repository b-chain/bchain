////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The bchain-go Authors.
//
// The bchain-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// @File: hash.go
// @Date: 2018/04/12 10:12:12
////////////////////////////////////////////////////////////////////////////////

package types

import (
	"fmt"
	"math/big"
	"reflect"
)

//go:generate msgp

const (
	HashLength = 32
)

var (
	hashType int8
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data. It's a hex string/num
type Hash [HashLength]byte

func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}
func StringToHash(s string) Hash { return BytesToHash([]byte(s)) }
func BigToHash(b *big.Int) Hash  { return BytesToHash(b.Bytes()) }
func HexToHash(s string) Hash    { return BytesToHash(FromHex(s)) }

// Get the string representation of the underlying hash
func (h Hash) Str() string   { return string(h[:]) }
func (h Hash) Bytes() []byte { return h[:] }
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }
func (h Hash) Hex() string   { return EncodeHex(h[:]) }

// Here, we'll pick an arbitrary number between
// 0 and 127 that isn't already in use
func (*Hash) ExtensionType() int8 {
	return hashType
}

func (h *Hash)Equal(in *Hash)bool{
	for i := 0;i < HashLength;i++{
		if h[i] != in[i]{
			return false
		}
	}
	return true
}

// We'll always use 32 bytes to encode the data
func (*Hash) Len() int {
	return HashLength
}

// MarshalBinaryTo simply copies the value
// of the bytes into 'b'
func (h *Hash) MarshalBinaryTo(b []byte) error {
	copy(b, h.Bytes())
	return nil
}

func (h Hash) TerminalString() string {
	return fmt.Sprintf("%x…%x", h[:3], h[29:])
}

func (h Hash) String() string {
	return h.Hex()
}


// UnmarshalBinary copies the value of 'b'
// into the Hash object. (We might want to add
// a sanity check here later that len(b) <= HashLength.)
func (h *Hash) UnmarshalBinary(b []byte) error {
	// TODO: check b, only hex, len <= HashLength
	if len(b) <= HashLength {
		*h = BytesToHash(b)
		return nil
	}

	return ErrBytesTooLong
}

// for json marshal
func (h Hash) MarshalText() ([]byte, error) {
	// TODO:
	return BytesForJson(h[:]).MarshalText()
}

// for json unmarshal
func (h *Hash) UnmarshalJSON(b []byte) error {
	return unmarshalFixedJSON(reflect.TypeOf(Hash{}), b, h[:])
}

// for json unmarshal
func (h *Hash) UnmarshalText(b []byte) error {
	// TODO:
	return unmarshalFixedText("Hash", b, h[:])
}

// for format print
func (h Hash) Format(s fmt.State, c rune) {
	switch c {
	case 'x' | 'X':
		fmt.Fprintf(s, "%#x", h[:])
	default:
		fmt.Fprintf(s, "%"+string(c), h[:])
	}
}

// Sets the hash to the value of b. If b is larger than len(h), 'b' will be cropped (from the left).
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

type Hashs struct {
	Hashs []*Hash
}
