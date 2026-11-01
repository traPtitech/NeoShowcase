package grpc

import (
	"context"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
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
		client: pbconnect.NewGiteaIntegrationServiceClient(connect.NewClient(connecthttp.NewTransport(web.NewH2CClient(), string(url), connecthttp.WithReadMaxBytes(0)))),
	}
}

func (c *GiteaIntegrationServiceClient) Sync(ctx context.Context) error {
	_, err := c.client.Sync(ctx, &emptypb.Empty{})
	return err
}

type GiteaIntegrationServiceClientNop struct{}

func NewGiteaIntegrationServiceClientNop() domain.GiteaIntegrationServiceClient {
	return &GiteaIntegrationServiceClientNop{}
}

func (c *GiteaIntegrationServiceClientNop) Sync(ctx context.Context) error {
	return nil
}
