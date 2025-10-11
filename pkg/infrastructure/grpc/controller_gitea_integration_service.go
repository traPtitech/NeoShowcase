package grpc

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	giteaintegration "github.com/traPtitech/neoshowcase/pkg/usecase/gitea-integration"
)

type GiteaIntegrationService struct {
	integration *giteaintegration.Integration
}

func NewGiteaIntegrationService(integration *giteaintegration.Integration) domain.GiteaIntegrationService {
	return &GiteaIntegrationService{
		integration: integration,
	}
}

func (s *GiteaIntegrationService) Sync(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	err := s.integration.Sync(ctx)

	return connect.NewResponse(&emptypb.Empty{}), err
}
