package session

import (
	"database/sql"
)

const DEFAULT_TIME int64 = 1800

type session interface {
	getID() string
	updateLastAccessedTime(transaction *sql.Tx) error
}

type provider interface {
	initSession(transaction *sql.Tx, sid string, maxAge int64) (session, error)
	getSession(transaction *sql.Tx, sid string) (session, error)
	set(transaction *sql.Tx, sid string, key string, value any) error
	get(transaction *sql.Tx, sid string, key string) ([]byte, error)
	remove(transaction *sql.Tx, sid string, key string) error
	update(transaction *sql.Tx, sid string) error
	destroySession(transaction *sql.Tx, sid string) error
	gcSession() bool
}

func newProvider(db *sql.DB) provider {
	return newFromDatabase(db)
}
