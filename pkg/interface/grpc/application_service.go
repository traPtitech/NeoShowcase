package grpc

import (
	"context"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pbconvert"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/samber/lo"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func handleUseCaseError(err error) error {
	underlying, typ, ok := usecase.DecomposeError(err)
	if ok {
		switch typ {
		case usecase.ErrorTypeBadRequest:
			return connect.NewError(connect.CodeInvalidArgument, underlying)
		case usecase.ErrorTypeNotFound:
			return connect.NewError(connect.CodeNotFound, underlying)
		case usecase.ErrorTypeAlreadyExists:
			return connect.NewError(connect.CodeAlreadyExists, underlying)
		case usecase.ErrorTypeForbidden:
			return connect.NewError(connect.CodePermissionDenied, underlying)
		}
	}
	return connect.NewError(connect.CodeInternal, err)
}

type ApplicationService struct {
	svc    *usecase.APIServerService
	pubKey *ssh.PublicKeys
}

func NewApplicationServiceServer(
	svc *usecase.APIServerService,
	pubKey *ssh.PublicKeys,
) pbconnect.ApplicationServiceHandler {
	return &ApplicationService{
		svc:    svc,
		pubKey: pubKey,
	}
}

func (s *ApplicationService) GetMe(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.User], error) {
	user := web.GetUser(ctx)
	res := connect.NewResponse(pbconvert.ToPBUser(user))
	return res, nil
}

func (s *ApplicationService) GetRepositories(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[pb.GetRepositoriesResponse], error) {
	repositories, err := s.svc.GetRepositories(ctx)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.GetRepositoriesResponse{
		Repositories: lo.Map(repositories, func(repo *domain.Repository, i int) *pb.Repository {
			return pbconvert.ToPBRepository(repo)
		}),
	})
	return res, nil
}

func (s *ApplicationService) CreateRepository(ctx context.Context, req *connect.Request[pb.CreateRepositoryRequest]) (*connect.Response[pb.Repository], error) {
	msg := req.Msg
	user := web.GetUser(ctx)
	repo := &domain.Repository{
		ID:       domain.NewID(),
		Name:     msg.Name,
		URL:      msg.Url,
		Auth:     pbconvert.FromPBRepositoryAuth(msg.Auth),
		OwnerIDs: []string{user.ID},
	}
	err := s.svc.CreateRepository(ctx, repo)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(pbconvert.ToPBRepository(repo))
	return res, nil
}

func (s *ApplicationService) UpdateRepository(ctx context.Context, req *connect.Request[pb.UpdateRepositoryRequest]) (*connect.Response[emptypb.Empty], error) {
	msg := req.Msg
	args := &domain.UpdateRepositoryArgs{
		Name:     optional.From(msg.Name),
		URL:      optional.From(msg.Url),
		Auth:     optional.From(pbconvert.FromPBRepositoryAuth(msg.Auth)),
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
			return pbconvert.ToPBApplication(app)
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
		Domains: lo.Map(domains, func(ad *domain.AvailableDomain, i int) *pb.AvailableDomain {
			return pbconvert.ToPBAvailableDomain(ad)
		}),
	})
	return res, nil
}

func (s *ApplicationService) AddAvailableDomain(ctx context.Context, req *connect.Request[pb.AvailableDomain]) (*connect.Response[emptypb.Empty], error) {
	err := s.svc.AddAvailableDomain(ctx, pbconvert.FromPBAvailableDomain(req.Msg))
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

func (s *ApplicationService) CreateApplication(ctx context.Context, req *connect.Request[pb.CreateApplicationRequest]) (*connect.Response[pb.Application], error) {
	msg := req.Msg
	user := web.GetUser(ctx)
	now := time.Now()
	config := pbconvert.FromPBApplicationConfig(msg.Config)
	app := &domain.Application{
		ID:            domain.NewID(),
		Name:          msg.Name,
		RepositoryID:  msg.RepositoryId,
		RefName:       msg.RefName,
		DeployType:    config.BuildConfig.BuildType().DeployType(),
		Running:       msg.StartOnCreate,
		Container:     domain.ContainerStateMissing,
		CurrentCommit: domain.EmptyCommit,
		WantCommit:    domain.EmptyCommit,
		CreatedAt:     now,
		UpdatedAt:     now,
		Config:        config,
		Websites: lo.Map(msg.Websites, func(website *pb.CreateWebsiteRequest, i int) *domain.Website {
			return pbconvert.FromPBCreateWebsiteRequest(website)
		}),
		OwnerIDs: []string{user.ID},
	}
	app, err := s.svc.CreateApplication(ctx, app)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(pbconvert.ToPBApplication(app))
	return res, nil
}

func (s *ApplicationService) GetApplication(ctx context.Context, req *connect.Request[pb.ApplicationIdRequest]) (*connect.Response[pb.Application], error) {
	application, err := s.svc.GetApplication(ctx, req.Msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(pbconvert.ToPBApplication(application))
	return res, nil
}

func (s *ApplicationService) UpdateApplication(ctx context.Context, req *connect.Request[pb.UpdateApplicationRequest]) (*connect.Response[emptypb.Empty], error) {
	msg := req.Msg
	app, err := s.svc.GetApplication(ctx, msg.Id)
	if err != nil {
		return nil, handleUseCaseError(err)
	}

	websites := app.Websites
	for _, createReq := range msg.NewWebsites {
		websites = append(websites, pbconvert.FromPBCreateWebsiteRequest(createReq))
	}
	for _, deleteReq := range msg.DeleteWebsites {
		websites = lo.Reject(websites, func(w *domain.Website, i int) bool { return w.ID == deleteReq.Id })
	}

	err = s.svc.UpdateApplication(ctx, msg.Id, &domain.UpdateApplicationArgs{
		Name:      optional.From(msg.Name),
		RefName:   optional.From(msg.RefName),
		UpdatedAt: optional.From(time.Now()),
		Config: optional.From(domain.ApplicationConfig{
			UseMariaDB:  app.Config.UseMariaDB,
			UseMongoDB:  app.Config.UseMongoDB,
			BuildType:   pbconvert.BuildTypeMapper.FromMust(msg.Config.BuildType),
			BuildConfig: pbconvert.FromPBBuildConfig(msg.Config.BuildConfig),
		}),
		Websites: optional.From(websites),
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
			return pbconvert.ToPBBuild(build)
		}),
	})
	return res, nil
}

func (s *ApplicationService) GetBuild(ctx context.Context, req *connect.Request[pb.BuildIdRequest]) (*connect.Response[pb.Build], error) {
	build, err := s.svc.GetBuild(ctx, req.Msg.BuildId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(pbconvert.ToPBBuild(build))
	return res, nil
}

func (s *ApplicationService) GetBuildLog(ctx context.Context, req *connect.Request[pb.BuildIdRequest]) (*connect.Response[pb.BuildLog], error) {
	log, err := s.svc.GetBuildLog(ctx, req.Msg.BuildId)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.BuildLog{Log: log})
	return res, nil
}

func (s *ApplicationService) GetBuildLogStream(ctx context.Context, req *connect.Request[pb.BuildIdRequest], st *connect.ServerStream[pb.BuildLog]) error {
	err := s.svc.GetBuildLogStream(ctx, req.Msg.BuildId, func(b []byte) error {
		return st.Send(&pb.BuildLog{Log: b})
	})
	if err != nil {
		return handleUseCaseError(err)
	}
	return nil
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
			return pbconvert.ToPBEnvironment(env)
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

func (s *ApplicationService) GetOutput(ctx context.Context, req *connect.Request[pb.GetOutputRequest]) (*connect.Response[pb.GetOutputResponse], error) {
	msg := req.Msg
	before := time.Now()
	if req.Msg.Before != nil {
		before = msg.Before.AsTime()
	}
	logs, err := s.svc.GetOutput(ctx, msg.ApplicationId, before)
	if err != nil {
		return nil, handleUseCaseError(err)
	}
	res := connect.NewResponse(&pb.GetOutputResponse{
		Outputs: lo.Map(logs, func(l *domain.ContainerLog, i int) *pb.ApplicationOutput {
			return pbconvert.ToPBApplicationOutput(l)
		}),
	})
	return res, nil
}

func (s *ApplicationService) GetOutputStream(ctx context.Context, req *connect.Request[pb.GetOutputStreamRequest], st *connect.ServerStream[pb.ApplicationOutput]) error {
	msg := req.Msg
	after := time.Now()
	if req.Msg.After != nil {
		after = msg.After.AsTime()
	}
	err := s.svc.GetOutputStream(ctx, msg.ApplicationId, after, func(l *domain.ContainerLog) error {
		return st.Send(pbconvert.ToPBApplicationOutput(l))
	})
	if err != nil {
		return handleUseCaseError(err)
	}
	return nil
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
