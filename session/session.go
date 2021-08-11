package session

import (
	"database/sql"
	"time"
)

const DEFAULT_TIME uint64 = 1800

type Session interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Remove(key string) error
	GetLastAccessedTime() (time.Time, error)
	UpdateLastAccessedTime()
	GetMaxAge() uint64
	SetMaxAge(age uint64)
	GetId() string
	Destroy() bool
}

type Provider interface {
	InitSession(sid string, maxAge uint64) (Session, error)
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
