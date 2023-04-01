package grpc

import (
	"context"
	"io"

	"github.com/bufbuild/connect-go"
	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb/pbconnect"
)

type ComponentServiceClientConfig struct {
	URL string `mapstructure:"url" yaml:"url"`
}

type ComponentServiceClient struct {
	client pbconnect.ComponentServiceClient
}

func NewComponentServiceClient(c ComponentServiceClientConfig) domain.ComponentServiceClient {
	return &ComponentServiceClient{
		client: pbconnect.NewComponentServiceClient(newInsecureClient(), c.URL),
	}
}

func (c *ComponentServiceClient) ConnectBuilder(ctx context.Context, onRequest func(req *pb.BuilderRequest), response <-chan *pb.BuilderResponse) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	st := c.client.ConnectBuilder(ctx)
	defer st.CloseResponse()

	go func() {
		defer cancel()
		defer st.CloseRequest()

		// Need to send one arbitrary event to actually start the connection
		// not sure if this is a bug with connect protocol or something
		err := st.Send(&pb.BuilderResponse{Type: pb.BuilderResponse_CONNECTED})
		if err != nil {
			log.Errorf("failed to send connected event: %+v", err)
			return
		}

		for {
			select {
			case res, ok := <-response:
				if !ok {
					return
				}
				err := st.Send(res)
				if err != nil {
					log.Errorf("failed to send builder response: %+v", err)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		req, err := st.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		onRequest(req)
	}
}

func (c *ComponentServiceClient) ConnectSSGen(ctx context.Context, onRequest func(req *pb.SSGenRequest)) error {
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
