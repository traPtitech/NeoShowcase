package grpc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pbconvert"
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
		Domains: lo.Map(domains, func(ad *domain.AvailableDomain, i int) *pb.AvailableDomain {
			return pbconvert.ToPBAvailableDomain(ad)
		}),
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
