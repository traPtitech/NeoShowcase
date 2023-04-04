package usecase

import (
	"context"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type AppDeployService interface {
	Synchronize(appID string, restart bool) (started bool)
	SynchronizeSS(ctx context.Context) error
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
	component domain.ComponentService

	image builder.ImageConfig
}

func NewAppDeployService(
	bus domain.Bus,
	backend domain.Backend,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	envRepo domain.EnvironmentRepository,
	component domain.ComponentService,
	imageConfig builder.ImageConfig,
) AppDeployService {
	return &appDeployService{
		deployLock: ds.NewMutex[string](),
		bus:        bus,
		backend:    backend,
		appRepo:    appRepo,
		buildRepo:  buildRepo,
		envRepo:    envRepo,
		component:  component,
		image:      imageConfig,
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
			log.Errorf("failed to synchronize app: %+v", err)
			err = s.appRepo.UpdateApplication(ctx, appID, &domain.UpdateApplicationArgs{State: optional.From(domain.ApplicationStateErrored)})
			if err != nil {
				log.Errorf("failed to update application state: %+v", err)
			}
		}
	}()

	return true
}

func (s *appDeployService) SynchronizeSS(ctx context.Context) error {
	s.component.BroadcastSSGen(&pb.SSGenRequest{Type: pb.SSGenRequest_RELOAD})
	if err := s.backend.ReloadSSIngress(ctx); err != nil {
		return errors.Wrap(err, "failed to reload static site ingress")
	}
	return nil
}

func (s *appDeployService) synchronize(ctx context.Context, appID string, restart bool) error {
	start := time.Now()

	app, err := s.appRepo.GetApplication(ctx, appID)
	if err != nil {
		return err
	}

	// Mark application as started if idle
	if app.State == domain.ApplicationStateIdle {
		err = s.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{State: optional.From(domain.ApplicationStateDeploying)})
		if err != nil {
			return err
		}
	}

	build, err := s.getSuccessBuild(ctx, app)
	if err != nil {
		return err
	}
	if build == nil {
		s.bus.Publish(event.CDServiceRegisterBuildRequest, nil)
		return nil
	}

	doDeploy := restart || (!restart && app.WantCommit != app.CurrentCommit)

	if doDeploy && app.BuildType == domain.BuildTypeRuntime {
		err = s.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{State: optional.From(domain.ApplicationStateDeploying)})
		if err != nil {
			return err
		}

		err = s.recreateContainer(ctx, app, build)
		if err != nil {
			return err
		}
	}

	err = s.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{
		State:         optional.From(domain.ApplicationStateRunning),
		CurrentCommit: optional.From(build.Commit),
	})
	if err != nil {
		return errors.Wrap(err, "failed to update application")
	}

	if doDeploy && app.BuildType == domain.BuildTypeStatic {
		if err = s.SynchronizeSS(ctx); err != nil {
			return err
		}
	}

	log.WithField("application", app.ID).Infof("app deploy suceeded in %v", time.Since(start))
	return nil
}

func (s *appDeployService) getSuccessBuild(ctx context.Context, app *domain.Application) (*domain.Build, error) {
	builds, err := s.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{Commit: optional.From(app.WantCommit), Status: optional.From(domain.BuildStatusSucceeded)})
	if err != nil {
		return nil, err
	}
	slices.SortFunc(builds, func(a, b *domain.Build) bool { return a.StartedAt.ValueOrZero().After(b.StartedAt.ValueOrZero()) })
	if len(builds) == 0 {
		return nil, nil
	}
	return builds[0], nil
}

func (s *appDeployService) recreateContainer(ctx context.Context, app *domain.Application, build *domain.Build) error {
	envs, err := s.envRepo.GetEnv(ctx, domain.GetEnvCondition{ApplicationID: optional.From(app.ID)})
	if err != nil {
		return err
	}
	err = s.backend.CreateContainer(ctx, app, domain.ContainerCreateArgs{
		ImageName: s.image.FullImageName(app.ID),
		ImageTag:  build.Commit,
		Envs:      lo.SliceToMap(envs, func(env *domain.Environment) (string, string) { return env.Key, env.Value }),
	})
	if err != nil {
		return errors.Wrap(err, "failed to create container")
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
			log.Errorf("failed to stop app: %+v", err)
			err = s.appRepo.UpdateApplication(ctx, appID, &domain.UpdateApplicationArgs{State: optional.From(domain.ApplicationStateErrored)})
			if err != nil {
				log.Errorf("failed to update application state: %+v", err)
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

	if app.BuildType == domain.BuildTypeRuntime {
		err = s.backend.DestroyContainer(ctx, app)
		if err != nil {
			return err
		}
	}

	err = s.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{State: optional.From(domain.ApplicationStateIdle)})
	if err != nil {
		return errors.Wrap(err, "failed to update application state")
	}

	if app.BuildType == domain.BuildTypeStatic {
		if err = s.SynchronizeSS(ctx); err != nil {
			return err
		}
	}

	return nil
}
