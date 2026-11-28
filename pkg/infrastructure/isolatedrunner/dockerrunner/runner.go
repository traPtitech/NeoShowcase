// Package dockerrunner implements builder.IsolatedRunner with a Docker daemon.
// Each step runs in its own container, and the shared directories are named volumes created for each job.
package dockerrunner

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"time"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
	"github.com/samber/oops"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/isolatedrunner"
	"github.com/traPtitech/neoshowcase/pkg/util/random"
)

const (
	// jobLabel records the ID of the job the resource was created for.
	jobLabel = "ns.trap.jp/isolated-job"
	// hostLabel records the hostname of the builder that created the resource.
	hostLabel = "ns.trap.jp/builder-host"

	cleanupTimeout = 30 * time.Second
)

type Config struct {
	// Network is the Docker network that step containers are connected to. Empty means the default bridge network.
	Network string `mapstructure:"network" yaml:"network"`
}

var _ builder.IsolatedRunner = (*Runner)(nil)

type Runner struct {
	c          *client.Client
	config     Config
	pullPolicy isolatedrunner.PullPolicy
	hostname   string
}

// Options are the settings of a Runner other than its client.
type Options struct {
	Config Config
	// Hostname labels the resources of the runner and must be unique among the builders sharing the Docker daemon.
	Hostname   string
	PullPolicy isolatedrunner.PullPolicy
}

// NewRunner fails if the pull policy is unknown.
func NewRunner(c *client.Client, opts Options) (*Runner, error) {
	if err := opts.PullPolicy.Validate(); err != nil {
		return nil, err
	}
	return &Runner{
		c:          c,
		config:     opts.Config,
		pullPolicy: opts.PullPolicy,
		hostname:   opts.Hostname,
	}, nil
}

// containerSpec is what differs between the containers of a job.
type containerSpec struct {
	stepName string
	user     string
	cmd      []string
	workDir  string
	env      []string
	stdin    io.Reader
}

// jobRun holds the resources created for one job.
type jobRun struct {
	r          *Runner
	job        builder.IsolatedJob
	namePrefix string
	labels     map[string]string
	mounts     []mount.Mount
	volumes    []string
	containers []string
}

func (r *Runner) Run(ctx context.Context, job builder.IsolatedJob, logWriter io.Writer) error {
	if err := isolatedrunner.ValidateJob(job); err != nil {
		return err
	}

	run := &jobRun{
		r:          r,
		job:        job,
		namePrefix: "ns-build-" + job.ID + "-" + random.SecureGenerateHex(8),
		labels:     map[string]string{jobLabel: job.ID, hostLabel: r.hostname},
	}
	defer run.cleanup(context.WithoutCancel(ctx))

	if err := r.ensureImage(ctx, job.Image); err != nil {
		return err
	}
	if err := run.createVolumes(ctx); err != nil {
		return err
	}

	err := run.runContainer(ctx, containerSpec{
		stepName: isolatedrunner.ReceiveStepName,
		user:     "0:0",
		cmd:      []string{"sh", "-c", isolatedrunner.ReceiveScript(job)},
		stdin:    job.Input,
	}, logWriter)
	if err != nil {
		return err
	}

	user := fmt.Sprintf("%d:%d", job.RunAs.UID, job.RunAs.GID)
	for _, step := range job.Steps {
		isolatedrunner.WriteStepStart(logWriter, step.Name)
		err := run.runContainer(ctx, containerSpec{
			stepName: step.Name,
			user:     user,
			cmd:      step.Cmd,
			workDir:  step.WorkDir,
			env:      envList(step.Env, step.SecretEnv),
		}, logWriter)
		if err != nil {
			return err
		}
	}
	return nil
}

func envList(envs ...map[string]string) []string {
	var list []string
	for _, env := range envs {
		for k, v := range env {
			list = append(list, k+"="+v)
		}
	}
	return list
}

func (r *Runner) ensureImage(ctx context.Context, image string) error {
	switch r.pullPolicy {
	case isolatedrunner.PullNever:
		return nil
	case isolatedrunner.PullIfNotPresent:
		if _, err := r.c.ImageInspect(ctx, image); err == nil {
			return nil
		}
	}
	res, err := r.c.ImagePull(ctx, image, client.ImagePullOptions{})
	if err != nil {
		return oops.Wrapf(err, "pulling image %s", image)
	}
	defer res.Close()
	if err := res.Wait(ctx); err != nil {
		return oops.Wrapf(err, "pulling image %s", image)
	}
	return nil
}

func (run *jobRun) createVolumes(ctx context.Context) error {
	for i, dir := range run.job.Shared {
		name := run.namePrefix + "-" + strconv.Itoa(i)
		_, err := run.r.c.VolumeCreate(ctx, client.VolumeCreateOptions{Name: name, Labels: run.labels})
		if err != nil {
			return oops.Wrapf(err, "creating volume for %s", dir)
		}
		run.volumes = append(run.volumes, name)
		run.mounts = append(run.mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: name,
			Target: dir,
			// Start from an empty directory as an emptyDir of k8s does, instead of copying the image contents.
			VolumeOptions: &mount.VolumeOptions{NoCopy: true},
		})
	}
	return nil
}

// runContainer runs a container to completion, writing its output to logWriter, and removes it.
func (run *jobRun) runContainer(ctx context.Context, spec containerSpec, logWriter io.Writer) error {
	created, err := run.r.c.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name: run.namePrefix + "-" + spec.stepName,
		Config: &container.Config{
			Image:        run.job.Image,
			User:         spec.user,
			Entrypoint:   spec.cmd,
			WorkingDir:   spec.workDir,
			Env:          spec.env,
			Labels:       run.labels,
			AttachStdout: true,
			AttachStderr: true,
			AttachStdin:  spec.stdin != nil,
			OpenStdin:    spec.stdin != nil,
			StdinOnce:    spec.stdin != nil,
		},
		HostConfig: &container.HostConfig{
			Mounts:      run.mounts,
			NetworkMode: container.NetworkMode(run.r.config.Network),
		},
	})
	if err != nil {
		return oops.Wrapf(err, "creating container for step %s", spec.stepName)
	}
	run.containers = append(run.containers, created.ID)
	defer run.removeContainer(context.WithoutCancel(ctx), created.ID)

	attached, err := run.r.c.ContainerAttach(ctx, created.ID, client.ContainerAttachOptions{
		Stream: true,
		Stdin:  spec.stdin != nil,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return oops.Wrapf(err, "attaching to container for step %s", spec.stepName)
	}
	defer attached.Close()
	stopClose := context.AfterFunc(ctx, attached.Close)
	defer stopClose()

	if _, err := run.r.c.ContainerStart(ctx, created.ID, client.ContainerStartOptions{}); err != nil {
		return oops.Wrapf(err, "starting container for step %s", spec.stepName)
	}

	stdinErr := make(chan error, 1)
	if spec.stdin != nil {
		go func() {
			_, err := io.Copy(attached.Conn, spec.stdin)
			if err == nil {
				err = attached.CloseWrite()
			} else {
				// The container would wait for the rest of the input forever, so stop reading its output.
				attached.Close()
			}
			stdinErr <- err
		}()
	} else {
		stdinErr <- nil
	}

	_, outErr := stdcopy.StdCopy(logWriter, logWriter, attached.Reader)
	if ctx.Err() != nil {
		return oops.Wrapf(ctx.Err(), "step %s", spec.stepName)
	}
	if err := <-stdinErr; err != nil {
		return oops.Wrapf(err, "writing input of step %s", spec.stepName)
	}
	if outErr != nil {
		return oops.Wrapf(outErr, "reading output of step %s", spec.stepName)
	}

	wait := run.r.c.ContainerWait(ctx, created.ID, client.ContainerWaitOptions{Condition: container.WaitConditionNotRunning})
	select {
	case res := <-wait.Result:
		if res.Error != nil {
			return oops.Errorf("waiting for step %s: %s", spec.stepName, res.Error.Message)
		}
		if res.StatusCode != 0 {
			return oops.Errorf("step %s exited with code %d", spec.stepName, res.StatusCode)
		}
		return nil
	case err := <-wait.Error:
		return oops.Wrapf(err, "waiting for step %s", spec.stepName)
	}
}

func (run *jobRun) removeContainer(ctx context.Context, id string) {
	ctx, cancel := context.WithTimeout(ctx, cleanupTimeout)
	defer cancel()
	_, err := run.r.c.ContainerRemove(ctx, id, client.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
	if err != nil && !cerrdefs.IsNotFound(err) {
		slog.ErrorContext(ctx, "failed to remove build container", "container", id, "error", err)
	}
}

func (run *jobRun) cleanup(ctx context.Context) {
	for _, id := range run.containers {
		run.removeContainer(ctx, id)
	}
	ctx, cancel := context.WithTimeout(ctx, cleanupTimeout)
	defer cancel()
	for _, name := range run.volumes {
		if _, err := run.r.c.VolumeRemove(ctx, name, client.VolumeRemoveOptions{Force: true}); err != nil && !cerrdefs.IsNotFound(err) {
			slog.ErrorContext(ctx, "failed to remove build volume", "volume", name, "error", err)
		}
	}
}

func (r *Runner) ReclaimOrphans(ctx context.Context) error {
	filters := make(client.Filters).Add("label", hostLabel+"="+r.hostname)
	containers, err := r.c.ContainerList(ctx, client.ContainerListOptions{All: true, Filters: filters})
	if err != nil {
		return oops.Wrapf(err, "listing orphaned build containers")
	}
	for _, ctr := range containers.Items {
		if _, err := r.c.ContainerRemove(ctx, ctr.ID, client.ContainerRemoveOptions{Force: true}); err != nil && !cerrdefs.IsNotFound(err) {
			return oops.Wrapf(err, "removing orphaned build container %s", ctr.ID)
		}
	}
	volumes, err := r.c.VolumeList(ctx, client.VolumeListOptions{Filters: filters})
	if err != nil {
		return oops.Wrapf(err, "listing orphaned build volumes")
	}
	for _, v := range volumes.Items {
		if _, err := r.c.VolumeRemove(ctx, v.Name, client.VolumeRemoveOptions{Force: true}); err != nil && !cerrdefs.IsNotFound(err) {
			return oops.Wrapf(err, "removing orphaned build volume %s", v.Name)
		}
	}
	return nil
}

// ReclaimEnded does nothing. Resources of a job whose builder is gone are left until a builder with the same hostname starts,
// since nothing tells whether the builder that created them is still running.
func (r *Runner) ReclaimEnded(context.Context) error {
	return nil
}
