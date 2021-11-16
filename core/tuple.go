package core

import lua "github.com/yuin/gopher-lua"

type Tuple struct {
	values []interface{}
}

func (t *Tuple) Len() int {
	return len(t.values)
}

func (t *Tuple) First() interface{} {
	if len(t.values) > 0 {
		return t.values[0]
	}
	return nil
}

func (t *Tuple) Last() interface{} {
	var size = len(t.values)
	if size > 0 {
		return t.values[size-1]
	}
	return nil
}

func (t *Tuple) Get(i int) interface{} {
	var size = t.Len()
	if size > 0 && i < size {
		return t.values[i]
	}
	return nil
}

func (t *Tuple) add(v interface{}) *Tuple {
	t.values = append(t.values, v)
	return t
}

func CreateTupleByLTable(t *lua.LTable) *Tuple {
	if t == nil {
		return nil
	}
	var (
		isTuple = true
		tu      = new(Tuple)
		passed  = &isTuple
	)
	t.ForEach(func(key lua.LValue, va lua.LValue) {
		if key == nil || !*passed {
			return
		}
		if _, ok := key.(lua.LNumber); ok {
			tu.add(va)
		} else {
			*passed = false
		}
	})
	if !*passed {
		return nil
	}
	return tu
}
