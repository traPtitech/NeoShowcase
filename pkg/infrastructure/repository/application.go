package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/repoconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
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
		qm.Load(models.ApplicationRels.PortPublications),
		qm.Load(models.ApplicationRels.Users),
	}

	if cond.IDIn.Valid {
		mods = append(mods, models.ApplicationWhere.ID.IN(cond.IDIn.V))
	}
	if cond.RepositoryID.Valid {
		mods = append(mods, models.ApplicationWhere.RepositoryID.EQ(cond.RepositoryID.V))
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
	if cond.DeployType.Valid {
		mods = append(mods, models.ApplicationWhere.DeployType.EQ(repoconvert.DeployTypeMapper.FromMust(cond.DeployType.V)))
	}
	if cond.Running.Valid {
		mods = append(mods, models.ApplicationWhere.Running.EQ(cond.Running.V))
	}
	if cond.InSync.Valid {
		if cond.InSync.V {
			mods = append(mods, qm.Where(fmt.Sprintf("%s == %s", models.ApplicationTableColumns.WantCommit, models.ApplicationTableColumns.CurrentCommit)))
		} else {
			mods = append(mods, qm.Where(fmt.Sprintf("%s != %s", models.ApplicationTableColumns.WantCommit, models.ApplicationTableColumns.CurrentCommit)))
		}
	}

	applications, err := models.Applications(mods...).All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get applications")
	}
	return ds.Map(applications, repoconvert.ToDomainApplication), nil
}

func (r *applicationRepository) getApplication(ctx context.Context, id string, forUpdate bool, ex boil.ContextExecutor) (*models.Application, error) {
	mods := []qm.QueryMod{
		models.ApplicationWhere.ID.EQ(id),
		qm.Load(models.ApplicationRels.ApplicationConfig),
		qm.Load(models.ApplicationRels.Repository),
		qm.Load(models.ApplicationRels.Websites),
		qm.Load(models.ApplicationRels.PortPublications),
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
	return repoconvert.ToDomainApplication(app), nil
}

func (r *applicationRepository) CreateApplication(ctx context.Context, app *domain.Application) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()

	ma := repoconvert.FromDomainApplication(app)
	if err = ma.Insert(ctx, tx, boil.Blacklist()); err != nil {
		return errors.Wrap(err, "failed to create application")
	}

	mc := repoconvert.FromDomainApplicationConfig(app.ID, &app.Config)
	err = mc.Insert(ctx, tx, boil.Blacklist())
	if err != nil {
		return fmt.Errorf("failed to create application config")
	}

	err = r.setWebsites(ctx, tx, ma, app.Websites)
	if err != nil {
		return err
	}

	err = r.setPortPublications(ctx, tx, ma, app.PortPublications)
	if err != nil {
		return err
	}

	err = r.setOwners(ctx, tx, ma, app.OwnerIDs)
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

	var cols []string
	if args.Name.Valid {
		app.Name = args.Name.V
		cols = append(cols, models.ApplicationColumns.Name)
	}
	if args.RepositoryID.Valid {
		app.RepositoryID = args.RepositoryID.V
		cols = append(cols, models.ApplicationColumns.RepositoryID)
	}
	if args.RefName.Valid {
		app.RefName = args.RefName.V
		cols = append(cols, models.ApplicationColumns.RefName)
	}
	if args.Config.Valid {
		app.DeployType = repoconvert.DeployTypeMapper.FromMust(
			args.Config.V.BuildConfig.BuildType().DeployType(),
		)
		cols = append(cols, models.ApplicationColumns.DeployType)
	}
	if args.Running.Valid {
		app.Running = args.Running.V
		cols = append(cols, models.ApplicationColumns.Running)
	}
	if args.Container.Valid {
		app.Container = repoconvert.ContainerStateMapper.FromMust(args.Container.V)
		cols = append(cols, models.ApplicationColumns.Container)
	}
	if args.CurrentCommit.Valid {
		app.CurrentCommit = args.CurrentCommit.V
		cols = append(cols, models.ApplicationColumns.CurrentCommit)
	}
	if args.WantCommit.Valid {
		app.WantCommit = args.WantCommit.V
		cols = append(cols, models.ApplicationColumns.WantCommit)
	}
	if args.UpdatedAt.Valid {
		app.UpdatedAt = args.UpdatedAt.V
		cols = append(cols, models.ApplicationColumns.UpdatedAt)
	}

	if len(cols) > 0 {
		_, err = app.Update(ctx, tx, boil.Whitelist(cols...))
		if err != nil {
			return errors.Wrap(err, "updating application")
		}
	}

	if args.Config.Valid {
		mac := repoconvert.FromDomainApplicationConfig(app.ID, &args.Config.V)
		err = mac.Upsert(ctx, tx, boil.Blacklist(), boil.Blacklist())
		if err != nil {
			return errors.Wrap(err, "failed to update config")
		}
	}
	if args.Websites.Valid {
		err = r.setWebsites(ctx, tx, app, args.Websites.V)
		if err != nil {
			return err
		}
	}
	if args.PortPublications.Valid {
		err = r.setPortPublications(ctx, tx, app, args.PortPublications.V)
		if err != nil {
			return err
		}
	}
	if args.OwnerIDs.Valid {
		err = r.setOwners(ctx, tx, app, args.OwnerIDs.V)
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

func (r *applicationRepository) BulkUpdateState(ctx context.Context, m map[string]domain.ContainerState) error {
	// NOTE: sqlboiler does not support bulk insert/update by default, could use custom templating
	for appID, state := range m {
		ma := models.Application{ID: appID, Container: repoconvert.ContainerStateMapper.FromMust(state)}
		_, err := ma.Update(ctx, r.db, boil.Whitelist(models.ApplicationColumns.Container))
		if err != nil {
			return errors.Wrap(err, "failed to update container state")
		}
	}
	return nil
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
	_, err = app.R.PortPublications.DeleteAll(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to delete port publications")
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

func (r *applicationRepository) setWebsites(ctx context.Context, ex boil.ContextExecutor, app *models.Application, websites []*domain.Website) error {
	if app.R != nil && app.R.Websites != nil {
		_, err := app.R.Websites.DeleteAll(ctx, ex)
		if err != nil {
			return errors.Wrap(err, "failed to delete existing app websites")
		}
	}
	for _, w := range websites {
		mw := repoconvert.FromDomainWebsite(app.ID, w)
		err := mw.Insert(ctx, ex, boil.Blacklist())
		if err != nil {
			return errors.Wrap(err, "failed to insert website")
		}
	}
	return nil
}

func (r *applicationRepository) setPortPublications(ctx context.Context, ex boil.ContextExecutor, app *models.Application, ports []*domain.PortPublication) error {
	if app.R != nil && app.R.PortPublications != nil {
		_, err := app.R.PortPublications.DeleteAll(ctx, ex)
		if err != nil {
			return errors.Wrap(err, "deleting existing port publications")
		}
	}
	for _, p := range ports {
		mp := repoconvert.FromDomainPortPublication(app.ID, p)
		err := mp.Insert(ctx, ex, boil.Blacklist())
		if err != nil {
			return errors.Wrap(err, "inserting port publication")
		}
	}
	return nil
}

func (r *applicationRepository) setOwners(ctx context.Context, ex boil.ContextExecutor, app *models.Application, ownerIDs []string) error {
	ownerIDs = lo.Uniq(ownerIDs)
	users, err := models.Users(models.UserWhere.ID.IN(ownerIDs)).All(ctx, ex)
	if len(users) < len(ownerIDs) {
		return ErrNotFound
	}
	err = app.SetUsers(ctx, ex, false, users...)
	if err != nil {
		return errors.Wrap(err, "failed to add owners")
	}
	return nil
}
