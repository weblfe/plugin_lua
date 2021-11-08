package migrate

import (
	"github.com/golang-migrate/migrate/v4"
		"github.com/weblfe/plugin_lua/core"
		lua "github.com/yuin/gopher-lua"
)

type (
	LuaSchemaBuilder struct {
		driver  *migrate.Migrate
		methods map[string]lua.LGFunction
	}
)

func createSchemaBuilder(L *lua.LState) int  {
		var args = core.GetArgs(L)
		if len(args) <= 0 {
				return 0
		}
		var (
				builder = NewSchemaBuilder()
				table  = L.NewTypeMetatable("schemaBuilder")
		)
		switch len(args) {
		case 1:

		case 2:

		case 3:

		}
		for k, fn := range builder.methods {
				table.RawSet(lua.LString(k), L.NewFunction(fn))
		}
		L.Push(table)
		return 1
}

func NewSchemaBuilder() *LuaSchemaBuilder {
	var builder = new(LuaSchemaBuilder)
	return builder.init()
}

func (builder *LuaSchemaBuilder) init() *LuaSchemaBuilder {
	builder.methods = builder.loads()
	return builder
}

func (builder *LuaSchemaBuilder) loads() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"string":             builder.String,
		"tinyint":            builder.TinyInt,
		"integer":            builder.Integer,
		"decimal":            builder.Decimal,
		"text":               builder.Text,
		"char":               builder.Char,
		"pk":                 builder.Pk,
		"bigPk":              builder.BigPk,
		"datetime":           builder.DateTime,
		"comment":            builder.Comment,
		"createTable":        builder.CreateTable,
		"addColumn":          builder.AddColumn,
		"renameColumn":       builder.RenameColumn,
		"alterColumnComment": builder.AlterColumnComment,
		"addColumnComment":   builder.AlterColumnComment,
		"alterColumn":        builder.AlterColumn,
		"createIndex":        builder.CreateIndex,
		"connDefault":        builder.ConnDefault,
		"dropTable":          builder.DropTable,
		"dropIndex":          builder.DropIndex,
		"dropColumn":         builder.DropColumn,
		"batchInsert":        builder.BatchInsert,
		"smallint":           builder.SmallInt,
		"float":              builder.FloatNumber,
		"double":             builder.DoubleNumber,
		"bigint":             builder.BigInt,
		"date":               builder.Date,
		"money":              builder.Money,
		"binary":             builder.Binary,
	}
}

func (builder *LuaSchemaBuilder) setDriver(d *migrate.Migrate) *LuaSchemaBuilder {
	if builder.driver == nil {
		builder.driver = d
	}
	return builder
}

func (builder *LuaSchemaBuilder) String(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) TinyInt(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) Integer(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) Decimal(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) Text(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) Char(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) Pk(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) BigPk(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) DateTime(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) Comment(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) CreateTable(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) AddColumn(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) RenameColumn(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) AlterColumnComment(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) AlterColumn(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) CreateIndex(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) ConnDefault(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) DropTable(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) DropIndex(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) DropColumn(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) BatchInsert(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) SmallInt(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) FloatNumber(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) DoubleNumber(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) BigInt(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) Date(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) Money(L *lua.LState) int {
	return 0
}

func (builder *LuaSchemaBuilder) Binary(L *lua.LState) int {
	return 0
}
