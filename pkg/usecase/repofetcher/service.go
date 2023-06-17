package repofetcher

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc/pool"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/usecase/cdservice"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

const (
	fetcherConcurrency = 10
	fetcherInterval    = 5 * time.Minute
)

// fetchInterval calculates the interval the repository should be fetched in, given the last activity time.
func fetchInterval(now time.Time, updatedAt time.Time) int {
	// note: to reach maximum interval of 1 day, 288 days needs to be elapsed from the last update
	const maxInterval = int(24 * time.Hour / fetcherInterval)
	epochsElapsed := float64(now.Sub(updatedAt)) / float64(fetcherInterval)
	if epochsElapsed < 0 {
		epochsElapsed = 0
	}
	interval := int(math.Ceil(math.Sqrt(epochsElapsed)))
	return lo.Clamp(interval, 1, maxInterval)
}

type Service interface {
	Run()
	Fetch(repositoryID string)
	Stop(ctx context.Context) error
}

type service struct {
	appRepo domain.ApplicationRepository
	gitRepo domain.GitRepositoryRepository
	pubKey  *ssh.PublicKeys
	cd      cdservice.Service

	fetcher   chan<- string
	run       func()
	runOnce   sync.Once
	close     func()
	closeOnce sync.Once
}

func NewService(
	appRepo domain.ApplicationRepository,
	gitRepo domain.GitRepositoryRepository,
	pubKey *ssh.PublicKeys,
	cd cdservice.Service,
) (Service, error) {
	r := &service{
		appRepo: appRepo,
		gitRepo: gitRepo,
		pubKey:  pubKey,
		cd:      cd,
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

func (r *service) resolveRefs(ctx context.Context, repo *domain.Repository) (refToCommit map[string]string, err error) {
	auth, err := domain.GitAuthMethod(repo, r.pubKey)
	if err != nil {
		return nil, err
	}
	remote := git.NewRemote(nil, &config.RemoteConfig{
		Name: "origin",
		URLs: []string{repo.URL},
	})
	refs, err := remote.ListContext(ctx, &git.ListOptions{
		Auth: auth,
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to list remote refs at %v", repo.URL))
	}

	refToCommit = make(map[string]string, 2*len(refs))
	for _, ref := range refs {
		if ref.Type() == plumbing.HashReference {
			refToCommit[ref.Name().String()] = ref.Hash().String()
			refToCommit[ref.Name().Short()] = ref.Hash().String()
		}
	}
	for _, ref := range refs {
		if ref.Type() == plumbing.SymbolicReference {
			commit, ok := refToCommit[ref.Target().String()]
			if ok {
				refToCommit[ref.Name().String()] = commit
			}
		}
	}
	return refToCommit, nil
}

func (r *service) updateApps(ctx context.Context, repo *domain.Repository, apps []*domain.Application) error {
	refToCommit, err := r.resolveRefs(ctx, repo)
	if err != nil {
		return err
	}

	for _, app := range apps {
		commit, ok := refToCommit[app.RefName]
		if !ok {
			log.Errorf("failed to get resolve ref %v for app %v", app.RefName, app.ID)
			commit = domain.EmptyCommit // Mark as empty commit to signal error
		}
		err = r.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{WantCommit: optional.From(commit)})
		if err != nil {
			return errors.Wrap(err, "failed to update application")
		}
		// Notify builds
		r.cd.RegisterBuild(app.ID)
	}
	return nil
}

func (r *service) Stop(_ context.Context) error {
	r.closeOnce.Do(r.close)
	return nil
}
