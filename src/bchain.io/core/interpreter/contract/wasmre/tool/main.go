package main

import (
	"bchain.io/core/interpreter/contract_parser"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println(`
		Notice: need 2 args
		ex:
		tool.exe in.file   out.file
		in.file: wasm code
		`)
		return
	}

	inPath := os.Args[1]
	outPath := os.Args[2]

	fin, err := os.Open(inPath)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	fileData, err := ioutil.ReadAll(fin)
	if err != nil {
		panic(err)
	}

	//packet result
	cc := new(contract_parser.ContractCode)
	cc.InterName = "wasmre.WasmRE"
	cc.Code = append(cc.Code, fileData...)

	outData, err := json.Marshal(cc)
	if err != nil {
		panic(err)
	}

	fout, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	_, err = fout.Write(outData)
	if err != nil {
		panic(err)
	}

	fmt.Println("Dealing Over.....")
}
