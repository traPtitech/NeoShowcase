package builder

import (
	"context"
	"fmt"
	"time"

	"github.com/friendsofgo/errors"
	buildkit "github.com/moby/buildkit/client"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pbconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/cli"
)

func (s *builderService) startBuild(req *domain.StartBuildRequest) error {
	s.statusLock.Lock()
	defer s.statusLock.Unlock()

	if s.state != nil {
		log.Warnf("Skipping build request for %v, builder busy - builder scheduling may be malfunctioning?", req.Build.ID)
		return nil // Builder busy - skip
	}

	log.Infof("Starting build for %v", req.Build.ID)

	st, err := newState(req.App, req.Envs, req.Build, req.Repo, s.client)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.state = st
	s.stateCancel = func() {
		cancel()
		st.Wait()
	}

	go func() {
		status := s.process(ctx, st)
		s.finalize(context.Background(), st) // don't want finalization tasks to be cancelled
		st.Done()

		cancel()
		s.statusLock.Lock()
		s.state = nil
		s.stateCancel = nil
		s.statusLock.Unlock()
		log.Infof("Build settled for %v", st.build.ID)
		// Send settled response *after* unlocking internal state for next build
		s.response <- &pb.BuilderResponse{Type: pb.BuilderResponse_BUILD_SETTLED, Body: &pb.BuilderResponse_Settled{Settled: &pb.BuildSettled{
			BuildId: st.build.ID,
			Status:  pbconvert.BuildStatusMapper.IntoMust(status),
		}}}
	}()

	return nil
}

type buildStep struct {
	desc string
	fn   func(ctx context.Context) error
}

func (s *builderService) buildSteps(st *state) ([]buildStep, error) {
	var steps []buildStep

	steps = append(steps, buildStep{"Repository Clone", func(ctx context.Context) error {
		return s.cloneRepository(ctx, st)
	}})

	switch bc := st.app.Config.BuildConfig.(type) {
	case *domain.BuildConfigRuntimeBuildpack:
		steps = append(steps, buildStep{"Build (Runtime Buildpack)", func(ctx context.Context) error {
			return s.buildRuntimeBuildpack(ctx, st, bc)
		}})
	case *domain.BuildConfigRuntimeCmd:
		steps = append(steps, buildStep{"Build (Runtime Command)", func(ctx context.Context) error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildRuntimeCmd(ctx, st, ch, bc)
			})
		}})
	case *domain.BuildConfigRuntimeDockerfile:
		steps = append(steps, buildStep{"Build (Runtime Dockerfile)", func(ctx context.Context) error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildRuntimeDockerfile(ctx, st, ch, bc)
			})
		}})
	case *domain.BuildConfigStaticBuildpack:
		steps = append(steps, buildStep{"Build (Static Buildpack)", func(ctx context.Context) error {
			return s.buildStaticBuildpackPack(ctx, st, bc)
		}})
		steps = append(steps, buildStep{"Extract from Temporary Image", func(ctx context.Context) error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildStaticExtract(ctx, st, ch)
			})
		}})
		steps = append(steps, buildStep{"Cleanup Temporary Image", func(ctx context.Context) error {
			return s.buildStaticCleanup(ctx, st)
		}})
		steps = append(steps, buildStep{"Save Artifact", func(ctx context.Context) error {
			return s.saveArtifact(ctx, st)
		}})
	case *domain.BuildConfigStaticCmd:
		steps = append(steps, buildStep{"Build (Static Command)", func(ctx context.Context) error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildStaticCmd(ctx, st, ch, bc)
			})
		}})
		steps = append(steps, buildStep{"Extract from Temporary Image", func(ctx context.Context) error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildStaticExtract(ctx, st, ch)
			})
		}})
		steps = append(steps, buildStep{"Cleanup Temporary Image", func(ctx context.Context) error {
			return s.buildStaticCleanup(ctx, st)
		}})
		steps = append(steps, buildStep{"Save Artifact", func(ctx context.Context) error {
			return s.saveArtifact(ctx, st)
		}})
	case *domain.BuildConfigStaticDockerfile:
		steps = append(steps, buildStep{"Build (Static Dockerfile)", func(ctx context.Context) error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildStaticDockerfile(ctx, st, ch, bc)
			})
		}})
		steps = append(steps, buildStep{"Extract from Temporary Image", func(ctx context.Context) error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildStaticExtract(ctx, st, ch)
			})
		}})
		steps = append(steps, buildStep{"Cleanup Temporary Image", func(ctx context.Context) error {
			return s.buildStaticCleanup(ctx, st)
		}})
		steps = append(steps, buildStep{"Save Artifact", func(ctx context.Context) error {
			return s.saveArtifact(ctx, st)
		}})
	default:
		return nil, errors.New("unknown build config type")
	}

	return steps, nil
}

func (s *builderService) process(ctx context.Context, st *state) domain.BuildStatus {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go s.buildPingLoop(ctx, st.build.ID)

	version, revision := cli.GetVersion()
	st.WriteLog(fmt.Sprintf("[ns-builder] Version %v (%v)", version, revision))

	steps, err := s.buildSteps(st)
	if err != nil {
		log.Errorf("calculating build steps: %+v", err)
		st.WriteLog(fmt.Sprintf("[ns-builder] Error calculating build steps: %v", err))
		return domain.BuildStatusFailed
	}

	for i, step := range steps {
		st.WriteLog(fmt.Sprintf("[ns-builder] ==> (%d/%d) %s", i+1, len(steps), step.desc))
		start := time.Now()

		// Execute step with timeout
		childCtx, cancel := context.WithTimeout(ctx, s.config.StepTimeout)
		err := step.fn(childCtx)
		cancel()

		// First, check if ctx was cancelled from parent
		// - this (usually) means the user pressed 'Cancel' button to cancel the build
		if cerr := ctx.Err(); cerr != nil {
			st.WriteLog(fmt.Sprintf("[ns-builder] Build cancelled by user: %v", cerr))
			return domain.BuildStatusCanceled
		}
		// Check if the childCtx timed out from the configured "StepTimeout"
		if errors.Is(err, context.DeadlineExceeded) {
			st.WriteLog(fmt.Sprintf("[ns-builder] Build timed out: %v", err))
			return domain.BuildStatusFailed
		}
		// Other reasons such as build script failure
		if err != nil {
			log.Errorf("Build failed for %v: %+v", st.build.ID, err)
			st.WriteLog(fmt.Sprintf("[ns-builder] Build failed: %v", err))
			return domain.BuildStatusFailed
		}

		st.WriteLog(fmt.Sprintf("[ns-builder] ==> (%d/%d) Done (%v).", i+1, len(steps), time.Since(start)))
	}

	return domain.BuildStatusSucceeded
}

func (s *builderService) buildPingLoop(ctx context.Context, buildID string) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := s.client.PingBuild(ctx, buildID)
			if err != nil {
				log.Errorf("pinging build: %+v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *builderService) finalize(ctx context.Context, st *state) {
	err := s.client.SaveBuildLog(ctx, st.build.ID, st.logWriter.buf.Bytes())
	if err != nil {
		log.Errorf("failed to save build log: %+v", err)
	}
}
