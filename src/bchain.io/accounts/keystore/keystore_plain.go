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
// @File: keystore_plain.go
// @Date: 2018/05/08 17:16:08
////////////////////////////////////////////////////////////////////////////////

package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"bchain.io/common/types"
)

type keyStorePlain struct {
	keysDirPath string
}

func (ks keyStorePlain) GetKey(addr types.Address, filename, auth string) (*Key, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	key := new(Key)
	if err := json.NewDecoder(fd).Decode(key); err != nil {
		return nil, err
	}
	if key.Address != addr {
		return nil, fmt.Errorf("key content mismatch: have address %x, want %x", key.Address, addr)
	}
	return key, nil
}

func (ks keyStorePlain) StoreKey(filename string, key *Key, auth string) error {
	content, err := json.Marshal(key)
	if err != nil {
		return err
	}
	return writeKeyFile(filename, content)
}

func (ks keyStorePlain) JoinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	} else {
		return filepath.Join(ks.keysDirPath, filename)
	}
}
