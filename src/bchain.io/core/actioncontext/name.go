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
// @File: name.go
// @Date: 2018/09/13 14:59:13
////////////////////////////////////////////////////////////////////////////////

package actioncontext

import (
	"strings"
)

type Name struct {
	value uint64
}


func (n Name) Str() string {
	charMap := []byte(".12345abcdefghijklmnopqrstuvwxyz")
	_ = charMap

	ret := make([]byte, 13)

	//str.(13,'.')

	tmp := n.value
	for i := 0; i <= 12; i++ {
		r1 := uint64(0x0f)
		r2 := uint64(4)
		if i != 0 {
			r1 = 0x1f
			r2 = 5
		}
		c := charMap[tmp & r1]
		ret[12-i] = c
		tmp >>= r2
	}

	return recombine(string(ret[:]))
}

func charToSymbol(c byte) uint64 {
	if c >= 'a' && c <= 'z' {
		return (uint64(c) - 'a') + 6
	}
	if c >= '1' && c <= '5' {
		return (uint64(c) - '1') + 1
	}
	return 0
}

func recombine(str string) string {
	vecStr := strings.Split(str, ".")
	str = ""
	for i, elem := range vecStr {
		if len(elem) != 0 {
			str += elem
			if i != len(vecStr)-1 {
				str += "."
			}
		}
	}
	pos := strings.LastIndex(str, ".")
	if pos == len(str) - 1 {
		str = str[0:pos]
	}
	return str
}

// Converts a base32 string to a uint64.
func StringToName(str string) (ret Name) {
	var value uint64 = 0

	for i := 0; i <= 12; i++ {
		var c uint64 = 0
		if i < len(str) && i <= 12 {
			c = charToSymbol(str[i])
		}

		if i < 12 {
			c &= 0x1f
			c <<= 64-5*(uint64(i)+1)
		} else {
			c &= 0x0f
		}

		value |= c
	}

	ret.value = value
	return
}

func IsNormativeName(str string) bool {
	name := StringToName(str)
	return strings.Compare(name.Str(), str) == 0
}
