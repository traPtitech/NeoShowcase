package builder

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	buildkit "github.com/moby/buildkit/client"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/buildpack"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/registry"
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

// These tests build the fixtures under testdata against real buildkitd, buildpack helper and registry,
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

	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_BUILD_TESTS")); !ok {
		t.SkipNow()
	}

	bk, err := buildkit.New(context.Background(), getEnvOrDefault("BUILDKIT_ADDRESS", "unix:///run/buildkit/buildkitd.sock"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = bk.Close() })

	var bpConfig buildpack.Config
	bpConfig.RemoteDir = "/workspace"
	bpConfig.PlatformAPI = "0.11"
	bp := buildpack.NewBuildpackBackend(bpConfig, grpc.NewBuildpackHelperServiceClient(getEnvOrDefault("BUILDPACK_HELPER_ADDRESS", "http://buildpack:1235")))

	imageConfig := builder.ImageConfig{
		Registry: builder.RegistryConfig{
			Scheme: "http",
			Addr:   getEnvOrDefault("REGISTRY_ADDRESS", "registry.local"),
		},
		NamePrefix:    "ns-build-test/",
		TmpNamePrefix: "ns-build-test-tmp/",
	}

	return &ServiceImpl{
		config:    &Config{StepTimeout: 30 * time.Minute},
		client:    &mocks.ControllerBuilderServiceClientMock{},
		buildkit:  bk,
		buildpack: bp,
		regclient: registry.NewClient(imageConfig),
		// Instead of cloning, copy the fixture directory given as the repository URL.
		gitsvc: &mocks.GitServiceMock{
			CloneRepositoryFunc: func(_ context.Context, dir string, repo *domain.Repository, _ string) error {
				return os.CopyFS(dir, os.DirFS(repo.URL))
			},
		},
		imageConfig: imageConfig,
	}
}

// runBuild builds the given fixture with the given build config and asserts that the build succeeds.
func runBuild(t *testing.T, s *ServiceImpl, fixture string, bc domain.BuildConfig, env map[string]string) {
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

// Buildpack builds share the single remote workspace of the buildpack helper, so the tests must not run in parallel.
func TestBuild_RuntimeBuildpack(t *testing.T) {
	s := prepareService(t)

	for fixture := range runtimeFixtures {
		t.Run(fixture, func(t *testing.T) {
			runBuild(t, s, fixture, &domain.BuildConfigRuntimeBuildpack{}, nil)
		})
	}
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
