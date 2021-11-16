package query

import "regexp"

var (
	quoteRegexp = regexp.MustCompile(`[,;^.& *'"]`)

	defaultTypeMap = map[string]string{
		`pk`:         `int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY`,
		`bigpk`:      `bigint(20) NOT NULL AUTO_INCREMENT PRIMARY KEY`,
		`ubigpk`:     `bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY`,
		`upk`:        `int(10) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY`,
		`char`:       `char(1)`,
		`string`:     `varchar(255)`,
		`varchar`:    `varchar(255)`,
		`text`:       `text`,
		`longtext`:   `longtext`,
		`smallint`:   `smallint(6)`,
		`integer`:    `int(11)`,
		`bigint`:     `bigint(20)`,
		`boolean`:    `tinyint(1)`,
		`bool`:       `tinyint(1)`,
		`float`:      `float`,
		`decimal`:    `decimal`,
		`datetime`:   `datetime`,
		`timestamp`:  `timestamp`,
		`time`:       `time`,
		`timestamps`: `int(11)`,
		`date`:       `date`,
		`money`:      `decimal(19,4)`,
		`binary`:     `blob`,
		`json`:       `json`,
	}

	expressContainer = newExpressContainer()
)

const (
	DriverType               = `mysql`
	quotePrefix              = "{{"
	quoteSuffix              = "}}"
	quoteReplacer            = "%"
	typeLengthReplacePattern = `/\(.+\)/`
	typeSpacePattern         = `/^(\w+)\s+/`
	typeReplacePattern       = `/^\w+/`
	typeFullPattern          = `/^(\w+)\((.+?)\)(.*)$/`
)

// RegisterExpressBuilder 注册支持的表达式
func RegisterExpressBuilder(expr ExpressionBuilderInterface) bool {
	return expressContainer.Register(expr)
}

func GetExpressBuilder(name string) ExpressionBuilderInterface {
	if builder, ok := expressContainer.GetBuilder(name); ok {
		return builder
	}
	return nil
}

// GetDriverType 支持的驱动类型名
func GetDriverType() string  {
		return DriverType
}