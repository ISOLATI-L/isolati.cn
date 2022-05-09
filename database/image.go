package database

import (
	"time"
)

type Image struct {
	Iid     uint64
	Isuffix string
	Itime   time.Time
}
