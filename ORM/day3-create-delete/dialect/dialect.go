package dialect

import "reflect"

// 针对数据库提供的数据类型支持，由于不同数据库支持的数据类型存在差异，需要针对性的设计映射
// 实现最大程度的复用和解耦
// dialect意为方言，即根据不同特征进行抽象

var dialectsMap = map[string]Dialect{}

type Dialect interface {

	// DataTypeOf 将go语言类型转化为该数据库类型
	DataTypeOf(typ reflect.Value) string

	// TableExistSQL 返回某个表是否存在的SQL，参数为表名
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
