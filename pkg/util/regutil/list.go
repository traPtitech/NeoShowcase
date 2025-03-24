package regutil

import (
	"context"

	"github.com/regclient/regclient"
	"github.com/regclient/regclient/scheme"
	"github.com/regclient/regclient/types/ref"
)

func RepoList(ctx context.Context, r *regclient.RegClient, regHost string) ([]string, error) {
	const limit = 100
	var repos []string
	for {
		opts := []scheme.RepoOpts{scheme.WithRepoLimit(limit)}
		if len(repos) > 0 {
			opts = append(opts, scheme.WithRepoLast(repos[len(repos)-1]))
		}
		repoList, err := r.RepoList(ctx, regHost, opts...)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repoList.Repositories...)
		if len(repoList.Repositories) < limit {
			return repos, nil
		}
	}
}

func TagList(ctx context.Context, r *regclient.RegClient, imageName string) ([]string, error) {
	const limit = 100
	var tags []string
	for {
		opts := []scheme.TagOpts{scheme.WithTagLimit(limit)}
		if len(tags) > 0 {
			opts = append(opts, scheme.WithTagLast(tags[len(tags)-1]))
		}
		repoRef, err := ref.New(imageName)
		if err != nil {
			return nil, err
		}
		tagList, err := r.TagList(ctx, repoRef, opts...)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tagList.Tags...)
		if len(tagList.Tags) < limit {
			return tags, nil
		}
	}
}
