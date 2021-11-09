package migrate

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/sirupsen/logrus"
	"github.com/weblfe/plugin_lua/core"
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
		length    string
		Comment    string
		Check      string
		Default    string
		After      string
		isFirst    bool
		Append     string
		isUnsigned bool
		isNotNull  bool
		isUnique   bool
		opt     *OptionKv
	}

	ColumnMap map[string]interface{}
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
	if argc <= 0 {
		return 0
	}
	switch argc {
	case 2:
	case 3:

	}
	return 0
}

func (l *LuaMigrateTable) createTable(table string, columns ColumnMap, append ...string) bool {
	var db = l.getDb()

	if err := db.Up(); err != nil {
		l.getLogger().Errorln("createTable.Error:", err)
		return false
	}
	return true
}

func (l *LuaMigrateTable) getLogger() *logrus.Logger {
	return logrus.New()
}

func (l *LuaMigrateTable) getDb() *migrate.Migrate {
	if m, ok := l.migrate[l.getArgOr("db", l.getConnName())]; ok {
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
