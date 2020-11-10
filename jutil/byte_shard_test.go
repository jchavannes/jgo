package jutil_test

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	"github.com/jchavannes/jgo/jutil"
	"testing"
)

type ShardTest struct {
	String string
	Int    uint
	Shard  uint
}

func (t ShardTest) Test(tst *testing.T) {
	i := jutil.GetByteMd5Int([]byte(t.String))
	if i != t.Int {
		tst.Error(jerr.Newf("error int (%d) does not match expected (%d)", i, t.Int))
		return
	}
	shard := i % 2
	if shard != t.Shard {
		tst.Error(jerr.Newf("error shard (%d) does not match expected (%d)", shard, t.Shard))
		return
	}
	jlog.Logf("%s test success (i: %d, shard: %d)\n", t.String, t.Int, t.Shard)
}

func TestHelloShard(t *testing.T) {
	ShardTest{
		String: "hello",
		Int:    1564557354,
		Shard:  0,
	}.Test(t)
}

func TestHiyaShard(t *testing.T) {
	ShardTest{
		String: "hiya",
		Int:    3804564973,
		Shard:  1,
	}.Test(t)
}

func TestYesShard(t *testing.T) {
	ShardTest{
		String: "yes",
		Int:    2786089994,
		Shard:  0,
	}.Test(t)
}

func TestNoShard(t *testing.T) {
	ShardTest{
		String: "no",
		Int:    2141435751,
		Shard:  1,
	}.Test(t)
}
