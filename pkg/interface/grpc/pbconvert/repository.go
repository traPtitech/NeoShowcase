package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/mapper"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func FromPBRepositoryAuth(req *pb.CreateRepositoryAuth) optional.Of[domain.RepositoryAuth] {
	switch v := req.Auth.(type) {
	case *pb.CreateRepositoryAuth_None:
		return optional.Of[domain.RepositoryAuth]{}
	case *pb.CreateRepositoryAuth_Basic:
		return optional.From(domain.RepositoryAuth{
			Method:   domain.RepositoryAuthMethodBasic,
			Username: v.Basic.Username,
			Password: v.Basic.Password,
		})
	case *pb.CreateRepositoryAuth_Ssh:
		return optional.From(domain.RepositoryAuth{
			Method: domain.RepositoryAuthMethodSSH,
			SSHKey: v.Ssh.SshKey,
		})
	default:
		panic("unknown auth type")
	}
}

var RepoAuthMethodMapper = mapper.NewValueMapper(map[domain.RepositoryAuthMethod]pb.Repository_AuthMethod{
	domain.RepositoryAuthMethodBasic: pb.Repository_BASIC,
	domain.RepositoryAuthMethodSSH:   pb.Repository_SSH,
})

func ToPBRepository(repo *domain.Repository) *pb.Repository {
	ret := &pb.Repository{
		Id:   repo.ID,
		Name: repo.Name,
		Url:  repo.URL,
	}
	if repo.Auth.Valid {
		ret.AuthMethod = RepoAuthMethodMapper.IntoMust(repo.Auth.V.Method)
	}
	return ret
}
