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
// @Date: 2018/07/31 14:48:31
////////////////////////////////////////////////////////////////////////////////

package jsre

import (
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
	"bchain.io/common/assert"
	"bchain.io/common/types"
	"bchain.io/core/actioncontext"
	"bchain.io/core/interpreter"
	"bchain.io/core/transaction"
	"strconv"
)

type bridge struct {
	rt   *Runtime
	apis *actioncontext.APIs
}

func newBridge(rt *Runtime) *bridge {
	apis := actioncontext.NewAPIs(rt.ctx)
	apis.SetInterpreter(interpreter.Singleton())
	return &bridge{
		rt:   rt,
		apis: apis,
	}
}

func (self *bridge) load() {
	self.bridgeVariable()
	self.bridgeApi()
}

// set the internal variable of contract context
func (self *bridge) bridgeVariable() {
	self.rt.vm.Set("_context_var_self", self.apis.Act.Contract().Hex())            // set contract address, hex string
	self.rt.vm.Set("_context_var_sender", self.apis.Act.Sender().Hex())            // set sender address, hex string
	self.rt.vm.Set("_context_var_creator", self.apis.Contract.Creator().Hex())     // set creator address, hex string
	self.rt.vm.Set("_context_var_number", self.apis.Sys.GetBlockNumber().String()) // set block number, string
	self.rt.vm.Set("_context_var_miner" , self.apis.Sys.GetBlockMiner().Hex())	  // set block miner , string
}

// bridge all api
func (self *bridge) bridgeApi() {
	self.bridgeAuthApi()
	self.bridgeSysApi()
	self.bridgeConsoleApi()
	self.bridgeAssertApi()
	self.bridgeCryptoApi()
	self.bridgeAssertApi()
	self.bridgeProducerApi()
	self.bridgeActApi()
	self.bridgeDbApi()
	self.bridgeCacheApi()
	self.bridgeContractApi()
	self.bridgeResultApi()
	self.bridgeCallApi()
}

// bridge Authorization Api
func (self *bridge) bridgeAuthApi() {
	self.rt.vm.Set("_context_api_auth", struct{}{})

	obj, _ := self.rt.vm.Get("_context_api_auth")
	isHexAddress := func(call otto.FunctionCall) otto.Value {
		val, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		assert.AssertEx(types.IsHexAddress(val), "not a hex address")
		//address := types.StringToAddress(val)
		oval, err := otto.ToValue(self.apis.Auth.IsHexAddress(val))
		assert.AsserErr(err)
		return oval
	}
	obj.Object().Set("isHexAddress", isHexAddress)

	requireAuth := func(call otto.FunctionCall) otto.Value {
		val, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		assert.AssertEx(types.IsHexAddress(val), "not a hex address")
		//address := types.StringToAddress(val)
		oval, err := otto.ToValue(self.apis.Auth.RequireAuth(val))
		assert.AsserErr(err)
		return oval
	}
	obj.Object().Set("requireAuth", requireAuth)

	IsAccount := func(call otto.FunctionCall) otto.Value {
		val, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		assert.AssertEx(types.IsHexAddress(val), "not a hex address")
		//address := types.StringToAddress(val)
		oval, err := otto.ToValue(self.apis.Auth.IsAccount(val))
		assert.AsserErr(err)
		return oval
	}
	obj.Object().Set("isAccount", IsAccount)

	IsContract := func(call otto.FunctionCall) otto.Value {
		val, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		assert.AssertEx(types.IsHexAddress(val), "not a hex address")
		//address := types.StringToAddress(val)
		oval, err := otto.ToValue(self.apis.Auth.IsContract(val))
		assert.AsserErr(err)
		return oval
	}
	obj.Object().Set("isContract", IsContract)
}

// bridge System Api
func (self *bridge) bridgeSysApi() {
	// TODO:
	self.rt.vm.Set("_context_api_sys", struct{}{})
	obj, _ := self.rt.vm.Get("_context_api_sys")

	hex2Int := func(call otto.FunctionCall) otto.Value {
		val, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		str := types.FromHex(val)

		i, err := strconv.Atoi(string(str))
		assert.AsserErr(err)
		oval, err := otto.ToValue(i)
		assert.AsserErr(err)
		return oval
	}
	obj.Object().Set("hex2int", hex2Int)
}

// bridge Console Api
func (self *bridge) bridgeConsoleApi() {
	// TODO:
	self.rt.vm.Set("_context_api_console", struct{}{})
	obj, _ := self.rt.vm.Get("_context_api_console")

	printf := func(call otto.FunctionCall) otto.Value {
		val, err := call.Argument(0).ToString()
		assert.AsserErr(err)

		self.apis.Console.Printf(val)
		return otto.Value{}
	}
	obj.Object().Set("print", printf)
}

// bridge Assert Api
func (self *bridge) bridgeAssertApi() {
	self.rt.vm.Set("_context_api_assert", struct{}{})
	obj, _ := self.rt.vm.Get("_context_api_assert")

	assertMethod := func(call otto.FunctionCall) otto.Value {
		test, err := call.Argument(0).ToBoolean()
		assert.AsserErr(err)

		msg, err := call.Argument(1).ToString()
		assert.AsserErr(err)

		self.apis.Assert.Assert(test, msg)
		return otto.Value{}
	}
	obj.Object().Set("assert", assertMethod)
}

// bridge Crypto Api
func (self *bridge) bridgeCryptoApi() {
	// TODO:
	self.rt.vm.Set("_context_api_crypto", struct{}{})
	obj, _ := self.rt.vm.Get("_context_api_crypto")

	sha1 := func(call otto.FunctionCall) otto.Value {
		val, err := call.Argument(0).ToString()
		assert.AsserErr(err)

		shaVal := self.apis.Crypto.Sha1([]byte(val))
		oval, err := otto.ToValue(types.ToHex(shaVal[:]))
		assert.AsserErr(err)
		return oval
	}
	obj.Object().Set("sha1", sha1)

	sha256 := func(call otto.FunctionCall) otto.Value {
		val, err := call.Argument(0).ToString()
		assert.AsserErr(err)

		shaVal := self.apis.Crypto.Sha256([]byte(val))
		oval, err := otto.ToValue(types.ToHex(shaVal[:]))
		assert.AsserErr(err)
		return oval
	}
	obj.Object().Set("sha256", sha256)

	sha512 := func(call otto.FunctionCall) otto.Value {
		val, err := call.Argument(0).ToString()
		assert.AsserErr(err)

		shaVal := self.apis.Crypto.Sha512([]byte(val))
		oval, err := otto.ToValue(types.ToHex(shaVal[:]))
		assert.AsserErr(err)
		return oval
	}
	obj.Object().Set("sha512", sha512)

	recover := func(call otto.FunctionCall) otto.Value {
		msg, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		assert.AssertEx(types.HasHexPrefix(msg), "no hex prefix")
		assert.AssertEx(types.IsHex(msg[2:]), "not a hex string")

		sig, err := call.Argument(1).ToString()
		assert.AsserErr(err)
		assert.AssertEx(types.HasHexPrefix(sig), "no hex prefix")
		assert.AssertEx(types.IsHex(sig[2:]), "not a hex string")

		key, err := self.apis.Crypto.Recover(types.FromHex(msg), types.FromHex(sig))
		assert.AsserErr(err)
		oval, err := otto.ToValue(types.Bytes2Hex(key))
		assert.AsserErr(err)
		return oval
	}
	obj.Object().Set("recover", recover)
}

// bridge Producer Api
func (self *bridge) bridgeProducerApi() {
	// TODO:
	self.rt.vm.Set("_context_api_producer", struct{}{})
}

// bridge Action Api
func (self *bridge) bridgeActApi() {
	// TODO:
	self.rt.vm.Set("_context_api_act", struct{}{})
}

func (self *bridge) bridgeResultApi() {
	// TODO:
	self.rt.vm.Set("_context_api_result", struct{}{})
	obj, _ := self.rt.vm.Get("_context_api_result")

	setResult := func(call otto.FunctionCall) otto.Value {

		valStr, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		fmt.Println("Len valStr:", len(valStr))
		fmt.Println(">>>>>>>>>>>>>action setResult:", valStr)
		//assert.AssertEx(types.HasHexPrefix(valStr), "no hex prefix")
		//assert.AssertEx(types.IsHex(valStr[2:]), "not a hex string")
		//val := types.FromHex(valStr)
		//todo:here use string to bytes just for test
		self.apis.Result.SetActionResult([]byte(valStr))
		return otto.Value{}
	}
	obj.Object().Set("setResult", setResult)

}

// bridge Database Api
func (self *bridge) bridgeDbApi() {
	// TODO:
	self.rt.vm.Set("_context_api_db", struct{}{})
	obj, _ := self.rt.vm.Get("_context_api_db")

	emplace := func(call otto.FunctionCall) otto.Value {
		keyStr, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		assert.AssertEx(types.HasHexPrefix(keyStr), "no hex prefix")
		assert.AssertEx(types.IsHex(keyStr[2:]), "not a hex string")
		key := types.FromHex(keyStr)

		valStr, err := call.Argument(1).ToString()
		assert.AsserErr(err)
		assert.AssertEx(types.HasHexPrefix(valStr), "no hex prefix")
		assert.AssertEx(types.IsHex(valStr[2:]), "not a hex string")
		val := types.FromHex(valStr)

		self.apis.Db.Emplace(key, val)
		return otto.Value{}
	}
	obj.Object().Set("emplace", emplace)

	//Set value
	set := func(call otto.FunctionCall) otto.Value {
		keyStr, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		//assert.AssertEx(types.HasHexPrefix(keyStr), "no hex prefix")
		//assert.AssertEx(types.IsHex(keyStr[2:]), "not a hex string")
		//now: "0x123xxxxx"+"balance"-->"0x123xxxxbalance",can not use from hex
		//now use string directly
		//key := types.FromHex(keyStr)

		valStr, err := call.Argument(1).ToString()
		assert.AsserErr(err)
		val := []byte(valStr)
		//self.apis.Db.Set(key , val)
		self.apis.Db.Set([]byte(keyStr), val)
		return otto.Value{}
	}
	obj.Object().Set("set", set)

	//Get value
	get := func(call otto.FunctionCall) otto.Value {
		keyStr, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		//assert.AssertEx(types.HasHexPrefix(keyStr), "no hex prefix")
		//assert.AssertEx(types.IsHex(keyStr[2:]), "not a hex string")
		//now: "0x123xxxxx"+"balance"-->"0x123xxxxbalance",can not use from hex
		//now use string directly
		//key := types.FromHex(keyStr)

		var retStr string
		ret := self.apis.Db.Get([]byte(keyStr))
		if nil == ret {
			fmt.Println("api.Db.Get ret == nil...")
			retStr = ""
		} else {
			//retStr = fmt.Sprintf("%x", ret)
			retStr = string(ret)
		}

		getResult, err := otto.ToValue(retStr)
		assert.AsserErr(err)
		return getResult

	}
	obj.Object().Set("get", get)

}

// bridge Database Api
func (self *bridge) bridgeCacheApi() {
	// TODO:
	self.rt.vm.Set("_context_api_cache", struct{}{})
	obj, _ := self.rt.vm.Get("_context_api_cache")

	//Set value
	emplace := func(call otto.FunctionCall) otto.Value {
		keyStr, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		//assert.AssertEx(types.HasHexPrefix(keyStr), "no hex prefix")
		//assert.AssertEx(types.IsHex(keyStr[2:]), "not a hex string")
		//now: "0x123xxxxx"+"balance"-->"0x123xxxxbalance",can not use from hex
		//now use string directly
		//key := types.FromHex(keyStr)

		valStr, err := call.Argument(1).ToString()
		assert.AsserErr(err)
		val := []byte(valStr)
		//self.apis.Db.Set(key , val)
		self.apis.MemDb.Emplace([]byte(keyStr), val)
		return otto.Value{}
	}
	obj.Object().Set("emplace", emplace)

	//Get value
	get := func(call otto.FunctionCall) otto.Value {
		keyStr, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		//assert.AssertEx(types.HasHexPrefix(keyStr), "no hex prefix")
		//assert.AssertEx(types.IsHex(keyStr[2:]), "not a hex string")
		//now: "0x123xxxxx"+"balance"-->"0x123xxxxbalance",can not use from hex
		//now use string directly
		//key := types.FromHex(keyStr)

		var retStr string
		ret := self.apis.MemDb.Get([]byte(keyStr))
		if nil == ret {
			fmt.Println("api.Db.Get ret == nil...")
			retStr = ""
		} else {
			//retStr = fmt.Sprintf("%x", ret)
			retStr = string(ret)
		}

		getResult, err := otto.ToValue(retStr)
		assert.AsserErr(err)
		return getResult

	}
	obj.Object().Set("get", get)

}

// bridge Contract Api
func (self *bridge) bridgeContractApi() {
	self.rt.vm.Set("_context_api_contract", struct{}{})
	obj, _ := self.rt.vm.Get("_context_api_contract")

	create := func(call otto.FunctionCall) otto.Value {
		creatorStr, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		codeStr, err := call.Argument(1).ToString()
		assert.AsserErr(err)

		addr := self.apis.Contract.Create(creatorStr, codeStr)
		addrStr := addr.Hex()
		oval, err := otto.ToValue(addrStr)
		assert.AsserErr(err)

		//standardApi := GetJsStandardApis()

		//act := &transaction.Action{addr, []byte(standardApi)}
		//fmt.Println("bridgeRunApi CreateContract Addr", addr.Hex())
		//self.apis.Call.Call(act)
		//fmt.Println("After Call Action")
		return oval
	}
	obj.Object().Set("create", create)

	//parameters must be even number, parameter pair <first, second>, first is the data for record, second is a boolean
	//value, true indicate the first to be a topic, false to be a data
	emitEvent := func(call otto.FunctionCall) otto.Value {
		paramCount := len(call.ArgumentList)
		if paramCount%2 != 0 {
			assert.AsserErr(errors.New("the api of emitEvent must have even number of parameters"))
		}
		var topics []string
		var data [][]byte
		for i := 0; i < paramCount; i = i + 2 {
			value, err := call.Argument(i).ToString()
			assert.AsserErr(err)
			isIndexed, err := call.Argument(i + 1).ToBoolean()
			assert.AsserErr(err)
			if isIndexed {
				topics = append(topics, value)
			} else {
				data = append(data, []byte(value))
			}
		}
		if len(topics) > 4 {
			assert.AsserErr(errors.New("the api of emitEvent must have no more than four parameters to indicate true"))
		}
		self.apis.Contract.EmitEvent(topics, data)
		return otto.Value{}
	}
	obj.Object().Set("emitEvent", emitEvent)
}

// bridge Call Api
func (self *bridge) bridgeCallApi() {
	self.rt.vm.Set("_context_api_call", struct{}{})
	obj, _ := self.rt.vm.Get("_context_api_call")

	call := func(call otto.FunctionCall) otto.Value {
		addressStr, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		address := types.HexToAddress(addressStr)
		nilAddress := types.Address{}
		assert.AssertEx(address != nilAddress, "not a valid address string")
		paraStr, err := call.Argument(1).ToString()
		assert.AsserErr(err)

		act := &transaction.Action{address, []byte(paraStr)}
		fmt.Println("bridgeCallApi", address.Hex(), paraStr)
		self.apis.Call.Call(act)
		return otto.Value{}
	}
	obj.Object().Set("call", call)

	innerCall := func(call otto.FunctionCall) otto.Value {
		addressStr, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		address := types.HexToAddress(addressStr)
		nilAddress := types.Address{}
		assert.AssertEx(address != nilAddress, "not a valid address string")
		paraStr, err := call.Argument(1).ToString()
		assert.AsserErr(err)

		act := &transaction.Action{address, []byte(paraStr)}
		fmt.Println("bridgeCallApi", address.Hex(), paraStr)
		self.apis.Call.InnerCall(act)
		return otto.Value{}
	}
	obj.Object().Set("innerCall", innerCall)

	run := func(call otto.FunctionCall) otto.Value {
		addressStr, err := call.Argument(0).ToString()
		assert.AsserErr(err)
		address := types.HexToAddress(addressStr)
		nilAddress := types.Address{}
		assert.AssertEx(address != nilAddress, "not a valid address string")

		act := &transaction.Action{address, []byte("")}
		fmt.Println("bridgeRunApi", address.Hex())
		self.apis.Call.Call(act)
		return otto.Value{}
	}
	obj.Object().Set("run", run)
}
