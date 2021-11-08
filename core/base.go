package core

import "github.com/yuin/gopher-lua"

type (
	LuaRegistryFunction struct {
		LName     string
		LFunction lua.LGFunction
	}
)


func GetArgs(L *lua.LState) []interface{} {
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
						args = append(args, L.CheckUserData(i).String())
				case lua.LTThread:
						args = append(args, L.CheckThread(i).String())
				case lua.LTTable:
						args = append(args, L.CheckTable(i).String())
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
