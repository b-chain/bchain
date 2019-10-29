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
// @File: api.go
// @Date: 2018/07/30 16:14:30
////////////////////////////////////////////////////////////////////////////////

package actioncontext

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"math/big"
	"bchain.io/common/assert"
	"bchain.io/common/types"
	"bchain.io/core/interpreter/contract_parser"
	"bchain.io/core/transaction"
	"bchain.io/params"
	"bchain.io/utils/crypto"
	"bchain.io/utils/database"
	"runtime"
)

type APIs struct {
	Auth     AuthorizationApi
	Sys      SystemApi
	Assert   AssertApi
	Console  ConsoleApi
	Crypto   CryptoApi
	Producer ProducerApi
	Act      ActionApi
	Db       DatabaseApi
	MemDb    BlockMemDbApi
	Contract ContractApi
	Result   ResultApi
	Call     CallApi
	Types    TypesApi
}

func (self *APIs) SetInterpreter(interpreter Interpreter) {
	self.Call.setInterpreter(interpreter)
}

func NewAPIs(ctx *Context) *APIs {
	apis := APIs{}
	apis.Auth.setCtx(ctx)
	apis.Sys.setCtx(ctx)
	apis.Assert.setCtx(ctx)
	apis.Console.setCtx(ctx)
	apis.Crypto.setCtx(ctx)
	apis.Producer.setCtx(ctx)
	apis.Act.setCtx(ctx)
	apis.Db.setCtx(ctx)
	apis.MemDb.setCtx(ctx)
	apis.Contract.setCtx(ctx)
	apis.Result.setCtx(ctx)
	apis.Call.setCtx(ctx)
	return &apis
}

type api interface {
	setCtx(ctx *Context)
	assert(ret bool, msg string)
}

type baseApi struct {
	ctx *Context
}

func (self *baseApi) setCtx(ctx *Context) {
	self.ctx = ctx
}

func (self *baseApi) SetCtx(ctx *Context) {
	self.ctx = ctx
}

func (self *baseApi) assert(test bool, msg string) {
	if test {
		return
	}

	pc, file, line, ok := runtime.Caller(1)
	_ = file
	_ = line
	if ok {
		f := runtime.FuncForPC(pc)
		if msg == "" {
			msg = "unknown"
		}
		assert.AssertEx(false, fmt.Sprintf("%s error : %s", f.Name(), msg))
	} else {
		assert.Assert(false)
	}
}

type AuthorizationApi struct {
	baseApi
}

func (self *AuthorizationApi) IsHexAddress(address string) bool {
	return types.IsHexAddress(address)
}

func (self *AuthorizationApi) RequireAuth(address string) bool {
	if !self.IsHexAddress(address) {
		return false
	}
	if self.ctx == nil {
		return false
	}
	if self.ctx.blkCtx.state == nil {
		return false
	}
	addr := types.HexToAddress(address)
	if self.ctx.con.creator == addr {
		return true
	} else {
		return false
	}
}

func (self *AuthorizationApi) IsAccount(address string) bool {
	if !self.IsHexAddress(address) {
		return false
	}
	if self.ctx == nil {
		return false
	}
	if self.ctx.blkCtx.state == nil {
		return false
	}
	addr := types.HexToAddress(address)
	if !self.ctx.blkCtx.state.Exist(addr) {
		return false
	}
	if self.ctx.blkCtx.state.Empty(addr) {
		return false
	}
	if self.ctx.blkCtx.state.GetCodeSize(addr) > 0 {
		return false
	} else {
		return true
	}
}

func (self *AuthorizationApi) IsContract(address string) bool {
	if !self.IsHexAddress(address) {
		return false
	}
	if self.ctx == nil {
		return false
	}
	if self.ctx.blkCtx.state == nil {
		return false
	}
	addr := types.HexToAddress(address)
	if !self.ctx.blkCtx.state.Exist(addr) {
		return false
	}
	if self.ctx.blkCtx.state.Empty(addr) {
		return false
	}
	if self.ctx.blkCtx.state.GetCodeSize(addr) > 0 {
		return true
	} else {
		return false
	}
}

/*func (self *AuthorizationApi) pushRecipient(address types.Address) {

}*/

type SystemApi struct {
	baseApi
}

func (self *SystemApi) BytesToInt(src []byte) uint64 {
	a := new(big.Int).SetBytes(src)
	return a.Uint64()
}

func (self *SystemApi) IntToBytes(i uint64) []byte {
	a := new(big.Int).SetUint64(i)
	return a.Bytes()
}

func (self *SystemApi) GetBlockNumber() *big.Int {
	return self.ctx.blkCtx.number
}

func (self *SystemApi) GetBlockMiner() types.Address {
	return self.ctx.blkCtx.miner
}

type AssertApi struct {
	baseApi
}

func (self *AssertApi) Assert(test bool, msg string) {
	assert.AssertEx(test, msg)
}

type ConsoleApi struct {
	baseApi
}

func (self *ConsoleApi) Printf(str string) {
	fmt.Printf(str)
}

type CryptoApi struct {
	baseApi
}

func (self *CryptoApi) Sha1(val []byte) [20]byte {
	return sha1.Sum(val)
}

func (self *CryptoApi) Sha256(val []byte) [32]byte {
	return sha256.Sum256(val)
}

func (self *CryptoApi) Sha512(val []byte) [64]byte {
	return sha512.Sum512(val)
}

//Recover recovers public key through signature and msg data.
//var testAddrHex = "970e8128ab834e8eac17ab8e3812f010678cf791"
//var testPrivHex = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
//key, _ := crypto.HexToECDSA(testPrivHex)
//msg := crypto.Keccak256([]byte("foo"))
//sig, err := crypto.Sign(msg, key)
func (self *CryptoApi) Recover(msg []byte, sig []byte) (pubkey []byte, err error) {
	recoveredPub, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return nil, err
	}
	addr := crypto.PubkeyToAddress(*recoveredPub)
	return addr[:], nil
}

type ProducerApi struct {
	baseApi
}

func (self *ProducerApi) Producer() types.Address {
	return types.Address{}
}

type ActionApi struct {
	baseApi
}

func (self *ActionApi) Contract() types.Address {
	return self.ctx.con.self
}

func (self *ActionApi) Sender() types.Address {
	return self.ctx.sender
}

type ResultApi struct {
	baseApi
}

func (self *ResultApi) SetActionResult(data []byte) {
	self.ctx.SetActionResult(data)
}

type DatabaseApi struct {
	baseApi
}

func (self *DatabaseApi) Emplace(key []byte, val []byte) {
	keyHash := types.BytesToHash(crypto.Keccak256(key))
	valHash := types.BytesToHash(crypto.Keccak256(val))

	ret := self.ctx.blkCtx.state.GetState(self.ctx.con.self, keyHash)
	self.assert(bytes.Compare(ret.Bytes(), types.Hash{}.Bytes()) == 0, "key is exist")

	self.ctx.blkCtx.state.SetState(self.ctx.con.self, keyHash, valHash)
	self.ctx.blkCtx.state.AddPreimage(valHash, val)
}

func (self *DatabaseApi) Modify(key []byte, val []byte) {
	keyHash := types.BytesToHash(crypto.Keccak256(key))
	valHash := types.BytesToHash(crypto.Keccak256(val))

	ret := self.ctx.blkCtx.state.GetState(self.ctx.con.self, keyHash)
	self.assert(bytes.Compare(ret.Bytes(), types.Hash{}.Bytes()) != 0, "key is not exist")

	self.ctx.blkCtx.state.SetState(self.ctx.con.self, keyHash, valHash)
	self.ctx.blkCtx.state.AddPreimage(valHash, val)
}

func (self *DatabaseApi) Set(key []byte, val []byte) {
	//fmt.Println("DatabaseApi", string(key), val)
	//if len(val) >=8 {
	//	v := binary.LittleEndian.Uint64(val)
	//	fmt.Printf("DatabaseApi %v %x %v\n", string(key), string(val), v)
	//}

	keyHash := types.BytesToHash(crypto.Keccak256(key))
	valHash := types.BytesToHash(crypto.Keccak256(val))

	self.ctx.blkCtx.state.SetState(self.ctx.con.self, keyHash, valHash)
	//fmt.Printf("AddPreimage: Hash:%x , Val:%x\n", valHash, val)
	self.ctx.blkCtx.state.AddPreimage(valHash, val)
}

func (self *DatabaseApi) Erase(key []byte) {
	// TODO: ??
	// nothing to do
	self.assert(false, NotSupport)
}

func (self *DatabaseApi) Find(key []byte) []byte {
	keyHash := types.BytesToHash(crypto.Keccak256(key))
	ret := self.ctx.blkCtx.state.GetState(self.ctx.con.self, keyHash)
	if bytes.Compare(ret.Bytes(), types.Hash{}.Bytes()) == 0 {
		//fmt.Println("Find return nil")
		return nil
	}
	//fmt.Printf("Find ValHash:%x\n", ret)
	// TODO: maybe use self.ctx.db.Get(preimagePrefix + ret);  Refactor statedb later, may be adjustments here.
	//step 1:get data from preimage
	preimage := self.ctx.blkCtx.state.GetPreimage(ret)
	if preimage != nil {
		return preimage
	} else {
		//fmt.Println("getPreimage == nil")
	}
	//step 2: get data from db with prefix
	table := database.NewTable(self.ctx.blkCtx.db, "secure-key-")
	value, err := table.Get(ret.Bytes())
	if err != nil {
		fmt.Println("db.Get return nil")
		return nil
	}
	return value
}

func (self *DatabaseApi) Get(key []byte) []byte {
	ret := self.Find(key)
	return ret
	//self.assert(ret != nil, "unable to find key")
	//return ret
}

// block memory database api
type BlockMemDbApi struct {
	baseApi
}

func (self *BlockMemDbApi) Emplace(key []byte, val []byte) {
	_, err := self.ctx.blkCtx.tmpDb.Get(key)
	self.assert(err != nil, "key is exist")

	err = self.ctx.blkCtx.tmpDb.Put(key, val)
	self.assert(err == nil, "put db error")
}

func (self *BlockMemDbApi) Modify(key []byte, val []byte) {
	_, err := self.ctx.blkCtx.tmpDb.Get(key)
	self.assert(err == nil, "key is not exist")

	err = self.ctx.blkCtx.tmpDb.Put(key, val)
	self.assert(err == nil, "put db error")
}

func (self *BlockMemDbApi) Set(key []byte, val []byte) {
	err := self.ctx.blkCtx.tmpDb.Put(key, val)
	self.assert(err == nil, "put db error")
}

func (self *BlockMemDbApi) Erase(key []byte) {
	err := self.ctx.blkCtx.tmpDb.Delete(key)
	self.assert(err == nil, "block mem db delete error")
}

func (self *BlockMemDbApi) Find(key []byte) []byte {
	value, err := self.ctx.blkCtx.tmpDb.Get(key)
	if err != nil {
		return nil
	}
	return value
}

func (self *BlockMemDbApi) Get(key []byte) []byte {
	ret := self.Find(key)
	return ret
}

type ContractApi struct {
	baseApi
}

func (self *ContractApi) Create(creator string, code string) (addr types.Address) {
	if len(creator) == 0 || len(code) == 0 {
		assert.AssertEx(false, "Parameter error.")
		return types.Address{}
	}
	if self.ctx == nil {
		assert.AssertEx(false, "context is nil.")
		return types.Address{}
	}
	if self.ctx.blkCtx.state == nil {
		assert.AssertEx(false, "context state is nil.")
		return types.Address{}
	}

	assert.AssertEx(len(code) <= params.MaxCodeSize, fmt.Sprintf("code is too long, max length is %d .", params.MaxCodeSize))

	// Ensure there's no existing contract already at the designated address
	nonce := self.ctx.blkCtx.state.GetNonce(self.ctx.con.self)
	contractAddr := crypto.CreateAddress(types.HexToAddress(creator), nonce)
	contractHash := self.ctx.blkCtx.state.GetCodeHash(contractAddr)
	if self.ctx.blkCtx.state.GetNonce(contractAddr) != 0 || (contractHash != (types.Hash{}) && contractHash != crypto.Keccak256Hash(nil)) {
		assert.AssertEx(false, "contract address is already exist.")
		return types.Address{}
	}

	cc := contract_parser.ParseCodeByData([]byte(code))

	self.ctx.blkCtx.state.SetNonce(self.ctx.con.self, nonce+1)
	// Create a new account on the state
	self.ctx.blkCtx.state.CreateAccount(contractAddr)
	self.ctx.blkCtx.state.SetNonce(contractAddr, 1)
	self.ctx.blkCtx.state.SetInterpreterID(contractAddr, crypto.Keccak256Hash([]byte(cc.InterName)))
	self.ctx.blkCtx.state.SetCreator(contractAddr, types.HexToAddress(creator))
	self.ctx.blkCtx.state.SetCode(contractAddr, []byte(cc.Code))
	self.ctx.AppendContract(contractAddr)

	//fmt.Println(contractAddr.HexLower(), cc.InterName)
	return contractAddr
}

func (self *ContractApi) Creator() types.Address {
	return self.ctx.con.creator
}

func (self *ContractApi) EmitEvent(topics []string, data [][]byte) {
	topic := []types.Hash{}
	for _, t := range topics {
		if len(t) > types.HashLength {
			hsh := crypto.Keccak256Hash([]byte(t))
			topic = append(topic, hsh)
		} else {
			topic = append(topic, types.StringToHash(t))
		}
	}
	self.ctx.blkCtx.GetState().AddLog(&transaction.Log{
		Address: self.ctx.con.self,
		Topics:  topic,
		Data:    data,
		// This is a non-consensus field, but assigned here because
		// core/state doesn't know the current block number.
		BlockNumber: self.ctx.blkCtx.number.Uint64(),
	})
}

type CallApi struct {
	baseApi
	interpreter Interpreter
}

func (self *CallApi) setInterpreter(interpreter Interpreter) {
	self.interpreter = interpreter
}

func (self *CallApi) SetInterpreter(interpreter Interpreter) {
	self.interpreter = interpreter
}

func (self *CallApi) Call(act *transaction.Action) {
	//fmt.Println("Call Sender:", self.ctx.sender.Hex())
	//fmt.Println("contract addr:", self.ctx.action.Contract.HexLower())
	ctx := NewContext(self.ctx.sender, act, self.ctx.blkCtx)
	self.assert(ctx != nil, "CallApi, new action context fail")

	self.assert(self.interpreter != nil, "CallApi, interpreter is nil")
	err := ctx.ExecAsync(self.interpreter)

	errMsg :=""
	if err != nil {
		errMsg = err.Error()
	}

	self.assert(err == nil, "CallApi, action exec fail "+errMsg)
}

func (self *CallApi) InnerCall(act *transaction.Action) {
	fmt.Println("InnerCall Sender:", self.ctx.con.self.Hex())
	ctx := NewContext(self.ctx.con.self, act, self.ctx.blkCtx)
	self.assert(ctx != nil, "CallApi, new action context fail")

	self.assert(self.interpreter != nil, "CallApi, interpreter is nil")
	err := ctx.ExecAsync(self.interpreter)
	errMsg :=""
	if err != nil {
		errMsg = err.Error()
	}
	self.assert(err == nil, "InnerCall Api, action exec fail "+ errMsg)
}

type TypesApi struct {
	baseApi
}

func (self *TypesApi) Name(str string) uint64 {
	self.assert(!IsNormativeName(str), "not a normative name("+str+")")
	return StringToName(str).value
}
