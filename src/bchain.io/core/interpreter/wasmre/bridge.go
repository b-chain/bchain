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
// @File: bridge.go
// @Date: 2018/12/05 16:43:05
//
////////////////////////////////////////////////////////////////////////////////

package wasmre

import (
	"bchain.io/common/assert"
	"bchain.io/common/types"
	"bchain.io/core/actioncontext"
	"bchain.io/core/interpreter"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"bchain.io/core/transaction"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
	"reflect"
	"strconv"
	"math"
	"math/big"
)

type bridge struct {
	rt *Runtime
}

func newBridge(rt *Runtime) *bridge {
	return &bridge{
		rt: rt,
	}
}

func (self *bridge) registerApi(fs *wasm.FunctionSig, apiName string, api interface{}) {
	self.rt.mApi.Types.Entries = append(self.rt.mApi.Types.Entries, *fs)
	f := wasm.Function{
		Sig:  fs,
		Host: reflect.ValueOf(api),
		Body: &wasm.FunctionBody{}, // create a dummy wasm body (the actual value will be taken from Host.)
	}
	self.rt.mApi.FunctionIndexSpace = append(self.rt.mApi.FunctionIndexSpace, f)
	ee := wasm.ExportEntry{
		FieldStr: apiName,
		Kind:     wasm.ExternalFunction,
		Index:    uint32(len(self.rt.mApi.FunctionIndexSpace) - 1),
	}
	self.rt.mApi.Export.Entries[apiName] = ee
}

//func (self *bridge) registerVariable(varType wasm.ValueType, varName string, initVal []byte) {
//	gv := wasm.GlobalVar{
//		Type:    varType,
//		Mutable: false,
//	}
//	ge := wasm.GlobalEntry{
//		Type: gv,
//		Init: initVal,
//	}
//	self.rt.m.GlobalIndexSpace = append(self.rt.m.GlobalIndexSpace, ge)
//	ee := wasm.ExportEntry{
//		FieldStr: varName,
//		Kind:     wasm.ExternalGlobal,
//		Index:    uint32(len(self.rt.m.GlobalIndexSpace) - 1),
//	}
//	self.rt.m.Export.Entries[varName] = ee
//}

func (self *bridge) load() {
	self.bridgeApi()
}

// bridge all api
func (self *bridge) bridgeApi() {
	self.bridgeMemoryApi()
	self.bridgeAuthApi()
	self.bridgeSysApi()
	self.bridgeConsoleApi()
	self.bridgeAssertApi()
	self.bridgeCryptoApi()
	self.bridgeProducerApi()
	self.bridgeActApi()
	self.bridgeDbApi()
	self.bridgeCacheApi()
	self.bridgeContractApi()
	self.bridgeResultApi()
	self.bridgeCallApi()
	self.bridgeStringApi()
	self.bridgeWeightApi()
	self.bridgeBigIntApi()
}

// bridge Console Api
func (self *bridge) bridgeMemoryApi() {
	// void *memset(void *s, int ch, size_t n);
	memset := func(proc *exec.Process, s_ptr int32, ch int32, n int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(s_ptr+n < mem_len, "s input pointer is exceed")
		for i := 0; i < int(n); i++ {
			mem[int(s_ptr)+i] = byte(ch)
		}
		return s_ptr
	}
	fs := &wasm.FunctionSig{
		ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "memset", memset)

	// void *memcpy(void *dest, const void *src, size_t n);
	memcpy := func(proc *exec.Process, dest_ptr int32, src_ptr int32, n int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(dest_ptr+n < mem_len, "dest input pointer is exceed")
		assert.AssertEx(src_ptr+n < mem_len, "src input pointer is exceed")
		assert.AssertEx(src_ptr-dest_ptr >= n || dest_ptr-src_ptr >= n, "memcpy can only accept non-aliasing pointers")
		copy(mem[dest_ptr:dest_ptr+n], mem[src_ptr:src_ptr+n])
		return dest_ptr
	}
	self.registerApi(fs, "memcpy", memcpy)

	// void *memmove( void* dest, const void* src, size_t count );
	memmove := func(proc *exec.Process, dest_ptr int32, src_ptr int32, n int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(dest_ptr+n < mem_len, "dest input pointer is exceed")
		assert.AssertEx(src_ptr+n < mem_len, "src input pointer is exceed")
		copy(mem[dest_ptr:dest_ptr+n], mem[src_ptr:src_ptr+n])
		return dest_ptr
	}
	self.registerApi(fs, "memmove", memmove)

	// int memcmp(const void *buf1, const void *buf2, unsigned int count);
	memcmp := func(proc *exec.Process, dest_ptr int32, src_ptr int32, n int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(dest_ptr+n < mem_len, "dest input pointer is exceed")
		assert.AssertEx(src_ptr+n < mem_len, "src input pointer is exceed")
		return int32(bytes.Compare(mem[dest_ptr:dest_ptr+n], mem[src_ptr:src_ptr+n]))
	}
	self.registerApi(fs, "memcmp", memcmp)
}

// bridge Authorization Api
func (self *bridge) bridgeAuthApi() {
	// extern "C" { bool isHexAddress(char* address);}
	isHexAddress := func(proc *exec.Process, addrss_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(addrss_ptr+43 < mem_len, "address input pointer is exceed")
		index := bytes.IndexByte(mem[addrss_ptr:], 0)
		assert.AssertEx(index == 42, "not a valid addr string")
		val := string(mem[addrss_ptr : addrss_ptr+42])
		assert.AssertEx(types.IsHexAddress(val), "not a hex address")
		return 0
	}
	fs := &wasm.FunctionSig{
		ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "isHexAddress", isHexAddress)

	// extern "C" { bool requireAuth(char* address);}
	requireAuth := func(proc *exec.Process, addrss_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(addrss_ptr+43 < mem_len, "address input pointer is exceed")
		assert.AssertEx(mem[addrss_ptr+42] == 0, "hex address too lang")
		val := string(mem[addrss_ptr : addrss_ptr+42])
		//assert.AssertEx(types.IsHexAddress(val), "not a hex address")
		authApi := actioncontext.AuthorizationApi{}
		authApi.SetCtx(self.rt.ctx)
		ret := authApi.RequireAuth(val)
		if ret {
			return 1
		}
		assert.AssertEx(false, "not have authority")
		return 0
	}
	self.registerApi(fs, "requireAuth", requireAuth)

	// extern "C" { bool isAccount(char* address);}
	isAccount := func(proc *exec.Process, addrss_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(addrss_ptr+43 < mem_len, "address input pointer is exceed")
		assert.AssertEx(mem[addrss_ptr+42] == 0, "hex address too lang")
		val := string(mem[addrss_ptr : addrss_ptr+42])
		assert.AssertEx(types.IsHexAddress(val), "not a hex address")
		authApi := actioncontext.AuthorizationApi{}
		authApi.SetCtx(self.rt.ctx)
		ret := authApi.IsAccount(val)
		if ret {
			return 1
		}
		return 0
	}
	self.registerApi(fs, "isAccount", isAccount)

	// extern "C" { bool isContract(char* address);}
	isContract := func(proc *exec.Process, addrss_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(addrss_ptr+43 < mem_len, "address input pointer is exceed")
		assert.AssertEx(mem[addrss_ptr+42] == 0, "hex address too lang")
		val := string(mem[addrss_ptr : addrss_ptr+42])
		assert.AssertEx(types.IsHexAddress(val), "not a hex address")
		authApi := actioncontext.AuthorizationApi{}
		authApi.SetCtx(self.rt.ctx)
		ret := authApi.IsContract(val)
		if ret {
			return 1
		}
		return 0
	}
	self.registerApi(fs, "IsContract", isContract)

	requireRewordAuth := func(proc *exec.Process) {
		memDbApi := actioncontext.BlockMemDbApi{}
		memDbApi.SetCtx(self.rt.ctx)
		memDbApi.Emplace([]byte("producer"), []byte("produced"))
	}
	fs = &wasm.FunctionSig{}
	self.registerApi(fs, "requireRewordAuth", requireRewordAuth)
}

// bridge System Api
func (self *bridge) bridgeSysApi() {
	// extern "C" { int hex2Int(char* hex);}
	hex2Int := func(proc *exec.Process, hex_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(hex_ptr < mem_len, "hex input pointer is exceed")
		index := bytes.IndexByte(mem[hex_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(hex_ptr+int32(index) < mem_len, "hex memory is exceed ")
		val := string(mem[hex_ptr : hex_ptr+int32(index)])
		str := types.FromHex(val)

		i, err := strconv.Atoi(string(str))
		assert.AsserErr(err)
		return int32(i)
	}
	fs := &wasm.FunctionSig{
		ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "hex2Int", hex2Int)
}

// bridge Console Api
func (self *bridge) bridgeConsoleApi() {
	// extern "C" { int log(char* msg);}
	printf := func(proc *exec.Process, msg_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(msg_ptr < mem_len, "msg input pointer is exceed")
		index := bytes.IndexByte(mem[msg_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(msg_ptr+int32(index) < mem_len, "msg memory is exceed ")
		fmt.Printf("WsamRE: %s", string(mem[msg_ptr:msg_ptr+int32(index)]))
		return 0
	}
	fs := &wasm.FunctionSig{
		ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "log", printf)
}

// bridge Assert Api
func (self *bridge) bridgeAssertApi() {
	// extern "C" { int assert(bool cond, char* assertMsg);}
	assertMethod := func(proc *exec.Process, cond int32, msg_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(msg_ptr < mem_len, "msg input pointer is exceed")
		index := bytes.IndexByte(mem[msg_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(msg_ptr+int32(index) < mem_len, "msg memory is exceed ")
		msg := string(mem[msg_ptr : msg_ptr+int32(index)])
		test := true
		if cond == 0 {
			test = false
		}
		assertApi := actioncontext.AssertApi{}
		assertApi.SetCtx(self.rt.ctx)
		assertApi.Assert(test, msg)
	}
	fs := &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
	}
	self.registerApi(fs, "assert", assertMethod)
}

// bridge Crypto Api
func (self *bridge) bridgeCryptoApi() {
	// extern "C" { void sha1(char* in, int len, char* out);}
	sha1 := func(proc *exec.Process, data_ptr int32, data_len int32, out_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(data_ptr+data_len < mem_len, "data input pointer is exceed")
		assert.AssertEx(out_ptr+43 < mem_len, "data output pointer is exceed")

		val := mem[data_ptr : data_ptr+data_len]

		cryptoApi := actioncontext.CryptoApi{}
		cryptoApi.SetCtx(self.rt.ctx)
		shaVal := cryptoApi.Sha1(val)
		hex := types.ToHex(shaVal[:])
		copy(mem[out_ptr:out_ptr+43], append([]byte(hex), 0))
	}
	fs := &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
	}
	self.registerApi(fs, "sha1", sha1)

	// extern "C" { char* sha256(char* msg);}
	sha256 := func(proc *exec.Process, data_ptr int32, data_len int32, out_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(data_ptr+data_len < mem_len, "data input pointer is exceed")
		assert.AssertEx(out_ptr+65 < mem_len, "data output pointer is exceed")

		val := mem[data_ptr : data_ptr+data_len]

		cryptoApi := actioncontext.CryptoApi{}
		cryptoApi.SetCtx(self.rt.ctx)
		shaVal := cryptoApi.Sha256([]byte(val))
		hex := types.ToHex(shaVal[:])
		copy(mem[out_ptr:out_ptr+65], append([]byte(hex), 0))
	}
	self.registerApi(fs, "sha256", sha256)

	sha512 := func(proc *exec.Process, data_ptr int32, data_len int32, out_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(data_ptr+data_len < mem_len, "data input pointer is exceed")
		assert.AssertEx(out_ptr+129 < mem_len, "data output pointer is exceed")

		val := mem[data_ptr : data_ptr+data_len]

		cryptoApi := actioncontext.CryptoApi{}
		cryptoApi.SetCtx(self.rt.ctx)
		shaVal := cryptoApi.Sha512([]byte(val))
		hex := types.ToHex(shaVal[:])
		copy(mem[out_ptr:out_ptr+129], append([]byte(hex), 0))
	}
	self.registerApi(fs, "sha512", sha512)

	recover := func(proc *exec.Process, data_ptr int32, sig_ptr int32, out_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(data_ptr < mem_len, "data input pointer is exceed")
		assert.AssertEx(sig_ptr < mem_len, "data output pointer is exceed")
		assert.AssertEx(out_ptr+43 < mem_len, "data output pointer is exceed")
		index := bytes.IndexByte(mem[data_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		msg := string(mem[data_ptr : data_ptr+int32(index)])
		assert.AssertEx(types.HasHexPrefix(msg), "no hex prefix")
		assert.AssertEx(types.IsHex(msg[2:]), "not a hex string")

		index = bytes.IndexByte(mem[data_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		sig := string(mem[data_ptr : data_ptr+int32(index)])
		assert.AssertEx(types.HasHexPrefix(sig), "no hex prefix")
		assert.AssertEx(types.IsHex(sig[2:]), "not a hex string")

		cryptoApi := actioncontext.CryptoApi{}
		cryptoApi.SetCtx(self.rt.ctx)
		address, err := cryptoApi.Recover(types.FromHex(msg), types.FromHex(sig))
		assert.AsserErr(err)
		hex := types.ToHex(address[:])
		copy(mem[out_ptr:out_ptr+43], append([]byte(hex), 0))
	}
	self.registerApi(fs, "recover", recover)
}

// bridge Producer Api
func (self *bridge) bridgeProducerApi() {
	// TODO:
	block_producer := func(proc *exec.Process, out_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(out_ptr+43 < mem_len, "data output pointer is exceed")
		sysApi := actioncontext.SystemApi{}
		sysApi.SetCtx(self.rt.ctx)
		hex := sysApi.GetBlockMiner().HexLower()
		copy(mem[out_ptr:out_ptr+43], append([]byte(hex), 0))
	}
	fs := &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "block_producer", block_producer)

	block_number := func(proc *exec.Process) int64 {
		sysApi := actioncontext.SystemApi{}
		sysApi.SetCtx(self.rt.ctx)
		return int64(sysApi.GetBlockNumber().Uint64())
	}

	fs_n := &wasm.FunctionSig{
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
	}
	self.registerApi(fs_n, "block_number", block_number)
}

// bridge weight Api
func (self *bridge) bridgeWeightApi() {
	getWeight := func(proc *exec.Process, amount int64) int64 {
		sysApi := actioncontext.SystemApi{}
		sysApi.SetCtx(self.rt.ctx)
		blkNum := sysApi.GetBlockNumber().Uint64()
		limit := int64(500000000000)
		if blkNum < 3650000 {
			limit = 10000000000
		}
		if amount < limit {
			return 0
		}
		return int64(math.Pow(float64(amount),0.33))
	}

	fs := &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI64},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
	}
	self.registerApi(fs, "getWeight", getWeight)
}

// bridge bigInt Api
func (self *bridge) bridgeBigIntApi() {
	big_add := func(proc *exec.Process, a_ptr int32, b_ptr int32, ret_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))

		assert.AssertEx(a_ptr < mem_len, "a_ptr input pointer is exceed")
		index := bytes.IndexByte(mem[a_ptr:], 0)
		assert.AssertEx(index != -1 && index<124, "a_ptr not a valid string")
		assert.AssertEx(a_ptr+int32(index) < mem_len, "a_ptr memory is exceed ")
		a_str := mem[a_ptr:a_ptr+int32(index)]
		a, ok := new(big.Int).SetString(string(a_str), 10)
		if !ok {
			a = big.NewInt(0)
		}

		assert.AssertEx(b_ptr < mem_len, "a_ptr input pointer is exceed")
		index = bytes.IndexByte(mem[b_ptr:], 0)
		assert.AssertEx(index != -1 && index<124, "a_ptr not a valid string")
		assert.AssertEx(b_ptr+int32(index) < mem_len, "a_ptr memory is exceed ")
		b_str := mem[b_ptr:b_ptr+int32(index)]
		b, ok := new(big.Int).SetString(string(b_str), 10)
		if !ok {
			b = big.NewInt(0)
		}

		assert.AssertEx(ret_ptr+128 < mem_len, "ret_ptr input pointer is exceed")

		c := new(big.Int).Add(a,b)
		ret := len(c.String())
		assert.AssertEx(ret < 124, "big_add ret len exceed")
		copy(mem[ret_ptr:], c.String())
		//fmt.Println(a, b,c)
		return int32(len(c.String()))
	}


	fs := &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "big_add", big_add)

	big_sub := func(proc *exec.Process, a_ptr int32, b_ptr int32, ret_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))

		assert.AssertEx(a_ptr < mem_len, "a_ptr input pointer is exceed")
		index := bytes.IndexByte(mem[a_ptr:], 0)
		assert.AssertEx(index != -1 && index<124, "a_ptr not a valid string")
		assert.AssertEx(a_ptr+int32(index) < mem_len, "a_ptr memory is exceed ")
		a_str := mem[a_ptr:a_ptr+int32(index)]
		a, ok := new(big.Int).SetString(string(a_str), 10)
		if !ok {
			a = big.NewInt(0)
		}

		assert.AssertEx(b_ptr < mem_len, "a_ptr input pointer is exceed")
		index = bytes.IndexByte(mem[b_ptr:], 0)
		assert.AssertEx(index != -1 && index<124, "a_ptr not a valid string")
		assert.AssertEx(b_ptr+int32(index) < mem_len, "a_ptr memory is exceed ")
		b_str := mem[b_ptr:b_ptr+int32(index)]
		b, ok := new(big.Int).SetString(string(b_str), 10)
		if !ok {
			b = big.NewInt(0)
		}

		assert.AssertEx(ret_ptr+128 < mem_len, "ret_ptr input pointer is exceed")

		c := new(big.Int).Sub(a,b)
		copy(mem[ret_ptr:], c.String())
		//fmt.Println(a, b,c)
		return int32(len(c.String()))
	}

	fs = &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "big_sub", big_sub)

	big_sub_safe := func(proc *exec.Process, a_ptr int32, b_ptr int32, ret_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))

		assert.AssertEx(a_ptr < mem_len, "a_ptr input pointer is exceed")
		index := bytes.IndexByte(mem[a_ptr:], 0)
		assert.AssertEx(index != -1 && index<124, "a_ptr not a valid string")
		assert.AssertEx(a_ptr+int32(index) < mem_len, "a_ptr memory is exceed ")
		a_str := mem[a_ptr:a_ptr+int32(index)]
		a, ok := new(big.Int).SetString(string(a_str), 10)
		if !ok {
			a = big.NewInt(0)
		}

		assert.AssertEx(b_ptr < mem_len, "a_ptr input pointer is exceed")
		index = bytes.IndexByte(mem[b_ptr:], 0)
		assert.AssertEx(index != -1 && index<124, "a_ptr not a valid string")
		assert.AssertEx(b_ptr+int32(index) < mem_len, "a_ptr memory is exceed ")
		b_str := mem[b_ptr:b_ptr+int32(index)]
		b, ok := new(big.Int).SetString(string(b_str), 10)
		if !ok {
			b = big.NewInt(0)
		}

		assert.AssertEx(ret_ptr+128 < mem_len, "ret_ptr input pointer is exceed")

		c := new(big.Int).Sub(a,b)
		assert.AssertEx(c.Sign() >= 0, "big_sub: a < b")
		copy(mem[ret_ptr:], c.String())
		//fmt.Println(a, b,c)
		return int32(len(c.String()))
	}

	fs = &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "big_sub_safe", big_sub_safe)

	big_exp_safe := func(proc *exec.Process, a_ptr int32, b_ptr int32, ret_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))

		assert.AssertEx(a_ptr < mem_len, "a_ptr input pointer is exceed")
		index := bytes.IndexByte(mem[a_ptr:], 0)
		assert.AssertEx(index != -1 && index<124, "a_ptr not a valid string")
		assert.AssertEx(a_ptr+int32(index) < mem_len, "a_ptr memory is exceed ")
		a_str := mem[a_ptr:a_ptr+int32(index)]
		a, ok := new(big.Int).SetString(string(a_str), 10)
		assert.AssertEx(ok, "a is not valid")

		assert.AssertEx(b_ptr < mem_len, "a_ptr input pointer is exceed")
		index = bytes.IndexByte(mem[b_ptr:], 0)
		assert.AssertEx(index != -1 && index<124, "a_ptr not a valid string")
		assert.AssertEx(b_ptr+int32(index) < mem_len, "a_ptr memory is exceed ")
		b_str := mem[b_ptr:b_ptr+int32(index)]
		b, ok := new(big.Int).SetString(string(b_str), 10)
		assert.AssertEx(ok, "b is not valid")
		assert.AssertEx(b.Sign()>0, "b is not valid")

		assert.AssertEx(ret_ptr+128 < mem_len, "ret_ptr input pointer is exceed")

		c := new(big.Int).Exp(a,b,nil)
		ret := len(c.String())
		assert.AssertEx(ret < 124, "big_exp_safe ret len exceed")
		copy(mem[ret_ptr:], c.String())
		//fmt.Println(a, b,c)
		return int32(ret)
	}

	fs = &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "big_exp_safe", big_exp_safe)

	big_mul_safe := func(proc *exec.Process, a_ptr int32, b_ptr int32, ret_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))

		assert.AssertEx(a_ptr < mem_len, "a_ptr input pointer is exceed")
		index := bytes.IndexByte(mem[a_ptr:], 0)
		assert.AssertEx(index != -1 && index<124, "a_ptr not a valid string")
		assert.AssertEx(a_ptr+int32(index) < mem_len, "a_ptr memory is exceed ")
		a_str := mem[a_ptr:a_ptr+int32(index)]
		a, ok := new(big.Int).SetString(string(a_str), 10)
		assert.AssertEx(ok, "a is not valid")
		assert.AssertEx(a.Sign()>0, "a is not valid")

		assert.AssertEx(b_ptr < mem_len, "a_ptr input pointer is exceed")
		index = bytes.IndexByte(mem[b_ptr:], 0)
		assert.AssertEx(index != -1 && index<124, "a_ptr not a valid string")
		assert.AssertEx(b_ptr+int32(index) < mem_len, "a_ptr memory is exceed ")
		b_str := mem[b_ptr:b_ptr+int32(index)]
		b, ok := new(big.Int).SetString(string(b_str), 10)
		assert.AssertEx(ok, "b is not valid")
		assert.AssertEx(b.Sign()>0, "b is not valid")

		assert.AssertEx(ret_ptr+128 < mem_len, "ret_ptr input pointer is exceed")

		c := new(big.Int).Mul(a, b)
		ret := len(c.String())
		assert.AssertEx(ret < 124, "big_mul_safe ret len exceed")
		copy(mem[ret_ptr:], c.String())
		//fmt.Println(a, b,c)
		return int32(ret)
	}

	fs = &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "big_mul_safe", big_mul_safe)
}

// bridge Action Api
func (self *bridge) bridgeActApi() {
	// TODO:
	action_sender := func(proc *exec.Process, out_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(out_ptr+43 < mem_len, "data output pointer is exceed")
		actApi := actioncontext.ActionApi{}
		actApi.SetCtx(self.rt.ctx)
		hex := actApi.Sender().HexLower()
		copy(mem[out_ptr:out_ptr+43], append([]byte(hex), 0))
	}
	fs := &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "action_sender", action_sender)

	contract_address := func(proc *exec.Process, out_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(out_ptr+43 < mem_len, "data output pointer is exceed")
		actApi := actioncontext.ActionApi{}
		actApi.SetCtx(self.rt.ctx)
		hex := actApi.Contract().Hex()
		copy(mem[out_ptr:out_ptr+43], append([]byte(hex), 0))
	}
	self.registerApi(fs, "contract_address", contract_address)

	contract_creator := func(proc *exec.Process, out_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(out_ptr+43 < mem_len, "data output pointer is exceed")
		conApi := actioncontext.ContractApi{}
		conApi.SetCtx(self.rt.ctx)
		hex := conApi.Creator().Hex()
		copy(mem[out_ptr:out_ptr+43], append([]byte(hex), 0))
	}
	self.registerApi(fs, "contract_creator", contract_creator)
}

// bridge Database Api
func (self *bridge) bridgeDbApi() {
	emplace := func(proc *exec.Process, key_ptr int32, key_len int32, val_ptr int32, val_len int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(key_ptr+key_len < mem_len, "data input pointer is exceed")
		assert.AssertEx(val_ptr+val_len < mem_len, "data output pointer is exceed")
		assert.AssertEx(key_len != 0, "key len is zero")
		assert.AssertEx(val_len != 0, "val len is zero")
		key := mem[key_ptr : key_ptr+key_len]
		val := mem[val_ptr : val_ptr+val_len]

		dbApi := actioncontext.DatabaseApi{}
		dbApi.SetCtx(self.rt.ctx)
		dbApi.Emplace(key, val)
	}
	fs := &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
	}
	self.registerApi(fs, "db_emplace", emplace)

	//Set value
	set := func(proc *exec.Process, key_ptr int32, key_len int32, val_ptr int32, val_len int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(key_ptr+key_len < mem_len, "data input pointer is exceed")
		assert.AssertEx(val_ptr+val_len < mem_len, "data output pointer is exceed")
		assert.AssertEx(key_len != 0, "key len is zero")
		assert.AssertEx(val_len != 0, "val len is zero")
		key := mem[key_ptr : key_ptr+key_len]
		val := mem[val_ptr : val_ptr+val_len]

		dbApi := actioncontext.DatabaseApi{}
		dbApi.SetCtx(self.rt.ctx)
		dbApi.Set(key, val)
	}
	self.registerApi(fs, "db_set", set)

	//Get value
	get := func(proc *exec.Process, key_ptr int32, key_len int32, val_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(key_ptr+key_len < mem_len, "data input pointer is exceed")
		assert.AssertEx(val_ptr < mem_len, "data output pointer is exceed")
		assert.AssertEx(key_len != 0, "key len is zero")
		key := mem[key_ptr : key_ptr+key_len]

		dbApi := actioncontext.DatabaseApi{}
		dbApi.SetCtx(self.rt.ctx)
		val := dbApi.Get(key)
		if nil == val {
			return 0
		}
		val_len := int32(len(val))
		assert.AssertEx(val_ptr+val_len < mem_len, "data output pointer is exceed")
		copy(mem[val_ptr:val_ptr+val_len], val)

		return val_len
	}
	fs_get := &wasm.FunctionSig{
		ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs_get, "db_get", get)
}

// bridge Cache Database Api
func (self *bridge) bridgeCacheApi() {
	emplace := func(proc *exec.Process, key_ptr int32, key_len int32, val_ptr int32, val_len int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(key_ptr+key_len < mem_len, "data input pointer is exceed")
		assert.AssertEx(val_ptr+val_len < mem_len, "data output pointer is exceed")
		assert.AssertEx(key_len != 0, "key len is zero")
		assert.AssertEx(val_len != 0, "val len is zero")
		key := mem[key_ptr : key_ptr+key_len]
		val := mem[val_ptr : val_ptr+val_len]

		memDbApi := actioncontext.BlockMemDbApi{}
		memDbApi.SetCtx(self.rt.ctx)
		memDbApi.Emplace(key, val)
	}
	fs := &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
	}
	self.registerApi(fs, "cacheDb_emplace", emplace)

	//Get value
	get := func(proc *exec.Process, key_ptr int32, key_len int32, val_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(key_ptr+key_len < mem_len, "data input pointer is exceed")
		assert.AssertEx(val_ptr < mem_len, "data output pointer is exceed")
		assert.AssertEx(key_len != 0, "key len is zero")
		key := mem[key_ptr : key_ptr+key_len]

		memDbApi := actioncontext.BlockMemDbApi{}
		memDbApi.SetCtx(self.rt.ctx)
		val := memDbApi.Get(key)
		if nil == val {
			return 0
		}
		val_len := int32(len(val))
		assert.AssertEx(val_ptr+val_len < mem_len, "data output pointer is exceed")
		copy(mem[val_ptr:val_ptr+val_len], val)

		return val_len
	}
	fs_get := &wasm.FunctionSig{
		ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs_get, "cacheDb_get", get)
}

// bridge Contract Api
func (self *bridge) bridgeContractApi() {
	create := func(proc *exec.Process, creator_ptr int32, code_ptr int32, out_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(out_ptr+43 < mem_len, "output pointer is exceed")
		assert.AssertEx(creator_ptr < mem_len, "creator input pointer is exceed")
		assert.AssertEx(code_ptr < mem_len, "code input pointer is exceed")
		index := bytes.IndexByte(mem[creator_ptr:], 0)
		assert.AssertEx(index != -1, "creator not a valid string")
		assert.AssertEx(creator_ptr+int32(index) < mem_len, "creator memory is exceed ")
		index_c := bytes.IndexByte(mem[code_ptr:], 0)
		assert.AssertEx(index_c != -1, "code not a valid string")
		assert.AssertEx(code_ptr+int32(index_c) < mem_len, "code memory is exceed ")

		conApi := actioncontext.ContractApi{}
		conApi.SetCtx(self.rt.ctx)
		addr := conApi.Create(string(mem[creator_ptr:creator_ptr+int32(index)]), string(mem[code_ptr:code_ptr+int32(index_c)]))
		addrStr := addr.HexLower()
		copy(mem[out_ptr:out_ptr+43], append([]byte(addrStr), 0))
	}
	fs := &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
	}
	self.registerApi(fs, "contract_create", create)

	// void emitEvent(stTopic* topics, stData* datas);
	//typedef struct stTopic {
	//	char topic[32];
	//	stTopic* next;
	//} stTopic;
	//
	//	typedef struct stData {
	//	char*  data;
	//	int data_len;
	//	stData* next;
	//} stData;
	emitEvent := func(proc *exec.Process, topics_ptr int32, data_ptr int32) {
		var topics []string
		var data [][]byte
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))

		assert.AssertEx(topics_ptr != 0, "not exit a topic")

		for topicCc, next := 0, topics_ptr; topicCc < 4; topicCc++ {
			if next == 0 {
				break
			}
			assert.AssertEx(next+12 < mem_len, "stData struct pointer is exceed")
			val_ptr := binary.LittleEndian.Uint32(mem[next : next+4])
			assert.AssertEx(val_ptr != 0, "data v pointer is NULL")
			assert.AssertEx(int32(val_ptr)+4 < mem_len, "data v pointer is exceed")
			val_len := binary.LittleEndian.Uint32(mem[next+4 : next+8])
			assert.AssertEx(val_len+val_ptr < uint32(mem_len), "data len is exceed")
			topic := mem[val_ptr : val_ptr+val_len]
			topics = append(topics, string(topic))
			next = int32(binary.LittleEndian.Uint32(mem[next+8 : next+12]))
		}

		for dataCc, next := 0, data_ptr; dataCc < 6; dataCc++ {
			if next == 0 {
				break
			}
			assert.AssertEx(next+12 < mem_len, "stData struct pointer is exceed")
			val_ptr := binary.LittleEndian.Uint32(mem[next : next+4])
			assert.AssertEx(val_ptr != 0, "data v pointer is NULL")
			assert.AssertEx(int32(val_ptr)+4 < mem_len, "data v pointer is exceed")
			val_len := binary.LittleEndian.Uint32(mem[next+4 : next+8])
			assert.AssertEx(val_len+val_ptr < uint32(mem_len), "data len is exceed")
			val := mem[val_ptr : val_ptr+val_len]
			data = append(data, val)
			next = int32(binary.LittleEndian.Uint32(mem[next+8 : next+12]))
		}

		conApi := actioncontext.ContractApi{}
		conApi.SetCtx(self.rt.ctx)
		conApi.EmitEvent(topics, data)
	}
	fs_emit := &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
	}
	self.registerApi(fs_emit, "emitEvent", emitEvent)
}

func (self *bridge) bridgeResultApi() {
	setResult := func(proc *exec.Process, data_ptr int32, data_len int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))

		assert.AssertEx(data_ptr+data_len < mem_len, "data input pointer is exceed")
		assert.AssertEx(data_len != 0, "data len is zero")

		data := mem[data_ptr : data_ptr+data_len]

		rltApi := actioncontext.ResultApi{}
		rltApi.SetCtx(self.rt.ctx)
		rltApi.SetActionResult(data)
	}
	fs := &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
	}
	self.registerApi(fs, "setResult", setResult)
}

func convertType(apiType uint32) string {
	if apiType == 0 {
		return para_paser.TypeI32
	} else if apiType == 1 {
		return para_paser.TypeI64
	} else if apiType == 2 {
		return para_paser.TypeF32
	} else if apiType == 3 {
		return para_paser.TypeF64
	} else if apiType == 4 {
		return para_paser.TypeAddress
	} else {
		panic("apiType is not in range")
	}
}

// bridge Call Api
func (self *bridge) bridgeCallApi() {
	call := func(proc *exec.Process, addr_ptr int32, para_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(addr_ptr < mem_len, "address input pointer is exceed")
		index := bytes.IndexByte(mem[addr_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(addr_ptr+int32(index) < mem_len, "address memory is exceed ")
		addrStr := mem[addr_ptr : addr_ptr+int32(index)]
		address := types.HexToAddress(string(addrStr))
		nilAddress := types.Address{}
		assert.AssertEx(address != nilAddress, "not a valid address string")

		assert.AssertEx(para_ptr < mem_len, "address input pointer is exceed")
		index = bytes.IndexByte(mem[para_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(para_ptr+int32(index) < mem_len, "address memory is exceed ")
		para := mem[para_ptr : para_ptr+int32(index)]

		act := &transaction.Action{address, para}
		callApi := actioncontext.CallApi{}
		callApi.SetCtx(self.rt.ctx)
		callApi.SetInterpreter(interpreter.Singleton())
		callApi.Call(act)
	}
	fs := &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
	}
	self.registerApi(fs, "action_call", call)

	contractCall := func(proc *exec.Process, addr_ptr int32, para_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(addr_ptr < mem_len, "address input pointer is exceed")
		index := bytes.IndexByte(mem[addr_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(addr_ptr+int32(index) < mem_len, "address memory is exceed ")
		addrStr := mem[addr_ptr : addr_ptr+int32(index)]
		address := types.HexToAddress(string(addrStr))
		nilAddress := types.Address{}
		assert.AssertEx(address != nilAddress, "not a valid address string")

		assert.AssertEx(para_ptr < mem_len, "address input pointer is exceed")
		index = bytes.IndexByte(mem[para_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(para_ptr+int32(index) < mem_len, "address memory is exceed ")
		para := mem[para_ptr : para_ptr+int32(index)]

		act := &transaction.Action{address, para}
		callApi := actioncontext.CallApi{}
		callApi.SetCtx(self.rt.ctx)
		callApi.SetInterpreter(interpreter.Singleton())
		callApi.InnerCall(act)
	}
	fs = &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
	}
	self.registerApi(fs, "contract_call", contractCall)

	getCallPara := func(proc *exec.Process, funcName_ptr int32, para_ptr int32, out_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))

		assert.AssertEx(para_ptr != 0, "para_ptr is null")
		assert.AssertEx(out_ptr != 0, "out_ptr is null")
		assert.AssertEx(funcName_ptr < mem_len, "str input pointer is exceed")
		index := bytes.IndexByte(mem[funcName_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(funcName_ptr+int32(index) < mem_len, "str memory is exceed ")
		funcName := string(mem[funcName_ptr : funcName_ptr+int32(index)])

		args := []para_paser.Arg{}
		for dataCc, next := 0, para_ptr; dataCc < 16; dataCc++ {
			if next == 0 {
				break
			}
			arg := para_paser.Arg{}
			assert.AssertEx(next+16 < mem_len, "stCallPara struct pointer is exceed")
			paraType := binary.LittleEndian.Uint32(mem[next : next+4])
			arg.Type = convertType(paraType)
			val_ptr := binary.LittleEndian.Uint32(mem[next+4 : next+8])
			assert.AssertEx(val_ptr != 0, "data v pointer is NULL")
			assert.AssertEx(int32(val_ptr)+4 < mem_len, "data v pointer is exceed")
			val_len := binary.LittleEndian.Uint32(mem[next+8 : next+12])
			assert.AssertEx(val_len+val_ptr < uint32(mem_len), "data len is exceed")
			val := mem[val_ptr : val_ptr+val_len]
			arg.Data = val
			args = append(args, arg)
			next = int32(binary.LittleEndian.Uint32(mem[next+12 : next+16]))
		}

		wp := &para_paser.WasmPara{
			FuncName: funcName,
			Args:     args,
		}
		paraBytes, err := json.Marshal(wp)
		assert.AssertEx(err == nil, err.Error())
		outLen := int32(len(paraBytes))
		assert.AssertEx(para_ptr+outLen < mem_len, "out_ptr input pointer is exceed")
		copy(mem[out_ptr:out_ptr+outLen], paraBytes)
	}

	fs = &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
	}
	self.registerApi(fs, "getCallPara", getCallPara)


	callWithPara := func(proc *exec.Process, addr_ptr int32, funcName_ptr int32, para_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(addr_ptr < mem_len, "address input pointer is exceed")
		index := bytes.IndexByte(mem[addr_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(addr_ptr+int32(index) < mem_len, "address memory is exceed ")
		addrStr := mem[addr_ptr : addr_ptr+int32(index)]
		address := types.HexToAddress(string(addrStr))
		nilAddress := types.Address{}
		assert.AssertEx(address != nilAddress, "not a valid address string")

		assert.AssertEx(funcName_ptr != 0, "funcName_ptr is null")
		assert.AssertEx(funcName_ptr < mem_len, "str input pointer is exceed")
		index = bytes.IndexByte(mem[funcName_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(funcName_ptr+int32(index) < mem_len, "str memory is exceed ")
		funcName := string(mem[funcName_ptr : funcName_ptr+int32(index)])

		assert.AssertEx(para_ptr != 0, "para_ptr is null")
		args := []para_paser.Arg{}
		for dataCc, next := 0, para_ptr; dataCc < 16; dataCc++ {
			if next == 0 {
				break
			}
			arg := para_paser.Arg{}
			assert.AssertEx(next+16 < mem_len, "stCallPara struct pointer is exceed")
			paraType := binary.LittleEndian.Uint32(mem[next : next+4])
			arg.Type = convertType(paraType)
			val_ptr := binary.LittleEndian.Uint32(mem[next+4 : next+8])
			assert.AssertEx(val_ptr != 0, "data v pointer is NULL")
			assert.AssertEx(int32(val_ptr)+4 < mem_len, "data v pointer is exceed")
			val_len := binary.LittleEndian.Uint32(mem[next+8 : next+12])
			assert.AssertEx(val_len+val_ptr < uint32(mem_len), "data len is exceed")
			val := mem[val_ptr : val_ptr+val_len]
			arg.Data = val
			args = append(args, arg)
			next = int32(binary.LittleEndian.Uint32(mem[next+12 : next+16]))
		}
		//fmt.Println(funcName, args)
		wp := &para_paser.WasmPara{
			FuncName: funcName,
			Args:     args,
		}
		paraBytes, err := json.Marshal(wp)
		//fmt.Println(paraBytes)
		assert.AsserErr(err)

		act := &transaction.Action{address, paraBytes}
		callApi := actioncontext.CallApi{}
		callApi.SetCtx(self.rt.ctx)
		callApi.SetInterpreter(interpreter.Singleton())
		callApi.Call(act)
	}
	fs = &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
	}
	self.registerApi(fs, "action_callWithPara", callWithPara)

	contract_callWithPara := func(proc *exec.Process, addr_ptr int32, funcName_ptr int32, para_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(addr_ptr < mem_len, "address input pointer is exceed")
		index := bytes.IndexByte(mem[addr_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(addr_ptr+int32(index) < mem_len, "address memory is exceed ")
		addrStr := mem[addr_ptr : addr_ptr+int32(index)]
		address := types.HexToAddress(string(addrStr))
		nilAddress := types.Address{}
		assert.AssertEx(address != nilAddress, "not a valid address string")

		assert.AssertEx(funcName_ptr != 0, "funcName_ptr is null")
		assert.AssertEx(funcName_ptr < mem_len, "str input pointer is exceed")
		index = bytes.IndexByte(mem[funcName_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(funcName_ptr+int32(index) < mem_len, "str memory is exceed ")
		funcName := string(mem[funcName_ptr : funcName_ptr+int32(index)])

		assert.AssertEx(para_ptr != 0, "para_ptr is null")
		args := []para_paser.Arg{}
		for dataCc, next := 0, para_ptr; dataCc < 16; dataCc++ {
			if next == 0 {
				break
			}
			arg := para_paser.Arg{}
			assert.AssertEx(next+16 < mem_len, "stCallPara struct pointer is exceed")
			paraType := binary.LittleEndian.Uint32(mem[next : next+4])
			arg.Type = convertType(paraType)
			val_ptr := binary.LittleEndian.Uint32(mem[next+4 : next+8])
			assert.AssertEx(val_ptr != 0, "data v pointer is NULL")
			assert.AssertEx(int32(val_ptr)+4 < mem_len, "data v pointer is exceed")
			val_len := binary.LittleEndian.Uint32(mem[next+8 : next+12])
			assert.AssertEx(val_len+val_ptr < uint32(mem_len), "data len is exceed")
			val := mem[val_ptr : val_ptr+val_len]
			arg.Data = val
			args = append(args, arg)
			next = int32(binary.LittleEndian.Uint32(mem[next+12 : next+16]))
		}

		wp := &para_paser.WasmPara{
			FuncName: funcName,
			Args:     args,
		}
		paraBytes, err := json.Marshal(wp)
		assert.AsserErr(err)

		act := &transaction.Action{address, paraBytes}
		callApi := actioncontext.CallApi{}
		callApi.SetCtx(self.rt.ctx)
		callApi.SetInterpreter(interpreter.Singleton())
		callApi.InnerCall(act)
	}
	fs = &wasm.FunctionSig{
		ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
	}
	self.registerApi(fs, "contract_callWithPara", contract_callWithPara)
}

func (self *bridge) bridgeStringApi() {
	strlen := func(proc *exec.Process, msg_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(msg_ptr < mem_len, "str input pointer is exceed")
		index := bytes.IndexByte(mem[msg_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(msg_ptr+int32(index) < mem_len, "str memory is exceed ")
		return int32(index)
	}
	fs := &wasm.FunctionSig{
		ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "strlen", strlen)

	strjoint := func(proc *exec.Process, fitsrt_ptr int32, second_ptr int32, out_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(fitsrt_ptr < mem_len, "first input pointer is exceed")
		assert.AssertEx(second_ptr < mem_len, "second input pointer is exceed")
		index := bytes.IndexByte(mem[fitsrt_ptr:], 0)
		assert.AssertEx(index != -1, "first str not a valid string")
		assert.AssertEx(fitsrt_ptr+int32(index) < mem_len, "first str memory is exceed")
		index_c := bytes.IndexByte(mem[second_ptr:], 0)
		assert.AssertEx(index_c != -1, "second str not a valid string")
		assert.AssertEx(second_ptr+int32(index_c) < mem_len, "second str memory is exceed")
		assert.AssertEx(int32(index_c)+int32(index)+out_ptr < mem_len, "out_ptr memory is exceed")

		copy(mem[out_ptr:out_ptr+int32(index)], mem[fitsrt_ptr:fitsrt_ptr+int32(index)])
		copy(mem[out_ptr+int32(index):out_ptr+int32(index)+int32(index_c)], mem[second_ptr:second_ptr+int32(index_c)])
		mem[out_ptr+out_ptr+int32(index)+int32(index_c)] = 0
		return int32(index) + int32(index_c)
	}
	fs = &wasm.FunctionSig{
		ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "strjoint", strjoint)

	str2lower := func(proc *exec.Process, msg_ptr int32) {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(msg_ptr < mem_len, "str input pointer is exceed")
		index := bytes.IndexByte(mem[msg_ptr:], 0)
		assert.AssertEx(index != -1, "not a valid string")
		assert.AssertEx(msg_ptr+int32(index) < mem_len, "str memory is exceed")

		s := mem[msg_ptr:msg_ptr+int32(index)]
		for i := 0; i < index; i++ {
			c := s[i]
			if c >= 'A' && c <= 'Z' {
				s[i] += 'a' - 'A'
			}
		}
	}
	fs = &wasm.FunctionSig{
		ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "str2lower", str2lower)


	strcmp := func(proc *exec.Process, fitsrt_ptr int32, second_ptr int32) int32 {
		mem := self.rt.vm.Memory()
		mem_len := int32(len(mem))
		assert.AssertEx(fitsrt_ptr < mem_len, "first input pointer is exceed")
		assert.AssertEx(second_ptr < mem_len, "second input pointer is exceed")
		index := bytes.IndexByte(mem[fitsrt_ptr:], 0)
		assert.AssertEx(index != -1, "first str not a valid string")
		assert.AssertEx(fitsrt_ptr+int32(index) < mem_len, "first str memory is exceed")
		index_c := bytes.IndexByte(mem[second_ptr:], 0)
		assert.AssertEx(index_c != -1, "second str not a valid string")
		assert.AssertEx(second_ptr+int32(index_c) < mem_len, "second str memory is exceed")

		ret := bytes.Equal(mem[fitsrt_ptr:fitsrt_ptr+int32(index)], mem[second_ptr:second_ptr+int32(index_c)])
		if ret {
			return 0
		}
		return 1
	}
	fs = &wasm.FunctionSig{
		ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
		ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
	}
	self.registerApi(fs, "strcmp", strcmp)
}
