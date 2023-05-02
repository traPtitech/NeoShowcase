package k8simpl

import (
	"context"
	"io"

	"github.com/friendsofgo/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func (b *k8sBackend) ExecContainer(ctx context.Context, appID string, cmd []string, stdin io.Reader, stdout, stderr io.Writer) error {
	req := b.client.CoreV1().RESTClient().Post().
		Resource("pods").Name(generatedPodName(appID)).
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
	err = ex.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    true,
	})
	if err != nil {
		return errors.Wrap(err, "streaming")
	}
	return nil
}
