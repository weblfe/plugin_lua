package query

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unicode"
)

type (
	StringerAny struct {
		row  reflect.Value
		body string
	}
)

func NewStringerAny(v reflect.Value) *StringerAny {
	var stringer = new(StringerAny)
	stringer.row = v
	return stringer
}

func NewStringerT(v interface{}) *StringerAny {
		var stringer = new(StringerAny)
		stringer.row = reflect.ValueOf(v)
		return stringer
}

func (any *StringerAny) String() string {
	if any.body != "" {
		return any.body
	}
	var (
		v      = any.row
		kind   = v.Kind()
		refVal = any.row.Interface()
	)
	// 接口类型
	switch refVal.(type) {
	case string:
		any.body = refVal.(string)
	case *string:
		any.body = *refVal.(*string)
	case fmt.Stringer:
		any.body = refVal.(fmt.Stringer).String()
	case fmt.GoStringer:
		any.body = refVal.(fmt.GoStringer).GoString()
	}
	if any.body != "" {
		return any.body
	}
	if kind == reflect.Ptr {
		v = any.row.Elem()
		kind = v.Kind()
	}
	// 基础类型
	switch kind {
	case reflect.String:
		any.body = v.String()
	case reflect.Complex64, reflect.Complex128:
		any.body = fmt.Sprintf("%v", v.Interface())
	case reflect.Chan:
		any.body = fmt.Sprintf("Chan<%p>", v.Interface())
	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		if bytes, err := json.Marshal(v.Interface()); err == nil {
			any.body = string(bytes)
		}
	case reflect.Bool:
		any.body = fmt.Sprintf("%v", v.Bool())
	case reflect.Func:
		any.body = fmt.Sprintf("func<%p>", v.Interface())
	case reflect.Uintptr, reflect.UnsafePointer:
		any.body = fmt.Sprintf("%p", v.Interface())
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Float32, reflect.Float64:
		any.body = fmt.Sprintf("%v", v.Interface())
	}
	return any.body
}

func (any *StringerAny) IsDigit() bool {
	var str = any.String()
	if str == "" {
		return false
	}
	for _, v := range []rune(str) {
		if !unicode.IsDigit(v) {
				return false
		}
	}
	return true
}

func (any *StringerAny) IsNumber() bool {
		var str = any.String()
		if str == "" {
				return false
		}
		for _, v := range []rune(str) {
				if !unicode.IsNumber(v) {
						return false
				}
		}
		return true
}