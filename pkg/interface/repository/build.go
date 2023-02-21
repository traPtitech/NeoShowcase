package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
type BuildRepository interface {
	CreateBuild(ctx context.Context, applicationID string) (*domain.Build, error)
	UpdateBuild(ctx context.Context, args UpdateBuildArgs) error
}

type buildLogRepository struct {
	db *sql.DB
}

func NewBuildLogRepository(db *sql.DB) BuildRepository {
	return &buildLogRepository{
		db: db,
	}
}

func (r *buildLogRepository) CreateBuild(ctx context.Context, applicationID string) (*domain.Build, error) {
	const errMsg = "failed to CreateBuildLog: %w"

	build := &models.Build{
		ID:            domain.NewID(),
		Status:        builder.BuildStatusQueued.String(),
		StartedAt:     time.Now(),
		ApplicationID: applicationID,
	}

	if err := build.Insert(ctx, r.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return toDomainBuild(build), nil
}

type UpdateBuildArgs struct {
	ID     string
	Status builder.BuildStatus
}

func (r *buildLogRepository) UpdateBuild(ctx context.Context, args UpdateBuildArgs) error {
	const errMsg = "failed to UpdateBuildStatus: %w"

	buildLog, err := models.FindBuild(ctx, r.db, args.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf(errMsg, err)
	}

	buildLog.Status = args.Status.String()
	if args.Status.IsFinished() {
		buildLog.FinishedAt = null.TimeFrom(time.Now())
	}

	if _, err := buildLog.Update(ctx, r.db, boil.Infer()); err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}
