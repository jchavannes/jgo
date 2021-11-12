package jutil_test

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jutil"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	now := time.Now()
	nowBytes := jutil.GetTimeByte(now)
	now2 := jutil.GetByteTime(nowBytes)
	if now2.Unix() != now.Unix() {
		t.Error(jerr.Newf("now does not equal now2 (%s %s)", now, now2))
	}
	nowBytesNano := jutil.GetTimeByteNano(now)
	now2Nano := jutil.GetByteTimeNano(nowBytesNano)
	if !now2Nano.Equal(now) {
		t.Error(jerr.Newf("now does not equal now2nano (%s %s0", now, now2Nano))
	}
}
