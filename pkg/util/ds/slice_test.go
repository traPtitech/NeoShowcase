package ds

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasPrefix(t *testing.T) {
	type testCase[T comparable] struct {
		name   string
		s      []T
		prefix []T
		want   bool
	}
	tests := []testCase[string]{
		{
			name:   "empty",
			s:      []string{},
			prefix: []string{},
			want:   true,
		},
		{
			name:   "has prefix (same length)",
			s:      []string{"a"},
			prefix: []string{"a"},
			want:   true,
		},
		{
			name:   "has prefix (different length)",
			s:      []string{"a", "b"},
			prefix: []string{"a"},
			want:   true,
		},
		{
			name:   "has no prefix (same length)",
			s:      []string{"b"},
			prefix: []string{"a"},
			want:   false,
		},
		{
			name:   "has no prefix (shorter length)",
			s:      []string{"a"},
			prefix: []string{"a", "b"},
			want:   false,
		},
		{
			name:   "has no prefix (longer length)",
			s:      []string{"b", "a"},
			prefix: []string{"a"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, HasPrefix(tt.s, tt.prefix), "HasPrefix(%v, %v)", tt.s, tt.prefix)
		})
	}
}
