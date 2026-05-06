package utils

func MaxLimit(n int, limit int) int {
	if n <= limit {
		return n
	}
	return limit
}

func MinLimit(n int, limit int) int {
	if n > limit {
		return n
	}
	return limit
}
