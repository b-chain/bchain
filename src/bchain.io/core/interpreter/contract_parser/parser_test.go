package contract_parser

import (
	"fmt"
	"bchain.io/core/interpreter/contract/jsre/deps"
	"testing"
)

func TestParseJs(t *testing.T) {
	r := deps.MustAsset("consensus.md")
	c := ParseCodeByData(r)
	fmt.Println("name:", c.InterName)
	fmt.Println("code:", string(c.Code))
}
