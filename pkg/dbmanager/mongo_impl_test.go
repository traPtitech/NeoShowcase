package dbmanager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/neoshowcase/pkg/cliutil"
	"go.mongodb.org/mongo-driver/mongo"
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
	m, _ := initMongoManager(t)

	a := CreateArgs{
		Database: "test",
		Password: "test",
	}
	ctx := context.Background()
	m.Create(ctx, a)
}

func TestMongoManagerImpl_Delete(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, _ := initMongoManager(t)

	a := CreateArgs{
		Database: "test",
		Password: "test",
	}
	ctx := context.Background()
	m.Create(ctx, a)

	da := DeleteArgs{
		Database: "test",
	}
	m.Delete(ctx, da)
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
