package builder

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type logWriter struct {
	buildID string
	send    chan<- []byte
	buf     bytes.Buffer
}

func newLogWriter(buildID string, client domain.ControllerBuilderServiceClient) *logWriter {
	send := make(chan []byte, 50)
	go func() {
		err := client.StreamBuildLog(context.Background(), buildID, send)
		if err != nil {
			log.Errorf("sending build log: %+v", err)
		}
	}()
	return &logWriter{
		buildID: buildID,
		send:    send,
		buf:     bytes.Buffer{},
	}
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	n, err = l.buf.Write(p)
	if err != nil {
		return
	}
	// Make a copy to send to another goroutine
	// From io.Writer's document: Implementations must not retain p.
	toSend := make([]byte, len(p))
	copy(toSend, p)
	// Non-blocking send may drop logs - complete log is to be sent on finishing build
	select {
	case l.send <- toSend:
	default:
	}
	return
}

func (l *logWriter) Close() error {
	close(l.send)
	return nil
}

type state struct {
	app       *domain.Application
	envs      []*domain.Environment
	build     *domain.Build
	repo      *domain.Repository
	logWriter *logWriter

	repositoryTempDir string
	artifactTempFile  *os.File
	done              chan struct{}

	staticDest string
}

func newState(app *domain.Application, envs []*domain.Environment, build *domain.Build, repo *domain.Repository, client domain.ControllerBuilderServiceClient) (*state, error) {
	st := &state{
		app:       app,
		envs:      envs,
		build:     build,
		repo:      repo,
		logWriter: newLogWriter(build.ID, client),
		done:      make(chan struct{}),
	}
	var err error
	st.repositoryTempDir, err = os.MkdirTemp("", "repository-")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tmp repository dir")
	}
	st.artifactTempFile, err = os.CreateTemp("", "artifacts-")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tmp artifact file")
	}
	return st, nil
}

func (s *state) appEnv() map[string]string {
	return lo.SliceToMap(s.envs, (*domain.Environment).GetKV)
}

func (st *state) deployType() domain.DeployType {
	return st.app.Config.BuildConfig.BuildType().DeployType()
}

func (s *state) Done() {
	_ = s.logWriter.Close()
	_ = os.RemoveAll(s.repositoryTempDir)
	_ = os.Remove(s.artifactTempFile.Name())
	close(s.done)
}

func (s *state) Wait() {
	<-s.done
}

func (s *state) Logger() io.Writer {
	return s.logWriter
}

func (s *state) WriteLog(a ...interface{}) {
	_, _ = fmt.Fprintln(s.logWriter, a...)
}
