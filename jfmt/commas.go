package jfmt

import (
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

func AddCommasFloat(f float64) string {
	initPrinter()
	return strings.TrimRight(strings.TrimRight(printer.Sprintf("%f", f), "0"), ".")
}
