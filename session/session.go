package session

import (
	"time"
)

const DEFAULT_TIME uint64 = 1800

type Session interface {
	Set(key interface{}, value interface{})
	Get(key interface{}) interface{}
	Remove(key interface{}) error
	GetLastAccessedTime() (time.Time, error)
	UpdateLastAccessedTime()
	GetMaxAge() uint64
	SetMaxAge(age uint64)
	GetId() string
	Destroy() bool
}
