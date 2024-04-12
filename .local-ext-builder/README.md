# ext-builder

Join external builder instance from outside the production (cluster) deployment

## Usage

tl;dr:
Do "First time setup", and run `make up` / `make down`

### First time setup

1. Make config files and write secret contents inside them

| Path                                     | Content                                         | Example                                        |
|------------------------------------------|-------------------------------------------------|------------------------------------------------|
| `./manifest/config/controller-url.txt`   | Controller URL (as seen from local builder pod) | `http://host.k3d.internal:10000`               |
| `./manifest/config/controller-token.txt` | Controller Token                                | `abc...xyz` (The same one with the controller) |
| `./manifest/config/known_hosts`          | known_hosts file                                | Contents of `~/.ssh/known_hosts`               |

The controller token is given as [`components.controller.token`](https://github.com/traPtitech/NeoShowcase/blob/6456cc97c5e890440bd283542a73520beb17787c/cmd/config.go#L62) (or `NS_COMPONENTS_CONTROLLER_TOKEN` as environment variable) in config.
This secret token should match when a builder connects to a controller.

For example, in our production environment, the token resides [here](https://github.com/traPtitech/manifest/blob/20d2573c38a9c51727f0c9e3b558a6b08bd30da3/ns-system/secrets/ns.yaml#L15) in the secret.
The secret token is given to the controller at [here](https://github.com/traPtitech/manifest/blob/20d2573c38a9c51727f0c9e3b558a6b08bd30da3/ns-system/components/controller-stateful-set.yaml#L78-L82),
and given to the builder at [here](https://github.com/traPtitech/manifest/blob/20d2573c38a9c51727f0c9e3b558a6b08bd30da3/ns-system/components/builder-deployment.yaml#L129-L134).

2. Prepare local forward as necessary, if controller port is unreachable from the internet
   - example: `ssh -L 0.0.0.0:10000:10.43.193.98:10000 c1-203`

### Spin up

1. `make k3d-up` to spin up the k3d cluster
2. `make import` to build and import the builder image
3. `make apply` to spin up the builder

### Managing

- To scale the number of builder instances, `make scale REPLICAS=3`
- To tail cluster events (pulling image, creating container etc.), `make events`
- To tail builder pod / container logs, `make logs`

### Spin down

1. `make k3d-down` to take down the k3d cluster completely
