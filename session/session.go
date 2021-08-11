package session

import (
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
