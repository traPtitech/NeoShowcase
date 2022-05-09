package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ArtifactRepository interface {
	CreateArtifact(ctx context.Context, filename string, buildID string, sid string) error
}

type artifactRepository struct {
	db *sql.DB
}

func NewArtifactRepository(db *sql.DB) ArtifactRepository {
	return &artifactRepository{
		db: db,
	}
}

func (r *artifactRepository) CreateArtifact(ctx context.Context, filename string, buildID string, sid string) error {
	stat, _ := os.Stat(filename)
	artifact := models.Artifact{
		ID:         sid,
		BuildLogID: buildID,
		Size:       stat.Size(),
		CreatedAt:  time.Now(),
	}

	if err := artifact.Insert(context.Background(), r.db, boil.Infer()); err != nil {
		return fmt.Errorf("failed to insert artifact entry: %w", err)
	}

	return nil
}
