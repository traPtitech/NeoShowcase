package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE

type GetApplicationCondition struct {
	UserID    optional.Of[string]
	BuildType optional.Of[builder.BuildType]
	// InSync WantCommit が CurrentCommit に一致する
	InSync optional.Of[bool]
}

type CreateApplicationArgs struct {
	RepositoryID string
	BranchName   string
	BuildType    builder.BuildType
}

type UpdateApplicationArgs struct {
	State         optional.Of[domain.ApplicationState]
	CurrentCommit optional.Of[string]
	WantCommit    optional.Of[string]
}

type ApplicationRepository interface {
	GetApplications(ctx context.Context, cond GetApplicationCondition) ([]*domain.Application, error)
	GetApplication(ctx context.Context, id string) (*domain.Application, error)
	CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error)
	UpdateApplication(ctx context.Context, id string, args UpdateApplicationArgs) error
	RegisterApplicationOwner(ctx context.Context, applicationID string, userID string) error
	GetWebsites(ctx context.Context, applicationIDs []string) ([]*domain.Website, error)
	AddWebsite(ctx context.Context, applicationID string, fqdn string, httpPort int) error
	DeleteWebsite(ctx context.Context, applicationID string, websiteID string) error
}

type applicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository(db *sql.DB) ApplicationRepository {
	return &applicationRepository{
		db: db,
	}
}

func (r *applicationRepository) GetApplications(ctx context.Context, cond GetApplicationCondition) ([]*domain.Application, error) {
	mods := []qm.QueryMod{
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

func (r *applicationRepository) CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error) {
	app := &models.Application{
		ID:            domain.NewID(),
		RepositoryID:  args.RepositoryID,
		BranchName:    args.BranchName,
		BuildType:     args.BuildType.String(),
		State:         domain.ApplicationStateIdle.String(),
		CurrentCommit: domain.EmptyCommit,
		WantCommit:    domain.EmptyCommit,
	}
	if err := app.Insert(ctx, r.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	log.WithField("appID", app.ID).
		Info("app created")

	if err := app.L.LoadRepository(ctx, r.db, true, app, nil); err != nil {
		return nil, fmt.Errorf("failed to load repository: %w", err)
	}
	if err := app.L.LoadWebsites(ctx, r.db, true, app, nil); err != nil {
		return nil, fmt.Errorf("failed to load website: %w", err)
	}

	return toDomainApplication(app), nil
}

func (r *applicationRepository) UpdateApplication(ctx context.Context, id string, args UpdateApplicationArgs) error {
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

func (r *applicationRepository) AddWebsite(ctx context.Context, applicationID string, fqdn string, httpPort int) error {
	app, err := r.getApplication(ctx, applicationID)
	if err != nil {
		return err
	}
	website := &models.Website{
		ID:       domain.NewID(),
		FQDN:     fqdn,
		HTTPPort: httpPort,
	}
	err = app.AddWebsites(ctx, r.db, true, website)
	if err != nil {
		return fmt.Errorf("failed to add website: %w", err)
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
