package dbmanager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/cli"
)

func TestNewMongoDBManager(t *testing.T) {
	skipOrDo(t)
	t.Parallel()

	c := MongoDBConfig{
		Host:          cli.GetEnvOrDefault("TEST_APP_MONGODB_HOST", "localhost"),
		Port:          cli.GetIntEnvOrDefault("TEST_APP_MONGODB_PORT", 5010),
		AdminUser:     cli.GetEnvOrDefault("TEST_APP_MONGODB_USER", "root"),
		AdminPassword: cli.GetEnvOrDefault("TEST_APP_MONGODB_PASSWORD", "password"),
	}

	_, err := NewMongoDBManager(c)
	assert.NoError(t, err)
}

func TestMongoDBManagerImpl_CreateDeleteExist(t *testing.T) {
	//skipOrDo(t)
	t.Parallel()
	m, _ := initMongoDBManager(t)

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

func initMongoDBManager(t *testing.T) (*mongoDBManagerImpl, *mongo.Client) {
	t.Helper()

	c := MongoDBConfig{
		Host:          cli.GetEnvOrDefault("TEST_APP_MONGODB_HOST", "localhost"),
		Port:          cli.GetIntEnvOrDefault("TEST_APP_MONGODB_PORT", 5010),
		AdminUser:     cli.GetEnvOrDefault("TEST_APP_MONGODB_USER", "root"),
		AdminPassword: cli.GetEnvOrDefault("TEST_APP_MONGODB_PASSWORD", "password"),
	}

	m, err := NewMongoDBManager(c)
	assert.NoError(t, err)

	impl := m.(*mongoDBManagerImpl)
	return impl, impl.client
}
