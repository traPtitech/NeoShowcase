package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
type BuildRepository interface {
	GetBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error)
	GetBuild(ctx context.Context, buildID string) (*domain.Build, error)
	CreateBuild(ctx context.Context, applicationID string) (*domain.Build, error)
	UpdateBuild(ctx context.Context, args UpdateBuildArgs) error
}

type buildRepository struct {
	db *sql.DB
}

func NewBuildRepository(db *sql.DB) BuildRepository {
	return &buildRepository{
		db: db,
	}
}

func (r *buildRepository) GetBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error) {
	builds, err := models.Builds(
		models.BuildWhere.ApplicationID.EQ(applicationID),
	).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get builds: %w", err)
	}
	return lo.Map(builds, func(b *models.Build, i int) *domain.Build {
		return toDomainBuild(b)
	}), nil
}

func (r *buildRepository) GetBuild(ctx context.Context, buildID string) (*domain.Build, error) {
	build, err := models.FindBuild(ctx, r.db, buildID)
	if err != nil {
		if isNoRowsErr(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find build: %w", err)
	}
	return toDomainBuild(build), nil
}

func (r *buildRepository) CreateBuild(ctx context.Context, applicationID string) (*domain.Build, error) {
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

func (r *buildRepository) UpdateBuild(ctx context.Context, args UpdateBuildArgs) error {
	const errMsg = "failed to UpdateBuildStatus: %w"

	build, err := models.FindBuild(ctx, r.db, args.ID)

	if err != nil {
		if isNoRowsErr(err) {
			return ErrNotFound
		}
		return fmt.Errorf(errMsg, err)
	}

	build.Status = args.Status.String()
	if args.Status.IsFinished() {
		build.FinishedAt = null.TimeFrom(time.Now())
	}

	if _, err := build.Update(ctx, r.db, boil.Infer()); err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}
