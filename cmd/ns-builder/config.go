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
	"github.com/moby/buildkit/util/appdefaults"
	"github.com/spf13/viper"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/buildpack/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/buildpack/k8simpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
)

type Config struct {
	Buildkit struct {
		Address string `mapstructure:"address" yaml:"address"`
	} `mapstructure:"buildkit" yaml:"buildkit"`
	Buildpack struct {
		Backend string            `mapstructure:"backend" yaml:"backend"`
		Docker  dockerimpl.Config `mapstructure:"docker" yaml:"docker"`
		K8s     k8simpl.Config    `mapstructure:"k8s" yaml:"k8s"`
	}
	Repository struct {
		PrivateKeyFile string `mapstructure:"privateKeyFile" yaml:"privateKeyFile"`
	} `mapstructure:"repository" yaml:"repository"`
	Controller grpc.ControllerServiceClientConfig `mapstructure:"controller" yaml:"controller"`
	DB         repository.Config                  `mapstructure:"db" yaml:"db"`
	Storage    domain.StorageConfig               `mapstructure:"storage" yaml:"storage"`
	Image      builder.ImageConfig                `mapstructure:"image" yaml:"image"`
}

func init() {
	viper.SetDefault("buildkit.address", appdefaults.Address)

	viper.SetDefault("buildpack.backend", "docker")
	viper.SetDefault("buildpack.docker.containerName", "buildpack")
	viper.SetDefault("buildpack.docker.remoteDir", "/workspace")
	viper.SetDefault("buildpack.docker.user", "cnb")
	viper.SetDefault("buildpack.docker.group", "cnb")
	viper.SetDefault("buildpack.docker.platformAPI", "0.11")
	viper.SetDefault("buildpack.k8s.namespace", "ns-system")
	viper.SetDefault("buildpack.k8s.podName", "ns-builder")
	viper.SetDefault("buildpack.k8s.containerName", "buildpack")
	viper.SetDefault("buildpack.k8s.localDir", "/neoshowcase/buildpack")
	viper.SetDefault("buildpack.k8s.remoteDir", "/workspace")
	viper.SetDefault("buildpack.k8s.user", 1000)
	viper.SetDefault("buildpack.k8s.group", 1000)
	viper.SetDefault("buildpack.k8s.platformAPI", "0.11")

	viper.SetDefault("repository.privateKeyFile", "")

	viper.SetDefault("controller.url", "http://ns-controller:10000")

	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", 3306)
	viper.SetDefault("db.username", "root")
	viper.SetDefault("db.password", "password")
	viper.SetDefault("db.database", "neoshowcase")
	viper.SetDefault("db.connection.maxOpen", 0)
	viper.SetDefault("db.connection.maxIdle", 2)
	viper.SetDefault("db.connection.lifetime", 0)

	viper.SetDefault("storage.type", "local")
	viper.SetDefault("storage.local.dir", "/neoshowcase")
	viper.SetDefault("storage.s3.bucket", "neoshowcase")
	viper.SetDefault("storage.s3.endpoint", "")
	viper.SetDefault("storage.s3.region", "")
	viper.SetDefault("storage.s3.accessKey", "")
	viper.SetDefault("storage.s3.accessSecret", "")
	viper.SetDefault("storage.swift.username", "")
	viper.SetDefault("storage.swift.apiKey", "")
	viper.SetDefault("storage.swift.tenantName", "")
	viper.SetDefault("storage.swift.tenantId", "")
	viper.SetDefault("storage.swift.container", "neoshowcase")
	viper.SetDefault("storage.swift.authUrl", "")

	viper.SetDefault("image.registry.scheme", "https")
	viper.SetDefault("image.registry.addr", "localhost")
	viper.SetDefault("image.registry.username", "")
	viper.SetDefault("image.registry.password", "")
	viper.SetDefault("image.namePrefix", "ns-apps/")
}

func provideBuildpackBackend(c Config) (builder.BuildpackBackend, error) {
	switch c.Buildpack.Backend {
	case "docker":
		return dockerimpl.NewBuildpackBackend(c.Buildpack.Docker, c.Image)
	case "k8s":
		return k8simpl.NewBuildpackBackend(c.Buildpack.K8s, c.Image)
	default:
		return nil, errors.Errorf("invalid buildpack backend: %v", c.Buildpack.Backend)
	}
}

func provideRepositoryPublicKey(c Config) (*ssh.PublicKeys, error) {
	bytes, err := os.ReadFile(c.Repository.PrivateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open private key file")
	}
	return ssh.NewPublicKeys("", bytes, "")
}

func provideStorage(c domain.StorageConfig) (domain.Storage, error) {
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
