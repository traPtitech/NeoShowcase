package cleaner

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/regclient/regclient"
	"github.com/regclient/regclient/types/ref"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/loop"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
	"github.com/traPtitech/neoshowcase/pkg/util/regutil"
)

type Service interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type cleanerService struct {
	artifactRepo domain.ArtifactRepository
	appRepo      domain.ApplicationRepository
	buildRepo    domain.BuildRepository
	image        builder.ImageConfig
	storage      domain.Storage

	start        func()
	startOnce    sync.Once
	shutdown     func()
	shutdownOnce sync.Once
}

func NewService(
	artifactRepo domain.ArtifactRepository,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	image builder.ImageConfig,
	storage domain.Storage,
) (Service, error) {
	c := &cleanerService{
		artifactRepo: artifactRepo,
		appRepo:      appRepo,
		buildRepo:    buildRepo,
		image:        image,
		storage:      storage,
	}

	r := image.NewRegistry()

	ctx, cancel := context.WithCancel(context.Background())
	c.start = func() {
		go loop.Loop(ctx, func(ctx context.Context) {
			start := time.Now()
			err := c.pruneImages(ctx, r, image.Registry.Addr)
			if err != nil {
				log.Errorf("failed to prune images: %+v", err)
				return
			}
			log.Infof("Pruned images in %v", time.Since(start))
		}, 1*time.Hour, true)
		go loop.Loop(ctx, func(ctx context.Context) {
			start := time.Now()
			err := c.pruneArtifacts(ctx)
			if err != nil {
				log.Errorf("failed to prune artifacts: %+v", err)
				return
			}
			log.Infof("Pruned artifacts in %v", time.Since(start))
		}, 1*time.Hour, true)
	}
	c.shutdown = cancel

	return c, nil
}

func (c *cleanerService) Start(_ context.Context) error {
	c.startOnce.Do(c.start)
	return nil
}

func (c *cleanerService) Shutdown(_ context.Context) error {
	c.shutdownOnce.Do(c.shutdown)
	return nil
}

func (c *cleanerService) getOlderBuilds(ctx context.Context, appID string, targetBuildID string) ([]*domain.Build, error) {
	if targetBuildID == "" {
		return nil, nil
	}
	builds, err := c.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{ApplicationID: optional.From(appID)})
	if err != nil {
		return nil, err
	}
	current, ok := lo.Find(builds, func(b *domain.Build) bool { return b.ID == targetBuildID })
	if !ok {
		return nil, errors.Errorf("failed to find build %v in retrieved builds", targetBuildID)
	}
	return lo.Filter(builds, func(b *domain.Build, _ int) bool { return b.QueuedAt.Before(current.QueuedAt) }), nil
}

func (c *cleanerService) pruneImages(ctx context.Context, r *regclient.RegClient, regHost string) error {
	applications, err := c.appRepo.GetApplications(ctx, domain.GetApplicationCondition{DeployType: optional.From(domain.DeployTypeRuntime)})
	if err != nil {
		return err
	}
	appsMap := lo.SliceToMap(applications, func(app *domain.Application) (string, *domain.Application) { return app.ID, app })

	repos, err := regutil.RepoList(ctx, r, regHost)
	if err != nil {
		return errors.Wrap(err, "failed to get image repositories")
	}

	for _, imageName := range repos {
		if !strings.HasPrefix(imageName, c.image.NamePrefix) {
			continue
		}
		err = c.pruneImage(ctx, r, regHost, imageName, appsMap)
		if err != nil {
			log.Errorf("pruning image %v: %+v", imageName, err)
			// fail-safe for each image
		}
	}

	return nil
}

func (c *cleanerService) pruneImage(ctx context.Context, r *regclient.RegClient, regHost string, imageName string, appsMap map[string]*domain.Application) error {
	appID := strings.TrimPrefix(imageName, c.image.NamePrefix)

	tags, err := regutil.TagList(ctx, r, regHost, imageName)
	if err != nil {
		return errors.Wrap(err, "getting tags")
	}
	app, ok := appsMap[appID]
	var danglingTags []string
	if ok {
		// app still exists; compare by queued_at time, then delete any older builds
		olderBuilds, err := c.getOlderBuilds(ctx, app.ID, app.CurrentBuild)
		if err != nil {
			return err
		}
		olderBuildIDs := ds.Map(olderBuilds, func(b *domain.Build) string { return b.ID })
		danglingTags = lo.Filter(tags, func(tag string, _ int) bool { return lo.Contains(olderBuildIDs, tag) })
	} else {
		// app was deleted
		danglingTags = tags
	}
	for _, tag := range danglingTags {
		// NOTE: needs manual execution of "registry garbage-collect <config> --delete-untagged" in docker registry side
		// to actually delete the layers
		// https://docs.docker.com/registry/garbage-collection/
		tagRef, err := ref.New(regHost + "/" + imageName + ":" + tag)
		if err != nil {
			return err
		}
		err = r.TagDelete(ctx, tagRef)
		if err != nil {
			log.Errorf("deleting tag %s:%s: %+v", imageName, tag, err)
			// fail-safe and continue
		}
	}
	return nil
}

func (c *cleanerService) getArtifactsNoLongerInUse(ctx context.Context) ([]*domain.Artifact, error) {
	applications, err := c.appRepo.GetApplications(ctx, domain.GetApplicationCondition{
		DeployType: optional.From(domain.DeployTypeStatic),
	})
	if err != nil {
		return nil, err
	}

	artifacts := make([]*domain.Artifact, 0, len(applications))
	for _, app := range applications {
		olderBuilds, err := c.getOlderBuilds(ctx, app.ID, app.CurrentBuild)
		if err != nil {
			return nil, err
		}
		for _, b := range olderBuilds {
			artifacts = append(artifacts, b.Artifacts...)
		}
	}
	return artifacts, nil
}

func (c *cleanerService) pruneArtifacts(ctx context.Context) error {
	notInUse, err := c.getArtifactsNoLongerInUse(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get artifacts in use")
	}

	for _, artifact := range notInUse {
		err = domain.DeleteArtifact(c.storage, artifact.ID)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("deleting artifact %v", artifact.ID))
		}
		err = c.artifactRepo.UpdateArtifact(ctx, artifact.ID, domain.UpdateArtifactArgs{DeletedAt: optional.From(time.Now())})
		if err != nil {
			return errors.Wrap(err, "failed to mark artifact as deleted")
		}
	}
	return nil
}
