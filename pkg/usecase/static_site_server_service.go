package usecase

import (
	"context"
	"database/sql"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	branches, err := models.Branches(
		models.BranchWhere.BuildType.EQ(models.BranchesBuildTypeStatic),
		qm.Load(models.BranchRels.Website),
		qm.Load(qm.Rels(models.BranchRels.Build, models.BuildLogRels.Artifact)),
	).All(ctx, s.db)
	if err != nil {
		return err
	}

	var data []*domain.Site
	for _, branch := range branches {
		if branch.R.Website != nil {
			website := branch.R.Website
			if branch.R.Build != nil {
				build := branch.R.Build
				if build.R.Artifact != nil {
					artifact := build.R.Artifact
					data = append(data, &domain.Site{
						ID:            website.ID,
						FQDN:          website.FQDN,
						ArtifactID:    artifact.ID,
						ApplicationID: branch.ApplicationID,
					})
				}
			}
		}
	}

	return s.engine.Reconcile(data)
}
