package core

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"strconv"
)

type (
	LuaRegistryFunction struct {
		LName     string
		LFunction lua.LGFunction
	}
	LuaArguments []interface{}
)

func GetArgs(L *lua.LState) LuaArguments {
	var argc = L.GetTop()
	if argc <= 0 {
		return nil
	}
	var args []interface{}
	for i := 1; i <= argc; i++ {
		v := L.CheckAny(i)
		if v == lua.LNil {
			continue
		}
		switch v.Type() {
		case lua.LTNil:
			args = append(args, v.String())
		case lua.LTBool:
			args = append(args, L.CheckBool(i))
		case lua.LTNumber:
			args = append(args, L.CheckNumber(i).String())
		case lua.LTString:
			args = append(args, L.CheckString(i))
		case lua.LTFunction:
			args = append(args, L.CheckFunction(i).String())
		case lua.LTUserData:
			args = append(args, L.CheckUserData(i))
		case lua.LTThread:
			args = append(args, L.CheckThread(i))
		case lua.LTTable:
			args = append(args, L.CheckTable(i))
		case lua.LTChannel:
			args = append(args, L.CheckChannel(i))
		default:
			args = append(args, v.String())
		}
	}
	if len(args) <= 0 {
		return nil
	}
	return args
}

func (arguments LuaArguments) Len() int {
	return len(arguments)
}

func (arguments LuaArguments) Set(v interface{}, index ...int) LuaArguments {
	index = append(index, 0)
	arguments[index[0]] = v
	return arguments
}

func (arguments LuaArguments) IsEmpty() bool {
	if arguments == nil || len(arguments) <= 0 {
		return true
	}
	return false
}

func (arguments LuaArguments) IsLuaType(i uint) bool {
	var argc = uint(arguments.Len())
	if argc <= i {
		return false
	}
	if v := arguments[int(i)]; v != nil {
		if _, ok := v.(lua.LValue); ok {
			return true
		}
	}
	return false
}

func (arguments LuaArguments) Arg(i uint) (interface{}, bool) {
	var argc = uint(arguments.Len())
	if argc <= i {
		return nil, false
	}
	if v := arguments[int(i)]; v != nil {
		return v, true
	}
	return nil, false
}

func (arguments LuaArguments) Get(i uint) interface{} {
	var argc = uint(arguments.Len())
	if argc <= i {
		return nil
	}
	if v := arguments[int(i)]; v != nil {
		return v
	}
	return nil
}

func (arguments LuaArguments) GetInt(i uint) int {
	var argc = uint(arguments.Len())
	if argc <= i {
		return 0
	}
	if v := arguments[int(i)]; v != nil {
		switch v.(type) {
		case float64:
			return int(v.(float64))
		case string:
			str := v.(string)
			if n, err := strconv.Atoi(str); err == nil {
				return n
			}
			return 0
		case lua.LNumber:
			return int(v.(lua.LNumber))
		}
	}
	return 0
}

func (arguments LuaArguments) GetNumber(i uint) float64 {
	var argc = uint(arguments.Len())
	if argc <= i {
		return 0
	}
	if v := arguments[int(i)]; v != nil {
		switch v.(type) {
		case float64:
			return v.(float64)
		case string:
			str := v.(string)
			if n, err := strconv.ParseFloat(str, 64); err == nil {
				return n
			}
			return 0
		case lua.LNumber:
			return float64(v.(lua.LNumber))
		}
	}
	return 0
}

func (arguments LuaArguments) GetString(i uint) string {
	var argc = uint(arguments.Len())
	if argc <= i {
		return ""
	}
	if v := arguments[int(i)]; v != nil {
		switch v.(type) {
		case string:
			return v.(string)
		case lua.LString:
			return string(v.(lua.LString))
		case lua.LValue:
			return v.(lua.LValue).String()
		case fmt.Stringer:
			return v.(fmt.Stringer).String()
		case fmt.GoStringer:
			return v.(fmt.GoStringer).GoString()
		}
	}
	return ""
}

func (arguments LuaArguments) GetTable(i uint) *lua.LTable {
	var argc = uint(arguments.Len())
	if argc <= i {
		return nil
	}
	if v := arguments[int(i)]; v != nil {
		switch v.(type) {
		case lua.LTable:
			tab := v.(lua.LTable)
			return &tab
		case *lua.LTable:
			tab := v.(*lua.LTable)
			return tab
		}
	}
	return nil
}

func (arguments LuaArguments) GetBool(i uint) bool {
	var argc = uint(arguments.Len())
	if argc <= i {
		return false
	}
	if v := arguments[int(i)]; v != nil {
		switch v.(type) {
		case bool:
			return v.(bool)
		case lua.LBool:
			return bool(v.(lua.LBool))
		case lua.LNumber:
			n := v.(lua.LNumber)
			if float64(n) > 0 {
				return true
			}
			return false
		}
	}
	return false
}

func (arguments LuaArguments) GetChannel(i uint) lua.LChannel {
	var argc = uint(arguments.Len())
	if argc <= i {
		return nil
	}
	if v := arguments[int(i)]; v != nil {
		switch v.(type) {
		case lua.LChannel:
			return v.(lua.LChannel)
		case *lua.LChannel:
			return *v.(*lua.LChannel)
		}
	}
	return nil
}

func (arguments LuaArguments) GetFunc(i uint) *lua.LFunction {
	var argc = uint(arguments.Len())
	if argc <= i {
		return nil
	}
	if v := arguments[int(i)]; v != nil {
		switch v.(type) {
		case lua.LFunction:
			var fn = v.(lua.LFunction)
			return &fn
		case *lua.LFunction:
			return v.(*lua.LFunction)
		}
	}
	return nil
}

func (arguments LuaArguments) GetUserData(i uint) *lua.LUserData {
	var argc = uint(arguments.Len())
	if argc <= i {
		return nil
	}
	if v := arguments[int(i)]; v != nil {
		switch v.(type) {
		case lua.LUserData:
			var fn = v.(lua.LUserData)
			return &fn
		case *lua.LUserData:
			return v.(*lua.LUserData)
		case interface{}:
			var d = lua.LUserData{
				Value: v,
			}
			return &d
		}
	}
	return nil
}

func (arguments LuaArguments) GetThread(i uint) *lua.LState {
	var argc = uint(arguments.Len())
	if argc <= i {
		return nil
	}
	if v := arguments[int(i)]; v != nil {
		switch v.(type) {
		case lua.LState:
			var fn = v.(lua.LState)
			return &fn
		case *lua.LState:
			return v.(*lua.LState)
		}
	}
	return nil
}

func (arguments *LuaArguments) Pop() (interface{}, bool) {
	var argc = uint(arguments.Len())
	if argc <= 0 {
		return nil, false
	}
	if v := (*arguments)[0]; v != nil {
		if argc <= 1 {
			*arguments = nil
		} else {
			*arguments = (*arguments)[1:]
		}
		return v, true
	}
	return nil, false
}
