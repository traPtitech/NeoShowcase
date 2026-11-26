package repository_test

import (
	"database/sql"
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
	db        *sql.DB
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
		db:        db,
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

// createStandaloneApp creates an app in its implicit group owned by the owner of the environment.
func (e *appRepoTestEnv) createStandaloneApp() *domain.Application {
	e.t.Helper()
	app := e.newApp("")
	group := app.MakeStandalone([]string{e.ownerID})
	require.NoError(e.t, e.appRepo.CreateApplication(e.t.Context(), app, optional.From(group)))
	return app
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

	app := e.createStandaloneApp()
	assert.ElementsMatch(t, []string{explicit.ID, app.GroupID}, e.groupIDs(), "an implicit group is created with its app")
	created, err := e.groupRepo.GetApplicationGroup(ctx, app.GroupID)
	require.NoError(t, err)
	assert.True(t, created.Implicit)
	assert.Equal(t, []string{e.ownerID}, created.OwnerIDs)

	require.NoError(t, e.appRepo.SetApplicationGroup(ctx, app.ID, app.GroupID, explicit.ID, optional.None[*domain.ApplicationGroup]()))
	assert.ElementsMatch(t, []string{explicit.ID}, e.groupIDs(), "an app leaving its implicit group deletes the group")

	back := domain.NewImplicitApplicationGroup([]string{e.ownerID})
	backID := back.ID
	require.NoError(t, e.appRepo.SetApplicationGroup(ctx, app.ID, explicit.ID, backID, optional.From(back)))
	got, err := e.appRepo.GetApplication(ctx, app.ID)
	require.NoError(t, err)
	assert.Equal(t, backID, got.GroupID)
	gotBack, err := e.groupRepo.GetApplicationGroup(ctx, backID)
	require.NoError(t, err)
	assert.True(t, gotBack.Implicit)
	assert.ElementsMatch(t, []string{explicit.ID, backID}, e.groupIDs(), "an app leaving a group the user created leaves the group")

	require.NoError(t, e.appRepo.DeleteApplication(ctx, app.ID))
	assert.ElementsMatch(t, []string{explicit.ID}, e.groupIDs(), "deleting an app deletes its implicit group")

	inExplicit := e.newApp(explicit.ID)
	require.NoError(t, e.appRepo.CreateApplication(ctx, inExplicit, optional.None[*domain.ApplicationGroup]()))
	require.NoError(t, e.appRepo.DeleteApplication(ctx, inExplicit.ID))
	assert.ElementsMatch(t, []string{explicit.ID}, e.groupIDs(), "deleting an app keeps a group the user created")
}

func TestApplicationRepository_CreateApplicationFailureLeavesNoImplicitGroup(t *testing.T) {
	t.Parallel()
	e := newAppRepoTestEnv(t)
	ctx := t.Context()

	app := e.createStandaloneApp()

	duplicate := e.newApp("")
	duplicate.ID = app.ID
	duplicateGroup := duplicate.MakeStandalone([]string{e.ownerID})
	require.Error(t, e.appRepo.CreateApplication(ctx, duplicate, optional.From(duplicateGroup)))
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

func TestApplicationRepository_ImplicitGroupHoldsOnlyItsApp(t *testing.T) {
	t.Parallel()
	e := newAppRepoTestEnv(t)
	ctx := t.Context()
	explicit := e.createExplicitGroup()
	standalone := e.createStandaloneApp()

	joining := e.newApp(standalone.GroupID)
	require.Error(t, e.appRepo.CreateApplication(ctx, joining, optional.None[*domain.ApplicationGroup]()), "a new app cannot join the implicit group of another app")

	inExplicit := e.newApp(explicit.ID)
	require.NoError(t, e.appRepo.CreateApplication(ctx, inExplicit, optional.None[*domain.ApplicationGroup]()))
	err := e.appRepo.SetApplicationGroup(ctx, inExplicit.ID, explicit.ID, standalone.GroupID, optional.None[*domain.ApplicationGroup]())
	require.Error(t, err, "an app cannot move into the implicit group of another app")
	got, err := e.appRepo.GetApplication(ctx, inExplicit.ID)
	require.NoError(t, err)
	assert.Equal(t, explicit.ID, got.GroupID)
}

// The usecase checks the group of an app before it calls the repository, and the app can move in between.
func TestApplicationRepository_ChangesFailWhenTheAppHasMoved(t *testing.T) {
	t.Parallel()
	e := newAppRepoTestEnv(t)
	ctx := t.Context()
	other, err := repository.NewUserRepository(e.db).EnsureUser(ctx, "other")
	require.NoError(t, err)
	explicit := domain.NewApplicationGroup("service", []string{e.ownerID, other.ID})
	require.NoError(t, e.groupRepo.CreateApplicationGroup(ctx, explicit))
	app := e.createStandaloneApp()
	checked := app.GroupID
	require.NoError(t, e.appRepo.SetApplicationGroup(ctx, app.ID, checked, explicit.ID, optional.None[*domain.ApplicationGroup]()))

	err = e.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{}, optional.From(domain.ImplicitGroupOwners{
		GroupID:  checked,
		OwnerIDs: []string{e.ownerID},
	}))
	_, typ, _ := domain.DecomposeError(err)
	assert.Equal(t, domain.ErrorTypeFailedPrecondition, typ, "owners checked on the old implicit group are not set on the new group")
	got, err := e.groupRepo.GetApplicationGroup(ctx, explicit.ID)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{e.ownerID, other.ID}, got.OwnerIDs)

	back := domain.NewImplicitApplicationGroup([]string{e.ownerID})
	err = e.appRepo.SetApplicationGroup(ctx, app.ID, checked, back.ID, optional.From(back))
	_, typ, _ = domain.DecomposeError(err)
	assert.Equal(t, domain.ErrorTypeFailedPrecondition, typ, "a move checked on the old group does not move the app")
	moved, err := e.appRepo.GetApplication(ctx, app.ID)
	require.NoError(t, err)
	assert.Equal(t, explicit.ID, moved.GroupID)
	assert.ElementsMatch(t, []string{explicit.ID}, e.groupIDs())
}

func TestApplicationRepository_UpdateOwnersOfStandaloneApp(t *testing.T) {
	t.Parallel()
	e := newAppRepoTestEnv(t)
	ctx := t.Context()
	other, err := repository.NewUserRepository(e.db).EnsureUser(ctx, "other")
	require.NoError(t, err)
	app := e.createStandaloneApp()

	err = e.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{}, optional.From(domain.ImplicitGroupOwners{
		GroupID:  app.GroupID,
		OwnerIDs: []string{e.ownerID, other.ID},
	}))
	require.NoError(t, err)
	got, err := e.groupRepo.GetApplicationGroup(ctx, app.GroupID)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{e.ownerID, other.ID}, got.OwnerIDs)
}

func TestApplicationRepository_UpdateOwnersRefusesExplicitGroup(t *testing.T) {
	t.Parallel()
	e := newAppRepoTestEnv(t)
	ctx := t.Context()
	other, err := repository.NewUserRepository(e.db).EnsureUser(ctx, "other")
	require.NoError(t, err)
	explicit := e.createExplicitGroup()
	app := e.newApp(explicit.ID)
	require.NoError(t, e.appRepo.CreateApplication(ctx, app, optional.None[*domain.ApplicationGroup]()))

	err = e.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{}, optional.From(domain.ImplicitGroupOwners{
		GroupID:  explicit.ID,
		OwnerIDs: []string{other.ID},
	}))
	require.Error(t, err, "owners of a group a user created are not replaced through an app")
	got, err := e.groupRepo.GetApplicationGroup(ctx, explicit.ID)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{e.ownerID}, got.OwnerIDs)
}

func TestApplicationRepository_CreateApplicationRefusesAGroupOfAnotherID(t *testing.T) {
	t.Parallel()
	e := newAppRepoTestEnv(t)
	ctx := t.Context()

	app := e.newApp(domain.NewID())
	other := domain.NewImplicitApplicationGroup([]string{e.ownerID})
	require.Error(t, e.appRepo.CreateApplication(ctx, app, optional.From(other)), "the implicit group has to be the group of the app")
	assert.Empty(t, e.groupIDs())
}
