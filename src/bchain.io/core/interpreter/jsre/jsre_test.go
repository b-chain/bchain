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
// @File: jsre_test.go
// @Date: 2018/08/06 15:25:06
////////////////////////////////////////////////////////////////////////////////

package jsre

import (
	"testing"
	"bchain.io/core/actioncontext"
	"time"
	"fmt"
	"bchain.io/core/interpreter"
)

func TestJSRE_PluginName(t *testing.T) {
	pn := interpreter.PluginName{}
	pn.Set(interpreter.PluginImpl(&JSRE{}))
	t.Logf("%x, %s", pn.Id(), pn.Name())
}

func TestJsre_1(t *testing.T) {
	ctx := actioncontext.Context{}
	ctx.InitForTest()

	jsre := JSRE{}
	jsre.Initialize()
	jsre.Startup()

	rt := jsre.Generate(&ctx)
	/*code := `
//var apijs = require('context_api') // not support

var ctxapi = ctxApi();
;(function () {
	this.test = function(msg) {
		ctxapi.console.print(msg)
		ctxapi.console.print("\n")
		ctxapi.console.print("============================================================================\n");
		ctxapi.console.print("ctxapi._contract: " + ctxapi._contract + "\n");
		ctxapi.console.print("ctxapi._sender:   " + ctxapi._sender + "\n");
		ctxapi.console.print("============================================================================\n");
	}
})(this);`
	param := `
test("hi, this is a test");
`*/
	time := rt.Exec(ctx.Code(), ctx.Param(), 1 * time.Second)

	fmt.Printf("Stopping after: %d ns\n", time.Nanoseconds())

	jsre.Shutdown()
}
func TestJsre_speed(t *testing.T) {
	ctx := actioncontext.Context{}
	ctx.InitForTest()

	jsre := JSRE{}
	jsre.Initialize()
	jsre.Startup()

	rt := jsre.Generate(&ctx)

	time_start := time.Now();
	for i:=0; i<100; i++{
		rt.Exec(ctx.Code(), ctx.Param(), 1 * time.Second)

	}

	//time := rt.Exec(ctx.Code(), ctx.Param(), 1 * time.Second)
	fmt.Println("time:" , time.Since(time_start))

	//fmt.Printf("Stopping after: %d ns\n", time.Nanoseconds())

	jsre.Shutdown()
}
