package migrate

import (
	lua "github.com/yuin/gopher-lua"
	"runtime"
	"sync"
)

type (
		
	luaVmPool struct {
		safe sync.RWMutex
		pool map[string]*lua.LState
	}

	luaFsPool struct {
		safe sync.RWMutex
		pool map[string]FileSystem
	}
)

var (
	vmPoolImpl = newLuaVmPool()
	fsPoolImpl = newLuaFsPool()
)

func newLuaVmPool() *luaVmPool {
	var p = new(luaVmPool)
	return p.init()
}

func newLuaFsPool() *luaFsPool {
	var p = new(luaFsPool)
	return p.init()
}

func (p *luaFsPool) init() *luaFsPool {
	p.safe = sync.RWMutex{}
	p.pool = make(map[string]FileSystem)
	runtime.SetFinalizer(p, (*luaFsPool).destroy)
	return p
}

func (p *luaFsPool) add(name string, system FileSystem) bool {
	if name == "" || system == nil {
		return false
	}
	p.safe.Lock()
	defer p.safe.Unlock()
	if _, ok := p.pool[name]; ok {
		return false
	}
	p.pool[name] = system
	return true
}

func (p *luaFsPool) Get(name string) (FileSystem, bool) {
	if name == "" || p.pool == nil {
		return nil, false
	}
	p.safe.Lock()
	defer p.safe.Unlock()
	if system, ok := p.pool[name]; ok && system != nil {
		return system, true
	}
	return nil, false
}

// 默认lua 运行时
func createMigrateVm() *lua.LState {
	var vm = lua.NewState()
	vm.Push(vm.NewFunction(NewLuaMigrateTables()))
	vm.Push(lua.LString(Name))
	vm.Call(1, 0)
	return vm
}

func (p *luaVmPool) init() *luaVmPool {
	p.safe = sync.RWMutex{}
	p.pool = make(map[string]*lua.LState)
	runtime.SetFinalizer(p, (*luaVmPool).destroy)
	return p
}

func (p *luaVmPool) add(name string, vm *lua.LState) bool {
	if name == "" || vm == nil {
		return false
	}
	p.safe.Lock()
	defer p.safe.Unlock()
	if _, ok := p.pool[name]; ok {
		return false
	}
	p.pool[name] = vm
	return true
}

func (p *luaVmPool) Get(name string) (*lua.LState, bool) {
	if name == "" || p.pool == nil {
		return nil, false
	}
	p.safe.Lock()
	defer p.safe.Unlock()
	if v, ok := p.pool[name]; ok && v != nil {
		return v, true
	}
	return nil, false
}

func (p *luaVmPool) GetMust(name ...string) *lua.LState {
	if p.pool == nil {
		p.pool = make(map[string]*lua.LState)
	}
	if len(name) <= 0 {
		name = append(name, "default")
	}
	var id = name[0]
	if vm, ok := p.Get(id); ok && vm != nil {
		return vm
	}
	var vm = createMigrateVm()
	p.pool[id] = vm
	return vm
}

func (p *luaVmPool) destroy() {
	runtime.SetFinalizer(p, nil)
	p.safe.Lock()
	defer p.safe.Unlock()
	for _, v := range p.pool {
		if v != nil {
			v.Close()
		}
	}
	p.pool = nil
}

func (p *luaFsPool) destroy() {
	runtime.SetFinalizer(p, nil)
	p.safe.Lock()
	defer p.safe.Unlock()
	for _, v := range p.pool {
		if v != nil {
			_ = v.Close()
		}
	}
	p.pool = nil
}

func RegisterVm(name string, vm *lua.LState) bool {
	return vmPoolImpl.add(name, vm)
}

func GetVm(name string) *lua.LState {
	return vmPoolImpl.GetMust(name)
}

func RegisterFs(name string, system FileSystem) bool {
	return fsPoolImpl.add(name, system)
}

func GetFs(name string) (FileSystem, bool) {
	return fsPoolImpl.Get(name)
}
