package transaction

import (
	"testing"

	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/tinylib/msgp/msgp"
	"bchain.io/common"
	"bchain.io/common/types"
	"bchain.io/utils/crypto"
	"reflect"
	"encoding/binary"
	"bchain.io/core/interpreter/wasmre/para_paser"
)

var (
	testKey, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddress = crypto.PubkeyToAddress(testKey.PublicKey)
	mSigner     = NewMSigner(common.Big1)

	emptyTx       = NewTransaction(0, Actions{})
	rightvrsTx, _ = SignTx(emptyTx, mSigner, testKey)
)

func TestTransactionPrint(t *testing.T) {
	var nonce uint64 = 10

	actions := Actions{{
		Contract: types.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed"),
		Params:   []byte{1, 2, 3, 4},
	}, {
		Contract: types.HexToAddress("0x888eb6053f3e94c9b9a09f33669435e7ef1beaed"),
		Params:   []byte{5, 6, 7, 8},
	}}

	tx := newTransaction(nonce, actions)
	tx, _ = SignTx(tx, mSigner, testKey)

	fmt.Printf("trx: %s\n", tx)
}

func TestGetTxPriority(t *testing.T) {
	actions := Actions{{
		Contract: types.HexToAddress("0x192d52D8cE0c7bBAf0780EAb04860D6Ba012578B"),
		Params:   []byte("balanceFee(10)"),
	}}
	tx := newTransaction(0, actions)
	tx.ParsePriority()
}

func TestGetTxPriorityWasm(t *testing.T) {
	data1 := make([]byte, 8)
	binary.LittleEndian.PutUint64(data1, 100)
	amount := para_paser.Arg{para_paser.TypeI64, data1}

	wp := &para_paser.WasmPara{
		FuncName: "transferFee",
		Args:     append([]para_paser.Arg{}, amount),
	}
	paraBytes, _ := json.Marshal(wp)
	fmt.Println(string(paraBytes))
	actions := Actions{{
		Contract: types.HexToAddress("0xb78f12Cb3924607A8BC6a66799e159E3459097e9"),
		Params:   paraBytes,
	}}
	tx := newTransaction(0, actions)
	tx.ParsePriority()
	fmt.Println(tx.priority.Load())
}

func TestTransactionNew(t *testing.T) {
	var nonce uint64 = 10
	data := []byte{}
	data = append(data, 1, 4, 5)
	actions := Actions{{
		Contract: testAddress,
		Params:   data,
	}}

	tx := newTransaction(nonce, actions)

	sig := NewMSigner(common.Big1)

	//h := sig.Hash(tx)
	//t.Logf("transaction hash = %x", h)
	//if !reflect.DeepEqual(h, testAddress) {
	//	t.Errorf("Error: have hash: %x, want: %v", h, testAddress.Hex())
	//}

	txWithSig, err := SignTx(tx, sig, testKey)
	if err != nil {
		t.Errorf("SignTx error: %v", err)
	}
	txWithSig.PrintVSR()

	addr, err := sig.Sender(txWithSig)
	if err != nil {
		t.Errorf("Sender error: %v", err)
	}
	if !reflect.DeepEqual(addr, testAddress) {
		t.Errorf("Error: get addr: %v want addr: %v", addr, testAddress)
	}
}

func TestAsMessageGenerate(t *testing.T) {
	var nonce uint64 = 10
	data := []byte{}
	data = append(data, 1, 4, 5)
	address := types.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed")
	actions := Actions{{
		Contract: address,
		Params:   data,
	}}
	//new transaction
	tx := newTransaction(nonce, actions)
	//create key
	key, _ := crypto.GenerateKey()
	//Sign tx
	txSigned, _ := SignTx(tx, mSigner, key)
	_ = txSigned

	t.Skip("TODO AsMessage() ...")

	/*msg , err := txSigned.AsMessage(mSigner)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("msg:" , msg)*/
}

func TestTransactionSigHash(t *testing.T) {
	msig := NewMSigner(common.Big1)

	if msig.Hash(emptyTx) != types.HexToHash("9843af805c36dea5141abfa884a0e34c0946573910abf85dffddfe8b0c8a51fd") {
		t.Errorf("empty transaction hash mismatch, got %x", emptyTx.Hash())
	}

	if msig.Hash(rightvrsTx) != types.HexToHash("9843af805c36dea5141abfa884a0e34c0946573910abf85dffddfe8b0c8a51fd") {
		t.Errorf("RightVRS transaction hash mismatch, got %x", rightvrsTx.Hash())
	}
}

func TestTransactionEncode(t *testing.T) {
	txb := bytes.Buffer{}
	err := msgp.Encode(&txb, rightvrsTx)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	should := types.FromHex("81a44461746185a14881a54e6f6e636500a44163747390a15681a6626967696e74c4020125a15281a6626967696e74c42101408f967caedb89eca32cf1bda6610de633f8b070deeb8e01e8211d51b6764e99a15381a6626967696e74c4210106cfa16228a02a381ba1a930b4244245c8b9e600dc49c483752492a7cc2b0940")
	if !bytes.Equal(txb.Bytes(), should) {
		t.Errorf("encoded RLP mismatch, got %x", txb.Bytes())
	}
}

func decodeTx(data []byte) (*Transaction, error) {
	var tx Transaction
	err := msgp.Decode(bytes.NewBuffer(data), &tx)
	return &tx, err
}

func defaultTestKey() (*ecdsa.PrivateKey, types.Address) {
	return testKey, testAddress
}

func TestSender(t *testing.T) {
	_, addr := defaultTestKey()
	tx, err := decodeTx(types.Hex2Bytes("81a44461746185a14881a54e6f6e636500a44163747390a15681a6626967696e74c4020125a15281a6626967696e74c42101408f967caedb89eca32cf1bda6610de633f8b070deeb8e01e8211d51b6764e99a15381a6626967696e74c4210106cfa16228a02a381ba1a930b4244245c8b9e600dc49c483752492a7cc2b0940"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	from, err := Sender(NewMSigner(common.Big1), tx)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if addr != from {
		t.Errorf("derived address doesn't match")
	}
}

// TestTransactionJSON tests serializing/de-serializing to/from JSON.
func TestTransactionJSON(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("could not generate key: %v", err)
	}
	signer := mSigner

	for i := uint64(0); i < 25; i++ {
		var tx *Transaction
		tx = NewTransaction(i, Actions{})

		tx, err := SignTx(tx, signer, key)
		if err != nil {
			t.Fatalf("could not sign transaction: %v", err)
		}

		data, err := json.Marshal(tx)
		if err != nil {
			t.Errorf("json.Marshal failed: %v", err)
		}

		var parsedTx *Transaction
		if err := json.Unmarshal(data, &parsedTx); err != nil {
			t.Errorf("json.Unmarshal failed: %v", err)
		}

		// compare tx, parsedTx
		if tx.Hash() != parsedTx.Hash() {
			t.Errorf("parsed tx differs from original tx, want %v, got %v", tx, parsedTx)
		}
		if tx.ChainId().Cmp(parsedTx.ChainId()) != 0 {
			t.Errorf("invalid chain id, want %d, got %d", tx.ChainId(), parsedTx.ChainId())
		}
	}
}

func TestTransactionPrintJSON(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("could not generate key: %v", err)
	}
	signer := mSigner
	actions := Actions{{
		Contract: types.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed"),
		Params:   []byte{1, 2, 3, 4},
	}, {
		Contract: types.HexToAddress("0x888eb6053f3e94c9b9a09f33669435e7ef1beaed"),
		Params:   []byte{5, 6, 7, 8},
	}}

	tx := NewTransaction(0, actions)
	tx, _ = SignTx(tx, signer, key)

	data, err := json.Marshal(tx)
	fmt.Printf("%s\n", string(data))

	data, err = json.MarshalIndent(tx, "", "    ")
	fmt.Printf("%s\n", string(data))
}
