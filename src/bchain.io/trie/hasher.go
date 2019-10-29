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
// @File: hasher.go
// @Date: 2018/05/07 10:43:07
////////////////////////////////////////////////////////////////////////////////

package trie

import (
	"bytes"
	"hash"
	"sync"

	"bchain.io/utils/crypto/sha3"
	"github.com/tinylib/msgp/msgp"
	"bchain.io/common/types"
)

type hasher struct {
	tmp                  *bytes.Buffer
	sha                  hash.Hash
	cachegen, cachelimit uint16
}

// hashers live in a global pool.
var hasherPool = sync.Pool{
	New: func() interface{} {
		return &hasher{tmp: new(bytes.Buffer), sha: sha3.NewKeccak256()}
	},
}

func newHasher(cachegen, cachelimit uint16) *hasher {
	h	 := hasherPool.Get().(*hasher)
	h.cachegen, h.cachelimit = cachegen, cachelimit
	return h
}

func returnHasherToPool(h *hasher) {
	hasherPool.Put(h)
}

// hash collapses a node down into a hash node, also returning a copy of the
// original node initialized with the computed hash to replace the original one.
func (h *hasher) hash(n NodeIntf, db DatabaseWriter, force bool) (NodeIntf, NodeIntf, error) {
	// If we're not storing the node, just hashing, use available cached data
	if hash, dirty := n.cache(); hash != nil {
		if db == nil {
			return NodeIntf{hash}, n, nil
		}
		if n.canUnload(h.cachegen, h.cachelimit) {
			// Unload the node from cache. All of its subnodes will have a lower or equal
			// cache generation number.
			cacheUnloadCounter.Inc(1)
			return NodeIntf{hash}, NodeIntf{hash}, nil
		}
		if !dirty {
			return NodeIntf{hash}, n, nil
		}
	}
	// Trie not processed yet or needs storage, walk the children
	collapsed, cached, err := h.hashChildren(n, db)
	if err != nil {
		return NodeIntf{HashNode{}}, n, err
	}
	hashed, err := h.store(collapsed, db, force)
	if err != nil {
		return NodeIntf{HashNode{}}, n, err
	}
	// Cache the hash of the node for later reuse and remove
	// the dirty flag in commit mode. It's fine to assign these values directly
	// without copying the node first because hashChildren copies it.
	cachedHash, _ := hashed.Node.(HashNode)
	switch cn := cached.Node.(type) {
	case *ShortNode:
		cn.flags.hash = cachedHash
		if db != nil {
			cn.flags.dirty = false
		}
	case *FullNode:
		cn.flags.hash = cachedHash
		if db != nil {
			cn.flags.dirty = false
		}
	}
	return hashed, cached, nil
}

// hashChildren replaces the children of a node with their hashes if the encoded
// size of the child is larger than a hash, returning the collapsed node as well
// as a replacement for the original node with the child hashes cached in.
func (h *hasher) hashChildren(original NodeIntf, db DatabaseWriter) (NodeIntf, NodeIntf, error) {
	var err error

	switch n := original.Node.(type) {
	case *ShortNode:
		// Hash the short node's child, caching the newly hashed subtree
		collapsed, cached := n.copy(), n.copy()
		collapsed.Key = hexToCompact(n.Key)
		cached.Key = types.CopyBytes(n.Key)

		if _, ok := n.Val.Node.(ValueNode); !ok {
			collapsed.Val, cached.Val, err = h.hash(n.Val, db, false)
			if err != nil {
				return original, original, err
			}
		}
		if collapsed.Val.Node == nil {
			collapsed.Val.Node = ValueNode(nil) // Ensure that nil children are encoded as empty strings.
		}
		return NodeIntf{collapsed}, NodeIntf{cached}, nil

	case *FullNode:
		// Hash the full node's children, caching the newly hashed subtrees
		collapsed, cached := n.copy(), n.copy()

		for i := 0; i < 16; i++ {
			if n.Children[i].Node != nil {
				collapsed.Children[i], cached.Children[i], err = h.hash(n.Children[i], db, false)
				if err != nil {
					return original, original, err
				}
			} else {
				collapsed.Children[i].Node = ValueNode(nil) // Ensure that nil children are encoded as empty strings.
			}
		}
		cached.Children[16] = n.Children[16]
		if collapsed.Children[16].Node == nil {
			collapsed.Children[16].Node = ValueNode(nil)
		}
		return NodeIntf{collapsed}, NodeIntf{cached}, nil

	default:
		// Value and hash nodes don't have children so they're left as were
		return NodeIntf{n}, original, nil
	}
}

func (h *hasher) store(n NodeIntf, db DatabaseWriter, force bool) (NodeIntf, error) {
	// Don't store hashes or empty nodes.
	if _, isHash := n.Node.(HashNode); n.Node == nil || isHash {
		return n, nil
	}
	// Generate the Msgp encoding of the node
	h.tmp.Reset()
	wr := msgp.NewWriter(h.tmp)
	if err := wr.WriteIntf(n); err != nil {
		panic("encode error: " + err.Error())
	}
	wr.Flush()

	if h.tmp.Len() < 32 && !force {
		return n, nil // Nodes smaller than 32 bytes are stored inside their parent
	}
	// Larger nodes are replaced by their hash and stored in the database.
	hash, _ := n.cache()
	if hash == nil {
		h.sha.Reset()
		h.sha.Write(h.tmp.Bytes())
		hash = HashNode(h.sha.Sum(nil))
	}
	if db != nil {
		return NodeIntf{hash}, db.Put(hash, h.tmp.Bytes())
	}
	return NodeIntf{hash}, nil
}
