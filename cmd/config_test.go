package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/buildpack"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner/dockerrunner"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner/k8srunner"
	"github.com/traPtitech/neoshowcase/pkg/util/cli"
)

func loadConfig(t *testing.T, yaml string) Config {
	t.Helper()
	path := filepath.Join(t.TempDir(), "ns.yaml")
	require.NoError(t, os.WriteFile(path, []byte(yaml), 0o600))
	var c Config
	cli.CobraOnInitializeFunc(&path, &c)()
	return c
}

func TestConfig_EmptyBuildpackResourcesMeanNoLimits(t *testing.T) {
	c := loadConfig(t, `
components:
  builder:
    buildpack:
      runtime: k8s
      k8s:
        namespace: builds
        resources: {}
`)
	r := c.Components.Builder.Buildpack.K8s.Resources
	require.Empty(t, r.Requests.CPU)
	require.Empty(t, r.Requests.Memory)
	require.Empty(t, r.Limits.CPU)
	require.Empty(t, r.Limits.Memory)
}

func TestNewBuilder_RequiresRuntime(t *testing.T) {
	// The builder config of the docker setup before buildpack builds ran in containers of their own.
	c := loadConfig(t, `
components:
  builder:
    buildpack:
      helper:
        address: http://ns-builder:1235
`)
	_, err := NewBuilder(c)
	require.ErrorContains(t, err, "components.builder.buildpack.runtime is required")
}

func TestNewBuilder_RejectsUnknownRuntime(t *testing.T) {
	c := loadConfig(t, `
components:
  builder:
    buildpack:
      runtime: podman
`)
	_, err := NewBuilder(c)
	require.ErrorContains(t, err, "invalid components.builder.buildpack.runtime")
}

func TestNewBuilder_AcceptsTheRuntimeNamesOfTheController(t *testing.T) {
	// Building the k8s runtime fails outside a cluster on loading the in-cluster config, before connecting to the controller.
	t.Setenv("KUBERNETES_SERVICE_HOST", "")
	c := loadConfig(t, `
components:
  controller:
    token: token
  builder:
    buildpack:
      runtime: kubernetes
      k8s:
        namespace: builds
`)
	_, err := NewBuilder(c)
	require.ErrorIs(t, err, rest.ErrNotInCluster)
}

func TestConfig_BuildpackImagePullPolicyDefaultsToIfNotPresent(t *testing.T) {
	c := loadConfig(t, `
components:
  builder:
    buildpack:
      runtime: k8s
`)
	require.Equal(t, "IfNotPresent", c.Components.Builder.Buildpack.ImagePullPolicy)
}

func TestProvideK8sRunnerOptions(t *testing.T) {
	c := loadConfig(t, `
components:
  builder:
    stepTimeout: 20m
    buildpack:
      runtime: k8s
      imagePullPolicy: Never
      k8s:
        namespace: builds
        resources:
          limits:
            cpu: "3"
`)
	builderConfig, err := provideBuilderConfig(c)
	require.NoError(t, err)
	opts, err := provideK8sRunnerOptions(c, builderConfig)
	require.NoError(t, err)

	hostname, err := os.Hostname()
	require.NoError(t, err)
	var want k8srunner.Options
	want.Config.Namespace = "builds"
	want.Config.Resources.Limits.CPU = "3"
	want.Hostname = hostname
	want.PullPolicy = isolatedrunner.PullNever
	want.ActiveDeadline = 20 * time.Minute
	require.Equal(t, want, opts)
}

func TestProvideDockerRunnerOptions(t *testing.T) {
	c := loadConfig(t, `
components:
  builder:
    buildpack:
      runtime: docker
      imagePullPolicy: Never
      docker:
        network: ns_default
`)
	opts, err := provideDockerRunnerOptions(c)
	require.NoError(t, err)

	hostname, err := os.Hostname()
	require.NoError(t, err)
	require.Equal(t, dockerrunner.Options{
		Config:     dockerrunner.Config{Network: "ns_default"},
		Hostname:   hostname,
		PullPolicy: isolatedrunner.PullNever,
	}, opts)
}

func TestConfig_BuildpackImageDefaultsToBuiltInImage(t *testing.T) {
	c := loadConfig(t, `
components:
  builder:
    buildpack:
      runtime: k8s
`)
	require.Equal(t, buildpack.DefaultImage, c.Components.Builder.Buildpack.Image)
}

func TestConfig_BuildpackPodPlacement(t *testing.T) {
	c := loadConfig(t, `
components:
  builder:
    buildpack:
      runtime: k8s
      k8s:
        namespace: builds
        nodeSelector:
          - key: kubernetes.io/arch
            value: amd64
        tolerations:
          - key: ns.trap.jp/builder
            operator: Exists
            effect: NoExecute
            tolerationSeconds: 30
`)
	k := c.Components.Builder.Buildpack.K8s
	require.Equal(t, []k8srunner.NodeLabel{{Key: "kubernetes.io/arch", Value: "amd64"}}, k.NodeSelector)
	require.Equal(t, []k8srunner.Toleration{{
		Key: "ns.trap.jp/builder", Operator: "Exists", Effect: "NoExecute", TolerationSeconds: ptr.To[int64](30),
	}}, k.Tolerations)
}

func TestProvideK8sRunner_RejectsNodeSelectorWrittenAsMap(t *testing.T) {
	// A pod spec and Helm values write a node selector as a map.
	c := loadConfig(t, `
components:
  builder:
    buildpack:
      runtime: k8s
      k8s:
        namespace: builds
        nodeSelector:
          kubernetes.io/arch: amd64
`)
	builderConfig, err := provideBuilderConfig(c)
	require.NoError(t, err)
	opts, err := provideK8sRunnerOptions(c, builderConfig)
	require.NoError(t, err)
	restConfig := &rest.Config{Host: "https://k8s.invalid"}
	client, err := kubernetes.NewForConfig(restConfig)
	require.NoError(t, err)
	_, err = provideK8sRunner(restConfig, client, opts)
	require.ErrorContains(t, err, "list of key/value")
}

func TestProvideBuildpackConfig(t *testing.T) {
	c := loadConfig(t, `
components:
  builder:
    buildpack:
      runtime: k8s
      image: example.com/builder:1
      runAs:
        uid: 2000
        gid: 3000
      platformAPI: "0.12"
`)
	require.Equal(t, buildpack.Config{
		Image:       "example.com/builder:1",
		RunAs:       builder.UIDGID{UID: 2000, GID: 3000},
		PlatformAPI: "0.12",
	}, provideBuildpackConfig(c))
}

func TestProvideBuildpackConfig_Defaults(t *testing.T) {
	c := loadConfig(t, `
components:
  builder:
    buildpack:
      runtime: k8s
`)
	require.Equal(t, buildpack.Config{
		Image:       buildpack.DefaultImage,
		RunAs:       buildpack.DefaultRunAs,
		PlatformAPI: "0.11",
	}, provideBuildpackConfig(c))
}

func TestProvideRunners_WithDefaults(t *testing.T) {
	c := loadConfig(t, `
components:
  builder:
    buildpack:
      k8s:
        namespace: builds
`)
	builderConfig, err := provideBuilderConfig(c)
	require.NoError(t, err)

	t.Run("k8s", func(t *testing.T) {
		opts, err := provideK8sRunnerOptions(c, builderConfig)
		require.NoError(t, err)
		restConfig := &rest.Config{Host: "https://k8s.invalid"}
		client, err := kubernetes.NewForConfig(restConfig)
		require.NoError(t, err)
		runner, err := provideK8sRunner(restConfig, client, opts)
		require.NoError(t, err)
		require.NotNil(t, runner)
	})
	t.Run("docker", func(t *testing.T) {
		opts, err := provideDockerRunnerOptions(c)
		require.NoError(t, err)
		client, err := dockerimpl.NewClientFromEnv()
		require.NoError(t, err)
		runner, err := provideDockerRunner(client, opts)
		require.NoError(t, err)
		require.NotNil(t, runner)
	})
}
