package util

func MapDiff[K comparable, V any](m1, m2 map[K]V) map[K]V {
	ret := make(map[K]V, len(m1))
	for k, v := range m1 {
		ret[k] = v
	}
	for k := range m2 {
		delete(ret, k)
	}
	return ret
}
