package jutil

func ReverseStringSlice(slice []string) []string {
	last := len(slice) - 1
	for i := 0; i < len(slice)/2; i++ {
		slice[i], slice[last-i] = slice[last-i], slice[i]
	}
	return slice
}
