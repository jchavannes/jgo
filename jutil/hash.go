package jutil

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
