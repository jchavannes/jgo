package jutil

import (
	"strconv"
	"strings"
)

func GetBoolFromString(s string) bool {
	s = strings.ToLower(s)
	if s == "true" {
		return true
	} else if s == "false" {
		return false
	}
	i, _ := strconv.Atoi(s)
	return i == 1
}

func GetIntFromString(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func GetInt64FromString(s string) int64 {
	i, _ := strconv.ParseInt(s, 0, 64)
	return i
}

func GetUIntFromString(s string) uint {
	return uint(GetUInt64FromString(s))
}

func GetUInt64FromString(s string) uint64 {
	i, _ := strconv.ParseUint(s, 0, 64)
	return i
}

func GetFloatFromString(s string, size int) float64 {
	f, _ := strconv.ParseFloat(s, size)
	return float64(f)
}
