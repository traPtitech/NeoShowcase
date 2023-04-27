//go:build wireinject
// +build wireinject

package main

import (
	"fmt"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/google/wire"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefikio/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/k8simpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

var commonSet = wire.NewSet(
	admindb.New,
	dbmanager.NewMariaDBManager,
	dbmanager.NewMongoDBManager,
	repository.NewApplicationRepository,
	repository.NewAvailableDomainRepository,
	repository.NewGitRepositoryRepository,
	repository.NewEnvironmentRepository,
	repository.NewBuildRepository,
	repository.NewArtifactRepository,
	repository.NewUserRepository,
	grpc.NewAPIServiceServer,
	grpc.NewAuthInterceptor,
	grpc.NewControllerService,
	grpc.NewControllerBuilderService,
	grpc.NewControllerSSGenService,
	usecase.NewAPIServerService,
	usecase.NewAppBuildHelper,
	usecase.NewAppDeployHelper,
	usecase.NewContinuousDeploymentService,
	usecase.NewRepositoryFetcherService,
	usecase.NewCleanerService,
	usecase.NewLogStreamService,
	usecase.NewContainerStateMutator,
	provideRepositoryPublicKey,
	provideStorage,
	provideControllerServer,
	wire.FieldsOf(new(Config), "DB", "Storage", "Docker", "K8s", "Image"),
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
		dockerimpl.NewClientFromEnv,
		dockerimpl.NewDockerBackend,
	)
	return nil, nil
}

func NewWithK8S(c Config) (*Server, error) {
	wire.Build(
		commonSet,
		rest.InClusterConfig,
		kubernetes.NewForConfig,
		traefikv1alpha1.NewForConfig,
		certmanagerv1.NewForConfig,
		k8simpl.NewK8SBackend,
	)
	return nil, nil
}
