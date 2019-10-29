package test

import (
	"fmt"
	"time"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"encoding/json"
	"bchain.io/utils/database"
	"bchain.io/core/state"
	"bchain.io/common/types"
	"bchain.io/core/transaction"
	"bchain.io/core/actioncontext"
	"math/big"
	"testing"
	"bchain.io/core/interpreter/wasmre"
	"encoding/binary"
)

//go:generate go-bindata -nometadata -pkg test -o bindata.go ../account.wasm
//go:generate gofmt -w -s bindata.go

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
func TestAccout(t *testing.T) {
	code := MustAsset("../account.wasm")
	ctx := makeTestCtx()

	wasmre := wasmre.WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	rt := wasmre.Generate(ctx)

	fmt.Println("Testcase: set val")
	Addr := para_paser.Arg{para_paser.TypeAddress, append([]byte("0x16ff762e278abb68526d8752e7adfa0eec98c0ba"), 0)}

	data := []byte("this is a test k-v account test value")
	val_len := make([]byte, 4)
	binary.LittleEndian.PutUint32(val_len, uint32(len(data)))
	valLen := para_paser.Arg{para_paser.TypeI32, val_len}
	val := para_paser.Arg{para_paser.TypeAddress, data}
	wp := &para_paser.WasmPara{
		FuncName: "set",
		Args:     append([]para_paser.Arg{}, Addr, val, valLen),
	}
	paraBytes, _ := json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d := rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	fmt.Println("Testcase: val of 0x2e68b0583021D78c122f719fc82036529a90571d")
	wp = &para_paser.WasmPara{
		FuncName: "get",
		Args:     append([]para_paser.Arg{}, Addr),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	ret := ctx.ActionResult()[0]
	fmt.Println(string(ret))
}
