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
// @File: msgp_test.go
// @Date: 2018/05/17 13:14:17
////////////////////////////////////////////////////////////////////////////////

package common

import (
	"math/big"
	"bchain.io/common/types"
	"reflect"
	"testing"
)

type testst struct {
	a string
}

// type vectest []testst

func TestMsgpHash(t *testing.T) {
	want := "msgp: type \"common.testst\" not supported"
	_, err := MsgpHash([]interface{}{
		uint64(1),
		uint(0),
		types.NewBigInt(*big.NewInt(123)),
		[]testst{{"abc"}, {"111"}},
	})

	if err != nil {
		if !reflect.DeepEqual(want, err.Error()) {
			t.Errorf("have error: %v, want: %v", err, want)
		}
	}
}
