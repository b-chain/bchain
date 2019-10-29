////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The bchain-go Authors.
//
// The bchain-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// @File: account_test.go
// @Date: 2018/05/08 17:14:08
////////////////////////////////////////////////////////////////////////////////

package accountTest

import (
	"testing"
	"fmt"
	"bchain.io/accounts/keystore"
	"bchain.io/core/transaction"
	"math/big"
)

//test Account Read
func TestAccountRead(t *testing.T){
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
	fmt.Println("okAccountRead")
}

//test Account Create
func TestAccountCreate(t *testing.T){
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
	ac,err := myKeyStore.NewAccount("123")
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

//test Wallets
func TestWalletsShow(t *testing.T){
	fmt.Println("Wellcome to Wallets show.....")
	//init myKeyStore
	myKeyStore := keystore.NewKeyStore("./keystore",1 << 18,1)
	if myKeyStore == nil{
		panic("mykeystore== nil")
	}

	//check all the wallets
	aw := myKeyStore.Wallets()
	//show all wallets
	for _,w := range aw{
		fmt.Printf("Wallet URL:%s\n" , w.URL())
	}
}

//unlock All account
func TestUnlockAllAccount(t *testing.T){
	fmt.Println("Wellcome to Wallets show.....")
	//init myKeyStore
	myKeyStore := keystore.NewKeyStore("./keystore",1 << 18,1)
	if myKeyStore == nil{
		panic("mykeystore== nil")
	}


	//read accounts and print

	acExists := myKeyStore.Accounts()
	for _,ac := range acExists{
		fmt.Printf("After Address:%x,   Url:%s\n",ac.Address,ac.URL)
		//unlock
		err := myKeyStore.Unlock(ac,"123")
		if err != nil{
			fmt.Println("unlock Wrong....",err)
		}
	}
	//print accounts have been unlocked
	aw := myKeyStore.Wallets()
	//show all wallets
	for _,w := range aw{
		s,_ := w.Status()
		fmt.Printf("wallet: Url:%s,  Status:%s\n" , w.URL(),s)
	}

}

//test transaction signing
func TestTransactionSign(t *testing.T){
	fmt.Println("Wellcome to Wallets show.....")
	myKeyStore := keystore.NewKeyStore("./keystore",1 << 18,1)
	if myKeyStore == nil{
		panic("mykeystore== nil")
	}


	//read accounts and print

	acExists := myKeyStore.Accounts()
	for _,ac := range acExists{
		fmt.Printf("After Address:%x,   Url:%s\n",ac.Address,ac.URL)
		//unlock
		err := myKeyStore.Unlock(ac,"123")
		if err != nil{
			fmt.Println("unlock Wrong....",err)
		}
	}
	//print all accounts have been unlocked
	aw := myKeyStore.Wallets()
	//show all wallets
	for _,w := range aw{
		s,_ := w.Status()
		fmt.Printf("wallet: Url:%s,  Status:%s\n" , w.URL(),s)
	}

	fmt.Println("AccountExists len:" , len(acExists))
	if len(acExists) < 2 {
		t.Skip("Test Transaction Signing Should have 2 or more accounts,Please Run test : TestAccountCreate ")
	}
	//amount:=big.NewInt(int64(20))
	//
	//res:=big.NewInt(int64(10))
	//create transaction,ac[0]--->ac[1]
	newTx:=transaction.NewTransaction(1,nil)

	newTx.PrintVSR()
	//use walet[0] to sign
	tx,err:=aw[0].SignTxWithPassphrase(acExists[0],"123",newTx,big.NewInt(int64(1)))
	if err != nil{
		fmt.Println("err=",err)
		panic("signTxWithPassphrase failed")
	}
	tx.PrintVSR()
	_ = tx
	fmt.Println("over.....")
}
