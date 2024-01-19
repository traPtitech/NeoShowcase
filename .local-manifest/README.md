# manifest

Manifest files required to deploy NeoShowcase locally using k8s backend

## Usage

tl;dr:
Run `make up` / `make down`

### Spin up

1. `make k3d-up` to spin up the k3d cluster
2. `make import` to build and import the builder image
3. `make apply` to spin up the builder

### Managing

- To tail cluster events (pulling image, creating container etc.), `make events`
- To tail specific pod / container logs, `make logs NAMESPACE=ns-system APP=ns-controller`
  - For more, visit http://grafana.local.trapti.tech/ and see centralized Loki logs
- Go to http://localhost:8080/ to view traefik dashboard

#### Optional recommended tools

- k9s: https://k9scli.io/
- Lens: https://k8slens.dev/

### Spin down

1. `make k3d-down` to take down the k3d cluster completely
