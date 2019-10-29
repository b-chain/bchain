package nativere

import (
	"bchain.io/common/types"
	"bchain.io/core/actioncontext"
	"bytes"
	"github.com/tinylib/msgp/msgp"
	"fmt"
	"bchain.io/common/assert"
	"encoding/binary"
	"math"
)

func (p *pledge) sender() types.Address {
	actApi := actioncontext.ActionApi{}
	actApi.SetCtx(p.ctx)
	from := actApi.Sender()
	return from
}

func (p *pledge) producer() types.Address {
	actApi := actioncontext.SystemApi{}
	actApi.SetCtx(p.ctx)
	producer := actApi.GetBlockMiner()
	return producer
}

func (p *pledge) conAddr() types.Address {
	actApi := actioncontext.ActionApi{}
	actApi.SetCtx(p.ctx)
	producer := actApi.Contract()
	return producer
}

func (p *pledge) dbGet(addr types.Address) []byte {
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	val := actApi.Get(addr.Bytes())
	if val == nil {
		return make([]byte, 8)
	}

	return val
}
func (p *pledge) setResult(val string) {
	actApi := actioncontext.ResultApi{}
	actApi.SetCtx(p.ctx)
	actApi.SetActionResult([]byte(val))
}

func (p *pledge) dbSet(addr types.Address, val []byte) {
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	actApi.Set(addr.Bytes(), val)
}



var (
	pdPrefix = []byte("_pd")
	poolPrefix = []byte("_pl")
	producerPrefix = []byte("_pr")
	pdTotalKey = []byte("_PdTotal")
	pdTotalWeightKey = []byte("_PdTW")
)

func (p *pledge) getPledgeData(addr types.Address) *PledgeData {
	nilAddr := types.Address{}
	if addr == nilAddr {
		return nil
	}
	key := append(pdPrefix, addr.Bytes()...)
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	val := actApi.Get(key)
	if val == nil {
		return nil
	}
	pd := new(PledgeData)
	byteBuf := bytes.NewBuffer(val)
	err := msgp.Decode(byteBuf, pd)
	if err != nil {
		panic(fmt.Errorf("getPledgeData fail! Err:%s", err.Error()))
	}
	return pd
}

func (p *pledge) savePledgeData(addr types.Address, pdData *PledgeData)  {
	var encData bytes.Buffer
	err := msgp.Encode(&encData, pdData)
	if err != nil {
		assert.AsserErr(err)
	}
	key := append(pdPrefix, addr.Bytes()...)
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	actApi.Set(key, encData.Bytes())
}
func (p *pledge) getPoolData(addr types.Address) *PoolData {
	nilAddr := types.Address{}
	if addr == nilAddr {
		return nil
	}
	key := append(poolPrefix, addr.Bytes()...)
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	val := actApi.Get(key)
	if val == nil {
		return nil
	}
	pool := new(PoolData)
	byteBuf := bytes.NewBuffer(val)
	err := msgp.Decode(byteBuf, pool)
	if err != nil {
		panic(fmt.Errorf("getPoolData fail! Err:%s", err.Error()))
	}
	return pool
}

func (p *pledge) savePoolData(addr types.Address, pdData *PoolData)  {
	var encData bytes.Buffer
	err := msgp.Encode(&encData, pdData)
	if err != nil {
		assert.AsserErr(err)
	}
	key := append(poolPrefix, addr.Bytes()...)
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	actApi.Set(key, encData.Bytes())
}

func (p *pledge) getProducerData(addr types.Address) *ProducerData {
	nilAddr := types.Address{}
	if addr == nilAddr {
		return nil
	}
	key := append(producerPrefix, addr.Bytes()...)
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	val := actApi.Get(key)
	if val == nil {
		return nil
	}
	pd := new(ProducerData)
	byteBuf := bytes.NewBuffer(val)
	err := msgp.Decode(byteBuf, pd)
	if err != nil {
		panic(fmt.Errorf("getProducerData fail! Err:%s", err.Error()))
	}
	return pd
}

func (p *pledge) saveProducerData(addr types.Address, pdData *ProducerData)  {
	var encData bytes.Buffer
	err := msgp.Encode(&encData, pdData)
	if err != nil {
		assert.AsserErr(err)
	}
	key := append(producerPrefix, addr.Bytes()...)
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	actApi.Set(key, encData.Bytes())
}

func (p *pledge) getWeight(amount uint64) uint64 {
	if amount < 500000000000 {
		return 0
	}
	return uint64(math.Pow(float64(amount),0.33))
}
func (p *pledge) addTotalPledgeVal(amount uint64)  {
	key := pdTotalKey
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	val := actApi.Get(key)
	totalAmount := amount
	if val != nil {
		lastAmount := binary.LittleEndian.Uint64(val)
		totalAmount += lastAmount
	}
	newVal := make([]byte, 8)
	binary.LittleEndian.PutUint64(newVal, totalAmount)
	actApi.Set(key, newVal)
}

func (p *pledge) addTotalPledgeWeight(amount uint64)  {
	key := pdTotalWeightKey
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	val := actApi.Get(key)
	totalAmount := amount
	if val != nil {
		lastAmount := binary.LittleEndian.Uint64(val)
		totalAmount += lastAmount
	}
	newVal := make([]byte, 8)
	binary.LittleEndian.PutUint64(newVal, totalAmount)
	actApi.Set(key, newVal)
}

func (p *pledge) subTotalPledgeVal(amount uint64)  {
	key := pdTotalKey
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	val := actApi.Get(key)
	assert.AssertEx(val != nil, "subTotalPledgeVal get val err")
	lastAmount := binary.LittleEndian.Uint64(val)
	totalAmount := lastAmount - amount
	newVal := make([]byte, 8)
	binary.LittleEndian.PutUint64(newVal, totalAmount)
	actApi.Set(key, newVal)
}

func (p *pledge) subTotalPledgeWeight(amount uint64)  {
	key := pdTotalWeightKey
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	val := actApi.Get(key)
	assert.AssertEx(val != nil, "subTotalPledgeWeight get val err")
	lastAmount := binary.LittleEndian.Uint64(val)
	totalAmount := lastAmount - amount
	assert.AssertEx(totalAmount > 0, "subTotalPledgeWeight to 0")
	newVal := make([]byte, 8)
	binary.LittleEndian.PutUint64(newVal, totalAmount)
	actApi.Set(key, newVal)
}

func (p *pledge) getTotalAmount() uint64 {
	key := pdTotalKey
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	val := actApi.Get(key)
	totalAmount := uint64(0)
	if val != nil {
		lastAmount := binary.LittleEndian.Uint64(val)
		totalAmount += lastAmount
	}
	return totalAmount
}

func (p *pledge) getTotalWeight() uint64 {
	key := pdTotalWeightKey
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	val := actApi.Get(key)
	totalAmount := uint64(0)
	if val != nil {
		lastAmount := binary.LittleEndian.Uint64(val)
		totalAmount += lastAmount
	}
	return totalAmount
}