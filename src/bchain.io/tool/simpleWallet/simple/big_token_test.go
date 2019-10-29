package simple

import (
	"bchain.io/utils/crypto"
	"bchain.io/common/types"
	"bchain.io/core/transaction"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"encoding/binary"
	"encoding/json"
	"math/big"
	"bchain.io/common/assert"
	"fmt"
	"testing"
	"crypto/ecdsa"
	"bchain.io/utils/crypto/sha3"
	"github.com/tinylib/msgp/msgp"
	"bytes"
)

func TestEncodeBC(t *testing.T)  {
	key, err := crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	checkErr(err)

	fmt.Println("-------------------input para start")
	fmt.Println("key", "49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	fmt.Println("")
	xx := MakeBcTransactionTest(types.Address{}, key, 55, "0x778181dd9b0382b9f1c32604b9d9332111c33fda", 5656, 0, "test mono")
	fmt.Println("--tx with signature")
	fmt.Println(xx)
	fmt.Println("")

	SendRawTransactionTest(xx)

}

func MakeBcTransactionTest(addr types.Address, key *ecdsa.PrivateKey, nc uint64, to string, amount, txFee uint64, memo string) *transaction.Transaction {
	actions := transaction.Actions{}
	data1 := make([]byte, 8)
	binary.LittleEndian.PutUint64(data1, amount)
	amountArg := para_paser.Arg{para_paser.TypeI64, data1}
	toAddr := para_paser.Arg{para_paser.TypeAddress, append([]byte(to), 0)}
	memoArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(memo), 0)}
	wp := &para_paser.WasmPara{
		FuncName: "transfer",
		Args:     append([]para_paser.Arg{}, toAddr, amountArg, memoArg),
	}
	paraBytes, _ := json.Marshal(wp)
	fmt.Println("-------------------encode start")
	fmt.Println("--contract para")
	fmt.Println(string(paraBytes))
	fmt.Println("")
	action := transaction.Action{types.HexToAddress(BcContract), paraBytes}

	if txFee > 0 {
		datafee := make([]byte, 8)
		binary.LittleEndian.PutUint64(datafee, txFee)
		feeArg := para_paser.Arg{para_paser.TypeI64, datafee}
		wp = &para_paser.WasmPara{
			FuncName: "transferFee",
			Args:     append([]para_paser.Arg{}, feeArg),
		}
		paraBytes, _ = json.Marshal(wp)
		actionFee := transaction.Action{types.HexToAddress(BcContract), paraBytes}
		actions = append(actions, &actionFee, &action)
	} else {
		actions = append(actions, &action)
	}

	s := transaction.NewMSigner(big.NewInt(1))
	tx := transaction.NewTransaction(nc, actions)
	txSign, err := SignTx(tx, s, key)
	assert.AsserErr(err)
	return txSign
}

func TestEncode(t *testing.T)  {
	key, err := crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	checkErr(err)

	bt := Big_token_tr{
		Key: key,
		ConAddr:types.HexToAddress("0x678181dd9b0382b9fac32604b9d9332111c33fdb"),

		Nc:55,
		TxFee:uint64(0),

		To:"0x778181dd9b0382b9f1c32604b9d9332111c33fda",
		Amount: "5656",
		Symbol:"CNY",
		Memo:"cny test memo",

		BlkNumber: 998,
		Expiry: 100,
	}
	para , _ := json.MarshalIndent(&bt,"", "	")
	fmt.Println("-------------------input para start")
	fmt.Println("key", "49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	fmt.Println(string(para))
	fmt.Println("")
	xx := bt.MakeTransactionTest()
	fmt.Println("--tx with signature")
	fmt.Println(xx)
	fmt.Println("")

	SendRawTransactionTest(xx)

}

func SendRawTransactionTest(tx *transaction.Transaction) {
	var encData bytes.Buffer
	err := msgp.Encode(&encData, tx)
	if err != nil {
		panic(err)
	}

	fmt.Println("--msgp transaction")

	fmt.Printf("%x\n",encData)
	fmt.Println(encData.String())

	TxPost11("x11", "bchain_sendRawTransaction", types.BytesForJson(encData.Bytes()))

}

func TxPost11(url string, method string, paras ...interface{})  {
	jsonData := jsonRpc{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  paras,
		Id:      "1",
	}

	bytesData, _ := json.MarshalIndent(&jsonData, "", "	")
	fmt.Println(string(bytesData))

}

func (b *Big_token_tr) MakeTransactionTest() *transaction.Transaction {

	fmt.Println("-------------------encode start")
	actions := transaction.Actions{}

	amountArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Amount), 0)}

	toAddr := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.To), 0)}
	memoArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Memo), 0)}

	symbolArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Symbol), 0)}

	blkNumber := make([]byte, 8)
	binary.LittleEndian.PutUint64(blkNumber, b.BlkNumber)
	blkNumber_arg := para_paser.Arg{para_paser.TypeI64, blkNumber}

	expiry := make([]byte, 4)
	binary.LittleEndian.PutUint32(expiry, b.Expiry)
	expiry_arg := para_paser.Arg{para_paser.TypeI32, expiry}

	wp := &para_paser.WasmPara{
		FuncName: "transfer",
		Args:     append([]para_paser.Arg{}, toAddr, amountArg, symbolArg, memoArg, blkNumber_arg, expiry_arg),
	}
	paraBytes, _ := json.Marshal(wp)
	fmt.Println("--contract para")
	fmt.Println(string(paraBytes))
	fmt.Println("")
	action := transaction.Action{b.ConAddr, paraBytes}

	if b.TxFee > 0 {
		datafee := make([]byte, 8)
		binary.LittleEndian.PutUint64(datafee, b.TxFee)
		feeArg := para_paser.Arg{para_paser.TypeI64, datafee}
		wp = &para_paser.WasmPara{
			FuncName: "transferFee",
			Args:     append([]para_paser.Arg{}, feeArg),
		}
		paraBytes, _ = json.Marshal(wp)
		fmt.Println("--contract para")
		fmt.Println(string(paraBytes))
		fmt.Println("")
		actionFee := transaction.Action{types.HexToAddress(BcContract), paraBytes}
		actions = append(actions, &actionFee, &action)
	} else {
		actions = append(actions, &action)
	}

	s := transaction.NewMSigner(big.NewInt(1))
	tx := transaction.NewTransaction(b.Nc, actions)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("--Transaction for signature")
	fmt.Println(tx)
	fmt.Println("")
	txSign, err := SignTx(tx, s, b.Key)
	assert.AsserErr(err)
	return txSign
}

func SignTx(tx *transaction.Transaction, s transaction.MSigner, prv *ecdsa.PrivateKey) (*transaction.Transaction, error) {
	h := Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}

	fmt.Println("--Sign")
	fmt.Printf("%x\n",sig)
	fmt.Println("")
	return tx.WithSignature(s, sig)
}

func Hash(tx *transaction.Transaction) types.Hash {
	a := []interface{}{
		tx.Data.H,
		tx.Actions(),
		types.BigInt{*big.NewInt(1)},
		uint(0),
		uint(0),
	}
	fmt.Println("--hash data")
	fmt.Println(a)
	h, err := MsgpHash([]interface{}{
		tx.Data.H,
		tx.Actions(),
		types.BigInt{*big.NewInt(1)},
		uint(0),
		uint(0),
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("--hash")
	fmt.Printf("%x\n",h)
	fmt.Println("")

	return h
}

func MsgpHash(x interface{}) (h types.Hash, err error) {
	defer func() {
		panic := recover()
		if panic != nil {
			err = fmt.Errorf("%v", panic)
		}
	}()

	hw := sha3.NewKeccak256()

	xx := bytes.Buffer{}
	wr := msgp.NewWriter(&xx)
	err = wr.WriteIntf(x)
	if err != nil {
		return
	}

	err = wr.Flush()
	if err != nil {
		return
	}

	fmt.Println("")
	fmt.Println("--msgp for hash")
	fmt.Printf("%x\n",xx)
	fmt.Println(xx.String())
	fmt.Println("")

	hw.Write(xx.Bytes())

	hw.Sum(h[:0])
	return h, nil
}

func TestAxxxx(t *testing.T)  {
	v := big.NewInt(65533)
	neg := v.Sign()
	b := make([]byte, 1+len(v.Bytes()))
	b[0] = byte(neg)
	copy(b[1:], v.Bytes())
	fmt.Println(v.Bytes())

	aa := types.FromHex("81a44461746185a1488")
	fmt.Println(string(aa))
}

func PubkeyToAddress11(p ecdsa.PublicKey) types.Address {
	pubBytes := crypto.FromECDSAPub(&p)
	fmt.Printf("%x\n", pubBytes)
	return types.BytesToAddress(crypto.Keccak256(pubBytes[1:])[12:])
}
func TestYyyy(t *testing.T)  {
	key, err := crypto.HexToECDSA("aa63add32b57889a0dfe0f5c5f2386f148a7664ce8fdb59b739fde67c1cc2e52")
	checkErr(err)
	fmt.Println(key)

	addr := PubkeyToAddress11(key.PublicKey)
	fmt.Println(addr.HexLower())
}


func TestEncode11(t *testing.T)  {
	key, err := crypto.HexToECDSA("88d9e723ee0879315822c57dd0db655c02e62f4f2a2ea2be6acbcd9af2778981")
	checkErr(err)

	bt := Big_token_create{
		Key: key,
		ConAddr:types.HexToAddress("0x26ea1c5a38bb48bd58e62a5fa3ff06c9e328855e"),

		Nc:0,
		TxFee:uint64(0),

		Symbol:"BT",
		Name   :"BC-BT",
		Decimals:"18",
		Supply:"100000000",
		IsIssue    :0,

		BlkNumber: 1,
		Expiry: 100,
	}
	para , _ := json.MarshalIndent(&bt,"", "	")
	fmt.Println("-------------------input para start")
	fmt.Println("key", "49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	fmt.Println(string(para))
	fmt.Println("")
	xx := bt.MakeTransaction11()
	fmt.Println("--tx with signature")
	fmt.Println(xx)
	fmt.Println("")

	SendRawTransactionTest(xx)

}

func (b *Big_token_create) MakeTransaction11() *transaction.Transaction {
	actions := transaction.Actions{}

	NameArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Name), 0)}
	DecimalsArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Decimals), 0)}
	SupplyArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Supply), 0)}

	symbolArg := para_paser.Arg{para_paser.TypeAddress, append([]byte(b.Symbol), 0)}

	blkNumber := make([]byte, 8)
	binary.LittleEndian.PutUint64(blkNumber, b.BlkNumber)
	blkNumber_arg := para_paser.Arg{para_paser.TypeI64, blkNumber}

	expiry := make([]byte, 4)
	binary.LittleEndian.PutUint32(expiry, b.Expiry)
	expiry_arg := para_paser.Arg{para_paser.TypeI32, expiry}

	issule := make([]byte, 4)
	binary.LittleEndian.PutUint32(issule, b.IsIssue)
	issule_arg := para_paser.Arg{para_paser.TypeI32, issule}

	wp := &para_paser.WasmPara {
		FuncName: "create",
		Args:     append([]para_paser.Arg{}, symbolArg, NameArg, DecimalsArg, SupplyArg, issule_arg, blkNumber_arg, expiry_arg),
	}
	paraBytes, _ := json.Marshal(wp)
	action := transaction.Action{b.ConAddr, paraBytes}

	if b.TxFee > 0 {
		datafee := make([]byte, 8)
		binary.LittleEndian.PutUint64(datafee, b.TxFee)
		feeArg := para_paser.Arg{para_paser.TypeI64, datafee}
		wp = &para_paser.WasmPara {
			FuncName: "transferFee",
			Args:     append([]para_paser.Arg{}, feeArg),
		}
		paraBytes, _ = json.Marshal(wp)
		actionFee := transaction.Action{types.HexToAddress(BcContract), paraBytes}
		actions = append(actions, &actionFee, &action)
	} else {
		actions = append(actions, &action)
	}

	s := transaction.NewMSigner(big.NewInt(1))
	tx := transaction.NewTransaction(b.Nc, actions)
	txSign, err := SignTx(tx, s, b.Key)
	assert.AsserErr(err)
	return txSign
}

func TestX2121(t *testing.T)  {
	a:= types.FromHex("0x7b2266756e635f6e616d65223a227472616e73666572222c2261726773223a5b7b2274797065223a2261646472657373222c2276616c223a224d4868694d7a63794e6d55775a6d566b4d6d49335a574668595459774f5749334e32566a4e5441335a6a49325a44566d4f475a6c4e54597841413d3d227d2c7b2274797065223a2261646472657373222c2276616c223a224d5449794e41413d227d2c7b2274797065223a2261646472657373222c2276616c223a2256564e455641413d227d2c7b2274797065223a2261646472657373222c2276616c223a22564649784f5441314d4445784e5441774d6a4930524441334e7a5935526b4d344f55497a4d7a6c4741413d3d227d2c7b2274797065223a22696e743634222c2276616c223a226d616f42414141414141413d227d2c7b2274797065223a22696e743332222c2276616c223a225a41414141413d3d227d5d7d")
	fmt.Println(string(a))
	wp := &para_paser.WasmPara{}
	_ = json.Unmarshal(a, wp)
	//fmt.Println(wp,err)
	//fmt.Println(wp.Args)
	for _,arg := range wp.Args {
		//fmt.Println(arg.Type, arg.Data)
		if arg.Type=="address" {
			fmt.Println(arg.Type, string(arg.Data))
		} else if arg.Type=="int64" {
			fmt.Println(arg.Type, binary.LittleEndian.Uint64(arg.Data))
		}else if arg.Type=="int32" {
			fmt.Println(arg.Type, binary.LittleEndian.Uint32(arg.Data))
		}


	}
}

func TestX21212(t *testing.T)  {
	a:= types.FromHex("0x81a44461746185a14881a54e6f6e6365cc9aa4416374739182a8436f6e7472616374c41426ea1c5a38bb48bd58e62a5fa3ff06c9e328855ea6506172616d73c501617b2266756e635f6e616d65223a227472616e73666572222c2261726773223a5b7b2274797065223a2261646472657373222c2276616c223a224d4868694d7a63794e6d55775a6d566b4d6d49335a574668595459774f5749334e32566a4e5441335a6a49325a44566d4f475a6c4e54597841413d3d227d2c7b2274797065223a2261646472657373222c2276616c223a224d6a41774d4441774d4441774d4441774d4441774d4441774d41413d227d2c7b2274797065223a2261646472657373222c2276616c223a22516c5141227d2c7b2274797065223a2261646472657373222c2276616c223a22564649784f5441304d7a41784e4455344d446444525446435154457a52455646524455784f55457941413d3d227d2c7b2274797065223a22696e743634222c2276616c223a2257374142414141414141413d227d2c7b2274797065223a22696e743332222c2276616c223a225a41414141413d3d227d5d7da15681a6626967696e74c4020126a15281a6626967696e74c42101eb274e266acca532e3c457c552d7fda856f5d79c8d1ab7b1bf26b40e4a7b7156a15381a6626967696e74c421013417bb757cb98760c9d0e94d39fc0bc743e47bb3b37092ba2df21751faa20761")
	//fmt.Println(string(a))
	tr := &transaction.Transaction{}
	fmt.Println(string(a))
	byteBuf := bytes.NewBuffer(a)
	err := msgp.Decode(byteBuf, tr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tr)
}

func TestBig(t *testing.T)  {
	a := big.NewInt(162147)
	fmt.Println(a.Bytes())
}