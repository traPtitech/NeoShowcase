package observability

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type MetricsServerConfig struct {
	Port int `mapstructure:"port" yaml:"port"`
}

type MetricsServer struct {
	config MetricsServerConfig
	server *http.Server
	mux    *http.ServeMux
}

func NewMetricsServer(config MetricsServerConfig) *MetricsServer {
	mux := http.NewServeMux()
	return &MetricsServer{
		config: config,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.Port),
			Handler: mux,
		},
		mux: mux,
	}
}

func (s *MetricsServer) Start() error {
	s.mux.Handle("/metrics", promhttp.Handler())
	return s.server.ListenAndServe()
}

func (s *MetricsServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

type ControllerMetrics struct {
	buildCounter   *prometheus.CounterVec
	deployDuration *prometheus.HistogramVec
}

func NewControllerMetrics() *ControllerMetrics {
	return &ControllerMetrics{
		buildCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: "neoshowcase",
			Subsystem: "controller",
			Name:      "builds_total",
		}, []string{"result", "build_type"}),
		deployDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "neoshowcase",
			Subsystem: "controller",
			Name:      "deploy_duration_seconds",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 8), // 1s ~ 128s
		}, []string{}),
	}
}

func (s *ControllerMetrics) IncrementBuild(status domain.BuildStatus, buildType domain.BuildType) {
	s.buildCounter.WithLabelValues(status.String(), buildType.String()).Inc()
}

func (s *ControllerMetrics) ObserveDeployDuration(d time.Duration) {
	s.deployDuration.WithLabelValues().Observe(d.Seconds())
}
