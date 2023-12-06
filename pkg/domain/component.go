package domain

import (
	"context"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
)

type SystemInfo struct {
	PublicKey  string
	PrivateKey *ssh.PublicKeys
	SSHInfo    struct {
		Host string
		Port int
	}
	AvailableDomains AvailableDomainSlice
	AvailablePorts   AvailablePortSlice
	AdminerURL       string
}

type ControllerServiceClient interface {
	GetSystemInfo(ctx context.Context) (*SystemInfo, error)

	FetchRepository(ctx context.Context, repositoryID string) error
	RegisterBuild(ctx context.Context, appID string) error
	SyncDeployments(ctx context.Context) error
	StreamBuildLog(ctx context.Context, buildID string) (<-chan *pb.BuildLog, error)
	CancelBuild(ctx context.Context, buildID string) error
}

type ControllerBuilderService interface {
	pbconnect.ControllerBuilderServiceHandler
	ListenBuilderIdle() (sub <-chan struct{}, unsub func())
	ListenBuildSettled() (sub <-chan struct{}, unsub func())
	StartBuilds(buildIDs []string)
	CancelBuild(buildID string)
}

type ControllerBuilderServiceClient interface {
	ConnectBuilder(ctx context.Context, onRequest func(req *pb.BuilderRequest), response <-chan *pb.BuilderResponse) error
}

type ControllerSSGenService interface {
	pbconnect.ControllerSSGenServiceHandler
	BroadcastSSGen(req *pb.SSGenRequest)
}

type ControllerSSGenServiceClient interface {
	ConnectSSGen(ctx context.Context, onRequest func(req *pb.SSGenRequest)) error
}

type ControllerGiteaIntegrationService interface {
	pbconnect.ControllerGiteaIntegrationServiceHandler
	Broadcast(req *pb.GiteaIntegrationRequest)
}

type ControllerGiteaIntegrationServiceClient interface {
	Connect(ctx context.Context, onRequest func(req *pb.GiteaIntegrationRequest)) error
}
