package grpc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (s *APIService) GetSSHInfo(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.SSHInfo], error) {
	host, port, err := s.svc.GetSSHInfo(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.SSHInfo{
		Host: host,
		Port: int32(port),
	})
	return res, nil
}

func (s *APIService) GetSystemPublicKey(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.GetSystemPublicKeyResponse], error) {
	key := s.svc.GetSystemPublicKey(ctx)
	res := connect.NewResponse(&pb.GetSystemPublicKeyResponse{
		PublicKey: key,
	})
	return res, nil
}

func (s *APIService) GetAvailableDomains(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.AvailableDomains], error) {
	domains, err := s.svc.GetAvailableDomains(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	res := connect.NewResponse(&pb.AvailableDomains{
		Domains: ds.Map(domains, pbconvert.ToPBAvailableDomain),
	})
	return res, nil
}

func (s *APIService) GetAvailablePorts(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.AvailablePorts], error) {
	available, err := s.svc.GetAvailablePorts(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.AvailablePorts{
		AvailablePorts: ds.Map(available, pbconvert.ToPBAvailablePort),
	})
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
