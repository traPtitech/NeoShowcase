package builder

import (
	"context"
	"sync"
	"time"

	"github.com/friendsofgo/errors"

	buildkit "github.com/moby/buildkit/client"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/git"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/registry"
	"github.com/traPtitech/neoshowcase/pkg/util/loop"
	"github.com/traPtitech/neoshowcase/pkg/util/retry"
)

type Config struct {
	StepTimeout time.Duration
}

type Service interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type ServiceImpl struct {
	config    *Config
	client    domain.ControllerBuilderServiceClient
	buildkit  *buildkit.Client
	buildpack builder.BuildpackBackend
	regclient builder.RegistryClient
	gitsvc    domain.GitService

	imageConfig builder.ImageConfig

	state       *state
	stateCancel func()
	statusLock  sync.Mutex
	response    chan<- *pb.BuilderResponse
	cancel      func()
}

func NewService(
	config *Config,
	client domain.ControllerBuilderServiceClient,
	buildkit *buildkit.Client,
	buildpack builder.BuildpackBackend,
) (*ServiceImpl, error) {
	systemInfo, err := client.GetBuilderSystemInfo(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get builder system info")
	}
	// FIXME: git service should be injected via DI,
	// but it's currently created here because it requires a public key
	// derived from a runtime value (SSHKey from systemInfo).
	pubKey, err := domain.IntoPublicKey(systemInfo.SSHKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert into public key")
	}
	gitsvc := git.NewService(pubKey)
	return &ServiceImpl{
		config:    config,
		client:    client,
		buildkit:  buildkit,
		buildpack: buildpack,
		regclient: registry.NewClient(systemInfo.ImageConfig),
		gitsvc:    gitsvc,

		imageConfig: systemInfo.ImageConfig,
	}, nil
}

func (s *ServiceImpl) destImage(app *domain.Application, build *domain.Build) string {
	return s.imageConfig.ImageName(app.ID) + ":" + build.ID
}

func (s *ServiceImpl) tmpDestImage(app *domain.Application, build *domain.Build) string {
	return s.imageConfig.TmpImageName(app.ID) + ":" + build.ID
}

func (s *ServiceImpl) imageTag(build *domain.Build) string {
	return build.ID
}

func (s *ServiceImpl) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	response := make(chan *pb.BuilderResponse, 100)
	s.response = response

	go retry.Do(ctx, func(ctx context.Context) error {
		return s.client.ConnectBuilder(ctx, s.onRequest, response)
	}, "connect to controller")
	go loop.Loop(ctx, s.prune, 1*time.Hour, false)

	return nil
}

func (s *ServiceImpl) Shutdown(_ context.Context) error {
	s.cancel()
	s.statusLock.Lock()
	defer s.statusLock.Unlock()
	if s.stateCancel != nil {
		s.stateCancel()
	}
	return nil
}

func (s *ServiceImpl) prune(ctx context.Context) {
	err := s.buildkit.Prune(ctx, nil, buildkit.PruneAll)
	if err != nil {
		log.Errorf("failed to prune buildkit: %+v", err)
	}
}

func (s *ServiceImpl) cancelBuild(buildID string) {
	s.statusLock.Lock()
	defer s.statusLock.Unlock()

	if s.state != nil && s.stateCancel != nil && s.state.build.ID == buildID {
		s.stateCancel()
	} else {
		log.Warnf("Skipping cancel build request for %v - a race condition or builder scheduling malfunction?", buildID)
	}
}

func (s *ServiceImpl) onRequest(req *pb.BuilderRequest) {
	switch req.Type {
	case pb.BuilderRequest_START_BUILD:
		b := req.Body.(*pb.BuilderRequest_StartBuild).StartBuild
		err := s.startBuild(pbconvert.FromPBStartBuildRequest(b))
		if err != nil {
			log.Errorf("failed to start build: %+v", err)
		}
	case pb.BuilderRequest_CANCEL_BUILD:
		b := req.Body.(*pb.BuilderRequest_CancelBuild).CancelBuild
		s.cancelBuild(b.BuildId)
	default:
		log.Errorf("unknown builder request type: %v", req.Type)
	}
}
