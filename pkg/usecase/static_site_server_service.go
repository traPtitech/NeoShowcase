package usecase

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
)

type StaticSiteServerService interface {
	Reload(ctx context.Context) error
}

type staticSiteServerService struct {
	buildRepo repository.BuildRepository
	engine    domain.Engine

	// TODO 後で消す
	db *sql.DB
}

func NewStaticSiteServerService(
	buildRepo repository.BuildRepository,
	engine domain.Engine,
	db *sql.DB,
) StaticSiteServerService {
	return &staticSiteServerService{
		buildRepo: buildRepo,
		engine:    engine,
		db:        db,
	}
}

func (s *staticSiteServerService) Reload(ctx context.Context) error {
	applications, err := models.Applications(
		models.ApplicationWhere.BuildType.EQ(models.ApplicationsBuildTypeStatic),
		qm.Load(models.ApplicationRels.Website),
	).All(ctx, s.db)
	if err != nil {
		return err
	}

	var data []*domain.Site
	for _, app := range applications {
		website := app.R.Website
		if website == nil {
			continue
		}

		build, err := s.buildRepo.GetLastSuccessBuild(ctx, app.ID)
		if err != nil && err != repository.ErrNotFound {
			return err
		}
		if err == repository.ErrNotFound {
			continue
		}

		if !build.Artifact.Valid {
			continue
		}
		data = append(data, &domain.Site{
			ID:            website.ID,
			FQDN:          website.FQDN,
			ArtifactID:    build.Artifact.V.ID,
			ApplicationID: app.ID,
		})
	}

	return s.engine.Reconcile(data)
}
