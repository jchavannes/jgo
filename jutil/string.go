package jutil

import (
	"strconv"
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

func RemoveDupeStrings(strings []string) []string {
	for i := 0; i < len(strings); i++ {
		for g := i + 1; g < len(strings); g++ {
			if strings[i] == strings[g] {
				strings = append(strings[:g], strings[g+1:]...)
				g--
			}
		}
	}
	return strings
}
