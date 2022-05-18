//go:build wireinject
// +build wireinject

package main

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/google/wire"
	"github.com/leandro-lugaresi/hub"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/k8simpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
	"github.com/traPtitech/neoshowcase/pkg/interface/broker"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var commonSet = wire.NewSet(
	web.NewServer,
	dbmanager.NewMariaDBManager,
	dbmanager.NewMongoManager,
	usecase.NewGitPushWebhookService,
	usecase.NewAppBuildService,
	usecase.NewAppDeployService,
	usecase.NewContinuousDeploymentService,
	repository.NewApplicationRepository,
	repository.NewGitRepositoryRepository,
	repository.NewBuildLogRepository,
	broker.NewBuilderEventsBroker,
	eventbus.NewLocal,
	admindb.New,
	handlerSet,
	provideWebServerConfig,
	provideImagePrefix,
	provideImageRegistry,
	hub.New,
	grpc.NewBuilderServiceClientConn,
	grpc.NewStaticSiteServiceClientConn,
	grpc.NewBuilderServiceClient,
	grpc.NewStaticSiteServiceClient,
	wire.FieldsOf(new(Config), "Builder", "SSGen", "DB"),
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
		k8simpl.NewK8SBackend,
	)
	return nil, nil
}
