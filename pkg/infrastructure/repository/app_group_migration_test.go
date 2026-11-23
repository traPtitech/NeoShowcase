package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/test/testhelper"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type legacyApp struct {
	name     string
	fqdn     string
	path     string
	ownerIDs []string
}

// insertLegacyApp inserts an application as it was before groups: no group, owners in application_owners.
func insertLegacyApp(t *testing.T, db *sql.DB, repoID string, app legacyApp) string {
	t.Helper()
	ctx := t.Context()
	id := domain.NewID()
	_, err := db.ExecContext(ctx, `INSERT INTO applications
		(id, name, repository_id, ref_name, commit, deploy_type, running, container, container_message, current_build, created_at, updated_at)
		VALUES (?, ?, ?, 'main', ?, 'runtime', 0, 'missing', '', '', NOW(), NOW())`,
		id, app.name, repoID, domain.EmptyCommit)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO application_config
		(application_id, use_mariadb, use_mongodb, build_type, base_image, build_cmd, artifact_path, spa, dockerfile_name, context, entrypoint, command)
		VALUES (?, 0, 0, 'runtime-buildpack', '', '', '', 0, '', '.', '', '')`, id)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO websites
		(id, fqdn, path_prefix, strip_prefix, https, h2c, http_port, authentication, application_id)
		VALUES (?, ?, ?, 0, 1, 0, 80, 'off', ?)`, domain.NewID(), app.fqdn, app.path, id)
	require.NoError(t, err)
	for _, ownerID := range app.ownerIDs {
		_, err = db.ExecContext(ctx, `INSERT INTO application_owners (user_id, application_id) VALUES (?, ?)`, ownerID, id)
		require.NoError(t, err)
	}
	return id
}

func TestMigrateApplicationGroups(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	db := testhelper.OpenDB(t)
	userRepo := repository.NewUserRepository(db)
	gitRepo := repository.NewGitRepositoryRepository(db)
	appRepo := repository.NewApplicationRepository(db)
	groupRepo := repository.NewApplicationGroupRepository(db)

	ensureUser := func(name string) string {
		t.Helper()
		user, err := userRepo.EnsureUser(ctx, name)
		require.NoError(t, err)
		return user.ID
	}
	u1, u2, u3 := ensureUser("u1"), ensureUser("u2"), ensureUser("u3")
	repo := domain.NewRepository("repo", "https://example.com/repo", optional.None[domain.RepositoryAuth](), []string{u1})
	require.NoError(t, gitRepo.CreateRepository(ctx, repo))

	front := insertLegacyApp(t, db, repo.ID, legacyApp{name: "front", fqdn: "svc.example.com", path: "/", ownerIDs: []string{u1}})
	api := insertLegacyApp(t, db, repo.ID, legacyApp{name: "api", fqdn: "svc.example.com", path: "/api", ownerIDs: []string{u2}})
	// Nested only under api, so it joins the group through api.
	apiV2 := insertLegacyApp(t, db, repo.ID, legacyApp{name: "api-v2", fqdn: "svc.example.com", path: "/api/v2", ownerIDs: []string{u1}})
	other := insertLegacyApp(t, db, repo.ID, legacyApp{name: "other", fqdn: "other.example.com", path: "/", ownerIDs: []string{u3}})

	require.NoError(t, repository.MigrateApplicationGroups(ctx, db))

	groupOf := func(appID string) *domain.ApplicationGroup {
		t.Helper()
		app, err := appRepo.GetApplication(ctx, appID)
		require.NoError(t, err)
		group, err := groupRepo.GetApplicationGroup(ctx, app.GroupID)
		require.NoError(t, err)
		return group
	}

	service := groupOf(front)
	assert.False(t, service.Implicit, "apps with nested URLs are put into a group users see")
	assert.Equal(t, "front", service.Name, "the group is named after the app at the shortest path")
	assert.ElementsMatch(t, []string{u1, u2}, service.OwnerIDs, "owners are the union of the owners of the apps")
	assert.Equal(t, service.ID, groupOf(api).ID)
	assert.Equal(t, service.ID, groupOf(apiV2).ID)

	standalone := groupOf(other)
	assert.True(t, standalone.Implicit)
	assert.ElementsMatch(t, []string{u3}, standalone.OwnerIDs)

	t.Run("running again changes nothing", func(t *testing.T) {
		before, err := groupRepo.GetApplicationGroups(ctx, domain.GetApplicationGroupCondition{})
		require.NoError(t, err)

		require.NoError(t, repository.MigrateApplicationGroups(ctx, db))

		after, err := groupRepo.GetApplicationGroups(ctx, domain.GetApplicationGroupCondition{})
		require.NoError(t, err)
		assert.ElementsMatch(t, before, after)
		assert.Equal(t, service.ID, groupOf(front).ID)
	})

	t.Run("running again migrates apps left without a group", func(t *testing.T) {
		late := insertLegacyApp(t, db, repo.ID, legacyApp{name: "late", fqdn: "late.example.com", path: "/", ownerIDs: []string{u2}})

		require.NoError(t, repository.MigrateApplicationGroups(ctx, db))

		group := groupOf(late)
		assert.True(t, group.Implicit)
		assert.ElementsMatch(t, []string{u2}, group.OwnerIDs)
	})
}

func TestMigrateApplicationGroups_NoApps(t *testing.T) {
	t.Parallel()
	db := testhelper.OpenDB(t)
	require.NoError(t, repository.MigrateApplicationGroups(context.Background(), db))
}
