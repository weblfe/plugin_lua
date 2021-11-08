package logger

import "github.com/yuin/gopher-lua"

const (
		Name = "logger"
)

func  NewLuaLoggerTables() lua.LGFunction{
		return func(state *lua.LState) int {
				var mod lua.LValue
				if len(Funcs) <= 0 {
						return 0
				}
				mod = state.RegisterModule(Name, Funcs)
				state.Push(mod)
				return 1
		}
}