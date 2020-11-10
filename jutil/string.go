package jutil

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

func ReverseStringSlice(slice []string) []string {
	last := len(slice) - 1
	for i := 0; i < len(slice)/2; i++ {
		slice[i], slice[last-i] = slice[last-i], slice[i]
	}
	return slice
}

func UnescapeByteString(s string) []byte {
	s, _ = strconv.Unquote(`"` + s + `"`)
	return []byte(s)
}

func StringInSlice(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func RemoveDupeStrings(stringList []string) []string {
	for i := 0; i < len(stringList); i++ {
		if stringList[i] == "" {
			stringList = append(stringList[:i], stringList[i+1:]...)
			i--
			continue
		}
		for g := i + 1; g < len(stringList); g++ {
			if stringList[i] == stringList[g] {
				stringList = append(stringList[:g], stringList[g+1:]...)
				g--
			}
		}
	}
	return stringList
}

// https://stackoverflow.com/a/20403220/744298
func GetUtf8String(data []byte) string {
	s := fmt.Sprintf("%s", data)
	v := make([]rune, 0, len(s))
	for i, r := range s {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(s[i:])
			if size == 1 {
				continue
			}
		}
		v = append(v, r)
	}
	s = string(v)
	return s
}
