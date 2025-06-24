package random

import (
	"testing"
)

func TestSecureGenerateHex(t *testing.T) {
	t.Parallel()

	set := make(map[string]bool, 1000)
	for range 1000 {
		s := SecureGenerateHex(22) // 88 bit
		if set[s] {
			t.FailNow()
		}
		set[s] = true
	}
}

func TestSecureGeneratePassword(t *testing.T) {
	t.Parallel()
	for i := range 64 {
		s := SecureGeneratePassword(i)
		if len(s) != i {
			t.Errorf("expected length %d, got %d", i, len(s))
		}
	}
}
