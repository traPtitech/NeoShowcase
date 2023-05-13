package usecase

import (
	"bytes"
	"fmt"

	"github.com/friendsofgo/errors"
	"github.com/gliderlabs/ssh"
	ssh2 "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/fig"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

var (
	figWelcome string
)

func init() {
	orange, err := fig.NewTrueColorFromHexString("FF9900")
	if err != nil {
		panic(err)
	}
	var b fig.Builder
	err = b.Append("Neo", "larry3d", orange)
	if err != nil {
		panic(err)
	}
	err = b.Append("Showcase", "larry3d", fig.ColorWhite)
	if err != nil {
		panic(err)
	}
	figWelcome = b.String()
}

type SSHServer interface {
	Start() error
	Close() error

	isSSHServer()
}

type sshServer struct {
	config   domain.SSHConfig
	sshKey   *ssh2.PublicKeys
	backend  domain.Backend
	appRepo  domain.ApplicationRepository
	userRepo domain.UserRepository

	server *ssh.Server
}

func NewSSHServer(
	config domain.SSHConfig,
	sshKey *ssh2.PublicKeys,
	backend domain.Backend,
	appRepo domain.ApplicationRepository,
	userRepo domain.UserRepository,
) SSHServer {
	s := &sshServer{
		config:   config,
		sshKey:   sshKey,
		backend:  backend,
		appRepo:  appRepo,
		userRepo: userRepo,
	}
	s.server = &ssh.Server{
		Addr:             fmt.Sprintf(":%d", config.Port),
		Handler:          s.handler,
		PublicKeyHandler: s.publicKeyHandler,
	}
	s.server.AddHostKey(sshKey.Signer)
	return s
}

func (s *sshServer) isSSHServer() {}

func (s *sshServer) Start() error {
	go func() {
		_ = s.server.ListenAndServe()
	}()
	return nil
}

func (s *sshServer) Close() error {
	return s.server.Close()
}

func (s *sshServer) publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	appID := ctx.User()
	log.Debugf("authenticating ssh for app %v", appID)
	app, err := s.appRepo.GetApplication(ctx, appID)
	if err != nil {
		log.Errorf("retrieving app with id %v: %+v", appID, err)
		return false
	}
	admins, err := s.userRepo.GetUsers(ctx, domain.GetUserCondition{Admin: optional.From(true)})
	if err != nil {
		log.Errorf("retrieving admin list: %+v", err)
		return false
	}

	var eligibleUsers []string
	eligibleUsers = append(eligibleUsers, app.OwnerIDs...)
	eligibleUsers = append(eligibleUsers, ds.Map(admins, func(u *domain.User) string { return u.ID })...)
	keys, err := s.userRepo.GetUserKeys(ctx, domain.GetUserKeyCondition{UserIDs: optional.From(eligibleUsers)})
	if err != nil {
		log.Errorf("retrieving user keys of app id %v: %+v", appID, err)
		return false
	}

	marshaledKey := key.Marshal()
	return lo.ContainsBy(keys, func(userKey *domain.UserKey) bool {
		return bytes.Equal(userKey.MarshalKey(), marshaledKey)
	})
}

func (s *sshServer) handler(sess ssh.Session) {
	writeErrAndClose := func(err error) {
		log.Errorf("%+v", err)
		_, _ = sess.Write([]byte(err.Error() + "\n"))
		_ = sess.Exit(1)
	}

	appID := sess.User()
	sessID := domain.NewID()
	log.Infof("new ssh connection into app %s (session id: %v)", appID, sessID)
	defer log.Infof("closing ssh connecttion into app %s (session id: %v)", appID, sessID)

	app, err := s.appRepo.GetApplication(sess.Context(), appID)
	if err != nil {
		writeErrAndClose(errors.Wrapf(err, "retrieving app with id %v", appID))
		return
	}

	_, _ = sess.Write([]byte(figWelcome))
	_, _ = sess.Write([]byte{'\n'})
	_, _ = sess.Write([]byte(fmt.Sprintf("Welcome to NeoShowcase! Connecting to application %s (id: %s) ...\n", app.Name, appID)))

	cmd := sess.Command()
	if len(cmd) == 0 {
		cmd = []string{"/bin/sh"}
	}
	err = s.backend.ExecContainer(sess.Context(), appID, cmd, sess, sess, sess.Stderr())
	if err != nil {
		writeErrAndClose(err)
		return
	}
	_ = sess.Exit(0)
}
