# Architecture

## Overview

NeoShowcase is a PaaS (**Platform as a Service**) that traP members (users) can use.

Here, PaaS refers to a service that runs applications on behalf of users just with the source code of the application and a minimal configuration.
Users manage the source code of their applications on Git and push it to services like GitHub or Gitea.
NeoShowcase automatically detects updates and builds applications according to user settings.

## The Reconciliation Design

One major design principle of NeoShowcase is **reconciliation loop** (also known as reconciliation pattern, synchronization, eventual consistency, etc.).
This is heavily influenced by the processing method of Kubernetes controllers (in fact, "ns-controller" component is a kind of Kubernetes controller).

Let's compare this to event(-only)-driven design, where internal processes are *only* triggered by user actions or external events, and perform some processing which may fire more events.
When events are lost during communication or when processing fails, we need rigorous error handling and retry logic.

Whereas the reconciliation loop design periodically monitors the system's state and performs **reconciliation** which tries to bring the current state to the desired state.

While the reconciliation loop design tends to have slightly higher computational overhead, it doesn't need complex retry logic.

### Examples

Each component has its own reconciliation loop.
Each component starts processing on a timer or when receiving an event from another component.

- controller/repository_fetcher
  - Updates applications' build target commit hashes.
    - Retrieves a list of repositories and applications from the database, fetches the latest commit hash corresponding to the git ref specified by each application, and stores that value in the database.
    - This is influenced by the way ArgoCD updates. Argo CD polls Git repositories every three minutes to detect changes to the manifests. https://argo-cd.readthedocs.io/en/stable/operator-manual/webhook/
- controller/continuous_deployment/build_registerer
  - Registers builds as "queued."
    - Lists applications that have not been built yet (no build information queued in the database) and saves necessary build information.
- controller/continuous_deployment/build_starter
  - Schedules builds to connected idle builders.
- controller/continuous_deployment/build_crash_detector
  - Detects builder crash and marks corresponding build as failed.
- controller/continuous_deployment/deployments_synchronizer
  - Updates "current_build" field in application to the latest success build.
- controller/backend
  - Watches "current_build" field in application, and configures actual deployments and routing.
    - Connects to Docker or Kubernetes and configures the actual containers and network. It retrieves a list of applications whose desired state is "running" and ensures that the actual system state matches it by starting/terminating containers and configuring routing. It also handles routing to the static file server configured by ss-gen.
- ss-gen (static site generator)
  - Watches "current_build" field in application, and downloads built files in order to serve them.

Each component monitors only the state it is responsible for and focuses on bringing it to the desired state.
This results in a system that is robust against failures.

## Components

### traefik-forward-auth

https://github.com/traPtitech/traefik-forward-auth

Performs user authentication using the forward auth middleware of the traefik proxy.

It is a simple HTTP server that performs the following:

- If authenticated, return 200 OK.
   - In the case of "soft" authentication, login is possible at `/_oauth/login`.
- If not authenticated, perform a 307 Temporary Redirect to carry out OAuth/OIDC authentication set in the configuration.
   - It first attempts "prompt=none", so the authorization screen should only appear once per root domain.

For detailed behavior, please refer to the README.

### Gateway (ns-gateway)

Handles HTTP requests from dashboard.

- Gateway uses [Connect](https://connect.build/) to perform typed communication on HTTP/1.1.
  - This allows utilization of wide variety of existing proxy authentication.
- Reads and writes various necessary information from the Controller, DB, and storage.
- Also triggers events to the Controller, and even if these events were to be missed, the system would automatically recover to its desired state through the Controller's internal reconciliation loop.

While "API Gateways" normally refer to components that aggregate multiple microservices, NeoShowcase's "Gateway" handles all API operations by itself since NeoShowcase's API is not that complex.

### Controller (ns-controller)

Core of NeoShowcase: update build states, configures deployments and routing.

It monitors the state in the database, and each subcomponent brings it to the desired state and eventually deploys the application.

For functionality of important subcomponents, please refer to "The Reconciliation Design" part above.

### Builder (ns-builder)

Receives build instructions from the Controller and performs the build of OCI Images (Docker images).

Currently, there are six types of build methods:

- Runtime (buildpack): It uses [Cloud Native Buildpacks](https://buildpacks.io/) to build runtime applications. It can build most common applications with zero config. It's also used in Heroku.
- Runtime (command): This method directly specifies the base image and the commands (shell scripts) for building and running during build.
- Runtime (dockerfile): This method specifies a Dockerfile for building. It offers higher customization than the previous two methods.
- Runtime (buildpack): It uses [Cloud Native Buildpacks](https://buildpacks.io/) and a specified output path to build static applications.
- Static (command): This method directly specifies the base image, build-time command (shell script), and the path where build artifacts are generated, to build static sites.
- Static (dockerfile): This method specifies a Dockerfile for building static sites. It offers higher customization than the command-based approach.

It performs builds according to each build method and pushes the generated images to the registry.

### Static-Site Generator (ns-ssgen)

Places build artifacts and serves static sites.

It can be extended to configure settings for static delivery processes like Apache HTTPD, Nginx, Caddy, etc.

### Migrator (ns-migrate)

Performs database migrations.
Migrator consists of a single shell script that executes [sqldef](https://github.com/sqldef/sqldef).

If possible, you should first define a schema that is compatible with both the old and new versions.
Then, run this migrator to make schema changes.
Afterward, migrate the necessary data manually or from within the code to perform zero-downtime migration.

However, since NeoShowcase's gateway works fine with a short controller downtime, you don't necessarily have to perform migrations that are backward-compatible.

## Compatibility with Various Backends

NeoShowcase is designed to be cloud-agnostic, based on traefik.
While it's theoretically possible to adapt it to various cloud's Ingress Controllers, implementing it may very well be challenging.

|           | Docker(traefik)         | K8s(traefik)                          | K8s(Cloud, not implemented)     |
|-----------|-------------------------|---------------------------------------|----------------|
| Routing   | traefik docker provider | traefik CRD provider                  | Ingress |
| Certificate Acquisition     | traefik Let's encrypt   | traefik Let's encrypt or cert-manager | Depends on cloud |
| Member Authentication      | traefik middleware      | traefik middleware                    | Depends on cloud        |
| Networking | docker network          | NetworkPolicy                         | Depends on cloud        |
| Container      | docker container        | StatefulSet etc.                        | StatefulSet etc. |
