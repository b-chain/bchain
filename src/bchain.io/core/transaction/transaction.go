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
// @File: transaction.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package transaction

import (
	"errors"
	"math/big"
	"sync/atomic"

	"bytes"
	"container/heap"
	"fmt"
	"github.com/tinylib/msgp/msgp"
	"bchain.io/common"
	"bchain.io/common/types"
	"bchain.io/utils/crypto"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"encoding/json"
	"encoding/binary"
)

//go:generate msgp
//msgp:ignore Message TransactionsByPriorityAndNonce
//go:generate gencodec -type TxHeader -field-override txHeaderMarshaling -out gen_txheader_json.go
//go:generate gencodec -type Action -field-override actionMarshaling -out gen_action_json.go
//go:generate gencodec -type Txdata  -out gen_tx_json.go
//go:generate gofmt -w -s gen_txheader_json.go gen_action_json.go gen_tx_json.go

var (
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
	errNoSigner   = errors.New("missing signing methods")
)

// deriveSigner makes a *best* guess about which signer to use.
func deriveSigner(V *big.Int) Signer {
	return NewMSigner(deriveChainId(V))
}

type Transaction struct {
	Data Txdata

	// caches
	priority atomic.Value // big.Int
	hash     atomic.Value
	size     atomic.Value
	from     atomic.Value
}

//for test
func (this *Transaction) PrintVSR() {
	fmt.Printf("V:%v, S:%v, R:%v\n",
		this.Data.V.IntVal,
		this.Data.S.IntVal,
		this.Data.R.IntVal)
}

type Action struct {
	Contract types.Address `json:"contract" gencodec:"required"`
	Params   []byte        `json:"params"   gencodec:"required"`
}

type actionMarshaling struct {
	Params types.BytesForJson
}

func NewAction() *Action {
	return &Action{
		types.Address{},
		make([]byte, 0),
	}
}

//String just print action, simple is best
func (tx *Action) String() string {
	rStr := fmt.Sprintf(`
{
Contract:    %x
Params:      %s
}`,
		tx.Contract,
		types.ToHex(tx.Params),
	)

	return rStr
}

type Actions []*Action

type TxHeader struct {
	Nonce uint64 `json:"AccountNonce"   gencodec:"required"`
	//Expiration      *types.BigInt
	//Delay			*types.BigInt
}

type txHeaderMarshaling struct {
	Nonce types.Uint64ForJson
}

type Txdata struct {
	H    TxHeader `json:"header"  gencodec:"required"`
	Acts Actions  `json:"actions" gencodec:"required"`

	// Signature values
	V *types.BigInt `json:"v"       gencodec:"required"`
	R *types.BigInt `json:"r"       gencodec:"required"`
	S *types.BigInt `json:"s"       gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *types.Hash `json:"hash"    msg:"-"`
}

//All actions is made by interpreter
func NewTransaction(nonce uint64, actions Actions) *Transaction {
	return newTransaction(nonce, actions)
}

//the actions is right or not ,should be judged by interpreter,we have no right to do this
func newTransaction(nonce uint64, actions Actions) *Transaction {
	if len(actions) < 0 {
		return nil
	}

	d := Txdata{
		H: TxHeader{
			nonce,
		},
		Acts: actions,
		V:    new(types.BigInt),
		R:    new(types.BigInt),
		S:    new(types.BigInt),
	}

	return &Transaction{Data: d}
}

// ChainId returns which chain id this transaction was signed for (if at all)
func (tx *Transaction) ChainId() *big.Int {
	return deriveChainId(&tx.Data.V.IntVal)
}

// Protected returns whether the transaction is protected from replay protection.
func (tx *Transaction) Protected() bool {
	return isProtectedV(&tx.Data.V.IntVal)
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		//if v is 27 or 28,return false
		return v != 27 && v != 28
	}
	// anything not 27 or 28 are considered unprotected
	return true
}

// MarshalJSON encodes the web3 RPC transaction format.
func (tx *Transaction) MarshalJSON() ([]byte, error) {
	hash := tx.Hash()
	data := tx.Data
	data.Hash = &hash
	return data.MarshalJSON()
}

// UnmarshalJSON decodes the web3 RPC transaction format.
func (tx *Transaction) UnmarshalJSON(input []byte) error {
	var dec Txdata
	if err := dec.UnmarshalJSON(input); err != nil {
		return err
	}
	var V byte
	if isProtectedV(&dec.V.IntVal) {
		chainID := deriveChainId(&dec.V.IntVal).Uint64()
		V = byte(dec.V.IntVal.Uint64() - 35 - 2*chainID)
	} else {
		V = byte(dec.V.IntVal.Uint64() - 27)
	}
	if !crypto.ValidateSignatureValues(V, &dec.R.IntVal, &dec.S.IntVal, false) {
		return ErrInvalidSig
	}
	*tx = Transaction{Data: dec}
	return nil
}

func (tx *Transaction) Nonce() uint64    { return tx.Data.H.Nonce }
func (tx *Transaction) CheckNonce() bool { return true }
func (tx *Transaction) Actions() Actions { return tx.Data.Acts }

// Hash hashes the Msgp encoding of tx.
// It uniquely identifies the transaction.
func (tx *Transaction) Hash() types.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(types.Hash)
	}

	v, err := common.MsgpHash(tx)
	if err != nil {
		logger.Errorf("Transaction hash error: %v", err)
		return types.Hash{}
	}

	tx.hash.Store(v)
	return v
}

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}

// Size returns the true Msgp encoded storage size of the transaction, either by
// encoding and returning it, or returning a previsouly cached value.
func (tx *Transaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	var buf bytes.Buffer
	err := msgp.Encode(&buf, tx)
	if err != nil {
		c = writeCounter(0)
	} else {
		c = writeCounter(len(buf.Bytes()))
	}

	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// get priority from 1st action by interpreter
func (tx *Transaction) Priority() *big.Int {
	if priority := tx.priority.Load(); priority != nil {
		return priority.(*big.Int)
	}

	// TODO: get priority from 1st action by interpreter
	v := big.NewInt(0) // default v = 0
	tx.priority.Store(v)
	return v
}

func (tx *Transaction) SetPriority(newPriority *big.Int) {
	tx.priority.Store(newPriority)
}

func (tx *Transaction) ParsePriority() {
	setDefault := true
	defer func() {
		if setDefault {
			tx.SetPriority(big.NewInt(0))
		}
	}()

	feeContractAddress := types.HexToAddress("0xb78f12Cb3924607A8BC6a66799e159E3459097e9")
	if len(tx.Actions()) == 0 {
		return
	}
	act := tx.Actions()[0]
	if act.Contract != feeContractAddress {
		//logger.Error("feeContractAddress wrong")
		return
	}
	//check balanceFee string
	wp := &para_paser.WasmPara{}
	err := json.Unmarshal(act.Params, wp)
	if err != err {
		return
	}
	if wp.FuncName != "transferFee" {
		return
	}
	if len(wp.Args) == 0 {
		return
	}
	arg := wp.Args[0]
	var trFee uint64
	if len(arg.Data) == 4{
		trFee = uint64(binary.LittleEndian.Uint32(arg.Data))
	} else if len(arg.Data) == 8 {
		trFee = binary.LittleEndian.Uint64(arg.Data)
	} else {
		return
	}

	p := big.NewInt(int64(trFee))
	tx.SetPriority(p)

	setDefault = false
}

/*//In Bchain, all details of transaction dealing should not visiable for others except vm(interpreter)
func (tx *Transaction) AsMessage(s Signer) (Message, error) {
	newActions := Actions{}
	newActions = append(newActions , tx.Data.Acts...)
	msg := Message{
		nonce:      tx.Nonce(),
		actions:    newActions,
		checkNonce: true,
	}
	var err error
	msg.from, err = Sender(s, tx)
	return msg, err
}*/

// WithSignature returns a new transaction with the given signature.
// This signature needs to be formatted as described in the yellow paper (v+27).
func (tx *Transaction) WithSignature(signer Signer, sig []byte) (*Transaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	cpy := &Transaction{Data: tx.Data}
	cpy.Data.R, cpy.Data.S, cpy.Data.V = &types.BigInt{*r}, &types.BigInt{*s}, &types.BigInt{*v}
	return cpy, nil
}

func (tx *Transaction) RawSignatureValues() (*big.Int, *big.Int, *big.Int) {
	return &tx.Data.V.IntVal, &tx.Data.R.IntVal, &tx.Data.S.IntVal
}

//String just print transaction, simple is best
func (tx *Transaction) String() string {
	var from string
	if tx.Data.V != nil {
		signer := deriveSigner(&tx.Data.V.IntVal)
		if f, err := Sender(signer, tx); err != nil {
			from = "[invalid sender: invalid sig]"
		} else {
			from = fmt.Sprintf("0x%x", f[:])
		}
	} else {
		from = "[invalid sender: nil V field]"
	}

	rStr := fmt.Sprintf(`
TX(%x)
From:       %s
Nonce:      %v
ActionLen:  %d
Actions:    %s
V:          %v
S:          %v
R:          %v
`,
		tx.Hash(),
		from,
		tx.Nonce(),
		len(tx.Data.Acts),
		tx.Data.Acts,
		tx.Data.V,
		tx.Data.S,
		tx.Data.R,
	)

	return rStr
}

// Transactions is a Transaction slice type for basic sorting.
type Transactions []*Transaction

// Len returns the length of s.
func (s Transactions) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s Transactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetMsgp implements Msgpable and returns the i'th element of s in msgp.
func (s Transactions) GetMsgp(i int) []byte {
	var buf bytes.Buffer
	err := msgp.Encode(&buf, s[i])
	if err != nil {
		return nil
	}

	return buf.Bytes()
}

// TxDifference returns a new set t which is the difference between a to b.
func TxDifference(a, b Transactions) (keep Transactions) {
	keep = make(Transactions, 0, len(a))

	remove := make(map[types.Hash]struct{})
	for _, tx := range b {
		remove[tx.Hash()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[tx.Hash()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}

// TxByPriority implements the sort interface to allow sorting a list of transactions
// by their priority. TxByPriority implements both the sort and the heap interface,
// making it useful for all at once sorting as well as individually adding and removing elements.
type TxByPriority Transactions

func (s TxByPriority) Len() int           { return len(s) }
func (s TxByPriority) Less(i, j int) bool { return s[i].Priority().Cmp(s[j].Priority()) > 0 }
func (s TxByPriority) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s *TxByPriority) Push(x interface{}) {
	*s = append(*s, x.(*Transaction))
}

func (s *TxByPriority) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

// TransactionsByPriorityAndNonce represents a set of transactions that can return
// transactions in a priority sorted order, while supporting removing
// entire batches of transactions for non-executable accounts.
type TransactionsByPriorityAndNonce struct {
	txs    map[types.Address]Transactions // Per account nonce-sorted list of transactions
	heads  TxByPriority                   // Next transaction for each unique account (priority heap)
	signer Signer                         // Signer for the set of transactions
}

func (t *TransactionsByPriorityAndNonce) Length() int {
	txCnt := 0
	for _, txs := range t.txs {
		txCnt += len(txs)
	}

	return txCnt + len(t.heads)
}

// NewTransactionsByPriorityAndNonce creates a transaction set that can retrieve
// priority sorted transactions in a nonce-honouring way.
//
// Note, the input map is reowned so the caller should not interact any more with
// if after providing it to the constructor.
func NewTransactionsByPriorityAndNonce(signer Signer, txs map[types.Address]Transactions, txsReward *Transaction) *TransactionsByPriorityAndNonce {
	// Initialize a priority based heap with the head transactions
	heads := make(TxByPriority, 0, len(txs))
	for _, accTxs := range txs {
		heads = append(heads, accTxs[0])
		// Ensure the sender address is from the signer
		acc, _ := Sender(signer, accTxs[0])
		txs[acc] = accTxs[1:]
	}
	heap.Init(&heads)

	newHeads := make(TxByPriority, 0, len(txs)+1)
	if txsReward != nil {
		newHeads = append(newHeads, txsReward)
	}
	newHeads = append(newHeads, heads...)
	// Assemble and return the transaction set
	return &TransactionsByPriorityAndNonce{
		txs:    txs,
		heads:  newHeads,
		signer: signer,
	}
}

// Peek returns the next transaction by priority.
func (t *TransactionsByPriorityAndNonce) Peek() *Transaction {
	if len(t.heads) == 0 {
		return nil
	}
	return t.heads[0]
}

// Shift replaces the current best head with the next one from the same account.
func (t *TransactionsByPriorityAndNonce) Shift() {
	acc, _ := Sender(t.signer, t.heads[0])
	if txs, ok := t.txs[acc]; ok && len(txs) > 0 {
		t.heads[0], t.txs[acc] = txs[0], txs[1:]
		heap.Fix(&t.heads, 0)
	} else {
		heap.Pop(&t.heads)
	}
}

// Pop removes the best transaction, *not* replacing it with the next one from
// the same account. This should be used when a transaction cannot be executed
// and hence all subsequent ones should be discarded from the same account.
func (t *TransactionsByPriorityAndNonce) Pop() {
	heap.Pop(&t.heads)
}

// TxByNonce implements the sort interface to allow sorting a list of transactions
// by their nonces. This is usually only useful for sorting transactions from a
// single account, otherwise a nonce comparison doesn't make much sense.
type TxByNonce Transactions

func (s TxByNonce) Len() int           { return len(s) }
func (s TxByNonce) Less(i, j int) bool { return s[i].Nonce() < s[j].Nonce() }
func (s TxByNonce) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Message is a fully derived transaction and implements core.Message
//
// NOTE: In a future PR this will be removed.
type Message struct {
	from       types.Address
	nonce      uint64
	actions    Actions
	checkNonce bool
}

func NewMessage(from types.Address, nonce uint64, actions Actions, checkNonce bool) Message {
	return Message{
		from:       from,
		nonce:      nonce,
		actions:    actions,
		checkNonce: checkNonce,
	}
}

func (m Message) From() types.Address { return m.from }
func (m Message) Nonce() uint64       { return m.nonce }
func (m Message) Actions() Actions    { return m.actions }
func (m Message) CheckNonce() bool    { return m.checkNonce }
