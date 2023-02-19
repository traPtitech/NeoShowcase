package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

type ArtifactRepository interface {
	CreateArtifact(ctx context.Context, size int64, buildID string, sid string) error
}

type artifactRepository struct {
	db *sql.DB
}

func NewArtifactRepository(db *sql.DB) ArtifactRepository {
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
		return fmt.Errorf("failed to insert artifact entry: %w", err)
	}

	return nil
}
