// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: neoshowcase/protobuf/apiserver.proto

package pbconnect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	pb "github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// ApplicationServiceName is the fully-qualified name of the ApplicationService service.
	ApplicationServiceName = "neoshowcase.protobuf.ApplicationService"
)

// ApplicationServiceClient is a client for the neoshowcase.protobuf.ApplicationService service.
type ApplicationServiceClient interface {
	GetMe(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.User], error)
	GetRepositories(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.GetRepositoriesResponse], error)
	CreateRepository(context.Context, *connect_go.Request[pb.CreateRepositoryRequest]) (*connect_go.Response[pb.Repository], error)
	GetApplications(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.GetApplicationsResponse], error)
	GetSystemPublicKey(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.GetSystemPublicKeyResponse], error)
	GetAvailableDomains(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.AvailableDomains], error)
	AddAvailableDomain(context.Context, *connect_go.Request[pb.AvailableDomain]) (*connect_go.Response[emptypb.Empty], error)
	CreateApplication(context.Context, *connect_go.Request[pb.CreateApplicationRequest]) (*connect_go.Response[pb.Application], error)
	GetApplication(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.Application], error)
	UpdateApplication(context.Context, *connect_go.Request[pb.UpdateApplicationRequest]) (*connect_go.Response[emptypb.Empty], error)
	DeleteApplication(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[emptypb.Empty], error)
	GetBuilds(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.GetBuildsResponse], error)
	GetBuild(context.Context, *connect_go.Request[pb.BuildIdRequest]) (*connect_go.Response[pb.Build], error)
	GetBuildLogStream(context.Context, *connect_go.Request[pb.BuildIdRequest]) (*connect_go.ServerStreamForClient[pb.BuildLog], error)
	GetBuildLog(context.Context, *connect_go.Request[pb.BuildIdRequest]) (*connect_go.Response[pb.BuildLog], error)
	GetBuildArtifact(context.Context, *connect_go.Request[pb.ArtifactIdRequest]) (*connect_go.Response[pb.ArtifactContent], error)
	GetEnvVars(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.ApplicationEnvVars], error)
	SetEnvVar(context.Context, *connect_go.Request[pb.SetApplicationEnvVarRequest]) (*connect_go.Response[emptypb.Empty], error)
	GetApplicationOutput(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.ApplicationOutput], error)
	CancelBuild(context.Context, *connect_go.Request[pb.BuildIdRequest]) (*connect_go.Response[emptypb.Empty], error)
	RetryCommitBuild(context.Context, *connect_go.Request[pb.RetryCommitBuildRequest]) (*connect_go.Response[emptypb.Empty], error)
	StartApplication(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[emptypb.Empty], error)
	StopApplication(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[emptypb.Empty], error)
}

// NewApplicationServiceClient constructs a client for the neoshowcase.protobuf.ApplicationService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewApplicationServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ApplicationServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &applicationServiceClient{
		getMe: connect_go.NewClient[emptypb.Empty, pb.User](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/GetMe",
			opts...,
		),
		getRepositories: connect_go.NewClient[emptypb.Empty, pb.GetRepositoriesResponse](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/GetRepositories",
			opts...,
		),
		createRepository: connect_go.NewClient[pb.CreateRepositoryRequest, pb.Repository](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/CreateRepository",
			opts...,
		),
		getApplications: connect_go.NewClient[emptypb.Empty, pb.GetApplicationsResponse](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/GetApplications",
			opts...,
		),
		getSystemPublicKey: connect_go.NewClient[emptypb.Empty, pb.GetSystemPublicKeyResponse](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/GetSystemPublicKey",
			opts...,
		),
		getAvailableDomains: connect_go.NewClient[emptypb.Empty, pb.AvailableDomains](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/GetAvailableDomains",
			opts...,
		),
		addAvailableDomain: connect_go.NewClient[pb.AvailableDomain, emptypb.Empty](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/AddAvailableDomain",
			opts...,
		),
		createApplication: connect_go.NewClient[pb.CreateApplicationRequest, pb.Application](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/CreateApplication",
			opts...,
		),
		getApplication: connect_go.NewClient[pb.ApplicationIdRequest, pb.Application](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/GetApplication",
			opts...,
		),
		updateApplication: connect_go.NewClient[pb.UpdateApplicationRequest, emptypb.Empty](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/UpdateApplication",
			opts...,
		),
		deleteApplication: connect_go.NewClient[pb.ApplicationIdRequest, emptypb.Empty](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/DeleteApplication",
			opts...,
		),
		getBuilds: connect_go.NewClient[pb.ApplicationIdRequest, pb.GetBuildsResponse](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/GetBuilds",
			opts...,
		),
		getBuild: connect_go.NewClient[pb.BuildIdRequest, pb.Build](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/GetBuild",
			opts...,
		),
		getBuildLogStream: connect_go.NewClient[pb.BuildIdRequest, pb.BuildLog](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/GetBuildLogStream",
			opts...,
		),
		getBuildLog: connect_go.NewClient[pb.BuildIdRequest, pb.BuildLog](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/GetBuildLog",
			opts...,
		),
		getBuildArtifact: connect_go.NewClient[pb.ArtifactIdRequest, pb.ArtifactContent](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/GetBuildArtifact",
			opts...,
		),
		getEnvVars: connect_go.NewClient[pb.ApplicationIdRequest, pb.ApplicationEnvVars](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/GetEnvVars",
			opts...,
		),
		setEnvVar: connect_go.NewClient[pb.SetApplicationEnvVarRequest, emptypb.Empty](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/SetEnvVar",
			opts...,
		),
		getApplicationOutput: connect_go.NewClient[pb.ApplicationIdRequest, pb.ApplicationOutput](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/GetApplicationOutput",
			opts...,
		),
		cancelBuild: connect_go.NewClient[pb.BuildIdRequest, emptypb.Empty](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/CancelBuild",
			opts...,
		),
		retryCommitBuild: connect_go.NewClient[pb.RetryCommitBuildRequest, emptypb.Empty](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/RetryCommitBuild",
			opts...,
		),
		startApplication: connect_go.NewClient[pb.ApplicationIdRequest, emptypb.Empty](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/StartApplication",
			opts...,
		),
		stopApplication: connect_go.NewClient[pb.ApplicationIdRequest, emptypb.Empty](
			httpClient,
			baseURL+"/neoshowcase.protobuf.ApplicationService/StopApplication",
			opts...,
		),
	}
}

// applicationServiceClient implements ApplicationServiceClient.
type applicationServiceClient struct {
	getMe                *connect_go.Client[emptypb.Empty, pb.User]
	getRepositories      *connect_go.Client[emptypb.Empty, pb.GetRepositoriesResponse]
	createRepository     *connect_go.Client[pb.CreateRepositoryRequest, pb.Repository]
	getApplications      *connect_go.Client[emptypb.Empty, pb.GetApplicationsResponse]
	getSystemPublicKey   *connect_go.Client[emptypb.Empty, pb.GetSystemPublicKeyResponse]
	getAvailableDomains  *connect_go.Client[emptypb.Empty, pb.AvailableDomains]
	addAvailableDomain   *connect_go.Client[pb.AvailableDomain, emptypb.Empty]
	createApplication    *connect_go.Client[pb.CreateApplicationRequest, pb.Application]
	getApplication       *connect_go.Client[pb.ApplicationIdRequest, pb.Application]
	updateApplication    *connect_go.Client[pb.UpdateApplicationRequest, emptypb.Empty]
	deleteApplication    *connect_go.Client[pb.ApplicationIdRequest, emptypb.Empty]
	getBuilds            *connect_go.Client[pb.ApplicationIdRequest, pb.GetBuildsResponse]
	getBuild             *connect_go.Client[pb.BuildIdRequest, pb.Build]
	getBuildLogStream    *connect_go.Client[pb.BuildIdRequest, pb.BuildLog]
	getBuildLog          *connect_go.Client[pb.BuildIdRequest, pb.BuildLog]
	getBuildArtifact     *connect_go.Client[pb.ArtifactIdRequest, pb.ArtifactContent]
	getEnvVars           *connect_go.Client[pb.ApplicationIdRequest, pb.ApplicationEnvVars]
	setEnvVar            *connect_go.Client[pb.SetApplicationEnvVarRequest, emptypb.Empty]
	getApplicationOutput *connect_go.Client[pb.ApplicationIdRequest, pb.ApplicationOutput]
	cancelBuild          *connect_go.Client[pb.BuildIdRequest, emptypb.Empty]
	retryCommitBuild     *connect_go.Client[pb.RetryCommitBuildRequest, emptypb.Empty]
	startApplication     *connect_go.Client[pb.ApplicationIdRequest, emptypb.Empty]
	stopApplication      *connect_go.Client[pb.ApplicationIdRequest, emptypb.Empty]
}

// GetMe calls neoshowcase.protobuf.ApplicationService.GetMe.
func (c *applicationServiceClient) GetMe(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.User], error) {
	return c.getMe.CallUnary(ctx, req)
}

// GetRepositories calls neoshowcase.protobuf.ApplicationService.GetRepositories.
func (c *applicationServiceClient) GetRepositories(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.GetRepositoriesResponse], error) {
	return c.getRepositories.CallUnary(ctx, req)
}

// CreateRepository calls neoshowcase.protobuf.ApplicationService.CreateRepository.
func (c *applicationServiceClient) CreateRepository(ctx context.Context, req *connect_go.Request[pb.CreateRepositoryRequest]) (*connect_go.Response[pb.Repository], error) {
	return c.createRepository.CallUnary(ctx, req)
}

// GetApplications calls neoshowcase.protobuf.ApplicationService.GetApplications.
func (c *applicationServiceClient) GetApplications(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.GetApplicationsResponse], error) {
	return c.getApplications.CallUnary(ctx, req)
}

// GetSystemPublicKey calls neoshowcase.protobuf.ApplicationService.GetSystemPublicKey.
func (c *applicationServiceClient) GetSystemPublicKey(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.GetSystemPublicKeyResponse], error) {
	return c.getSystemPublicKey.CallUnary(ctx, req)
}

// GetAvailableDomains calls neoshowcase.protobuf.ApplicationService.GetAvailableDomains.
func (c *applicationServiceClient) GetAvailableDomains(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.AvailableDomains], error) {
	return c.getAvailableDomains.CallUnary(ctx, req)
}

// AddAvailableDomain calls neoshowcase.protobuf.ApplicationService.AddAvailableDomain.
func (c *applicationServiceClient) AddAvailableDomain(ctx context.Context, req *connect_go.Request[pb.AvailableDomain]) (*connect_go.Response[emptypb.Empty], error) {
	return c.addAvailableDomain.CallUnary(ctx, req)
}

// CreateApplication calls neoshowcase.protobuf.ApplicationService.CreateApplication.
func (c *applicationServiceClient) CreateApplication(ctx context.Context, req *connect_go.Request[pb.CreateApplicationRequest]) (*connect_go.Response[pb.Application], error) {
	return c.createApplication.CallUnary(ctx, req)
}

// GetApplication calls neoshowcase.protobuf.ApplicationService.GetApplication.
func (c *applicationServiceClient) GetApplication(ctx context.Context, req *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.Application], error) {
	return c.getApplication.CallUnary(ctx, req)
}

// UpdateApplication calls neoshowcase.protobuf.ApplicationService.UpdateApplication.
func (c *applicationServiceClient) UpdateApplication(ctx context.Context, req *connect_go.Request[pb.UpdateApplicationRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return c.updateApplication.CallUnary(ctx, req)
}

// DeleteApplication calls neoshowcase.protobuf.ApplicationService.DeleteApplication.
func (c *applicationServiceClient) DeleteApplication(ctx context.Context, req *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return c.deleteApplication.CallUnary(ctx, req)
}

// GetBuilds calls neoshowcase.protobuf.ApplicationService.GetBuilds.
func (c *applicationServiceClient) GetBuilds(ctx context.Context, req *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.GetBuildsResponse], error) {
	return c.getBuilds.CallUnary(ctx, req)
}

// GetBuild calls neoshowcase.protobuf.ApplicationService.GetBuild.
func (c *applicationServiceClient) GetBuild(ctx context.Context, req *connect_go.Request[pb.BuildIdRequest]) (*connect_go.Response[pb.Build], error) {
	return c.getBuild.CallUnary(ctx, req)
}

// GetBuildLogStream calls neoshowcase.protobuf.ApplicationService.GetBuildLogStream.
func (c *applicationServiceClient) GetBuildLogStream(ctx context.Context, req *connect_go.Request[pb.BuildIdRequest]) (*connect_go.ServerStreamForClient[pb.BuildLog], error) {
	return c.getBuildLogStream.CallServerStream(ctx, req)
}

// GetBuildLog calls neoshowcase.protobuf.ApplicationService.GetBuildLog.
func (c *applicationServiceClient) GetBuildLog(ctx context.Context, req *connect_go.Request[pb.BuildIdRequest]) (*connect_go.Response[pb.BuildLog], error) {
	return c.getBuildLog.CallUnary(ctx, req)
}

// GetBuildArtifact calls neoshowcase.protobuf.ApplicationService.GetBuildArtifact.
func (c *applicationServiceClient) GetBuildArtifact(ctx context.Context, req *connect_go.Request[pb.ArtifactIdRequest]) (*connect_go.Response[pb.ArtifactContent], error) {
	return c.getBuildArtifact.CallUnary(ctx, req)
}

// GetEnvVars calls neoshowcase.protobuf.ApplicationService.GetEnvVars.
func (c *applicationServiceClient) GetEnvVars(ctx context.Context, req *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.ApplicationEnvVars], error) {
	return c.getEnvVars.CallUnary(ctx, req)
}

// SetEnvVar calls neoshowcase.protobuf.ApplicationService.SetEnvVar.
func (c *applicationServiceClient) SetEnvVar(ctx context.Context, req *connect_go.Request[pb.SetApplicationEnvVarRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return c.setEnvVar.CallUnary(ctx, req)
}

// GetApplicationOutput calls neoshowcase.protobuf.ApplicationService.GetApplicationOutput.
func (c *applicationServiceClient) GetApplicationOutput(ctx context.Context, req *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.ApplicationOutput], error) {
	return c.getApplicationOutput.CallUnary(ctx, req)
}

// CancelBuild calls neoshowcase.protobuf.ApplicationService.CancelBuild.
func (c *applicationServiceClient) CancelBuild(ctx context.Context, req *connect_go.Request[pb.BuildIdRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return c.cancelBuild.CallUnary(ctx, req)
}

// RetryCommitBuild calls neoshowcase.protobuf.ApplicationService.RetryCommitBuild.
func (c *applicationServiceClient) RetryCommitBuild(ctx context.Context, req *connect_go.Request[pb.RetryCommitBuildRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return c.retryCommitBuild.CallUnary(ctx, req)
}

// StartApplication calls neoshowcase.protobuf.ApplicationService.StartApplication.
func (c *applicationServiceClient) StartApplication(ctx context.Context, req *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return c.startApplication.CallUnary(ctx, req)
}

// StopApplication calls neoshowcase.protobuf.ApplicationService.StopApplication.
func (c *applicationServiceClient) StopApplication(ctx context.Context, req *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return c.stopApplication.CallUnary(ctx, req)
}

// ApplicationServiceHandler is an implementation of the neoshowcase.protobuf.ApplicationService
// service.
type ApplicationServiceHandler interface {
	GetMe(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.User], error)
	GetRepositories(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.GetRepositoriesResponse], error)
	CreateRepository(context.Context, *connect_go.Request[pb.CreateRepositoryRequest]) (*connect_go.Response[pb.Repository], error)
	GetApplications(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.GetApplicationsResponse], error)
	GetSystemPublicKey(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.GetSystemPublicKeyResponse], error)
	GetAvailableDomains(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.AvailableDomains], error)
	AddAvailableDomain(context.Context, *connect_go.Request[pb.AvailableDomain]) (*connect_go.Response[emptypb.Empty], error)
	CreateApplication(context.Context, *connect_go.Request[pb.CreateApplicationRequest]) (*connect_go.Response[pb.Application], error)
	GetApplication(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.Application], error)
	UpdateApplication(context.Context, *connect_go.Request[pb.UpdateApplicationRequest]) (*connect_go.Response[emptypb.Empty], error)
	DeleteApplication(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[emptypb.Empty], error)
	GetBuilds(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.GetBuildsResponse], error)
	GetBuild(context.Context, *connect_go.Request[pb.BuildIdRequest]) (*connect_go.Response[pb.Build], error)
	GetBuildLogStream(context.Context, *connect_go.Request[pb.BuildIdRequest], *connect_go.ServerStream[pb.BuildLog]) error
	GetBuildLog(context.Context, *connect_go.Request[pb.BuildIdRequest]) (*connect_go.Response[pb.BuildLog], error)
	GetBuildArtifact(context.Context, *connect_go.Request[pb.ArtifactIdRequest]) (*connect_go.Response[pb.ArtifactContent], error)
	GetEnvVars(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.ApplicationEnvVars], error)
	SetEnvVar(context.Context, *connect_go.Request[pb.SetApplicationEnvVarRequest]) (*connect_go.Response[emptypb.Empty], error)
	GetApplicationOutput(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.ApplicationOutput], error)
	CancelBuild(context.Context, *connect_go.Request[pb.BuildIdRequest]) (*connect_go.Response[emptypb.Empty], error)
	RetryCommitBuild(context.Context, *connect_go.Request[pb.RetryCommitBuildRequest]) (*connect_go.Response[emptypb.Empty], error)
	StartApplication(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[emptypb.Empty], error)
	StopApplication(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[emptypb.Empty], error)
}

// NewApplicationServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewApplicationServiceHandler(svc ApplicationServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/neoshowcase.protobuf.ApplicationService/GetMe", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/GetMe",
		svc.GetMe,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/GetRepositories", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/GetRepositories",
		svc.GetRepositories,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/CreateRepository", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/CreateRepository",
		svc.CreateRepository,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/GetApplications", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/GetApplications",
		svc.GetApplications,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/GetSystemPublicKey", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/GetSystemPublicKey",
		svc.GetSystemPublicKey,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/GetAvailableDomains", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/GetAvailableDomains",
		svc.GetAvailableDomains,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/AddAvailableDomain", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/AddAvailableDomain",
		svc.AddAvailableDomain,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/CreateApplication", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/CreateApplication",
		svc.CreateApplication,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/GetApplication", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/GetApplication",
		svc.GetApplication,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/UpdateApplication", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/UpdateApplication",
		svc.UpdateApplication,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/DeleteApplication", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/DeleteApplication",
		svc.DeleteApplication,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/GetBuilds", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/GetBuilds",
		svc.GetBuilds,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/GetBuild", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/GetBuild",
		svc.GetBuild,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/GetBuildLogStream", connect_go.NewServerStreamHandler(
		"/neoshowcase.protobuf.ApplicationService/GetBuildLogStream",
		svc.GetBuildLogStream,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/GetBuildLog", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/GetBuildLog",
		svc.GetBuildLog,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/GetBuildArtifact", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/GetBuildArtifact",
		svc.GetBuildArtifact,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/GetEnvVars", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/GetEnvVars",
		svc.GetEnvVars,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/SetEnvVar", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/SetEnvVar",
		svc.SetEnvVar,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/GetApplicationOutput", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/GetApplicationOutput",
		svc.GetApplicationOutput,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/CancelBuild", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/CancelBuild",
		svc.CancelBuild,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/RetryCommitBuild", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/RetryCommitBuild",
		svc.RetryCommitBuild,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/StartApplication", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/StartApplication",
		svc.StartApplication,
		opts...,
	))
	mux.Handle("/neoshowcase.protobuf.ApplicationService/StopApplication", connect_go.NewUnaryHandler(
		"/neoshowcase.protobuf.ApplicationService/StopApplication",
		svc.StopApplication,
		opts...,
	))
	return "/neoshowcase.protobuf.ApplicationService/", mux
}

// UnimplementedApplicationServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedApplicationServiceHandler struct{}

func (UnimplementedApplicationServiceHandler) GetMe(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.User], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.GetMe is not implemented"))
}

func (UnimplementedApplicationServiceHandler) GetRepositories(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.GetRepositoriesResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.GetRepositories is not implemented"))
}

func (UnimplementedApplicationServiceHandler) CreateRepository(context.Context, *connect_go.Request[pb.CreateRepositoryRequest]) (*connect_go.Response[pb.Repository], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.CreateRepository is not implemented"))
}

func (UnimplementedApplicationServiceHandler) GetApplications(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.GetApplicationsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.GetApplications is not implemented"))
}

func (UnimplementedApplicationServiceHandler) GetSystemPublicKey(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.GetSystemPublicKeyResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.GetSystemPublicKey is not implemented"))
}

func (UnimplementedApplicationServiceHandler) GetAvailableDomains(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[pb.AvailableDomains], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.GetAvailableDomains is not implemented"))
}

func (UnimplementedApplicationServiceHandler) AddAvailableDomain(context.Context, *connect_go.Request[pb.AvailableDomain]) (*connect_go.Response[emptypb.Empty], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.AddAvailableDomain is not implemented"))
}

func (UnimplementedApplicationServiceHandler) CreateApplication(context.Context, *connect_go.Request[pb.CreateApplicationRequest]) (*connect_go.Response[pb.Application], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.CreateApplication is not implemented"))
}

func (UnimplementedApplicationServiceHandler) GetApplication(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.Application], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.GetApplication is not implemented"))
}

func (UnimplementedApplicationServiceHandler) UpdateApplication(context.Context, *connect_go.Request[pb.UpdateApplicationRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.UpdateApplication is not implemented"))
}

func (UnimplementedApplicationServiceHandler) DeleteApplication(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.DeleteApplication is not implemented"))
}

func (UnimplementedApplicationServiceHandler) GetBuilds(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.GetBuildsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.GetBuilds is not implemented"))
}

func (UnimplementedApplicationServiceHandler) GetBuild(context.Context, *connect_go.Request[pb.BuildIdRequest]) (*connect_go.Response[pb.Build], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.GetBuild is not implemented"))
}

func (UnimplementedApplicationServiceHandler) GetBuildLogStream(context.Context, *connect_go.Request[pb.BuildIdRequest], *connect_go.ServerStream[pb.BuildLog]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.GetBuildLogStream is not implemented"))
}

func (UnimplementedApplicationServiceHandler) GetBuildLog(context.Context, *connect_go.Request[pb.BuildIdRequest]) (*connect_go.Response[pb.BuildLog], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.GetBuildLog is not implemented"))
}

func (UnimplementedApplicationServiceHandler) GetBuildArtifact(context.Context, *connect_go.Request[pb.ArtifactIdRequest]) (*connect_go.Response[pb.ArtifactContent], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.GetBuildArtifact is not implemented"))
}

func (UnimplementedApplicationServiceHandler) GetEnvVars(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.ApplicationEnvVars], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.GetEnvVars is not implemented"))
}

func (UnimplementedApplicationServiceHandler) SetEnvVar(context.Context, *connect_go.Request[pb.SetApplicationEnvVarRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.SetEnvVar is not implemented"))
}

func (UnimplementedApplicationServiceHandler) GetApplicationOutput(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[pb.ApplicationOutput], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.GetApplicationOutput is not implemented"))
}

func (UnimplementedApplicationServiceHandler) CancelBuild(context.Context, *connect_go.Request[pb.BuildIdRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.CancelBuild is not implemented"))
}

func (UnimplementedApplicationServiceHandler) RetryCommitBuild(context.Context, *connect_go.Request[pb.RetryCommitBuildRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.RetryCommitBuild is not implemented"))
}

func (UnimplementedApplicationServiceHandler) StartApplication(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.StartApplication is not implemented"))
}

func (UnimplementedApplicationServiceHandler) StopApplication(context.Context, *connect_go.Request[pb.ApplicationIdRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("neoshowcase.protobuf.ApplicationService.StopApplication is not implemented"))
}
