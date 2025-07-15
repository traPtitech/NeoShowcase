package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/repoconvert"
)

type runtimeImageRepository struct {
	db *sql.DB
}

func NewRuntimeImageRepository(db *sql.DB) domain.RuntimeImageRepository {
	return &runtimeImageRepository{
		db: db,
	}
}

func (r *runtimeImageRepository) CreateRuntimeImage(ctx context.Context, image *domain.RuntimeImage) error {
	ri := repoconvert.FromDomainRuntimeImage(image)
	err := ri.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return errors.Wrap(err, "failed to insert runtime image")
	}
	return nil
}

func (r *runtimeImageRepository) DeleteRuntimeImagesByAppID(ctx context.Context, appID string) error {
	images, err := models.RuntimeImages(
		qm.InnerJoin(fmt.Sprintf(
			"%s ON %s = %s",
			models.TableNames.Builds,
			models.BuildTableColumns.ID,
			models.RuntimeImageTableColumns.BuildID,
		)),
		models.BuildWhere.ApplicationID.EQ(appID),
	).All(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to get runtime images")
	}
	_, err = images.DeleteAll(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to delete runtime images")
	}
	return nil
}
