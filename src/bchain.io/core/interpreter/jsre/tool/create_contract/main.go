package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"bchain.io/common/types"
	"bchain.io/core/interpreter/contract_parser"
	"bchain.io/core/transaction"
	"os"
)

var inPath = ""
var outPath = ""

//go:generate genAction test.js test.act jsre.JSRE
func main() {
	//check args
	if len(os.Args) != 4 {
		fmt.Println(`
		Notice: need 3 args
		ex:
		genAction	in.file	out.file	jsre.JSRE
		in.file: js code
		out.file: a json data
		`)
		return
	}

	inPath = os.Args[1]
	outPath = os.Args[2]
	/*******************************/
	//read js file
	fin, err := os.Open(inPath)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	fileData, err := ioutil.ReadAll(fin)
	if err != nil {
		panic(err)
	}
	/*******************************/
	//zip data
	//zipJs , err := MinJS(fileData)

	/*******************************/
	//packet result
	result := new(contract_parser.ContractCode)
	result.InterName = os.Args[3]
	result.Code = append(result.Code, fileData...)

	outData, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	paraStr := fmt.Sprintf("createContract('%s')", outData)
	fmt.Println(paraStr)
	act := &transaction.Action{
		Contract: types.HexToAddress("0xFE0604eF5A5D502A4f630bCd796AAA4bC82A813f"),
		Params:   []byte(paraStr),
	}

	actBytes, err := json.MarshalIndent(act, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(actBytes))

	//write result to file
	fout, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	_, err = fout.Write(actBytes)
	if err != nil {
		panic(err)
	}
	fmt.Println("genreate create contract action OK.....")
}
