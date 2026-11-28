package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type failingComponent struct{ shutdown bool }

func (c *failingComponent) Start(context.Context) error { return errors.New("start-failed-7e1b") }

func (c *failingComponent) Shutdown(context.Context) error {
	c.shutdown = true
	return nil
}

func TestRunComponent_ExitsWhenStartFails(t *testing.T) {
	c := &failingComponent{}
	// A context that is never done stands in for a signal that never comes.
	done := make(chan error, 1)
	go func() { done <- runComponent(context.Background(), "test", c) }()

	select {
	case err := <-done:
		require.ErrorContains(t, err, "start-failed-7e1b")
		require.True(t, c.shutdown)
	case <-time.After(5 * time.Second):
		t.Fatal("should exit without waiting for a signal")
	}
}

// blockingComponent runs until the context given to Start is done, as a server does.
type blockingComponent struct {
	started      chan struct{}
	startCtxDone atomic.Bool
	shutdown     atomic.Bool
	// shutdownBlocks makes Shutdown never return, as when it waits for a request that never completes.
	shutdownBlocks bool
}

func newBlockingComponent() *blockingComponent {
	return &blockingComponent{started: make(chan struct{})}
}

func (c *blockingComponent) Start(ctx context.Context) error {
	close(c.started)
	<-ctx.Done()
	c.startCtxDone.Store(true)
	return nil
}

func (c *blockingComponent) Shutdown(context.Context) error {
	c.shutdown.Store(true)
	if c.shutdownBlocks {
		fmt.Println("shutting-down")
		select {}
	}
	return nil
}

func TestRunComponent_ShutsDownWhenContextIsDone(t *testing.T) {
	c := newBlockingComponent()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- runComponent(ctx, "test", c) }()
	<-c.started

	cancel()
	select {
	case err := <-done:
		require.NoError(t, err)
		require.True(t, c.shutdown.Load())
		require.Eventually(t, c.startCtxDone.Load, time.Second, 10*time.Millisecond)
	case <-time.After(5 * time.Second):
		t.Fatal("should shut down when the context is done")
	}
}

func TestRunComponentUntilSignal_ShutsDownOnSignal(t *testing.T) {
	for _, sig := range []syscall.Signal{
		syscall.SIGTERM, // Kubernetes sends SIGTERM to stop a pod.
		syscall.SIGINT,  // Ctrl-C in a terminal.
	} {
		t.Run(sig.String(), func(t *testing.T) {
			c := newBlockingComponent()
			done := make(chan error, 1)
			go func() { done <- runComponentUntilSignal(context.Background(), "test", c) }()
			<-c.started

			require.NoError(t, syscall.Kill(os.Getpid(), sig))
			select {
			case err := <-done:
				require.NoError(t, err)
				require.True(t, c.shutdown.Load())
			case <-time.After(5 * time.Second):
				t.Fatalf("should shut down on %v", sig)
			}
		})
	}
}

const helperProcessEnv = "NS_TEST_STUCK_SHUTDOWN_PROCESS"

// TestRunComponentUntilSignal_SecondSignalTerminates runs a component whose Shutdown never returns in a child process,
// since the second signal terminates the whole process.
func TestRunComponentUntilSignal_SecondSignalTerminates(t *testing.T) {
	if os.Getenv(helperProcessEnv) != "" {
		c := newBlockingComponent()
		c.shutdownBlocks = true
		go func() {
			<-c.started
			fmt.Println("started")
		}()
		_ = runComponentUntilSignal(context.Background(), "test", c)
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestRunComponentUntilSignal_SecondSignalTerminates$")
	cmd.Env = append(os.Environ(), helperProcessEnv+"=1")
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)
	require.NoError(t, cmd.Start())
	// Killing the child ends the reads of its output, which otherwise wait forever when it neither shuts down nor terminates.
	deadline := time.AfterFunc(10*time.Second, func() { _ = cmd.Process.Kill() })
	t.Cleanup(func() {
		deadline.Stop()
		_ = cmd.Process.Kill()
	})
	lines := bufio.NewScanner(stdout)
	waitLine := func(want string) {
		t.Helper()
		for lines.Scan() {
			if strings.TrimSpace(lines.Text()) == want {
				return
			}
		}
		t.Fatalf("child process exited before printing %q", want)
	}

	waitLine("started")
	require.NoError(t, cmd.Process.Signal(syscall.SIGTERM))
	waitLine("shutting-down")

	exited := make(chan error, 1)
	go func() {
		_, _ = io.Copy(io.Discard, stdout)
		exited <- cmd.Wait()
	}()
	// The signal is sent again until the child exits, as a person presses Ctrl-C again,
	// since the signal handler is released right after the shutdown starts.
	timeout := time.After(5 * time.Second)
	for {
		require.NoError(t, cmd.Process.Signal(syscall.SIGTERM))
		select {
		case err := <-exited:
			var exitErr *exec.ExitError
			require.ErrorAs(t, err, &exitErr)
			status := exitErr.Sys().(syscall.WaitStatus)
			require.True(t, status.Signaled(), "should be terminated by the signal, got: %v", err)
			return
		case <-time.After(100 * time.Millisecond):
		case <-timeout:
			t.Fatal("another signal should terminate the process during a stuck shutdown")
		}
	}
}
