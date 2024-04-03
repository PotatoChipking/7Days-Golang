package session

import (
	"7days-golang-learn/ORM/day4-UD/log"
	"7days-golang-learn/ORM/day4-UD/schema"
	"fmt"
	"strings"

	"reflect"
)

func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable) {
		// 使用定义的dialect对schema进行解析
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// RefTable 判断是否解析过refTable
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

// 创建表，需要先通过Model方法建立表model
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))

	}
	// 表名，列名创建表
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

// 删除表
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXIST %s", s.RefTable().Name)).Exec()
	return err
}

// 表是否存在
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var temp string

	_ = row.Scan(&temp)
	return temp == s.RefTable().Name
}
