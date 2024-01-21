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

Recommended dev environment is docker compose.
k3d environment is mainly for testing k8s features.

### docker backend (docker compose)

You will need `/Makefile` and `/compose.yaml` at the project root.

1. `make init`: Install / update development tools
2. `make up`: Spin up development environment
3. `make down`: Tear down development environment

Everything should automatically start after running `make up`.

- Dashboard: http://ns.local.trapti.tech/
- For more, type `make` to display all commands help

### k8s-backend (k3d)

You will need manifest files in `/.local-manifest` directory.

1. `make init` (at project root): Install / update development tools
2. `make up` (at `/.local-manifest`): Spin up development environment
3. `make down` (at `/.local-manifest`): Tear down development environment
   - The use of k3d (k3s in docker) allows ease cleanup.

Everything should automatically start after running `make up`.

For more, see [.local-manifest/README.md](../.local-manifest/README.md).

## Testing

Run tests with `make test`.

See [Makefile](../Makefile) for more.
