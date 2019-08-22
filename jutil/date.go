package jutil

import "time"

func GetDaySeconds(t time.Time) time.Duration {
	return time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second +
		time.Duration(t.Nanosecond())
}

func DateMin(a, b time.Time) time.Time {
	if b.IsZero() || (! a.IsZero() && a.Before(b)) {
		return a
	}
	return b
}

func DateMax(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func DateGTE(a, b time.Time) bool {
	return a.After(b) || a.Equal(b)
}

func DateLTE(a, b time.Time) bool {
	return a.Before(b) || a.Equal(b)
}
