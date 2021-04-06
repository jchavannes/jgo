package jutil_test

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	"github.com/jchavannes/jgo/jutil"
	"testing"
)

type ParseIntTest struct {
	String string
	Int    int64
}

func (t ParseIntTest) Test(tst *testing.T) {
	i64 := jutil.GetInt64FromString(t.String)
	if i64 != t.Int {
		tst.Error(jerr.Newf("error int (%d) does not match expected (%d)", i64, t.Int))
		return
	}
	i64 = int64(jutil.GetIntFromString(t.String))
	if i64 != t.Int {
		tst.Error(jerr.Newf("error int64 (%d) does not match expected (%d)", i64, t.Int))
		return
	}
	ui64 := jutil.GetUInt64FromString(t.String)
	if ui64 != uint64(t.Int) {
		tst.Error(jerr.Newf("error uint64 (%d) does not match expected (%d)", ui64, t.Int))
		return
	}
	ui64 = uint64(jutil.GetUIntFromString(t.String))
	if ui64 != uint64(t.Int) {
		tst.Error(jerr.Newf("error uint (%d) does not match expected (%d)", ui64, t.Int))
		return
	}
	jlog.Logf("%s test success (int: %d)\n", t.String, t.Int)
}

func TestParseInt(t *testing.T) {
	ParseIntTest{
		String: "1000",
		Int:    1000,
	}.Test(t)
}

// Commas not supported (yet)
func TestParseIntCommas(t *testing.T) {
	ParseIntTest{
		String: "1,000",
		Int:    0,
	}.Test(t)
}
