package domain

import "github.com/traPtitech/neoshowcase/pkg/domain/builder"

type Application struct {
	ID         string
	Repository Repository
}

type Branch struct {
	ID            string
	ApplicationID string
	BranchName    string
	BuildType     builder.BuildType
}

// type Environment struct {
// 	ID       string
// 	BranchID string
// 	Key      string
// 	Value    string
// }
type Repository struct {
	ID        string
	RemoteURL string
	Provider  Provider
}

type Provider struct {
	ID     string
	Name   string
	Secret []byte
}
