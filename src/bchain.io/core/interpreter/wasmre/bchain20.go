package wasmre

import (
	"bchain.io/common/types"
	"bchain.io/core/actioncontext"
	"encoding/json"
	"bchain.io/common/assert"
	"encoding/binary"
	"strings"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"bytes"
	"time"
)

const (
	DecaimasBase  = uint64(100000000)
	BcPrefix = "BC"
)

var nilAddr = types.Address{}



func  newBc20Runtime(ctx *actioncontext.Context) *RuntimeBc20 {
	rt := &RuntimeBc20{
		ctx: ctx,
	}
	return rt
}

type RuntimeBc20 struct {
	ctx *actioncontext.Context
}

func (rt *RuntimeBc20) Exec(code []byte, param []byte, ttl time.Duration) time.Duration {
	now := time.Now()

	p := bchain{rt.ctx}
	p.run(param)

	return time.Since(now)
}


type bchain struct {
	ctx *actioncontext.Context
}

func (p *bchain) run(para []byte) {
	np := &para_paser.WasmPara{}
	err := json.Unmarshal(para, np)
	assert.AsserErr(err)
	paraLen := len(np.Args)
	ErrLen := "args length err"

	switch np.FuncName {
	case "reword":
		assert.AssertEx(paraLen == 0, ErrLen)
		p.rewords()
	case "transfer":
		assert.AssertEx(paraLen == 3, ErrLen)
		index := bytes.IndexByte(np.Args[0].Data, 0)
		to := np.Args[0].Data
		if index >= 0 {
			to = np.Args[0].Data[0:index]
		}
		index = bytes.IndexByte(np.Args[2].Data, 0)
		memo := np.Args[2].Data
		if index >= 0 {
			memo = np.Args[2].Data[0:index]
		}
		amount := make([]byte, 8)
		copy(amount, np.Args[1].Data)
		amountU := binary.LittleEndian.Uint64(amount)
		p.transfer(string(to), amountU, string(memo))
	case "transferFee":
		assert.AssertEx(paraLen == 1, ErrLen)
		amount := make([]byte, 8)
		copy(amount, np.Args[0].Data)
		amountU := binary.LittleEndian.Uint64(amount)
		p.transferFee(amountU)
	case "balanceOf", "balenceOf":
		assert.AssertEx(paraLen == 1, ErrLen)
		index := bytes.IndexByte(np.Args[0].Data, 0)
		addr := np.Args[0].Data
		if index >= 0 {
			addr = np.Args[0].Data[0:index]
		}
		p.balanceOf(string(addr))
	default:
		assert.AssertEx(false, "func name is invalid")
	}
}

func (p *bchain) transfer(to string, amount uint64, memo string) {
	assert.AssertEx(types.IsHexAddress(to), "to addr invalid")
	assert.AssertEx(len(memo) < 64, "memo length invalid")
	assert.AssertEx(amount > 0, "amount is invalid")

	from := p.sender()
	fromKey := append([]byte(BcPrefix), []byte(from.HexLower())...)
	fromVal := p.dbGet(fromKey)
	fromVal_uint := binary.LittleEndian.Uint64(fromVal)
	assert.AssertEx(fromVal_uint >= amount, "insufficient balance")
	fromVal_uint -= amount
	newFrom := make([]byte, 8)
	binary.LittleEndian.PutUint64(newFrom, fromVal_uint)
	p.dbSet(fromKey, newFrom)

	toKey := append([]byte(BcPrefix), []byte(strings.ToLower(to))...)
	toVal := p.dbGet(toKey)
	toVal_uint := binary.LittleEndian.Uint64(toVal)
	toVal_uint += amount
	newTo := make([]byte, 8)
	binary.LittleEndian.PutUint64(newTo, toVal_uint)
	p.dbSet(toKey, newTo)
}

func (p *bchain) transferFee(amount uint64) {
	assert.AssertEx(amount > 0, "amount is invalid")

	producer := p.producer()
	from := p.sender()
	if producer == from {
		return
	}

	fromKey := append([]byte(BcPrefix), []byte(from.HexLower())...)
	fromVal := p.dbGet(fromKey)
	fromVal_uint := binary.LittleEndian.Uint64(fromVal)

	toKey := append([]byte(BcPrefix), []byte(producer.HexLower())...)
	toVal := p.dbGet(toKey)
	toVal_uint := binary.LittleEndian.Uint64(toVal)

	assert.AssertEx(fromVal_uint >= amount, "insufficient balance")
	fromVal_uint -= amount
	toVal_uint += amount
	newFrom := make([]byte, 8)
	newTo := make([]byte, 8)
	binary.LittleEndian.PutUint64(newFrom, fromVal_uint)
	binary.LittleEndian.PutUint64(newTo, toVal_uint)
	p.dbSet(fromKey, newFrom)
	p.dbSet(toKey, newTo)
}

func (p *bchain) requireRewordAuth() {
	memDbApi := actioncontext.BlockMemDbApi{}
	memDbApi.SetCtx(p.ctx)
	memDbApi.Emplace([]byte("producer"), []byte("produced"))
}
func (p *bchain) getCurRewordsNumber() uint64 {
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	val := actApi.Get([]byte("_rNumber"))
	number := uint64(0)
	if val != nil {
		number = binary.LittleEndian.Uint64(val)
	}
	number++
	valInc := make([]byte, 8)
	binary.LittleEndian.PutUint64(valInc, number)
	actApi.Set([]byte("_rNumber"), valInc)
	return number
}
func (p *bchain) getRewordsValue(number uint64) uint64 {
	if number == 1 {
		return 17500000 * DecaimasBase
	} else if number < 250000 {
		return 2*DecaimasBase
	} else if number < 365000 {
		return 199600000
	} else {
		return 4000000
	}
}


func (p *bchain) rewords() {
	p.requireRewordAuth()
	producer := p.producer()
	number := p.getCurRewordsNumber()
	rewordsVal := p.getRewordsValue(number)

	toKey := append([]byte(BcPrefix), []byte(producer.HexLower())...)
	toVal := p.dbGet(toKey)
	toVal_uint := binary.LittleEndian.Uint64(toVal)

	toVal_uint += rewordsVal
	newTo := make([]byte, 8)
	binary.LittleEndian.PutUint64(newTo, toVal_uint)
	p.dbSet(toKey, newTo)
}


func (p *bchain) balanceOf(addr string) {
	assert.AssertEx(types.IsHexAddress(addr), "to addr invalid")
	Addr := types.HexToAddress(addr)
	addrKey := append([]byte(BcPrefix), []byte(Addr.HexLower())...)
	val := p.dbGet(addrKey)
	p.setResult(val)
}



