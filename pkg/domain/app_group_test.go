package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func TestApplicationGroup_VisibleID(t *testing.T) {
	t.Parallel()
	explicit := domain.NewApplicationGroup("service", []string{"user"})
	assert.Equal(t, optional.From(explicit.ID), explicit.VisibleID())

	implicit := domain.NewImplicitApplicationGroup([]string{"user"})
	assert.Equal(t, optional.None[string](), implicit.VisibleID(), "users do not see implicit groups")
}
