package mock

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/retry"
)

type BuilderServiceMock struct {
	client   domain.ControllerBuilderServiceClient
	response chan *pb.BuilderResponse

	close func()
}

func NewBuilderServiceMock(client domain.ControllerBuilderServiceClient) *BuilderServiceMock {
	return &BuilderServiceMock{
		client: client,
	}
}

func (bs *BuilderServiceMock) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	bs.close = func() {
		cancel()
	}
	bs.response = make(chan *pb.BuilderResponse)
	retry.Do(ctx, func(ctx context.Context) error {
		return bs.client.ConnectBuilder(ctx, bs.onRequest, bs.response)
	}, "connect to controller")
	return nil
}

func (bs *BuilderServiceMock) Shutdown(ctx context.Context) error {
	bs.close()
	return nil
}

func (bs *BuilderServiceMock) onRequest(req *pb.BuilderRequest) {
	// always cancel build
	switch req.Type {
	case pb.BuilderRequest_START_BUILD:
		b := req.Body.(*pb.BuilderRequest_StartBuild).StartBuild
		bs.response <- &pb.BuilderResponse{
			Type: pb.BuilderResponse_BUILD_SETTLED,
			Body: &pb.BuilderResponse_Settled{
				Settled: &pb.BuildSettled{
					BuildId: b.Build.Id,
					Status:  pb.BuildStatus_CANCELLED,
				},
			},
		}
		bs.client.SaveBuildLog(context.Background(), b.Build.Id, []byte("Build skipped by admin"))
	case pb.BuilderRequest_CANCEL_BUILD:
		// no-op
	default:
		panic("unexpected request type")
	}
}
