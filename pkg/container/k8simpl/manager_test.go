package k8simpl

import (
	"context"
	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/neoshowcase/pkg/container"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"os"
	"strconv"
	"testing"
	"time"
)

func prepareManager(t *testing.T) (*Manager, *kubernetes.Clientset, *hub.Hub) {
	t.Helper()
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_K8S_TESTS")); !ok {
		t.SkipNow()
	}
	bus := hub.New()

	// k8s接続
	var (
		kubeconf *rest.Config
		err      error
	)
	if ctx := os.Getenv("K8S_TESTS_CLUSTER_CONTEXT"); len(ctx) > 0 {
		kubeconf, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: clientcmd.RecommendedHomeFile},
			&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}, CurrentContext: ctx}).ClientConfig()
	} else {
		kubeconf, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	}
	require.NoError(t, err)
	clientset, err := kubernetes.NewForConfig(kubeconf)
	require.NoError(t, err)

	if _, err := clientset.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: appNamespace,
		},
	}, metav1.CreateOptions{}); err != nil && !errors.IsAlreadyExists(err) {
		t.Fatal(err)
	}

	m, err := NewManager(bus, clientset)
	require.NoError(t, err)

	return m, clientset, bus
}

func waitPodRunning(t *testing.T, c *kubernetes.Clientset, podName string) {
	t.Helper()

	for i := 0; i < 120; i++ {
		pod, err := c.CoreV1().Pods(appNamespace).Get(context.Background(), podName, metav1.GetOptions{})
		require.NoError(t, err)

		if getContainerState(pod.Status) == container.StateRunning {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("wait pod running timeout: %s", podName)
}

func TestNewManager(t *testing.T) {
	_, c, _ := prepareManager(t)

	_, err := c.ServerVersion()
	assert.NoError(t, err)
}
