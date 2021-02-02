package staticsitegen

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"github.com/traPtitech/neoshowcase/pkg/staticsitegen/api"
	"github.com/traPtitech/neoshowcase/pkg/staticsitegen/generator"
	"github.com/traPtitech/neoshowcase/pkg/storage"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
)

type Service struct {
	server  *grpc.Server
	engine  generator.Engine
	db      *sql.DB
	storage storage.Storage

	config Config
	api.UnimplementedStaticSiteGenServiceServer
}

func New(c Config) (*Service, error) {
	s := &Service{
		server: grpc.NewServer(),
		config: c,
	}
	api.RegisterStaticSiteGenServiceServer(s.server, s)
	reflection.Register(s.server)

	// Storageに接続
	storage, err := c.Storage.InitStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	s.storage = storage

	// Engine初期化
	engine, err := c.GetEngine()
	if err != nil {
		return nil, err
	}
	if err := engine.Init(storage); err != nil {
		return nil, fmt.Errorf("failed to init engine: %w", err)
	}
	s.engine = engine

	// DBに接続
	db, err := c.DB.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}
	s.db = db

	return s, nil
}

func (s *Service) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPC.GetPort()))
	if err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}
	return s.server.Serve(listener)
}

func (s *Service) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		s.server.GracefulStop()
		return nil
	})
	eg.Go(func() error {
		return s.db.Close()
	})

	return eg.Wait()
}

func (s *Service) Reload(ctx context.Context, _ *api.ReloadRequest) (*api.ReloadResponse, error) {
	err := s.reload(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &api.ReloadResponse{}, nil
}

func (s *Service) reload(ctx context.Context) error {
	envs, err := models.Environments(
		models.EnvironmentWhere.BuildType.EQ(models.EnvironmentsBuildTypeStatic),
		qm.Load(qm.Rels(models.EnvironmentRels.Website, models.WebsiteRels.Build, models.BuildLogRels.Artifact)),
	).All(ctx, s.db)
	if err != nil {
		return err
	}

	var data []*generator.Site
	for _, env := range envs {
		if env.R.Website != nil {
			website := env.R.Website
			if website.R.Build != nil {
				build := website.R.Build
				if build.R.Artifact != nil {
					artifact := build.R.Artifact
					data = append(data, &generator.Site{
						ID:            website.ID,
						FQDN:          website.FQDN,
						ArtifactID:    artifact.ID,
						ApplicationID: env.ApplicationID,
					})
				}
			}
		}
	}

	return s.engine.Reconcile(data)
}
