package wasmre

import (
	"bchain.io/core/actioncontext"
	"bchain.io/common/types"
)

func (p *bchain) sender() types.Address {
	actApi := actioncontext.ActionApi{}
	actApi.SetCtx(p.ctx)
	from := actApi.Sender()
	return from
}

func (p *bchain) producer() types.Address {
	actApi := actioncontext.SystemApi{}
	actApi.SetCtx(p.ctx)
	producer := actApi.GetBlockMiner()
	return producer
}

func (p *bchain) conAddr() types.Address {
	actApi := actioncontext.ActionApi{}
	actApi.SetCtx(p.ctx)
	producer := actApi.Contract()
	return producer
}

func (p *bchain) dbGet(key  []byte) []byte {
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	val := actApi.Get(key)
	if val == nil {
		return make([]byte, 8)
	}

	return val
}
func (p *bchain) setResult(val []byte) {
	actApi := actioncontext.ResultApi{}
	actApi.SetCtx(p.ctx)
	actApi.SetActionResult(val)
}

func (p *bchain) dbSet(key []byte, val []byte) {
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	actApi.Set(key, val)
}
