package day_1_log

import (
	"database/sql"
	"day1/log"
	"day1/session"
)

type Engine struct {
	db *sql.DB
}

// 1、Session的创建，db的封装(新建、关闭)
// 2、Session中封装了sql语句、sql值的调用，将sql与值合并，转化为Raw格式，通过Exec方法调用db执行
// 3、封装了简单的查询row、rows

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

	e = &Engine{
		db: db,
	}
	log.Info("Connect success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Fail to close")
	}
	log.Info("Close success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db)
}
