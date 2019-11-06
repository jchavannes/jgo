package jutil_test

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	"github.com/jchavannes/jgo/jutil"
	"testing"
)

const (
	Full = "This is a sentence."
	Part = "sentence"
	New  = "*$0*"

	Regex = `[A-Za-z]{4,}`

	ExpectedIsMatch           = true
	ExpectedIsMatchRegex      = true
	ExpectedCountMatches      = 1
	ExpectedCountMatchesRegex = 2
	ExpectedReplace           = "*This* is a *sentence*."
)

func TestRegex(t *testing.T) {
	isMatch := jutil.IsMatch(Full, Part)
	msg := fmt.Sprintf("isMatch is %t expected %t", isMatch, ExpectedIsMatch)
	handleMsg(t, isMatch != ExpectedIsMatch, msg)

	isMatchRegex := jutil.IsMatchRegex(Full, Regex)
	msg = fmt.Sprintf("isMatchRegex is %t expected %t", isMatchRegex, ExpectedIsMatchRegex)
	handleMsg(t, isMatchRegex != ExpectedIsMatchRegex, msg)

	countMatches := jutil.CountMatches(Full, Part)
	msg = fmt.Sprintf("countMatches is %d expected %d", countMatches, ExpectedCountMatches)
	handleMsg(t, countMatches != ExpectedCountMatches, msg)

	countMatchesRegex := jutil.CountMatchesRegex(Full, Regex)
	msg = fmt.Sprintf("countMatchesRegex is %d expected %d", countMatchesRegex, ExpectedCountMatchesRegex)
	handleMsg(t, countMatchesRegex != ExpectedCountMatchesRegex, msg)

	matches := jutil.GetMatches(Full, Regex)
	msg = fmt.Sprintf("matches is %v expected %v", matches, []string{"This", "sentence"})
	handleMsg(t, len(matches) != 2 || matches[0] != "This" || matches[1] != "sentence", msg)

	replace := jutil.Replace(Full, Regex, New)
	msg = fmt.Sprintf("replace is %s expected %s", replace, ExpectedReplace)
	handleMsg(t, replace != ExpectedReplace, msg)
}

func handleMsg(t *testing.T, isError bool, msg string) {
	if isError {
		t.Error(jerr.New(msg))
	}
	if testing.Verbose() {
		jlog.Log(msg)
	}
}
