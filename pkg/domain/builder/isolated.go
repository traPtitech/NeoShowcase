package builder

import (
	"context"
	"io"
)

// IsolatedRunner runs the steps of a job in a disposable environment created for each job.
//
// Every implementation guarantees the following:
//   - Steps run one at a time in the given order. A step starts only after every process of the previous step has exited.
//   - The SecretEnv of a step appears only in the environment of the processes of that step.
//     It never appears in other steps, in the Shared directories, or in the log.
//   - Only the contents of the Shared directories are passed between steps. Each step runs in a new container.
//   - Step processes run as RunAs. The Shared directories and the contents of Input are owned by RunAs when the first step starts.
//   - Step containers are given no credentials of the orchestrator (ServiceAccount tokens, the Docker socket).
type IsolatedRunner interface {
	// Run returns an error containing the step name as soon as a step exits with a non-zero code,
	// or an error wrapping ctx.Err() as soon as ctx is done.
	// It also fails with the error of job.Input when reading it fails, instead of running the steps on a partial input.
	// By the time Run returns, deletion of every resource created for the job has been requested.
	Run(ctx context.Context, job IsolatedJob, logWriter io.Writer) error
	// ReclaimOrphans deletes resources left behind by jobs of a previous process of this builder,
	// such as when the builder crashed in the middle of a job.
	// It must be called before the first Run, since it also deletes resources of jobs running at the time.
	ReclaimOrphans(ctx context.Context) error
	// ReclaimEnded deletes resources left behind by jobs that have ended, whichever builder ran them,
	// such as when a builder was removed in the middle of a job.
	// Unlike ReclaimOrphans, it leaves running jobs alone and may be called at any time.
	ReclaimEnded(ctx context.Context) error
}

// IsolatedJob is a unit of work for IsolatedRunner.
type IsolatedJob struct {
	// ID identifies the job in resource names and labels.
	// It must consist of lowercase alphanumerics and be at most 32 characters long.
	ID string
	// Image is used for every step. It must contain sh and tar.
	Image string
	// RunAs is the user that step processes run as.
	RunAs UIDGID
	// Shared lists absolute paths of the directories shared between steps.
	Shared []string
	// Input is a tar stream extracted at / before the first step. It is required; an empty job has an empty tar stream.
	// Every entry must be under one of Shared, with its name relative to / (e.g. "workspace/app/main.go").
	Input io.Reader
	Steps []IsolatedStep
}

type UIDGID struct {
	UID int `mapstructure:"uid" yaml:"uid"`
	GID int `mapstructure:"gid" yaml:"gid"`
}

// IsolatedStep is a command run in its own container.
type IsolatedStep struct {
	// Name identifies the step in the log and in errors.
	// It must be a DNS label (lowercase alphanumerics and '-') other than "receive".
	Name    string
	Cmd     []string
	WorkDir string
	Env     map[string]string
	// SecretEnv is passed only to this step.
	SecretEnv map[string]string
}
