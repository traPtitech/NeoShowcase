package ds

import (
	"golang.org/x/exp/constraints"
)

func Map[T, U any](s []T, mapper func(item T) U) []U {
	ret := make([]U, len(s))
	for i := range s {
		ret[i] = mapper(s[i])
	}
	return ret
}

func LessFunc[E any, K constraints.Ordered](key func(e E) K) func(e1, e2 E) bool {
	return func(e1, e2 E) bool {
		return key(e1) < key(e2)
	}
}

func SliceOfPtr[T any](s []T) []*T {
	ret := make([]*T, len(s))
	for i := range s {
		ret[i] = &s[i]
	}
	return ret
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

func FirstN[T any](s []T, n int) []T {
	if len(s) < n {
		return s
	}
	return s[:n]
}
