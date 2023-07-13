package k8simpl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config/configfile"
	types2 "github.com/docker/cli/cli/config/types"
	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

type k8sBackend struct {
	c          *kubernetes.Clientset
	restConfig *rest.Config
	config     Config
	image      builder.ImageConfig
}

func NewBuildpackBackend(
	config Config,
	image builder.ImageConfig,
) (builder.BuildpackBackend, error) {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	c, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	b := &k8sBackend{
		c:          c,
		restConfig: restConfig,
		config:     config,
		image:      image,
	}
	err = b.prepareAuth()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (k *k8sBackend) dockerAuth() (s string, ok bool) {
	if k.image.Registry.Username == "" && k.image.Registry.Password == "" {
		return "", false
	}
	c := configfile.ConfigFile{
		AuthConfigs: map[string]types2.AuthConfig{
			k.image.Registry.Addr: {
				Username: k.image.Registry.Username,
				Password: k.image.Registry.Password,
			},
		},
	}
	b, _ := json.Marshal(&c)
	return string(b), true
}

func escapeSingleQuote(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}

func (k *k8sBackend) prepareAuth() error {
	auth, ok := k.dockerAuth()
	if ok {
		err := k.exec(context.Background(), "/", "mkdir -p ~/.docker", nil, io.Discard, io.Discard)
		if err != nil {
			return errors.Wrap(err, "making ~/.docker directory")
		}
		err = k.exec(context.Background(), "/", fmt.Sprintf(`echo '%s' > ~/.docker/config.json`, escapeSingleQuote(auth)), nil, io.Discard, io.Discard)
		if err != nil {
			return errors.Wrap(err, "writing ~/.docker/config.json to builder")
		}
	}
	return nil
}

func (k *k8sBackend) exec(ctx context.Context, workDir string, cmd string, env map[string]string, stdout io.Writer, stderr io.Writer) error {
	req := k.c.CoreV1().RESTClient().Post().Resource("pods").Name(k.config.PodName).
		Namespace(k.config.Namespace).SubResource("exec")
	var shCmds []string
	for k, v := range env {
		shCmds = append(shCmds, fmt.Sprintf("export %v=\"%v\"", k, strings.ReplaceAll(v, `"`, `\"`)))
	}
	shCmds = append(shCmds,
		"cd "+workDir,
		cmd)
	option := &v1.PodExecOptions{
		Command:   []string{"sh", "-c", strings.Join(shCmds, " && ")},
		Stdout:    true,
		Stderr:    true,
		Container: k.config.ContainerName,
	}
	req.VersionedParams(option, scheme.ParameterCodec)
	ex, err := remotecommand.NewSPDYExecutor(k.restConfig, "POST", req.URL())
	if err != nil {
		return err
	}
	err = ex.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		return err
	}
	return nil
}

func (k *k8sBackend) prepareDir(ctx context.Context, localName, remoteName string) (remotePath string, cleanup func(), err error) {
	localVolumePath := filepath.Join(k.config.LocalDir, remoteName)
	remotePath = filepath.Join(k.config.RemoteDir, remoteName)
	err = exec.CommandContext(ctx, "cp", "-r", localName, localVolumePath).Run()
	if err != nil {
		return "", nil, errors.Wrap(err, "copying files")
	}
	cleanup = func() {
		err := exec.Command("rm", "-r", localVolumePath).Run()
		if err != nil {
			log.Errorf("failed to rm tmp dir: %+v", err)
		}
	}

	err = exec.CommandContext(ctx, "chown", "-R", fmt.Sprintf("%d:%d", k.config.User, k.config.Group), localVolumePath).Run()
	if err != nil {
		cleanup()
		return "", nil, errors.Wrap(err, "setting remote dir owner")
	}
	return remotePath, cleanup, nil
}

func (k *k8sBackend) Pack(
	ctx context.Context,
	repoDir string,
	imageDest string,
	env map[string]string,
	logWriter io.Writer,
) (path string, err error) {
	// NOTE: safe to use the same remote name between builds, under the assumption that buildpack pod is not shared
	remoteRepoPath, cleanupRepoDir, err := k.prepareDir(ctx, repoDir, "ns-repo")
	if err != nil {
		return "", err
	}
	defer cleanupRepoDir()

	localEnvTmp, err := os.MkdirTemp("", "env-")
	if err != nil {
		return "", errors.Wrap(err, "creating env temp dir")
	}
	defer os.RemoveAll(localEnvTmp)
	err = os.Mkdir(filepath.Join(localEnvTmp, "env"), 0700)
	if err != nil {
		return "", errors.Wrap(err, "creating platform env dir")
	}
	for k, v := range env {
		err = os.WriteFile(filepath.Join(localEnvTmp, "env", k), []byte(v), 0600)
		if err != nil {
			return "", errors.Wrap(err, "creating env file")
		}
	}
	remoteEnvPath, cleanupEnvDir, err := k.prepareDir(ctx, localEnvTmp, "ns-env")
	if err != nil {
		return "", err
	}
	defer cleanupEnvDir()

	// TODO: support pushing to insecure registry for local development
	// https://github.com/buildpacks/lifecycle/issues/524
	// https://github.com/buildpacks/rfcs/blob/main/text/0111-support-insecure-registries.md
	// Workaround: use registry host "*.local" to allow google/go-containerregistry to detect as http protocol
	// see: https://github.com/traPtitech/NeoShowcase/issues/493
	err = k.exec(ctx,
		remoteRepoPath,
		fmt.Sprintf("/cnb/lifecycle/creator -skip-restore -platform=%s -app=. %s", remoteEnvPath, imageDest),
		ds.MergeMap(env, map[string]string{"CNB_PLATFORM_API": k.config.PlatformAPI}),
		logWriter, logWriter)
	if err != nil {
		return "", err
	}
	return remoteRepoPath, nil
}
