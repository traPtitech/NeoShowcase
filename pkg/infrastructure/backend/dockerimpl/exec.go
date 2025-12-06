package dockerimpl

import (
	"context"
	"io"

	"github.com/friendsofgo/errors"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
	"golang.org/x/sync/errgroup"
)

func streamHijackedResp(ctx context.Context, res client.HijackedResponse, stdin io.Reader, stdout, stderr io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		defer cancel()
		_, err := io.Copy(res.Conn, stdin)
		if err != nil {
			return errors.Wrap(err, "writing into exec conn")
		}
		return nil
	})
	eg.Go(func() error {
		defer cancel()
		_, err := stdcopy.StdCopy(stdout, stderr, res.Reader)
		if err != nil {
			return errors.Wrap(err, "reading exec response")
		}
		return nil
	})
	eg.Go(func() error {
		<-ctx.Done()
		res.Close()
		return res.CloseWrite()
	})
	return eg.Wait()
}

func (b *Backend) AttachContainer(ctx context.Context, appID string, stdin io.Reader, stdout, stderr io.Writer) error {
	res, err := b.c.ContainerAttach(ctx, containerName(appID), client.ContainerAttachOptions{
		Stream:     true,
		Stdin:      true,
		Stdout:     true,
		Stderr:     true,
		DetachKeys: "",
		Logs:       true,
	})
	if err != nil {
		return errors.Wrap(err, "attaching to container")
	}
	return streamHijackedResp(ctx, res.HijackedResponse, stdin, stdout, stderr)
}

func (b *Backend) ExecContainer(ctx context.Context, appID string, cmd []string, stdin io.Reader, stdout, stderr io.Writer) error {
	execConf := client.ExecCreateOptions{
		TTY:          true,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		WorkingDir:   "/srv",
		Cmd:          cmd,
	}
	execID, err := b.c.ExecCreate(ctx, containerName(appID), execConf)
	if err != nil {
		return errors.Wrap(err, "creating exec")
	}

	res, err := b.c.ExecAttach(ctx, execID.ID, client.ExecAttachOptions{})
	if err != nil {
		return errors.Wrap(err, "attaching exec process")
	}

	return streamHijackedResp(ctx, res.HijackedResponse, stdin, stdout, stderr)
}
