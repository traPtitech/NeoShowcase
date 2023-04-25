package dockerimpl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

const (
	appLabel            = "neoshowcase.trap.jp/app"
	appIDLabel          = "neoshowcase.trap.jp/appId"
	appRestartedAtLabel = "neoshowcase.trap.jp/restartedAt"
)

const (
	traefikRuntimeFilename = "apps.yaml"
	traefikSSFilename      = "ss.yaml"
	traefikSSServiceName   = "ss"
)

type dockerBackend struct {
	c         *client.Client
	conf      Config
	image     builder.ImageConfig
	eventSubs domain.PubSub[*domain.ContainerEvent]

	eventCh     <-chan events.Message
	eventErr    <-chan error
	eventCancel func()

	reloadLock sync.Mutex
}

func NewClientFromEnv() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv)
}

func NewDockerBackend(
	c *client.Client,
	conf Config,
	image builder.ImageConfig,
) (domain.Backend, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}
	return &dockerBackend{
		c:     c,
		conf:  conf,
		image: image,
	}, nil
}

func (b *dockerBackend) Start(ctx context.Context) error {
	// showcase用のネットワークを用意
	if err := b.initNetworks(ctx); err != nil {
		return errors.Wrap(err, "failed to init networks")
	}

	eventCtx, eventCancel := context.WithCancel(context.Background())
	b.eventCancel = eventCancel
	go b.eventListenerLoop(eventCtx)

	return nil
}

func (b *dockerBackend) eventListenerLoop(ctx context.Context) {
	for {
		err := b.eventListener(ctx)
		if err == nil {
			return
		}
		log.Errorf("docker event listner errored, retrying in 1s: %+v", err)
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return
		}
	}
}

func (b *dockerBackend) eventListener(ctx context.Context) error {
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

func (b *dockerBackend) Dispose(_ context.Context) error {
	b.eventCancel()
	return nil
}

func (b *dockerBackend) AuthAllowed(fqdn string) bool {
	for _, ac := range b.conf.Middlewares.Auth {
		if domain.MatchDomain(ac.Domain, fqdn) {
			return true
		}
	}
	return false
}

func (b *dockerBackend) targetAuth(fqdn string) *authConf {
	for _, ac := range b.conf.Middlewares.Auth {
		if domain.MatchDomain(ac.Domain, fqdn) {
			return ac
		}
	}
	return nil
}

func (b *dockerBackend) ListenContainerEvents() (sub <-chan *domain.ContainerEvent, unsub func()) {
	return b.eventSubs.Subscribe()
}

func (b *dockerBackend) initNetworks(ctx context.Context) error {
	networks, err := b.c.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to list networks")
	}
	for _, network := range networks {
		if network.Name == b.conf.Network {
			return nil
		}
	}

	_, err = b.c.NetworkCreate(ctx, b.conf.Network, types.NetworkCreate{})
	return err
}

func (b *dockerBackend) authConfig() (string, error) {
	if b.image.Registry.Username == "" && b.image.Registry.Password == "" {
		return "", nil
	}
	c := types.AuthConfig{
		Username: b.image.Registry.Username,
		Password: b.image.Registry.Password,
	}
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (b *dockerBackend) containerLabels(app *domain.Application) map[string]string {
	return ds.MergeMap(b.conf.labels(), map[string]string{
		appLabel:            "true",
		appIDLabel:          app.ID,
		appRestartedAtLabel: app.UpdatedAt.Format(time.RFC3339),
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
