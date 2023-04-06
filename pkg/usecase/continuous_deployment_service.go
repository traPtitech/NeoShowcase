package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/util/coalesce"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type ContinuousDeploymentService interface {
	Run()
	Stop(ctx context.Context) error
}

type continuousDeploymentService struct {
	bus       domain.Bus
	appRepo   domain.ApplicationRepository
	buildRepo domain.BuildRepository
	builder   *AppBuildHelper
	deployer  *AppDeployHelper

	run       func()
	runOnce   sync.Once
	close     func()
	closeOnce sync.Once
}

func NewContinuousDeploymentService(
	bus domain.Bus,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	builder *AppBuildHelper,
	deployer *AppDeployHelper,
) (ContinuousDeploymentService, error) {
	cd := &continuousDeploymentService{
		bus:       bus,
		appRepo:   appRepo,
		buildRepo: buildRepo,
		builder:   builder,
		deployer:  deployer,
	}

	ctx, cancel := context.WithCancel(context.Background())
	startBuilds := make(chan struct{})
	cd.run = func() {
		go cd.registerBuildLoop(ctx, startBuilds)
		go cd.startBuildLoop(ctx, startBuilds)
		go cd.detectBuildCrashLoop(ctx)
		go cd.syncDeployLoop(ctx)
	}
	cd.close = func() {
		cancel()
	}

	return cd, nil
}

func (cd *continuousDeploymentService) Run() {
	cd.runOnce.Do(cd.run)
}

func (cd *continuousDeploymentService) Stop(_ context.Context) error {
	cd.closeOnce.Do(cd.close)
	return nil
}

func (cd *continuousDeploymentService) registerBuildLoop(ctx context.Context, startBuilds chan<- struct{}) {
	sub := cd.bus.Subscribe(event.FetcherFetchComplete, event.CDServiceRegisterBuildRequest)
	defer sub.Unsubscribe()

	doSync := func() {
		start := time.Now()
		if err := cd.registerBuilds(ctx); err != nil {
			log.Errorf("failed to kickoff builds: %+v", err)
			return
		}
		select {
		case startBuilds <- struct{}{}:
		default:
		}
		log.Infof("Synced builds in %v", time.Since(start))
	}

	for {
		select {
		case <-sub.Chan():
			doSync()
		case <-ctx.Done():
			return
		}
	}
}

func (cd *continuousDeploymentService) startBuildLoop(ctx context.Context, syncer <-chan struct{}) {
	sub := cd.bus.Subscribe(event.BuilderConnected, event.BuilderBuildSettled)
	defer sub.Unsubscribe()

	doSync := func() {
		start := time.Now()
		if err := cd.startBuilds(ctx); err != nil {
			log.Errorf("failed to start builds: %+v", err)
			return
		}
		log.Infof("Started builds in %v", time.Since(start))
	}

	doSync()

	for {
		select {
		case <-syncer:
			doSync()
		case <-sub.Chan():
			doSync()
		case <-ctx.Done():
			return
		}
	}
}

func (cd *continuousDeploymentService) detectBuildCrashLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	doSync := func() {
		start := time.Now()
		if err := cd.detectBuildCrash(ctx); err != nil {
			log.Errorf("failed to detect build crash: %+v", err)
		}
		log.Debugf("Build crash detection complete in %v", time.Since(start))
	}

	doSync()

	for {
		select {
		case <-ticker.C:
			doSync()
		case <-ctx.Done():
			return
		}
	}
}

func (cd *continuousDeploymentService) syncDeployLoop(ctx context.Context) {
	sub := cd.bus.Subscribe(event.BuilderBuildSettled, event.CDServiceSyncDeployRequest)
	defer sub.Unsubscribe()
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	coalescer := coalesce.NewCoalescer(func(ctx context.Context) error {
		start := time.Now()
		if err := cd.syncDeployments(ctx); err != nil {
			log.Errorf("failed to sync deployments: %+v", err)
			return nil
		}
		log.Infof("Synced deployments in %v", time.Since(start))
		return nil
	})

	_ = coalescer.Do(ctx)

	for {
		select {
		case <-ticker.C:
			_ = coalescer.Do(ctx)
		case <-sub.Chan():
			_ = coalescer.Do(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (cd *continuousDeploymentService) registerBuilds(ctx context.Context) error {
	applications, err := cd.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
	if err != nil {
		return err
	}
	commits := lo.Map(applications, func(app *domain.Application, i int) string { return app.WantCommit })
	builds, err := cd.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{CommitIn: optional.From(commits), Retriable: optional.From(false)})
	if err != nil {
		return err
	}

	buildExistsForCommit := lo.SliceToMap(builds, func(b *domain.Build) (string, bool) { return b.Commit, true })
	for _, app := range applications {
		if buildExistsForCommit[app.WantCommit] {
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
	appIDs := lo.Map(builds, func(b *domain.Build, i int) string { return b.ApplicationID })
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
	crashed := lo.Filter(builds, func(build *domain.Build, i int) bool {
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

func (cd *continuousDeploymentService) syncDeployments(ctx context.Context) error {
	// Get out-of-sync applications
	apps, err := cd.appRepo.GetApplications(ctx, domain.GetApplicationCondition{
		Running: optional.From(true),
		InSync:  optional.From(false),
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
	buildExists := lo.SliceToMap(builds, func(b *domain.Build) (string, bool) { return b.Commit, true })

	// Check if build has succeeded, and if so save as synced
	for _, app := range apps {
		if buildExists[app.WantCommit] {
			err = cd.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{CurrentCommit: optional.From(app.WantCommit)})
			if err != nil {
				return errors.Wrap(err, "failed to sync application commit")
			}
		}
	}

	// Synchronize
	err = cd.deployer.synchronize(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to synchronize app deployments")
	}
	return cd.deployer.synchronizeSS(ctx)
}
