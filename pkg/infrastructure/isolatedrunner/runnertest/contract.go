// Package runnertest checks that an implementation of builder.IsolatedRunner keeps the contract of the interface.
package runnertest

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
)

// Image is used by the contract tests. It contains sh, tar, stat and grep.
// Its default user is not root, so that the tests also check that the input is extracted as root.
const Image = "nginxinc/nginx-unprivileged:1.27-alpine"

// Target is the implementation under test.
type Target struct {
	Runner builder.IsolatedRunner
	// Restarted stands in for the next process of the builder running Runner: another runner given the same hostname.
	Restarted builder.IsolatedRunner
	// Leftovers returns the names of resources that still exist for the job.
	Leftovers func(t *testing.T, jobID string) []string
}

// TarOf returns a tar stream holding the given files, keyed by path relative to /.
func TarOf(t *testing.T, files map[string]string) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for name, content := range files {
		require.NoError(t, tw.WriteHeader(&tar.Header{Name: name, Mode: 0o644, Size: int64(len(content)), Typeflag: tar.TypeReg}))
		_, err := tw.Write([]byte(content))
		require.NoError(t, err)
	}
	require.NoError(t, tw.Close())
	return &buf
}

func sh(script string) []string {
	return []string{"sh", "-c", script}
}

func newJob(t *testing.T, steps ...builder.IsolatedStep) builder.IsolatedJob {
	return builder.IsolatedJob{
		ID:     domain.NewID(),
		Image:  Image,
		RunAs:  builder.UIDGID{UID: 1234, GID: 2345},
		Shared: []string{"/data", "/more"},
		Input:  TarOf(t, map[string]string{"data/dir/in.txt": "hello-input"}),
		Steps:  steps,
	}
}

// Run runs the contract tests against the target.
func Run(t *testing.T, target Target) {
	t.Run("steps run in order as RunAs and share only the Shared directories", func(t *testing.T) {
		job := newJob(t,
			builder.IsolatedStep{Name: "first", Cmd: sh(`
				set -e
				test "$(id -u):$(id -g)" = 1234:2345
				test "$(stat -c %u:%g /data /data/dir /data/dir/in.txt /more | sort -u)" = 1234:2345
				test "$(pwd)" = /data/dir
				test "$GREETING" = hi
				cat in.txt > /data/from-first
				echo outside > /tmp/outside
				echo first-done`),
				WorkDir: "/data/dir",
				Env:     map[string]string{"GREETING": "hi"},
			},
			builder.IsolatedStep{Name: "second", Cmd: sh(`
				set -e
				test ! -e /tmp/outside
				test -z "$GREETING"
				echo "second-read-$(cat /data/from-first)"`)},
		)
		var log bytes.Buffer
		err := target.Runner.Run(context.Background(), job, &log)
		require.NoError(t, err, log.String())

		out := log.String()
		require.Contains(t, out, "first-done")
		require.Contains(t, out, "second-read-hello-input")
		require.Less(t, strings.Index(out, "first-done"), strings.Index(out, "second-read"))
		require.Regexp(t, `(?s)\[ns-builder\] Running step "first"\n.*first-done.*\[ns-builder\] Running step "second"\n.*second-read`, out)
		require.Empty(t, target.Leftovers(t, job.ID))
	})

	t.Run("secret env is visible only to its step", func(t *testing.T) {
		const secret = "contract-secret-6f1d2c"
		const search = `
			if env | grep -q contract-secret || grep -rqs contract-secret /data /more /proc/*/environ; then
				echo "found-in-$0"
			else
				echo "not-found-in-$0"
			fi`
		job := newJob(t,
			builder.IsolatedStep{Name: "before", Cmd: []string{"sh", "-c", search, "before"}},
			builder.IsolatedStep{Name: "with-secret", Cmd: sh(`test "$TOKEN" = ` + secret + ` && echo token-ok`),
				SecretEnv: map[string]string{"TOKEN": secret}},
			builder.IsolatedStep{Name: "after", Cmd: []string{"sh", "-c", search, "after"}},
		)
		var log bytes.Buffer
		err := target.Runner.Run(context.Background(), job, &log)
		require.NoError(t, err, log.String())

		out := log.String()
		require.Contains(t, out, "token-ok")
		require.Contains(t, out, "not-found-in-before")
		require.Contains(t, out, "not-found-in-after")
		require.NotContains(t, out, secret)
		require.Empty(t, target.Leftovers(t, job.ID))
	})

	t.Run("a failing step stops the job and is named in the error", func(t *testing.T) {
		job := newJob(t,
			builder.IsolatedStep{Name: "broken", Cmd: sh(`exit 3`)},
			builder.IsolatedStep{Name: "never", Cmd: sh(`echo never-ran`)},
		)
		var log bytes.Buffer
		err := target.Runner.Run(context.Background(), job, &log)
		require.ErrorContains(t, err, "broken")
		require.NotContains(t, log.String(), "never-ran")
		require.Empty(t, target.Leftovers(t, job.ID))
	})

	t.Run("an input that fails to be read fails the job with its error", func(t *testing.T) {
		// A tar stream cut off at an entry boundary, which tar extracts without an error.
		var partial bytes.Buffer
		tw := tar.NewWriter(&partial)
		require.NoError(t, tw.WriteHeader(&tar.Header{Name: "data/in.txt", Mode: 0o644, Size: 2, Typeflag: tar.TypeReg}))
		_, err := tw.Write([]byte("hi"))
		require.NoError(t, err)
		require.NoError(t, tw.Flush())

		job := newJob(t, builder.IsolatedStep{Name: "after-input", Cmd: sh(`echo after-input-ran`)})
		job.Input = io.MultiReader(&partial, failingReader{errors.New("input-broken-3c9a")})
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		var log bytes.Buffer
		err = target.Runner.Run(ctx, job, &log)
		require.ErrorContains(t, err, "input-broken-3c9a")
		require.NoError(t, ctx.Err(), "should fail before the timeout")
		require.NotContains(t, log.String(), "after-input-ran")
		require.Empty(t, target.Leftovers(t, job.ID))
	})

	t.Run("an invalid job fails without creating resources", func(t *testing.T) {
		job := newJob(t, builder.IsolatedStep{Name: "never", Cmd: sh(`echo never-ran`)})
		job.Input = nil
		var log bytes.Buffer
		err := target.Runner.Run(context.Background(), job, &log)
		require.ErrorContains(t, err, "input")
		require.NotContains(t, log.String(), "never-ran")
		require.Empty(t, target.Leftovers(t, job.ID))
	})

	t.Run("steps are given no credentials of the orchestrator", func(t *testing.T) {
		job := newJob(t, builder.IsolatedStep{Name: "look", Cmd: sh(`
			set -e
			test ! -e /var/run/secrets/kubernetes.io/serviceaccount/token
			test ! -S /var/run/docker.sock
			echo no-credentials`)})
		var log bytes.Buffer
		err := target.Runner.Run(context.Background(), job, &log)
		require.NoError(t, err, log.String())
		require.Contains(t, log.String(), "no-credentials")
	})

	t.Run("the deadline of ctx fails the job with context.DeadlineExceeded", func(t *testing.T) {
		job := newJob(t, builder.IsolatedStep{Name: "slow", Cmd: sh(`sleep 600`)})
		// Short enough to expire before the step starts on k8s, where the pod takes a while to be scheduled.
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()
		err := target.Runner.Run(ctx, job, &bytes.Buffer{})
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Empty(t, target.Leftovers(t, job.ID))
	})

	t.Run("cancelling ctx stops the job", func(t *testing.T) {
		job := newJob(t,
			builder.IsolatedStep{Name: "slow", Cmd: sh(`echo slow-started; sleep 600`)},
		)
		var log SyncBuffer
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			waitForLog(&log, "slow-started")
			cancel()
		}()

		start := time.Now()
		err := target.Runner.Run(ctx, job, &log)
		require.ErrorIs(t, err, context.Canceled)
		require.Less(t, time.Since(start), 5*time.Minute)
		require.Empty(t, target.Leftovers(t, job.ID))
	})

	t.Run("the deadline of ctx expiring in a step fails the job with context.DeadlineExceeded", func(t *testing.T) {
		job := newJob(t,
			builder.IsolatedStep{Name: "slow", Cmd: sh(`echo slow-started; sleep 600`)},
		)
		var log SyncBuffer
		ctx := newExpirableContext()
		go func() {
			waitForLog(&log, "slow-started")
			ctx.expire()
		}()

		start := time.Now()
		err := target.Runner.Run(ctx, job, &log)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Less(t, time.Since(start), 5*time.Minute)
		require.Empty(t, target.Leftovers(t, job.ID))
	})

	t.Run("the next process of the builder reclaims the resources of a running job", func(t *testing.T) {
		job := newJob(t,
			builder.IsolatedStep{Name: "slow", Cmd: sh(`echo slow-started; sleep 600`)},
			builder.IsolatedStep{Name: "never", Cmd: sh(`echo never-ran`)},
		)
		var log SyncBuffer
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		reclaimErr := make(chan error, 1)
		go func() {
			waitForLog(&log, "slow-started")
			reclaimErr <- target.Restarted.ReclaimOrphans(context.Background())
		}()

		err := target.Runner.Run(ctx, job, &log)
		require.Error(t, err)
		require.NoError(t, ctx.Err(), "should fail as soon as the resources are reclaimed")
		require.NoError(t, <-reclaimErr)
		require.NotContains(t, log.String(), "never-ran")
		require.Empty(t, target.Leftovers(t, job.ID))
	})
}

func waitForLog(log *SyncBuffer, s string) {
	for !strings.Contains(log.String(), s) {
		time.Sleep(100 * time.Millisecond)
	}
}

// expirableContext is a context whose deadline is reached when expire is called,
// so that a test can make the deadline pass at a chosen point of a job.
type expirableContext struct {
	context.Context
	done chan struct{}
	once sync.Once
}

func newExpirableContext() *expirableContext {
	return &expirableContext{Context: context.Background(), done: make(chan struct{})}
}

func (c *expirableContext) expire() { c.once.Do(func() { close(c.done) }) }

func (c *expirableContext) Done() <-chan struct{} { return c.done }

func (c *expirableContext) Err() error {
	select {
	case <-c.done:
		return context.DeadlineExceeded
	default:
		return nil
	}
}

type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }
