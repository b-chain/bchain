package test

import (
	"bchain.io/core/interpreter/wasmre/para_paser"
	"encoding/json"
	"fmt"
	"time"
	"testing"
	"bchain.io/core/interpreter/wasmre"
	"bchain.io/core/actioncontext"
	"bchain.io/utils/database"
	"bchain.io/core/state"
	"bchain.io/common/types"
	"bchain.io/core/transaction"
	"math/big"
	"encoding/binary"
)

func makeTestBlkCtx(blkNumber int64,producer types.Address) *actioncontext.BlockContext{
	db, _ := database.OpenMemDB()
	stateDb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	tmpdb, _ := database.OpenMemDB()

	tr := transaction.Action{}
	tr.Contract = types.Address{}

	return actioncontext.NewBlockContext(stateDb, db, tmpdb, big.NewInt(blkNumber), producer)
}

func TestBcToken(t *testing.T) {
	code := MustAsset("../bchain.wasm")
	sender := types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	blkctx := makeTestBlkCtx(250000*200-1, sender)
	tr := transaction.Action{}
	tr.Contract = types.Address{}
	ctx := actioncontext.NewContext(sender, &tr, blkctx)

	wasmre := wasmre.WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	rt := wasmre.Generate(ctx)

	fmt.Println("Testcase: bchain reword")
	wp := &para_paser.WasmPara{
		FuncName: "reword",
		Args:     append([]para_paser.Arg{}),
	}
	paraBytes, _ := json.Marshal(wp)
	fmt.Println(string(paraBytes))

	d := rt.Exec(code, paraBytes, 1*time.Second)
	fmt.Println("exec time:", d)

	addr := para_paser.Arg{para_paser.TypeAddress, append([]byte("0x2e68b0583021d78c122f719fc82036529a90571d"), 0)}
	wp = &para_paser.WasmPara{
		FuncName: "balenceOf",
		Args:     append([]para_paser.Arg{}, addr),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	ret := ctx.ActionResult()[0]
	v := binary.LittleEndian.Uint64(ret)
	fmt.Printf("balenceOf 0x2e68b0583021d78c122f719fc82036529a90571d %v\n",  v)

	fmt.Println("Testcase: transfer")
	data1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(data1, 1000)
	amount := para_paser.Arg{para_paser.TypeI32, data1}
	toAddr := para_paser.Arg{para_paser.TypeAddress, append([]byte("0x16ff762e278abb68526d8752e7adfa0eec98c0ba"), 0)}
	memo := para_paser.Arg{para_paser.TypeAddress, append([]byte("this is a test transfer"), 0)}
	wp = &para_paser.WasmPara{
		FuncName: "transfer",
		Args:     append([]para_paser.Arg{}, toAddr, amount, memo),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	wp = &para_paser.WasmPara{
		FuncName: "balenceOf",
		Args:     append([]para_paser.Arg{}, addr),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	ret = ctx.ActionResult()[1]
	v = binary.LittleEndian.Uint64(ret)
	fmt.Printf("balenceOf 0x2e68b0583021d78c122f719fc82036529a90571d %v\n",  v)

	wp = &para_paser.WasmPara{
		FuncName: "balenceOf",
		Args:     append([]para_paser.Arg{}, toAddr),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	ret = ctx.ActionResult()[2]
	v = binary.LittleEndian.Uint64(ret)
	fmt.Printf("balenceOf 0x16ff762e278abb68526d8752e7adfa0eec98c0ba %v\n",  v)

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = wasmre.Generate(ctx)

	data1 = make([]byte, 4)
	binary.LittleEndian.PutUint32(data1, 100)
	amount = para_paser.Arg{para_paser.TypeI32, data1}

	wp = &para_paser.WasmPara{
		FuncName: "transferFee",
		Args:     append([]para_paser.Arg{}, amount),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	wp = &para_paser.WasmPara{
		FuncName: "balenceOf",
		Args:     append([]para_paser.Arg{}, addr),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	ret = ctx.ActionResult()[0]
	v = binary.LittleEndian.Uint64(ret)
	fmt.Printf("balenceOf 0x2e68b0583021d78c122f719fc82036529a90571d %v\n",  v)

	wp = &para_paser.WasmPara{
		FuncName: "balenceOf",
		Args:     append([]para_paser.Arg{}, toAddr),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	ret = ctx.ActionResult()[1]
	v = binary.LittleEndian.Uint64(ret)
	fmt.Printf("balenceOf 0x16ff762e278abb68526d8752e7adfa0eec98c0ba %v\n",  v)
}
