# Deployment

How to run a production environment

NeoShowcase is backed by either docker or kubernetes.
Both environment works, but docker backend is limited to only one node (machine) and has a limited scalability.
Using kubernetes backend is recommended for large-scale set-up.

NeoShowcase heavily makes use of traefik reverse-proxy.
Combining with other reverse-proxies like nginx "may" work, but is not tested.
Use and expose traefik when possible.

Authentication is done by proxy authentication.
By default, it uses [traefik-forward-auth](https://github.com/traPtitech/traefik-forward-auth).

## Using docker

See [../compose.yaml](../compose.yaml) for required components.

- ns-gateway, dashboard, ns-controller, ns-builder, ns-ssgen, ns-migrate
  - NeoShowcase itself
- ns-gitea-integration (optional)
  - Allows auto-sync of Gitea repositories, if you own an instance.
- ns-auth
  - Proxy authentication with traefik
- static-server
  - Used by "ns-ssgen", serves applications' static files.
- registry, registry-frontend
  - Hosting your own docker registry is recommended.
  - Public registries "may" work, but is not tested.
- buildpack, buildkitd
  - Used by "ns-builder", this is where applications are actually built.
- mysql, adminer
  - Used by NeoShowcase itself
- mysql, mongo, adminer
  - Used by applications
- grafana, loki, promtail, victoria-metrics (or prometheus), cadvisor
  - Used for displaying application metrics and logs

## Using k8s

NeoShowcase is NOT built against some specific cloud vendor, it is a cloud-agnostic application; it uses traefik reverse-proxy for both Ingress Controller and for routing components / deployed applications.

Use of [k3s](https://k3s.io/) is recommended whenever possible as it is well tested against.
Using other distributions "may" very well work, but it is not tested.

See [../.local-dev/manifest](../.local-dev/manifest) for a complete list of require components to deploy NeoShowcase.
Required components are (almost) the same with docker-backed deployment.

Our production manifest is also available at [traPtitech/manifest](https://github.com/traPtitech/manifest).
