package block

import (
	"bytes"
	"fmt"
	"math/big"
	"bchain.io/common/types"
	"bchain.io/utils/crypto"
	"testing"
	"github.com/tinylib/msgp/msgp"
	"bchain.io/core/transaction"
)

func TestHeaderSignatureRamdomkey(t *testing.T) {
	header := &Header{Number: types.NewBigInt(*big.NewInt(334)), Time: types.NewBigInt(*big.NewInt(1212121))}
	chainId := big.NewInt(101)
	singner := NewBlockSigner(chainId)

	var (
		key, _  = crypto.GenerateKey()
		address = crypto.PubkeyToAddress(key.PublicKey)
	)
	signHeaer, err := SignHeader(header, singner, key)
	if err != nil {
		t.Fatalf("SignHeader fail")
	}

	getaddress, err := singner.Sender(signHeaer)
	if err != nil {
		t.Fatalf("cann't get senser form header %v", err)
	}
	fmt.Println(signHeaer)

	if !bytes.Equal(getaddress.Bytes(), address.Bytes()) {
		t.Fatalf("address is not same got:%v, want:%v", getaddress.Hex(), address.Hex())
	}
}

func TestHeaderSigantureFixkey(t *testing.T) {
	conData := make([]byte, 10)
	conData[4] = 7
	Producer := types.Address{}
	Producer[10] = 1

	header := &Header{Number: types.NewBigInt(*big.NewInt(333)), Time: types.NewBigInt(*big.NewInt(1212121)), Producer: Producer, Cdata: ConsensusData{"test", conData}, Extra:[]byte("asdfasd")}
	chainId := big.NewInt(101)
	singner := NewBlockSigner(chainId)

	var (
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f292")
		address = crypto.PubkeyToAddress(key.PublicKey)
	)

	//fmt.Println(header)

	signHeaer, err := SignHeader(header, singner, key)
	if err != nil {
		t.Fatalf("SignHeader fail")
	}

	getaddress, err := singner.Sender(signHeaer)
	if err != nil {
		t.Fatalf("cann't get senser form header %v", err)
	}
	fmt.Println(signHeaer)

	if !bytes.Equal(getaddress.Bytes(), address.Bytes()) {
		t.Fatalf("address is not same got:%v, want:%v", getaddress.Hex(), address.Hex())
	}
}

func TestBlock_Size(t *testing.T) {
	header := &Header{Number: types.NewBigInt(*big.NewInt(334)), Time: types.NewBigInt(*big.NewInt(1212121))}
	chainId := big.NewInt(101)
	singner := NewBlockSigner(chainId)

	var (
		key, _  = crypto.GenerateKey()
	)
	signHeaer, err := SignHeader(header, singner, key)
	if err != nil {
		t.Fatalf("SignHeader fail")
	}
	fmt.Println(signHeaer)
	fmt.Println(signHeaer.Msgsize())
	fmt.Println(header.Msgsize())
	var buf bytes.Buffer
	msgp.Encode(&buf, signHeaer)
	c := buf.Len()
	fmt.Println(c)
	tx1 := transaction.NewTransaction(1, transaction.Actions{&transaction.Action{Contract: types.BytesToAddress([]byte{0x11}), Params: []byte{0x11, 0x11, 0x11, 44}}})
	tx2 := transaction.NewTransaction(2, transaction.Actions{&transaction.Action{Contract: types.BytesToAddress([]byte{0x22}), Params: []byte{0x22, 0x22, 0x22}}})
	tx3 := transaction.NewTransaction(3, transaction.Actions{&transaction.Action{Contract: types.BytesToAddress([]byte{0x33}), Params: []byte{0x33, 0x33, 0x33}}})
	txs := []*transaction.Transaction{tx1, tx2, tx3}
	fmt.Println(tx1.Size())
	fmt.Println(tx1.Msgsize())

	block :=NewBlock(signHeaer, txs, nil)
	fmt.Println(block.Size())

	blockE :=NewBlock(signHeaer, nil, nil)
	fmt.Println(blockE.Size())
}

