package dbmanager

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/cli"
)

func TestNewMariaDBManager(t *testing.T) {
	skipOrDo(t)
	t.Parallel()

	c := MariaDBConfig{
		Host:          cli.GetEnvOrDefault("TEST_APP_MARIADB_HOST", "localhost"),
		Port:          cli.GetIntEnvOrDefault("TEST_APP_MARIADB_PORT", 5004),
		AdminUser:     cli.GetEnvOrDefault("TEST_APP_MARIADB_USER", "root"),
		AdminPassword: cli.GetEnvOrDefault("TEST_APP_MARIADB_PASSWORD", "password"),
	}

	_, err := NewMariaDBManager(c)
	assert.NoError(t, err)
}

func TestMariaDBManagerImpl_CreateDeleteExist(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, _ := initMariaDBManager(t)

	ctx := context.Background()

	dbName := "testCreateDeleteExist"

	dbExists, err := m.IsExist(ctx, dbName)
	assert.NoError(t, err)
	assert.Equal(t, false, dbExists)

	cArgs := domain.CreateArgs{
		Database: dbName,
		Password: dbName,
	}
	err = m.Create(ctx, cArgs)
	assert.NoError(t, err)

	dbExists, err = m.IsExist(ctx, dbName)
	assert.NoError(t, err)
	assert.Equal(t, true, dbExists)

	dArgs := domain.DeleteArgs{
		Database: dbName,
	}
	err = m.Delete(ctx, dArgs)
	assert.NoError(t, err)

	dbExists, err = m.IsExist(ctx, dbName)
	assert.NoError(t, err)
	assert.Equal(t, false, dbExists)
}

func initMariaDBManager(t *testing.T) (*mariaDBManagerImpl, *sql.DB) {
	t.Helper()

	c := MariaDBConfig{
		Host:          cli.GetEnvOrDefault("TEST_APP_MARIADB_HOST", "localhost"),
		Port:          cli.GetIntEnvOrDefault("TEST_APP_MARIADB_PORT", 5004),
		AdminUser:     cli.GetEnvOrDefault("TEST_APP_MARIADB_USER", "root"),
		AdminPassword: cli.GetEnvOrDefault("TEST_APP_MARIADB_PASSWORD", "password"),
	}

	m, err := NewMariaDBManager(c)
	assert.NoError(t, err)

	impl := m.(*mariaDBManagerImpl)
	return impl, impl.db
}
