package main

import (
	"github.com/moby/buildkit/util/appdefaults"
	"github.com/spf13/viper"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/k8simpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/buildpack"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/log/loki"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/metrics/prometheus"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver/builtin"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver/caddy"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/webhook"
	"github.com/traPtitech/neoshowcase/pkg/usecase/healthcheck"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
)

type Config struct {
	PrivateKeyFile  string                   `mapstructure:"privateKeyFile" yaml:"privateKeyFile"`
	AdditionalLinks []*domain.AdditionalLink `mapstructure:"additionalLinks" yaml:"additionalLinks"`

	DB      repository.Config    `mapstructure:"db" yaml:"db"`
	Storage domain.StorageConfig `mapstructure:"storage" yaml:"storage"`
	Image   builder.ImageConfig  `mapstructure:"image" yaml:"image"`

	Components ComponentsConfig `mapstructure:"components" yaml:"components"`
}

type ComponentsConfig struct {
	AuthDev          AuthDevConfig          `mapstructure:"authDev" yaml:"authDev"`
	Builder          BuilderConfig          `mapstructure:"builder" yaml:"builder"`
	Controller       ControllerConfig       `mapstructure:"controller" yaml:"controller"`
	Gateway          GatewayConfig          `mapstructure:"gateway" yaml:"gateway"`
	GiteaIntegration GiteaIntegrationConfig `mapstructure:"giteaIntegration" yaml:"giteaIntegration"`
	SSGen            SSGenConfig            `mapstructure:"ssgen" yaml:"ssgen"`
}

type AuthDevConfig struct {
	Header string `mapstructure:"header" yaml:"header"`
	Port   int    `mapstructure:"port" yaml:"port"`
	User   string `mapstructure:"user" yaml:"user"`
}

type BuilderConfig struct {
	Buildkit struct {
		Address string `mapstructure:"address" yaml:"address"`
	} `mapstructure:"buildkit" yaml:"buildkit"`
	Buildpack   buildpack.Config                   `mapstructure:"buildpack" yaml:"buildpack"`
	Controller  grpc.ControllerServiceClientConfig `mapstructure:"controller" yaml:"controller"`
	Priority    int                                `mapstructure:"priority" yaml:"priority"`
	StepTimeout string                             `mapstructure:"stepTimeout" yaml:"stepTimeout"`
	Mock        bool                               `mapstructure:"mock" yaml:"mock"`
}

type ControllerConfig struct {
	Port        int                    `mapstructure:"port" yaml:"port"`
	TokenHeader string                 `mapstructure:"token" yaml:"port"`
	Token       string                 `mapstructure:"token" yaml:"port"`
	Mode        string                 `mapstructure:"mode" yaml:"mode"`
	Docker      dockerimpl.Config      `mapstructure:"docker" yaml:"docker"`
	K8s         k8simpl.Config         `mapstructure:"k8s" yaml:"k8s"`
	SSH         domain.SSHConfig       `mapstructure:"ssh" yaml:"ssh"`
	Webhook     webhook.ReceiverConfig `mapstructure:"webhook" yaml:"webhook"`
}

type GatewayConfig struct {
	Port          int                                `mapstructure:"port" yaml:"port"`
	AvatarBaseURL domain.AvatarBaseURL               `mapstructure:"avatarBaseURL" yaml:"avatarBaseURL"`
	AuthHeader    grpc.AuthHeader                    `mapstructure:"authHeader" yaml:"authHeader"`
	Controller    grpc.ControllerServiceClientConfig `mapstructure:"controller" yaml:"controller"`
	MariaDB       dbmanager.MariaDBConfig            `mapstructure:"mariadb" yaml:"mariadb"`
	MongoDB       dbmanager.MongoDBConfig            `mapstructure:"mongodb" yaml:"mongodb"`
	Log           struct {
		Type string      `mapstructure:"type" yaml:"type"`
		Loki loki.Config `mapstructure:"loki" yaml:"loki"`
	} `mapstructure:"log" yaml:"log"`
	Metrics struct {
		Type       string            `mapstructure:"type" yaml:"type"`
		Prometheus prometheus.Config `mapstructure:"prometheus" yaml:"prometheus"`
	}
}

type GiteaIntegrationConfig struct {
	URL             string                             `mapstructure:"url" yaml:"url"`
	Token           string                             `mapstructure:"token" yaml:"token"`
	IntervalSeconds int                                `mapstructure:"intervalSeconds" yaml:"intervalSeconds"`
	Concurrency     int                                `mapstructure:"concurrency" yaml:"concurrency"`
	Controller      grpc.ControllerServiceClientConfig `mapstructure:"controller" yaml:"controller"`
}

type SSGenConfig struct {
	ArtifactsRoot string           `mapstructure:"artifactsRoot" yaml:"artifactsRoot"`
	HealthPort    healthcheck.Port `mapstructure:"healthPort" yaml:"healthPort"`
	Server        struct {
		Type    string         `mapstructure:"type" yaml:"type"`
		BuiltIn builtin.Config `mapstructure:"builtIn" yaml:"builtIn"`
		Caddy   caddy.Config   `mapstructure:"caddy" yaml:"caddy"`
	} `mapstructure:"server" yaml:"server"`
	Controller grpc.ControllerServiceClientConfig `mapstructure:"controller" yaml:"controller"`
}

func init() {
	viper.SetDefault("privateKeyFile", "")
	viper.SetDefault("additionalLinks", nil)

	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", 3306)
	viper.SetDefault("db.username", "root")
	viper.SetDefault("db.password", "password")
	viper.SetDefault("db.database", "neoshowcase")
	viper.SetDefault("db.connection.maxOpen", 0)
	viper.SetDefault("db.connection.maxIdle", 2)
	viper.SetDefault("db.connection.lifetime", 0)

	viper.SetDefault("storage.type", "local")
	viper.SetDefault("storage.local.dir", "/neoshowcase")
	viper.SetDefault("storage.s3.bucket", "neoshowcase")
	viper.SetDefault("storage.s3.endpoint", "")
	viper.SetDefault("storage.s3.region", "")
	viper.SetDefault("storage.s3.accessKey", "")
	viper.SetDefault("storage.s3.accessSecret", "")
	viper.SetDefault("storage.swift.username", "")
	viper.SetDefault("storage.swift.apiKey", "")
	viper.SetDefault("storage.swift.tenantName", "")
	viper.SetDefault("storage.swift.tenantId", "")
	viper.SetDefault("storage.swift.container", "neoshowcase")
	viper.SetDefault("storage.swift.authUrl", "")

	viper.SetDefault("image.registry.scheme", "https")
	viper.SetDefault("image.registry.addr", "localhost")
	viper.SetDefault("image.registry.username", "")
	viper.SetDefault("image.registry.password", "")
	viper.SetDefault("image.namePrefix", "ns-apps/")
	viper.SetDefault("image.tmpNamePrefix", "ns-apps-tmp/")

	viper.SetDefault("components.authDev.header", "X-Showcase-User")
	viper.SetDefault("components.authDev.port", 4181)
	viper.SetDefault("components.authDev.user", "toki")

	viper.SetDefault("components.builder.buildkit.address", appdefaults.Address)

	viper.SetDefault("components.builder.buildpack.helper.address", "http://localhost:1235")
	viper.SetDefault("components.builder.buildpack.helper.listenPort", 1235)
	viper.SetDefault("components.builder.buildpack.remoteDir", "/workspace")
	viper.SetDefault("components.builder.buildpack.platformAPI", "0.11")

	viper.SetDefault("components.builder.controller.url", "http://ns-controller:10000")

	viper.SetDefault("components.builder.priority", 0)
	viper.SetDefault("components.builder.stepTimeout", "1h")

	viper.SetDefault("components.controller.port", 10000)
	viper.SetDefault("components.controller.tokenHeader", "X-NS-Controller-Token")
	viper.SetDefault("components.controller.token", "")
	viper.SetDefault("components.controller.mode", "docker")

	viper.SetDefault("components.controller.docker.confDir", "/opt/traefik/conf")
	viper.SetDefault("components.controller.docker.domains", nil)
	viper.SetDefault("components.controller.docker.ports", nil)

	viper.SetDefault("components.controller.docker.ss.url", "")
	viper.SetDefault("components.controller.docker.routing.type", "traefik")
	viper.SetDefault("components.controller.docker.routing.traefik.priorityOffset", 0)
	viper.SetDefault("components.controller.docker.tls.certResolver", "nsresolver")
	viper.SetDefault("components.controller.docker.tls.wildcard.domains", nil)

	viper.SetDefault("components.controller.docker.network", "neoshowcase_apps")
	viper.SetDefault("components.controller.docker.labels", nil)
	viper.SetDefault("components.controller.docker.resources.cpus", 1.6)
	viper.SetDefault("components.controller.docker.resources.memory", 1e9 /* 1GB */)
	viper.SetDefault("components.controller.docker.resources.memorySwap", -1 /* unlimited swap */)
	viper.SetDefault("components.controller.docker.resources.memoryReservation", 256*1e6 /* 256MB */)

	viper.SetDefault("components.controller.k8s.domains", nil)
	viper.SetDefault("components.controller.k8s.ports", nil)

	viper.SetDefault("components.controller.k8s.ss.namespace", "default")
	viper.SetDefault("components.controller.k8s.ss.kind", "Service")
	viper.SetDefault("components.controller.k8s.ss.name", "")
	viper.SetDefault("components.controller.k8s.ss.port", 80)
	viper.SetDefault("components.controller.k8s.ss.scheme", "http")

	viper.SetDefault("components.controller.k8s.routing.type", "traefik")
	viper.SetDefault("components.controller.k8s.routing.traefik.priorityOffset", 0)

	viper.SetDefault("components.controller.k8s.service.ipFamilies", nil)
	viper.SetDefault("components.controller.k8s.service.ipFamilyPolicy", "PreferDualStack")

	viper.SetDefault("components.controller.k8s.middleware.sablier.enable", true)
	viper.SetDefault("components.controller.k8s.middleware.sablier.url", "http://sablier:10000")
	viper.SetDefault("components.controller.k8s.middleware.sablier.sessionDuration", "1h")
	viper.SetDefault("components.controller.k8s.middleware.sablier.dynamic.theme", "ghost")
	viper.SetDefault("components.controller.k8s.middleware.sablier.blocking.timeout", "1m")

	viper.SetDefault("components.controller.k8s.tls.type", "traefik")
	viper.SetDefault("components.controller.k8s.tls.traefik.certResolver", "nsresolver")
	viper.SetDefault("components.controller.k8s.tls.traefik.wildcard.domains", nil)
	viper.SetDefault("components.controller.k8s.tls.certManager.issuer.name", "cert-issuer")
	viper.SetDefault("components.controller.k8s.tls.certManager.issuer.kind", "ClusterIssuer")
	viper.SetDefault("components.controller.k8s.tls.certManager.wildcard.domains", nil)

	viper.SetDefault("components.controller.k8s.namespace", "neoshowcase-apps")
	viper.SetDefault("components.controller.k8s.imagePullSecret", "")
	viper.SetDefault("components.controller.k8s.labels", nil)
	viper.SetDefault("components.controller.k8s.scheduling.nodeSelector", nil)
	viper.SetDefault("components.controller.k8s.scheduling.tolerations", nil)
	viper.SetDefault("components.controller.k8s.scheduling.forceHosts", nil)
	viper.SetDefault("components.controller.k8s.scheduling.spreadConstraints", nil)
	viper.SetDefault("components.controller.k8s.resources.requests.cpu", "")
	viper.SetDefault("components.controller.k8s.resources.requests.memory", "")
	viper.SetDefault("components.controller.k8s.resources.limits.cpu", "")
	viper.SetDefault("components.controller.k8s.resources.limits.memory", "")

	viper.SetDefault("components.controller.ssh.host", "localhost")
	viper.SetDefault("components.controller.ssh.port", 2201)

	viper.SetDefault("components.controller.webhook.basePath", "/api/webhook")
	viper.SetDefault("components.controller.webhook.port", 8080)

	viper.SetDefault("components.gateway.port", 8080)
	viper.SetDefault("components.gateway.avatarBaseURL", "https://q.trap.jp/api/v3/public/icon/")
	viper.SetDefault("components.gateway.authHeader", "X-Showcase-User")

	viper.SetDefault("components.gateway.controller.url", "http://ns-controller:10000")

	viper.SetDefault("components.gateway.mariadb.host", "mariadb")
	viper.SetDefault("components.gateway.mariadb.port", 3306)
	viper.SetDefault("components.gateway.mariadb.adminUser", "root")
	viper.SetDefault("components.gateway.mariadb.adminPassword", "password")

	viper.SetDefault("components.gateway.mongodb.host", "mongo")
	viper.SetDefault("components.gateway.mongodb.port", 27017)
	viper.SetDefault("components.gateway.mongodb.adminUser", "root")
	viper.SetDefault("components.gateway.mongodb.adminPassword", "password")

	viper.SetDefault("components.gateway.log.type", "loki")
	viper.SetDefault("components.gateway.log.loki.endpoint", "http://loki:3100")
	viper.SetDefault("components.gateway.log.loki.queryTemplate", loki.DefaultQueryTemplate())

	viper.SetDefault("components.gateway.metrics.type", "prometheus")
	viper.SetDefault("components.gateway.metrics.endpoint", "http://prometheus:9090")
	viper.SetDefault("components.gateway.metric.queries", prometheus.DefaultQueriesConfig())

	viper.SetDefault("components.giteaIntegration.url", "https://git.trap.jp")
	viper.SetDefault("components.giteaIntegration.token", "")
	viper.SetDefault("components.giteaIntegration.intervalSeconds", 86400)
	viper.SetDefault("components.giteaIntegration.concurrency", 10)
	viper.SetDefault("components.giteaIntegration.controller.url", "http://ns-controller:10000")

	viper.SetDefault("components.ssgen.artifactsRoot", "/srv/artifacts")
	viper.SetDefault("components.ssgen.healthPort", 8081)

	viper.SetDefault("components.ssgen.server.type", "builtIn")
	viper.SetDefault("components.ssgen.server.builtIn.port", 8080)
	viper.SetDefault("components.ssgen.server.caddy.adminAPI", "http://localhost:2019")
	viper.SetDefault("components.ssgen.server.caddy.docsRoot", "/srv/artifacts")

	viper.SetDefault("components.ssgen.controller.url", "http://ns-controller:10000")
}
