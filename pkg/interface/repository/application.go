package repository

import (
	"context"
	"database/sql"
	"fmt"

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
				"%s ON %s.%s = %s.%s",
				models.TableNames.ApplicationOwners,
				models.TableNames.ApplicationOwners,
				"application_id",
				models.TableNames.Applications,
				models.ApplicationColumns.ID,
			)),
			qm.Where(fmt.Sprintf("%s.%s = ?", models.TableNames.ApplicationOwners, "user_id"), cond.UserID.V),
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

func (r *applicationRepository) getApplication(ctx context.Context, id string) (*models.Application, error) {
	app, err := models.Applications(
		models.ApplicationWhere.ID.EQ(id),
		qm.Load(models.ApplicationRels.ApplicationConfig),
		qm.Load(models.ApplicationRels.Repository),
		qm.Load(models.ApplicationRels.Websites),
		qm.Load(models.ApplicationRels.Users),
	).One(ctx, r.db)
	if err != nil {
		if isNoRowsErr(err) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to get application")
	}
	return app, nil
}

func (r *applicationRepository) GetApplication(ctx context.Context, id string) (*domain.Application, error) {
	app, err := r.getApplication(ctx, id)
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

	// Validate and insert from the most specific
	slices.SortFunc(app.Websites, func(a, b *domain.Website) bool { return len(b.PathPrefix) < len(a.PathPrefix) })
	for _, website := range app.Websites {
		if err = r.validateAndInsertWebsite(ctx, tx, ma, website); err != nil {
			return err
		}
	}

	app.OwnerIDs = lo.Uniq(app.OwnerIDs)
	users, err := models.Users(models.UserWhere.ID.IN(app.OwnerIDs)).All(ctx, tx)
	if len(users) < len(app.OwnerIDs) {
		return ErrNotFound
	}
	err = ma.AddUsers(ctx, tx, false, users...)
	if err != nil {
		return errors.Wrap(err, "failed to add owners")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (r *applicationRepository) UpdateApplication(ctx context.Context, id string, args domain.UpdateApplicationArgs) error {
	app, err := r.getApplication(ctx, id)
	if err != nil {
		return err
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

	_, err = app.Update(ctx, r.db, boil.Infer())
	return err
}

func (r *applicationRepository) RegisterApplicationOwner(ctx context.Context, applicationID string, userID string) error {
	app, err := r.getApplication(ctx, applicationID)
	if err != nil {
		return err
	}
	user, err := models.Users(models.UserWhere.ID.EQ(userID)).One(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to find user")
	}
	if err := app.AddUsers(ctx, r.db, false, user); err != nil {
		return errors.Wrap(err, "failed to register owner")
	}
	return nil
}

func (r *applicationRepository) GetWebsites(ctx context.Context, applicationIDs []string) ([]*domain.Website, error) {
	websites, err := models.Websites(models.WebsiteWhere.ApplicationID.IN(applicationIDs)).All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get websites")
	}
	return lo.Map(websites, func(website *models.Website, i int) *domain.Website {
		return toDomainWebsite(website)
	}), nil
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

func (r *applicationRepository) AddWebsite(ctx context.Context, applicationID string, website *domain.Website) error {
	app, err := r.getApplication(ctx, applicationID)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()

	err = r.validateAndInsertWebsite(ctx, tx, app, website)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (r *applicationRepository) DeleteWebsite(ctx context.Context, applicationID string, websiteID string) error {
	app, err := r.getApplication(ctx, applicationID)
	if err != nil {
		return err
	}

	for _, website := range app.R.Websites {
		if website.ID == websiteID {
			_, err := website.Delete(ctx, r.db)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return ErrNotFound
}
