package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/heroku/docker-registry-client/registry"
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
	builder   AppBuildService
	deployer  AppDeployService
	image     builder.ImageConfig

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
	image builder.ImageConfig,
) (ContinuousDeploymentService, error) {
	cd := &continuousDeploymentService{
		bus:       bus,
		appRepo:   appRepo,
		buildRepo: buildRepo,
		envRepo:   envRepo,
		builder:   builder,
		deployer:  deployer,
		image:     image,
	}

	r, err := registry.New(image.Registry.Scheme+"://"+image.Registry.Addr, image.Registry.Username, image.Registry.Password)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	startBuilds := make(chan struct{})
	cd.run = func() {
		go cd.registerBuildLoop(ctx, startBuilds)
		go cd.startBuildLoop(ctx, startBuilds)
		go cd.syncDeployLoop(ctx)
		go cd.pruneImagesLoop(ctx, r)
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
	sub := cd.bus.Subscribe(event.BuilderBuildSettled)
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

func (cd *continuousDeploymentService) syncDeployLoop(ctx context.Context) {
	sub := cd.bus.Subscribe(event.BuilderBuildSettled)
	defer sub.Unsubscribe()
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	doSync := func() {
		start := time.Now()
		if err := cd.syncDeployments(ctx); err != nil {
			log.Errorf("failed to sync deployments: %+v", err)
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
		case <-ctx.Done():
			return
		}
	}
}

func (cd *continuousDeploymentService) pruneImagesLoop(ctx context.Context, r *registry.Registry) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	doPrune := func() {
		start := time.Now()
		err := cd.pruneImages(ctx, r)
		if err != nil {
			log.Errorf("failed to prune images: %+v", err)
			return
		}
		log.Infof("Pruned images in %v", time.Since(start))
	}

	for {
		select {
		case <-ticker.C:
			doPrune()
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
			log.Errorf("failed to mark crashed build as errored: %+v", err)
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

func (cd *continuousDeploymentService) startBuilds(ctx context.Context) error {
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

func (cd *continuousDeploymentService) syncDeployments(ctx context.Context) error {
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

func (cd *continuousDeploymentService) pruneImages(ctx context.Context, r *registry.Registry) error {
	applications, err := cd.appRepo.GetApplications(ctx, domain.GetApplicationCondition{BuildType: optional.From(builder.BuildTypeRuntime)})
	if err != nil {
		return err
	}
	commits := make(map[string]struct{}, 2*len(applications))
	for _, app := range applications {
		commits[app.WantCommit] = struct{}{}
		commits[app.CurrentCommit] = struct{}{}
	}

	for _, app := range applications {
		imageName := cd.image.ImageName(app.ID)
		tags, err := r.Tags(imageName)
		if err != nil {
			return errors.Wrap(err, "failed to get tags for image")
		}
		danglingTags := lo.Reject(tags, func(tag string, i int) bool { _, ok := commits[tag]; return ok })
		for _, tag := range danglingTags {
			digest, err := r.ManifestDigest(imageName, tag)
			if err != nil {
				return errors.Wrap(err, "failed to get manifest digest")
			}
			// NOTE: needs manual execution of "registry garbage-collect <config> --delete-untagged" in docker registry side
			// to actually delete the layers
			// https://docs.docker.com/registry/garbage-collection/
			err = r.DeleteManifest(imageName, digest)
			if err != nil {
				return errors.Wrap(err, "failed to delete manifest")
			}
		}
	}

	return nil
}
