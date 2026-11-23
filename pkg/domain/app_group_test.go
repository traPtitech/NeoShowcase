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

func TestApplicationGroup_IsOwner(t *testing.T) {
	t.Parallel()
	g := domain.NewApplicationGroup("service", []string{"owner"})
	tests := []struct {
		name string
		user *domain.User
		want bool
	}{
		{"owner", &domain.User{ID: "owner"}, true},
		{"not owner", &domain.User{ID: "other"}, false},
		{"admin who is not an owner", &domain.User{ID: "admin", Admin: true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, g.IsOwner(tt.user))
		})
	}
}
