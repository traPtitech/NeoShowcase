package pbconvert

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func ToPBBuilderSystemInfo(info *domain.BuilderSystemInfo) *pb.BuilderSystemInfo {
	return &pb.BuilderSystemInfo{
		PrivateKey: info.SSHKey,
		ImageConfig: &pb.ImageConfig{
			Registry: &pb.ImageConfig_RegistryConfig{
				Scheme:   info.ImageConfig.Registry.Scheme,
				Addr:     info.ImageConfig.Registry.Addr,
				Username: info.ImageConfig.Registry.Username,
				Password: info.ImageConfig.Registry.Password,
			},
			NamePrefix:    info.ImageConfig.NamePrefix,
			TmpNamePrefix: info.ImageConfig.TmpNamePrefix,
		},
	}
}

func FromPBBuilderSystemInfo(info *pb.BuilderSystemInfo) *domain.BuilderSystemInfo {
	return &domain.BuilderSystemInfo{
		SSHKey: info.PrivateKey,
		ImageConfig: builder.ImageConfig{
			Registry: builder.RegistryConfig{
				Scheme:   info.ImageConfig.Registry.Scheme,
				Addr:     info.ImageConfig.Registry.Addr,
				Username: info.ImageConfig.Registry.Username,
				Password: info.ImageConfig.Registry.Password,
			},
			NamePrefix:    info.ImageConfig.NamePrefix,
			TmpNamePrefix: info.ImageConfig.TmpNamePrefix,
		},
	}
}

func ToPBRepositoryPrivate(repo *domain.Repository) *pb.RepositoryPrivate {
	return &pb.RepositoryPrivate{
		Repo:     ToPBRepository(repo),
		Username: repo.Auth.V.Username,
		Password: repo.Auth.V.Password,
		SshKey:   repo.Auth.V.SSHKey,
	}
}

func FromPBRepositoryPrivate(repo *pb.RepositoryPrivate) *domain.Repository {
	ret := FromPBRepository(repo.Repo)
	ret.Auth.V.Username = repo.Username
	ret.Auth.V.Password = repo.Password
	ret.Auth.V.SSHKey = repo.SshKey
	return ret
}

func ToPBStartBuildRequest(req *domain.StartBuildRequest) *pb.StartBuildRequest {
	return &pb.StartBuildRequest{
		Repo: ToPBRepositoryPrivate(req.Repo),
		App:  ToPBApplication(req.App, nil /* builder does not use this field */),
		AppEnvs: &pb.ApplicationEnvVars{
			Variables: ds.Map(req.Envs, ToPBEnvironment),
		},
		Build: ToPBBuild(req.Build),
	}
}

func FromPBStartBuildRequest(req *pb.StartBuildRequest) *domain.StartBuildRequest {
	return &domain.StartBuildRequest{
		Repo:  FromPBRepositoryPrivate(req.Repo),
		App:   FromPBApplication(req.App),
		Envs:  ds.Map(req.AppEnvs.Variables, FromPBEnvironment),
		Build: FromPBBuild(req.Build),
	}
}
