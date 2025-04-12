package registry

import (
	"context"

	"github.com/pkg/errors"
	"github.com/regclient/regclient"
	"github.com/regclient/regclient/config"
	"github.com/regclient/regclient/scheme"
	"github.com/regclient/regclient/types/descriptor"
	"github.com/regclient/regclient/types/manifest"
	"github.com/regclient/regclient/types/ref"
	"github.com/samber/lo"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
)

var _ builder.RegistryClient = &client{}

type client struct {
	regclient *regclient.RegClient
}

func NewClient(conf builder.ImageConfig) builder.RegistryClient {
	var opts []regclient.Opt

	host := config.HostNewName(conf.Registry.Scheme + "://" + conf.Registry.Addr)
	// RepoAuth should be set to true, because by default regclient internally merges scopes for all repositories
	// it accesses, resulting in a bloating "Authorization" header when accessing large number of repositories at once.
	// also see: https://distribution.github.io/distribution/spec/auth/jwt/
	host.RepoAuth = true
	if conf.Registry.Username != "" {
		host.User = conf.Registry.Username
	}
	if conf.Registry.Password != "" {
		host.Pass = conf.Registry.Password
	}
	opts = append(opts, regclient.WithConfigHost(*host))

	c := regclient.New(opts...)

	return &client{regclient: c}
}

func (c *client) DeleteImage(ctx context.Context, image, tag string) error {
	ref, err := ref.New(image + ":" + tag)
	if err != nil {
		return errors.Wrap(err, "invalid image reference")
	}
	err = c.regclient.TagDelete(ctx, ref)
	if err != nil {
		return errors.Wrap(err, "delete image")
	}
	return nil
}

func (c *client) GetTags(ctx context.Context, image string) ([]string, error) {
	ref, err := ref.New(image)
	if err != nil {
		return nil, errors.Wrap(err, "invalid image reference")
	}

	const limit = 100
	tags := []string{}

	for {
		opts := []scheme.TagOpts{scheme.WithTagLimit(limit)}
		if len(tags) > 0 {
			opts = append(opts, scheme.WithTagLast(tags[len(tags)-1]))
		}
		tagList, err := c.regclient.TagList(ctx, ref, opts...)
		if err != nil {
			return nil, errors.Wrap(err, "list tags")
		}
		tags = append(tags, tagList.Tags...)
		if len(tagList.Tags) < limit {
			break
		}
	}

	return tags, nil
}

func (c *client) GetImageSize(ctx context.Context, image, tag string) (int64, error) {
	ref, err := ref.New(image + ":" + tag)
	if err != nil {
		return 0, errors.Wrap(err, "invalid image reference")
	}
	m, err := c.regclient.ManifestGet(ctx, ref)
	if err != nil {
		return 0, errors.Wrap(err, "get manifest")
	}
	imager, ok := m.(manifest.Imager)
	if !ok {
		return 0, errors.New("interface conversion failed: manifest is not Imager")
	}
	layers, err := imager.GetLayers()
	if err != nil {
		return 0, errors.Wrap(err, "get image layers")
	}

	size := lo.SumBy(layers, func(l descriptor.Descriptor) int64 { return l.Size })

	return size, nil

}
