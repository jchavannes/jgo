package jfmt

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strings"
)

var printer *message.Printer

func initPrinter() {
	if printer == nil {
		printer = message.NewPrinter(language.English)
	}
}

func AddCommas(i int64) string {
	initPrinter()
	return printer.Sprintf("%d", i)
}

func AddCommasUint(i uint64) string {
	initPrinter()
	return printer.Sprintf("%d", i)
}

func AddCommasFloat(f float64, decimals ...int) string {
	initPrinter()
	var str string
	if len(decimals) > 0 {
		str = printer.Sprintf(fmt.Sprintf("%%.%df", decimals[0]), f)
	} else {
		str = printer.Sprintf("%.9f", f)
	}
	if len(decimals) == 0 || (len(decimals) >= 2 && decimals[1] == 0) {
		str = strings.TrimRight(strings.TrimRight(str, "0"), ".")
	}
	return str
}
