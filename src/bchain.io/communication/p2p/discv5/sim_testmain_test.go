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
// @File: sim_testmain_test.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

// +build go1.4,nacl,faketime_simulation

package discv5

import (
	"os"
	"runtime"
	"testing"
	"unsafe"
)

// Enable fake time mode in the runtime, like on the go playground.
// There is a slight chance that this won't work because some go code
// might have executed before the variable is set.

//go:linkname faketime runtime.faketime
var faketime = 1

func TestMain(m *testing.M) {
	// We need to use unsafe somehow in order to get access to go:linkname.
	_ = unsafe.Sizeof(0)

	// Run the actual test. runWithPlaygroundTime ensures that the only test
	// that runs is the one calling it.
	runtime.GOMAXPROCS(8)
	os.Exit(m.Run())
}
