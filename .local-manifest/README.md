# manifest

Manifest files required to deploy NeoShowcase locally using k8s backend

## Usage

tl;dr:
Run `mise run up` / `mise run down`

### Spin up

1. `mise run k3d:up` to spin up the k3d cluster
2. `mise run import` to build and import NeoShowcase images
3. `mise run apply` to apply the manifest files

### Managing

- To tail cluster events (pulling image, creating container etc.), `kubectl get events --watch`
- To tail specific pod / container logs, `kubectl logs -f <pod-name> -n <namespace>`
  - For more, visit http://grafana.local.trapti.tech/ and see centralized Loki logs
- Go to http://localhost:8080/ to view traefik dashboard

#### Optional recommended tools

- k9s: https://k9scli.io/
- Lens: https://k8slens.dev/

### Spin down

1. `mise run k3d:down` to take down the k3d cluster completely
