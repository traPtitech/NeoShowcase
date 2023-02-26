package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
)

var (
	ErrNotFound = errors.New("not found")
)

func handleRepoError[T any](entity T, err error) (T, error) {
	switch err {
	case repository.ErrNotFound:
		return entity, ErrNotFound
	default:
		return entity, err
	}
}

type CreateApplicationArgs struct {
	UserID        string
	RepositoryURL string
	BranchName    string
	BuildType     builder.BuildType
}

type APIServerService interface {
	GetApplicationsByUserID(ctx context.Context, userID string) ([]*domain.Application, error)
	CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error)
	GetApplication(ctx context.Context, id string) (*domain.Application, error)
	DeleteApplication(ctx context.Context, id string) error
	GetApplicationBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error)
	GetApplicationBuild(ctx context.Context, buildID string) (*domain.Build, error)
	SetApplicationEnvironmentVariable(ctx context.Context, applicationID string, key string, value string) error
	GetApplicationEnvironmentVariables(ctx context.Context, applicationID string) ([]*domain.Environment, error)
}

type apiServerService struct {
	appRepo   repository.ApplicationRepository
	buildRepo repository.BuildRepository
	envRepo   repository.EnvironmentRepository
	gitRepo   repository.GitRepositoryRepository
}

func NewAPIServerService(
	appRepo repository.ApplicationRepository,
	buildRepo repository.BuildRepository,
	envRepo repository.EnvironmentRepository,
	gitRepo repository.GitRepositoryRepository,
) APIServerService {
	return &apiServerService{
		appRepo:   appRepo,
		buildRepo: buildRepo,
		envRepo:   envRepo,
		gitRepo:   gitRepo,
	}
}

func (s *apiServerService) GetApplicationsByUserID(ctx context.Context, userID string) ([]*domain.Application, error) {
	return s.appRepo.GetApplicationsByUserID(ctx, userID)
}

func (s *apiServerService) CreateApplication(ctx context.Context, args CreateApplicationArgs) (*domain.Application, error) {
	repo, err := s.gitRepo.GetRepository(ctx, args.RepositoryURL)
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	if err == repository.ErrNotFound {
		repoName, err := domain.ExtractNameFromRepositoryURL(args.RepositoryURL)
		if err != nil {
			return nil, fmt.Errorf("malformed repository url: %w", err)
		}
		repo, err = s.gitRepo.RegisterRepository(ctx, repository.RegisterRepositoryArgs{
			Name: repoName,
			URL:  args.RepositoryURL,
		})
		if err != nil {
			return nil, err
		}
	}

	application, err := s.appRepo.CreateApplication(ctx, repository.CreateApplicationArgs{
		RepositoryID: repo.ID,
		BranchName:   args.BranchName,
		BuildType:    args.BuildType,
	})

	err = s.appRepo.RegisterApplicationOwner(ctx, application.ID, args.UserID)
	if err != nil {
		return nil, err
	}

	return application, nil
}

func (s *apiServerService) GetApplication(ctx context.Context, id string) (*domain.Application, error) {
	application, err := s.appRepo.GetApplicationByID(ctx, id)
	return handleRepoError(application, err)
}

func (s *apiServerService) DeleteApplication(ctx context.Context, id string) error {
	// TODO implement me
	panic("implement me")
}

func (s *apiServerService) GetApplicationBuilds(ctx context.Context, applicationID string) ([]*domain.Build, error) {
	return s.buildRepo.GetBuilds(ctx, applicationID)
}

func (s *apiServerService) GetApplicationBuild(ctx context.Context, buildID string) (*domain.Build, error) {
	build, err := s.buildRepo.GetBuild(ctx, buildID)
	return handleRepoError(build, err)
}

func (s *apiServerService) GetApplicationEnvironmentVariables(ctx context.Context, applicationID string) ([]*domain.Environment, error) {
	return s.envRepo.GetEnv(ctx, applicationID)
}

func (s *apiServerService) SetApplicationEnvironmentVariable(ctx context.Context, applicationID string, key string, value string) error {
	return s.envRepo.SetEnv(ctx, applicationID, key, value)
}
