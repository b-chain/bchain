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
// @File: interface.go
// @Date: 2018/03/19 09:54:19
////////////////////////////////////////////////////////////////////////////////

package database

type IDatabaseGetter interface {
	Get(key []byte) ([]byte, error)
}

type IDatabasePutter interface {
	Put(key []byte, value []byte) error
}

// Database wraps all database operations. All methods are safe for concurrent use.
type IDatabase interface {
	IDatabaseGetter
	IDatabasePutter
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	Close()
	NewBatch() IBatch
}

// Batch is a write-only database that commits changes to its host database
// when Write is called. Batch cannot be used concurrently.
type IBatch interface {
	IDatabasePutter
	ValueSize() int 	// amount of data in the batch
	Write() error
	Reset() 			// Reset resets the batch for reuse
}
