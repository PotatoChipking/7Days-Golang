package session

import (
	"7days-golang-learn/ORM/day6-transaction/clause"
	"7days-golang-learn/ORM/day6-transaction/dialect"
	"7days-golang-learn/ORM/day6-transaction/log"
	"7days-golang-learn/ORM/day6-transaction/schema"
	"database/sql"
	"strings"
)

type Session struct {
	db *sql.DB
	// 新增数据库类型转化dialect与模式映射
	dialect  dialect.Dialect
	refTable *schema.Schema

	sql     strings.Builder
	sqlVars []interface{}
	// SQL
	clause clause.Clause

	tx *sql.Tx
}

// 这里根据DB与TX的共性，提出接口方法
// 后续可以根据不同需求，返回同一类型commenDB的不同实现：db、tx
// 之前的db每个SQL单独执行，TX的区别在于sqlite对事务的支持
// 因此提出共性后，可以按需选择db、TX

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

// DB returns tx if a tx begins. otherwise return *sql.DB
// 返回CommonDB对象，实例可以是db也可以是tx
func (s *Session) DB() CommonDB {
	// 选择tx实例
	if s.tx != nil {
		return s.tx
	}
	// 选择db实例
	return s.db
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)

	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
