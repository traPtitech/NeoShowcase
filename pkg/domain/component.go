package domain

import (
	"context"
	"io"

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
	AdditionalLinks  []*AdditionalLink
	Version          string
	Revision         string
}

type ControllerServiceClient interface {
	GetSystemInfo(ctx context.Context) (*SystemInfo, error)

	FetchRepository(ctx context.Context, repositoryID string) error
	RegisterBuild(ctx context.Context, appID string) error
	SyncDeployments(ctx context.Context) error
	DiscoverBuildLogInstance(ctx context.Context, buildID string) (*pb.AddressInfo, error)
	StreamBuildLog(ctx context.Context, address string, buildID string) (<-chan *pb.BuildLog, error)
	StartBuild(ctx context.Context) error
	CancelBuild(ctx context.Context, buildID string) error

	DiscoverBuildLogLocal(ctx context.Context, buildID string) (*pb.AddressInfo, error)
	StartBuildLocal(ctx context.Context) error
	SyncDeploymentsLocal(ctx context.Context) error
	CancelBuildLocal(ctx context.Context, buildID string) error
}

type ControllerBuilderService interface {
	pbconnect.ControllerBuilderServiceHandler
	ListenBuilderIdle() (sub <-chan struct{}, unsub func())
	ListenBuildSettled() (sub <-chan struct{}, unsub func())
	StartBuilds(ctx context.Context, buildIDs []string)
	CancelBuild(buildID string)
}

type ControllerBuilderServiceClient interface {
	GetBuilderSystemInfo(ctx context.Context) (*BuilderSystemInfo, error)
	PingBuild(ctx context.Context, buildID string) error
	StreamBuildLog(ctx context.Context, buildID string, send <-chan []byte) error
	SaveArtifact(ctx context.Context, artifact *Artifact, body []byte) error
	SaveBuildLog(ctx context.Context, buildID string, body []byte) error
	SaveRuntimeImage(ctx context.Context, buildId string, size int64) error
	ConnectBuilder(ctx context.Context, onRequest func(req *pb.BuilderRequest), response <-chan *pb.BuilderResponse) error
}

type BuildpackHelperServiceClient interface {
	CopyFileTree(ctx context.Context, destination string, tarStream io.Reader) error
	Exec(
		ctx context.Context,
		workDir string,
		cmd []string,
		envs map[string]string,
		logWriter io.Writer,
	) (code int, err error)
}

type ControllerSSGenService interface {
	pbconnect.ControllerSSGenServiceHandler
	BroadcastSSGen(req *pb.SSGenRequest)
}

type ControllerSSGenServiceClient interface {
	ConnectSSGen(ctx context.Context, onRequest func(req *pb.SSGenRequest)) error
}

type GiteaIntegrationService interface {
	pbconnect.GiteaIntegrationServiceHandler
}

type GiteaIntegrationServiceClient interface {
	Sync(ctx context.Context) error
}
