package pbconvert_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
)

func TestToPBApplicationGroup(t *testing.T) {
	t.Parallel()
	got := pbconvert.ToPBApplicationGroup(&domain.ApplicationGroup{ID: "g1", Name: "service", OwnerIDs: []string{"u1"}})
	assert.Equal(t, "g1", got.Id)
	assert.Equal(t, "service", got.Name)
	assert.Equal(t, []string{"u1"}, got.OwnerIds)
}

func TestFromPBUpdateApplicationGroupRequest(t *testing.T) {
	t.Parallel()
	name := "renamed"
	args := pbconvert.FromPBUpdateApplicationGroupRequest(&pb.UpdateApplicationGroupRequest{
		Id:       "g1",
		Name:     &name,
		OwnerIds: &pb.UpdateApplicationGroupRequest_UpdateOwners{OwnerIds: []string{"u1", "u2"}},
	})
	require.True(t, args.Name.Valid)
	assert.Equal(t, "renamed", args.Name.V)
	require.True(t, args.OwnerIDs.Valid)
	assert.Equal(t, []string{"u1", "u2"}, args.OwnerIDs.V)

	args = pbconvert.FromPBUpdateApplicationGroupRequest(&pb.UpdateApplicationGroupRequest{Id: "g1"})
	assert.False(t, args.Name.Valid, "a name not given is left as it is")
	assert.False(t, args.OwnerIDs.Valid, "owners not given are left as they are")
}
