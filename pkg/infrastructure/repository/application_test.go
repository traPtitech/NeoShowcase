package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/test/testhelper"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type appRepoTestEnv struct {
	t         *testing.T
	appRepo   domain.ApplicationRepository
	groupRepo domain.ApplicationGroupRepository
	repoID    string
	ownerID   string
}

func newAppRepoTestEnv(t *testing.T) *appRepoTestEnv {
	t.Helper()
	ctx := t.Context()
	db := testhelper.OpenDB(t)
	user, err := repository.NewUserRepository(db).EnsureUser(ctx, "owner")
	require.NoError(t, err)
	repo := domain.NewRepository("repo", "https://example.com/repo", optional.None[domain.RepositoryAuth](), []string{user.ID})
	require.NoError(t, repository.NewGitRepositoryRepository(db).CreateRepository(ctx, repo))
	return &appRepoTestEnv{
		t:         t,
		appRepo:   repository.NewApplicationRepository(db),
		groupRepo: repository.NewApplicationGroupRepository(db),
		repoID:    repo.ID,
		ownerID:   user.ID,
	}
}

func (e *appRepoTestEnv) newApp(groupID string) *domain.Application {
	return &domain.Application{
		ID:           domain.NewID(),
		Name:         "app",
		GroupID:      groupID,
		RepositoryID: e.repoID,
		RefName:      "main",
		Commit:       domain.EmptyCommit,
		DeployType:   domain.DeployTypeRuntime,
		Container:    domain.ContainerStateMissing,
		Config: domain.ApplicationConfig{
			BuildConfig: &domain.BuildConfigRuntimeBuildpack{Context: "."},
		},
	}
}

func (e *appRepoTestEnv) createExplicitGroup() *domain.ApplicationGroup {
	e.t.Helper()
	group := domain.NewApplicationGroup("service", []string{e.ownerID})
	require.NoError(e.t, e.groupRepo.CreateApplicationGroup(e.t.Context(), group))
	return group
}

func (e *appRepoTestEnv) groupIDs() []string {
	e.t.Helper()
	groups, err := e.groupRepo.GetApplicationGroups(e.t.Context(), domain.GetApplicationGroupCondition{})
	require.NoError(e.t, err)
	ids := make([]string, 0, len(groups))
	for _, g := range groups {
		ids = append(ids, g.ID)
	}
	return ids
}

func TestApplicationRepository_ImplicitGroupLivesWithItsApp(t *testing.T) {
	t.Parallel()
	e := newAppRepoTestEnv(t)
	ctx := t.Context()
	explicit := e.createExplicitGroup()

	implicit := domain.NewImplicitApplicationGroup([]string{e.ownerID})
	app := e.newApp(implicit.ID)
	require.NoError(t, e.appRepo.CreateApplication(ctx, app, implicit))
	assert.ElementsMatch(t, []string{explicit.ID, app.GroupID}, e.groupIDs(), "an implicit group is created with its app")
	created, err := e.groupRepo.GetApplicationGroup(ctx, app.GroupID)
	require.NoError(t, err)
	assert.True(t, created.Implicit)
	assert.Equal(t, []string{e.ownerID}, created.OwnerIDs)

	require.NoError(t, e.appRepo.SetApplicationGroup(ctx, app.ID, explicit))
	assert.ElementsMatch(t, []string{explicit.ID}, e.groupIDs(), "an app leaving its implicit group deletes the group")

	back := domain.NewImplicitApplicationGroup([]string{e.ownerID})
	require.NoError(t, e.appRepo.SetApplicationGroup(ctx, app.ID, back))
	got, err := e.appRepo.GetApplication(ctx, app.ID)
	require.NoError(t, err)
	assert.Equal(t, back.ID, got.GroupID)
	assert.ElementsMatch(t, []string{explicit.ID, back.ID}, e.groupIDs(), "an app leaving a group the user created leaves the group")

	require.NoError(t, e.appRepo.DeleteApplication(ctx, app.ID))
	assert.ElementsMatch(t, []string{explicit.ID}, e.groupIDs(), "deleting an app deletes its implicit group")

	inExplicit := e.newApp(explicit.ID)
	require.NoError(t, e.appRepo.CreateApplication(ctx, inExplicit, explicit))
	require.NoError(t, e.appRepo.DeleteApplication(ctx, inExplicit.ID))
	assert.ElementsMatch(t, []string{explicit.ID}, e.groupIDs(), "deleting an app keeps a group the user created")
}

func TestApplicationRepository_CreateApplicationFailureLeavesNoImplicitGroup(t *testing.T) {
	t.Parallel()
	e := newAppRepoTestEnv(t)
	ctx := t.Context()

	group := domain.NewImplicitApplicationGroup([]string{e.ownerID})
	app := e.newApp(group.ID)
	require.NoError(t, e.appRepo.CreateApplication(ctx, app, group))

	duplicateGroup := domain.NewImplicitApplicationGroup([]string{e.ownerID})
	duplicate := e.newApp(duplicateGroup.ID)
	duplicate.ID = app.ID
	require.Error(t, e.appRepo.CreateApplication(ctx, duplicate, duplicateGroup))
	assert.ElementsMatch(t, []string{app.GroupID}, e.groupIDs())
}

// Applications created before groups keep their owners in application_owners, which references the application.
func TestApplicationRepository_DeleteMigratedApplication(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	db := testhelper.OpenDB(t)
	user, err := repository.NewUserRepository(db).EnsureUser(ctx, "owner")
	require.NoError(t, err)
	repo := domain.NewRepository("repo", "https://example.com/repo", optional.None[domain.RepositoryAuth](), []string{user.ID})
	require.NoError(t, repository.NewGitRepositoryRepository(db).CreateRepository(ctx, repo))
	appID := insertLegacyApp(t, db, repo.ID, legacyApp{name: "legacy", fqdn: "legacy.example.com", path: "/", ownerIDs: []string{user.ID}})
	require.NoError(t, repository.MigrateApplicationGroups(ctx, db))
	appRepo := repository.NewApplicationRepository(db)

	require.NoError(t, appRepo.DeleteApplication(ctx, appID))

	_, err = appRepo.GetApplication(ctx, appID)
	_, typ, _ := domain.DecomposeError(err)
	assert.Equal(t, domain.ErrorTypeNotFound, typ)
}

func TestApplicationRepository_CreateApplicationRejectsAnotherGroup(t *testing.T) {
	t.Parallel()
	e := newAppRepoTestEnv(t)
	ctx := t.Context()
	explicit := e.createExplicitGroup()

	group := domain.NewImplicitApplicationGroup([]string{e.ownerID})
	app := e.newApp(explicit.ID)
	require.Error(t, e.appRepo.CreateApplication(ctx, app, group))
	assert.ElementsMatch(t, []string{explicit.ID}, e.groupIDs())
}
