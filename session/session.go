package session

import (
	"database/sql"
	"time"
)

const DEFAULT_TIME int64 = 1800

type Session interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Remove(key string) error
	GetLastAccessedTime() (time.Time, error)
	UpdateLastAccessedTime()
	GetMaxAge() int64
	SetMaxAge(age int64)
	GetId() string
}

type Provider interface {
	InitSession(sid string, maxAge int64) (Session, error)
	GetSession(sid string) Session
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
