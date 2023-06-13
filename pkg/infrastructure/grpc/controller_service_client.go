package grpc

import (
	"context"

	"github.com/bufbuild/connect-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

type ControllerServiceClientConfig struct {
	URL string `mapstructure:"url" yaml:"url"`
}

type ControllerServiceClient struct {
	client pbconnect.ControllerServiceClient
}

func NewControllerServiceClient(
	c ControllerServiceClientConfig,
) domain.ControllerServiceClient {
	return &ControllerServiceClient{
		client: pbconnect.NewControllerServiceClient(web.NewH2CClient(), c.URL),
	}
}

func (c *ControllerServiceClient) GetSSHInfo(ctx context.Context) (host string, port int, err error) {
	req := connect.NewRequest(&emptypb.Empty{})
	res, err := c.client.GetSSHInfo(ctx, req)
	if err != nil {
		return "", 0, err
	}
	return res.Msg.Host, int(res.Msg.Port), nil
}

func (c *ControllerServiceClient) GetAvailableDomains(ctx context.Context) (domain.AvailableDomainSlice, error) {
	req := connect.NewRequest(&emptypb.Empty{})
	res, err := c.client.GetAvailableDomains(ctx, req)
	if err != nil {
		return nil, err
	}
	return ds.Map(res.Msg.Domains, pbconvert.FromPBAvailableDomain), nil
}

func (c *ControllerServiceClient) GetAvailablePorts(ctx context.Context) (domain.AvailablePortSlice, error) {
	req := connect.NewRequest(&emptypb.Empty{})
	res, err := c.client.GetAvailablePorts(ctx, req)
	if err != nil {
		return nil, err
	}
	return ds.Map(res.Msg.AvailablePorts, pbconvert.FromPBAvailablePort), nil
}

func (c *ControllerServiceClient) FetchRepository(ctx context.Context, repositoryID string) error {
	req := connect.NewRequest(&pb.RepositoryIdRequest{RepositoryId: repositoryID})
	_, err := c.client.FetchRepository(ctx, req)
	return err
}

func (c *ControllerServiceClient) RegisterBuilds(ctx context.Context) error {
	req := connect.NewRequest(&emptypb.Empty{})
	_, err := c.client.RegisterBuilds(ctx, req)
	return err
}

func (c *ControllerServiceClient) SyncDeployments(ctx context.Context) error {
	req := connect.NewRequest(&emptypb.Empty{})
	_, err := c.client.SyncDeployments(ctx, req)
	return err
}

func (c *ControllerServiceClient) StreamBuildLog(ctx context.Context, buildID string) (<-chan *pb.BuildLog, error) {
	req := connect.NewRequest(&pb.BuildIdRequest{BuildId: buildID})
	st, err := c.client.StreamBuildLog(ctx, req)
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
			log.Errorf("failed to receive build log stream: %+v", err)
		}
	}()
	return ch, nil
}

func (c *ControllerServiceClient) CancelBuild(ctx context.Context, buildID string) error {
	req := connect.NewRequest(&pb.BuildIdRequest{BuildId: buildID})
	_, err := c.client.CancelBuild(ctx, req)
	return err
}
