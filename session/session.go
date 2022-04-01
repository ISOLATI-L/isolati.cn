package session

import (
	"database/sql"
)

const DEFAULT_TIME int64 = 1800

type session interface {
	getID() string
	updateLastAccessedTime()
}

type provider interface {
	initSession(sid string, maxAge int64) (session, error)
	getSession(sid string) session
	set(sid string, key string, value any) error
	get(sid string, key string) ([]byte, error)
	remove(sid string, key string) error
	destroySession(sid string) error
	gcSession() bool
}

func newProvider(db *sql.DB) provider {
	if db != nil {
		return newFromDatabase(db)
	} else {
		return newFromMemory()
	}
}
