package repoconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
)

func ToDomainRepositoryCommit(c *models.RepositoryCommit) *domain.RepositoryCommit {
	return &domain.RepositoryCommit{
		Hash: c.Hash,
		Author: domain.CommitAuthorSignature{
			Name:  c.AuthorName,
			Email: c.AuthorEmail,
			Date:  c.AuthorDate,
		},
		Committer: domain.CommitAuthorSignature{
			Name:  c.CommitterName,
			Email: c.CommitterEmail,
			Date:  c.CommitterDate,
		},
		Message: c.Message,
		Error:   c.Error,
	}
}

func FromDomainRepositoryCommit(c *domain.RepositoryCommit) *models.RepositoryCommit {
	return &models.RepositoryCommit{
		Hash:           c.Hash,
		AuthorName:     c.Author.Name,
		AuthorEmail:    c.Author.Email,
		AuthorDate:     c.Author.Date,
		CommitterName:  c.Committer.Name,
		CommitterEmail: c.Committer.Email,
		CommitterDate:  c.Committer.Date,
		Message:        c.Message,
		Error:          c.Error,
	}
}
