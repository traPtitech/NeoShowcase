package grpc

import (
	"context"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func handleUseCaseError(err error) error {
	typ, ok := usecase.GetErrorType(err)
	if ok {
		switch typ {
		case usecase.ErrorTypeBadRequest:
			return connect.NewError(connect.CodeInvalidArgument, err)
		case usecase.ErrorTypeNotFound:
			return connect.NewError(connect.CodeNotFound, err)
		case usecase.ErrorTypeAlreadyExists:
			return connect.NewError(connect.CodeAlreadyExists, err)
		}
	}
	return connect.NewError(connect.CodeInternal, err)
}

type ApplicationService struct {
	svc    *usecase.APIServerService
	logSvc *usecase.LogStreamService
	pubKey *ssh.PublicKeys
}

func NewApplicationServiceServer(
	svc *usecase.APIServerService,
	logSvc *usecase.LogStreamService,
	pubKey *ssh.PublicKeys,
) pbconnect.ApplicationServiceHandler {
	return &ApplicationService{
		svc:    svc,
		logSvc: logSvc,
		pubKey: pubKey,
	}
}

func (s *ApplicationService) GetMe(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.User], error) {
	user := getUser(ctx)
	res := connect.NewResponse(toPBUser(user))
	return res, nil
}

func (s *ApplicationService) GetRepositories(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.GetRepositoriesResponse], error) {
	repositories, err := s.svc.GetRepositories(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.GetRepositoriesResponse{
		Repositories: lo.Map(repositories, func(repo *domain.Repository, i int) *pb.Repository {
			return toPBRepository(repo)
		}),
	})
	return res, nil
}

func (s *ApplicationService) CreateRepository(ctx context.Context, req *connect.Request[pb.CreateRepositoryRequest]) (*connect.Response[pb.Repository], error) {
	msg := req.Msg
	user := getUser(ctx)
	repo := &domain.Repository{
		ID:       domain.NewID(),
		Name:     msg.Name,
		URL:      msg.Url,
		Auth:     fromPBRepositoryAuth(msg.Auth),
		OwnerIDs: []string{user.ID},
	}
	err := s.svc.CreateRepository(ctx, repo)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(toPBRepository(repo))
	return res, nil
}

func (s *ApplicationService) UpdateRepository(ctx context.Context, req *connect.Request[pb.UpdateRepositoryRequest]) (*connect.Response[emptypb.Empty], error) {
	msg := req.Msg
	args := &domain.UpdateRepositoryArgs{
		Name:     optional.From(msg.Name),
		URL:      optional.From(msg.Url),
		Auth:     optional.From(fromPBRepositoryAuth(msg.Auth)),
		OwnerIDs: optional.From(msg.OwnerIds),
	}
	err := s.svc.UpdateRepository(ctx, msg.Id, args)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ApplicationService) DeleteRepository(ctx context.Context, req *connect.Request[pb.RepositoryIdRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.DeleteRepository(ctx, req.Msg.RepositoryId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ApplicationService) GetApplications(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.GetApplicationsResponse], error) {
	applications, err := s.svc.GetApplications(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.GetApplicationsResponse{
		Applications: lo.Map(applications, func(app *domain.Application, i int) *pb.Application {
			return toPBApplication(app)
		}),
	})
	return res, nil
}

func (s *ApplicationService) GetSystemPublicKey(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.GetSystemPublicKeyResponse], error) {
	encoded := domain.Base64EncodedPublicKey(s.pubKey)
	res := connect.NewResponse(&pb.GetSystemPublicKeyResponse{
		PublicKey: encoded,
	})
	return res, nil
}

func (s *ApplicationService) GetAvailableDomains(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.AvailableDomains], error) {
	domains, err := s.svc.GetAvailableDomains(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	res := connect.NewResponse(&pb.AvailableDomains{
		Domains: lo.Map(domains, func(ad *domain.AvailableDomain, i int) *pb.AvailableDomain { return toPBAvailableDomain(ad) }),
	})
	return res, nil
}

func (s *ApplicationService) AddAvailableDomain(ctx context.Context, req *connect.Request[pb.AvailableDomain]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.AddAvailableDomain(ctx, fromPBAvailableDomain(req.Msg))
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ApplicationService) CreateApplication(ctx context.Context, req *connect.Request[pb.CreateApplicationRequest]) (*connect.Response[pb.Application], error) {
	msg := req.Msg
	user := getUser(ctx)
	now := time.Now()
	app := &domain.Application{
		ID:            domain.NewID(),
		Name:          msg.Name,
		RepositoryID:  msg.RepositoryId,
		RefName:       msg.RefName,
		BuildType:     buildTypeMapper.FromMust(msg.BuildType),
		Running:       msg.StartOnCreate,
		Container:     domain.ContainerStateMissing,
		CurrentCommit: domain.EmptyCommit,
		WantCommit:    domain.EmptyCommit,
		CreatedAt:     now,
		UpdatedAt:     now,
		Config:        fromPBApplicationConfig(msg.Config),
		Websites: lo.Map(msg.Websites, func(website *pb.CreateWebsiteRequest, i int) *domain.Website {
			return &domain.Website{
				ID:          domain.NewID(),
				FQDN:        website.Fqdn,
				PathPrefix:  website.PathPrefix,
				StripPrefix: website.StripPrefix,
				HTTPS:       website.Https,
				HTTPPort:    int(website.HttpPort),
			}
		}),
		OwnerIDs: []string{user.ID},
	}
	app, err := s.svc.CreateApplication(ctx, app)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(toPBApplication(app))
	return res, nil
}

func (s *ApplicationService) GetApplication(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[pb.Application], error) {
	application, err := s.svc.GetApplication(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(toPBApplication(application))
	return res, nil
}

func (s *ApplicationService) UpdateApplication(ctx context.Context, req *connect.Request[pb.UpdateApplicationRequest]) (*connect.Response[emptypb.Empty], error) {
	msg := req.Msg
	app, err := s.svc.GetApplication(ctx, msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}

	for _, createReq := range msg.NewWebsites {
		app.Websites = append(app.Websites, fromPBCreateWebsiteRequest(createReq))
	}
	for _, deleteReq := range msg.DeleteWebsites {
		app.Websites = lo.Reject(app.Websites, func(w *domain.Website, i int) bool { return w.ID == deleteReq.Id })
	}

	err = s.svc.UpdateApplication(ctx, app, &domain.UpdateApplicationArgs{
		Name:      optional.From(msg.Name),
		RefName:   optional.From(msg.RefName),
		UpdatedAt: optional.From(time.Now()),
		Config: optional.From(domain.ApplicationConfig{
			UseMariaDB:     app.Config.UseMariaDB,
			UseMongoDB:     app.Config.UseMongoDB,
			BaseImage:      msg.Config.BaseImage,
			DockerfileName: msg.Config.DockerfileName,
			ArtifactPath:   msg.Config.ArtifactPath,
			BuildCmd:       msg.Config.BuildCmd,
			EntrypointCmd:  msg.Config.EntrypointCmd,
			Authentication: authTypeMapper.FromMust(msg.Config.Authentication),
		}),
		Websites: optional.From(app.Websites),
		OwnerIDs: optional.From(msg.OwnerIds),
	})
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ApplicationService) DeleteApplication(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.DeleteApplication(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ApplicationService) GetBuilds(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[pb.GetBuildsResponse], error) {
	builds, err := s.svc.GetBuilds(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.GetBuildsResponse{
		Builds: lo.Map(builds, func(build *domain.Build, i int) *pb.Build {
			return toPBBuild(build)
		}),
	})
	return res, nil
}

func (s *ApplicationService) GetBuild(ctx context.Context, req *connect.Request[pb.BuildIdRequest]) (*connect.Response[pb.Build], error) {
	build, err := s.svc.GetBuild(ctx, req.Msg.BuildId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(toPBBuild(build))
	return res, nil
}

func (s *ApplicationService) GetBuildLogStream(_ context.Context, req *connect.Request[pb.BuildIdRequest], st *connect.ServerStream[pb.BuildLog]) error {
	sub := make(chan []byte, 100)
	ok, unsubscribe := s.logSvc.SubscribeBuildLog(req.Msg.BuildId, sub)
	if !ok {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("build log stream not available"))
	}
	defer unsubscribe()

	for log := range sub {
		err := st.Send(&pb.BuildLog{Log: log})
		if err != nil {
			return connect.NewError(connect.CodeInternal, errors.New("failed to send log"))
		}
	}
	return nil
}

func (s *ApplicationService) GetBuildLog(ctx context.Context, req *connect.Request[pb.BuildIdRequest]) (*connect.Response[pb.BuildLog], error) {
	log, err := s.svc.GetBuildLog(ctx, req.Msg.BuildId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.BuildLog{Log: log})
	return res, nil
}

func (s *ApplicationService) GetBuildArtifact(ctx context.Context, req *connect.Request[pb.ArtifactIdRequest]) (*connect.Response[pb.ArtifactContent], error) {
	content, err := s.svc.GetArtifact(ctx, req.Msg.ArtifactId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.ArtifactContent{
		Filename: req.Msg.ArtifactId + ".tar",
		Content:  content,
	})
	return res, nil
}

func (s *ApplicationService) GetEnvVars(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[pb.ApplicationEnvVars], error) {
	environments, err := s.svc.GetEnvironmentVariables(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.ApplicationEnvVars{
		Variables: lo.Map(environments, func(env *domain.Environment, i int) *pb.ApplicationEnvVar {
			return toPBEnvironment(env)
		}),
	})
	return res, nil
}

func (s *ApplicationService) SetEnvVar(ctx context.Context, req *connect.Request[pb.SetApplicationEnvVarRequest]) (*connect.Response[emptypb.Empty], error) {
	msg := req.Msg
	err := s.svc.SetEnvironmentVariable(ctx, msg.ApplicationId, msg.Key, msg.Value)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ApplicationService) GetApplicationOutput(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[pb.ApplicationOutput], error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetApplicationOutput not implemented")
}

func (s *ApplicationService) CancelBuild(ctx context.Context, req *connect.Request[pb.BuildIdRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.CancelBuild(ctx, req.Msg.BuildId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ApplicationService) RetryCommitBuild(ctx context.Context, req *connect.Request[pb.RetryCommitBuildRequest]) (*connect.Response[emptypb.Empty], error) {
	msg := req.Msg
	err := s.svc.RetryCommitBuild(ctx, msg.ApplicationId, msg.Commit)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ApplicationService) StartApplication(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.StartApplication(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ApplicationService) StopApplication(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.StopApplication(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func fromPBAvailableDomain(ad *pb.AvailableDomain) *domain.AvailableDomain {
	return &domain.AvailableDomain{
		Domain:    ad.Domain,
		Available: ad.Available,
	}
}

func toPBAvailableDomain(ad *domain.AvailableDomain) *pb.AvailableDomain {
	return &pb.AvailableDomain{
		Domain:    ad.Domain,
		Available: ad.Available,
	}
}

func fromPBRepositoryAuth(req *pb.CreateRepositoryAuth) optional.Of[domain.RepositoryAuth] {
	switch v := req.Auth.(type) {
	case *pb.CreateRepositoryAuth_None:
		return optional.Of[domain.RepositoryAuth]{}
	case *pb.CreateRepositoryAuth_Basic:
		return optional.From(domain.RepositoryAuth{
			Method:   domain.RepositoryAuthMethodBasic,
			Username: v.Basic.Username,
			Password: v.Basic.Password,
		})
	case *pb.CreateRepositoryAuth_Ssh:
		return optional.From(domain.RepositoryAuth{
			Method: domain.RepositoryAuthMethodSSH,
			SSHKey: v.Ssh.SshKey,
		})
	default:
		panic("unknown auth type")
	}
}

var repoAuthMethodMapper = mapper.NewValueMapper(map[domain.RepositoryAuthMethod]pb.Repository_AuthMethod{
	domain.RepositoryAuthMethodBasic: pb.Repository_BASIC,
	domain.RepositoryAuthMethodSSH:   pb.Repository_SSH,
})

func toPBRepository(repo *domain.Repository) *pb.Repository {
	ret := &pb.Repository{
		Id:   repo.ID,
		Name: repo.Name,
		Url:  repo.URL,
	}
	if repo.Auth.Valid {
		ret.AuthMethod = repoAuthMethodMapper.IntoMust(repo.Auth.V.Method)
	}
	return ret
}

var buildTypeMapper = mapper.NewValueMapper(map[domain.BuildType]pb.BuildType{
	domain.BuildTypeRuntime: pb.BuildType_RUNTIME,
	domain.BuildTypeStatic:  pb.BuildType_STATIC,
})

var authTypeMapper = mapper.NewValueMapper(map[domain.AuthenticationType]pb.AuthenticationType{
	domain.AuthenticationTypeOff:  pb.AuthenticationType_OFF,
	domain.AuthenticationTypeSoft: pb.AuthenticationType_SOFT,
	domain.AuthenticationTypeHard: pb.AuthenticationType_HARD,
})

func toPBApplicationConfig(c domain.ApplicationConfig) *pb.ApplicationConfig {
	return &pb.ApplicationConfig{
		UseMariadb:     c.UseMariaDB,
		UseMongodb:     c.UseMongoDB,
		BaseImage:      c.BaseImage,
		DockerfileName: c.DockerfileName,
		ArtifactPath:   c.ArtifactPath,
		BuildCmd:       c.BuildCmd,
		EntrypointCmd:  c.EntrypointCmd,
		Authentication: authTypeMapper.IntoMust(c.Authentication),
	}
}

func fromPBApplicationConfig(c *pb.ApplicationConfig) domain.ApplicationConfig {
	return domain.ApplicationConfig{
		UseMariaDB:     c.UseMariadb,
		UseMongoDB:     c.UseMongodb,
		BaseImage:      c.BaseImage,
		DockerfileName: c.DockerfileName,
		ArtifactPath:   c.ArtifactPath,
		BuildCmd:       c.BuildCmd,
		EntrypointCmd:  c.EntrypointCmd,
		Authentication: authTypeMapper.FromMust(c.Authentication),
	}
}

func fromPBCreateWebsiteRequest(req *pb.CreateWebsiteRequest) *domain.Website {
	return &domain.Website{
		ID:          domain.NewID(),
		FQDN:        req.Fqdn,
		PathPrefix:  req.PathPrefix,
		StripPrefix: req.StripPrefix,
		HTTPS:       req.Https,
		HTTPPort:    int(req.HttpPort),
	}
}

func toPBWebsite(website *domain.Website) *pb.Website {
	return &pb.Website{
		Id:         website.ID,
		Fqdn:       website.FQDN,
		PathPrefix: website.PathPrefix,
		Https:      website.HTTPS,
		HttpPort:   int32(website.HTTPPort),
	}
}

var containerStateMapper = mapper.NewValueMapper(map[domain.ContainerState]pb.Application_ContainerState{
	domain.ContainerStateMissing:  pb.Application_MISSING,
	domain.ContainerStateStarting: pb.Application_STARTING,
	domain.ContainerStateRunning:  pb.Application_RUNNING,
	domain.ContainerStateExited:   pb.Application_EXITED,
	domain.ContainerStateErrored:  pb.Application_ERRORED,
	domain.ContainerStateUnknown:  pb.Application_UNKNOWN,
})

func toPBApplication(app *domain.Application) *pb.Application {
	return &pb.Application{
		Id:            app.ID,
		Name:          app.Name,
		RepositoryId:  app.RepositoryID,
		RefName:       app.RefName,
		BuildType:     buildTypeMapper.IntoMust(app.BuildType),
		Running:       app.Running,
		Container:     containerStateMapper.IntoMust(app.Container),
		CurrentCommit: app.CurrentCommit,
		WantCommit:    app.WantCommit,
		Config:        toPBApplicationConfig(app.Config),
		Websites:      lo.Map(app.Websites, func(website *domain.Website, i int) *pb.Website { return toPBWebsite(website) }),
		OwnerIds:      app.OwnerIDs,
	}
}

var buildStatusMapper = mapper.NewValueMapper(map[domain.BuildStatus]pb.Build_BuildStatus{
	domain.BuildStatusQueued:    pb.Build_QUEUED,
	domain.BuildStatusBuilding:  pb.Build_BUILDING,
	domain.BuildStatusSucceeded: pb.Build_SUCCEEDED,
	domain.BuildStatusFailed:    pb.Build_FAILED,
	domain.BuildStatusCanceled:  pb.Build_CANCELLED,
	domain.BuildStatusSkipped:   pb.Build_SKIPPED,
})

func toPBNullTimestamp(t optional.Of[time.Time]) *pb.NullTimestamp {
	return &pb.NullTimestamp{Timestamp: timestamppb.New(t.V), Valid: t.Valid}
}

func toPBArtifact(artifact *domain.Artifact) *pb.Artifact {
	return &pb.Artifact{
		Id:        artifact.ID,
		Size:      artifact.Size,
		CreatedAt: timestamppb.New(artifact.CreatedAt),
		DeletedAt: toPBNullTimestamp(artifact.DeletedAt),
	}
}

func toPBBuild(build *domain.Build) *pb.Build {
	b := &pb.Build{
		Id:         build.ID,
		Commit:     build.Commit,
		Status:     buildStatusMapper.IntoMust(build.Status),
		StartedAt:  toPBNullTimestamp(build.StartedAt),
		UpdatedAt:  toPBNullTimestamp(build.UpdatedAt),
		FinishedAt: toPBNullTimestamp(build.FinishedAt),
		Retriable:  build.Retriable,
	}
	if build.Artifact.Valid {
		b.Artifact = toPBArtifact(&build.Artifact.V)
	}
	return b
}

func toPBEnvironment(env *domain.Environment) *pb.ApplicationEnvVar {
	return &pb.ApplicationEnvVar{
		Key:    env.Key,
		Value:  env.Value,
		System: env.System,
	}
}

func toPBUser(user *domain.User) *pb.User {
	return &pb.User{
		Id:    user.ID,
		Name:  user.Name,
		Admin: user.Admin,
	}
}
