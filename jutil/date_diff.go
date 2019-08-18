package jutil

import (
	"fmt"
	"time"
)

type DateDiff struct {
	Year  int
	Month int
	Day   int
	Hour  int
	Min   int
	Sec   int
	a     time.Time
	b     time.Time
}

func (d DateDiff) String() string {
	return fmt.Sprintf("%d years, %d months, %d days, %d hours, %d mins, %d secs",
		d.Year, d.Month, d.Day, d.Hour, d.Min, d.Sec)
}

func (d DateDiff) Months() float32 {
	return float32(d.Year*12+d.Month) + float32(d.Day)/float32(d.DaysInLastMonth())
}

// https://github.com/jinzhu/now/blob/928c32c8eb60e699b591de5911a1c8f50d11d15a/now.go#L44
func (d DateDiff) DaysInLastMonth() int {
	lastDayOfMonth := time.
		Date(d.b.Year(), d.b.Month(), 1, 0, 0, 0, 0, d.b.Location()).
		AddDate(0, 1, 0).
		Add(-time.Nanosecond)
	return lastDayOfMonth.Day()
}

// https://stackoverflow.com/a/36531443/744298
func GetDateDiff(a, b time.Time) (diff DateDiff) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	diff.a = a
	diff.b = b

	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	diff.Year = int(y2 - y1)
	diff.Month = int(M2 - M1)
	diff.Day = int(d2 - d1)
	diff.Hour = int(h2 - h1)
	diff.Min = int(m2 - m1)
	diff.Sec = int(s2 - s1)

	if diff.Sec < 0 {
		diff.Sec += 60
		diff.Min--
	}
	if diff.Min < 0 {
		diff.Min += 60
		diff.Hour--
	}
	if diff.Hour < 0 {
		diff.Hour += 24
		diff.Day--
	}
	if diff.Day < 0 {
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		diff.Day += 32 - t.Day()
		diff.Month--
	}
	if diff.Month < 0 {
		diff.Month += 12
		diff.Year--
	}
	return
}
