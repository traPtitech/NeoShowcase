package builder

import (
	"context"
	"fmt"
	"time"

	"github.com/friendsofgo/errors"
	buildkit "github.com/moby/buildkit/client"
	log "github.com/sirupsen/logrus"
	gstatus "google.golang.org/grpc/status"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func (s *builderService) tryStartBuild(buildID string) error {
	s.statusLock.Lock()
	defer s.statusLock.Unlock()

	if s.state != nil {
		log.Infof("skipping build request for %v, builder busy", buildID)
		return nil // Builder busy - skip
	}

	now := time.Now()
	n, err := s.buildRepo.UpdateBuild(context.Background(), domain.GetBuildCondition{
		ID:     optional.From(buildID),
		Status: optional.From(domain.BuildStatusQueued),
	}, domain.UpdateBuildArgs{
		Status:    optional.From(domain.BuildStatusBuilding),
		StartedAt: optional.From(now),
		UpdatedAt: optional.From(now),
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return nil // other builder has acquired the build lock - skip
	}

	// Acquired build lock
	log.Infof("Starting build for %v", buildID)

	build, err := s.buildRepo.GetBuild(context.Background(), buildID)
	if err != nil {
		return err
	}
	app, err := s.appRepo.GetApplication(context.Background(), build.ApplicationID)
	if err != nil {
		return err
	}
	repo, err := s.gitRepo.GetRepository(context.Background(), app.RepositoryID)
	if err != nil {
		return err
	}

	st, err := newState(app, build, repo, s.response)
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
		s.response <- &pb.BuilderResponse{Type: pb.BuilderResponse_BUILD_STARTED, Body: &pb.BuilderResponse_Started{Started: &pb.BuildStarted{
			BuildId: buildID,
		}}}

		status := s.process(ctx, st)
		s.finalize(context.Background(), st, status) // don't want finalization tasks to be cancelled
		st.Done()

		cancel()
		s.statusLock.Lock()
		s.state = nil
		s.stateCancel = nil
		s.statusLock.Unlock()
		log.Infof("Build settled for %v", buildID)
		// Send settled response *after* unlocking internal state for next build
		s.response <- &pb.BuilderResponse{Type: pb.BuilderResponse_BUILD_SETTLED, Body: &pb.BuilderResponse_Settled{Settled: &pb.BuildSettled{
			BuildId: buildID,
			Reason:  toPBSettleReason(status),
		}}}
	}()

	return nil
}

func toPBSettleReason(status domain.BuildStatus) pb.BuildSettled_Reason {
	switch status {
	case domain.BuildStatusSucceeded:
		return pb.BuildSettled_SUCCESS
	case domain.BuildStatusFailed:
		return pb.BuildSettled_FAILED
	case domain.BuildStatusCanceled:
		return pb.BuildSettled_CANCELLED
	default:
		panic(fmt.Sprintf("unexpected settled status: %v", status))
	}
}

type buildStep struct {
	desc string
	fn   func() error
}

func (s *builderService) buildSteps(ctx context.Context, st *state) ([]buildStep, error) {
	var steps []buildStep

	steps = append(steps, buildStep{"Repository Clone", func() error {
		return s.cloneRepository(ctx, st)
	}})

	switch bc := st.app.Config.BuildConfig.(type) {
	case *domain.BuildConfigRuntimeBuildpack:
		steps = append(steps, buildStep{"Build (Runtime Buildpack)", func() error {
			return s.buildRuntimeBuildpack(ctx, st, bc)
		}})
	case *domain.BuildConfigRuntimeCmd:
		steps = append(steps, buildStep{"Build (Runtime Command)", func() error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildRuntimeCmd(ctx, st, ch, bc)
			})
		}})
	case *domain.BuildConfigRuntimeDockerfile:
		steps = append(steps, buildStep{"Build (Runtime Dockerfile)", func() error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildRuntimeDockerfile(ctx, st, ch, bc)
			})
		}})
	case *domain.BuildConfigStaticBuildpack:
		steps = append(steps, buildStep{"Build (Static Buildpack)", func() error {
			return s.buildStaticBuildpackPack(ctx, st, bc)
		}})
		steps = append(steps, buildStep{"Extract from Temporary Image", func() error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildStaticExtract(ctx, st, ch)
			})
		}})
		steps = append(steps, buildStep{"Cleanup Temporary Image", func() error {
			return s.buildStaticCleanup(ctx, st)
		}})
		steps = append(steps, buildStep{"Save Artifact", func() error {
			return s.saveArtifact(ctx, st)
		}})
	case *domain.BuildConfigStaticCmd:
		steps = append(steps, buildStep{"Build (Static Command)", func() error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildStaticCmd(ctx, st, ch, bc)
			})
		}})
		steps = append(steps, buildStep{"Extract from Temporary Image", func() error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildStaticExtract(ctx, st, ch)
			})
		}})
		steps = append(steps, buildStep{"Cleanup Temporary Image", func() error {
			return s.buildStaticCleanup(ctx, st)
		}})
		steps = append(steps, buildStep{"Save Artifact", func() error {
			return s.saveArtifact(ctx, st)
		}})
	case *domain.BuildConfigStaticDockerfile:
		steps = append(steps, buildStep{"Build (Static Dockerfile)", func() error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildStaticDockerfile(ctx, st, ch, bc)
			})
		}})
		steps = append(steps, buildStep{"Extract from Temporary Image", func() error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildStaticExtract(ctx, st, ch)
			})
		}})
		steps = append(steps, buildStep{"Cleanup Temporary Image", func() error {
			return s.buildStaticCleanup(ctx, st)
		}})
		steps = append(steps, buildStep{"Save Artifact", func() error {
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
	go s.updateStatusLoop(ctx, st.build.ID)

	steps, err := s.buildSteps(ctx, st)
	if err != nil {
		log.Errorf("calculating build steps: %+v", err)
		st.WriteLog(fmt.Sprintf("[ns-builder] Error calculating build steps: %v", err))
		return domain.BuildStatusFailed
	}

	for i, step := range steps {
		st.WriteLog(fmt.Sprintf("[ns-builder] ==> (%d/%d) %s", i+1, len(steps), step.desc))
		start := time.Now()
		err := step.fn()
		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, gstatus.FromContextError(context.Canceled).Err()) {
			st.WriteLog("[ns-builder] Build cancelled.")
			return domain.BuildStatusCanceled
		}
		if err != nil {
			log.Errorf("Build failed for %v: %+v", st.build.ID, err)
			st.WriteLog(fmt.Sprintf("[ns-builder] Build failed: %v", err))
			return domain.BuildStatusFailed
		}
		st.WriteLog(fmt.Sprintf("[ns-builder] ==> (%d/%d) Done (%v).", i+1, len(steps), time.Since(start)))
	}

	return domain.BuildStatusSucceeded
}

func (s *builderService) updateStatusLoop(ctx context.Context, buildID string) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_, err := s.buildRepo.UpdateBuild(ctx,
				domain.GetBuildCondition{ID: optional.From(buildID)},
				domain.UpdateBuildArgs{UpdatedAt: optional.From(time.Now())})
			if err != nil {
				log.Errorf("failed to update build time: %+v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *builderService) finalize(ctx context.Context, st *state, status domain.BuildStatus) {
	err := domain.SaveBuildLog(s.storage, st.build.ID, st.logWriter.LogReader())
	if err != nil {
		log.Errorf("failed to save build log: %+v", err)
	}

	now := time.Now()
	updateCond := domain.GetBuildCondition{ID: optional.From(st.build.ID)}
	updateArgs := domain.UpdateBuildArgs{
		Status:     optional.From(status),
		UpdatedAt:  optional.From(now),
		FinishedAt: optional.From(now),
	}
	if _, err = s.buildRepo.UpdateBuild(ctx, updateCond, updateArgs); err != nil {
		log.Errorf("failed to update build: %+v", err)
	}
}
