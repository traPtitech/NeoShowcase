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
- buildkitd
  - Used by "ns-builder", this is where applications are actually built.
  - Buildpack builds run in containers that "ns-builder" creates for each build through the Docker socket mounted in it.
- mysql, adminer (or any other db admin UI)
  - Used by NeoShowcase itself
- mysql, mongo, adminer (or any other db admin UI)
  - Used by applications
- grafana, loki, promtail, victoria-metrics (or prometheus), cadvisor
  - Used for displaying application metrics and logs

"ns-builder" requires `components.builder.buildpack.runtime: docker` and the Docker socket (`/var/run/docker.sock`) mounted in it.
It labels build containers with its hostname and removes the ones left by a builder with the same hostname when it starts,
so give it a fixed hostname (`hostname:` in compose).

## Using k8s

NeoShowcase is NOT built against some specific cloud vendor, it is a cloud-agnostic application; it uses traefik reverse-proxy for both Ingress Controller and for routing components / deployed applications.

Use of [k3s](https://k3s.io/) is recommended whenever possible as it is well tested against.
Using other distributions "may" very well work, but it is not tested.

See [../.local-manifest](../.local-manifest) for a complete list of require components to deploy NeoShowcase.
Required components are (almost) the same with docker-backed deployment.

Our production manifest is also available at [traPtitech/manifest](https://github.com/traPtitech/manifest).

"ns-builder" runs each buildpack build in a pod in the namespace set by `components.builder.buildpack.k8s.namespace`.
Dedicate this namespace to build pods and bind the builder's Role there (see [../charts/neoshowcase/templates/builder](../charts/neoshowcase/templates/builder)), not in the namespace of NeoShowcase:
whoever can create pods in a namespace can read every secret in it through the pod spec.
Set `components.builder.buildpack.runtime: k8s`.
`components.builder.buildpack.k8s.resources` applies to each container of build pods; a value left out sets no request or limit.
`components.builder.buildpack.k8s.nodeSelector` (a list of `key` / `value`) and `tolerations` choose the nodes build pods run on.
The Helm chart sets them from `builder.nodeSelector` and `builder.tolerations`.
A build fails as soon as the scheduler reports its pod unschedulable, and the build log shows the reason that the scheduler gives.

## Upgrading from versions with the buildpack sidecar

Earlier versions ran buildpack builds in a "buildpack" sidecar of "ns-builder" (the `buildpack-helper` command).
The sidecar and `components.builder.buildpack.helper` / `remoteDir` are no longer used. To upgrade:

- Set `components.builder.buildpack.runtime` to `docker` or `k8s`. The builder does not start without it.
- With docker, mount the Docker socket in "ns-builder", give it a fixed hostname, and remove the "buildpack" service.
- With k8s, create the namespace for build pods and the builder's ServiceAccount, Role and RoleBinding, as the Helm chart does.
- With k8s, builds no longer run on the node of "ns-builder". If "ns-builder" is placed on particular nodes (such as build nodes with a taint, or nodes of an architecture),
  set the same placement in `components.builder.buildpack.k8s.nodeSelector` and `tolerations`.
