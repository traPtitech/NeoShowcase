package builder

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/client/llb/sourceresolver"
	gateway "github.com/moby/buildkit/frontend/gateway/client"
	dockerclient "github.com/moby/moby/client"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/buildpack"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner/dockerrunner"
	"github.com/traPtitech/neoshowcase/pkg/test/mocks"
)

// runtimeFixtures maps each fixture used for runtime builds to its build command for command builds.
// A smoke test with a few lightweight languages, rather than covering all supported ones.
var runtimeFixtures = map[string]string{
	"go":     "go build -o app .",
	"nodejs": "npm ci",
	"python": "uv sync --locked",
}

// staticFixture is used for static builds; its build script outputs dist/index.html.
const staticFixture = "nodejs"

// These tests build the fixtures under testdata against real buildkitd, Docker daemon and registry,
// to check that builds keep working after upgrading container-related dependencies.
// They are expected to run inside the compose network; see the builder-test service in compose.yaml.

func getEnvOrDefault(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func prepareService(t *testing.T) *ServiceImpl {
	t.Helper()
	return prepareServiceWithRegistry(t, builder.RegistryConfig{})
}

// prepareServiceWithRegistry prepares the service with the given registry credentials.
// The test registry accepts any credentials.
func prepareServiceWithRegistry(t *testing.T, credentials builder.RegistryConfig) *ServiceImpl {
	t.Helper()

	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_BUILD_TESTS")); !ok {
		t.SkipNow()
	}

	bk, err := buildkit.New(context.Background(), getEnvOrDefault("BUILDKIT_ADDRESS", "unix:///run/buildkit/buildkitd.sock"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = bk.Close() })

	docker, err := dockerclient.New(dockerclient.FromEnv)
	require.NoError(t, err)
	t.Cleanup(func() { _ = docker.Close() })
	hostname, err := os.Hostname()
	require.NoError(t, err)
	runner, err := dockerrunner.NewRunner(docker, dockerrunner.Options{
		Config:     dockerrunner.Config{Network: getEnvOrDefault("BUILDPACK_DOCKER_NETWORK", "neoshowcase_default")},
		Hostname:   hostname,
		PullPolicy: isolatedrunner.PullIfNotPresent,
	})
	require.NoError(t, err)
	bp := buildpack.NewBuildpackBackend(buildpack.Config{
		Image:       getEnvOrDefault("BUILDPACK_IMAGE", buildpack.DefaultImage),
		RunAs:       buildpack.DefaultRunAs,
		PlatformAPI: "0.11",
	}, runner)

	imageConfig := builder.ImageConfig{
		Registry: builder.RegistryConfig{
			Scheme:   "http",
			Addr:     getEnvOrDefault("REGISTRY_ADDRESS", "registry.local"),
			Username: credentials.Username,
			Password: credentials.Password,
		},
		NamePrefix:    "ns-build-test/",
		TmpNamePrefix: "ns-build-test-tmp/",
	}

	// Instead of cloning, copy the fixture directory given as the repository URL.
	gitsvc := &mocks.GitServiceMock{
		CloneRepositoryFunc: func(_ context.Context, dir string, repo *domain.Repository, _ string) error {
			return os.CopyFS(dir, os.DirFS(repo.URL))
		},
	}
	return newService(&Config{StepTimeout: 30 * time.Minute}, &mocks.ControllerBuilderServiceClientMock{}, bk, bp, gitsvc, imageConfig)
}

// runBuild builds the given fixture with the given build config, asserts that the build succeeds, and returns the state of the build.
func runBuild(t *testing.T, s *ServiceImpl, fixture string, bc domain.BuildConfig, env map[string]string) *state {
	t.Helper()

	fixtureDir, err := filepath.Abs(filepath.Join("testdata", fixture))
	require.NoError(t, err)

	app := &domain.Application{
		ID:     domain.NewID(),
		Config: domain.ApplicationConfig{BuildConfig: bc},
	}
	build := &domain.Build{ID: domain.NewID(), ApplicationID: app.ID}
	repo := &domain.Repository{URL: fixtureDir}
	var envs []*domain.Environment
	for k, v := range env {
		envs = append(envs, &domain.Environment{ApplicationID: app.ID, Key: k, Value: v})
	}

	st, err := newState(app, envs, build, repo, s.client)
	require.NoError(t, err)
	t.Cleanup(st.Done)

	status := s.process(context.Background(), st)
	if status != domain.BuildStatusSucceeded {
		t.Log(st.logWriter.buf.String())
	}
	require.Equal(t, domain.BuildStatusSucceeded, status)

	if st.deployType() == domain.DeployTypeRuntime {
		_, err := s.fetchImageSize(context.Background(), st)
		require.NoError(t, err, "built image should be pushed to the registry")
		t.Cleanup(func() {
			_ = s.regclient.DeleteImage(context.Background(), s.imageConfig.ImageName(app.ID), s.imageTag(build))
		})
	}
	if st.deployType() == domain.DeployTypeStatic {
		require.Contains(t, savedArtifactFiles(t, s, build.ID), "index.html")
	}
	return st
}

// launch runs cmd through the launcher of the image built for st, with the user, environment variables and
// working directory the image is configured with, as a container runtime would run the image.
// It returns the output of the run, and the error if the image has no user configured, or cmd or the launcher fails.
func launch(t *testing.T, s *ServiceImpl, st *state, cmd ...string) (string, error) {
	t.Helper()

	statusCh := make(chan *buildkit.SolveStatus)
	var output bytes.Buffer
	collected := make(chan struct{})
	go func() {
		defer close(collected)
		for status := range statusCh {
			for _, l := range status.Logs {
				output.Write(l.Data)
			}
		}
	}()

	image := s.destImage(st.app, st.build)
	_, err := s.buildkit.Build(context.Background(), buildkit.SolveOpt{Session: s.authSessions()}, "",
		func(ctx context.Context, c gateway.Client) (*gateway.Result, error) {
			_, _, dt, err := c.ResolveImageConfig(ctx, image, sourceresolver.Opt{})
			if err != nil {
				return nil, err
			}
			var img ocispecs.Image
			if err := json.Unmarshal(dt, &img); err != nil {
				return nil, err
			}
			// The callback runs in a goroutine of buildkit, from which require cannot stop the test.
			if img.Config.User == "" {
				return nil, errors.New("the image should be configured to run as the user of the run image")
			}
			base, err := llb.Image(image).WithImageConfig(dt)
			if err != nil {
				return nil, err
			}
			def, err := base.User(img.Config.User).
				Run(llb.Args(append([]string{"/cnb/lifecycle/launcher"}, cmd...)), llb.IgnoreCache).
				Root().Marshal(ctx)
			if err != nil {
				return nil, err
			}
			return c.Solve(ctx, gateway.SolveRequest{Definition: def.ToPB(), Evaluate: true})
		}, statusCh)
	<-collected
	return output.String(), err
}

// savedArtifactFiles returns the file names in the artifact saved for the build.
func savedArtifactFiles(t *testing.T, s *ServiceImpl, buildID string) []string {
	t.Helper()

	for _, call := range s.client.(*mocks.ControllerBuilderServiceClientMock).SaveArtifactCalls() {
		if call.Artifact.BuildID != buildID {
			continue
		}

		gr, err := gzip.NewReader(bytes.NewReader(call.Body))
		require.NoError(t, err)

		tr := tar.NewReader(gr)

		var names []string
		for {
			h, err := tr.Next()
			if errors.Is(err, io.EOF) {
				break
			}
			require.NoError(t, err)
			names = append(names, filepath.Clean(h.Name))
		}

		return names
	}

	t.Fatalf("artifact for build %s was not saved", buildID)
	return nil
}

func TestBuild_RuntimeBuildpack(t *testing.T) {
	s := prepareService(t)

	for fixture := range runtimeFixtures {
		t.Run(fixture, func(t *testing.T) {
			runBuild(t, s, fixture, &domain.BuildConfigRuntimeBuildpack{}, nil)
		})
	}
}

// The repository is cloned into a directory only its owner can enter, and may hold files only their owner can read,
// while the built image runs as a user other than the one running the build.
func TestBuild_RuntimeBuildpackAppIsAccessibleToRunUser(t *testing.T) {
	s := prepareService(t)
	clone := s.gitsvc.(*mocks.GitServiceMock).CloneRepositoryFunc
	s.gitsvc = &mocks.GitServiceMock{
		CloneRepositoryFunc: func(ctx context.Context, dir string, repo *domain.Repository, ref string) error {
			if err := clone(ctx, dir, repo, ref); err != nil {
				return err
			}
			if err := os.Chmod(dir, 0o700); err != nil {
				return err
			}
			if err := os.Mkdir(filepath.Join(dir, "private"), 0o700); err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(dir, "private", "file"), []byte("private"), 0o600)
		},
	}

	// The Go buildpack removes the source files from the app directory except those kept by BP_KEEP_FILES.
	st := runBuild(t, s, "go", &domain.BuildConfigRuntimeBuildpack{}, map[string]string{"BP_KEEP_FILES": "private/*"})

	out, err := launch(t, s, st, "cat", "private/file")
	require.NoError(t, err, out)
	require.Contains(t, out, "private")
}

// The build script of the fixture searches for the registry credentials and writes everything it finds to the build log,
// then rewrites /layers/analyzed.toml to point the run image at a nonexistent image.
func TestBuild_RuntimeBuildpackIsolatesCredentials(t *testing.T) {
	const username = "ns-build-test-robot"
	password := "ns-build-test-" + domain.NewID()
	s := prepareServiceWithRegistry(t, builder.RegistryConfig{Username: username, Password: password})

	log := runBuild(t, s, "malicious-buildpack", &domain.BuildConfigRuntimeBuildpack{}, map[string]string{"BP_NODE_RUN_SCRIPTS": "build"}).logWriter.buf.String()

	require.Contains(t, log, "[malicious] done", "the build script should have run")
	require.NotContains(t, log, password)
	require.NotContains(t, log, base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
}

func TestBuild_RuntimeDockerfile(t *testing.T) {
	t.Parallel()

	s := prepareService(t)

	for fixture := range runtimeFixtures {
		t.Run(fixture, func(t *testing.T) {
			t.Parallel()
			runBuild(t, s, fixture, &domain.BuildConfigRuntimeDockerfile{DockerfileName: "Dockerfile"}, nil)
		})
	}
}

// firstBaseImage returns the image of the first FROM instruction in the fixture's Dockerfile,
// so that the command build uses the same (Renovate-managed) image as the Dockerfile build.
func firstBaseImage(t *testing.T, fixture string) string {
	t.Helper()

	b, err := os.ReadFile(filepath.Join("testdata", fixture, "Dockerfile"))
	require.NoError(t, err)
	for line := range strings.Lines(string(b)) {
		fields := strings.Fields(line)
		if len(fields) >= 2 && strings.EqualFold(fields[0], "FROM") {
			return fields[1]
		}
	}

	t.Fatalf("no FROM instruction found in %s/Dockerfile", fixture)
	return ""
}

func TestBuild_RuntimeCmd(t *testing.T) {
	t.Parallel()

	s := prepareService(t)

	for fixture, buildCmd := range runtimeFixtures {
		t.Run(fixture, func(t *testing.T) {
			t.Parallel()

			runBuild(t, s, fixture, &domain.BuildConfigRuntimeCmd{
				BaseImage: firstBaseImage(t, fixture),
				BuildCmd:  buildCmd,
			}, nil)
		})
	}
}

func TestBuild_StaticBuildpack(t *testing.T) {
	s := prepareService(t)
	runBuild(t, s, staticFixture, &domain.BuildConfigStaticBuildpack{
		StaticConfig: domain.StaticConfig{ArtifactPath: "dist"},
	}, map[string]string{"BP_NODE_RUN_SCRIPTS": "build"})
}

func TestBuild_StaticDockerfile(t *testing.T) {
	t.Parallel()

	s := prepareService(t)
	runBuild(t, s, staticFixture, &domain.BuildConfigStaticDockerfile{
		StaticConfig:   domain.StaticConfig{ArtifactPath: "/app/dist"},
		DockerfileName: "Dockerfile",
	}, nil)
}

func TestBuild_StaticCmd(t *testing.T) {
	t.Parallel()

	s := prepareService(t)
	runBuild(t, s, staticFixture, &domain.BuildConfigStaticCmd{
		StaticConfig: domain.StaticConfig{ArtifactPath: "dist"},
		BaseImage:    firstBaseImage(t, staticFixture),
		BuildCmd:     "npm ci && npm run build",
	}, nil)
}
