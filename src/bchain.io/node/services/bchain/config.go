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
// @File: config.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package bchain

import (
	"os"
	"os/user"


	"bchain.io/common/types"
	"bchain.io/node/services/bchain/downloader"
	"bchain.io/core/txprocessor"
	"bchain.io/core/genesis"
	"bchain.io/node"
	"bchain.io/params"
)

// DefaultConfig contains default settings for use on the Bchain main net.
var DefaultConfig = Config{
	SyncMode: downloader.FastSync,

	NetworkId:     1,
	LightPeers:    20,
	DatabaseCache: 128,

	TxPool: txprocessor.DefaultTxPoolConfig,

}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}

}

//go:generate gencodec -type Config -field-override configMarshaling -formats toml -out gen_config.go

type Config struct {
	// The genesis block, which is inserted if the database is empty.
	// If nil, the Bchain main net block is used.
	Genesis *genesis.Genesis `toml:"-"`

	// Protocol options
	NetworkId uint64 // Network ID to use for selecting peers to connect to
	SyncMode  downloader.SyncMode

	// Light client options
	LightServ  int `toml:",omitempty"` // Maximum percentage of time allowed for serving LES requests
	LightPeers int `toml:",omitempty"` // Maximum number of LES client peers

	// Database options
	SkipBcVersionCheck bool `toml:"-"`
	DatabaseHandles    int  `toml:"-"`
	DatabaseCache      int

	// Producing-related options
	Coinbase    types.Address `toml:",omitempty"`
	BlockproducerThreads int  `toml:",omitempty"`
	ExtraData    []byte       `toml:",omitempty"`


	// Transaction pool options
	TxPool txprocessor.TxPoolConfig



	// Enables tracking of SHA3 preimages in the VM
	EnablePreimageRecording bool

	// Miscellaneous options
	DocRoot string `toml:"-"`

	//should we start blockproducer at first
	StartBlockproducerAtStart bool

	// block options
	MaxBlockTxSize uint64
}

func (c *Config) SetDefaultConfig() error{
	c.Genesis = genesis.DefaultGenesisBlock()
	c.NetworkId = params.DefaultChainConfig.ChainId.Uint64()
	c.TxPool = txprocessor.DefaultTxPoolConfig
	c.StartBlockproducerAtStart = true
	c.MaxBlockTxSize = 10*1024*1024
	return nil
}

func (c *Config) SetSpecialConfig() error{
	c.Genesis = genesis.DefaultGenesisBlock()
	return nil
}

type configMarshaling struct {
	ExtraData types.BytesForJson
}


type BchaindConfig struct{
	Bchain Config
	Node node.Config
}