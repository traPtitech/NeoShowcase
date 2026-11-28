package isolatedrunner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPullPolicy_Validate(t *testing.T) {
	for _, p := range []PullPolicy{PullAlways, PullIfNotPresent, PullNever} {
		require.NoError(t, p.Validate(), p)
	}
	for _, p := range []PullPolicy{"", "always", "Sometimes"} {
		require.Error(t, p.Validate(), p)
	}
}
