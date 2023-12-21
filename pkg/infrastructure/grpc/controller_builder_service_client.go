package grpc

import (
	"connectrpc.com/connect"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"

	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
)

type ControllerBuilderServiceClient struct {
	client   pbconnect.ControllerBuilderServiceClient
	priority int
}

func NewControllerBuilderServiceClient(
	c ControllerServiceClientConfig,
	priority int,
	auth *TokenAuthInterceptor,
) domain.ControllerBuilderServiceClient {
	return &ControllerBuilderServiceClient{
		client: pbconnect.NewControllerBuilderServiceClient(
			web.NewH2CClient(),
			c.URL,
			connect.WithInterceptors(auth),
		),
		priority: priority,
	}
}

func (c *ControllerBuilderServiceClient) GetBuilderSystemInfo(ctx context.Context) (*domain.BuilderSystemInfo, error) {
	si, err := c.client.GetBuilderSystemInfo(ctx, connect.NewRequest(&emptypb.Empty{}))
	if err != nil {
		return nil, err
	}
	return pbconvert.FromPBBuilderSystemInfo(si.Msg), nil
}

func (c *ControllerBuilderServiceClient) PingBuild(ctx context.Context, buildID string) error {
	req := connect.NewRequest(&pb.BuildIdRequest{BuildId: buildID})
	_, err := c.client.PingBuild(ctx, req)
	return err
}

func (c *ControllerBuilderServiceClient) StreamBuildLog(ctx context.Context, buildID string, send <-chan []byte) error {
	st := c.client.StreamBuildLog(ctx)
	for logBytes := range send {
		req := &pb.BuildLogPortion{
			BuildId: buildID,
			Log:     logBytes,
		}
		err := st.Send(req)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ControllerBuilderServiceClient) SaveArtifact(ctx context.Context, artifact *domain.Artifact, body []byte) error {
	req := connect.NewRequest(&pb.SaveArtifactRequest{
		Artifact: pbconvert.ToPBArtifact(artifact),
		Body:     body,
	})
	_, err := c.client.SaveArtifact(ctx, req)
	return err
}

func (c *ControllerBuilderServiceClient) SaveBuildLog(ctx context.Context, buildID string, body []byte) error {
	req := connect.NewRequest(&pb.SaveBuildLogRequest{
		BuildId: buildID,
		Log:     body,
	})
	_, err := c.client.SaveBuildLog(ctx, req)
	return err
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
		err := st.Send(&pb.BuilderResponse{
			Type: pb.BuilderResponse_CONNECTED,
			Body: &pb.BuilderResponse_Connected{Connected: &pb.ConnectedBody{
				Priority: int64(c.priority),
			}},
		})
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
