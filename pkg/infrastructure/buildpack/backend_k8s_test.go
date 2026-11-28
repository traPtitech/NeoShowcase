package buildpack

import (
	"bytes"
	"context"
	"encoding/base64"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner/k8srunner"
)

// TestPack_K8sIsolatesCredentials builds the malicious fixture of the builder tests in a k8s cluster.
// The cluster needs the builder image and an http registry reachable from pods, given by the environment variables below.
// Its build script dumps everything that may hold the credentials to the build log,
// and rewrites /layers/analyzed.toml to point the run image at a nonexistent image.
func TestPack_K8sIsolatesCredentials(t *testing.T) {
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_BUILDPACK_K8S_TESTS")); !ok {
		t.SkipNow()
	}
	registry := os.Getenv("BUILDPACK_K8S_TESTS_REGISTRY") // e.g. lc-reg:5000
	require.NotEmpty(t, registry)
	image := os.Getenv("BUILDPACK_K8S_TESTS_IMAGE")
	if image == "" {
		image = DefaultImage
	}

	kubeconf, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: clientcmd.RecommendedHomeFile},
		&clientcmd.ConfigOverrides{CurrentContext: os.Getenv("K8S_TESTS_CLUSTER_CONTEXT")}).ClientConfig()
	require.NoError(t, err)
	client, err := kubernetes.NewForConfig(kubeconf)
	require.NoError(t, err)

	ctx := context.Background()
	// The registry is resolved by its short name, which works only from its namespace.
	const namespace = "default"
	runner, err := k8srunner.NewRunner(kubeconf, client, k8srunner.Options{
		Config:         k8srunner.Config{Namespace: namespace},
		Hostname:       "fake-builder-" + domain.NewID(),
		PullPolicy:     isolatedrunner.PullIfNotPresent,
		ActiveDeadline: 30 * time.Minute,
	})
	require.NoError(t, err)
	b := NewBuildpackBackend(Config{Image: image, RunAs: DefaultRunAs, PlatformAPI: "0.11"}, runner)

	const username = "ns-build-test-robot"
	password := "ns-build-test-" + domain.NewID()
	buildID := domain.NewID()
	var log bytes.Buffer
	_, err = b.Pack(ctx, buildID, "../../usecase/builder/testdata/malicious-buildpack", registry+"/ns-build-test/malicious:"+buildID,
		builder.ImageConfig{Registry: builder.RegistryConfig{Scheme: "http", Addr: registry, Username: username, Password: password}},
		map[string]string{"BP_NODE_RUN_SCRIPTS": "build"}, &log)
	require.NoError(t, err, log.String())

	require.Contains(t, log.String(), "[malicious] done", "the build script should have run")
	require.NotContains(t, log.String(), password)
	require.NotContains(t, log.String(), base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
}
