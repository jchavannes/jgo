package jutil_test

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jutil"
	"testing"
)

const (
	Group0 = 0
	Flag1  = 1 << iota
	Flag2
	Flag3
	Flag4
	Flag5
	Flag6
	Flag7
	Flag8
)

const (
	Group1 = 1
	Flag9  = 1 << iota
	Flag10
)

func TestBitwise(t *testing.T) {
	var bitwise = jutil.Bitwise{}
	bitwise.Set(Group0, Flag1)
	bitwise.Set(Group1, Flag9)
	if !bitwise.Has(Group0, Flag1) {
		t.Error(jerr.New("flag 1 is not set when expected to be set"))
	}
	if bitwise.Has(Group0, Flag2) {
		t.Error(jerr.New("flag 2 is set when expected not to be set"))
	}
	if !bitwise.Has(Group1, Flag9) {
		t.Error(jerr.New("flag 9 is not set when expected to be set"))
	}
}
