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

func Get(message string, err error) error {
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

func New(message string) error {
	return JError{
		Messages: []string{message},
	}
}
