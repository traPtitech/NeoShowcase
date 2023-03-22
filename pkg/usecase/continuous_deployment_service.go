package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
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
	envRepo   domain.EnvironmentRepository
	deployer  AppDeployService
	builder   AppBuildService

	run       func()
	runOnce   sync.Once
	close     func()
	closeOnce sync.Once
}

func NewContinuousDeploymentService(
	bus domain.Bus,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	envRepo domain.EnvironmentRepository,
	deployer AppDeployService,
	builder AppBuildService,
) ContinuousDeploymentService {
	cd := &continuousDeploymentService{
		bus:       bus,
		appRepo:   appRepo,
		buildRepo: buildRepo,
		envRepo:   envRepo,
		deployer:  deployer,
		builder:   builder,
	}

	syncBuildCloser := make(chan struct{})
	syncDeployCloser := make(chan struct{})
	cd.run = func() {
		go cd.syncBuildLoop(syncBuildCloser)
		go cd.syncDeployLoop(syncDeployCloser)
	}
	cd.close = func() {
		close(syncBuildCloser)
		close(syncDeployCloser)
	}

	return cd
}

func (cd *continuousDeploymentService) Run() {
	cd.runOnce.Do(cd.run)
}

func (cd *continuousDeploymentService) Stop(_ context.Context) error {
	cd.closeOnce.Do(cd.close)
	return nil
}

func (cd *continuousDeploymentService) syncBuildLoop(closer <-chan struct{}) {
	sub := cd.bus.Subscribe(event.FetcherFetchComplete, event.CDServiceSyncBuildRequest)
	defer sub.Unsubscribe()

	doSync := func() {
		start := time.Now()
		if err := cd.kickoffBuilds(); err != nil {
			log.WithError(err).Error("failed to kickoff builds")
			return
		}
		log.Infof("Synced builds in %v", time.Since(start))
	}

	for {
		select {
		case <-sub.Chan():
			doSync()
		case <-closer:
			return
		}
	}
}

func (cd *continuousDeploymentService) syncDeployLoop(closer <-chan struct{}) {
	sub := cd.bus.Subscribe(event.BuilderBuildSucceeded)
	defer sub.Unsubscribe()
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	doSync := func() {
		start := time.Now()
		if err := cd.syncDeployments(); err != nil {
			log.WithError(err).Error("failed to sync deployments")
			return
		}
		log.Infof("Synced deployments in %v", time.Since(start))
	}

	doSync()

	for {
		select {
		case <-ticker.C:
			doSync()
		case <-sub.Chan():
			doSync()
		case <-closer:
			return
		}
	}
}

func (cd *continuousDeploymentService) kickoffBuilds() error {
	ctx := context.Background()
	applications, err := cd.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
	if err != nil {
		return err
	}
	commits := lo.Map(applications, func(app *domain.Application, i int) string { return app.WantCommit })
	builds, err := cd.buildRepo.GetBuildsInCommit(ctx, commits)
	if err != nil {
		return err
	}
	builds = lo.Filter(builds, func(b *domain.Build, i int) bool { return !(b.Status == builder.BuildStatusFailed && b.Retriable) })
	buildExistsForCommit := lo.SliceToMap(builds, func(b *domain.Build) (string, bool) { return b.Commit, true })
	for _, app := range applications {
		if buildExistsForCommit[app.WantCommit] {
			continue
		}
		if app.WantCommit == domain.EmptyCommit {
			continue
		}
		if app.State == domain.ApplicationStateIdle {
			continue
		}
		_, err := cd.builder.QueueBuild(ctx, app, app.WantCommit)
		if err != nil {
			log.WithError(err).WithField("application", app.ID).WithField("commit", app.WantCommit).Error("failed to queue build")
			continue // Continue even if each app errors
		}
	}
	return nil
}

func (cd *continuousDeploymentService) syncDeployments() error {
	ctx := context.Background()

	// Get out-of-sync and non-idle applications
	applications, err := cd.appRepo.GetApplications(ctx, domain.GetApplicationCondition{InSync: optional.From(false)})
	if err != nil {
		return err
	}
	applications = lo.Filter(applications, func(app *domain.Application, i int) bool {
		return app.State != domain.ApplicationStateIdle && app.State != domain.ApplicationStateErrored
	})

	for _, app := range applications {
		_ = cd.deployer.Synchronize(app.ID, false)
	}
	return nil
}
