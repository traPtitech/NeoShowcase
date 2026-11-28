package buildpack

import (
	"archive/tar"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
)

// recordingRunner records the job it is given, reading its input as it would be extracted.
type recordingRunner struct {
	job   builder.IsolatedJob
	files map[string]string
	modes map[string]int64
}

func (r *recordingRunner) Run(_ context.Context, job builder.IsolatedJob, _ io.Writer) error {
	r.job = job
	r.files = map[string]string{}
	r.modes = map[string]int64{}
	tr := tar.NewReader(job.Input)
	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		b, err := io.ReadAll(tr)
		if err != nil {
			return err
		}
		r.files[h.Name] = string(b)
		r.modes[h.Name] = h.Mode
	}
}

func (r *recordingRunner) ReclaimOrphans(context.Context) error { return nil }

func (r *recordingRunner) ReclaimEnded(context.Context) error { return nil }

func (r *recordingRunner) step(t *testing.T, name string) builder.IsolatedStep {
	t.Helper()
	for _, s := range r.job.Steps {
		if s.Name == name {
			return s
		}
	}
	t.Fatalf("step %s not found", name)
	return builder.IsolatedStep{}
}

func pack(t *testing.T, registry builder.RegistryConfig) (*recordingRunner, string) {
	t.Helper()
	repo := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(repo, "src"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(repo, "src", "main.go"), []byte("package main"), 0o644))

	runner := &recordingRunner{}
	b := NewBuildpackBackend(Config{
		Image:       "builder-image:1",
		RunAs:       builder.UIDGID{UID: 1001, GID: 1000},
		PlatformAPI: "0.11",
	}, runner)
	path, err := b.Pack(context.Background(), "buildid", repo, registry.Addr+"/app:buildid",
		builder.ImageConfig{Registry: registry}, map[string]string{"FOO": "bar"}, io.Discard)
	require.NoError(t, err)
	return runner, path
}

func TestPack_RunsLifecyclePhasesInSeparateSteps(t *testing.T) {
	runner, path := pack(t, builder.RegistryConfig{Scheme: "https", Addr: "registry.example.com"})

	require.Equal(t, "/workspace/ns-repo", path)
	require.Equal(t, "builder-image:1", runner.job.Image)
	require.Equal(t, builder.UIDGID{UID: 1001, GID: 1000}, runner.job.RunAs)
	require.ElementsMatch(t, []string{"/workspace", "/layers", "/platform"}, runner.job.Shared)

	var names []string
	for _, s := range runner.job.Steps {
		names = append(names, s.Name)
		require.Equal(t, "0.11", s.Env["CNB_PLATFORM_API"])
	}
	require.Equal(t, []string{"detect", "build", "export"}, names)

	require.Equal(t, "package main", runner.files["workspace/ns-repo/src/main.go"])
	require.Equal(t, "bar", runner.files["platform/env/FOO"])
}

func TestPack_ExportReanalyzesOutsideSharedDirectories(t *testing.T) {
	runner, _ := pack(t, builder.RegistryConfig{Scheme: "https", Addr: "registry.example.com"})

	export := strings.Join(runner.step(t, "export").Cmd, " ")
	analyzer := strings.Index(export, "/cnb/lifecycle/analyzer")
	exporter := strings.Index(export, "/cnb/lifecycle/exporter")
	require.GreaterOrEqual(t, analyzer, 0)
	require.Greater(t, exporter, analyzer)
	require.Contains(t, export, "-skip-layers")
	require.Equal(t, 2, strings.Count(export, "-analyzed=/tmp/analyzed.toml"))

	for _, name := range []string{"detect", "build"} {
		require.NotContains(t, strings.Join(runner.step(t, name).Cmd, " "), "analyzer")
	}
}

func TestPack_PassesCredentialsOnlyToExport(t *testing.T) {
	runner, _ := pack(t, builder.RegistryConfig{Scheme: "https", Addr: "registry.example.com", Username: "robot", Password: "s3cret"})

	encoded := base64.StdEncoding.EncodeToString([]byte("robot:s3cret"))
	want := `{"registry.example.com":"Basic ` + encoded + `"}`
	require.Equal(t, map[string]string{"CNB_REGISTRY_AUTH": want}, runner.step(t, "export").SecretEnv)
	for _, s := range runner.job.Steps {
		require.NotContains(t, s.Env, "CNB_REGISTRY_AUTH", s.Name)
		if s.Name == "export" {
			continue
		}
		require.Empty(t, s.SecretEnv, s.Name)
		passed := strings.Join(append(s.Cmd, mapValues(s.Env)...), " ")
		for _, credential := range []string{"s3cret", encoded} {
			require.NotContains(t, passed, credential, s.Name)
		}
	}
	for name, content := range runner.files {
		for _, credential := range []string{"s3cret", encoded} {
			require.NotContains(t, content, credential, name)
		}
	}
}

// The built image runs as a user other than the one the build runs as, and only the build user owns the app directory.
func TestPack_AppDirectoryIsAccessibleToAnyUser(t *testing.T) {
	repo := t.TempDir()
	require.NoError(t, os.Chmod(repo, 0o700))
	require.NoError(t, os.Mkdir(filepath.Join(repo, "private"), 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(repo, "private", "file"), nil, 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(repo, "run.sh"), nil, 0o700))

	runner := &recordingRunner{}
	b := NewBuildpackBackend(Config{Image: "builder-image:1", PlatformAPI: "0.11"}, runner)
	_, err := b.Pack(context.Background(), "buildid", repo, "registry.example.com/app:buildid",
		builder.ImageConfig{}, nil, io.Discard)
	require.NoError(t, err)

	require.Equal(t, map[string]int64{
		"workspace/ns-repo/":             0o755,
		"workspace/ns-repo/private/":     0o755,
		"workspace/ns-repo/private/file": 0o644,
		"workspace/ns-repo/run.sh":       0o755,
	}, lo.PickBy(runner.modes, func(name string, _ int64) bool { return strings.HasPrefix(name, "workspace/") }))
}

func TestPack_RejectsEnvNamesOutsideTheEnvDirectory(t *testing.T) {
	for _, name := range []string{"", ".", "..", "a/b"} {
		t.Run(name, func(t *testing.T) {
			runner := &recordingRunner{}
			b := NewBuildpackBackend(Config{Image: "builder-image:1", PlatformAPI: "0.11"}, runner)
			_, err := b.Pack(context.Background(), "buildid", t.TempDir(), "registry.example.com/app:buildid",
				builder.ImageConfig{}, map[string]string{name: "v"}, io.Discard)
			require.ErrorContains(t, err, "invalid environment variable name")
			require.Nil(t, runner.files, "should fail before running the build")
		})
	}
}

func TestPack_NoCredentialsWithoutRegistryAuth(t *testing.T) {
	runner, _ := pack(t, builder.RegistryConfig{Scheme: "http", Addr: "registry.local"})

	for _, s := range runner.job.Steps {
		require.Empty(t, s.SecretEnv)
		require.Equal(t, "registry.local", s.Env["CNB_INSECURE_REGISTRIES"])
	}
}

func mapValues(m map[string]string) []string {
	var vs []string
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

// failingRunner fails without reading the input, as when the build pod cannot start.
type failingRunner struct{}

func (failingRunner) Run(context.Context, builder.IsolatedJob, io.Writer) error {
	return errors.New("build pod cannot start")
}

func (failingRunner) ReclaimOrphans(context.Context) error { return nil }

func (failingRunner) ReclaimEnded(context.Context) error { return nil }

// A failed build used to leave the goroutine writing the input blocked forever, with a file of the repository open.
func TestPack_LeavesNothingRunningWhenInputIsNotRead(t *testing.T) {
	repo := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(repo, "big.bin"), make([]byte, 1<<20), 0o644))
	b := NewBuildpackBackend(Config{Image: "builder-image:1", PlatformAPI: "0.11"}, failingRunner{})

	before := runtime.NumGoroutine()
	for range 5 {
		_, err := b.Pack(context.Background(), "buildid", repo, "registry.example.com/app:buildid",
			builder.ImageConfig{}, nil, io.Discard)
		require.ErrorContains(t, err, "build pod cannot start")
	}
	require.LessOrEqual(t, runtime.NumGoroutine(), before)
}

// A repository that fails to be read must fail the build. The input would otherwise end at an entry boundary,
// which tar extracts without an error, and the build would go on with files missing.
func TestPack_FailsWhenTheRepositoryCannotBeRead(t *testing.T) {
	// The header of a directory has no content, so failing to list the directory ends the input at an entry boundary.
	unreadable := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(unreadable, "a.txt"), []byte("a"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(unreadable, "sub"), 0o000))
	if _, err := os.ReadDir(filepath.Join(unreadable, "sub")); err == nil {
		t.Skip("the directory is readable regardless of its mode, as when running as root")
	}

	for name, repo := range map[string]string{
		"missing directory":    filepath.Join(t.TempDir(), "missing"),
		"unreadable directory": unreadable,
	} {
		t.Run(name, func(t *testing.T) {
			b := NewBuildpackBackend(Config{Image: "builder-image:1", PlatformAPI: "0.11"}, &recordingRunner{})
			_, err := b.Pack(context.Background(), "buildid", repo, "registry.example.com/app:buildid",
				builder.ImageConfig{}, nil, io.Discard)
			require.Error(t, err)
		})
	}
}
