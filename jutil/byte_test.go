package jutil_test

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jutil"
	"testing"
)

func TestHasPrefix(t *testing.T) {
	var b = []byte{0x6A, 0x02, 0x6d, 0x04}
	var prefix = []byte{0x6A}
	if ! jutil.HasPrefix(b, prefix) {
		t.Error(jerr.Newf("error prefix did not match (b: %x, prefix: %x)", b, prefix))
	}
	if jutil.HasPrefix(prefix, b) {
		t.Error(jerr.Newf("error prefix unexpectedly matched (b: %x, prefix: %x)", prefix, b))
	}
}

func TestGetUint32Data(t *testing.T) {
	jutil.GetUint32Data(23)
	jutil.GetUint32Data(1234567890)
}
