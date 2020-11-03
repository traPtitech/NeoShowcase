package builder

import (
	"context"
	"fmt"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/progress/progressui"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
)

type BuildkitWrapper struct {
	client *client.Client
}

type BuildImageArgs struct {
	ImageName  string
	ContextDir string
}

func (bw *BuildkitWrapper) BuildImage(ctx context.Context, args BuildImageArgs, logWriter io.Writer) error {
	if logWriter == nil {
		logWriter = ioutil.Discard // ログを破棄
	}

	ch := make(chan *client.SolveStatus)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// buildkitdでDockerfileを実行
		_, err := bw.client.Solve(ctx, nil, client.SolveOpt{
			Exports: []client.ExportEntry{{
				Type: "image",
				Attrs: map[string]string{
					"name": args.ImageName,
					"push": "true",
				},
			}},
			LocalDirs: map[string]string{
				"context":    args.ContextDir,
				"dockerfile": args.ContextDir,
			},
			Frontend:      "dockerfile.v0",
			FrontendAttrs: map[string]string{"filename": "Dockerfile"},
		}, ch)
		return err
	})
	eg.Go(func() error {
		// ビルドログを収集
		return progressui.DisplaySolveStatus(ctx, "", nil, logWriter, ch)
	})

	return eg.Wait()
}

type BuildStaticArgs struct {
	Output     io.WriteCloser
	ContextDir string
	LLB        *llb.Definition
}

func (bw *BuildkitWrapper) BuildStatic(ctx context.Context, args BuildStaticArgs, logWriter io.Writer) error {
	if logWriter == nil {
		logWriter = ioutil.Discard // ログを破棄
	}
	if args.Output == nil {
		return fmt.Errorf("no output")
	}
	if args.LLB == nil {
		return fmt.Errorf("no llb")
	}

	ch := make(chan *client.SolveStatus)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// buildkitdでDockerfileを実行
		_, err := bw.client.Solve(ctx, args.LLB, client.SolveOpt{
			Exports: []client.ExportEntry{{
				Type:   "tar",
				Output: func(_ map[string]string) (io.WriteCloser, error) { return args.Output, nil },
			}},
			LocalDirs: map[string]string{"local-src": args.ContextDir},
		}, ch)
		return err
	})
	eg.Go(func() error {
		// ビルドログを収集
		return progressui.DisplaySolveStatus(ctx, "", nil, logWriter, ch)
	})

	return eg.Wait()
}
