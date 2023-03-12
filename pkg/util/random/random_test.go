package random

import (
	"testing"
)

func TestSecureGenerateHex(t *testing.T) {
	t.Parallel()

	set := make(map[string]bool, 1000)
	for i := 0; i < 1000; i++ {
		s := SecureGenerateHex(22) // 88 bit
		if set[s] {
			t.FailNow()
		}
		set[s] = true
	}
}

func TestSecureGeneratePassword(t *testing.T) {
	t.Parallel()

	set := make(map[string]bool, 1000)
	for i := 0; i < 1000; i++ {
		s := SecureGeneratePassword(32) // log2(69) * 32 =~ 195 bit
		if set[s] {
			t.FailNow()
		}
		set[s] = true
	}
}
