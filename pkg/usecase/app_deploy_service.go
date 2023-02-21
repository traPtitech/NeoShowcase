package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
)

type AppDeployService interface {
	QueueDeployment(ctx context.Context, applicationID string, buildID string) error
}

type appDeployService struct {
	backend domain.Backend
	ss      pb.StaticSiteServiceClient

	imageRegistry   string
	imageNamePrefix string

	// TODO 後で消す
	db *sql.DB
}

func NewAppDeployService(backend domain.Backend, ss pb.StaticSiteServiceClient, registry builder.DockerImageRegistryString, prefix builder.DockerImageNamePrefixString, db *sql.DB) AppDeployService {
	return &appDeployService{
		backend:         backend,
		ss:              ss,
		imageRegistry:   string(registry),
		imageNamePrefix: string(prefix),
		db:              db,
	}
}

func (s *appDeployService) QueueDeployment(ctx context.Context, applicationID string, buildID string) error {
	app, err := models.Applications(
		qm.Load(models.ApplicationRels.Website),
		models.ApplicationWhere.ID.EQ(applicationID),
	).One(ctx, s.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("application (%s) not found", applicationID)
		}
		return fmt.Errorf("failed to get Application: %w", err)
	}

	ok, err := app.Builds(
		models.BuildWhere.ID.EQ(buildID),
		models.BuildWhere.Status.EQ(builder.BuildStatusSucceeded.String()),
	).Exists(ctx, s.db)
	if err != nil {
		return fmt.Errorf("failed to BuildExists: %w", err)
	}
	if !ok {
		return fmt.Errorf("build (%s) not found", buildID)
	}

	entry := &appDeployment{
		ApplicationID: applicationID,
		BuildID:       buildID,
		QueuedAt:      time.Now(),
		app:           app,
	}
	// TODO ちゃんとキューに入れて非同期処理する
	go s.deploy(entry)
	return nil
}

type appDeployment struct {
	ApplicationID string
	BuildID       string
	QueuedAt      time.Time
	app           *models.Application
}

func (s *appDeployService) deploy(entry *appDeployment) {
	app := entry.app
	website := app.R.Website

	// TODO Ingressの設定がここでする必要があるかどうか確認する
	switch builder.BuildTypeFromString(app.BuildType) {
	case builder.BuildTypeImage:
		var httpProxy *domain.ContainerHTTPProxy
		if website != nil {
			httpProxy = &domain.ContainerHTTPProxy{
				Domain: website.FQDN,
				Port:   website.HTTPPort,
			}
		}

		err := s.backend.CreateContainer(context.Background(), domain.ContainerCreateArgs{
			ApplicationID: app.ID,
			ImageName:     builder.GetImageName(s.imageRegistry, s.imageNamePrefix, app.ID),
			ImageTag:      entry.BuildID,
			HTTPProxy:     httpProxy,
			Recreate:      true,
		})
		if err != nil {
			log.WithField("applicationID", entry.ApplicationID).
				WithField("buildID", entry.BuildID).
				WithField("queuedAt", entry.QueuedAt).
				WithError(err).
				Errorf("failed to create container")
			return
		}

	case builder.BuildTypeStatic:
		if _, err := s.ss.Reload(context.Background(), &pb.ReloadRequest{}); err != nil {
			log.WithField("applicationID", entry.ApplicationID).
				WithField("buildID", entry.BuildID).
				WithField("queuedAt", entry.QueuedAt).
				WithError(err).
				Errorf("failed to reload StaticSiteServer")
			return
		}
	}
}
