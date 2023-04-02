package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/exp/slices"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

type applicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository(db *sql.DB) domain.ApplicationRepository {
	return &applicationRepository{
		db: db,
	}
}

func (r *applicationRepository) GetApplications(ctx context.Context, cond domain.GetApplicationCondition) ([]*domain.Application, error) {
	mods := []qm.QueryMod{
		qm.Load(models.ApplicationRels.ApplicationConfig),
		qm.Load(models.ApplicationRels.Repository),
		qm.Load(models.ApplicationRels.Websites),
		qm.Load(models.ApplicationRels.Users),
	}

	if cond.IDIn.Valid {
		mods = append(mods, models.ApplicationWhere.ID.IN(cond.IDIn.V))
	}
	if cond.UserID.Valid {
		mods = append(mods,
			qm.InnerJoin(fmt.Sprintf(
				"%s ON %s.application_id = %s",
				models.TableNames.ApplicationOwners,
				models.TableNames.ApplicationOwners,
				models.ApplicationTableColumns.ID,
			)),
			qm.Where(fmt.Sprintf("%s.user_id = ?", models.TableNames.ApplicationOwners), cond.UserID.V),
		)
	}
	if cond.BuildType.Valid {
		mods = append(mods, models.ApplicationWhere.BuildType.EQ(cond.BuildType.V.String()))
	}
	if cond.State.Valid {
		mods = append(mods, models.ApplicationWhere.State.EQ(cond.State.V.String()))
	}
	if cond.InSync.Valid {
		if cond.InSync.V {
			mods = append(mods, qm.Where(fmt.Sprintf("%s == %s", models.ApplicationColumns.WantCommit, models.ApplicationColumns.CurrentCommit)))
		} else {
			mods = append(mods, qm.Where(fmt.Sprintf("%s != %s", models.ApplicationColumns.WantCommit, models.ApplicationColumns.CurrentCommit)))
		}
	}

	applications, err := models.Applications(mods...).All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get applications")
	}
	return lo.Map(applications, func(app *models.Application, i int) *domain.Application {
		return toDomainApplication(app)
	}), nil
}

func (r *applicationRepository) getApplication(ctx context.Context, id string, forUpdate bool, ex boil.ContextExecutor) (*models.Application, error) {
	mods := []qm.QueryMod{
		models.ApplicationWhere.ID.EQ(id),
		qm.Load(models.ApplicationRels.ApplicationConfig),
		qm.Load(models.ApplicationRels.Repository),
		qm.Load(models.ApplicationRels.Websites),
		qm.Load(models.ApplicationRels.Users),
	}
	if forUpdate {
		mods = append(mods, qm.For("UPDATE"))
	}
	app, err := models.Applications(mods...).One(ctx, ex)
	if err != nil {
		if isNoRowsErr(err) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to get application")
	}
	return app, nil
}

func (r *applicationRepository) GetApplication(ctx context.Context, id string) (*domain.Application, error) {
	app, err := r.getApplication(ctx, id, false, r.db)
	if err != nil {
		return nil, err
	}
	return toDomainApplication(app), nil
}

func (r *applicationRepository) CreateApplication(ctx context.Context, app *domain.Application) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()

	ma := fromDomainApplication(app)
	if err = ma.Insert(ctx, tx, boil.Infer()); err != nil {
		return errors.Wrap(err, "failed to create application")
	}

	mc := fromDomainApplicationConfig(app.ID, &app.Config)
	if err = ma.SetApplicationConfig(ctx, tx, true, mc); err != nil {
		return fmt.Errorf("failed to set application config")
	}

	err = r.validateAndInsertWebsites(ctx, tx, ma, app.Websites)
	if err != nil {
		return err
	}

	err = r.insertOwners(ctx, tx, ma, app.OwnerIDs)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (r *applicationRepository) UpdateApplication(ctx context.Context, id string, args *domain.UpdateApplicationArgs) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()

	app, err := r.getApplication(ctx, id, true, tx)
	if err != nil {
		return err
	}

	if args.Name.Valid {
		app.Name = args.Name.V
	}
	if args.RefName.Valid {
		app.RefName = args.RefName.V
	}
	if args.State.Valid {
		app.State = args.State.V.String()
	}
	if args.CurrentCommit.Valid {
		app.CurrentCommit = args.CurrentCommit.V
	}
	if args.WantCommit.Valid {
		app.WantCommit = args.WantCommit.V
	}

	app.UpdatedAt = time.Now()
	_, err = app.Update(ctx, tx, boil.Infer())

	if args.Config.Valid {
		err = app.SetApplicationConfig(ctx, tx, false, fromDomainApplicationConfig(app.ID, &args.Config.V))
		if err != nil {
			return errors.Wrap(err, "failed to update config")
		}
	}
	if args.Websites.Valid {
		_, err = app.R.Websites.DeleteAll(ctx, tx)
		if err != nil {
			return errors.Wrap(err, "failed to delete all websites")
		}
		err = r.validateAndInsertWebsites(ctx, tx, app, args.Websites.V)
		if err != nil {
			return err
		}
	}
	if args.OwnerIDs.Valid {
		_, err = app.R.Users.DeleteAll(ctx, tx)
		if err != nil {
			return errors.Wrap(err, "failed to delete all owners")
		}
		err = r.insertOwners(ctx, tx, app, args.OwnerIDs.V)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return err
}

func (r *applicationRepository) DeleteApplication(ctx context.Context, id string) error {
	app, err := r.getApplication(ctx, id, false, r.db)
	if err != nil {
		return err
	}
	err = app.SetUsers(ctx, r.db, false)
	if err != nil {
		return errors.Wrap(err, "failed to delete application owners")
	}
	_, err = app.R.Websites.DeleteAll(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to delete websites")
	}
	_, err = app.R.ApplicationConfig.Delete(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to delete application config")
	}
	_, err = app.Delete(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to delete application")
	}
	return nil
}

func (r *applicationRepository) validateAndInsertWebsites(ctx context.Context, ex boil.ContextExecutor, app *models.Application, websites []*domain.Website) error {
	// Validate and insert from the most specific
	slices.SortFunc(websites, func(a, b *domain.Website) bool { return len(b.PathPrefix) < len(a.PathPrefix) })
	for _, website := range websites {
		if err := r.validateAndInsertWebsite(ctx, ex, app, website); err != nil {
			return err
		}
	}
	return nil
}

func (r *applicationRepository) validateAndInsertWebsite(ctx context.Context, ex boil.ContextExecutor, app *models.Application, website *domain.Website) error {
	websites, err := models.Websites(models.WebsiteWhere.FQDN.EQ(website.FQDN), qm.For("UPDATE")).All(ctx, ex)
	if err != nil {
		return errors.Wrap(err, "failed to get existing websites")
	}
	existing := lo.Map(websites, func(website *models.Website, i int) *domain.Website { return toDomainWebsite(website) })
	if website.ConflictsWith(existing) {
		return ErrDuplicate
	}
	err = app.AddWebsites(ctx, ex, true, fromDomainWebsite(website))
	if err != nil {
		return errors.Wrap(err, "failed to add website")
	}
	return nil
}

func (r *applicationRepository) insertOwners(ctx context.Context, ex boil.ContextExecutor, app *models.Application, ownerIDs []string) error {
	ownerIDs = lo.Uniq(ownerIDs)
	users, err := models.Users(models.UserWhere.ID.IN(ownerIDs)).All(ctx, ex)
	if len(users) < len(ownerIDs) {
		return ErrNotFound
	}
	err = app.AddUsers(ctx, ex, false, users...)
	if err != nil {
		return errors.Wrap(err, "failed to add owners")
	}
	return nil
}
