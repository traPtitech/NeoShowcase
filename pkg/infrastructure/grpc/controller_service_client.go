package grpc

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/motoki317/sc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
)

type ControllerServiceClientConfig struct {
	URL string `mapstructure:"url" yaml:"url"`
}

type ControllerServiceClient struct {
	client      pbconnect.ControllerServiceClient
	clientCache *sc.Cache[string, pbconnect.ControllerServiceClient]
}

func NewControllerServiceClient(
	c ControllerServiceClientConfig,
) domain.ControllerServiceClient {
	return &ControllerServiceClient{
		client: pbconnect.NewControllerServiceClient(web.NewH2CClient(), c.URL),
		clientCache: sc.NewMust(func(ctx context.Context, address string) (pbconnect.ControllerServiceClient, error) {
			return pbconnect.NewControllerServiceClient(web.NewH2CClient(), address), nil
		}, time.Hour, time.Hour),
	}
}

func (c *ControllerServiceClient) GetSystemInfo(ctx context.Context) (*domain.SystemInfo, error) {
	req := connect.NewRequest(&emptypb.Empty{})
	res, err := c.client.GetSystemInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return pbconvert.FromPBSystemInfo(res.Msg), nil
}

func (c *ControllerServiceClient) FetchRepository(ctx context.Context, repositoryID string) error {
	req := connect.NewRequest(&pb.RepositoryIdRequest{RepositoryId: repositoryID})
	_, err := c.client.FetchRepository(ctx, req)
	return err
}

func (c *ControllerServiceClient) RegisterBuild(ctx context.Context, appID string) error {
	req := connect.NewRequest(&pb.ApplicationIdRequest{Id: appID})
	_, err := c.client.RegisterBuild(ctx, req)
	return err
}

func (c *ControllerServiceClient) SyncDeployments(ctx context.Context) error {
	req := connect.NewRequest(&emptypb.Empty{})
	_, err := c.client.SyncDeployments(ctx, req)
	return err
}

func (c *ControllerServiceClient) DiscoverBuildLogInstance(ctx context.Context, buildID string) (*pb.AddressInfo, error) {
	req := connect.NewRequest(&pb.BuildIdRequest{BuildId: buildID})
	res, err := c.client.DiscoverBuildLogInstance(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Msg, nil
}

func (c *ControllerServiceClient) StreamBuildLog(ctx context.Context, address string, buildID string) (<-chan *pb.BuildLog, error) {
	req := connect.NewRequest(&pb.BuildIdRequest{BuildId: buildID})
	client, _ := c.clientCache.Get(ctx, address)
	st, err := client.StreamBuildLog(ctx, req)
	if err != nil {
		return nil, err
	}

	ch := make(chan *pb.BuildLog, 100)
	go func() {
		defer close(ch)
		defer st.Close()

		for st.Receive() {
			ch <- st.Msg()
		}
		if err := st.Err(); err != nil {
			slog.ErrorContext(ctx, "failed to receive build log stream", "error", err)
		}
	}()
	return ch, nil
}

func (c *ControllerServiceClient) StartBuild(ctx context.Context) error {
	req := connect.NewRequest(&emptypb.Empty{})
	_, err := c.client.StartBuild(ctx, req)
	return err
}

func (c *ControllerServiceClient) CancelBuild(ctx context.Context, buildID string) error {
	req := connect.NewRequest(&pb.BuildIdRequest{BuildId: buildID})
	_, err := c.client.CancelBuild(ctx, req)
	return err
}

func (c *ControllerServiceClient) DiscoverBuildLogLocal(ctx context.Context, buildID string) (*pb.AddressInfo, error) {
	req := connect.NewRequest(&pb.BuildIdRequest{BuildId: buildID})
	res, err := c.client.DiscoverBuildLogLocal(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Msg, nil
}

func (c *ControllerServiceClient) StartBuildLocal(ctx context.Context) error {
	req := connect.NewRequest(&emptypb.Empty{})
	_, err := c.client.StartBuildLocal(ctx, req)
	return err
}

func (c *ControllerServiceClient) SyncDeploymentsLocal(ctx context.Context) error {
	req := connect.NewRequest(&emptypb.Empty{})
	_, err := c.client.SyncDeploymentsLocal(ctx, req)
	return err
}

func (c *ControllerServiceClient) CancelBuildLocal(ctx context.Context, buildID string) error {
	req := connect.NewRequest(&pb.BuildIdRequest{BuildId: buildID})
	_, err := c.client.CancelBuildLocal(ctx, req)
	return err
}
