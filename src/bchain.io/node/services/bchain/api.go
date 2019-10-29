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
// @File: api.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package bchain

import (
	"compress/gzip"
	"fmt"
	"io"

	"os"
	"strings"

	"bchain.io/common/types"

	"bchain.io/core/blockchain/block"
	"github.com/tinylib/msgp/msgp"
	"bchain.io/core/blockchain"
)

// PublicBchainAPI provides an API to access Bchain full node-related
// information.
type PublicBchainAPI struct {
	e *Bchain
}

// NewPublicBchainAPI creates a new Bchain protocol API for full nodes.
func NewPublicBchainAPI(e *Bchain) *PublicBchainAPI {
	return &PublicBchainAPI{e}
}


// Coinbase is the address that producing rewards will be send to (alias for Coinbase)
func (api *PublicBchainAPI) Coinbase() (types.Address, error) {
	return api.Coinbase()
}

// Hashrate returns the POW hashrate
func (api *PublicBchainAPI) Hashrate() types.Uint64ForJson{
	return types.Uint64ForJson(api.e.Blockproducer().HashRate())
}



// PrivateBlockproducerAPI provides private RPC methods to control the blockproducer.
// These methods can be abused by external users and must be considered insecure for use by untrusted users.
type PrivateBlockproducerAPI struct {
	e *Bchain
}

// NewPrivateBlockproducerAPI create a new RPC service which controls the blockproducer of this node.
func NewPrivateBlockproducerAPI(e *Bchain) *PrivateBlockproducerAPI {
	return &PrivateBlockproducerAPI{e: e}
}

// Start the blockproducer with the given number of threads. If threads is nil the number
// of workers started is equal to the number of logical CPUs that are usable by
// this process. If producing is already running, this method adjust the number of
// threads allowed to use.
func (api *PrivateBlockproducerAPI) Start(threads *int, password string) error {
	// Set the number of threads if the seal engine supports it
	if threads == nil {
		threads = new(int)
	} else if *threads == 0 {
		*threads = -1 // Disable the blockproducer from within
	}
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := api.e.engine.(threaded); ok {
		logger.Info("Updated producing threads", "threads", *threads)
		th.SetThreads(*threads)
	}
	// Start the blockproducer and return
	if !api.e.IsProducing() {
		// Propagate the initial price point to the transaction pool
		api.e.lock.RLock()
		api.e.lock.RUnlock()

		return api.e.StartProducing(true, password)
	}
	return nil
}

// Stop the blockproducer
func (api *PrivateBlockproducerAPI) Stop() bool {
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := api.e.engine.(threaded); ok {
		th.SetThreads(-1)
	}
	api.e.StopProducing()
	return true
}

// SetExtra sets the extra data string that is included when this blockproducer produces a block.
func (api *PrivateBlockproducerAPI) SetExtra(extra string) (bool, error) {
	return false, nil
}

// SetCoinbase sets the coinbase of the blockproducer
func (api *PrivateBlockproducerAPI) SetCoinbase(coinbase types.Address) bool {
	api.e.SetCoinbase(coinbase)
	return true
}

// GetHashrate returns the current hashrate of the blockproducer.
func (api *PrivateBlockproducerAPI) GetHashrate() uint64 {
	return uint64(api.e.blockproducer.HashRate())
}

// PrivateAdminAPI is the collection of Bchain full node-related APIs
// exposed over the private admin endpoint.
type PrivateAdminAPI struct {
	bchain *Bchain
}

// NewPrivateAdminAPI creates a new API definition for the full node private
// admin methods of the Bchain service.
func NewPrivateAdminAPI(bchain *Bchain) *PrivateAdminAPI {
	return &PrivateAdminAPI{bchain: bchain}
}

// ExportChain exports the current blockchain into a local file.
func (api *PrivateAdminAPI) ExportChain(file string) (bool, error) {
	// Make sure we can create the file to export into
	out, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return false, err
	}
	defer out.Close()

	var writer io.Writer = out
	if strings.HasSuffix(file, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}

	// Export the blockchain
	if err := api.bchain.BlockChain().Export(writer); err != nil {
		return false, err
	}
	return true, nil
}

func hasAllBlocks(chain *blockchain.BlockChain, bs []*block.Block) bool {
	for _, b := range bs {
		if !chain.HasBlock(b.Hash(), b.NumberU64()) {
			return false
		}
	}

	return true
}

// ImportChain imports a blockchain from a local file.
func (api *PrivateAdminAPI) ImportChain(file string) (bool, error) {
	// Make sure the can access the file to import
	in, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer in.Close()

	var reader io.Reader = in
	if strings.HasSuffix(file, ".gz") {
		if reader, err = gzip.NewReader(reader); err != nil {
			return false, err
		}
	}


	blocks, index := make([]*block.Block, 0, 2500), 0
	for batch := 0; ; batch++ {
		// Load a batch of blocks from the input file
		for len(blocks) < cap(blocks) {
			block := new(block.Block)

			if err := msgp.Decode(reader , block); err == io.EOF {
				break
			} else if err != nil {
				return false, fmt.Errorf("block %d: failed to parse: %v", index, err)
			}
			blocks = append(blocks, block)
			index++
		}
		if len(blocks) == 0 {
			break
		}

		if hasAllBlocks(api.bchain.BlockChain(), blocks) {
			blocks = blocks[:0]
			continue
		}
		// Import the batch and reset the buffer
		if _, err := api.bchain.BlockChain().InsertChain(blocks); err != nil {
			return false, fmt.Errorf("batch %d: failed to insert: %v", batch, err)
		}
		blocks = blocks[:0]
	}
	return true, nil
}


// StorageRangeResult is the result of a debug_storageRangeAt API call.
type StorageRangeResult struct {
	Storage storageMap  `json:"storage"`
	NextKey *types.Hash `json:"nextKey"` // nil if Storage includes the last key in the trie.
}

type storageMap map[types.Hash]storageEntry

type storageEntry struct {
	Key   *types.Hash `json:"key"`
	Value types.Hash  `json:"value"`
}



