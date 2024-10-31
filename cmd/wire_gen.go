// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/friendsofgo/errors"
	builder2 "github.com/traPtitech/neoshowcase/cmd/builder"
	"github.com/traPtitech/neoshowcase/cmd/buildpack-helper"
	"github.com/traPtitech/neoshowcase/cmd/controller"
	"github.com/traPtitech/neoshowcase/cmd/gateway"
	giteaintegration2 "github.com/traPtitech/neoshowcase/cmd/gitea-integration"
	ssgen2 "github.com/traPtitech/neoshowcase/cmd/ssgen"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/k8simpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/buildpack"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/webhook"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/usecase/builder"
	"github.com/traPtitech/neoshowcase/pkg/usecase/cdservice"
	"github.com/traPtitech/neoshowcase/pkg/usecase/cleaner"
	"github.com/traPtitech/neoshowcase/pkg/usecase/commit-fetcher"
	"github.com/traPtitech/neoshowcase/pkg/usecase/gitea-integration"
	"github.com/traPtitech/neoshowcase/pkg/usecase/healthcheck"
	"github.com/traPtitech/neoshowcase/pkg/usecase/logstream"
	"github.com/traPtitech/neoshowcase/pkg/usecase/repofetcher"
	"github.com/traPtitech/neoshowcase/pkg/usecase/ssgen"
	"github.com/traPtitech/neoshowcase/pkg/usecase/sshserver"
	"github.com/traPtitech/neoshowcase/pkg/usecase/systeminfo"
	"github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefikio/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

// Injectors from wire.go:

func NewAuthDev(c Config) (component, error) {
	server := provideAuthDevServer(c)
	return server, nil
}

func NewBuilder(c Config) (component, error) {
	client, err := provideBuildkitClient(c)
	if err != nil {
		return nil, err
	}
	builderConfig, err := provideBuilderConfig(c)
	if err != nil {
		return nil, err
	}
	tokenAuthInterceptor, err := provideTokenAuthInterceptor(c)
	if err != nil {
		return nil, err
	}
	controllerBuilderServiceClient := provideControllerBuilderServiceClient(c, tokenAuthInterceptor)
	componentsConfig := c.Components
	mainBuilderConfig := componentsConfig.Builder
	buildpackConfig := mainBuilderConfig.Buildpack
	buildpackHelperServiceClient := provideBuildpackHelperClient(c)
	buildpackBackend := buildpack.NewBuildpackBackend(buildpackConfig, buildpackHelperServiceClient)
	service, err := builder.NewService(builderConfig, controllerBuilderServiceClient, client, buildpackBackend)
	if err != nil {
		return nil, err
	}
	server := &builder2.Server{
		Buildkit: client,
		Builder:  service,
	}
	return server, nil
}

func NewBuildpackHelper(c Config) (component, error) {
	buildpackHelperServiceHandler := grpc.NewBuildpackHelperService()
	apiServer := provideBuildpackHelperServer(c, buildpackHelperServiceHandler)
	server := &buildpackhelper.Server{
		Helper: apiServer,
	}
	return server, nil
}

func NewControllerDocker(c Config) (component, error) {
	serviceConfig := provideSystemInfoConfig(c)
	client, err := dockerimpl.NewClientFromEnv()
	if err != nil {
		return nil, err
	}
	componentsConfig := c.Components
	controllerConfig := componentsConfig.Controller
	dockerimplConfig := controllerConfig.Docker
	imageConfig := c.Image
	backend, err := dockerimpl.NewDockerBackend(client, dockerimplConfig, imageConfig)
	if err != nil {
		return nil, err
	}
	repositoryConfig := c.DB
	db, err := repository.New(repositoryConfig)
	if err != nil {
		return nil, err
	}
	applicationRepository := repository.NewApplicationRepository(db)
	sshConfig := controllerConfig.SSH
	privateKey, err := provideRepositoryPrivateKey(c)
	if err != nil {
		return nil, err
	}
	publicKeys, err := domain.IntoPublicKey(privateKey)
	if err != nil {
		return nil, err
	}
	service := systeminfo.NewService(serviceConfig, backend, applicationRepository, sshConfig, publicKeys)
	gitRepositoryRepository := repository.NewGitRepositoryRepository(db)
	buildRepository := repository.NewBuildRepository(db)
	environmentRepository := repository.NewEnvironmentRepository(db)
	logstreamService := logstream.NewService()
	storageConfig := c.Storage
	storage, err := provideStorage(storageConfig)
	if err != nil {
		return nil, err
	}
	artifactRepository := repository.NewArtifactRepository(db)
	runtimeImageRepository := repository.NewRuntimeImageRepository(db)
	controllerBuilderService := grpc.NewControllerBuilderService(logstreamService, privateKey, imageConfig, storage, applicationRepository, artifactRepository, runtimeImageRepository, buildRepository, environmentRepository, gitRepositoryRepository)
	controllerSSGenService := grpc.NewControllerSSGenService()
	appDeployHelper := cdservice.NewAppDeployHelper(backend, applicationRepository, buildRepository, environmentRepository, controllerSSGenService, imageConfig)
	containerStateMutator := cdservice.NewContainerStateMutator(applicationRepository, backend)
	cdserviceService, err := cdservice.NewService(applicationRepository, buildRepository, environmentRepository, backend, controllerBuilderService, appDeployHelper, containerStateMutator)
	if err != nil {
		return nil, err
	}
	repositoryCommitRepository := repository.NewRepositoryCommitRepository(db)
	commitfetcherService, err := commitfetcher.NewService(applicationRepository, buildRepository, gitRepositoryRepository, repositoryCommitRepository, publicKeys)
	if err != nil {
		return nil, err
	}
	repofetcherService, err := repofetcher.NewService(applicationRepository, gitRepositoryRepository, publicKeys, cdserviceService, commitfetcherService)
	if err != nil {
		return nil, err
	}
	controllerServiceHandler := grpc.NewControllerService(service, repofetcherService, cdserviceService, controllerBuilderService, logstreamService)
	controllerGiteaIntegrationService := grpc.NewControllerGiteaIntegrationService()
	tokenAuthInterceptor, err := provideTokenAuthInterceptor(c)
	if err != nil {
		return nil, err
	}
	apiServer := provideControllerServer(c, controllerServiceHandler, controllerBuilderService, controllerSSGenService, controllerGiteaIntegrationService, tokenAuthInterceptor)
	userRepository := repository.NewUserRepository(db)
	sshServer := sshserver.NewSSHServer(sshConfig, publicKeys, backend, applicationRepository, userRepository)
	receiverConfig := controllerConfig.Webhook
	receiver := webhook.NewReceiver(receiverConfig, gitRepositoryRepository, repofetcherService, controllerGiteaIntegrationService)
	cleanerService, err := cleaner.NewService(artifactRepository, applicationRepository, buildRepository, imageConfig, storage)
	if err != nil {
		return nil, err
	}
	server := &controller.Server{
		APIServer:      apiServer,
		DB:             db,
		Backend:        backend,
		SSHServer:      sshServer,
		Webhook:        receiver,
		CDService:      cdserviceService,
		CommitFetcher:  commitfetcherService,
		FetcherService: repofetcherService,
		CleanerService: cleanerService,
	}
	return server, nil
}

func NewControllerK8s(c Config) (component, error) {
	serviceConfig := provideSystemInfoConfig(c)
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	traefikV1alpha1Client, err := v1alpha1.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	versionedClientset, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	componentsConfig := c.Components
	controllerConfig := componentsConfig.Controller
	k8simplConfig := controllerConfig.K8s
	backend, err := k8simpl.NewK8SBackend(restConfig, clientset, traefikV1alpha1Client, versionedClientset, k8simplConfig)
	if err != nil {
		return nil, err
	}
	repositoryConfig := c.DB
	db, err := repository.New(repositoryConfig)
	if err != nil {
		return nil, err
	}
	applicationRepository := repository.NewApplicationRepository(db)
	sshConfig := controllerConfig.SSH
	privateKey, err := provideRepositoryPrivateKey(c)
	if err != nil {
		return nil, err
	}
	publicKeys, err := domain.IntoPublicKey(privateKey)
	if err != nil {
		return nil, err
	}
	service := systeminfo.NewService(serviceConfig, backend, applicationRepository, sshConfig, publicKeys)
	gitRepositoryRepository := repository.NewGitRepositoryRepository(db)
	buildRepository := repository.NewBuildRepository(db)
	environmentRepository := repository.NewEnvironmentRepository(db)
	logstreamService := logstream.NewService()
	imageConfig := c.Image
	storageConfig := c.Storage
	storage, err := provideStorage(storageConfig)
	if err != nil {
		return nil, err
	}
	artifactRepository := repository.NewArtifactRepository(db)
	runtimeImageRepository := repository.NewRuntimeImageRepository(db)
	controllerBuilderService := grpc.NewControllerBuilderService(logstreamService, privateKey, imageConfig, storage, applicationRepository, artifactRepository, runtimeImageRepository, buildRepository, environmentRepository, gitRepositoryRepository)
	controllerSSGenService := grpc.NewControllerSSGenService()
	appDeployHelper := cdservice.NewAppDeployHelper(backend, applicationRepository, buildRepository, environmentRepository, controllerSSGenService, imageConfig)
	containerStateMutator := cdservice.NewContainerStateMutator(applicationRepository, backend)
	cdserviceService, err := cdservice.NewService(applicationRepository, buildRepository, environmentRepository, backend, controllerBuilderService, appDeployHelper, containerStateMutator)
	if err != nil {
		return nil, err
	}
	repositoryCommitRepository := repository.NewRepositoryCommitRepository(db)
	commitfetcherService, err := commitfetcher.NewService(applicationRepository, buildRepository, gitRepositoryRepository, repositoryCommitRepository, publicKeys)
	if err != nil {
		return nil, err
	}
	repofetcherService, err := repofetcher.NewService(applicationRepository, gitRepositoryRepository, publicKeys, cdserviceService, commitfetcherService)
	if err != nil {
		return nil, err
	}
	controllerServiceHandler := grpc.NewControllerService(service, repofetcherService, cdserviceService, controllerBuilderService, logstreamService)
	controllerGiteaIntegrationService := grpc.NewControllerGiteaIntegrationService()
	tokenAuthInterceptor, err := provideTokenAuthInterceptor(c)
	if err != nil {
		return nil, err
	}
	apiServer := provideControllerServer(c, controllerServiceHandler, controllerBuilderService, controllerSSGenService, controllerGiteaIntegrationService, tokenAuthInterceptor)
	userRepository := repository.NewUserRepository(db)
	sshServer := sshserver.NewSSHServer(sshConfig, publicKeys, backend, applicationRepository, userRepository)
	receiverConfig := controllerConfig.Webhook
	receiver := webhook.NewReceiver(receiverConfig, gitRepositoryRepository, repofetcherService, controllerGiteaIntegrationService)
	cleanerService, err := cleaner.NewService(artifactRepository, applicationRepository, buildRepository, imageConfig, storage)
	if err != nil {
		return nil, err
	}
	server := &controller.Server{
		APIServer:      apiServer,
		DB:             db,
		Backend:        backend,
		SSHServer:      sshServer,
		Webhook:        receiver,
		CDService:      cdserviceService,
		CommitFetcher:  commitfetcherService,
		FetcherService: repofetcherService,
		CleanerService: cleanerService,
	}
	return server, nil
}

func NewGateway(c Config) (component, error) {
	repositoryConfig := c.DB
	db, err := repository.New(repositoryConfig)
	if err != nil {
		return nil, err
	}
	artifactRepository := repository.NewArtifactRepository(db)
	runtimeImageRepository := repository.NewRuntimeImageRepository(db)
	applicationRepository := repository.NewApplicationRepository(db)
	buildRepository := repository.NewBuildRepository(db)
	environmentRepository := repository.NewEnvironmentRepository(db)
	gitRepositoryRepository := repository.NewGitRepositoryRepository(db)
	repositoryCommitRepository := repository.NewRepositoryCommitRepository(db)
	userRepository := repository.NewUserRepository(db)
	storageConfig := c.Storage
	storage, err := provideStorage(storageConfig)
	if err != nil {
		return nil, err
	}
	componentsConfig := c.Components
	gatewayConfig := componentsConfig.Gateway
	mariaDBConfig := gatewayConfig.MariaDB
	mariaDBManager, err := dbmanager.NewMariaDBManager(mariaDBConfig)
	if err != nil {
		return nil, err
	}
	mongoDBConfig := gatewayConfig.MongoDB
	mongoDBManager, err := dbmanager.NewMongoDBManager(mongoDBConfig)
	if err != nil {
		return nil, err
	}
	metricsService, err := provideMetricsService(c)
	if err != nil {
		return nil, err
	}
	containerLogger, err := provideContainerLogger(c)
	if err != nil {
		return nil, err
	}
	controllerServiceClientConfig := gatewayConfig.Controller
	controllerServiceClient := grpc.NewControllerServiceClient(controllerServiceClientConfig)
	imageConfig := c.Image
	privateKey, err := provideRepositoryPrivateKey(c)
	if err != nil {
		return nil, err
	}
	publicKeys, err := domain.IntoPublicKey(privateKey)
	if err != nil {
		return nil, err
	}
	service, err := apiserver.NewService(artifactRepository, runtimeImageRepository, applicationRepository, buildRepository, environmentRepository, gitRepositoryRepository, repositoryCommitRepository, userRepository, storage, mariaDBManager, mongoDBManager, metricsService, containerLogger, controllerServiceClient, imageConfig, publicKeys)
	if err != nil {
		return nil, err
	}
	avatarBaseURL := gatewayConfig.AvatarBaseURL
	apiServiceHandler := grpc.NewAPIServiceServer(service, avatarBaseURL)
	authHeader := gatewayConfig.AuthHeader
	authInterceptor := grpc.NewAuthInterceptor(userRepository, authHeader)
	cacheInterceptor := grpc.NewCacheInterceptor()
	apiServer := provideGatewayServer(c, apiServiceHandler, authInterceptor, cacheInterceptor)
	server := &gateway.Server{
		APIServer: apiServer,
		DB:        db,
	}
	return server, nil
}

func NewGiteaIntegration(c Config) (component, error) {
	giteaintegrationConfig := provideGiteaIntegrationConfig(c)
	componentsConfig := c.Components
	giteaIntegrationConfig := componentsConfig.GiteaIntegration
	controllerServiceClientConfig := giteaIntegrationConfig.Controller
	controllerGiteaIntegrationServiceClient := grpc.NewControllerGiteaIntegrationServiceClient(controllerServiceClientConfig)
	repositoryConfig := c.DB
	db, err := repository.New(repositoryConfig)
	if err != nil {
		return nil, err
	}
	gitRepositoryRepository := repository.NewGitRepositoryRepository(db)
	applicationRepository := repository.NewApplicationRepository(db)
	userRepository := repository.NewUserRepository(db)
	integration, err := giteaintegration.NewIntegration(giteaintegrationConfig, controllerGiteaIntegrationServiceClient, gitRepositoryRepository, applicationRepository, userRepository)
	if err != nil {
		return nil, err
	}
	server := &giteaintegration2.Server{
		Integration: integration,
		DB:          db,
	}
	return server, nil
}

func NewSSGen(c Config) (component, error) {
	repositoryConfig := c.DB
	db, err := repository.New(repositoryConfig)
	if err != nil {
		return nil, err
	}
	componentsConfig := c.Components
	ssGenConfig := componentsConfig.SSGen
	controllerServiceClientConfig := ssGenConfig.Controller
	controllerSSGenServiceClient := grpc.NewControllerSSGenServiceClient(controllerServiceClientConfig)
	applicationRepository := repository.NewApplicationRepository(db)
	buildRepository := repository.NewBuildRepository(db)
	storageConfig := c.Storage
	storage, err := provideStorage(storageConfig)
	if err != nil {
		return nil, err
	}
	staticServer, err := provideStaticServer(c)
	if err != nil {
		return nil, err
	}
	staticServerDocumentRootPath := provideStaticServerDocumentRootPath(c)
	generatorService := ssgen.NewGeneratorService(controllerSSGenServiceClient, applicationRepository, buildRepository, storage, staticServer, staticServerDocumentRootPath)
	port := ssGenConfig.HealthPort
	healthcheckFunc := provideHealthCheckFunc(generatorService)
	server := healthcheck.NewServer(port, healthcheckFunc)
	ssgenServer := &ssgen2.Server{
		DB:      db,
		Service: generatorService,
		Health:  server,
		Engine:  staticServer,
	}
	return ssgenServer, nil
}

// wire.go:

func NewController(c Config) (component, error) {
	switch c.Components.Controller.Mode {
	case "docker":
		return NewControllerDocker(c)
	case "k8s", "kubernetes":
		return NewControllerK8s(c)
	}
	return nil, errors.New("unknown mode: " + c.Components.Controller.Mode)
}
