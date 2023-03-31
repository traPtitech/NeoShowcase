package repository

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

type environmentRepository struct {
	db *sql.DB
}

func NewEnvironmentRepository(db *sql.DB) domain.EnvironmentRepository {
	return &environmentRepository{db: db}
}

func (r *environmentRepository) buildMods(cond domain.GetEnvCondition) []qm.QueryMod {
	var mods []qm.QueryMod
	if cond.ApplicationID.Valid {
		mods = append(mods, models.EnvironmentWhere.ApplicationID.EQ(cond.ApplicationID.V))
	}
	if cond.Key.Valid {
		mods = append(mods, models.EnvironmentWhere.Key.EQ(cond.Key.V))
	}
	return mods
}

func (r *environmentRepository) GetEnv(ctx context.Context, cond domain.GetEnvCondition) ([]*domain.Environment, error) {
	environments, err := models.Environments(r.buildMods(cond)...).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return lo.Map(environments, func(env *models.Environment, i int) *domain.Environment {
		return toDomainEnvironment(env)
	}), nil
}

func (r *environmentRepository) SetEnv(ctx context.Context, applicationID string, env *domain.Environment) error {
	me := fromDomainEnvironment(applicationID, env)
	if err := me.Upsert(ctx, r.db, boil.Infer(), boil.Infer()); err != nil {
		return errors.Wrap(err, "failed to upsert env")
	}
	return nil
}

func (r *environmentRepository) DeleteEnv(ctx context.Context, cond domain.GetEnvCondition) error {
	_, err := models.Environments(r.buildMods(cond)...).DeleteAll(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to delete env")
	}
	return nil
}
