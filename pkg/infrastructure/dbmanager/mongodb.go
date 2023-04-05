package dbmanager

import (
	"context"
	"fmt"
	"time"

	"github.com/friendsofgo/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type mongoDBManagerImpl struct {
	client *mongo.Client
	c      MongoDBConfig
}

type MongoDBConfig struct {
	Host          string `mapstructure:"host" yaml:"host"`
	Port          int    `mapstructure:"port" yaml:"port"`
	AdminUser     string `mapstructure:"adminUser" yaml:"adminUser"`
	AdminPassword string `mapstructure:"adminPassword" yaml:"adminPassword"`
}

func NewMongoDBManager(config MongoDBConfig) (domain.MongoDBManager, error) {
	// DB接続
	client, err := mongo.NewClient(
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d", config.AdminUser, config.AdminPassword, config.Host, config.Port)),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new client")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}

	return &mongoDBManagerImpl{client: client, c: config}, nil
}

func (m *mongoDBManagerImpl) GetHost() (host string, port int) {
	return m.c.Host, m.c.Port
}

func (m *mongoDBManagerImpl) Create(ctx context.Context, args domain.CreateArgs) error {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	cmd := bson.D{
		{Key: "createUser", Value: args.Database},
		{Key: "pwd", Value: args.Password},
		{Key: "roles", Value: []bson.M{{"role": "dbOwner", "db": args.Database}}},
	}

	// NOTE: the database is created only after first write operation
	if r := m.client.Database(args.Database).RunCommand(ctx, cmd); r.Err() != nil {
		return r.Err()
	}

	return nil
}

func (m *mongoDBManagerImpl) Delete(ctx context.Context, args domain.DeleteArgs) error {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	cmd := bson.D{{Key: "dropUser", Value: args.Database}}

	if r := m.client.Database(args.Database).RunCommand(ctx, cmd); r.Err() != nil {
		return r.Err()
	}
	if err := m.client.Database(args.Database).Drop(ctx); err != nil {
		return err
	}
	return nil
}

func (m *mongoDBManagerImpl) IsExist(ctx context.Context, name string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	cur, err := m.client.Database("admin").Collection("system.users").Find(ctx, bson.D{{"db", name}})
	if err != nil {
		return false, err
	}
	for cur.TryNext(ctx) {
		return true, nil
	}
	return false, nil
}

func (m *mongoDBManagerImpl) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
