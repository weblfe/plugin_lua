package migrate

import (
	"fmt"
	"github.com/weblfe/plugin_lua/core"
	lua "github.com/yuin/gopher-lua"
	"io"
	"strconv"
	"strings"
)

type (
	Kv struct {
		K string
		V string
	}
)

// ColumnOf 构建字段
func ColumnOf(args []interface{}) *LuaMigrateColumn {
	var argc = len(args)
	if argc <= 0 {
		return nil
	}
	var column = new(LuaMigrateColumn)
	column.isNotNull = true
	switch argc {
	case 1:
		if str, ok := args[0].(string); ok && str != "" {
			column.setType(str)
		}
		if column.check() {
			return column
		}
	case 2:
		if str, ok := args[0].(string); ok && str != "" {
			column.setType(str)
		}
		if num := args[1]; num != nil {
			column.setLength(num)
		}
		if !column.check() {
			return nil
		}
		return column
	case 3:
		if str, ok := args[0].(string); ok && str != "" {
			column.setType(str)
		}
		if num := args[1]; num != nil {
			column.setLength(num)
		}
		if v, ok := args[2].(string); ok && v != "" {
			column.setDefault(v)
		}
		if !column.check() {
			return nil
		}
		return column
	}
	return nil
}

// ColumnNew 创建
func ColumnNew(ty ColumnType) *LuaMigrateColumn {
	var column = new(LuaMigrateColumn)
	column.Type = ty
	column.check()
	column.isNotNull = true
	return column
}

func (c *LuaMigrateColumn) check() bool {
	if c.Type == "" {
		return false
	}
	if !c.Type.Check() {
		return false
	}
	if c.Len == nil {
		c.Len = c.Type.DefaultSize()
	}
	return true
}

func (c *LuaMigrateColumn) setCheck(check string) *LuaMigrateColumn {
	c.Check = check
	return c
}

func (c *LuaMigrateColumn) SetArgs(args []interface{}) *LuaMigrateColumn {
	var argc = len(args)
	switch argc {
	case 1:
		if num := args[0]; num != nil {
			c.setLength(num)
		}
		if !c.check() {
			return nil
		}
		return c
	case 2:
		if num := args[0]; num != nil {
			c.setLength(num)
		}
		if v, ok := args[1].(string); ok && v != "" {
			c.setDefault(v)
		}
		if !c.check() {
			return nil
		}
		return c
	}
	return c
}

func (c *LuaMigrateColumn) setLength(value interface{}) *LuaMigrateColumn {
	if c.Type.DefaultSize() == nil {
		return c
	}
	var tmp []interface{}
	// lua table array
	if t, ok := value.(*lua.LTable); ok {
		t.ForEach(func(key lua.LValue, v lua.LValue) {
			if key.Type() == lua.LTNumber {
				tmp = append(tmp, v.String())
			}
		})
		if len(tmp) > 0 {
			value = tmp
		}
	}
	// singe value/ go types
	switch value.(type) {
	case string, lua.LString:
		if num, ok := value.(string); ok && num != "" {
			if n, err := strconv.ParseFloat(num, 64); err == nil {
				c.Len = []int{int(n)}
			}
		}
	case []interface{}:
		var arr []int
		for _, v := range value.([]interface{}) {
			switch v.(type) {
			case string:
				if num, ok := v.(string); ok && num != "" {
					if n, err := strconv.ParseFloat(num, 64); err == nil {
						arr = append(arr, int(n))
					}
				}
			case int:
				if num, ok := v.(int); ok {
					arr = append(arr, num)
				}
			case float64:
				if num, ok := v.(float64); ok {
					arr = append(arr, int(num))
				}
			}
		}
		if len(arr) > 0 {
			c.Len = arr
		}
	case float64, lua.LNumber:
		if num, ok := value.(float64); ok && num > 0 {
			c.Len = []int{int(num)}
		}
	case int:
		if num, ok := value.(int); ok && num > 0 {
			c.Len = []int{num}
		}
	}
	return c
}

func (c *LuaMigrateColumn) setDefault(v string) *LuaMigrateColumn {
	c.Default = v
	return c
}

func (c *LuaMigrateColumn) setType(t string) *LuaMigrateColumn {
	c.Type = ColumnType(t)
	return c.initDefault()
}

func (c *LuaMigrateColumn) initDefault() *LuaMigrateColumn {

	return c
}

func (c *LuaMigrateColumn) Write(writer io.Writer) error {
	if _, err := writer.Write(c.Bytes()); err != nil {
		return err
	}
	return nil
}

func (c *LuaMigrateColumn) Bytes() []byte {
	var format = ``
	switch c.Type.GetTypeCategory() {
	case CategoryPk:
		format = `{type}{check}{comment}{append}`
	default:
		format = `{type}{length}{notnull}{unique}{default}{check}{comment}{append}`
	}
	return []byte(c.buildCompleteString(format))
}

func (c *LuaMigrateColumn) LuaObject(L *lua.LState) lua.LValue {
	var table = L.NewTypeMetatable("column")
	for k, v := range c.methods() {
		table.RawSet(lua.LString(k), L.NewClosure(v, table))
	}
	return table
}

func (c *LuaMigrateColumn) toString(L *lua.LState) int {
	var str = c.Bytes()
	L.Push(lua.LString(str))
	return 1
}

func (c *LuaMigrateColumn) methods() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"null":     c.SetNull,
		"notNull":  c.SetNotNull,
		"nullable":  c.SetNull,
		"size":     c.SetSize,
		"comment":  c.SetComment,
		"default":  c.SetDefault,
		"toString": c.toString,
	}
}

func (c *LuaMigrateColumn) SetNotNull(state *lua.LState) int {
	var value = state.Get(lua.UpvalueIndex(1))
	c.nullable(true)
	if value != nil {
		state.Push(value)
	}
	return 1
}

func (c *LuaMigrateColumn) SetNull(state *lua.LState) int {
	var value = state.Get(lua.UpvalueIndex(1))
	c.nullable(false)
	if value != nil {
		state.Push(value)
	}
	return 1
}

func (c *LuaMigrateColumn) SetSize(state *lua.LState) int {
	var value = state.Get(lua.UpvalueIndex(1))
	if value != nil {
		state.Push(value)
	}
	return 1
}

func (c *LuaMigrateColumn) SetComment(state *lua.LState) int {
	var (
		value = state.Get(lua.UpvalueIndex(1))
		args  = core.GetArgs(state)
	)
	c.comment(args.GetString(0))
	if value != nil {
		state.Push(value)
	}
	return 1
}

func (c *LuaMigrateColumn) comment(comment string) *LuaMigrateColumn {
	c.Comment = comment
	return c
}

func (c *LuaMigrateColumn) SetDefault(state *lua.LState) int {
	var (
		value = state.Get(lua.UpvalueIndex(1))
		args  = core.GetArgs(state)
	)
	c.defaultValue(args.GetString(0))
	if value != nil {
		state.Push(value)
	}
	return 1
}

func (c *LuaMigrateColumn) defaultValue(value string) *LuaMigrateColumn {
	c.Default = value
	return c
}

func (c *LuaMigrateColumn) SetUnsigned(state *lua.LState) int {
	var value = state.Get(lua.UpvalueIndex(1))
	c.unsigned()
	if value != nil {
		state.Push(value)
	}
	return 1
}

func (c *LuaMigrateColumn) unsigned() *LuaMigrateColumn {
	switch c.Type {
	case Pk:
		c.Type = UPk
	case BigPk:
		c.Type = UBigPk
	}
	c.isUnsigned = true
	return c
}

func (c *LuaMigrateColumn) nullable(v bool) *LuaMigrateColumn {
	c.isNotNull = v
	return c
}

func (c *LuaMigrateColumn) buildCompleteString(format string) string {
	var placeholderValues = []Kv{
		{"{type}", c.buildTypeString()},
		{"{length}", c.buildLengthString()},
		{"{unsigned}", c.buildUnsignedString()},
		{"{notnull}", c.buildNotNullString()},
		{"{unique}", c.buildUniqueString()},
		{"{default}", c.buildDefaultString()},
		{"{check}", c.buildCheckString()},
		{"{comment}", c.buildCommentString()},
		{"{pos}", c.buildPosString()},
		{"{append}", c.buildAppendString()},
	}
	return c.formatReplace(format, placeholderValues)
}

func (c *LuaMigrateColumn) buildTypeString() string {
	if c.Type.Check() {
		return c.Type.String()
	}
	return ""
}

func (c *LuaMigrateColumn) buildLengthString() string {
	if c.length != "" {
		return c.length
	}
	if c.Len != nil {
		var size = len(c.Len)
		if size == 1 {
			c.length = fmt.Sprintf("(%d)", c.Len[0])
		} else {
			c.length = fmt.Sprintf("(%d,%d)", c.Len[0], c.Len[1])
		}
	}
	return c.length
}

func (c *LuaMigrateColumn) buildUnsignedString() string {
	return ""
}

func (c *LuaMigrateColumn)GetCategoryMap() map[ColumnType]string  {
		return categoryMap
}

func (c *LuaMigrateColumn) buildNotNullString() string {
	if c.isNotNull {
		return ` NOT NULL`
	}
	return ` NULL`
}

func (c *LuaMigrateColumn) buildUniqueString() string {
	if c.isUnique {
		return " UNIQUE"
	}
	return ""
}

func (c *LuaMigrateColumn) buildDefaultString() string {
		var defaultValue = c.buildDefaultValue()
		if defaultValue == "" {
				return ""
		}
		return fmt.Sprintf(` DEFAULT %s` ,defaultValue)
}

func (c *LuaMigrateColumn)buildDefaultValue() string  {
		if c.Default == "" {
				if !c.isNotNull {
						return "NULL"
				}
				return ""
		}
		return c.Default
}

func (c *LuaMigrateColumn) buildCheckString() string {
	if c.Check != "" {
		return fmt.Sprintf(` CHECK (%s)`, c.Check)
	}
	return ""
}

func (c *LuaMigrateColumn) buildCommentString() string {
	if c.Comment != "" {
		return fmt.Sprintf(` COMMENT "%s"`, c.Comment)
	}
	return ""
}

func (c *LuaMigrateColumn) buildFirstString() string {
	if c.isFirst {
		return ` FIRST`
	}
	return ""
}
func (c *LuaMigrateColumn) buildAfterString() string {
	if c.After != "" {
		return fmt.Sprintf(` AFTER %s `, c.opt.quoteColumnName(c.After))
	}
	return ""
}

func (c *LuaMigrateColumn) buildPosString() string {
	if c.isFirst {
		return c.buildFirstString()
	}
	return c.buildAfterString()
}

func (c *LuaMigrateColumn) buildAppendString() string {
	if c.Append != "" {
		if strings.HasSuffix(c.Append, ",") || strings.Contains(c.Append, ";") {
			return c.Append
		}
		return c.Append + ","
	}
	return ","
}

func (c *LuaMigrateColumn) formatReplace(format string, arr []Kv) string {
	if format == "" {
		return ""
	}
	for _, item := range arr {
		var (
			key   = item.K
			value = item.V
		)
		if key != "" && strings.Contains(format, key) {
			format = strings.ReplaceAll(format, key, value)
		}
	}
	return format
}
