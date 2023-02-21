package usecase

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

type StaticSiteServerService interface {
	Reload(ctx context.Context) error
}

type staticSiteServerService struct {
	engine domain.Engine

	// TODO 後で消す
	db *sql.DB
}

func NewStaticSiteServerService(engine domain.Engine, db *sql.DB) StaticSiteServerService {
	return &staticSiteServerService{
		engine: engine,
		db:     db,
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

		build, err := models.Builds(
			models.BuildWhere.ApplicationID.EQ(app.ID),
			models.BuildWhere.Status.EQ(builder.BuildStatusSucceeded.String()),
			qm.OrderBy(models.BuildColumns.FinishedAt+" desc"),
			qm.Load(models.BuildRels.Artifact),
		).One(ctx, s.db)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			continue
		}

		artifact := build.R.Artifact
		if artifact == nil {
			continue
		}
		data = append(data, &domain.Site{
			ID:            website.ID,
			FQDN:          website.FQDN,
			ArtifactID:    artifact.ID,
			ApplicationID: app.ID,
		})
	}

	return s.engine.Reconcile(data)
}
