package dockerimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config/configfile"
	types2 "github.com/docker/cli/cli/config/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/tarfs"
)

type dockerBackend struct {
	c      *client.Client
	config Config
	image  builder.ImageConfig
}

func NewBuildpackBackend(
	config Config,
	image builder.ImageConfig,
) (builder.BuildpackBackend, error) {
	c, err := client.NewClientWithOpts(
		client.FromEnv,
		// Using github.com/moby/moby of v25 master@8e51b8b59cb8 (2023-07-18), required by github.com/moby/buildkit@v0.12.2,
		// defaults to API version 1.44, which currently available docker installation does not support.
		client.WithVersion("1.43"),
	)
	if err != nil {
		return nil, err
	}
	b := &dockerBackend{
		c:      c,
		config: config,
		image:  image,
	}
	err = b.prepareAuth()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (d *dockerBackend) dockerAuth() (s string, ok bool) {
	if d.image.Registry.Username == "" && d.image.Registry.Password == "" {
		return "", false
	}
	c := configfile.ConfigFile{
		AuthConfigs: map[string]types2.AuthConfig{
			d.image.Registry.Addr: {
				Username: d.image.Registry.Username,
				Password: d.image.Registry.Password,
			},
		},
	}
	b, _ := json.Marshal(&c)
	return string(b), true
}

func escapeSingleQuote(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}

func (d *dockerBackend) prepareAuth() error {
	auth, ok := d.dockerAuth()
	if ok {
		err := d.exec(context.Background(), "/", []string{"sh", "-c", "mkdir -p ~/.docker"}, nil, io.Discard, io.Discard)
		if err != nil {
			return errors.Wrap(err, "making ~/.docker directory")
		}
		err = d.exec(context.Background(), "/", []string{"sh", "-c", fmt.Sprintf(`echo '%s' > ~/.docker/config.json`, escapeSingleQuote(auth))}, nil, io.Discard, io.Discard)
		if err != nil {
			return errors.Wrap(err, "writing ~/.docker/config.json to builder")
		}
	}
	return nil
}

func (d *dockerBackend) execRoot(ctx context.Context, workDir string, cmd []string, env map[string]string, outWriter, errWriter io.Writer) error {
	return d._exec(ctx, true, workDir, cmd, env, outWriter, errWriter)
}

func (d *dockerBackend) exec(ctx context.Context, workDir string, cmd []string, env map[string]string, outWriter, errWriter io.Writer) error {
	return d._exec(ctx, false, workDir, cmd, env, outWriter, errWriter)
}

func (d *dockerBackend) _exec(ctx context.Context, root bool, workDir string, cmd []string, env map[string]string, outWriter, errWriter io.Writer) error {
	execConf := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		WorkingDir:   workDir,
		Cmd:          cmd,
		Env:          lo.MapToSlice(env, func(k, v string) string { return k + "=" + v }),
	}
	if root {
		execConf.User = "root"
	} else {
		execConf.User = d.config.User
	}
	execID, err := d.c.ContainerExecCreate(ctx, d.config.ContainerName, execConf)
	if err != nil {
		return errors.Wrap(err, "creating exec")
	}

	ex, err := d.c.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return errors.Wrap(err, "attaching exec process")
	}
	defer ex.Close()

	_, err = stdcopy.StdCopy(outWriter, errWriter, ex.Reader)
	if err != nil {
		return errors.Wrap(err, "reading exec response")
	}

	res, err := d.c.ContainerExecInspect(ctx, execID.ID)
	if err != nil {
		return errors.Wrap(err, "inspecting exec result")
	}
	if res.ExitCode != 0 {
		return errors.Errorf("exec %s failed with exit code %d", strings.Join(cmd, " "), res.ExitCode)
	}

	return nil
}

func (d *dockerBackend) prepareDir(ctx context.Context, localName, remoteName string) (remotePath string, cleanup func(), err error) {
	remotePath = filepath.Join(d.config.RemoteDir, remoteName)
	err = d.exec(ctx, d.config.RemoteDir, []string{"mkdir", remotePath}, nil, io.Discard, io.Discard)
	if err != nil {
		return "", nil, errors.Wrap(err, "making remote tmp dir")
	}
	cleanup = func() {
		err := d.exec(context.Background(), d.config.RemoteDir, []string{"rm", "-r", remotePath}, nil, io.Discard, io.Discard)
		if err != nil {
			log.Errorf("failed to remove tmp repo dir: %+v", err)
		}
	}

	err = d.c.CopyToContainer(ctx, d.config.ContainerName, remotePath, tarfs.Compress(localName), types.CopyToContainerOptions{})
	if err != nil {
		cleanup()
		return "", nil, errors.Wrap(err, "copying files to container")
	}
	err = d.execRoot(ctx, d.config.RemoteDir, []string{"chown", "-R", d.config.User + ":" + d.config.Group, remotePath}, nil, io.Discard, io.Discard)
	if err != nil {
		cleanup()
		return "", nil, errors.Wrap(err, "setting remote dir owner")
	}
	return remotePath, cleanup, nil
}

func (d *dockerBackend) Pack(
	ctx context.Context,
	repoDir string,
	imageDest string,
	env map[string]string,
	logWriter io.Writer,
) (path string, err error) {
	remoteRepoPath, cleanupRepoDir, err := d.prepareDir(ctx, repoDir, "ns-repo")
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
	remoteEnvPath, cleanupEnvDir, err := d.prepareDir(ctx, localEnvTmp, "ns-env")
	if err != nil {
		return "", err
	}
	defer cleanupEnvDir()

	// TODO: support pushing to insecure registry for local development
	// https://github.com/buildpacks/lifecycle/issues/524
	// https://github.com/buildpacks/rfcs/blob/main/text/0111-support-insecure-registries.md
	// Workaround: use registry host "*.local" to allow google/go-containerregistry to detect as http protocol
	// see: https://github.com/traPtitech/NeoShowcase/issues/493
	err = d.exec(ctx,
		remoteRepoPath,
		[]string{"/cnb/lifecycle/creator", "-skip-restore", "-platform=" + remoteEnvPath, "-app=.", imageDest},
		ds.MergeMap(env, map[string]string{"CNB_PLATFORM_API": d.config.PlatformAPI}),
		logWriter, logWriter)
	if err != nil {
		return "", err
	}
	return remoteRepoPath, nil
}
