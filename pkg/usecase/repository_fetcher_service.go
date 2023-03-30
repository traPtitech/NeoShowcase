package usecase

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/friendsofgo/errors"

	"github.com/go-git/go-git/v5"
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

type RepositoryFetcherCacheDir string

type repositoryFetcherService struct {
	bus      domain.Bus
	appRepo  domain.ApplicationRepository
	gitRepo  domain.GitRepositoryRepository
	cacheDir string
	pubKey   *ssh.PublicKeys

	run       func()
	runOnce   sync.Once
	close     func()
	closeOnce sync.Once
}

func NewRepositoryFetcherService(
	bus domain.Bus,
	appRepo domain.ApplicationRepository,
	gitRepo domain.GitRepositoryRepository,
	cacheDir RepositoryFetcherCacheDir,
	pubKey *ssh.PublicKeys,
) (RepositoryFetcherService, error) {
	r := &repositoryFetcherService{
		bus:     bus,
		appRepo: appRepo,
		gitRepo: gitRepo,
		pubKey:  pubKey,
	}

	if cacheDir == "" {
		tmp, err := os.MkdirTemp("", "repo-fetcher-cache-")
		if err != nil {
			return nil, err
		}
		r.cacheDir = tmp
	} else {
		r.cacheDir = string(cacheDir)
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
		gitRepo, err := r.fetchRepository(ctx, repo)
		if err != nil {
			log.Errorf("failed to fetch repository: %+v", err)
			continue // fail-safe
		}
		for _, app := range reposToApps[repo.ID] {
			err := r.updateApplication(ctx, app, gitRepo)
			if err != nil {
				return errors.Wrap(err, "failed to update app")
			}
		}
	}
	return nil
}

func (r *repositoryFetcherService) fetchRepository(ctx context.Context, repo *domain.Repository) (*git.Repository, error) {
	auth, err := domain.GitAuthMethod(repo, r.pubKey)
	if err != nil {
		return nil, err
	}

	repoDir := filepath.Join(r.cacheDir, repo.ID)
	_, err = os.Stat(repoDir)
	if errors.Is(err, fs.ErrNotExist) {
		gitRepo, err := git.PlainCloneContext(ctx, repoDir, true, &git.CloneOptions{
			URL:        repo.URL,
			Auth:       auth,
			RemoteName: "origin",
		})
		if err != nil {
			return nil, err
		}
		return gitRepo, nil
	} else {
		gitRepo, err := git.PlainOpen(repoDir)
		if err != nil {
			return nil, err
		}
		err = gitRepo.FetchContext(ctx, &git.FetchOptions{
			Auth:       auth,
			RemoteName: "origin",
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return nil, err
		}
		return gitRepo, nil
	}
}

func (r *repositoryFetcherService) updateApplication(ctx context.Context, app *domain.Application, repo *git.Repository) error {
	branch, err := repo.Branch(app.BranchName)
	if err == git.ErrBranchNotFound {
		log.WithField("app", app.ID).Errorf("branch %s not found", app.BranchName)
		return nil // skip app update
	}
	if err != nil {
		return err
	}
	ref, err := repo.Reference(branch.Merge, true)
	if err != nil {
		return err
	}
	hash := ref.Hash().String()
	if app.WantCommit == hash {
		return nil
	}
	return r.appRepo.UpdateApplication(ctx, app.ID, domain.UpdateApplicationArgs{WantCommit: optional.From(hash)})
}

func (r *repositoryFetcherService) Stop(_ context.Context) error {
	r.closeOnce.Do(r.close)
	return nil
}
