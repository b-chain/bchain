package simple

import (
	"bchain.io/common/assert"
	"bchain.io/common/types"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"bchain.io/core/transaction"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/json"
	"math/big"
	"bchain.io/communication/rpc/bchainapi"
	"fmt"
)




type Big_token_tr struct {
	Key     *ecdsa.PrivateKey
	ConAddr types.Address

	Nc    uint64
	TxFee uint64

	To     string
	Amount string
	Symbol string
	Memo   string

	BlkNumber uint64
	Expiry    uint32
}

func (b *Big_token_tr) MakeTransaction() *transaction.Transaction {
	actions := transaction.Actions{}

	amountArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Amount), 0)}

	toAddr := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.To), 0)}
	memoArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Memo), 0)}

	symbolArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Symbol), 0)}

	blkNumber := make([]byte, 8)
	binary.LittleEndian.PutUint64(blkNumber, b.BlkNumber)
	blkNumber_arg := para_paser.Arg{para_paser.TypeI64, blkNumber}

	expiry := make([]byte, 4)
	binary.LittleEndian.PutUint32(expiry, b.Expiry)
	expiry_arg := para_paser.Arg{para_paser.TypeI32, expiry}

	wp := &para_paser.WasmPara{
		FuncName: "transfer",
		Args:     append([]para_paser.Arg{}, toAddr, amountArg, symbolArg, memoArg, blkNumber_arg, expiry_arg),
	}
	paraBytes, _ := json.Marshal(wp)
	action := transaction.Action{b.ConAddr, paraBytes}

	if b.TxFee > 0 {
		datafee := make([]byte, 8)
		binary.LittleEndian.PutUint64(datafee, b.TxFee)
		feeArg := para_paser.Arg{para_paser.TypeI64, datafee}
		wp = &para_paser.WasmPara{
			FuncName: "transferFee",
			Args:     append([]para_paser.Arg{}, feeArg),
		}
		paraBytes, _ = json.Marshal(wp)
		actionFee := transaction.Action{types.HexToAddress(BcContract), paraBytes}
		actions = append(actions, &actionFee, &action)
	} else {
		actions = append(actions, &action)
	}

	s := transaction.NewMSigner(big.NewInt(1))
	tx := transaction.NewTransaction(b.Nc, actions)
	txSign, err := transaction.SignTx(tx, s, b.Key)
	assert.AsserErr(err)
	return txSign
}

type Big_token_create struct {
	Key     *ecdsa.PrivateKey
	ConAddr types.Address

	Nc    uint64
	TxFee uint64


	Symbol     string
	Name       string
	Decimals   string
	Supply     string
	IsIssue    uint32

	BlkNumber uint64
	Expiry    uint32
}

func (b *Big_token_create) MakeTransaction() *transaction.Transaction {
	actions := transaction.Actions{}

	NameArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Name), 0)}
	DecimalsArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Decimals), 0)}
	SupplyArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Supply), 0)}

	symbolArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Symbol), 0)}

	blkNumber := make([]byte, 8)
	binary.LittleEndian.PutUint64(blkNumber, b.BlkNumber)
	blkNumber_arg := para_paser.Arg{para_paser.TypeI64, blkNumber}

	expiry := make([]byte, 4)
	binary.LittleEndian.PutUint32(expiry, b.Expiry)
	expiry_arg := para_paser.Arg{para_paser.TypeI32, expiry}

	issule := make([]byte, 4)
	binary.LittleEndian.PutUint32(issule, b.IsIssue)
	issule_arg := para_paser.Arg{para_paser.TypeI32, issule}

	wp := &para_paser.WasmPara {
		FuncName: "create",
		Args:     append([]para_paser.Arg{}, symbolArg, NameArg, DecimalsArg, SupplyArg, issule_arg, blkNumber_arg, expiry_arg),
	}
	paraBytes, _ := json.Marshal(wp)
	action := transaction.Action{b.ConAddr, paraBytes}

	if b.TxFee > 0 {
		datafee := make([]byte, 8)
		binary.LittleEndian.PutUint64(datafee, b.TxFee)
		feeArg := para_paser.Arg{para_paser.TypeI64, datafee}
		wp = &para_paser.WasmPara {
			FuncName: "transferFee",
			Args:     append([]para_paser.Arg{}, feeArg),
		}
		paraBytes, _ = json.Marshal(wp)
		actionFee := transaction.Action{types.HexToAddress(BcContract), paraBytes}
		actions = append(actions, &actionFee, &action)
	} else {
		actions = append(actions, &action)
	}

	s := transaction.NewMSigner(big.NewInt(1))
	tx := transaction.NewTransaction(b.Nc, actions)
	txSign, err := transaction.SignTx(tx, s, b.Key)
	assert.AsserErr(err)
	return txSign
}

type Big_token_issue struct {
	Key     *ecdsa.PrivateKey
	ConAddr types.Address

	Nc    uint64
	TxFee uint64


	Symbol     string
	Amount     string
	Memo       string

	BlkNumber uint64
	Expiry    uint32
}

func (b *Big_token_issue) MakeTransaction() *transaction.Transaction {
	actions := transaction.Actions{}

	symbolArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Symbol), 0)}
	AmountArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Amount), 0)}
	memoArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Memo), 0)}

	blkNumber := make([]byte, 8)
	binary.LittleEndian.PutUint64(blkNumber, b.BlkNumber)
	blkNumber_arg := para_paser.Arg{para_paser.TypeI64, blkNumber}

	expiry := make([]byte, 4)
	binary.LittleEndian.PutUint32(expiry, b.Expiry)
	expiry_arg := para_paser.Arg{para_paser.TypeI32, expiry}

	wp := &para_paser.WasmPara {
		FuncName: "issue",
		Args:     append([]para_paser.Arg{}, symbolArg, AmountArg, memoArg, blkNumber_arg, expiry_arg),
	}
	paraBytes, _ := json.Marshal(wp)
	action := transaction.Action{b.ConAddr, paraBytes}

	if b.TxFee > 0 {
		datafee := make([]byte, 8)
		binary.LittleEndian.PutUint64(datafee, b.TxFee)
		feeArg := para_paser.Arg{para_paser.TypeI64, datafee}
		wp = &para_paser.WasmPara {
			FuncName: "transferFee",
			Args:     append([]para_paser.Arg{}, feeArg),
		}
		paraBytes, _ = json.Marshal(wp)
		actionFee := transaction.Action{types.HexToAddress(BcContract), paraBytes}
		actions = append(actions, &actionFee, &action)
	} else {
		actions = append(actions, &action)
	}

	s := transaction.NewMSigner(big.NewInt(1))
	tx := transaction.NewTransaction(b.Nc, actions)
	txSign, err := transaction.SignTx(tx, s, b.Key)
	assert.AsserErr(err)
	return txSign
}


func ActionCallBalanceOfBt(url string, addr, symbol string, contract types.Address) {
	Addr := para_paser.Arg{para_paser.TypeAddress, append([]byte(addr), 0)}
	smbloArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(symbol), 0)}
	wp := &para_paser.WasmPara{
		FuncName: "balanceOf",
		Args:     append([]para_paser.Arg{}, Addr, smbloArg),
	}
	paraBytes, _ := json.Marshal(wp)
	hexbyte := make(types.BytesForJson, len(paraBytes))
	copy(hexbyte, paraBytes)

	apiAction := bchainapi.SendTxAction{&contract, &hexbyte}

	ret := TxPost(url, "bchain_actionCall", apiAction, "latest")
	jsonRet := &jsonRpcRet{}
	err := json.Unmarshal(ret, jsonRet)
	assert.AsserErr(err)
	val := jsonRet.Rlt
	fmt.Println("balanceOf", addr, "is", string(val[0]), symbol, "unit")
}

func ActionCallGetSupplyBt(url string, symbol string, contract types.Address) {
	smbloArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(symbol), 0)}
	wp := &para_paser.WasmPara{
		FuncName: "getSupply",
		Args:     append([]para_paser.Arg{}, smbloArg),
	}
	paraBytes, _ := json.Marshal(wp)
	hexbyte := make(types.BytesForJson, len(paraBytes))
	copy(hexbyte, paraBytes)

	apiAction := bchainapi.SendTxAction{&contract, &hexbyte}

	ret := TxPost(url, "bchain_actionCall", apiAction, "latest")
	jsonRet := &jsonRpcRet{}
	err := json.Unmarshal(ret, jsonRet)
	assert.AsserErr(err)
	val := jsonRet.Rlt
	fmt.Println("supply of", symbol, "is", string(val[0]))
}

func ActionCallGetDecimalsBt(url string, symbol string, contract types.Address) {
	smbloArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(symbol), 0)}
	wp := &para_paser.WasmPara{
		FuncName: "getDecimals",
		Args:     append([]para_paser.Arg{}, smbloArg),
	}
	paraBytes, _ := json.Marshal(wp)
	hexbyte := make(types.BytesForJson, len(paraBytes))
	copy(hexbyte, paraBytes)

	apiAction := bchainapi.SendTxAction{&contract, &hexbyte}

	ret := TxPost(url, "bchain_actionCall", apiAction, "latest")
	jsonRet := &jsonRpcRet{}
	err := json.Unmarshal(ret, jsonRet)
	assert.AsserErr(err)
	val := jsonRet.Rlt
	fmt.Println("decimals of", symbol, "is", string(val[0]))
}

func ActionCallGetNameBt(url string, symbol string, contract types.Address) {
	smbloArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(symbol), 0)}
	wp := &para_paser.WasmPara{
		FuncName: "getName",
		Args:     append([]para_paser.Arg{}, smbloArg),
	}
	paraBytes, _ := json.Marshal(wp)
	hexbyte := make(types.BytesForJson, len(paraBytes))
	copy(hexbyte, paraBytes)

	apiAction := bchainapi.SendTxAction{&contract, &hexbyte}

	ret := TxPost(url, "bchain_actionCall", apiAction, "latest")
	jsonRet := &jsonRpcRet{}
	err := json.Unmarshal(ret, jsonRet)
	assert.AsserErr(err)
	val := jsonRet.Rlt
	fmt.Println("name of", symbol, "is", string(val[0]))
}