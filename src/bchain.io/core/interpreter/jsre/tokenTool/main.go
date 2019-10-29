package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
)

type TemplateJs struct {
	Symbol           string `json:"symbol"`
	TokenName        string `json:"token_name"`
	Decimals         string `json:"decimals"`
	TotalSupply      string `json:"total_supply"`
	BchainContract 	 string	`json:"bchain_contract"`
	Ratio			 string	`json:"ratio"`
}

func (this *TemplateJs)Assert(){
	if len(this.Symbol) == 0 {
		panic("this.Symbol == 0")
	}

	if len(this.TokenName) == 0 {
		panic("this.TokenName == 0")
	}

	if len(this.Decimals) == 0 {
		panic("this.Decimals == 0")
	}

	if len(this.TotalSupply) == 0 {
		panic("this.TotalSupply == 0")
	}

}

var Info = `
	.exe Tmplate.config  out.md
`

func openConfig(path string)*TemplateJs{

	f,err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("Open config failed: %s\n" , err.Error()))
	}
	defer f.Close()

	data , err := ioutil.ReadAll(f)
	if err != nil {
		panic(fmt.Errorf("Read config data failed: %s\n" , err.Error()))
	}

	t := new(TemplateJs)

	err = json.Unmarshal(data , t)
	if err != nil {
		panic(fmt.Errorf("Unmarshal config failed: %s\n" , err.Error()))
	}

	t.Assert()

	return t
}
/*
	Symbol           string `json:"symbol"`
	TokenName        string `json:"token_name"`
	Decimals         string `json:"decimals"`
	TotalSupply      string `json:"total_supply"`
	TotalSupplyTimes string `json:"total_supply_times"`
*/
func replaceData(src string , t *TemplateJs)string{
	//Symbol
	src = strings.Replace(src , "template_Symbol" , t.Symbol , -1)
	//TokenName
	src = strings.Replace(src , "template_tokenName" , t.TokenName , -1)
	//decimals
	src = strings.Replace(src , "template_Decimals" , t.Decimals , -1)
	//totalSupply
	src = strings.Replace(src , "template_totalSupply" , t.TotalSupply , -1)
	//BchainContractAddress
	src = strings.Replace(src , "BchainContractAddress" , t.BchainContract , -1)
	//ratio
	src = strings.Replace(src , "template_Ratio" , t.Ratio , -1)

	return src



}

func writeResult(path  , resultData string){
	f , err := os.Create(path)
	if err != nil {
		panic(fmt.Errorf("Create resultFile failed: %s\n" , err.Error()))
	}
	defer f.Close()

	_ , err = f.Write([]byte(resultData))
	if err != nil {
		panic(fmt.Errorf("Write resultFile failed: %s\n" , err.Error()))
	}
}

func loadSrcJsCode()string{
	f , err := os.Open("templateJsCode.js")
	if err != nil {
		panic(err)
	}
	allData , err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(allData)
}

func main(){
	if len(os.Args) < 2 {
		fmt.Println(Info)
		return
	}

	configFilePath := os.Args[1]
	outputFilePath := os.Args[2]

	srcJsCode := loadSrcJsCode()

	t := openConfig(configFilePath)
	result := replaceData(srcJsCode , t)

	writeResult(outputFilePath , result)

}



