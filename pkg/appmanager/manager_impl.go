package appmanager

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/leandro-lugaresi/hub"
	log "github.com/sirupsen/logrus"
	builderApi "github.com/traPtitech/neoshowcase/pkg/builder/api"
	"github.com/traPtitech/neoshowcase/pkg/container"
	"github.com/traPtitech/neoshowcase/pkg/event"
	"github.com/traPtitech/neoshowcase/pkg/models"
	ssgenApi "github.com/traPtitech/neoshowcase/pkg/staticsitegen/api"
	"github.com/traPtitech/neoshowcase/pkg/util"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
)

type managerImpl struct {
	db      *sql.DB
	bus     *hub.Hub
	builder builderApi.BuilderServiceClient
	ssgen   ssgenApi.StaticSiteGenServiceClient
	cm      container.Manager

	config Config

	stream builderApi.BuilderService_ConnectEventStreamClient
}

type Config struct {
	DB              *sql.DB
	Hub             *hub.Hub
	Builder         builderApi.BuilderServiceClient
	SSGen           ssgenApi.StaticSiteGenServiceClient
	CM              container.Manager
	ImageRegistry   string
	ImageNamePrefix string
}

func NewManager(config Config) (Manager, error) {
	stream, err := config.Builder.ConnectEventStream(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	m := &managerImpl{
		db:      config.DB,
		bus:     config.Hub,
		builder: config.Builder,
		ssgen:   config.SSGen,
		cm:      config.CM,
		config:  config,
		stream:  stream,
	}
	go m.receiveBuilderEvents()
	return m, nil
}

// receiveBuilderEvents　builderから届くイベントを内部イベントに変換してpublish
func (m *managerImpl) receiveBuilderEvents() error {
	for {
		ev, err := m.stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.WithField("payload", ev.String()).Debug("builder event was received")
		switch ev.Type {
		case builderApi.Event_BUILD_STARTED:
			m.bus.Publish(hub.Message{
				Name:   event.BuilderBuildStarted,
				Fields: util.FromJSON(ev.Body),
			})
		case builderApi.Event_BUILD_SUCCEEDED:
			m.bus.Publish(hub.Message{
				Name:   event.BuilderBuildSucceeded,
				Fields: util.FromJSON(ev.Body),
			})
		case builderApi.Event_BUILD_FAILED:
			m.bus.Publish(hub.Message{
				Name:   event.BuilderBuildFailed,
				Fields: util.FromJSON(ev.Body),
			})
		case builderApi.Event_BUILD_CANCELED:
			m.bus.Publish(hub.Message{
				Name:   event.BuilderBuildCanceled,
				Fields: util.FromJSON(ev.Body),
			})
		}
	}
	return nil
}

// getFullImageName registryのhost付きのイメージ名を返す
func (m *managerImpl) getFullImageName(app App) string {
	if m.config.ImageRegistry == "" {
		return m.getImageName(app)
	}
	return m.config.ImageRegistry + "/" + m.getImageName(app)
}

// getImageName イメージ名を返す
func (m *managerImpl) getImageName(app App) string {
	return m.config.ImageNamePrefix + app.GetID()
}

func (m *managerImpl) GetApp(appID string) (App, error) {
	app, err := models.Applications(
		qm.Load(models.ApplicationRels.Repository),
		models.ApplicationWhere.DeletedAt.IsNull(),
		models.ApplicationWhere.ID.EQ(appID),
	).One(context.Background(), m.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to GetApp: %w", err)
	}

	return &appImpl{
		m:       m,
		dbmodel: app,
	}, nil
}

func (m *managerImpl) GetAppByRepository(repo string) (App, error) {
	repoModel, err := models.Repositories(models.RepositoryWhere.Remote.EQ(repo)).One(context.Background(), m.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to GetAppByRepository: %w", err)
	}

	app, err := models.Applications(
		qm.Load(models.ApplicationRels.Repository),
		models.ApplicationWhere.DeletedAt.IsNull(),
		models.ApplicationWhere.RepositoryID.EQ(repoModel.ID),
	).One(context.Background(), m.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to GetApp: %w", err)
	}

	return &appImpl{
		m:       m,
		dbmodel: app,
	}, nil
}

func (m *managerImpl) Shutdown(ctx context.Context) error {
	return nil
}
