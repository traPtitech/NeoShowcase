package dockerimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

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

func (d *dockerBackend) prepareAuth() error {
	auth, ok := d.dockerAuth()
	if ok {
		err := d.exec(context.Background(), "/", []string{"sh", "-c", fmt.Sprintf(`echo '%s' > ~/.docker/config.json`, auth)}, io.Discard, io.Discard)
		if err != nil {
			return errors.Wrap(err, "writing ~/.docker/config.json to builder")
		}
	}
	return nil
}

func (d *dockerBackend) exec(ctx context.Context, workDir string, cmd []string, outWriter, errWriter io.Writer) error {
	execID, err := d.c.ContainerExecCreate(ctx, d.config.ContainerName, types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Tty:          true,
		WorkingDir:   workDir,
		Cmd:          cmd,
	})
	if err != nil {
		return errors.Wrap(err, "creating exec")
	}
	res, err := d.c.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return errors.Wrap(err, "attaching exec process")
	}
	defer res.Close()

	_, err = stdcopy.StdCopy(outWriter, errWriter, res.Reader)
	if err != nil {
		return errors.Wrap(err, "reading exec response")
	}
	return nil
}

func (d *dockerBackend) Pack(ctx context.Context, repoDir string, imageDest string, logWriter io.Writer) error {
	tmpID := domain.NewID()
	baseDir := "/work"
	dstRepoPath := fmt.Sprintf("repo-%s", tmpID)
	dstPath := filepath.Join(baseDir, dstRepoPath)

	err := d.c.CopyToContainer(ctx, d.config.ContainerName, dstPath, tarfs.Compress(repoDir), types.CopyToContainerOptions{})
	if err != nil {
		return errors.Wrap(err, "copying file to container")
	}
	defer func() {
		err := d.exec(ctx, baseDir, []string{"rm", "-r", dstRepoPath}, io.Discard, io.Discard)
		if err != nil {
			log.Errorf("failed to remove remote repo dir: %+v", err)
		}
	}()

	err = d.exec(ctx, dstPath, []string{"/cnb/lifecycle/creator", "-app=.", imageDest}, logWriter, logWriter)
	if err != nil {
		return err
	}
	return nil
}
