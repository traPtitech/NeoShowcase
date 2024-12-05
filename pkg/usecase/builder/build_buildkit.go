package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	"github.com/friendsofgo/errors"
	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/tonistiigi/fsutil"
	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

const (
	buildScriptName           = "neoshowcase_internal_build.sh"
	temporaryDockerignoreName = "neoshowcase_temporary_dockerignore"
)

func withBuildkitProgress(ctx context.Context, logger io.Writer, buildFn func(ctx context.Context, ch chan *buildkit.SolveStatus) error) error {
	ch := make(chan *buildkit.SolveStatus)
	disp, err := progressui.NewDisplay(logger, progressui.PlainMode)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return buildFn(ctx, ch)
	})
	eg.Go(func() error {
		// TODO: VertexWarningを使う (LLBのどのvertexに問題があったか)
		// NOTE: https://github.com/moby/buildkit/pull/1721#issuecomment-703937866
		// progress-ui context should not be cancelled, in order to receive 'cancelled' events from buildkit API call.
		_, err := disp.UpdateFrom(context.WithoutCancel(ctx), ch)
		return err
	})

	return eg.Wait()
}

func createTempFile(pattern string, content string) (name string, cleanup func(), err error) {
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", nil, errors.Wrap(err, "creating temp "+pattern+" file")
	}
	defer f.Close()
	cleanup = func() {
		err := os.Remove(f.Name())
		if err != nil {
			log.Errorf("removing temp file "+f.Name()+": %+v", err)
		}
	}
	_, err = f.WriteString(content)
	if err != nil {
		cleanup()
		return "", nil, errors.Wrap(err, "writing to temp file "+f.Name())
	}
	return f.Name(), cleanup, nil
}

func createFile(filename string, content string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func createScriptFile(filename string, script string) error {
	return createFile(filename, "#!/bin/sh\nset -eux\n"+script)
}

func dockerignoreExists(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, ".dockerignore"))
	return err == nil && !info.IsDir()
}

func (s *ServiceImpl) authSessions() []session.Attachable {
	if s.imageConfig.Registry.Username == "" && s.imageConfig.Registry.Password == "" {
		return nil
	}
	return []session.Attachable{authprovider.NewDockerAuthProvider(&configfile.ConfigFile{
		AuthConfigs: map[string]types.AuthConfig{
			s.imageConfig.Registry.Addr: {
				Username: s.imageConfig.Registry.Username,
				Password: s.imageConfig.Registry.Password,
			},
		},
	}, nil)}
}

func (s *ServiceImpl) solveDockerfile(
	ctx context.Context,
	dest string,
	contextDir string,
	dockerfileDir, dockerfileName string,
	env map[string]string,
	ch chan *buildkit.SolveStatus,
) error {
	// ch must be closed when this function returns because it is listened by progress ui display
	var channelClosed bool = false
	defer func() {
		if !channelClosed {
			close(ch)
		}
	}()

	ctxMount, err := fsutil.NewFS(contextDir)
	if err != nil {
		return errors.Wrap(err, "invalid context mount dir")
	}
	dockerfileMount, err := fsutil.NewFS(dockerfileDir)
	if err != nil {
		return errors.Wrap(err, "invalid dockerfile mount dir")
	}

	opts := buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type: buildkit.ExporterImage,
			Attrs: map[string]string{
				"name": dest,
				"push": "true",
			},
		}},
		LocalMounts: map[string]fsutil.FS{
			"context":    ctxMount,
			"dockerfile": dockerfileMount,
		},
		Frontend: "dockerfile.v0",
		FrontendAttrs: ds.MergeMap(
			map[string]string{"filename": dockerfileName},
			lo.MapEntries(env, func(k string, v string) (string, string) {
				return "build-arg:" + k, v
			}),
		),
		Session: s.authSessions(),
	}

	_, err = s.buildkit.Solve(ctx, nil, opts, ch)
	// ch is closed by buildkit.Solve
	channelClosed = true

	return err
}

func (s *ServiceImpl) buildRuntimeCmd(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigRuntimeCmd,
) error {
	// If .dockerignore exists, rename to prevent it from being picked up by buildkitd,
	// as this is not the behavior we want in 'Command' build which is supposed to execute commands against raw repository files.
	// See https://github.com/traPtitech/NeoShowcase/issues/877 for more details.
	dockerignoreExists := dockerignoreExists(st.repositoryTempDir)
	if dockerignoreExists {
		err := os.Rename(filepath.Join(st.repositoryTempDir, ".dockerignore"), filepath.Join(st.repositoryTempDir, temporaryDockerignoreName))
		if err != nil {
			return errors.Wrap(err, "renaming .dockerignore")
		}
	}

	var dockerfile strings.Builder
	if bc.BaseImage == "" {
		dockerfile.WriteString("FROM scratch\n")
	} else {
		dockerfile.WriteString(fmt.Sprintf("FROM %v\n", bc.BaseImage))
	}

	for key := range st.appEnv() {
		dockerfile.WriteString(fmt.Sprintf("ARG %v\n", key))
		dockerfile.WriteString(fmt.Sprintf("ENV %v=$%v\n", key, key))
	}

	dockerfile.WriteString("WORKDIR /srv\n")
	dockerfile.WriteString("COPY . .\n")
	if dockerignoreExists {
		dockerfile.WriteString(fmt.Sprintf("RUN mv %s .dockerignore\n", temporaryDockerignoreName))
	}

	if bc.BuildCmd != "" {
		err := createScriptFile(filepath.Join(st.repositoryTempDir, buildScriptName), bc.BuildCmd)
		if err != nil {
			return err
		}
		dockerfile.WriteString(fmt.Sprintf("RUN ./%v\n", buildScriptName))
		dockerfile.WriteString(fmt.Sprintf("RUN rm ./%v\n", buildScriptName))
	}

	tmpName, cleanup, err := createTempFile("dockerfile-*", dockerfile.String())
	if err != nil {
		return err
	}
	defer cleanup()

	return s.solveDockerfile(
		ctx,
		s.destImage(st.app, st.build),
		st.repositoryTempDir,
		filepath.Dir(tmpName),
		filepath.Base(tmpName),
		st.appEnv(),
		ch,
	)
}

func (s *ServiceImpl) buildRuntimeDockerfile(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigRuntimeDockerfile,
) error {
	contextDir := lo.Ternary(bc.Context != "", bc.Context, ".")
	return s.solveDockerfile(
		ctx,
		s.destImage(st.app, st.build),
		filepath.Join(st.repositoryTempDir, contextDir),
		filepath.Join(st.repositoryTempDir, contextDir),
		bc.DockerfileName,
		st.appEnv(),
		ch,
	)
}

func (s *ServiceImpl) buildStaticCmd(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigStaticCmd,
) error {
	// If .dockerignore exists, rename to prevent it from being picked up by buildkitd,
	// as this is not the behavior we want in 'Command' build which is supposed to execute commands against raw repository files.
	// See https://github.com/traPtitech/NeoShowcase/issues/877 for more details.
	dockerignoreExists := dockerignoreExists(st.repositoryTempDir)
	if dockerignoreExists {
		err := os.Rename(filepath.Join(st.repositoryTempDir, ".dockerignore"), filepath.Join(st.repositoryTempDir, temporaryDockerignoreName))
		if err != nil {
			return errors.Wrap(err, "renaming .dockerignore")
		}
	}

	var dockerfile strings.Builder

	dockerfile.WriteString(fmt.Sprintf(
		"FROM %s\n",
		lo.Ternary(bc.BaseImage == "", "scratch", bc.BaseImage),
	))

	for key := range st.appEnv() {
		dockerfile.WriteString(fmt.Sprintf("ARG %v\n", key))
		dockerfile.WriteString(fmt.Sprintf("ENV %v=$%v\n", key, key))
	}

	dockerfile.WriteString("WORKDIR /srv\n")
	dockerfile.WriteString("COPY . .\n")
	if dockerignoreExists {
		dockerfile.WriteString(fmt.Sprintf("RUN mv %s .dockerignore\n", temporaryDockerignoreName))
	}

	if bc.BuildCmd != "" {
		err := createScriptFile(filepath.Join(st.repositoryTempDir, buildScriptName), bc.BuildCmd)
		if err != nil {
			return err
		}
		dockerfile.WriteString("RUN ./" + buildScriptName + "\n")
		dockerfile.WriteString("RUN rm ./" + buildScriptName + "\n")
	}

	tmpName, cleanup, err := createTempFile("dockerfile-*", dockerfile.String())
	if err != nil {
		return err
	}
	defer cleanup()

	st.staticDest = filepath.Join("/srv", bc.ArtifactPath)
	return s.solveDockerfile(
		ctx,
		s.tmpDestImage(st.app, st.build),
		st.repositoryTempDir,
		filepath.Dir(tmpName),
		filepath.Base(tmpName),
		st.appEnv(),
		ch,
	)
}

func (s *ServiceImpl) buildStaticDockerfile(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigStaticDockerfile,
) error {
	contextDir := lo.Ternary(bc.Context != "", bc.Context, ".")
	st.staticDest = bc.ArtifactPath
	return s.solveDockerfile(
		ctx,
		s.tmpDestImage(st.app, st.build),
		filepath.Join(st.repositoryTempDir, contextDir),
		filepath.Join(st.repositoryTempDir, contextDir),
		bc.DockerfileName,
		st.appEnv(),
		ch,
	)
}
