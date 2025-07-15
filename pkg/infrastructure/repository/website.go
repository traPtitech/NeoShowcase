package repository

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/repoconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

type websiteRepository struct {
	db *sql.DB
}

func NewWebsiteRepository(db *sql.DB) domain.WebsiteRepository {
	return &websiteRepository{
		db: db,
	}
}

func (w *websiteRepository) GetWebsites(ctx context.Context) ([]*domain.Website, error) {
	websites, err := models.Websites().All(ctx, w.db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get websites")
	}
	return ds.Map(websites, repoconvert.ToDomainWebsite), nil
}
