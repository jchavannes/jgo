package jutil

func Abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

func Abs64(x int64) int64 {
	if x >= 0 {
		return x
	}
	return -x
}

func MinInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func MaxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func MaxInt64(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

func MinInt64(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}

func MaxFloat64(x, y float64) float64 {
	if x < y {
		return y
	}
	return x
}

func InUintSlice(needle uint, haystack []uint) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
