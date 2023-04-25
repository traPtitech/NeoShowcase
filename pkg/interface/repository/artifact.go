package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository/repoconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type artifactRepository struct {
	db *sql.DB
}

func NewArtifactRepository(db *sql.DB) domain.ArtifactRepository {
	return &artifactRepository{
		db: db,
	}
}

func (r *artifactRepository) buildMods(cond domain.GetArtifactCondition) []qm.QueryMod {
	var mods []qm.QueryMod
	if cond.ApplicationID.Valid {
		mods = append(mods,
			qm.InnerJoin(fmt.Sprintf(
				"%s ON %s = %s",
				models.TableNames.Builds,
				models.BuildTableColumns.ID,
				models.ArtifactTableColumns.BuildID,
			)),
			models.BuildWhere.ApplicationID.EQ(cond.ApplicationID.V),
		)
	}
	if cond.IsDeleted.Valid {
		if cond.IsDeleted.V {
			mods = append(mods, models.ArtifactWhere.DeletedAt.IsNotNull())
		} else {
			mods = append(mods, models.ArtifactWhere.DeletedAt.IsNull())
		}
	}
	return mods
}

func (r *artifactRepository) GetArtifacts(ctx context.Context, cond domain.GetArtifactCondition) ([]*domain.Artifact, error) {
	artifacts, err := models.Artifacts(r.buildMods(cond)...).All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get artifacts")
	}
	return lo.Map(artifacts, func(a *models.Artifact, i int) *domain.Artifact {
		return repoconvert.ToDomainArtifact(a)
	}), nil
}

func (r *artifactRepository) CreateArtifact(ctx context.Context, artifact *domain.Artifact) error {
	ma := repoconvert.FromDomainArtifact(artifact)
	if err := ma.Insert(ctx, r.db, boil.Blacklist()); err != nil {
		return errors.Wrap(err, "failed to insert artifact")
	}
	return nil
}

func (r *artifactRepository) UpdateArtifact(ctx context.Context, id string, args domain.UpdateArtifactArgs) error {
	artifact, err := models.Artifacts(models.ArtifactWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to get artifact")
	}

	var cols []string
	if args.DeletedAt.Valid {
		artifact.DeletedAt = optional.IntoTime(args.DeletedAt)
		cols = append(cols, models.ArtifactColumns.DeletedAt)
	}

	_, err = artifact.Update(ctx, r.db, boil.Whitelist(cols...))
	if err != nil {
		return errors.Wrap(err, "failed to update artifact")
	}
	return nil
}

func (r *artifactRepository) HardDeleteArtifacts(ctx context.Context, cond domain.GetArtifactCondition) error {
	artifacts, err := models.Artifacts(r.buildMods(cond)...).All(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to get artifacts")
	}
	_, err = artifacts.DeleteAll(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to delete artifacts")
	}
	return nil
}
