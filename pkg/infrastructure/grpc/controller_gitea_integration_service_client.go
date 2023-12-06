package grpc

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
)

type ControllerGiteaIntegrationServiceClient struct {
	client pbconnect.ControllerGiteaIntegrationServiceClient
}

func NewControllerGiteaIntegrationServiceClient(
	c ControllerServiceClientConfig,
) domain.ControllerGiteaIntegrationServiceClient {
	return &ControllerGiteaIntegrationServiceClient{
		client: pbconnect.NewControllerGiteaIntegrationServiceClient(web.NewH2CClient(), c.URL),
	}
}

func (c *ControllerGiteaIntegrationServiceClient) Connect(ctx context.Context, onRequest func(req *pb.GiteaIntegrationRequest)) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	st, err := c.client.Connect(ctx, connect.NewRequest(&emptypb.Empty{}))
	if err != nil {
		return err
	}
	defer st.Close()

	for st.Receive() {
		msg := st.Msg()
		onRequest(msg)
	}
	return st.Err()
}
