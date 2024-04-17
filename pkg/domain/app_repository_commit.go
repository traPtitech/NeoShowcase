package domain

import (
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
)

type CommitAuthorSignature struct {
	Name  string
	Email string
	Date  time.Time
}

type RepositoryCommit struct {
	Hash      string
	Author    CommitAuthorSignature
	Committer CommitAuthorSignature
	Message   string
	Error     bool
}

func ToRepositoryCommit(c *object.Commit) *RepositoryCommit {
	return &RepositoryCommit{
		Hash: c.Hash.String(),
		Author: CommitAuthorSignature{
			Name:  c.Author.Name,
			Email: c.Author.Email,
			Date:  c.Author.When,
		},
		Committer: CommitAuthorSignature{
			Name:  c.Committer.Name,
			Email: c.Committer.Email,
			Date:  c.Committer.When,
		},
		Message: c.Message,
		Error:   false,
	}
}

func ToErroredRepositoryCommit(hash string) *RepositoryCommit {
	return &RepositoryCommit{
		Hash:  hash,
		Error: true,
	}
}
