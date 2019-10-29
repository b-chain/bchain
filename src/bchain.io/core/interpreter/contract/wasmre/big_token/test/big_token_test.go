package test

import (
	"bchain.io/core/actioncontext"
	"bchain.io/utils/database"
	"bchain.io/core/state"
	"bchain.io/common/types"
	"bchain.io/core/transaction"
	"math/big"
	"testing"
	"bchain.io/core/interpreter/wasmre"
	"fmt"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"encoding/json"
	"time"
	"encoding/binary"
)

//go:generate go-bindata -nometadata -pkg test -o bindata.go ../big_token.wasm
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
func TestToken(t *testing.T) {
	code := MustAsset("../big_token.wasm")
	ctx := makeTestCtx()

	wasmre := wasmre.WasmRE{}
	wasmre.Initialize()
	wasmre.Startup()

	rt := wasmre.Generate(ctx)

	fmt.Println("Testcase: token create")
	blkNumber := make([]byte, 8)
	binary.LittleEndian.PutUint64(blkNumber, 998)
	blkNumber_arg := para_paser.Arg{para_paser.TypeI64, blkNumber}

	expiry := make([]byte, 4)
	binary.LittleEndian.PutUint32(expiry, 10)
	expiry_arg := para_paser.Arg{para_paser.TypeI32, expiry}

	isIssue := make([]byte, 4)
	binary.LittleEndian.PutUint32(isIssue, 10)
	isIssue_arg := para_paser.Arg{para_paser.TypeI32, isIssue}

	decimals := para_paser.Arg{para_paser.TypeAddress, append([]byte("18"), 0)}
	sybmol := para_paser.Arg{para_paser.TypeAddress, append([]byte("CNB"), 0)}
	name := para_paser.Arg{para_paser.TypeAddress, append([]byte("china NB"), 0)}
	supply := para_paser.Arg{para_paser.TypeAddress, append([]byte("10000000000"), 0)}
	wp := &para_paser.WasmPara{
		FuncName: "create",
		Args:     append([]para_paser.Arg{}, sybmol, name, decimals, supply, isIssue_arg, blkNumber_arg, expiry_arg),
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
	ret := ctx.ActionResult()[0]
	fmt.Println(string(ret))

	fmt.Println("Testcase: balance tansfer")
	amount := para_paser.Arg{para_paser.TypeAddress, append([]byte("100"), 0)}
	toAddr := para_paser.Arg{para_paser.TypeAddress, append([]byte("0x16ff762e278abb68526d8752e7adfa0eec98c0ba"), 0)}
	memo := para_paser.Arg{para_paser.TypeAddress, append([]byte("this is a test transfer"), 0)}
	wp = &para_paser.WasmPara{
		FuncName: "transfer",
		Args:     append([]para_paser.Arg{}, toAddr, amount, sybmol, memo, blkNumber_arg, expiry_arg),
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
	ret = ctx.ActionResult()[1]
	fmt.Println(string(ret))

	fmt.Println("Testcase: balance of 0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	wp = &para_paser.WasmPara{
		FuncName: "balanceOf",
		Args:     append([]para_paser.Arg{}, toAddr, sybmol),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	ret = ctx.ActionResult()[2]
	fmt.Println(string(ret))

	fmt.Println("Testcase: getSupply")
	wp = &para_paser.WasmPara{
		FuncName: "getSupply",
		Args:     append([]para_paser.Arg{}, sybmol),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	ret = ctx.ActionResult()[3]
	fmt.Println(string(ret))

	fmt.Println("Testcase: getDecimals")
	wp = &para_paser.WasmPara{
		FuncName: "getDecimals",
		Args:     append([]para_paser.Arg{}, sybmol),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	ret = ctx.ActionResult()[4]
	fmt.Println(string(ret))

	wp = &para_paser.WasmPara{
		FuncName: "issue",
		Args:     append([]para_paser.Arg{}, sybmol, supply, memo, blkNumber_arg, expiry_arg),
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
	ret = ctx.ActionResult()[5]
	fmt.Println(string(ret))

	wp = &para_paser.WasmPara{
		FuncName: "balanceOf",
		Args:     append([]para_paser.Arg{}, creator, sybmol),
	}
	paraBytes, _ = json.Marshal(wp)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	ret = ctx.ActionResult()[6]
	fmt.Println(string(ret))
}

func TestXX(t *testing.T)  {
	xx := []byte{1,2,3}
	dd := xx[0:0]
	w, ok := new(big.Int).SetString(string(dd), 10)
	fmt.Println(w.String(), ok)


}