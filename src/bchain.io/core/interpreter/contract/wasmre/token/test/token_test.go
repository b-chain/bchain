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
// @File: token_test.go
// @Date: 2019/01/08 16:30:08
////////////////////////////////////////////////////////////////////////////////

package test

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"bchain.io/common/types"
	"bchain.io/core/actioncontext"
	"bchain.io/core/interpreter/wasmre"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
	"bchain.io/utils/database"
	"testing"
	"time"
)

func makeTestCtx() *actioncontext.Context {
	db, _ := database.OpenMemDB()
	stateDb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	tmpdb, _ := database.OpenMemDB()

	tr := transaction.Action{}
	tr.Contract = types.Address{}
	blkctx := actioncontext.NewBlockContext(stateDb, db, tmpdb, big.NewInt(998), types.Address{})
	sender := types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	ctx := actioncontext.NewContext(sender, &tr, blkctx)
	ctx.SetCreatorForTest(types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d"))
	return ctx
}
func TestToken(t *testing.T) {
	code := MustAsset("../token.wasm")
	ctx := makeTestCtx()

	wasmre := wasmre.WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	rt := wasmre.Generate(ctx)

	fmt.Println("Testcase: token create")
	supply_b := make([]byte, 8)
	binary.LittleEndian.PutUint64(supply_b, 1000000000000)
	supply := para_paser.Arg{para_paser.TypeI64, supply_b}
	decimals_b := make([]byte, 4)
	binary.LittleEndian.PutUint32(decimals_b, 4)
	decimals := para_paser.Arg{para_paser.TypeI32, decimals_b}
	sybmol := para_paser.Arg{para_paser.TypeAddress, append([]byte("CNB"), 0)}
	name := para_paser.Arg{para_paser.TypeAddress, append([]byte("china NB"), 0)}
	wp := &para_paser.WasmPara{
		FuncName: "create",
		Args:     append([]para_paser.Arg{}, sybmol, name, decimals, supply),
	}
	paraBytes, _ := json.Marshal(wp)
	fmt.Println(string(paraBytes))

	d := rt.Exec(code, paraBytes, 1*time.Second)
	fmt.Println("exec time:", d)

	fmt.Println("Testcase: token creater balence")
	creator := para_paser.Arg{para_paser.TypeAddress, append([]byte("0x2e68b0583021d78c122f719fc82036529a90571d"), 0)}
	wp = &para_paser.WasmPara{
		FuncName: "balanceOf",
		Args:     append([]para_paser.Arg{}, creator, sybmol),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	fmt.Println("Testcase: balence tansfer")
	data1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(data1, 1000)
	amount := para_paser.Arg{para_paser.TypeI32, data1}
	toAddr := para_paser.Arg{para_paser.TypeAddress, append([]byte("0x16ff762e278abb68526d8752e7adfa0eec98c0ba"), 0)}
	memo := para_paser.Arg{para_paser.TypeAddress, append([]byte("this is a test transfer"), 0)}
	wp = &para_paser.WasmPara{
		FuncName: "transfer",
		Args:     append([]para_paser.Arg{}, toAddr, amount, sybmol, memo),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	fmt.Println("Testcase: balance of 0x2e68b0583021D78c122f719fc82036529a90571d")
	wp = &para_paser.WasmPara{
		FuncName: "balanceOf",
		Args:     append([]para_paser.Arg{}, creator, sybmol),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	fmt.Println("Testcase: balance of 0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	wp = &para_paser.WasmPara{
		FuncName: "balanceOf",
		Args:     append([]para_paser.Arg{}, toAddr, sybmol),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	fmt.Println("Testcase: getSupply")
	wp = &para_paser.WasmPara{
		FuncName: "getSupply",
		Args:     append([]para_paser.Arg{}, sybmol),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	fmt.Println("Testcase: getDecimals")
	wp = &para_paser.WasmPara{
		FuncName: "getDecimals",
		Args:     append([]para_paser.Arg{}, sybmol),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
}

func TestSpeedToken(t *testing.T) {
	code := MustAsset("../token.wasm")
	wasmre := wasmre.WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	supply_b := make([]byte, 8)
	binary.LittleEndian.PutUint64(supply_b, 1000000000000)
	supply := para_paser.Arg{para_paser.TypeI64, supply_b}
	decimals_b := make([]byte, 4)
	binary.LittleEndian.PutUint32(decimals_b, 4)
	decimals := para_paser.Arg{para_paser.TypeI32, decimals_b}
	sybmol := para_paser.Arg{para_paser.TypeAddress, append([]byte("￥"), 0)}
	name := para_paser.Arg{para_paser.TypeAddress, append([]byte("CNY"), 0)}
	wp := &para_paser.WasmPara{
		FuncName: "create",
		Args:     append([]para_paser.Arg{}, sybmol, name, decimals, supply),
	}
	paraCreate, _ := json.Marshal(wp)

	data1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(data1, 1000)
	amount := para_paser.Arg{para_paser.TypeI32, data1}
	toAddr := para_paser.Arg{para_paser.TypeAddress, append([]byte("0x16ff762e278abb68526d8752e7adfa0eec98c0ba"), 0)}
	memo := para_paser.Arg{para_paser.TypeAddress, append([]byte("this is a test transfer"), 0)}
	wp = &para_paser.WasmPara{
		FuncName: "transfer",
		Args:     append([]para_paser.Arg{}, toAddr, amount, sybmol, memo),
	}
	paraTransfer, _ := json.Marshal(wp)

	wp = &para_paser.WasmPara{
		FuncName: "balanceOf",
		Args:     append([]para_paser.Arg{}, toAddr, sybmol),
	}
	paraQuery, _ := json.Marshal(wp)

	start := time.Now()
	count := 1000
	for i := 0; i < count; i++ {
		ctx1 := makeTestCtx()
		rt := wasmre.Generate(ctx1)
		rt.Exec(code, paraCreate, 1*time.Second)
		rt.Exec(code, paraTransfer, 1*time.Second)
		rt.Exec(code, paraQuery, 1*time.Second)
	}
	fmt.Println("exec", count, "times, time:", time.Since(start))
}

func TestTokenxx(t *testing.T) {
	fmt.Println("Testcase: balance tansfer")
	data1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(data1, 1000)
	amount := para_paser.Arg{para_paser.TypeI32, data1}
	toAddr := para_paser.Arg{para_paser.TypeAddress, append([]byte("0x16ff762e278abb68526d8752e7adfa0eec98c0ba"), 0)}
	memo := para_paser.Arg{para_paser.TypeAddress, append([]byte("this is a test transfer"), 0)}
	wp := &para_paser.WasmPara{
		FuncName: "transfer",
		Args:     append([]para_paser.Arg{}, toAddr, amount, memo),
	}
	paraBytes, _ := json.Marshal(wp)
	fmt.Println(string(paraBytes))

	addr := para_paser.Arg{para_paser.TypeAddress, append([]byte("0x07667bb16451840c55201487def17fecd94894f5"), 0)}
	wp = &para_paser.WasmPara{
		FuncName: "balanceOf",
		Args:     append([]para_paser.Arg{}, addr),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
}
