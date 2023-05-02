package k8simpl

import (
	"bytes"
	"fmt"

	"github.com/friendsofgo/errors"
	"github.com/gliderlabs/ssh"
	ssh2 "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type sshConfig struct {
	Port int `mapstructure:"port" yaml:"port"`
}

type sshServer struct {
	s *ssh.Server
}

func newSSHServer(
	b *k8sBackend,
	config sshConfig,
	key *ssh2.PublicKeys,
) *sshServer {
	s := &ssh.Server{
		Addr:             fmt.Sprintf(":%d", config.Port),
		Handler:          b.sshHandler,
		PublicKeyHandler: b.sshPublicKeyHandler,
	}
	s.AddHostKey(key.Signer)
	return &sshServer{s: s}
}

func (s *sshServer) Start() {
	go func() {
		_ = s.s.ListenAndServe()
	}()
}

func (s *sshServer) Close() error {
	return s.s.Close()
}

func (b *k8sBackend) sshPublicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	appID := ctx.User()
	log.Debugf("authenticating ssh for app %v", appID)
	app, err := b.appRepo.GetApplication(ctx, appID)
	if err != nil {
		log.Errorf("retrieving app with id %v: %+v", appID, err)
		return false
	}
	admins, err := b.userRepo.GetUsers(ctx, domain.GetUserCondition{Admin: optional.From(true)})
	if err != nil {
		log.Errorf("retrieving admin list: %+v", err)
		return false
	}

	var eligibleUsers []string
	eligibleUsers = append(eligibleUsers, app.OwnerIDs...)
	eligibleUsers = append(eligibleUsers, lo.Map(admins, func(u *domain.User, _ int) string { return u.ID })...)
	keys, err := b.userRepo.GetUserKeys(ctx, domain.GetUserKeyCondition{UserIDs: optional.From(eligibleUsers)})
	if err != nil {
		log.Errorf("retrieving user keys of app id %v: %+v", appID, err)
		return false
	}

	marshaledKey := key.Marshal()
	return lo.ContainsBy(keys, func(userKey *domain.UserKey) bool {
		return bytes.Equal(userKey.MarshalKey(), marshaledKey)
	})
}

func (b *k8sBackend) sshHandler(s ssh.Session) {
	writeErrAndClose := func(err error) {
		log.Errorf("%+v", err)
		_, _ = s.Write([]byte(err.Error() + "\n"))
		_ = s.Exit(1)
	}

	appID := s.User()
	log.Infof("new ssh connection into app %s", appID)

	app, err := b.appRepo.GetApplication(s.Context(), appID)
	if err != nil {
		writeErrAndClose(errors.Wrapf(err, "retrieving app with id %v", appID))
		return
	}

	_, _ = s.Write([]byte(fmt.Sprintf("Welcome to NeoShowcase! Connecting to application %s (id: %s) ...\n", app.Name, appID)))

	err = b.exec(s, app)
	if err != nil {
		writeErrAndClose(err)
		return
	}
	_ = s.Exit(0)
}

func (b *k8sBackend) exec(s ssh.Session, app *domain.Application) error {
	cmd := s.Command()
	if len(cmd) == 0 {
		cmd = []string{"/bin/sh"}
	}
	req := b.client.CoreV1().RESTClient().Post().
		Resource("pods").Name(generatedPodName(app.ID)).
		Namespace(b.config.Namespace).SubResource("exec")
	option := &v1.PodExecOptions{
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
		Container: podContainerName,
	}
	req.VersionedParams(option, scheme.ParameterCodec)
	ex, err := remotecommand.NewSPDYExecutor(b.restConfig, "POST", req.URL())
	if err != nil {
		return errors.Wrap(err, "creating SPDY executor")
	}
	err = ex.StreamWithContext(s.Context(), remotecommand.StreamOptions{
		Stdin:  s,
		Stdout: s,
		Stderr: s.Stderr(),
		Tty:    true,
	})
	if err != nil {
		return errors.Wrap(err, "streaming")
	}
	return nil
}
