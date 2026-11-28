package dockerrunner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRunner_RejectsInvalidPullPolicy(t *testing.T) {
	_, err := NewRunner(nil, Options{Hostname: "host", PullPolicy: "Sometimes"})
	require.ErrorContains(t, err, "Sometimes")
}
