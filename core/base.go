package core

import "github.com/yuin/gopher-lua"

type (
	LuaRegistryFunction struct {
		LName     string
		LFunction lua.LGFunction
	}
)
