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
// @File: block.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package block

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/tinylib/msgp/msgp"
	"math/big"
	"bchain.io/common"
	"bchain.io/common/types"
	"bchain.io/core/transaction"
	"bchain.io/trie"
	"bchain.io/utils/bloom"
	"sort"
	"sync/atomic"
	"time"
)

//go:generate msgp
//msgp:ignore Blocks
//go:generate gencodec -type Header -field-override headerMarshaling -out gen_header_json.go
//go:generate gencodec -type ConsensusData -field-override consensusDataMarshaling -out gen_consensus_json.go
//go:generate gofmt -w -s gen_header_json.go gen_consensus_json.go

type DerivableList interface {
	Len() int
	GetMsgp(i int) []byte
}

func DeriveSha(list DerivableList) types.Hash {
	keyBytesBuf := bytes.NewBuffer([]byte{})
	trie := new(trie.Trie)
	for i := 0; i < list.Len(); i++ {
		keyBytesBuf.Reset()
		binary.Write(keyBytesBuf, binary.BigEndian, i)
		trie.Update(keyBytesBuf.Bytes(), list.GetMsgp(i))
	}
	return trie.Hash()
}

type ConsensusData struct {
	Id   string `json:"consensus_id"         gencodec:"required"`
	Para []byte `json:"consensus_param"      gencodec:"required"`
}

func (cd ConsensusData) String() string {
	return cd.Id + fmt.Sprintf(" 0x%x", cd.Para)
}

type consensusDataMarshaling struct {
	Para types.BytesForJson
}

// block header
type Header struct {
	ParentHash      types.Hash    `json:"parentHash"         gencodec:"required"`
	StateRootHash   types.Hash    `json:"stateRoot"          gencodec:"required"`
	TxRootHash      types.Hash    `json:"transactionsRoot"   gencodec:"required"`
	ReceiptRootHash types.Hash    `json:"receiptsRoot"       gencodec:"required"`
	Bloom           types.Bloom   `json:"logsBloom"          gencodec:"required"`
	Number          *types.BigInt `json:"number"             gencodec:"required"`
	Time            *types.BigInt `json:"timestamp"          gencodec:"required"`
	Cdata           ConsensusData `json:"consensusData"      gencodec:"required"`
	Extra           []byte        `json:"extraData"          gencodec:"required"`

	//Signature values
	V *types.BigInt `json:"v"`
	R *types.BigInt `json:"r"`
	S *types.BigInt `json:"s"`

	//BlockProducer is not used in protocol, get from the signature(v,r,s)
	Producer types.Address `json:"blockProducer"      gencodec:"required"   msg:"-"`
}

// field type overrides for gencodec
type headerMarshaling struct {
	Extra types.BytesForJson
	Hash  types.Hash `json:"hash"` // adds call to Hash() in MarshalJSON
}

type Body struct {
	Transactions []*transaction.Transaction
}

type Block struct {
	H *Header // block header
	B Body    // all transactions in this block

	// caches
	hash atomic.Value
	size atomic.Value

	// These fields are used by package bchain to track
	// inter-peer block relay.
	ReceivedAt   time.Time   `msg:"-"`
	ReceivedFrom interface{} `msg:"-"`
}

func (h *Header) Hash() types.Hash {
	hash, err := common.MsgpHash(h)
	if err != nil {
		return types.Hash{}
	}
	return hash
}

var (
	EmptyRootHash = DeriveSha(transaction.Transactions{})
)

// NewBlock creates a new block. The input data is copied,
// changes to header and to the field values will not affect the
// block.
//
// The values of TxHash, ReceiptHash and Bloom in header
// are ignored and set to values derived from the given txs and receipts.
func NewBlock(header *Header, txs []*transaction.Transaction, receipts []*transaction.Receipt) *Block {
	b := &Block{H: CopyHeader(header)}

	// TODO: panic if len(txs) != len(receipts)
	if len(txs) == 0 {
		b.H.TxRootHash = EmptyRootHash
	} else {
		b.H.TxRootHash = DeriveSha(transaction.Transactions(txs))
		b.B.Transactions = make(transaction.Transactions, len(txs))
		copy(b.B.Transactions, txs)
	}

	if len(receipts) == 0 {
		b.H.ReceiptRootHash = EmptyRootHash
	} else {
		b.H.ReceiptRootHash = DeriveSha(transaction.Receipts(receipts))
		bloomIn := []bloom.BloomByte{}
		for _, receipt := range receipts {
			for _, log := range receipt.Logs {
				bloomIn = append(bloomIn, log.Address)
				for _, b := range log.Topics {
					bloomIn = append(bloomIn, b)
				}
			}
		}
		b.H.Bloom = bloom.CreateBloom(bloomIn)
	}
	return b
}

// CopyHeader creates a deep copy of a block header to prevent side effects from
// modifying a header variable.
func CopyHeader(h *Header) *Header {
	cpy := *h
	if cpy.Time = new(types.BigInt); h.Time != nil {
		cpy.Time.Put(h.Time.IntVal)
	}

	if cpy.Number = new(types.BigInt); h.Number != nil {
		cpy.Number.Put(h.Number.IntVal)
	}

	if cpy.V = new(types.BigInt); h.V != nil {
		cpy.V.Put(h.V.IntVal)
	}

	if cpy.R = new(types.BigInt); h.R != nil {
		cpy.R.Put(h.R.IntVal)
	}

	if cpy.S = new(types.BigInt); h.S != nil {
		cpy.S.Put(h.S.IntVal)
	}

	return &cpy
}

func NewBlockWithHeader(header *Header) *Block {
	return &Block{H: CopyHeader(header)}
}

func (header *Header) HashNoSig() types.Hash {
	v := &HeaderNoSig{
		header.ParentHash,
		header.StateRootHash,
		header.TxRootHash,
		header.ReceiptRootHash,
		header.Bloom,
		header.Number,
		header.Time,
		header.Cdata,
		header.Extra,
		header.Producer,
	}
	return v.Hash()
}

// WithSignature returns a new header with the given signature.
func (h *Header) WithSignature(signer Signer, sig []byte) (*Header, error) {
	r, s, v, err := signer.SignatureValues(h, sig)
	if err != nil {
		return nil, err
	}
	cpy := CopyHeader(h)
	cpy.R, cpy.S, cpy.V = &types.BigInt{*r}, &types.BigInt{*s}, &types.BigInt{*v}
	return cpy, nil
}

// AddSignature returns modify the  header( R, S, V) with the given signature.
func (h *Header) AddSignature(signer Signer, sig []byte) error {
	r, s, v, err := signer.SignatureValues(h, sig)
	if err != nil {
		return err
	}
	h.R, h.S, h.V = &types.BigInt{*r}, &types.BigInt{*s}, &types.BigInt{*v}
	return nil
}

func (b *Block) Transactions() transaction.Transactions { return b.B.Transactions }

func (b *Block) Transaction(hash types.Hash) *transaction.Transaction {
	for _, transaction := range b.B.Transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}
	return nil
}

func (b *Block) Number() *big.Int { return new(big.Int).Set(&b.H.Number.IntVal) }
func (b *Block) Time() *big.Int   { return new(big.Int).Set(&b.H.Time.IntVal) }

func (b *Block) NumberU64() uint64       { return b.H.Number.IntVal.Uint64() }
func (b *Block) Bloom() types.Bloom      { return b.H.Bloom }
func (b *Block) Producer() types.Address { return b.H.Producer }
func (b *Block) Root() types.Hash        { return b.H.StateRootHash }
func (b *Block) ParentHash() types.Hash  { return b.H.ParentHash }
func (b *Block) TxHash() types.Hash      { return b.H.TxRootHash }
func (b *Block) ReceiptHash() types.Hash { return b.H.ReceiptRootHash }

func (b *Block) Header() *Header { return CopyHeader(b.H) }

// Body returns the non-header content of the block.
func (b *Block) Body() *Body { return &Body{b.B.Transactions} }

func (b *Block) HashNoSig() types.Hash {
	return b.H.HashNoSig()
}

func (b *Block) Size() common.StorageSize {
	if size := b.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := 0
	var buf bytes.Buffer
	msgp.Encode(&buf, b)
	c = buf.Len()
	b.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// WithSeal returns a new block with the data from b but the header replaced with
// the sealed one.
func (b *Block) WithSeal(header *Header) *Block {
	cpy := *header

	return &Block{
		H: &cpy,
		B: b.B,
	}
}

// WithBody returns a new block with the given transaction  contents.
func (b *Block) WithBody(body *Body) *Block {
	block := &Block{
		H: CopyHeader(b.H),
	}
	block.B.Transactions = make([]*transaction.Transaction, len(body.Transactions))
	copy(block.B.Transactions, body.Transactions)
	return block
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() types.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(types.Hash)
	}
	v := b.H.Hash()
	b.hash.Store(v)
	return v
}

func (b *Block) String() string {
	str := fmt.Sprintf(`Block(#%v): Size: %v {
ProducerHash: %x
%v
Tx-Cnt:%d

Transactions:
%v
}
`, b.Number(), b.Size(), b.H.HashNoSig(), b.H,len(b.B.Transactions) , b.B.Transactions)
	return str
}

func (h *Header) String() string {
	return fmt.Sprintf(`Header(%x):
[
	ParentHash:         %x
	BlockProducer:      %x
	StateRootHash:      %x
	TxRootHash          %x
	ReceiptRootHash:    %x
	Bloom:              %x
	Number:	            %v
	Time:               %v
	ConsensusData:      %v
	ExtraData:          %s
	R:                  %v
	S:                  %v
	V:                  %v
]`, h.Hash(), h.ParentHash, h.Producer, h.StateRootHash, h.TxRootHash, h.ReceiptRootHash, h.Bloom, h.Number, h.Time, h.Cdata, h.Extra, h.R, h.S, h.V)
}

//blocks part

type Blocks []*Block

type BlockBy func(b1, b2 *Block) bool

func (self BlockBy) Sort(blocks Blocks) {
	bs := blockSorter{
		blocks: blocks,
		by:     self,
	}
	sort.Sort(bs)
}

type blockSorter struct {
	blocks Blocks
	by     func(b1, b2 *Block) bool
}

func (self blockSorter) Len() int { return len(self.blocks) }
func (self blockSorter) Swap(i, j int) {
	self.blocks[i], self.blocks[j] = self.blocks[j], self.blocks[i]
}
func (self blockSorter) Less(i, j int) bool { return self.by(self.blocks[i], self.blocks[j]) }

func Number(b1, b2 *Block) bool { return b1.H.Number.IntVal.Cmp(&b2.H.Number.IntVal) < 0 }

// header wihtout signature
type HeaderNoSig struct {
	ParentHash      types.Hash
	StateRootHash   types.Hash
	TxRootHash      types.Hash
	ReceiptRootHash types.Hash
	Bloom           types.Bloom
	Number          *types.BigInt
	Time            *types.BigInt
	Cdata    		ConsensusData
	Extra           []byte

	//BlockProducer is not used in protocol.
	Producer types.Address `msg:"-" `
}

func (h *HeaderNoSig) Hash() types.Hash {
	hash, err := common.MsgpHash(h)
	if err != nil {
		return types.Hash{}
	}
	return hash
}

//type Headers []*Header
type Headers struct {
	Headers []*Header
}
