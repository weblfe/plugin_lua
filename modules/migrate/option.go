package migrate

import (
	"github.com/weblfe/plugin_lua/core"
	"strings"
)

type (
	OptionKv struct {
		Source     string
		ConnUrl    string
		Prefix     string
		Suffix     string
		Properties map[string]string
	}
)

func WithEnvOptions(env string) *OptionKv {
	return &OptionKv{
		Source:     core.GetEnvOr(core.SprintfEnv("%s_source", env)),
		ConnUrl:    core.GetEnvOr(core.SprintfEnv("%s_db_conn", env)),
		Prefix:     core.GetEnvOr(core.SprintfEnv("%s_table_prefix", env)),
		Suffix:     core.GetEnvOr(core.SprintfEnv("%s_table_suffix", env)),
		Properties: core.GetEnvJsonKvOr(core.SprintfEnv("%s_extras", env)),
	}
}

func (o *OptionKv) quoteColumnName(column string) string {
	if strings.HasPrefix(column, "{{") && strings.HasSuffix(column, "}}") {
		var (
			prefix = o.getPropertyOr("column_prefix")
			suffix = o.getPropertyOr("column_suffix")
		)
		column = strings.TrimSuffix(strings.TrimPrefix(column, "{{"), "}}")
		if strings.HasPrefix(column, "%") && prefix != "" {
			column = strings.Replace(column, "%", prefix, 1)
		}
		if strings.HasSuffix(column, "%") && suffix != "" {
			column = strings.Replace(column, "%", suffix, 1)
		}
		return column
	}
	return column
}

func (o *OptionKv) getPropertyOr(key string, or ...string) string {
	or = append(or, "")
	if o.Properties != nil {
		if v, ok := o.Properties[key]; ok {
			return v
		}
	}
	return or[0]
}

func (o *OptionKv) quoteTableName(table string) string {
	if strings.HasPrefix(table, "{{") && strings.HasSuffix(table, "}}") {
		var (
			prefix = o.Prefix
			suffix = o.Suffix
		)
		table = strings.TrimSuffix(strings.TrimPrefix(table, "{{"), "}}")
		if strings.HasPrefix(table, "%") && prefix != "" {
			table = strings.Replace(table, "%", prefix, 1)
		}
		if strings.HasSuffix(table, "%") && suffix != "" {
			table = strings.Replace(table, "%", suffix, 1)
		}
		return table
	}
	return table
}
