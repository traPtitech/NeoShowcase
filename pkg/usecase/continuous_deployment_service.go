package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
)

type ContinuousDeploymentService interface {
	Run()
	Stop(ctx context.Context) error
}

type continuousDeploymentService struct {
	bus       domain.Bus
	appRepo   repository.ApplicationRepository
	buildRepo repository.BuildRepository
	envRepo   repository.EnvironmentRepository
	deployer  AppDeployService
	builder   AppBuildService

	run       func()
	runOnce   sync.Once
	close     func()
	closeOnce sync.Once
}

func NewContinuousDeploymentService(
	bus domain.Bus,
	appRepo repository.ApplicationRepository,
	buildRepo repository.BuildRepository,
	envRepo repository.EnvironmentRepository,
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
	sub := cd.bus.Subscribe(event.FetcherFetchComplete)
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
	applications, err := cd.appRepo.GetApplications(ctx)
	if err != nil {
		return err
	}
	commits := lo.Map(applications, func(app *domain.Application, i int) string { return app.WantCommit })
	builds, err := cd.buildRepo.GetBuildsInCommit(ctx, commits)
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
		_, err := cd.builder.QueueBuild(ctx, app, app.WantCommit)
		if err != nil {
			return fmt.Errorf("failed to queue build: %w", err)
		}
	}
	return nil
}

func (cd *continuousDeploymentService) syncDeployments() error {
	ctx := context.Background()
	applications, err := cd.appRepo.GetApplicationsOutOfSync(ctx)
	if err != nil {
		return err
	}
	applications = lo.Filter(applications, func(app *domain.Application, i int) bool { return app.State != domain.ApplicationStateDeploying })
	commits := lo.Map(applications, func(app *domain.Application, i int) string { return app.WantCommit })
	builds, err := cd.buildRepo.GetBuildsInCommit(ctx, commits)
	if err != nil {
		return err
	}

	// Last succeeded builds for each commit
	builds = lo.Filter(builds, func(build *domain.Build, i int) bool { return build.Status == builder.BuildStatusSucceeded })
	slices.SortFunc(builds, func(a, b *domain.Build) bool { return a.StartedAt.Before(b.StartedAt) })
	commitToBuild := lo.SliceToMap(builds, func(b *domain.Build) (string, *domain.Build) { return b.Commit, b })

	for _, app := range applications {
		if build, ok := commitToBuild[app.WantCommit]; ok {
			err := cd.deployer.QueueDeployment(ctx, build.ApplicationID, build.ID)
			if err != nil {
				return fmt.Errorf("failed to queue deployment: %w", err)
			}
		}
	}
	return nil
}
