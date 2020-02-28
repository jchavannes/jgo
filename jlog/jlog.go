package jlog

import (
	"fmt"
	"io"
	"runtime/debug"
	"time"
)

var (
	logLevel  LogLevel = DEFAULT
	logWriter io.Writer
)

func SetLogLevel(level LogLevel) {
	logLevel = level
}

func SetLogWriter(writer io.Writer) {
	logWriter = writer
}

func LogStack() {
	Logf("Stack:\n%s", debug.Stack())
}

func Log(a ...interface{}) {
	Logf(fmt.Sprintln(a...))
}

func Logf(message string, a ...interface{}) {
	message = fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339), message)
	if logWriter != nil {
		logWriter.Write([]byte(fmt.Sprintf(message, a...)))
	} else {
		fmt.Printf(message, a...)
	}
}

func Debug(message string) {
	Debugf(message)
}

func Debugf(message string, a ...interface{}) {
	if logLevel == DEBUG {
		Logf(message, a...)
	}
}
