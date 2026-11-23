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

func newTestApp(repoID string, group *domain.ApplicationGroup) *domain.Application {
	return &domain.Application{
		ID:           domain.NewID(),
		Name:         "app",
		GroupID:      group.ID,
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

// Writers of one repository used to become owners of whole groups, including apps from other repositories.
func TestSyncRepository_AddsWritersOnlyToStandaloneApps(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	db := testhelper.OpenDB(t)
	userRepo := repository.NewUserRepository(db)
	i := &Integration{
		gitRepo:   repository.NewGitRepositoryRepository(db),
		appRepo:   repository.NewApplicationRepository(db),
		groupRepo: repository.NewApplicationGroupRepository(db),
		userRepo:  userRepo,
	}

	owner, err := userRepo.EnsureUser(ctx, "owner")
	require.NoError(t, err)
	writer, err := userRepo.EnsureUser(ctx, "writer")
	require.NoError(t, err)

	repoX := domain.NewRepository("org/x", "ssh://git@example.com/org/x.git", optional.None[domain.RepositoryAuth](), []string{owner.ID})
	repoY := domain.NewRepository("org/y", "ssh://git@example.com/org/y.git", optional.None[domain.RepositoryAuth](), []string{owner.ID})
	require.NoError(t, i.gitRepo.CreateRepository(ctx, repoX))
	require.NoError(t, i.gitRepo.CreateRepository(ctx, repoY))

	shared := domain.NewApplicationGroup("service", []string{owner.ID})
	require.NoError(t, i.groupRepo.CreateApplicationGroup(ctx, shared))
	require.NoError(t, i.appRepo.CreateApplication(ctx, newTestApp(repoX.ID, shared), optional.None[[]string]()))
	require.NoError(t, i.appRepo.CreateApplication(ctx, newTestApp(repoY.ID, shared), optional.None[[]string]()))
	standalone := domain.NewImplicitApplicationGroup([]string{owner.ID})
	require.NoError(t, i.appRepo.CreateApplication(ctx, newTestApp(repoX.ID, standalone), optional.From(standalone.OwnerIDs)))

	err = i.syncRepository(ctx, "org", []string{writer.ID}, &gitea.Repository{Name: "x", SSHURL: repoX.URL})
	require.NoError(t, err)

	gotShared, err := i.groupRepo.GetApplicationGroup(ctx, shared.ID)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{owner.ID}, gotShared.OwnerIDs, "owners of a group users created are left to its owners")
	gotStandalone, err := i.groupRepo.GetApplicationGroup(ctx, standalone.ID)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{owner.ID, writer.ID}, gotStandalone.OwnerIDs)
}
