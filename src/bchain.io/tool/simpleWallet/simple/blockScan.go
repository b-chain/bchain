package simple

import (
	"encoding/json"
	"bchain.io/communication/rpc/bchainapi"
	"fmt"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"encoding/binary"
	"os"
	"bchain.io/common/types"
	"bchain.io/common/assert"
)

type TrScanRlt struct {
	From   string
	To     string
	Amount uint64
	Memo   string
}


func GetBlockByNumer(url, blkNumber string) []TrScanRlt {
	resp := TxPost(url,"bchain_getBlockByNumber", blkNumber, true)
	fields := make(map[string]interface{})
	err := json.Unmarshal(resp, &fields)
	if err != nil {
		panic(err)
	}
	rlt := fields["result"]
	rltMap := rlt.(map[string]interface{})

	tr := rltMap["transactions"]
	trs := tr.([]interface{})

	trRlts := []TrScanRlt{}
	for _, tr := range trs {
		rpcTr := &bchainapi.RPCTransaction{}
		trJson,err := json.Marshal(&tr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = json.Unmarshal(trJson, rpcTr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, aa := range rpcTr.Actions {
			if aa.Address.HexLower() == "0xb78f12cb3924607a8bc6a66799e159e3459097e9" {
				//fmt.Println(aa.Address.HexLower(), string(*aa.Params))
				wp := &para_paser.WasmPara{}
				err = json.Unmarshal(*aa.Params, wp)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if wp.FuncName == "transfer" {
					trRlt := TrScanRlt{
						From: rpcTr.From.HexLower(),
						To:string(wp.Args[0].Data[0:len(wp.Args[0].Data)-1]),
						Amount: binary.LittleEndian.Uint64(wp.Args[1].Data),
						Memo:string(wp.Args[2].Data[0:len(wp.Args[2].Data)-1]),
					}
					trRlts = append(trRlts, trRlt)
				}

			}
		}
	}
	trRltsStr, err := json.MarshalIndent(&trRlts, "", "	")
	fmt.Println(string(trRltsStr))
	return trRlts
}

func GetBlocCertificateByNumer(url, blkNumber string) {
	resp := TxPost(url,"bchain_getBlockCertificateByNumber", blkNumber)
	fields := make(map[string]interface{})
	err := json.Unmarshal(resp, &fields)
	if err != nil {
		panic(err)
	}
	rlt := fields["result"]
	dataJson,err := json.Marshal(&rlt)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	cert := []bchainapi.Certificate{}
	err = json.Unmarshal(dataJson, &cert)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	trRltsStr, err := json.MarshalIndent(&cert, "", "	")
	fmt.Println(string(trRltsStr))
}

func GetBlockNumer(url string) uint64 {
	resp := TxPost(url,"bchain_blockNumber")
	xx := make(map[string]string)
	nc := new(types.Uint64ForJson)
	err := json.Unmarshal(resp, &xx)
	checkErr(err)
	rlt, ok := xx["result"]
	assert.AssertEx(ok, "result is no exist")
	fmt.Println("current block number is", rlt)
	err = nc.UnmarshalText([]byte(rlt))
	checkErr(err)
	return uint64(*nc)
}