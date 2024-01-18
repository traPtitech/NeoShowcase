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

1. `make k3d-up` to spin up the k3d cluster
2. `make import` to build and import the builder image
3. `make apply` to spin up the builder

### Managing

- To scale the number of builder instances, `make scale REPLICAS=3`
- To tail cluster events (pulling image, creating container etc.), `make events`
- To tail builder pod / container logs, `make logs`

### Spin down

1. `make k3d-down` to take down the k3d cluster completely
