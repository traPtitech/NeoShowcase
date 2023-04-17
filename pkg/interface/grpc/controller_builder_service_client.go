package grpc

import (
	"context"
	"io"

	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb/pbconnect"
)

type ControllerBuilderServiceClient struct {
	client pbconnect.ControllerBuilderServiceClient
}

func NewControllerBuilderServiceClient(
	c ControllerServiceClientConfig,
) domain.ControllerBuilderServiceClient {
	return &ControllerBuilderServiceClient{
		client: pbconnect.NewControllerBuilderServiceClient(web.NewH2CClient(), c.URL),
	}
}

func (c *ControllerBuilderServiceClient) ConnectBuilder(ctx context.Context, onRequest func(req *pb.BuilderRequest), response <-chan *pb.BuilderResponse) error {
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
