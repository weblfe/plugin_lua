package migrate

import (
	"bytes"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/sirupsen/logrus"
	"github.com/weblfe/plugin_lua/core"
	"github.com/weblfe/plugin_lua/query"
	lua "github.com/yuin/gopher-lua"
	"sort"
	"sync"
)

type (
	LuaMigrateTable struct {
		loggerName  string
		constructor sync.Once
		safe        sync.RWMutex
		args        map[string]string
		options     map[string]*OptionKv
		migrate     map[string]*migrate.Migrate
	}

	LuaConn struct {
		Name string
		Conn *migrate.Migrate
	}

	LuaMigrateColumn struct {
		Type       ColumnType
		Value      string
		Len        []int
		length     string
		Comment    string
		Check      string
		Default    string
		After      string
		isFirst    bool
		Append     string
		isUnsigned bool
		isNotNull  bool
		isUnique   bool
		opt        *OptionKv
	}

	ColumnMap []*Column

	Column struct {
		Name    string
		RawBody fmt.Stringer
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

func createMigrate(L *lua.LState) int {
	var args = core.GetArgs(L)
	if len(args) <= 0 {
		return 0
	}
	var (
		m     = NewLuaMigrate()
		table = L.NewTypeMetatable(Name)
	)
	switch len(args) {
	case 1:
	case 2:
	case 3:
	}
	table.Metatable = &lua.LUserData{
		Value: m,
		Env:   L.Env,
	}
	for k, fn := range m.methods() {
		table.RawSet(lua.LString(k), L.NewFunction(fn))
	}

	L.Push(table)
	return 1
}

func (c *Column) Body() string {
	if c.RawBody != nil {
		return c.RawBody.String()
	}
	return ""
}

func (c *Column) Key() string {
	return c.Name
}

func (c *Column) Value() fmt.Stringer {
	return query.NewString(c.Body())
}

func (c *Column) String() string {
	return fmt.Sprintf(`%s %s`, c.Key(), c.Value())
}

func (c *Column) IsString() bool {
	return false
}

func (l *LuaMigrateTable) init() *LuaMigrateTable {
	l.safe = sync.RWMutex{}
	l.constructor = sync.Once{}
	l.options = make(map[string]*OptionKv)
	l.migrate = make(map[string]*migrate.Migrate)
	l.loggerName = core.GetEnvOr("migrate_logger", "migrate")
	return l
}

func (l *LuaMigrateTable) Boot() {
	if l.migrate == nil || l.options == nil {
		return
	}
	l.constructor.Do(func() {
		for k, v := range l.options {
			l.migrate[k] = l.conn(v)
		}
	})
}

func (l *LuaMigrateTable) conn(opt *OptionKv) *migrate.Migrate {
	if v, err := migrate.New(opt.Source, opt.ConnUrl); err == nil {
		return v
	}
	return nil
}

func (l *LuaMigrateTable) methods() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"new":         createMigrate,
		"connection":  l.Connection,
		"connDefault": l.ConnDefault,
	}
}

func (l *LuaMigrateTable) Connection(state *lua.LState) int {
	var (
		name string
		args = core.GetArgs(state)
	)
	l.Boot()
	if args.IsEmpty() {
		name = l.getConnName()
	} else {
		name = args.GetString(0)
	}
	if name == "" {
		name = l.getConnName()
	}
	if v, ok := l.migrate[name]; ok {
		var (
			opt     = l.getOption(name)
			builder = NewSchemaBuilder().setDriver(v)
		)
		if opt != nil {
			builder.setPrefix(opt.Prefix)
		}
		state.Push(builder.Module(state))
		return 1
	}
	state.Push(NewSchemaBuilder().Module(state))
	return 1
}

func (l *LuaMigrateTable) getOption(key string) *OptionKv {
	if v, ok := l.options[key]; ok {
		return v
	}
	var opt = OptionKv{
		Source:  core.GetEnvOr(core.SprintfEnv("%s_source", key)),
		ConnUrl: core.GetEnvOr(core.SprintfEnv("%s_conn_url", key)),
		Prefix:  core.GetEnvOr(core.SprintfEnv("%s_table_prefix", key)),
	}
	return &opt
}

func (l *LuaMigrateTable) setPrefix(p string) *LuaMigrateTable {
	if v, ok := l.options["default"]; ok {
		v.Prefix = p
	}
	return l
}

func (l *LuaMigrateTable) Comment(state *lua.LState) int {
	var (
		code = `comment("%s")`
		args = core.GetArgs(state)
		str  = args.GetString(0)
	)
	if args.Len() <= 0 || str == "" {
		state.Push(lua.LString(""))
	} else {
		state.Push(lua.LString(fmt.Sprintf(code, str)))
	}
	return 1
}

func (l *LuaMigrateTable) ConnDefault(state *lua.LState) int {
	var conn = l.getConnName()
	state.Push(lua.LString(conn))
	if conn == "" {
		return 0
	}
	return 1
}

func (l *LuaMigrateTable) getConnName() string {
	if v, ok := l.options["default"]; ok && v != nil {
		return "default"
	}
	var conn = l.getConns()
	if len(conn) <= 0 {
		return ""
	}
	return conn[0].Name
}

func (l *LuaMigrateTable) getConns() []*LuaConn {
	var (
		keys  []string
		conns []*LuaConn
	)
	for k := range l.migrate {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		conns = append(conns, &LuaConn{Name: k, Conn: l.migrate[k]})
	}
	return conns
}

func (l *LuaMigrateTable) CreateTable(state *lua.LState) int {
	var (
		args = core.GetArgs(state)
		argc = args.Len()
	)
	if argc < 2 {
		// @todo lua error
		logrus.Error("createTable require less 2 params ")
		return 0
	}
	var (
		tableName  = args.GetString(0)        // table
		obj        = args.GetTable(1)         // columns
		reader     = state.GetGlobal(GBuffer) // buffer
		appendSql  = args.GetString(2)        // append
		columnsMap = CreateColumnMapByTable(obj)
	)
	// 执行sql 解析逻辑
	switch reader.Type() {
	case lua.LTNil:
		logrus.Error("global sql buffer is nil")
		return 0
	case lua.LTUserData:
		var sql = l.createTable(tableName, columnsMap, appendSql)
		if u, ok := reader.(*lua.LUserData); ok {
			if buffer, ok := u.Value.(*bytes.Buffer); ok {
				buffer.Write([]byte(sql))
				return 1
			}
		}
	}
	return 0
}

// createTable 构造创建表 command sql
func (l *LuaMigrateTable) createTable(table string, columns ColumnMap, append ...string) string {
	var (
		size      = len(columns)
		db        = l.getDbOption()
		realTable = db.quoteTableName(table)
		buffer    = bytes.NewBufferString(fmt.Sprintf(`CREATE TABLE %s`, realTable))
	)
	if columns == nil {
		return ""
	}
	for i, column := range columns {
		var sql = `%s %s`
		if i < size-1 {
			sql = sql + `,\n`
		} else {
			sql = sql + `\n`
		}
		buffer.Write([]byte(fmt.Sprintf(`%s %s \n`, column.Key(), column.Body())))
	}
	if len(append) > 0 {
		for _, v := range append {
			if v == "" {
				continue
			}
			buffer.Write([]byte(v))
		}
	}
	buffer.Write([]byte(`;`))
	return buffer.String()
}

func (l *LuaMigrateTable) getLogger() *logrus.Logger {
	return logrus.New()
}

func (l *LuaMigrateTable) getDbOption() *OptionKv {
	if m, ok := l.options[l.getArgOr("db", l.getConnName())]; ok {
		return m
	}
	return nil
}

func (l *LuaMigrateTable) getArgOr(key string, or ...string) string {
	or = append(or, "")
	if key == "" {
		return or[0]
	}
	if v, ok := l.args[key]; ok {
		return v
	}
	return or[0]
}

func (l *LuaMigrateTable) DropIndex(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) BatchInsert(state *lua.LState) int {
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

func (l *LuaMigrateTable) DropTable(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) DropColumn(state *lua.LState) int {
	return 0
}

func (l *LuaMigrateTable) CreateIndex(state *lua.LState) int {
	return 0
}

func CreateColumnMapByTable(t *lua.LTable) ColumnMap {
	if t == nil {
		return nil
	}
	var cMap []*Column
	// 取出column
	t.ForEach(func(key lua.LValue, v lua.LValue) {
		if key == nil || v == nil {
			return
		}
		var c *Column
		if s, ok := key.(lua.LString); ok && s != "" {
			c = NewColumn(string(s), v)
		}
		if c != nil {
			cMap = append(cMap, c)
		}
	})
	return cMap
}

func NewColumn(key string, v lua.LValue) *Column {
	if v == nil || key == "" {
		return nil
	}
	if m, ok := v.(*lua.LTable); ok {
		var ref = m.Metatable
		if ref == nil {
			return nil
		}
		u, ok := ref.(*lua.LUserData)
		if !ok || u.Value == nil {
			return nil
		}
		if body, ok := u.Value.(fmt.Stringer); ok {
			return &Column{
				Name:    key,
				RawBody: body,
			}
		}
	}
	return nil
}
