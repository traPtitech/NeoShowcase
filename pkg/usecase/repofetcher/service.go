package repofetcher

import (
	"context"
	"fmt"
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

const fetcherConcurrency = 10

type Service interface {
	Run()
	Fetch(repositoryIDs []string)
	Stop(ctx context.Context) error
}

type service struct {
	appRepo domain.ApplicationRepository
	gitRepo domain.GitRepositoryRepository
	pubKey  *ssh.PublicKeys
	cd      cdservice.Service

	fetcher   chan<- []string
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

	fetcher := make(chan []string, 100)
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

func (r *service) Fetch(repositoryIDs []string) {
	// non-blocking; fetch operation is eventually consistent
	select {
	case r.fetcher <- repositoryIDs:
	default:
	}
}

func (r *service) fetchLoop(ctx context.Context, fetcher <-chan []string) {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	doSync := func(repositoryIDs optional.Of[[]string]) {
		start := time.Now()
		if err := r.fetch(repositoryIDs); err != nil {
			log.Errorf("failed to fetch repositories: %+v", err)
			return
		}
		log.Infof("Fetched repositories in %v", time.Since(start))

		r.cd.RegisterBuilds()
	}

	doSync(optional.None[[]string]())

	for {
		select {
		case ids := <-fetcher:
			// coalesce specific repository fetch request
			additional := lo.Flatten(ds.ReadAll(fetcher))
			ids = append(ids, additional...)
			doSync(optional.From(ids))
		case <-ticker.C:
			doSync(optional.None[[]string]())
		case <-ctx.Done():
			return
		}
	}
}

func (r *service) fetch(repositoryIDs optional.Of[[]string]) error {
	ctx := context.Background()

	var getCond domain.GetRepositoryCondition
	if repositoryIDs.Valid {
		getCond.IDs = repositoryIDs
	}
	repositories, err := r.gitRepo.GetRepositories(ctx, getCond)
	if err != nil {
		return err
	}
	repos := lo.SliceToMap(repositories, func(repo *domain.Repository) (string, *domain.Repository) { return repo.ID, repo })

	applications, err := r.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
	if err != nil {
		return err
	}
	reposToApps := make(map[string][]*domain.Application)
	for _, app := range applications {
		reposToApps[app.RepositoryID] = append(reposToApps[app.RepositoryID], app)
	}

	p := pool.New().WithMaxGoroutines(fetcherConcurrency)
	for _, repo := range repos {
		repo := repo
		p.Go(func() {
			err := r.updateApps(ctx, repo, reposToApps[repo.ID])
			if err != nil {
				log.Warnf("failed to update repo: %+v", err)
			}
		})
	}
	p.Wait()
	return nil
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
	if len(apps) == 0 {
		return nil
	}

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
	}
	return nil
}

func (r *service) Stop(_ context.Context) error {
	r.closeOnce.Do(r.close)
	return nil
}
