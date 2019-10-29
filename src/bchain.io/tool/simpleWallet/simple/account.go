package simple

import (
	"fmt"
	"bchain.io/accounts/keystore"
	"crypto/ecdsa"
	"bchain.io/common/types"
	"bchain.io/common/assert"
	"errors"
	"bchain.io/utils/crypto"
	"os"
	"encoding/hex"
	"crypto/rand"
)

func AccountCreate(passwd string)  {
	fmt.Println("WellCome to CreateAccount")
	//init myKeyStore
	myKeyStore := keystore.NewKeyStore("./keystore",1 << 18,1)
	if myKeyStore == nil{
		panic("mykeystore== nil")
	}
	//read accounts and print
	acExists := myKeyStore.Accounts()
	for _,ac := range acExists{
		fmt.Printf("Address:%x,   Url:%s\n",ac.Address,ac.URL)
	}
	//create account
	ac,err := myKeyStore.NewAccount(passwd)
	if err != nil{
		panic(err)
	}
	fmt.Printf("Print NewAccount Address:%x,   Url:%s\n",ac.Address,ac.URL)
	//read accounts again and print
	acExists = acExists[:0]
	acExists = myKeyStore.Accounts()
	for _,ac := range acExists{
		fmt.Printf("After Address:%x,   Url:%s\n",ac.Address,ac.URL)
	}
}

func AccountCreateByKey(keyStr, passwd string)  {
	fmt.Println("WellCome to CreateAccount")
	//init myKeyStore
	myKeyStore := keystore.NewKeyStore("./keystore",1 << 18,1)
	if myKeyStore == nil{
		panic("mykeystore== nil")
	}
	//read accounts and print
	acExists := myKeyStore.Accounts()
	for _,ac := range acExists{
		fmt.Printf("Address:%x,   Url:%s\n",ac.Address,ac.URL)
	}

	key, err := crypto.HexToECDSA(keyStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	ac,err := myKeyStore.ImportECDSA(key, passwd)

	fmt.Printf("import Account Address:%x,   Url:%s\n",ac.Address,ac.URL)
	//read accounts again and print
	acExists = acExists[:0]
	acExists = myKeyStore.Accounts()
	for _,ac := range acExists{
		fmt.Printf("After Address:%x,   Url:%s\n",ac.Address,ac.URL)
	}
}

func GetPriKey(passord string) (*ecdsa.PrivateKey, types.Address, error) {
	myKeyStore := keystore.NewKeyStore("./keystore", 1<<18, 1)
	if myKeyStore == nil {
		return nil, types.Address{}, errors.New("mykeystore== nil")
	}
	acExists := myKeyStore.Accounts()
	assert.AssertEx(len(acExists) > 0, "accounts is not exist, please create it")
	key, err := myKeyStore.GetKeyWithPassphrase(acExists[0], passord)
	if err != nil {
		return nil, types.Address{}, err
	}
	addr := acExists[0].Address
	fmt.Println(addr.Hex())
	assert.AsserErr(err)
	return key, addr, nil
}

func NewKey() {
	r := rand.Reader
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), r)
	if err != nil {
		fmt.Println("GenerateKey err", err)
		return
	}
	ret := hex.EncodeToString(crypto.FromECDSA(privateKeyECDSA))
	fmt.Println("privateKey is" ,ret)
}