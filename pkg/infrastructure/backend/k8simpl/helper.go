package k8simpl

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"

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

func setResourceHashAnnotation[T apiResource](rc T, hash string) {
	labels := rc.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[resourceHashAnnotation] = hash
	rc.SetLabels(labels)
}

type syncer[T apiResource] interface {
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subResources ...string) (result T, err error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
}

type deletionNotifier struct {
	mu      sync.Mutex
	waiters map[string]chan struct{}
}

func newDeletionNotifier() *deletionNotifier {
	return &deletionNotifier{
		waiters: make(map[string]chan struct{}),
	}
}

func (n *deletionNotifier) add(name string) <-chan struct{} {
	n.mu.Lock()
	defer n.mu.Unlock()
	ch := make(chan struct{})
	n.waiters[name] = ch
	return ch
}

func (n *deletionNotifier) notify(name string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if ch, ok := n.waiters[name]; ok {
		close(ch)
		delete(n.waiters, name)
	}
}

func syncResources[T apiResource](ctx context.Context, cluster *discovery.Cluster, rcName string, existing []T, next []T, s syncer[T]) error {
	var patched int
	oldHashes := lo.SliceToMap(existing, func(rc T) (string, string) {
		return rc.GetName(), rc.GetLabels()[resourceHashAnnotation]
	})

	// Apply new / existing resources
	for _, rc := range next {
		h, err := hashResource(rc)
		if err != nil {
			return err
		}
		if h == oldHashes[rc.GetName()] {
			continue
		}
		setResourceHashAnnotation(rc, h)

		b, err := marshalResource(rc)
		if err != nil {
			return err
		}
		_, err = s.Patch(ctx, rc.GetName(), types.ApplyPatchType, b, metav1.PatchOptions{Force: lo.ToPtr(true), FieldManager: fieldManager})
		if err != nil {
			slog.ErrorContext(ctx, "failed to patch", "resource", rcName+"/"+rc.GetName(), "error", err)
			continue // skip this resource if patch fails
		}
		patched++
	}

	pruned := pruneResources(ctx, cluster, rcName, existing, next, s)

	if patched > 0 || pruned > 0 {
		slog.DebugContext(ctx, "patched and pruned resources", "resource", rcName, "patched", patched, "pruned", pruned)
	}
	return nil
}

func syncResourcesWithReplace[T apiResource](ctx context.Context, cluster *discovery.Cluster, rcName string, existing []T, next []T, s syncer[T]) error {
	var patched int
	var replaced atomic.Int64
	replacePool := pool.New().WithErrors().WithContext(ctx).WithMaxGoroutines(10)
	oldHashes := lo.SliceToMap(existing, func(rc T) (string, string) {
		return rc.GetName(), rc.GetLabels()[resourceHashAnnotation]
	})

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	notifier := newDeletionNotifier()
	watcher, err := s.Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	defer watcher.Stop()
	go func() {
		for {
			select {
			case event, ok := <-watcher.ResultChan():
				if !ok {
					slog.WarnContext(ctx, "watcher channel closed")
					return
				}
				if event.Type == watch.Deleted {
					obj, ok := event.Object.(apiResource)
					if !ok {
						continue
					}
					notifier.notify(obj.GetName())
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Apply new / existing resources
	for _, rc := range next {
		h, err := hashResource(rc)
		if err != nil {
			return err
		}
		if h == oldHashes[rc.GetName()] {
			continue
		}
		setResourceHashAnnotation(rc, h)

		b, err := marshalResource(rc)
		if err != nil {
			return err
		}
		// For StatefulSets, delete the resource before applying again - StatefulSet has many immutable fields
		// Example: StatefulSet.apps "nsapp-add177a080c4c78936e192" is invalid: spec: Forbidden: updates to statefulset spec for fields other than 'replicas', 'ordinals', 'template', 'updateStrategy', 'revisionHistoryLimit', 'persistentVolumeClaimRetentionPolicy' and 'minReadySeconds' are forbidden
		_, err = s.Patch(ctx, rc.GetName(), types.ApplyPatchType, b, metav1.PatchOptions{Force: lo.ToPtr(true), FieldManager: fieldManager})
		if err != nil {
			// Try to replace the resource if patch fails.
			// This may take a while, so run it in a goroutine.
			replacePool.Go(func(ctx context.Context) error {
				err := replaceResource(ctx, rc, s, b, notifier)
				if err != nil {
					return err
				}
				replaced.Add(1)
				return nil
			})
		} else {
			patched++
		}
	}

	if err := replacePool.Wait(); err != nil {
		slog.ErrorContext(ctx, "error occurred while waiting for replace", "error", err)
		// no return here, continue to prune old resources
	}

	pruned := pruneResources(ctx, cluster, rcName, existing, next, s)

	if patched > 0 || replaced.Load() > 0 || pruned > 0 {
		slog.DebugContext(ctx, "patched, replaced, and pruned resources", "resource", rcName, "patched", patched, "replaced", replaced.Load(), "pruned", pruned)
	}
	return nil
}

func pruneResources[T apiResource](ctx context.Context, cluster *discovery.Cluster, rcName string, existing []T, next []T, s syncer[T]) int {
	pruned := 0
	for _, rc := range diff(existing, next) {
		appID, ok := rc.GetLabels()[appIDLabel]
		if ok && cluster.AssignedShardIndex(appID) != cluster.MyShardIndex() {
			continue
		}
		err := s.Delete(ctx, rc.GetName(), metav1.DeleteOptions{PropagationPolicy: lo.ToPtr(metav1.DeletePropagationForeground)})
		if err != nil {
			slog.ErrorContext(ctx, "failed to delete resource", "resource", rcName+"/"+rc.GetName(), "error", err)
			continue // skip this resource if delete fails
		}
		pruned++
	}
	return pruned
}

func replaceResource[T apiResource](
	ctx context.Context,
	rc T,
	s syncer[T],
	data []byte,
	notifier *deletionNotifier,
) error {
	ch := notifier.add(rc.GetName())
	_ = s.Delete(ctx, rc.GetName(), metav1.DeleteOptions{PropagationPolicy: lo.ToPtr(metav1.DeletePropagationForeground)})

	select {
	case <-ch: // Wait for the resource to be deleted
	// 2 minutes timeout
	// This timeout provides sufficient buffer for most deletion scenarios while
	// preventing indefinite waits for stuck resources.
	case <-time.After(2 * time.Minute):
		return errors.New("timeout while waiting for resource to be deleted")
	case <-ctx.Done():
		return ctx.Err()
	}

	_, err := s.Patch(ctx, rc.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{Force: lo.ToPtr(true), FieldManager: fieldManager})
	return err
}
