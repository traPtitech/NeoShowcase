package k8simpl

import (
	"context"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/friendsofgo/errors"
	traefikv1alpha1 "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

type resources struct {
	statefulSets  []*appsv1.StatefulSet
	secrets       []*v1.Secret
	services      []*v1.Service
	middlewares   []*traefikv1alpha1.Middleware
	ingressRoutes []*traefikv1alpha1.IngressRoute
}

type sharedResources struct {
	certificates []*certmanagerv1.Certificate
}

func (b *Backend) listCurrentResources(ctx context.Context) (*resources, error) {
	var rsc resources
	listOpt := metav1.ListOptions{LabelSelector: toSelectorString(b.shardedAllSelector())}

	ss, err := b.client.AppsV1().StatefulSets(b.config.Namespace).List(ctx, listOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stateful sets")
	}
	rsc.statefulSets = ds.SliceOfPtr(ss.Items)

	secrets, err := b.client.CoreV1().Secrets(b.config.Namespace).List(ctx, listOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get secrets")
	}
	rsc.secrets = ds.SliceOfPtr(secrets.Items)

	svc, err := b.client.CoreV1().Services(b.config.Namespace).List(ctx, listOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get services")
	}
	rsc.services = ds.SliceOfPtr(svc.Items)

	mw, err := b.traefikClient.Middlewares(b.config.Namespace).List(ctx, listOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get middlewares")
	}
	rsc.middlewares = ds.SliceOfPtr(mw.Items)

	ir, err := b.traefikClient.IngressRoutes(b.config.Namespace).List(ctx, listOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ingress routes")
	}
	rsc.ingressRoutes = ds.SliceOfPtr(ir.Items)

	return &rsc, nil
}

func (b *Backend) listCurrentSharedResources(ctx context.Context) (*sharedResources, error) {
	var rsc sharedResources
	listOpt := metav1.ListOptions{LabelSelector: toSelectorString(allSelector())}

	if b.config.TLS.Type == tlsTypeCertManager {
		certs, err := b.certManagerClient.CertmanagerV1().Certificates(b.config.Namespace).List(ctx, listOpt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get certificates")
		}
		rsc.certificates = ds.SliceOfPtr(certs.Items)
	}

	return &rsc, nil
}

func (b *Backend) Synchronize(ctx context.Context, s *domain.DesiredState) error {
	b.reloadLock.Lock()
	defer b.reloadLock.Unlock()

	// Calculate next resources to apply
	var next resources
	b.runtimeResources(&next, s.Runtime)
	b.ssResources(&next, s.StaticSites)

	// List old resources
	old, err := b.listCurrentResources(ctx)
	if err != nil {
		return err
	}

	// Synchronize resources
	err = syncResources[*appsv1.StatefulSet](ctx, b.cluster, "statefulsets", old.statefulSets, next.statefulSets, b.client.AppsV1().StatefulSets(b.config.Namespace), true)
	if err != nil {
		return errors.Wrap(err, "failed to sync stateful sets")
	}
	err = syncResources[*v1.Secret](ctx, b.cluster, "secrets", old.secrets, next.secrets, b.client.CoreV1().Secrets(b.config.Namespace), false)
	if err != nil {
		return errors.Wrap(err, "failed to sync secrets")
	}
	err = syncResources[*v1.Service](ctx, b.cluster, "services", old.services, next.services, b.client.CoreV1().Services(b.config.Namespace), false)
	if err != nil {
		return errors.Wrap(err, "failed to sync services")
	}
	err = syncResources[*traefikv1alpha1.Middleware](ctx, b.cluster, "middlewares", old.middlewares, next.middlewares, b.traefikClient.Middlewares(b.config.Namespace), false)
	if err != nil {
		return errors.Wrap(err, "failed to sync middlewares")
	}
	err = syncResources[*traefikv1alpha1.IngressRoute](ctx, b.cluster, "ingressroutes", old.ingressRoutes, next.ingressRoutes, b.traefikClient.IngressRoutes(b.config.Namespace), false)
	if err != nil {
		return errors.Wrap(err, "failed to sync ingressroutes")
	}

	return nil
}

func (b *Backend) SynchronizeShared(ctx context.Context, s *domain.DesiredStateLeader) error {
	// Calculate next resources to apply
	var next sharedResources
	next.certificates = ds.Map(s.TLSTargetDomains, b.certificate)

	// List old resources
	old, err := b.listCurrentSharedResources(ctx)
	if err != nil {
		return err
	}

	// Synchronize resources
	err = syncResources[*certmanagerv1.Certificate](ctx, b.cluster, "certificates", old.certificates, next.certificates, b.certManagerClient.CertmanagerV1().Certificates(b.config.Namespace), false)
	if err != nil {
		return errors.Wrap(err, "failed to sync certificates")
	}

	return nil
}
