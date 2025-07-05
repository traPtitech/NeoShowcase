package ds

import (
	"maps"

	"github.com/samber/lo"
)

func MergeMap[K comparable, V any](mm ...map[K]V) map[K]V {
	ret := make(map[K]V, lo.SumBy(mm, len))
	for _, m := range mm {
		maps.Copy(ret, m)
	}
	return ret
}

// AppendMap appends key to the map, optionally initializing map if nil.
func AppendMap[K comparable, V any, M ~map[K]V](m *M, key K, value V) {
	if *m == nil {
		*m = make(map[K]V)
	}
	(*m)[key] = value
}
