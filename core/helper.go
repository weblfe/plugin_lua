package core

import (
	"fmt"
	"os"
	"strings"
)

func GetEnvOr(key string, or ...string) string {
	var v = os.Getenv(key)
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
