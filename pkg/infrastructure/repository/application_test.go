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

func (e *appRepoTestEnv) newApp(group *domain.ApplicationGroup) *domain.Application {
	return &domain.Application{
		ID:           domain.NewID(),
		Name:         "app",
		GroupID:      group.ID,
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
	app := e.newApp(implicit)
	require.NoError(t, e.appRepo.CreateApplication(ctx, app, implicit))
	assert.ElementsMatch(t, []string{explicit.ID, implicit.ID}, e.groupIDs(), "an implicit group is created with its app")

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

	inExplicit := e.newApp(explicit)
	require.NoError(t, e.appRepo.CreateApplication(ctx, inExplicit, explicit))
	require.NoError(t, e.appRepo.DeleteApplication(ctx, inExplicit.ID))
	assert.ElementsMatch(t, []string{explicit.ID}, e.groupIDs(), "deleting an app keeps a group the user created")
}

func TestApplicationRepository_CreateApplicationFailureLeavesNoImplicitGroup(t *testing.T) {
	t.Parallel()
	e := newAppRepoTestEnv(t)
	ctx := t.Context()

	first := domain.NewImplicitApplicationGroup([]string{e.ownerID})
	app := e.newApp(first)
	require.NoError(t, e.appRepo.CreateApplication(ctx, app, first))

	second := domain.NewImplicitApplicationGroup([]string{e.ownerID})
	duplicate := e.newApp(second)
	duplicate.ID = app.ID
	require.Error(t, e.appRepo.CreateApplication(ctx, duplicate, second))
	assert.ElementsMatch(t, []string{first.ID}, e.groupIDs())
}
