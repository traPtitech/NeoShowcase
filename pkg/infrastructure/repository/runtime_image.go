package repository

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/traPtitech/neoshowcase/pkg/domain"
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
	return errors.New("not implemented")
}
