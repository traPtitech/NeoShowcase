package apiserver_test

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/test/mocks"
	"github.com/traPtitech/neoshowcase/pkg/test/testhelper"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type appDatabases struct {
	MariaDB bool
	MongoDB bool
}

// serviceFixture is a service with a repository to create applications from.
type serviceFixture struct {
	c       *testhelper.Container
	svc     *apiserver.Service
	ctx     context.Context
	repoID  string
	mariaDB *mocks.MariaDBManagerMock
	mongoDB *mocks.MongoDBManagerMock
}

// appFixture is a serviceFixture with one application created through the API.
type appFixture struct {
	*serviceFixture
	app *domain.Application
}

// newAppFixture creates an application through the API and fails the test if
// the creation fails.
func newAppFixture(t *testing.T, appName string, dbs appDatabases) *appFixture {
	t.Helper()
	f, err := newServiceFixture(t).createApp(appName, dbs)
	require.NoError(t, err)
	return f
}

// newServiceFixture prepares a service and a repository to create applications
// from, without creating any application.
func newServiceFixture(t *testing.T) *serviceFixture {
	t.Helper()

	gitMock := &mocks.GitServiceMock{
		ResolveRefsFunc: func(ctx context.Context, repo *domain.Repository) (map[string]string, error) {
			return map[string]string{"main": exampleCommitHash}, nil
		},
	}
	registryMock := &mocks.RegistryClientMock{
		GetTagsFunc: func(ctx context.Context, image string) ([]string, error) {
			return nil, nil
		},
	}
	mariaDBMock := &mocks.MariaDBManagerMock{
		GetHostFunc: func() (string, int) { return "mariadb.example.com", 3306 },
		CreateFunc:  func(ctx context.Context, args domain.CreateArgs) error { return nil },
		DeleteFunc:  func(ctx context.Context, args domain.DeleteArgs) error { return nil },
	}
	mongoDBMock := &mocks.MongoDBManagerMock{
		GetHostFunc: func() (string, int) { return "mongodb.example.com", 27017 },
		CreateFunc:  func(ctx context.Context, args domain.CreateArgs) error { return nil },
		DeleteFunc:  func(ctx context.Context, args domain.DeleteArgs) error { return nil },
	}
	c := testhelper.NewContainer(
		apiserver.DefaultOption(t),
		apiserver.WithGitMock(gitMock),
		apiserver.WithRegistryMock(registryMock),
		apiserver.WithMariaDBManagerMock(mariaDBMock),
		apiserver.WithMongoDBManagerMock(mongoDBMock),
	)

	svc := testhelper.Resolve[*apiserver.Service](c)
	ctx := t.Context()

	user := apiserver.CreateUser(c, "test-user")
	web.SetUser(&ctx, user)

	repo, err := svc.CreateRepository(ctx, "test-repo", "https://example.com/user/repo", optional.From(apiserver.CreateRepositoryAuth{
		Method:   domain.RepositoryAuthMethodBasic,
		Username: "test-user",
		Password: "test-password",
	}))
	require.NoError(t, err)

	return &serviceFixture{c: c, svc: svc, ctx: ctx, repoID: repo.ID, mariaDB: mariaDBMock, mongoDB: mongoDBMock}
}

func (f *serviceFixture) createApp(appName string, dbs appDatabases) (*appFixture, error) {
	app, err := f.svc.CreateApplication(f.ctx, &domain.Application{
		ID:           domain.NewID(),
		Name:         appName,
		RepositoryID: f.repoID,
		RefName:      "main",
		Commit:       domain.EmptyCommit,
		DeployType:   domain.DeployTypeRuntime,
		Config: domain.ApplicationConfig{
			BuildConfig: &domain.BuildConfigRuntimeBuildpack{
				UseMariaDB: dbs.MariaDB,
				UseMongoDB: dbs.MongoDB,
				Context:    ".",
			},
		},
		Websites: []*domain.Website{
			{
				ID:             domain.NewID(),
				FQDN:           appName + ".example.com",
				PathPrefix:     "/",
				HTTPPort:       80,
				Authentication: domain.AuthenticationTypeOff,
			},
		},
		OwnerIDs: []string{web.GetUser(f.ctx).ID},
	})
	if err != nil {
		return nil, err
	}
	return &appFixture{serviceFixture: f, app: app}, nil
}

func (f *appFixture) env(t *testing.T, key string) (*domain.Environment, bool) {
	t.Helper()
	envs, err := f.svc.GetEnvironmentVariables(f.ctx, f.app.ID)
	require.NoError(t, err)
	return lo.Find(envs, func(e *domain.Environment) bool { return e.Key == key })
}
