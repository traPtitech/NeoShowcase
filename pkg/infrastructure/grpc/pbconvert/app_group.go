package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

var AppGroupScopeMapper = mapper.MustNewValueMapper(map[apiserver.GetAppGroupScope]pb.GetApplicationGroupsRequest_Scope{
	apiserver.GetAppGroupScopeMine: pb.GetApplicationGroupsRequest_MINE,
	apiserver.GetAppGroupScopeAll:  pb.GetApplicationGroupsRequest_ALL,
})

func ToPBApplicationGroup(group *domain.ApplicationGroup) *pb.ApplicationGroup {
	return &pb.ApplicationGroup{
		Id:       group.ID,
		Name:     group.Name,
		OwnerIds: group.OwnerIDs,
	}
}

func FromPBUpdateApplicationGroupRequest(req *pb.UpdateApplicationGroupRequest) *domain.UpdateApplicationGroupArgs {
	return &domain.UpdateApplicationGroupArgs{
		Name: optional.FromPtr(req.Name),
		OwnerIDs: optional.FromNonZero(req.OwnerIds).Map(func(o *pb.UpdateApplicationGroupRequest_UpdateOwners) []string {
			return o.OwnerIds
		}),
	}
}
