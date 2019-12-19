package jbuf_test

import (
	"github.com/jchavannes/jgo/jbuf"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewBuffer(t *testing.T) {
	var expected = [][]string{
		{"1"},
		{"2", "3"},
		{"4"},
	}
	var wg sync.WaitGroup
	updater := jbuf.NewBuffer(func(items []interface{}) {
		var nextExpected = expected[0]
		expected = expected[1:]
		var ss []string
		for _, item := range items {
			s, ok := item.(string)
			if ! ok {
				t.Error(jerr.Newf("unexpected non-string: %#v", s))
			} else {
				if len(nextExpected) == 0 {
					t.Error(jerr.Newf("unexpected item counts don't match: %s", s))
				} else if s != nextExpected[0] {
					t.Error(jerr.Newf("unexpected item: %s (expected %s)", s, nextExpected[0]))
				} else {
					nextExpected = nextExpected[1:]
				}
				ss = append(ss, s)
			}
		}
		jlog.Logf("processing: %s\n", strings.Join(ss, ", "))
		time.Sleep(100 * time.Millisecond)
		for range items {
			wg.Done()
		}
	})
	wg.Add(4)
	updater.Buffer("1")
	updater.Buffer("2")
	updater.Buffer("3")
	time.Sleep(110 * time.Millisecond)
	updater.Buffer("4")
	wg.Wait()
}
