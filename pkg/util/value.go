package util

func ValueOr[T comparable](v T, fallback T) T {
	var zero T
	if v == zero {
		return fallback
	}
	return v
}
