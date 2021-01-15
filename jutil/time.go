package jutil

import (
	"fmt"
	"time"
)

func GetTimeAgo(ts time.Time) string {
	delta := time.Now().Sub(ts)
	hours := int(delta.Hours())
	if hours > 0 {
		if hours >= 24 {
			if hours < 48 {
				return "1 day ago"
			}
			return fmt.Sprintf("%d days ago", hours/24)
		}
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	minutes := int(delta.Minutes())
	if minutes > 0 {
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	return fmt.Sprintf("%d seconds ago", int(delta.Seconds()))
}

func GetTimeAgoShort(ts time.Time) string {
	delta := time.Now().Sub(ts)
	hours := int(delta.Hours())
	if hours > 0 {
		if hours >= 24 {
			return fmt.Sprintf("%dd", hours/24)
		}
		return fmt.Sprintf("%dh", hours)
	}
	minutes := int(delta.Minutes())
	if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	return fmt.Sprintf("%ds", int(delta.Seconds()))
}

func GetTime(ts time.Time) string {
	return GetTimezoneTime(ts, "")
}

func GetTimezoneTime(ts time.Time, timezone string) string {
	timeLayout := "2006-01-02 15:04:05"
	if len(timezone) > 0 {
		displayLocation, err := time.LoadLocation(timezone)
		if err != nil {
			return ts.Format(timeLayout)
		}
		return ts.In(displayLocation).Format(timeLayout)
	} else {
		return ts.Format(timeLayout)
	}
}

func GetTimezoneTimeShort(ts time.Time, timezone string) string {
	timeLayout := "Jan 2006"
	if len(timezone) > 0 {
		displayLocation, err := time.LoadLocation(timezone)
		if err != nil {
			return ts.Format(timeLayout)
		}
		return ts.In(displayLocation).Format(timeLayout)
	} else {
		return ts.Format(timeLayout)
	}
}

func TimeInSlice(item time.Time, times []time.Time) bool {
	for _, a := range times {
		if a.Equal(item) {
			return true
		}
	}
	return false
}
func RoundTime(t time.Time, days int) time.Time {
	year, month, day := t.Local().Date()
	t = time.Date(year, month, roundInt(day, 1), 0, 0, 0, 0, time.Local)
	if days == 7 {
		t = t.AddDate(0, 0, int(-t.Weekday()))
	}
	return t
}

func RoundWeek(t time.Time) time.Time {
	return FirstDayOfISOWeek(t.ISOWeek())
}

func RoundMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
}

func roundInt(x, y int) int {
	x = x / y
	return x * y
}

// https://stackoverflow.com/a/18632496/744298
func FirstDayOfISOWeek(year int, week int) time.Time {
	time.Now()
	date := time.Date(year, 0, 0, 0, 0, 0, 0, time.Local)
	isoYear, isoWeek := date.ISOWeek()
	for date.Weekday() != time.Monday { // iterate back to Monday
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoYear < year { // iterate forward to the first day of the first week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoWeek < week { // iterate forward to the first day of the given week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	return date
}

func GetTimeByte(t time.Time) []byte {
	return GetInt64Data(t.Unix())
}

func GetByteTime(b []byte) time.Time {
	return time.Unix(GetInt64(b), 0)
}

func GetDurationByte(d time.Duration) []byte {
	return GetInt64Data(int64(d))
}

func GetByteDuration(b []byte) time.Duration {
	return time.Duration(GetInt64(b))
}
