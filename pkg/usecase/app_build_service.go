package usecase

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"google.golang.org/protobuf/types/known/emptypb"
)

const queueBufferSize = 10

type AppBuildService interface {
	QueueBuild(ctx context.Context, env *domain.Environment) error
}

type appBuildService struct {
	repo    repository.ApplicationRepository
	builder pb.BuilderServiceClient

	queue           chan *buildJob
	imageRegistry   string
	imageNamePrefix string
}

type buildJob struct {
	App *domain.Application
	Env *domain.Environment
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

func (s *appBuildService) QueueBuild(ctx context.Context, env *domain.Environment) error {
	app, err := s.repo.GetApplicationByID(ctx, env.ApplicationID)
	if err != nil {
		return fmt.Errorf("failed to QueueBuild: %w", err)
	}

	s.queue <- &buildJob{
		App: app,
		Env: env,
	}
	return nil
}

func (s *appBuildService) startQueueManager() {
	for v := range s.queue {
		stat, err := s.builder.GetStatus(context.Background(), &emptypb.Empty{})
		if err != nil {
			log.WithError(err).Error("failed to get status")
		}
		for {
			if stat.GetStatus() == pb.BuilderStatus_WAITING {
				err := s.requestBuild(context.Background(), v.App, v.Env)
				if err != nil {
					log.WithError(err).Error("failed to request build")
				}
				break
			}
			time.Sleep(10 * time.Second)
		}
	}
}

func (s *appBuildService) requestBuild(ctx context.Context, app *domain.Application, env *domain.Environment) error {
	switch env.BuildType {
	case builder.BuildTypeImage:
		_, err := s.builder.StartBuildImage(ctx, &pb.StartBuildImageRequest{
			ImageName: builder.GetImageName(s.imageNamePrefix, s.imageNamePrefix, env.ApplicationID),
			Source: &pb.BuildSource{
				RepositoryUrl: app.Repository.RemoteURL, // TODO ブランチ・タグ指定に対応
			},
			Options:       &pb.BuildOptions{}, // TODO 汎用ベースイメージビルドに対応させる
			EnvironmentId: env.ID,
		})
		if err != nil {
			return fmt.Errorf("builder failed to start build image: %w", err)
		}

	case builder.BuildTypeStatic:
		_, err := s.builder.StartBuildStatic(ctx, &pb.StartBuildStaticRequest{
			Source: &pb.BuildSource{
				RepositoryUrl: app.Repository.RemoteURL, // TODO ブランチ・タグ指定に対応
			},
			Options:       &pb.BuildOptions{}, // TODO 汎用ベースイメージビルドに対応させる
			EnvironmentId: env.ID,
		})
		if err != nil {
			return fmt.Errorf("builder failed to start build static: %w", err)
		}

	default:
		return fmt.Errorf("unknown build type: %s", env.BuildType)
	}

	log.WithField("envID", env.ID).
		Info("build requested")
	return nil
}
