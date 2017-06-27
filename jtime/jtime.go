package jtime

const SecondsInDay = 86400

func RoundTimeToDay(inTime int64) int64 {
	outTime := inTime / SecondsInDay
	return outTime * SecondsInDay
}
