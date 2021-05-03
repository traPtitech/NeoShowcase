package appmanager

import (
	"context"
	"database/sql"
	"io"

	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

type managerImpl struct {
	db      *sql.DB
	bus     eventbus.Bus
	builder pb.BuilderServiceClient
	ss      pb.StaticSiteServiceClient
	backend backend.Backend

	config Config

	stream pb.BuilderService_ConnectEventStreamClient
}

type Config struct {
	DB              *sql.DB
	Hub             eventbus.Bus
	Builder         pb.BuilderServiceClient
	SS              pb.StaticSiteServiceClient
	Backend         backend.Backend
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
		ss:      config.SS,
		backend: config.Backend,
		config:  config,
		stream:  stream,
	}
	go m.appDeployLoop()
	go m.receiveBuilderEvents()
	return m, nil
}

// receiveBuilderEvents　builderから届くイベントを内部イベントに変換してpublish
func (m *managerImpl) receiveBuilderEvents() error {
	for {
		ev, err := m.stream.Recv()
		if err == io.EOF {
			log.Debug("builder event stream was closed: EOF")
			break
		}
		if err != nil {
			log.WithError(err).
				Debug("builder event stream was disconnected with error")
			return err
		}

		payload := util.FromJSON(ev.Body)

		log.WithField("type", ev.Type).
			WithField("payload", payload).
			Info("builder event received")

		switch ev.Type {
		case pb.Event_BUILD_STARTED:
			m.bus.Publish(event.BuilderBuildStarted, payload)
		case pb.Event_BUILD_SUCCEEDED:
			m.bus.Publish(event.BuilderBuildSucceeded, payload)
		case pb.Event_BUILD_FAILED:
			m.bus.Publish(event.BuilderBuildFailed, payload)
		case pb.Event_BUILD_CANCELED:
			m.bus.Publish(event.BuilderBuildCanceled, payload)
		}
	}
	return nil
}

func (m *managerImpl) appDeployLoop() {
	sub := m.bus.Subscribe(event.BuilderBuildSucceeded, event.WebhookRepositoryPush)
	defer sub.Unsubscribe()
	for ev := range sub.Chan() {
		switch ev.Type {
		case event.WebhookRepositoryPush:
			repoURL := ev.Body["repository_url"].(string)
			branch := ev.Body["branch"].(string)

			log.WithField("repo", repoURL).
				WithField("refs", branch).
				Info("repository push event received")

			app, err := m.GetAppByRepository(repoURL)
			if err != nil {
				if err != ErrNotFound {
					log.WithError(err).WithField("repoURL", repoURL).Error("failed to GetAppByRepository")
				}
				continue
			}

			env, err := app.GetEnvByBranchName(branch)
			if err != nil {
				if err != ErrNotFound {
					log.WithError(err).WithField("repoURL", repoURL).Error("failed to GetAppByRepository")
				}
				continue
			}

			if err := app.RequestBuild(context.Background(), env.GetID()); err != nil {
				log.WithError(err).
					WithField("appID", app.GetID()).
					WithField("envID", env.GetID()).
					Error("failed to RequestBuild")
			}

		case event.BuilderBuildSucceeded:
			envID := ev.Body["environment_id"].(string)
			buildID := ev.Body["build_id"].(string)
			if len(envID) == 0 {
				// envIDが無い場合はテストビルド
				continue
			}

			app, err := m.GetAppByEnvironment(envID)
			if err != nil {
				log.WithError(err).WithField("envID", envID).Error("failed to GetAppByEnvironment")
				continue
			}

			// 自動デプロイ
			log.WithField("envID", envID).
				WithField("buildID", buildID).
				Error("starting application")
			err = app.Start(AppStartArgs{
				EnvironmentID: envID,
				BuildID:       buildID,
			})
			if err != nil {
				log.WithError(err).
					WithField("envID", envID).
					WithField("buildID", buildID).
					Error("failed to Start Application")
			}
		}
	}
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

func (m *managerImpl) Shutdown(ctx context.Context) error {
	return nil
}
