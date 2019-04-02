package jfmt

import (
	"fmt"
	"regexp"
	"strconv"
)

var commaRex *regexp.Regexp

func initCommaRex() {
	if commaRex == nil {
		commaRex = regexp.MustCompile("(\\d+)(\\d{3})")
	}
}

func AddCommas(i int64) string {
	initCommaRex()
	str := fmt.Sprintf("%d", i)
	for n := ""; n != str; {
		n = str
		str = commaRex.ReplaceAllString(str, "$1,$2")
	}
	return str
}

func AddCommasUint(i uint64) string {
	initCommaRex()
	str := fmt.Sprintf("%d", i)
	for n := ""; n != str; {
		n = str
		str = commaRex.ReplaceAllString(str, "$1,$2")
	}
	return str
}

func AddCommasFloat(f float64) string {
	initCommaRex()
	str := strconv.FormatFloat(f, 'f', 8, 64)
	for n := ""; n != str; {
		n = str
		str = commaRex.ReplaceAllString(str, "$1,$2")
	}
	return str
}
