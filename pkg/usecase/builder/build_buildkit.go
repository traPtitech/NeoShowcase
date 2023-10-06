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
	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

const (
	buildScriptName = "neoshowcase_internal_build.sh"
)

func withBuildkitProgress(ctx context.Context, logger io.Writer, buildFn func(ctx context.Context, ch chan *buildkit.SolveStatus) error) error {
	ch := make(chan *buildkit.SolveStatus)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return buildFn(ctx, ch)
	})
	eg.Go(func() error {
		// TODO: VertexWarningを使う (LLBのどのvertexに問題があったか)
		// NOTE: https://github.com/moby/buildkit/pull/1721#issuecomment-703937866
		// DisplaySolveStatus's context should not be cancelled, in order to receive 'cancelled' events from buildkit API call.
		_, err := progressui.DisplaySolveStatus(context.WithoutCancel(ctx), nil, logger, ch)
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

func (s *builderService) authSessions() []session.Attachable {
	if s.config.Registry.Username == "" && s.config.Registry.Password == "" {
		return nil
	}
	return []session.Attachable{authprovider.NewDockerAuthProvider(&configfile.ConfigFile{
		AuthConfigs: map[string]types.AuthConfig{
			s.config.Registry.Addr: {
				Username: s.config.Registry.Username,
				Password: s.config.Registry.Password,
			},
		},
	})}
}

func (s *builderService) solveDockerfile(
	ctx context.Context,
	dest string,
	contextDir string,
	dockerfileDir, dockerfileName string,
	env map[string]string,
	ch chan *buildkit.SolveStatus,
) error {
	opts := buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type: buildkit.ExporterImage,
			Attrs: map[string]string{
				"name": dest,
				"push": "true",
			},
		}},
		LocalDirs: map[string]string{
			"context":    contextDir,
			"dockerfile": dockerfileDir,
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
	_, err := s.buildkit.Solve(ctx, nil, opts, ch)
	return err
}

func (s *builderService) buildRuntimeCmd(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigRuntimeCmd,
) error {
	var dockerfile strings.Builder
	if bc.BaseImage == "" {
		dockerfile.WriteString("FROM scratch\n")
	} else {
		dockerfile.WriteString(fmt.Sprintf("FROM %v\n", bc.BaseImage))
	}

	env, err := s.appEnv(ctx, st.app)
	if err != nil {
		return err
	}
	for key := range env {
		dockerfile.WriteString(fmt.Sprintf("ARG %v\n", key))
		dockerfile.WriteString(fmt.Sprintf("ENV %v=$%v\n", key, key))
	}

	dockerfile.WriteString("WORKDIR /srv\n")
	dockerfile.WriteString("COPY . .\n")

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
		env,
		ch,
	)
}

func (s *builderService) buildRuntimeDockerfile(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigRuntimeDockerfile,
) error {
	env, err := s.appEnv(ctx, st.app)
	if err != nil {
		return err
	}
	contextDir := lo.Ternary(bc.Context != "", bc.Context, ".")
	return s.solveDockerfile(
		ctx,
		s.destImage(st.app, st.build),
		filepath.Join(st.repositoryTempDir, contextDir),
		filepath.Join(st.repositoryTempDir, contextDir),
		bc.DockerfileName,
		env,
		ch,
	)
}

func (s *builderService) buildStaticCmd(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigStaticCmd,
) error {
	var dockerfile strings.Builder

	dockerfile.WriteString(fmt.Sprintf(
		"FROM %s\n",
		lo.Ternary(bc.BaseImage == "", "scratch", bc.BaseImage),
	))

	env, err := s.appEnv(ctx, st.app)
	if err != nil {
		return err
	}
	for key := range env {
		dockerfile.WriteString(fmt.Sprintf("ARG %v\n", key))
		dockerfile.WriteString(fmt.Sprintf("ENV %v=$%v\n", key, key))
	}

	dockerfile.WriteString("WORKDIR /srv\n")
	dockerfile.WriteString("COPY . .\n")

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
		env,
		ch,
	)
}

func (s *builderService) buildStaticDockerfile(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigStaticDockerfile,
) error {
	env, err := s.appEnv(ctx, st.app)
	if err != nil {
		return err
	}
	contextDir := lo.Ternary(bc.Context != "", bc.Context, ".")
	st.staticDest = bc.ArtifactPath
	return s.solveDockerfile(
		ctx,
		s.tmpDestImage(st.app, st.build),
		filepath.Join(st.repositoryTempDir, contextDir),
		filepath.Join(st.repositoryTempDir, contextDir),
		bc.DockerfileName,
		env,
		ch,
	)
}
