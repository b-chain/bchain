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
// @File: interpreter_test.go
// @Date: 2018/08/08 15:40:08
////////////////////////////////////////////////////////////////////////////////

package interpreter_test

import (
	"bchain.io/common/types"
	"bchain.io/core/actioncontext"
	"bchain.io/core/interpreter"
	"bchain.io/core/interpreter/jsre"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
	"bchain.io/utils/database"
	"testing"
	"time"
)

//go:generate jsmarshal test_contract.js test_contract_code.md jsre.JSRE

func TestPlugin_1(t *testing.T) {
	interpreter.Singleton().Register(func() interpreter.PluginImpl {
		return &jsre.JSRE{}
	})

	interpreter.Singleton().Initialize()
	interpreter.Singleton().Startup()

	// Create an empty state database
	db, _ := database.OpenMemDB()
	stateDb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	tmpdb, _ := database.OpenMemDB()

	tr := transaction.Action{}
	blkctx := actioncontext.NewBlockContext(stateDb, db, tmpdb,nil,types.Address{})
	ctx := actioncontext.NewContext(types.Address{}, &tr, blkctx)

	ctx.InitForTest()
	ctx.Exec(interpreter.Singleton())

	time.Sleep(1 * time.Second)
	t.Logf("shutdown ...\n")
	interpreter.Singleton().Shutdown()

	time.Sleep(1 * time.Second)
}
