package giteaintegration

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"code.gitea.io/sdk/gitea"
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"

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
	slog.InfoContext(ctx, "Starting sync")
	defer func() {
		slog.InfoContext(ctx, "Sync finished", "duration", time.Since(start))
	}()

	// Sync users
	slog.InfoContext(ctx, "Retrieving users from Gitea")
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

	slog.InfoContext(ctx, "Syncing users", "count", len(userNames))
	users, err := i.userRepo.EnsureUsers(ctx, userNames)
	if err != nil {
		return err
	}
	usersMap := lo.SliceToMap(users, func(user *domain.User) (string, *domain.User) { return user.Name, user })
	slog.InfoContext(ctx, "Synced users", "count", len(userNames))

	// Get repositories
	slog.InfoContext(ctx, "Retrieving repositories from Gitea")
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
	slog.InfoContext(ctx, "Syncing repositories", "count", len(repos))
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
					slog.Warn("failed to find user", "user_name", member.UserName)
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
		slog.InfoContext(ctx, "New repository found", "name", repo.Name, "id", repo.ID)
		return i.gitRepo.CreateRepository(ctx, repo)
	}

	// Already exists
	repo := repos[0]

	// Sync owners
	// Are all repository owners on Gitea also an owner of the repository on NeoShowcase?
	allOwnersAdded := lo.EveryBy(giteaOwnerIDs, func(ownerID string) bool { return slices.Contains(repo.OwnerIDs, ownerID) })
	if !allOwnersAdded {
		newOwners := ds.UniqMergeSlice(repo.OwnerIDs, giteaOwnerIDs)
		slog.InfoContext(ctx, "Syncing repository owners", "repo_name", repo.Name, "repo_id", repo.ID, "old_count", len(repo.OwnerIDs), "new_count", len(newOwners))
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
		slog.InfoContext(ctx, "Syncing application owners", "app_name", app.Name, "app_id", app.ID, "old_count", len(app.OwnerIDs), "new_count", len(newOwners))
		err := i.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{OwnerIDs: optional.From(newOwners)})
		if err != nil {
			return err
		}
	}
	return nil
}
