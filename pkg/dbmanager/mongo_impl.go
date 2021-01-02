package dbmanager

import (
	"context"
	"fmt"
	"time"

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

func NewMongoManager(config MongoConfig) (MongoManager, error) {
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

func (m *mongoManagerImpl) Create(ctx context.Context, args CreateArgs) error {
	client := m.client
	_, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	err := client.Connect(ctx)
	if err != nil {
		return err
	}
	r := client.Database(args.Database).RunCommand(ctx, bson.D{{"createUser", args.Database}, {"pwd", args.Password}, {"roles", []bson.M{{"role": "dbAdminAnyDatabase", "db": args.Database}}}})
	if r.Err() != nil {
		return r.Err()
	}
	defer client.Disconnect(ctx)
	return nil
}

func (m *mongoManagerImpl) Delete(ctx context.Context, args DeleteArgs) error {
	client := m.client
	_, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	err := client.Connect(ctx)
	if err != nil {
		return err
	}
	r := client.Database(args.Database).RunCommand(ctx, bson.D{{"removeUser", args.Database}})
	if r.Err() != nil {
		return r.Err()
	}
	r = client.Database(args.Database).RunCommand(ctx, "dropDatabase")
	if r.Err() != nil {
		return r.Err()
	}
	defer client.Disconnect(ctx)
	return nil
}

func (m *mongoManagerImpl) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
