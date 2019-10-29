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
// @File: api_backend.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package bchain

import (
	"context"

	"bchain.io/params"
	"bchain.io/core/state"
	"bchain.io/core"
	"bchain.io/utils/event"
	"bchain.io/node/services/bchain/downloader"
	"bchain.io/accounts"
	"bchain.io/core/blockchain/block"
	"bchain.io/communication/rpc"
	"bchain.io/common/types"
	"bchain.io/core/transaction"
	"bchain.io/core/blockchain"
	"bchain.io/utils/database"
	"bchain.io/utils/bloom"
	"bytes"
	"github.com/tinylib/msgp/msgp"
	"bchain.io/consensus/apos"
)

// BchainApiBackend implements bchainapi.Backend for full nodes
type BchainApiBackend struct {
	bchain *Bchain
}

func (b *BchainApiBackend) ChainConfig() *params.ChainConfig {
	return b.bchain.chainConfig
}

func (b *BchainApiBackend) CurrentBlock() *block.Block {
	return b.bchain.blockchain.CurrentBlock()
}

func (b *BchainApiBackend)CurrentBlockNum()uint64{
	return b.bchain.blockchain.CurrentBlockNum()
}

func (b *BchainApiBackend) SetHead(number uint64) {
	b.bchain.protocolManager.downloader.Cancel()
	b.bchain.blockchain.SetHead(number)
}

func (b *BchainApiBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*block.Header, error) {

	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.bchain.blockchain.CurrentBlock().Header(), nil
	}
	return b.bchain.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *BchainApiBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*block.Block, error) {

	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.bchain.blockchain.CurrentBlock(), nil
	}
	return b.bchain.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *BchainApiBackend) StatInfoByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*blockchain.BlockStat, error) {
	var hash types.Hash
	if blockNr == rpc.LatestBlockNumber {
		hash = b.bchain.blockchain.CurrentBlock().Hash()
	} else {
		blk := b.bchain.blockchain.GetBlockByNumber(uint64(blockNr))
		if blk == nil {
			return nil, nil
		}
		hash = blk.Hash()
	}
	return b.bchain.blockchain.GetBlockStat(hash), nil
}

func (b *BchainApiBackend) GetBlockCertificate(ctx context.Context, blockNr rpc.BlockNumber) apos.BlockCertificate {
	var hash types.Hash
	if blockNr == rpc.LatestBlockNumber {
		hash = b.bchain.blockchain.CurrentBlock().Hash()
	} else {
		blk := b.bchain.blockchain.GetBlockByNumber(uint64(blockNr))
		if blk == nil {
			return nil
		}
		hash = blk.Hash()
	}

	msgpData := b.bchain.blockchain.GetConsensusData(hash)
	if msgpData == nil || len(msgpData) == 0 {
		return nil
	}
	var blockCertificate apos.BlockCertificate
	byteBuf := bytes.NewBuffer(msgpData)
	err := msgp.Decode(byteBuf, &blockCertificate)
	if err != nil {
		logger.Error("GetBlockCertificate.Decode err", "hash", hash.String(), "err", err)
		return nil
	}
	return blockCertificate
}

func (b *BchainApiBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *block.Header, error) {

	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.bchain.BlockChain().StateAt(header.StateRootHash)
	return stateDb, header, err
}

func (b *BchainApiBackend) GetBlock(ctx context.Context, blockHash types.Hash) (*block.Block, error) {
	return b.bchain.blockchain.GetBlockByHash(blockHash), nil
}

func (b *BchainApiBackend) GetReceipts(ctx context.Context, blockHash types.Hash) (transaction.Receipts, error) {
	return blockchain.GetBlockReceipts(b.bchain.chainDb, blockHash, blockchain.GetBlockNumber(b.bchain.chainDb, blockHash)), nil
}

func (b *BchainApiBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.bchain.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *BchainApiBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.bchain.BlockChain().SubscribeChainEvent(ch)
}

func (b *BchainApiBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.bchain.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *BchainApiBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.bchain.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *BchainApiBackend) SubscribeLogsEvent(ch chan<- []*transaction.Log) event.Subscription {
	return b.bchain.BlockChain().SubscribeLogsEvent(ch)
}

func (b *BchainApiBackend) SendTx(ctx context.Context, signedTx *transaction.Transaction) error {
	return b.bchain.txPool.AddLocal(signedTx)
}

func (b *BchainApiBackend) GetPoolTransactions() (transaction.Transactions, error) {
	pending, err := b.bchain.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs transaction.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *BchainApiBackend) GetPoolTransaction(hash types.Hash) *transaction.Transaction {
	return b.bchain.txPool.Get(hash)
}

func (b *BchainApiBackend) GetPoolNonce(ctx context.Context, addr types.Address) (uint64, error) {
	return b.bchain.txPool.State().GetNonce(addr), nil
}

func (b *BchainApiBackend) Stats() (pending int, queued int) {
	return b.bchain.txPool.Stats()
}

func (b *BchainApiBackend) TxPoolContent() (map[types.Address]transaction.Transactions, map[types.Address]transaction.Transactions) {
	return b.bchain.TxPool().Content()
}

func (b *BchainApiBackend) SubscribeTxPreEvent(ch chan<- core.TxPreEvent) event.Subscription {
	return b.bchain.TxPool().SubscribeTxPreEvent(ch)
}

func (b *BchainApiBackend) Downloader() *downloader.Downloader {
	return b.bchain.Downloader()
}

func (b *BchainApiBackend) ProtocolVersion() int {
	return b.bchain.BchainVersion()
}

func (b *BchainApiBackend) ChainDb() database.IDatabase {
	return b.bchain.ChainDb()
}

func (b *BchainApiBackend) EventMux() *event.TypeMux {
	return b.bchain.EventMux()
}

func (b *BchainApiBackend) AccountManager() *accounts.Manager {
	return b.bchain.AccountManager()
}

func (b *BchainApiBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.bchain.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *BchainApiBackend) ServiceFilter(ctx context.Context, session *bloom.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.bchain.bloomRequests)
	}
}
