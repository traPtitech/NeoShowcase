package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

func ToPBUser(user *domain.User) *pb.User {
	return &pb.User{
		Id:    user.ID,
		Name:  user.Name,
		Admin: user.Admin,
	}
}

func ToPBUserKey(key *domain.UserKey) *pb.UserKey {
	return &pb.UserKey{
		Id:        key.ID,
		UserId:    key.UserID,
		PublicKey: key.PublicKey,
	}
}