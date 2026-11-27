package giteaintegration

import (
	"testing"

	"code.gitea.io/sdk/gitea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/test/testhelper"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func newTestApp(repoID string, groupID string) *domain.Application {
	return &domain.Application{
		ID:           domain.NewID(),
		Name:         "app",
		GroupID:      groupID,
		RepositoryID: repoID,
		RefName:      "main",
		Commit:       domain.EmptyCommit,
		DeployType:   domain.DeployTypeRuntime,
		Container:    domain.ContainerStateMissing,
		Config: domain.ApplicationConfig{
			BuildConfig: &domain.BuildConfigRuntimeBuildpack{Context: "."},
		},
	}
}

func newSyncTestIntegration(t *testing.T) *Integration {
	t.Helper()
	db := testhelper.OpenDB(t)
	return &Integration{
		gitRepo:   repository.NewGitRepositoryRepository(db),
		appRepo:   repository.NewApplicationRepository(db),
		groupRepo: repository.NewApplicationGroupRepository(db),
		userRepo:  repository.NewUserRepository(db),
	}
}

// createStandaloneApp creates a standalone app of the repository and returns the ID of its implicit group.
func (i *Integration) createStandaloneApp(t *testing.T, repoID string, ownerIDs []string) string {
	t.Helper()
	app := newTestApp(repoID, "")
	group := app.MakeStandalone(ownerIDs)
	require.NoError(t, i.appRepo.CreateApplication(t.Context(), app, optional.From(group)))
	return app.GroupID
}

func (i *Integration) ownersOfGroup(t *testing.T, groupID string) []string {
	t.Helper()
	group, err := i.groupRepo.GetApplicationGroup(t.Context(), groupID)
	require.NoError(t, err)
	return group.OwnerIDs
}

// Writers of one repository used to become owners of whole groups, including apps from other repositories.
func TestSyncRepository_AddsWritersOnlyToStandaloneApps(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	i := newSyncTestIntegration(t)

	owner, err := i.userRepo.EnsureUser(ctx, "owner")
	require.NoError(t, err)
	writer, err := i.userRepo.EnsureUser(ctx, "writer")
	require.NoError(t, err)

	repoX := domain.NewRepository("org/x", "ssh://git@example.com/org/x.git", optional.None[domain.RepositoryAuth](), []string{owner.ID})
	repoY := domain.NewRepository("org/y", "ssh://git@example.com/org/y.git", optional.None[domain.RepositoryAuth](), []string{owner.ID})
	require.NoError(t, i.gitRepo.CreateRepository(ctx, repoX))
	require.NoError(t, i.gitRepo.CreateRepository(ctx, repoY))

	shared := domain.NewApplicationGroup("service", []string{owner.ID})
	require.NoError(t, i.groupRepo.CreateApplicationGroup(ctx, shared))
	require.NoError(t, i.appRepo.CreateApplication(ctx, newTestApp(repoX.ID, shared.ID), optional.None[*domain.ApplicationGroup]()))
	require.NoError(t, i.appRepo.CreateApplication(ctx, newTestApp(repoY.ID, shared.ID), optional.None[*domain.ApplicationGroup]()))
	standalone := i.createStandaloneApp(t, repoX.ID, []string{owner.ID})
	standaloneOfY := i.createStandaloneApp(t, repoY.ID, []string{owner.ID})

	err = i.syncRepository(ctx, "org", []string{writer.ID}, &gitea.Repository{Name: "x", SSHURL: repoX.URL})
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{owner.ID}, i.ownersOfGroup(t, shared.ID), "owners of a group users created are left to its owners")
	assert.ElementsMatch(t, []string{owner.ID, writer.ID}, i.ownersOfGroup(t, standalone))
	assert.ElementsMatch(t, []string{owner.ID}, i.ownersOfGroup(t, standaloneOfY), "standalone apps of other repositories are left untouched")
}

func TestSyncRepository_RepositoryWithoutApps(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	i := newSyncTestIntegration(t)

	owner, err := i.userRepo.EnsureUser(ctx, "owner")
	require.NoError(t, err)
	writer, err := i.userRepo.EnsureUser(ctx, "writer")
	require.NoError(t, err)

	empty := domain.NewRepository("org/empty", "ssh://git@example.com/org/empty.git", optional.None[domain.RepositoryAuth](), []string{owner.ID})
	other := domain.NewRepository("org/other", "ssh://git@example.com/org/other.git", optional.None[domain.RepositoryAuth](), []string{owner.ID})
	require.NoError(t, i.gitRepo.CreateRepository(ctx, empty))
	require.NoError(t, i.gitRepo.CreateRepository(ctx, other))
	standaloneOfOther := i.createStandaloneApp(t, other.ID, []string{owner.ID})

	err = i.syncRepository(ctx, "org", []string{writer.ID}, &gitea.Repository{Name: "empty", SSHURL: empty.URL})
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{owner.ID}, i.ownersOfGroup(t, standaloneOfOther))
}
