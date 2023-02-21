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
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
type ApplicationRepository interface {
	GetApplicationsByUserID(ctx context.Context, userID string) ([]*domain.Application, error)
	CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error)
	RegisterApplicationOwner(ctx context.Context, applicationID string, userID string) error
	GetApplicationByID(ctx context.Context, id string) (*domain.Application, error)
	SetWebsite(ctx context.Context, applicationID string, fqdn string, httpPort int) error
}

type applicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository(db *sql.DB) ApplicationRepository {
	return &applicationRepository{
		db: db,
	}
}

type CreateApplicationArgs struct {
	RepositoryID string
	BranchName   string
	BuildType    builder.BuildType
}

func (r *applicationRepository) GetApplicationsByUserID(ctx context.Context, userID string) ([]*domain.Application, error) {
	applications, err := models.Applications(
		qm.Load(models.ApplicationRels.Repository),
		qm.Load(models.ApplicationRels.Users),
		models.UserWhere.ID.EQ(userID),
	).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications: %w", err)
	}

	return lo.Map(applications, func(app *models.Application, i int) *domain.Application {
		return toDomainApplication(app, toDomainRepository(app.R.Repository))
	}), nil
}

func (r *applicationRepository) CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error) {
	app := &models.Application{
		ID:           domain.NewID(),
		RepositoryID: args.RepositoryID,
		BranchName:   args.BranchName,
		BuildType:    args.BuildType.String(),
	}
	if err := app.Insert(ctx, r.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	log.WithField("appID", app.ID).
		Info("app created")

	if err := app.L.LoadRepository(ctx, r.db, true, app, nil); err != nil {
		return nil, fmt.Errorf("failed to load repository: %w", err)
	}

	return toDomainApplication(app, toDomainRepository(app.R.Repository)), nil
}

func (r *applicationRepository) RegisterApplicationOwner(ctx context.Context, applicationID string, userID string) error {
	app, err := models.Applications(models.ApplicationWhere.ID.EQ(applicationID)).One(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to find application: %w", err)
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

func (r *applicationRepository) GetApplicationByID(ctx context.Context, id string) (*domain.Application, error) {
	app, err := models.Applications(
		models.ApplicationWhere.ID.EQ(id),
		qm.Load(models.ApplicationRels.Repository),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	return toDomainApplication(app, toDomainRepository(app.R.Repository)), nil
}

func (r *applicationRepository) SetWebsite(ctx context.Context, applicationID string, fqdn string, httpPort int) error {
	const errMsg = "failed to SetWebsite: %w"

	branch, err := models.Applications(
		qm.Load(models.ApplicationRels.Website),
		models.ApplicationWhere.ID.EQ(applicationID),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf(errMsg, err)
	}

	ws := branch.R.Website
	if ws != nil {
		// テーブルの情報を更新
		ws.FQDN = fqdn
		ws.HTTPPort = httpPort
		if _, err := ws.Update(ctx, r.db, boil.Infer()); err != nil {
			return fmt.Errorf(errMsg, err)
		}
		return nil
	}

	// Websiteをテーブルに挿入
	ws = &models.Website{
		ID:       domain.NewID(),
		FQDN:     fqdn,
		HTTPPort: httpPort,
	}
	if err := branch.SetWebsite(ctx, r.db, true, ws); err != nil {
		return fmt.Errorf(errMsg, err)
	}
	return nil
}
