package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/samber/oops"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/repoconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
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
	}

	if cond.IDIn.Valid {
		mods = append(mods, models.ApplicationWhere.ID.IN(cond.IDIn.V))
	}
	if cond.RepositoryID.Valid {
		mods = append(mods, models.ApplicationWhere.RepositoryID.EQ(cond.RepositoryID.V))
	}
	if cond.GroupID.Valid {
		mods = append(mods, models.ApplicationWhere.GroupID.EQ(null.StringFrom(cond.GroupID.V)))
	}
	if cond.UserID.Valid {
		mods = append(mods,
			qm.InnerJoin(fmt.Sprintf(
				"%s ON %s.group_id = %s",
				models.TableNames.ApplicationGroupOwners,
				models.TableNames.ApplicationGroupOwners,
				models.ApplicationTableColumns.GroupID,
			)),
			qm.Where(fmt.Sprintf("%s.user_id = ?", models.TableNames.ApplicationGroupOwners), cond.UserID.V),
		)
	}
	if cond.DeployType.Valid {
		mods = append(mods, models.ApplicationWhere.DeployType.EQ(repoconvert.DeployTypeMapper.FromMust(cond.DeployType.V)))
	}
	if cond.Running.Valid {
		mods = append(mods, models.ApplicationWhere.Running.EQ(cond.Running.V))
	}

	applications, err := models.Applications(mods...).All(ctx, r.db)
	if err != nil {
		return nil, oops.Wrapf(err, "getting applications")
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
	}
	if forUpdate {
		mods = append(mods, qm.For("UPDATE"))
	}
	app, err := models.Applications(mods...).One(ctx, ex)
	if err != nil {
		if isNoRowsErr(err) {
			return nil, notFound("application")
		}
		return nil, oops.Wrapf(err, "getting application")
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

func (r *applicationRepository) CreateApplication(ctx context.Context, app *domain.Application, newImplicitGroupOwnerIDs optional.Of[[]string]) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return oops.Wrapf(err, "starting transaction")
	}
	defer tx.Rollback()

	if err = insertImplicitApplicationGroup(ctx, tx, app.GroupID, newImplicitGroupOwnerIDs); err != nil {
		return err
	}
	ma := repoconvert.FromDomainApplication(app)
	if err = ma.Insert(ctx, tx, boil.Blacklist()); err != nil {
		return oops.Wrapf(err, "creating application")
	}

	mc := repoconvert.FromDomainApplicationConfig(app.ID, &app.Config)
	err = mc.Insert(ctx, tx, boil.Blacklist())
	if err != nil {
		return oops.Wrapf(err, "creating application config")
	}

	err = r.setWebsites(ctx, tx, ma, app.Websites)
	if err != nil {
		return err
	}

	err = r.setPortPublications(ctx, tx, ma, app.PortPublications)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return oops.Wrapf(err, "committing transaction")
	}

	return nil
}

func (r *applicationRepository) UpdateApplication(ctx context.Context, id string, args *domain.UpdateApplicationArgs) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return oops.Wrapf(err, "starting transaction")
	}
	defer tx.Rollback()

	app, err := r.getApplication(ctx, id, true, tx)
	if err != nil {
		return err
	}

	if args.GroupOwnerIDs.Valid {
		group, err := models.FindApplicationGroup(ctx, tx, app.GroupID.String)
		if err != nil {
			return oops.Wrapf(err, "getting application group")
		}
		if err = setApplicationGroupOwners(ctx, tx, group, args.GroupOwnerIDs.V); err != nil {
			return err
		}
	}

	var cols []string
	if args.Name.Valid {
		app.Name = args.Name.V
		cols = append(cols, models.ApplicationColumns.Name)
	}
	if args.RefName.Valid {
		app.RefName = args.RefName.V
		cols = append(cols, models.ApplicationColumns.RefName)
	}
	if args.Commit.Valid {
		app.Commit = args.Commit.V
		cols = append(cols, models.ApplicationColumns.Commit)
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
	if args.ContainerMessage.Valid {
		app.ContainerMessage = args.ContainerMessage.V
		cols = append(cols, models.ApplicationColumns.ContainerMessage)
	}
	if args.CurrentBuild.Valid {
		app.CurrentBuild = args.CurrentBuild.V
		cols = append(cols, models.ApplicationColumns.CurrentBuild)
	}
	if args.UpdatedAt.Valid {
		app.UpdatedAt = args.UpdatedAt.V
		cols = append(cols, models.ApplicationColumns.UpdatedAt)
	}

	if len(cols) > 0 {
		_, err = app.Update(ctx, tx, boil.Whitelist(cols...))
		if err != nil {
			return oops.Wrapf(err, "updating application")
		}
	}

	if args.Config.Valid {
		mac := repoconvert.FromDomainApplicationConfig(app.ID, &args.Config.V)
		err = mac.Upsert(ctx, tx, boil.Blacklist(), boil.Blacklist())
		if err != nil {
			return oops.Wrapf(err, "updating config")
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
	err = tx.Commit()
	if err != nil {
		return oops.Wrapf(err, "committing transaction")
	}

	return err
}

func (r *applicationRepository) SetApplicationGroup(ctx context.Context, id string, groupID string, newImplicitGroupOwnerIDs optional.Of[[]string]) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return oops.Wrapf(err, "starting transaction")
	}
	defer tx.Rollback()

	app, err := r.getApplication(ctx, id, true, tx)
	if err != nil {
		return err
	}
	if err = insertImplicitApplicationGroup(ctx, tx, groupID, newImplicitGroupOwnerIDs); err != nil {
		return err
	}
	prevGroupID := app.GroupID
	app.GroupID = null.StringFrom(groupID)
	if _, err = app.Update(ctx, tx, boil.Whitelist(models.ApplicationColumns.GroupID)); err != nil {
		return oops.Wrapf(err, "updating application group")
	}
	if err = deleteImplicitApplicationGroup(ctx, tx, prevGroupID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return oops.Wrapf(err, "committing transaction")
	}
	return nil
}

// insertImplicitApplicationGroup inserts an implicit group of groupID with the owners, if the owners are given.
func insertImplicitApplicationGroup(ctx context.Context, ex boil.ContextExecutor, groupID string, ownerIDs optional.Of[[]string]) error {
	if !ownerIDs.Valid {
		return nil
	}
	return insertApplicationGroup(ctx, ex, &domain.ApplicationGroup{ID: groupID, Implicit: true, OwnerIDs: ownerIDs.V})
}

// deleteImplicitApplicationGroup deletes the group of groupID if it is implicit.
// Call it after the application in the group has left it.
func deleteImplicitApplicationGroup(ctx context.Context, ex boil.ContextExecutor, groupID null.String) error {
	if !groupID.Valid {
		return nil
	}
	group, err := models.ApplicationGroups(
		models.ApplicationGroupWhere.ID.EQ(groupID.String),
		models.ApplicationGroupWhere.Implicit.EQ(true),
	).One(ctx, ex)
	if err != nil {
		if isNoRowsErr(err) {
			return nil
		}
		return oops.Wrapf(err, "getting application group")
	}
	return deleteApplicationGroup(ctx, ex, group)
}

func (r *applicationRepository) BulkUpdateState(ctx context.Context, states []*domain.Container) error {
	// NOTE: sqlboiler does not support bulk insert/update by default, could use custom templating
	// Run all updates in one transaction: with autocommit, every UPDATE pays a commit and the total grows with the number of apps.
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return oops.Wrapf(err, "starting transaction")
	}
	defer tx.Rollback()

	for _, state := range states {
		ma := models.Application{
			ID:               state.ApplicationID,
			Container:        repoconvert.ContainerStateMapper.FromMust(state.State),
			ContainerMessage: state.Message,
		}
		_, err := ma.Update(ctx, tx, boil.Whitelist(
			models.ApplicationColumns.Container,
			models.ApplicationColumns.ContainerMessage,
		))
		if err != nil {
			return oops.Wrapf(err, "updating container state")
		}
	}

	if err = tx.Commit(); err != nil {
		return oops.Wrapf(err, "committing transaction")
	}
	return nil
}

func (r *applicationRepository) DeleteApplication(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return oops.Wrapf(err, "starting transaction")
	}
	defer tx.Rollback()

	app, err := r.getApplication(ctx, id, true, tx)
	if err != nil {
		return err
	}
	// application_owners keeps the owners before groups as the source of the group migration.
	// Its rows reference the application, so they have to go first.
	err = app.SetUsers(ctx, tx, false)
	if err != nil {
		return oops.Wrapf(err, "deleting legacy application owners")
	}
	_, err = app.R.PortPublications.DeleteAll(ctx, tx)
	if err != nil {
		return oops.Wrapf(err, "deleting port publications")
	}
	_, err = app.R.Websites.DeleteAll(ctx, tx)
	if err != nil {
		return oops.Wrapf(err, "deleting websites")
	}
	_, err = app.R.ApplicationConfig.Delete(ctx, tx)
	if err != nil {
		return oops.Wrapf(err, "deleting application config")
	}
	_, err = app.Delete(ctx, tx)
	if err != nil {
		return oops.Wrapf(err, "deleting application")
	}
	if err = deleteImplicitApplicationGroup(ctx, tx, app.GroupID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return oops.Wrapf(err, "committing transaction")
	}
	return nil
}

func (r *applicationRepository) setWebsites(ctx context.Context, ex boil.ContextExecutor, app *models.Application, websites []*domain.Website) error {
	if app.R != nil && app.R.Websites != nil {
		_, err := app.R.Websites.DeleteAll(ctx, ex)
		if err != nil {
			return oops.Wrapf(err, "deleting existing app websites")
		}
	}
	for _, w := range websites {
		mw := repoconvert.FromDomainWebsite(app.ID, w)
		err := mw.Insert(ctx, ex, boil.Blacklist())
		if err != nil {
			return oops.Wrapf(err, "inserting website")
		}
	}
	return nil
}

func (r *applicationRepository) setPortPublications(ctx context.Context, ex boil.ContextExecutor, app *models.Application, ports []*domain.PortPublication) error {
	if app.R != nil && app.R.PortPublications != nil {
		_, err := app.R.PortPublications.DeleteAll(ctx, ex)
		if err != nil {
			return oops.Wrapf(err, "deleting existing port publications")
		}
	}
	for _, p := range ports {
		mp := repoconvert.FromDomainPortPublication(app.ID, p)
		err := mp.Insert(ctx, ex, boil.Blacklist())
		if err != nil {
			return oops.Wrapf(err, "inserting port publication")
		}
	}
	return nil
}
