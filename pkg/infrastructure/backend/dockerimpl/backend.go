package dockerimpl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/traPtitech/neoshowcase/pkg/util/retry"

	clitypes "github.com/docker/cli/cli/config/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

const (
	appLabel            = "ns.trap.jp/app"
	appIDLabel          = "ns.trap.jp/app-id"
	appRestartedAtLabel = "ns.trap.jp/restarted-at"
)

const (
	traefikRuntimeFilename = "apps.yaml"
	traefikSSFilename      = "ss.yaml"
	traefikSSServiceName   = "ss"
)

var _ domain.Backend = (*Backend)(nil)

type Backend struct {
	c      *client.Client
	config Config
	image  builder.ImageConfig

	eventSubs   domain.PubSub[*domain.ContainerEvent]
	stopWatcher func()

	reloadLock sync.Mutex
}

func NewClientFromEnv() (*client.Client, error) {
	return client.NewClientWithOpts(
		client.FromEnv,
		// Using github.com/moby/moby of v25 master@032797ea4bcb (2023-09-05), required by github.com/moby/buildkit@v0.12.3,
		// defaults to API version 1.44, which currently available docker installation does not support.
		client.WithVersion("1.43"),
	)
}

func NewDockerBackend(
	c *client.Client,
	config Config,
	image builder.ImageConfig,
) (*Backend, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}
	b := &Backend{
		c:      c,
		config: config,
		image:  image,
	}
	return b, nil
}

func (b *Backend) Start(ctx context.Context) error {
	// showcase用のネットワークを用意
	if err := b.initNetworks(ctx); err != nil {
		return errors.Wrap(err, "failed to init networks")
	}

	eventCtx, eventCancel := context.WithCancel(context.Background())
	b.stopWatcher = eventCancel
	go retry.Do(eventCtx, b.eventListener, "container watcher")

	return nil
}

func (b *Backend) eventListener(ctx context.Context) error {
	// https://docs.docker.com/engine/reference/commandline/events/
	ch, errCh := b.c.Events(ctx, types.EventsOptions{Filters: filters.NewArgs(filters.Arg("type", "container"))})
	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				return nil
			}
			switch ev.Type {
			case "container":
				appID, ok := ev.Actor.Attributes[appIDLabel]
				if !ok {
					continue
				}
				b.eventSubs.Publish(&domain.ContainerEvent{ApplicationID: appID})
			}
		case err, ok := <-errCh:
			if !ok {
				return nil
			}
			return err
		}
	}
}

func (b *Backend) Dispose(_ context.Context) error {
	b.stopWatcher()
	return nil
}

func (b *Backend) AvailableDomains() domain.AvailableDomainSlice {
	return ds.Map(b.config.Domains, (*domainConf).toDomainAD)
}

func (b *Backend) targetAuth(fqdn string) *domainAuthConf {
	for _, dc := range b.config.Domains {
		if dc.Auth.Available && dc.toDomainAD().Match(fqdn) {
			return dc.Auth
		}
	}
	return nil
}

func (b *Backend) AvailablePorts() domain.AvailablePortSlice {
	return ds.Map(b.config.Ports, (*portConf).toDomainAP)
}

func (b *Backend) ListenContainerEvents() (sub <-chan *domain.ContainerEvent, unsub func()) {
	return b.eventSubs.Subscribe()
}

func (b *Backend) initNetworks(ctx context.Context) error {
	networks, err := b.c.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to list networks")
	}
	for _, network := range networks {
		if network.Name == b.config.Network {
			return nil
		}
	}

	_, err = b.c.NetworkCreate(ctx, b.config.Network, types.NetworkCreate{})
	return err
}

func (b *Backend) authConfig() (string, error) {
	if b.image.Registry.Username == "" && b.image.Registry.Password == "" {
		return "", nil
	}
	c := clitypes.AuthConfig{
		Username: b.image.Registry.Username,
		Password: b.image.Registry.Password,
	}
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func getRestartedAt(c *types.Container) time.Time {
	if c == nil {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339Nano, c.Labels[appRestartedAtLabel])
	if err != nil {
		return time.Time{}
	}
	return t
}

func (b *Backend) containerLabels(app *domain.Application) map[string]string {
	return ds.MergeMap(b.config.labels(), map[string]string{
		appLabel:            "true",
		appIDLabel:          app.ID,
		appRestartedAtLabel: app.UpdatedAt.Format(time.RFC3339Nano),
	})
}

func containerName(appID string) string {
	return fmt.Sprintf("nsapp-%s", appID)
}

func networkName(appID string) string {
	return fmt.Sprintf("%s.nsapp.internal", appID)
}

func traefikName(website *domain.Website) string {
	return fmt.Sprintf("nsapp-%s", website.ID)
}

func stripMiddlewareName(website *domain.Website) string {
	return traefikName(website) + "-strip"
}

func ssHeaderMiddlewareName(ss *domain.StaticSite) string {
	return fmt.Sprintf("nsapp-ss-header-%s", ss.Application.ID)
}
