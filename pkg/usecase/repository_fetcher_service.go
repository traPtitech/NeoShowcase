package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type RepositoryFetcherService interface {
	Run()
	Stop(ctx context.Context) error
}

type RepositoryFetcherCacheDir string

type repositoryFetcherService struct {
	bus     domain.Bus
	appRepo repository.ApplicationRepository

	cacheDir string

	run       func()
	runOnce   sync.Once
	close     func()
	closeOnce sync.Once
}

func NewRepositoryFetcherService(
	bus domain.Bus,
	appRepo repository.ApplicationRepository,
	cacheDir RepositoryFetcherCacheDir,
) (RepositoryFetcherService, error) {
	r := &repositoryFetcherService{
		bus:     bus,
		appRepo: appRepo,
	}

	if cacheDir == "" {
		tmp, err := os.MkdirTemp("", "repo-fetcher")
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
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	doSync := func() {
		start := time.Now()
		if err := r.fetchAll(); err != nil {
			log.WithError(err).Error("failed to fetch repositories")
			return
		}
		log.Infof("Fetched repositories in %v", time.Since(start))
		r.bus.Publish(event.FetcherFetchComplete, nil)
	}

	doSync()

	for {
		select {
		case <-ticker.C:
			doSync()
		case <-closer:
			return
		}
	}
}

func (r *repositoryFetcherService) fetchAll() error {
	ctx := context.Background()
	applications, err := r.appRepo.GetApplications(ctx)
	if err != nil {
		return err
	}
	repos := lo.SliceToMap(applications, func(app *domain.Application) (string, domain.Repository) { return app.Repository.ID, app.Repository })
	reposToApps := make(map[string][]*domain.Application)
	for _, app := range applications {
		reposToApps[app.Repository.ID] = append(reposToApps[app.Repository.ID], app)
	}

	for _, repo := range repos {
		gitRepo, err := r.fetchRepository(ctx, repo)
		if err != nil {
			log.WithError(err).
				WithField("repository", repo.URL).
				Error("failed to fetch repository")
			continue // fail-safe
		}
		for _, app := range reposToApps[repo.ID] {
			err := r.updateApplication(ctx, app, gitRepo)
			if err != nil {
				return fmt.Errorf("failed to update app: %w", err)
			}
		}
	}
	return nil
}

func (r *repositoryFetcherService) fetchRepository(ctx context.Context, repo domain.Repository) (*git.Repository, error) {
	repoDir := filepath.Join(r.cacheDir, repo.ID)
	_, err := os.Stat(repoDir)
	if os.IsNotExist(err) {
		gitRepo, err := git.PlainCloneContext(ctx, repoDir, true, &git.CloneOptions{
			URL:        repo.URL,
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
			RemoteName: "origin",
		})
		if err != nil {
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
	return r.appRepo.UpdateApplication(ctx, app.ID, repository.UpdateApplicationArgs{WantCommit: optional.From(hash)})
}

func (r *repositoryFetcherService) Stop(_ context.Context) error {
	r.closeOnce.Do(r.close)
	return nil
}