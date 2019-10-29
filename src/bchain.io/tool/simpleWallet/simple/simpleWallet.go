package simple

import (
	"bchain.io/common/assert"
	"bchain.io/common/types"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"bchain.io/core/transaction"
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/tinylib/msgp/msgp"
	"math/big"
	"bchain.io/communication/rpc/bchainapi"
	"github.com/shopspring/decimal"
	"strconv"
)

func GetAccountNonce(url, addr string) uint64 {
	resp := TxPost(url,"bchain_getAccountNonce", addr, "latest")
	xx := make(map[string]string)
	nc := new(types.Uint64ForJson)
	err := json.Unmarshal(resp, &xx)
	assert.AsserErr(err)
	rlt, ok := xx["result"]
	assert.AssertEx(ok, "result is no exist")
	fmt.Println("nonceOf", addr, "is", rlt)
	err = nc.UnmarshalText([]byte(rlt))
	assert.AsserErr(err)
	return uint64(*nc)
}

var BcContract = "0xb78f12Cb3924607A8BC6a66799e159E3459097e9"
var BcPledgeContract = "0xFa58d9f83D1D86DF22435e67D5F7422337624737"
func MakeBcTransaction(addr types.Address, key *ecdsa.PrivateKey, nc uint64, to string, amount, txFee uint64, memo string) *transaction.Transaction {
	actions := transaction.Actions{}
	data1 := make([]byte, 8)
	binary.LittleEndian.PutUint64(data1, amount)
	amountArg := para_paser.Arg{para_paser.TypeI64, data1}
	toAddr := para_paser.Arg{para_paser.TypeAddress, append([]byte(to), 0)}
	memoArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(memo), 0)}
	wp := &para_paser.WasmPara{
		FuncName: "transfer",
		Args:     append([]para_paser.Arg{}, toAddr, amountArg, memoArg),
	}
	paraBytes, _ := json.Marshal(wp)
	action := transaction.Action{types.HexToAddress(BcContract), paraBytes}

	if txFee > 0 {
		datafee := make([]byte, 8)
		binary.LittleEndian.PutUint64(datafee, txFee)
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
	tx := transaction.NewTransaction(nc, actions)
	txSign, err := transaction.SignTx(tx, s, key)
	assert.AsserErr(err)
	return txSign
}

func MakeBcPledgeTransaction(addr types.Address, key *ecdsa.PrivateKey, nc uint64, amount, txFee uint64) *transaction.Transaction {
	actions := transaction.Actions{}
	data1 := make([]byte, 8)
	binary.LittleEndian.PutUint64(data1, amount)
	amountArg := para_paser.Arg{para_paser.TypeI64, data1}
	wp := &para_paser.WasmPara{
		FuncName: "pledge",
		Args:     append([]para_paser.Arg{}, amountArg),
	}
	paraBytes, _ := json.Marshal(wp)
	action := transaction.Action{types.HexToAddress(BcPledgeContract), paraBytes}

	if txFee > 0 {
		datafee := make([]byte, 8)
		binary.LittleEndian.PutUint64(datafee, txFee)
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
	tx := transaction.NewTransaction(nc, actions)
	txSign, err := transaction.SignTx(tx, s, key)
	assert.AsserErr(err)
	return txSign
}

func MakeBcRedeemTransaction(addr types.Address, key *ecdsa.PrivateKey, nc uint64, amount, txFee uint64) *transaction.Transaction {
	actions := transaction.Actions{}
	data1 := make([]byte, 8)
	binary.LittleEndian.PutUint64(data1, amount)
	amountArg := para_paser.Arg{para_paser.TypeI64, data1}
	wp := &para_paser.WasmPara{
		FuncName: "redeem",
		Args:     append([]para_paser.Arg{}, amountArg),
	}
	paraBytes, _ := json.Marshal(wp)
	action := transaction.Action{types.HexToAddress(BcPledgeContract), paraBytes}

	if txFee > 0 {
		datafee := make([]byte, 8)
		binary.LittleEndian.PutUint64(datafee, txFee)
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
	tx := transaction.NewTransaction(nc, actions)
	txSign, err := transaction.SignTx(tx, s, key)
	assert.AsserErr(err)
	return txSign
}

func SendRawTransaction(url string, tx *transaction.Transaction) {
	var encData bytes.Buffer
	err := msgp.Encode(&encData, tx)
	if err != nil {
		panic(err)
	}
	TxPost(url, "bchain_sendRawTransaction", types.BytesForJson(encData.Bytes()))
	//fmt.Println(string(rsp))
}

func Blockproductor_start(url string, pwd string) {
	TxPost(url, "blockproducer_start", 10, pwd)
}

type jsonRpcRet struct {
	Jsonrpc string        `json:"jsonrpc" `
	Rlt     []types.BytesForJson       `json:"result" `
	Params  []interface{} `json:"params" `
}
func ActionCallBalenceofBc(url string, addr string) {
	Addr := para_paser.Arg{para_paser.TypeAddress, append([]byte(addr), 0)}
	wp := &para_paser.WasmPara{
		FuncName: "balenceOf",
		Args:     append([]para_paser.Arg{}, Addr),
	}
	paraBytes, _ := json.Marshal(wp)
	hexbyte := make(types.BytesForJson, len(paraBytes))
	copy(hexbyte, paraBytes)
	conAddr := types.HexToAddress(BcContract)
	apiAction := bchainapi.SendTxAction{&conAddr, &hexbyte}

	ret := TxPost(url, "bchain_actionCall", apiAction, "latest")
	jsonRet := &jsonRpcRet{}
	err := json.Unmarshal(ret, jsonRet)
	assert.AsserErr(err)
	val := jsonRet.Rlt
	if len(val[0]) >=8 {
		v := binary.LittleEndian.Uint64(val[0])
		vStr := strconv.FormatInt(int64(v),10)
		xx, _ := decimal.NewFromString(vStr)
		deci := decimal.NewFromFloat(100000000)
		xx = xx.Div(deci)
		fmt.Println("balenceOf", addr, "is", xx, "BC")
	}
}

func ActionCallPledgeofBc(url string, addr string) {
	Addr := para_paser.Arg{para_paser.TypeAddress, append([]byte(addr), 0)}
	wp := &para_paser.WasmPara{
		FuncName: "pledgeOfExt",
		Args:     append([]para_paser.Arg{}, Addr),
	}
	paraBytes, _ := json.Marshal(wp)
	hexbyte := make(types.BytesForJson, len(paraBytes))
	copy(hexbyte, paraBytes)
	conAddr := types.HexToAddress(BcPledgeContract)
	apiAction := bchainapi.SendTxAction{&conAddr, &hexbyte}

	ret := TxPost(url, "bchain_actionCall", apiAction, "latest")
	jsonRet := &jsonRpcRet{}
	err := json.Unmarshal(ret, jsonRet)
	assert.AsserErr(err)
	val := jsonRet.Rlt

	v := binary.LittleEndian.Uint64(val[0])
	vStr := strconv.FormatInt(int64(v),10)
	xx, _ := decimal.NewFromString(vStr)
	deci := decimal.NewFromFloat(100000000)
	xx = xx.Div(deci)

	v1 := binary.LittleEndian.Uint64(val[1])
	vStr1 := strconv.FormatInt(int64(v1),10)
	xx1, _ := decimal.NewFromString(vStr1)
	xx1 = xx1.Div(deci)
	fmt.Println("pledgeOf", addr, "is", xx, "BC")
	fmt.Println("pledge pool total is", xx1, "BC")
}