package commitfetcher

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"

	"github.com/traPtitech/neoshowcase/pkg/domain"
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
	appRepo     domain.ApplicationRepository
	buildRepo   domain.BuildRepository
	gitRepo     domain.GitRepositoryRepository
	commitsRepo domain.RepositoryCommitRepository
	fallbackKey *ssh.PublicKeys

	queue chan<- queueItem

	run       func()
	runOnce   sync.Once
	close     func()
	closeOnce sync.Once
}

func NewService(
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	gitRepo domain.GitRepositoryRepository,
	commitsRepo domain.RepositoryCommitRepository,
	fallbackKey *ssh.PublicKeys,
) (Service, error) {
	s := &service{
		appRepo:     appRepo,
		buildRepo:   buildRepo,
		gitRepo:     gitRepo,
		commitsRepo: commitsRepo,
		fallbackKey: fallbackKey,
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
		log.Errorf("failed to get applications: %+v", err)
		return
	}

	for _, app := range apps {
		builds, err := s.buildRepo.GetBuilds(ctx, domain.GetBuildCondition{ApplicationID: optional.From(app.ID)})
		if err != nil {
			log.Errorf("failed to get builds: %+v", err)
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
		log.Warnf("commit fetcher: queue is full, skipping request for repository %s and %d hashes", repositoryID, len(hashes))
	}
}

func (s *service) fetchLoop(ctx context.Context, fetcher <-chan queueItem) {
	for {
		select {
		case item := <-fetcher:
			err := s.fetchOne(ctx, item.repositoryID, item.hashes)
			if err != nil {
				log.Errorf("failed to fetch %d commits for repository %v: %v", len(item.hashes), item.repositoryID, err)
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

	// Get repository auth
	repo, err := s.gitRepo.GetRepository(ctx, repositoryID)
	if err != nil {
		return errors.Wrap(err, "failed to get repository")
	}
	auth, err := domain.GitAuthMethod(repo, s.fallbackKey)
	if err != nil {
		return errors.Wrap(err, "failed to calculate git auth")
	}

	// Init local git directory
	tmpDir, err := os.MkdirTemp("", "commit-fetcher-")
	if err != nil {
		return errors.Wrap(err, "failed to create temp dir")
	}
	defer os.RemoveAll(tmpDir)

	// Set up remote...
	localRepo, err := git.PlainInit(tmpDir, true)
	if err != nil {
		return errors.Wrap(err, "failed to init repo")
	}
	remote, err := localRepo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{repo.URL},
	})
	if err != nil {
		return errors.Wrap(err, "failed to create remote")
	}

	// Shallow clone at the commits, in order to get the commit object
	refSpecs := lo.Map(hashes, func(hash string, idx int) config.RefSpec {
		targetRef := plumbing.NewRemoteReferenceName("origin", fmt.Sprintf("target-%d", idx))
		return config.RefSpec(fmt.Sprintf("+%s:%s", hash, targetRef))
	})
	err = remote.FetchContext(ctx, &git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   refSpecs,
		Depth:      1,
		Auth:       auth,
	})
	if err != nil {
		return errors.Wrap(err, "failed to fetch")
	}

	// Get commit objects and record, in a *fail-safe* manner -
	// this prevents spam-cloning of remote repository
	for _, hash := range hashes {
		commitObj, err := localRepo.CommitObject(plumbing.NewHash(hash))

		// Convert commit object
		var commit *domain.RepositoryCommit
		if err == nil {
			commit = domain.ToRepositoryCommit(commitObj)
		} else {
			log.Errorf("failed to fetch commit %v for repository %v: %+v", hash, repositoryID, err)
			commit = domain.ToErroredRepositoryCommit(hash)
		}

		// Record commit object
		err = s.commitsRepo.RecordCommit(ctx, commit)
		if err != nil {
			return errors.Wrap(err, "failed to record commit")
		}
	}

	log.Debugf("commit fetcher: fetched %v commit(s) for repository %v in %v", len(hashes), repositoryID, time.Since(start))
	return nil
}

func (s *service) Stop(_ context.Context) error {
	s.closeOnce.Do(s.close)
	return nil
}
