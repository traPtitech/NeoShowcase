package grpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"slices"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/observability"
	"github.com/traPtitech/neoshowcase/pkg/usecase/logstream"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type builderConnection struct {
	reqSender chan<- *pb.BuilderRequest
	priority  int64
	buildID   string
}

func (c *builderConnection) Send(req *pb.BuilderRequest) {
	select {
	case c.reqSender <- req:
	default:
	}
}

func (c *builderConnection) SetBuildID(id string) {
	c.buildID = id
}

func (c *builderConnection) ClearBuildID() {
	c.buildID = ""
}

func (c *builderConnection) Busy() bool {
	return c.buildID != ""
}

type ControllerBuilderService struct {
	logStream  *logstream.Service
	systemInfo *domain.BuilderSystemInfo

	storage          domain.Storage
	appRepo          domain.ApplicationRepository
	artifactRepo     domain.ArtifactRepository
	runtimeImageRepo domain.RuntimeImageRepository
	buildRepo        domain.BuildRepository
	envRepo          domain.EnvironmentRepository
	gitRepo          domain.GitRepositoryRepository

	idle    domain.PubSub[struct{}]
	settled domain.PubSub[struct{}]

	builderConnections []*builderConnection
	lock               sync.Mutex

	metrics *observability.ControllerMetrics
}

func NewControllerBuilderService(
	logStream *logstream.Service,
	privateKey domain.PrivateKey,
	imageConfig builder.ImageConfig,
	storage domain.Storage,
	appRepo domain.ApplicationRepository,
	artifactRepo domain.ArtifactRepository,
	runtimeImageRepo domain.RuntimeImageRepository,
	buildRepo domain.BuildRepository,
	envRepo domain.EnvironmentRepository,
	gitRepo domain.GitRepositoryRepository,
	metrics *observability.ControllerMetrics,
) domain.ControllerBuilderService {
	return &ControllerBuilderService{
		logStream: logStream,
		systemInfo: &domain.BuilderSystemInfo{
			SSHKey:      privateKey,
			ImageConfig: imageConfig,
		},
		storage:          storage,
		appRepo:          appRepo,
		artifactRepo:     artifactRepo,
		runtimeImageRepo: runtimeImageRepo,
		buildRepo:        buildRepo,
		envRepo:          envRepo,
		gitRepo:          gitRepo,
		metrics:          metrics,
	}
}

func (s *ControllerBuilderService) GetBuilderSystemInfo(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.BuilderSystemInfo], error) {
	si := pbconvert.ToPBBuilderSystemInfo(s.systemInfo)
	res := connect.NewResponse(si)
	return res, nil
}

func (s *ControllerBuilderService) PingBuild(ctx context.Context, req *connect.Request[pb.BuildIdRequest]) (*connect.Response[emptypb.Empty], error) {
	now := time.Now()
	updateCond := domain.GetBuildCondition{ID: optional.From(req.Msg.BuildId)}
	updateArgs := domain.UpdateBuildArgs{UpdatedAt: optional.From(now)}
	_, err := s.buildRepo.UpdateBuild(ctx, updateCond, updateArgs)
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerBuilderService) StreamBuildLog(_ context.Context, st *connect.ClientStream[pb.BuildLogPortion]) (*connect.Response[emptypb.Empty], error) {
	for st.Receive() {
		msg := st.Msg()
		s.logStream.AppendBuildLog(msg.BuildId, msg.Log)
	}
	if err := st.Err(); err != nil {
		log.Errorf("receiving build log: %+v", err)
		return nil, errors.Wrap(err, "receiving build log")
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerBuilderService) SaveArtifact(ctx context.Context, req *connect.Request[pb.SaveArtifactRequest]) (*connect.Response[emptypb.Empty], error) {
	artifact := pbconvert.FromPBArtifact(req.Msg.Artifact)

	err := s.artifactRepo.CreateArtifact(ctx, artifact)
	if err != nil {
		return nil, err
	}
	err = domain.SaveArtifact(s.storage, artifact.ID, bytes.NewReader(req.Msg.Body))
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerBuilderService) SaveBuildLog(_ context.Context, req *connect.Request[pb.SaveBuildLogRequest]) (*connect.Response[emptypb.Empty], error) {
	err := domain.SaveBuildLog(s.storage, req.Msg.BuildId, bytes.NewReader(req.Msg.Log))
	if err != nil {
		return nil, err
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerBuilderService) SaveRuntimeImage(ctx context.Context, req *connect.Request[pb.SaveRuntimeImageRequest]) (*connect.Response[emptypb.Empty], error) {
	image := domain.NewRuntimeImage(req.Msg.BuildId, req.Msg.Size)
	err := s.runtimeImageRepo.CreateRuntimeImage(ctx, image)
	if err != nil {
		return nil, err
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ControllerBuilderService) ConnectBuilder(ctx context.Context, st *connect.BidiStream[pb.BuilderResponse, pb.BuilderRequest]) error {
	id := domain.NewID()
	log.WithField("id", id).Info("new builder connection")
	defer log.WithField("id", id).Info("builder connection closed")

	s.idle.Publish(struct{}{})

	ctx, cancel := context.WithCancel(ctx)
	reqSender := make(chan *pb.BuilderRequest)
	conn := &builderConnection{reqSender: reqSender}
	s.lock.Lock()
	s.builderConnections = append(s.builderConnections, conn)
	s.lock.Unlock()

	defer func() {
		s.lock.Lock()
		defer s.lock.Unlock()
		s.builderConnections = lo.Without(s.builderConnections, conn)
	}()

	go func() {
		defer cancel()

		for {
			res, err := st.Receive()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				log.Errorf("error receiving builder event: %+v", err)
				return
			}

			s.lock.Lock()
			switch res.Type {
			case pb.BuilderResponse_CONNECTED:
				payload := res.Body.(*pb.BuilderResponse_Connected).Connected
				conn.priority = payload.Priority

			case pb.BuilderResponse_BUILD_SETTLED:
				payload := res.Body.(*pb.BuilderResponse_Settled).Settled
				status := pbconvert.BuildStatusMapper.FromMust(payload.Status)
				err := s.finishBuild(ctx, payload.BuildId, status)
				if err != nil {
					log.Errorf("error finishing build: %+v", err)
				}
				conn.ClearBuildID()
				s.idle.Publish(struct{}{})
				s.settled.Publish(struct{}{})
				s.logStream.CloseBuildLog(payload.BuildId)
			}
			s.lock.Unlock()
		}
	}()

loop:
	for {
		select {
		case req := <-reqSender:
			err := st.Send(req)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			break loop
		}
	}

	return nil
}

func (s *ControllerBuilderService) ListenBuilderIdle() (sub <-chan struct{}, unsub func()) {
	return s.idle.Subscribe()
}

func (s *ControllerBuilderService) ListenBuildSettled() (sub <-chan struct{}, unsub func()) {
	return s.settled.Subscribe()
}

func (s *ControllerBuilderService) startBuildPayload(ctx context.Context, buildID string) (*domain.StartBuildRequest, error) {
	build, err := s.buildRepo.GetBuild(ctx, buildID)
	if err != nil {
		return nil, err
	}
	app, err := s.appRepo.GetApplication(ctx, build.ApplicationID)
	if err != nil {
		return nil, err
	}
	envs, err := s.envRepo.GetEnv(ctx, domain.GetEnvCondition{ApplicationID: optional.From(app.ID)})
	if err != nil {
		return nil, err
	}
	repo, err := s.gitRepo.GetRepository(ctx, app.RepositoryID)
	if err != nil {
		return nil, err
	}
	return &domain.StartBuildRequest{
		Repo:  repo,
		App:   app,
		Envs:  envs,
		Build: build,
	}, nil
}

func (s *ControllerBuilderService) startBuild(ctx context.Context, conn *builderConnection, buildID string) error {
	// Change build status in order to acquire lock
	now := time.Now()
	updateCond := domain.GetBuildCondition{
		ID:     optional.From(buildID),
		Status: optional.From(domain.BuildStatusQueued),
	}
	updateArgs := domain.UpdateBuildArgs{
		Status:    optional.From(domain.BuildStatusBuilding),
		StartedAt: optional.From(now),
		UpdatedAt: optional.From(now),
	}
	n, err := s.buildRepo.UpdateBuild(ctx, updateCond, updateArgs)
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("failed to acquire build lock for %v, skipping", buildID)
	}

	// Construct payload to send to builder
	req, err := s.startBuildPayload(ctx, buildID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to construct start build payload for %v", buildID))
	}
	// Send payload to builder
	conn.Send(&pb.BuilderRequest{
		Type: pb.BuilderRequest_START_BUILD,
		Body: &pb.BuilderRequest_StartBuild{StartBuild: pbconvert.ToPBStartBuildRequest(req)},
	})
	// Mark connection as busy
	conn.SetBuildID(buildID)

	// Start log stream service
	s.logStream.StartBuildLog(buildID)

	return nil
}

func (s *ControllerBuilderService) finishBuild(ctx context.Context, buildID string, status domain.BuildStatus) error {
	now := time.Now()
	updateCond := domain.GetBuildCondition{
		ID:     optional.From(buildID),
		Status: optional.From(domain.BuildStatusBuilding),
	}
	updateArgs := domain.UpdateBuildArgs{
		Status:     optional.From(status),
		UpdatedAt:  optional.From(now),
		FinishedAt: optional.From(now),
	}
	n, err := s.buildRepo.UpdateBuild(ctx, updateCond, updateArgs)
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("failed to change build status from builing to finished for %v - builder scheduling may be malfunctioning", buildID)
	}

	// metrics
	// errors are ignored
	build, err := s.buildRepo.GetBuild(ctx, buildID)
	if err != nil {
		log.WithError(err).Warn("getting build for metrics")
		return nil
	}
	app, err := s.appRepo.GetApplication(ctx, build.ApplicationID)
	if err != nil {
		log.WithError(err).Warn("getting application for metrics")
		return nil
	}
	s.metrics.IncrementBuild(status, app.Config.BuildConfig.BuildType())

	return nil
}

func (s *ControllerBuilderService) StartBuilds(ctx context.Context, buildIDs []string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Select available builders (and copy the slice)
	availableBuilders := lo.Filter(s.builderConnections, func(c *builderConnection, _ int) bool { return !c.Busy() })
	// Select from higher priority builders
	slices.SortFunc(availableBuilders, ds.MoreFunc(func(c *builderConnection) int64 { return c.priority }))

	// Send builds to available builders
	builderIdx := 0
	for _, buildID := range buildIDs {
		if builderIdx >= len(availableBuilders) {
			break
		}
		conn := availableBuilders[builderIdx]
		err := s.startBuild(ctx, conn, buildID)
		if err == nil {
			builderIdx++
		} else {
			// It is possible that some other controller has acquired build lock first
			// - in that case, skip and try the next build ID
			log.Errorf("error starting build: %+v", err)
		}
	}
}

func (s *ControllerBuilderService) CancelBuild(buildID string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	conns := lo.Filter(s.builderConnections, func(c *builderConnection, _ int) bool { return c.buildID == buildID })
	// assert len(conns) <= 1
	for _, conn := range conns {
		conn.Send(&pb.BuilderRequest{
			Type: pb.BuilderRequest_CANCEL_BUILD,
			Body: &pb.BuilderRequest_CancelBuild{CancelBuild: &pb.BuildIdRequest{
				BuildId: buildID,
			}},
		})
	}
}
