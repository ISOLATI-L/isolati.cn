package session

import (
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
