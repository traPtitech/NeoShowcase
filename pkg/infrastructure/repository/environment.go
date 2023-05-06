package repository

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/repoconvert"
)

type environmentRepository struct {
	db *sql.DB
}

func NewEnvironmentRepository(db *sql.DB) domain.EnvironmentRepository {
	return &environmentRepository{db: db}
}

func (r *environmentRepository) buildMods(cond domain.GetEnvCondition) []qm.QueryMod {
	var mods []qm.QueryMod
	if cond.ApplicationIDIn.Valid {
		mods = append(mods, models.EnvironmentWhere.ApplicationID.IN(cond.ApplicationIDIn.V))
	}
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
		return repoconvert.ToDomainEnvironment(env)
	}), nil
}

func (r *environmentRepository) SetEnv(ctx context.Context, env *domain.Environment) error {
	// NOTE: sqlboiler does not recognize multiple column unique keys: https://github.com/volatiletech/sqlboiler/issues/328
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}

	me, err := models.Environments(
		models.EnvironmentWhere.ApplicationID.EQ(env.ApplicationID),
		models.EnvironmentWhere.Key.EQ(env.Key),
		qm.For("UPDATE"),
	).One(ctx, tx)
	if err != nil && !isNoRowsErr(err) {
		return errors.Wrap(err, "failed to get environment")
	}
	exists := !isNoRowsErr(err)

	me = repoconvert.FromDomainEnvironment(env)

	if exists {
		_, err = me.Update(ctx, tx, boil.Blacklist())
	} else {
		err = me.Insert(ctx, tx, boil.Blacklist())
	}
	if err != nil {
		return errors.Wrap(err, "failed to upsert environment")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit")
	}

	return nil
}

func (r *environmentRepository) DeleteEnv(ctx context.Context, cond domain.GetEnvCondition) error {
	envs, err := models.Environments(r.buildMods(cond)...).All(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to get env")
	}
	_, err = envs.DeleteAll(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to delete env")
	}
	return nil
}
