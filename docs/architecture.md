# Architecture

## Overview

NeoShowcase is a PaaS (Platform as a Service) that traP members (users) can use.
Here, PaaS refers to a service that runs applications on behalf of users just with the source code of the application and a minimal configuration.
Users manage the source code of their applications on Git and push it to services like GitHub or Gitea.
NeoShowcase automatically detects updates and builds applications according to user settings.

One of the major design principles of NeoShowcase is the "reconciliation loop" (also known as the reconciliation pattern, synchronization, eventual consistency, etc.).
This is heavily influenced by the processing method of Kubernetes controllers (in fact, NeoShowcase can also be considered a Kubernetes controller).
In contrast, there is an "event-driven" design, where internal processes are triggered by user actions or external events, and perform some processing which may fire more events.
However, event-driven designs can become complex when events are lost during communication or when processes triggered by events fail or crash.
In contrast, the "reconciliation loop" periodically monitors the system's state and performs "reconciliation" only when it is not in the desired state.

While the event-driven approach tends to have slightly higher computational overhead, it eliminates the need for complex retry logic.
Each process only needs to bring its responsible state to the desired state, making the logic simpler.

Each component has its own "reconciliation loop."

- controller/repository_fetcher: Every 3 minutes or when an event occurs, it retrieves a list of repositories and applications from the database, fetches the latest commit hash corresponding to the git ref specified by each application, and stores that value in the database.
   - This is influenced by the way ArgoCD updates. Argo CD polls Git repositories every three minutes to detect changes to the manifests. https://argo-cd.readthedocs.io/en/stable/operator-manual/webhook/
- controller/continuous_deployment/build_registerer: Every 3 minutes or when an event occurs, it lists applications that have not been built yet (no build information queued in the database) and saves the necessary build information in the database as "queued."
- controller/continuous_deployment/build_starter: Every 3 minutes or when an event occurs, it instructs connected builders to process the next queued build.
- controller/continuous_deployment/build_crash_detector: Every 1 minute, it detects whether the builder crashed or became unresponsive and records the corresponding build as a failure.
- controller/continuous_deployment/deployments_synchronizer: Every 3 minutes or when an event occurs, it checks whether the latest build has completed for each application whose desired state is "running" and, if so, requests synchronization with backend and ssgen.
- controller/backend: It connects to Docker or Kubernetes and manages the actual containers and network. It receives a list of applications whose desired state is "running" and ensures that the actual system state matches it by starting/terminating containers and configuring routing. It also handles routing to the static file server configured by ss-gen.
- ss-gen: Every 3 minutes or when an event occurs, it downloads static site files from storage and arranges them for delivery.

In this way, each component monitors only the state it is responsible for and focuses on bringing it to the desired state, resulting in a system that is robust against failures.

## Components

### traefik-forward-auth

https://github.com/traPtitech/traefik-forward-auth

This component performs user authentication using the forward auth middleware of the traefik proxy.

It is a simple HTTP server that performs the following:

- If authenticated, return 200 OK.
   - In the case of "soft" authentication, login is possible at `/_oauth/login`.
- If not authenticated, perform a 307 Temporary Redirect to carry out OAuth/OIDC authentication set in the configuration.
   - It first attempts "prompt=none", so the authorization screen should only appear once per root domain.

For detailed behavior, please refer to the README.

### Gateway (ns-gateway)

This is the part where users operate directly from the front-end (dashboard).
It uses [Connect](https://connect.build/) to perform typed communication while using existing proxy authentication on HTTP/1.1.

Normally, "API Gateways" often refer to components that aggregate multiple microservices, NeoShowcase's Gateway handles all API operations by itself since NeoShowcase's API is not that complex.
It reads and writes various necessary information from the Controller, DB, and storage.
It also triggers events to the Controller, and even if these events were to be missed, the system would automatically recover to its desired state through the Controller's internal reconciliation loop.

### Controller (ns-controller)

Core of NeoShowcase

It monitors the state in the database, and each subcomponent brings it to the desired state and eventually deploys the application.

For the functionality of important subcomponents, please refer to the description above.

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

It places the build artifacts of static sites and configures them for delivery.

It can be extended to configure settings for static delivery processes like Apache HTTPD, Nginx, Caddy, etc.

### Migrator (ns-migrate)

Performs database migrations.
Migrator consists of a single shell script that executes [sqldef](https://github.com/k0kubun/sqldef).

If possible, you should first define a schema that is compatible with both the old and new versions.
Then, run this migrator to make schema changes.
Afterward, migrate the necessary data manually or from within the code to perform zero-downtime migration.

However, since NeoShowcase's gateway works fine with a short controller downtime, you don't necessarily have to perform migrations that are backward-compatible.

## Compatibility with Various Backends

NeoShowcase is designed to be cloud-agnostic, based on traefik. While it's theoretically possible to adapt it to various cloud's Ingress Controllers, implementing it extensively can become challenging.

|           | Docker(traefik)         | K8s(traefik)                          | K8s(Cloud, not implemented)     |
|-----------|-------------------------|---------------------------------------|----------------|
| Routing   | traefik docker provider | traefik CRD provider                  | Ingress |
| Certificate Acquisition     | traefik Let's encrypt   | traefik Let's encrypt or cert-manager | Depends on cloud |
| Member Authentication      | traefik middleware      | traefik middleware                    | Depends on cloud        |
| Networking | docker network          | NetworkPolicy                         | Depends on cloud        |
| Container      | docker container        | StatefulSet etc.                        | StatefulSet etc. |
