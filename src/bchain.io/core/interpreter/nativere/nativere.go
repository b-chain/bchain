package nativere

import (
	"bchain.io/core/actioncontext"
	"bchain.io/core/interpreter"
	"time"
)

func NewNativeRe() interpreter.PluginImpl {
	return &NativeRe{}
}

type NativeRe struct {
}

func (self *NativeRe) Initialize() {
}

func (self *NativeRe) Startup() {

}

func (self *NativeRe) Shutdown() {
}

func (self *NativeRe) Generate(ctx *actioncontext.Context) interpreter.Interpreter {
	rt := self.newRuntime(ctx)
	return rt
}

func (self *NativeRe) newRuntime(ctx *actioncontext.Context) *Runtime {
	rt := &Runtime{
		ctx: ctx,
	}
	return rt
}

type Runtime struct {
	ctx *actioncontext.Context
}

func (rt *Runtime) Exec(code []byte, param []byte, ttl time.Duration) time.Duration {
	now := time.Now()
	codeNative := code[0]
	switch codeNative {
	case 0:
		p := pledge{rt.ctx}
		p.run(param)
	default:
		panic("code is invalid")
	}
	return time.Since(now)
}
