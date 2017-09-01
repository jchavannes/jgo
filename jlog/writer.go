package jlog

import "io"

type LogWriter struct {
	Writer io.Writer
	Levels []LogLevel
}
