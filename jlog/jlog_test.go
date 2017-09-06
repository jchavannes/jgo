package jlog_test

import (
	"bytes"
	"fmt"
	"github.com/jchavannes/jgo/jlog"
	"testing"
)

func TestLog(t *testing.T) {
	writer := new(bytes.Buffer)
	var logger jlog.Logger
	logger.SetLogWriter(jlog.LogWriter{
		Writer: writer,
		Levels: []jlog.LogLevel{
			jlog.DEBUG,
		},
	})
	message := "test message"
	logger.Debug(message)
	logged := writer.String()
	fmt.Printf("message: %s, logged: %s\n", message, logged)
	if logged != message {
		t.Errorf("logged message does not match")
	}
}
