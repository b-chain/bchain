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
// @File: jsre.go
// @Date: 2018/07/31 09:00:31
////////////////////////////////////////////////////////////////////////////////

package jsre

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/parser"
	"bchain.io/common/assert"
	"bchain.io/core/actioncontext"
	"bchain.io/core/interpreter"
	"bchain.io/core/interpreter/jsre/deps"
	"os"
	"time"
)

var (
	BigNumber_JS = deps.MustAsset("bignumber.js")
	ctxApi_JS    = deps.MustAsset("context_api.js")
)

//func init() {
//	interpreter.Singleton().Register(func() interpreter.PluginImpl {
//		return &JSRE{}
//	})
//}

/*
JSRE is a generic JS Runtime environment embedding the otto JS interpreter.
*/
type JSRE struct {
	b *bridge
}

func (self *JSRE) Initialize() {

}

func (self *JSRE) Startup() {

}

func (self *JSRE) Shutdown() {

}

var g_vm *otto.Otto

func init() {
	fmt.Println("JSRE INIT..........")
	g_vm = otto.New()

	if _, err := compileAndRun(g_vm, "bignumber.js", BigNumber_JS); err != nil {
		assert.AssertEx(false, fmt.Sprintf("bignumber.js: %v", err))
	}
	if _, err := compileAndRun(g_vm, "context_api.js", ctxApi_JS); err != nil {
		assert.AssertEx(false, fmt.Sprintf("context_api.js: %v", err))
	}
}

func (self *JSRE) Generate(ctx *actioncontext.Context) interpreter.Interpreter {
	// new Runtime
	rt := Runtime{
		ctx:   ctx,
		vm:    g_vm.Copy(),
		close: make(chan struct{}),
		timer: time.NewTimer(0),
	}
	rt.vm.Interrupt = make(chan func(), 1)
	<-rt.timer.C

	// bridge variable and api
	self.b = newBridge(&rt)
	self.b.load()

	// environment configuration
	// Load all the internal utility JavaScript libraries
	//if err := rt.compile("bignumber.js", bigNumber_JS); err != nil {
	//	assert.AssertEx(false, fmt.Sprintf("bignumber.js: %v", err))
	//}
	//if err := rt.compile("context_api.js", ctxApi_JS); err != nil {
	//	assert.AssertEx(false, fmt.Sprintf("context_api.js: %v", err))
	//}

	return &rt
}

// the Runtime object
type Runtime struct {
	ctx   *actioncontext.Context
	vm    *otto.Otto // otto vm
	close chan struct{}
	timer *time.Timer
}

func (rt *Runtime) Exec(code []byte, param []byte, ttl time.Duration) time.Duration {
	start := time.Now()
	var duration time.Duration
	halt := interpreter.TimeoutInterrupter
	stopChan := make(chan interface{})

	defer func() {
		duration = time.Since(start)
		if caught := recover(); caught != nil {
			if caught == halt {
				fmt.Fprintf(os.Stderr, "Some code took to long! Stopping after: %v\n", duration)
				//return
			}
			panic(caught) // repanic!
		}
		rt.timer.Stop() // TODO: timer, need to correct shutdown !!!
		close(stopChan)
	}()

	rt.timer.Reset(ttl)
	go func() {
		select {
		case now := <-rt.timer.C:
			_ = now
			rt.vm.Interrupt <- func() {
				panic(halt)
			}
		case <-stopChan:
			return
		}
	}()

	program, err := parser.ParseFile(nil, "", code, 0)
	assert.AssertEx(err == nil, fmt.Sprintf("Exec: %v", err))
	_, err = rt.vm.Run(program)
	assert.AssertEx(err == nil, fmt.Sprintf("Exec: %v", err))
	resultValue, err := rt.vm.Run(param)
	assert.AsserErr(err)
	_ = resultValue
	//todo:resultValue not use

	//resultInt , err := resultValue.ToInteger()
	//assert.AsserErr(err)
	//fmt.Println("+++++++++++++++Vm return:" , resultInt)
	return time.Since(start)
}

// Compile compiles and then runs a piece of JS code.
func (rt *Runtime) compile(filename string, src interface{}) (err error) {
	_, err = compileAndRun(rt.vm, filename, src)
	return err
}

func compileAndRun(vm *otto.Otto, filename string, src interface{}) (otto.Value, error) {
	script, err := vm.Compile(filename, src)
	if err != nil {
		return otto.Value{}, err
	}
	fmt.Printf("CompileAndRun file path:%s \n", filename)
	return vm.Run(script)
}
