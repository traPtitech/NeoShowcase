package builder

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
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
func runBuild(t *testing.T, s *ServiceImpl, fixture string, bc domain.BuildConfig) {
	t.Helper()

	fixtureDir, err := filepath.Abs(filepath.Join("testdata", fixture))
	require.NoError(t, err)

	app := &domain.Application{
		ID:     domain.NewID(),
		Config: domain.ApplicationConfig{BuildConfig: bc},
	}
	build := &domain.Build{ID: domain.NewID(), ApplicationID: app.ID}
	repo := &domain.Repository{URL: fixtureDir}

	st, err := newState(app, nil, build, repo, s.client)
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
}

// Buildpack builds share the single remote workspace of the buildpack helper, so the tests must not run in parallel.
func TestBuild_RuntimeBuildpack(t *testing.T) {
	s := prepareService(t)

	// A smoke test with a few lightweight languages, rather than covering all supported ones.
	fixtures := []string{
		"go",
		"nodejs",
		"python",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			runBuild(t, s, fixture, &domain.BuildConfigRuntimeBuildpack{})
		})
	}
}

func TestBuild_RuntimeDockerfile(t *testing.T) {
	t.Parallel()
	s := prepareService(t)

	fixtures := []string{
		"go",
		"nodejs",
		"python",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			t.Parallel()
			runBuild(t, s, fixture, &domain.BuildConfigRuntimeDockerfile{DockerfileName: "Dockerfile"})
		})
	}
}
