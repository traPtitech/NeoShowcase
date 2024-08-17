package repository

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/repoconvert"
	"github.com/volatiletech/sqlboiler/v4/boil"
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
		errors.Wrap(err, "failed to insert runtime image")
	}
	return nil
}
