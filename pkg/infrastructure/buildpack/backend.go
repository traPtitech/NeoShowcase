package buildpack

import (
	"archive/tar"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"path"
	"strings"

	"github.com/samber/oops"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/tarfs"
)

const (
	workspaceDir = "/workspace"
	layersDir    = "/layers"
	platformDir  = "/platform"
	// appDir is also the working directory of the built image, which existing apps may refer to by absolute path.
	appDir = workspaceDir + "/ns-repo"
	// exportAnalyzedPath is outside the shared directories, so that the export step never reads analyzed.toml
	// written by the user code of earlier steps.
	exportAnalyzedPath = "/tmp/analyzed.toml"
)

type backend struct {
	config Config
	runner builder.IsolatedRunner
}

func NewBuildpackBackend(
	config Config,
	runner builder.IsolatedRunner,
) builder.BuildpackBackend {
	return &backend{
		config: config,
		runner: runner,
	}
}

func (b *backend) ReclaimOrphans(ctx context.Context) error {
	return b.runner.ReclaimOrphans(ctx)
}

func (b *backend) ReclaimEnded(ctx context.Context) error {
	return b.runner.ReclaimEnded(ctx)
}

// registryAuth returns the value of CNB_REGISTRY_AUTH for the registry, or false if no credentials are configured.
func registryAuth(registry builder.RegistryConfig) (string, bool) {
	if registry.Username == "" && registry.Password == "" {
		return "", false
	}
	basic := base64.StdEncoding.EncodeToString([]byte(registry.Username + ":" + registry.Password))
	b, _ := json.Marshal(map[string]string{registry.Addr: "Basic " + basic})
	return string(b), true
}

// inputStream is a tar stream written by a goroutine while it is read.
type inputStream struct {
	*io.PipeReader
	done chan struct{}
}

// Close stops writing the stream and waits for the goroutine writing it to exit.
func (s *inputStream) Close() error {
	_ = s.PipeReader.Close()
	<-s.done
	return nil
}

// input returns a tar stream holding the source under appDir and the app environment variables under platformDir/env.
// The caller must close it.
func input(repoDir string, env map[string]string) (*inputStream, error) {
	for k := range env {
		if k == "" || k == "." || k == ".." || strings.Contains(k, "/") {
			return nil, oops.Errorf("invalid environment variable name %q", k)
		}
	}

	pr, pw := io.Pipe()
	done := make(chan struct{})
	go func() {
		defer close(done)
		tw := tar.NewWriter(pw)
		err := tarfs.AddDir(tw, repoDir, strings.TrimPrefix(appDir, "/"), setAppMode)
		if err == nil {
			err = writeEnvFiles(tw, env)
		}
		if err == nil {
			err = tw.Close()
		}
		_ = pw.CloseWithError(err)
	}()
	return &inputStream{PipeReader: pr, done: done}, nil
}

// setAppMode makes an entry of the app directory accessible to any user, keeping only whether a file is executable,
// as git does. The built image runs as a user other than the one owning the app directory.
func setAppMode(h *tar.Header) {
	switch {
	case h.Typeflag == tar.TypeSymlink:
	case h.Typeflag == tar.TypeDir || h.Mode&0o111 != 0:
		h.Mode = 0o755
	default:
		h.Mode = 0o644
	}
}

func writeEnvFiles(tw *tar.Writer, env map[string]string) error {
	envDir := path.Join(strings.TrimPrefix(platformDir, "/"), "env")
	if err := tw.WriteHeader(&tar.Header{Name: envDir + "/", Mode: 0o700, Typeflag: tar.TypeDir}); err != nil {
		return err
	}
	for k, v := range env {
		err := tw.WriteHeader(&tar.Header{Name: path.Join(envDir, k), Mode: 0o600, Size: int64(len(v)), Typeflag: tar.TypeReg})
		if err != nil {
			return err
		}
		if _, err := tw.Write([]byte(v)); err != nil {
			return err
		}
	}
	return nil
}

func (b *backend) Pack(
	ctx context.Context,
	buildID string,
	repoDir string,
	imageDest string,
	imageConfig builder.ImageConfig,
	env map[string]string,
	logWriter io.Writer,
) (path string, err error) {
	in, err := input(repoDir, env)
	if err != nil {
		return "", err
	}
	defer in.Close()

	lifecycleEnv := map[string]string{"CNB_PLATFORM_API": b.config.PlatformAPI}
	// Allow pushing to insecure (http) registry for local development
	// https://github.com/buildpacks/rfcs/blob/main/text/0111-support-insecure-registries.md
	// see also: https://github.com/traPtitech/NeoShowcase/issues/493
	if imageConfig.Registry.Scheme == "http" {
		lifecycleEnv["CNB_INSECURE_REGISTRIES"] = imageConfig.Registry.Addr
	}
	exportSecretEnv := map[string]string{}
	if auth, ok := registryAuth(imageConfig.Registry); ok {
		exportSecretEnv["CNB_REGISTRY_AUTH"] = auth
	}

	dirArgs := []string{"-app=" + appDir, "-layers=" + layersDir, "-platform=" + platformDir}
	// The analyzer runs only in the export step, which runs no user code and is the only step given the credentials.
	// The image references it resolves (the run image and the previous image) are written outside the shared directories
	// and handed to the exporter, so user code cannot make the exporter pull other images with the credentials.
	// The image name is passed as $0 so that it is not interpreted by the shell.
	exportScript := `/cnb/lifecycle/analyzer -layers=` + layersDir + ` -analyzed=` + exportAnalyzedPath + ` -skip-layers "$0" && ` +
		`exec /cnb/lifecycle/exporter -app=` + appDir + ` -layers=` + layersDir + ` -analyzed=` + exportAnalyzedPath + ` "$0"`

	job := builder.IsolatedJob{
		ID:     buildID,
		Image:  b.config.Image,
		RunAs:  b.config.RunAs,
		Shared: []string{workspaceDir, layersDir, platformDir},
		Input:  in,
		Steps: []builder.IsolatedStep{
			{Name: "detect", Cmd: append([]string{"/cnb/lifecycle/detector"}, dirArgs...), WorkDir: appDir, Env: lifecycleEnv},
			{Name: "build", Cmd: append([]string{"/cnb/lifecycle/builder"}, dirArgs...), WorkDir: appDir, Env: lifecycleEnv},
			{Name: "export", Cmd: []string{"sh", "-c", exportScript, imageDest}, WorkDir: appDir, Env: lifecycleEnv, SecretEnv: exportSecretEnv},
		},
	}
	if err := b.runner.Run(ctx, job, logWriter); err != nil {
		return "", err
	}
	return appDir, nil
}
