package repofetcher

import (
	"context"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc/pool"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	commitfetcher "github.com/traPtitech/neoshowcase/pkg/usecase/commit-fetcher"
	"github.com/traPtitech/neoshowcase/pkg/util/discovery"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

const (
	fetcherConcurrency = 10
	fetcherInterval    = 3 * time.Minute
)

// fetchInterval calculates the interval the repository should be fetched in, given the last activity time.
func fetchInterval(now time.Time, updatedAt time.Time) int {
	const intervalIncreasePeriod = 24 * time.Hour
	const maxInterval = 24 * time.Hour
	const maxPeriods = maxInterval / fetcherInterval

	fullPeriodsElapsed := int(now.Sub(updatedAt)) / int(intervalIncreasePeriod)
	fullPeriodsElapsed = max(fullPeriodsElapsed, 0)

	return min(fullPeriodsElapsed+1, int(maxPeriods))
}

// Service fetches commit metadata sequentially.
type Service interface {
	Run()
	Fetch(repositoryID string)
	Stop(ctx context.Context) error
}

type service struct {
	cluster       *discovery.Cluster
	appRepo       domain.ApplicationRepository
	gitRepo       domain.GitRepositoryRepository
	gitsvc        domain.GitService
	cd            domain.CDService
	commitFetcher commitfetcher.Service

	fetcher   chan<- string
	run       func()
	runOnce   sync.Once
	close     func()
	closeOnce sync.Once
}

func NewService(
	cluster *discovery.Cluster,
	appRepo domain.ApplicationRepository,
	gitRepo domain.GitRepositoryRepository,
	cd domain.CDService,
	commitFetcher commitfetcher.Service,
	gitsvc domain.GitService,
) (Service, error) {
	r := &service{
		cluster:       cluster,
		appRepo:       appRepo,
		gitRepo:       gitRepo,
		cd:            cd,
		commitFetcher: commitFetcher,
		gitsvc:        gitsvc,
	}

	fetcher := make(chan string, 100)
	r.fetcher = fetcher
	ctx, cancel := context.WithCancel(context.Background())
	r.run = func() {
		go r.fetchLoop(ctx, fetcher)
	}
	r.close = cancel

	return r, nil
}

func (r *service) Run() {
	r.runOnce.Do(r.run)
}

func (r *service) Fetch(repositoryID string) {
	// non-blocking; fetch operation is eventually consistent
	select {
	case r.fetcher <- repositoryID:
	default:
	}
}

func (r *service) fetchLoop(ctx context.Context, fetcher <-chan string) {
	ticker := time.NewTicker(fetcherInterval)
	defer ticker.Stop()

	epoch := 0
	runFetchEpoch := func() {
		start := time.Now()
		n, err := r.doFetchEpoch(ctx, epoch)
		if err != nil {
			log.Errorf("failed to fetch repositories: %+v", err)
		}
		log.Infof("Fetched %v repositories in %v", n, time.Since(start))
		epoch++
	}

	runFetchEpoch()

	for {
		select {
		case id := <-fetcher:
			err := r.fetchOne(ctx, id)
			if err != nil {
				log.Errorf("failed to fetch repository %v: %+v", id, err)
			}
		case <-ticker.C:
			runFetchEpoch()
		case <-ctx.Done():
			return
		}
	}
}

func (r *service) doFetchEpoch(ctx context.Context, epoch int) (int, error) {
	repos, err := r.gitRepo.GetRepositories(ctx, domain.GetRepositoryCondition{})
	if err != nil {
		return 0, err
	}
	apps, err := r.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
	if err != nil {
		return 0, err
	}
	repoMap := lo.SliceToMap(repos, func(repo *domain.Repository) (string, *domain.Repository) { return repo.ID, repo })
	repoToApps := make(map[string][]*domain.Application, len(repos))
	for _, app := range apps {
		repoToApps[app.RepositoryID] = append(repoToApps[app.RepositoryID], app)
	}
	repoLastUpdated := lo.MapValues(repoToApps, func(apps []*domain.Application, repoID string) time.Time {
		// assert len(apps) > 0
		updatedAts := ds.Map(apps, func(app *domain.Application) time.Time { return app.UpdatedAt })
		return lo.MaxBy(updatedAts, func(a, b time.Time) bool { return a.Compare(b) < 0 })
	})

	now := time.Now()
	count := 0
	p := pool.New().WithMaxGoroutines(fetcherConcurrency)
	for repoID, lastUpdated := range repoLastUpdated {
		// Shard by repo ID
		if !r.cluster.IsAssigned(repoID) {
			continue
		}

		interval := fetchInterval(now, lastUpdated)
		if epoch%interval == 0 {
			repo := repoMap[repoID]
			count++
			p.Go(func() { // careful not to capture loop variable
				apps := repoToApps[repo.ID]
				someAppRunning := lo.ContainsBy(apps, func(app *domain.Application) bool { return app.Running })
				if !someAppRunning {
					return
				}
				err := r.updateApps(ctx, repo, apps)
				if err != nil {
					log.Warnf("failed to update repo: %+v", err)
				}
			})
		}
	}
	p.Wait()
	return count, nil
}

func (r *service) fetchOne(ctx context.Context, repositoryID string) error {
	repo, err := r.gitRepo.GetRepository(ctx, repositoryID)
	if err != nil {
		return err
	}
	apps, err := r.appRepo.GetApplications(ctx, domain.GetApplicationCondition{
		RepositoryID: optional.From(repositoryID),
	})
	if err != nil {
		return err
	}
	return r.updateApps(ctx, repo, apps)
}

func (r *service) updateApps(ctx context.Context, repo *domain.Repository, apps []*domain.Application) error {
	refToCommit, err := r.gitsvc.ResolveRefs(ctx, repo)
	if err != nil {
		return err
	}

	var hashes []string
	for _, app := range apps {
		commit, ok := refToCommit[app.RefName]
		if ok {
			hashes = append(hashes, commit)
		} else {
			log.Errorf("failed to get resolve ref %v for app %v", app.RefName, app.ID)
			commit = domain.EmptyCommit // Mark as empty commit to signal error
		}

		err = r.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{Commit: optional.From(commit)})
		if err != nil {
			return errors.Wrap(err, "failed to update application")
		}
		// Notify builds
		r.cd.RegisterBuild(app.ID)
	}

	// Notify to fetch commit metadata, if missing
	r.commitFetcher.Fetch(repo.ID, hashes)

	return nil
}

func (r *service) Stop(_ context.Context) error {
	r.closeOnce.Do(r.close)
	return nil
}
