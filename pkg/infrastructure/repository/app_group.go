package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/samber/lo"
	"github.com/samber/oops"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/repoconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

type applicationGroupRepository struct {
	db *sql.DB
}

func NewApplicationGroupRepository(db *sql.DB) domain.ApplicationGroupRepository {
	return &applicationGroupRepository{
		db: db,
	}
}

func (r *applicationGroupRepository) GetApplicationGroups(ctx context.Context, cond domain.GetApplicationGroupCondition) ([]*domain.ApplicationGroup, error) {
	mods := []qm.QueryMod{
		qm.Load(models.ApplicationGroupRels.Users),
	}
	if cond.IDIn.Valid {
		mods = append(mods, models.ApplicationGroupWhere.ID.IN(cond.IDIn.V))
	}
	if cond.UserID.Valid {
		mods = append(mods,
			qm.InnerJoin(fmt.Sprintf(
				"%s ON %s.group_id = %s",
				models.TableNames.ApplicationGroupOwners,
				models.TableNames.ApplicationGroupOwners,
				models.ApplicationGroupTableColumns.ID,
			)),
			qm.Where(fmt.Sprintf("%s.user_id = ?", models.TableNames.ApplicationGroupOwners), cond.UserID.V),
		)
	}
	if cond.Implicit.Valid {
		mods = append(mods, models.ApplicationGroupWhere.Implicit.EQ(cond.Implicit.V))
	}

	groups, err := models.ApplicationGroups(mods...).All(ctx, r.db)
	if err != nil {
		return nil, oops.Wrapf(err, "getting application groups")
	}
	return ds.Map(groups, repoconvert.ToDomainApplicationGroup), nil
}

func (r *applicationGroupRepository) getApplicationGroup(ctx context.Context, ex boil.ContextExecutor, id string) (*models.ApplicationGroup, error) {
	group, err := models.ApplicationGroups(
		models.ApplicationGroupWhere.ID.EQ(id),
		qm.Load(models.ApplicationGroupRels.Users),
	).One(ctx, ex)
	if err != nil {
		if isNoRowsErr(err) {
			return nil, notFound("application group")
		}
		return nil, oops.Wrapf(err, "getting application group")
	}
	return group, nil
}

func (r *applicationGroupRepository) GetApplicationGroup(ctx context.Context, id string) (*domain.ApplicationGroup, error) {
	group, err := r.getApplicationGroup(ctx, r.db, id)
	if err != nil {
		return nil, err
	}
	return repoconvert.ToDomainApplicationGroup(group), nil
}

func (r *applicationGroupRepository) CreateApplicationGroup(ctx context.Context, group *domain.ApplicationGroup) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return oops.Wrapf(err, "starting transaction")
	}
	defer tx.Rollback()

	if err = insertApplicationGroup(ctx, tx, group); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return oops.Wrapf(err, "committing transaction")
	}
	return nil
}

func (r *applicationGroupRepository) UpdateApplicationGroup(ctx context.Context, id string, args *domain.UpdateApplicationGroupArgs) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return oops.Wrapf(err, "starting transaction")
	}
	defer tx.Rollback()

	group, err := r.getApplicationGroup(ctx, tx, id)
	if err != nil {
		return err
	}

	if args.Name.Valid {
		group.Name = args.Name.V
		_, err = group.Update(ctx, tx, boil.Whitelist(models.ApplicationGroupColumns.Name))
		if err != nil {
			return oops.Wrapf(err, "updating application group")
		}
	}
	if args.OwnerIDs.Valid {
		if err = setApplicationGroupOwners(ctx, tx, group, args.OwnerIDs.V); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return oops.Wrapf(err, "committing transaction")
	}
	return nil
}

func (r *applicationGroupRepository) DeleteApplicationGroup(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return oops.Wrapf(err, "starting transaction")
	}
	defer tx.Rollback()

	group, err := r.getApplicationGroup(ctx, tx, id)
	if err != nil {
		return err
	}
	if err = deleteApplicationGroup(ctx, tx, group); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return oops.Wrapf(err, "committing transaction")
	}
	return nil
}

func insertApplicationGroup(ctx context.Context, ex boil.ContextExecutor, group *domain.ApplicationGroup) error {
	mg := repoconvert.FromDomainApplicationGroup(group)
	if err := mg.Insert(ctx, ex, boil.Blacklist()); err != nil {
		return oops.Wrapf(err, "inserting application group")
	}
	return setApplicationGroupOwners(ctx, ex, mg, group.OwnerIDs)
}

func deleteApplicationGroup(ctx context.Context, ex boil.ContextExecutor, group *models.ApplicationGroup) error {
	if err := group.SetUsers(ctx, ex, false); err != nil {
		return oops.Wrapf(err, "deleting application group owners")
	}
	if _, err := group.Delete(ctx, ex); err != nil {
		return oops.Wrapf(err, "deleting application group")
	}
	return nil
}

func setApplicationGroupOwners(ctx context.Context, ex boil.ContextExecutor, group *models.ApplicationGroup, ownerIDs []string) error {
	ownerIDs = lo.Uniq(ownerIDs)
	users, err := models.Users(models.UserWhere.ID.IN(ownerIDs)).All(ctx, ex)
	if err != nil {
		return oops.Wrapf(err, "getting owners")
	}
	if len(users) < len(ownerIDs) {
		return notFound("user")
	}
	if err = group.SetUsers(ctx, ex, false, users...); err != nil {
		return oops.Wrapf(err, "setting application group owners")
	}
	return nil
}
