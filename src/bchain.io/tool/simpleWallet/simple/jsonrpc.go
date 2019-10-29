package simple

import (
	"encoding/json"
	"bchain.io/common/assert"
	"bytes"
	"net/http"
	"io/ioutil"
	"unsafe"
	"fmt"
)

type jsonRpc struct {
	Jsonrpc string        `json:"jsonrpc" `
	Method  string        `json:"method" `
	Params  []interface{} `json:"params" `
	Id      string        `json:"id" `
}

func TxPost(url string, method string, paras ...interface{}) []byte {
	jsonData := jsonRpc{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  paras,
		Id:      "1",
	}

	bytesData, err := json.Marshal(jsonData)
	assert.AsserErr(err)
	reader := bytes.NewReader(bytesData)
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
