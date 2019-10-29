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
// @File: context.go
// @Date: 2018/07/30 16:16:30
////////////////////////////////////////////////////////////////////////////////

package actioncontext

import (
	"math/big"
	"bchain.io/common/types"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
	"bchain.io/utils/database"
	"time"
)

type BlockContext struct {
	state  *state.StateDB
	db     database.IDatabase
	tmpDb  database.IDatabase
	number *big.Int
	miner 	types.Address
}

func (ctx *BlockContext) GetState() *state.StateDB {
	return ctx.state
}

func NewBlockContext(state *state.StateDB, db, tmpDb database.IDatabase, number *big.Int , minerAddress types.Address) *BlockContext {
	blkCtx := &BlockContext{
		state:  state,
		db:     db,
		tmpDb:  tmpDb,
		number: number,
		miner:minerAddress,
	}
	return blkCtx
}

type Context struct {
	action *transaction.Action
	sender types.Address
	con    contract

	blkCtx *BlockContext

	actionResult [][]byte //action result

	ttl          time.Duration
	resultCh     chan error
	contracts    []types.Address //contracts that action created
}

type Visitor interface {
	InterpreterId() types.Hash
	InterpreterName() string
	Code() []byte
	Param() []byte
	TTL() time.Duration
}

type Interpreter interface {
	Exec(ctx *Context) error
	ExecAsync(ctx *Context) error
}

type contract struct {
	iid     types.Hash    // interpreter id
	iname   string        // interpreter name
	creator types.Address // creator address
	self    types.Address // contract address
	code    []byte        // code
}

func NewContext(sender types.Address, action *transaction.Action, blkCtx *BlockContext) *Context {
	//if action == nil || blkCtx == nil || blkCtx.db == nil {
	//	return nil
	//}

	ctx := Context{
		action: action,
		sender: sender, // TODO: get by input transaction
		blkCtx: blkCtx,
		con: contract{
			iid:     blkCtx.state.GetInterpreterID(action.Contract), // TODO: get by action.Contract
			creator: blkCtx.state.GetCreator(action.Contract),       // TODO: get by action.Contract
			self:    action.Contract,
			code:    blkCtx.state.GetCode(action.Contract), // TODO: get by action.Contract
		},
		resultCh: make(chan error),
		ttl:      1 * time.Second, //TODO:
	}

	//nilHash := types.Hash{}
	//if ctx.con.iid == nilHash {
	//	return nil
	//}
	return &ctx
}

func (ctx *Context) Exec(in Interpreter) error {
	// TODO:
	return in.Exec(ctx)
}

func (ctx *Context) ExecAsync(in Interpreter) error {
	// TODO:
	return in.ExecAsync(ctx)
}

func (ctx *Context) InterpreterId() types.Hash {
	return ctx.con.iid
}

func (ctx *Context) InterpreterName() string {
	return ctx.con.iname
}

func (ctx *Context) Code() []byte {
	return ctx.con.code
}

func (ctx *Context) ContractAddress() types.Address {
	return ctx.con.self
}

func (ctx *Context) ActionResult() [][]byte {
	return ctx.actionResult
}

func (ctx *Context) Param() []byte {
	return ctx.action.Params
}

func (ctx *Context) TTL() time.Duration {
	return ctx.ttl
}

func (self *Context) SetActionResult(data []byte) {
	ret := make([]byte, len(data))
	copy(ret, data)
	self.actionResult = append(self.actionResult, ret)
}

func (ctx *Context) ResultCh() chan error {
	return ctx.resultCh
}

func (ctx *Context) AppendContract(addr types.Address) {
	ctx.contracts = append(ctx.contracts, addr)
}

func (ctx *Context) GetContracts() []types.Address {
	return ctx.contracts
}

func (ctx *Context) InitForTest() {
	code := `
//var apijs = require('context_api') // not support

var ctxapi = ctxApi();
function PrintHello(){
	ctxapi.console.print("================================PrintHello======================================\n");
}
;(function () {


	function PrintPrivateInfo(){
		ctxapi.console.print("===========Print Private Info===========\n")
	}

	function getRewordByNumber(number){
		ctxapi.console.print("===========getRewordByNumber===========\n")
		idx = number.dividedToIntegerBy(blockRewordAdjustNumber)
		val = (initReword.times(q.pow(idx))).toFixed(0)
		return new BigNumber(val)
	}

	//var decimals = 18
	//var totalSupply = new BigNumber(1e+9)
	blockRewordAdjustNumber = new BigNumber(6250000)
	initReword = new BigNumber(80e+18)
	//block reword adjust fraction
	q = new BigNumber(0.5)


	this.test = function(param) {
		bckNumber = new BigNumber('6250001')
		val = getRewordByNumber(bckNumber)
		ctxapi.console.print(val)
		ctxapi.console.print("\n")
		ctxapi.console.print("============================================================================\n");
		ctxapi.console.print("ctxapi._contract: " + ctxapi._contract + "\n");
		ctxapi.console.print("ctxapi._sender:   " + ctxapi._sender + "\n");
		ctxapi.console.print("ctxapi._blockNumber:   " + ctxapi._number + "\n");
		ctxapi.console.print("============================================================================\n");
	},
	this.crypto = function(param) {
		ctxapi.console.print("raw param: " + param)
		ctxapi.console.print("ctxapi.cypto.sha1: " + ctxapi.crypto.sha1(param) + "\n")
		ctxapi.console.print("ctxapi.cypto.sha256: " + ctxapi.crypto.sha256(param) + "\n")
		ctxapi.console.print("ctxapi.cypto.sha512: " + ctxapi.crypto.sha512(param) + "\n")
		var msg="0x41b1a0649752af1b28b3dc29a1556eee781e4a4c3a1f7f53f90fa834de098c4d";
		var sig="0xd155e94305af7e07dd8c32873e5c03cb95c9e05960ef85be9c07f671da58c73718c19adc397a211aa9e87e519e2038c5a3b658618db335f74f800b8e0cfeef4401";
		ctxapi.console.print("msg data: " + msg + "\n")
		ctxapi.console.print("sig data: " + sig + "\n")
		ctxapi.console.print("ctxapi.cypto.recover pubkey:0x" + ctxapi.crypto.recover(msg, sig) + "\n")
	},
	this.auth = function(param) {
		ctxapi.console.print("======================================================\n")
		ctxapi.console.print("Auth raw param:" + param + "\n")
		var addr="970e8128ab834e8eac17ab8e3812f010678c970e";
		ctxapi.console.print("ctxapi.auth.isHexAddress: " + ctxapi.auth.isHexAddress(addr) + "\n")
		ctxapi.console.print("ctxapi.auth.requireAuth: " + ctxapi.auth.requireAuth(addr) + "\n")
		ctxapi.console.print("ctxapi.auth.isAccount: " + ctxapi.auth.isAccount(addr) + "\n")
		ctxapi.console.print("ctxapi.auth.isContract: " + ctxapi.auth.isContract(addr) + "\n")
	},
	this.contractCreate = function(interpreter, creator, code) {
		ctxapi.console.print("======================================================\n")
		ctxapi.console.print("Contract Create Test: \n interpreter:" + interpreter + "\n creator:" + creator +"\n code:" + code +"\n")
		ctxapi.console.print("ctxapi.contract.create:" + ctxapi.contract.create(interpreter, creator, code) + "\n")
	},
	this.getRewardByBlockNum = function(blockNum) {
		ctxapi.console.print("========================getRewardByBlockNum==============================\n")
		totalSupply  = new BigNumber(1e+9)
		precision  = new BigNumber(1e+8)
		q = 0.5 //ratio
		a1BlockNumPerStep = new BigNumber(6250000)

		SnAllReward = totalSupply.times(precision)
		initReward = SnAllReward.times(1 - q)
		initReward = initReward.div(a1BlockNumPerStep)
		blk = new BigNumber(blockNum)
		idx = blk.dividedToIntegerBy(a1BlockNumPerStep)
		currentReward = initReward.dividedToIntegerBy(new BigNumber(2).pow(idx))

		ctxapi.console.print("current block num is " + blockNum + " , reward is:" + currentReward + " TOKEN\n")
		return currentReward
	}
})(this);`
	param := `
PrintHello();
test("hi, this is a test");
crypto("hello world\n");
auth("hello auth api");
//contractCreate("jsre.JSRE", "0x2a3cb462b299491f960891f8cb88675cdf5705ba", "this is code.");
getRewardByBlockNum(44150400);
`

	ctx.action = transaction.NewAction()
	ctx.action.Contract = types.HexToAddress("0x0011223344556677889900112233445566778899")
	ctx.action.Params = []byte(param)
	ctx.sender = types.HexToAddress("0x0000000000000000000000000000000000000011")
	ctx.con.iid = types.Hash{}
	ctx.con.iid = types.HexToHash("0xcfde91645b17a17d158271ae73a0c18b803d41d2072cdef747e8bbf135e6fe6b")
	ctx.con.iname = "jsre.JSRE"
	ctx.con.creator = types.HexToAddress("0x0000000000000000000000000000000000000033")
	ctx.con.self = types.HexToAddress("0x0011223344556677889900112233445566778899")
	ctx.con.code = []byte(code)
	ctx.ttl = 1 * time.Second
	ctx.blkCtx = new(BlockContext)
	ctx.blkCtx.number = big.NewInt(100)
}

func (ctx *Context) SetCreatorForTest(addr types.Address) {
	ctx.con.creator = addr
}
