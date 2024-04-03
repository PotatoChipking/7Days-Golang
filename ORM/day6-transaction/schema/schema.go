package schema

import (
	"7days-golang-learn/ORM/day6-transaction/dialect"
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

// RecordValues 对于某个类，将其所有属性平铺到[]interface{}中
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	// 遍历所有属性
	for _, field := range schema.Fields {
		// 转换为接口类型，内部记录Type、原有的成员变量为value
		// Person struct{name age} => Type:Person, Value: p.name, p.age
		// 对于每个属性，将value保存到fieldValues中
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
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
