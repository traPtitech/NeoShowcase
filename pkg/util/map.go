package util

func MergeLabels(m1, m2 map[string]string) map[string]string {
	r := map[string]string{}
	for k, v := range m1 {
		r[k] = v
	}
	for k, v := range m2 {
		r[k] = v
	}
	return r
}
