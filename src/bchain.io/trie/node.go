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
// @File: node.go
// @Date: 2018/05/07 10:44:07
////////////////////////////////////////////////////////////////////////////////

package trie

import (
	"fmt"
	"io"
	"github.com/tinylib/msgp/msgp"
	"bytes"
)

//go:generate msgp

var indices = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "[17]"}

type Node interface {
	fstring(string) string
	cache() (HashNode, bool)
	canUnload(cachegen, cachelimit uint16) bool
}

//msgp:shim NodeIntf as:interface{} using:fromNode/toNode mode:convert
type NodeIntf struct {
	Node
}

func fromNode(v NodeIntf) (interface{}, error) {
	node := v.Node
	if node == nil {
		return nil, nil
	}

	buf := bytes.Buffer{}
	wr := msgp.NewWriter(&buf)
	var err error
	switch node.(type) {
	case *FullNode:
		wr.WriteString("F")
		n, _ := node.(*FullNode)
		err = n.EncodeMsg(wr)
	case *ShortNode:
		wr.WriteString("S")
		n, _ := node.(*ShortNode)
		err = n.EncodeMsg(wr)
	case HashNode:
		wr.WriteString("H")
		n, _ := node.(HashNode)
		err = n.EncodeMsg(wr)
	case ValueNode:
		wr.WriteString("V")
		n, _ := node.(ValueNode)
		err = n.EncodeMsg(wr)
	default:
		panic("v.Node type is invalid")
	}
	wr.Flush()
	return buf.Bytes(), err
}

func toNode(s interface{}) (NodeIntf, error) {
	if s == nil {
		return NodeIntf{nil}, nil
	}

	buf, ok := s.([]byte)
	if !ok {
		return NodeIntf{nil}, fmt.Errorf("can't convert")
	}

	var err error

	//rd := bytes.NewReader(buf)
	//var prefix string
	prefix, buf, err := msgp.ReadStringBytes(buf)
	if err != nil {
		return NodeIntf{nil}, fmt.Errorf("error prefix")
	}

	switch prefix {
	case "F":
		n := FullNode{}
		err = msgp.Decode(bytes.NewReader(buf), &n)
		if err == nil {
			return NodeIntf{&n}, nil
		}
	case "S":
		n := ShortNode{}
		err = msgp.Decode(bytes.NewReader(buf), &n)
		if err == nil {
			n.Key = compactToHex(n.Key)
			return NodeIntf{&n}, nil
		}
	case "H":
		n := HashNode{}
		err = msgp.Decode(bytes.NewReader(buf), &n)
		if err == nil {
			return NodeIntf{n}, nil
		}
	case "V":
		n := ValueNode{}
		err = msgp.Decode(bytes.NewReader(buf), &n)
		if err == nil {
			if len(n) != 0 {
				return NodeIntf{n}, nil
			}
			return NodeIntf{nil}, nil
		}
	}

	return NodeIntf{nil}, fmt.Errorf("s type is unknown")
}

type (
	FullNode struct {
		Children [17]NodeIntf // Actual trie node data to encode/decode (needs custom encoder)
		flags    nodeFlag
	}
	ShortNode struct {
		Key   []byte
		Val   NodeIntf
		flags nodeFlag
	}
	HashNode  []byte
	ValueNode []byte
)

func (n *FullNode) copy() *FullNode   { copy := *n; return &copy }
func (n *ShortNode) copy() *ShortNode { copy := *n; return &copy }

// nodeFlag contains caching-related metadata about a node.
type nodeFlag struct {
	hash  HashNode // cached hash of the node (may be nil)
	gen   uint16   // cache generation counter
	dirty bool     // whether the node has changes that must be written to the database
}

// canUnload tells whether a node can be unloaded.
func (n *nodeFlag) canUnload(cachegen, cachelimit uint16) bool {
	return !n.dirty && cachegen-n.gen >= cachelimit
}

func (n *FullNode) canUnload(gen, limit uint16) bool  { return n.flags.canUnload(gen, limit) }
func (n *ShortNode) canUnload(gen, limit uint16) bool { return n.flags.canUnload(gen, limit) }
func (n HashNode) canUnload(uint16, uint16) bool      { return false }
func (n ValueNode) canUnload(uint16, uint16) bool     { return false }

func (n *FullNode) cache() (HashNode, bool)  { return n.flags.hash, n.flags.dirty }
func (n *ShortNode) cache() (HashNode, bool) { return n.flags.hash, n.flags.dirty }
func (n HashNode) cache() (HashNode, bool)   { return nil, true }
func (n ValueNode) cache() (HashNode, bool)  { return nil, true }

// Pretty printing.
func (n *FullNode) String() string  { return n.fstring("") }
func (n *ShortNode) String() string { return n.fstring("") }
func (n HashNode) String() string   { return n.fstring("") }
func (n ValueNode) String() string  { return n.fstring("") }

func (n *FullNode) fstring(ind string) string {
	resp := fmt.Sprintf("[\n%s  ", ind)
	for i, node := range n.Children {
		if node.Node == nil {
			resp += fmt.Sprintf("%s: <nil> ", indices[i])
		} else {
			resp += fmt.Sprintf("%s: %v", indices[i], node.fstring(ind+"  "))
		}
	}
	return resp + fmt.Sprintf("\n%s] ", ind)
}
func (n *ShortNode) fstring(ind string) string {
	return fmt.Sprintf("{%x: %v} ", n.Key, n.Val.fstring(ind+"  "))
}
func (n HashNode) fstring(ind string) string {
	return fmt.Sprintf("<%x> ", []byte(n))
}
func (n ValueNode) fstring(ind string) string {
	return fmt.Sprintf("%x ", []byte(n))
}

func mustDecodeNode(hash, buf []byte, cachegen uint16) NodeIntf {
	n, err := decodeNode(hash, buf, cachegen)
	if err != nil {
		panic(fmt.Sprintf("node %x: %v", hash, err))
	}
	return n
}

// decodeNode parses the Msgp encoding of a trie node.
func decodeNode(hash, buf []byte, cachegen uint16) (NodeIntf, error) {
	if len(buf) == 0 {
		return NodeIntf{}, io.ErrUnexpectedEOF
	}

	var node NodeIntf
	err := node.DecodeMsg(msgp.NewReader(bytes.NewReader(buf)))
	if err != nil {
		return NodeIntf{}, err
	}

	// set nodeFlag
	node.setNodeflag(hash, cachegen)

	return node, err
}

func (v *NodeIntf) setNodeflag(hash HashNode, cachegen uint16) {
	switch n := v.Node.(type) {
	case *ShortNode:
		n.flags.hash = hash
		n.flags.gen = cachegen
		n.Val.setNodeflag(nil, cachegen)
	case *FullNode:
		n.flags.hash = hash
		n.flags.gen = cachegen
		for _, ch := range n.Children {
			ch.setNodeflag(nil, cachegen)
		}
	default:
	}
}
