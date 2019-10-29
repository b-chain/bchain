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
// @File: json.go
// @Date: 2018/05/08 17:26:08
////////////////////////////////////////////////////////////////////////////////

package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

var (
	bytesT  = reflect.TypeOf(BytesForJson(nil))
	uintT   = reflect.TypeOf(UintForJson(0))
	uint64T = reflect.TypeOf(Uint64ForJson(0))
)

// Bytes marshals/unmarshals as a JSON string with 0x prefix.
// The empty slice marshals as "0x".
type BytesForJson []byte

// MarshalText implements encoding.TextMarshaler
func (b BytesForJson) MarshalText() ([]byte, error) {
	result := make([]byte, len(b)*2+2)
	copy(result, `0x`)
	hex.Encode(result[2:], b)
	return result, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BytesForJson) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return errNonString(bytesT)
	}
	return wrapTypeError(b.UnmarshalText(input[1:len(input)-1]), bytesT)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *BytesForJson) UnmarshalText(input []byte) error {
	raw, err := checkText(input, true)
	if err != nil {
		return err
	}
	dec := make([]byte, len(raw)/2)
	if _, err = hex.Decode(dec, raw); err != nil {
		err = mapError(err)
	} else {
		*b = dec
	}
	return err
}

// String returns the hex encoding of b.
func (b BytesForJson) String() string {
	return EncodeHex(b)
}

// UnmarshalFixedJSON decodes the input as a string with 0x prefix. The length of out
// determines the required input length. This function is commonly used to implement the
// UnmarshalJSON method for fixed-size types.
func unmarshalFixedJSON(typ reflect.Type, input, out []byte) error {
	if !isString(input) {
		return errNonString(typ)
	}
	return wrapTypeError(unmarshalFixedText(typ.String(), input[1:len(input)-1], out), typ)
}

// unmarshalFixedText decodes the input as a string with 0x prefix. The length of out
// determines the required input length. This function is commonly used to implement the
// UnmarshalText method for fixed-size types.
func unmarshalFixedText(typname string, input, out []byte) error {
	raw, err := checkText(input, true)
	if err != nil {
		return err
	}
	if len(raw)/2 != len(out) {
		return fmt.Errorf("hex string has length %d, want %d for %s", len(raw), len(out)*2, typname)
	}
	// Pre-verify syntax before modifying out.
	for _, b := range raw {
		if decodeNibble(b) == badNibble {
			return ErrSyntax
		}
	}
	hex.Decode(out, raw)
	return nil
}

// unmarshalFixedUnprefixedText decodes the input as a string with optional 0x prefix. The
// length of out determines the required input length. This function is commonly used to
// implement the UnmarshalText method for fixed-size types.
func unmarshalFixedUnprefixedText(typname string, input, out []byte) error {
	raw, err := checkText(input, false)
	if err != nil {
		return err
	}
	if len(raw)/2 != len(out) {
		return fmt.Errorf("hex string has length %d, want %d for %s", len(raw), len(out)*2, typname)
	}
	// Pre-verify syntax before modifying out.
	for _, b := range raw {
		if decodeNibble(b) == badNibble {
			return ErrSyntax
		}
	}
	hex.Decode(out, raw)
	return nil
}

// Uint64 marshals/unmarshals as a JSON string with 0x prefix.
// The zero value marshals as "0x0".
type Uint64ForJson uint64

// MarshalText implements encoding.TextMarshaler.
func (b Uint64ForJson) MarshalText() ([]byte, error) {
	buf := make([]byte, 2, 10)
	copy(buf, `0x`)
	buf = strconv.AppendUint(buf, uint64(b), 16)
	return buf, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Uint64ForJson) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return errNonString(uint64T)
	}
	return wrapTypeError(b.UnmarshalText(input[1:len(input)-1]), uint64T)
}

// UnmarshalText implements encoding.TextUnmarshaler
func (b *Uint64ForJson) UnmarshalText(input []byte) error {
	raw, err := checkNumberText(input)
	if err != nil {
		return err
	}
	if len(raw) > 16 {
		return ErrUint64Range
	}
	var dec uint64
	for _, byte := range raw {
		nib := decodeNibble(byte)
		if nib == badNibble {
			return ErrSyntax
		}
		dec *= 16
		dec += nib
	}
	*b = Uint64ForJson(dec)
	return nil
}

// String returns the hex encoding of b.
func (b Uint64ForJson) String() string {
	return EncodeHexUint64(uint64(b))
}

// Uint marshals/unmarshals as a JSON string with 0x prefix.
// The zero value marshals as "0x0".
type UintForJson uint

// MarshalText implements encoding.TextMarshaler.
func (b UintForJson) MarshalText() ([]byte, error) {
	return Uint64ForJson(b).MarshalText()
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *UintForJson) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return errNonString(uintT)
	}
	return wrapTypeError(b.UnmarshalText(input[1:len(input)-1]), uintT)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *UintForJson) UnmarshalText(input []byte) error {
	var u64 Uint64ForJson
	err := u64.UnmarshalText(input)
	if u64 > Uint64ForJson(^uint(0)) || err == ErrUint64Range {
		return ErrUintRange
	} else if err != nil {
		return err
	}
	*b = UintForJson(u64)
	return nil
}

// String returns the hex encoding of b.
func (b UintForJson) String() string {
	return EncodeHexUint64(uint64(b))
}

func isString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

func bytesHave0xPrefix(input []byte) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func checkText(input []byte, wantPrefix bool) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil // empty strings are allowed
	}
	if bytesHave0xPrefix(input) {
		input = input[2:]
	} else if wantPrefix {
		return nil, ErrMissingPrefix
	}
	if len(input)%2 != 0 {
		return nil, ErrOddLength
	}
	return input, nil
}

func checkNumberText(input []byte) (raw []byte, err error) {
	if len(input) == 0 {
		return nil, nil // empty strings are allowed
	}
	if !bytesHave0xPrefix(input) {
		return nil, ErrMissingPrefix
	}
	input = input[2:]
	if len(input) == 0 {
		return nil, ErrEmptyNumber
	}
	if len(input) > 1 && input[0] == '0' {
		return nil, ErrLeadingZero
	}
	return input, nil
}

func wrapTypeError(err error, typ reflect.Type) error {
	if _, ok := err.(*decError); ok {
		return &json.UnmarshalTypeError{Value: err.Error(), Type: typ}
	}
	return err
}

func errNonString(typ reflect.Type) error {
	return &json.UnmarshalTypeError{Value: "non-string", Type: typ}
}
