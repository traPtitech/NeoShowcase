package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
type BuildRepository interface {
	GetBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error)
	GetBuildsInCommit(ctx context.Context, commits []string) ([]*domain.Build, error)
	GetBuild(ctx context.Context, buildID string) (*domain.Build, error)
	CreateBuild(ctx context.Context, applicationID string, commit string) (*domain.Build, error)
	UpdateBuild(ctx context.Context, args UpdateBuildArgs) error
	MarkCommitAsRetriable(ctx context.Context, commit string) error
}

type buildRepository struct {
	db *sql.DB
}

func NewBuildRepository(db *sql.DB) BuildRepository {
	return &buildRepository{
		db: db,
	}
}

func (r *buildRepository) getBuild(ctx context.Context, id string) (*models.Build, error) {
	return models.Builds(
		models.BuildWhere.ID.EQ(id),
		qm.Load(models.BuildRels.Artifact),
	).One(ctx, r.db)
}

func (r *buildRepository) GetBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error) {
	builds, err := models.Builds(
		models.BuildWhere.ApplicationID.EQ(applicationID),
		qm.Load(models.BuildRels.Artifact),
	).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get builds: %w", err)
	}
	return lo.Map(builds, func(b *models.Build, i int) *domain.Build {
		return toDomainBuild(b)
	}), nil
}

func (r *buildRepository) GetBuildsInCommit(ctx context.Context, commits []string) ([]*domain.Build, error) {
	builds, err := models.Builds(
		models.BuildWhere.Commit.IN(commits),
		qm.Load(models.BuildRels.Artifact),
	).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get builds: %w", err)
	}
	return lo.Map(builds, func(b *models.Build, i int) *domain.Build {
		return toDomainBuild(b)
	}), nil
}

func (r *buildRepository) GetBuild(ctx context.Context, buildID string) (*domain.Build, error) {
	build, err := r.getBuild(ctx, buildID)
	if err != nil {
		if isNoRowsErr(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find build: %w", err)
	}
	return toDomainBuild(build), nil
}

func (r *buildRepository) CreateBuild(ctx context.Context, applicationID string, commit string) (*domain.Build, error) {
	const errMsg = "failed to CreateBuildLog: %w"

	build := &models.Build{
		ID:            domain.NewID(),
		Commit:        commit,
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

	build, err := r.getBuild(ctx, args.ID)
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

func (r *buildRepository) MarkCommitAsRetriable(ctx context.Context, commit string) error {
	_, err := models.Builds(
		models.BuildWhere.Commit.EQ(commit),
	).UpdateAll(ctx, r.db, models.M{
		models.BuildColumns.Retriable: true,
	})
	if err != nil {
		return fmt.Errorf("failed to mark commit as retriable: %w", err)
	}
	return nil
}
