package dbmanager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/cli"
)

func TestNewMongoDBManager(t *testing.T) {
	skipOrDo(t)
	t.Parallel()

	c := testMongoDBConfig()

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

// A deletion that dropped the user but failed to drop the database must be retryable.
func TestMongoDBManagerImpl_DeleteAfterUserDropped(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, client := initMongoDBManager(t)

	ctx := context.Background()

	dbName := "testDeleteAfterUserDropped"
	dropMongoDBOnCleanup(t, client, dbName)

	err := m.Create(ctx, domain.CreateArgs{Database: dbName, Password: dbName})
	require.NoError(t, err)
	_, err = client.Database(dbName).Collection("c").InsertOne(ctx, bson.D{{Key: "k", Value: "v"}})
	require.NoError(t, err)
	err = client.Database(dbName).RunCommand(ctx, bson.D{{Key: "dropUser", Value: dbName}}).Err()
	require.NoError(t, err)

	err = m.Delete(ctx, domain.DeleteArgs{Database: dbName})
	assert.NoError(t, err)

	names, err := client.ListDatabaseNames(ctx, bson.D{{Key: "name", Value: dbName}})
	require.NoError(t, err)
	assert.Empty(t, names)
}

func TestMongoDBManagerImpl_DeleteTwice(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, client := initMongoDBManager(t)

	ctx := context.Background()

	dbName := "testDeleteTwice"
	dropMongoDBOnCleanup(t, client, dbName)

	err := m.Create(ctx, domain.CreateArgs{Database: dbName, Password: dbName})
	require.NoError(t, err)

	err = m.Delete(ctx, domain.DeleteArgs{Database: dbName})
	assert.NoError(t, err)
	err = m.Delete(ctx, domain.DeleteArgs{Database: dbName})
	assert.NoError(t, err)
}

// A deletion that fails to drop the database must leave the user, so that IsExist reports the database as not yet deleted.
func TestMongoDBManagerImpl_DeleteKeepsUserWhenDropFails(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, client := initMongoDBManager(t)

	ctx := context.Background()

	dbName := "testDeleteKeepsUserWhenDropFails"
	dropMongoDBOnCleanup(t, client, dbName)

	err := m.Create(ctx, domain.CreateArgs{Database: dbName, Password: dbName})
	require.NoError(t, err)
	_, err = client.Database(dbName).Collection("c").InsertOne(ctx, bson.D{{Key: "k", Value: "v"}})
	require.NoError(t, err)

	restricted := newRestrictedMongoDBManager(t, client, dbName, "dropUser")
	err = restricted.Delete(ctx, domain.DeleteArgs{Database: dbName})
	assert.Error(t, err)

	dbExists, err := m.IsExist(ctx, dbName)
	require.NoError(t, err)
	assert.True(t, dbExists)
}

func TestMongoDBManagerImpl_DeleteReturnsDropUserError(t *testing.T) {
	skipOrDo(t)
	t.Parallel()
	m, client := initMongoDBManager(t)

	ctx := context.Background()

	dbName := "testDeleteReturnsDropUserError"
	dropMongoDBOnCleanup(t, client, dbName)

	err := m.Create(ctx, domain.CreateArgs{Database: dbName, Password: dbName})
	require.NoError(t, err)

	restricted := newRestrictedMongoDBManager(t, client, dbName, "dropDatabase")
	err = restricted.Delete(ctx, domain.DeleteArgs{Database: dbName})
	assert.Error(t, err)

	dbExists, err := m.IsExist(ctx, dbName)
	require.NoError(t, err)
	assert.True(t, dbExists)
}

func testMongoDBConfig() MongoDBConfig {
	return MongoDBConfig{
		Host:          cli.GetEnvOrDefault("TEST_APP_MONGODB_HOST", "localhost"),
		Port:          cli.GetIntEnvOrDefault("TEST_APP_MONGODB_PORT", 5010),
		AdminUser:     cli.GetEnvOrDefault("TEST_APP_MONGODB_USER", "root"),
		AdminPassword: cli.GetEnvOrDefault("TEST_APP_MONGODB_PASSWORD", "password"),
	}
}

func initMongoDBManager(t *testing.T) (*mongoDBManagerImpl, *mongo.Client) {
	t.Helper()

	c := testMongoDBConfig()

	m, err := NewMongoDBManager(c)
	assert.NoError(t, err)

	impl := m.(*mongoDBManagerImpl)
	return impl, impl.client
}

// dropMongoDBOnCleanup drops the user and the database named dbName now and again when the test ends,
// so that leftovers of a failed run do not break the next one.
func dropMongoDBOnCleanup(t *testing.T, client *mongo.Client, dbName string) {
	t.Helper()

	drop := func() {
		ctx := context.Background()
		_ = client.Database(dbName).RunCommand(ctx, bson.D{{Key: "dropUser", Value: dbName}}).Err()
		require.NoError(t, client.Database(dbName).Drop(ctx))
	}
	drop()
	t.Cleanup(drop)
}

// newRestrictedMongoDBManager returns a manager whose admin user may only perform the given actions on dbName.
// The admin user and its role are removed when the test ends.
func newRestrictedMongoDBManager(t *testing.T, client *mongo.Client, dbName string, actions ...string) *mongoDBManagerImpl {
	t.Helper()

	ctx := context.Background()
	admin := client.Database("admin")
	name := dbName + "Admin"
	password := "password"

	dropAdmin := func() {
		_ = admin.RunCommand(ctx, bson.D{{Key: "dropUser", Value: name}}).Err()
		_ = admin.RunCommand(ctx, bson.D{{Key: "dropRole", Value: name}}).Err()
	}
	dropAdmin()
	t.Cleanup(dropAdmin)

	err := admin.RunCommand(ctx, bson.D{
		{Key: "createRole", Value: name},
		{Key: "privileges", Value: bson.A{bson.M{
			"resource": bson.M{"db": dbName, "collection": ""},
			"actions":  actions,
		}}},
		{Key: "roles", Value: bson.A{}},
	}).Err()
	require.NoError(t, err)
	err = admin.RunCommand(ctx, bson.D{
		{Key: "createUser", Value: name},
		{Key: "pwd", Value: password},
		{Key: "roles", Value: bson.A{bson.M{"role": name, "db": "admin"}}},
	}).Err()
	require.NoError(t, err)

	c := testMongoDBConfig()
	c.AdminUser = name
	c.AdminPassword = password
	m, err := NewMongoDBManager(c)
	require.NoError(t, err)
	t.Cleanup(func() { _ = m.Close(context.Background()) })
	return m.(*mongoDBManagerImpl)
}
