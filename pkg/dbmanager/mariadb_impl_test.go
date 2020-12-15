package dbmanager

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/neoshowcase/pkg/cliutil"
	"testing"
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
	m, db := initMariaDBManager(t)

	t.Skip("TODO") // TODO
	_ = m
	_ = db
}

func TestMariaDBManagerImpl_Delete(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, db := initMariaDBManager(t)

	t.Skip("TODO") // TODO
	_ = m
	_ = db
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
