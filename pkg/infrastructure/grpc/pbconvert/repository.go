package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

var RepoScopeMapper = mapper.MustNewValueMapper(map[apiserver.GetRepoScope]pb.GetRepositoriesRequest_Scope{
	apiserver.GetRepoScopeMine:   pb.GetRepositoriesRequest_MINE,
	apiserver.GetRepoScopePublic: pb.GetRepositoriesRequest_PUBLIC,
	apiserver.GetRepoScopeAll:    pb.GetRepositoriesRequest_ALL,
})

func FromPBRepositoryAuth(req *pb.CreateRepositoryAuth) optional.Of[apiserver.CreateRepositoryAuth] {
	switch v := req.Auth.(type) {
	case *pb.CreateRepositoryAuth_None:
		return optional.Of[apiserver.CreateRepositoryAuth]{}
	case *pb.CreateRepositoryAuth_Basic:
		return optional.From(apiserver.CreateRepositoryAuth{
			Method:   domain.RepositoryAuthMethodBasic,
			Username: v.Basic.Username,
			Password: v.Basic.Password,
		})
	case *pb.CreateRepositoryAuth_Ssh:
		return optional.From(apiserver.CreateRepositoryAuth{
			Method: domain.RepositoryAuthMethodSSH,
			KeyID:  v.Ssh.KeyId,
		})
	default:
		panic("unknown auth type")
	}
}

func FromPBUpdateRepositoryOwners(req *pb.UpdateRepositoryRequest_UpdateOwners) []string {
	return req.OwnerIds
}

var RepoAuthMethodMapper = mapper.MustNewValueMapper(map[domain.RepositoryAuthMethod]pb.Repository_AuthMethod{
	domain.RepositoryAuthMethodBasic: pb.Repository_BASIC,
	domain.RepositoryAuthMethodSSH:   pb.Repository_SSH,
})

func ToPBRepository(repo *domain.Repository) *pb.Repository {
	ret := &pb.Repository{
		Id:       repo.ID,
		Name:     repo.Name,
		Url:      repo.URL,
		OwnerIds: repo.OwnerIDs,
	}
	if repo.Auth.Valid {
		ret.AuthMethod = RepoAuthMethodMapper.IntoMust(repo.Auth.V.Method)
	}
	return ret
}
