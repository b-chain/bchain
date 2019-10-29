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
// @File: db_upgrade.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

// Package bchain implements the bchain protocol.
//go:generate msgp
package bchain

import (
	"bytes"
	"time"

	"bchain.io/common/types"
	"bchain.io/utils/database"
	"github.com/tinylib/msgp/msgp"
	"bchain.io/core/blockchain"
)

type Meta struct {
	BlockHash  types.Hash
	BlockIndex uint64
	Index      uint64
}

var deduplicateData = []byte("dbUpgrade_20170714deduplicateData")
func UpgradeDeduplicateData(db database.IDatabase)func()error{
	return upgradeDeduplicateData(db)
}
// upgradeDeduplicateData checks the chain database version and
// starts a background process to make upgrades if necessary.
// Returns a stop function that blocks until the process has
// been safely stopped.
func upgradeDeduplicateData(db database.IDatabase) func() error {
	// If the database is already converted or empty, bail out
	data, _ := db.Get(deduplicateData)
	if len(data) > 0 && data[0] == 42 {
		return nil
	}
	if data, _ := db.Get([]byte("LastHeader")); len(data) == 0 {
		db.Put(deduplicateData, []byte{42})
		return nil
	}
	// Start the deduplication upgrade on a new goroutine
	logger.Warn("Upgrading database to use lookup entries")
	stop := make(chan chan error)

	go func() {
		// Create an iterator to read the entire database and covert old lookup entires
		it := db.(*database.LDatabase).NewIterator()
		defer func() {
			if it != nil {
				it.Release()
			}
		}()

		var (
			converted uint64
			failed    error
		)
		for failed == nil && it.Next() {
			// Skip any entries that don't look like old transaction meta entires (<hash>0x01)
			key := it.Key()
			if len(key) != types.HashLength+1 || key[types.HashLength] != 0x01 {
				continue
			}
			// Skip any entries that don't contain metadata (name clash between <hash>0x01 and <some-prefix><hash>)
			var meta Meta
			rd := bytes.NewReader(it.Value())

			if err := msgp.Decode(rd , &meta); err != nil {
				continue
			}
			// Skip any already upgraded entries (clash due to <hash> ending with 0x01 (old suffix))
			hash := key[:types.HashLength]

			if hash[0] == byte('l') {
				// Potential clash, the "old" `hash` must point to a live transaction.
				if tx, _, _, _ := blockchain.GetTransaction(db, types.BytesToHash(hash)); tx == nil || !bytes.Equal(tx.Hash().Bytes(), hash) {
					continue
				}
			}
			// Convert the old metadata to a new lookup entry, delete duplicate data
			if failed = db.Put(append([]byte("l"), hash...), it.Value()); failed == nil { // Write the new looku entry
				if failed = db.Delete(hash); failed == nil { // Delete the duplicate transaction data
					if failed = db.Delete(append([]byte("receipts-"), hash...)); failed == nil { // Delete the duplicate receipt data
						if failed = db.Delete(key); failed != nil { // Delete the old transaction metadata
							break
						}
					}
				}
			}
			// Bump the conversion counter, and recreate the iterator occasionally to
			// avoid too high memory consumption.
			converted++
			if converted%100000 == 0 {
				it.Release()
				it = db.(*database.LDatabase).NewIterator()
				it.Seek(key)

				logger.Info("Deduplicating database entries", "deduped", converted)
			}
			// Check for termination, or continue after a bit of a timeout
			select {
			case errc := <-stop:
				errc <- nil
				return
			case <-time.After(time.Microsecond * 100):
			}
		}
		// Upgrade finished, mark a such and terminate
		if failed == nil {
			logger.Info("Database deduplication successful", "deduped", converted)
			db.Put(deduplicateData, []byte{42})
		} else {
			logger.Error("Database deduplication failed", "deduped", converted, "err", failed)
		}
		it.Release()
		it = nil

		errc := <-stop
		errc <- failed
	}()
	// Assembly the cancellation callback
	return func() error {
		errc := make(chan error)
		stop <- errc
		return <-errc
	}
}
