package grpc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb/pbconnect"
)

type ControllerSSGenServiceClient struct {
	client pbconnect.ControllerSSGenServiceClient
}

func NewControllerSSGenServiceClient(
	c ControllerServiceClientConfig,
) domain.ControllerSSGenServiceClient {
	return &ControllerSSGenServiceClient{
		client: pbconnect.NewControllerSSGenServiceClient(web.NewH2CClient(), c.URL),
	}
}

func (c *ControllerSSGenServiceClient) ConnectSSGen(ctx context.Context, onRequest func(req *pb.SSGenRequest)) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	st, err := c.client.ConnectSSGen(ctx, connect.NewRequest(&emptypb.Empty{}))
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
