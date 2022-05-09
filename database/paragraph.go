package database

import (
	"time"
)

type Paragraph struct {
	Pid    uint64
	Ptitle string
	Ptime  time.Time
}
