package builder

import (
	"context"
	"sync"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/heroku/docker-registry-client/registry"
	buildkit "github.com/moby/buildkit/client"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/loop"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
	"github.com/traPtitech/neoshowcase/pkg/util/retry"
)

type Service interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type builderService struct {
	client    domain.ControllerBuilderServiceClient
	buildkit  *buildkit.Client
	buildpack builder.BuildpackBackend
	storage   domain.Storage
	pubKey    *ssh.PublicKeys
	config    builder.ImageConfig
	registry  *registry.Registry

	appRepo      domain.ApplicationRepository
	artifactRepo domain.ArtifactRepository
	buildRepo    domain.BuildRepository
	envRepo      domain.EnvironmentRepository
	gitRepo      domain.GitRepositoryRepository

	state       *state
	stateCancel func()
	statusLock  sync.Mutex
	response    chan<- *pb.BuilderResponse
	cancel      func()
}

func NewService(
	client domain.ControllerBuilderServiceClient,
	buildkit *buildkit.Client,
	buildpack builder.BuildpackBackend,
	storage domain.Storage,
	pubKey *ssh.PublicKeys,
	config builder.ImageConfig,
	appRepo domain.ApplicationRepository,
	artifactRepo domain.ArtifactRepository,
	buildRepo domain.BuildRepository,
	envRepo domain.EnvironmentRepository,
	gitRepo domain.GitRepositoryRepository,
) (Service, error) {
	r, err := config.NewRegistry()
	if err != nil {
		return nil, err
	}
	return &builderService{
		client:       client,
		buildkit:     buildkit,
		buildpack:    buildpack,
		storage:      storage,
		pubKey:       pubKey,
		config:       config,
		appRepo:      appRepo,
		artifactRepo: artifactRepo,
		buildRepo:    buildRepo,
		envRepo:      envRepo,
		gitRepo:      gitRepo,
		registry:     r,
	}, nil
}

func (s *builderService) destImage(app *domain.Application, build *domain.Build) string {
	return s.config.ImageName(app.ID) + ":" + build.Commit
}

func (s *builderService) tmpDestImage(app *domain.Application, build *domain.Build) string {
	return s.config.TmpImageName(app.ID) + ":" + build.Commit
}

func (s *builderService) appEnv(ctx context.Context, app *domain.Application) (map[string]string, error) {
	env, err := s.envRepo.GetEnv(ctx, domain.GetEnvCondition{ApplicationID: optional.From(app.ID)})
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(env, (*domain.Environment).GetKV), nil
}

func (s *builderService) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	response := make(chan *pb.BuilderResponse, 100)
	s.response = response

	go retry.Do(ctx, func(ctx context.Context) error {
		return s.client.ConnectBuilder(ctx, s.onRequest, response)
	}, 1*time.Second, 60*time.Second)
	go loop.Loop(ctx, s.prune, 1*time.Hour, false)

	return nil
}

func (s *builderService) Shutdown(_ context.Context) error {
	s.cancel()
	s.statusLock.Lock()
	defer s.statusLock.Unlock()
	if s.stateCancel != nil {
		s.stateCancel()
	}
	return nil
}

func (s *builderService) prune(ctx context.Context) {
	err := s.buildkit.Prune(ctx, nil, buildkit.PruneAll)
	if err != nil {
		log.Errorf("failed to prune buildkit: %+v", err)
	}
}

func (s *builderService) cancelBuild(buildID string) {
	s.statusLock.Lock()
	defer s.statusLock.Unlock()

	if s.state != nil && s.stateCancel != nil {
		if s.state.build.ID == buildID {
			s.stateCancel()
		}
	}
}

func (s *builderService) onRequest(req *pb.BuilderRequest) {
	switch req.Type {
	case pb.BuilderRequest_START_BUILD:
		b := req.Body.(*pb.BuilderRequest_StartBuild).StartBuild
		err := s.tryStartBuild(b.BuildId)
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
