package plugins

import (
	"fmt"
	"testing"
)

func TestNewLua(t *testing.T) {
	var (
		method    = "safeUp"
		module    = "create_queue_info_table"
		file      = `testdata/create_queue_info_table.lua`
		luaParser = NewLua().SetLoader(CreateExtendsLoader)
	)
	luaParser.Boot()
	if fnv, err2 := luaParser.LoadFile(file); err2 == nil || fnv != nil {
		var (
			expr = `migrate=require("%s")
migrate.%s`
			code = fmt.Sprintf(expr, module, method)
		)
		if err3 := luaParser.EvalExpr(code); err3 != nil {
			t.Error(err3)
		}
	}
}

func TestLuaPluginImpl_LoadFile(t *testing.T) {
	var (
		file      = `./testdata/logger_test.lua`
		luaParser = NewLua().SetLoader(CreateExtendsLoader)
	)
	luaParser.Boot()
	err2 := luaParser.DoFile(file)
	if err2 != nil {
		t.Error(err2)
	}
}
