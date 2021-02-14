package k8simpl

import (
	"context"
	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/neoshowcase/pkg/container"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
	kubeconf, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	require.NoError(t, err)
	clientset, err := kubernetes.NewForConfig(kubeconf)
	require.NoError(t, err)

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
	t.Fatalf("wait pod timeout: %s", podName)
}

func TestNewManager(t *testing.T) {
	_, c, _ := prepareManager(t)

	_, err := c.ServerVersion()
	assert.NoError(t, err)
}
