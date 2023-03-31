package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type CleanerService interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type cleanerService struct {
	appRepo domain.ApplicationRepository
	image   builder.ImageConfig

	start        func()
	startOnce    sync.Once
	shutdown     func()
	shutdownOnce sync.Once
}

func NewCleanerService(
	appRepo domain.ApplicationRepository,
	image builder.ImageConfig,
) (CleanerService, error) {
	c := &cleanerService{
		appRepo: appRepo,
		image:   image,
	}

	r, err := registry.New(image.Registry.Scheme+"://"+image.Registry.Addr, image.Registry.Username, image.Registry.Password)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.start = func() {
		go c.pruneImagesLoop(ctx, r)
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

func (c *cleanerService) pruneImagesLoop(ctx context.Context, r *registry.Registry) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	doPrune := func() {
		start := time.Now()
		err := c.pruneImages(ctx, r)
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

func (c *cleanerService) pruneImages(ctx context.Context, r *registry.Registry) error {
	applications, err := c.appRepo.GetApplications(ctx, domain.GetApplicationCondition{BuildType: optional.From(builder.BuildTypeRuntime)})
	if err != nil {
		return err
	}
	commits := make(map[string]struct{}, 2*len(applications))
	for _, app := range applications {
		commits[app.WantCommit] = struct{}{}
		commits[app.CurrentCommit] = struct{}{}
	}

	for _, app := range applications {
		imageName := c.image.ImageName(app.ID)
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
