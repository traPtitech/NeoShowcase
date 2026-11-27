// Command datamigrate migrates data that the declarative schema migration cannot.
// It runs after the schema is applied, and every migration in it is safe to run again.
// It reads the same DB_* environment variables as entrypoint.sh.
package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
)

func main() {
	if err := run(context.Background()); err != nil {
		slog.Error("data migration failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	port, err := strconv.Atoi(getEnv("DB_PORT", "3306"))
	if err != nil {
		return err
	}
	db, err := repository.New(repository.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		Username: getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASS", "password"),
		Database: getEnv("DB_NAME", "neoshowcase"),
	})
	if err != nil {
		return err
	}
	defer db.Close()

	slog.InfoContext(ctx, "migrating application groups")
	return repository.MigrateApplicationGroups(ctx, db)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
