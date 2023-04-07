package k8simpl

import (
	"context"
	"encoding/json"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type patcher[T any] interface {
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subResources ...string) (result T, err error)
}

func patch[T any](ctx context.Context, name string, resource T, patcher patcher[T]) error {
	b, err := json.Marshal(resource)
	if err != nil {
		return errors.Wrap(err, "failed to marshal resource")
	}
	_, err = patcher.Patch(ctx, name, types.ApplyPatchType, b, metav1.PatchOptions{Force: lo.ToPtr(true), FieldManager: fieldManager})
	return err
}

func strategicPatch[T any](ctx context.Context, name string, resource any, patcher patcher[T]) error {
	b, err := json.Marshal(resource)
	if err != nil {
		return errors.Wrap(err, "failed to marshal resource")
	}
	_, err = patcher.Patch(ctx, name, types.StrategicMergePatchType, b, metav1.PatchOptions{FieldManager: fieldManager})
	return err
}

type namedResource interface {
	GetName() string
}

func diff[T namedResource](existing []T, current []T) []T {
	var ret []T
	currentMap := lo.SliceToMap(current, func(r T) (string, struct{}) { return r.GetName(), struct{}{} })
	for _, ex := range existing {
		if _, ok := currentMap[ex.GetName()]; !ok {
			ret = append(ret, ex)
		}
	}
	return ret
}

type deleter[T namedResource] interface {
	patcher[T] // embedded just to actually catch type errors
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
}

func prune[T namedResource](ctx context.Context, resources []T, deleter deleter[T]) error {
	for _, resource := range resources {
		if err := deleter.Delete(ctx, resource.GetName(), metav1.DeleteOptions{PropagationPolicy: lo.ToPtr(metav1.DeletePropagationForeground)}); err != nil {
			return err
		}
	}
	return nil
}
