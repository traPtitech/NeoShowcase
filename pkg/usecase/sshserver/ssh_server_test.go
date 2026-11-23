package sshserver

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gossh "golang.org/x/crypto/ssh"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/test/testhelper"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

// loginContext is the context of a login as the user, which names the application to enter.
// Methods of ssh.Context that the handler does not use are left nil.
type loginContext struct {
	ssh.Context
	ctx  context.Context
	user string
}

func (c loginContext) User() string                { return c.user }
func (c loginContext) Deadline() (time.Time, bool) { return c.ctx.Deadline() }
func (c loginContext) Done() <-chan struct{}       { return c.ctx.Done() }
func (c loginContext) Err() error                  { return c.ctx.Err() }
func (c loginContext) Value(key any) any           { return c.ctx.Value(key) }

type sshTestEnv struct {
	t      *testing.T
	server *sshServer
	repoID string
}

func newSSHTestEnv(t *testing.T) *sshTestEnv {
	t.Helper()
	db := testhelper.OpenDB(t)
	e := &sshTestEnv{
		t: t,
		server: &sshServer{
			appRepo:   repository.NewApplicationRepository(db),
			groupRepo: repository.NewApplicationGroupRepository(db),
			userRepo:  repository.NewUserRepository(db),
		},
	}
	owner := e.newUser("repo-owner")
	repo := domain.NewRepository("repo", "https://example.com/repo", optional.None[domain.RepositoryAuth](), []string{owner.ID})
	require.NoError(t, repository.NewGitRepositoryRepository(db).CreateRepository(t.Context(), repo))
	e.repoID = repo.ID
	return e
}

func (e *sshTestEnv) newUser(name string) *domain.User {
	e.t.Helper()
	user, err := e.server.userRepo.EnsureUser(e.t.Context(), name)
	require.NoError(e.t, err)
	return user
}

// newKey registers a new key of the user and returns it.
func (e *sshTestEnv) newKey(user *domain.User) gossh.PublicKey {
	e.t.Helper()
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(e.t, err)
	key, err := gossh.NewPublicKey(pub)
	require.NoError(e.t, err)
	userKey, err := domain.NewUserKey(user.ID, string(gossh.MarshalAuthorizedKey(key)), "key")
	require.NoError(e.t, err)
	require.NoError(e.t, e.server.userRepo.CreateUserKey(e.t.Context(), userKey))
	return key
}

func (e *sshTestEnv) newApp(groupID string, standaloneOwnerIDs optional.Of[[]string]) *domain.Application {
	e.t.Helper()
	app := &domain.Application{
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
	require.NoError(e.t, e.server.appRepo.CreateApplication(e.t.Context(), app, standaloneOwnerIDs))
	return app
}

func (e *sshTestEnv) canLogIn(app *domain.Application, key gossh.PublicKey) bool {
	return e.server.publicKeyHandler(loginContext{ctx: e.t.Context(), user: app.ID}, key)
}

func TestPublicKeyHandler_LetsInOwnersOfTheGroupOfTheApp(t *testing.T) {
	t.Parallel()
	e := newSSHTestEnv(t)
	groupOwner := e.newUser("group-owner")
	standaloneOwner := e.newUser("standalone-owner")
	groupOwnerKey := e.newKey(groupOwner)
	standaloneOwnerKey := e.newKey(standaloneOwner)

	group := domain.NewApplicationGroup("service", []string{groupOwner.ID})
	require.NoError(t, e.server.groupRepo.CreateApplicationGroup(t.Context(), group))
	inGroup := e.newApp(group.ID, optional.None[[]string]())
	standalone := e.newApp(domain.NewID(), optional.From([]string{standaloneOwner.ID}))

	assert.True(t, e.canLogIn(inGroup, groupOwnerKey))
	assert.False(t, e.canLogIn(inGroup, standaloneOwnerKey), "owning another app does not let a user in")
	assert.True(t, e.canLogIn(standalone, standaloneOwnerKey))
	assert.False(t, e.canLogIn(standalone, groupOwnerKey))
}
