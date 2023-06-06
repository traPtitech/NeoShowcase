package builder

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

type logWriter struct {
	buildID  string
	response chan<- *pb.BuilderResponse
	buf      bytes.Buffer
}

func (l *logWriter) toBuilderResponse(p []byte) *pb.BuilderResponse {
	return &pb.BuilderResponse{Type: pb.BuilderResponse_BUILD_LOG, Body: &pb.BuilderResponse_Log{
		Log: &pb.BuildLogPortion{BuildId: l.buildID, Log: p},
	}}
}

func (l *logWriter) LogReader() io.Reader {
	return &l.buf
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	n, err = l.buf.Write(p)
	if err != nil {
		return
	}
	select {
	case l.response <- l.toBuilderResponse(p):
	default:
	}
	return
}

type state struct {
	app       *domain.Application
	build     *domain.Build
	repo      *domain.Repository
	logWriter *logWriter

	repositoryTempDir string
	artifactTempFile  *os.File
	done              chan struct{}
}

func newState(app *domain.Application, build *domain.Build, repo *domain.Repository, response chan<- *pb.BuilderResponse) (*state, error) {
	st := &state{
		app:   app,
		build: build,
		repo:  repo,
		logWriter: &logWriter{
			buildID:  build.ID,
			response: response,
		},
		done: make(chan struct{}),
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

func (s *state) Done() {
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
