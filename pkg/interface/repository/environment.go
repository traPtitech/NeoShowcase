package repository

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

type environmentRepository struct {
	db *sql.DB
}

func NewEnvironmentRepository(db *sql.DB) domain.EnvironmentRepository {
	return &environmentRepository{db: db}
}

func (r *environmentRepository) GetEnv(ctx context.Context, applicationID string) ([]*domain.Environment, error) {
	environments, err := models.Environments(
		models.EnvironmentWhere.ApplicationID.EQ(applicationID),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return lo.Map(environments, func(env *models.Environment, i int) *domain.Environment {
		return toDomainEnvironment(env)
	}), nil
}

func (r *environmentRepository) SetEnv(ctx context.Context, applicationID, key, value string) error {
	env := models.Environment{
		ID:            domain.NewID(),
		ApplicationID: applicationID,
		Key:           key,
		Value:         value,
	}
	if err := env.Upsert(ctx, r.db, boil.Infer(), boil.Infer()); err != nil {
		return errors.Wrap(err, "failed to upsert env")
	}
	return nil
}
