package dbmanager

import (
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/neoshowcase/pkg/cliutil"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestNewMongoManager(t *testing.T) {
	skipOrDo(t)
	t.Parallel()

	c := MongoConfig{
		Host:          cliutil.GetEnvOrDefault("TEST_APP_MONGO_HOST", "localhost"),
		Port:          cliutil.GetIntEnvOrDefault("TEST_APP_MONGO_PORT", 5010),
		AdminUser:     cliutil.GetEnvOrDefault("TEST_APP_MONGO_USER", "root"),
		AdminPassword: cliutil.GetEnvOrDefault("TEST_APP_MONGO_PASSWORD", "password"),
	}

	_, err := NewMongoManager(c)
	assert.NoError(t, err)
}

func TestMongoManagerImpl_Create(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, client := initMongoManager(t)

	t.Skip("TODO") // TODO
	_ = m
	_ = client
}

func TestMongoManagerImpl_Delete(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, client := initMongoManager(t)

	t.Skip("TODO") // TODO
	_ = m
	_ = client
}

func initMongoManager(t *testing.T) (*mongoManagerImpl, *mongo.Client) {
	t.Helper()

	c := MongoConfig{
		Host:          cliutil.GetEnvOrDefault("TEST_APP_MONGO_HOST", "localhost"),
		Port:          cliutil.GetIntEnvOrDefault("TEST_APP_MONGO_PORT", 5010),
		AdminUser:     cliutil.GetEnvOrDefault("TEST_APP_MONGO_USER", "root"),
		AdminPassword: cliutil.GetEnvOrDefault("TEST_APP_MONGO_PASSWORD", "password"),
	}

	m, err := NewMongoManager(c)
	assert.NoError(t, err)

	impl := m.(*mongoManagerImpl)
	return impl, impl.client
}
