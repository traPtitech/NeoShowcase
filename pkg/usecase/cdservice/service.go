package cdservice

import (
	"context"
	"log/slog"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/observability"
	"github.com/traPtitech/neoshowcase/pkg/util/discovery"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/loop"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
	"github.com/traPtitech/neoshowcase/pkg/util/scutil"
)

type service struct {
	cluster   *discovery.Cluster
	appRepo   domain.ApplicationRepository
	buildRepo domain.BuildRepository
	envRepo   domain.EnvironmentRepository
	backend   domain.Backend
	builder   domain.ControllerBuilderService
	deployer  *AppDeployHelper
	mutator   *ContainerStateMutator

	localClient domain.ControllerServiceClient

	doLocalStartBuild   func()
	doClusterStartBuild func()
	doLocalSyncDeploy   func()
	doClusterSyncDeploy func()
	run                 func()
	runOnce             sync.Once
	close               func()
	closeOnce           sync.Once
}

func NewService(
	cluster *discovery.Cluster,
	port grpc.ControllerPort,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	envRepo domain.EnvironmentRepository,
	backend domain.Backend,
	builder domain.ControllerBuilderService,
	deployer *AppDeployHelper,
	mutator *ContainerStateMutator,
	metrics *observability.ControllerMetrics,
) (domain.CDService, error) {
	cd := &service{
		cluster:   cluster,
		appRepo:   appRepo,
		buildRepo: buildRepo,
		envRepo:   envRepo,
		backend:   backend,
		builder:   builder,
		deployer:  deployer,
		mutator:   mutator,

		// そのまま現在のgRPCレイヤーを参照すると循環参照になってしまうため、ネットワークから呼び出しする
		// https://github.com/traPtitech/NeoShowcase/pull/1071#discussion_r2193711878
		localClient: grpc.NewControllerServiceClient(grpc.ControllerServiceClientConfig{
			URL: "http://127.0.0.1:" + strconv.Itoa(int(port)),
		}),
	}

	ctx, cancel := context.WithCancel(context.Background())

	doLocalStartBuild := scutil.NewCoalescer(func(ctx context.Context) error {
		start := time.Now()
		if err := cd.startBuildsLocal(ctx); err != nil {
			slog.Error("failed to start builds", "error", err)
			return nil
		}
		slog.Info("Started builds", "duration", time.Since(start))
		time.Sleep(1 * time.Second) // 1 second throttle between build checks to account for quick succession of repo checks
		return nil
	})
	cd.doLocalStartBuild = func() {
		_ = doLocalStartBuild.Do(context.Background())
	}
	cd.doClusterStartBuild = func() {
		err := cd.localClient.StartBuild(context.Background())
		if err != nil {
			slog.Error("failed to broadcast StartBuild", "error", err)
		}
	}

	doLocalSyncDeploy := scutil.NewCoalescer(func(ctx context.Context) error {
		start := time.Now()
		if err := cd.syncDeployments(ctx); err != nil {
			slog.Error("failed to sync deployments", "error", err)
			return nil
		}
		elapsed := time.Since(start)
		slog.Info("Synced deployments", "duration", elapsed)
		metrics.ObserveDeployDuration(elapsed)
		return nil
	})
	cd.doLocalSyncDeploy = func() {
		_ = doLocalSyncDeploy.Do(context.Background())
	}
	cd.doClusterSyncDeploy = func() {
		err := cd.localClient.SyncDeployments(context.Background())
		if err != nil {
			slog.Error("failed to broadcast SyncDeployments", "error", err)
		}
	}

	doDetectBuildCrash := func(ctx context.Context) {
		start := time.Now()
		if err := cd.detectBuildCrash(ctx); err != nil {
			slog.Error("failed to detect build crash", "error", err)
		}
		slog.Debug("Build crash detection complete", "duration", time.Since(start))
	}

	cd.run = func() {
		go func() {
			sub, _ := builder.ListenBuilderIdle()
			for range sub {
				go cd.doLocalStartBuild()
			}
		}()
		go func() {
			sub, _ := builder.ListenBuildSettled()
			for range sub {
				go cd.doClusterSyncDeploy()
			}
		}()
		go loop.Loop(ctx, func(ctx context.Context) {
			_ = doLocalSyncDeploy.Do(ctx)
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
			slog.Error("failed to kickoff build for app", "app_id", appID, "error", err)
			return
		}
		go cd.doClusterStartBuild()
	}()
}

func (cd *service) StartBuildLocal() {
	go cd.doLocalStartBuild()
}

func (cd *service) SyncDeploymentsLocal() {
	go cd.doLocalSyncDeploy()
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
	if !app.Running {
		// Skip if app not running
		return nil
	}
	if app.Commit == domain.EmptyCommit {
		// Skip if still fetching commit
		return nil
	}

	env, err := cd.envRepo.GetEnv(ctx, domain.GetEnvCondition{
		ApplicationID: optional.From(appID),
	})
	if err != nil {
		return err
	}

	// Check if already queued
	builds, err := cd.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{
		ApplicationID: optional.From(appID),
		Commit:        optional.From(app.Commit),
		ConfigHash:    optional.From(app.Config.Hash(env)),
		// Do not count retriable build as 'exists' - enqueue a new build if only retriable builds exist
		Retriable: optional.From(false),
	})
	if err != nil {
		return err
	}
	if len(builds) > 0 {
		// Already queued for the commit / config
		return nil
	}

	// Cancel any previously queued builds
	_, err = cd.buildRepo.UpdateBuild(ctx, domain.GetBuildCondition{
		ApplicationID: optional.From(appID),
		Status:        optional.From(domain.BuildStatusQueued),
	}, domain.UpdateBuildArgs{
		Status: optional.From(domain.BuildStatusCanceled),
	})
	if err != nil {
		return err
	}
	// Queue new build
	return cd.buildRepo.CreateBuild(ctx, domain.NewBuild(app, env))
}

func (cd *service) startBuildsLocal(ctx context.Context) error {
	builds, err := cd.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{Status: optional.From(domain.BuildStatusQueued)})
	if err != nil {
		return err
	}
	slices.SortFunc(builds, ds.LessFunc(func(a *domain.Build) int64 { return a.QueuedAt.UnixNano() }))
	buildIDs := ds.Map(builds, func(build *domain.Build) string { return build.ID })
	cd.builder.StartBuilds(ctx, buildIDs)
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
		_, err = cd.buildRepo.UpdateBuild(ctx, domain.GetBuildCondition{
			ID:     optional.From(build.ID),
			Status: optional.From(domain.BuildStatusBuilding),
		}, domain.UpdateBuildArgs{
			Status: optional.From(domain.BuildStatusFailed),
		})
		if err != nil {
			slog.ErrorContext(ctx, "failed to mark crashed build as errored", "error", err)
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
	// Shard by app ID
	apps = lo.Filter(apps, func(app *domain.Application, _ int) bool {
		return cd.cluster.IsAssigned(app.ID)
	})

	commits := lo.Map(apps, func(app *domain.Application, _ int) string { return app.Commit })
	builds, err := cd.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{
		CommitIn: optional.From(commits),
		Status:   optional.From(domain.BuildStatusSucceeded),
	})
	if err != nil {
		return err
	}
	// get last succeeded builds for each app
	slices.SortFunc(builds, ds.LessFunc(func(b *domain.Build) int64 { return b.QueuedAt.UnixNano() }))
	buildsMap := lo.SliceToMap(builds, func(b *domain.Build) (string, *domain.Build) { return b.ApplicationID, b })
	for _, app := range apps {
		nextBuild, ok := buildsMap[app.ID]
		if !ok {
			continue
		}
		toUpdate := func() bool {
			if app.CurrentBuild == nextBuild.ID {
				return false
			}
			if app.CurrentBuild == "" {
				return true
			}
			beforeBuild, err := cd.buildRepo.GetBuild(ctx, app.CurrentBuild)
			if err != nil {
				slog.Warn("failed to retrieve build", "build_id", app.CurrentBuild)
				return false
			}
			return beforeBuild.QueuedAt.Before(nextBuild.QueuedAt)
		}()
		if toUpdate {
			err = cd.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{
				CurrentBuild: optional.From(nextBuild.ID),
				UpdatedAt:    optional.From(nextBuild.FinishedAt.V),
			})
			if err != nil {
				return errors.Wrap(err, "failed to sync application commit")
			}
		}
	}
	return nil
}

func (cd *service) syncDeployments(ctx context.Context) error {
	// Sync app fields from build result in an idempotent way
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
