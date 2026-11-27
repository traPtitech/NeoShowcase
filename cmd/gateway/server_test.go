package gateway

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/test/testhelper"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

// The previous version keeps creating applications without a group until it stops,
// which can be after the data migration of ns-migrate has run.
func TestServerStart_MigratesApplicationsWithoutGroup(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	db := testhelper.OpenDB(t)
	user, err := repository.NewUserRepository(db).EnsureUser(ctx, "owner")
	require.NoError(t, err)
	repo := domain.NewRepository("repo", "https://example.com/repo", optional.None[domain.RepositoryAuth](), []string{user.ID})
	require.NoError(t, repository.NewGitRepositoryRepository(db).CreateRepository(ctx, repo))

	appID := domain.NewID()
	_, err = db.ExecContext(ctx, `INSERT INTO applications
		(id, name, repository_id, ref_name, commit, deploy_type, running, container, container_message, current_build, created_at, updated_at)
		VALUES (?, 'legacy', ?, 'main', ?, 'runtime', 0, 'missing', '', '', NOW(), NOW())`,
		appID, repo.ID, domain.EmptyCommit)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO application_config
		(application_id, use_mariadb, use_mongodb, build_type, base_image, build_cmd, artifact_path, spa, dockerfile_name, context, entrypoint, command)
		VALUES (?, 0, 0, 'runtime-buildpack', '', '', '', 0, '', '.', '', '')`, appID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO application_owners (user_id, application_id) VALUES (?, ?)`, user.ID, appID)
	require.NoError(t, err)

	s := &Server{
		APIServer: &APIServer{H2CServer: web.NewH2CServer(web.H2CConfig{SetupRoute: func(*http.ServeMux) {}})},
		DB:        db,
	}
	go func() { _ = s.Start(ctx) }()
	t.Cleanup(func() { _ = s.APIServer.Shutdown(ctx) })

	appRepo := repository.NewApplicationRepository(db)
	groupRepo := repository.NewApplicationGroupRepository(db)
	var app *domain.Application
	require.Eventually(t, func() bool {
		app, err = appRepo.GetApplication(ctx, appID)
		require.NoError(t, err)
		return app.GroupID != ""
	}, 10*time.Second, 50*time.Millisecond)
	group, err := groupRepo.GetApplicationGroup(ctx, app.GroupID)
	require.NoError(t, err)
	assert.Equal(t, []string{user.ID}, group.OwnerIDs)
}
