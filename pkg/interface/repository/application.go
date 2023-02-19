package repository

import (
	"context"
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
type ApplicationRepository interface {
	CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error)
	GetApplicationByID(ctx context.Context, id string) (*domain.Application, error)
	SetWebsite(ctx context.Context, branchID string, fqdn string, httpPort int) error
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
	OwnerID       string // TODO: UserID型にする
	RepositoryURL string
	BranchName    string
	BuildType     builder.BuildType
}

func (r *applicationRepository) CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error) {
	const errMsg = "failed to CreateApplication: %w"

	// リポジトリ情報を設定
	repo, err := models.Repositories(models.RepositoryWhere.URL.EQ(args.RepositoryURL)).One(ctx, r.db)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf(errMsg, err)
	} else if repo == nil {
		repo = &models.Repository{
			ID:  domain.NewID(),
			URL: args.RepositoryURL,
		}
		if err := repo.Insert(ctx, r.db, boil.Infer()); err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
	}

	// アプリケーション作成
	app := &models.Application{
		ID:           domain.NewID(),
		RepositoryID: repo.ID,
		BranchName:   args.BranchName,
		BuildType:    args.BuildType.String(),
	}
	if err := app.Insert(ctx, r.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	// 初期Ownerを設定
	user, err := models.Users(models.UserWhere.ID.EQ(args.OwnerID)).One(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	if err := app.AddUsers(ctx, r.db, false, user); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	log.WithField("appID", app.ID).
		Info("app created")

	return toDomainApplication(app, toDomainRepository(repo)), nil
}

func (r *applicationRepository) GetApplicationByID(ctx context.Context, id string) (*domain.Application, error) {
	const errMsg = "failed to GetApplicationByID: %w"

	app, err := models.Applications(
		qm.Load(models.ApplicationRels.Repository),
		models.ApplicationWhere.ID.EQ(id),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(errMsg, err)
	}

	return toDomainApplication(app, toDomainRepository(app.R.Repository)), nil
}

func (r *applicationRepository) SetWebsite(ctx context.Context, branchID string, fqdn string, httpPort int) error {
	const errMsg = "failed to SetWebsite: %w"

	branch, err := models.Applications(
		qm.Load(models.ApplicationRels.Website),
		models.ApplicationWhere.ID.EQ(branchID),
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
