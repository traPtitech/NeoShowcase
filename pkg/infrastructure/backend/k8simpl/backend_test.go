package k8simpl

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	ssh2 "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/stretchr/testify/require"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefikio/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

const key = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACAC1iAC54T1ooCQN545XcXDPdTxJEEDdt9TsO3MwoPMwwAAAJCX+efxl/nn
8QAAAAtzc2gtZWQyNTUxOQAAACAC1iAC54T1ooCQN545XcXDPdTxJEEDdt9TsO3MwoPMww
AAAEA+FzwWKIYduEDOqkEOZ2wmxZWPc2wpZeWj+J8e3Q6x0QLWIALnhPWigJA3njldxcM9
1PEkQQN231Ow7czCg8zDAAAADG1vdG9AbW90by13cwE=
-----END OPENSSH PRIVATE KEY-----`

func prepareManager(t *testing.T) (*k8sBackend, *kubernetes.Clientset, *traefikv1alpha1.TraefikV1alpha1Client) {
	const appsNamespace = "neoshowcase-apps"

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
	certManagerClient, err := certmanagerv1.NewForConfig(kubeconf)
	require.NoError(t, err)

	if _, err := client.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: appsNamespace,
		},
	}, metav1.CreateOptions{}); err != nil && !errors.IsAlreadyExists(err) {
		t.Fatal(err)
	}

	sshKey, err := ssh2.NewPublicKeys("", []byte(key), "")
	require.NoError(t, err)
	var config Config
	config.Namespace = appsNamespace
	config.TLS.Type = tlsTypeTraefik
	b, err := NewK8SBackend(kubeconf, client, traefikClient, certManagerClient, config, sshKey, nil, nil)
	require.NoError(t, err)

	err = b.Start(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = b.Dispose(context.Background())
	})

	return b.(*k8sBackend), client, traefikClient
}

func waitCondition(t *testing.T, cond func() (ok bool, err error)) {
	t.Helper()

	for i := 0; i < 120; i++ {
		ok, err := cond()
		if err != nil {
			t.Fatalf(err.Error())
		}
		if ok {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func waitPodRunning(t *testing.T, b *k8sBackend, appID string) {
	t.Helper()
	waitCondition(t, func() (ok bool, err error) {
		status, err := b.GetContainer(context.Background(), appID)
		if err != nil {
			return false, err
		}
		return status.State == domain.ContainerStateRunning, nil
	})
}

func waitPodDeleted(t *testing.T, b *k8sBackend, appID string) {
	t.Helper()
	waitCondition(t, func() (ok bool, err error) {
		status, err := b.GetContainer(context.Background(), appID)
		if err != nil {
			return false, err
		}
		return status.State == domain.ContainerStateMissing, nil
	})
}

type getter[T any] interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (T, error)
}

func exists[T any](t *testing.T, name string, getter getter[T]) {
	t.Helper()
	_, err := getter.Get(context.Background(), name, metav1.GetOptions{})
	require.NoError(t, err)
}

func notExists[T any](t *testing.T, name string, getter getter[T]) {
	t.Helper()
	_, err := getter.Get(context.Background(), name, metav1.GetOptions{})
	require.Error(t, err)
	require.True(t, errors.IsNotFound(err))
}

func waitNotExists[T any](t *testing.T, name string, getter getter[T]) {
	t.Helper()
	waitCondition(t, func() (ok bool, err error) {
		_, err = getter.Get(context.Background(), name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
}
