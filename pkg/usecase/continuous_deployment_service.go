package usecase

import (
	"context"
	"fmt"
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

type ContinuousDeploymentService interface {
	Run()
	RegisterBuilds()
	StartBuilds()
	SyncDeployments()
	Stop(ctx context.Context) error
}

type continuousDeploymentService struct {
	appRepo   domain.ApplicationRepository
	buildRepo domain.BuildRepository
	backend   domain.Backend
	builder   *AppBuildHelper
	deployer  *AppDeployHelper
	mutator   *ContainerStateMutator

	doRegisterBuild func()
	doStartBuild    func()
	doSyncDeploy    func()
	run             func()
	runOnce         sync.Once
	close           func()
	closeOnce       sync.Once
}

func NewContinuousDeploymentService(
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	backend domain.Backend,
	builder *AppBuildHelper,
	builderSvc domain.ControllerBuilderService,
	deployer *AppDeployHelper,
	mutator *ContainerStateMutator,
) (ContinuousDeploymentService, error) {
	cd := &continuousDeploymentService{
		appRepo:   appRepo,
		buildRepo: buildRepo,
		backend:   backend,
		builder:   builder,
		deployer:  deployer,
		mutator:   mutator,
	}

	ctx, cancel := context.WithCancel(context.Background())

	doRegisterBuild := coalesce.NewCoalescer(func(ctx context.Context) error {
		start := time.Now()
		if err := cd.registerBuilds(ctx); err != nil {
			log.Errorf("failed to kickoff builds: %+v", err)
			return nil
		}
		log.Infof("Synced builds in %v", time.Since(start))

		go cd.doStartBuild()
		return nil
	})
	cd.doRegisterBuild = func() {
		_ = doRegisterBuild.Do(context.Background())
	}

	doStartBuild := coalesce.NewCoalescer(func(ctx context.Context) error {
		start := time.Now()
		if err := cd.startBuilds(ctx); err != nil {
			log.Errorf("failed to start builds: %+v", err)
			return nil
		}
		log.Infof("Started builds in %v", time.Since(start))
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
			sub, _ := builderSvc.ListenBuilderIdle()
			for range sub {
				go cd.doStartBuild()
			}
		}()
		go func() {
			sub, _ := builderSvc.ListenBuildSettled()
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

func (cd *continuousDeploymentService) Run() {
	cd.runOnce.Do(cd.run)
}

func (cd *continuousDeploymentService) RegisterBuilds() {
	go cd.doRegisterBuild()
}

func (cd *continuousDeploymentService) StartBuilds() {
	go cd.doStartBuild()
}

func (cd *continuousDeploymentService) SyncDeployments() {
	go cd.doSyncDeploy()
}

func (cd *continuousDeploymentService) Stop(_ context.Context) error {
	cd.closeOnce.Do(cd.close)
	return nil
}

func (cd *continuousDeploymentService) registerBuilds(ctx context.Context) error {
	applications, err := cd.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
	if err != nil {
		return err
	}
	commits := ds.Map(applications, func(app *domain.Application) string { return app.WantCommit })
	builds, err := cd.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{CommitIn: optional.From(commits)})
	if err != nil {
		return err
	}

	buildExists := make(map[string]bool, len(applications)) // app id + commit -> bool
	for _, build := range builds {
		if build.Retriable {
			continue // Do not count retriable build as 'exists'
		}
		buildExists[build.ApplicationID+build.Commit] = true
	}
	for _, app := range applications {
		if buildExists[app.ID+app.WantCommit] {
			continue
		}
		if app.WantCommit == domain.EmptyCommit {
			continue
		}
		if !app.Running {
			continue
		}
		build := domain.NewBuild(app.ID, app.WantCommit)
		err = cd.buildRepo.CreateBuild(ctx, build)
		if err != nil {
			return errors.Wrap(err, "failed to create build")
		}
	}
	return nil
}

func (cd *continuousDeploymentService) startBuilds(ctx context.Context) error {
	builds, err := cd.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{Status: optional.From(domain.BuildStatusQueued)})
	if err != nil {
		return err
	}
	appIDs := ds.Map(builds, func(b *domain.Build) string { return b.ApplicationID })
	apps, err := cd.appRepo.GetApplications(ctx, domain.GetApplicationCondition{IDIn: optional.From(appIDs)})
	if err != nil {
		return err
	}
	appByID := lo.SliceToMap(apps, func(app *domain.Application) (string, *domain.Application) { return app.ID, app })
	for _, build := range builds {
		app, ok := appByID[build.ApplicationID]
		if !ok {
			return fmt.Errorf("app %v not found", build.ApplicationID)
		}
		cd.builder.tryStartBuild(app, build)
	}
	return nil
}

func (cd *continuousDeploymentService) detectBuildCrash(ctx context.Context) error {
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

func (cd *continuousDeploymentService) _syncAppFields(ctx context.Context) error {
	// Get all running applications
	apps, err := cd.appRepo.GetApplications(ctx, domain.GetApplicationCondition{
		Running: optional.From(true),
	})
	if err != nil {
		return err
	}
	commits := lo.SliceToMap(apps, func(app *domain.Application) (string, struct{}) { return app.WantCommit, struct{}{} })
	builds, err := cd.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{
		CommitIn: optional.From(lo.Keys(commits)),
		Status:   optional.From(domain.BuildStatusSucceeded),
	})
	if err != nil {
		return err
	}

	slices.SortFunc(builds, func(b1, b2 *domain.Build) bool { return b1.FinishedAt.V.Before(b2.FinishedAt.V) })
	latestBuild := lo.SliceToMap(builds, func(b *domain.Build) (string, *domain.Build) { return b.ApplicationID + b.Commit, b })

	for _, app := range apps {
		// Check if build has succeeded, if so sync 'current_commit' and 'updated_at' fields for restarting app
		if b, ok := latestBuild[app.ID+app.WantCommit]; ok {
			isSynced := app.CurrentCommit == app.WantCommit && !app.UpdatedAt.Before(b.FinishedAt.V)
			if isSynced {
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

func (cd *continuousDeploymentService) syncDeployments(ctx context.Context) error {
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
