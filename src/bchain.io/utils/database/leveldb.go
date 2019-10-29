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
// @File: leveldb.go
// @Date: 2018/05/07 09:30:07
////////////////////////////////////////////////////////////////////////////////

package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
	"bchain.io/utils/metrics"
	"bchain.io/log"
	"fmt"
	"os"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"time"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"strings"
	"strconv"
)

// database struct
type LDatabase struct {
	fileName string
	db   *leveldb.DB // LevelDB instance

	getTimer       metrics.Timer // Timer for measuring the database get request counts and latencies
	putTimer       metrics.Timer // Timer for measuring the database put request counts and latencies
	delTimer       metrics.Timer // Timer for measuring the database delete request counts and latencies
	missMeter      metrics.Meter // Meter for measuring the missed database get requests
	readMeter      metrics.Meter // Meter for measuring the database get request data usage
	writeMeter     metrics.Meter // Meter for measuring the database put request data usage
	compTimeMeter  metrics.Meter // Meter for measuring the total time spent in database compaction
	compReadMeter  metrics.Meter // Meter for measuring the data read during compaction
	compWriteMeter metrics.Meter // Meter for measuring the data written during compaction

	quitLock sync.Mutex      // Mutex protecting the quit channel access
	quitChan chan chan error // Quit channel to stop the metrics collection before closing the database
}

var (
	logTag = "utils.database"
	logger log.Logger
)



func init() {
	logger = log.GetLogger(logTag)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", logTag)
		os.Exit(1)
	}
}

// Open returns a LevelDB wrapped object.
func OpenLDB(file string, blockCache int, fileCache int) (*LDatabase, error) {
	// Ensure we have some minimal caching and file guarantees
	if fileCache < 16 {
		fileCache = 16
	}
	if blockCache < 16 {
		blockCache = 16
	}
	logger.Info("Allocated cache and file handles", "blockCache", blockCache, "fileCache", fileCache)

	// Open the db and recover any potential corruptions
	db, err := leveldb.OpenFile(file, &opt.Options{
		OpenFilesCacheCapacity: fileCache,
		BlockCacheCapacity:     blockCache / 2 * opt.MiB,
		WriteBuffer:            blockCache / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(file, nil)
	}
	// (Re)check for errors and abort if opening of the db failed
	if err != nil {
		return nil, err
	}
	return &LDatabase{
		fileName:  file,
		db:  db,
	}, nil
}

func (db *LDatabase) Path() string {
	return db.fileName
}

func (db *LDatabase) Put(key []byte, value []byte) error {
	// Measure the database put latency, if requested
	if db.putTimer != nil {
		defer db.putTimer.UpdateSince(time.Now())
	}

	if db.writeMeter != nil {
		db.writeMeter.Mark(int64(len(value)))
	}
	return db.db.Put(key, value, nil)
}

func (db *LDatabase) Get(key []byte) ([]byte, error) {
	// Measure the database get latency, if requested
	if db.getTimer != nil {
		defer db.getTimer.UpdateSince(time.Now())
	}
	// Retrieve the key and increment the miss counter if not found
	dat, err := db.db.Get(key, nil)
	if err != nil {
		if db.missMeter != nil {
			db.missMeter.Mark(1)
		}
		return nil, err
	}
	// Otherwise update the actually retrieved amount of data
	if db.readMeter != nil {
		db.readMeter.Mark(int64(len(dat)))
	}
	return dat, nil
}

func (db *LDatabase) Has(key []byte) (bool, error) {
	return db.db.Has(key, nil)
}

func (db *LDatabase) Delete(key []byte) error {
	// Measure the database delete latency, if requested
	if db.delTimer != nil {
		defer db.delTimer.UpdateSince(time.Now())
	}
	// Execute the actual operation
	return db.db.Delete(key, nil)
}

func (db *LDatabase) NewIterator() iterator.Iterator {
	return db.db.NewIterator(nil, nil)
}

func (db *LDatabase) Close() {
	// Stop the metrics collection to avoid internal database races
	db.quitLock.Lock()
	defer db.quitLock.Unlock()

	if db.quitChan != nil {
		errc := make(chan error)
		db.quitChan <- errc
		if err := <-errc; err != nil {
			logger.Error("Metrics collection failed", "err", err)
		}
	}
	err := db.db.Close()
	if err == nil {
		logger.Info("Database closed")
	} else {
		logger.Error("Failed to close database", "err", err)
	}
}

func (db *LDatabase) NewBatch() IBatch {
	return &ldbBatch{db: db.db, b: new(leveldb.Batch)}
}

// batch db struct
type ldbBatch struct {
	db   *leveldb.DB
	b    *leveldb.Batch
	size int
}

func (batch *ldbBatch) Put(key []byte, value []byte) error {
	batch.b.Put(key, value)
	batch.size += len(value)
	return nil
}

func (batch *ldbBatch) ValueSize() int {
	return batch.size
}

func (batch *ldbBatch) Write() error {
	return batch.db.Write(batch.b, nil)
}

func (batch *ldbBatch) Reset() {
	batch.b.Reset()
	batch.size = 0
}


// Meter configures the database metrics collectors and
func (db *LDatabase) Meter(prefix string) {
	// Short circuit metering if the metrics system is disabled
	if !metrics.Enabled {
		return
	}
	// Initialize all the metrics collector at the requested prefix
	db.getTimer = metrics.NewRegisteredTimer(prefix+"user/gets", nil)
	db.putTimer = metrics.NewRegisteredTimer(prefix+"user/puts", nil)
	db.delTimer = metrics.NewRegisteredTimer(prefix+"user/dels", nil)
	db.missMeter = metrics.NewRegisteredMeter(prefix+"user/misses", nil)
	db.readMeter = metrics.NewRegisteredMeter(prefix+"user/reads", nil)
	db.writeMeter = metrics.NewRegisteredMeter(prefix+"user/writes", nil)
	db.compTimeMeter = metrics.NewRegisteredMeter(prefix+"compact/time", nil)
	db.compReadMeter = metrics.NewRegisteredMeter(prefix+"compact/input", nil)
	db.compWriteMeter = metrics.NewRegisteredMeter(prefix+"compact/output", nil)

	// Create a quit channel for the periodic collector and run it
	db.quitLock.Lock()
	db.quitChan = make(chan chan error)
	db.quitLock.Unlock()

	go db.meter(3 * time.Second)
}

// meter periodically retrieves internal leveldb counters and reports them to
// the metrics subsystem.
//
// This is how a stats table look like (currently):
//   Compactions
//    Level |   Tables   |    Size(MB)   |    Time(sec)  |    Read(MB)   |   Write(MB)
//   -------+------------+---------------+---------------+---------------+---------------
//      0   |          0 |       0.00000 |       1.27969 |       0.00000 |      12.31098
//      1   |         85 |     109.27913 |      28.09293 |     213.92493 |     214.26294
//      2   |        523 |    1000.37159 |       7.26059 |      66.86342 |      66.77884
//      3   |        570 |    1113.18458 |       0.00000 |       0.00000 |       0.00000
func (db *LDatabase) meter(refresh time.Duration) {
	// Create the counters to store current and previous values
	counters := make([][]float64, 2)
	for i := 0; i < 2; i++ {
		counters[i] = make([]float64, 3)
	}
	// Iterate ad infinitum and collect the stats
	for i := 1; ; i++ {
		// Retrieve the database stats
		stats, err := db.db.GetProperty("leveldb.stats")
		if err != nil {
			logger.Error("Failed to read database stats", "err", err)
			return
		}
		// Find the compaction table, skip the header
		lines := strings.Split(stats, "\n")
		for len(lines) > 0 && strings.TrimSpace(lines[0]) != "Compactions" {
			lines = lines[1:]
		}
		if len(lines) <= 3 {
			logger.Error("Compaction table not found")
			return
		}
		lines = lines[3:]

		// Iterate over all the table rows, and accumulate the entries
		for j := 0; j < len(counters[i%2]); j++ {
			counters[i%2][j] = 0
		}
		for _, line := range lines {
			parts := strings.Split(line, "|")
			if len(parts) != 6 {
				break
			}
			for idx, counter := range parts[3:] {
				value, err := strconv.ParseFloat(strings.TrimSpace(counter), 64)
				if err != nil {
					logger.Error("Compaction entry parsing failed", "err", err)
					return
				}
				counters[i%2][idx] += value
			}
		}
		// Update all the requested meters
		if db.compTimeMeter != nil {
			db.compTimeMeter.Mark(int64((counters[i%2][0] - counters[(i-1)%2][0]) * 1000 * 1000 * 1000))
		}
		if db.compReadMeter != nil {
			db.compReadMeter.Mark(int64((counters[i%2][1] - counters[(i-1)%2][1]) * 1024 * 1024))
		}
		if db.compWriteMeter != nil {
			db.compWriteMeter.Mark(int64((counters[i%2][2] - counters[(i-1)%2][2]) * 1024 * 1024))
		}
		// Sleep a bit, then repeat the stats collection
		select {
		case errc := <-db.quitChan:
			// Quit requesting, stop hammering the database
			errc <- nil
			return

		case <-time.After(refresh):
			// Timeout, gather a new set of stats
		}
	}
}
