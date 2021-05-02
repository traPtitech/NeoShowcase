package userdb

import (
	"context"
	"fmt"
	"time"

	"github.com/traPtitech/neoshowcase/pkg/interface/userdb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoManagerImpl struct {
	client *mongo.Client
}

type MongoConfig struct {
	Host          string
	Port          int
	AdminUser     string
	AdminPassword string
}

func NewMongoManager(config MongoConfig) (userdb.MongoManager, error) {
	// DB接続
	client, err := mongo.NewClient(
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d", config.AdminUser, config.AdminPassword, config.Host, config.Port)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &mongoManagerImpl{client: client}, nil
}

func (m *mongoManagerImpl) Create(ctx context.Context, args userdb.CreateArgs) error {
	client := m.client
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	r := client.Database(args.Database).RunCommand(ctx, bson.D{{Key: "createUser", Value: args.Database}, {Key: "pwd", Value: args.Password}, {Key: "roles", Value: []bson.M{{"role": "dbOwner", "db": args.Database}}}})
	if r.Err() != nil {
		return r.Err()
	}
	return nil
}

func (m *mongoManagerImpl) Delete(ctx context.Context, args userdb.DeleteArgs) error {
	client := m.client
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	r := client.Database(args.Database).RunCommand(ctx, bson.D{{Key: "dropUser", Value: args.Database}})
	if r.Err() != nil {
		return r.Err()
	}
	err := client.Database(args.Database).Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *mongoManagerImpl) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
