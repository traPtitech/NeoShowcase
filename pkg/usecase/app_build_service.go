package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

const (
	queueCheckInterval = 1 * time.Second
	requestInterval    = 10 * time.Second
)

type AppBuildService interface {
	QueueBuild(ctx context.Context, application *domain.Application, commit string) (string, error)
	CancelBuild(ctx context.Context, buildID string) error
	Shutdown()
}

type appBuildService struct {
	appRepo   domain.ApplicationRepository
	buildRepo domain.BuildRepository
	builder   pb.BuilderServiceClient

	queue           *ds.Queue[*buildJob]
	queueWait       sync.WaitGroup
	cancel          context.CancelFunc
	imageRegistry   string
	imageNamePrefix string
}

type buildJob struct {
	buildID string
	app     *domain.Application
	commit  string
}

func NewAppBuildService(appRepo domain.ApplicationRepository, buildRepo domain.BuildRepository, builder pb.BuilderServiceClient, registry builder.DockerImageRegistryString, prefix builder.DockerImageNamePrefixString) AppBuildService {
	ctx, cancel := context.WithCancel(context.Background())

	s := &appBuildService{
		appRepo:         appRepo,
		buildRepo:       buildRepo,
		builder:         builder,
		queue:           ds.NewQueue[*buildJob](),
		cancel:          cancel,
		imageRegistry:   string(registry),
		imageNamePrefix: string(prefix),
	}
	go s.startQueueManager(ctx)
	return s
}

func (s *appBuildService) QueueBuild(ctx context.Context, application *domain.Application, commit string) (string, error) {
	app, err := s.appRepo.GetApplication(ctx, application.ID)
	if err != nil {
		return "", fmt.Errorf("failed to QueueBuild: %w", err)
	}

	// ビルドのエントリをDBに挿入
	buildLog, err := s.buildRepo.CreateBuild(ctx, application.ID, commit)
	if err != nil {
		log.WithError(err).Errorf("failed to create build: %s", application.ID)
		return "", fmt.Errorf("failed to create build: %w", err)
	}

	s.queueWait.Add(1)
	s.queue.Push(&buildJob{
		buildID: buildLog.ID,
		app:     app,
		commit:  commit,
	})

	return buildLog.ID, nil
}

func (s *appBuildService) Shutdown() {
	s.queueWait.Wait()
	s.cancel()
}

func (s *appBuildService) CancelBuild(_ context.Context, buildID string) error {
	deleted := s.queue.DeleteIf(func(j *buildJob) bool { return j.buildID == buildID })
	if deleted == 0 {
		return fmt.Errorf("job is already canceled")
	}

	s.queueWait.Add(-deleted)
	return nil
}

func (s *appBuildService) startQueueManager(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			v := s.queue.Pop()
			if v == nil {
				time.Sleep(queueCheckInterval)
				continue
			}

		requestLoop:
			for {
				select {
				case <-ctx.Done():
					return
				default:
					res, err := s.builder.GetStatus(context.Background(), &emptypb.Empty{})
					if err != nil {
						log.WithError(err).Error("failed to get status")
						break requestLoop
					}
					if res.GetStatus() == pb.BuilderStatus_WAITING {
						err := s.requestBuild(context.Background(), v.app, v.buildID, v.commit)
						if err != nil {
							log.WithError(err).Error("failed to request build")
						}
						s.queueWait.Done()
						break requestLoop
					}

					time.Sleep(requestInterval)
				}
			}
		}
	}
}

func (s *appBuildService) requestBuild(ctx context.Context, app *domain.Application, buildID string, commit string) error {
	switch app.BuildType {
	case builder.BuildTypeRuntime:
		_, err := s.builder.StartBuildImage(ctx, &pb.StartBuildImageRequest{
			ImageName: builder.GetImageName(s.imageRegistry, s.imageNamePrefix, app.ID),
			ImageTag:  buildID,
			Source: &pb.BuildSource{
				RepositoryUrl: app.Repository.URL,
				Commit:        commit,
			},
			Options: &pb.BuildOptions{
				BaseImageName:  app.Config.BaseImage,
				DockerfileName: app.Config.DockerfileName,
				ArtifactPath:   app.Config.ArtifactPath,
				BuildCmd:       app.Config.BuildCmd,
				EntrypointCmd:  app.Config.EntrypointCmd,
			},
			BuildId:       buildID,
			ApplicationId: app.ID,
		})
		if err != nil {
			return fmt.Errorf("builder failed to start build image: %w", err)
		}

	case builder.BuildTypeStatic:
		_, err := s.builder.StartBuildStatic(ctx, &pb.StartBuildStaticRequest{
			Source: &pb.BuildSource{
				RepositoryUrl: app.Repository.URL,
				Commit:        commit,
			},
			Options: &pb.BuildOptions{
				BaseImageName:  app.Config.BaseImage,
				DockerfileName: app.Config.DockerfileName,
				ArtifactPath:   app.Config.ArtifactPath,
				BuildCmd:       app.Config.BuildCmd,
				EntrypointCmd:  app.Config.EntrypointCmd,
			},
			BuildId:       buildID,
			ApplicationId: app.ID,
		})
		if err != nil {
			return fmt.Errorf("builder failed to start build static: %w", err)
		}

	default:
		return fmt.Errorf("unknown build type: %s", app.BuildType)
	}

	log.WithField("applicationID", app.ID).
		Info("build requested")
	return nil
}
