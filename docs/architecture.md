Here is the translation of the provided text:

# Architecture

## Overview

NeoShowcase is a PaaS (Platform as a Service) that traP members (users) can use. In this context, PaaS refers to a service that runs applications on behalf of users with just the code of the application and minimal configuration brought in by the user. Users manage the code of their applications on Git and push it to services like GitHub or Gitea. NeoShowcase automatically detects updates and builds applications according to user settings.

One of the major design principles of NeoShowcase is the "reconciliation loop" (also known as the reconciliation pattern, synchronization, eventual consistency, etc.). This is heavily influenced by the processing method of Kubernetes controllers (in fact, NeoShowcase can also be considered a type of controller). In contrast, there is an "event-driven" design, where internal processes are triggered by user actions or external events and perform some processing, and may fire events as needed. However, event-driven designs can become complex when events are lost during communication or when processes triggered by events fail or crash. In contrast, the "reconciliation loop" periodically monitors the system's state and only performs "reconciliation" processing if the current state is not as desired.

While the event-driven approach tends to have slightly higher computational overhead, it eliminates the need for complex retry logic. Each process only needs to manage bringing its responsible state to the desired state, making the processing simpler.

Each component has its own "reconciliation loop."

- controller/repository_fetcher: Every 3 minutes or when an event occurs, it retrieves a list of repositories and applications from the database, fetches the latest commit hash corresponding to the git ref specified by each application, and stores that value in the database.
   - This is influenced by the way ArgoCD updates. Argo CD polls Git repositories every three minutes to detect changes to the manifests.

- controller/continuous_deployment/build_registerer: Every 3 minutes or when an event occurs, it lists applications that have not been built yet (no build information queued in the database) and saves the necessary build information in the database as "queued."

- controller/continuous_deployment/build_starter: Every 3 minutes or when an event occurs, it instructs the connected builder to perform the next queued build.

- controller/continuous_deployment/build_crash_detector: Every 1 minute, it detects when the builder crashes or becomes unresponsive and records the corresponding build as a failure.

- controller/continuous_deployment/deployments_synchronizer: Every 3 minutes or when an event occurs, it checks whether the latest build has completed for each "should be running" application and, if so, requests synchronization with backend and ssgen according to the user's settings.

- controller/backend: It connects to Docker or Kubernetes systems and manages the actual launching of containers and network management. It receives a list of "should be running" applications and ensures that the actual system state matches it by starting/terminating containers and routing as needed. It also handles routing to the delivery server set by ss-gen.

- ss-gen: Every 3 minutes or when an event occurs, it downloads static site files from storage and arranges them for delivery.

In this way, each component monitors only the state it is responsible for and focuses on bringing it to the desired state, resulting in a system that is robust against failures.

## Components

### traefik-forward-auth

https://github.com/traPtitech/traefik-forward-auth

This component performs user authentication using the forward auth middleware of the traefik proxy.

Basically, it does the following:

- If authenticated, return 200 OK.
   - In the case of "soft," login is possible at `/_oauth/login`.
- If not authenticated, perform a 307 Temporary Redirect to carry out OAuth/OIDC authentication set in the configuration.
   - OAuth2 requests first attempt "prompt=none," so the authorization screen only appears once per root domain.

This is an HTTP server that performs only these actions. For detailed behavior, please refer to the README.

### Gateway (ns-gateway)

This is the part that users operate directly from the front end (dashboard). It uses [Connect](https://connect.build/) to perform typed and secure communication while using existing proxy authentication on HTTP/1.1.

While API Gateways often aggregate multiple microservices, this Gateway handles all API operations since NeoShowcase's API is not that complex. It receives requests and reads and writes various necessary information from Controller, DB, and Storage. It also triggers events to the Controller, and even if these events were to be missed, the system would automatically recover its state through the Controller's internal reconciliation loop.

### Controller (ns-controller)

This is the core of NeoShowcase. It monitors the state recorded in the database, and each subcomponent brings it to the desired state and eventually deploys the application.

For the functionality of important subcomponents, please refer to the description above.

### Builder (ns-builder)

It receives build instructions from the Controller and actually performs the build of OCI Images (Docker images).

Currently, there are five types of build methods:

- Runtime (buildpack): It uses [Cloud Native Buildpacks](https://buildpacks.io/) to build runtime applications. It can build most common applications with zero config. It's also used in Heroku.
- Runtime (command): This method directly specifies the base image and the commands (shell scripts) for building and running during build.
- Runtime (dockerfile): This method specifies a Dockerfile for building. It offers higher customization than the previous two methods.
- Static (command): This method directly specifies the base image, build-time command (shell script), and the path where build artifacts are generated, to build static sites.
- Static (dockerfile): This method specifies a Dockerfile for building static sites. It offers higher customization than the command-based approach.

It performs builds according to each build method and pushes the generated images to the registry.

### Static-Site Generator (ns-ssgen)

It places the build artifacts of static sites and configures them for delivery.

It can be extended to configure settings for static delivery processes like Apache HTTPD, Nginx, Caddy, etc.

### Migrator (ns-migrate)

It performs database migrations. It relies solely on scripts executed from [sqldef](https://github.com/k0kubun/sqldef), without any Go code.

During migration, it first defines a schema that is compatible with both the old and new versions. Then, it uses sqldef to make schema changes. Afterward, it complements the necessary data manually or from within the code to enable zero-downtime migration. However, since NeoShowcase's application itself continues to run without controller intervention, it's not necessary to perform migrations that are not backwards-compatible with the schema.

## Compatibility with Various Backends

NeoShowcase is designed to be cloud-agnostic, based on traefik. While it's theoretically possible to adapt it to various cloud's Ingress Controllers, implementing it extensively can become challenging.

|           | Docker(traefik)         | K8s(traefik)                          | K8s(Cloud)     |
|-----------|-------------------------|---------------------------------------|----------------|
| Routing   | traefik docker provider | traefik CRD provider                  | Ingress (not implemented) |
| Certificate Acquisition | traefik
