# Development

## Workaround Notes

### Local (docker)

Add following to your `/etc/hosts` before executing `make up`
(workaround to issue #493)

```
127.0.0.1 registry.local
```

### Accessing wildcard local domains

`*.local.trapti.tech`, `local.trapti.tech` which points to `127.0.0.1` are used during development.

Accessing wildcard domains when running under WSL environment may require configuring Windows `C:\Windows\System32\drivers\etc\hosts` file like following.

```plaintext
127.0.0.1       ns.local.trapti.tech
::1             ns.local.trapti.tech    localhost
```

## Environments

### docker (main dev environment)

- `make`: Display Makefile help
- `make up`: Spin up local docker development environment
- `make down`: Tear down local docker development environment

Everything should automatically start after running `make up`.

- Dashboard: http://ns.local.trapti.tech/
- Gateway debug: `make ns-evans`

### k3s

Manifest files and instructions for debugging k8s backend implementation using k3s are available at [.local-manifest](../.local-manifest).

## Testing

### dind for docker test

dind (Docker in Docker) allows separation of docker environment from the host.

Run docker backend implementation tests as follows:

1. `make dind-up`
2. `make docker-test`
3. `make dind-down`

### k3d for k8s test

k3d (k3s in docker) allows separation of k3s environment from the host.
k3s nodes will be available as docker containers.

Run k8s backend implementation tests as follows:

1. `make k3d-up`
2. `make k8s-test`
3. `make k3d-down`
