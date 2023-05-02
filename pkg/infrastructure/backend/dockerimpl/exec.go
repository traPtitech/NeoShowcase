package dockerimpl

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/friendsofgo/errors"
	"golang.org/x/sync/errgroup"
)

func (b *dockerBackend) ExecContainer(ctx context.Context, appID string, cmd []string, stdin io.Reader, stdout, stderr io.Writer) error {
	execConf := types.ExecConfig{
		Tty:          true,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		WorkingDir:   "/srv",
		Cmd:          cmd,
	}
	execID, err := b.c.ContainerExecCreate(ctx, containerName(appID), execConf)
	if err != nil {
		return errors.Wrap(err, "creating exec")
	}

	ex, err := b.c.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return errors.Wrap(err, "attaching exec process")
	}

	ctx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		defer cancel()
		_, err := io.Copy(ex.Conn, stdin)
		if err != nil {
			return errors.Wrap(err, "writing into exec conn")
		}
		return nil
	})
	eg.Go(func() error {
		defer cancel()
		_, err := stdcopy.StdCopy(stdout, stderr, ex.Reader)
		if err != nil {
			return errors.Wrap(err, "reading exec response")
		}
		return nil
	})
	eg.Go(func() error {
		<-ctx.Done()
		ex.Close()
		return ex.CloseWrite()
	})
	return eg.Wait()
}
