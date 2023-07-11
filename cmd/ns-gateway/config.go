package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/spf13/viper"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/log/loki"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/metrics/prometheus"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
)

type Config struct {
	Port           int                                `mapstructure:"port" yaml:"port"`
	Debug          bool                               `mapstructure:"debug" yaml:"debug"`
	PrivateKeyFile string                             `mapstructure:"privateKeyFile" yaml:"privateKeyFile"`
	AvatarBaseURL  domain.AvatarBaseURL               `mapstructure:"avatarBaseURL" yaml:"avatarBaseURL"`
	AuthHeader     grpc.AuthHeader                    `mapstructure:"authHeader" yaml:"authHeader"`
	Controller     grpc.ControllerServiceClientConfig `mapstructure:"controller" yaml:"controller"`
	DB             repository.Config                  `mapstructure:"db" yaml:"db"`
	MariaDB        dbmanager.MariaDBConfig            `mapstructure:"mariadb" yaml:"mariadb"`
	MongoDB        dbmanager.MongoDBConfig            `mapstructure:"mongodb" yaml:"mongodb"`
	Storage        domain.StorageConfig               `mapstructure:"storage" yaml:"storage"`
	Log            struct {
		Type string      `mapstructure:"type" yaml:"type"`
		Loki loki.Config `mapstructure:"loki" yaml:"loki"`
	} `mapstructure:"log" yaml:"log"`
	Metrics struct {
		Type       string            `mapstructure:"type" yaml:"type"`
		Prometheus prometheus.Config `mapstructure:"prometheus" yaml:"prometheus"`
	}
}

func init() {
	viper.SetDefault("port", 8080)
	viper.SetDefault("debug", false)
	viper.SetDefault("privateKeyFile", "")
	viper.SetDefault("avatarBaseURL", "https://q.trap.jp/api/v3/public/icon/")
	viper.SetDefault("authHeader", "X-Showcase-User")

	viper.SetDefault("controller.url", "http://ns-controller:10000")

	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", 3306)
	viper.SetDefault("db.username", "root")
	viper.SetDefault("db.password", "password")
	viper.SetDefault("db.database", "neoshowcase")
	viper.SetDefault("db.connection.maxOpen", 0)
	viper.SetDefault("db.connection.maxIdle", 2)
	viper.SetDefault("db.connection.lifetime", 0)

	viper.SetDefault("mariadb.host", "mariadb")
	viper.SetDefault("mariadb.port", 3306)
	viper.SetDefault("mariadb.adminUser", "root")
	viper.SetDefault("mariadb.adminPassword", "password")

	viper.SetDefault("mongodb.host", "mongo")
	viper.SetDefault("mongodb.port", 27017)
	viper.SetDefault("mongodb.adminUser", "root")
	viper.SetDefault("mongodb.adminPassword", "password")

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

	viper.SetDefault("log.type", "loki")
	viper.SetDefault("log.loki.endpoint", "http://loki:3100")
	viper.SetDefault("log.loki.queryTemplate", loki.DefaultQueryTemplate())

	viper.SetDefault("metrics.type", "prometheus")
	viper.SetDefault("metrics.endpoint", "http://prometheus:9090")
	viper.SetDefault("metric.queries", prometheus.DefaultQueriesConfig())
}

func providePublicKey(c Config) (*ssh.PublicKeys, error) {
	bytes, err := os.ReadFile(c.PrivateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open private key file")
	}
	return domain.NewPublicKey(bytes)
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

func provideContainerLogger(c Config) (domain.ContainerLogger, error) {
	switch c.Log.Type {
	case "loki":
		return loki.NewLokiStreamer(c.Log.Loki)
	default:
		return nil, errors.Errorf("invalid log type: %v (supported values: loki)", c.Log.Type)
	}
}

func provideMetricsService(c Config) (domain.MetricsService, error) {
	switch c.Metrics.Type {
	case "prometheus":
		return prometheus.NewPromClient(c.Metrics.Prometheus)
	default:
		return nil, errors.Errorf("invalid metrics type: %v (supported values: prometheus)", c.Metrics.Type)
	}
}

type gatewayServer struct {
	*web.H2CServer
}

func provideGatewayServer(
	c Config,
	appService pbconnect.APIServiceHandler,
	authInterceptor *grpc.AuthInterceptor,
) *gatewayServer {
	wc := web.H2CConfig{
		Port: c.Port,
		SetupRoute: func(mux *http.ServeMux) {
			mux.Handle(pbconnect.NewAPIServiceHandler(
				appService,
				connect.WithInterceptors(authInterceptor),
			))
		},
	}
	return &gatewayServer{web.NewH2CServer(wc)}
}
