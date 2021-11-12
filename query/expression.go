package query

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type (
	ExpressionInterface interface {
		fmt.Stringer
		GetClass() string
	}

	ExpressionBuilderInterface interface {
		ExpressionInterface
		Build(expression ExpressionInterface, arr ArrayAble) string
	}

	ArrayExpress struct {
		className string
		values    []interface{}
	}

	ArrayExpressBuilder struct {
		ArrayExpress
	}

	expressBuilderContainer struct {
		safe  sync.RWMutex
		cache map[string]ExpressionBuilderInterface
	}
)

func newExpressContainer() *expressBuilderContainer {
	var container = new(expressBuilderContainer)
	return container
}

func (c *expressBuilderContainer) init() *expressBuilderContainer {
	c.safe = sync.RWMutex{}
	c.cache = make(map[string]ExpressionBuilderInterface)
	return c
}

func (c *expressBuilderContainer) Get(name string) (ExpressionInterface, bool) {
	c.safe.Lock()
	defer c.safe.Unlock()
	if v, ok := c.cache[name]; ok && v != nil {
		return v, true
	}
	return nil, false
}

func (c *expressBuilderContainer) Exists(name string) bool {
	c.safe.Lock()
	defer c.safe.Unlock()
	if _, ok := c.cache[name]; ok {
		return true
	}
	return false
}

func (c *expressBuilderContainer) GetBuilder(name string) (ExpressionBuilderInterface, bool) {
	c.safe.Lock()
	defer c.safe.Unlock()
	if v, ok := c.cache[name]; ok && v != nil {
		return v, true
	}
	return nil, false
}

func (c *expressBuilderContainer) Register(express ExpressionBuilderInterface) bool {
	c.safe.Lock()
	defer c.safe.Unlock()
	var name = express.GetClass()
	if v, ok := c.cache[name]; ok && v != nil {
		return false
	}
	c.cache[name] = express
	return true
}

func NewArrayExpress(value []interface{}) *ArrayExpress {
	var express = new(ArrayExpress)
	express.values = value
	express.className = `Mysql::ArrayExpress`
	return express
}

func (express *ArrayExpress) GetClass() string {
	return express.className
}

func (express *ArrayExpress) String() string {
	var arr []string
	for _, v := range express.values {
		var strAble = NewStringerAny(reflect.ValueOf(v))
		arr = append(arr, strAble.String())
	}
	return strings.Join(arr, ",")
}

func (express *ArrayExpressBuilder) Build(expression ExpressionInterface, arr ArrayAble) string {
	if express.GetClass() != expression.GetClass() {
		return ``
	}
	if expr, isArrExpr := expression.(*ArrayExpress); isArrExpr {
		_ = expr.values[0]
		for _, it := range arr.Kvs() {
			_ = it
		}
	}
	return ``
}
