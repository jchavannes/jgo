package db_util

import "strings"

func GetQuestionMarks(i int) []string {
	var s = make([]string, i)
	for h := 0; h < i; h++ {
		s[h] = "?"
	}
	return s
}

func GetQuestionMarksCombined(i int) string {
	return strings.Join(GetQuestionMarks(i), ", ")
}
