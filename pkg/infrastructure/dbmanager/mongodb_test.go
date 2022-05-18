package dbmanager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/cli"
	"go.mongodb.org/mongo-driver/mongo"
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

func TestMongoDBManagerImpl_Create(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, _ := initMongoDBManager(t)

	a := domain.CreateArgs{
		Database: "testCreate",
		Password: "testCreate",
	}
	ctx := context.Background()
	err := m.Create(ctx, a)
	assert.NoError(t, err)
}

func TestMongoDBManagerImpl_Delete(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, _ := initMongoDBManager(t)

	a := domain.CreateArgs{
		Database: "testDelete",
		Password: "testDelete",
	}
	ctx := context.Background()
	_ = m.Create(ctx, a)

	da := domain.DeleteArgs{
		Database: "testDelete",
	}
	err := m.Delete(ctx, da)
	assert.NoError(t, err)
}

func TestMongoDBManagerImpl_Exist(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, _ := initMongoDBManager(t)

	a := domain.CreateArgs{
		Database: "testExist",
		Password: "testExist",
	}

	dbExists, _ := m.IsExist(context.Background(), a.Database)
	assert.Equal(t, false, dbExists)

	ctx := context.Background()
	_ = m.Create(ctx, a)

	dbExists, _ = m.IsExist(ctx, a.Database)
	assert.Equal(t, true, dbExists)
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
