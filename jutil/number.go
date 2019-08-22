package jutil

func MaxInt64(x, y int64) int64 {
	if x < y {
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
