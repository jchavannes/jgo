package jutil_test

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	"github.com/jchavannes/jgo/jutil"
	"math"
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

func TestTimeRoundSeconds(t *testing.T) {
	type TimeRoundTest struct {
		Start  time.Time
		Round  uint
		Expect time.Time
	}
	defaultStart := time.Date(2020, 1, 2, 3, 4, 56, 10, time.Local)
	var tests = []TimeRoundTest{{
		Start:  defaultStart,
		Round:  60 * 60,
		Expect: time.Date(2020, 1, 2, 3, 0, 0, 0, time.Local),
	}, {
		Start:  defaultStart,
		Round:  60,
		Expect: time.Date(2020, 1, 2, 3, 4, 0, 0, time.Local),
	}, {
		Start:  defaultStart,
		Round:  1,
		Expect: time.Date(2020, 1, 2, 3, 4, 56, 0, time.Local),
	}, {
		Start:  defaultStart,
		Round:  0,
		Expect: time.Date(2020, 1, 2, 3, 4, 56, 0, time.Local),
	}, {
		Start:  defaultStart,
		Round:  math.MaxUint / 100,
		Expect: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
	}}
	for _, test := range tests {
		rounded := jutil.TimeRoundSeconds(test.Start, test.Round)
		if rounded != test.Expect {
			t.Error(jerr.Newf("rounded does not equal expected (%s %s)", rounded, test.Expect))
		}
		jlog.Logf("start: %s rounded: %s (round: %d)\n",
			test.Start.Format(time.RFC3339Nano), rounded.Format(time.RFC3339Nano), test.Round)
	}
}
