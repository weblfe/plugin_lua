package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

type (
	String string
	Array  []string
	Args   map[string]interface{}

	ArrayAble interface {
		fmt.Stringer
		Kvs() []KV
		Array() []string
		Empty() bool
	}

	KV interface {
		Key() string
		Value() fmt.Stringer
		IsString() bool
		fmt.Stringer
	}

	arrayAbleImpl struct {
		body   interface{}
		pairs  []KV
		array  []string
		isStr  bool
		err    error
		parsed *sync.Once
	}

	KvPairs struct {
		K string
		V fmt.Stringer
	}
)

const (
	RegexpSplitNoEmpty       = 0
	RegexpSplitDelimCapture  = 1
	RegexpSplitOffsetCapture = 2
	RegexpQuoteChars         = `\.\\\+\*\?\[\^\]\$\(\)\{\}=\!<>\|:-`
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
func RegexpMatches(pattern string, subject string) Array {
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

// RegexpSplit 正则分隔
func RegexpSplit(pattern, subject string, limit, flag int) Array {
	if limit == 0 {
		limit = -1
	}
	pattern = regexpPattern(pattern)
	var reg, err = regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	var (
		list Array
		arr  = reg.Split(subject, limit)
	)
	switch flag {
	case RegexpSplitNoEmpty:
		for _, v := range arr {
			if v == "" {
				continue
			}
			list = append(list, v)
		}
		if len(list) <= 0 {
			return nil
		}
		return list
	case RegexpSplitDelimCapture:
		return arr
	case RegexpSplitOffsetCapture:
		return arr
	}
	return nil
}

func RegexpQuote(word string, delimiter ...string) string {
	var pattern = RegexpQuoteChars
	if len(delimiter) > 0 {
		pattern = pattern + delimiter[0]
	}
	pattern = fmt.Sprintf(`[%s]`, pattern)
	var reg, err = regexp.Compile(regexpPattern(pattern))
	if err != nil {
		return word
	}
	return string(reg.ReplaceAllFunc([]byte(word), func(bytes []byte) []byte {
		return append([]byte(`\`), bytes...)
	}))
}

func RegexpStrAllQuote(d string) string {
	var (
		arr  = []rune(d)
		list []string
	)
	for _, v := range arr {
		list = append(list, `\`+string(v))
	}
	return strings.Join(list, "")
}

// 过滤PHP JS 正则开头(/) 与结尾(/)
func regexpPattern(pattern string) string {
	return strings.TrimSuffix(strings.TrimPrefix(pattern, `/`), `/`)
}

func (arr *Array) Add(v string) Array {
	*arr = append(*arr, v)
	return *arr
}

func (arr *Array) Len() int {
	return len(*arr)
}

func (arr *Array) Pop() string {
	if arr.Len() > 0 {
		var v = (*arr)[0]
		*arr = (*arr)[1:]
		return v
	}
	return ""
}

func (arr *Array) Strings() []string {
	return *arr
}

func (arr *Array) Include(v string) bool {
	for _, value := range *arr {
		if value == v {
			return true
		}
	}
	return false
}

func (arr *Array) Index(v string) int {
	for i, value := range *arr {
		if value == v {
			return i
		}
	}
	return -1
}

func (arr *Array) Empty() bool {
	if arr == nil || len(*arr) <= 0 {
		return true
	}
	return false
}

func (arr *Array) String() string {
	return `[` + strings.Join(*arr, ",") + `]`
}

func NewArrayAble(v interface{}) *arrayAbleImpl {
	var arr = new(arrayAbleImpl)
	arr.body = v
	return arr
}

func (arr *arrayAbleImpl) parse() error {
	if arr.body == nil {
		arr.array = []string{}
		return nil
	}
	switch arr.body.(type) {
	case string:
		var str = arr.body.(string)
		arr.isStr = true
		arr.array, arr.pairs = arr.strSplit(str)
		return nil
	case *string:
		var str = arr.body.(*string)
		arr.isStr = true
		arr.array, arr.pairs = arr.strSplit(*str)
		return nil
	case fmt.Stringer:
		arr.array, arr.pairs = arr.strSplit(arr.body.(fmt.Stringer).String())
		return nil
	case fmt.GoStringer:
		arr.array, arr.pairs = arr.strSplit(arr.body.(fmt.GoStringer).GoString())
		return nil
	case []string:
		var strArr = arr.body.([]string)
		arr.array = strArr
		return nil
	case []fmt.Stringer:
		var items = arr.body.([]fmt.Stringer)
		for _, v := range items {
			arr.array = append(arr.array, v.String())
		}
		return nil
	case []KV:
		arr.pairs = arr.body.([]KV)
	case []fmt.GoStringer:
		var items = arr.body.([]fmt.GoStringer)
		for _, v := range items {
			arr.array = append(arr.array, v.GoString())
		}
		return nil
	}
	var (
		v    = reflect.ValueOf(arr.body)
		kind = v.Kind()
	)
	if kind == reflect.Ptr {
		v = v.Elem()
		kind = v.Kind()
	}
	// array
	if kind == reflect.Array || kind == reflect.Slice {
		var size = v.Len()
		for i := 0; i < size; i++ {
			var (
				it = v.Index(i)
				k  = it.Kind()
			)
			if IsNumber(k) {
				arr.array = append(arr.array, fmt.Sprintf("%v", it.Interface()))
				continue
			}
			if IsString(k) {
				arr.array = append(arr.array, it.String())
				continue
			}
			if k == reflect.Map {
				for _, k := range v.MapKeys() {
					var value = v.MapIndex(k)
					arr.pairs = append(arr.pairs, NewKvPairs(NewStringerAny(k), NewStringerAny(value)))
				}
				continue
			}
			var d = it.Interface()
			switch d.(type) {
			case fmt.Stringer:
				arr.array = append(arr.array, d.(fmt.Stringer).String())
			case fmt.GoStringer:
				arr.array = append(arr.array, d.(fmt.GoStringer).GoString())
			}
		}
		return nil
	}
	if kind == reflect.Map {
		for _, k := range v.MapKeys() {
			var value = v.MapIndex(k)
			arr.pairs = append(arr.pairs, NewKvPairs(NewStringerAny(k), NewStringerAny(value)))
		}
		return nil
	}
	arr.array = []string{}
	return nil
}

func (arr *arrayAbleImpl) Empty() bool {
	if arr == nil || arr.body == "" || arr.body == nil {
		return true
	}
	if err := arr.done(); err != nil {
		return false
	}
	if len(arr.array) == 0 && len(arr.pairs) == 0 {
		return true
	}
	return false
}

func (arr *arrayAbleImpl) done() error {
	if arr.parsed == nil {
		arr.parsed = &sync.Once{}
	}
	arr.parsed.Do(func() {
		if err := arr.parse(); err != nil {
			arr.err = err
		}
	})
	return arr.err
}

func (arr *arrayAbleImpl) IsPair() bool {
	if arr == nil || arr.body == nil {
		return false
	}
	if err := arr.done(); err != nil {
		return false
	}
	if len(arr.array) <= 0 && len(arr.pairs) > 0 {
		return true
	}
	return false
}

func (arr *arrayAbleImpl) IsArray() bool {
	if arr == nil || arr.body == nil {
		return false
	}
	if err := arr.done(); err != nil {
		return false
	}
	if len(arr.array) > 0 && len(arr.pairs) <= 0 {
		return true
	}
	return false
}

func (arr *arrayAbleImpl) IsMix() bool {
	if err := arr.done(); err != nil {
		return false
	}
	if len(arr.array) > 0 && len(arr.pairs) > 0 {
		return true
	}
	return false
}

func (arr *arrayAbleImpl) strSplit(str string) ([]string, []KV) {
	if json.Valid([]byte(str)) {
		var isPairs = strings.HasPrefix(str, `[`) && strings.HasSuffix(str, `]`) &&
			strings.Contains(str, `{`) && strings.Contains(str, `}`)
		if isPairs {
			var (
				arrPairs []KV
				tmpMap   = new([]map[string]interface{})
			)
			if err := json.Unmarshal([]byte(str), tmpMap); err == nil {
				for _, it := range *tmpMap {
					for k, v := range it {
						arrPairs = append(arrPairs, NewKvPairs(NewString(k), NewStringerT(v)))
					}
				}
				return nil, arrPairs
			}
		}
	}
	// 剔除 []
	if strings.HasPrefix(str, `[`) && strings.HasSuffix(str, `]`) {
		str = strings.TrimSuffix(strings.TrimPrefix(str, `[`), `]`)
	}
	if items := RegexpSplit(`/\s*,\s*/`, str, -1, RegexpSplitNoEmpty); items.Len() > 0 {
		return items, nil
	}
	if strings.HasPrefix(str, "(") {
		if strings.HasSuffix(str, ")") && strings.Contains(str, `,`) {
			str = strings.TrimSuffix(strings.TrimPrefix(str, `(`), `)`)
			return strings.Split(str, ","), nil
		}
		return []string{str}, nil
	}
	return strings.Split(str, ","), nil
}

func (arr *arrayAbleImpl) String() string {
	if err := arr.done(); err != nil {
		return ""
	}
	if arr.pairs != nil {
		var (
			size   = len(arr.pairs)
			buffer = bytes.NewBufferString(`[`)
		)
		for i, v := range arr.pairs {
			var (
				k     = v.Key()
				value = v.Value().String()
			)
			buffer.WriteString(fmt.Sprintf(`{"%s":"%s"}`, k, value))
			if i < size {
				buffer.WriteString(`,`)
			}
		}
		buffer.WriteString(`]`)
		return buffer.String()
	}
	if arr.array == nil {
		return `[]`
	}
	// 单字符串
	if arr.isStr && len(arr.array) == 1 {
		return arr.array[0]
	}
	return `[` + strings.Join(arr.array, ",") + `]`
}

func (arr *arrayAbleImpl) Array() []string {
	if err := arr.done(); err != nil {
		return nil
	}
	if arr.pairs != nil {
		var items []string
		for _, v := range arr.pairs {
			items = append(items, v.String())
		}
		return items
	}
	return arr.array
}

func (arr *arrayAbleImpl) Kvs() []KV {
	if err := arr.done(); err != nil {
		return nil
	}
	if arr.pairs != nil {
		return arr.pairs
	}
	var kvs []KV
	for i, v := range arr.array {
		var (
			value = NewString(v)
			k     = NewString(fmt.Sprintf(`%d`, i))
		)
		kvs = append(kvs, NewKvPairs(k, value))
	}
	return kvs
}

func NewKvPairs(k, v fmt.Stringer) *KvPairs {
	var pairs = new(KvPairs)
	pairs.V = v
	pairs.K = k.String()
	return pairs
}

func KvPairsWithMap(m map[string]string) *KvPairs {
	var pairs = new(KvPairs)
	if len(m) == 1 {
		for k, v := range m {
			pairs.K = k
			pairs.V = NewString(v)
		}
	}
	return pairs
}

func KvPairsWithMapT(m map[string]fmt.Stringer) *KvPairs {
	var pairs = new(KvPairs)
	if len(m) == 1 {
		for k, v := range m {
			pairs.K = k
			pairs.V = v
		}
	}
	return pairs
}

func (p *KvPairs) Key() string {
	return p.K
}

func (p *KvPairs) Value() fmt.Stringer {
	return p.V
}

func (p *KvPairs) IsString() bool {
	switch p.V.(type) {
	case *String, String:
		return true
	}
	return false
}

func (p *KvPairs) String() string {
	return fmt.Sprintf(`{"%s":"%s"}`, p.Key(), p.Value().String())
}

func IsNumber(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr, reflect.Float32,
		reflect.Float64:
		return true
	default:
		return false
	}
}

func IsString(k reflect.Kind) bool {
	if k == reflect.String {
		return true
	}
	return false
}
