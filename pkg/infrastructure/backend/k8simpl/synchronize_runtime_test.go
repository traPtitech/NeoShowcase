package k8simpl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/discovery"
)

func newSpecTestBackend(t *testing.T) *Backend {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	c := discovery.NewCluster(discovery.NewSingleDiscoverer("127.0.0.1"))
	go c.Start(ctx)
	return &Backend{config: Config{Namespace: "neoshowcase-apps"}, cluster: c}
}

func loadBalancerServices(rsc *resources) []*v1.Service {
	var svcs []*v1.Service
	for _, svc := range rsc.services {
		if svc.Spec.Type == v1.ServiceTypeLoadBalancer {
			svcs = append(svcs, svc)
		}
	}
	return svcs
}

func TestBackend_runtimeResources_portPublications(t *testing.T) {
	b := newSpecTestBackend(t)

	t.Run("1アプリの公開ポートを1つのLoadBalancer Serviceにまとめる", func(t *testing.T) {
		app := &domain.Application{
			ID:     "app1",
			Config: domain.ApplicationConfig{BuildConfig: &domain.BuildConfigRuntimeBuildpack{}},
			PortPublications: []*domain.PortPublication{
				{InternetPort: 25565, ApplicationPort: 8080, Protocol: domain.PortPublicationProtocolTCP},
				{InternetPort: 25565, ApplicationPort: 8080, Protocol: domain.PortPublicationProtocolUDP},
				{InternetPort: 25566, ApplicationPort: 9090, Protocol: domain.PortPublicationProtocolTCP},
			},
		}
		var rsc resources
		b.runtimeResources(&rsc, []*domain.RuntimeDesiredState{{App: app}})

		svcs := loadBalancerServices(&rsc)
		require.Len(t, svcs, 1)
		svc := svcs[0]
		assert.Equal(t, portServiceName(app.ID), svc.Name)
		assert.Equal(t, appSelector(app.ID), svc.Spec.Selector)
		assert.Equal(t, []v1.ServicePort{
			{Name: "tcp-25565", Protocol: v1.ProtocolTCP, Port: 25565, TargetPort: intstr.FromInt(8080)},
			{Name: "udp-25565", Protocol: v1.ProtocolUDP, Port: 25565, TargetPort: intstr.FromInt(8080)},
			{Name: "tcp-25566", Protocol: v1.ProtocolTCP, Port: 25566, TargetPort: intstr.FromInt(9090)},
		}, svc.Spec.Ports)
	})

	t.Run("公開ポートがないアプリにはLoadBalancer Serviceを作らない", func(t *testing.T) {
		app := &domain.Application{
			ID:     "app2",
			Config: domain.ApplicationConfig{BuildConfig: &domain.BuildConfigRuntimeBuildpack{}},
		}
		var rsc resources
		b.runtimeResources(&rsc, []*domain.RuntimeDesiredState{{App: app}})

		assert.Empty(t, loadBalancerServices(&rsc))
	})

	t.Run("アプリごとに別のLoadBalancer Serviceを作る", func(t *testing.T) {
		apps := []*domain.RuntimeDesiredState{
			{App: &domain.Application{
				ID:               "app3",
				Config:           domain.ApplicationConfig{BuildConfig: &domain.BuildConfigRuntimeBuildpack{}},
				PortPublications: []*domain.PortPublication{{InternetPort: 30000, ApplicationPort: 80, Protocol: domain.PortPublicationProtocolTCP}},
			}},
			{App: &domain.Application{
				ID:               "app4",
				Config:           domain.ApplicationConfig{BuildConfig: &domain.BuildConfigRuntimeBuildpack{}},
				PortPublications: []*domain.PortPublication{{InternetPort: 30001, ApplicationPort: 80, Protocol: domain.PortPublicationProtocolTCP}},
			}},
		}
		var rsc resources
		b.runtimeResources(&rsc, apps)

		svcs := loadBalancerServices(&rsc)
		require.Len(t, svcs, 2)
		assert.Equal(t, portServiceName("app3"), svcs[0].Name)
		assert.Equal(t, portServiceName("app4"), svcs[1].Name)
	})
}

// Sablier の blocking で起動を待ったリクエストが、起動したアプリに届くための設定。
func TestBackend_runtimeResources_routingToScaledFromZero(t *testing.T) {
	b := newSpecTestBackend(t)
	app := &domain.Application{
		ID:       "app1",
		Config:   domain.ApplicationConfig{BuildConfig: &domain.BuildConfigRuntimeBuildpack{}},
		Websites: []*domain.Website{{ID: "web1", FQDN: "app1.example.com", PathPrefix: "/", HTTPPort: 8080}},
	}
	var rsc resources
	b.runtimeResources(&rsc, []*domain.RuntimeDesiredState{{App: app}})

	t.Run("IngressRouteはPodのIPではなくServiceのClusterIPに転送する", func(t *testing.T) {
		require.Len(t, rsc.ingressRoutes, 1)
		services := rsc.ingressRoutes[0].Spec.Routes[0].Services
		require.Len(t, services, 1)
		assert.Equal(t, new(true), services[0].NativeLB)
	})

	t.Run("ServiceはReadyになる前のPodも転送先に含める", func(t *testing.T) {
		require.Len(t, rsc.services, 1)
		assert.Equal(t, v1.ServiceTypeClusterIP, rsc.services[0].Spec.Type)
		assert.True(t, rsc.services[0].Spec.PublishNotReadyAddresses)
	})
}
