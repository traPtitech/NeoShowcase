# Development

## Prerequisites

This project uses mise to manage development tools and tasks.

1. Install mise: If you don't have it yet, follow the [Official Getting Started Guide](https://mise.jdx.dev/getting-started.html).
2. Setup Tools: Run the following command at the project root to install all required dependencies:
    ```
    mise install
    ```

## Workaround Notes

### Local (docker)

Add following to your `/etc/hosts` before executing `mise run up`
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

You will need `/compose.yaml` and `/mise.toml` at the project root.

1. `mise run up`: Spin up development environment
2. `mise run down`: Tear down development environment

Everything should automatically start after running `mise run up`.

- Dashboard: http://ns.local.trapti.tech/
- For more, run `mise tasks` to display all available commands

If you use Docker Desktop for Windows and WSL2...

- (recommended) install and use Docker on WSL2 instead of Docker Desktop
- (workaround) run `chmod o+rw /var/run/docker.sock`

### k8s-backend (k3d)

You will need manifest files in `/.local-manifest` directory.

1. `mise run up` (at `/.local-manifest`): Spin up development environment
2. `mise run down` (at `/.local-manifest`): Tear down development environment
   - The use of k3d (k3s in docker) allows ease cleanup.

Everything should automatically start after running `mise run up`.

For more, see [.local-manifest/README.md](../.local-manifest/README.md).

## Testing

Run tests with `mise run test`.

See [mise.toml](../mise.toml) for more available tasks.

## Changing Database Schema

To change the database schema, do the following:

1. Change schema in `./migrations/schema.sql`. This definition file is the source of truth for all generated tables / codes.
2. Run `mise run migrate` to apply schema changes to local db container in an idempotent manner.
3. Run `mise run gen` (or individually, `mise run gen:go && mise run gen:db-docs`) to update generated codes and docs via SQLBoiler and tbls.
4. Write your code.
   - Don't forget to modify fields in `./pkg/domain` structs, `./pkg/infrastructure/repository/repoconvert` functions, etc.

## Changing API

To change the API schema, do the following:

1. Change schema in `./api/proto/neoshowcase/protobuf/*.proto` files. These files are the source of truth for all generated codes.
2. Run `mise run gen` (or individually, `mise run gen:proto`) to generate both server (Go) and client (TypeScript) codes.
3. Write your code.
   - Don't forget to modify fields in `./pkg/infrastructure/grpc/pbconvert` etc.

## Adding Internal Component

To add a new internal component, we are using [github.com/google/wire](https://github.com/google/wire).

Example: A new repository in `./pkg/infrastructure/repository`, a new use-case service in `./pkg/usecase` etc.

1. Write a new component.
2. Add its constructor method (`New...()`) to `./cmd/providers.go`.
   - Reference the component from needed component. Example: add the component as a member in `Server` struct in `./cmd/controller/server.go`. See `./cmd/wire.go` to see how each component references multiple internal components.
3. Add config and its default to `./cmd/config.go`, if necessary.
4. Run `mise run gen` (or individually, `mise run gen:go`) to generate DI (dependency injection) codes.
