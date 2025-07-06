package cdservice

import (
	"context"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/discovery"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type AppDeployHelper struct {
	cluster     *discovery.Cluster
	backend     domain.Backend
	appRepo     domain.ApplicationRepository
	buildRepo   domain.BuildRepository
	envRepo     domain.EnvironmentRepository
	websiteRepo domain.WebsiteRepository
	ssgen       domain.ControllerSSGenService
	image       builder.ImageConfig
}

func NewAppDeployHelper(
	cluster *discovery.Cluster,
	backend domain.Backend,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	envRepo domain.EnvironmentRepository,
	websiteRepo domain.WebsiteRepository,
	ssgen domain.ControllerSSGenService,
	imageConfig builder.ImageConfig,
) *AppDeployHelper {
	return &AppDeployHelper{
		cluster:     cluster,
		backend:     backend,
		appRepo:     appRepo,
		buildRepo:   buildRepo,
		envRepo:     envRepo,
		websiteRepo: websiteRepo,
		ssgen:       ssgen,
		image:       imageConfig,
	}
}

func (s *AppDeployHelper) _getEnv(ctx context.Context, apps []*domain.Application) (map[string]map[string]string, error) {
	appIDs := ds.Map(apps, func(app *domain.Application) string { return app.ID })
	envs, err := s.envRepo.GetEnv(ctx, domain.GetEnvCondition{ApplicationIDIn: optional.From(appIDs)})
	if err != nil {
		return nil, err
	}
	ret := make(map[string]map[string]string, len(appIDs))
	for _, env := range envs {
		if _, ok := ret[env.ApplicationID]; !ok {
			ret[env.ApplicationID] = make(map[string]string)
		}
		ret[env.ApplicationID][env.Key] = env.Value
	}
	return ret, nil
}

func (s *AppDeployHelper) _runtimeDesiredStates(ctx context.Context) ([]*domain.RuntimeDesiredState, error) {
	// Get all 'running' state applications
	apps, err := s.appRepo.GetApplications(ctx, domain.GetApplicationCondition{
		DeployType: optional.From(domain.DeployTypeRuntime),
		Running:    optional.From(true),
	})
	if err != nil {
		return nil, err
	}
	// Shard by app ID
	apps = lo.Filter(apps, func(app *domain.Application, _ int) bool {
		return s.cluster.Assigned(app.ID)
	})

	syncableApps := lo.Filter(apps, func(app *domain.Application, _ int) bool { return app.CurrentBuild != "" })
	envs, err := s._getEnv(ctx, syncableApps)
	if err != nil {
		return nil, err
	}
	desiredStates := ds.Map(syncableApps, func(app *domain.Application) *domain.RuntimeDesiredState {
		return &domain.RuntimeDesiredState{
			App:       app,
			ImageName: s.image.ImageName(app.ID),
			ImageTag:  app.CurrentBuild,
			Envs:      envs[app.ID],
		}
	})
	return desiredStates, nil
}

func (s *AppDeployHelper) _collectSharedResources(ctx context.Context) (*domain.DesiredStateLeader, error) {
	var st domain.DesiredStateLeader
	websites, err := s.websiteRepo.GetWebsites(ctx)
	if err != nil {
		return nil, err
	}
	tlsTargetDomains := make(map[string]struct{})
	for _, website := range websites {
		host, ok := s.backend.TLSTargetDomain(website)
		if ok {
			tlsTargetDomains[host] = struct{}{}
		}
	}
	st.TLSTargetDomains = lo.Keys(tlsTargetDomains)
	return &st, nil
}

func (s *AppDeployHelper) synchronize(ctx context.Context) error {
	// Synchronize sharded resource
	var st domain.DesiredState
	var err error
	st.Runtime, err = s._runtimeDesiredStates(ctx)
	if err != nil {
		return err
	}
	st.StaticSites, err = domain.GetActiveStaticSites(ctx, s.cluster, s.appRepo, s.buildRepo)
	if err != nil {
		return err
	}
	log.Infof("[shard %d/%d] %v runtime, %v static sites active", s.cluster.Me(), s.cluster.Size(), len(st.Runtime), len(st.StaticSites))

	s.ssgen.BroadcastSSGen(&pb.SSGenRequest{Type: pb.SSGenRequest_RELOAD})
	err = s.backend.Synchronize(ctx, &st)
	if err != nil {
		return err
	}

	// Only let leader synchronize shared resources (certificates)
	if s.cluster.IsLeader() {
		st, err := s._collectSharedResources(ctx)
		if err != nil {
			return err
		}
		log.Infof("[shard leader] %v tls targets active", len(st.TLSTargetDomains))
		err = s.backend.SynchronizeShared(ctx, st)
		if err != nil {
			return err
		}
	}

	return nil
}
