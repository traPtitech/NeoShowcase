package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (s *APIService) GetMe(ctx context.Context, _ *emptypb.Empty) (*pb.User, error) {
	user := s.svc.GetMe(ctx)
	res := pbconvert.ToPBUser(user, s.avatarBaseURL)
	return res, nil
}

func (s *APIService) GetUsers(ctx context.Context, _ *emptypb.Empty) (*pb.GetUsersResponse, error) {
	users, err := s.svc.GetUsers(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &pb.GetUsersResponse{
		Users: ds.Map(users, func(user *domain.User) *pb.User {
			return pbconvert.ToPBUser(user, s.avatarBaseURL)
		}),
	}
	return res, nil
}

func (s *APIService) CreateUserKey(ctx context.Context, c *pb.CreateUserKeyRequest) (*pb.UserKey, error) {
	key, err := s.svc.CreateUserKey(ctx, c.PublicKey, c.Name)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := pbconvert.ToPBUserKey(key)
	return res, nil
}

func (s *APIService) GetUserKeys(ctx context.Context, _ *emptypb.Empty) (*pb.GetUserKeysResponse, error) {
	keys, err := s.svc.GetUserKeys(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &pb.GetUserKeysResponse{
		Keys: ds.Map(keys, pbconvert.ToPBUserKey),
	}
	return res, nil
}

func (s *APIService) DeleteUserKey(ctx context.Context, c *pb.DeleteUserKeyRequest) (*emptypb.Empty, error) {
	err := s.svc.DeleteUserKey(ctx, c.KeyId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := &emptypb.Empty{}
	return res, nil
}
