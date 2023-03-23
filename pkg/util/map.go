package util

func MergeMap(m1, m2 map[string]string) map[string]string {
	r := map[string]string{}
	for k, v := range m1 {
		r[k] = v
	}
	for k, v := range m2 {
		r[k] = v
	}
	return r
}

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
