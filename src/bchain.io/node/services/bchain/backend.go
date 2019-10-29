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
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

// Package bchain implements the Bchain protocol.
package bchain

import (
	"errors"
	"fmt"
	"bchain.io/consensus"
	"sync"
	"sync/atomic"

	"bchain.io/accounts"
	"bchain.io/common/types"
	"bchain.io/communication/p2p"
	"bchain.io/communication/rpc"
	"bchain.io/core/blockchain"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/chainindexer"
	"bchain.io/core/txprocessor"
	"bchain.io/node"
	"bchain.io/node/services/bchain/downloader"
	"bchain.io/params"
	"bchain.io/utils/bloom"
	"bchain.io/utils/database"
	"bchain.io/utils/event"
	//"bchain.io/core/genesis"

	"crypto/ecdsa"
	"bchain.io/accounts/keystore"
	"bchain.io/blockproducer"
	"bchain.io/communication/rpc/bchainapi"
	"bchain.io/consensus/apos"
	"bchain.io/core/genesis"
	"bchain.io/core/interpreter"
	"bchain.io/core/interpreter/jsre"
	"bchain.io/core/interpreter/wasmre"
	"bchain.io/bchaind/config"
	"bchain.io/node/services/bchain/filters"
)

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *chainindexer.ChainIndexer)
}

// Bchain implements the Bchain full node service.
type Bchain struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan  chan bool    // Channel for shutting down the bchain
	stopDbUpgrade func() error // stop chain db sequential key upgrade

	// Handlers
	txPool          *txprocessor.TxPool
	blockchain      *blockchain.BlockChain
	protocolManager *ProtocolManager
	lesServer       LesServer

	// DB interfaces
	chainDb database.IDatabase // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests chan chan *bloom.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *chainindexer.ChainIndexer // Bloom indexer operating during block imports

	ApiBackend *BchainApiBackend

	blockproducer *blockproducer.Blockproducer //create newBlock
	aposConsensus *apos.Apos                   //judge the last Block
	//interVm             *interpreter.Vms

	coinbase types.Address

	networkId uint64
	//netRPCService *bchainapi.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. coinbase)
}

func (s *Bchain) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

type SetupGenesisResult struct {
	ChainConfig  *params.ChainConfig
	GennesisHash *types.Hash
	GenesisErr   error
	ChainDb      *database.IDatabase
}

// New creates a new Bchain object (including the
// initialisation of the common Bchain object)----Move to bchain2.go
/**/
func New(ctx *node.ServiceContext) (*Bchain, error) {
	c := config.GetConfigInstance()
	var config = &Config{}
	err := c.Register("bchain", config)
	if err != nil {
		logger.Error("get config fail", "err", err)
	}

	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run bchain.Bchain in light sync mode, use les.LightBchain")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}
	stopDbUpgrade := upgradeDeduplicateData(chainDb)
	chainConfig, genesisHash, genesisErr := genesis.SetupGenesisBlock(chainDb, config.Genesis)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}
	logger.Info("Initialised chain configuration", "config", chainConfig)

	bchain := &Bchain{
		config:         config,
		chainDb:        chainDb,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		shutdownChan:   make(chan bool),
		stopDbUpgrade:  stopDbUpgrade,
		networkId:      config.NetworkId,
		coinbase:       config.Coinbase,
		bloomRequests:  make(chan chan *bloom.Retrieval),
		bloomIndexer:   chainindexer.NewBloomIndexer(chainDb, params.BloomBitsBlocks),
	}

	logger.Info("Initialising Bchain protocol", "versions", ProtocolVersions, "network", config.NetworkId)
	bchain.engine = CreateConsensusEngine(bchain)
	if !config.SkipBcVersionCheck {
		bcVersion := blockchain.GetBlockChainVersion(chainDb)
		if bcVersion != blockchain.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d). Run bchaind upgradedb.\n", bcVersion, blockchain.BlockChainVersion)
		}
		blockchain.WriteBlockChainVersion(chainDb, blockchain.BlockChainVersion)
	}

	bchain.blockchain, err = blockchain.NewBlockChain(chainDb, bchain.chainConfig, bchain.engine)
	if err != nil {
		return nil, err
	}
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		logger.Warn("Rewinding chain to upgrade configuration", "err", compat)
		bchain.blockchain.SetHead(compat.RewindTo)
		blockchain.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	bchain.bloomIndexer.Start(bchain.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}

	interpreter.Singleton().Register(func() interpreter.PluginImpl {
		return &jsre.JSRE{}
	})

	interpreter.Singleton().Register(wasmre.NewWasmRe)

	interpreter.Singleton().Initialize()
	interpreter.Singleton().Startup()

	bchain.txPool = txprocessor.NewTxPool(config.TxPool, bchain.chainConfig, bchain.blockchain)

	//Init blockProducer
	bchain.blockproducer = blockproducer.New(bchain, bchain.chainConfig, bchain.EventMux(), bchain.engine, config.MaxBlockTxSize)

	//Init Consensus
	bchain.aposConsensus = apos.NewApos(bchain.blockchain, bchain.blockproducer, bchain.EventMux())

	if bchain.protocolManager, err = NewProtocolManager(bchain.chainConfig, config.SyncMode, config.NetworkId, bchain.eventMux, bchain.txPool, bchain.engine, bchain.blockchain, chainDb, bchain.aposConsensus); err != nil {
		return nil, err
	}

	bchain.ApiBackend = &BchainApiBackend{bchain}

	fmt.Println("New......Bchain")
	return bchain, nil
}

func makeExtraData(extra []byte) []byte {
	return make([]byte, 0)
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (database.IDatabase, error) {

	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if db, ok := db.(*database.LDatabase); ok {
		db.Meter("bchain/db/chaindata/")
	}
	return db, nil
}

// CreateConsensusEngine creates the required type of consensus engine instance for an Bchain service
func CreateConsensusEngine(bchain *Bchain) consensus.Engine {
	engine := apos.NewAposEngine(nil)
	return engine
}

func (s *Bchain) SetEngineKey(pri *ecdsa.PrivateKey) {
	switch v := s.engine.(type) {
	case *consensus.Engine_basic:
		v.SetKey(pri)
	case *apos.EngineApos:
		v.SetKey(pri)
	}
}

// APIs returns the collection of RPC services the bchain package offers.
// NOTE, some of these services probably need to be moved to somewhere else.

func (s *Bchain) APIs() []rpc.API {
	apis := bchainapi.GetAPIs(s.ApiBackend)

	//create New
	//apis := make([]rpc.API , 0)

	return append(apis, []rpc.API{
		{
			Namespace: "bchain",
			Version:   "1.0",
			Service:   NewPublicBchainAPI(s),
			Public:    true,
		}, {
			Namespace: "blockproducer",
			Version:   "1.0",
			Service:   NewPrivateBlockproducerAPI(s),
			Public:    true,
		}, {
			Namespace: "bchain",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, false),
			Public:    true,
		}, {
			Namespace: "bchain",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		},
	}...)

}

func (s *Bchain) ResetWithGenesisBlock(gb *block.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *Bchain) Coinbase() (eb accounts.Account, err error) {

	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			coinbase := accounts[0].Address

			s.lock.Lock()
			s.coinbase = coinbase
			s.lock.Unlock()

			logger.Infof("Coinbase automatically configured address:0x%x\n", coinbase)
			return accounts[0], nil
		}
	}
	return accounts.Account{}, fmt.Errorf("Coinbase must be explicitly specified")
}

// set in js console via admin interface or wrapper from cli flags
func (self *Bchain) SetCoinbase(coinbase types.Address) {
	self.lock.Lock()
	self.coinbase = coinbase
	self.lock.Unlock()

	self.blockproducer.SetCoinbase(coinbase)
}

func (s *Bchain) StartProducing(local bool, password string) error {
	eb, err := s.Coinbase()
	if err != nil {
		logger.Error("Cannot start producing without coinbase", "err", err)
		return fmt.Errorf("coinbase missing: %v", err)
	}

	//get key
	ks := s.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	key, err := ks.GetKeyWithPassphrase(eb, password)
	if err != nil {
		logger.Error("Cannot start producing without coinbase, get sign key err ", "err", err)
		return fmt.Errorf("get sign key err: %v", err)
	}

	s.SetEngineKey(key)
	s.aposConsensus.SetPriKey(key)
	s.aposConsensus.SetCoinBase(eb.Address)
	s.blockproducer.SetPriKey(key)
	if local {
		// If local (CPU) producing is started, we can disable the transaction rejection
		// mechanism introduced to speed sync times. CPU producing on mainnet is ludicrous
		// so noone will ever hit this path, whereas marking sync done on CPU producing
		// will ensure that private networks work in single blockproducer mode too.
		atomic.StoreUint32(&s.protocolManager.acceptTxs, 1)
	}
	go s.blockproducer.Start(eb.Address)

	//start the apos now
	go s.aposConsensus.Start()

	return nil
}

func (s *Bchain) StopProducing()                              { s.blockproducer.Stop() }
func (s *Bchain) IsProducing() bool                           { return s.blockproducer.Producing() }
func (s *Bchain) Blockproducer() *blockproducer.Blockproducer { return s.blockproducer }

func (s *Bchain) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *Bchain) BlockChain() *blockchain.BlockChain { return s.blockchain }
func (s *Bchain) TxPool() *txprocessor.TxPool        { return s.txPool }
func (s *Bchain) EventMux() *event.TypeMux           { return s.eventMux }
func (s *Bchain) Engine() consensus.Engine           { return s.engine }
func (s *Bchain) ChainDb() database.IDatabase        { return s.chainDb }
func (s *Bchain) IsListening() bool                  { return true } // Always listening
func (s *Bchain) BchainVersion() int                   { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *Bchain) NetVersion() uint64                 { return s.networkId }
func (s *Bchain) Downloader() *downloader.Downloader { return s.protocolManager.downloader }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *Bchain) Protocols() []p2p.Protocol {
	if s.lesServer == nil {
		return s.protocolManager.SubProtocols
	}
	return append(s.protocolManager.SubProtocols, s.lesServer.Protocols()...)
}

// Start implements node.Service, starting all internal goroutines needed by the
// Bchain protocol implementation.
func (s *Bchain) Start(srvr *p2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers
	if s.config.LightServ > 0 {
		maxPeers -= s.config.LightPeers
		if maxPeers < srvr.MaxPeers/2 {
			maxPeers = srvr.MaxPeers / 2
		}
	}
	// Start the networking layer and the light server if requested
	s.protocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}

	_, err := s.Coinbase()
	if err != nil {
		fmt.Println("[Warn]No CoinBase Do Not Start Producing Block!!!!!!!!!!!!!!!!!s")
	}
	//when start bchain service,not start blockproducer,except the cmd order we should start it
	if s.config.StartBlockproducerAtStart {
		fmt.Println("Start Blockproducer At Service Start.......................New")
		//s.blockproducer.Start(eb.Address)
	} else {
		fmt.Println("Not Start Blockproducer At Service Start........................")
	}

	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Bchain protocol.
func (s *Bchain) Stop() error {
	if s.stopDbUpgrade != nil {
		s.stopDbUpgrade()
	}
	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}
	s.txPool.Stop()
	s.blockproducer.Stop()
	s.eventMux.Stop()

	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
