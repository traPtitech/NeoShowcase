package ds

import (
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
)

func Map[T, U any](s []T, mapper func(item T) U) []U {
	ret := make([]U, len(s))
	for i := range s {
		ret[i] = mapper(s[i])
	}
	return ret
}

func LessFunc[E any, K constraints.Ordered](key func(e E) K) func(e1, e2 E) int {
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

func MoreFunc[E any, K constraints.Ordered](key func(e E) K) func(e1, e2 E) int {
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

func FirstN[T any](s []T, n int) []T {
	if len(s) < n {
		return s
	}
	return s[:n]
}

func UniqMergeSlice[T comparable](s1, s2 []T) []T {
	s := make([]T, 0, len(s1)+len(s2))
	for _, elt := range s1 {
		s = append(s, elt)
	}
	for _, elt := range s2 {
		s = append(s, elt)
	}
	lo.Uniq(s)
	return s
}
