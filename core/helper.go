package core

import (
	"encoding/json"
	"fmt"
		lua "github.com/yuin/gopher-lua"
		"net/url"
	"os"
	"strings"
)

func GetEnvOr(key string, or ...string) string {
	var v = os.Getenv(key)
	or = append(or, "")
	if len(or) > 0 && v == "" {
		return or[0]
	}
	return v
}

func SprintfEnv(format string, args ...interface{}) string {
	var key = format
	if len(args) > 0 {
		key = fmt.Sprintf(format, args...)
	}
	return strings.ToUpper(key)
}

func GetEnvJsonKvOr(key string, or ...map[string]string) map[string]string {
	or = append(or, map[string]string{})
	var data = GetEnvOr(key)
	if data != "" {
		var m = map[string]string{}
		// 非json 格式
		if !json.Valid([]byte(data)) {
			if !strings.Contains(data, "=") {
				return or[0]
			}
			// query params
			values, err := url.ParseQuery(data)
			if err != nil {
				return or[0]
			}
			for k, v := range values {
				if len(v) <= 1 {
					m[k] = v[0]
				} else {
					m[k] = strings.Join(v, ",")
				}
			}
			return m
		}
		// json 格式
		if err := json.Unmarshal([]byte(data), &m); err == nil {
			return m
		}
	}
	return or[0]
}

func CreateArrByTable(t *lua.LTable) []string {
		if t == nil {
				return nil
		}
		var arr []string
		// 取出 arr
		t.ForEach(func(key lua.LValue, v lua.LValue) {
				if key == nil || v == nil {
						return
				}
				if key.Type() != lua.LTNumber {
						return
				}
				arr = append(arr, v.String())
		})
		return arr
}