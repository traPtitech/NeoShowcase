package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/friendsofgo/errors"
	buildkit "github.com/moby/buildkit/client"

	authdev "github.com/traPtitech/neoshowcase/cmd/auth-dev"
	buildpackhelper "github.com/traPtitech/neoshowcase/cmd/buildpack-helper"
	"github.com/traPtitech/neoshowcase/cmd/controller"
	"github.com/traPtitech/neoshowcase/cmd/gateway"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/log/loki"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/log/victorialogs"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/metrics/prometheus"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver/builtin"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver/caddy"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
	ubuilder "github.com/traPtitech/neoshowcase/pkg/usecase/builder"
	giteaintegration "github.com/traPtitech/neoshowcase/pkg/usecase/gitea-integration"
	"github.com/traPtitech/neoshowcase/pkg/usecase/healthcheck"
	"github.com/traPtitech/neoshowcase/pkg/usecase/ssgen"
	"github.com/traPtitech/neoshowcase/pkg/usecase/systeminfo"
)

func provideRepositoryPrivateKey(c Config) (domain.PrivateKey, error) {
	bytes, err := os.ReadFile(c.PrivateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open private key file")
	}
	return bytes, nil
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

func provideTokenAuthInterceptor(c Config) (*grpc.TokenAuthInterceptor, error) {
	return grpc.NewTokenAuthInterceptor(
		c.Components.Controller.TokenHeader,
		c.Components.Controller.Token,
	)
}

func provideControllerBuilderServiceClient(c Config, auth *grpc.TokenAuthInterceptor) domain.ControllerBuilderServiceClient {
	return grpc.NewControllerBuilderServiceClient(
		c.Components.Builder.Controller,
		c.Components.Builder.Priority,
		auth,
	)
}

func provideBuildpackHelperClient(c Config) domain.BuildpackHelperServiceClient {
	return grpc.NewBuildpackHelperServiceClient(c.Components.Builder.Buildpack.Helper.Address)
}

func provideBuildpackHelperServer(
	c Config,
	helperServiceHandler pbconnect.BuildpackHelperServiceHandler,
) *buildpackhelper.APIServer {
	wc := web.H2CConfig{
		Port: c.Components.Builder.Buildpack.Helper.ListenPort,
		SetupRoute: func(mux *http.ServeMux) {
			mux.Handle(pbconnect.NewBuildpackHelperServiceHandler(helperServiceHandler))
		},
	}
	return &buildpackhelper.APIServer{H2CServer: web.NewH2CServer(wc)}
}

func provideBuilderConfig(c Config) (*ubuilder.Config, error) {
	stepTimeoutStr := c.Components.Builder.StepTimeout
	stepTimeout, err := time.ParseDuration(stepTimeoutStr)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse components.builder.stepTimeout value: %s", stepTimeoutStr))
	}
	if stepTimeout <= 0 {
		return nil, errors.Errorf("components.builder.stepTimeout must be positive: %s", stepTimeoutStr)
	}
	return &ubuilder.Config{
		StepTimeout: stepTimeout,
	}, nil
}

func provideBuildkitClient(c Config) (*buildkit.Client, error) {
	cc := c.Components.Builder
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := buildkit.New(ctx, cc.Buildkit.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize Buildkit Client")
	}
	return client, nil
}

func provideSystemInfoConfig(c Config) *systeminfo.ServiceConfig {
	return &systeminfo.ServiceConfig{
		AdditionalLinks: c.AdditionalLinks,
	}
}

func provideControllerServer(
	c Config,
	controllerHandler pbconnect.ControllerServiceHandler,
	builderHandler domain.ControllerBuilderService,
	ssgenHandler domain.ControllerSSGenService,
	giteaIntegrationHandler domain.ControllerGiteaIntegrationService,
	tokenAuth *grpc.TokenAuthInterceptor,
) *controller.APIServer {
	wc := web.H2CConfig{
		Port: c.Components.Controller.Port,
		SetupRoute: func(mux *http.ServeMux) {
			mux.Handle(pbconnect.NewControllerServiceHandler(controllerHandler))
			mux.Handle(pbconnect.NewControllerBuilderServiceHandler(builderHandler, connect.WithInterceptors(tokenAuth)))
			mux.Handle(pbconnect.NewControllerSSGenServiceHandler(ssgenHandler))
			mux.Handle(pbconnect.NewControllerGiteaIntegrationServiceHandler(giteaIntegrationHandler))
		},
	}
	return &controller.APIServer{H2CServer: web.NewH2CServer(wc)}
}

func provideContainerLogger(c Config) (domain.ContainerLogger, error) {
	cc := c.Components.Gateway
	switch cc.Log.Type {
	case "loki":
		return loki.NewLokiStreamer(cc.Log.Loki)
	case "victorialogs":
		return victorialogs.NewVictoriaLogsStreamer(cc.Log.VictoriaLogs)
	default:
		return nil, errors.Errorf("invalid log type: %v (supported values: loki, victorialogs)", cc.Log.Type)
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
	cacheInterceptor *grpc.CacheInterceptor,
) *gateway.APIServer {
	wc := web.H2CConfig{
		Port: c.Components.Gateway.Port,
		SetupRoute: func(mux *http.ServeMux) {
			mux.Handle(pbconnect.NewAPIServiceHandler(
				appService,
				connect.WithInterceptors(
					authInterceptor, // Make sure auth is the outermost interceptor
					cacheInterceptor,
				),
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
		Concurrency:     cc.Concurrency,
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
