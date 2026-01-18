package buildpack

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config/configfile"
	types2 "github.com/docker/cli/cli/config/types"
	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/tarfs"
)

type backend struct {
	config Config
	client domain.BuildpackHelperServiceClient
}

func NewBuildpackBackend(
	config Config,
	client domain.BuildpackHelperServiceClient,
) builder.BuildpackBackend {
	return &backend{
		config: config,
		client: client,
	}
}

func (b *backend) dockerAuth(imageConfig builder.ImageConfig) (s string, ok bool) {
	if imageConfig.Registry.Username == "" && imageConfig.Registry.Password == "" {
		return "", false
	}
	c := configfile.ConfigFile{
		AuthConfigs: map[string]types2.AuthConfig{
			imageConfig.Registry.Addr: {
				Username: imageConfig.Registry.Username,
				Password: imageConfig.Registry.Password,
			},
		},
	}
	bytes, _ := json.Marshal(&c)
	return string(bytes), true
}

func escapeSingleQuote(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}

func (b *backend) exec(ctx context.Context, workDir string, cmd []string, env map[string]string, logWriter io.Writer) error {
	code, err := b.client.Exec(ctx, workDir, cmd, env, logWriter)
	if err != nil {
		return err
	}
	if code != 0 {
		return fmt.Errorf("command exited with code %d", code)
	}
	return nil
}

func (b *backend) prepareAuth(imageConfig builder.ImageConfig) error {
	auth, ok := b.dockerAuth(imageConfig)
	if ok {
		err := b.exec(context.Background(), "/", []string{"sh", "-c", "mkdir -p ~/.docker"}, nil, io.Discard)
		if err != nil {
			return errors.Wrap(err, "making ~/.docker directory")
		}
		err = b.exec(context.Background(), "/", []string{"sh", "-c", fmt.Sprintf(`echo '%s' > ~/.docker/config.json`, escapeSingleQuote(auth))}, nil, io.Discard)
		if err != nil {
			return errors.Wrap(err, "writing ~/.docker/config.json to builder")
		}
	}
	return nil
}

func (b *backend) prepareDir(ctx context.Context, localName, remoteName string) (remotePath string, cleanup func(), err error) {
	remotePath = filepath.Join(b.config.RemoteDir, remoteName)
	err = b.exec(ctx, b.config.RemoteDir, []string{"mkdir", remotePath}, nil, io.Discard)
	if err != nil {
		return "", nil, errors.Wrap(err, "making remote tmp dir")
	}
	cleanup = func() {
		err := b.exec(context.Background(), b.config.RemoteDir, []string{"rm", "-r", remotePath}, nil, io.Discard)
		if err != nil {
			slog.ErrorContext(ctx, "failed to remove tmp repo dir", "error", err)
		}
	}

	err = b.client.CopyFileTree(ctx, remotePath, tarfs.Compress(localName))
	if err != nil {
		cleanup()
		return "", nil, errors.Wrap(err, "copying files to container")
	}
	return remotePath, cleanup, nil
}

func (b *backend) Pack(
	ctx context.Context,
	repoDir string,
	imageDest string,
	imageConfig builder.ImageConfig,
	env map[string]string,
	logWriter io.Writer,
) (path string, err error) {
	err = b.prepareAuth(imageConfig)
	if err != nil {
		return "", err
	}

	remoteRepoPath, cleanupRepoDir, err := b.prepareDir(ctx, repoDir, "ns-repo")
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
	remoteEnvPath, cleanupEnvDir, err := b.prepareDir(ctx, localEnvTmp, "ns-env")
	if err != nil {
		return "", err
	}
	defer cleanupEnvDir()

	// TODO: support pushing to insecure registry for local development
	// https://github.com/buildpacks/lifecycle/issues/524
	// https://github.com/buildpacks/rfcs/blob/main/text/0111-support-insecure-registries.md
	// Workaround: use registry host "*.local" to allow google/go-containerregistry to detect as http protocol
	// see: https://github.com/traPtitech/NeoShowcase/issues/493
	err = b.exec(ctx,
		remoteRepoPath,
		[]string{"/cnb/lifecycle/creator", "-skip-restore", "-platform=" + remoteEnvPath, "-app=.", imageDest},
		map[string]string{"CNB_PLATFORM_API": b.config.PlatformAPI},
		logWriter)
	if err != nil {
		return "", err
	}
	return remoteRepoPath, nil
}
