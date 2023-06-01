package optional

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	type testCase[T any, U any] struct {
		name   string
		o      Of[T]
		mapper func(T) U
		want   Of[U]
	}
	isNotEmpty := func(s string) bool { return s != "" }
	tests := []testCase[string, bool]{
		{
			"valid 1",
			From("aaa"),
			isNotEmpty,
			From(true),
		},
		{
			"valid 2",
			From(""),
			isNotEmpty,
			From(false),
		},
		{
			"empty",
			None[string](),
			isNotEmpty,
			None[bool](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.o, tt.mapper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

type nonZeroTestCase[T comparable] struct {
	name string
	v    T
	want Of[T]
}

func runTestFromNonZero[T comparable](t *testing.T, tests []nonZeroTestCase[T]) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromNonZero(tt.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromNonZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromNonZero(t *testing.T) {
	s := "test"
	testsPtr := []nonZeroTestCase[*string]{
		{
			"non zero",
			&s,
			From(&s),
		},
		{
			"zero",
			nil,
			None[*string](),
		},
	}
	runTestFromNonZero(t, testsPtr)

	testsStr := []nonZeroTestCase[string]{
		{
			"non zero",
			"aa",
			From("aa"),
		},
		{
			"zero",
			"",
			None[string](),
		},
	}
	runTestFromNonZero(t, testsStr)
}

func TestFromNonZeroSlice(t *testing.T) {
	type testCase[T any] struct {
		name string
		s    []T
		want Of[[]T]
	}
	tests := []testCase[string]{
		{
			"non zero",
			[]string{"a"},
			From([]string{"a"}),
		},
		{
			"non zero (empty slice)",
			[]string{},
			From([]string{}),
		},
		{
			"zero",
			nil,
			None[[]string](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromNonZeroSlice(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromNonZeroSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
