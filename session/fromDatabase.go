package session

import (
	"database/sql"
	"sync"
	"time"
)

type SessionFromDatabase struct {
	sid              string
	lock             sync.Mutex
	db               *sql.DB
	lastAccessedTime time.Time
	maxAge           int64
	data             map[interface{}]interface{}
}

func newSessionFromDatabase(sid string, db *sql.DB) *SessionFromDatabase {
	return &SessionFromDatabase{
		sid:    sid,
		data:   make(map[interface{}]interface{}),
		db:     db,
		maxAge: DEFAULT_TIME,
	}
}

func (si *SessionFromDatabase) Set(key interface{}, value interface{}) {
	si.lock.Lock()
	defer si.lock.Unlock()
	si.data[key] = value
}

func (si *SessionFromDatabase) Get(key interface{}) interface{} {
	return si.data[key]
}

func (si *SessionFromDatabase) Remove(key interface{}) error {
	if si.data[key] != nil {
		delete(si.data, key)
	}
	return nil
}

func (si *SessionFromDatabase) GetLastAccessedTime() time.Time {
	return si.lastAccessedTime
}

func (si *SessionFromDatabase) UpdateLastAccessedTime() {
	si.lastAccessedTime = time.Now()
}

func (si *SessionFromDatabase) GetMaxAge() int64 {
	return si.maxAge
}

func (si *SessionFromDatabase) SetMaxAge(age int64) {
	si.maxAge = age
}

func (si *SessionFromDatabase) GetId() string {
	return si.sid
}
