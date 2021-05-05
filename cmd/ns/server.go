package main

import (
	"context"
	"database/sql"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/appmanager"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/k8simpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Server struct {
	webserver  *web.Server
	db         *sql.DB
	bus        eventbus.Bus
	appmanager appmanager.Manager

	builderConn *grpc.ClientConn
	ssgenConn   *grpc.ClientConn
	backend     backend.Backend

	k8sCSet *kubernetes.Clientset
}

func NewServer(c Config, webserver *web.Server, db *sql.DB, bus eventbus.Bus) (*Server, error) {
	s := &Server{
		webserver: webserver,
		db:        db,
		bus:       bus,
	}

	switch c.GetMode() {
	case ModeDocker:
		// Builderに接続
		builderConn, err := c.Builder.Connect()
		if err != nil {
			return nil, fmt.Errorf("failed to connect builder service: %w", err)
		}
		s.builderConn = builderConn

		// SSGenに接続
		ssgenConn, err := c.SSGen.Connect()
		if err != nil {
			return nil, fmt.Errorf("failed to connect ssgen service: %w", err)
		}
		s.ssgenConn = ssgenConn

		// Dockerデーモンに接続 (DooD)
		dc, err := docker.NewClientFromEnv()
		if err != nil {
			return nil, fmt.Errorf("failed to create docker client: %w", err)
		}

		// コンテナマネージャー生成
		connM, err := dockerimpl.NewDockerBackend(dc, s.bus, "/opt/traefik/conf")
		if err != nil {
			return nil, fmt.Errorf("failed to init container manager: %w", err)
		}
		s.backend = connM

		// appmanger生成
		am, err := appmanager.NewManager(appmanager.Config{
			DB:              db,
			Hub:             s.bus,
			Builder:         pb.NewBuilderServiceClient(builderConn),
			SS:              pb.NewStaticSiteServiceClient(ssgenConn),
			Backend:         connM,
			ImageRegistry:   c.Image.Registry,
			ImageNamePrefix: c.Image.NamePrefix,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to init app manager: %w", err)
		}
		s.appmanager = am

	case ModeK8s:
		// Builderに接続
		builderConn, err := c.Builder.Connect()
		if err != nil {
			return nil, fmt.Errorf("failed to connect builder service: %w", err)
		}
		s.builderConn = builderConn

		// SSGenに接続
		ssgenConn, err := c.SSGen.Connect()
		if err != nil {
			return nil, fmt.Errorf("failed to connect ssgen service: %w", err)
		}
		s.ssgenConn = ssgenConn

		// k8s接続
		kubeconf, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
		}
		clientset, err := kubernetes.NewForConfig(kubeconf)
		if err != nil {
			return nil, fmt.Errorf("failed to create clientset: %w", err)
		}
		s.k8sCSet = clientset

		// コンテナマネージャー生成
		connM, err := k8simpl.NewK8SBackend(s.bus, clientset)
		if err != nil {
			return nil, fmt.Errorf("failed to init container manager: %w", err)
		}
		s.backend = connM

		// appmanger生成
		am, err := appmanager.NewManager(appmanager.Config{
			DB:              db,
			Hub:             s.bus,
			Builder:         pb.NewBuilderServiceClient(builderConn),
			SS:              pb.NewStaticSiteServiceClient(ssgenConn),
			Backend:         connM,
			ImageRegistry:   c.Image.Registry,
			ImageNamePrefix: c.Image.NamePrefix,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to init app manager: %w", err)
		}
		s.appmanager = am

	default:
		log.Fatalf("unknown mode: %s", c.Mode)
	}

	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	return s.webserver.Start(ctx)
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.webserver.Shutdown(ctx)
	})
	eg.Go(func() error {
		return s.db.Close()
	})
	eg.Go(func() error {
		return s.builderConn.Close()
	})
	eg.Go(func() error {
		return s.ssgenConn.Close()
	})
	eg.Go(func() error {
		return s.backend.Dispose(ctx)
	})
	eg.Go(func() error {
		return s.appmanager.Shutdown(ctx)
	})
	eg.Go(func() error {
		return s.bus.Close(ctx)
	})

	return eg.Wait()
}
