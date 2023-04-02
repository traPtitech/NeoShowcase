package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type RepositoryFetcherService interface {
	Run()
	Stop(ctx context.Context) error
}

type repositoryFetcherService struct {
	bus     domain.Bus
	appRepo domain.ApplicationRepository
	gitRepo domain.GitRepositoryRepository
	pubKey  *ssh.PublicKeys

	run       func()
	runOnce   sync.Once
	close     func()
	closeOnce sync.Once
}

func NewRepositoryFetcherService(
	bus domain.Bus,
	appRepo domain.ApplicationRepository,
	gitRepo domain.GitRepositoryRepository,
	pubKey *ssh.PublicKeys,
) (RepositoryFetcherService, error) {
	r := &repositoryFetcherService{
		bus:     bus,
		appRepo: appRepo,
		gitRepo: gitRepo,
		pubKey:  pubKey,
	}

	closer := make(chan struct{})
	r.run = func() {
		go r.fetchLoop(closer)
	}
	r.close = func() {
		close(closer)
	}

	return r, nil
}

func (r *repositoryFetcherService) Run() {
	r.runOnce.Do(r.run)
}

func (r *repositoryFetcherService) fetchLoop(closer <-chan struct{}) {
	sub := r.bus.Subscribe(event.FetcherFetchRequest)
	defer sub.Unsubscribe()
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	doSync := func() {
		start := time.Now()
		if err := r.fetchAll(); err != nil {
			log.Errorf("failed to fetch repositories: %+v", err)
			return
		}
		log.Infof("Fetched repositories in %v", time.Since(start))
		r.bus.Publish(event.FetcherFetchComplete, nil)
	}

	doSync()

	for {
		select {
		case <-sub.Chan():
			doSync()
		case <-ticker.C:
			doSync()
		case <-closer:
			return
		}
	}
}

func (r *repositoryFetcherService) fetchAll() error {
	ctx := context.Background()
	repositories, err := r.gitRepo.GetRepositories(ctx, domain.GetRepositoryCondition{})
	if err != nil {
		return err
	}
	applications, err := r.appRepo.GetApplications(ctx, domain.GetApplicationCondition{})
	if err != nil {
		return err
	}

	repos := lo.SliceToMap(repositories, func(repo *domain.Repository) (string, *domain.Repository) { return repo.ID, repo })
	reposToApps := make(map[string][]*domain.Application)
	for _, app := range applications {
		reposToApps[app.RepositoryID] = append(reposToApps[app.RepositoryID], app)
	}

	for _, repo := range repos {
		err = r.updateApps(ctx, repo, reposToApps[repo.ID])
		if err != nil {
			return errors.Wrap(err, "failed to update repo")
		}
	}
	return nil
}

func (r *repositoryFetcherService) resolveRefs(ctx context.Context, repo *domain.Repository) (refToCommit map[string]string, err error) {
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
		return nil, errors.Wrap(err, "failed to list remote refs")
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

func (r *repositoryFetcherService) updateApps(ctx context.Context, repo *domain.Repository, apps []*domain.Application) error {
	if len(apps) == 0 {
		return nil
	}

	refToCommit, err := r.resolveRefs(ctx, repo)
	if err != nil {
		return err
	}

	for _, app := range apps {
		commit, ok := refToCommit[app.BranchName]
		if !ok {
			// TODO: log error and present to user?
			log.Errorf("failed to get resolve ref %v for app %v", app.BranchName, app.ID)
			continue // fail-safe
		}
		err = r.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{WantCommit: optional.From(commit)})
		if err != nil {
			return errors.Wrap(err, "failed to update application")
		}
	}
	return nil
}

func (r *repositoryFetcherService) Stop(_ context.Context) error {
	r.closeOnce.Do(r.close)
	return nil
}
