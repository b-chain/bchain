package tpsTest

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"bchain.io/accounts/keystore"
	"bchain.io/common/assert"
	"bchain.io/common/types"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"bchain.io/core/transaction"
	"net/http"
	"os"
	"unsafe"
)

var g_addr types.Address

func GetPriKey() (*ecdsa.PrivateKey, types.Address) {
	myKeyStore := keystore.NewKeyStore("./keystore", 1<<18, 1)
	if myKeyStore == nil {
		panic("mykeystore== nil")
	}
	acExists := myKeyStore.Accounts()
	assert.AssertEx(len(acExists) > 0, "accounts is not exist, please create it")
	key, err := myKeyStore.GetKeyWithPassphrase(acExists[0], "123")
	g_addr = acExists[0].Address
	fmt.Println(g_addr.Hex())
	assert.AsserErr(err)
	return key, g_addr
}
func MakeTransaction(addr types.Address, key *ecdsa.PrivateKey, nc uint64) *transaction.Transaction {
	actions := transaction.Actions{}

	data1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(data1, uint32(GetConfig().Money))
	amount := para_paser.Arg{para_paser.TypeI32, data1}
	toAddr := para_paser.Arg{para_paser.TypeAddress, append([]byte(GetConfig().To), 0)}
	memo := para_paser.Arg{para_paser.TypeAddress, append([]byte("tps test"), 0)}
	sybmol := para_paser.Arg{para_paser.TypeAddress, append([]byte("CNB"), 0)}
	wp := &para_paser.WasmPara{
		FuncName: "transfer",
		Args:     append([]para_paser.Arg{}, toAddr, amount, sybmol, memo),
	}
	paraBytes, _ := json.Marshal(wp)

	action := transaction.Action{types.HexToAddress(GetConfig().Contract), paraBytes}
	actions = append(actions, &action)

	s := transaction.NewMSigner(big.NewInt(1))
	tx := transaction.NewTransaction(nc, actions)
	txSign, err := transaction.SignTx(tx, s, key)
	assert.AsserErr(err)
	return txSign
}

type jsonRpc struct {
	Jsonrpc string        `json:"jsonrpc" `
	Method  string        `json:"method" `
	Params  []interface{} `json:"params" `
	Id      string        `json:"id" `
}

func GetAccountNonce() uint64 {
	resp := TxPost("bchain_getAccountNonce", g_addr.Hex(), "latest")
	xx := make(map[string]string)
	nc := new(types.Uint64ForJson)
	err := json.Unmarshal(resp, &xx)
	assert.AsserErr(err)
	rlt, ok := xx["result"]
	assert.AssertEx(ok, "result is no exist")
	fmt.Println("get nonce", rlt)
	err = nc.UnmarshalText([]byte(rlt))
	assert.AsserErr(err)
	return uint64(*nc)
}

func TxPost(method string, paras ...interface{}) []byte {
	jsonData := jsonRpc{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  paras,
		Id:      "1",
	}

	bytesData, err := json.Marshal(jsonData)
	assert.AsserErr(err)
	reader := bytes.NewReader(bytesData)
	//url := "http://localhost:7980/"
	url := GetConfig().Url
	request, err := http.NewRequest("POST", url, reader)
	assert.AsserErr(err)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Connection", "keep-alive")
	client := http.Client{}
	resp, err := client.Do(request)
	assert.AsserErr(err)
	respBytes, err := ioutil.ReadAll(resp.Body)
	assert.AsserErr(err)
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println(*str)
	return respBytes
}

type config struct {
	Url      string `json:"url" `
	Contract string `json:"contract" `
	Money    int    `json:"money" `
	To       string `json:"to" `
	Tps      int    `json:"tps" `
}

var g_config *config

func GetConfig() *config {
	if g_config != nil {
		return g_config
	}
	file, err := os.Open("config.txt")
	assert.AsserErr(err)

	all, err := ioutil.ReadAll(file)
	assert.AsserErr(err)
	fmt.Println(string(all))

	c := &config{}
	err = json.Unmarshal(all, c)
	assert.AsserErr(err)
	g_config = c
	return c
}
