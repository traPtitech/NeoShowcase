package apiserver

import (
	"context"
	"testing"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
	"github.com/traPtitech/neoshowcase/pkg/test/mocks"
	"github.com/traPtitech/neoshowcase/pkg/util/testhelper"
)

func DefaultOption(t *testing.T) testhelper.ContainerOption {
	return func(c *testhelper.Container) {
		db := testhelper.OpenDB(t)
		c.Provide(wrapValue(db))

		c.Provide(repository.NewArtifactRepository)
		c.Provide(repository.NewRuntimeImageRepository)
		c.Provide(repository.NewApplicationRepository)
		c.Provide(repository.NewBuildRepository)
		c.Provide(repository.NewEnvironmentRepository)
		c.Provide(repository.NewGitRepositoryRepository)
		c.Provide(repository.NewRepositoryCommitRepository)
		c.Provide(repository.NewUserRepository)

		c.Provide(func() domain.ControllerServiceClient {
			return &mocks.ControllerServiceClientMock{
				GetSystemInfoFunc: DefaultGetSystemInfoFunc,
			}
		})
		c.Provide(func() domain.GitService { return &mocks.GitServiceMock{} })
		c.Provide(func() builder.RegistryClient { return &mocks.RegistryClientMock{} })
		c.Provide(func() domain.MariaDBManager { return &mocks.MariaDBManagerMock{} })
		c.Provide(func() domain.MongoDBManager { return &mocks.MongoDBManagerMock{} })

		dir := t.TempDir()
		c.Provide(func() (domain.Storage, error) {
			return storage.NewLocalStorage(dir)
		})

		c.Provide(wrapValue[domain.MetricsService](nil))
		c.Provide(wrapValue[domain.ContainerLogger](nil))

		c.Provide(wrapValue(builder.ImageConfig{
			Registry: builder.RegistryConfig{
				Scheme:   "https",
				Addr:     "registry.example.com",
				Username: "user",
				Password: "password",
			},
			NamePrefix:    "ns-apps/",
			TmpNamePrefix: "ns-apps-tmp/",
		}))

		c.Provide(NewService)
	}
}

func WithGitMock(gitMock *mocks.GitServiceMock) testhelper.ContainerOption {
	return func(c *testhelper.Container) {
		c.Decorate(func(_ domain.GitService) domain.GitService { return gitMock })
	}
}

func WithRegistryMock(registryMock *mocks.RegistryClientMock) testhelper.ContainerOption {
	return func(c *testhelper.Container) {
		c.Decorate(func(_ builder.RegistryClient) builder.RegistryClient { return registryMock })
	}
}

func WithMariaDBManagerMock(mariaDBManagerMock *mocks.MariaDBManagerMock) testhelper.ContainerOption {
	return func(c *testhelper.Container) {
		c.Decorate(func(_ domain.MariaDBManager) domain.MariaDBManager { return mariaDBManagerMock })
	}
}

func CreateUser(c *testhelper.Container, name string) *domain.User {
	var user *domain.User
	if err := c.Invoke(func(repo domain.UserRepository) (err error) {
		user, err = repo.EnsureUser(context.Background(), name)
		return
	}); err != nil {
		panic(err)
	}
	return user
}

var DefaultGetSystemInfoFunc = func(_ context.Context) (*domain.SystemInfo, error) {
	return &domain.SystemInfo{
		AvailableDomains: domain.AvailableDomainSlice{{Domain: "*.example.com", AuthAvailable: true}},
	}, nil
}

func wrapValue[T any](v T) func() T {
	return func() T {
		return v
	}
}
