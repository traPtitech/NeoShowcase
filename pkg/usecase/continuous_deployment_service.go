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
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type ContinuousDeploymentService interface {
	Run()
	Stop(ctx context.Context) error
}

type continuousDeploymentService struct {
	bus             domain.Bus
	appRepo         domain.ApplicationRepository
	buildRepo       domain.BuildRepository
	envRepo         domain.EnvironmentRepository
	builder         AppBuildService
	deployer        AppDeployService
	imageRegistry   string
	imageNamePrefix string

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
	builder AppBuildService,
	deployer AppDeployService,
	registry builder.DockerImageRegistryString,
	prefix builder.DockerImageNamePrefixString,
) ContinuousDeploymentService {
	cd := &continuousDeploymentService{
		bus:             bus,
		appRepo:         appRepo,
		buildRepo:       buildRepo,
		envRepo:         envRepo,
		builder:         builder,
		deployer:        deployer,
		imageRegistry:   string(registry),
		imageNamePrefix: string(prefix),
	}

	registerBuildCloser := make(chan struct{})
	startBuilds := make(chan struct{})
	startBuildCloser := make(chan struct{})
	syncDeployCloser := make(chan struct{})
	cd.run = func() {
		go cd.registerBuildLoop(startBuilds, registerBuildCloser)
		go cd.startBuildLoop(startBuilds, startBuildCloser)
		go cd.syncDeployLoop(syncDeployCloser)
	}
	cd.close = func() {
		close(registerBuildCloser)
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

func (cd *continuousDeploymentService) registerBuildLoop(startBuilds chan<- struct{}, closer <-chan struct{}) {
	sub := cd.bus.Subscribe(event.FetcherFetchComplete, event.CDServiceRegisterBuildRequest)
	defer sub.Unsubscribe()

	doSync := func() {
		start := time.Now()
		if err := cd.registerBuilds(); err != nil {
			log.WithError(err).Error("failed to kickoff builds")
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
		case <-closer:
			return
		}
	}
}

func (cd *continuousDeploymentService) startBuildLoop(syncer <-chan struct{}, closer <-chan struct{}) {
	sub := cd.bus.Subscribe(event.BuilderBuildSettled)
	defer sub.Unsubscribe()

	doSync := func() {
		start := time.Now()
		if err := cd.startBuilds(); err != nil {
			log.WithError(err).Error("failed to start builds")
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
		case <-closer:
			return
		}
	}
}

func (cd *continuousDeploymentService) syncDeployLoop(closer <-chan struct{}) {
	sub := cd.bus.Subscribe(event.BuilderBuildSettled)
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

func (cd *continuousDeploymentService) registerBuilds() error {
	ctx := context.Background()
	applications, err := cd.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
	if err != nil {
		return err
	}
	commits := lo.Map(applications, func(app *domain.Application, i int) string { return app.WantCommit })
	builds, err := cd.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{CommitIn: optional.From(commits), Retriable: optional.From(false)})
	if err != nil {
		return err
	}

	// Detect builder crash and mark builds as errored
	const crashDetectThreshold = 60 * time.Second
	crashThreshold := time.Now().Add(-crashDetectThreshold)
	crashed := lo.Filter(builds, func(build *domain.Build, i int) bool {
		return build.Status == builder.BuildStatusBuilding && build.UpdatedAt.ValueOrZero().Before(crashThreshold)
	})
	for _, build := range crashed {
		err = cd.buildRepo.UpdateBuild(ctx, build.ID, domain.UpdateBuildArgs{
			FromStatus: optional.From(builder.BuildStatusBuilding),
			Status:     optional.From(builder.BuildStatusFailed),
		})
		if err != nil {
			log.WithError(err).Error("failed to mark crashed build as errored")
		}
	}

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
		build := domain.NewBuild(app.ID, app.WantCommit)
		err = cd.buildRepo.CreateBuild(ctx, build)
		if err != nil {
			return errors.Wrap(err, "failed to create build")
		}
	}
	return nil
}

func (cd *continuousDeploymentService) startBuilds() error {
	ctx := context.Background()
	builds, err := cd.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{Status: optional.From(builder.BuildStatusQueued)})
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
		cd.builder.TryStartBuild(app, build)
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
	return cd.deployer.SynchronizeSS(ctx)
}
