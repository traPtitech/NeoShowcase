package k8srunner

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/utils/ptr"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner/runnertest"
	"github.com/traPtitech/neoshowcase/pkg/util/random"
)

const (
	// builderNamespace stands in for the namespace the builder runs in, which holds secrets of other components.
	builderNamespace = "ns-isolated-runner-test"
	// buildNamespace is the namespace dedicated to build pods.
	buildNamespace = "ns-isolated-runner-test-build"
)

func prepareClient(t *testing.T) (*rest.Config, *kubernetes.Clientset) {
	t.Helper()
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_K8S_TESTS")); !ok {
		t.SkipNow()
	}

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

	for _, ns := range []string{builderNamespace, buildNamespace} {
		_, err = client.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: ns},
		}, metav1.CreateOptions{})
		if err != nil && !errors.IsAlreadyExists(err) {
			t.Fatal(err)
		}
	}
	return kubeconf, client
}

// builderRole has the same rules as the builder Role in charts/neoshowcase/templates/builder/role.yaml.
var builderRole = []rbacv1.PolicyRule{
	{APIGroups: []string{""}, Resources: []string{"pods"}, Verbs: []string{"create", "get", "list", "watch", "delete", "deletecollection"}},
	{APIGroups: []string{""}, Resources: []string{"pods/log"}, Verbs: []string{"get"}},
	{APIGroups: []string{""}, Resources: []string{"pods/exec"}, Verbs: []string{"create", "get"}},
	{APIGroups: []string{""}, Resources: []string{"secrets"}, Verbs: []string{"create", "delete"}},
}

// builderClient returns a client acting as a ServiceAccount in builderNamespace,
// bound to builderRole in buildNamespace as the chart does.
func builderClient(t *testing.T, kubeconf *rest.Config, client *kubernetes.Clientset) (*rest.Config, *kubernetes.Clientset) {
	t.Helper()
	return serviceAccountClient(t, kubeconf, client, "ns-builder", builderRole)
}

// serviceAccountClient returns a client acting as the ServiceAccount of the name in builderNamespace,
// bound to a Role with the rules in buildNamespace.
func serviceAccountClient(t *testing.T, kubeconf *rest.Config, client *kubernetes.Clientset, name string, rules []rbacv1.PolicyRule) (*rest.Config, *kubernetes.Clientset) {
	t.Helper()
	ctx := context.Background()
	ignoreExists := func(_ any, err error) {
		if err != nil && !errors.IsAlreadyExists(err) {
			t.Fatal(err)
		}
	}
	ignoreExists(client.CoreV1().ServiceAccounts(builderNamespace).Create(ctx,
		&v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: name}}, metav1.CreateOptions{}))
	_, err := client.RbacV1().Roles(buildNamespace).Update(ctx,
		&rbacv1.Role{ObjectMeta: metav1.ObjectMeta{Name: name}, Rules: rules}, metav1.UpdateOptions{})
	if errors.IsNotFound(err) {
		_, err = client.RbacV1().Roles(buildNamespace).Create(ctx,
			&rbacv1.Role{ObjectMeta: metav1.ObjectMeta{Name: name}, Rules: rules}, metav1.CreateOptions{})
	}
	require.NoError(t, err)
	ignoreExists(client.RbacV1().RoleBindings(buildNamespace).Create(ctx, &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		RoleRef:    rbacv1.RoleRef{APIGroup: "rbac.authorization.k8s.io", Kind: "Role", Name: name},
		Subjects:   []rbacv1.Subject{{Kind: "ServiceAccount", Name: name, Namespace: builderNamespace}},
	}, metav1.CreateOptions{}))

	conf := rest.CopyConfig(kubeconf)
	conf.Impersonate.UserName = "system:serviceaccount:" + builderNamespace + ":" + name
	c, err := kubernetes.NewForConfig(conf)
	require.NoError(t, err)
	return conf, c
}

// newTestRunner returns a runner with the permissions of the builder, and an admin client to inspect the cluster.
func newTestRunner(t *testing.T) (*Runner, *kubernetes.Clientset) {
	t.Helper()
	kubeconf, client := prepareClient(t)
	bConf, bClient := builderClient(t, kubeconf, client)
	hostname := "fake-builder-" + random.SecureGenerateHex(4)
	r, err := NewRunner(bConf, bClient, Options{Config: Config{Namespace: buildNamespace}, Hostname: hostname, PullPolicy: isolatedrunner.PullIfNotPresent, ActiveDeadline: 10 * time.Minute})
	require.NoError(t, err)
	return r, client
}

// liveObjects returns the names of the pods and secrets in buildNamespace with the given label that are not being deleted.
func liveObjects(t *testing.T, client *kubernetes.Clientset, selector string) []string {
	t.Helper()
	var names []string
	pods, err := client.CoreV1().Pods(buildNamespace).List(context.Background(), metav1.ListOptions{LabelSelector: selector})
	require.NoError(t, err)
	for _, p := range pods.Items {
		if p.DeletionTimestamp == nil {
			names = append(names, "pod/"+p.Name)
		}
	}
	secrets, err := client.CoreV1().Secrets(buildNamespace).List(context.Background(), metav1.ListOptions{LabelSelector: selector})
	require.NoError(t, err)
	for _, s := range secrets.Items {
		if s.DeletionTimestamp == nil {
			names = append(names, "secret/"+s.Name)
		}
	}
	return names
}

func TestRunner_Contract(t *testing.T) {
	r, client := newTestRunner(t)
	restarted, err := NewRunner(r.restConfig, r.client, Options{Config: r.config, Hostname: r.hostname, PullPolicy: r.pullPolicy, ActiveDeadline: r.activeDeadline})
	require.NoError(t, err)
	runnertest.Run(t, runnertest.Target{
		Runner:    r,
		Restarted: restarted,
		Leftovers: func(t *testing.T, jobID string) []string {
			return liveObjects(t, client, jobLabel+"="+jobID)
		},
	})
}

func TestRunner_SecretIsNotWrittenInPodSpec(t *testing.T) {
	r, client := newTestRunner(t)
	const secret = "pod-spec-secret-0b9e"
	job := builder.IsolatedJob{
		ID:    domain.NewID(),
		Image: runnertest.Image,
		RunAs: builder.UIDGID{UID: 1000, GID: 1000},
		Input: runnertest.TarOf(t, map[string]string{}),
		Steps: []builder.IsolatedStep{
			{Name: "with-secret", Cmd: []string{"sh", "-c", "sleep 5"}, SecretEnv: map[string]string{"TOKEN": secret}},
		},
	}

	var checked atomic.Bool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for ctx.Err() == nil {
			pods, err := client.CoreV1().Pods(buildNamespace).List(ctx, metav1.ListOptions{LabelSelector: jobLabel + "=" + job.ID})
			if err == nil && len(pods.Items) > 0 {
				b, _ := json.Marshal(pods.Items)
				if bytes.Contains(b, []byte(secret)) {
					panic("secret found in pod spec")
				}
				checked.Store(true)
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	var log bytes.Buffer
	require.NoError(t, r.Run(context.Background(), job, &log), log.String())
	require.True(t, checked.Load())
}

func TestRunner_ImagePullFailureFailsFast(t *testing.T) {
	for _, policy := range []isolatedrunner.PullPolicy{isolatedrunner.PullIfNotPresent, isolatedrunner.PullNever} {
		t.Run(string(policy), func(t *testing.T) {
			r, client := newTestRunner(t)
			r.pullPolicy = policy
			job := builder.IsolatedJob{
				ID:    domain.NewID(),
				Image: "registry.invalid/no-such-image:1",
				Input: runnertest.TarOf(t, map[string]string{}),
				Steps: []builder.IsolatedStep{{Name: "never", Cmd: []string{"true"}}},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			err := r.Run(ctx, job, &bytes.Buffer{})
			require.Error(t, err)
			require.NoError(t, ctx.Err(), "should fail before the timeout")
			require.Empty(t, liveObjects(t, client, jobLabel+"="+job.ID))
		})
	}
}

func TestRunner_ReclaimOrphans(t *testing.T) {
	r, client := newTestRunner(t)
	ctx := context.Background()

	create := func(host string) {
		pod, err := client.CoreV1().Pods(buildNamespace).Create(ctx, &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{GenerateName: "orphan-", Labels: map[string]string{jobLabel: domain.NewID(), hostLabel: host}},
			Spec: v1.PodSpec{
				Containers: []v1.Container{{Name: "c", Image: runnertest.Image, Command: []string{"sleep", "3600"}}},
			},
		}, metav1.CreateOptions{})
		require.NoError(t, err)
		t.Cleanup(func() { _ = client.CoreV1().Pods(buildNamespace).Delete(ctx, pod.Name, metav1.DeleteOptions{}) })
	}
	create(r.hostname)
	create("other-builder")

	require.NoError(t, r.ReclaimOrphans(ctx))
	require.Empty(t, liveObjects(t, client, hostLabel+"="+r.hostname))
	require.NotEmpty(t, liveObjects(t, client, hostLabel+"=other-builder"))
}

// TestBuilderCannotReadSecretsOutsideBuildNamespace checks the reason build pods have their own namespace:
// whoever can create pods in a namespace can read every secret in it through the pod spec.
func TestBuilderCannotReadSecretsOutsideBuildNamespace(t *testing.T) {
	kubeconf, client := prepareClient(t)
	_, bClient := builderClient(t, kubeconf, client)
	ctx := context.Background()

	secret, err := client.CoreV1().Secrets(builderNamespace).Create(ctx, &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{GenerateName: "db-password-"},
		StringData: map[string]string{"password": "do-not-read"},
	}, metav1.CreateOptions{})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.CoreV1().Secrets(builderNamespace).Delete(ctx, secret.Name, metav1.DeleteOptions{}) })

	_, err = bClient.CoreV1().Pods(builderNamespace).Create(ctx, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{GenerateName: "steal-"},
		Spec: v1.PodSpec{
			RestartPolicy: v1.RestartPolicyNever,
			Containers: []v1.Container{{
				Name: "c", Image: runnertest.Image, Command: []string{"sh", "-c", "echo $PASSWORD"},
				Env: []v1.EnvVar{{Name: "PASSWORD", ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: secret.Name}, Key: "password",
				}}}},
			}},
		},
	}, metav1.CreateOptions{})
	require.True(t, errors.IsForbidden(err), "creating the pod should be forbidden, got: %v", err)
}

func TestRunner_ReclaimEnded(t *testing.T) {
	r, client := newTestRunner(t)
	ctx := context.Background()

	host := "gone-builder-" + random.SecureGenerateHex(4)
	// createWithInit creates a pod running initScript in an init container, if not empty, and then script, as build pods run their steps.
	createWithInit := func(initScript, script string) string {
		spec := v1.PodSpec{
			RestartPolicy: v1.RestartPolicyNever,
			Containers:    []v1.Container{{Name: "c", Image: runnertest.Image, Command: []string{"sh", "-c", script}}},
		}
		if initScript != "" {
			spec.InitContainers = []v1.Container{{Name: "init", Image: runnertest.Image, Command: []string{"sh", "-c", initScript}}}
		}
		pod, err := client.CoreV1().Pods(buildNamespace).Create(ctx, &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{GenerateName: "ended-", Labels: map[string]string{jobLabel: domain.NewID(), hostLabel: host}},
			Spec:       spec,
		}, metav1.CreateOptions{})
		require.NoError(t, err)
		t.Cleanup(func() { _ = client.CoreV1().Pods(buildNamespace).Delete(ctx, pod.Name, metav1.DeleteOptions{}) })
		return "pod/" + pod.Name
	}
	create := func(script string) string { return createWithInit("", script) }
	succeeded := create("true")
	failed := create("exit 1")
	running := create("sleep 3600")
	// A pod that ran long, which ended well after it was created.
	ranLong := create("sleep 10")
	// A pod whose step before the last ran long and failed, which leaves only the end of its init container.
	failedLongInInit := createWithInit("sleep 10; exit 1", "true")

	phase := func(name string) v1.PodPhase {
		pod, err := client.CoreV1().Pods(buildNamespace).Get(ctx, strings.TrimPrefix(name, "pod/"), metav1.GetOptions{})
		require.NoError(t, err)
		return pod.Status.Phase
	}
	ended := func(name string) bool { p := phase(name); return p == v1.PodSucceeded || p == v1.PodFailed }
	require.Eventually(t, func() bool { return ended(succeeded) && ended(failed) }, 2*time.Minute, time.Second)

	// Pods that have just ended may still be watched by the builder that created them.
	require.NoError(t, r.ReclaimEnded(ctx))
	require.ElementsMatch(t, []string{succeeded, failed, running, ranLong, failedLongInInit}, liveObjects(t, client, hostLabel+"="+host))

	require.Eventually(t, func() bool { return ended(ranLong) && ended(failedLongInInit) }, 2*time.Minute, 200*time.Millisecond)
	// The grace is shorter than the time since the long pods were created, and longer than the time since they ended.
	r.endedPodGrace = 5 * time.Second
	require.NoError(t, r.ReclaimEnded(ctx))
	require.ElementsMatch(t, []string{running, ranLong, failedLongInInit}, liveObjects(t, client, hostLabel+"="+host))

	r.endedPodGrace = 0
	require.NoError(t, r.ReclaimEnded(ctx))
	require.Equal(t, []string{running}, liveObjects(t, client, hostLabel+"="+host))
}

func TestRunner_StepsGetNoServiceLinks(t *testing.T) {
	r, client := newTestRunner(t)
	ctx := context.Background()
	svc, err := client.CoreV1().Services(buildNamespace).Create(ctx, &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "probe-" + random.SecureGenerateHex(4)},
		Spec:       v1.ServiceSpec{Ports: []v1.ServicePort{{Port: 80}}},
	}, metav1.CreateOptions{})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.CoreV1().Services(buildNamespace).Delete(ctx, svc.Name, metav1.DeleteOptions{}) })

	envName := strings.ToUpper(strings.ReplaceAll(svc.Name, "-", "_")) + "_SERVICE_HOST"
	job := builder.IsolatedJob{
		ID:    domain.NewID(),
		Image: runnertest.Image,
		Input: runnertest.TarOf(t, map[string]string{}),
		Steps: []builder.IsolatedStep{{Name: "look", Cmd: []string{"sh", "-c", `test -z "$(printenv ` + envName + `)" && echo no-service-links`}}},
	}
	var log bytes.Buffer
	require.NoError(t, r.Run(ctx, job, &log), log.String())
	require.Contains(t, log.String(), "no-service-links")
}

func TestRunner_ResourcesApplyToBuildPods(t *testing.T) {
	r, client := newTestRunner(t)
	// No node has this much CPU, so the pod is never scheduled.
	r.config.Resources.Requests.CPU = "1000"
	r.config.Resources.Limits.Memory = "64Mi"
	job := builder.IsolatedJob{
		ID:    domain.NewID(),
		Image: runnertest.Image,
		Input: runnertest.TarOf(t, map[string]string{}),
		Steps: []builder.IsolatedStep{{Name: "never", Cmd: []string{"true"}}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := r.Run(ctx, job, &bytes.Buffer{})
	require.ErrorContains(t, err, "unschedulable")
	require.ErrorContains(t, err, "nodes are available", "should tell why the scheduler cannot place the pod")
	require.NoError(t, ctx.Err(), "should fail as soon as the pod is unschedulable")
	require.Empty(t, liveObjects(t, client, jobLabel+"="+job.ID))
}

func TestRunner_PodSettingsApplyToEveryContainer(t *testing.T) {
	r, client := newTestRunner(t)
	r.config.Resources.Requests.CPU = "10m"
	r.config.Resources.Requests.Memory = "16Mi"
	r.config.Resources.Limits.CPU = "500m"
	r.config.Resources.Limits.Memory = "128Mi"
	// The node of the test cluster has this label.
	r.config.NodeSelector = []NodeLabel{{Key: "kubernetes.io/os", Value: "linux"}}
	r.config.Tolerations = []Toleration{
		{Key: "ns.trap.jp/builder", Operator: "Exists", Effect: "NoSchedule"},
		{Key: "ns.trap.jp/other", Operator: "Equal", Value: "v", Effect: "NoExecute", TolerationSeconds: ptr.To[int64](30)},
	}
	job := builder.IsolatedJob{
		ID:    domain.NewID(),
		Image: runnertest.Image,
		Input: runnertest.TarOf(t, map[string]string{}),
		Steps: []builder.IsolatedStep{
			{Name: "first", Cmd: []string{"true"}},
			{Name: "slow", Cmd: []string{"sh", "-c", "echo slow-started; sleep 600"}},
		},
	}
	var log runnertest.SyncBuffer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- r.Run(ctx, job, &log) }()
	require.Eventually(t, func() bool { return strings.Contains(log.String(), "slow-started") }, 3*time.Minute, 100*time.Millisecond)

	pods, err := client.CoreV1().Pods(buildNamespace).List(context.Background(), metav1.ListOptions{LabelSelector: jobLabel + "=" + job.ID})
	require.NoError(t, err)
	require.Len(t, pods.Items, 1)
	pod := pods.Items[0]
	cancel()
	<-done

	want := v1.ResourceRequirements{
		Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("10m"), v1.ResourceMemory: resource.MustParse("16Mi")},
		Limits:   v1.ResourceList{v1.ResourceCPU: resource.MustParse("500m"), v1.ResourceMemory: resource.MustParse("128Mi")},
	}
	containers := append(pod.Spec.InitContainers, pod.Spec.Containers...)
	require.Len(t, containers, 3, "receive and the two steps")
	for _, c := range containers {
		require.Truef(t, equality.Semantic.DeepEqual(want, c.Resources), "resources of %s: %v", c.Name, c.Resources)
	}
	require.Equal(t, map[string]string{"kubernetes.io/os": "linux"}, pod.Spec.NodeSelector)
	require.Contains(t, pod.Spec.Tolerations, v1.Toleration{Key: "ns.trap.jp/builder", Operator: v1.TolerationOpExists, Effect: v1.TaintEffectNoSchedule})
	require.Contains(t, pod.Spec.Tolerations, v1.Toleration{
		Key: "ns.trap.jp/other", Operator: v1.TolerationOpEqual, Value: "v", Effect: v1.TaintEffectNoExecute, TolerationSeconds: ptr.To[int64](30),
	})
}

func TestRunner_ActiveDeadlineEndsTheBuildPod(t *testing.T) {
	r, client := newTestRunner(t)
	r.activeDeadline = 10 * time.Second
	job := builder.IsolatedJob{
		ID:    domain.NewID(),
		Image: runnertest.Image,
		Input: runnertest.TarOf(t, map[string]string{}),
		Steps: []builder.IsolatedStep{{Name: "slow", Cmd: []string{"sleep", "600"}}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	err := r.Run(ctx, job, &bytes.Buffer{})
	require.Error(t, err)
	require.NoError(t, ctx.Err(), "should fail by the active deadline, before the timeout")
	require.Empty(t, liveObjects(t, client, jobLabel+"="+job.ID))
}

func TestRunner_DeletingTheBuildPodFailsTheJob(t *testing.T) {
	r, client := newTestRunner(t)
	job := builder.IsolatedJob{
		ID:    domain.NewID(),
		Image: runnertest.Image,
		Input: runnertest.TarOf(t, map[string]string{}),
		Steps: []builder.IsolatedStep{
			{Name: "slow", Cmd: []string{"sh", "-c", "echo slow-started; sleep 600"}},
			{Name: "never", Cmd: []string{"echo", "never-ran"}},
		},
	}
	var log runnertest.SyncBuffer
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	go func() {
		for ctx.Err() == nil && !strings.Contains(log.String(), "slow-started") {
			time.Sleep(100 * time.Millisecond)
		}
		// Deleting without the grace period removes the pod before the runner waits for it again,
		// so the runner sees no event of the deletion.
		_ = client.CoreV1().Pods(buildNamespace).DeleteCollection(context.Background(), metav1.DeleteOptions{GracePeriodSeconds: new(int64)},
			metav1.ListOptions{LabelSelector: jobLabel + "=" + job.ID})
	}()

	err := r.Run(ctx, job, &log)
	require.Error(t, err)
	require.NoError(t, ctx.Err(), "should fail as soon as the pod is deleted")
	require.NotContains(t, log.String(), "never-ran")
}

// TestRunner_SecretGoesWithItsPod checks that the secret of a build is deleted with its pod
// even when the builder does not delete the secret, as when ReclaimOrphans deletes only the pods of a crashed builder.
func TestRunner_SecretGoesWithItsPod(t *testing.T) {
	kubeconf, client := prepareClient(t)
	var rules []rbacv1.PolicyRule
	for _, rule := range builderRole {
		if slices.Contains(rule.Resources, "secrets") {
			rule.Verbs = []string{"create"}
		}
		rules = append(rules, rule)
	}
	conf, c := serviceAccountClient(t, kubeconf, client, "ns-builder-no-secret-delete", rules)
	r, err := NewRunner(conf, c, Options{Config: Config{Namespace: buildNamespace}, Hostname: "fake-builder-" + random.SecureGenerateHex(4), PullPolicy: isolatedrunner.PullIfNotPresent, ActiveDeadline: 10 * time.Minute})
	require.NoError(t, err)

	job := builder.IsolatedJob{
		ID:    domain.NewID(),
		Image: runnertest.Image,
		Input: runnertest.TarOf(t, map[string]string{}),
		Steps: []builder.IsolatedStep{{Name: "with-secret", Cmd: []string{"true"}, SecretEnv: map[string]string{"TOKEN": "t"}}},
	}
	require.NoError(t, r.Run(context.Background(), job, &bytes.Buffer{}))
	require.Eventually(t, func() bool {
		secrets, err := client.CoreV1().Secrets(buildNamespace).List(context.Background(), metav1.ListOptions{LabelSelector: jobLabel + "=" + job.ID})
		require.NoError(t, err)
		return len(secrets.Items) == 0
	}, 2*time.Minute, time.Second)
}
