package grpc

import (
	"context"
	"errors"
	"io"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
)

type ControllerSSGenServiceClient struct {
	client pbconnect.ControllerSSGenServiceClient
}

func NewControllerSSGenServiceClient(
	c ControllerServiceClientConfig,
) domain.ControllerSSGenServiceClient {
	return &ControllerSSGenServiceClient{
		client: pbconnect.NewControllerSSGenServiceClient(connect.NewClient(connecthttp.NewTransport(web.NewH2CClient(), c.URL, connecthttp.WithReadMaxBytes(0)))),
	}
}

func (c *ControllerSSGenServiceClient) ConnectSSGen(ctx context.Context, onRequest func(req *pb.SSGenRequest)) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	st, err := c.client.ConnectSSGen(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	defer st.Close()

	for {
		msg, err := st.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		onRequest(msg)
	}
}
