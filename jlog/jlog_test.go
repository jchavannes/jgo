package jlog_test

import (
	"bytes"
	"fmt"
	"github.com/jchavannes/jgo/jlog"
	"testing"
)

const test_message = "test message"

func TestLog(t *testing.T) {
	writer := new(bytes.Buffer)
	var logger jlog.Logger
	logger.SetLogWriter(jlog.LogWriter{
		Writer: writer,
		Levels: []jlog.LogLevel{jlog.DEBUG},
	})
	logger.Debug(test_message)
	logged := writer.String()
	fmt.Printf("TestLog -\n test_message: %s, logged: %s\n", test_message, logged)
	if logged != test_message {
		t.Errorf("logged message does not match")
	}
}

func TestDontLog(t *testing.T) {
	writer := new(bytes.Buffer)
	var logger jlog.Logger
	logger.SetLogWriter(jlog.LogWriter{
		Writer: writer,
		Levels: []jlog.LogLevel{jlog.DEBUG},
	})
	logger.Trace(test_message)
	logged := writer.String()
	fmt.Printf("TestDontLog -\n test_message: %s, logged: %s\n", test_message, logged)
	if logged != "" {
		t.Errorf("logged trace message when log level only set to debug")
	}
}
