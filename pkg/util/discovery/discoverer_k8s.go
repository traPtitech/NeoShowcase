package discovery

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/retry"
)

type k8sDiscoverer struct {
	client    *kubernetes.Clientset
	namespace string
	svcName   string
}

func NewK8sDiscoverer(svcName string) (Discoverer, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	nsBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return nil, err
	}
	namespace := string(nsBytes)

	d := &k8sDiscoverer{
		client:    client,
		namespace: namespace,
		svcName:   svcName,
	}
	_, err = d.findService()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find service %s, is configuration done right?", svcName)
	}
	return d, nil
}

// findService finds Service itself. Used for configuration sanity check.
func (k *k8sDiscoverer) findService() (*corev1.Service, error) {
	return k.client.CoreV1().Services(k.namespace).Get(context.Background(), k.svcName, metav1.GetOptions{})
}

func (k *k8sDiscoverer) isMe(ep *discoveryv1.Endpoint) bool {
	return ep.TargetRef.Kind == "Pod" &&
		ep.TargetRef.Namespace == k.namespace &&
		ep.TargetRef.Name == os.Getenv("POD_NAME")
}

func (k *k8sDiscoverer) discover() ([]Target, error) {
	// EndpointSlices use a label "kubernetes.io/service-name=<svcName>"
	labelSelector := fmt.Sprintf("kubernetes.io/service-name=%s", k.svcName)

	// List all EndpointSlices in the namespace for that Service
	epSlices, err := k.client.DiscoveryV1().
		EndpointSlices(k.namespace).
		List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
	if err != nil {
		return nil, err
	}

	// Iterate each slice → each endpoint → each address
	var targets []Target
	for _, slice := range epSlices.Items {
		for _, ep := range slice.Endpoints {
			for _, addr := range ep.Addresses {
				// optionally: filter out empty string or non-Ready endpoints
				if addr == "" {
					continue
				}
				// if you only want Ready endpoints:
				if ep.Conditions.Ready != nil && !*ep.Conditions.Ready {
					continue
				}
				targets = append(targets, Target{
					IP: addr,
					Me: k.isMe(&ep),
				})
			}
		}
	}

	err = validateTargets(targets)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(targets, ds.LessFunc(func(e Target) string { return e.IP }))
	log.Infof("[k8s discoverer] Discovered %d targets", len(targets))
	return targets, nil
}

func (k *k8sDiscoverer) watch(ctx context.Context, updates chan<- []Target) error {
	labelSelector := fmt.Sprintf("kubernetes.io/service-name=%s", k.svcName)

	watcher, err := k.client.DiscoveryV1().
		EndpointSlices(k.namespace).
		Watch(ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
		})
	if err != nil {
		return err
	}

	go func() {
		defer watcher.Stop()

		// Send initial state
		res, err := k.discover()
		if err != nil {
			log.Errorf("failed to discover targets: %v", err)
			return
		}
		updates <- res

		for range watcher.ResultChan() {
			res, err = k.discover()
			if err != nil {
				log.Errorf("failed to discover targets: %v", err)
				return
			}
			updates <- res
		}
	}()

	return nil
}

func (k *k8sDiscoverer) Watch(ctx context.Context) (<-chan []Target, error) {
	updates := make(chan []Target)
	go func() {
		defer close(updates)
		retry.Do(ctx, func(ctx context.Context) error {
			return k.watch(ctx, updates)
		}, "k8s discoverer")
	}()
	return updates, nil
}
