package grpc

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"google.golang.org/protobuf/types/known/emptypb"

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
		client: pbconnect.NewControllerBuilderServiceClient(connect.NewClient(connecthttp.NewTransport(
			web.NewH2CClient(),
			c.URL,
			connecthttp.WithReadMaxBytes(0),
		), auth.Client())),
		priority: priority,
	}
}

func (c *ControllerBuilderServiceClient) GetBuilderSystemInfo(ctx context.Context) (*domain.BuilderSystemInfo, error) {
	si, err := c.client.GetBuilderSystemInfo(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return pbconvert.FromPBBuilderSystemInfo(si), nil
}

func (c *ControllerBuilderServiceClient) PingBuild(ctx context.Context, buildID string) error {
	req := &pb.BuildIdRequest{BuildId: buildID}
	_, err := c.client.PingBuild(ctx, req)
	return err
}

func (c *ControllerBuilderServiceClient) StreamBuildLog(ctx context.Context, buildID string, send <-chan []byte) error {
	st, err := c.client.StreamBuildLog(ctx)
	if err != nil {
		return err
	}
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
	_, err = st.CloseAndReceive()
	return err
}

func (c *ControllerBuilderServiceClient) SaveArtifact(ctx context.Context, artifact *domain.Artifact, body []byte) error {
	req := &pb.SaveArtifactRequest{
		Artifact: pbconvert.ToPBArtifact(artifact),
		Body:     body,
	}
	_, err := c.client.SaveArtifact(ctx, req)
	return err
}

func (c *ControllerBuilderServiceClient) SaveBuildLog(ctx context.Context, buildID string, body []byte) error {
	req := &pb.SaveBuildLogRequest{
		BuildId: buildID,
		Log:     body,
	}
	_, err := c.client.SaveBuildLog(ctx, req)
	return err
}

func (c *ControllerBuilderServiceClient) SaveRuntimeImage(ctx context.Context, buildID string, size int64) error {
	req := &pb.SaveRuntimeImageRequest{
		BuildId: buildID,
		Size:    size,
	}
	_, err := c.client.SaveRuntimeImage(ctx, req)
	return err
}

func (c *ControllerBuilderServiceClient) ConnectBuilder(ctx context.Context, onRequest func(req *pb.BuilderRequest), response <-chan *pb.BuilderResponse) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	st, err := c.client.ConnectBuilder(ctx)
	if err != nil {
		return err
	}
	defer st.Close()

	go func() {
		defer cancel()
		defer st.CloseSend()

		// Need to send one arbitrary event to actually start the connection
		// not sure if this is a bug with connect protocol or something
		err := st.Send(&pb.BuilderResponse{
			Type: pb.BuilderResponse_CONNECTED,
			Body: &pb.BuilderResponse_Connected{Connected: &pb.ConnectedBody{
				Priority: int64(c.priority),
			}},
		})
		if err != nil {
			slog.ErrorContext(ctx, "failed to send connected event", "error", err)
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
					slog.ErrorContext(ctx, "failed to send builder response", "error", err)
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
