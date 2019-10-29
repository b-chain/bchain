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
// @File: plugin.go
// @Date: 2018/07/30 17:34:30
////////////////////////////////////////////////////////////////////////////////

package interpreter

import (
	"errors"
	"fmt"
	"bchain.io/common/assert"
	"bchain.io/common/types"
	"bchain.io/core/actioncontext"
	"bchain.io/utils/crypto"
	"reflect"
	"sync"
)

type State int

const (
	Registered  = iota ///< the plugin is constructed but doesn't do anything
	Initialized        ///< the plugin has initialized any state required but is idle
	Started            ///< the plugin is actively running
	Stopped            ///< the plugin is no longer running
)

// a interpreter must implement PluginImpl interface
type PluginImpl interface {
	Initialize()
	Startup()
	Shutdown()
	Generator
}

type Generator interface {
	Generate(ctx *actioncontext.Context) Interpreter
}

type Plugin interface {
	GetState() State
	Name() string
	Id() types.Hash

	PluginImpl
	//Register( /*p *Plugin*/ )
}

type PluginObj struct {
	pImpl PluginImpl
	state State
	name  PluginName
}

func newPluginObj(pImpl PluginImpl) *PluginObj {
	plugObj := &PluginObj{
		pImpl,
		Registered,
		PluginName{},
	}
	plugObj.name.Set(pImpl)
	return plugObj
}

func (obj *PluginObj) Initialize() {
	assert.Assert(obj.pImpl != nil)

	if obj.state == Registered {
		obj.state = Initialized
		obj.pImpl.Initialize()

		Singleton().pluginInitialized(obj)
	}
	assert.Assert(obj.state == Initialized)
}

func (obj *PluginObj) Startup() {
	assert.Assert(obj.pImpl != nil)

	if obj.state == Initialized {
		obj.state = Started
		obj.pImpl.Startup()

		Singleton().pluginStarted(obj)
	}
	assert.Assert(obj.state == Started)
}

func (obj *PluginObj) Shutdown() {
	//assert.Assert(obj.state == Started)

	if obj.state == Started {
		obj.state = Stopped
		obj.pImpl.Shutdown()
	}
}

func (obj PluginObj) Generate(ctx *actioncontext.Context) Interpreter {
	return obj.pImpl.Generate(ctx)
}

func (obj PluginObj) GetState() State {
	return obj.state
}

func (obj PluginObj) Name() string {
	return obj.name.Name()
}

func (obj PluginObj) Id() types.Hash {
	return obj.name.Id()
}

type PluginName struct {
	pImpl PluginImpl
	name  string
	id    types.Hash
}

func (pn *PluginName) Set(pImpl interface{}) {
	ppImpl := reflect.TypeOf(pImpl) //.Elem()
	required := reflect.TypeOf((*PluginImpl)(nil)).Elem()
	assert.Assert(ppImpl.Implements(required))
	pn.pImpl, _ = ppImpl.(PluginImpl)
	pn.name = ppImpl.Elem().String()
	pn.id = crypto.Keccak256Hash([]byte(pn.name))
	// TODO: The name only needs the string after the last ".". i.e., main.pkgname => pkgname

	fmt.Printf("%s, %s\n", required.String(), pn.name)
}

func (pn PluginName) Name() string {
	return pn.name
}

func (pn PluginName) Id() types.Hash {
	return pn.id
}

// Interpreters
// Interpreters is a singleton to manage interpreters
var (
	instance *Interpreters
	once     sync.Once
)

// get the interpreter singleton
func Singleton() *Interpreters {
	once.Do(func() {
		instance = &Interpreters{
			plugins:            make(map[types.Hash]Plugin),
			initializedPlugins: make([]Plugin, 0),
			runningPlugins:     make([]Plugin, 0),
			ctxCh:              make(chan *actioncontext.Context),
			wg:                 &sync.WaitGroup{},
			exitCh:             make(chan struct{}),
		}
	})
	return instance
}

type Interpreters struct {
	pgLock             sync.RWMutex
	plugins            map[types.Hash]Plugin ///< all registered plugins
	initializedPlugins []Plugin
	runningPlugins     []Plugin

	ctxCh  chan *actioncontext.Context
	wg     *sync.WaitGroup
	exitCh chan struct{}
}

// register plugin, need input the constructor of plugin
func (self *Interpreters) Register(plugin func() PluginImpl) Plugin {
	plug := newPluginObj(plugin())
	if self.FindById(plug.Id()) != nil {
		return nil
	}

	self.pgLock.Lock()
	self.plugins[plug.Id()] = plug
	self.pgLock.Unlock()
	logger.Infof("Register plugin : %x, %s\n", plug.Id(), plug.Name())
	return plug
}

func (self Interpreters) FindById(id types.Hash) Plugin {
	self.pgLock.Lock()
	defer self.pgLock.Unlock()
	if plugin, ok := self.plugins[id]; ok {
		return plugin
	}
	return nil
}

func (self Interpreters) Find(plugin func() PluginImpl) Plugin {
	name := PluginName{}
	name.Set(plugin())
	return self.FindById(name.Id())
}

func (self *Interpreters) Initialize() {
	self.pgLock.Lock()
	defer self.pgLock.Unlock()
	for _, plug := range self.plugins {
		plug.Initialize()
	}
}

func (self *Interpreters) Startup() {
	defer func() {
		// TODO: assert process
		if err := recover(); err != nil {
			logger.Errorf("%v\n", err)
			self.Shutdown()
			assert.Assert(false)
		}
	}()

	// startup all plugin
	for _, plug := range self.initializedPlugins {
		plug.Startup()
		self.wg.Add(1)
	}

	// process routine
	restartCh := make(chan bool)
	go func() {
		go self.process(restartCh)
		for {
			select {
			case restart := <-restartCh:
				if restart {
					go self.process(restartCh)
				} else {
					logger.Debugf("process routine return\n")
					return
				}
			}
		}
	}()

	// wg wait
	go func() {
		self.wg.Wait()
		self.exitCh <- struct{}{}
	}()
}

// process action context
// synchronous function, all action is synchronous
func (self *Interpreters) process(restartCh chan bool) {
	var ctx *actioncontext.Context = nil
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("%v\n", err)

			// do not panic and shutdown, continue to run
			//panic(err)
			//self.Shutdown()
			if ctx != nil {
				ctx.ResultCh() <- errors.New(fmt.Sprintf("%v", err))
			}
			restartCh <- true
		}
	}()

	for {
		select {
		case ctx = <-self.ctxCh:
			if ctx == nil {
				logger.Warnf("get a null context \n")
			} else {
				plug := self.FindById(ctx.InterpreterId())
				if plug != nil {
					in := plug.Generate(ctx)
					in.Exec(ctx.Code(), ctx.Param(), ctx.TTL())
					ctx.ResultCh() <- nil
				} else {
					logger.Warn("can't get plugin by id", ctx.InterpreterId().String(), "contract Addr", ctx.ContractAddress().HexLower())
					ctx.ResultCh() <- errors.New("interpreter is not registered")
				}
			}
		case <-self.exitCh:
			restartCh <- false
			return
		}
	}
}

func (self *Interpreters) Shutdown() {
	self.pgLock.Lock()
	defer self.pgLock.Unlock()
	for _, plug := range self.runningPlugins {
		plug.Shutdown()
		self.wg.Done()
	}
	self.runningPlugins = nil
	self.initializedPlugins = nil
	self.plugins = nil
}

func (self *Interpreters) Exec(ctx *actioncontext.Context) error {
	if ctx == nil {
		logger.Warnf("action context is null\n")
		return errors.New("action context is null")
	}
	self.ctxCh <- ctx
	err := <-ctx.ResultCh()
	return err
}

func (self *Interpreters) ExecAsync(ctx *actioncontext.Context) error {
	if ctx == nil {
		logger.Warnf("action context is null\n")
		return errors.New("action context is null")
	}
	errCh := make(chan error, 1)
	self.execAsync(ctx, errCh)
	err := <-errCh
	return err
}

func (self *Interpreters) execAsync(ctx *actioncontext.Context, errorCh chan error) {
	defer func() {
		if err := recover(); err != nil {
			//debug.PrintStack()
			//fmt.Println("execAsync", err)
			errorCh <- errors.New(fmt.Sprintf("%v", err))
		}
	}()
	plug := self.FindById(ctx.InterpreterId())
	if plug != nil {
		in := plug.Generate(ctx)
		in.Exec(ctx.Code(), ctx.Param(), ctx.TTL())
		errorCh <- nil
		return
	} else {
		logger.Warnf("can't get plugin by id: %x", ctx.InterpreterId())
		errorCh <- errors.New("can't get plugin")
		return
	}
}

func (self *Interpreters) pluginInitialized(plug Plugin) {
	self.initializedPlugins = append(self.initializedPlugins, plug)
}

func (self *Interpreters) pluginStarted(plug Plugin) {
	self.runningPlugins = append(self.runningPlugins, plug)
}
