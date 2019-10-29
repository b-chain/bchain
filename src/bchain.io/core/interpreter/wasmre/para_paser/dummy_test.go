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
// @File: dummy_test.go
// @Date: 2018/12/06 11:28:06
////////////////////////////////////////////////////////////////////////////////

package para_paser

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"testing"
)

func generateWasmPara() []byte {
	data1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(data1, 100)
	arg1 := Arg{TypeI32, data1}
	arg2 := Arg{TypeAddress, append([]byte("hello world!\n"), 0)}
	arg3 := Arg{TypeAddress, append([]byte("bchain!\n"), 0)}
	args := []Arg{}
	args = append(args, arg1, arg2, arg3, arg1, arg2, arg1, arg3)
	wp := &WasmPara{
		FuncName: "test",
		Args:     args,
	}

	paraBytes, err := json.Marshal(wp)
	//paraBytes, err := json.MarshalIndent(wp,"","    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(len(paraBytes))
	fmt.Println(string(paraBytes))
	return paraBytes
}

func TestDummyPara(t *testing.T) {
	paraBytes := generateWasmPara()
	dpp := DummyParaPaser{}
	wp := dpp.parseWasmPara(paraBytes)
	fmt.Println(wp)

	name, args, mem := dpp.ParseInputPara(paraBytes, 32768, 65536)
	fmt.Println(name, args, mem)
}
