package database

import (
	"time"
)

type Image struct {
	Iid   uint64
	Iname string
	Itime time.Time
}
