# loki

Generated using loki chart

https://grafana.com/docs/loki/latest/installation/helm/install-monolithic/

In root directory:
```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
rm -rf loki/templates
helm template --values loki/values.yaml --output-dir . --include-crds --namespace loki loki grafana/loki
```

## patches

- `templates/configmap.yaml`
  - `retention_deletes_enabled: true`
  - `retention_period: 672h`
