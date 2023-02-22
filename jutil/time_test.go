package jutil_test

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	"github.com/jchavannes/jgo/jutil"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	now := time.Now()

	nowBytes := jutil.GetTimeByte(now)
	jlog.Logf("nowBytes: %x\n", nowBytes)
	now2 := jutil.GetByteTime(nowBytes)
	if now2.Unix() != now.Unix() {
		t.Error(jerr.Newf("now does not equal now2 (%s %s)", now, now2))
	}

	nowBytesNano := jutil.GetTimeByteNano(now)
	jlog.Logf("nowBytesNano: %x\n", nowBytesNano)
	now2Nano := jutil.GetByteTimeNano(nowBytesNano)
	if !now2Nano.Equal(now) {
		t.Error(jerr.Newf("now does not equal now2nano (%s %s0", now, now2Nano))
	}

	nowBytesBig := jutil.GetTimeByteBig(now)
	jlog.Logf("nowBytesBig: %x\n", nowBytesBig)
	now2Big := jutil.GetByteTimeBig(nowBytesBig)
	if now2Big.Unix() != now.Unix() {
		t.Error(jerr.Newf("now does not equal now2Big (%s %s)", now, now2Big))
	}

	nowBytesNanoBig := jutil.GetTimeByteNanoBig(now)
	jlog.Logf("nowBytesNanoBig: %x\n", nowBytesNanoBig)
	now2NanoBig := jutil.GetByteTimeNanoBig(nowBytesNanoBig)
	if now2NanoBig.Unix() != now.Unix() {
		t.Error(jerr.Newf("now does not equal now2NanoBig (%s %s)", now, now2NanoBig))
	}
}
