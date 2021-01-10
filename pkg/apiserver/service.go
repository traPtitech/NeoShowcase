package apiserver

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/leandro-lugaresi/hub"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/apiserver/httpserver"
	"github.com/traPtitech/neoshowcase/pkg/appmanager"
	builderApi "github.com/traPtitech/neoshowcase/pkg/builder/api"
	"github.com/traPtitech/neoshowcase/pkg/container"
	"github.com/traPtitech/neoshowcase/pkg/container/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/container/k8simpl"
	ssgenApi "github.com/traPtitech/neoshowcase/pkg/staticsitegen/api"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Service struct {
	server     *httpserver.Server
	db         *sql.DB
	bus        *hub.Hub
	appmanager appmanager.Manager

	builderConn      *grpc.ClientConn
	ssgenConn        *grpc.ClientConn
	containerManager container.Manager

	k8sCSet *kubernetes.Clientset

	config Config
}

func New(c Config) (*Service, error) {
	s := &Service{
		config: c,
		bus:    hub.New(),
	}

	// DBに接続
	db, err := c.DB.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}
	s.db = db

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
		am, err := appmanager.NewManager(
			db,
			s.bus,
			builderApi.NewBuilderServiceClient(builderConn),
			ssgenApi.NewStaticSiteGenServiceClient(ssgenConn),
			connM,
		)
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
		am, err := appmanager.NewManager(
			db,
			s.bus,
			builderApi.NewBuilderServiceClient(builderConn),
			ssgenApi.NewStaticSiteGenServiceClient(ssgenConn),
			connM,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to init app manager: %w", err)
		}
		s.appmanager = am

	default:
		log.Fatalf("unknown mode: %s", c.Mode)
	}

	// HTTP APIサーバー生成
	s.server = httpserver.New(httpserver.Config{
		Debug:      c.HTTP.Debug,
		Port:       c.HTTP.Port,
		Bus:        s.bus,
		AppManager: s.appmanager,
	})

	return s, nil
}

func (s *Service) Start(ctx context.Context) error {
	return s.server.Start()
}

func (s *Service) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.server.Shutdown(ctx)
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

func (s *Service) builder() builderApi.BuilderServiceClient {
	return builderApi.NewBuilderServiceClient(s.builderConn)
}

func (s *Service) ssgen() ssgenApi.StaticSiteGenServiceClient {
	return ssgenApi.NewStaticSiteGenServiceClient(s.ssgenConn)
}
