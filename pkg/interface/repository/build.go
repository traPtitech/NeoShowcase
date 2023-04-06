package repository

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type buildRepository struct {
	db *sql.DB
}

func NewBuildRepository(db *sql.DB) domain.BuildRepository {
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

func (r *buildRepository) buildMods(cond domain.GetBuildCondition) []qm.QueryMod {
	var mods []qm.QueryMod
	if cond.ApplicationID.Valid {
		mods = append(mods, models.BuildWhere.ApplicationID.EQ(cond.ApplicationID.V))
	}
	if cond.Commit.Valid {
		mods = append(mods, models.BuildWhere.Commit.EQ(cond.Commit.V))
	}
	if cond.CommitIn.Valid {
		mods = append(mods, models.BuildWhere.Commit.IN(cond.CommitIn.V))
	}
	if cond.Status.Valid {
		mods = append(mods, models.BuildWhere.Status.EQ(buildStatusMapper.FromMust(cond.Status.V)))
	}
	if cond.Retriable.Valid {
		mods = append(mods, models.BuildWhere.Retriable.EQ(cond.Retriable.V))
	}
	return mods
}

func (r *buildRepository) GetBuilds(ctx context.Context, cond domain.GetBuildCondition) ([]*domain.Build, error) {
	mods := []qm.QueryMod{qm.Load(models.BuildRels.Artifact)}
	mods = append(mods, r.buildMods(cond)...)
	builds, err := models.Builds(mods...).All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get builds")
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
		return nil, errors.Wrap(err, "failed to find build")
	}
	return toDomainBuild(build), nil
}

func (r *buildRepository) CreateBuild(ctx context.Context, build *domain.Build) error {
	mb := fromDomainBuild(build)
	err := mb.Insert(ctx, r.db, boil.Blacklist())
	if err != nil {
		return errors.Wrap(err, "failed to insert build")
	}
	return nil
}

func (r *buildRepository) UpdateBuild(ctx context.Context, id string, args domain.UpdateBuildArgs) error {
	mods := []qm.QueryMod{
		models.BuildWhere.ID.EQ(id),
		qm.For("UPDATE"),
	}

	if args.FromStatus.Valid {
		mods = append(mods, models.BuildWhere.Status.EQ(buildStatusMapper.FromMust(args.FromStatus.V)))
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()

	build, err := models.Builds(mods...).One(ctx, tx)
	if err != nil {
		if isNoRowsErr(err) {
			return ErrNotFound
		}
		return errors.Wrap(err, "failed to get build")
	}

	var cols []string
	if args.Status.Valid {
		build.Status = buildStatusMapper.FromMust(args.Status.V)
		cols = append(cols, models.BuildColumns.Status)
	}
	if args.StartedAt.Valid {
		build.StartedAt = optional.IntoTime(args.StartedAt)
		cols = append(cols, models.BuildColumns.StartedAt)
	}
	if args.UpdatedAt.Valid {
		build.UpdatedAt = optional.IntoTime(args.UpdatedAt)
		cols = append(cols, models.BuildColumns.UpdatedAt)
	}
	if args.FinishedAt.Valid {
		build.FinishedAt = optional.IntoTime(args.FinishedAt)
		cols = append(cols, models.BuildColumns.FinishedAt)
	}

	_, err = build.Update(ctx, tx, boil.Whitelist(cols...))
	if err != nil {
		return errors.Wrap(err, "failed to update build")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit")
	}

	return nil
}

func (r *buildRepository) MarkCommitAsRetriable(ctx context.Context, applicationID string, commit string) error {
	_, err := models.Builds(
		models.BuildWhere.ApplicationID.EQ(applicationID),
		models.BuildWhere.Commit.EQ(commit),
	).UpdateAll(ctx, r.db, models.M{
		models.BuildColumns.Retriable: true,
	})
	if err != nil {
		return errors.Wrap(err, "failed to mark commit as retriable")
	}
	return nil
}

func (r *buildRepository) DeleteBuilds(ctx context.Context, cond domain.GetBuildCondition) error {
	builds, err := models.Builds(r.buildMods(cond)...).All(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to get builds")
	}
	_, err = builds.DeleteAll(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to delete builds")
	}
	return nil
}
