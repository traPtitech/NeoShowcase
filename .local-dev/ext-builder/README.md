# ext-builder

Join external builder instance from outside the production (cluster) deployment

## Usage

tl;dr:
Do "First time setup", and run `make up` / `make down`

### First time setup

1. Set up config files
   - `./config/controller-url.txt` for controller URL
   - `./config/controller-token.txt` for controller token
   - `./config/known_hosts` for known hosts configuration
     - `cp ~/.ssh/known_hosts ./config` should be enough
2. Prepare local forward if controller port is unreachable from the internet
   - e.g. `ssh -L 0.0.0.0:10000:10.43.193.98:10000 c1-203`

### Importing images

To import ns images to the k3d cluster, `make import`

### Spin up

1. `make k3d-up`
2. `make apply`

### Workaround for local registry

ref: https://zenn.dev/toshikish/articles/7f555dbf1b4b7d

Edit and add `rewrite name registry.local host.k3d.internal` inside the `.:53 {}` section:
`kubectl edit cm -n kube-system coredns`

Restart CoreDNS afterwards:
`kubectl restart deployment/coredns -n kube-system`

### Managing

- To scale the number of builder instances, `make scale REPLICAS=3`
- To tail builder logs, `make logs`

### Spin down

1. `make k3d-down`
