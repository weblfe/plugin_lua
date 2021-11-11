package query

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type (
	CommandQueryBuilder struct {
		options map[string]string
		exprMap map[string]string
		typeMap *map[string]string
	}
)

func NewCommandBuilder(opts ...map[string]string) *CommandQueryBuilder {
	var builder = new(CommandQueryBuilder)
	if len(opts) > 0 {
		builder.options = opts[0]
	}
	return builder.init()
}

func (builder *CommandQueryBuilder) init() *CommandQueryBuilder {
	builder.typeMap = &defaultTypeMap
	builder.exprMap = make(map[string]string)
	return builder
}

// AddCheck public function addCheck($name, $table, $expression)
func (builder *CommandQueryBuilder) AddCheck(name, table string, expression fmt.Stringer) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.Write([]byte(builder.quoteTableName(table) + ` ADD CONSTRAINT `))
	buffer.Write([]byte(builder.quoteColumnName(name)))
	buffer.Write([]byte(` CHECK (` + builder.quoteSql(expression.String()) + `)`))
	return buffer.String()
}

// DropCheck public function dropCheck($name, $table)
func (builder *CommandQueryBuilder) DropCheck(name, table string) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.Write([]byte(builder.quoteTableName(table)))
	buffer.Write([]byte(` DROP CONSTRAINT ` + builder.quoteColumnName(name)))
	return buffer.String()
}

// CreateIndex public function createIndex(name, table, columns, unique)
func (builder *CommandQueryBuilder) CreateIndex(name, table string, columns interface{}, unique ...bool) string {
	if len(unique) <= 0 {
		unique = append(unique, false)
	}
	var (
		u      = unique[0]
		buffer = bytes.NewBufferString(``)
	)
	if !u {
		buffer.Write([]byte(`CREATE INDEX `))
	} else {
		buffer.Write([]byte(`CREATE UNIQUE INDEX `))
	}
	buffer.Write([]byte(builder.quoteTableName(name) + ` ON `))
	buffer.Write([]byte(builder.quoteTableName(table)))
	buffer.Write([]byte(fmt.Sprintf(` (%s)`, builder.buildColumns(columns))))
	return buffer.String()
}

// DropIndex public function dropIndex($name, $table)
func (builder *CommandQueryBuilder) DropIndex(name, table string) string {
	var (
		buffer = bytes.NewBufferString(`DROP INDEX `)
	)
	buffer.Write([]byte(builder.quoteTableName(name) + ` ON `))
	buffer.Write([]byte(builder.quoteTableName(table)))
	return buffer.String()
}

// AddForeignKey public function addForeignKey($name, $table, $columns, $refTable, $refColumns, $delete = null, $update = null)
func (builder *CommandQueryBuilder) AddForeignKey(name, table string, columns interface{}, refTable string, refColumns interface{}, extras ...Args) string {
	if len(extras) <= 0 {
		extras = append(extras, Args{})
	}
	var (
		opt    = extras[0]
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.Write([]byte(builder.quoteTableName(table)))
	buffer.Write([]byte(` ADD CONSTRAINT ` + builder.quoteTableName(name)))
	buffer.Write([]byte(` FOREIGN KEY (` + builder.buildColumns(columns)))
	buffer.Write([]byte(` REFERENCES ` + builder.quoteTableName(refTable)))
	buffer.Write([]byte(` (` + builder.buildColumns(refColumns)))
	if opt.NotEmpty(`delete`) {
		buffer.Write([]byte(` ON DELETE ` + opt.Str(`delete`)))
	}
	if opt.NotEmpty(`update`) {
		buffer.Write([]byte(` ON UPDATE ` + opt.Str(`update`)))
	}
	return buffer.String()
}

// DropForeignKey   public function dropForeignKey($name, $table)
func (builder *CommandQueryBuilder) DropForeignKey(name, table string) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.Write([]byte(builder.quoteTableName(table)))
	buffer.Write([]byte(` DROP CONSTRAINT ` + builder.quoteTableName(name)))
	return buffer.String()
}

// AlterColumn public function alterColumn($table, $column, $type)
func (builder *CommandQueryBuilder) AlterColumn(table, column string, T fmt.Stringer) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.Write([]byte(builder.quoteTableName(table) + ` CHANGE `))
	buffer.Write([]byte(builder.quoteColumnName(column) + ` `))
	buffer.Write([]byte(builder.quoteColumnName(column) + ` `))
	buffer.Write([]byte(builder.GetColumnType(T)))
	return buffer.String()
}

// AddCommentOnColumn  public function addCommentOnColumn($table, $column, $comment)
func (builder *CommandQueryBuilder) AddCommentOnColumn(table, column, comment string) string {
	var (
		buffer = bytes.NewBufferString(`COMMENT ON COLUMN `)
	)
	buffer.Write([]byte(builder.quoteTableName(table) + `.`))
	buffer.Write([]byte(builder.quoteColumnName(column) + ` IS `))
	buffer.Write([]byte(builder.quoteValue(comment)))
	return buffer.String()
}

// AddCommentOnTable  public function addCommentOnTable($table, $comment)
func (builder *CommandQueryBuilder) AddCommentOnTable(table, comment string) string {
	var (
		buffer = bytes.NewBufferString(`COMMENT ON TABLE `)
	)
	buffer.Write([]byte(builder.quoteTableName(table) + ` IS `))
	buffer.Write([]byte(builder.quoteValue(comment)))
	return buffer.String()
}

// DropCommentFromTable  public function dropCommentFromTable($table)
func (builder *CommandQueryBuilder) DropCommentFromTable(table string) string {
	var (
		buffer = bytes.NewBufferString(`COMMENT ON TABLE `)
	)
	buffer.Write([]byte(builder.quoteTableName(table) + ` IS NULL`))
	return buffer.String()
}

// DropCommentFromColumn public function dropCommentFromColumn($table, $column)
func (builder *CommandQueryBuilder) DropCommentFromColumn(table, column string) string {
	var (
		buffer = bytes.NewBufferString(`COMMENT ON COLUMN `)
	)
	buffer.Write([]byte(builder.quoteTableName(table) + `.`))
	buffer.Write([]byte(builder.quoteColumnName(column) + ` IS NULL`))
	return buffer.String()
}

func (builder *CommandQueryBuilder) quoteTableName(name string) string {
	return builder.quoteFor(name, "table_prefix", "table_suffix", "tableMap")
}

func (builder *CommandQueryBuilder) quoteColumnName(name string) string {
	return builder.quoteFor(name, "column_prefix", "column_suffix", "columnMap")
}

func (builder *CommandQueryBuilder) quoteValue(value string) string {
	return builder.quoteFor(value, "comment_prefix", "comment_suffix", "")
}

func (builder *CommandQueryBuilder) quoteFor(value, prefixKey, suffixKey, mapGroup string) string {
	if value == "" {
		return ""
	}
	var (
		prefix = builder.getOptOr(prefixKey)
		suffix = builder.getOptOr(suffixKey)
	)
	// 前后缀 替换
	if strings.HasPrefix(value, quotePrefix) && strings.HasSuffix(value, quoteSuffix) {
		value = strings.TrimSuffix(strings.TrimPrefix(value, quotePrefix), quoteSuffix)
		if strings.HasPrefix(value, quoteReplacer) {
			value = strings.Replace(value, quoteReplacer, prefix, 1)
		}
		if strings.HasSuffix(value, quoteReplacer) {
			value = strings.Replace(value, quoteReplacer, suffix, 1)
		}
	}
	value = builder.quoteSql(value)
	// 键名-隐射
	if mapGroup != "" {
		var mapKey = fmt.Sprintf(`%s.%s`, mapGroup, value)
		if v := builder.getOptOr(mapKey); v != "" {
			return v
		}
	}
	return fmt.Sprintf("`%s`", value)
}

// GetColumnType 获取类型
func (builder *CommandQueryBuilder) GetColumnType(T fmt.Stringer) string {
	var ty = T.String()
	if ty == "" {
		return ""
	}
	var typesMap = builder.GetTypesMap()
	if t, ok := typesMap[ty]; ok {
		return t
	}
	if matches := RegexpMatches(typeFullPattern, ty); len(matches) >= 4 {
		var (
			v      = matches[1]
			tv, ok = typesMap[v]
		)
		if !ok {
			var replacement = fmt.Sprintf(`(%s)`, matches[2])
			return RegexpReplace(typeLengthReplacePattern, replacement, tv) + matches[3]
		}
		return ty
	}
	if matches := RegexpMatches(typeSpacePattern, ty); len(matches) >= 1 {
		var (
			v      = matches[1]
			tv, ok = typesMap[v]
		)
		if !ok {
			return RegexpReplace(typeReplacePattern, tv, ty)
		}
	}
	return ty
}

func (builder *CommandQueryBuilder) getOptOr(key string, v ...string) string {
	if builder.options == nil || key == "" {
		return ""
	}
	v = append(v, "")
	if v, ok := builder.options[key]; ok {
		return v
	}
	return ""
}

func (builder *CommandQueryBuilder) quoteSql(sql string) string {
	return quoteRegexp.ReplaceAllStringFunc(template.HTMLEscapeString(sql), func(s string) string {
		if s == " " {
			return ""
		}
		return fmt.Sprintf(`\%s`, s)
	})
}

func (builder *CommandQueryBuilder) buildColumns(columns interface{}) string {
	if columns == nil {
		return ""
	}
	return builder.quoteAny(columns, builder.quoteColumnName)
}

func (builder *CommandQueryBuilder) GetTypesMap() map[string]string {
	if builder.typeMap != nil {
		return *builder.typeMap
	}
	return defaultTypeMap
}

func (builder *CommandQueryBuilder) quoteAny(value interface{}, quoteTransform func(string) string) string {
	switch value.(type) {
	case string:
		var str = value.(string)
		// array
		if strings.Contains(str, ",") && !strings.Contains(str, "(") {
			arr := strings.Split(str, ",")
			for i, v := range arr {
				arr[i] = quoteTransform(v)
			}
			return strings.Join(arr, ",")
		}
		return str
	case []string:
		var (
			all []string
			arr = value.([]string)
		)
		for _, v := range arr {
			if v == "" {
				continue
			}
			all = append(all, quoteTransform(v))
		}
		if len(all) > 0 {
			return strings.Join(all, ",")
		}
	case []interface{}:
		var (
			all []string
			arr = value.([]interface{})
		)
		for _, v := range arr {
			if v == "" {
				continue
			}
			var s = ""
			switch v.(type) {
			case string:
				s = v.(string)
			case fmt.Stringer:
				s = v.(fmt.Stringer).String()
			case fmt.GoStringer:
				s = v.(fmt.GoStringer).GoString()
			}
			if s != "" {
				all = append(all, quoteTransform(s))
			}
		}
		if len(all) > 0 {
			return strings.Join(all, ",")
		}
	case []fmt.Stringer:
		var (
			all []string
			arr = value.([]fmt.Stringer)
		)
		for _, v := range arr {
			if v == nil {
				continue
			}
			var s = v.(fmt.Stringer).String()
			if s != "" {
				all = append(all, quoteTransform(s))
			}
		}
		if len(all) > 0 {
			return strings.Join(all, ",")
		}
	case []fmt.GoStringer:
		var (
			all []string
			arr = value.([]fmt.GoStringer)
		)
		for _, v := range arr {
			if v == nil {
				continue
			}
			var s = v.(fmt.GoStringer).GoString()
			if s != "" {
				all = append(all, quoteTransform(s))
			}
		}
		if len(all) > 0 {
			return strings.Join(all, ",")
		}
	}
	return ""
}
