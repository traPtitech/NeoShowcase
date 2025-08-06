package k8simpl

import (
	"context"
	"encoding/json"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/traPtitech/neoshowcase/pkg/util/discovery"
	"github.com/traPtitech/neoshowcase/pkg/util/hash"
)

type apiResource interface {
	GetName() string
	GetLabels() map[string]string
	SetLabels(labels map[string]string)
}

// diff returns A - B
func diff[T apiResource](a []T, b []T) []T {
	var ret []T
	bMap := lo.SliceToMap(b, func(r T) (string, struct{}) { return r.GetName(), struct{}{} })
	for _, aa := range a {
		if _, ok := bMap[aa.GetName()]; !ok {
			ret = append(ret, aa)
		}
	}
	return ret
}

func marshalResource(rc apiResource) ([]byte, error) {
	return json.Marshal(rc)
}

func hashResource(rc apiResource) (string, error) {
	b, err := marshalResource(rc)
	if err != nil {
		return "", err
	}
	return hash.XXH3Hex(b), nil
}

type syncer[T apiResource] interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (T, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subResources ...string) (result T, err error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
}

func syncResources[T apiResource](ctx context.Context, cluster *discovery.Cluster, rcName string, existing []T, next []T, s syncer[T], replace bool) error {
	var patched, pruned int
	oldHashes := lo.SliceToMap(existing, func(rc T) (string, string) {
		return rc.GetName(), rc.GetLabels()[resourceHashAnnotation]
	})

	// Apply new / existing resources
	for _, rc := range next {
		// Compute hash before applying annotation
		h, err := hashResource(rc)
		if err != nil {
			return err
		}
		if h == oldHashes[rc.GetName()] {
			// No need to apply
			continue
		}

		// Set label, and apply
		labels := rc.GetLabels()
		if labels == nil {
			labels = make(map[string]string)
		}
		labels[resourceHashAnnotation] = h
		rc.SetLabels(labels)

		b, err := marshalResource(rc)
		if err != nil {
			return err
		}
		_, err = s.Patch(ctx, rc.GetName(), types.ApplyPatchType, b, metav1.PatchOptions{Force: lo.ToPtr(true), FieldManager: fieldManager})
		// For StatefulSets, delete the resource before applying again - StatefulSet has many immutable fields
		// Example: StatefulSet.apps "nsapp-add177a080c4c78936e192" is invalid: spec: Forbidden: updates to statefulset spec for fields other than 'replicas', 'ordinals', 'template', 'updateStrategy', 'revisionHistoryLimit', 'persistentVolumeClaimRetentionPolicy' and 'minReadySeconds' are forbidden
		if err != nil && replace {
			err = replaceResource(ctx, rc, b, s)
		}
		if err != nil {
			log.WithError(err).Errorf("failed to patch %s/%s", rcName, rc.GetName())
			continue // Skip if error occurred
		}

		patched++
	}

	// Prune old resources
	for _, rc := range diff(existing, next) {
		// On the controller cluster scale-up, ensure this shard does not delete resources controlled by a new shard
		appID, ok := rc.GetLabels()[appIDLabel]
		if ok && cluster.AssignedShardIndex(appID) != cluster.MyShardIndex() {
			continue
		}

		err := s.Delete(ctx, rc.GetName(), metav1.DeleteOptions{PropagationPolicy: lo.ToPtr(metav1.DeletePropagationForeground)})
		if err != nil {
			log.WithError(err).Errorf("failed to delete %s/%s", rcName, rc.GetName())
			continue // Skip if error occurred
		}
		pruned++
	}

	if patched > 0 || pruned > 0 {
		log.Debugf("patched %v %v, pruned %v %v", patched, rcName, pruned, rcName)
	}
	return nil
}

func replaceResource[T apiResource](ctx context.Context, rc T, data []byte, s syncer[T]) error {
	err := s.Delete(ctx, rc.GetName(), metav1.DeleteOptions{PropagationPolicy: lo.ToPtr(metav1.DeletePropagationForeground)})
	if err != nil {
		return errors.Wrapf(err, "delete %s/%s before re-apply", rc.GetName(), rc.GetLabels()[resourceHashAnnotation])
	}

	// Wait for the resource to be deleted
	interval := 100 * time.Millisecond
	timeout := 30 * time.Second
	err = wait.PollUntilContextTimeout(ctx, interval, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := s.Get(ctx, rc.GetName(), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			return true, nil // Resource is deleted
		}
		return false, err
	})
	if err != nil {
		return errors.Wrapf(err, "waiting for %s/%s to be deleted", rc.GetName(), rc.GetLabels()[resourceHashAnnotation])
	}

	// Try again
	_, err = s.Patch(ctx, rc.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{Force: lo.ToPtr(true), FieldManager: fieldManager})
	if err != nil {
		return errors.Wrapf(err, "re-apply %s/%s after deletion", rc.GetName(), rc.GetLabels()[resourceHashAnnotation])
	}

	return nil
}
