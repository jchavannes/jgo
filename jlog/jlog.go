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

var _process string

func SetProcess(process string) {
	_process = process
}

func Logf(message string, a ...interface{}) {
	var info = time.Now().Format(time.RFC3339)
	if _process != "" {
		info = info + " " + _process
	}
	message = fmt.Sprintf("[%s] %s", info, message)
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
