package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
)

func (s *APIService) GetSystemInfo(ctx context.Context, _ *emptypb.Empty) (*pb.SystemInfo, error) {
	i, err := s.svc.GetSystemInfo(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return pbconvert.ToPBSystemInfo(i), nil
}

func (s *APIService) GenerateKeyPair(ctx context.Context, _ *emptypb.Empty) (*pb.GenerateKeyPairResponse, error) {
	keyID, pubKey, err := s.svc.GenerateKeyPair(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	return &pb.GenerateKeyPairResponse{
		KeyId:     keyID,
		PublicKey: pubKey,
	}, nil
}
