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
// @File: types.go
// @Date: 2019/01/07 10:58:07
////////////////////////////////////////////////////////////////////////////////

package para_paser

import "encoding/binary"

const (
	TypeI32     string = "int32"
	TypeI64     string = "int64"
	TypeF32     string = "float32"
	TypeF64     string = "float64"
	TypeAddress string = "address"
)
//go:generate msgp
type WasmPara struct {
	FuncName string `json:"func_name" msg:"funcName"`
	Args     []Arg  `json:"args"      msg:"args"`
}

type Arg struct {
	Type string `json:"type" msg:"type"`
	Data []byte `json:"val"  msg:"val"`
}

func safeAppend(mem, data []byte, max int) []byte {
	mem = append(mem, data...)
	if len(mem) > max {
		panic("memory exceed!")
	}
	return mem
}

func formatLinerMemory(base, max int, wp *WasmPara) (string, []uint64, []byte) {
	mem := []byte{}
	args := make([]uint64, 0)
	realMax := max - base
	cur := base
	for _, arg := range wp.Args {
		switch arg.Type {
		case TypeI32:
			args = append(args, uint64(binary.LittleEndian.Uint32(arg.Data)))
		case TypeI64:
			args = append(args, binary.LittleEndian.Uint64(arg.Data))
		case TypeF32:
			args = append(args, uint64(binary.LittleEndian.Uint32(arg.Data)))
		case TypeF64:
			args = append(args, binary.LittleEndian.Uint64(arg.Data))
		case TypeAddress:
			mem = safeAppend(mem, arg.Data, realMax)
			args = append(args, uint64(cur))
			cur = base + len(mem)
		default:
			panic("invalid wasm para arg type")
		}
	}
	return wp.FuncName, args, mem
}