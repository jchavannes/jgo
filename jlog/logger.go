package jlog

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"io"
	"os"
	"time"
)

type Logger struct {
	LogLines []LogLine
	writers  []LogWriter
}

func (l *Logger) log(message string, level LogLevel) error {
	logLine := LogLine{
		Message: message,
		Level:   level,
		Time:    time.Now(),
	}
	var returnError error
	l.LogLines = append(l.LogLines, logLine)
	for _, logWriter := range l.writers {
		for _, writerLogLevel := range logWriter.Levels {
			if logLine.Level == writerLogLevel {
				_, err := logWriter.Writer.Write([]byte(logLine.Message))
				if err != nil {
					errorMessage := err.Error()
					// Different writers can have different errors.  We don't want to cancel the rest of the writers
					// because one of them fails. So instead, aggregate any errors.
					if returnError == nil {
						returnError = jerr.New(errorMessage)
					} else {
						returnError = jerr.Get(errorMessage, returnError)
					}
				}
			}
		}
	}
	return returnError
}

func (l *Logger) Log(message string, a ...interface{}) {
	l.log(fmt.Sprintf(message, a...), DEFAULT)
}

func (l *Logger) Debug(message string, a ...interface{}) {
	l.log(fmt.Sprintf(message, a...), DEBUG)
}

func (l *Logger) Trace(message string, a ...interface{}) {
	l.log(fmt.Sprintf(message, a...), TRACE)
}

func (l *Logger) SetLogWriter(writer LogWriter) *Logger {
	l.writers = []LogWriter{writer}
	return l
}

func (l *Logger) AddLogWriter(writer LogWriter) *Logger {
	l.writers = append(l.writers, writer)
	return l
}

func (l *Logger) SetDefaultLogLevel() *Logger {
	l.writers = []LogWriter{{
		Writer: os.Stdout,
		Levels: []LogLevel{
			DEFAULT,
		},
	}, {
		Writer: os.Stderr,
		Levels: []LogLevel{
			ERROR,
		},
	}}
	return l
}

func (l *Logger) SetDebugLogLevel() *Logger {
	l.writers = []LogWriter{{
		Writer: os.Stdout,
		Levels: []LogLevel{
			DEFAULT,
			DEBUG,
		},
	}, {
		Writer: os.Stderr,
		Levels: []LogLevel{
			ERROR,
		},
	}}
	return l
}

func (l *Logger) SetTraceLogLevel() *Logger {
	l.writers = []LogWriter{{
		Writer: os.Stdout,
		Levels: []LogLevel{
			DEFAULT,
			DEBUG,
			TRACE,
		},
	}, {
		Writer: os.Stderr,
		Levels: []LogLevel{
			ERROR,
		},
	}}
	return l
}

func GetLogger(writer io.Writer, levels []LogLevel) *Logger {
	var logger Logger
	logger.SetLogWriter(LogWriter{
		Writer: writer,
		Levels: levels,
	})
	return &logger
}
