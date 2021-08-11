package sqlStruct

import (
	"html/template"
	"time"
)

type Video struct {
	Vid      uint64
	Vtitle   string
	Vcontent template.HTML
	Vcover   string
	Vtime    time.Time
}
