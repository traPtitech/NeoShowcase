package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"connectrpc.com/connect"
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/google/wire"
	buildkit "github.com/moby/buildkit/client"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefikio/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	authdev "github.com/traPtitech/neoshowcase/cmd/auth-dev"
	"github.com/traPtitech/neoshowcase/cmd/controller"
	"github.com/traPtitech/neoshowcase/cmd/gateway"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/k8simpl"
	bdockerimpl "github.com/traPtitech/neoshowcase/pkg/infrastructure/buildpack/dockerimpl"
	bk8simpl "github.com/traPtitech/neoshowcase/pkg/infrastructure/buildpack/k8simpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/log/loki"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/metrics/prometheus"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver/builtin"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver/caddy"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/webhook"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	ubuilder "github.com/traPtitech/neoshowcase/pkg/usecase/builder"
	"github.com/traPtitech/neoshowcase/pkg/usecase/cdservice"
	"github.com/traPtitech/neoshowcase/pkg/usecase/cleaner"
	giteaintegration "github.com/traPtitech/neoshowcase/pkg/usecase/gitea-integration"
	"github.com/traPtitech/neoshowcase/pkg/usecase/healthcheck"
	"github.com/traPtitech/neoshowcase/pkg/usecase/logstream"
	"github.com/traPtitech/neoshowcase/pkg/usecase/repofetcher"
	"github.com/traPtitech/neoshowcase/pkg/usecase/ssgen"
	"github.com/traPtitech/neoshowcase/pkg/usecase/sshserver"
)

var providers = wire.NewSet(
	apiserver.NewService,
	cdservice.NewAppDeployHelper,
	cdservice.NewContainerStateMutator,
	cdservice.NewService,
	certmanagerv1.NewForConfig,
	cleaner.NewService,
	dbmanager.NewMariaDBManager,
	dbmanager.NewMongoDBManager,
	dockerimpl.NewClientFromEnv,
	dockerimpl.NewDockerBackend,
	giteaintegration.NewIntegration,
	grpc.NewAPIServiceServer,
	grpc.NewAuthInterceptor,
	grpc.NewControllerService,
	grpc.NewControllerServiceClient,
	grpc.NewControllerBuilderService,
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
	repository.NewBuildRepository,
	repository.NewEnvironmentRepository,
	repository.NewGitRepositoryRepository,
	repository.NewUserRepository,
	rest.InClusterConfig,
	traefikv1alpha1.NewForConfig,
	ssgen.NewGeneratorService,
	sshserver.NewSSHServer,
	ubuilder.NewService,
	webhook.NewReceiver,
	provideRepositoryPublicKey,
	provideStorage,
	provideAuthDevServer,
	provideBuildpackBackend,
	provdeBuildkitClient,
	provideControllerServer,
	provideContainerLogger,
	provideMetricsService,
	provideGatewayServer,
	provideGiteaIntegrationConfig,
	provideHealthCheckFunc,
	provideStaticServer,
	provideStaticServerDocumentRootPath,
	wire.FieldsOf(new(Config), "AdminerURL", "DB", "Storage", "Image", "Components"),
	wire.FieldsOf(new(ComponentsConfig), "Builder", "Controller", "Gateway", "GiteaIntegration", "SSGen"),
)

func provideRepositoryPublicKey(c Config) (*ssh.PublicKeys, error) {
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

func provideAuthDevServer(c Config) *authdev.Server {
	cc := c.Components.AuthDev
	return authdev.NewServer(cc.Header, cc.Port, cc.User)
}

func provideControllerBuilderServiceClient(c Config) domain.ControllerBuilderServiceClient {
	return grpc.NewControllerBuilderServiceClient(
		c.Components.Builder.Controller,
		c.Components.Builder.Priority,
	)
}

func provideBuildpackBackend(c Config) (builder.BuildpackBackend, error) {
	cc := c.Components.Builder
	switch cc.Buildpack.Backend {
	case "docker":
		return bdockerimpl.NewBuildpackBackend(cc.Buildpack.Docker, c.Image)
	case "k8s":
		return bk8simpl.NewBuildpackBackend(cc.Buildpack.K8s, c.Image)
	default:
		return nil, errors.Errorf("invalid buildpack backend: %v", cc.Buildpack.Backend)
	}
}

func provdeBuildkitClient(c Config) (*buildkit.Client, error) {
	cc := c.Components.Builder
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := buildkit.New(ctx, cc.Buildkit.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize Buildkit Client")
	}
	return client, nil
}

func provideControllerServer(
	c Config,
	controllerHandler pbconnect.ControllerServiceHandler,
	builderHandler domain.ControllerBuilderService,
	ssgenHandler domain.ControllerSSGenService,
) *controller.APIServer {
	wc := web.H2CConfig{
		Port: c.Components.Controller.Port,
		SetupRoute: func(mux *http.ServeMux) {
			mux.Handle(pbconnect.NewControllerServiceHandler(controllerHandler))
			mux.Handle(pbconnect.NewControllerBuilderServiceHandler(builderHandler))
			mux.Handle(pbconnect.NewControllerSSGenServiceHandler(ssgenHandler))
		},
	}
	return &controller.APIServer{H2CServer: web.NewH2CServer(wc)}
}

func provideContainerLogger(c Config) (domain.ContainerLogger, error) {
	cc := c.Components.Gateway
	switch cc.Log.Type {
	case "loki":
		return loki.NewLokiStreamer(cc.Log.Loki)
	default:
		return nil, errors.Errorf("invalid log type: %v (supported values: loki)", cc.Log.Type)
	}
}

func provideMetricsService(c Config) (domain.MetricsService, error) {
	cc := c.Components.Gateway
	switch cc.Metrics.Type {
	case "prometheus":
		return prometheus.NewPromClient(cc.Metrics.Prometheus)
	default:
		return nil, errors.Errorf("invalid metrics type: %v (supported values: prometheus)", cc.Metrics.Type)
	}
}

func provideGatewayServer(
	c Config,
	appService pbconnect.APIServiceHandler,
	authInterceptor *grpc.AuthInterceptor,
) *gateway.APIServer {
	wc := web.H2CConfig{
		Port: c.Components.Gateway.Port,
		SetupRoute: func(mux *http.ServeMux) {
			mux.Handle(pbconnect.NewAPIServiceHandler(
				appService,
				connect.WithInterceptors(authInterceptor),
			))
		},
	}
	return &gateway.APIServer{H2CServer: web.NewH2CServer(wc)}
}

func provideGiteaIntegrationConfig(c Config) giteaintegration.Config {
	cc := c.Components.GiteaIntegration
	return giteaintegration.Config{
		URL:             cc.URL,
		Token:           cc.Token,
		IntervalSeconds: cc.IntervalSeconds,
		ListIntervalMs:  cc.ListIntervalMs,
	}
}

func provideHealthCheckFunc(gen ssgen.GeneratorService) healthcheck.Func {
	return gen.Healthy
}

func provideStaticServer(c Config) (domain.StaticServer, error) {
	cc := c.Components.SSGen
	switch cc.Server.Type {
	case "builtIn":
		return builtin.NewServer(cc.Server.BuiltIn, cc.ArtifactsRoot), nil
	case "caddy":
		return caddy.NewServer(cc.Server.Caddy), nil
	default:
		return nil, errors.Errorf("invalid static server type: %v", cc.Server.Type)
	}
}

func provideStaticServerDocumentRootPath(c Config) domain.StaticServerDocumentRootPath {
	cc := c.Components.SSGen
	return domain.StaticServerDocumentRootPath(cc.ArtifactsRoot)
}
