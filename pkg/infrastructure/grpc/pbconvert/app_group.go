package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
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

func FromPBUpdateApplicationGroupOwners(req *pb.UpdateApplicationGroupRequest_UpdateOwners) []string {
	return req.OwnerIds
}
