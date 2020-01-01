package jcli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var validYeses = []string{
	"yes",
	"y",
}

func ConfirmContinue(overrideMessage ...string) bool {
	var message = "Are you sure you want to continue?"
	if len(overrideMessage) > 0 {
		message = overrideMessage[0]
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s ", message)
	text, _ := reader.ReadString('\n')
	text = strings.ToLower(strings.TrimSpace(text))
	for _, validYes := range validYeses {
		if text == validYes {
			return true
		}
	}
	return false
}
