package ds

import "maps"

func MergeMap[K comparable, V any](m1, m2 map[K]V) map[K]V {
	ret := make(map[K]V, len(m1)+len(m2))
	maps.Copy(ret, m1)
	maps.Copy(ret, m2)
	return ret
}

// AppendMap appends key to the map, optionally initializing map if nil.
func AppendMap[K comparable, V any, M ~map[K]V](m *M, key K, value V) {
	if *m == nil {
		*m = make(map[K]V)
	}
	(*m)[key] = value
}
