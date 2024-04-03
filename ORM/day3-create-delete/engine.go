package day3_create_delete

import (
	"7days-golang-learn/ORM/day3-create-delete/dialect"
	"7days-golang-learn/ORM/day3-create-delete/session"
	"database/sql"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}
