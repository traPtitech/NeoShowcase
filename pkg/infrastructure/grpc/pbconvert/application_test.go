package pbconvert_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
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

func TestFromPBUpdateApplicationRequest_Owners(t *testing.T) {
	t.Parallel()
	now := time.Now()

	t.Run("given owners replace the owners", func(t *testing.T) {
		t.Parallel()
		args, ownerIDs := pbconvert.FromPBUpdateApplicationRequest(&pb.UpdateApplicationRequest{
			Id:       "app",
			OwnerIds: &pb.UpdateApplicationRequest_UpdateOwners{OwnerIds: []string{"u1", "u2"}},
		}, now)
		require.True(t, ownerIDs.Valid)
		assert.Equal(t, []string{"u1", "u2"}, ownerIDs.V)
		assert.Equal(t, now, args.UpdatedAt.V)
	})

	t.Run("omitted owners are left as they are", func(t *testing.T) {
		t.Parallel()
		_, ownerIDs := pbconvert.FromPBUpdateApplicationRequest(&pb.UpdateApplicationRequest{Id: "app"}, now)
		assert.False(t, ownerIDs.Valid)
	})
}
