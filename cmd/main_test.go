package main

import (
	"context"
	"errors"
	"os"
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
