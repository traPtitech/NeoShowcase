package usecase

import (
	"context"
	"database/sql"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type StaticSiteService interface {
	Reload(ctx context.Context) error
}

type staticSiteService struct {
	engine staticserver.Engine

	// TODO 後で消す
	db *sql.DB
}

func NewStaticSiteService(engine staticserver.Engine, db *sql.DB) StaticSiteService {
	return &staticSiteService{
		engine: engine,
		db:     db,
	}
}

func (s *staticSiteService) Reload(ctx context.Context) error {
	envs, err := models.Environments(
		models.EnvironmentWhere.BuildType.EQ(models.EnvironmentsBuildTypeStatic),
		qm.Load(models.EnvironmentRels.Website),
		qm.Load(qm.Rels(models.EnvironmentRels.Build, models.BuildLogRels.Artifact)),
	).All(ctx, s.db)
	if err != nil {
		return err
	}

	var data []*domain.Site
	for _, env := range envs {
		if env.R.Website != nil {
			website := env.R.Website
			if env.R.Build != nil {
				build := env.R.Build
				if build.R.Artifact != nil {
					artifact := build.R.Artifact
					data = append(data, &domain.Site{
						ID:            website.ID,
						FQDN:          website.FQDN,
						ArtifactID:    artifact.ID,
						ApplicationID: env.ApplicationID,
					})
				}
			}
		}
	}

	return s.engine.Reconcile(data)
}
