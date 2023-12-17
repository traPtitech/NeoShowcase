package builder

import (
	"context"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"sync"
	"time"

	buildkit "github.com/moby/buildkit/client"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/loop"
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

	pubKey      *ssh.PublicKeys
	imageConfig builder.ImageConfig

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
) (Service, error) {
	systemInfo, err := client.GetBuilderSystemInfo(context.Background())
	if err != nil {
		return nil, err
	}
	pubKey, err := domain.IntoPublicKey(systemInfo.SSHKey)
	if err != nil {
		return nil, err
	}
	return &builderService{
		client:    client,
		buildkit:  buildkit,
		buildpack: buildpack,

		pubKey:      pubKey,
		imageConfig: systemInfo.ImageConfig,
	}, nil
}

func (s *builderService) destImage(app *domain.Application, build *domain.Build) string {
	return s.imageConfig.ImageName(app.ID) + ":" + build.ID
}

func (s *builderService) tmpDestImage(app *domain.Application, build *domain.Build) string {
	return s.imageConfig.TmpImageName(app.ID) + ":" + build.ID
}

func (s *builderService) Start(_ context.Context) error {
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

	if s.state != nil && s.stateCancel != nil && s.state.build.ID == buildID {
		s.stateCancel()
	} else {
		log.Warnf("Skipping cancel build request for %v - a race condition or builder scheduling malfunction?", buildID)
	}
}

func (s *builderService) onRequest(req *pb.BuilderRequest) {
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
