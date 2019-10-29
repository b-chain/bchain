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
// @File: dump.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package state

import (
	"encoding/json"
	"fmt"
	"bchain.io/common/types"
	"bchain.io/trie"
	"bytes"
	"github.com/tinylib/msgp/msgp"
)

type DumpAccount struct {
	Nonce    uint64            `json:"nonce"`
	Root     string            `json:"root"`
	CodeHash string            `json:"codeHash"`
	Code     string            `json:"code"`
	Storage  map[string]string `json:"storage"`
}

type Dump struct {
	Root     string                 `json:"root"`
	Accounts map[string]DumpAccount `json:"accounts"`
}

func (self *StateDB) RawDump() Dump {
	dump := Dump{
		Root:     fmt.Sprintf("%x", self.trie.Hash()),
		Accounts: make(map[string]DumpAccount),
	}

	it := trie.NewIterator(self.trie.NodeIterator(nil))
	for it.Next() {
		addr := self.trie.GetKey(it.Key)
		var data Account
		byteBuf := bytes.NewBuffer(it.Value)
		if err := msgp.Decode(byteBuf, &data); err!= nil{
			panic(err)
		}

		obj := newObject(nil, types.BytesToAddress(addr), data, nil)
		account := DumpAccount{
			Nonce:    data.Nonce,
			
			Root:     types.Bytes2Hex(data.Root[:]),
			CodeHash: types.Bytes2Hex(data.CodeHash),
			Code:     types.Bytes2Hex(obj.Code(self.db)),
			Storage:  make(map[string]string),
		}
		storageIt := trie.NewIterator(obj.getTrie(self.db).NodeIterator(nil))
		for storageIt.Next() {
			account.Storage[types.Bytes2Hex(self.trie.GetKey(storageIt.Key))] = types.Bytes2Hex(storageIt.Value)
		}
		dump.Accounts[types.Bytes2Hex(addr)] = account
	}
	return dump
}

func (self *StateDB) Dump() []byte {
	json, err := json.MarshalIndent(self.RawDump(), "", "    ")
	if err != nil {
		fmt.Println("dump err", err)
	}

	return json
}
