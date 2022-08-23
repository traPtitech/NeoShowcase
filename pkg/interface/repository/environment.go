package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
type EnvironmentRepository interface {
	SetEnv(ctx context.Context, branchID, key, value string) error
}

type environmentRepository struct {
	db *sql.DB
}

func NewEnvironmentRepository(db *sql.DB) EnvironmentRepository {
	return &environmentRepository{db: db}
}

func (r *environmentRepository) SetEnv(ctx context.Context, branchID, key, value string) error {
	const errMsg = "failed to SetEnv: %w"

	env := models.Environment{
		ID:       domain.NewID(),
		BranchID: branchID,
		Key:      key,
		Value:    value,
	}

	if err := env.Upsert(ctx, r.db, boil.Infer(), boil.Infer()); err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}
