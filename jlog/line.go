package jlog

import (
	"time"
)

type LogLine struct {
	Level   LogLevel
	Message string
	Time    time.Time
}
