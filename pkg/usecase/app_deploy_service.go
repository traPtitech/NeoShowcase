package usecase

import (
	"context"
	"sync"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type AppDeployService interface {
	Synchronize(ctx context.Context, restartIDs []string) error
	SynchronizeSS(ctx context.Context) error
}

// appDeployService アプリのデプロイ操作を行う
type appDeployService struct {
	syncLock sync.Mutex

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
		bus:       bus,
		backend:   backend,
		appRepo:   appRepo,
		buildRepo: buildRepo,
		envRepo:   envRepo,
		component: component,
		image:     imageConfig,
	}
}

func (s *appDeployService) getEnv(ctx context.Context, apps []*domain.Application) (map[string]map[string]string, error) {
	appIDs := lo.Map(apps, func(app *domain.Application, i int) string { return app.ID })
	envs, err := s.envRepo.GetEnv(ctx, domain.GetEnvCondition{ApplicationIDIn: optional.From(appIDs)})
	if err != nil {
		return nil, err
	}
	ret := make(map[string]map[string]string, len(appIDs))
	for _, env := range envs {
		if _, ok := ret[env.ApplicationID]; !ok {
			ret[env.ApplicationID] = make(map[string]string)
		}
		ret[env.ApplicationID][env.Key] = env.Value
	}
	return ret, nil
}

func (s *appDeployService) Synchronize(ctx context.Context, restartIDs []string) error {
	s.syncLock.Lock()
	defer s.syncLock.Unlock()

	// Get all 'running' state applications
	apps, err := s.appRepo.GetApplications(ctx, domain.GetApplicationCondition{
		Running: optional.From(true),
	})
	if err != nil {
		return err
	}

	// Calculate deploy-able applications
	commits := lo.SliceToMap(apps, func(app *domain.Application) (string, struct{}) { return app.CurrentCommit, struct{}{} })
	builds, err := s.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{
		CommitIn: optional.From(lo.Keys(commits)),
		Status:   optional.From(domain.BuildStatusSucceeded),
	})
	if err != nil {
		return err
	}
	buildExists := lo.SliceToMap(builds, func(b *domain.Build) (string, bool) { return b.Commit, true })
	syncableApps := lo.Filter(apps, func(app *domain.Application, i int) bool { return buildExists[app.WantCommit] })

	// Deploy
	envs, err := s.getEnv(ctx, syncableApps)
	if err != nil {
		return err
	}
	desiredStates := lo.Map(syncableApps, func(app *domain.Application, i int) *domain.AppDesiredState {
		return &domain.AppDesiredState{
			App:       app,
			ImageName: s.image.FullImageName(app.ID),
			ImageTag:  app.CurrentCommit,
			Envs:      envs[app.ID],
			Restart:   lo.Contains(restartIDs, app.ID),
		}
	})
	return s.backend.Synchronize(ctx, desiredStates)
}

func (s *appDeployService) SynchronizeSS(ctx context.Context) error {
	s.component.BroadcastSSGen(&pb.SSGenRequest{Type: pb.SSGenRequest_RELOAD})
	if err := s.backend.SynchronizeSSIngress(ctx); err != nil {
		return errors.Wrap(err, "failed to reload static site ingress")
	}
	return nil
}
