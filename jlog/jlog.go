package jlog

import "fmt"

type LogLevel string

const (
	DEBUG   LogLevel = "debug"
	DEFAULT LogLevel = "default"
)

var logLevel = DEFAULT

func SetLogLevel(level LogLevel) {
	logLevel = level
}

func Log(message string) {
	Logf(message)
}

func Logf(message string, a ...interface{}) {
	fmt.Printf(message, a)
}

func Debug(message string) {
	Debugf(message)
}

func Debugf(message string, a ...interface{}) {
	if logLevel == DEBUG {
		Logf(message, a)
	}
}
