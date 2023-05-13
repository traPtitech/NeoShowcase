package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/repoconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

type availableDomainRepository struct {
	db *sql.DB
}

func NewAvailableDomainRepository(db *sql.DB) domain.AvailableDomainRepository {
	return &availableDomainRepository{
		db: db,
	}
}

func (r *availableDomainRepository) GetAvailableDomains(ctx context.Context) (domain.AvailableDomainSlice, error) {
	domains, err := models.AvailableDomains().All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get available domains")
	}
	dDomains := ds.Map(domains, repoconvert.ToDomainAvailableDomain)
	return dDomains, nil
}

func (r *availableDomainRepository) AddAvailableDomain(ctx context.Context, ad *domain.AvailableDomain) error {
	mad := repoconvert.FromDomainAvailableDomain(ad)
	err := mad.Insert(ctx, r.db, boil.Blacklist())
	if err != nil {
		return fmt.Errorf("failed to insert available domain")
	}
	return nil
}

func (r *availableDomainRepository) DeleteAvailableDomain(ctx context.Context, domain string) error {
	d := models.AvailableDomain{Domain: domain}
	_, err := d.Delete(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to delete available domain")
	}
	return nil
}
