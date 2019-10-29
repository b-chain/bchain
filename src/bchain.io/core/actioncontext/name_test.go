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
// @File: name_test.go
// @Date: 2018/09/13 15:50:13
////////////////////////////////////////////////////////////////////////////////

package actioncontext

import "testing"

func Test_1(t *testing.T) {
	vecStr := []string{
		"abcdeABCDfg.7",
		"xyz.1234567890",
		"12.34.51.231.45",
		"12.abcdxyz",
	}

	for _, str := range vecStr {
		name := StringToName(str)
		t.Logf("name = %v, %v", name.Str(), IsNormativeName(str))
	}
}
