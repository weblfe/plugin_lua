package migrate

func NewColumn(ty string) *LuaMigrateColumn {
	var column = new(LuaMigrateColumn)
	column.TypeName = ty
	return column
}

// ColumnOf 构建字段
func ColumnOf(args []interface{}) *LuaMigrateColumn {
	var argc = len(args)
	if argc <= 0 {
		return nil
	}
	var column = new(LuaMigrateColumn)
	switch argc {
	case 1:
		if str, ok := args[0].(string); ok && str != "" {
			column.setType(str)
		}
		if column.check() {
			return column
		}
	case 2:
		if str, ok := args[0].(string); ok && str != "" {
			column.setType(str)
		}
		if num, ok := args[1].(float64); ok && num > 0 {
			column.setLength(num)
		}
		if !column.check() {
			return nil
		}
		return column
	case 3:
		if str, ok := args[0].(string); ok && str != "" {
			column.setType(str)
		}
		if num, ok := args[1].(float64); ok && num > 0 {
			column.setLength(num)
		}
		if v, ok := args[2].(string); ok && v != "" {
			column.setDefault(v)
		}
		if !column.check() {
			return nil
		}
		return column
	}
	return nil
}

func (c *LuaMigrateColumn) check() bool {
	if c.TypeName == "" {
		return false
	}
	switch c.TypeName {

	}
	return false
}

func (c *LuaMigrateColumn) setCheck(check string) *LuaMigrateColumn {
	c.Check = check
	return c
}

func (c *LuaMigrateColumn) setLength(len float64) *LuaMigrateColumn {
	c.Len = uint(len)
	return c
}

func (c *LuaMigrateColumn) setDefault(v string) *LuaMigrateColumn {
	c.Default = v
	return c
}

func (c *LuaMigrateColumn) setType(t string) *LuaMigrateColumn {
	c.TypeName = t
	return c.initDefault()
}

func (c *LuaMigrateColumn) initDefault() *LuaMigrateColumn {

	return c
}
