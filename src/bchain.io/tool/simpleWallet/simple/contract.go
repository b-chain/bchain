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
	"os"
	"io/ioutil"
	"bchain.io/core/interpreter/contract_parser"
	"fmt"
)

var SystemContract = "0x2ba8A6318fb0390e8693af78c8086C086D923A96"

func checkErr(err error)  {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func getWasmFile(path string) []byte {
	f, err := os.Open(path)
	checkErr(err)
	defer f.Close()

	fdata, err := ioutil.ReadAll(f)
	checkErr(err)

	cc := new(contract_parser.ContractCode)
	cc.InterName = "wasmre.WasmRE"
	cc.Code = fdata

	jsonData, err := json.Marshal(cc)
	checkErr(err)
	return jsonData
}

func MakeSystemTransaction(key *ecdsa.PrivateKey, nc uint64, txFee uint64, path string) *transaction.Transaction {
	code := getWasmFile(path)
	actions := transaction.Actions{}
	codeArg := para_paser.Arg{para_paser.TypeAddress, append(code, 0)}
	wp := &para_paser.WasmPara{
		FuncName: "createContract",
		Args:     append([]para_paser.Arg{}, codeArg),
	}
	paraBytes, _ := json.Marshal(wp)
	action := transaction.Action{types.HexToAddress(SystemContract), paraBytes}

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
