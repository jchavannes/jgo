package db_util

import (
	"encoding/hex"
	"fmt"
	"github.com/jchavannes/jgo/jutil"
)

type Config struct {
	Username string
	Password string
	Host     string
	Database string
	CryptKey string
	QueryLog bool
	TestDb   bool
	Shards   uint
}

func (m Config) GetIdHash() []byte {
	return jutil.FastHash32(fmt.Sprintf("h:%s-d:%s", m.Host, m.Database))
}

func (m Config) GetKey() []byte {
	key, _ := hex.DecodeString(m.CryptKey)
	return key
}

func (m Config) GetShardString(shard uint) string {
	return GetShardString(shard, m.Shards)
}

func GetShardPlaceHolder(count uint) string {
	if count < 10 {
		return "1"
	} else if count < 100 {
		return "2"
	} else if count < 1000 {
		return "3"
	} else if count < 10000 {
		return "4"
	}
	return "5"
}

func GetShardString(shard, count uint) string {
	if count == 0 {
		return "single"
	}
	return fmt.Sprintf("%0"+GetShardPlaceHolder(count)+"d", shard)
}
