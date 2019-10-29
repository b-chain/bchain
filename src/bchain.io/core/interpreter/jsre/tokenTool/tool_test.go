package main

import (
	"testing"
	"encoding/json"
	"fmt"
)

func TestGenerateJson(t *testing.T){
	temp := new(TemplateJs)
	temp.Symbol = "ABC"
	temp.TokenName = "ABC Token"
	temp.Decimals = "18"
	temp.TotalSupply = "1e+10"
	temp.BchainContract = "0x192d52D8cE0c7bBAf0780EAb04860D6Ba012578B"
	temp.Ratio = "10"

	data , err := json.Marshal(temp)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	/*
	{
	"symbol":"ABC",
	"token_name":"ABC Token",
	"decimals":"18",
	"total_supply":"1e+10",
	"bchain_contract":"0x192d52D8cE0c7bBAf0780EAb04860D6Ba012578B",
	"ratio":"10"
	}

	*/
}


func TestModifyJsTemplate(t *testing.T){
	configFilePath := "test.config"
	outputFilePath := "testOutput.js"


	temp := openConfig(configFilePath)
	srcJsCode := loadSrcJsCode()
	result := replaceData(srcJsCode , temp)

	writeResult(outputFilePath , result)
}


