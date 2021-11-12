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

//BatchInsert public function batchInsert($table, $columns, $rows, &$params = [])
//    {
//        if (empty($rows)) {
//            return '';
//        }
//
//        $schema = $this->db->getSchema();
//        if (($tableSchema = $schema->getTableSchema($table)) !== null) {
//            $columnSchemas = $tableSchema->columns;
//        } else {
//            $columnSchemas = [];
//        }
//
//        $values = [];
//        foreach ($rows as $row) {
//            $vs = [];
//            foreach ($row as $i => $value) {
//                if (isset($columns[$i], $columnSchemas[$columns[$i]])) {
//                    $value = $columnSchemas[$columns[$i]]->dbTypecast($value);
//                }
//                if (is_string($value)) {
//                    $value = $schema->quoteValue($value);
//                } elseif (is_float($value)) {
//                    // ensure type cast always has . as decimal separator in all locales
//                    $value = StringHelper::floatToString($value);
//                } elseif ($value === false) {
//                    $value = 0;
//                } elseif ($value === null) {
//                    $value = 'NULL';
//                } elseif ($value instanceof ExpressionInterface) {
//                    $value = $this->buildExpression($value, $params);
//                }
//                $vs[] = $value;
//            }
//            $values[] = '(' . implode(', ', $vs) . ')';
//        }
//        if (empty($values)) {
//            return '';
//        }
//
//        foreach ($columns as $i => $name) {
//            $columns[$i] = $schema->quoteColumnName($name);
//        }
//
//        return 'INSERT INTO ' . $schema->quoteTableName($table)
//            . ' (' . implode(', ', $columns) . ') VALUES ' . implode(', ', $values);
//    }
func (builder *CommandQueryBuilder) BatchInsert(table string, columns, rows ArrayAble, params ...ArrayAble) string {
		if rows.Empty() {
				return ""
		}
		var (
				buffer = bytes.NewBufferString(`INSERT INTO `)
		)
		buffer.WriteString(builder.quoteTableName(table))
		buffer.WriteString(` (`+strings.Join(columns.Array(),",")+`) VALUES `)
		buffer.WriteString(builder.quoteAny(columns, builder.quoteColumnName) + `)`)
		return buffer.String()
}

// AddUnique public function addUnique($name, $table, $columns)
func (builder *CommandQueryBuilder) AddUnique(name, table string, columns ArrayAble) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table) + ` ADD CONSTRAINT `)
	buffer.WriteString(builder.quoteColumnName(name) + ` UNIQUE (`)
	buffer.WriteString(builder.quoteAny(columns, builder.quoteColumnName) + `)`)
	return buffer.String()
}

// DropUnique public function dropUnique($name, $table)
func (builder *CommandQueryBuilder) DropUnique(name, table string) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table))
	buffer.WriteString(` DROP CONSTRAINT ` + builder.quoteColumnName(name))
	return buffer.String()
}

// AddCheck public function addCheck($name, $table, $expression)
func (builder *CommandQueryBuilder) AddCheck(name, table string, expression fmt.Stringer) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table) + ` ADD CONSTRAINT `)
	buffer.WriteString(builder.quoteColumnName(name))
	buffer.WriteString(` CHECK (` + builder.quoteSql(expression.String()) + `)`)
	return buffer.String()
}

// DropCheck public function dropCheck($name, $table)
func (builder *CommandQueryBuilder) DropCheck(name, table string) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table))
	buffer.WriteString(` DROP CONSTRAINT ` + builder.quoteColumnName(name))
	return buffer.String()
}

// CreateIndex public function createIndex(name, table, columns, unique)
func (builder *CommandQueryBuilder) CreateIndex(name, table string, columns ArrayAble, unique ...bool) string {
	if len(unique) <= 0 {
		unique = append(unique, false)
	}
	var (
		u      = unique[0]
		buffer = bytes.NewBufferString(``)
	)
	if !u {
		buffer.WriteString(`CREATE INDEX `)
	} else {
		buffer.WriteString(`CREATE UNIQUE INDEX `)
	}
	buffer.WriteString(builder.quoteTableName(name) + ` ON `)
	buffer.WriteString(builder.quoteTableName(table))
	buffer.WriteString(fmt.Sprintf(` (%s)`, builder.buildColumns(columns)))
	return buffer.String()
}

// DropIndex public function dropIndex($name, $table)
func (builder *CommandQueryBuilder) DropIndex(name, table string) string {
	var (
		buffer = bytes.NewBufferString(`DROP INDEX `)
	)
	buffer.WriteString(builder.quoteTableName(name) + ` ON `)
	buffer.WriteString(builder.quoteTableName(table))
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
	buffer.WriteString(builder.quoteTableName(table))
	buffer.WriteString(` ADD CONSTRAINT ` + builder.quoteTableName(name))
	buffer.WriteString(` FOREIGN KEY (` + builder.buildColumns(columns))
	buffer.WriteString(` REFERENCES ` + builder.quoteTableName(refTable))
	buffer.WriteString(` (` + builder.buildColumns(refColumns))
	if opt.NotEmpty(`delete`) {
		buffer.WriteString(` ON DELETE ` + opt.Str(`delete`))
	}
	if opt.NotEmpty(`update`) {
		buffer.WriteString(` ON UPDATE ` + opt.Str(`update`))
	}
	return buffer.String()
}

// DropForeignKey   public function dropForeignKey($name, $table)
func (builder *CommandQueryBuilder) DropForeignKey(name, table string) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table))
	buffer.WriteString(` DROP CONSTRAINT ` + builder.quoteTableName(name))
	return buffer.String()
}

//Delete public function delete($table, $condition, &$params)
func (builder *CommandQueryBuilder) Delete(table string, condition ArrayAble, params ArrayAble) string {
	var (
		buffer = bytes.NewBufferString(`DELETE FROM `)
		where  = builder.buildWhere(condition, params)
	)
	buffer.WriteString(builder.quoteTableName(table))
	if where != "" {
		buffer.WriteString(` ` + where)
	}
	return buffer.String()
}

//CreateTable public function createTable($table, $columns, $options = null)
func (builder *CommandQueryBuilder) CreateTable(table string, columns ArrayAble, options ...string) string {
	var (
		cols   []string
		buffer = bytes.NewBufferString(`CREATE TABLE `)
	)
	for _, kv := range columns.Kvs() {
		if !kv.IsString() {
			cols = append(cols, kv.String())
		} else {
			cols = append(cols, builder.quoteColumnName(kv.Key())+` `+builder.GetColumnType(kv.Value()))
		}
	}
	buffer.WriteString(builder.quoteTableName(table) + ` (\n`)
	buffer.WriteString(strings.Join(cols, `,\n`) + `\n)`)
	if len(options) > 0 && options[0] != "" {
		buffer.WriteString(` ` + options[0])
	}
	return buffer.String()
}

//RenameTable public function renameTable($oldName, $newName)
func (builder *CommandQueryBuilder) RenameTable(table, newName string) string {
	var (
		buffer = bytes.NewBufferString(`RENAME TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table))
	buffer.WriteString(` TO ` + builder.quoteTableName(newName))
	return buffer.String()
}

// DropTable  public function dropTable($table)
func (builder *CommandQueryBuilder) DropTable(table string) string {
	var (
		buffer = bytes.NewBufferString(`DROP TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table))
	return buffer.String()
}

// AddPrimaryKey public function addPrimaryKey($name, $table, $columns)
func (builder *CommandQueryBuilder) AddPrimaryKey(name, table string, columns ArrayAble) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table) + ` ADD CONSTRAINT `)
	buffer.WriteString(builder.quoteColumnName(name) + ` PRIMARY KEY (`)
	buffer.WriteString(builder.quoteAny(columns, builder.quoteColumnName) + `)`)
	return buffer.String()
}

// DropPrimaryKey public function DropPrimaryKey($name, $table)
func (builder *CommandQueryBuilder) DropPrimaryKey(name, table string) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table))
	buffer.WriteString(` DROP CONSTRAINT ` + builder.quoteColumnName(name))
	return buffer.String()
}

//TruncateTable public function truncateTable($table)
func (builder *CommandQueryBuilder) TruncateTable(table string) string {
	var (
		buffer = bytes.NewBufferString(`TRUNCATE TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table))
	return buffer.String()
}

// AddColumn public function addColumn($table, $column, $type)
func (builder *CommandQueryBuilder) AddColumn(table, column string, typeName fmt.Stringer) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table))
	buffer.WriteString(` ADD ` + builder.quoteColumnName(column))
	buffer.WriteString(builder.GetColumnType(typeName))
	return buffer.String()
}

//DropColumn public function dropColumn($table, $column)
func (builder *CommandQueryBuilder) DropColumn(table, column string) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table))
	buffer.WriteString(` DROP COLUMN ` + builder.quoteColumnName(column))
	return buffer.String()
}

// RenameColumn public function renameColumn($table, $oldName, $newName)
func (builder *CommandQueryBuilder) RenameColumn(table, oldName, newName string) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table))
	buffer.WriteString(` RENAME COLUMN ` + builder.quoteColumnName(oldName))
	buffer.WriteString(` TO ` + builder.quoteColumnName(newName))
	return buffer.String()
}

// AlterColumn public function alterColumn($table, $column, $type)
func (builder *CommandQueryBuilder) AlterColumn(table, column string, typeName fmt.Stringer) string {
	var (
		buffer = bytes.NewBufferString(`ALTER TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table) + ` CHANGE `)
	buffer.WriteString(builder.quoteColumnName(column) + ` `)
	buffer.WriteString(builder.quoteColumnName(column) + ` `)
	buffer.WriteString(builder.GetColumnType(typeName))
	return buffer.String()
}

// AddCommentOnColumn  public function addCommentOnColumn($table, $column, $comment)
func (builder *CommandQueryBuilder) AddCommentOnColumn(table, column, comment string) string {
	var (
		buffer = bytes.NewBufferString(`COMMENT ON COLUMN `)
	)
	buffer.WriteString(builder.quoteTableName(table) + `.`)
	buffer.WriteString(builder.quoteColumnName(column) + ` IS `)
	buffer.WriteString(builder.quoteValue(comment))
	return buffer.String()
}

// AddCommentOnTable  public function addCommentOnTable($table, $comment)
func (builder *CommandQueryBuilder) AddCommentOnTable(table, comment string) string {
	var (
		buffer = bytes.NewBufferString(`COMMENT ON TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table) + ` IS `)
	buffer.WriteString(builder.quoteValue(comment))
	return buffer.String()
}

// DropCommentFromTable  public function dropCommentFromTable($table)
func (builder *CommandQueryBuilder) DropCommentFromTable(table string) string {
	var (
		buffer = bytes.NewBufferString(`COMMENT ON TABLE `)
	)
	buffer.WriteString(builder.quoteTableName(table) + ` IS NULL`)
	return buffer.String()
}

// DropCommentFromColumn public function dropCommentFromColumn($table, $column)
func (builder *CommandQueryBuilder) DropCommentFromColumn(table, column string) string {
	var (
		buffer = bytes.NewBufferString(`COMMENT ON COLUMN `)
	)
	buffer.WriteString(builder.quoteTableName(table) + `.`)
	buffer.WriteString(builder.quoteColumnName(column) + ` IS NULL`)
	return buffer.String()
}

func (builder *CommandQueryBuilder) buildWhere(condition, params ArrayAble) string {
	var where = builder.buildCondition(condition, params)
	if where == "" {
		return ``
	}
	return `WHERE ` + where
}

func (builder *CommandQueryBuilder) buildCondition(condition, params ArrayAble) string {
	if condition == nil {
		return ``
	}
	var (
		arr  = condition.Array()
		cond = condition.String()
	)
	if cond == "" && len(arr) == 0 {
		return ""
	}
	if cond != "" && len(arr) == 1 {
		return cond
	}
	var conditionExpress = builder.createConditionFromArray(condition)
	switch conditionExpress.(type) {
	case string:
		return conditionExpress.(string)
	case ExpressionInterface:
		var express = conditionExpress.(ExpressionInterface)
		return builder.buildExpression(express, params)
	case fmt.Stringer:
		return conditionExpress.(fmt.Stringer).String()
	case fmt.GoStringer:
		return conditionExpress.(fmt.GoStringer).GoString()
	}
	return ``
}

func (builder *CommandQueryBuilder) createConditionFromArray(condition ArrayAble) interface{} {
  //@todo
	return nil
}

func (builder *CommandQueryBuilder) getExpressionBuilder(express ExpressionInterface) ExpressionBuilderInterface {
	if express == nil {
		return nil
	}
	return GetExpressBuilder(express.GetClass())
}

func (builder *CommandQueryBuilder) buildExpression(express ExpressionInterface, param ArrayAble) string {
	if express == nil {
		return ``
	}
	if param == nil {
		return express.String()
	}
	var b = builder.getExpressionBuilder(express)
	if b != nil {
		return b.Build(express, param)
	}
	return express.String()
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
	case ArrayAble:
		var (
			all []string
			arr = value.(ArrayAble)
		)
		for _, v := range arr.Array() {
			if v == "" {
				continue
			}
			all = append(all, quoteTransform(v))
		}
		if len(all) > 0 {
			return strings.Join(all, ",")
		}
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
