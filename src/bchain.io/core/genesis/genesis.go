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
// @File: genesis.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package genesis

import (
	"errors"
	"fmt"
	"math/big"
	"bchain.io/common/types"
	"bchain.io/consensus/apos/aposgensis"
	"bchain.io/core/blockchain"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/interpreter/contract_parser"
	"bchain.io/core/interpreter/contract/wasmre/deps"
	"bchain.io/core/state"
	"bchain.io/params"
	"bchain.io/utils/crypto"
	"bchain.io/utils/database"
)

var (
	bchain_cc          = deps.MustAsset("bchain.json")
	pledge_cc          = deps.MustAsset("pledge.json")
	system_cc          = deps.MustAsset("system.json")
)

var errGenesisNoConfig = errors.New("genesis has no chain configuration")

// Genesis specifies the header fields, state of a genesis block. It also defines hard
// fork switch-over blocks through the chain configuration.
type Genesis struct {
	Config    *params.ChainConfig `json:"config"`
	Timestamp uint64              `json:"timestamp"`
	Alloc     GenesisAlloc        `json:"alloc"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number        uint64              `json:"number"`
	ParentHash    types.Hash          `json:"parentHash"`
	ConsensusData block.ConsensusData `json:"consensusData"`
	Extra         []ExtraStorage      `json:"extra"`
}

type GenesisAlloc map[types.Address]GenesisAccount

// GenesisAccount is an account in the state of the genesis block.
type GenesisAccount struct {
	Code            []byte                    `json:"code,omitempty"`
	Storage         map[types.Hash]types.Hash `json:"storage,omitempty"`
	Nonce           uint64                    `json:"nonce,omitempty"`
	InterpreterName string                    `json:"interpreterName,omitempty"`
	Creator         string                    `json:"creator,omitempty"`
}

type ExtraStorage struct {
	Key []byte `json:"key"`
	Val []byte `json:"val"`
}

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored types.Hash
	New    types.Hash
}

func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database already contains an incompatible genesis block (have %x, new %x)", e.Stored[:8], e.New[:8])
}

// SetupGenesisBlock writes or updates the genesis block in db.
// The block that will be used is:
//
//                          genesis == nil       genesis != nil
//                       +------------------------------------------
//     db has no genesis |  main-net default  |  genesis
//     db has genesis    |  from DB           |  genesis (if compatible)
//
// The stored chain configuration will be updated if it is compatible (i.e. does not
// specify a fork block below the local head block). In case of a conflict, the
// error is a *params.ConfigCompatError and the new, unwritten config is returned.
//
// The returned chain configuration is never nil.
func SetupGenesisBlock(db database.IDatabase, genesis *Genesis) (*params.ChainConfig, types.Hash, error) {
	if genesis != nil && genesis.Config == nil {
		return params.DefaultChainConfig, types.Hash{}, errGenesisNoConfig
	}

	// Just commit the new block if there is no stored genesis block.
	stored := blockchain.GetCanonicalHash(db, 0)
	if (stored == types.Hash{}) {
		if genesis == nil {
			logger.Info("Writing default test-net genesis block")
			genesis = DefaultGenesisBlock()
		} else {
			logger.Info("Writing custom genesis block")
		}
		block, err := genesis.Commit(db)
		return genesis.Config, block.Hash(), err
	}

	// Check whether the genesis block is already written.
	if genesis != nil {
		block, _ := genesis.ToBlock()
		hash := block.Hash()
		if hash != stored {
			return genesis.Config, block.Hash(), &GenesisMismatchError{stored, hash}
		} else {
			return genesis.Config, block.Hash(), nil
		}
	}

	//genesis == nil, return the dabatase genesis block
	storedcfg, err := blockchain.GetChainConfig(db, stored)
	if err != nil {
		if err == blockchain.ErrChainConfigNotFound {
			// This case happens if a genesis write was interrupted.
			logger.Warn("Found genesis block without chain config")
			err = blockchain.WriteChainConfig(db, stored, params.DefaultChainConfig)
		}
		return params.DefaultChainConfig, stored, err
	}
	return storedcfg, stored, nil

}

func (g *Genesis) configOrDefault(ghash types.Hash) *params.ChainConfig {
	switch {
	case g != nil:
		return g.Config
	default:
		return params.DefaultChainConfig
	}
}

// DefaultGenesisBlock returns the bchain main net genesis block.
func DefaultGenesisBlock() *Genesis {
	//in := hexutil.MustDecodeHex("0x571be45c5e74ddafe129d20fe3e21c4c76be925aa32967cda0b073a7cb51cb9e")
	//hash := types.Hash{}
	//hash.SetBytes(in)
	defaultBchainCode := contract_parser.ParseCodeByData(bchain_cc)
	defaultSystemCode := contract_parser.ParseCodeByData(system_cc)
	defaultPledgeCode := contract_parser.ParseCodeByData(pledge_cc)
	cdata, key, val := aposgensis.MakeAposGenesisConsensusData()
	exData := []ExtraStorage{}
	exData = append(exData, ExtraStorage{key, val})

	return &Genesis{
		Config: params.DefaultChainConfig,
		Alloc: map[types.Address]GenesisAccount{
			crypto.CreateAddress(types.Address{}, 0): {
				Code:            []byte(defaultBchainCode.Code),
				InterpreterName: defaultBchainCode.InterName,
				Creator:         types.StringToAddress("bchain").Hex(),
			},
			crypto.CreateAddress(types.Address{}, 1): {
				Code:            []byte(defaultSystemCode.Code),
				InterpreterName: defaultSystemCode.InterName,
				Creator:         types.StringToAddress("system").Hex(),
			},
			crypto.CreateAddress(types.Address{}, 2): {
				Code:            []byte(defaultPledgeCode.Code),
				InterpreterName: defaultPledgeCode.InterName,
				Creator:         types.StringToAddress("pledge").Hex(),
			},
		},
		ConsensusData: *cdata,
		Extra:         exData,
	}
}

// ToBlock creates the block and state of a genesis specification.
func (g *Genesis) ToBlock() (*block.Block, *state.StateDB) {
	db, _ := database.OpenMemDB()
	statedb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	for key, account := range g.Alloc {
		contractAddr := crypto.CreateAddress(types.HexToAddress(account.Creator), account.Nonce)
		logger.Info("Genesis contractAddr", contractAddr.Hex(),"crate by", account.Creator)
		logger.Info("Genesis keyAddr:", key.Hex())
		statedb.SetInterpreterID(contractAddr, crypto.Keccak256Hash([]byte(account.InterpreterName)))
		statedb.SetCreator(contractAddr, types.HexToAddress(account.Creator))
		statedb.SetCode(contractAddr, account.Code)
		statedb.SetNonce(contractAddr, account.Nonce)
		for key, value := range account.Storage {
			statedb.SetState(contractAddr, key, value)
		}
	}
	//for i:=0;i<3;i++{
	//	logger.Infof("i=%d addr=%s\n" , i , crypto.CreateAddress(types.Address{} , uint64(i)).Hex())
	//}
	root := statedb.IntermediateRoot()
	head := &block.Header{
		Number:        types.NewBigInt(*new(big.Int).SetUint64(g.Number)),
		Time:          types.NewBigInt(*new(big.Int).SetUint64(g.Timestamp)),
		ParentHash:    g.ParentHash,
		StateRootHash: root,
		Cdata:         g.ConsensusData,
	}

	return block.NewBlock(head, nil, nil), statedb
}

// Commit writes the block and state of a genesis specification to the database.
// The block is committed as the canonical head block.
func (g *Genesis) Commit(db database.IDatabase) (*block.Block, error) {
	block, statedb := g.ToBlock()
	if block.Number().Sign() != 0 {
		return nil, fmt.Errorf("can't commit genesis block with number > 0")
	}
	_, si, err := statedb.CommitTo(db, false)
	if err != nil {
		return nil, fmt.Errorf("cannot write state: %v", err)
	}
	stat := &blockchain.BlockStat{
		Ttxs:        types.NewBigInt(*new(big.Int)),
		TsoContract: types.NewBigInt(*big.NewInt(int64(si.TnewsoContract))),
		TsoNormal:   types.NewBigInt(*big.NewInt(int64(si.TnewsoNormal))),
		TstateNum:   types.NewBigInt(*big.NewInt(int64(si.TnewState))),
	}
	for _, exdata := range g.Extra {
		if err := blockchain.WriteExtra(db, exdata.Key, exdata.Val); err != nil {
			return nil, err
		}
	}
	if err := blockchain.WriteBlockStat(db, block.Hash(), stat); err != nil {
		return nil, err
	}
	if err := blockchain.WriteBlock(db, block); err != nil {
		return nil, err
	}
	if err := blockchain.WriteBlockReceipts(db, block.Hash(), block.NumberU64(), nil); err != nil {
		return nil, err
	}
	if err := blockchain.WriteCanonicalHash(db, block.Hash(), block.NumberU64()); err != nil {
		return nil, err
	}
	if err := blockchain.WriteHeadBlockHash(db, block.Hash()); err != nil {
		return nil, err
	}
	if err := blockchain.WriteHeadHeaderHash(db, block.Hash()); err != nil {
		return nil, err
	}
	config := g.Config
	if config == nil {
		config = params.DefaultChainConfig
	}
	return block, blockchain.WriteChainConfig(db, block.Hash(), config)
}

/*for test*/
// MustCommit writes the genesis block and state to db, panicking on error.
// The block is committed as the canonical head block.
func (g *Genesis) MustCommit(db database.IDatabase) *block.Block {
	block, err := g.Commit(db)
	if err != nil {
		panic(err)
	}
	return block
}

func GenesisBlockForTesting(db database.IDatabase, addr types.Address, balance *big.Int) *block.Block {
	//todo
	g := Genesis{Alloc: GenesisAlloc{addr: {}}}
	return g.MustCommit(db)
}
