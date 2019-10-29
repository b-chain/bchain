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
// @File: genesis_test.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package genesis

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"math/big"
	"bchain.io/common/types"
	"bchain.io/core/blockchain"
	"bchain.io/params"
	"bchain.io/utils/crypto"
	"bchain.io/utils/database"
	"reflect"
	"testing"
)

var defaultGenesisHexHash = "a0623059b648f7b592166b5c165e5c0c79a018d11ef739c049d451884f8c23eb"

func TestDefaultGenesisBlock(t *testing.T) {
	block, _ := DefaultGenesisBlock().ToBlock()

	if hexHash := types.Bytes2Hex(block.Hash().Bytes()); hexHash != defaultGenesisHexHash {
		t.Errorf("wrong mainnet genesis hash, got %v,", hexHash)
	}
}

func TestSetupGenesis(t *testing.T) {
	var (
		customghash = types.HexToHash("0x11873b8c6b93ee8e86d78734a089f5f2b546cf29d3e9ec785fe84253eca91e6b")
		customg     = Genesis{
			Config: &params.ChainConfig{big.NewInt(500)},
			Alloc: GenesisAlloc{
				{1}: {Storage: map[types.Hash]types.Hash{{1}: {1}}},
			},
		}
		oldcustomg = customg

		customghash2 = types.HexToHash("0xab0d510f06569cb8a0a2e60eaa6a0bd1ba1a0fd5d94d5f49cc771bcae08ef648")
		customg2     = Genesis{
			Config: &params.ChainConfig{big.NewInt(700)},
			Alloc: GenesisAlloc{
				{1}: {Storage: map[types.Hash]types.Hash{{2}: {2}}},
			},
		}
	)
	oldcustomg.Config = &params.ChainConfig{ChainId: big.NewInt(2)}
	tests := []struct {
		name       string
		fn         func(database.IDatabase) (*params.ChainConfig, types.Hash, error)
		wantConfig *params.ChainConfig
		wantHash   types.Hash
		wantErr    error
	}{
		{
			name: "genesis without ChainConfig",
			fn: func(db database.IDatabase) (*params.ChainConfig, types.Hash, error) {
				return SetupGenesisBlock(db, new(Genesis))
			},
			wantErr:    errGenesisNoConfig,
			wantConfig: params.DefaultChainConfig,
		},
		{
			name: "no block in DB, genesis == nil",
			fn: func(db database.IDatabase) (*params.ChainConfig, types.Hash, error) {
				return SetupGenesisBlock(db, nil)
			},
			wantHash:   types.HexToHash(defaultGenesisHexHash),
			wantConfig: params.DefaultChainConfig,
		},
		{
			name: "test block in DB, genesis == nil",
			fn: func(db database.IDatabase) (*params.ChainConfig, types.Hash, error) {
				DefaultGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, nil)
			},
			wantHash:   types.HexToHash(defaultGenesisHexHash),
			wantConfig: params.DefaultChainConfig,
		},
		{
			name: "custom block in DB, genesis == nil",
			fn: func(db database.IDatabase) (*params.ChainConfig, types.Hash, error) {
				customg.MustCommit(db)
				return SetupGenesisBlock(db, nil)
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
		},
		{
			name: "custom block in DB, genesis == custom2",
			fn: func(db database.IDatabase) (*params.ChainConfig, types.Hash, error) {
				customg.MustCommit(db)
				return SetupGenesisBlock(db, &customg2)
			},
			wantErr:    &GenesisMismatchError{Stored: customghash, New: customghash2},
			wantHash:   customghash2,
			wantConfig: customg2.Config,
		},
		{
			name: "custom block in DB, genesis == nil",
			fn: func(db database.IDatabase) (*params.ChainConfig, types.Hash, error) {
				customg.MustCommit(db)
				return SetupGenesisBlock(db, nil)
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
		},
	}

	for _, test := range tests {
		db, _ := database.OpenMemDB()
		config, hash, err := test.fn(db)
		// Check the return values.
		if !reflect.DeepEqual(err, test.wantErr) {
			spew := spew.ConfigState{DisablePointerAddresses: true, DisableCapacities: true}
			t.Errorf("%s: returned error %#v, want %#v", test.name, spew.NewFormatter(err), spew.NewFormatter(test.wantErr))
		}
		if !reflect.DeepEqual(config, test.wantConfig) {
			t.Errorf("%s:\nreturned %v\nwant     %v", test.name, config, test.wantConfig)
		}
		if hash != test.wantHash {
			t.Errorf("%s: returned hash %s, want %s", test.name, hash.Hex(), test.wantHash.Hex())
		} else if err == nil {
			// Check database content.
			stored := blockchain.GetBlock(db, test.wantHash, 0)
			if stored.Hash() != test.wantHash {
				t.Errorf("%s: block in DB has hash %s, want %s", test.name, stored.Hash(), test.wantHash)
			}
		}
	}
}

func TestGenerateInnerContractAddress(t *testing.T) {
	type pair struct {
		contractName string
		creatorName  string
	}
	pairs := []pair{
		{"bchain.js", "0x55fda7601ffa55f61b819642816460aa24883f7f"},
		{"consensus.js", "bchain1"},
		{"system.js", "bchain2"},
		{"pledge.js", "bchain3"},
	}

	var acc string
	for _, pair := range pairs {
		if types.IsHexAddress(pair.creatorName) {
			acc = pair.creatorName
		} else {
			acc = types.StringToAddress(pair.creatorName).Hex()
		}
		contractAddr := crypto.CreateAddress(types.HexToAddress(acc), 0)
		fmt.Printf("contract name: %s  address: %s \n", pair.contractName, contractAddr.Hex())
	}
}
