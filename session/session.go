package session

import (
	"sync"
	"time"
)

const DEFAULT_TIME int64 = 1800

type Session interface {
	Set(key interface{}, value interface{})
	Get(key interface{}) interface{}
	Remove(key interface{}) error
	GetLastAccessedTime() time.Time
	UpdateLastAccessedTime()
	GetMaxAge() int64
	SetMaxAge(age int64)
	GetId() string
}

type SessionFromMemory struct {
	sid              string
	lock             sync.Mutex
	lastAccessedTime time.Time
	maxAge           int64
	data             map[interface{}]interface{}
}

func newSessionFromMemory(sid string) *SessionFromMemory {
	return &SessionFromMemory{
		sid:    sid,
		data:   make(map[interface{}]interface{}),
		maxAge: DEFAULT_TIME,
	}
}

func (si *SessionFromMemory) Set(key interface{}, value interface{}) {
	si.lock.Lock()
	defer si.lock.Unlock()
	si.data[key] = value
}

func (si *SessionFromMemory) Get(key interface{}) interface{} {
	return si.data[key]
}

func (si *SessionFromMemory) Remove(key interface{}) error {
	if si.data[key] != nil {
		delete(si.data, key)
	}
	return nil
}

func (si *SessionFromMemory) GetLastAccessedTime() time.Time {
	return si.lastAccessedTime
}

func (si *SessionFromMemory) UpdateLastAccessedTime() {
	si.lastAccessedTime = time.Now()
}

func (si *SessionFromMemory) GetMaxAge() int64 {
	return si.maxAge
}

func (si *SessionFromMemory) SetMaxAge(age int64) {
	si.maxAge = age
}

func (si *SessionFromMemory) GetId() string {
	return si.sid
}
