package builder

import (
	"context"
	"errors"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/test/mocks"
)

type reclaimingBuildpack struct {
	reclaimErr   error
	reclaimEnded atomic.Int32
}

func (r *reclaimingBuildpack) Pack(context.Context, string, string, string, builder.ImageConfig, map[string]string, io.Writer) (string, error) {
	return "", errors.New("not expected to pack")
}

func (r *reclaimingBuildpack) ReclaimOrphans(context.Context) error { return r.reclaimErr }

func (r *reclaimingBuildpack) ReclaimEnded(context.Context) error {
	r.reclaimEnded.Add(1)
	return nil
}

// blockingConnect returns a ConnectBuilder that reports the connection and holds it until ctx is done.
func blockingConnect(connected *atomic.Bool) func(context.Context, func(*pb.BuilderRequest), <-chan *pb.BuilderResponse) error {
	return func(ctx context.Context, _ func(*pb.BuilderRequest), _ <-chan *pb.BuilderResponse) error {
		connected.Store(true)
		<-ctx.Done()
		return ctx.Err()
	}
}

func TestService_StartTakesBuildsOnlyAfterReclaimingOrphans(t *testing.T) {
	t.Run("reclaim fails", func(t *testing.T) {
		var connected atomic.Bool
		client := &mocks.ControllerBuilderServiceClientMock{ConnectBuilderFunc: blockingConnect(&connected)}
		s := newService(nil, client, nil, &reclaimingBuildpack{reclaimErr: errors.New("k8s api is down")}, nil, builder.ImageConfig{})

		err := s.Start(context.Background())
		require.ErrorContains(t, err, "k8s api is down")
		time.Sleep(100 * time.Millisecond)
		require.False(t, connected.Load(), "should not take builds, which may be deleted by a later reclaim")
		// The process shuts the service down when it fails to start.
		require.NoError(t, s.Shutdown(context.Background()))
	})

	t.Run("reclaim succeeds", func(t *testing.T) {
		var connected atomic.Bool
		client := &mocks.ControllerBuilderServiceClientMock{ConnectBuilderFunc: blockingConnect(&connected)}
		s := newService(nil, client, nil, &reclaimingBuildpack{}, nil, builder.ImageConfig{})
		require.NoError(t, s.Start(context.Background()))
		t.Cleanup(func() { _ = s.Shutdown(context.Background()) })
		require.Eventually(t, connected.Load, 5*time.Second, 10*time.Millisecond)
	})
}

// slowReclaimBuildpack holds ReclaimOrphans until released, ignoring the cancellation of its context.
type slowReclaimBuildpack struct {
	reclaimingBuildpack
	entered, release chan struct{}
}

func (r *slowReclaimBuildpack) ReclaimOrphans(context.Context) error {
	close(r.entered)
	<-r.release
	return nil
}

func TestService_ShutdownDuringStartStopsTakingBuilds(t *testing.T) {
	var connected atomic.Bool
	client := &mocks.ControllerBuilderServiceClientMock{ConnectBuilderFunc: blockingConnect(&connected)}
	bp := &slowReclaimBuildpack{entered: make(chan struct{}), release: make(chan struct{})}
	s := newService(nil, client, nil, bp, nil, builder.ImageConfig{})
	done := make(chan error, 1)
	go func() { done <- s.Start(context.Background()) }()
	<-bp.entered

	require.NoError(t, s.Shutdown(context.Background()))
	close(bp.release)
	require.NoError(t, <-done)
	time.Sleep(200 * time.Millisecond)
	require.False(t, connected.Load(), "should not take builds after Shutdown returned")
}

// blockedReclaimBuildpack holds ReclaimOrphans until its context is done.
type blockedReclaimBuildpack struct {
	reclaimingBuildpack
	entered chan struct{}
}

func (r *blockedReclaimBuildpack) ReclaimOrphans(ctx context.Context) error {
	close(r.entered)
	<-ctx.Done()
	return ctx.Err()
}

func TestService_ShutdownInterruptsReclaimInStart(t *testing.T) {
	var connected atomic.Bool
	client := &mocks.ControllerBuilderServiceClientMock{ConnectBuilderFunc: blockingConnect(&connected)}
	bp := &blockedReclaimBuildpack{entered: make(chan struct{})}
	s := newService(nil, client, nil, bp, nil, builder.ImageConfig{})
	done := make(chan struct{})
	go func() {
		_ = s.Start(context.Background()) // Fails or not, depending on how the reclaim reports the interruption.
		close(done)
	}()
	<-bp.entered

	require.NoError(t, s.Shutdown(context.Background()))
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Start should return once Shutdown interrupts the reclaim")
	}
	time.Sleep(200 * time.Millisecond)
	require.False(t, connected.Load(), "should not take builds after Shutdown returned")
}

func TestService_ShutdownBeforeStartStopsTakingBuilds(t *testing.T) {
	var connected atomic.Bool
	client := &mocks.ControllerBuilderServiceClientMock{ConnectBuilderFunc: blockingConnect(&connected)}
	s := newService(nil, client, nil, &reclaimingBuildpack{}, nil, builder.ImageConfig{})

	// The process may shut the service down before the goroutine running Start is scheduled.
	require.NoError(t, s.Shutdown(context.Background()))
	require.NoError(t, s.Start(context.Background()))
	time.Sleep(200 * time.Millisecond)
	require.False(t, connected.Load(), "should not take builds after Shutdown returned")
}

func TestService_StartReclaimsEndedBuilds(t *testing.T) {
	client := &mocks.ControllerBuilderServiceClientMock{
		ConnectBuilderFunc: func(ctx context.Context, _ func(*pb.BuilderRequest), _ <-chan *pb.BuilderResponse) error {
			<-ctx.Done()
			return ctx.Err()
		},
	}
	bp := &reclaimingBuildpack{}
	s := newService(nil, client, nil, bp, nil, builder.ImageConfig{})
	require.NoError(t, s.Start(context.Background()))
	t.Cleanup(func() { _ = s.Shutdown(context.Background()) })
	require.Eventually(t, func() bool { return bp.reclaimEnded.Load() > 0 }, 5*time.Second, 10*time.Millisecond)
}

func TestService_ShutdownWaitsForRunningBuildToCleanUp(t *testing.T) {
	var cleanedUp atomic.Bool
	started := make(chan struct{})
	client := &mocks.ControllerBuilderServiceClientMock{
		ConnectBuilderFunc: func(ctx context.Context, _ func(*pb.BuilderRequest), _ <-chan *pb.BuilderResponse) error {
			<-ctx.Done()
			return ctx.Err()
		},
		StreamBuildLogFunc: func(_ context.Context, _ string, send <-chan []byte) error {
			for range send {
			}
			return nil
		},
		PingBuildFunc:    func(context.Context, string) error { return nil },
		SaveBuildLogFunc: func(context.Context, string, []byte) error { return nil },
	}
	// Stands in for a build step whose resources take a while to be deleted after the build is cancelled.
	gitsvc := &mocks.GitServiceMock{
		CloneRepositoryFunc: func(ctx context.Context, _ string, _ *domain.Repository, _ string) error {
			close(started)
			<-ctx.Done()
			time.Sleep(500 * time.Millisecond)
			cleanedUp.Store(true)
			return ctx.Err()
		},
	}
	s := newService(&Config{StepTimeout: time.Hour}, client, nil, &reclaimingBuildpack{}, gitsvc, builder.ImageConfig{})
	require.NoError(t, s.Start(context.Background()))

	app := &domain.Application{ID: domain.NewID(), Config: domain.ApplicationConfig{BuildConfig: &domain.BuildConfigRuntimeBuildpack{}}}
	require.NoError(t, s.startBuild(&domain.StartBuildRequest{
		App:   app,
		Build: &domain.Build{ID: domain.NewID(), ApplicationID: app.ID},
		Repo:  &domain.Repository{},
	}))
	<-started

	require.NoError(t, s.Shutdown(context.Background()))
	require.True(t, cleanedUp.Load(), "Shutdown should return only after the running build has cleaned up")
}

// A request received by ConnectBuilder just before Shutdown may reach the service after Shutdown returned.
func TestService_BuildRequestAfterShutdownIsNotTaken(t *testing.T) {
	var cloned atomic.Bool
	client := &mocks.ControllerBuilderServiceClientMock{
		ConnectBuilderFunc: func(ctx context.Context, _ func(*pb.BuilderRequest), _ <-chan *pb.BuilderResponse) error {
			<-ctx.Done()
			return ctx.Err()
		},
		StreamBuildLogFunc: func(_ context.Context, _ string, send <-chan []byte) error {
			for range send {
			}
			return nil
		},
		PingBuildFunc:    func(context.Context, string) error { return nil },
		SaveBuildLogFunc: func(context.Context, string, []byte) error { return nil },
	}
	gitsvc := &mocks.GitServiceMock{
		CloneRepositoryFunc: func(ctx context.Context, _ string, _ *domain.Repository, _ string) error {
			cloned.Store(true)
			<-ctx.Done()
			return ctx.Err()
		},
	}
	s := newService(&Config{StepTimeout: time.Hour}, client, nil, &reclaimingBuildpack{}, gitsvc, builder.ImageConfig{})
	require.NoError(t, s.Start(context.Background()))
	require.NoError(t, s.Shutdown(context.Background()))

	app := &domain.Application{ID: domain.NewID(), Config: domain.ApplicationConfig{BuildConfig: &domain.BuildConfigRuntimeBuildpack{}}}
	err := s.startBuild(&domain.StartBuildRequest{
		App:   app,
		Build: &domain.Build{ID: domain.NewID(), ApplicationID: app.ID},
		Repo:  &domain.Repository{},
	})
	require.Error(t, err, "should refuse a build request after Shutdown returned")
	time.Sleep(300 * time.Millisecond)
	require.False(t, cloned.Load(), "should not start a build after Shutdown returned")
}
