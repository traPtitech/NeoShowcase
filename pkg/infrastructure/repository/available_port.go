package repository

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/repoconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

type availablePortRepository struct {
	db *sql.DB
}

func NewAvailablePortRepository(db *sql.DB) domain.AvailablePortRepository {
	return &availablePortRepository{
		db: db,
	}
}

func (r *availablePortRepository) GetAvailablePorts(ctx context.Context) (domain.AvailablePortSlice, error) {
	ports, err := models.AvailablePorts().All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "getting available ports")
	}
	dPorts := ds.Map(ports, repoconvert.ToDomainAvailablePort)
	return dPorts, nil
}

func (r *availablePortRepository) AddAvailablePort(ctx context.Context, ap *domain.AvailablePort) error {
	modelAp := repoconvert.FromDomainAvailablePort(ap)
	err := modelAp.Insert(ctx, r.db, boil.Blacklist())
	if err != nil {
		return errors.New("inserting available port")
	}
	return nil
}

func (r *availablePortRepository) DeleteAvailablePort(ctx context.Context, ap *domain.AvailablePort) error {
	modelAp := repoconvert.FromDomainAvailablePort(ap)
	_, err := modelAp.Delete(ctx, r.db)
	if err != nil {
		return errors.New("deleting available port")
	}
	return nil
}
