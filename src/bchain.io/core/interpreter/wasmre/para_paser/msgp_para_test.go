////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2019 The bchain-go Authors.
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
// @File: msgp_para_test.go
// @Date: 2019/01/07 11:01:07
////////////////////////////////////////////////////////////////////////////////

package para_paser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/tinylib/msgp/msgp"
	"testing"
	"encoding/json"
	"math/big"
)

func generateWasmParabyMsgp() []byte {
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

	var encData bytes.Buffer
	err := msgp.Encode(&encData, wp)
	if err != nil {
		panic(err)
	}

	paraBytes, err := json.Marshal(wp)
	//paraBytes, err := json.MarshalIndent(wp,"","    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(len(paraBytes))
	fmt.Println(len(encData.Bytes()))
	return encData.Bytes()
}

func generateWasmParabyMsgp1() []byte {
	data1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(data1, 100)
	arg1 := Arg{TypeI32, data1}
	arg2 := Arg{TypeAddress, append([]byte("0x2e68b0583021d78c122f719fc82036529a90571d"), 0)}
	arg3 := Arg{TypeAddress, append([]byte("test transfer"), 0)}
	args := []Arg{}
	args = append(args, arg1, arg2, arg3)
	wp := &WasmPara{
		FuncName: "transfer",
		Args:     args,
	}

	var encData bytes.Buffer
	err := msgp.Encode(&encData, wp)
	if err != nil {
		panic(err)
	}

	paraBytes, err := json.Marshal(wp)
	//paraBytes, err := json.MarshalIndent(wp,"","    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(len(paraBytes))

	fmt.Println(len(encData.Bytes()))
	return encData.Bytes()
}

func TestMsgpPara(t *testing.T) {
	paraBytes := generateWasmParabyMsgp1()
	mpp := MsgpParaPaser{}
	wp := mpp.parseWasmPara(paraBytes)
	fmt.Println(wp)

	name, args, mem := mpp.ParseInputPara(paraBytes, 32768, 65536)
	fmt.Println(name, args, mem)
}

func TestMs(t *testing.T) {
	a := big.NewInt(11)
	b := big.NewInt(22)

	a.Add(a,b)
	fmt.Println(a)
}