package nativere

import (
	"bchain.io/common/types"
	"bchain.io/core/actioncontext"
	"bchain.io/utils/database"
	"bchain.io/core/state"
	"math/big"
	"testing"
	"bchain.io/core/transaction"
	"fmt"
	"encoding/json"
	"time"
)

func TestRewordVal(t *testing.T) {
	sender := types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	blkctx := makeTestBlkCtx(1, sender)
	tr := transaction.Action{}
	tr.Contract = types.Address{}
	ctx := actioncontext.NewContext(sender, &tr, blkctx)

	nativeRe := NativeRe{}
	nativeRe.Initialize()
	nativeRe.Startup()

	p := pledge{ctx}
	ret := p.getRewordsValue(24000000000000)
	fmt.Println(ret)

}

func makeTestBlkCtx(blkNumber int64,producer types.Address) *actioncontext.BlockContext{
	db, _ := database.OpenMemDB()
	stateDb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	tmpdb, _ := database.OpenMemDB()


	return actioncontext.NewBlockContext(stateDb, db, tmpdb, big.NewInt(blkNumber), producer)
}

func printBalance(blkctx *actioncontext.BlockContext, addr string)  {
	nativeRe := NativeRe{}
	nativeRe.Initialize()
	nativeRe.Startup()
	sender := types.HexToAddress(addr)
	tr := transaction.Action{}
	tr.Contract = types.Address{}
	ctx := actioncontext.NewContext(sender, &tr, blkctx)
	rt := nativeRe.Generate(ctx)
	np := &NativePara {
		FuncName: "pledgeOfExt",
		Args:     append([]string{}, addr),
	}
	paraBytes, _ := json.Marshal(np)

	fmt.Println(string(paraBytes))
	code := []byte{0,0,0}
	rt.Exec(code, paraBytes, 1000*time.Second)

	ret := ctx.ActionResult()
	fmt.Println(string(ret[0]), string(ret[1]), string(ret[2]), string(ret[3]), string(ret[4]), string(ret[5]))
}

func TestNativeContract(t *testing.T) {
	sender := types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	blkctx := makeTestBlkCtx(1, sender)
	tr := transaction.Action{}
	tr.Contract = types.Address{}
	ctx := actioncontext.NewContext(sender, &tr, blkctx)

	nativeRe := NativeRe{}
	nativeRe.Initialize()
	nativeRe.Startup()

	rt := nativeRe.Generate(ctx)

	fmt.Println("Testcase: create contract")
	np := &NativePara {
		FuncName: "rewords",
		Args:     append([]string{}),
	}
	paraBytes, _ := json.Marshal(np)
	fmt.Println(string(paraBytes))

	code := []byte{0,0,0}
	d := rt.Exec(code, paraBytes, 1*time.Second)
	fmt.Println("exec time:", d)

	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")


	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98c0ba", "1000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98aaaa", "100000000000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98aaaa")


	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98bbbb", "100000000000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)

	np = &NativePara {
		FuncName: "transferFee",
		Args:     append([]string{},  "99"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98c0ba")

	//pledge
	sender = types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "500000000000", "0x3e68b0583021d78c122f719fc82036529a903333", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	p := pledge{ctx}
	p.dumpPool("0x3e68b0583021d78c122f719fc82036529a903333")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98aaaa")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "1500000000", "0x3e68b0583021d78c122f719fc82036529a903333", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98aaaa")
	printBalance(blkctx,"0x3e68b0583021d78c122f719fc82036529a903333")
	p.dumpPool("0x3e68b0583021d78c122f719fc82036529a903333")

	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "1500000000", "0x3e68b0583021d78c122f719fc82036529a903333", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98aaaa")
	printBalance(blkctx,"0x3e68b0583021d78c122f719fc82036529a903333")
	p.dumpPool("0x3e68b0583021d78c122f719fc82036529a903333")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "2500000000", "0x3e68b0583021d78c122f719fc82036529a903333", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	p.dumpPool("0x3e68b0583021d78c122f719fc82036529a903333")

	time.Sleep(time.Second)


}

func TestNativeContract11(t *testing.T) {
	sender := types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	blkctx := makeTestBlkCtx(1, sender)
	tr := transaction.Action{}
	tr.Contract = types.Address{}
	ctx := actioncontext.NewContext(sender, &tr, blkctx)

	nativeRe := NativeRe{}
	nativeRe.Initialize()
	nativeRe.Startup()

	rt := nativeRe.Generate(ctx)

	fmt.Println("Testcase: create contract")
	np := &NativePara {
		FuncName: "rewords",
		Args:     append([]string{}),
	}
	paraBytes, _ := json.Marshal(np)
	fmt.Println(string(paraBytes))

	code := []byte{0,0,0}
	d := rt.Exec(code, paraBytes, 1*time.Second)
	fmt.Println("exec time:", d)

	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")


	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98c0ba", "1000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98aaaa", "100000000000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98aaaa")


	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98bbbb", "100000000000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)

	//pledge
	sender = types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "500000000000", "0x2e68b0583021d78c122f719fc82036529a90571d", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	p := pledge{ctx}
	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98aaaa")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "1500000000", "0x2e68b0583021d78c122f719fc82036529a90571d", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")

	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "1500000000", "0x2e68b0583021d78c122f719fc82036529a90571d", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "2500000000", "0x2e68b0583021d78c122f719fc82036529a90571d", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	np = &NativePara {
		FuncName: "makeProducer",
		Args:     append([]string{},  "50000000000", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	sender = types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)

	np = &NativePara {
		FuncName: "proxy",
		Args:     append([]string{},  "0x16ff762e278abb68526d8752e7adfa0eec98bbbb"),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	p.dumpProducer("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	time.Sleep(time.Second)
}


func TestNativeContractProxy(t *testing.T) {
	sender := types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	blkctx := makeTestBlkCtx(1, sender)
	tr := transaction.Action{}
	tr.Contract = types.Address{}
	ctx := actioncontext.NewContext(sender, &tr, blkctx)

	nativeRe := NativeRe{}
	nativeRe.Initialize()
	nativeRe.Startup()

	rt := nativeRe.Generate(ctx)

	fmt.Println("Testcase: create contract")
	np := &NativePara {
		FuncName: "rewords",
		Args:     append([]string{}),
	}
	paraBytes, _ := json.Marshal(np)
	fmt.Println(string(paraBytes))

	code := []byte{0,0,0}
	d := rt.Exec(code, paraBytes, 1*time.Second)
	fmt.Println("exec time:", d)

	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")


	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98c0ba", "1000000000000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98aaaa", "1000000000000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98aaaa")


	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98bbbb", "1000000000000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)

	//pledge
	sender = types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "500000000000", "0x2e68b0583021d78c122f719fc82036529a90571d", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	p := pledge{ctx}
	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98aaaa")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "500000000000", "0x16ff762e278abb68526d8752e7adfa0eec98aaaa", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "500000000000", "0x16ff762e278abb68526d8752e7adfa0eec98bbbb", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	np = &NativePara {
		FuncName: "makeProducer",
		Args:     append([]string{},  "50000000000", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	sender = types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)

	np = &NativePara {
		FuncName: "proxy",
		Args:     append([]string{},  "0x16ff762e278abb68526d8752e7adfa0eec98bbbb"),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	p.dumpProducer("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)

	np = &NativePara {
		FuncName: "proxy",
		Args:     append([]string{},  "0x16ff762e278abb68526d8752e7adfa0eec98bbbb"),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	p.dumpProducer("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98aaaa")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)

	np = &NativePara {
		FuncName: "proxy",
		Args:     append([]string{},  "0x16ff762e278abb68526d8752e7adfa0eec98bbbb"),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	p.dumpProducer("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")
	time.Sleep(time.Second)
}

func TestNativeContract113(t *testing.T) {
	sender := types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	blkctx := makeTestBlkCtx(1, sender)
	tr := transaction.Action{}
	tr.Contract = types.Address{}
	ctx := actioncontext.NewContext(sender, &tr, blkctx)

	nativeRe := NativeRe{}
	nativeRe.Initialize()
	nativeRe.Startup()

	rt := nativeRe.Generate(ctx)

	fmt.Println("Testcase: create contract")
	np := &NativePara {
		FuncName: "rewords",
		Args:     append([]string{}),
	}
	paraBytes, _ := json.Marshal(np)
	fmt.Println(string(paraBytes))

	code := []byte{0,0,0}
	d := rt.Exec(code, paraBytes, 1*time.Second)
	fmt.Println("exec time:", d)

	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")


	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98c0ba", "1000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98aaaa", "100000000000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98aaaa")


	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98bbbb", "100000000000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)

	//pledge
	sender = types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "500000000000", "0x2e68b0583021d78c122f719fc82036529a90571d", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	p := pledge{ctx}
	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98aaaa")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "1500000000", "0x2e68b0583021d78c122f719fc82036529a90571d", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")

	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "1500000000", "0x2e68b0583021d78c122f719fc82036529a90571d", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "2500000000", "0x2e68b0583021d78c122f719fc82036529a90571d", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	np = &NativePara {
		FuncName: "makeProducer",
		Args:     append([]string{},  "50000000000", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	sender = types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)

	np = &NativePara {
		FuncName: "proxy",
		Args:     append([]string{},  "0x16ff762e278abb68526d8752e7adfa0eec98bbbb"),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	p.dumpProducer("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "2500000000", "0x2e68b0583021d78c122f719fc82036529a90571d", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")
	p.dumpProducer("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	time.Sleep(time.Second)
}

func TestNativeCancelProxy(t *testing.T) {
	sender := types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	blkctx := makeTestBlkCtx(1, sender)
	tr := transaction.Action{}
	tr.Contract = types.Address{}
	ctx := actioncontext.NewContext(sender, &tr, blkctx)

	nativeRe := NativeRe{}
	nativeRe.Initialize()
	nativeRe.Startup()

	rt := nativeRe.Generate(ctx)

	fmt.Println("Testcase: create contract")
	np := &NativePara {
		FuncName: "rewords",
		Args:     append([]string{}),
	}
	paraBytes, _ := json.Marshal(np)
	fmt.Println(string(paraBytes))

	code := []byte{0,0,0}
	d := rt.Exec(code, paraBytes, 1*time.Second)
	fmt.Println("exec time:", d)

	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")


	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98c0ba", "1000000000000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98aaaa", "1000000000000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98aaaa")


	np = &NativePara {
		FuncName: "transfer",
		Args:     append([]string{}, "0x16ff762e278abb68526d8752e7adfa0eec98bbbb", "1000000000000", "xx"),
	}
	paraBytes, _ = json.Marshal(np)

	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98c0ba")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)

	//pledge
	sender = types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "500000000000", "0x2e68b0583021d78c122f719fc82036529a90571d", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	p := pledge{ctx}
	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98aaaa")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "500000000000", "0x16ff762e278abb68526d8752e7adfa0eec98aaaa", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "pledge",
		Args:     append([]string{},  "500000000000", "0x16ff762e278abb68526d8752e7adfa0eec98bbbb", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)

	p.dumpPool("0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	np = &NativePara {
		FuncName: "makeProducer",
		Args:     append([]string{},  "50000000000", ),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	printBalance(blkctx,"0x2e68b0583021d78c122f719fc82036529a90571d")
	printBalance(blkctx,"0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	sender = types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)

	np = &NativePara {
		FuncName: "proxy",
		Args:     append([]string{},  "0x16ff762e278abb68526d8752e7adfa0eec98bbbb"),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	p.dumpProducer("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)

	np = &NativePara {
		FuncName: "proxy",
		Args:     append([]string{},  "0x16ff762e278abb68526d8752e7adfa0eec98bbbb"),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	p.dumpProducer("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	sender = types.HexToAddress("0x16ff762e278abb68526d8752e7adfa0eec98aaaa")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)

	np = &NativePara {
		FuncName: "proxy",
		Args:     append([]string{},  "0x16ff762e278abb68526d8752e7adfa0eec98bbbb"),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	p.dumpProducer("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")

	sender = types.HexToAddress("0x2e68b0583021d78c122f719fc82036529a90571d")
	ctx = actioncontext.NewContext(sender, &tr, blkctx)
	p = pledge{ctx}
	rt = nativeRe.Generate(ctx)
	np = &NativePara {
		FuncName: "cancelProxy",
		Args:     append([]string{}),
	}
	paraBytes, _ = json.Marshal(np)
	fmt.Println(string(paraBytes))
	d = rt.Exec(code, paraBytes, 1000*time.Second)
	fmt.Println("exec time:", d)
	p.dumpProducer("0x16ff762e278abb68526d8752e7adfa0eec98bbbb")
	time.Sleep(time.Second)
}