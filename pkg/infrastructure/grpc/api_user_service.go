package grpc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (s *APIService) GetMe(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.GetMeResponse], error) {
	user, avatarURL := s.svc.GetMe(ctx)
	res := connect.NewResponse(&pb.GetMeResponse{
		User:      pbconvert.ToPBUser(user),
		AvatarUrl: avatarURL,
	})
	return res, nil
}

func (s *APIService) CreateUserKey(ctx context.Context, c *connect.Request[pb.CreateUserKeyRequest]) (*connect.Response[pb.UserKey], error) {
	key, err := s.svc.CreateUserKey(ctx, c.Msg.PublicKey)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(pbconvert.ToPBUserKey(key))
	return res, nil
}

func (s *APIService) GetUserKeys(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.GetUserKeysResponse], error) {
	keys, err := s.svc.GetUserKeys(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.GetUserKeysResponse{
		Keys: ds.Map(keys, pbconvert.ToPBUserKey),
	})
	return res, nil
}

func (s *APIService) DeleteUserKey(ctx context.Context, c *connect.Request[pb.DeleteUserKeyRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.DeleteUserKey(ctx, c.Msg.KeyId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}
