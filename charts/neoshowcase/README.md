# neoshowcase

![Version: 1.9.4](https://img.shields.io/badge/Version-1.9.4-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.9.4](https://img.shields.io/badge/AppVersion-1.9.4-informational?style=flat-square)

NeoShowcase is a PaaS application for Digital Creators Club traP.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| app | object | `{"imagePullSecret":"","labels":[],"namespace":"ns-apps","nodeSelector":{},"resources":{},"routing":{"traefik":{"priorityOffset":0},"type":"traefik"},"service":{"ipFamilies":[],"ipFamilyPolicy":""},"tolerations":[],"topologySpreadConstraints":[]}` | app defines user app pod configurations. |
| auth | object | `{"avatarBaseURL":"https://q.trap.jp/api/v3/public/icon/","header":"X-Forwarded-User"}` | auth defines user authentication. |
| auth.avatarBaseURL | string | `"https://q.trap.jp/api/v3/public/icon/"` | avatarBaseURL is used to display user icons in the dashboard. |
| auth.header | string | `"X-Forwarded-User"` | header defines the name of the header auth. |
| builder | object | `{"buildkit":{"buildkitd.toml":"","image":"moby/buildkit:rootless","resources":{}},"buildpack":{"image":"paketobuildpacks/builder:full","resources":{}},"nodeSelector":{},"replicas":1,"resources":{},"tolerations":[],"topologySpreadConstraints":[]}` | ns-builder component |
| common | object | `{"additionalLinks":[{"name":"Wiki","url":"https://wiki.trap.jp/services/NeoShowcase"},{"name":"DB Admin","url":"https://adminer.ns.trap.jp/"}],"additionalVolumeMounts":[],"additionalVolumes":[],"db":{"database":"neoshowcase","host":"mariadb.db.svc.cluster.local","port":3306,"username":"root"},"image":{"namePrefix":"ns-apps/","registry":{"addr":"registry.ns.trap.jp","scheme":"https","username":"robot$neoshowcase"},"tmpNamePrefix":"ns-apps-tmp/"},"logLevel":"info","storage":{"local":{"dir":"/data"},"s3":{"bucket":"neoshowcase","region":"ap-northeast-1"},"type":"s3"}}` | common defines various settings for use by all NeoShowcase components. |
| common.additionalLinks | list | `[{"name":"Wiki","url":"https://wiki.trap.jp/services/NeoShowcase"},{"name":"DB Admin","url":"https://adminer.ns.trap.jp/"}]` | Links to be displayed in the dashboard. |
| common.additionalVolumes | list | `[]` | Additional mounts to NeoShowcase component containers |
| common.db | object | `{"database":"neoshowcase","host":"mariadb.db.svc.cluster.local","port":3306,"username":"root"}` | db is used by NeoShowcase system components |
| common.image | object | `{"namePrefix":"ns-apps/","registry":{"addr":"registry.ns.trap.jp","scheme":"https","username":"robot$neoshowcase"},"tmpNamePrefix":"ns-apps-tmp/"}` | image is used to store user app images. |
| common.logLevel | string | `"info"` | Available log levels: trace, debug, info, warn, or error |
| common.storage | object | `{"local":{"dir":"/data"},"s3":{"bucket":"neoshowcase","region":"ap-northeast-1"},"type":"s3"}` | storage is used by NeoShowcase system components |
| controller | object | `{"nodeSelector":{},"replicas":1,"resources":{},"ssh":{"host":"ns.trap.jp","port":2201},"tolerations":[],"topologySpreadConstraints":[]}` | ns-controller component |
| dashboard | object | `{"nodeSelector":{},"replicas":1,"resources":{},"tolerations":[],"topologySpreadConstraints":[]}` | ns-dashboard component |
| domains | list | `[]` | domains define available domains to be used by user apps. For more, see pkg/infrastructure/k8simpl/config.go. |
| gateway | object | `{"nodeSelector":{},"replicas":1,"resources":{},"tolerations":[],"topologySpreadConstraints":[]}` | ns-gateway component |
| giteaIntegration | object | `{"enabled":false,"nodeSelector":{},"resources":{},"tolerations":[],"topologySpreadConstraints":[],"url":"https://git.trap.jp"}` | ns-gitea-integration component |
| global.appVersionOverride | string | `""` | If specified, overrides the chart app version. Used by NeoShowcase image tags. |
| global.image | object | `{"prefix":"ns-","repository":"ghcr.io/traptitech/"}` | image defines NeoShowcase image config. |
| ingressRoute | object | `{"enabled":true,"entrypoints":["web"],"host":"ns.trap.jp","middlewares":[],"tls":{"secretName":""}}` | ingressRoute renders IngressRoute resource, if enabled. |
| known_hosts | object | `{"additionalContent":""}` | known_hosts is mounted into builder, controller, and gateway to clone user repositories. |
| observability | object | `{"log":{"loki":{"endpoint":"http://loki.monitor.svc.cluster.local:3100","queryTemplate":"{namespace=\"ns-apps\",pod=\"nsapp-{{ .App.ID }}-0\"}"},"type":"victorialogs","victorialogs":{"endpoint":"http://vl-victoria-logs-single-server.victoria-logs.svc.cluster.local:9428","queryTemplate":"{namespace=\"ns-apps\",pod=\"nsapp-{{ .App.ID }}-0\"}"}},"metrics":{"prometheus":{"endpoint":"http://victoria-metrics.monitor.svc.cluster.local:8428","queries":[{"name":"CPU","template":"rate(container_cpu_user_seconds_total{namespace=\"ns-apps\", pod=\"nsapp-{{ .App.ID }}-0\", container=\"app\"}[5m]) + rate(container_cpu_system_seconds_total{namespace=\"ns-apps\", pod=\"nsapp-{{ .App.ID }}-0\", container=\"app\"}[5m])"},{"name":"Memory","template":"container_memory_usage_bytes{namespace=\"ns-apps\", pod=\"nsapp-{{ .App.ID }}-0\", container=\"app\"} + container_memory_swap{namespace=\"ns-apps\", pod=\"nsapp-{{ .App.ID }}-0\", container=\"app\"}"}]},"type":"prometheus"}}` | observability (o11y) defines user apps' o11y configuration to be viewed from the dashboard. |
| ports | list | `[]` | ports define available port-forward ports to be used by user apps. For more, see pkg/infrastructure/k8simpl/config.go. |
| sablier | object | `{"blocking":{"timeout":"1m"},"dynamic":{"theme":"neoshowcase"},"enabled":true,"nodeSelector":{},"resources":{},"sessionDuration":"1h","tolerations":[],"topologySpreadConstraints":[]}` | sablier component starts user pods on demand. |
| secret | object | `{"keys":{"existingName":"ns-keys","keyName":"id_ed25519"},"ns":{"existingName":"ns"}}` | secret defines secret names to be used by NeoShowcase components. |
| secret.keys | object | `{"existingName":"ns-keys","keyName":"id_ed25519"}` | Keys are used by gateway and controller to clone user repositories. The corresponding public key is intended to be set to an admin deploy-key of an external Gitea instance. |
| secret.keys.keyName | string | `"id_ed25519"` | Only ed25519 type is supported for now. |
| ssgen | object | `{"caddy":{"image":"caddy:2-alpine","resources":{}},"nodeSelector":{},"pvc":{"storage":"1Gi","storageClassName":""},"replicas":2,"resources":{},"tolerations":[],"topologySpreadConstraints":[]}` | ns-ssgen component |
| tls | object | `{"certManager":{"issuer":{"kind":"ClusterIssuer","name":"cluster-issuer"},"wildcard":{"domains":[]}},"type":"cert-manager"}` | tls defines tls setting for user app ingress. For more, see pkg/infrastructure/k8simpl/config.go. |
| userMariaDB | object | `{"adminUser":"root","host":"mariadb.db.svc.cluster.local","port":3306}` | userMariaDB is used by user apps. |
| userMongoDB | object | `{"adminUser":"root","host":"mongo.db.svc.cluster.local","port":27017}` | userMongoDB is used by user apps. |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.14.2](https://github.com/norwoodj/helm-docs/releases/v1.14.2)
