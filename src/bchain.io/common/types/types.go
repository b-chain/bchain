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
// @File: types.go
// @Date: 2018/05/07 15:09:07
////////////////////////////////////////////////////////////////////////////////

package types

import (
	"fmt"
	"bchain.io/log"
	"github.com/tinylib/msgp/msgp"
	"os"
	"sync"
)

var (
	globalTypeIndex = int8(10)
	mu              sync.Mutex

	logTag = "common.types"
	logger log.Logger
)

func init() {
	// get a logger
	logger = log.GetLogger(logTag)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", logTag)
		os.Exit(1)
	}

	// Registering an extension is as simple as matching the
	// appropriate type number with a function that initializes
	// a freshly-allocated object of that type
	RegisterExtension(&hashType, func() msgp.Extension { return new(Hash) })
	RegisterExtension(&addressType, func() msgp.Extension { return new(Address) })
	RegisterExtension(&ipType, func() msgp.Extension { return new(IP) })
	RegisterExtension(&bigIntType, func() msgp.Extension { return new(BigInt) })
	RegisterExtension(&bloomType, func() msgp.Extension { return new(Bloom) })
}

type decError struct{ msg string }

func (err decError) Error() string { return err.msg }

var (
	ErrBytesTooLong = &decError{"bytes too long"}

	ErrRegisterFull    = &decError{"Can't register more type"}
	ErrRegisterFailure = &decError{"Register is failure"}
)

func registerExtension(typ *int8, f func() msgp.Extension) error {
	mu.Lock()
	defer func() (err error) {
		if p := recover(); p != nil {
			fmt.Printf("panic recover! p: %v", p)
			err = ErrRegisterFailure
		}

		mu.Unlock()
		return err
	}()

	if globalTypeIndex == -128 {
		return ErrRegisterFull
	}
	msgp.RegisterExtension(globalTypeIndex, f)
	*typ = globalTypeIndex
	globalTypeIndex++

	return nil
}

func RegisterExtension(typ *int8, f func() msgp.Extension) error {
	for {
		err := registerExtension(typ, f)
		if err == ErrRegisterFailure {
			continue
		} else {
			return err
		}
	}
}
