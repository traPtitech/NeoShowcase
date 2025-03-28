package domain

import (
	"time"
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

func ToErroredRepositoryCommit(hash string) *RepositoryCommit {
	return &RepositoryCommit{
		Hash:  hash,
		Error: true,
	}
}
