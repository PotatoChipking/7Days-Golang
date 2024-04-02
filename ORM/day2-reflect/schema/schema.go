package schema

import (
	"7days-golang-learn/ORM/day2-reflect/dialect"
	"go/ast"
	"reflect"
)

// 类与表的映射

// type User struct {
//    Name string `geeorm:"PRIMARY KEY"`
//    Age  int
//}

// CREATE TABLE `User` (`Name` text PRIMARY KEY, `Age` integer);

// Field 代表数据库中的一列
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema 代表数据库的一个表
type Schema struct {
	// 映射的model
	Model interface{}
	// 表名
	Name string
	// 所有的列
	Fields []*Field
	// 列名
	FieldNames []string
	// 列名与列的映射
	fieldMap map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// 将任意对象解析为Schema
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// reflect.ValueOf() 函数返回一个包含指定值的reflect.Value类型 返回某个接口类型的实例
	// reflect.Indirect() 函数用于获取指针指向的值。如果传入的参数是一个指针，它会返回指针指向的值；如果传入的参数不是指针，它会返回参数本身。
	// .type返回该值的类型。
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()

	schema := &Schema{
		Model: dest,
		// .Name返回结构体的类型名称
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	// 遍历每一个Field列
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				// reflect.Indirect(reflect.ValueOf(p)).Type()
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
