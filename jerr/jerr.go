package jerr

import (
	"encoding/json"
	"fmt"
	"github.com/jchavannes/jgo/jutil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

type JError struct {
	Messages []string
}

const (
	boldStart    = "\x1b[1m"
	formatEnd    = "\x1b[0m"
	colorDefault = "\x1b[39m"
	colorRed     = "\x1b[31m"
	colorYellow  = "\x1b[33m"
)

func (e JError) Error() string {
	return e.getText(false)
}

func (e JError) Print() {
	fmt.Println(e.getText(false))
}

func (e JError) PrintWithStack() {
	Getf(e, "Stack:\n%s", debug.Stack()).Print()
}

func (e JError) Warn() {
	fmt.Println(e.getText(true))
}

func (e JError) Fatal() {
	e.Print()
	os.Exit(1)
}

func (e JError) getText(warn bool) string {
	returnString := ""
	for i := len(e.Messages) - 1; i >= 0; i-- {
		returnString += "\n " + boldStart + "[" + fmt.Sprintf("%d", len(e.Messages)-i) + "]" + formatEnd + " " + e.Messages[i]
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	if warn {
		return boldStart + colorYellow + "Warning at" + timeStr + ":" + colorDefault + formatEnd + returnString
	} else {
		return boldStart + colorRed + "Error at " + timeStr + ":" + colorDefault + formatEnd + returnString
	}
}

func (e JError) JSON() string {
	data, _ := json.Marshal(jutil.ReverseStringSlice(e.Messages))
	return string(data)
}

func Get(message string, err error) JError {
	if err == nil {
		return Get(message, New("nil error!"))
	}
	returnError := Create(err)
	returnError.Messages = append(returnError.Messages, message)
	return returnError
}

func Create(err error) JError {
	if err == nil {
		return Create(New("nil error!"))
	}
	var returnError JError
	switch t := err.(type) {
	case JError:
		returnError = t
	default:
		returnError = JError{
			Messages: []string{err.Error()},
		}
	}
	return returnError
}

func Getf(err error, format string, a ...interface{}) JError {
	return Get(fmt.Sprintf(format, a...), err)
}

func New(message string) JError {
	return JError{
		Messages: []string{message},
	}
}

func Newf(format string, a ...interface{}) JError {
	return New(fmt.Sprintf(format, a...))
}

func Combine(errorArray ...error) JError {
	var returnError JError
	for _, err := range errorArray {
		switch t := err.(type) {
		case JError:
			returnError.Messages = append(returnError.Messages, t.Messages...)
		default:
			returnError.Messages = append(returnError.Messages, t.Error())
		}
	}
	return returnError
}

func HasError(e error, s string) bool {
	if e == nil {
		return false
	}
	err, ok := e.(JError)
	if !ok {
		return e.Error() == s
	}
	for _, errMessage := range err.Messages {
		if errMessage == s {
			return true
		}
	}
	return false
}

func HasErrorPart(e error, s string) bool {
	if e == nil {
		return false
	}
	err, ok := e.(JError)
	if !ok {
		return strings.Contains(e.Error(), s)
	}
	for _, errMessage := range err.Messages {
		if strings.Contains(errMessage, s) {
			return true
		}
	}
	return false
}
