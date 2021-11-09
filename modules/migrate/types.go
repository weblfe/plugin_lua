package migrate

type (
		ColumnType string
)

const (
		Text       ColumnType = "text"
		Pk         ColumnType = "pk"
		UPk        ColumnType = "upk"
		UBigPk     ColumnType = "ubigpk"
		BigPk      ColumnType = "bigPk"
		String     ColumnType = "string"
		Date       ColumnType = "date"
		DateTime   ColumnType = "datetime"
		Integer    ColumnType = "integer"
		TinyInt    ColumnType = "tinyint"
		Decimal    ColumnType = "decimal"
		Char       ColumnType = "char"
		BigInteger ColumnType = "bigint"
		SmallInt   ColumnType = "smallint"
		Float      ColumnType = "float"
		Double     ColumnType = "double"
		Binary     ColumnType = "binary"
		Bool       ColumnType = "bool"
		Money      ColumnType = "money"
		Time       ColumnType = "time"
		Timestamp  ColumnType = "timestamp"

		CategoryPk      = "pk"
		CategoryString  = "string"
		CategoryNumeric = "numeric"
		CategoryTime    = "time"
		CategoryOther   = "other"
)

var (
		columns = []ColumnType{
				Pk,
				UPk,
				BigPk,
				UBigPk,
				Char,
				String,
				Text,
				TinyInt,
				SmallInt,
				Integer,
				BigInteger,
				Float,
				Double,
				Decimal,
				DateTime,
				Timestamp,
				Time,
				Date,
				Binary,
				Bool,
				Money,
		}

		categoryMap = map[ColumnType]string{
				Pk:         CategoryPk,
				UPk:        CategoryPk,
				BigPk:      CategoryPk,
				UBigPk:     CategoryPk,
				Char:       CategoryString,
				String:     CategoryString,
				Text:       CategoryString,
				TinyInt:    CategoryNumeric,
				SmallInt:   CategoryNumeric,
				Integer:    CategoryNumeric,
				BigInteger: CategoryNumeric,
				Float:      CategoryNumeric,
				Double:     CategoryNumeric,
				Decimal:    CategoryNumeric,
				DateTime:   CategoryTime,
				Timestamp:  CategoryTime,
				Time:       CategoryTime,
				Date:       CategoryTime,
				Binary:     CategoryOther,
				Bool:       CategoryNumeric,
				Money:      CategoryNumeric,
		}
)

func (t ColumnType) String() string {
		return string(t)
}

func (t ColumnType) Check() bool {
		for _, v := range columns {
				if t == v {
						return true
				}
		}
		return false
}

func (t ColumnType) DefaultSize() []int {
		switch t {
		case Text:
				return nil
		case String:
				return []int{255}
		case Date:
				return nil
		case DateTime:
				return nil
		case Integer:
				return []int{12}
		case Decimal:
				return []int{12}
		case Char:
				return  nil
		case BigInteger:
				return []int{20}
		case SmallInt:
				return []int{2}
		case Float:
				return nil
		case Double:
				return nil
		case Binary:
				return nil
		case Money:
				return []int{12}
		case Timestamp:
				return []int{11}
		}
		return nil
}

func (t ColumnType) IsString() bool {
		if v, ok := categoryMap[t]; ok {
				if v == CategoryString {
						return true
				}
		}
		return false
}

func (t ColumnType) IsNumerical() bool {
		if v, ok := categoryMap[t]; ok {
				if v == CategoryNumeric {
						return true
				}
		}
		return false
}

func (t ColumnType) IsBytes() bool {
		return t == Binary
}

func (t ColumnType) IsTime() bool {
		if v, ok := categoryMap[t]; ok {
				if v == CategoryTime {
						return true
				}
		}
		return false
}

func (t ColumnType) GetTypeCategory() string {
		if v, ok := categoryMap[t]; ok {
				return v
		}
		return CategoryOther
}
