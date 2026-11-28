// Package k8srunner implements builder.IsolatedRunner on k8s.
// Each job is a Pod whose steps are init containers in order, except for the last step, which is the regular container.
package k8srunner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/samber/oops"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/remotecommand"
	watchtools "k8s.io/client-go/tools/watch"
	"k8s.io/utils/ptr"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner"
	"github.com/traPtitech/neoshowcase/pkg/util/random"
)

const (
	// jobLabel records the ID of the job the resource was created for.
	jobLabel = "ns.trap.jp/isolated-job"
	// hostLabel records the hostname of the builder that created the resource, which is the name of the builder pod.
	hostLabel = "ns.trap.jp/builder-host"

	// receiveDoneMarker is created in the receive container, outside the shared directories, once the input is in place.
	receiveDoneMarker = "/tmp/ns-receive-done"

	cleanupTimeout = 30 * time.Second

	// defaultEndedPodGrace is how long ReclaimEnded leaves an ended build pod for the builder that created it.
	defaultEndedPodGrace = 5 * time.Minute
)

type Config struct {
	// Namespace is where build pods and their secrets are created.
	// It must be dedicated to them, since the builder can read every secret in it through the pod spec.
	Namespace string `mapstructure:"namespace" yaml:"namespace"`
	// Resources apply to every container of a build pod. An empty quantity sets no request or limit for the resource.
	Resources struct {
		Requests struct {
			CPU    string `mapstructure:"cpu" yaml:"cpu"`
			Memory string `mapstructure:"memory" yaml:"memory"`
		} `mapstructure:"requests" yaml:"requests"`
		Limits struct {
			CPU    string `mapstructure:"cpu" yaml:"cpu"`
			Memory string `mapstructure:"memory" yaml:"memory"`
		} `mapstructure:"limits" yaml:"limits"`
	} `mapstructure:"resources" yaml:"resources"`
	// NodeSelector and Tolerations choose the nodes that build pods run on.
	NodeSelector []NodeLabel  `mapstructure:"nodeSelector" yaml:"nodeSelector"`
	Tolerations  []Toleration `mapstructure:"tolerations" yaml:"tolerations"`
}

func (c *Config) validate() error {
	if c.Namespace == "" {
		return oops.New("namespace is required")
	}
	quantities := map[string]string{
		"requests.cpu":    c.Resources.Requests.CPU,
		"requests.memory": c.Resources.Requests.Memory,
		"limits.cpu":      c.Resources.Limits.CPU,
		"limits.memory":   c.Resources.Limits.Memory,
	}
	for name, q := range quantities {
		if q == "" {
			continue
		}
		if _, err := resource.ParseQuantity(q); err != nil {
			return oops.Wrapf(err, "resources.%s: invalid quantity", name)
		}
	}
	if err := validateNodeSelector(c.NodeSelector); err != nil {
		return err
	}
	return validateTolerations(c.Tolerations)
}

func (c *Config) resourceRequirements() v1.ResourceRequirements {
	var r v1.ResourceRequirements
	set := func(list *v1.ResourceList, name v1.ResourceName, q string) {
		if q == "" {
			return
		}
		if *list == nil {
			*list = v1.ResourceList{}
		}
		(*list)[name] = resource.MustParse(q)
	}
	set(&r.Requests, v1.ResourceCPU, c.Resources.Requests.CPU)
	set(&r.Requests, v1.ResourceMemory, c.Resources.Requests.Memory)
	set(&r.Limits, v1.ResourceCPU, c.Resources.Limits.CPU)
	set(&r.Limits, v1.ResourceMemory, c.Resources.Limits.Memory)
	return r
}

var _ builder.IsolatedRunner = (*Runner)(nil)

type Runner struct {
	restConfig     *rest.Config
	client         kubernetes.Interface
	config         Config
	hostname       string
	pullPolicy     isolatedrunner.PullPolicy
	activeDeadline time.Duration
	endedPodGrace  time.Duration
}

// Options are the settings of a Runner other than its clients.
type Options struct {
	Config Config
	// Hostname labels the build pods and must be unique among the builders sharing the namespace.
	Hostname   string
	PullPolicy isolatedrunner.PullPolicy
	// ActiveDeadline kills a build pod even if the builder stops watching it.
	ActiveDeadline time.Duration
}

// NewRunner fails if the namespace is empty, a resource quantity cannot be parsed, the API would reject the node selector or a toleration,
// the pull policy is unknown, or the active deadline is shorter than a second.
func NewRunner(restConfig *rest.Config, client kubernetes.Interface, opts Options) (*Runner, error) {
	if err := opts.Config.validate(); err != nil {
		return nil, err
	}
	if err := opts.PullPolicy.Validate(); err != nil {
		return nil, err
	}
	// The deadline of a pod is in whole seconds, and the API rejects 0.
	if opts.ActiveDeadline < time.Second {
		return nil, oops.Errorf("active deadline %s is shorter than a second", opts.ActiveDeadline)
	}
	return &Runner{
		restConfig:     restConfig,
		client:         client,
		config:         opts.Config,
		hostname:       opts.Hostname,
		pullPolicy:     opts.PullPolicy,
		activeDeadline: opts.ActiveDeadline,
		endedPodGrace:  defaultEndedPodGrace,
	}, nil
}

func (r *Runner) labels(job builder.IsolatedJob) map[string]string {
	return map[string]string{
		jobLabel:  job.ID,
		hostLabel: r.hostname,
	}
}

func secretKey(step builder.IsolatedStep, envName string) string {
	return step.Name + "." + envName
}

func (r *Runner) podSpec(name string, job builder.IsolatedJob) *v1.Pod {
	var volumes []v1.Volume
	var mounts []v1.VolumeMount
	for i, dir := range job.Shared {
		volName := fmt.Sprintf("shared-%d", i)
		volumes = append(volumes, v1.Volume{Name: volName, VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}}})
		mounts = append(mounts, v1.VolumeMount{Name: volName, MountPath: dir})
	}

	resources := r.config.resourceRequirements()
	receive := v1.Container{
		Name:            isolatedrunner.ReceiveStepName,
		Image:           job.Image,
		ImagePullPolicy: v1.PullPolicy(r.pullPolicy),
		Command:         []string{"sh", "-c", fmt.Sprintf("until [ -f %s ]; do sleep 0.2; done", receiveDoneMarker)},
		VolumeMounts:    mounts,
		Resources:       resources,
		SecurityContext: &v1.SecurityContext{RunAsUser: ptr.To[int64](0), RunAsGroup: ptr.To[int64](0)},
	}
	containers := []v1.Container{receive}
	for _, step := range job.Steps {
		var env []v1.EnvVar
		for _, k := range slices.Sorted(maps.Keys(step.Env)) {
			env = append(env, v1.EnvVar{Name: k, Value: step.Env[k]})
		}
		for _, k := range slices.Sorted(maps.Keys(step.SecretEnv)) {
			env = append(env, v1.EnvVar{Name: k, ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{Name: name},
				Key:                  secretKey(step, k),
			}}})
		}
		containers = append(containers, v1.Container{
			Name:            step.Name,
			Image:           job.Image,
			ImagePullPolicy: v1.PullPolicy(r.pullPolicy),
			Command:         step.Cmd,
			WorkingDir:      step.WorkDir,
			Env:             env,
			VolumeMounts:    mounts,
			Resources:       resources,
			SecurityContext: &v1.SecurityContext{
				RunAsUser:  ptr.To(int64(job.RunAs.UID)),
				RunAsGroup: ptr.To(int64(job.RunAs.GID)),
			},
		})
	}

	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: r.config.Namespace,
			Labels:    r.labels(job),
		},
		Spec: v1.PodSpec{
			RestartPolicy:                 v1.RestartPolicyNever,
			AutomountServiceAccountToken:  ptr.To(false),
			EnableServiceLinks:            ptr.To(false),
			ActiveDeadlineSeconds:         ptr.To(int64(r.activeDeadline.Seconds())),
			TerminationGracePeriodSeconds: ptr.To[int64](5),
			NodeSelector:                  podNodeSelector(r.config.NodeSelector),
			Tolerations:                   podTolerations(r.config.Tolerations),
			Volumes:                       volumes,
			InitContainers:                containers[:len(containers)-1],
			Containers:                    containers[len(containers)-1:],
		},
	}
}

func (r *Runner) secretSpec(pod *v1.Pod, job builder.IsolatedJob) *v1.Secret {
	data := map[string]string{}
	for _, step := range job.Steps {
		for k, v := range step.SecretEnv {
			data[secretKey(step, k)] = v
		}
	}
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Labels:    r.labels(job),
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: "v1",
				Kind:       "Pod",
				Name:       pod.Name,
				UID:        pod.UID,
			}},
		},
		StringData: data,
	}
}

func hasSecretEnv(job builder.IsolatedJob) bool {
	return slices.ContainsFunc(job.Steps, func(s builder.IsolatedStep) bool { return len(s.SecretEnv) > 0 })
}

func (r *Runner) Run(ctx context.Context, job builder.IsolatedJob, logWriter io.Writer) error {
	if err := isolatedrunner.ValidateJob(job); err != nil {
		return err
	}

	name := "ns-build-" + job.ID + "-" + random.SecureGenerateHex(8)
	pods := r.client.CoreV1().Pods(r.config.Namespace)
	pod, err := pods.Create(ctx, r.podSpec(name, job), metav1.CreateOptions{})
	if err != nil {
		return oops.Wrapf(err, "creating build pod")
	}
	createdSecret := false
	defer func() { r.cleanup(context.WithoutCancel(ctx), name, createdSecret) }()

	// The secret is created after the pod so that it can be owned by the pod.
	// The containers referring to it start only after the receive step, which waits for the input sent below.
	if hasSecretEnv(job) {
		_, err := r.client.CoreV1().Secrets(r.config.Namespace).Create(ctx, r.secretSpec(pod, job), metav1.CreateOptions{})
		if err != nil {
			return oops.Wrapf(err, "creating build secret")
		}
		createdSecret = true
	}

	if err := r.receive(ctx, name, job, logWriter); err != nil {
		return err
	}
	for _, step := range job.Steps {
		isolatedrunner.WriteStepStart(logWriter, step.Name)
		if err := r.followStep(ctx, name, step.Name, logWriter); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) receive(ctx context.Context, podName string, job builder.IsolatedJob, logWriter io.Writer) error {
	_, err := r.waitPod(ctx, podName, stepRunning(isolatedrunner.ReceiveStepName))
	if err != nil {
		return oops.Wrapf(err, "waiting for step %s to start", isolatedrunner.ReceiveStepName)
	}

	req := r.client.CoreV1().RESTClient().Post().
		Resource("pods").Name(podName).Namespace(r.config.Namespace).SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: isolatedrunner.ReceiveStepName,
			Command:   []string{"sh", "-c", isolatedrunner.ReceiveScript(job) + " && touch " + receiveDoneMarker},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
		}, scheme.ParameterCodec)
	ex, err := remotecommand.NewSPDYExecutor(r.restConfig, "POST", req.URL())
	if err != nil {
		return oops.Wrapf(err, "creating SPDY executor")
	}
	input := &errRecordingReader{r: job.Input}
	err = ex.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  input,
		Stdout: logWriter,
		Stderr: logWriter,
	})
	// The executor closes stdin on a read error as if the input had ended, and tar may then succeed with missing files.
	if readErr := input.Err(); readErr != nil {
		return oops.Wrapf(readErr, "step %s: reading input", isolatedrunner.ReceiveStepName)
	}
	if err != nil {
		return oops.Wrapf(err, "step %s", isolatedrunner.ReceiveStepName)
	}
	return nil
}

type errRecordingReader struct {
	r   io.Reader
	mu  sync.Mutex
	err error
}

func (e *errRecordingReader) Read(p []byte) (int, error) {
	n, err := e.r.Read(p)
	if err != nil && !errors.Is(err, io.EOF) {
		e.mu.Lock()
		e.err = err
		e.mu.Unlock()
	}
	return n, err
}

// Err returns the error that reading failed with, other than io.EOF.
func (e *errRecordingReader) Err() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.err
}

// followStep writes the log of the step until the step exits, and returns an error if it exits with a non-zero code.
func (r *Runner) followStep(ctx context.Context, podName, stepName string, logWriter io.Writer) error {
	_, err := r.waitPod(ctx, podName, stepStarted(stepName))
	if err != nil {
		return oops.Wrapf(err, "waiting for step %s to start", stepName)
	}

	logs, err := r.client.CoreV1().Pods(r.config.Namespace).GetLogs(podName, &v1.PodLogOptions{Container: stepName, Follow: true}).Stream(ctx)
	if err != nil {
		return oops.Wrapf(err, "streaming log of step %s", stepName)
	}
	_, err = io.Copy(logWriter, logs)
	_ = logs.Close()
	if err != nil {
		return oops.Wrapf(err, "streaming log of step %s", stepName)
	}

	pod, err := r.waitPod(ctx, podName, stepExited(stepName))
	if err != nil {
		return oops.Wrapf(err, "waiting for step %s to exit", stepName)
	}
	if code := containerStatus(pod, stepName).State.Terminated.ExitCode; code != 0 {
		return oops.Errorf("step %s exited with code %d", stepName, code)
	}
	return nil
}

func containerStatus(pod *v1.Pod, name string) *v1.ContainerStatus {
	for _, statuses := range [][]v1.ContainerStatus{pod.Status.InitContainerStatuses, pod.Status.ContainerStatuses} {
		for i := range statuses {
			if statuses[i].Name == name {
				return &statuses[i]
			}
		}
	}
	return nil
}

func stepRunning(name string) func(pod *v1.Pod) bool {
	return func(pod *v1.Pod) bool {
		s := containerStatus(pod, name)
		return s != nil && s.State.Running != nil
	}
}

// stepStarted holds once the step is running or has exited.
func stepStarted(name string) func(pod *v1.Pod) bool {
	return func(pod *v1.Pod) bool {
		s := containerStatus(pod, name)
		return s != nil && (s.State.Running != nil || s.State.Terminated != nil)
	}
}

func stepExited(name string) func(pod *v1.Pod) bool {
	return func(pod *v1.Pod) bool {
		s := containerStatus(pod, name)
		return s != nil && s.State.Terminated != nil
	}
}

// failFastReasons are waiting reasons of a container that do not resolve by waiting.
var failFastReasons = []string{"ErrImagePull", "ImagePullBackOff", "ErrImageNeverPull", "InvalidImageName"}

// podFailure returns why the pod cannot proceed, or nil if it may still proceed.
func podFailure(pod *v1.Pod) error {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == v1.PodScheduled && cond.Status == v1.ConditionFalse && cond.Reason == v1.PodReasonUnschedulable {
			return oops.Errorf("build pod is unschedulable: %s", cond.Message)
		}
	}
	for _, statuses := range [][]v1.ContainerStatus{pod.Status.InitContainerStatuses, pod.Status.ContainerStatuses} {
		for _, s := range statuses {
			if s.State.Waiting != nil && slices.Contains(failFastReasons, s.State.Waiting.Reason) {
				return oops.Errorf("container %s: %s: %s", s.Name, s.State.Waiting.Reason, s.State.Waiting.Message)
			}
		}
	}
	return nil
}

// podReached reports whether cond holds for the pod, or fails when the pod can no longer reach it.
func podReached(pod *v1.Pod, cond func(pod *v1.Pod) bool) (bool, error) {
	if cond(pod) {
		return true, nil
	}
	if err := podFailure(pod); err != nil {
		return false, err
	}
	if pod.Status.Phase == v1.PodFailed || pod.Status.Phase == v1.PodSucceeded {
		return false, oops.Errorf("build pod ended (%s): %s %s", pod.Status.Phase, pod.Status.Reason, pod.Status.Message)
	}
	return false, nil
}

// waitPod waits until cond holds for the pod, and fails when the pod cannot reach it.
func (r *Runner) waitPod(ctx context.Context, podName string, cond func(pod *v1.Pod) bool) (*v1.Pod, error) {
	lw := cache.NewListWatchFromClient(r.client.CoreV1().RESTClient(), "pods", r.config.Namespace, fields.OneTermEqualSelector("metadata.name", podName))
	// The pod is created before it is waited for, so its absence means that it was deleted, which no later event reports.
	exists := func(store cache.Store) (bool, error) {
		if len(store.List()) == 0 {
			return false, oops.New("build pod was deleted")
		}
		return false, nil
	}
	ev, err := watchtools.UntilWithSync(ctx, lw, &v1.Pod{}, exists, func(ev watch.Event) (bool, error) {
		if ev.Type == watch.Deleted {
			return false, oops.New("build pod was deleted")
		}
		pod, ok := ev.Object.(*v1.Pod)
		if !ok {
			return false, nil
		}
		return podReached(pod, cond)
	})
	if err != nil {
		// UntilWithSync reports the end of ctx with an error of its own, which does not wrap ctx.Err().
		if ctx.Err() != nil {
			return nil, oops.Wrapf(ctx.Err(), "waiting for the build pod")
		}
		return nil, err
	}
	return ev.Object.(*v1.Pod), nil
}

// cleanup requests deletion of the pod and the secret without waiting for it to complete.
func (r *Runner) cleanup(ctx context.Context, name string, deleteSecret bool) {
	ctx, cancel := context.WithTimeout(ctx, cleanupTimeout)
	defer cancel()
	if deleteSecret {
		// The secret is garbage collected with the pod, but deleting it directly removes it without waiting for the pod to terminate.
		err := r.client.CoreV1().Secrets(r.config.Namespace).Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil && !k8serrors.IsNotFound(err) {
			slog.ErrorContext(ctx, "failed to delete build secret", "name", name, "error", err)
		}
	}
	err := r.client.CoreV1().Pods(r.config.Namespace).Delete(ctx, name, metav1.DeleteOptions{
		PropagationPolicy: ptr.To(metav1.DeletePropagationBackground),
	})
	if err != nil && !k8serrors.IsNotFound(err) {
		slog.ErrorContext(ctx, "failed to delete build pod", "name", name, "error", err)
	}
}

func (r *Runner) ReclaimOrphans(ctx context.Context) error {
	err := r.client.CoreV1().Pods(r.config.Namespace).DeleteCollection(ctx,
		metav1.DeleteOptions{PropagationPolicy: ptr.To(metav1.DeletePropagationBackground)},
		metav1.ListOptions{LabelSelector: hostLabel + "=" + r.hostname},
	)
	if err != nil {
		return oops.Wrapf(err, "deleting orphaned build pods")
	}
	return nil
}

// ReclaimEnded deletes the build pods in the namespace that have ended, whichever builder created them.
// A build pod whose builder is gone ends by its active deadline. Its secret is garbage collected with it.
func (r *Runner) ReclaimEnded(ctx context.Context) error {
	pods, err := r.client.CoreV1().Pods(r.config.Namespace).List(ctx, metav1.ListOptions{LabelSelector: jobLabel})
	if err != nil {
		return oops.Wrapf(err, "listing build pods")
	}
	now := time.Now()
	for _, pod := range pods.Items {
		if pod.DeletionTimestamp != nil || (pod.Status.Phase != v1.PodSucceeded && pod.Status.Phase != v1.PodFailed) {
			continue
		}
		if now.Sub(endedAt(&pod)) < r.endedPodGrace {
			continue
		}
		err := r.client.CoreV1().Pods(r.config.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{
			PropagationPolicy: ptr.To(metav1.DeletePropagationBackground),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return oops.Wrapf(err, "deleting ended build pod %s", pod.Name)
		}
	}
	return nil
}

// endedAt returns when the last container of the ended pod exited, or when the pod was created if no container ran.
func endedAt(pod *v1.Pod) time.Time {
	t := pod.CreationTimestamp.Time
	for _, statuses := range [][]v1.ContainerStatus{pod.Status.InitContainerStatuses, pod.Status.ContainerStatuses} {
		for _, s := range statuses {
			if s.State.Terminated != nil && s.State.Terminated.FinishedAt.After(t) {
				t = s.State.Terminated.FinishedAt.Time
			}
		}
	}
	return t
}
