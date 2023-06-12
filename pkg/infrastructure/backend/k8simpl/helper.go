package k8simpl

import (
	"context"
	"encoding/json"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

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
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subResources ...string) (result T, err error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
}

func syncResources[T apiResource](ctx context.Context, rcName string, existing []T, next []T, s syncer[T]) error {
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
		labels[resourceHashAnnotation] = h
		rc.SetLabels(labels)

		b, err := marshalResource(rc)
		if err != nil {
			return err
		}
		_, err = s.Patch(ctx, rc.GetName(), types.ApplyPatchType, b, metav1.PatchOptions{Force: lo.ToPtr(true), FieldManager: fieldManager})
		if err != nil {
			return err
		}

		patched++
	}

	// Prune old resources
	for _, rc := range diff(existing, next) {
		err := s.Delete(ctx, rc.GetName(), metav1.DeleteOptions{PropagationPolicy: lo.ToPtr(metav1.DeletePropagationForeground)})
		if err != nil {
			return err
		}
		pruned++
	}

	if patched > 0 || pruned > 0 {
		log.Debugf("patched %v %v, pruned %v %v", patched, rcName, pruned, rcName)
	}
	return nil
}
