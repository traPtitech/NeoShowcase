package ds

func SliceOfPtr[T any](s []T) []*T {
	ret := make([]*T, len(s))
	for i := range s {
		ret[i] = &s[i]
	}
	return ret
}

func HasSuffix[T comparable](s, suffix []T) bool {
	if len(s) < len(suffix) {
		return false
	}
	return Equals(s[len(s)-len(suffix):], suffix)
}

func Equals[T comparable](s, t []T) bool {
	if len(s) != len(t) {
		return false
	}
	for i := range s {
		if s[i] != t[i] {
			return false
		}
	}
	return true
}
