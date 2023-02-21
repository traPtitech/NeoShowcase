package usecase

import (
	"context"
	"fmt"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
)

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
}

type apiServerService struct {
	appRepo repository.ApplicationRepository
	gitRepo repository.GitRepositoryRepository
}

func NewAPIServerService(appRepo repository.ApplicationRepository, gitRepo repository.GitRepositoryRepository) APIServerService {
	return &apiServerService{
		appRepo: appRepo,
		gitRepo: gitRepo,
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
	return s.appRepo.GetApplicationByID(ctx, id)
}

func (s *apiServerService) DeleteApplication(ctx context.Context, id string) error {
	// TODO implement me
	panic("implement me")
}
