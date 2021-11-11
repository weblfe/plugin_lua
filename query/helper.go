package query

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	String string
	Args   map[string]interface{}
)

func (arg *Args) NotNil(key string) bool {
	if arg == nil || key == "" {
		return false
	}
	if v, ok := (*arg)[key]; ok && v != nil {
		return true
	}
	return false
}

func (arg *Args) NotEmpty(key string) bool {
	if arg == nil || key == "" {
		return false
	}
	if v, ok := (*arg)[key]; ok && v != nil && v != "" {
		return true
	}
	return false
}

func (arg *Args) Value(key string) (interface{}, bool) {
	if arg == nil || key == "" {
		return nil, false
	}
	if v, ok := (*arg)[key]; ok && v != nil {
		return v, true
	}
	return nil, false
}

func (arg *Args) Str(key string) string {
	if v, ok := arg.Value(key); ok {
		switch v.(type) {
		case string:
			return v.(string)
		case fmt.Stringer:
			return v.(fmt.Stringer).String()
		case fmt.GoStringer:
			return v.(fmt.GoStringer).GoString()
		case int, uint, uint32, uint8, int16, float64, float32, int32, int64, int8:
			return fmt.Sprintf("%v", v)
		default:
			return ""
		}
	}
	return ""
}

func (arg *Args) Len() int {
	if arg == nil {
		return 0
	}
	return len(*arg)
}

func (s String) String() string {
	return string(s)
}

func (s String) GoString() string {
		return string(s)
}

func NewString(s string) String {
	return String(s)
}

// RegexpReplace 正则替换
func RegexpReplace(pattern string, replacement string, subject string) string {
	var reg, err = regexp.Compile(regexpPattern(pattern))
	if err != nil {
		return subject
	}
	return string(reg.ReplaceAll([]byte(subject), []byte(replacement)))
}

// RegexpMatches 正则获取匹配子字符串
func RegexpMatches(pattern string, subject string) []string {
	var reg, err = regexp.Compile(regexpPattern(pattern))
	if err != nil {
		return nil
	}
	if sub := reg.FindAllStringSubmatch(subject, -1); len(sub) > 0 {
		var matches []string
		for _, v := range sub {
			matches = append(matches, v...)
		}
		return matches
	}
	return nil
}

// 过滤PHP JS 正则开头(/) 与结尾(/)
func regexpPattern(pattern string) string {
	return strings.TrimSuffix(strings.TrimPrefix(pattern, `/`), `/`)
}
