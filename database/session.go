package database

import (
	"time"
)

type Session struct {
	Sid               string
	SlastAccessedTime time.Time
	SmaxAge           uint64
	Sdata             map[string]any
}
