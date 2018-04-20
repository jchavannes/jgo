package jerr

import "fmt"

type JError struct {
	Messages []string
}

const (
	ITALIC_START = "\x1b[3m"
	BOLD_START = "\x1b[1m"
	FORMAT_END = "\x1b[0m"
	COLOR_DEFAULT = "\x1b[39m"
	COLOR_RED = "\x1b[31m"
)

func (e JError) Error() string {
	returnString := ""
	for i := len(e.Messages) - 1; i >= 0; i-- {
		returnString += "\n " + BOLD_START + "[" + fmt.Sprintf("%d", len(e.Messages) - i) + "]" + FORMAT_END + " " + e.Messages[i]
	}
	return BOLD_START + COLOR_RED + "Error:" + COLOR_DEFAULT + FORMAT_END + returnString
}

func (e JError) Print() {
	fmt.Println(e.Error())
}

func Get(message string, err error) JError {
	var returnError JError
	switch t := err.(type) {
	case JError:
		returnError = t
		returnError.Messages = append(returnError.Messages, message)
	default:
		returnError = JError{
			Messages: []string{
				err.Error(),
				message,
			},
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

func Combine(errorArray ...error) error {
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
