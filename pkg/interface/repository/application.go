package repository

import (
	"context"
	"database/sql"
	"fmt"

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
	}

	if cond.UserID.Valid {
		mods = append(mods,
			qm.InnerJoin(fmt.Sprintf(
				"%s ON %s.%s = %s.%s",
				models.TableNames.Owners,
				models.TableNames.Owners,
				"application_id",
				models.TableNames.Applications,
				models.ApplicationColumns.ID,
			)),
			qm.Where(fmt.Sprintf("%s.%s = ?", models.TableNames.Owners, "user_id"), cond.UserID.V),
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
		return nil, fmt.Errorf("failed to get applications: %w", err)
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
	).One(ctx, r.db)
	if err != nil {
		if isNoRowsErr(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
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

func (r *applicationRepository) CreateApplication(ctx context.Context, args domain.CreateApplicationArgs) (*domain.Application, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	app := &models.Application{
		ID:            domain.NewID(),
		Name:          args.Name,
		RepositoryID:  args.RepositoryID,
		BranchName:    args.BranchName,
		BuildType:     args.BuildType.String(),
		State:         args.State.String(),
		CurrentCommit: domain.EmptyCommit,
		WantCommit:    domain.EmptyCommit,
	}
	if err := app.Insert(ctx, tx, boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	config := &models.ApplicationConfig{
		ApplicationID:  app.ID,
		UseMariadb:     args.Config.UseMariaDB,
		UseMongodb:     args.Config.UseMongoDB,
		BaseImage:      args.Config.BaseImage,
		DockerfileName: args.Config.DockerfileName,
		ArtifactPath:   args.Config.ArtifactPath,
		BuildCMD:       args.Config.BuildCmd,
		EntrypointCMD:  args.Config.EntrypointCmd,
		Authentication: args.Config.Authentication.String(),
	}
	if err := app.SetApplicationConfig(ctx, tx, true, config); err != nil {
		return nil, fmt.Errorf("failed to set application config")
	}

	// Validate and insert from the most specific
	slices.SortFunc(args.Websites, func(a, b *domain.Website) bool { return len(b.PathPrefix) < len(a.PathPrefix) })
	for _, website := range args.Websites {
		if err = r.validateAndInsertWebsite(ctx, tx, app, website); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.GetApplication(ctx, app.ID)
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
		return fmt.Errorf("failed to find user: %w", err)
	}
	if err := app.AddUsers(ctx, r.db, false, user); err != nil {
		return fmt.Errorf("failed to register owner: %w", err)
	}
	return nil
}

func (r *applicationRepository) GetWebsites(ctx context.Context, applicationIDs []string) ([]*domain.Website, error) {
	websites, err := models.Websites(models.WebsiteWhere.ApplicationID.IN(applicationIDs)).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get websites: %w", err)
	}
	return lo.Map(websites, func(website *models.Website, i int) *domain.Website {
		return toDomainWebsite(website)
	}), nil
}

func (r *applicationRepository) validateAndInsertWebsite(ctx context.Context, ex boil.ContextExecutor, app *models.Application, website *domain.Website) error {
	websites, err := models.Websites(models.WebsiteWhere.FQDN.EQ(website.FQDN), qm.For("UPDATE")).All(ctx, ex)
	if err != nil {
		return fmt.Errorf("failed to get existing websites: %w", err)
	}
	existing := lo.Map(websites, func(website *models.Website, i int) *domain.Website { return toDomainWebsite(website) })
	if website.ConflictsWith(existing) {
		return ErrDuplicate
	}
	err = app.AddWebsites(ctx, ex, true, fromDomainWebsite(website))
	if err != nil {
		return fmt.Errorf("failed to add website: %w", err)
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
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	err = r.validateAndInsertWebsite(ctx, tx, app, website)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
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
