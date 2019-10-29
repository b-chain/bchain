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
// @File: backend.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

// Package bchainapi implements the general BCHAIN API functions.
package bchainapi

import (
	"context"

	"bchain.io/node/services/bchain/downloader"
	"bchain.io/utils/event"
	"bchain.io/accounts"
	"bchain.io/core/state"
	"bchain.io/core"
	"bchain.io/params"
	"bchain.io/utils/database"
	"bchain.io/communication/rpc"
	"bchain.io/core/blockchain/block"
	"bchain.io/common/types"
	"bchain.io/core/transaction"
	"bchain.io/core/blockchain"
	"bchain.io/consensus/apos"
)

// Backend interface provides the common API services (that are provided by
// both full and light clients) with access to necessary functions.
type Backend interface {
	// General bchain API
	Downloader() *downloader.Downloader
	ProtocolVersion() int
	ChainDb() database.IDatabase
	EventMux() *event.TypeMux
	AccountManager() *accounts.Manager

	// BlockChain API
	SetHead(number uint64)
	HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*block.Header, error)
	BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*block.Block, error)
	StatInfoByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*blockchain.BlockStat, error)
	StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *block.Header, error)
	GetBlock(ctx context.Context, blockHash types.Hash) (*block.Block, error)
	GetReceipts(ctx context.Context, blockHash types.Hash) (transaction.Receipts, error)
	SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription
	SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription
	SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription
	GetBlockCertificate(ctx context.Context, blockNr rpc.BlockNumber) apos.BlockCertificate

	// TxPool API
	SendTx(ctx context.Context, signedTx *transaction.Transaction) error
	GetPoolTransactions() (transaction.Transactions, error)
	GetPoolTransaction(txHash types.Hash) *transaction.Transaction
	GetPoolNonce(ctx context.Context, addr types.Address) (uint64, error)
	Stats() (pending int, queued int)
	TxPoolContent() (map[types.Address]transaction.Transactions, map[types.Address]transaction.Transactions)
	SubscribeTxPreEvent(chan<- core.TxPreEvent) event.Subscription

	ChainConfig() *params.ChainConfig
	CurrentBlock() *block.Block
	CurrentBlockNum()uint64
}

func GetAPIs(apiBackend Backend) []rpc.API {
	nonceLock := new(AddrLocker)
	return []rpc.API{
		{
			Namespace: "bchain",
			Version:   "1.0",
			Service:   NewPublicBchainAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "bchain",
			Version:   "1.0",
			Service:   NewPublicBlockChainAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "bchain",
			Version:   "1.0",
			Service:   NewPublicTransactionPoolAPI(apiBackend, nonceLock),
			Public:    true,
		}, {
			Namespace: "txpool",
			Version:   "1.0",
			Service:   NewPublicTxPoolAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "bchain",
			Version:   "1.0",
			Service:   NewPublicAccountAPI(apiBackend.AccountManager()),
			Public:    true,
		}, {
			Namespace: "personal",
			Version:   "1.0",
			Service:   NewPrivateAccountAPI(apiBackend, nonceLock),
			Public:    false,
		},
	}
}
