package k8simpl

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefik/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func prepareManager(t *testing.T, bus domain.Bus) (*k8sBackend, *kubernetes.Clientset, *traefikv1alpha1.TraefikV1alpha1Client) {
	t.Helper()
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_K8S_TESTS")); !ok {
		t.SkipNow()
	}

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
	client, err := kubernetes.NewForConfig(kubeconf)
	require.NoError(t, err)
	traefikClient, err := traefikv1alpha1.NewForConfig(kubeconf)
	require.NoError(t, err)

	if _, err := client.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: appNamespace,
		},
	}, metav1.CreateOptions{}); err != nil && !errors.IsAlreadyExists(err) {
		t.Fatal(err)
	}

	b := NewK8SBackend(bus, client, traefikClient, nil, nil, domain.StaticServerConnectivityConfig{})
	err = b.Start(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = b.Dispose(context.Background())
	})

	return b.(*k8sBackend), client, traefikClient
}

func waitPodRunning(t *testing.T, b *k8sBackend, appID string) {
	t.Helper()

	for i := 0; i < 120; i++ {
		status, err := b.GetContainer(context.Background(), appID)
		if err != nil && err != domain.ErrContainerNotFound {
			t.Fatalf("error in get container: %v", err)
		}
		if err == nil && status.State == domain.ContainerStateRunning {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("wait pod running timeout: %s", appID)
}

func waitPodDeleted(t *testing.T, b *k8sBackend, appID string) {
	t.Helper()

	for i := 0; i < 120; i++ {
		_, err := b.GetContainer(context.Background(), appID)
		if err == domain.ErrContainerNotFound {
			return
		}
		if err != nil {
			t.Fatalf("error in get container: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("wait pod running timeout: %s", appID)
}
