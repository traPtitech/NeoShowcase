//go:build wireinject
// +build wireinject

package main

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/google/wire"
	"github.com/leandro-lugaresi/hub"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefik/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/k8simpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
	"github.com/traPtitech/neoshowcase/pkg/interface/broker"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

var commonSet = wire.NewSet(
	web.NewServer,
	hub.New,
	eventbus.NewLocal,
	admindb.New,
	dbmanager.NewMariaDBManager,
	dbmanager.NewMongoDBManager,
	repository.NewApplicationRepository,
	repository.NewAvailableDomainRepository,
	repository.NewGitRepositoryRepository,
	repository.NewEnvironmentRepository,
	repository.NewBuildRepository,
	grpc.NewServer,
	grpc.NewApplicationServiceServer,
	grpc.NewBuilderServiceClientConn,
	grpc.NewStaticSiteServiceClientConn,
	grpc.NewBuilderServiceClient,
	grpc.NewStaticSiteServiceClient,
	broker.NewBuilderEventsBroker,
	usecase.NewAPIServerService,
	usecase.NewAppBuildService,
	usecase.NewAppDeployService,
	usecase.NewContinuousDeploymentService,
	usecase.NewRepositoryFetcherService,
	handlerSet,
	provideGRPCPort,
	provideWebServerConfig,
	provideImagePrefix,
	provideImageRegistry,
	provideRepositoryFetcherCacheDir,
	wire.FieldsOf(new(Config), "Builder", "SSGen", "DB", "MariaDB", "MongoDB"),
	wire.Struct(new(Router), "*"),
	wire.Bind(new(web.Router), new(*Router)),
	wire.Struct(new(Server), "*"),
)

func New(c Config) (*Server, error) {
	switch c.GetMode() {
	case ModeDocker:
		return NewWithDocker(c)
	case ModeK8s:
		return NewWithK8S(c)
	default:
		return nil, fmt.Errorf("unknown mode: %s", c.Mode)
	}
}

func NewWithDocker(c Config) (*Server, error) {
	wire.Build(
		commonSet,
		docker.NewClientFromEnv,
		dockerimpl.NewDockerBackend,
		wire.Value(dockerimpl.IngressConfDirPath("/opt/traefik/conf")),
	)
	return nil, nil
}

func NewWithK8S(c Config) (*Server, error) {
	wire.Build(
		commonSet,
		rest.InClusterConfig,
		kubernetes.NewForConfig,
		traefikv1alpha1.NewForConfig,
		k8simpl.NewK8SBackend,
	)
	return nil, nil
}
