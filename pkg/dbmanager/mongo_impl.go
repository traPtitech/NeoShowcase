package dbmanager

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
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
	panic("implement me") // TODO
}

func (m *mongoManagerImpl) Delete(ctx context.Context, args DeleteArgs) error {
	panic("implement me") // TODO
}

func (m *mongoManagerImpl) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
