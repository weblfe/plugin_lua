package migrate

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/weblfe/plugin_lua/core"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

type (
	LuaSchemaBuilder struct {
		prefix  string
		driver  *migrate.Migrate
		methods map[string]lua.LGFunction
	}
)

func createSchemaBuilder(L *lua.LState) int {
	var args = core.GetArgs(L)
	if len(args) <= 0 {
		return 0
	}
	var (
		builder = NewSchemaBuilder()
		table   = L.NewTypeMetatable("schemaBuilder")
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
		"string":             builder.Str,
		"tinyint":            builder.TinyInt,
		"integer":            builder.Integer,
		"decimal":            builder.Decimal,
		"text":               builder.Text,
		"char":               builder.Char,
		"pk":                 builder.Pk,
		"upk":                builder.UPk,
		"bigpk":              builder.BigPk,
		"ubigpk":             builder.UBigPk,
		"datetime":           builder.DateTime,
		"comment":            builder.Comment,
		"smallint":           builder.SmallInt,
		"float":              builder.FloatNumber,
		"double":             builder.DoubleNumber,
		"bigint":             builder.BigInt,
		"date":               builder.Date,
		"money":              builder.Money,
		"binary":             builder.Binary,
		"createTable":        builder.CreateTable,
		"addColumn":          builder.AddColumn,
		"renameColumn":       builder.RenameColumn,
		"alterColumnComment": builder.AlterColumnComment,
		"addColumnComment":   builder.AlterColumnComment,
		"alterColumn":        builder.AlterColumn,
		"createIndex":        builder.CreateIndex,
		"dropTable":          builder.DropTable,
		"dropIndex":          builder.DropIndex,
		"dropColumn":         builder.DropColumn,
		"batchInsert":        builder.BatchInsert,
	}
}

func (builder *LuaSchemaBuilder) setDriver(d *migrate.Migrate) *LuaSchemaBuilder {
	if builder.driver == nil {
		builder.driver = d
	}
	return builder
}

func (builder *LuaSchemaBuilder) setPrefix(p string) *LuaSchemaBuilder {
	if builder.prefix == "" {
		builder.prefix = p
	}
	return builder
}

func (builder *LuaSchemaBuilder) table(t string) string {
	if strings.HasPrefix(t, "{{%") && strings.HasSuffix(t, "}}") {
		if builder.prefix == "" {
			return t
		}
		t = strings.TrimSuffix(t, "}}")
		t = strings.TrimPrefix(t, "{{%")
		if !strings.HasSuffix(builder.prefix, "_") {
			return fmt.Sprintf("%s_%s", builder.prefix, t)
		}
		return fmt.Sprintf("%s%s", builder.prefix, t)
	}
	return t
}

func (builder *LuaSchemaBuilder) Module(L *lua.LState) lua.LValue {
	var table = L.NewTypeMetatable("schemaBuilder")
	for k, v := range builder.methods {
		table.RawSet(lua.LString(k), L.NewFunction(v))
	}
	return table
}

func (builder *LuaSchemaBuilder) Str(L *lua.LState) int {
	var column = ColumnNew(String).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) TinyInt(L *lua.LState) int {
	var column = ColumnNew(TinyInt).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) Integer(L *lua.LState) int {
	var column = ColumnNew(Integer).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) Decimal(L *lua.LState) int {
	var column = ColumnNew(Decimal).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) Text(L *lua.LState) int {
	var column = ColumnNew(Text).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) Char(L *lua.LState) int {
	var column = ColumnNew(Char).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) Pk(L *lua.LState) int {
	var column = ColumnNew(Pk).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) BigPk(L *lua.LState) int {
	var column = ColumnNew(BigPk).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) UBigPk(L *lua.LState) int {
	var column = ColumnNew(UBigPk).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) UPk(L *lua.LState) int {
	var column = ColumnNew(UPk).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) DateTime(L *lua.LState) int {
	var column = ColumnNew(DateTime).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) Comment(L *lua.LState) int {
	var (
		code = `COMMENT("%s")`
		args = core.GetArgs(L)
		str  = args.GetString(0)
	)
	if args.Len() <= 0 || str == "" {
		L.Push(lua.LString(""))
	} else {
		L.Push(lua.LString(fmt.Sprintf(code, str)))
	}
	return 1
}

func (builder *LuaSchemaBuilder) SmallInt(L *lua.LState) int {
	var column = ColumnNew(SmallInt).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) FloatNumber(L *lua.LState) int {
	var column = ColumnNew(Float).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) DoubleNumber(L *lua.LState) int {
	var column = ColumnNew(Double).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) BigInt(L *lua.LState) int {
	var column = ColumnNew(BigInteger).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) Date(L *lua.LState) int {
	var column = ColumnNew(Date).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) Money(L *lua.LState) int {
	var column = ColumnNew(Money).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
}

func (builder *LuaSchemaBuilder) Binary(L *lua.LState) int {
	var column = ColumnNew(Binary).SetArgs(core.GetArgs(L))
	L.Push(column.LuaObject(L))
	return 1
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

func (builder *LuaSchemaBuilder) String() string {
		return ``
}