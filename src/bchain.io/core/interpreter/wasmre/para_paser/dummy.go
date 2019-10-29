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
// @File: dummy.go
// @Date: 2018/12/06 11:05:06
////////////////////////////////////////////////////////////////////////////////

package para_paser

import (
	"encoding/json"
	"fmt"
)



type DummyParaPaser struct {
}

func (dpp *DummyParaPaser) ParseInputPara(para []byte, base, max int) (string, []uint64, []byte) {
	wp := dpp.parseWasmPara(para)
	return formatLinerMemory(base, max, wp)
}

func (dpp *DummyParaPaser) parseWasmPara(data []byte) *WasmPara {
	wp := new(WasmPara)
	err := json.Unmarshal(data, wp)
	if err != nil {
		panic(fmt.Errorf("Parse wasm para fail! Err:%s", err.Error()))
	}
	return wp
}
