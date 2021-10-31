package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	queueBufferSize = 10
	requestInterval = 10 * time.Second
)

type AppBuildService interface {
	QueueBuild(ctx context.Context, branch *domain.Branch) error
	Shutdown()
}

type appBuildService struct {
	repo    repository.ApplicationRepository
	builder pb.BuilderServiceClient

	queue           chan *buildJob
	queueWait       sync.WaitGroup
	imageRegistry   string
	imageNamePrefix string
}

type buildJob struct {
	App    *domain.Application
	Branch *domain.Branch
}

func NewAppBuildService(repo repository.ApplicationRepository, builder pb.BuilderServiceClient, registry builder.DockerImageRegistryString, prefix builder.DockerImageNamePrefixString) AppBuildService {
	s := &appBuildService{
		repo:            repo,
		builder:         builder,
		queue:           make(chan *buildJob, queueBufferSize),
		imageRegistry:   string(registry),
		imageNamePrefix: string(prefix),
	}
	go s.startQueueManager()
	return s
}

func (s *appBuildService) QueueBuild(ctx context.Context, branch *domain.Branch) error {
	app, err := s.repo.GetApplicationByID(ctx, branch.ApplicationID)
	if err != nil {
		return fmt.Errorf("failed to QueueBuild: %w", err)
	}
	s.queueWait.Add(1)
	s.queue <- &buildJob{
		App:    app,
		Branch: branch,
	}
	return nil
}

func (s *appBuildService) Shutdown() {
	s.queueWait.Wait()
	close(s.queue)
}

func (s *appBuildService) startQueueManager() {
	for v := range s.queue {
		for {
			res, err := s.builder.GetStatus(context.Background(), &emptypb.Empty{})
			if err != nil {
				log.WithError(err).Error("failed to get status")
				break
			}
			if res.GetStatus() == pb.BuilderStatus_WAITING {
				err := s.requestBuild(context.Background(), v.App, v.Branch)
				if err != nil {
					log.WithError(err).Error("failed to request build")
				}
				s.queueWait.Done()
				break
			}
			time.Sleep(requestInterval)
		}
	}
}

func (s *appBuildService) requestBuild(ctx context.Context, app *domain.Application, branch *domain.Branch) error {
	switch branch.BuildType {
	case builder.BuildTypeImage:
		_, err := s.builder.StartBuildImage(ctx, &pb.StartBuildImageRequest{
			ImageName: builder.GetImageName(s.imageRegistry, s.imageNamePrefix, branch.ApplicationID),
			Source: &pb.BuildSource{
				RepositoryUrl: app.Repository.RemoteURL, // TODO ブランチ・タグ指定に対応
			},
			Options:  &pb.BuildOptions{}, // TODO 汎用ベースイメージビルドに対応させる
			BranchId: branch.ID,
		})
		if err != nil {
			return fmt.Errorf("builder failed to start build image: %w", err)
		}

	case builder.BuildTypeStatic:
		_, err := s.builder.StartBuildStatic(ctx, &pb.StartBuildStaticRequest{
			Source: &pb.BuildSource{
				RepositoryUrl: app.Repository.RemoteURL, // TODO ブランチ・タグ指定に対応
			},
			Options:  &pb.BuildOptions{}, // TODO 汎用ベースイメージビルドに対応させる
			BranchId: branch.ID,
		})
		if err != nil {
			return fmt.Errorf("builder failed to start build static: %w", err)
		}

	default:
		return fmt.Errorf("unknown build type: %s", branch.BuildType)
	}

	log.WithField("branchID", branch.ID).
		Info("build requested")
	return nil
}
