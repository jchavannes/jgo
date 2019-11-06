package jutil

import (
	"regexp"
	"strings"
)

func IsMatch(full, part string) bool {
	return strings.Contains(full, part)
}

func IsMatchRegex(full, regex string) bool {
	match, _ := regexp.MatchString(regex, full)
	return match
}

func CountMatches(full, part string) int {
	return strings.Count(full, part)
}

func CountMatchesRegex(full, regex string) int {
	matches := GetMatches(full, regex)
	return len(matches)
}

func GetMatches(full, regex string) []string {
	re := regexp.MustCompile(regex)
	return re.FindAllString(full, -1)
}

func Replace(full, regex, replace string) string {
	re := regexp.MustCompile(regex)
	return re.ReplaceAllString(full, replace)
}
