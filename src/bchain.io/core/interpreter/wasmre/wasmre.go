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
// @File: wasmre.go
// @Date: 2018/12/05 13:58:05
//
////////////////////////////////////////////////////////////////////////////////

package wasmre

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
	"github.com/hashicorp/golang-lru"
	"bchain.io/common/assert"
	"bchain.io/core/actioncontext"
	"bchain.io/core/interpreter"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"os"
	"sync"
	"time"
	"bchain.io/common/types"
)

/*
WasmRE is a generic WebAssembly Runtime environment embedding the wasm interpreter.
*/

const (
	vmCacheLimit = 1024
)

type ParaPaser interface {
	ParseInputPara(para []byte, base, max int) (funcName string, args []uint64, mem []byte)
}

func NewWasmRe() interpreter.PluginImpl {
	return &WasmRE{}
}

type WasmRE struct {
	vmCache *lru.Cache
}

func (self *WasmRE) Initialize() {
	self.vmCache, _ = lru.New(vmCacheLimit)
}

func (self *WasmRE) Startup() {

}

func (self *WasmRE) Shutdown() {
	self.vmCache.Purge()
}

func (self *WasmRE) newRuntime(ctx *actioncontext.Context) *Runtime {
	// new Runtime
	mApi := wasm.NewModule()
	mApi.Export.Entries = make(map[string]wasm.ExportEntry)
	rt := &Runtime{
		ctx:   ctx,
		m:     nil,
		mApi:  mApi,
		vm:    nil,
		close: make(chan struct{}),
		timer: time.NewTimer(0),
		pp:    new(para_paser.DummyParaPaser), //todo :: parse para, should use abi
		running:false,
	}
	<-rt.timer.C

	// bridge variable and api
	//self.b = newBridge(rt)
	//self.b.load()
	b := newBridge(rt)
	b.load()

	return rt
}
var BcContract = "0xb78f12Cb3924607A8BC6a66799e159E3459097e9"
func (self *WasmRE) Generate(ctx *actioncontext.Context) interpreter.Interpreter {
	conAddr := ctx.ContractAddress()
	//bc20 fork point
	if conAddr == types.HexToAddress(BcContract) {
		return newBc20Runtime(ctx)
	}
	if cached, ok := self.vmCache.Get(conAddr); ok {
		rt := cached.(*Runtime)
		if rt.running == false  {
			rt.ctx = ctx
			return rt
		}
	}
	rt := self.newRuntime(ctx)
	self.vmCache.Add(conAddr, rt)
	return rt
}

// the Runtime object
type Runtime struct {
	ctx   *actioncontext.Context
	m     *wasm.Module
	mApi  *wasm.Module
	vm    *exec.VM
	close chan struct{}
	timer *time.Timer
	pp    ParaPaser
	running bool
}

func (rt *Runtime) importer(name string) (*wasm.Module, error) {
	return rt.mApi, nil
}

func (rt *Runtime) createVm(code []byte) {
	if rt.vm == nil {
		// create wasm vm
		buf := bytes.NewBuffer(code)
		m, err := wasm.ReadModule(buf, rt.importer)
		//fmt.Println("load time:",time.Since(start))
		assert.AssertEx(err == nil, fmt.Sprintf("Exec: %v", err))
		assert.AssertEx(m != nil, "Exec: wasm read module return nil")
		assert.AssertEx(m.Export != nil, "Exec: module has no export section")

		vm, errVm := exec.NewVM(m)
		//fmt.Println("vm time:",time.Since(start))
		assert.AssertEx(errVm == nil, fmt.Sprintf("Exec: could not create VM %v", err))
		rt.vm = vm
		rt.m = m
	} else {
		// should reset memory
		mem := make([]byte, len(rt.vm.Memory()))
		copy(rt.vm.Memory(), mem)
		copy(rt.vm.Memory(), rt.m.LinearMemoryIndexSpace[0])

		// exec the start function
		if rt.m.Start != nil {
			_, err := rt.vm.ExecCode(int64(rt.m.Start.Index))
			assert.AssertEx(err == nil, fmt.Sprintf("Exec: exec vm start function err: %v", err))
		}
	}
}

func (rt *Runtime) Exec(code []byte, param []byte, ttl time.Duration) time.Duration {
	start := time.Now()
	var duration time.Duration
	halt := interpreter.TimeoutInterrupter
	stopChan := make(chan interface{})
	var wg sync.WaitGroup

	defer func() {
		rt.running = false
		duration = time.Since(start)
		rt.timer.Stop()
		close(stopChan)
		wg.Wait()
		if caught := recover(); caught != nil {
			if caught == halt {
				fmt.Fprintf(os.Stderr, "Some code took to long! Stopping after: %v\n", duration)
			}
			panic(caught) // repanic!
		}
	}()
	rt.running = true
	rt.timer.Reset(ttl)
	timeOut := false
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-rt.timer.C:
			timeOut = true
			if rt.vm != nil {
				proc := exec.NewProcess(rt.vm)
				proc.Terminate()
				// todo, because of after VM Terminate, no api to restart it.
				// Here just make it nil to restart the vm in cache
				rt.vm = nil
			}
		case <-stopChan:
			return
		}
	}()

	rt.createVm(code)
	assert.AssertEx(timeOut == false, halt)

	m := rt.m
	vm := rt.vm

	vmMem := vm.Memory()
	oriStack := binary.LittleEndian.Uint32(vmMem[4:8])
	base := (len(m.LinearMemoryIndexSpace[0]) + 15) / 16 * 16
	// the stack must >= 4 k
	assert.AssertEx(oriStack < 65536 && uint32(base)+4096 <= oriStack, "stack val error!")
	maxMem := 65536 - oriStack + uint32(base)
	entryName, args, mem := rt.pp.ParseInputPara(param, base, int(maxMem))

	newStack := oriStack + uint32((len(mem)+15)/16*16)
	assert.AssertEx(newStack < 65536, "new stack exceed!")
	encNewStack := make([]byte, 4)
	binary.LittleEndian.PutUint32(encNewStack, newStack)
	copy(vmMem[4:], encNewStack)
	copy(vmMem[base:], mem)

	ee, ok := m.Export.Entries[entryName]
	assert.AssertEx(ok == true, "Exec: can not find para")

	entryIdx := int64(ee.Index)
	_, errExec := vm.ExecCode(entryIdx, args...)
	//fmt.Println("WasmRE return val:", o)
	assert.AssertEx(errExec == nil, fmt.Sprintf("Exec: exec code err: %v", errExec))
	assert.AssertEx(timeOut == false, halt)

	//fmt.Println("exec time:",time.Since(start))
	return time.Since(start)
}
