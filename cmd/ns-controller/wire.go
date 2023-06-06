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

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/k8simpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/webhook"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/usecase/cdservice"
	"github.com/traPtitech/neoshowcase/pkg/usecase/cleaner"
	"github.com/traPtitech/neoshowcase/pkg/usecase/logstream"
	"github.com/traPtitech/neoshowcase/pkg/usecase/repofetcher"
	"github.com/traPtitech/neoshowcase/pkg/usecase/sshserver"
)

var commonSet = wire.NewSet(
	dbmanager.NewMariaDBManager,
	dbmanager.NewMongoDBManager,
	repository.New,
	repository.NewApplicationRepository,
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
	webhook.NewReceiver,
	apiserver.NewService,
	cdservice.NewAppDeployHelper,
	cdservice.NewContainerStateMutator,
	cdservice.NewService,
	repofetcher.NewService,
	cleaner.NewService,
	logstream.NewService,
	sshserver.NewSSHServer,
	providePublicKey,
	provideStorage,
	provideControllerServer,
	wire.FieldsOf(new(Config), "Docker", "K8s", "SSH", "Webhook", "DB", "Storage", "Image"),
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
