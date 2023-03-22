package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type AppDeployService interface {
	Synchronize(appID string, restart bool) (started bool)
	Stop(appID string) (started bool)
}

// appDeployService アプリのデプロイ操作を行う
// 正しく状態をロックするため、アプリの State の操作はここでしか行わないようにする
type appDeployService struct {
	deployLock *ds.Mutex[string]

	bus       domain.Bus
	backend   domain.Backend
	appRepo   domain.ApplicationRepository
	buildRepo domain.BuildRepository
	envRepo   domain.EnvironmentRepository
	ss        pb.StaticSiteServiceClient

	imageRegistry   string
	imageNamePrefix string
}

func NewAppDeployService(
	bus domain.Bus,
	backend domain.Backend,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	envRepo domain.EnvironmentRepository,
	ss pb.StaticSiteServiceClient,
	registry builder.DockerImageRegistryString,
	prefix builder.DockerImageNamePrefixString,
) AppDeployService {
	return &appDeployService{
		deployLock:      ds.NewMutex[string](),
		bus:             bus,
		backend:         backend,
		appRepo:         appRepo,
		buildRepo:       buildRepo,
		envRepo:         envRepo,
		ss:              ss,
		imageRegistry:   string(registry),
		imageNamePrefix: string(prefix),
	}
}

func (s *appDeployService) Synchronize(appID string, restart bool) (started bool) {
	if ok := s.deployLock.TryLock(appID); !ok {
		return false
	}

	go func() {
		defer s.deployLock.Unlock(appID)

		ctx := context.Background()
		err := s.synchronize(ctx, appID, restart)
		if err != nil {
			log.WithError(err).WithField("application", appID).Error("failed to synchronize app")
			err = s.appRepo.UpdateApplication(ctx, appID, domain.UpdateApplicationArgs{State: optional.From(domain.ApplicationStateErrored)})
			if err != nil {
				log.WithError(err).Error("failed to update application state")
			}
		}
	}()

	return true
}

func (s *appDeployService) synchronize(ctx context.Context, appID string, restart bool) error {
	start := time.Now()

	app, err := s.appRepo.GetApplication(ctx, appID)
	if err != nil {
		return err
	}

	// Mark application as started if idle
	if app.State == domain.ApplicationStateIdle {
		err = s.appRepo.UpdateApplication(ctx, app.ID, domain.UpdateApplicationArgs{State: optional.From(domain.ApplicationStateDeploying)})
		if err != nil {
			return err
		}
	}

	build, err := s.getSuccessBuild(ctx, app)
	if err == ErrNotFound {
		s.bus.Publish(event.CDServiceSyncBuildRequest, nil)
		return nil
	}
	if err != nil {
		return err
	}

	doDeploy := restart || (!restart && app.WantCommit != app.CurrentCommit)

	if doDeploy && app.BuildType == builder.BuildTypeRuntime {
		err = s.appRepo.UpdateApplication(ctx, app.ID, domain.UpdateApplicationArgs{State: optional.From(domain.ApplicationStateDeploying)})
		if err != nil {
			return err
		}

		err = s.recreateContainer(ctx, app, build)
		if err != nil {
			return err
		}
	}

	err = s.appRepo.UpdateApplication(ctx, app.ID, domain.UpdateApplicationArgs{
		State:         optional.From(domain.ApplicationStateRunning),
		CurrentCommit: optional.From(build.Commit),
	})
	if err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}

	if doDeploy && app.BuildType == builder.BuildTypeStatic {
		if _, err = s.ss.Reload(ctx, &emptypb.Empty{}); err != nil {
			return fmt.Errorf("failed to reload static site server: %w", err)
		}
	}

	log.WithField("application", app.ID).Infof("app deploy suceeded in %v", time.Since(start))
	return nil
}

func (s *appDeployService) getSuccessBuild(ctx context.Context, app *domain.Application) (*domain.Build, error) {
	builds, err := s.buildRepo.GetBuildsInCommit(ctx, []string{app.WantCommit})
	if err != nil {
		return nil, err
	}
	builds = lo.Filter(builds, func(build *domain.Build, i int) bool { return build.Status == builder.BuildStatusSucceeded })
	slices.SortFunc(builds, func(a, b *domain.Build) bool { return a.StartedAt.After(b.StartedAt) })
	if len(builds) == 0 {
		return nil, ErrNotFound
	}
	return builds[0], nil
}

func (s *appDeployService) recreateContainer(ctx context.Context, app *domain.Application, build *domain.Build) error {
	envs, err := s.envRepo.GetEnv(ctx, app.ID)
	if err != nil {
		return err
	}
	err = s.backend.CreateContainer(ctx, app, domain.ContainerCreateArgs{
		ImageName: builder.GetImageName(s.imageRegistry, s.imageNamePrefix, app.ID),
		ImageTag:  build.ID,
		Recreate:  true,
		Envs:      lo.SliceToMap(envs, func(env *domain.Environment) (string, string) { return env.Key, env.Value }),
	})
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}
	return nil
}

func (s *appDeployService) Stop(appID string) (started bool) {
	if ok := s.deployLock.TryLock(appID); !ok {
		return false
	}

	go func() {
		defer s.deployLock.Unlock(appID)

		ctx := context.Background()
		err := s.stop(ctx, appID)
		if err != nil {
			log.WithError(err).WithField("application", appID).Error("failed to stop app")
			err = s.appRepo.UpdateApplication(ctx, appID, domain.UpdateApplicationArgs{State: optional.From(domain.ApplicationStateErrored)})
			if err != nil {
				log.WithError(err).Error("failed to update application state")
			}
		}
	}()

	return true
}

func (s *appDeployService) stop(ctx context.Context, appID string) error {
	app, err := s.appRepo.GetApplication(ctx, appID)
	if err != nil {
		return err
	}

	if app.BuildType == builder.BuildTypeRuntime {
		err = s.backend.DestroyContainer(ctx, app)
		if err != nil {
			return err
		}
	}

	err = s.appRepo.UpdateApplication(ctx, app.ID, domain.UpdateApplicationArgs{State: optional.From(domain.ApplicationStateIdle)})
	if err != nil {
		return fmt.Errorf("failed to update application state: %w", err)
	}

	if app.BuildType == builder.BuildTypeStatic {
		if _, err = s.ss.Reload(ctx, &emptypb.Empty{}); err != nil {
			return fmt.Errorf("failed to reload static site server: %w", err)
		}
	}

	return nil
}
