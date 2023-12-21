# ext-builder

Join external builder instance from outside the production (cluster) deployment

## Usage

1. Set config `./config.yaml` accordingly
   - Fetch configuration from production if necessary
2. Set token (`NS_COMPONENTS_CONTROLLER_TOKEN`) in `.env` file
3. Prepare local forward if controller port is unreachable from the internet
   - `ssh -L 0.0.0.0:10000:10.43.193.98:10000 c1-203`
4. `docker compose pull`
5. `docker compose up -d`
