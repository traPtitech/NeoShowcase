package ssgen

import (
	"compress/gzip"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util/discovery"
	"github.com/traPtitech/neoshowcase/pkg/util/retry"
	"github.com/traPtitech/neoshowcase/pkg/util/scutil"
	"github.com/traPtitech/neoshowcase/pkg/util/tarfs"
)

type SiteReloadTarget struct {
	ApplicationID string
	BuildID       string
}

type GeneratorService interface {
	Start(ctx context.Context) error
	Healthy() bool
	Shutdown(ctx context.Context) error
}

type generatorService struct {
	cluster   *discovery.Cluster
	client    domain.ControllerSSGenServiceClient
	appRepo   domain.ApplicationRepository
	buildRepo domain.BuildRepository
	storage   domain.Storage
	engine    domain.StaticServer
	docsRoot  string

	cancel   func()
	reloaded atomic.Bool
	reloader *scutil.Coalescer
}

func NewGeneratorService(
	client domain.ControllerSSGenServiceClient,
	appRepo domain.ApplicationRepository,
	buildRepo domain.BuildRepository,
	storage domain.Storage,
	engine domain.StaticServer,
	path domain.StaticServerDocumentRootPath,
) GeneratorService {
	g := &generatorService{
		cluster:   discovery.NewCluster(discovery.NewSingleDiscoverer("127.0.0.1")), // No sharding for now
		client:    client,
		appRepo:   appRepo,
		buildRepo: buildRepo,
		storage:   storage,
		engine:    engine,
		docsRoot:  string(path),
	}
	g.reloader = scutil.NewCoalescer(g._reload)
	return g
}

func (s *generatorService) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	go func() {
		err := s.cluster.Start(ctx)
		if err != nil {
			slog.Error("failed to start cluster", "error", err)
		}
	}()
	go retry.Do(ctx, func(ctx context.Context) error {
		return s.client.ConnectSSGen(ctx, s.onRequest)
	}, "connect to controller")
	go func() {
		for range 300 {
			s.reload()
			if s.reloaded.Load() {
				break
			}
			<-time.After(1 * time.Second)
		}
	}()

	return nil
}

func (s *generatorService) Healthy() bool {
	return s.reloaded.Load()
}

func (s *generatorService) Shutdown(_ context.Context) error {
	s.cancel()
	return nil
}

func (s *generatorService) onRequest(req *pb.SSGenRequest) {
	switch req.Type {
	case pb.SSGenRequest_RELOAD:
		go s.reload()
	}
}

func (s *generatorService) reload() {
	err := s.reloader.Do(context.Background())
	if err != nil {
		slog.Error("failed to reload static server", "error", err)
	}
}

func (s *generatorService) _reload(ctx context.Context) error {
	start := time.Now()
	// Calculate active sites
	sites, err := domain.GetActiveStaticSites(ctx, s.cluster, s.appRepo, s.buildRepo)
	if err != nil {
		return err
	}
	// Sync artifacts on disk (download)
	err = s.syncArtifacts(sites)
	if err != nil {
		return err
	}
	// Reconcile server config
	err = s.engine.Reconcile(sites)
	if err != nil {
		return err
	}
	s.reloaded.Store(true)
	slog.InfoContext(ctx, "reloaded static server", "duration", time.Since(start), "sites_count", len(sites))
	return nil
}

func (s *generatorService) syncArtifacts(sites []*domain.StaticSite) error {
	entries, err := os.ReadDir(s.docsRoot)
	if err != nil {
		return errors.Wrap(err, "failed to read docs root")
	}
	artifactsOnDisk := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			artifactsOnDisk[e.Name()] = struct{}{}
		}
	}

	wantArtifacts := lo.SliceToMap(sites, func(site *domain.StaticSite) (string, struct{}) { return site.ArtifactID, struct{}{} })

	// Download artifacts to disk
	for artifactID := range wantArtifacts {
		if _, ok := artifactsOnDisk[artifactID]; ok {
			continue
		}
		err = s.extractArtifact(artifactID)
		if err != nil {
			return err
		}
	}

	// Prune old artifacts on disk
	for artifactID := range artifactsOnDisk {
		if _, ok := wantArtifacts[artifactID]; ok {
			continue
		}
		artifactDir := filepath.Join(s.docsRoot, artifactID)
		err = os.RemoveAll(artifactDir)
		if err != nil {
			return errors.Wrap(err, "failed to delete unused artifact directory")
		}
	}

	return nil
}

func (s *generatorService) extractArtifact(artifactID string) error {
	destDir := filepath.Join(s.docsRoot, artifactID)
	r, err := domain.GetArtifact(s.storage, artifactID)
	if err != nil {
		return errors.Wrap(err, "getting artifact")
	}
	defer r.Close()
	tarReader, err := gzip.NewReader(r)
	if err != nil {
		return errors.Wrap(err, "preparing gzip reader")
	}
	defer tarReader.Close()
	err = tarfs.Extract(tarReader, destDir)
	if err != nil {
		return errors.Wrap(err, "failed to extract artifact tar")
	}
	return nil
}
