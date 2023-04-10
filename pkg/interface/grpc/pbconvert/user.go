package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
)

func ToPBUser(user *domain.User) *pb.User {
	return &pb.User{
		Id:    user.ID,
		Name:  user.Name,
		Admin: user.Admin,
	}
}
