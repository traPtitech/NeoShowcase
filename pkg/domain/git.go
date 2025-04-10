package domain

import (
	"context"
)

type GitService interface {
	// ResolveRefs resolves refs in the repository and returns a map of reference name to commit hash.
	//
	// e.g. refs/heads/master -> 0123456789abcdef0123456789abcdef01234567
	ResolveRefs(ctx context.Context, repo *Repository) (refToCommit map[string]string, err error)
	CloneRepository(ctx context.Context, dir string, repo *Repository, commitHash string) error
	CreateBareRepository(dir string, repo *Repository) (GitRepository, error)
}

type GitRepository interface {
	Fetch(ctx context.Context, hashes []string) error
	GetCommit(hash string) (*RepositoryCommit, error)
}
