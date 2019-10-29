package contract_parser

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type ContractCode struct {
	InterName string `json:"inter_name"`
	Code      []byte `json:"code"`
}

func ParseCodeByData(data []byte) *ContractCode {
	if data == nil {
		panic(fmt.Errorf("ParseCode Err:data == nil"))
	}

	c := new(ContractCode)
	err := json.Unmarshal(data, c)
	if err != nil {
		panic(fmt.Errorf("ParseCode Err:%s", err.Error()))
	}
	return c
}

func ParseCodeByFile(inPath string) *ContractCode {
	//open file
	fin, err := os.Open(inPath)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	//read all data
	fileData, err := ioutil.ReadAll(fin)
	if err != nil {
		panic(err)
	}

	return ParseCodeByData(fileData)

}

func ParseCodeByReader(rd io.Reader) *ContractCode {

	//read all data
	fileData, err := ioutil.ReadAll(rd)
	if err != nil {
		panic(err)
	}

	return ParseCodeByData(fileData)

}
