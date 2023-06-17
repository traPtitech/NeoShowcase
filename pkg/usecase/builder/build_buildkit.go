package builder

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	"github.com/friendsofgo/errors"
	"github.com/mattn/go-shellwords"
	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

const (
	buildScriptName   = "neoshowcase_internal_build.sh"
	tmpDockerfileName = "neoshowcase_internal_dockerfile"
)

func withBuildkitProgress(ctx context.Context, logger io.Writer, buildFn func(ctx context.Context, ch chan *buildkit.SolveStatus) error) error {
	ch := make(chan *buildkit.SolveStatus)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return buildFn(ctx, ch)
	})
	eg.Go(func() error {
		// TODO: VertexWarningを使う (LLBのどのvertexに問題があったか)
		_, err := progressui.DisplaySolveStatus(ctx, "", nil, logger, ch)
		return err
	})
	return eg.Wait()
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

func toJSONExecFormat(args []string) string {
	return strings.Join(ds.Map(args, func(s string) string {
		return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}), ", ")
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

func (s *builderService) buildRuntimeCmd(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigRuntimeCmd,
) error {
	var dockerfile bytes.Buffer
	if bc.BaseImage == "" {
		dockerfile.WriteString("FROM scratch\n")
	} else {
		dockerfile.WriteString(fmt.Sprintf("FROM %v\n", bc.BaseImage))
	}
	dockerfile.WriteString("WORKDIR /srv\n")
	dockerfile.WriteString("COPY . .\n")

	if bc.BuildCmd != "" {
		if bc.BuildCmdShell {
			err := createScriptFile(filepath.Join(st.repositoryTempDir, buildScriptName), bc.BuildCmd)
			if err != nil {
				return err
			}
			dockerfile.WriteString(fmt.Sprintf("RUN ./%v\n", buildScriptName))
		} else {
			args, err := shellwords.Parse(bc.BuildCmd)
			if err != nil {
				return err
			}
			dockerfile.WriteString(fmt.Sprintf("RUN [%v]\n", toJSONExecFormat(args)))
		}
	}

	err := createFile(filepath.Join(st.repositoryTempDir, tmpDockerfileName), dockerfile.String())
	if err != nil {
		return err
	}
	_, err = s.buildkit.Solve(ctx, nil, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type: buildkit.ExporterImage,
			Attrs: map[string]string{
				"name": s.destImage(st.app, st.build),
				"push": "true",
			},
		}},
		LocalDirs: map[string]string{
			"context":    st.repositoryTempDir,
			"dockerfile": st.repositoryTempDir,
		},
		Frontend:      "dockerfile.v0",
		FrontendAttrs: map[string]string{"filename": tmpDockerfileName},
		Session:       s.authSessions(),
	}, ch)
	return err
}

func (s *builderService) buildRuntimeDockerfile(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigRuntimeDockerfile,
) error {
	contextDir := lo.Ternary(bc.Context != "", bc.Context, ".")
	_, err := s.buildkit.Solve(ctx, nil, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type: buildkit.ExporterImage,
			Attrs: map[string]string{
				"name": s.destImage(st.app, st.build),
				"push": "true",
			},
		}},
		LocalDirs: map[string]string{
			"context":    filepath.Join(st.repositoryTempDir, contextDir),
			"dockerfile": filepath.Join(st.repositoryTempDir, contextDir),
		},
		Frontend:      "dockerfile.v0",
		FrontendAttrs: map[string]string{"filename": bc.DockerfileName},
		Session:       s.authSessions(),
	}, ch)
	return err
}

func (s *builderService) buildStaticCmd(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigStaticCmd,
) error {
	var ls llb.State
	if bc.BaseImage == "" {
		ls = llb.Scratch()
	} else {
		ls = llb.Image(bc.BaseImage)
	}
	ls = ls.
		Dir("/srv").
		File(llb.Copy(llb.Local("local-src"), ".", ".", &llb.CopyInfo{
			CopyDirContentsOnly: true,
			AllowWildcard:       true,
			CreateDestPath:      true,
		}))

	if bc.BuildCmd != "" {
		if bc.BuildCmdShell {
			err := createScriptFile(filepath.Join(st.repositoryTempDir, buildScriptName), bc.BuildCmd)
			if err != nil {
				return err
			}
			ls = ls.Run(llb.Args([]string{"./" + buildScriptName})).Root()
		} else {
			args, err := shellwords.Parse(bc.BuildCmd)
			if err != nil {
				return err
			}
			ls = ls.Run(llb.Args(args)).Root()
		}
	}

	// ビルドで生成された静的ファイルのみを含むScratchイメージを構成
	def, err := llb.
		Scratch().
		File(llb.Copy(ls, filepath.Join("/srv", bc.ArtifactPath), "/", &llb.CopyInfo{
			CopyDirContentsOnly: true,
			CreateDestPath:      true,
			AllowWildcard:       true,
		})).
		Marshal(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to marshal llb")
	}

	_, err = s.buildkit.Solve(ctx, def, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:   buildkit.ExporterTar,
			Output: func(_ map[string]string) (io.WriteCloser, error) { return st.artifactTempFile, nil },
		}},
		LocalDirs: map[string]string{
			"local-src": st.repositoryTempDir,
		},
		Session: s.authSessions(),
	}, ch)
	return err
}

func (s *builderService) buildStaticDockerfile(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigStaticDockerfile,
) error {
	contextDir := lo.Ternary(bc.Context != "", bc.Context, ".")
	dockerfile, err := os.ReadFile(filepath.Join(st.repositoryTempDir, contextDir, bc.DockerfileName))
	if err != nil {
		return err
	}

	b, _, _, err := dockerfile2llb.Dockerfile2LLB(ctx, dockerfile, dockerfile2llb.ConvertOpt{})
	if err != nil {
		return err
	}

	def, err := llb.
		Scratch().
		File(llb.Copy(*b, bc.ArtifactPath, "/", &llb.CopyInfo{
			CopyDirContentsOnly: true,
			CreateDestPath:      true,
			AllowWildcard:       true,
		})).
		Marshal(ctx)
	if err != nil {
		return err
	}

	_, err = s.buildkit.Solve(ctx, def, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:   buildkit.ExporterTar,
			Output: func(_ map[string]string) (io.WriteCloser, error) { return st.artifactTempFile, nil },
		}},
		LocalDirs: map[string]string{
			"context": filepath.Join(st.repositoryTempDir, contextDir),
		},
		Session: s.authSessions(),
	}, ch)
	return err
}
