package dockerrunner

import (
	"bytes"
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner/runnertest"
)

func prepareClient(t *testing.T) *client.Client {
	t.Helper()
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_DOCKER_TESTS")); !ok {
		t.SkipNow()
	}
	c, err := client.New(client.FromEnv)
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func leftovers(c *client.Client, filters client.Filters) func(t *testing.T) []string {
	return func(t *testing.T) []string {
		t.Helper()
		var names []string
		containers, err := c.ContainerList(context.Background(), client.ContainerListOptions{All: true, Filters: filters})
		require.NoError(t, err)
		for _, ctr := range containers.Items {
			names = append(names, ctr.Names...)
		}
		volumes, err := c.VolumeList(context.Background(), client.VolumeListOptions{Filters: filters})
		require.NoError(t, err)
		for _, v := range volumes.Items {
			names = append(names, v.Name)
		}
		return names
	}
}

func TestRunner_Contract(t *testing.T) {
	c := prepareClient(t)
	hostname := "contract-test-" + strconv.Itoa(os.Getpid())
	r, err := NewRunner(c, Options{Hostname: hostname, PullPolicy: isolatedrunner.PullIfNotPresent})
	require.NoError(t, err)
	restarted, err := NewRunner(c, Options{Hostname: hostname, PullPolicy: isolatedrunner.PullIfNotPresent})
	require.NoError(t, err)
	runnertest.Run(t, runnertest.Target{
		Runner:    r,
		Restarted: restarted,
		Leftovers: func(t *testing.T, jobID string) []string {
			return leftovers(c, make(client.Filters).Add("label", jobLabel+"="+jobID))(t)
		},
	})
}

func TestRunner_ReclaimOrphans(t *testing.T) {
	c := prepareClient(t)
	ctx := context.Background()
	mine := "reclaim-test-mine-" + strconv.Itoa(os.Getpid())
	other := "reclaim-test-other-" + strconv.Itoa(os.Getpid())

	for _, host := range []string{mine, other} {
		labels := map[string]string{hostLabel: host}
		_, err := c.VolumeCreate(ctx, client.VolumeCreateOptions{Name: host, Labels: labels})
		require.NoError(t, err)
		_, err = c.ContainerCreate(ctx, client.ContainerCreateOptions{
			Name:   host,
			Config: &container.Config{Image: runnertest.Image, Labels: labels},
		})
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		_, _ = c.ContainerRemove(ctx, other, client.ContainerRemoveOptions{Force: true})
		_, _ = c.VolumeRemove(ctx, other, client.VolumeRemoveOptions{Force: true})
	})

	r, err := NewRunner(c, Options{Hostname: mine, PullPolicy: isolatedrunner.PullIfNotPresent})
	require.NoError(t, err)
	require.NoError(t, r.ReclaimOrphans(ctx))

	require.Empty(t, leftovers(c, make(client.Filters).Add("label", hostLabel+"="+mine))(t))
	require.NotEmpty(t, leftovers(c, make(client.Filters).Add("label", hostLabel+"="+other))(t))
}

func TestRunner_PullNeverDoesNotPull(t *testing.T) {
	c := prepareClient(t)
	r, err := NewRunner(c, Options{Hostname: "pull-test-" + strconv.Itoa(os.Getpid()), PullPolicy: isolatedrunner.PullNever})
	require.NoError(t, err)
	job := builder.IsolatedJob{
		ID:    domain.NewID(),
		Image: "registry.invalid/no-such-image:1",
		Input: runnertest.TarOf(t, map[string]string{}),
		Steps: []builder.IsolatedStep{{Name: "never", Cmd: []string{"true"}}},
	}
	err = r.Run(context.Background(), job, &bytes.Buffer{})
	require.ErrorContains(t, err, "No such image")
	require.NotContains(t, err.Error(), "pull")
}

func TestRunner_SharedDirectoriesStartEmpty(t *testing.T) {
	c := prepareClient(t)
	r, err := NewRunner(c, Options{Hostname: "empty-test-" + strconv.Itoa(os.Getpid()), PullPolicy: isolatedrunner.PullIfNotPresent})
	require.NoError(t, err)
	// /var has files in the image, which must not be copied into the volume as an emptyDir of k8s does not.
	job := builder.IsolatedJob{
		ID:     domain.NewID(),
		Image:  runnertest.Image,
		Shared: []string{"/var"},
		Input:  runnertest.TarOf(t, map[string]string{}),
		Steps:  []builder.IsolatedStep{{Name: "look", Cmd: []string{"sh", "-c", `test -z "$(ls -A /var)" && echo var-is-empty`}}},
	}
	var log bytes.Buffer
	require.NoError(t, r.Run(context.Background(), job, &log), log.String())
	require.Contains(t, log.String(), "var-is-empty")
}

func TestRunner_PullIfNotPresentUsesTheLocalImage(t *testing.T) {
	c := prepareClient(t)
	ctx := context.Background()
	res, err := c.ImagePull(ctx, runnertest.Image, client.ImagePullOptions{})
	require.NoError(t, err)
	require.NoError(t, res.Wait(ctx))
	// The image exists only in the daemon, so pulling it fails.
	localOnly := "ns-local-only/" + domain.NewID() + ":1"
	_, err = c.ImageTag(ctx, client.ImageTagOptions{Source: runnertest.Image, Target: localOnly})
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = c.ImageRemove(ctx, localOnly, client.ImageRemoveOptions{}) })

	r, err := NewRunner(c, Options{Hostname: "local-image-test-" + strconv.Itoa(os.Getpid()), PullPolicy: isolatedrunner.PullIfNotPresent})
	require.NoError(t, err)
	job := builder.IsolatedJob{
		ID:    domain.NewID(),
		Image: localOnly,
		Input: runnertest.TarOf(t, map[string]string{}),
		Steps: []builder.IsolatedStep{{Name: "look", Cmd: []string{"echo", "local-image-ran"}}},
	}
	var log bytes.Buffer
	require.NoError(t, r.Run(ctx, job, &log), log.String())
	require.Contains(t, log.String(), "local-image-ran")
}
