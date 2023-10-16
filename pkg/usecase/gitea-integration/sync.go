package giteaintegration

import (
	"context"
	"fmt"
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

func listAllPages[T any](fn func(page, perPage int) ([]T, error), listInterval time.Duration) ([]T, error) {
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
		time.Sleep(listInterval)
	}
	return items, nil
}

func (i *Integration) sync(ctx context.Context) {
	start := time.Now()
	err := i._sync(ctx)
	if err != nil {
		log.Errorf("failed to sync with gitea: %+v", err)
		return
	}
	log.Infof("Synced with gitea in %v", time.Since(start))
}

func (i *Integration) _sync(ctx context.Context) error {
	// Sync users
	giteaUsers, err := listAllPages(func(page, perPage int) ([]*gitea.User, error) {
		users, _, err := i.c.AdminListUsers(gitea.AdminListUsersOptions{
			ListOptions: gitea.ListOptions{Page: page, PageSize: perPage},
		})
		return users, err
	}, i.listInterval)
	if err != nil {
		return errors.Wrap(err, "listing users")
	}
	userNames := ds.Map(giteaUsers, func(user *gitea.User) string { return user.UserName })
	log.Infof("Syncing %v users", len(userNames))
	users, err := i.userRepo.EnsureUsers(ctx, userNames)
	if err != nil {
		return err
	}
	usersMap := lo.SliceToMap(users, func(user *domain.User) (string, *domain.User) { return user.Name, user })

	// Sync repositories for each user
	for _, user := range users {
		giteaRepos, err := listAllPages(func(page, perPage int) ([]*gitea.Repository, error) {
			repos, _, err := i.c.ListUserRepos(user.Name, gitea.ListReposOptions{
				ListOptions: gitea.ListOptions{Page: page, PageSize: perPage},
			})
			return repos, err
		}, i.listInterval)
		if err != nil {
			return errors.Wrap(err, "listing user repositories")
		}

		for _, giteaRepo := range giteaRepos {
			err = i.syncRepository(ctx, user.Name, []string{user.ID}, giteaRepo)
			if err != nil {
				return errors.Wrap(err, "syncing user repository")
			}
		}

		time.Sleep(i.listInterval)
	}

	// Sync repositories for each org
	giteaOrgs, err := listAllPages(func(page, perPage int) ([]*gitea.Organization, error) {
		orgs, _, err := i.c.AdminListOrgs(gitea.AdminListOrgsOptions{
			ListOptions: gitea.ListOptions{Page: page, PageSize: perPage},
		})
		return orgs, err
	}, i.listInterval)
	if err != nil {
		return errors.Wrap(err, "listing organizations")
	}
	for _, giteaOrg := range giteaOrgs {
		giteaRepos, err := listAllPages(func(page, perPage int) ([]*gitea.Repository, error) {
			repos, _, err := i.c.ListOrgRepos(giteaOrg.UserName, gitea.ListOrgReposOptions{
				ListOptions: gitea.ListOptions{Page: page, PageSize: perPage},
			})
			return repos, err
		}, i.listInterval)
		if err != nil {
			return errors.Wrap(err, "listing org repositories")
		}

		giteaOrgMembers, err := listAllPages(func(page, perPage int) ([]*gitea.User, error) {
			members, _, err := i.c.ListOrgMembership(giteaOrg.UserName, gitea.ListOrgMembershipOption{
				ListOptions: gitea.ListOptions{Page: page, PageSize: perPage},
			})
			return members, err
		}, i.listInterval)
		if err != nil {
			return errors.Wrap(err, "listing org members")
		}
		memberIDs := lo.Flatten(ds.Map(giteaOrgMembers, func(member *gitea.User) []string {
			user, ok := usersMap[member.UserName]
			if ok {
				return []string{user.ID}
			} else {
				log.Warnf("failed to find user %v", member.UserName)
				return nil
			}
		}))

		for _, giteaRepo := range giteaRepos {
			err = i.syncRepository(ctx, giteaOrg.UserName, memberIDs, giteaRepo)
		}

		time.Sleep(i.listInterval)
	}

	return nil
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
		log.Infof("Syncing repository %v -> id: %v", repo.Name, repo.ID)
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
		log.Infof("Syncing application %v (id: %v) owners, %v users -> %v users", app.Name, app.ID, len(app.OwnerIDs), newOwners)
		err := i.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{OwnerIDs: optional.From(newOwners)})
		if err != nil {
			return err
		}
	}
	return nil
}
