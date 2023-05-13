package grpc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func (s *APIService) GetSystemPublicKey(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.GetSystemPublicKeyResponse], error) {
	encoded := domain.Base64EncodedPublicKey(s.pubKey)
	res := connect.NewResponse(&pb.GetSystemPublicKeyResponse{
		PublicKey: encoded,
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

func (s *APIService) AddAvailableDomain(ctx context.Context, req *connect.Request[pb.AvailableDomain]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.AddAvailableDomain(ctx, pbconvert.FromPBAvailableDomain(req.Msg))
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *APIService) GetAvailablePorts(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.AvailablePorts], error) {
	available, unavailable, err := s.svc.GetAvailablePorts(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.AvailablePorts{
		AvailablePorts:   ds.Map(available, pbconvert.ToPBAvailablePort),
		UnavailablePorts: ds.Map(unavailable, pbconvert.ToPBUnavailablePort),
	})
	return res, nil
}

func (s *APIService) AddAvailablePort(ctx context.Context, c *connect.Request[pb.AvailablePort]) (*connect.Response[emptypb.Empty], error) {
	ap := pbconvert.FromPBAvailablePort(c.Msg)
	err := s.svc.AddAvailablePort(ctx, ap)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}
