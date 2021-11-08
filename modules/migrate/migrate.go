package migrate

import (
		"fmt"
		"github.com/golang-migrate/migrate/v4"
		_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
		_ "github.com/golang-migrate/migrate/v4/database/mysql"
		_ "github.com/golang-migrate/migrate/v4/database/postgres"
		_ "github.com/golang-migrate/migrate/v4/source/file"
		_ "github.com/golang-migrate/migrate/v4/source/github"
		"github.com/weblfe/plugin_lua/core"
		lua "github.com/yuin/gopher-lua"
		"sync"
)

type (

	LuaMigrateTable struct {
		loggerName string
		safe       sync.RWMutex
		options    map[string]*OptionKv
		migrate    map[string]*migrate.Migrate
	}

	LuaMigrateColumn struct {
		TypeName string
		Value    string
		Len      uint
		Comment  string
		Check    string
		Default  string
	}

	OptionKv struct {
		Source  string
		ConnUrl string
	}
)

var (
	defaultMigrate = NewLuaMigrate()
	Funcs          = defaultMigrate.methods()
)

func NewLuaMigrate() *LuaMigrateTable {
	var table = new(LuaMigrateTable)
	return table.init()
}

func createMigrate(L *lua.LState) int  {
		var args = core.GetArgs(L)
		if len(args) <= 0 {
				return 0
		}
		var (
				m = NewLuaMigrate()
				table  = L.NewTypeMetatable(Name)
		)
		switch len(args) {
		case 1:

		case 2:

		case 3:

		}
		for k, fn := range m.methods() {
				table.RawSet(lua.LString(k), L.NewFunction(fn))
		}
		L.Push(table)
		return 1
}

func (c *LuaMigrateColumn) String() string {
	return fmt.Sprintf("")
}

func (l *LuaMigrateTable) init() *LuaMigrateTable {
	return l
}

func (l *LuaMigrateTable) methods() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"new":                createMigrate,
		"connection":         l.Connection,
	}
}

func (l *LuaMigrateTable) Connection(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) Comment(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) ConnDefault(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) CreateTable(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) DropIndex(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) BatchInsert(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) Integer(state *lua.LState) int {
	return 0
}

// AddColumn 添加字段
func (l *LuaMigrateTable) AddColumn(state *lua.LState) int {
	return 0
}

// RenameColumn 重命名字段
func (l *LuaMigrateTable) RenameColumn(state *lua.LState) int {
	return 0
}

// AlterColumn 修改字段类型
func (l *LuaMigrateTable) AlterColumn(state *lua.LState) int {
	return 0
}

// AlterColumnComment 添加字段备注
func (l *LuaMigrateTable) AlterColumnComment(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) Char(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) FloatNumber(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) DoubleNumber(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) Money(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) BigInt(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) Decimal(state *lua.LState) int {
	return 0
}

// Binary 二进制
func (l *LuaMigrateTable) Binary(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) SmallInt(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) DropTable(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) DropColumn(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) CreateIndex(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) String(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) TinyInt(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) Date(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) Text(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) Pk(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) BigPk(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) DateTime(state *lua.LState) int {
	return 0
}
