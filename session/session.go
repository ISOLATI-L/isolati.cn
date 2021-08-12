package session

import (
	"database/sql"
)

const DEFAULT_TIME int64 = 1800

type Session interface {
	UpdateLastAccessedTime()
}

type Provider interface {
	InitSession(sid string, maxAge int64) (Session, error)
	GetSession(sid string) Session
	Set(sid string, key string, value interface{}) error
	Get(sid string, key string) interface{}
	Remove(sid string, key string) error
	DestroySession(sid string) error
	GCSession() bool
}

func newProvider(db *sql.DB) Provider {
	if db != nil {
		return newFromDatabase(db)
	} else {
		return newFromMemory()
	}
}
