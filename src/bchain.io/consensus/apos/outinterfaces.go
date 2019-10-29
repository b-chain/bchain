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
// @File: outinterfaces.go
// @Date: 2018/06/15 10:26:15
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"crypto/ecdsa"
	"bchain.io/common/types"
	"bchain.io/core"
	"bchain.io/core/blockchain/block"
	"bchain.io/utils/event"
)

/*
For out caller
*/

type dataPack interface {
}

type OutMsger interface {
	GetDataMsg() <-chan dataPack
	SendInner(dataPack) error
	Send2Apos(dataPack)
	PropagateMsg(dataPack) error //just send propagate msg to other nodes
	setAposRunning(bool)
}

//some out tools offered by Bchain,such as signer and blockInfo getter
type CommonTools interface {
	//Note: apos no right to hold a privateKey , so the Sig in commontools just though a privateKey to apos
	Sig(pCs *CredentialSign) error
	Esig(pEphemeralSign *EphemeralSign) error
	SigHash(hash types.Hash) []byte
	SeedSig(psd *SeedData) error

	//SigVerify(hash types.Hash, sig *SignatureVal) error
	//Sender(hash types.Hash, sig *SignatureVal) (types.Address, error)

	ESigVerify(hash types.Hash, sig []byte) error
	ESender(hash types.Hash, sig []byte) (types.Address, error)

	GetLastQrSignature() []byte
	GetQrSignature(round uint64) []byte
	GetNowBlockNum() uint64
	GetNextRound() int
	GetNowBlockHash() types.Hash

	SetPriKey(priKey *ecdsa.PrivateKey)
	SetCoinBase(coinbase types.Address)
	GetCoinBase() types.Address

	CreateTmpPriKey(step int)
	DelTmpKey(step int)
	ClearTmpKeys()

	GetProducerNewBlock(data *block.ConsensusData, timeLimit int64) *block.Block //get a new block from block producer
	MakeEmptyBlock(data *block.ConsensusData) *block.Block
	InsertChain(chain block.Blocks) (int, error)
	GetCurrentBlock() *block.Block
	GetBlockByNum(num uint64) *block.Block
	VerifyNextRoundBlock(block *block.Block)bool

	//version 1.1
	GetBlockByHash(hash types.Hash) *block.Block
	SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription

	WriteBlockCertificate(blk *block.Block, certificate BlockCertificate) error
	GetBlockCertificate(blockHash types.Hash) BlockCertificate
	GetWeight(r uint64, addr types.Address) (int64, int64)
}
