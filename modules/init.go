package modules

import (
	"github.com/weblfe/plugin_lua/core"
	"github.com/weblfe/plugin_lua/modules/logger"
	"github.com/weblfe/plugin_lua/modules/migrate"
)

var (
	_modules []*core.LuaRegistryFunction
)

func init() {
	if _modules == nil {
		_modules = []*core.LuaRegistryFunction{
			{
				LName:     logger.Name,
				LFunction: logger.NewLuaLoggerTables(),
			},
			{
				LName:     migrate.Name,
				LFunction: migrate.NewLuaMigrateTables(),
			},
		}
	}
}

func GetModules() []*core.LuaRegistryFunction {
	return _modules
}
