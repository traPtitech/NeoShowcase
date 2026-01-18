package commitfetcher

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/discovery"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type Service interface {
	Run()
	Fetch(repositoryID string, hashes []string)
	Stop(ctx context.Context) error
}

const queueMax = 1000

type queueItem struct {
	repositoryID string
	hashes       []string
}

type service struct {
	cluster     *discovery.Cluster
	appRepo     domain.ApplicationRepository
	buildRepo   domain.BuildRepository
	gitRepo     domain.GitRepositoryRepository
	commitsRepo domain.RepositoryCommitRepository
	gitsvc      domain.GitService

	queue chan<- queueItem

	run       func()
	runOnce   sync.Once
	close     func()
	closeOnce sync.Once
}

func NewService(
	cluster *discovery.Cluster,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	gitRepo domain.GitRepositoryRepository,
	commitsRepo domain.RepositoryCommitRepository,
	gitsvc domain.GitService,
) (Service, error) {
	s := &service{
		cluster:     cluster,
		appRepo:     appRepo,
		buildRepo:   buildRepo,
		gitRepo:     gitRepo,
		commitsRepo: commitsRepo,
		gitsvc:      gitsvc,
	}

	q := make(chan queueItem, queueMax)
	s.queue = q

	ctx, cancel := context.WithCancel(context.Background())
	s.run = func() {
		go s.resolveCommits(ctx)
		go s.fetchLoop(ctx, q)
	}
	s.close = cancel

	return s, nil
}

func (s *service) Run() {
	s.runOnce.Do(s.run)
}

// resolveCommits resolves all recorded commits in current database - i.e. applications and builds table
func (s *service) resolveCommits(ctx context.Context) {
	apps, err := s.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
	if err != nil {
		slog.ErrorContext(ctx, "failed to get applications", "error", err)
		return
	}

	for _, app := range apps {
		// Shard by repository ID
		if !s.cluster.IsAssigned(app.RepositoryID) {
			continue
		}

		builds, err := s.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{ApplicationID: optional.From(app.ID)})
		if err != nil {
			slog.ErrorContext(ctx, "failed to get builds", "error", err)
			return
		}

		hashMap := lo.SliceToMap(builds, func(b *domain.Build) (string, struct{}) {
			return b.Commit, struct{}{}
		})
		hashMap[app.Commit] = struct{}{}
		hashes := lo.Keys(hashMap)

		s.Fetch(app.RepositoryID, hashes)
		time.Sleep(time.Second)
	}
}

func (s *service) Fetch(repositoryID string, hashes []string) {
	select {
	case s.queue <- queueItem{repositoryID, hashes}:
	default:
		slog.Warn("commit fetcher: queue is full, skipping request for repository and hashes", "repository_id", repositoryID, "hash_count", len(hashes))
	}
}

func (s *service) fetchLoop(ctx context.Context, fetcher <-chan queueItem) {
	for {
		select {
		case item := <-fetcher:
			err := s.fetchOne(ctx, item.repositoryID, item.hashes)
			if err != nil {
				slog.ErrorContext(ctx, "failed to fetch commits for repository", "count", len(item.hashes), "repository_id", item.repositoryID, "error", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *service) fetchOne(ctx context.Context, repositoryID string, hashes []string) error {
	start := time.Now()

	// Filter out errored app hashes
	hashes = lo.Filter(hashes, func(hash string, _ int) bool { return hash != domain.EmptyCommit })
	hashes = lo.Uniq(hashes) // Make them unique just in case

	// Check if we have already tried
	recordedCommits, err := s.commitsRepo.GetCommits(ctx, hashes)
	if err != nil {
		return errors.Wrap(err, "failed to get recorded commits")
	}
	recordedCommitMap := lo.SliceToMap(recordedCommits, func(c *domain.RepositoryCommit) (string, bool) {
		return c.Hash, true
	})
	hashes = lo.Filter(hashes, func(hash string, _ int) bool {
		return !recordedCommitMap[hash]
	})
	if len(hashes) == 0 {
		return nil
	}

	repo, err := s.gitRepo.GetRepository(ctx, repositoryID)
	if err != nil {
		return errors.Wrap(err, "failed to get repository")
	}

	// Init local git directory
	tmpDir, err := os.MkdirTemp("", "commit-fetcher-")
	if err != nil {
		return errors.Wrap(err, "failed to create temp dir")
	}
	defer os.RemoveAll(tmpDir)

	localRepo, err := s.gitsvc.CreateBareRepository(tmpDir, repo)
	if err != nil {
		return errors.Wrap(err, "failed to initialize git repo")
	}

	err = localRepo.Fetch(ctx, hashes)
	if err != nil {
		return errors.Wrap(err, "failed to fetch")
	}

	// Get commit objects and record, in a *fail-safe* manner -
	// this prevents spam-cloning of remote repository
	for _, hash := range hashes {
		commit, err := localRepo.GetCommit(hash)
		if err != nil {
			slog.ErrorContext(ctx, "failed to fetch commit for repository", "hash", hash, "repository_id", repositoryID, "error", err)
			commit = domain.ToErroredRepositoryCommit(hash)
		}

		err = s.commitsRepo.RecordCommit(ctx, commit)
		if err != nil {
			return errors.Wrap(err, "failed to record commit")
		}
	}

	slog.DebugContext(ctx, "commit fetcher: fetched commits for repository", "count", len(hashes), "repository_id", repositoryID, "duration", time.Since(start))
	return nil
}

func (s *service) Stop(_ context.Context) error {
	s.closeOnce.Do(s.close)
	return nil
}
