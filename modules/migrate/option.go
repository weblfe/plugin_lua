package migrate

import "strings"

type (
	OptionKv struct {
		Source     string
		ConnUrl    string
		Prefix     string
		Properties map[string]string
	}
)

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
