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
// @File: bigint.go
// @Date: 2018/05/08 17:27:08
////////////////////////////////////////////////////////////////////////////////

package types

import (
	"math/big"
	"reflect"
)

//go:generate msgp
//msgp:shim big.Int as:interface{} using:bigToBytes/bigFromBytes

var (
	bigIntType int8

	bigT = reflect.TypeOf((*BigInt)(nil))
)

func bigToBytes(v big.Int) interface{} {
	neg := v.Sign()
	b := make([]byte, 1+len(v.Bytes()))
	b[0] = byte(neg)
	copy(b[1:], v.Bytes())
	return b
}

func bigFromBytes(b interface{}) big.Int {
	if b == nil {
		return big.Int{}
	}

	buf, ok := b.([]byte)
	if !ok {
		return big.Int{}
	}
	neg := buf[0]
	//v := new(big.Int)
	var v big.Int
	v.SetBytes(buf[1:])

	if neg == 255 {
		v.Neg(&v)
	}
	return v
}

type BigInt struct {
	IntVal big.Int `msg:"bigint"`
}

func (bigInt BigInt) Get() big.Int {
	return bigInt.IntVal
}

func (bigInt *BigInt) Put(in big.Int) *BigInt {
	bigInt.IntVal = in
	return bigInt
}

func NewBigInt(in big.Int) *BigInt {
	bigInt := new(BigInt)
	bigInt.IntVal = in
	return bigInt
}

// Here, we'll pick an arbitrary number between
// 0 and 127 that isn't already in use
func (bigInt *BigInt) ExtensionType() int8 {
	return bigIntType
}

// We'll always use 1 + len(big.int.x) bytes to encode the data
func (bigInt *BigInt) Len() int {
	//return 1 + len(bigInt.intVal.Bytes())
	return 10
}

func (bigInt *BigInt) MarshalBinaryTo(b []byte) error {
	neg := bigInt.IntVal.Sign()
	b[0] = byte(neg)
	copy(b[1:], bigInt.IntVal.Bytes())
	return nil
}

func (bigInt *BigInt) UnmarshalBinary(b []byte) error {

	neg := b[0]
	bigInt.IntVal.SetBytes(b[1:])

	if neg == 255 {
		bigInt.IntVal.Neg(&bigInt.IntVal)
	}
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (b BigInt) MarshalText() ([]byte, error) {
	return []byte(EncodeHexBig(b.ToInt())), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BigInt) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return errNonString(bigT)
	}
	return wrapTypeError(b.UnmarshalText(input[1:len(input)-1]), bigT)
}

// UnmarshalText implements encoding.TextUnmarshaler
func (b *BigInt) UnmarshalText(input []byte) error {
	raw, err := checkNumberText(input)
	if err != nil {
		return err
	}
	if len(raw) > 64 {
		return ErrBig256Range
	}
	words := make([]big.Word, len(raw)/bigWordNibbles+1)
	end := len(raw)
	for i := range words {
		start := end - bigWordNibbles
		if start < 0 {
			start = 0
		}
		for ri := start; ri < end; ri++ {
			nib := decodeNibble(raw[ri])
			if nib == badNibble {
				return ErrSyntax
			}
			words[i] *= 16
			words[i] += big.Word(nib)
		}
		end = start
	}
	var dec big.Int
	dec.SetBits(words)
	b.IntVal = dec
	return nil
}

// ToInt converts b to a big.Int.
func (b *BigInt) ToInt() *big.Int {
	return &b.IntVal
}

// String returns the hex encoding of b.
func (b *BigInt) String() string {
	return EncodeHexBig(b.ToInt())
}
