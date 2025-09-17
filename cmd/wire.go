//go:build wireinject
// +build wireinject

package main

import (
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/friendsofgo/errors"
	"github.com/google/wire"
	traefikv1alpha1 "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefikio/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	authdev "github.com/traPtitech/neoshowcase/cmd/auth-dev"
	"github.com/traPtitech/neoshowcase/cmd/builder"
	buildpackhelper "github.com/traPtitech/neoshowcase/cmd/buildpack-helper"
	"github.com/traPtitech/neoshowcase/cmd/controller"
	"github.com/traPtitech/neoshowcase/cmd/gateway"
	giteaintegration "github.com/traPtitech/neoshowcase/cmd/gitea-integration"
	"github.com/traPtitech/neoshowcase/cmd/ssgen"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/k8simpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/buildpack"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/git"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/observability"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/registry"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/webhook"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	ubuilder "github.com/traPtitech/neoshowcase/pkg/usecase/builder"
	buildermock "github.com/traPtitech/neoshowcase/pkg/usecase/builder/mock"
	"github.com/traPtitech/neoshowcase/pkg/usecase/cdservice"
	"github.com/traPtitech/neoshowcase/pkg/usecase/cleaner"
	commitfetcher "github.com/traPtitech/neoshowcase/pkg/usecase/commit-fetcher"
	ugiteaintegration "github.com/traPtitech/neoshowcase/pkg/usecase/gitea-integration"
	"github.com/traPtitech/neoshowcase/pkg/usecase/healthcheck"
	"github.com/traPtitech/neoshowcase/pkg/usecase/logstream"
	"github.com/traPtitech/neoshowcase/pkg/usecase/repofetcher"
	ussgen "github.com/traPtitech/neoshowcase/pkg/usecase/ssgen"
	"github.com/traPtitech/neoshowcase/pkg/usecase/sshserver"
	"github.com/traPtitech/neoshowcase/pkg/usecase/systeminfo"
	"github.com/traPtitech/neoshowcase/pkg/util/discovery"
)

var providers = wire.NewSet(
	apiserver.NewService,
	cdservice.NewAppDeployHelper,
	cdservice.NewContainerStateMutator,
	cdservice.NewService,
	certmanagerv1.NewForConfig,
	cleaner.NewService,
	commitfetcher.NewService,
	dbmanager.NewMariaDBManager,
	dbmanager.NewMongoDBManager,
	dockerimpl.NewClientFromEnv,
	dockerimpl.NewDockerBackend,
	ugiteaintegration.NewIntegration,
	grpc.NewAPIServiceServer,
	grpc.NewAuthInterceptor,
	grpc.NewBuildpackHelperService,
	provideBuildpackHelperClient,
	grpc.NewCacheInterceptor,
	grpc.NewControllerService,
	grpc.NewControllerServiceClient,
	grpc.NewControllerBuilderService,
	grpc.NewGiteaIntegrationService,
	provideTokenAuthInterceptor,
	provideControllerBuilderServiceClient,
	grpc.NewControllerSSGenService,
	grpc.NewControllerSSGenServiceClient,
	healthcheck.NewServer,
	k8simpl.NewK8SBackend,
	kubernetes.NewForConfig,
	logstream.NewService,
	repofetcher.NewService,
	repository.New,
	repository.NewApplicationRepository,
	repository.NewArtifactRepository,
	repository.NewRuntimeImageRepository,
	repository.NewBuildRepository,
	repository.NewEnvironmentRepository,
	repository.NewGitRepositoryRepository,
	repository.NewRepositoryCommitRepository,
	repository.NewUserRepository,
	repository.NewWebsiteRepository,
	rest.InClusterConfig,
	traefikv1alpha1.NewForConfig,
	ussgen.NewGeneratorService,
	sshserver.NewSSHServer,
	systeminfo.NewService,
	ubuilder.NewService,
	webhook.NewReceiver,
	provideRepositoryPrivateKey,
	domain.IntoPublicKey,
	git.NewService,
	registry.NewClient,
	observability.NewMetricsServer,
	observability.NewControllerMetrics,
	provideStorage,
	provideAuthDevServer,
	provideBuildpackHelperServer,
	buildpack.NewBuildpackBackend,
	provideDiscoverer,
	discovery.NewCluster,
	provideBuilderConfig,
	provideBuildkitClient,
	provideSystemInfoConfig,
	provideControllerServer,
	provideContainerLogger,
	provideMetricsService,
	provideGatewayServer,
	provideGiteaIntegrationConfig,
	provideGiteaIntegrationServiceClient,
	provideGiteaIntegrationAPIServer,
	provideHealthCheckFunc,
	provideStaticServer,
	provideStaticServerDocumentRootPath,
	wire.FieldsOf(new(Config), "DB", "Storage", "Image", "Components"),
	wire.FieldsOf(new(ComponentsConfig), "Builder", "Controller", "Gateway", "GiteaIntegration", "SSGen"),
)

func NewAuthDev(c Config) (component, error) {
	wire.Build(
		providers,
		wire.Bind(new(component), new(*authdev.Server)),
	)
	return nil, nil
}

func NewBuilder(c Config) (component, error) {
	if c.Components.Builder.Mock {
		return newMockBuilder(c)
	} else {
		return newBuilder(c)
	}
}

func newBuilder(c Config) (component, error) {
	wire.Build(
		providers,
		wire.FieldsOf(new(BuilderConfig), "Buildpack"),
		wire.Bind(new(ubuilder.Service), new(*ubuilder.ServiceImpl)),
		wire.Bind(new(component), new(*builder.Server)),
		wire.Struct(new(builder.Server), "*"),
	)
	return nil, nil
}

func newMockBuilder(c Config) (component, error) {
	wire.Build(
		providers,
		buildermock.NewBuilderServiceMock,
		wire.Bind(new(ubuilder.Service), new(*buildermock.BuilderServiceMock)),
		wire.Bind(new(component), new(*builder.Server)),
		wire.Struct(new(builder.Server), "*"),
	)
	return nil, nil
}

func NewBuildpackHelper(c Config) (component, error) {
	wire.Build(
		providers,
		wire.Bind(new(component), new(*buildpackhelper.Server)),
		wire.Struct(new(buildpackhelper.Server), "*"),
	)
	return nil, nil
}

func NewController(c Config) (component, error) {
	switch c.Components.Controller.Mode {
	case "docker":
		return NewControllerDocker(c)
	case "k8s", "kubernetes":
		return NewControllerK8s(c)
	}
	return nil, errors.New("unknown mode: " + c.Components.Controller.Mode)
}

func NewControllerDocker(c Config) (component, error) {
	wire.Build(
		providers,
		wire.FieldsOf(new(ControllerConfig), "Port", "Docker", "SSH", "Webhook", "Metrics"),
		wire.Bind(new(domain.Backend), new(*dockerimpl.Backend)),
		wire.Bind(new(component), new(*controller.Server)),
		wire.Struct(new(controller.Server), "*"),
	)
	return nil, nil
}

func NewControllerK8s(c Config) (component, error) {
	wire.Build(
		providers,
		wire.FieldsOf(new(ControllerConfig), "Port", "K8s", "SSH", "Webhook", "Metrics"),
		wire.Bind(new(domain.Backend), new(*k8simpl.Backend)),
		wire.Bind(new(component), new(*controller.Server)),
		wire.Struct(new(controller.Server), "*"),
	)
	return nil, nil
}

func NewGateway(c Config) (component, error) {
	wire.Build(
		providers,
		wire.FieldsOf(new(GatewayConfig), "AvatarBaseURL", "AuthHeader", "Controller", "MariaDB", "MongoDB"),
		wire.Bind(new(component), new(*gateway.Server)),
		wire.Struct(new(gateway.Server), "*"),
	)
	return nil, nil
}

func NewGiteaIntegration(c Config) (component, error) {
	wire.Build(
		providers,
		wire.Bind(new(component), new(*giteaintegration.Server)),
		wire.Struct(new(giteaintegration.Server), "*"),
	)
	return nil, nil
}

func NewSSGen(c Config) (component, error) {
	wire.Build(
		providers,
		wire.FieldsOf(new(SSGenConfig), "HealthPort", "Controller"),
		wire.Bind(new(component), new(*ssgen.Server)),
		wire.Struct(new(ssgen.Server), "*"),
	)
	return nil, nil
}
