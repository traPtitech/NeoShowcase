package dockerimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config/configfile"
	types2 "github.com/docker/cli/cli/config/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
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
	c, err := client.NewClientWithOpts(client.FromEnv)
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
		err := d.exec(context.Background(), "/", []string{"sh", "-c", "mkdir -p ~/.docker"}, io.Discard, io.Discard)
		if err != nil {
			return errors.Wrap(err, "making ~/.docker directory")
		}
		err = d.exec(context.Background(), "/", []string{"sh", "-c", fmt.Sprintf(`echo '%s' > ~/.docker/config.json`, escapeSingleQuote(auth))}, io.Discard, io.Discard)
		if err != nil {
			return errors.Wrap(err, "writing ~/.docker/config.json to builder")
		}
	}
	return nil
}

func (d *dockerBackend) execRoot(ctx context.Context, workDir string, cmd []string, outWriter, errWriter io.Writer) error {
	return d._exec(ctx, true, workDir, cmd, outWriter, errWriter)
}

func (d *dockerBackend) exec(ctx context.Context, workDir string, cmd []string, outWriter, errWriter io.Writer) error {
	return d._exec(ctx, false, workDir, cmd, outWriter, errWriter)
}

func (d *dockerBackend) _exec(ctx context.Context, root bool, workDir string, cmd []string, outWriter, errWriter io.Writer) error {
	execConf := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		WorkingDir:   workDir,
		Cmd:          cmd,
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

func (d *dockerBackend) Pack(ctx context.Context, repoDir string, logWriter io.Writer, imageDest string) error {
	dstRepoPath := fmt.Sprintf("repo-%s", domain.NewID())
	remoteDstPath := filepath.Join(d.config.RemoteDir, dstRepoPath)

	err := d.exec(ctx, d.config.RemoteDir, []string{"mkdir", dstRepoPath}, io.Discard, io.Discard)
	if err != nil {
		return errors.Wrap(err, "making remote repo tmp dir")
	}
	err = d.c.CopyToContainer(ctx, d.config.ContainerName, remoteDstPath, tarfs.Compress(repoDir), types.CopyToContainerOptions{})
	if err != nil {
		return errors.Wrap(err, "copying file to container")
	}
	err = d.execRoot(ctx, d.config.RemoteDir, []string{"chown", "-R", d.config.User + ":" + d.config.Group, dstRepoPath}, io.Discard, io.Discard)
	if err != nil {
		return errors.Wrap(err, "setting remote repo owner")
	}
	defer func() {
		err := d.exec(ctx, d.config.RemoteDir, []string{"rm", "-r", dstRepoPath}, io.Discard, io.Discard)
		if err != nil {
			log.Errorf("failed to remove remote repo dir: %+v", err)
		}
	}()

	// TODO: support pushing to insecure registry for local development
	// https://github.com/buildpacks/lifecycle/issues/524
	// https://github.com/buildpacks/rfcs/blob/main/text/0111-support-insecure-registries.md
	err = d.exec(ctx, remoteDstPath, []string{"/cnb/lifecycle/creator", "-skip-restore", "-app=.", imageDest}, logWriter, logWriter)
	if err != nil {
		return err
	}
	return nil
}
