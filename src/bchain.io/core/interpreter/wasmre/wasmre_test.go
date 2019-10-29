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
// @File: wasmre_test.go
// @Date: 2018/12/06 15:59:06
////////////////////////////////////////////////////////////////////////////////

package wasmre

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"bchain.io/common/types"
	"bchain.io/core/actioncontext"
	"bchain.io/core/interpreter"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"bchain.io/core/interpreter/wasmre/test_deps"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
	"bchain.io/utils/database"
	"testing"
	"time"
	"bytes"
	"strings"
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
	return ctx
}
func genTestCtx(number int) *actioncontext.Context {
	db, _ := database.OpenMemDB()
	stateDb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	tmpdb, _ := database.OpenMemDB()

	tr := transaction.Action{}
	con := types.Address{}
	con[1] = 1
	tr.Contract = con
	blkctx := actioncontext.NewBlockContext(stateDb, db, tmpdb, big.NewInt(int64(number)), types.Address{})
	sender := types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	ctx := actioncontext.NewContext(sender, &tr, blkctx)
	return ctx
}

func test1Para() []byte {
	args := []para_paser.Arg{}
	data1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(data1, 100)
	arg := para_paser.Arg{para_paser.TypeI32, data1}
	args = append(args, arg)
	wp := &para_paser.WasmPara{
		FuncName: "test1",
		Args:     args,
	}
	paraBytes, err := json.Marshal(wp)
	if err != nil {
		panic(err)
	}
	return paraBytes
}

func TestWasmExec(t *testing.T) {
	code := test_deps.MustAsset("test.wasm")
	ctx := makeTestCtx()

	wasmre := WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	rt := wasmre.Generate(ctx)

	d := rt.Exec(code, test1Para(), 1*time.Second)
	fmt.Println("exec time:", d)
}

func timoutPanic(t *testing.T, code []byte, para []byte, rt interpreter.Interpreter) {
	defer func() {
		if caught := recover(); caught != nil {
			if caught.(string) != "timeout interrupter" {
				t.Errorf("wrong assert msg, got %v,", caught.(string))
			}
		} else {
			t.Errorf("can not caught panic !")
		}
	}()
	d := rt.Exec(code, para, 1*time.Second)
	fmt.Println("first exec time:", d)
}

func TestTimeout(t *testing.T) {
	code := test_deps.MustAsset("test.wasm")
	wasmre := WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	ctx1 := genTestCtx(1000)
	rt := wasmre.Generate(ctx1)
	args := []para_paser.Arg{}
	wp := &para_paser.WasmPara{
		FuncName: "loopForever",
		Args:     args,
	}
	paraBytes, _ := json.Marshal(wp)
	timoutPanic(t, code, paraBytes, rt)

	ctx2 := genTestCtx(2000)
	rt = wasmre.Generate(ctx2)
	wp = &para_paser.WasmPara{
		FuncName: "testBlockNumber",
		Args:     args,
	}
	paraBytes, _ = json.Marshal(wp)
	d := rt.Exec(code, paraBytes, 1*time.Second)
	fmt.Println("exec time:", d.Nanoseconds(), "ns")
}

func TestApi_assert(t *testing.T) {
	code := test_deps.MustAsset("test.wasm")
	ctx := makeTestCtx()

	wasmre := WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	rt := wasmre.Generate(ctx)

	args := []para_paser.Arg{}
	wp := &para_paser.WasmPara{
		FuncName: "testAssert1",
		Args:     args,
	}
	paraBytes, _ := json.Marshal(wp)

	d := rt.Exec(code, paraBytes, 1*time.Second)
	fmt.Println("exec time:", d)

	wp = &para_paser.WasmPara{
		FuncName: "testAssert2",
		Args:     args,
	}
	paraBytes, _ = json.Marshal(wp)
	defer func() {
		if caught := recover(); caught != nil {
			if caught.(string) != "test assert" {
				t.Errorf("wrong assert msg, got %v,", caught.(string))
			}
		} else {
			t.Errorf("can not caught panic !")
		}
	}()
	d = rt.Exec(code, paraBytes, 1*time.Second)
	fmt.Println("exec time:", d)
}

func TestApi_sha(t *testing.T) {
	code := test_deps.MustAsset("test.wasm")
	ctx := makeTestCtx()

	wasmre := WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	rt := wasmre.Generate(ctx)

	args := []para_paser.Arg{}
	wp := &para_paser.WasmPara{
		FuncName: "testCryotoApi",
		Args:     args,
	}
	paraBytes, _ := json.Marshal(wp)

	d := rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
}

// for function call, the stack in memory grow direction is negative(form big to little)
// when function exit, the stack address(mem[4:8]) will be reset to caller function value
// clang -o0 is needed.
func TestStack(t *testing.T) {
	code := test_deps.MustAsset("test.wasm")
	ctx := makeTestCtx()

	wasmre := WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	rt := wasmre.Generate(ctx)

	args := []para_paser.Arg{}
	wp := &para_paser.WasmPara{
		//FuncName: "_Z13testCryotoApiv",
		FuncName: "testStack",
		Args:     args,
	}
	paraBytes, _ := json.Marshal(wp)

	d := rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
}

func TestApi_BlockNumber(t *testing.T) {
	code := test_deps.MustAsset("test.wasm")
	ctx := makeTestCtx()

	wasmre := WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	rt := wasmre.Generate(ctx)

	args := []para_paser.Arg{}
	wp := &para_paser.WasmPara{
		FuncName: "testBlockNumber",
		Args:     args,
	}
	paraBytes, _ := json.Marshal(wp)

	d := rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
}

func TestApi_mem(t *testing.T) {
	code := test_deps.MustAsset("test.wasm")
	ctx := makeTestCtx()

	wasmre := WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	rt := wasmre.Generate(ctx)

	args := []para_paser.Arg{}
	wp := &para_paser.WasmPara{
		FuncName: "memTest",
		Args:     args,
	}
	paraBytes, _ := json.Marshal(wp)

	d := rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
}

func TestApi_Token(t *testing.T) {
	code := test_deps.MustAsset("test.wasm")
	ctx := makeTestCtx()

	wasmre := WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	rt := wasmre.Generate(ctx)

	fmt.Println("Testcase: token create")
	wp := &para_paser.WasmPara{
		FuncName: "create",
		Args:     []para_paser.Arg{},
	}
	paraBytes, _ := json.Marshal(wp)

	d := rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	fmt.Println("Testcase: token creater balence")
	creator := para_paser.Arg{para_paser.TypeAddress, append([]byte("0x2e68b0583021D78c122f719fc82036529a90571d"), 0)}
	wp = &para_paser.WasmPara{
		FuncName: "balenceOf",
		Args:     append([]para_paser.Arg{}, creator),
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
	wp = &para_paser.WasmPara{
		FuncName: "transer",
		Args:     append([]para_paser.Arg{}, toAddr, amount),
	}
	paraBytes, _ = json.Marshal(wp)
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	fmt.Println("Testcase: balence of 0x2e68b0583021D78c122f719fc82036529a90571d")
	wp = &para_paser.WasmPara{
		FuncName: "balenceOf",
		Args:     append([]para_paser.Arg{}, creator),
	}
	paraBytes, _ = json.Marshal(wp)
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	fmt.Println("Testcase: balence of 0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	wp = &para_paser.WasmPara{
		FuncName: "balenceOf",
		Args:     append([]para_paser.Arg{}, toAddr),
	}
	paraBytes, _ = json.Marshal(wp)
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
}

func TestCache(t *testing.T) {
	code := test_deps.MustAsset("test.wasm")
	wasmre := WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	ctx1 := genTestCtx(1000)
	rt := wasmre.Generate(ctx1)
	args := []para_paser.Arg{}
	wp := &para_paser.WasmPara{
		FuncName: "testBlockNumber",
		Args:     args,
	}
	paraBytes, _ := json.Marshal(wp)
	d := rt.Exec(code, paraBytes, 1*time.Second)
	fmt.Println("first exec time:", d)

	ctx2 := genTestCtx(2000)
	rt = wasmre.Generate(ctx2)
	wp = &para_paser.WasmPara{
		FuncName: "testBlockNumber",
		Args:     args,
	}
	paraBytes, _ = json.Marshal(wp)
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("optimazation exec time:", d.Nanoseconds(), "ns")

	d = rt.Exec(code, test1Para(), 1*time.Second)
	fmt.Println("optimazation exec time:", d.Nanoseconds(), "ns")
}

// testBlockNumber is very simple, very fast
func TestSpeed(t *testing.T) {
	code := test_deps.MustAsset("test.wasm")
	wasmre := WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	args := []para_paser.Arg{}
	wp := &para_paser.WasmPara{
		FuncName: "testBlockNumber",
		Args:     args,
	}
	paraBytes, _ := json.Marshal(wp)

	start := time.Now()
	count := 2000
	for i := 0; i < count; i++ {
		ctx1 := genTestCtx(1000 + i)
		rt := wasmre.Generate(ctx1)
		rt.Exec(code, paraBytes, 1*time.Second)
	}
	fmt.Println("exec", count, "times, time:", time.Since(start))
}

// token contract is more complex, need more exec time
func TestSpeedToken(t *testing.T) {
	code := test_deps.MustAsset("test.wasm")
	wasmre := WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	wp := &para_paser.WasmPara{
		FuncName: "create",
		Args:     []para_paser.Arg{},
	}
	paraCreate, _ := json.Marshal(wp)

	data1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(data1, 1000)
	amount := para_paser.Arg{para_paser.TypeI32, data1}
	toAddr := para_paser.Arg{para_paser.TypeAddress, append([]byte("0x16ff762e278abb68526d8752e7adfa0eec98c0ba"), 0)}
	wp = &para_paser.WasmPara{
		FuncName: "transer",
		Args:     append([]para_paser.Arg{}, toAddr, amount),
	}
	paraTransfer, _ := json.Marshal(wp)

	wp = &para_paser.WasmPara{
		FuncName: "balenceOf",
		Args:     append([]para_paser.Arg{}, toAddr),
	}
	paraQuery, _ := json.Marshal(wp)

	start := time.Now()
	count := 1000
	for i := 0; i < count; i++ {
		ctx1 := genTestCtx(1000 + i)
		rt := wasmre.Generate(ctx1)
		rt.Exec(code, paraCreate, 1*time.Second)
		rt.Exec(code, paraTransfer, 1*time.Second)
		rt.Exec(code, paraQuery, 1*time.Second)
	}
	fmt.Println("exec", count, "times, time:", time.Since(start))
}

func TestXx(t *testing.T)  {
	xxx := make([]byte, 10)
	//xxx[0] = 1
	xxx[0] = 1
	xxx[1] =2
 	index := bytes.IndexByte(xxx, 0)
	fmt.Println(index)
 	fmt.Println(xxx[0:2])

	aa := "1100xxGGFs"
	cc := append([]byte(aa),5)
	bb := strings.ToLower(string(cc))
	fmt.Println(aa,bb)

	addr := types.HexToAddress("f466859ead1932d743d622cb74fc058882e8648a")
	fmt.Println(addr.Hex(),  addr.HexLower())

}
