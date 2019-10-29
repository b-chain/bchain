package test

import (
	"fmt"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"encoding/json"
	"time"
	"bchain.io/common/types"
	"bchain.io/core/actioncontext"
	"bchain.io/utils/database"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
	"math/big"
	"testing"
	"bchain.io/core/interpreter/wasmre"
	"bchain.io/core/interpreter/contract/wasmre/deps"
)

//go:generate go-bindata -nometadata -pkg test -o bindata.go ../system.wasm
//go:generate gofmt -w -s bindata.go

func makeTestBlkCtx(blkNumber int64,producer types.Address) *actioncontext.BlockContext{
	db, _ := database.OpenMemDB()
	stateDb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	tmpdb, _ := database.OpenMemDB()

	tr := transaction.Action{}
	tr.Contract = types.Address{}

	return actioncontext.NewBlockContext(stateDb, db, tmpdb, big.NewInt(blkNumber), producer)
}

func TestSystemContract(t *testing.T) {
	code := MustAsset("../system.wasm")
	sender := types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	blkctx := makeTestBlkCtx(250000*200-1, sender)
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


}