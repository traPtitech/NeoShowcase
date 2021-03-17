package appmanager

import (
	"container/list"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"

	"github.com/leandro-lugaresi/hub"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/builder/api"
	builderApi "github.com/traPtitech/neoshowcase/pkg/builder/api"
	"github.com/traPtitech/neoshowcase/pkg/container"
	"github.com/traPtitech/neoshowcase/pkg/event"
	"github.com/traPtitech/neoshowcase/pkg/models"
	ssgenApi "github.com/traPtitech/neoshowcase/pkg/staticsitegen/api"
	"github.com/traPtitech/neoshowcase/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type managerImpl struct {
	db      *sql.DB
	bus     *hub.Hub
	builder builderApi.BuilderServiceClient
	ssgen   ssgenApi.StaticSiteGenServiceClient
	cm      container.Manager

	config Config

	queue buildQueue

	stream builderApi.BuilderService_ConnectEventStreamClient
}

type buildTask struct {
	ctx       context.Context
	in        interface{}
	opts      []grpc.CallOption
	buildType string
}

type buildQueue struct {
	queue list.List
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
		case builderApi.Event_BUILD_STARTED:
			m.bus.Publish(hub.Message{
				Name:   event.BuilderBuildStarted,
				Fields: payload,
			})
		case builderApi.Event_BUILD_SUCCEEDED:
			m.bus.Publish(hub.Message{
				Name:   event.BuilderBuildSucceeded,
				Fields: payload,
			})
		case builderApi.Event_BUILD_FAILED:
			m.bus.Publish(hub.Message{
				Name:   event.BuilderBuildFailed,
				Fields: payload,
			})
		case builderApi.Event_BUILD_CANCELED:
			m.bus.Publish(hub.Message{
				Name:   event.BuilderBuildCanceled,
				Fields: payload,
			})
		}
	}
	return nil
}

func (m *managerImpl) appDeployLoop() {
	sub := m.bus.Subscribe(10, event.BuilderBuildSucceeded, event.WebhookRepositoryPush)
	for ev := range sub.Receiver {
		switch ev.Name {
		case event.WebhookRepositoryPush:
			repoURL := ev.Fields["repository_url"].(string)
			branch := ev.Fields["branch"].(string)

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
			stat, err := m.builder.GetStatus(context.Background(), &emptypb.Empty{})
			if err != nil {
				log.WithError(err).
					Error("failed to builder.Getstatus")
			}
			if m.queue.queue.Len() > 0 {
				if stat.GetStatus() == api.BuilderStatus_WAITING { //TODO:複数builderに対応
					_, err = m.sendBuildRequest()
					if err != nil {
						log.WithError(err).
							Error("failed to sendRequestBuild")
					}
				}
			}
		case event.BuilderBuildSucceeded:
			envID := ev.Fields["environment_id"].(string)
			buildID := ev.Fields["build_id"].(string)

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
			stat, err := m.builder.GetStatus(context.Background(), &emptypb.Empty{})
			if err != nil {
				log.WithError(err).
					Error("failed to builder.Getstatus")
			}
			if m.queue.queue.Len() > 0 { //TODO:複数builderに対応
				if stat.GetStatus() == api.BuilderStatus_WAITING {
					_, err = m.sendBuildRequest()
					if err != nil {
						log.WithError(err).
							Error("failed to sendRequestBuild")
					}
				}
			}
		}
	}
}

func (m *managerImpl) sendBuildRequest() (interface{}, error) {
	req, err := m.queue.PopQueue()
	if err != nil {
		return nil, err
	}
	v := req.Value.(buildTask)
	switch v.buildType {
	case models.EnvironmentsBuildTypeImage:
		res, err := m.builder.StartBuildImage(v.ctx, v.in.(*builderApi.StartBuildImageRequest), v.opts...)
		if err != nil {
			return nil, err
		}
		return res, nil
	case models.EnvironmentsBuildTypeStatic:
		res, err := m.builder.StartBuildStatic(v.ctx, v.in.(*builderApi.StartBuildStaticRequest), v.opts...)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		return nil, fmt.Errorf("unknown build type: %s", v.buildType)
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

// PushQueue ビルドキューにアイテムを追加する
func (b *buildQueue) PushQueue(ctx context.Context, in interface{}, buildType string, opts ...grpc.CallOption) (*list.Element, error) {
	//TODO:staticbuildも対応
	t := &buildTask{
		ctx:       ctx,
		in:        in,
		opts:      opts,
		buildType: buildType,
	}
	// 重複チェック
	for e := b.queue.Front(); e != nil; e = e.Next() {
		if e.Value.(*buildTask) == t {
			log.Error("Already Existed")
			return e, errors.New("Already Existed")
		}
	}
	r := b.queue.PushBack(t)
	return r, nil
}

// PushQueue ビルドキューから先頭のアイテムを取り出す
func (b *buildQueue) PopQueue() (*list.Element, error) {
	if b.queue.Len() == 0 {
		return nil, errors.New("No Elements")
	}
	r := b.queue.Front()
	b.queue.Remove(r)
	return r, nil
}
