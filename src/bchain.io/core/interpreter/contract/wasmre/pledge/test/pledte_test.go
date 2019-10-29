package test

import (
	"bchain.io/core/interpreter/wasmre/para_paser"
	"fmt"
	"encoding/json"
	"time"
	"bchain.io/core/transaction"
	"bchain.io/common/types"
	"bchain.io/core/actioncontext"
	"math/big"
	"testing"
	"bchain.io/core/interpreter/wasmre"
	"bchain.io/utils/database"
	"bchain.io/core/state"
	"bchain.io/core/interpreter/contract/wasmre/deps"
	"encoding/binary"
	"bchain.io/core/interpreter"
)

//go:generate go-bindata -nometadata -pkg test -o bindata.go ../pledge.wasm ../../bchain/bchain.wasm ../../system/system.wasm

func makeTestBlkCtx(blkNumber int64,producer types.Address) *actioncontext.BlockContext{
	db, _ := database.OpenMemDB()
	stateDb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	tmpdb, _ := database.OpenMemDB()


	return actioncontext.NewBlockContext(stateDb, db, tmpdb, big.NewInt(blkNumber), producer)
}


func TestSystemContract(t *testing.T) {
	interpreter.Singleton().Register(wasmre.NewWasmRe)
	interpreter.Singleton().Initialize()
	interpreter.Singleton().Startup()
	code := MustAsset("../../system/system.wasm")
	sender := types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	blkctx := makeTestBlkCtx(5, sender)
	tr := transaction.Action{}
	tr.Contract = types.Address{}
	ctx := actioncontext.NewContext(sender, &tr, blkctx)

	wasmre := wasmre.WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	rt := wasmre.Generate(ctx)

	contractCode := para_paser.Arg{para_paser.TypeAddress, append(deps.MustAsset("bchain.json"), 0)}
	fmt.Println("Testcase: create contract")
	wp := &para_paser.WasmPara{
		FuncName: "createContract",
		Args:     append([]para_paser.Arg{}, contractCode),
	}
	paraBytes, _ := json.Marshal(wp)
	fmt.Println(string(paraBytes))

	d := rt.Exec(code, paraBytes, 1*time.Second)
	fmt.Println("exec time:", d)

	ret := ctx.ActionResult()[0]
	fmt.Println(string(ret))

	code = MustAsset("../pledge.wasm")
	tr = transaction.Action{}
	tr.Contract = types.HexToAddress("0x11b33e3fe72d536a43ddcb84f0db776382069c51")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = wasmre.Generate(ctx)

	data1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(data1, 1000)
	amount := para_paser.Arg{para_paser.TypeI32, data1}

	wp = &para_paser.WasmPara{
		FuncName: "pledge",
		Args:     append([]para_paser.Arg{}, amount),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))

	d = rt.Exec(code, paraBytes, 1*time.Second)
	fmt.Println("exec time:", d)
}

func TestXX(t *testing.T)  {
	data1 := make([]byte, 8)
	binary.LittleEndian.PutUint64(data1, 10000)
	amount := para_paser.Arg{para_paser.TypeI64, data1}
	wp := &para_paser.WasmPara{
		FuncName: "redeem",
		Args:     append([]para_paser.Arg{}, amount),
	}
	paraBytes, _ := json.Marshal(wp)
	fmt.Println(string(paraBytes))
}

func TestX111(t *testing.T)  {
	x := make([]byte, 8)
	W := int64(binary.LittleEndian.Uint64(x))
	fmt.Println(W)
}