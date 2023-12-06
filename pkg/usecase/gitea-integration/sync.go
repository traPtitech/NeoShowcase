package giteaintegration

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"

	"code.gitea.io/sdk/gitea"
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func listAllPages[T any](fn func(page, perPage int) ([]T, error)) ([]T, error) {
	items := make([]T, 0)
	for page := 1; ; page++ {
		const perPage = 50 // max per page
		pageItems, err := fn(page, perPage)
		if err != nil {
			return nil, err
		}
		items = append(items, pageItems...)
		if len(pageItems) < perPage {
			break
		}
	}
	return items, nil
}

func (i *Integration) sync(ctx context.Context) error {
	start := time.Now()
	log.Infof("Starting sync ...")
	defer func() {
		log.Infof("Sync finished in %v.", time.Since(start))
	}()

	// Sync users
	log.Infof("Retrieving users from Gitea ...")
	giteaUsers, err := listAllPages(func(page, perPage int) ([]*gitea.User, error) {
		users, _, err := i.c.AdminListUsers(gitea.AdminListUsersOptions{
			ListOptions: gitea.ListOptions{Page: page, PageSize: perPage},
		})
		return users, err
	})
	if err != nil {
		return errors.Wrap(err, "listing users")
	}
	userNames := ds.Map(giteaUsers, func(user *gitea.User) string { return user.UserName })

	log.Infof("Syncing %v users ...", len(userNames))
	users, err := i.userRepo.EnsureUsers(ctx, userNames)
	if err != nil {
		return err
	}
	usersMap := lo.SliceToMap(users, func(user *domain.User) (string, *domain.User) { return user.Name, user })
	log.Infof("Synced %v users.", len(userNames))

	// Get repositories
	log.Infof("Retrieving repositories from Gitea ...")
	repos, err := listAllPages(func(page, perPage int) ([]*gitea.Repository, error) {
		repos, _, err := i.c.SearchRepos(gitea.SearchRepoOptions{
			ListOptions: gitea.ListOptions{
				Page:     page,
				PageSize: perPage,
			},
			Sort:  "created",
			Order: "desc",
		})
		return repos, err
	})
	if err != nil {
		return errors.Wrap(err, "listing repositories")
	}

	// Sync repositories
	log.Infof("Syncing %v repositories ...", len(repos))
	var eg errgroup.Group
	eg.SetLimit(i.concurrency)
	for _, repo := range repos {
		repo := repo
		eg.Go(func() error {
			// Get users with write access
			members, _, err := i.c.GetAssignees(repo.Owner.UserName, repo.Name)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("syncing repository %v/%v: getting users", repo.Owner.UserName, repo.Name))
			}
			memberIDs := lo.Flatten(ds.Map(members, func(member *gitea.User) []string {
				user, ok := usersMap[member.UserName]
				if ok {
					return []string{user.ID}
				} else {
					log.Warnf("failed to find user %v", member.UserName)
					return nil
				}
			}))
			// Sync repository and apps
			return i.syncRepository(ctx, repo.Owner.UserName, memberIDs, repo)
		})
	}
	return eg.Wait()
}

func (i *Integration) syncRepository(ctx context.Context, username string, giteaOwnerIDs []string, giteaRepo *gitea.Repository) error {
	// NOTE: no transaction, creating repository is assumed rare
	repos, err := i.gitRepo.GetRepositories(ctx, domain.GetRepositoryCondition{URLs: optional.From([]string{giteaRepo.SSHURL})})
	if err != nil {
		return err
	}

	if len(repos) == 0 {
		// Does not exist, sync repository metadata
		repo := domain.NewRepository(
			fmt.Sprintf("%v/%v", username, giteaRepo.Name),
			giteaRepo.SSHURL,
			optional.From(domain.RepositoryAuth{Method: domain.RepositoryAuthMethodSSH}),
			giteaOwnerIDs,
		)
		log.Infof("New repository %v (id: %v)", repo.Name, repo.ID)
		return i.gitRepo.CreateRepository(ctx, repo)
	}

	// Already exists
	repo := repos[0]

	// Sync owners
	// Are all repository owners on Gitea also an owner of the repository on NeoShowcase?
	allOwnersAdded := lo.EveryBy(giteaOwnerIDs, func(ownerID string) bool { return slices.Contains(repo.OwnerIDs, ownerID) })
	if !allOwnersAdded {
		newOwners := ds.UniqMergeSlice(repo.OwnerIDs, giteaOwnerIDs)
		log.Infof("Syncing repository %v (id: %v) owners, %v users -> %v users", repo.Name, repo.ID, len(repo.OwnerIDs), len(newOwners))
		err = i.gitRepo.UpdateRepository(ctx, repo.ID, &domain.UpdateRepositoryArgs{OwnerIDs: optional.From(newOwners)})
		if err != nil {
			return err
		}
	}

	// Sync owners of generated applications
	apps, err := i.appRepo.GetApplications(ctx, domain.GetApplicationCondition{RepositoryID: optional.From(repo.ID)})
	if err != nil {
		return err
	}
	for _, app := range apps {
		err = i.syncApplication(ctx, app, giteaOwnerIDs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Integration) syncApplication(ctx context.Context, app *domain.Application, giteaOwnerIDs []string) error {
	// Are all repository owners on Gitea also an owner of generated application on NeoShowcase?
	allOwnersAdded := lo.EveryBy(giteaOwnerIDs, func(ownerID string) bool { return slices.Contains(app.OwnerIDs, ownerID) })
	if !allOwnersAdded {
		newOwners := ds.UniqMergeSlice(app.OwnerIDs, giteaOwnerIDs)
		log.Infof("Syncing application %v (id: %v) owners, %v users -> %v users", app.Name, app.ID, len(app.OwnerIDs), len(newOwners))
		err := i.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{OwnerIDs: optional.From(newOwners)})
		if err != nil {
			return err
		}
	}
	return nil
}
