package k8simpl

import (
	"context"
	"encoding/json"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
)

type patcher[T any] interface {
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result T, err error)
}

func patch[T any](ctx context.Context, name string, resource T, patcher patcher[T]) error {
	b, err := json.Marshal(resource)
	if err != nil {
		return errors.Wrap(err, "failed to marshal resource")
	}
	_, err = patcher.Patch(ctx, name, types.ApplyPatchType, b, metav1.PatchOptions{Force: pointer.Bool(true), FieldManager: fieldManager})
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

func diff[T namedResource](existing []T, current []T) []string {
	var ret []string
	currentMap := lo.SliceToMap(current, func(r T) (string, struct{}) { return r.GetName(), struct{}{} })
	for _, ex := range existing {
		if _, ok := currentMap[ex.GetName()]; !ok {
			ret = append(ret, ex.GetName())
		}
	}
	return ret
}

func names[T namedResource](resources []T) []string {
	ret := make([]string, len(resources))
	for i := range resources {
		ret[i] = resources[i].GetName()
	}
	return ret
}

type deleter interface {
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
}

func prune(ctx context.Context, names []string, deleter deleter) error {
	for _, name := range names {
		if err := deleter.Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}
