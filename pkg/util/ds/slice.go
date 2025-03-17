package ds

import (
	"cmp"
	"slices"

	"github.com/samber/lo"
)

func Map[T, U any](s []T, mapper func(item T) U) []U {
	ret := make([]U, len(s))
	for i := range s {
		ret[i] = mapper(s[i])
	}
	return ret
}

func LessFunc[E any, K cmp.Ordered](key func(e E) K) func(e1, e2 E) int {
	return func(e1, e2 E) int {
		k1, k2 := key(e1), key(e2)
		if k1 < k2 {
			return -1
		} else if k1 == k2 {
			return 0
		} else {
			return 1
		}
	}
}

func MoreFunc[E any, K cmp.Ordered](key func(e E) K) func(e1, e2 E) int {
	return func(e1, e2 E) int {
		k1, k2 := key(e1), key(e2)
		if k1 < k2 {
			return 1
		} else if k1 == k2 {
			return 0
		} else {
			return -1
		}
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

func HasPrefix[T comparable](s, prefix []T) bool {
	if len(s) < len(prefix) {
		return false
	}
	return Equals(s[:len(prefix)], prefix)
}

func FirstN[T any](s []T, n int) []T {
	if len(s) < n {
		return s
	}
	return s[:n]
}

func UniqMergeSlice[T comparable](s1, s2 []T) []T {
	return lo.Uniq(slices.Concat(s1, s2))
}
