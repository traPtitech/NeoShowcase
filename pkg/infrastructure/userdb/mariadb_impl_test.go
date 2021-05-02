package userdb

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/neoshowcase/pkg/cliutil"
)

func TestNewMariaDBManager(t *testing.T) {
	skipOrDo(t)
	t.Parallel()

	c := MariaDBConfig{
		Host:          cliutil.GetEnvOrDefault("TEST_APP_MARIADB_HOST", "localhost"),
		Port:          cliutil.GetIntEnvOrDefault("TEST_APP_MARIADB_PORT", 5004),
		AdminUser:     cliutil.GetEnvOrDefault("TEST_APP_MARIADB_USER", "root"),
		AdminPassword: cliutil.GetEnvOrDefault("TEST_APP_MARIADB_PASSWORD", "password"),
	}

	_, err := NewMariaDBManager(c)
	assert.NoError(t, err)
}

func TestMariaDBManagerImpl_Create(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, _ := initMariaDBManager(t)
	a := CreateArgs{
		Database: "TestMariaCreate",
		Password: "TestMariaCreate",
	}
	ctx := context.Background()
	err := m.Create(ctx, a)
	assert.NoError(t, err)
}

func TestMariaDBManagerImpl_Delete(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, _ := initMariaDBManager(t)

	a := CreateArgs{
		Database: "TestMariaDelete",
		Password: "TestMariaDelete",
	}
	ctx := context.Background()
	err := m.Create(ctx, a)
	assert.NoError(t, err)

	da := DeleteArgs{
		Database: "TestMariaDelete",
	}
	err = m.Delete(ctx, da)
	assert.NoError(t, err)
}

func initMariaDBManager(t *testing.T) (*mariaDBManagerImpl, *sql.DB) {
	t.Helper()

	c := MariaDBConfig{
		Host:          cliutil.GetEnvOrDefault("TEST_APP_MARIADB_HOST", "localhost"),
		Port:          cliutil.GetIntEnvOrDefault("TEST_APP_MARIADB_PORT", 5004),
		AdminUser:     cliutil.GetEnvOrDefault("TEST_APP_MARIADB_USER", "root"),
		AdminPassword: cliutil.GetEnvOrDefault("TEST_APP_MARIADB_PASSWORD", "password"),
	}

	m, err := NewMariaDBManager(c)
	assert.NoError(t, err)

	impl := m.(*mariaDBManagerImpl)
	return impl, impl.db
}
