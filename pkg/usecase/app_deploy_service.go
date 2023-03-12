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
		// Rollback state
		err = s.appRepo.UpdateApplication(ctx, app.ID, repository.UpdateApplicationArgs{State: optional.From(domain.ApplicationStateErrored)})
		if err != nil {
			log.WithError(err).Error("failed to update application state back")
		}
		return
	}
	log.WithField("application", app.ID).WithField("build", build.ID).Infof("deploy succeeded in %v", time.Since(start))
}

func (s *appDeployService) deploy(ctx context.Context, app *domain.Application, build *domain.Build) error {
	err := s.recreate(ctx, app, build)
	if err != nil {
		return err
	}

	err = s.appRepo.UpdateApplication(ctx, app.ID, repository.UpdateApplicationArgs{
		State:         optional.From(domain.ApplicationStateRunning),
		CurrentCommit: optional.From(build.Commit),
	})
	if err != nil {
		return fmt.Errorf("failed to update application")
	}
	return nil
}

func (s *appDeployService) recreate(ctx context.Context, app *domain.Application, build *domain.Build) error {
	website, err := s.appRepo.GetWebsite(ctx, app.ID)
	if err != nil && err != repository.ErrNotFound {
		return fmt.Errorf("failed to get website: %w", err)
	}

	// TODO Ingressの設定がここでする必要があるかどうか確認する
	switch app.BuildType {
	case builder.BuildTypeImage:
		var httpProxy *domain.ContainerHTTPProxy
		if website != nil {
			httpProxy = &domain.ContainerHTTPProxy{
				Domain: website.FQDN,
				Port:   website.Port,
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
	case builder.BuildTypeStatic:
		if _, err := s.ss.Reload(ctx, &pb.ReloadRequest{}); err != nil {
			return fmt.Errorf("failed to reload static site server: %w", err)
		}
	}
	return nil
}
