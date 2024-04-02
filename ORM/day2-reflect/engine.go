package day2_reflect

import (
	"7days-golang-learn/ORM/day2-reflect/dialect"
	"7days-golang-learn/ORM/day2-reflect/log"
	"7days-golang-learn/ORM/day2-reflect/session"
	"database/sql"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// 1、新增数据库的类型转化，对不同数据库提出dialect，对共性进行抽象
// 2、将db的table与业务的类进行映射，主要包括表名、列名(类型、标识)
// 3、封装了数据库的建表、删表、查询某个表是否存在的方法

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)

	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	e = &Engine{
		db:      db,
		dialect: dial,
	}
	log.Info("Connect database success")
	return
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}
