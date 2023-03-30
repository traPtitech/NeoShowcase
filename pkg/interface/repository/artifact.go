package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

type artifactRepository struct {
	db *sql.DB
}

func NewArtifactRepository(db *sql.DB) domain.ArtifactRepository {
	return &artifactRepository{
		db: db,
	}
}

func (r *artifactRepository) CreateArtifact(ctx context.Context, size int64, buildID string, sid string) error {
	artifact := models.Artifact{
		ID:        sid,
		BuildID:   buildID,
		Size:      size,
		CreatedAt: time.Now(),
	}

	if err := artifact.Insert(context.Background(), r.db, boil.Infer()); err != nil {
		return errors.Wrap(err, "failed to insert artifact entry")
	}

	return nil
}
