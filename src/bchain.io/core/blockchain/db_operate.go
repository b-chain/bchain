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
// @File: db_operate.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package blockchain

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	msgp "github.com/tinylib/msgp/msgp"
	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/transaction"
	"bchain.io/params"
	"bchain.io/utils/database"
	"bchain.io/utils/metrics"
	"math/big"
)

// DatabaseReader wraps the Get method of a backing data store.
type DatabaseReader interface {
	Get(key []byte) (value []byte, err error)
}

// DatabaseDeleter wraps the Delete method of a backing data store.
type DatabaseDeleter interface {
	Delete(key []byte) error
}

//go:generate msgp
var (
	headHeaderKey = []byte("LastHeader")
	headBlockKey  = []byte("LastBlock")
	headFastKey   = []byte("LastFast")

	// Data item prefixes (use single byte to avoid mixing data types, avoid `i`).
	headerPrefix        = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
	numSuffix           = []byte("n") // headerPrefix + num (uint64 big endian) + numSuffix -> hash
	blockHashPrefix     = []byte("H") // blockHashPrefix + hash -> num (uint64 big endian)
	bodyPrefix          = []byte("b") // bodyPrefix + num (uint64 big endian) + hash -> block body
	blockReceiptsPrefix = []byte("r") // blockReceiptsPrefix + num (uint64 big endian) + hash -> block receipts
	lookupPrefix        = []byte("l") // lookupPrefix + hash -> transaction/receipt lookup metadata
	trAddrNoncePrefix   = []byte("a") // trAddrNoncePrefix + addr + nonce -> transaction/receipt lookup metadata
	bloomBitsPrefix     = []byte("B") // bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash -> bloom bits
	consensusPrefix     = []byte("consensus-")
	totalTxPrefix       = []byte("t-tx")

	preimagePrefix = "secure-key-"          // preimagePrefix + hash -> preimage
	configPrefix   = []byte("bchain-config-") // config prefix for the db

	// Chain index prefixes (use `i` + single byte to avoid mixing data types).
	BloomBitsIndexPrefix = []byte("iB") // BloomBitsIndexPrefix is the data table of a chain indexer to track its progress

	// used by old db, now only used for conversion
	oldReceiptsPrefix = []byte("receipts-")
	oldTxMetaSuffix   = []byte{0x01}

	ErrChainConfigNotFound = errors.New("ChainConfig not found") // general config not found error

	preimageCounter    = metrics.NewRegisteredCounter("db/preimage/total", nil)
	preimageHitCounter = metrics.NewRegisteredCounter("db/preimage/hits", nil)
)

// TxLookupEntry is a positional metadata to help looking up the data content of
// a transaction or receipt given only its hash.
type TxLookupEntry struct {
	BlockHash  types.Hash
	BlockIndex uint64
	Index      uint64
}

// block statistics info
type BlockStat struct {
	Ttxs        *types.BigInt
	TsoNormal   *types.BigInt
	TsoContract *types.BigInt
	TstateNum   *types.BigInt
}

// encodeBlockNumber encodes a block number as big endian uint64
func encodeBlockNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// GetCanonicalHash retrieves a hash assigned to a canonical block number.
func GetCanonicalHash(db DatabaseReader, number uint64) types.Hash {
	data, _ := db.Get(append(append(headerPrefix, encodeBlockNumber(number)...), numSuffix...))
	if len(data) == 0 {
		return types.Hash{}
	}
	return types.BytesToHash(data)
}

const missingNumber = uint64(0xffffffffffffffff)

// GetBlockNumber returns the block number assigned to a block hash
// if the corresponding header is present in the database
func GetBlockNumber(db DatabaseReader, hash types.Hash) uint64 {
	data, _ := db.Get(append(blockHashPrefix, hash.Bytes()...))
	if len(data) != 8 {
		return missingNumber
	}
	return binary.BigEndian.Uint64(data)
}

// GetHeadHeaderHash retrieves the hash of the current canonical head block's
// header. The difference between this and GetHeadBlockHash is that whereas the
// last block hash is only updated upon a full block import, the last header
// hash is updated already at header import, allowing head tracking for the
// light synchronization mechanism.
func GetHeadHeaderHash(db DatabaseReader) types.Hash {
	data, _ := db.Get(headHeaderKey)
	if len(data) == 0 {
		return types.Hash{}
	}
	return types.BytesToHash(data)
}

// GetHeadBlockHash retrieves the hash of the current canonical head block.
func GetHeadBlockHash(db DatabaseReader) types.Hash {
	data, _ := db.Get(headBlockKey)
	if len(data) == 0 {
		return types.Hash{}
	}
	return types.BytesToHash(data)
}

// GetHeadFastBlockHash retrieves the hash of the current canonical head block during
// fast synchronization. The difference between this and GetHeadBlockHash is that
// whereas the last block hash is only updated upon a full block import, the last
// fast hash is updated when importing pre-processed blocks.
func GetHeadFastBlockHash(db DatabaseReader) types.Hash {
	data, _ := db.Get(headFastKey)
	if len(data) == 0 {
		return types.Hash{}
	}
	return types.BytesToHash(data)
}

// GetHeaderMsgp retrieves a block header in its raw MSGP database encoding, or nil
// if the header's not found.
func GetHeaderMsgp(db DatabaseReader, hash types.Hash, number uint64) []byte {
	data, _ := db.Get(headerKey(hash, number))
	return data
}

// GetHeader retrieves the block header corresponding to the hash, nil if none
// found.
func GetHeader(db DatabaseReader, hash types.Hash, number uint64) *block.Header {
	data := GetHeaderMsgp(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	header := new(block.Header)
	byteBuf := bytes.NewBuffer(data)
	err := msgp.Decode(byteBuf, header)
	if err != nil {
		logger.Error("header.Decode err", "hash", hash, "err", err)
		return nil
	}
	return header
}

// GetBodyMsgp retrieves the block body in Msgp encoding.
func GetBodyMsgp(db DatabaseReader, hash types.Hash, number uint64) []byte {
	data, _ := db.Get(blockBodyKey(hash, number))
	return data
}

func headerKey(hash types.Hash, number uint64) []byte {
	return append(append(headerPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
}

func blockBodyKey(hash types.Hash, number uint64) []byte {
	return append(append(bodyPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
}

// GetBody retrieves the block body corresponding to the
// hash, nil if none found.
func GetBody(db DatabaseReader, hash types.Hash, number uint64) *block.Body {
	data := GetBodyMsgp(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	body := new(block.Body)
	byteBuf := bytes.NewBuffer(data)
	err := msgp.Decode(byteBuf, body)
	if err != nil {
		logger.Error("Invalid block body msgp", "hash", hash, "err", err)
		return nil
	}
	return body
}

// GetBlock retrieves an entire block corresponding to the hash, assembling it
// back from the stored header and body. If either the header or body could not
// be retrieved nil is returned.
//
// Note, due to concurrent download of header and block body the header and thus
// canonical hash can be stored in the database but the body data not (yet).
func GetBlock(db DatabaseReader, hash types.Hash, number uint64) *block.Block {
	// Retrieve the block header and body contents
	header := GetHeader(db, hash, number)
	if header == nil {
		return nil
	}
	body := GetBody(db, hash, number)
	if body == nil {
		return nil
	}
	// Reassemble the block and return
	return block.NewBlockWithHeader(header).WithBody(body)
}

// GetBlockReceipts retrieves the receipts generated by the transactions included
// in a block given by its hash.
func GetBlockReceipts(db DatabaseReader, hash types.Hash, number uint64) transaction.Receipts {
	data, _ := db.Get(append(append(blockReceiptsPrefix, encodeBlockNumber(number)...), hash[:]...))
	if len(data) == 0 {
		return nil
	}
	var buf bytes.Buffer
	//buf.Write(data)
	left := data
	receipts := []*transaction.Receipt{}

	for {
		if len(left) == 0 {
			break
		}
		var r transaction.Receipt
		buf.Reset()
		buf.Write(left)
		err := msgp.Decode(&buf, &r)
		if err != nil {
			logger.Error("Invalid receipt array MSGP", "hash", hash, "err", err)
			return nil
		}
		receipts = append(receipts, &r)
		left, _ = msgp.Skip(left)
	}
	return receipts
}

// GetTxLookupEntry retrieves the positional metadata associated with a transaction
// hash to allow retrieving the transaction or receipt by hash.
func GetTxLookupEntry(db DatabaseReader, hash types.Hash) (types.Hash, uint64, uint64) {
	// Load the positional metadata from disk and bail if it fails
	data, _ := db.Get(append(lookupPrefix, hash.Bytes()...))
	if len(data) == 0 {
		return types.Hash{}, 0, 0
	}
	// Parse and return the contents of the lookup entry
	var entry TxLookupEntry
	byteBuf := bytes.NewBuffer(data)
	err := msgp.Decode(byteBuf, &entry)
	if err != nil {
		logger.Error("Invalid lookup entry msgp", "hash", hash, "err", err)
		return types.Hash{}, 0, 0
	}

	return entry.BlockHash, entry.BlockIndex, entry.Index
}

func GetAddrNonceTxLookupEntry(db DatabaseReader, sender types.Address, nonce uint64) (types.Hash, uint64, uint64) {
	// Load the positional metadata from disk and bail if it fails
	nonceByte := big.NewInt(int64(nonce))
	data, _ := db.Get(append(append(trAddrNoncePrefix, sender.Bytes()...), nonceByte.Bytes()...))
	if len(data) == 0 {
		return types.Hash{}, 0, 0
	}
	// Parse and return the contents of the lookup entry
	var entry TxLookupEntry
	byteBuf := bytes.NewBuffer(data)
	err := msgp.Decode(byteBuf, &entry)
	if err != nil {
		logger.Error("Invalid lookup entry msgp", "sender", sender.HexLower(), "err", err)
		return types.Hash{}, 0, 0
	}

	return entry.BlockHash, entry.BlockIndex, entry.Index
}

// GetTransaction retrieves a specific transaction from the database, along with
// its added positional metadata.
func GetTransaction(db DatabaseReader, hash types.Hash) (*transaction.Transaction, types.Hash, uint64, uint64) {
	// Retrieve the lookup metadata and resolve the transaction from the body
	blockHash, blockNumber, txIndex := GetTxLookupEntry(db, hash)

	if blockHash != (types.Hash{}) {
		body := GetBody(db, blockHash, blockNumber)
		if body == nil || len(body.Transactions) <= int(txIndex) {
			logger.Error("Transaction referenced missing", "number", blockNumber, "hash", blockHash, "index", txIndex)
			return nil, types.Hash{}, 0, 0
		}
		return body.Transactions[txIndex], blockHash, blockNumber, txIndex
	}
	// Old transaction representation, load the transaction and it's metadata separately
	data, _ := db.Get(hash.Bytes())
	if len(data) == 0 {
		return nil, types.Hash{}, 0, 0
	}
	var tx transaction.Transaction
	byteBuf := bytes.NewBuffer(data)
	if err := msgp.Decode(byteBuf, &tx); err != nil {
		return nil, types.Hash{}, 0, 0
	}
	// Retrieve the blockchain positional metadata
	data, _ = db.Get(append(hash.Bytes(), oldTxMetaSuffix...))
	if len(data) == 0 {
		return nil, types.Hash{}, 0, 0
	}
	var entry TxLookupEntry
	byteBuf.Reset()
	byteBuf.Write(data)
	if err := msgp.Decode(byteBuf, &entry); err != nil {
		return nil, types.Hash{}, 0, 0
	}
	return &tx, entry.BlockHash, entry.BlockIndex, entry.Index
}

func GetTransactionByAddress(db DatabaseReader, addr types.Address, nonce uint64) (*transaction.Transaction, types.Hash, uint64, uint64) {
	// Retrieve the lookup metadata and resolve the transaction from the body
	blockHash, blockNumber, txIndex := GetAddrNonceTxLookupEntry(db, addr, nonce)
	if blockHash != (types.Hash{}) {
		body := GetBody(db, blockHash, blockNumber)
		if body == nil || len(body.Transactions) <= int(txIndex) {
			logger.Error("Transaction referenced missing", "number", blockNumber, "hash", blockHash, "index", txIndex)
			return nil, types.Hash{}, 0, 0
		}
		return body.Transactions[txIndex], blockHash, blockNumber, txIndex
	}
	return nil, types.Hash{}, 0, 0
}

// GetReceipt retrieves a specific transaction receipt from the database, along with
// its added positional metadata.
func GetReceipt(db DatabaseReader, hash types.Hash) (*transaction.Receipt, types.Hash, uint64, uint64) {
	// Retrieve the lookup metadata and resolve the receipt from the receipts
	blockHash, blockNumber, receiptIndex := GetTxLookupEntry(db, hash)

	if blockHash != (types.Hash{}) {
		receipts := GetBlockReceipts(db, blockHash, blockNumber)
		if len(receipts) <= int(receiptIndex) {
			logger.Error("Receipt refereced missing", "number", blockNumber, "hash", blockHash, "index", receiptIndex)
			return nil, types.Hash{}, 0, 0
		}
		//normal transaction receipt offset 1
		receiptIndex += 1
		return receipts[receiptIndex], blockHash, blockNumber, receiptIndex
	}
	// Old receipt representation, load the receipt and set an unknown metadata
	data, _ := db.Get(append(oldReceiptsPrefix, hash[:]...))
	if len(data) == 0 {
		return nil, types.Hash{}, 0, 0
	}
	var receipt transaction.Receipt
	byteBuf := bytes.NewBuffer(data)
	err := msgp.Decode(byteBuf, &receipt)
	if err != nil {
		logger.Error("Invalid receipt msgp", "hash", hash, "err", err)
	}
	return &receipt, types.Hash{}, 0, 0
}

// GetCoinBaseReceipt  returns coin base transaction receipt
func GetCoinBaseReceipt(db DatabaseReader, blockHash types.Hash, blockNumber uint64) (*transaction.Receipt, types.Hash, uint64, uint64) {
	receiptIndex := uint64(0)
	receipts := GetBlockReceipts(db, blockHash, blockNumber)
	if len(receipts) <= int(receiptIndex) {
		logger.Error("GetCoinBaseLog refereced missing", "number", blockNumber, "hash", blockHash, "index", receiptIndex)
		return nil, types.Hash{}, 0, 0
	}

	return receipts[receiptIndex], blockHash, blockNumber, receiptIndex
}

// GetBloomBits retrieves the compressed bloom bit vector belonging to the given
// section and bit index from the.
func GetBloomBits(db DatabaseReader, bit uint, section uint64, head types.Hash) ([]byte, error) {
	key := append(append(bloomBitsPrefix, make([]byte, 10)...), head.Bytes()...)

	binary.BigEndian.PutUint16(key[1:], uint16(bit))
	binary.BigEndian.PutUint64(key[3:], section)

	return db.Get(key)
}

// WriteCanonicalHash stores the canonical hash for the given block number.
func WriteCanonicalHash(db database.IDatabasePutter, hash types.Hash, number uint64) error {
	key := append(append(headerPrefix, encodeBlockNumber(number)...), numSuffix...)
	if err := db.Put(key, hash.Bytes()); err != nil {
		logger.Critical("Failed to store number to hash mapping", "err", err)
	}
	return nil
}

// WriteHeadHeaderHash stores the head header's hash.
func WriteHeadHeaderHash(db database.IDatabasePutter, hash types.Hash) error {
	if err := db.Put(headHeaderKey, hash.Bytes()); err != nil {
		logger.Critical("Failed to store last header's hash", "err", err)
	}
	return nil
}

// WriteHeadBlockHash stores the head block's hash.
func WriteHeadBlockHash(db database.IDatabasePutter, hash types.Hash) error {
	if err := db.Put(headBlockKey, hash.Bytes()); err != nil {
		logger.Critical("Failed to store last block's hash", "err", err)
	}
	return nil
}

// WriteHeadFastBlockHash stores the fast head block's hash.
func WriteHeadFastBlockHash(db database.IDatabasePutter, hash types.Hash) error {
	if err := db.Put(headFastKey, hash.Bytes()); err != nil {
		logger.Critical("Failed to store last fast block's hash", "err", err)
	}
	return nil
}

// WriteHeader serializes a block header into the database.
func WriteHeader(db database.IDatabasePutter, header *block.Header) error {
	var encData bytes.Buffer
	err := msgp.Encode(&encData, header)
	if err != nil {
		return err
	}
	hash := header.Hash().Bytes()
	num := header.Number.IntVal.Uint64()
	encNum := encodeBlockNumber(num)
	key := append(blockHashPrefix, hash...)
	if err := db.Put(key, encNum); err != nil {
		logger.Critical("Failed to store hash to number mapping", "err", err)
	}
	key = append(append(headerPrefix, encNum...), hash...)
	if err := db.Put(key, encData.Bytes()); err != nil {
		logger.Critical("Failed to store header", "err", err)
	}
	return nil
}

// WriteBody serializes the body of a block into the database.
func WriteBody(db database.IDatabasePutter, hash types.Hash, number uint64, body *block.Body) error {
	var buf bytes.Buffer
	err := msgp.Encode(&buf, body)
	if err != nil {
		return err
	}
	return WriteBodyMsgp(db, hash, number, buf.Bytes())
}

// WriteBodyMsgp writes a serialized body of a block into the database.
func WriteBodyMsgp(db database.IDatabasePutter, hash types.Hash, number uint64, data []byte) error {
	key := append(append(bodyPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
	if err := db.Put(key, data); err != nil {
		logger.Critical("Failed to store block body", "err", err)
	}
	return nil
}

// WriteBlock serializes a block into the database, header and body separately.
func WriteBlock(db database.IDatabasePutter, block *block.Block) error {
	// Store the body first to retain database consistency
	if err := WriteBody(db, block.Hash(), block.NumberU64(), block.Body()); err != nil {
		return err
	}
	// Store the header too, signaling full block ownership
	if err := WriteHeader(db, block.Header()); err != nil {
		return err
	}
	return nil
}

// WriteBlockReceipts stores all the transaction receipts belonging to a block
// as a single receipt slice. This is used during chain reorganisations for
// rescheduling dropped transactions.
func WriteBlockReceipts(db database.IDatabasePutter, hash types.Hash, number uint64, receipts transaction.Receipts) error {
	// Convert the receipts into their storage form and serialize them
	var buf bytes.Buffer
	for _, receipt := range receipts {
		err := msgp.Encode(&buf, receipt)
		if err != nil {
			logger.Critical("Failed to Encode receipts", "err", err)
			return err
		}
	}
	// Store the flattened receipt slice
	key := append(append(blockReceiptsPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
	if err := db.Put(key, buf.Bytes()); err != nil {
		logger.Critical("Failed to store block receipts", "err", err)
	}
	return nil
}

// WriteTxLookupEntries stores a positional metadata for every transaction from
// a block, enabling hash based transaction and receipt lookups.
func WriteTxLookupEntries(db database.IDatabasePutter, block *block.Block) error {
	// Iterate over each transaction and encode its metadata
	for i, tx := range block.Transactions() {
		entry := TxLookupEntry{
			BlockHash:  block.Hash(),
			BlockIndex: block.NumberU64(),
			Index:      uint64(i),
		}

		var data bytes.Buffer
		err := msgp.Encode(&data, &entry)
		if err != nil {
			return err
		}
		if err := db.Put(append(lookupPrefix, tx.Hash().Bytes()...), data.Bytes()); err != nil {
			return err
		}
		signer := transaction.NewMSigner(tx.ChainId())
		sender, _ := transaction.Sender(signer, tx)
		nonce := big.NewInt(int64(tx.Nonce()))

		if err := db.Put(append(append(trAddrNoncePrefix, sender.Bytes()...), nonce.Bytes()...), data.Bytes()); err != nil {
			return err
		}

	}
	return nil
}

// WriteBloomBits writes the compressed bloom bits vector belonging to the given
// section and bit index.
func WriteBloomBits(db database.IDatabasePutter, bit uint, section uint64, head types.Hash, bits []byte) {
	key := append(append(bloomBitsPrefix, make([]byte, 10)...), head.Bytes()...)

	binary.BigEndian.PutUint16(key[1:], uint16(bit))
	binary.BigEndian.PutUint64(key[3:], section)

	if err := db.Put(key, bits); err != nil {
		logger.Critical("Failed to store bloom bits", "err", err)
	}
}

func WriteBlockStat(db database.IDatabasePutter, hash types.Hash, stat *BlockStat) error {
	var data bytes.Buffer
	err := msgp.Encode(&data, stat)
	if err != nil {
		return err
	}
	key := append(totalTxPrefix, hash.Bytes()...)
	if err := db.Put(key, data.Bytes()); err != nil {
		logger.Critical("Failed to store block total txs", "err", err)
	}
	return nil
}

func GetBlockStat(db DatabaseReader, hash types.Hash) *BlockStat {
	key := append(totalTxPrefix, hash.Bytes()...)
	data, _ := db.Get(key)
	if len(data) == 0 {
		return nil
	}
	stat := new(BlockStat)
	byteBuf := bytes.NewBuffer(data)
	err := msgp.Decode(byteBuf, stat)
	if err != nil {
		logger.Error("GetTtxs Invalid block body msgp", "hash", hash, "err", err)
		return nil
	}
	return stat
}

func WriteExtra(db database.IDatabasePutter, key, val []byte) error {
	if err := db.Put(key, val); err != nil {
		logger.Critical("Failed to store block total txs", "err", err)
	}
	return nil
}

func GetExtra(db DatabaseReader, key []byte) []byte {
	data, _ := db.Get(key)
	if len(data) == 0 {
		return nil
	}
	return data
}

// DeleteCanonicalHash removes the number to hash canonical mapping.
func DeleteCanonicalHash(db DatabaseDeleter, number uint64) {
	db.Delete(append(append(headerPrefix, encodeBlockNumber(number)...), numSuffix...))
}

// DeleteHeader removes all block header data associated with a hash.
func DeleteHeader(db DatabaseDeleter, hash types.Hash, number uint64) {
	db.Delete(append(blockHashPrefix, hash.Bytes()...))
	db.Delete(append(append(headerPrefix, encodeBlockNumber(number)...), hash.Bytes()...))
}

// DeleteBody removes all block body data associated with a hash.
func DeleteBody(db DatabaseDeleter, hash types.Hash, number uint64) {
	db.Delete(append(append(bodyPrefix, encodeBlockNumber(number)...), hash.Bytes()...))
}

// DeleteBlock removes all block data associated with a hash.
func DeleteBlock(db DatabaseDeleter, hash types.Hash, number uint64) {
	DeleteBlockReceipts(db, hash, number)
	DeleteHeader(db, hash, number)
	DeleteBody(db, hash, number)
}

// DeleteBlockReceipts removes all receipt data associated with a block hash.
func DeleteBlockReceipts(db DatabaseDeleter, hash types.Hash, number uint64) {
	db.Delete(append(append(blockReceiptsPrefix, encodeBlockNumber(number)...), hash.Bytes()...))
}

// DeleteTxLookupEntry removes all transaction data associated with a hash.
func DeleteTxLookupEntry(db DatabaseDeleter, hash types.Hash) {
	db.Delete(append(lookupPrefix, hash.Bytes()...))
}

// DeleteTxAddrNonceLookupEntry removes all transaction data associated with a hash.
func DeleteTxAddrNonceLookupEntry(db DatabaseDeleter, tx *transaction.Transaction) {
	signer := transaction.NewMSigner(tx.ChainId())
	sender, _ := transaction.Sender(signer, tx)
	nonce := big.NewInt(int64(tx.Nonce()))
	db.Delete(append(append(trAddrNoncePrefix, sender.Bytes()...), nonce.Bytes()...))
}

// PreimageTable returns a Database instance with the key prefix for preimage entries.
func PreimageTable(db database.IDatabase) database.IDatabase {
	return database.NewTable(db, preimagePrefix)
}

// WritePreimages writes the provided set of preimages to the database. `number` is the
// current block number, and is used for debug messages only.
func WritePreimages(db database.IDatabase, number uint64, preimages map[types.Hash][]byte) error {
	table := PreimageTable(db)
	batch := table.NewBatch()
	hitCount := 0
	for hash, preimage := range preimages {
		_, err := table.Get(hash.Bytes())
		if err != nil {
			batch.Put(hash.Bytes(), preimage)
			hitCount++
		}
	}
	preimageCounter.Inc(int64(len(preimages)))
	preimageHitCounter.Inc(int64(hitCount))
	if hitCount > 0 {
		if err := batch.Write(); err != nil {
			return fmt.Errorf("preimage write fail for block %d: %v", number, err)
		}
	}
	return nil
}

// GetBlockChainVersion reads the version number from db.
func GetBlockChainVersion(db DatabaseReader) int {
	var vsn uint32
	enc, _ := db.Get([]byte("BlockchainVersion"))
	if len(enc) <= 0 {
		return 0
	}
	vsn = binary.BigEndian.Uint32(enc)
	return int(vsn)
}

// WriteBlockChainVersion writes vsn as the version number to db.
func WriteBlockChainVersion(db database.IDatabasePutter, vsn int) {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint32(enc, uint32(vsn))
	db.Put([]byte("BlockchainVersion"), enc)
}

// WriteChainConfig writes the chain config settings to the database.
func WriteChainConfig(db database.IDatabasePutter, hash types.Hash, cfg *params.ChainConfig) error {
	// short circuit and ignore if nil config. GetChainConfig
	// will return a default.
	if cfg == nil {
		return nil
	}

	jsonChainConfig, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return db.Put(append(configPrefix, hash[:]...), jsonChainConfig)
}

// GetChainConfig will fetch the network settings based on the given hash.
func GetChainConfig(db DatabaseReader, hash types.Hash) (*params.ChainConfig, error) {
	jsonChainConfig, _ := db.Get(append(configPrefix, hash[:]...))
	if len(jsonChainConfig) == 0 {
		return nil, ErrChainConfigNotFound
	}

	var config params.ChainConfig
	if err := json.Unmarshal(jsonChainConfig, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// FindCommonAncestor returns the last common ancestor of two block headers
func FindCommonAncestor(db DatabaseReader, a, b *block.Header) *block.Header {
	for bn := b.Number.IntVal.Uint64(); a.Number.IntVal.Uint64() > bn; {
		a = GetHeader(db, a.ParentHash, a.Number.IntVal.Uint64()-1)
		if a == nil {
			return nil
		}
	}
	for an := a.Number.IntVal.Uint64(); an < b.Number.IntVal.Uint64(); {
		b = GetHeader(db, b.ParentHash, b.Number.IntVal.Uint64()-1)
		if b == nil {
			return nil
		}
	}
	for a.Hash() != b.Hash() {
		a = GetHeader(db, a.ParentHash, a.Number.IntVal.Uint64()-1)
		if a == nil {
			return nil
		}
		b = GetHeader(db, b.ParentHash, b.Number.IntVal.Uint64()-1)
		if b == nil {
			return nil
		}
	}
	return a
}
