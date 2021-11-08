package migrate

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	lua "github.com/yuin/gopher-lua"
	"sync"
)

type (
	LuaMigrateTable struct {
		loggerName string
		safe       sync.RWMutex
		migrate    map[string]*migrate.Migrate
	}

	LuaMigrateColumn struct {
		TypeName string
		Value    string
		Comment  string
		Check    string
		Default  string
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

func NewColumn(ty string) *LuaMigrateColumn {
	var column = new(LuaMigrateColumn)
	column.TypeName = ty
	return column
}

func (c *LuaMigrateColumn) String() string {
	return fmt.Sprintf("")
}

func (l *LuaMigrateTable) init() *LuaMigrateTable {
	return l
}

func (l *LuaMigrateTable) methods() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"connection":  l.Connection,
		"string":      l.String,
		"tinyint":     l.TinyInt,
		"integer":     l.Integer,
		"decimal":     l.Decimal,
		"text":        l.Text,
		"char":        l.Char,
		"pk":          l.Pk,
		"bigPk":       l.BigPk,
		"datetime":    l.DateTime,
		"comment":     l.Comment,
		"createTable": l.CreateTable,
		"createIndex": l.CreateIndex,
		"connDefault": l.ConnDefault,
		"dropTable":   l.DropTable,
		"dropIndex":   l.DropIndex,
		"dropColumn":  l.DropColumn,
		"batchInsert": l.BatchInsert,
		"smallint":    l.SmallInt,
		"float":       l.FloatNumber,
		"double":      l.DoubleNumber,
		"bigint":      l.BigInt,
		"date":        l.Date,
		"money":       l.Money,
		"binary":      l.Binary,
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
