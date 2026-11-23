package apiserver_test

import (
	"context"
	"strings"
	"testing"

	"github.com/samber/lo"
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

// asAdmin returns a context of the admin user with the name, creating the user if needed.
func (e *groupTestEnv) asAdmin(t *testing.T, name string) context.Context {
	t.Helper()
	ctx := t.Context()
	user := apiserver.CreateUser(e.c, name)
	user.Admin = true
	web.SetUser(&ctx, user)
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
	err = e.svc.UpdateApplication(bob, created.App.ID, &domain.UpdateApplicationArgs{Name: optional.From("renamed")}, optional.None[[]string]())
	assertErrorType(t, domain.ErrorTypeForbidden, err)
}

func assertErrorType(t *testing.T, want domain.ErrorType, err error) {
	t.Helper()
	require.Error(t, err)
	_, typ, _ := domain.DecomposeError(err)
	assert.Equal(t, want, typ, "error: %v", err)
}

func TestCreateApplicationGroup_NameLength(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")

	_, err := e.svc.CreateApplicationGroup(alice, strings.Repeat("あ", 101))
	assertErrorType(t, domain.ErrorTypeBadRequest, err)

	group, err := e.svc.CreateApplicationGroup(alice, strings.Repeat("あ", 100))
	require.NoError(t, err)
	assert.Equal(t, strings.Repeat("あ", 100), group.Name)
}

func TestCreateApplication_InGroup(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")
	bob := e.as(t, "bob")

	group, err := e.svc.CreateApplicationGroup(alice, "service")
	require.NoError(t, err)

	_, err = e.svc.CreateApplication(bob, e.newApp("intruder", "in-group.example.com", "/"), optional.From(group.ID))
	assertErrorType(t, domain.ErrorTypeForbidden, err)

	created, err := e.svc.CreateApplication(alice, e.newApp("frontend", "in-group.example.com", "/"), optional.From(group.ID))
	require.NoError(t, err)
	assert.Equal(t, group.ID, created.Group.ID)

	apps, err := e.svc.GetApplications(alice, apiserver.GetAppScope{Scope: apiserver.GetAppScopeGroup, GroupID: optional.From(group.ID)})
	require.NoError(t, err)
	require.Len(t, apps, 1)
	assert.Equal(t, created.App.ID, apps[0].App.ID)

	// Adding an owner to the group lets the user operate the apps in it.
	err = e.svc.UpdateApplicationGroup(alice, group.ID, &domain.UpdateApplicationGroupArgs{
		OwnerIDs: optional.From([]string{web.GetUser(alice).ID, web.GetUser(bob).ID}),
	})
	require.NoError(t, err)
	err = e.svc.UpdateApplication(bob, created.App.ID, &domain.UpdateApplicationArgs{Name: optional.From("renamed")}, optional.None[[]string]())
	require.NoError(t, err)

	mine, err := e.svc.GetApplications(bob, apiserver.GetAppScope{Scope: apiserver.GetAppScopeMine})
	require.NoError(t, err)
	assert.Len(t, mine, 1, "apps of the groups the user owns are the user's apps")
}

func TestCreateApplication_NestedWebsites(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")

	group, err := e.svc.CreateApplicationGroup(alice, "service")
	require.NoError(t, err)
	_, err = e.svc.CreateApplication(alice, e.newApp("frontend", "nested.example.com", "/"), optional.From(group.ID))
	require.NoError(t, err)

	_, err = e.svc.CreateApplication(alice, e.newApp("other", "nested.example.com", "/api"), optional.None[string]())
	assertErrorType(t, domain.ErrorTypeBadRequest, err)

	_, err = e.svc.CreateApplication(alice, e.newApp("same", "nested.example.com", "/"), optional.From(group.ID))
	assertErrorType(t, domain.ErrorTypeBadRequest, err)

	_, err = e.svc.CreateApplication(alice, e.newApp("backend", "nested.example.com", "/api"), optional.From(group.ID))
	require.NoError(t, err)
}

func TestMoveApplication(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")
	bob := e.as(t, "bob")

	aliceGroup, err := e.svc.CreateApplicationGroup(alice, "alice-service")
	require.NoError(t, err)
	bobGroup, err := e.svc.CreateApplicationGroup(bob, "bob-service")
	require.NoError(t, err)
	_, err = e.svc.CreateApplication(alice, e.newApp("frontend", "move.example.com", "/"), optional.From(aliceGroup.ID))
	require.NoError(t, err)
	backend, err := e.svc.CreateApplication(alice, e.newApp("backend", "move-api.example.com", "/"), optional.None[string]())
	require.NoError(t, err)

	t.Run("needs ownership of the destination", func(t *testing.T) {
		err := e.svc.MoveApplication(alice, backend.App.ID, optional.From(bobGroup.ID))
		assertErrorType(t, domain.ErrorTypeForbidden, err)
	})

	t.Run("needs ownership of the source", func(t *testing.T) {
		err := e.svc.MoveApplication(bob, backend.App.ID, optional.From(bobGroup.ID))
		assertErrorType(t, domain.ErrorTypeForbidden, err)
	})

	t.Run("moves a standalone app into a group", func(t *testing.T) {
		err := e.svc.MoveApplication(alice, backend.App.ID, optional.From(aliceGroup.ID))
		require.NoError(t, err)

		got, err := e.svc.GetApplication(alice, backend.App.ID)
		require.NoError(t, err)
		assert.Equal(t, aliceGroup.ID, got.Group.ID)
	})

	t.Run("websites are checked in the destination", func(t *testing.T) {
		err := e.svc.UpdateApplication(alice, backend.App.ID, &domain.UpdateApplicationArgs{
			Websites: optional.From([]*domain.Website{{
				ID: domain.NewID(), FQDN: "move.example.com", PathPrefix: "/api", HTTPPort: 80,
			}}),
		}, optional.None[[]string]())
		require.NoError(t, err, "nesting is allowed in the same group")

		err = e.svc.MoveApplication(alice, backend.App.ID, optional.None[string]())
		assertErrorType(t, domain.ErrorTypeBadRequest, err)
	})

	t.Run("moves an app out of a group to be standalone", func(t *testing.T) {
		err := e.svc.UpdateApplication(alice, backend.App.ID, &domain.UpdateApplicationArgs{
			Websites: optional.From([]*domain.Website{{
				ID: domain.NewID(), FQDN: "move-api.example.com", PathPrefix: "/", HTTPPort: 80,
			}}),
		}, optional.None[[]string]())
		require.NoError(t, err)

		err = e.svc.MoveApplication(alice, backend.App.ID, optional.None[string]())
		require.NoError(t, err)

		got, err := e.svc.GetApplication(alice, backend.App.ID)
		require.NoError(t, err)
		assert.True(t, got.Group.Implicit)
		assert.ElementsMatch(t, aliceGroup.OwnerIDs, got.Group.OwnerIDs, "the owners of the group it left stay owners")
	})
}

func TestUpdateApplication_StandaloneOwners(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")
	bob := e.as(t, "bob")

	standalone, err := e.svc.CreateApplication(alice, e.newApp("standalone", "owners.example.com", "/"), optional.None[string]())
	require.NoError(t, err)

	err = e.svc.UpdateApplication(alice, standalone.App.ID, &domain.UpdateApplicationArgs{}, optional.From([]string{web.GetUser(alice).ID, web.GetUser(bob).ID}))
	require.NoError(t, err)
	err = e.svc.UpdateApplication(bob, standalone.App.ID, &domain.UpdateApplicationArgs{Name: optional.From("renamed")}, optional.None[[]string]())
	require.NoError(t, err, "a new owner can operate the app")

	err = e.svc.UpdateApplication(alice, standalone.App.ID, &domain.UpdateApplicationArgs{}, optional.From([]string{}))
	assertErrorType(t, domain.ErrorTypeBadRequest, err)

	group, err := e.svc.CreateApplicationGroup(alice, "service")
	require.NoError(t, err)
	grouped, err := e.svc.CreateApplication(alice, e.newApp("grouped", "owners-grouped.example.com", "/"), optional.From(group.ID))
	require.NoError(t, err)
	err = e.svc.UpdateApplication(alice, grouped.App.ID, &domain.UpdateApplicationArgs{}, optional.From([]string{web.GetUser(bob).ID}))
	assertErrorType(t, domain.ErrorTypeBadRequest, err)
}

func TestDeleteApplicationGroup(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")
	bob := e.as(t, "bob")

	group, err := e.svc.CreateApplicationGroup(alice, "service")
	require.NoError(t, err)
	app, err := e.svc.CreateApplication(alice, e.newApp("app", "delete-group.example.com", "/"), optional.From(group.ID))
	require.NoError(t, err)

	err = e.svc.DeleteApplicationGroup(bob, group.ID)
	assertErrorType(t, domain.ErrorTypeForbidden, err)

	err = e.svc.DeleteApplicationGroup(alice, group.ID)
	assertErrorType(t, domain.ErrorTypeBadRequest, err)

	require.NoError(t, e.svc.MoveApplication(alice, app.App.ID, optional.None[string]()))
	require.NoError(t, e.svc.DeleteApplicationGroup(alice, group.ID))

	_, err = e.svc.GetApplicationGroup(alice, group.ID)
	assertErrorType(t, domain.ErrorTypeNotFound, err)
}

func TestGetApplicationGroup_ImplicitIsHidden(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")

	app, err := e.svc.CreateApplication(alice, e.newApp("app", "hidden.example.com", "/"), optional.None[string]())
	require.NoError(t, err)

	_, err = e.svc.GetApplicationGroup(alice, app.Group.ID)
	assertErrorType(t, domain.ErrorTypeNotFound, err)
	_, err = e.svc.CreateApplication(alice, e.newApp("app2", "hidden2.example.com", "/"), optional.From(app.Group.ID))
	assertErrorType(t, domain.ErrorTypeNotFound, err)
	err = e.svc.MoveApplication(alice, app.App.ID, optional.From(app.Group.ID))
	assertErrorType(t, domain.ErrorTypeNotFound, err)
	err = e.svc.UpdateApplicationGroup(alice, app.Group.ID, &domain.UpdateApplicationGroupArgs{Name: optional.From("visible")})
	assertErrorType(t, domain.ErrorTypeNotFound, err)
	err = e.svc.DeleteApplicationGroup(alice, app.Group.ID)
	assertErrorType(t, domain.ErrorTypeNotFound, err)
}

// A rejected owner change used to leave the other changes of the same update applied.
func TestUpdateApplication_OwnersAreValidatedBeforeAnyChange(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")
	bob := e.as(t, "bob")

	group, err := e.svc.CreateApplicationGroup(alice, "service")
	require.NoError(t, err)
	grouped, err := e.svc.CreateApplication(alice, e.newApp("grouped", "atomic-grouped.example.com", "/"), optional.From(group.ID))
	require.NoError(t, err)
	standalone, err := e.svc.CreateApplication(alice, e.newApp("standalone", "atomic-standalone.example.com", "/"), optional.None[string]())
	require.NoError(t, err)

	tests := []struct {
		name     string
		app      *apiserver.TopAppInfo
		ownerIDs []string
	}{
		{"owners of an app in a group", grouped, []string{web.GetUser(bob).ID}},
		{"empty owners of a standalone app", standalone, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := e.svc.UpdateApplication(alice, tt.app.App.ID, &domain.UpdateApplicationArgs{Name: optional.From("renamed")}, optional.From(tt.ownerIDs))
			assertErrorType(t, domain.ErrorTypeBadRequest, err)

			got, err := e.svc.GetApplication(alice, tt.app.App.ID)
			require.NoError(t, err)
			assert.Equal(t, tt.app.App.Name, got.App.Name)
		})
	}

	err = e.svc.UpdateApplication(alice, standalone.App.ID, &domain.UpdateApplicationArgs{Name: optional.From("renamed")},
		optional.From([]string{web.GetUser(bob).ID}))
	require.NoError(t, err, "an owner may hand the app over to others in the same update")
	got, err := e.svc.GetApplication(bob, standalone.App.ID)
	require.NoError(t, err)
	assert.Equal(t, "renamed", got.App.Name)
	assert.Equal(t, []string{web.GetUser(bob).ID}, got.Group.OwnerIDs)
}

func TestUpdateApplicationGroup(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")
	bob := e.as(t, "bob")

	group, err := e.svc.CreateApplicationGroup(alice, "service")
	require.NoError(t, err)

	t.Run("needs ownership", func(t *testing.T) {
		err := e.svc.UpdateApplicationGroup(bob, group.ID, &domain.UpdateApplicationGroupArgs{Name: optional.From("taken")})
		assertErrorType(t, domain.ErrorTypeForbidden, err)
	})

	t.Run("rejects invalid values", func(t *testing.T) {
		tests := []struct {
			name string
			args *domain.UpdateApplicationGroupArgs
		}{
			{"empty name", &domain.UpdateApplicationGroupArgs{Name: optional.From("")}},
			{"name longer than 100 characters", &domain.UpdateApplicationGroupArgs{Name: optional.From(strings.Repeat("あ", 101))}},
			{"empty owners", &domain.UpdateApplicationGroupArgs{OwnerIDs: optional.From([]string{})}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := e.svc.UpdateApplicationGroup(alice, group.ID, tt.args)
				assertErrorType(t, domain.ErrorTypeBadRequest, err)
			})
		}
	})

	t.Run("renames the group", func(t *testing.T) {
		err := e.svc.UpdateApplicationGroup(alice, group.ID, &domain.UpdateApplicationGroupArgs{Name: optional.From(strings.Repeat("あ", 100))})
		require.NoError(t, err, "a name of 100 characters is accepted")

		err = e.svc.UpdateApplicationGroup(alice, group.ID, &domain.UpdateApplicationGroupArgs{Name: optional.From("renamed")})
		require.NoError(t, err)

		got, err := e.svc.GetApplicationGroup(alice, group.ID)
		require.NoError(t, err)
		assert.Equal(t, "renamed", got.Name)
		assert.Equal(t, []string{web.GetUser(alice).ID}, got.OwnerIDs)
	})
}

func TestGetApplicationGroups_Mine(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")
	bob := e.as(t, "bob")

	aliceGroup, err := e.svc.CreateApplicationGroup(alice, "alice-service")
	require.NoError(t, err)
	bobGroup, err := e.svc.CreateApplicationGroup(bob, "bob-service")
	require.NoError(t, err)

	ids := func(groups []*domain.ApplicationGroup) []string {
		return lo.Map(groups, func(g *domain.ApplicationGroup, _ int) string { return g.ID })
	}

	mine, err := e.svc.GetApplicationGroups(alice, apiserver.GetAppGroupScopeMine)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{aliceGroup.ID}, ids(mine), "only the groups the user owns are listed")

	all, err := e.svc.GetApplicationGroups(alice, apiserver.GetAppGroupScopeAll)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{aliceGroup.ID, bobGroup.ID}, ids(all))

	err = e.svc.UpdateApplicationGroup(alice, aliceGroup.ID, &domain.UpdateApplicationGroupArgs{
		OwnerIDs: optional.From([]string{web.GetUser(alice).ID, web.GetUser(bob).ID}),
	})
	require.NoError(t, err)
	bobMine, err := e.svc.GetApplicationGroups(bob, apiserver.GetAppGroupScopeMine)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{aliceGroup.ID, bobGroup.ID}, ids(bobMine), "a group is listed for an owner added later")
}

func TestAdminOperatesGroupsOfOthers(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")
	admin := e.asAdmin(t, "admin")

	aliceGroup, err := e.svc.CreateApplicationGroup(alice, "alice-service")
	require.NoError(t, err)
	app, err := e.svc.CreateApplication(alice, e.newApp("app", "admin.example.com", "/"), optional.From(aliceGroup.ID))
	require.NoError(t, err)
	adminGroup, err := e.svc.CreateApplicationGroup(admin, "admin-service")
	require.NoError(t, err)

	err = e.svc.UpdateApplication(admin, app.App.ID, &domain.UpdateApplicationArgs{Name: optional.From("renamed")}, optional.None[[]string]())
	require.NoError(t, err)
	err = e.svc.UpdateApplicationGroup(admin, aliceGroup.ID, &domain.UpdateApplicationGroupArgs{Name: optional.From("renamed-service")})
	require.NoError(t, err)
	err = e.svc.MoveApplication(admin, app.App.ID, optional.From(adminGroup.ID))
	require.NoError(t, err)

	got, err := e.svc.GetApplication(admin, app.App.ID)
	require.NoError(t, err)
	assert.Equal(t, "renamed", got.App.Name)
	assert.Equal(t, adminGroup.ID, got.Group.ID)
}

func TestUnknownOwnerIDsAreRejected(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")
	aliceID := web.GetUser(alice).ID

	t.Run("standalone app", func(t *testing.T) {
		app, err := e.svc.CreateApplication(alice, e.newApp("standalone", "unknown-owner.example.com", "/"), optional.None[string]())
		require.NoError(t, err)

		err = e.svc.UpdateApplication(alice, app.App.ID, &domain.UpdateApplicationArgs{Name: optional.From("renamed")},
			optional.From([]string{aliceID, "no-such-user"}))
		assertErrorType(t, domain.ErrorTypeNotFound, err)

		got, err := e.svc.GetApplication(alice, app.App.ID)
		require.NoError(t, err)
		assert.Equal(t, app.App.Name, got.App.Name)
		assert.ElementsMatch(t, app.Group.OwnerIDs, got.Group.OwnerIDs)
	})

	t.Run("application group", func(t *testing.T) {
		group, err := e.svc.CreateApplicationGroup(alice, "service")
		require.NoError(t, err)

		err = e.svc.UpdateApplicationGroup(alice, group.ID, &domain.UpdateApplicationGroupArgs{
			Name:     optional.From("renamed"),
			OwnerIDs: optional.From([]string{aliceID, "no-such-user"}),
		})
		assertErrorType(t, domain.ErrorTypeNotFound, err)

		got, err := e.svc.GetApplicationGroup(alice, group.ID)
		require.NoError(t, err)
		assert.Equal(t, "service", got.Name)
		assert.Equal(t, []string{aliceID}, got.OwnerIDs)
	})
}

func TestGetApplications_Scopes(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")
	bob := e.as(t, "bob")

	aliceGroup, err := e.svc.CreateApplicationGroup(alice, "alice-service")
	require.NoError(t, err)
	bobGroup, err := e.svc.CreateApplicationGroup(bob, "bob-service")
	require.NoError(t, err)
	aliceGrouped, err := e.svc.CreateApplication(alice, e.newApp("alice-grouped", "scope-alice-grouped.example.com", "/"), optional.From(aliceGroup.ID))
	require.NoError(t, err)
	aliceStandalone, err := e.svc.CreateApplication(alice, e.newApp("alice-standalone", "scope-alice-standalone.example.com", "/"), optional.None[string]())
	require.NoError(t, err)
	bobGrouped, err := e.svc.CreateApplication(bob, e.newApp("bob-grouped", "scope-bob-grouped.example.com", "/"), optional.From(bobGroup.ID))
	require.NoError(t, err)
	bobStandalone, err := e.svc.CreateApplication(bob, e.newApp("bob-standalone", "scope-bob-standalone.example.com", "/"), optional.None[string]())
	require.NoError(t, err)

	appIDs := func(apps []*apiserver.TopAppInfo) []string {
		return lo.Map(apps, func(a *apiserver.TopAppInfo, _ int) string { return a.App.ID })
	}

	tests := []struct {
		name  string
		scope apiserver.GetAppScope
		want  []*apiserver.TopAppInfo
	}{
		{"mine", apiserver.GetAppScope{Scope: apiserver.GetAppScopeMine}, []*apiserver.TopAppInfo{aliceGrouped, aliceStandalone}},
		{"all", apiserver.GetAppScope{Scope: apiserver.GetAppScopeAll}, []*apiserver.TopAppInfo{aliceGrouped, aliceStandalone, bobGrouped, bobStandalone}},
		{"group", apiserver.GetAppScope{Scope: apiserver.GetAppScopeGroup, GroupID: optional.From(bobGroup.ID)}, []*apiserver.TopAppInfo{bobGrouped}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := e.svc.GetApplications(alice, tt.scope)
			require.NoError(t, err)
			assert.ElementsMatch(t, appIDs(tt.want), appIDs(got))
		})
	}

	t.Run("each app comes with its group", func(t *testing.T) {
		got, err := e.svc.GetApplications(alice, apiserver.GetAppScope{Scope: apiserver.GetAppScopeAll})
		require.NoError(t, err)
		byID := lo.SliceToMap(got, func(a *apiserver.TopAppInfo) (string, *apiserver.TopAppInfo) { return a.App.ID, a })
		for _, want := range []*apiserver.TopAppInfo{aliceGrouped, aliceStandalone, bobGrouped, bobStandalone} {
			require.Contains(t, byID, want.App.ID)
			gotGroup := byID[want.App.ID].Group
			require.NotNil(t, gotGroup, want.App.Name)
			assert.Equal(t, want.Group.ID, gotGroup.ID, want.App.Name)
			assert.Equal(t, want.Group.Implicit, gotGroup.Implicit, want.App.Name)
			assert.ElementsMatch(t, want.Group.OwnerIDs, gotGroup.OwnerIDs, want.App.Name)
		}
		assert.Equal(t, aliceGroup.ID, byID[aliceGrouped.App.ID].Group.ID)
		assert.Equal(t, []string{web.GetUser(bob).ID}, byID[bobGrouped.App.ID].Group.OwnerIDs)
	})

	t.Run("group scope needs a group", func(t *testing.T) {
		_, err := e.svc.GetApplications(alice, apiserver.GetAppScope{Scope: apiserver.GetAppScopeGroup})
		assertErrorType(t, domain.ErrorTypeBadRequest, err)
	})
}

func TestOwnerOnlyOperationsFollowTheGroup(t *testing.T) {
	t.Parallel()
	e := newGroupTestEnv(t)
	alice := e.as(t, "alice")
	bob := e.as(t, "bob")

	group, err := e.svc.CreateApplicationGroup(alice, "service")
	require.NoError(t, err)
	app, err := e.svc.CreateApplication(alice, e.newApp("app", "owner-only.example.com", "/"), optional.From(group.ID))
	require.NoError(t, err)
	appID := app.App.ID

	operations := []struct {
		name string
		call func(ctx context.Context) error
	}{
		{"get environment variables", func(ctx context.Context) error {
			_, err := e.svc.GetEnvironmentVariables(ctx, appID)
			return err
		}},
		{"set an environment variable", func(ctx context.Context) error {
			return e.svc.SetEnvironmentVariable(ctx, appID, "KEY", "value")
		}},
		{"delete an environment variable", func(ctx context.Context) error {
			return e.svc.DeleteEnvironmentVariable(ctx, appID, "KEY")
		}},
		{"start the app", func(ctx context.Context) error { return e.svc.StartApplication(ctx, appID) }},
		{"stop the app", func(ctx context.Context) error { return e.svc.StopApplication(ctx, appID) }},
		{"retry a build", func(ctx context.Context) error { return e.svc.RetryCommitBuild(ctx, appID, exampleCommitHash) }},
		{"delete the app", func(ctx context.Context) error { return e.svc.DeleteApplication(ctx, appID) }},
	}
	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			assertErrorType(t, domain.ErrorTypeForbidden, op.call(bob))
		})
	}

	err = e.svc.UpdateApplicationGroup(alice, group.ID, &domain.UpdateApplicationGroupArgs{
		OwnerIDs: optional.From([]string{web.GetUser(alice).ID, web.GetUser(bob).ID}),
	})
	require.NoError(t, err)
	require.NoError(t, e.svc.SetEnvironmentVariable(bob, appID, "KEY", "value"), "an owner added to the group can operate the app")
	envs, err := e.svc.GetEnvironmentVariables(bob, appID)
	require.NoError(t, err)
	assert.Contains(t, lo.Map(envs, func(env *domain.Environment, _ int) string { return env.Key }), "KEY")
}
