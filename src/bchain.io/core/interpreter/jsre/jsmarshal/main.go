package main

import (
"os"
"io/ioutil"
"fmt"
"encoding/json"
"bchain.io/core/interpreter/contract_parser"
)

var inPath = ""
var outPath = ""


func main(){

	//check args
	if len(os.Args) != 4{
		fmt.Println(`
		Notice: need 3 args
		ex:
		tool.exe in.file   out.file jsre.JSRE
		in.file: js code
		out.file: a json data
		`)
		return
	}

	inPath = os.Args[1]
	outPath = os.Args[2]


	/*******************************/

	//read js file
	fin , err := os.Open(inPath)
	if err != nil {
		panic(err)
	}
	defer fin.Close()


	fileData , err := ioutil.ReadAll(fin)
	if err != nil {
		panic(err)
	}
	/*******************************/
	//zip data
	zipJs , err := MinJS(fileData)

	/*******************************/
	//packet result
	result := new(contract_parser.ContractCode)
	result.InterName = os.Args[3]
	result.Code = append(result.Code , zipJs...)

	outData ,err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	//write result to file
	fout ,err  := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	_ , err  = fout.Write(outData)
	if err != nil {
		panic(err)
	}

	fmt.Println("Dealing Over.....")

}

