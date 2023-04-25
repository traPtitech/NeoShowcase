package mapper

import (
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
)

func reverse[V1, V2 comparable](m map[V1]V2) map[V2]V1 {
	ret := make(map[V2]V1, len(m))
	for k, v := range m {
		ret[v] = k
	}
	return ret
}

func must[V any](v V, ok bool) V {
	if !ok {
		panic("not ok")
	}
	return v
}

type ValueMapper[V1, V2 comparable] struct {
	m1 map[V1]V2
	m2 map[V2]V1
}

func NewValueMapper[V1, V2 comparable](m map[V1]V2) (*ValueMapper[V1, V2], error) {
	m2 := reverse(m)
	if len(m2) != len(m) {
		return nil, errors.New("reverse map len does not match: possible typo")
	}
	return &ValueMapper[V1, V2]{
		m1: m,
		m2: m2,
	}, nil
}

func MustNewValueMapper[V1, V2 comparable](m map[V1]V2) *ValueMapper[V1, V2] {
	return lo.Must(NewValueMapper(m))
}

func (m *ValueMapper[V1, V2]) Into(v1 V1) (v2 V2, ok bool) {
	v2, ok = m.m1[v1]
	return
}

func (m *ValueMapper[V1, V2]) IntoMust(v1 V1) V2 {
	return must(m.Into(v1))
}

func (m *ValueMapper[V1, V2]) From(v2 V2) (v1 V1, ok bool) {
	v1, ok = m.m2[v2]
	return
}

func (m *ValueMapper[V1, V2]) FromMust(v2 V2) V1 {
	return must(m.From(v2))
}
