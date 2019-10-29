package main

import (
	"crypto/rand"
	"crypto/ecdsa"
	"bchain.io/utils/crypto"
	"fmt"
	"encoding/hex"
)


func newKey() {
	r := rand.Reader
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), r)
	if err != nil {
		fmt.Println("GenerateKey err", err)
		return
	}
	ret := hex.EncodeToString(crypto.FromECDSA(privateKeyECDSA))
	fmt.Println(ret)
}
func main() {
	newKey()
}
