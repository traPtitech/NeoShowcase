package grpc

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
)

type GiteaIntegrationServiceClient struct {
	client pbconnect.GiteaIntegrationServiceClient
}

type GiteaIntegrationServiceURL string

func NewGiteaIntegrationServiceClient(
	url GiteaIntegrationServiceURL,
) domain.GiteaIntegrationServiceClient {
	return &GiteaIntegrationServiceClient{
		client: pbconnect.NewGiteaIntegrationServiceClient(web.NewH2CClient(), string(url)),
	}
}

func (c *GiteaIntegrationServiceClient) Sync(ctx context.Context) error {
	_, err := c.client.Sync(ctx, connect.NewRequest(&emptypb.Empty{}))
	return err
}

type GiteaIntegrationServiceClientNop struct{}

func NewGiteaIntegrationServiceClientNop() domain.GiteaIntegrationServiceClient {
	return &GiteaIntegrationServiceClientNop{}
}

func (c *GiteaIntegrationServiceClientNop) Sync(ctx context.Context) error {
	return nil
}
