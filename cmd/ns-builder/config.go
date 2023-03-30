package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	buildkit "github.com/moby/buildkit/client"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
)

type Config struct {
	Buildkit struct {
		Address string `mapstructure:"address" yaml:"address"`
	} `mapstructure:"buildkit" yaml:"buildkit"`
	Repository struct {
		PrivateKeyFile string `mapstructure:"privateKeyFile" yaml:"privateKeyFile"`
	} `mapstructure:"repository" yaml:"repository"`
	NS      grpc.ComponentServiceClientConfig `mapstructure:"ns" yaml:"ns"`
	DB      admindb.Config                    `mapstructure:"db" yaml:"db"`
	Storage domain.StorageConfig              `mapstructure:"storage" yaml:"storage"`
}

func provideRepositoryPublicKey(c Config) (*ssh.PublicKeys, error) {
	bytes, err := os.ReadFile(c.Repository.PrivateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open private key file")
	}
	return ssh.NewPublicKeys("", bytes, "")
}

func initStorage(c domain.StorageConfig) (domain.Storage, error) {
	switch strings.ToLower(c.Type) {
	case "local":
		return storage.NewLocalStorage(c.Local.Dir)
	case "s3":
		return storage.NewS3Storage(c.S3.Bucket, c.S3.AccessKey, c.S3.AccessSecret, c.S3.Region, c.S3.Endpoint)
	case "swift":
		return storage.NewSwiftStorage(c.Swift.Container, c.Swift.UserName, c.Swift.APIKey, c.Swift.TenantName, c.Swift.TenantID, c.Swift.AuthURL)
	default:
		return nil, fmt.Errorf("unknown storage: %s", c.Type)
	}
}

func initBuildkitClient(c Config) (*buildkit.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := buildkit.New(ctx, c.Buildkit.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize Buildkit Client")
	}
	return client, nil
}
