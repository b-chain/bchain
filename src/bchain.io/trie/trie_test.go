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
// @File: trie_test.go
// @Date: 2018/05/07 10:38:07
////////////////////////////////////////////////////////////////////////////////

package trie

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/davecgh/go-spew/spew"

	"github.com/tinylib/msgp/msgp"
	"bchain.io/common"
	"bchain.io/common/types"
	"bchain.io/utils/crypto"
	"bchain.io/utils/database"
)

func init() {
	spew.Config.Indent = "    "
	spew.Config.DisableMethods = false
}

// Used for testing
func newEmpty() *Trie {
	db, _ := database.OpenMemDB()
	trie, _ := New(types.Hash{}, db)
	return trie
}

func TestEmptyTrie(t *testing.T) {
	var trie Trie
	res := trie.Hash()
	exp := emptyRoot
	if res != types.Hash(exp) {
		t.Errorf("expected %x got %x", exp, res)
	}
}

func TestNull(t *testing.T) {
	var trie Trie
	key := make([]byte, 32)
	value := []byte("test")
	trie.Update(key, value)
	if !bytes.Equal(trie.Get(key), value) {
		t.Fatal("wrong value")
	}
}

func TestMissingRoot(t *testing.T) {
	db, _ := database.OpenMemDB()
	trie, err := New(types.HexToHash("0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"), db)
	if trie != nil {
		t.Error("New returned non-nil trie for invalid root")
	}
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("New returned wrong error: %v", err)
	}
}

func TestMissingNode(t *testing.T) {
	db, _ := database.OpenMemDB()

	trie, _ := New(types.Hash{}, db)
	updateString(trie, "120000", "qwerqwerqwerqwerqwerqwerqwerqwer")
	updateString(trie, "123456", "asdfasdfasdfasdfasdfasdfasdfasdf")
	root, _ := trie.Commit()

	trie, _ = New(root, db)
	_, err := trie.TryGet([]byte("120000"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	trie, _ = New(root, db)
	_, err = trie.TryGet([]byte("120099"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	trie, _ = New(root, db)
	_, err = trie.TryGet([]byte("123456"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	trie, _ = New(root, db)
	err = trie.TryUpdate([]byte("120099"), []byte("zxcvzxcvzxcvzxcvzxcvzxcvzxcvzxcv"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	trie, _ = New(root, db)
	err = trie.TryDelete([]byte("123456"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	db.Delete(types.FromHex("e04ce98667e4bb2f315ddf35cd3db821e2b7ec04687a66fe4384771753fd2780"))

	trie, _ = New(root, db)
	_, err = trie.TryGet([]byte("120000"))
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("Wrong error: %v", err)
	}

	trie, _ = New(root, db)
	_, err = trie.TryGet([]byte("120099"))
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("Wrong error: %v", err)
	}

	trie, _ = New(root, db)
	_, err = trie.TryGet([]byte("123456"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	trie, _ = New(root, db)
	err = trie.TryUpdate([]byte("120099"), []byte("zxcv"))
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("Wrong error: %v", err)
	}

	trie, _ = New(root, db)
	err = trie.TryDelete([]byte("123456"))
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("Wrong error: %v", err)
	}
}

func TestInsert(t *testing.T) {
	trie := newEmpty()

	updateString(trie, "doe", "reindeer")
	updateString(trie, "dog", "puppy")
	updateString(trie, "dogglesworth", "cat")

	exp := types.HexToHash("93a32dffc094f65f2a0aa145c481f07760bb325f0b3fd05330a21f6fc6307ea5")
	root := trie.Hash()
	if root != exp {
		t.Errorf("exp %x got %x", exp, root)
	}

	trie = newEmpty()
	updateString(trie, "A", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	exp = types.HexToHash("ebd78bdf3820096cb028d432552907affd60b333be3d142ebd115106ec471e68")
	root, err := trie.Commit()
	if err != nil {
		t.Fatalf("commit error: %v", err)
	}
	if root != exp {
		t.Errorf("exp %x got %x", exp, root)
	}
}

func TestGet(t *testing.T) {
	trie := newEmpty()
	updateString(trie, "doe", "reindeer")
	updateString(trie, "dog", "puppy")
	updateString(trie, "dogglesworth", "cat")

	for i := 0; i < 100; i++ {
		res := getString(trie, "dog")
		if !bytes.Equal(res, []byte("puppy")) {
			t.Errorf("expected puppy got %x", res)
		}

		unknown := getString(trie, "unknown")
		if unknown != nil {
			t.Errorf("expected nil got %x", unknown)
		}

		/*if i == 1 {
			return
		}*/

		if i == 1 {
			deleteString(trie, "dogglesworth")
			//return
		}

		if i == 2 {
			deleteString(trie, "doe")
			//return
		}

		if i == 5 {
			deleteString(trie, "d")
		}

		if i == 6 {
			updateString(trie, "dogg", "puppy")
		}

		if i == 7 {
			trie.Commit()
		}

		//trie.Commit()
	}
}

func TestDelete(t *testing.T) {
	trie := newEmpty()
	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"bchain", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"bchain", ""},
		{"dog", "puppy"},
		{"shaman", ""},
	}
	for _, val := range vals {
		if val.v != "" {
			updateString(trie, val.k, val.v)
		} else {
			deleteString(trie, val.k)
		}
	}

	hash := trie.Hash()
	exp := types.HexToHash("bf432ccc90de7fea74124337d31b105640a23b541e64f2dc005d2ae15ceaf769")
	if hash != exp {
		t.Errorf("expected %x got %x", exp, hash)
	}
}

func TestEmptyValues(t *testing.T) {
	trie := newEmpty()

	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"bchain", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"bchain", ""},
		{"dog", "puppy"},
		{"shaman", ""},
	}
	for _, val := range vals {
		updateString(trie, val.k, val.v)
	}

	hash := trie.Hash()
	exp := types.HexToHash("bf432ccc90de7fea74124337d31b105640a23b541e64f2dc005d2ae15ceaf769")
	if hash != exp {
		t.Errorf("expected %x got %x", exp, hash)
	}
}

func TestReplication(t *testing.T) {
	trie := newEmpty()
	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"bchain", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"dog", "puppy"},
		{"somethingveryoddindeedthis is", "myothernodedata"},
	}
	for _, val := range vals {
		updateString(trie, val.k, val.v)
	}
	exp, err := trie.Commit()
	if err != nil {
		t.Fatalf("commit error: %v", err)
	}

	// create a new trie on top of the database and check that lookups work.
	trie2, err := New(exp, trie.db)
	if err != nil {
		t.Fatalf("can't recreate trie at %x: %v", exp, err)
	}
	for _, kv := range vals {
		if string(getString(trie2, kv.k)) != kv.v {
			t.Errorf("trie2 doesn't have %q => %q", kv.k, kv.v)
		}
	}
	hash, err := trie2.Commit()
	if err != nil {
		t.Fatalf("commit error: %v", err)
	}
	if hash != exp {
		t.Errorf("root failure. expected %x got %x", exp, hash)
	}

	// perform some insertions on the new trie.
	vals2 := []struct{ k, v string }{
		{"do", "verb"},
		{"bchain", "wookiedoo"},
		{"horse", "stallion"},
		// {"shaman", "horse"},
		// {"doge", "coin"},
		// {"bchain", ""},
		// {"dog", "puppy"},
		// {"somethingveryoddindeedthis is", "myothernodedata"},
		// {"shaman", ""},
	}
	for _, val := range vals2 {
		updateString(trie2, val.k, val.v)
	}
	if hash := trie2.Hash(); hash != exp {
		t.Errorf("root failure. expected %x got %x", exp, hash)
	}
}

func TestLargeValue(t *testing.T) {
	trie := newEmpty()
	trie.Update([]byte("key1"), []byte{99, 99, 99, 99})
	trie.Update([]byte("key2"), bytes.Repeat([]byte{1}, 32))
	trie.Hash()
}

type countingDB struct {
	Database
	gets map[string]int
}

func (db *countingDB) Get(key []byte) ([]byte, error) {
	db.gets[string(key)]++
	return db.Database.Get(key)
}

// TestCacheUnload checks that decoded nodes are unloaded after a
// certain number of commit operations.
func TestCacheUnload(t *testing.T) {
	// Create test trie with two branches.
	trie := newEmpty()
	key1 := "---------------------------------"
	key2 := "---some other branch"
	updateString(trie, key1, "this is the branch of key1.")
	updateString(trie, key2, "this is the branch of key2.")
	root, _ := trie.Commit()

	// Commit the trie repeatedly and access key1.
	// The branch containing it is loaded from DB exactly two times:
	// in the 0th and 6th iteration.
	db := &countingDB{Database: trie.db, gets: make(map[string]int)}
	trie, _ = New(root, db)
	trie.SetCacheLimit(5)
	for i := 0; i < 12; i++ {
		getString(trie, key1)
		trie.Commit()
	}

	// Check that it got loaded two times.
	for dbkey, count := range db.gets {
		if count != 2 {
			t.Errorf("db key %x loaded %d times, want %d times", []byte(dbkey), count, 2)
		}
	}
}

// randTest performs random trie operations.
// Instances of this test are created by Generate.
type randTest []randTestStep

type randTestStep struct {
	op    int
	key   []byte // for opUpdate, opDelete, opGet
	value []byte // for opUpdate
	err   error  // for debugging
}

const (
	opUpdate = iota
	opDelete
	opGet
	opCommit
	opHash
	opReset
	opItercheckhash
	opCheckCacheInvariant
	opMax // boundary value, not an actual op
)

func (randTest) Generate(r *rand.Rand, size int) reflect.Value {
	var allKeys [][]byte
	genKey := func() []byte {
		if len(allKeys) < 2 || r.Intn(100) < 10 {
			// new key
			key := make([]byte, r.Intn(50))
			r.Read(key)
			allKeys = append(allKeys, key)
			return key
		}
		// use existing key
		return allKeys[r.Intn(len(allKeys))]
	}

	var steps randTest
	for i := 0; i < size; i++ {
		step := randTestStep{op: r.Intn(opMax)}
		//step := randTestStep{op: r.Intn(4)}
		switch step.op {
		case opUpdate:
			step.key = genKey()
			step.value = make([]byte, 8)
			binary.BigEndian.PutUint64(step.value, uint64(i))
		case opGet, opDelete:
			step.key = genKey()
		}
		steps = append(steps, step)
	}
	return reflect.ValueOf(steps)
}

func runRandTest(rt randTest) bool {
	db, _ := database.OpenMemDB()
	tr, _ := New(types.Hash{}, db)
	values := make(map[string]string) // tracks content of the trie

	for i, step := range rt {

		myTest := func() {
			for k, v := range values {
				val := tr.Get([]byte(k))
				if string(val) != v {
					fmt.Printf("ERROR - (%x, %x), %x\n", k, v, val)
				}
			}

			return
		}

		fmt.Printf("%d step: %d\n", i, step.op)
		myTest()

		switch step.op {
		case opUpdate:
			tr.Update(step.key, step.value)
			values[string(step.key)] = string(step.value)
		case opDelete:
			tr.Delete(step.key)
			delete(values, string(step.key))
		case opGet:
			v := tr.Get(step.key)
			want := values[string(step.key)]
			fmt.Printf("v, want: %x, %x\n", v, want)
			if string(v) != want {
				rt[i].err = fmt.Errorf("mismatch for key 0x%x, got 0x%x want 0x%x", step.key, v, want)
			}
		case opCommit:
			_, rt[i].err = tr.Commit()
		case opHash:
			tr.Hash()
		case opReset:
			hash, err := tr.Commit()
			if err != nil {
				rt[i].err = err
				return false
			}
			newtr, err := New(hash, db)
			if err != nil {
				rt[i].err = err
				return false
			}
			tr = newtr
		case opItercheckhash:
			checktr, _ := New(types.Hash{}, nil)
			it := NewIterator(tr.NodeIterator(nil))
			for it.Next() {
				checktr.Update(it.Key, it.Value)
			}
			if tr.Hash() != checktr.Hash() {
				rt[i].err = fmt.Errorf("hash mismatch in opItercheckhash")
			}
		case opCheckCacheInvariant:
			rt[i].err = checkCacheInvariant(tr.root, NodeIntf{nil}, tr.cachegen, false, 0)
		}

		fmt.Printf("========  ")
		myTest()
		fmt.Printf("\n")

		// Abort the test on error.
		if rt[i].err != nil {
			return false
		}
	}
	return true
}

func checkCacheInvariant(n, parent NodeIntf, parentCachegen uint16, parentDirty bool, depth int) error {
	var children []NodeIntf
	var flag nodeFlag
	switch n := n.Node.(type) {
	case *ShortNode:
		flag = n.flags
		children = []NodeIntf{n.Val}
	case *FullNode:
		flag = n.flags
		children = n.Children[:]
	default:
		return nil
	}

	errorf := func(format string, args ...interface{}) error {
		msg := fmt.Sprintf(format, args...)
		msg += fmt.Sprintf("\nat depth %d node %s", depth, spew.Sdump(n))
		msg += fmt.Sprintf("parent: %s", spew.Sdump(parent))
		return errors.New(msg)
	}
	if flag.gen > parentCachegen {
		return errorf("cache invariant violation: %d > %d\n", flag.gen, parentCachegen)
	}
	if depth > 0 && !parentDirty && flag.dirty {
		return errorf("cache invariant violation: %d > %d\n", flag.gen, parentCachegen)
	}
	for _, child := range children {
		if err := checkCacheInvariant(child, n, flag.gen, flag.dirty, depth+1); err != nil {
			return err
		}
	}
	return nil
}

func TestRandom(t *testing.T) {
	if err := quick.Check(runRandTest, nil); err != nil {
		if cerr, ok := err.(*quick.CheckError); ok {
			t.Fatalf("random test iteration %d failed: %s", cerr.Count, spew.Sdump(cerr.In))
		}
		t.Fatal(err)
	}
}

func TestRandomCustom(t *testing.T) {
	rt := randTest{
		randTestStep{
			op: 2,
			key: []byte{
				0xea, 0x0d, 0xad, 0xa8, 0xb4, 0x56, 0x1a, 0x2a, 0x70, 0x5d, 0xdd, 0x3d, 0xa8,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 0,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
			},
			err: nil,
		},
		randTestStep{
			op:    3,
			key:   nil,
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 0,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
			},
			err: nil,
		},
		randTestStep{
			op: 2,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op:    3,
			key:   nil,
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 1,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 0,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07,
			},
			err: nil,
		},
		randTestStep{
			op: 2,
			key: []byte{
				0xea, 0x0d, 0xad, 0xa8, 0xb4, 0x56, 0x1a, 0x2a, 0x70, 0x5d, 0xdd, 0x3d, 0xa8,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op:    3,
			key:   nil,
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 1,
			key: []byte{
				0xea, 0x0d, 0xad, 0xa8, 0xb4, 0x56, 0x1a, 0x2a, 0x70, 0x5d, 0xdd, 0x3d, 0xa8,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 2,
			key: []byte{
				0xea, 0x0d, 0xad, 0xa8, 0xb4, 0x56, 0x1a, 0x2a, 0x70, 0x5d, 0xdd, 0x3d, 0xa8,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 2,
			key: []byte{
				0xea, 0x0d, 0xad, 0xa8, 0xb4, 0x56, 0x1a, 0x2a, 0x70, 0x5d, 0xdd, 0x3d, 0xa8,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 2,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 0,
			key: []byte{
				0x26, 0xbb, 0x26, 0x02, 0xbb, 0xfd, 0xcf, 0x25, 0xb3, 0xb7, 0xae, 0x8d, 0x7a, 0x7e, 0xd6, 0x96,
				0xc7, 0x3c, 0x60, 0x5e, 0x5a, 0xfe, 0xc7, 0x07, 0x5d, 0x53, 0x8b, 0x47, 0x15, 0xee, 0xd8, 0x41,
				0x9e, 0x08, 0x2c, 0xa8, 0xe9, 0x30, 0xf5, 0xb3, 0x3f, 0xc4,
			},
			value: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0e,
			},
			err: nil,
		},
		randTestStep{
			op:    3,
			key:   nil,
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 0,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10,
			},
			err: nil,
		},
		randTestStep{
			op: 2,
			key: []byte{
				0xea, 0x0d, 0xad, 0xa8, 0xb4, 0x56, 0x1a, 0x2a, 0x70, 0x5d, 0xdd, 0x3d, 0xa8,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 0,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12,
			},
			err: nil,
		},
		randTestStep{
			op: 1,
			key: []byte{
				0x26, 0xbb, 0x26, 0x02, 0xbb, 0xfd, 0xcf, 0x25, 0xb3, 0xb7, 0xae, 0x8d, 0x7a, 0x7e, 0xd6, 0x96,
				0xc7, 0x3c, 0x60, 0x5e, 0x5a, 0xfe, 0xc7, 0x07, 0x5d, 0x53, 0x8b, 0x47, 0x15, 0xee, 0xd8, 0x41,
				0x9e, 0x08, 0x2c, 0xa8, 0xe9, 0x30, 0xf5, 0xb3, 0x3f, 0xc4,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 1,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 0,
			key: []byte{
				0xc1, 0xaa,
			},
			value: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x15,
			},
			err: nil,
		},
		randTestStep{
			op: 1,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 0,
			key: []byte{
				0xea, 0x0d, 0xad, 0xa8, 0xb4, 0x56, 0x1a, 0x2a, 0x70, 0x5d, 0xdd, 0x3d, 0xa8,
			},
			value: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x17,
			},
			err: nil,
		},
		randTestStep{
			op: 2,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 2,
			key: []byte{
				0x26, 0xbb, 0x26, 0x02, 0xbb, 0xfd, 0xcf, 0x25, 0xb3, 0xb7, 0xae, 0x8d, 0x7a, 0x7e, 0xd6, 0x96,
				0xc7, 0x3c, 0x60, 0x5e, 0x5a, 0xfe, 0xc7, 0x07, 0x5d, 0x53, 0x8b, 0x47, 0x15, 0xee, 0xd8, 0x41,
				0x9e, 0x08, 0x2c, 0xa8, 0xe9, 0x30, 0xf5, 0xb3, 0x3f, 0xc4,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op:    3,
			key:   nil,
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 2,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 0,
			key: []byte{
				0xc1, 0xaa,
			},
			value: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1c,
			},
			err: nil,
		},
		randTestStep{
			op:    3,
			key:   nil,
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 1,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op:    3,
			key:   nil,
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 0,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20,
			},
			err: nil,
		},
		randTestStep{
			op: 2,
			key: []byte{
				0x7b, 0x22, 0xd8, 0x93, 0x47, 0x4a, 0xa4, 0x94, 0x4e, 0x67, 0xff, 0xb8, 0xe3, 0xe0, 0x94, 0x50,
				0x0a, 0x4c, 0xa6, 0x24, 0x4d, 0xac, 0xa1, 0xbb, 0x97, 0xf4, 0xc1,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op:    3,
			key:   nil,
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 1,
			key: []byte{
				0x26, 0xbb, 0x26, 0x02, 0xbb, 0xfd, 0xcf, 0x25, 0xb3, 0xb7, 0xae, 0x8d, 0x7a, 0x7e, 0xd6, 0x96,
				0xc7, 0x3c, 0x60, 0x5e, 0x5a, 0xfe, 0xc7, 0x07, 0x5d, 0x53, 0x8b, 0x47, 0x15, 0xee, 0xd8, 0x41,
				0x9e, 0x08, 0x2c, 0xa8, 0xe9, 0x30, 0xf5, 0xb3, 0x3f, 0xc4,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 2,
			key: []byte{
				0xea, 0x0d, 0xad, 0xa8, 0xb4, 0x56, 0x1a, 0x2a, 0x70, 0x5d, 0xdd, 0x3d, 0xa8,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 0,
			key: []byte{
				0x26, 0xbb, 0x26, 0x02, 0xbb, 0xfd, 0xcf, 0x25, 0xb3, 0xb7, 0xae, 0x8d, 0x7a, 0x7e, 0xd6, 0x96,
				0xc7, 0x3c, 0x60, 0x5e, 0x5a, 0xfe, 0xc7, 0x07, 0x5d, 0x53, 0x8b, 0x47, 0x15, 0xee, 0xd8, 0x41,
				0x9e, 0x08, 0x2c, 0xa8, 0xe9, 0x30, 0xf5, 0xb3, 0x3f, 0xc4,
			},
			value: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x25,
			},
			err: nil,
		},
		randTestStep{
			op: 2,
			key: []byte{
				0x26, 0xbb, 0x26, 0x02, 0xbb, 0xfd, 0xcf, 0x25, 0xb3, 0xb7, 0xae, 0x8d, 0x7a, 0x7e, 0xd6, 0x96,
				0xc7, 0x3c, 0x60, 0x5e, 0x5a, 0xfe, 0xc7, 0x07, 0x5d, 0x53, 0x8b, 0x47, 0x15, 0xee, 0xd8, 0x41,
				0x9e, 0x08, 0x2c, 0xa8, 0xe9, 0x30, 0xf5, 0xb3, 0x3f, 0xc4,
			},
			value: nil,
			err:   nil,
		},
		randTestStep{
			op: 2,
			key: []byte{
				0xc1, 0xaa,
			},
			value: nil,
			err:   nil,
		},
	}

	if !runRandTest(rt) {
		t.Fatalf("runRandTest failded.\n")
	}
}

func BenchmarkGet(b *testing.B)      { benchGet(b, false) }
func BenchmarkGetDB(b *testing.B)    { benchGet(b, true) }
func BenchmarkUpdateBE(b *testing.B) { benchUpdate(b, binary.BigEndian) }
func BenchmarkUpdateLE(b *testing.B) { benchUpdate(b, binary.LittleEndian) }

const benchElemCount = 20000

func benchGet(b *testing.B, commit bool) {
	trie := new(Trie)
	if commit {
		_, tmpdb := tempDB()
		trie, _ = New(types.Hash{}, tmpdb)
	}
	k := make([]byte, 32)
	for i := 0; i < benchElemCount; i++ {
		binary.LittleEndian.PutUint64(k, uint64(i))
		trie.Update(k, k)
	}
	binary.LittleEndian.PutUint64(k, benchElemCount/2)
	if commit {
		trie.Commit()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Get(k)
	}
	b.StopTimer()

	if commit {
		ldb := trie.db.(*database.LDatabase)
		ldb.Close()
		os.RemoveAll(ldb.Path())
	}
}

func benchUpdate(b *testing.B, e binary.ByteOrder) *Trie {
	trie := newEmpty()
	k := make([]byte, 32)
	for i := 0; i < b.N; i++ {
		e.PutUint64(k, uint64(i))
		trie.Update(k, k)
	}
	return trie
}

// Benchmarks the trie hashing. Since the trie caches the result of any operation,
// we cannot use b.N as the number of hashing rouns, since all rounds apart from
// the first one will be NOOP. As such, we'll use b.N as the number of account to
// insert into the trie before measuring the hashing.
func BenchmarkHash(b *testing.B) {
	// Make the random benchmark deterministic
	random := rand.New(rand.NewSource(0))

	// Create a realistic account trie to hash
	addresses := make([][20]byte, b.N)
	for i := 0; i < len(addresses); i++ {
		for j := 0; j < len(addresses[i]); j++ {
			addresses[i][j] = byte(random.Intn(256))
		}
	}
	accounts := make([][]byte, len(addresses))
	for i := 0; i < len(accounts); i++ {
		var (
			nonce   = uint64(random.Int63())
			balance = new(big.Int).Rand(random, new(big.Int).Exp(common.Big2, common.Big256, nil))
			root    = emptyRoot
			code    = crypto.Keccak256(nil)
		)

		buf := bytes.Buffer{}
		wr := msgp.NewWriter(&buf)
		wr.WriteIntf([]interface{}{nonce, balance, root, code})
		wr.Flush()
		accounts[i] = buf.Bytes()
	}
	// Insert the accounts into the trie and hash it
	trie := newEmpty()
	for i := 0; i < len(addresses); i++ {
		trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
	}
	b.ResetTimer()
	b.ReportAllocs()
	trie.Hash()
}

func tempDB() (string, Database) {
	dir, err := ioutil.TempDir("", "trie-bench")
	if err != nil {
		panic(fmt.Sprintf("can't create temporary directory: %v", err))
	}
	db, err := database.OpenLDB(dir, 256, 0)
	if err != nil {
		panic(fmt.Sprintf("can't create temporary database: %v", err))
	}
	return dir, db
}

func getString(trie *Trie, k string) []byte {
	return trie.Get([]byte(k))
}

func updateString(trie *Trie, k, v string) {
	trie.Update([]byte(k), []byte(v))
}

func deleteString(trie *Trie, k string) {
	trie.Delete([]byte(k))
}
