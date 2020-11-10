package jutil

import (
	"encoding/binary"
	"hash/fnv"
)

const (
	DefaultPrefixSuffixLen = 8
)

func ShortHash(hash string) string {
	return ShortHashLen(hash, DefaultPrefixSuffixLen)
}

func ShortHashLen(hash string, prefixSuffixLen int) string {
	if len(hash) < prefixSuffixLen {
		return hash
	}
	hashRunes := []rune(hash)
	return string(hashRunes[:prefixSuffixLen]) + "..." + string(hashRunes[64-prefixSuffixLen:])
}

func FastHash32(s string) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, FastHash32Uint(s))
	return bs
}

func FastHash32Uint(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
