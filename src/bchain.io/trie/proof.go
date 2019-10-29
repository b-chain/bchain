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
// @File: proof.go
// @Date: 2018/05/04 11:03:04
////////////////////////////////////////////////////////////////////////////////

package trie

import (
	"bytes"
	"fmt"

	"github.com/tinylib/msgp/msgp"
	"bchain.io/utils/crypto"
	"bchain.io/common/types"
)

// Prove constructs a merkle proof for key. The result contains all
// encoded nodes on the path to the value at key. The value itself is
// also included in the last node and can be retrieved by verifying
// the proof.
//
// If the trie does not contain a value for key, the returned proof
// contains all nodes of the longest existing prefix of the key
// (at least the root node), ending with the node that proves the
// absence of the key.
func (t *Trie) Prove(key []byte, fromLevel uint, proofDb DatabaseWriter) error {
	// Collect all nodes on the path to key.
	key = keybytesToHex(key)
	nodes := []NodeIntf{}
	tn := t.root
	for len(key) > 0 && tn.Node != nil {
		switch n := tn.Node.(type) {
		case *ShortNode:
			if len(key) < len(n.Key) || !bytes.Equal(n.Key, key[:len(n.Key)]) {
				// The trie doesn't contain the key.
				tn.Node = nil
			} else {
				tn = n.Val
				key = key[len(n.Key):]
			}
			nodes = append(nodes, NodeIntf{n})
		case *FullNode:
			tn = n.Children[key[0]]
			key = key[1:]
			nodes = append(nodes, NodeIntf{n})
		case HashNode:
			var err error
			tn, err = t.resolveHash(n, nil)
			if err != nil {
				logger.Error(fmt.Sprintf("Unhandled trie error: %v", err))
				return err
			}
		default:
			panic(fmt.Sprintf("%T: invalid node: %v", tn, tn))
		}
	}
	hasher := newHasher(0, 0)
	for i, n := range nodes {
		// Don't bother checking for errors here since hasher panics
		// if encoding doesn't work and we're not writing to any database.
		n, _, _ = hasher.hashChildren(n, nil)
		hn, _ := hasher.store(n, nil, false)
		if hash, ok := hn.Node.(HashNode); ok || i == 0 {
			// If the node's database encoding is a hash (or is the
			// root node), it becomes a proof element.
			if fromLevel > 0 {
				fromLevel--
			} else {
				enc := bytes.Buffer{}
				wr := msgp.NewWriter(&enc)
				err := wr.WriteIntf(&n)
				if err != nil {
					panic(fmt.Sprintf("Encode failed."))
				}
				wr.Flush()
				if !ok {
					hash = crypto.Keccak256(enc.Bytes())
				}
				proofDb.Put(hash, enc.Bytes())
			}
		}
	}
	return nil
}

// VerifyProof checks merkle proofs. The given proof must contain the
// value for key in a trie with the given root hash. VerifyProof
// returns an error if the proof contains invalid trie nodes or the
// wrong value.
func VerifyProof(rootHash types.Hash, key []byte, proofDb DatabaseReader) (value []byte, err error, nodes int) {
	key = keybytesToHex(key)
	wantHash := rootHash[:]
	for i := 0; ; i++ {
		buf, _ := proofDb.Get(wantHash)
		if buf == nil {
			return nil, fmt.Errorf("proof node %d (hash %064x) missing", i, wantHash[:]), i
		}
		n, err := decodeNode(wantHash, buf, 0)
		if err != nil {
			return nil, fmt.Errorf("bad proof node %d: %v", i, err), i
		}
		keyrest, cld := get(n, key)
		switch cld := cld.Node.(type) {
		case nil:
			// The trie doesn't contain the key.
			return nil, nil, i
		case HashNode:
			key = keyrest
			wantHash = cld
		case ValueNode:
			return cld, nil, i + 1
		}
	}
}

func get(tn NodeIntf, key []byte) ([]byte, NodeIntf) {
	for {
		switch n := tn.Node.(type) {
		case *ShortNode:
			if len(key) < len(n.Key) || !bytes.Equal(n.Key, key[:len(n.Key)]) {
				return nil, NodeIntf{nil}
			}
			tn = n.Val
			key = key[len(n.Key):]
		case *FullNode:
			tn = n.Children[key[0]]
			key = key[1:]
		case HashNode:
			return key, NodeIntf{n}
		case nil:
			return key, NodeIntf{nil}
		case ValueNode:
			return nil, NodeIntf{n}
		default:
			panic(fmt.Sprintf("%T: invalid node: %v", tn, tn))
		}
	}
}
