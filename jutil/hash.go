package jutil

const (
	DefaultPrefixSuffixLen = 8
)

func ShortHash(hash string) string {
	return ShortHashLen(hash, DefaultPrefixSuffixLen)
}

func ShortHashLen(hash string, prefixSuffixLen int) string {
	hashRunes := []rune(hash)
	return string(hashRunes[:prefixSuffixLen]) + "..." + string(hashRunes[64-prefixSuffixLen:])
}
