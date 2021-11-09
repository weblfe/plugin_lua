package plugins

import (
	"fmt"
	"testing"
)

func TestNewLua(t *testing.T) {
	var (
		method    = "safeUp"
		module    = "create_queue_info_table"
		file      = `./testdata/create_queue_info_table.lua`
		luaParser = NewLua().SetLoader(CreateExtendsLoader)
	)
	luaParser.Boot()
	if err2 := luaParser.DoFile(file); err2 == nil {
		var (
			expr = `migration=module("%s");migration.%s();`
			code = fmt.Sprintf(expr, module, method)
		// libs = luaParser.Libs()
		)
		if err3 := luaParser.EvalExpr(code); err3 != nil {
			t.Error(err3)
		}
	}
}

func TestLuaPluginImpl_Libs(t *testing.T) {
	var (
		luaParser = NewLua().SetLoader(CreateExtendsLoader)
	)
	luaParser.Boot()
	libs := luaParser.Libs()
	if len(libs) <= 0 {
		t.Error("库加载失败")
	}
}

func TestLuaPluginImpl_LoadFile(t *testing.T) {
	var (
		file      = `./testdata/logger_test.lua`
		luaParser = NewLua().SetLoader(CreateExtendsLoader)
	)
	// migrate.New("lua://./testdata/migrates","mysql://root@123:127.0.0.1:3306/test?charset=utf8mb")
	luaParser.Boot()
	err2 := luaParser.DoFile(file)
	if err2 != nil {
		t.Error(err2)
	}
}
