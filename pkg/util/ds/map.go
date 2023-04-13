package ds

func MergeMap[K comparable, V any](m1, m2 map[K]V) map[K]V {
	ret := make(map[K]V, len(m1)+len(m2))
	for k, v := range m1 {
		ret[k] = v
	}
	for k, v := range m2 {
		ret[k] = v
	}
	return ret
}
