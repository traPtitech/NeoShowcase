package usecase

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
)

type AppBuildService interface {
	QueueBuild(ctx context.Context, env *domain.Environment) error
}

type appBuildService struct {
	repo    repository.ApplicationRepository
	builder pb.BuilderServiceClient

	imageRegistry   string
	imageNamePrefix string
}

func NewAppBuildService(repo repository.ApplicationRepository, builder pb.BuilderServiceClient, registry builder.DockerImageRegistryString, prefix builder.DockerImageNamePrefixString) AppBuildService {
	return &appBuildService{
		repo:            repo,
		builder:         builder,
		imageRegistry:   string(registry),
		imageNamePrefix: string(prefix),
	}
}

func (s *appBuildService) QueueBuild(ctx context.Context, env *domain.Environment) error {
	app, err := s.repo.GetApplicationByID(ctx, env.ApplicationID)
	if err != nil {
		return fmt.Errorf("failed to QueueBuild: %w", err)
	}

	// TODO キューに入れて、非同期で処理する
	return s.requestBuild(ctx, app, env)
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
