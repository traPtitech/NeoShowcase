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
	return pbconvert.ToPBUser(user, s.avatarBaseURL), nil
}

func (s *APIService) GetUsers(ctx context.Context, _ *emptypb.Empty) (*pb.GetUsersResponse, error) {
	users, err := s.svc.GetUsers(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &pb.GetUsersResponse{
		Users: ds.Map(users, func(user *domain.User) *pb.User {
			return pbconvert.ToPBUser(user, s.avatarBaseURL)
		}),
	}, nil
}

func (s *APIService) CreateUserKey(ctx context.Context, req *pb.CreateUserKeyRequest) (*pb.UserKey, error) {
	key, err := s.svc.CreateUserKey(ctx, req.PublicKey, req.Name)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return pbconvert.ToPBUserKey(key), nil
}

func (s *APIService) GetUserKeys(ctx context.Context, _ *emptypb.Empty) (*pb.GetUserKeysResponse, error) {
	keys, err := s.svc.GetUserKeys(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &pb.GetUserKeysResponse{
		Keys: ds.Map(keys, pbconvert.ToPBUserKey),
	}, nil
}

func (s *APIService) DeleteUserKey(ctx context.Context, req *pb.DeleteUserKeyRequest) (*emptypb.Empty, error) {
	err := s.svc.DeleteUserKey(ctx, req.KeyId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &emptypb.Empty{}, nil
}
