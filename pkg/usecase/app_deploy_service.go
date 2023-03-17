package usecase

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type AppDeployService interface {
	StartDeployment(ctx context.Context, app *domain.Application, build *domain.Build) error
}

type appDeployService struct {
	backend domain.Backend
	appRepo repository.ApplicationRepository
	ss      pb.StaticSiteServiceClient

	imageRegistry   string
	imageNamePrefix string
}

func NewAppDeployService(
	backend domain.Backend,
	appRepo repository.ApplicationRepository,
	ss pb.StaticSiteServiceClient,
	registry builder.DockerImageRegistryString,
	prefix builder.DockerImageNamePrefixString,
) AppDeployService {
	return &appDeployService{
		backend:         backend,
		appRepo:         appRepo,
		ss:              ss,
		imageRegistry:   string(registry),
		imageNamePrefix: string(prefix),
	}
}

func (s *appDeployService) StartDeployment(ctx context.Context, app *domain.Application, build *domain.Build) error {
	// Lock app state
	err := s.appRepo.UpdateApplication(ctx, app.ID, repository.UpdateApplicationArgs{State: optional.From(domain.ApplicationStateDeploying)})
	if err != nil {
		return err
	}
	go s.deployAndHandleError(app, build)
	return nil
}

func (s *appDeployService) deployAndHandleError(app *domain.Application, build *domain.Build) {
	start := time.Now()
	ctx := context.Background()
	err := s.deploy(ctx, app, build)
	if err != nil {
		log.WithError(err).WithField("application", app.ID).WithField("build", build.ID).Error("failed to deploy")
		err = s.appRepo.UpdateApplication(ctx, app.ID, repository.UpdateApplicationArgs{State: optional.From(domain.ApplicationStateErrored)})
		if err != nil {
			log.WithError(err).Error("failed to update application state back")
		}
		return
	}
	log.WithField("application", app.ID).WithField("build", build.ID).Infof("deploy succeeded in %v", time.Since(start))
}

func (s *appDeployService) deploy(ctx context.Context, app *domain.Application, build *domain.Build) error {
	if app.BuildType == builder.BuildTypeImage {
		err := s.recreateContainer(ctx, app, build)
		if err != nil {
			return err
		}
	}

	err := s.appRepo.UpdateApplication(ctx, app.ID, repository.UpdateApplicationArgs{
		State:         optional.From(domain.ApplicationStateRunning),
		CurrentCommit: optional.From(build.Commit),
	})
	if err != nil {
		return fmt.Errorf("failed to update application")
	}

	if app.BuildType == builder.BuildTypeStatic {
		if _, err := s.ss.Reload(ctx, &emptypb.Empty{}); err != nil {
			return fmt.Errorf("failed to reload static site server: %w", err)
		}
	}
	return nil
}

func (s *appDeployService) recreateContainer(ctx context.Context, app *domain.Application, build *domain.Build) error {
	// TODO Ingressの設定がここでする必要があるかどうか確認する
	var httpProxy *domain.ContainerHTTPProxy
	if app.Website.Valid {
		httpProxy = &domain.ContainerHTTPProxy{
			Domain: app.Website.V.FQDN,
			Port:   app.Website.V.Port,
		}
	}

	err := s.backend.CreateContainer(ctx, domain.ContainerCreateArgs{
		ApplicationID: app.ID,
		ImageName:     builder.GetImageName(s.imageRegistry, s.imageNamePrefix, app.ID),
		ImageTag:      build.ID,
		HTTPProxy:     httpProxy,
		Recreate:      true,
	})
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}
	return nil
}
