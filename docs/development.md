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

## Changing Database Schema

To change the database schema, do the following:

1. Change schema in `./migrations/schema.sql`. This definition file is the source of truth for all generated tables / codes.
2. Run `make migrate` to apply schema changes to local db container in an idempotent manner.
3. Run `make gen` (or individually, `make gen-go && make gen-db-docs`) to update generated codes and docs via SQLBoiler and tbls.
4. Write your code.
   - Don't forget to modify fields in `./pkg/domain` structs, `./pkg/infrastructure/repository/repoconvert` functions, etc.

## Changing API

To change the API schema, do the following:

1. Change schema in `./api/proto/neoshowcase/protobuf/*.proto` files. These files are the source of truth for all generated codes.
2. Run `make gen` (or individually, `make gen-proto`) to generate both server (Go) and client (TypeScript) codes.
3. Write your code.
   - Don't forget to modify fields in `./pkg/infrastructure/grpc/pbconvert` etc.
