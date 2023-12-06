package pbconvert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

func ToPBUser(user *domain.User, avatarBaseURL domain.AvatarBaseURL) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Name:      user.Name,
		Admin:     user.Admin,
		AvatarUrl: user.AvatarURL(avatarBaseURL),
	}
}

func ToPBUserKey(key *domain.UserKey) *pb.UserKey {
	return &pb.UserKey{
		Id:        key.ID,
		UserId:    key.UserID,
		PublicKey: key.PublicKey,
		Name:      key.Name,
		CreatedAt: timestamppb.New(key.CreatedAt),
	}
}
