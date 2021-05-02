package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/leandro-lugaresi/hub"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/appmanager"
	builderApi "github.com/traPtitech/neoshowcase/pkg/builder/api"
	"github.com/traPtitech/neoshowcase/pkg/container"
	"github.com/traPtitech/neoshowcase/pkg/container/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/container/k8simpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/web"
	ssgenApi "github.com/traPtitech/neoshowcase/pkg/staticsitegen/api"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Server struct {
	webserver  *web.Server
	db         *sql.DB
	bus        *hub.Hub
	appmanager appmanager.Manager

	builderConn      *grpc.ClientConn
	ssgenConn        *grpc.ClientConn
	containerManager container.Manager

	k8sCSet *kubernetes.Clientset
}

func NewServer(c Config, webserver *web.Server, db *sql.DB, bus *hub.Hub) (*Server, error) {
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

		// コンテナマネージャー生成
		connM, err := dockerimpl.NewManager(s.bus)
		if err != nil {
			return nil, fmt.Errorf("failed to init container manager: %w", err)
		}
		s.containerManager = connM

		// appmanger生成
		am, err := appmanager.NewManager(appmanager.Config{
			DB:              db,
			Hub:             s.bus,
			Builder:         builderApi.NewBuilderServiceClient(builderConn),
			SSGen:           ssgenApi.NewStaticSiteGenServiceClient(ssgenConn),
			CM:              connM,
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
		connM, err := k8simpl.NewManager(s.bus, clientset)
		if err != nil {
			return nil, fmt.Errorf("failed to init container manager: %w", err)
		}
		s.containerManager = connM

		// appmanger生成
		am, err := appmanager.NewManager(appmanager.Config{
			DB:              db,
			Hub:             s.bus,
			Builder:         builderApi.NewBuilderServiceClient(builderConn),
			SSGen:           ssgenApi.NewStaticSiteGenServiceClient(ssgenConn),
			CM:              connM,
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
		return s.containerManager.Dispose(ctx)
	})
	eg.Go(func() error {
		return s.appmanager.Shutdown(ctx)
	})

	return eg.Wait()
}
