package pbconvert_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
)

func TestToPBApplication_Group(t *testing.T) {
	t.Parallel()
	app := &domain.Application{
		ID:         "app",
		DeployType: domain.DeployTypeRuntime,
		Container:  domain.ContainerStateMissing,
		Config: domain.ApplicationConfig{
			BuildConfig: &domain.BuildConfigRuntimeBuildpack{Context: "."},
		},
	}

	t.Run("an implicit group is hidden and gives its owners", func(t *testing.T) {
		t.Parallel()
		group := &domain.ApplicationGroup{ID: "g1", Implicit: true, OwnerIDs: []string{"u1"}}
		got := pbconvert.ToPBApplication(app, group, nil)
		assert.Nil(t, got.GroupId)
		assert.Equal(t, []string{"u1"}, got.OwnerIds)
	})

	t.Run("an explicit group gives its ID and owners", func(t *testing.T) {
		t.Parallel()
		group := &domain.ApplicationGroup{ID: "g1", Name: "service", OwnerIDs: []string{"u1"}}
		got := pbconvert.ToPBApplication(app, group, nil)
		require.NotNil(t, got.GroupId)
		assert.Equal(t, "g1", *got.GroupId)
		assert.Equal(t, []string{"u1"}, got.OwnerIds)
	})
}
