package grpc

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
)

func (s *APIService) GetSystemInfo(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.SystemInfo], error) {
	i, err := s.svc.GetSystemInfo(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(pbconvert.ToPBSystemInfo(i))
	return res, nil
}

func (s *APIService) GenerateKeyPair(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.GenerateKeyPairResponse], error) {
	keyID, pubKey, err := s.svc.GenerateKeyPair(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.GenerateKeyPairResponse{
		KeyId:     keyID,
		PublicKey: pubKey,
	})
	return res, nil
}
