package appmanager

import (
	"context"
	"database/sql"
	"github.com/leandro-lugaresi/hub"
	builderApi "github.com/traPtitech/neoshowcase/pkg/builder/api"
	"github.com/traPtitech/neoshowcase/pkg/container"
	ssgenApi "github.com/traPtitech/neoshowcase/pkg/staticsitegen/api"
)

type managerImpl struct {
	db      *sql.DB
	bus     *hub.Hub
	builder builderApi.BuilderServiceClient
	ssgen   ssgenApi.StaticSiteGenServiceClient
	cm      container.Manager
}

func NewManager(db *sql.DB, bus *hub.Hub, builder builderApi.BuilderServiceClient, ssgen ssgenApi.StaticSiteGenServiceClient, cm container.Manager) (Manager, error) {
	return &managerImpl{
		db:      db,
		bus:     bus,
		builder: builder,
		ssgen:   ssgen,
		cm:      cm,
	}, nil
}

func (m *managerImpl) GetApp(appID string) (App, error) {
	panic("not implemented") // TODO
}

func (m *managerImpl) Shutdown(ctx context.Context) error {
	return nil
}
