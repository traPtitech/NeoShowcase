package userdb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/neoshowcase/pkg/util/cli"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewMongoManager(t *testing.T) {
	skipOrDo(t)
	t.Parallel()

	c := MongoConfig{
		Host:          cli.GetEnvOrDefault("TEST_APP_MONGO_HOST", "localhost"),
		Port:          cli.GetIntEnvOrDefault("TEST_APP_MONGO_PORT", 5010),
		AdminUser:     cli.GetEnvOrDefault("TEST_APP_MONGO_USER", "root"),
		AdminPassword: cli.GetEnvOrDefault("TEST_APP_MONGO_PASSWORD", "password"),
	}

	_, err := NewMongoManager(c)
	assert.NoError(t, err)
}

func TestMongoManagerImpl_Create(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, _ := initMongoManager(t)

	a := CreateArgs{
		Database: "testCreate",
		Password: "testCreate",
	}
	ctx := context.Background()
	err := m.Create(ctx, a)
	assert.NoError(t, err)
}

func TestMongoManagerImpl_Delete(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, _ := initMongoManager(t)

	a := CreateArgs{
		Database: "testDelete",
		Password: "testDelete",
	}
	ctx := context.Background()
	_ = m.Create(ctx, a)

	da := DeleteArgs{
		Database: "testDelete",
	}
	err := m.Delete(ctx, da)
	assert.NoError(t, err)
}

func initMongoManager(t *testing.T) (*mongoManagerImpl, *mongo.Client) {
	t.Helper()

	c := MongoConfig{
		Host:          cli.GetEnvOrDefault("TEST_APP_MONGO_HOST", "localhost"),
		Port:          cli.GetIntEnvOrDefault("TEST_APP_MONGO_PORT", 5010),
		AdminUser:     cli.GetEnvOrDefault("TEST_APP_MONGO_USER", "root"),
		AdminPassword: cli.GetEnvOrDefault("TEST_APP_MONGO_PASSWORD", "password"),
	}

	m, err := NewMongoManager(c)
	assert.NoError(t, err)

	impl := m.(*mongoManagerImpl)
	return impl, impl.client
}
