package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
type BuildLogRepository interface {
	CreateBuildLog(ctx context.Context, branchID string) (*domain.BuildLog, error)
	UpdateBuildLog(ctx context.Context, args UpdateBuildLogArgs) error
}

type buildLogRepository struct {
	db *sql.DB
}

func NewBuildLogRepository(db *sql.DB) BuildLogRepository {
	return &buildLogRepository{
		db: db,
	}
}

func (r *buildLogRepository) CreateBuildLog(ctx context.Context, branchID string) (*domain.BuildLog, error) {
	const errMsg = "failed to CreateBuildLog: %w"

	buildLog := &models.BuildLog{
		ID:        domain.NewID(),
		Result:    builder.BuildStatusQueued.String(),
		StartedAt: time.Now(),
		BranchID:  branchID,
	}

	if err := buildLog.Insert(ctx, r.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return &domain.BuildLog{
		ID:       buildLog.ID,
		Result:   builder.BuildStatusQueued,
		BranchID: buildLog.BranchID,
	}, nil
}

type UpdateBuildLogArgs struct {
	ID       string
	Result   builder.BuildStatus
	Finished bool
}

func (r *buildLogRepository) UpdateBuildLog(ctx context.Context, args UpdateBuildLogArgs) error {
	const errMsg = "failed to UpdateBuildStatus: %w"

	buildLog, err := models.FindBuildLog(ctx, r.db, args.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf(errMsg, err)
	}

	buildLog.Result = args.Result.String()
	if args.Finished {
		buildLog.FinishedAt = null.TimeFrom(time.Now())
	}

	if _, err := buildLog.Update(ctx, r.db, boil.Infer()); err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}
