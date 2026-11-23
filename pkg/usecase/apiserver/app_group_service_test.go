package apiserver_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/test/mocks"
	"github.com/traPtitech/neoshowcase/pkg/test/testhelper"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type groupTestEnv struct {
	c    *testhelper.Container
	svc  *apiserver.Service
	repo *domain.Repository
}

func newGroupTestEnv(t *testing.T) *groupTestEnv {
	t.Helper()
	gitMock := &mocks.GitServiceMock{
		ResolveRefsFunc: func(ctx context.Context, repo *domain.Repository) (map[string]string, error) {
			return map[string]string{"main": exampleCommitHash}, nil
		},
	}
	c := testhelper.NewContainer(
		apiserver.DefaultOption(t),
		apiserver.WithGitMock(gitMock),
	)
	svc := testhelper.Resolve[*apiserver.Service](c)
	env := &groupTestEnv{c: c, svc: svc}

	// A public repository lets any user create applications from it.
	ctx := env.as(t, "repo-owner")
	repo, err := svc.CreateRepository(ctx, "test-repo", "https://example.com/user/repo", optional.None[apiserver.CreateRepositoryAuth]())
	require.NoError(t, err)
	env.repo = repo
	return env
}

// as returns a context of the user with the name, creating the user if needed.
func (e *groupTestEnv) as(t *testing.T, name string) context.Context {
	t.Helper()
	ctx := t.Context()
	web.SetUser(&ctx, apiserver.CreateUser(e.c, name))
	return ctx
}

func (e *groupTestEnv) newApp(name, fqdn, pathPrefix string) *domain.Application {
	return &domain.Application{
		ID:           domain.NewID(),
		Name:         name,
		RepositoryID: e.repo.ID,
		RefName:      "main",
		Commit:       domain.EmptyCommit,
		DeployType:   domain.DeployTypeRuntime,
		Config: domain.ApplicationConfig{
			BuildConfig: &domain.BuildConfigRuntimeBuildpack{Context: "."},
		},
		Websites: []*domain.Website{{
			ID:             domain.NewID(),
			FQDN:           fqdn,
			PathPrefix:     pathPrefix,
			HTTPPort:       80,
			Authentication: domain.AuthenticationTypeOff,
		}},
	}
}

func TestCreateApplication_Standalone(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")

	created, err := e.svc.CreateApplication(alice, e.newApp("app", "standalone.example.com", "/"), optional.None[string]())
	require.NoError(t, err)

	got, err := e.svc.GetApplication(alice, created.App.ID)
	require.NoError(t, err)
	assert.True(t, got.Group.Implicit, "an app created without a group is standalone")
	assert.ElementsMatch(t, []string{web.GetUser(alice).ID, e.repo.OwnerIDs[0]}, got.Group.OwnerIDs, "owners are the creator and the repository owners")

	groups, err := e.svc.GetApplicationGroups(alice, apiserver.GetAppGroupScopeAll)
	require.NoError(t, err)
	assert.Empty(t, groups, "standalone apps are not listed as groups")

	bob := e.as(t, "bob")
	err = e.svc.UpdateApplication(bob, created.App.ID, &domain.UpdateApplicationArgs{Name: optional.From("renamed")})
	assertErrorType(t, domain.ErrorTypeForbidden, err)
}

func assertErrorType(t *testing.T, want domain.ErrorType, err error) {
	t.Helper()
	require.Error(t, err)
	_, typ, _ := domain.DecomposeError(err)
	assert.Equal(t, want, typ, "error: %v", err)
}
