package cdservice

import (
	"context"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/coalesce"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/loop"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type Service interface {
	Run()
	RegisterBuild(appID string)
	SyncDeployments()
	Stop(ctx context.Context) error
}

type service struct {
	appRepo   domain.ApplicationRepository
	buildRepo domain.BuildRepository
	backend   domain.Backend
	builder   domain.ControllerBuilderService
	deployer  *AppDeployHelper
	mutator   *ContainerStateMutator

	doStartBuild func()
	doSyncDeploy func()
	run          func()
	runOnce      sync.Once
	close        func()
	closeOnce    sync.Once
}

func NewService(
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	backend domain.Backend,
	builder domain.ControllerBuilderService,
	deployer *AppDeployHelper,
	mutator *ContainerStateMutator,
) (Service, error) {
	cd := &service{
		appRepo:   appRepo,
		buildRepo: buildRepo,
		backend:   backend,
		builder:   builder,
		deployer:  deployer,
		mutator:   mutator,
	}

	ctx, cancel := context.WithCancel(context.Background())

	doStartBuild := coalesce.NewCoalescer(func(ctx context.Context) error {
		start := time.Now()
		if err := cd.startBuilds(ctx); err != nil {
			log.Errorf("failed to start builds: %+v", err)
			return nil
		}
		log.Infof("Started builds in %v", time.Since(start))
		time.Sleep(1 * time.Second) // 1 second throttle between build checks to account for quick succession of repo checks
		return nil
	})
	cd.doStartBuild = func() {
		_ = doStartBuild.Do(context.Background())
	}

	doSyncDeploy := coalesce.NewCoalescer(func(ctx context.Context) error {
		start := time.Now()
		if err := cd.syncDeployments(ctx); err != nil {
			log.Errorf("failed to sync deployments: %+v", err)
			return nil
		}
		log.Infof("Synced deployments in %v", time.Since(start))
		return nil
	})
	cd.doSyncDeploy = func() {
		_ = doSyncDeploy.Do(context.Background())
	}

	doDetectBuildCrash := func(ctx context.Context) {
		start := time.Now()
		if err := cd.detectBuildCrash(ctx); err != nil {
			log.Errorf("failed to detect build crash: %+v", err)
		}
		log.Debugf("Build crash detection complete in %v", time.Since(start))
	}

	cd.run = func() {
		go func() {
			sub, _ := builder.ListenBuilderIdle()
			for range sub {
				go cd.doStartBuild()
			}
		}()
		go func() {
			sub, _ := builder.ListenBuildSettled()
			for range sub {
				go cd.doSyncDeploy()
			}
		}()
		go loop.Loop(ctx, func(ctx context.Context) {
			_ = doSyncDeploy.Do(ctx)
		}, 3*time.Minute, true)
		go loop.Loop(ctx, doDetectBuildCrash, 1*time.Minute, true)
	}
	cd.close = cancel

	return cd, nil
}

func (cd *service) Run() {
	cd.runOnce.Do(cd.run)
}

func (cd *service) RegisterBuild(appID string) {
	go func() {
		if err := cd.registerBuild(context.Background(), appID); err != nil {
			log.Errorf("failed to kickoff build for app %v: %+v", appID, err)
			return
		}
		go cd.doStartBuild()
	}()
}

func (cd *service) SyncDeployments() {
	go cd.doSyncDeploy()
}

func (cd *service) Stop(_ context.Context) error {
	cd.closeOnce.Do(cd.close)
	return nil
}

func (cd *service) registerBuild(ctx context.Context, appID string) error {
	app, err := cd.appRepo.GetApplication(ctx, appID)
	if err != nil {
		return err
	}
	builds, err := cd.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{
		ApplicationID: optional.From(appID),
		Commit:        optional.From(app.WantCommit),
		// Do not count retriable build as 'exists' - enqueue a new build if only retriable builds exist
		Retriable: optional.From(false),
	})
	if err != nil {
		return err
	}
	if len(builds) > 0 {
		// Already queued for the commit
		return nil
	}
	if app.WantCommit == domain.EmptyCommit {
		// Skip if still fetching commit
		return nil
	}
	if !app.Running {
		// Skip if app not running
		return nil
	}
	build := domain.NewBuild(app.ID, app.WantCommit)
	return cd.buildRepo.CreateBuild(ctx, build)
}

func (cd *service) startBuilds(ctx context.Context) error {
	builds, err := cd.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{Status: optional.From(domain.BuildStatusQueued)})
	if err != nil {
		return err
	}
	slices.SortFunc(builds, func(a, b *domain.Build) bool { return a.QueuedAt.Before(b.QueuedAt) })
	buildIDs := ds.Map(builds, func(build *domain.Build) string { return build.ID })
	cd.builder.StartBuilds(buildIDs)
	return nil
}

func (cd *service) detectBuildCrash(ctx context.Context) error {
	const crashDetectThreshold = 60 * time.Second
	now := time.Now()

	builds, err := cd.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{Status: optional.From(domain.BuildStatusBuilding)})
	if err != nil {
		return errors.Wrap(err, "failed to get running builds")
	}
	crashed := lo.Filter(builds, func(build *domain.Build, _ int) bool {
		return now.Sub(build.UpdatedAt.ValueOrZero()) > crashDetectThreshold
	})
	for _, build := range crashed {
		err = cd.buildRepo.UpdateBuild(ctx, build.ID, domain.UpdateBuildArgs{
			FromStatus: optional.From(domain.BuildStatusBuilding),
			Status:     optional.From(domain.BuildStatusFailed),
		})
		if err != nil {
			log.Errorf("failed to mark crashed build as errored: %+v", err)
		}
	}

	return nil
}

func (cd *service) _syncAppFields(ctx context.Context) error {
	// Get all running applications
	apps, err := cd.appRepo.GetApplications(ctx, domain.GetApplicationCondition{
		Running: optional.From(true),
	})
	if err != nil {
		return err
	}

	builds, err := domain.GetSuccessBuilds(ctx, cd.buildRepo, apps)
	if err != nil {
		return err
	}
	for _, app := range apps {
		// Check if build has succeeded, if so sync 'current_commit' and 'updated_at' fields for restarting app
		if b, ok := builds[app.ID+app.WantCommit]; ok {
			toSync := app.CurrentCommit != app.WantCommit || app.UpdatedAt.Before(b.FinishedAt.V)
			if !toSync {
				continue
			}
			err = cd.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{
				CurrentCommit: optional.From(app.WantCommit),
				UpdatedAt:     optional.From(b.FinishedAt.V),
			})
			if err != nil {
				return errors.Wrap(err, "failed to sync application commit")
			}
		}
	}
	return nil
}

func (cd *service) syncDeployments(ctx context.Context) error {
	// Sync app fields from build result
	err := cd._syncAppFields(ctx)
	if err != nil {
		return err
	}

	// Synchronize
	err = cd.deployer.synchronize(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to synchronize deployments")
	}

	// Update container states
	err = cd.mutator.updateAll(ctx)
	if err != nil {
		return err
	}
	return nil
}
