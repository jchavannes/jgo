package jutil_test

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	"github.com/jchavannes/jgo/jutil"
	"testing"
)

func TestShortHash(t *testing.T) {
	const hash = "1e940b96340761df2357ca4c57569d8558f2db4203907cf25f3d5de3c542694b"
	const expected = "1e940b96...c542694b"
	shortHash := jutil.ShortHash(hash)
	if shortHash != expected {
		t.Error(jerr.Newf("short hash (%s) does not match expected %s", shortHash, expected))
	} else {
		jlog.Logf("short hash (%s) matches expected (%s)\n", shortHash, expected)
	}
}
