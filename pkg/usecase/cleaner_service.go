package usecase

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/loop"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type CleanerService interface {
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

func NewCleanerService(
	artifactRepo domain.ArtifactRepository,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	image builder.ImageConfig,
	storage domain.Storage,
) (CleanerService, error) {
	c := &cleanerService{
		artifactRepo: artifactRepo,
		appRepo:      appRepo,
		buildRepo:    buildRepo,
		image:        image,
		storage:      storage,
	}

	r, err := registry.New(image.Registry.Scheme+"://"+image.Registry.Addr, image.Registry.Username, image.Registry.Password)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.start = func() {
		go loop.Loop(ctx, func(ctx context.Context) {
			start := time.Now()
			err := c.pruneImages(ctx, r)
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

func (c *cleanerService) pruneImages(ctx context.Context, r *registry.Registry) error {
	applications, err := c.appRepo.GetApplications(ctx, domain.GetApplicationCondition{DeployType: optional.From(domain.DeployTypeRuntime)})
	if err != nil {
		return err
	}
	appsMap := lo.SliceToMap(applications, func(app *domain.Application) (string, *domain.Application) { return app.ID, app })

	imageNames, err := r.Repositories()
	if err != nil {
		return errors.Wrap(err, "failed to get image repositories")
	}

	for _, imageName := range imageNames {
		if !strings.HasPrefix(imageName, c.image.NamePrefix) {
			continue
		}
		appID := strings.TrimPrefix(imageName, c.image.NamePrefix)

		tags, err := r.Tags(imageName)
		if err != nil {
			return errors.Wrap(err, "failed to get tags for image")
		}
		app, ok := appsMap[appID]
		var danglingTags []string
		if ok {
			danglingTags = lo.Reject(tags, func(tag string, _ int) bool { return app.WantCommit == tag || app.CurrentCommit == tag })
		} else {
			danglingTags = tags
		}
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

func (c *cleanerService) pruneArtifacts(ctx context.Context) error {
	inUse, err := domain.GetArtifactsInUse(ctx, c.appRepo, c.buildRepo)
	if err != nil {
		return errors.Wrap(err, "failed to get artifacts in use")
	}
	inUseIDs := lo.SliceToMap(inUse, func(a *domain.Artifact) (string, bool) { return a.ID, true })

	artifacts, err := c.artifactRepo.GetArtifacts(ctx, domain.GetArtifactCondition{IsDeleted: optional.From(false)})
	if err != nil {
		return errors.Wrap(err, "failed to get existing artifacts")
	}

	for _, artifact := range artifacts {
		if inUseIDs[artifact.ID] {
			continue
		}
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
