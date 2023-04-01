package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

const (
	ModeDocker = iota
	ModeK8s
)

type Config struct {
	Debug   bool                                  `mapstructure:"debug" yaml:"debug"`
	Mode    string                                `mapstructure:"mode" yaml:"mode"`
	SS      domain.StaticServerConnectivityConfig `mapstructure:"ss" yaml:"ss"`
	DB      admindb.Config                        `mapstructure:"db" yaml:"db"`
	MariaDB dbmanager.MariaDBConfig               `mapstructure:"mariadb" yaml:"mariadb"`
	MongoDB dbmanager.MongoDBConfig               `mapstructure:"mongodb" yaml:"mongodb"`
	Storage domain.StorageConfig                  `mapstructure:"storage" yaml:"storage"`
	Docker  struct {
		ConfDir string `mapstructure:"confDir" yaml:"confDir"`
	} `mapstructure:"docker" yaml:"docker"`
	Web struct {
		App struct {
			Port int `mapstructure:"port" yaml:"port"`
		} `mapstructure:"app" yaml:"app"`
		Component struct {
			Port int `mapstructure:"port" yaml:"port"`
		} `mapstructure:"component" yaml:"component"`
	} `mapstructure:"web" yaml:"web"`
	Repository struct {
		CacheDir       string `mapstructure:"cacheDir" yaml:"cacheDir"`
		PrivateKeyFile string `mapstructure:"privateKeyFile" yaml:"privateKeyFile"`
	} `mapstructure:"repository" yaml:"repository"`
	Image builder.ImageConfig `mapstructure:"image" yaml:"image"`
}

func (c *Config) GetMode() int {
	switch strings.ToLower(c.Mode) {
	case "k8s", "kubernetes":
		return ModeK8s
	case "docker":
		return ModeDocker
	default:
		return ModeDocker
	}
}

func provideIngressConfDirPath(c Config) dockerimpl.IngressConfDirPath {
	return dockerimpl.IngressConfDirPath(c.Docker.ConfDir)
}

func provideRepositoryFetcherCacheDir(c Config) usecase.RepositoryFetcherCacheDir {
	return usecase.RepositoryFetcherCacheDir(c.Repository.CacheDir)
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

type webAppServer struct {
	*web.H2CServer
}
type webComponentServer struct {
	*web.H2CServer
}

func provideWebAppServer(c Config, appService pbconnect.ApplicationServiceHandler) *webAppServer {
	wc := web.H2CConfig{
		Port: c.Web.App.Port,
		SetupRoute: func(mux *http.ServeMux) {
			mux.Handle(pbconnect.NewApplicationServiceHandler(appService))
		},
	}
	return &webAppServer{web.NewH2CServer(wc)}
}

func provideWebComponentServer(c Config, componentService domain.ComponentService) *webComponentServer {
	wc := web.H2CConfig{
		Port: c.Web.Component.Port,
		SetupRoute: func(mux *http.ServeMux) {
			mux.Handle(pbconnect.NewComponentServiceHandler(componentService))
		},
	}
	return &webComponentServer{web.NewH2CServer(wc)}
}
