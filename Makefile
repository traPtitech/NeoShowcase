PROTOC_VERSION := 23.2
TBLS_VERSION := 1.65.3
SPECTRAL_VERSION := 6.4.0

GO_REPO_ROOT_PACKAGE := "github.com/traPtitech/neoshowcase"
PROTOC_OPTS := -I ./api/proto --go_out=. --go_opt=module=$(GO_REPO_ROOT_PACKAGE) --connect-go_out=. --connect-go_opt=module=$(GO_REPO_ROOT_PACKAGE)
PROTOC_OPTS_CLIENT := -I ./api/proto --es_out=./dashboard/src/api --es_opt=target=ts --connect-es_out=./dashboard/src/api --connect-es_opt=target=ts
PROTOC_SOURCES ?= $(shell find ./api/proto/neoshowcase -type f -name "*.proto" -print)
PROTOC_SOURCES_CLIENT := ./api/proto/neoshowcase/protobuf/gateway.proto ./api/proto/neoshowcase/protobuf/null.proto

TBLS_CMD := docker run --rm --net=host -v $$(pwd):/work --workdir /work -u $$(id -u):$$(id -g) ghcr.io/k1low/tbls:v$(TBLS_VERSION)
SQLDEF_CMD := APP_VERSION=local APP_REVISION=makefile mysqldef --port=5004 --user=root --password=password neoshowcase
EVANS_CMD := evans

.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: init-protoc ## Install protoc
init-protoc:
	@PROTOC_VERSION=$(PROTOC_VERSION) ./.local-dev/install-protoc.sh

.PHONY: init-protoc-tools
init-protoc-tools: ## Install other protoc tools
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
	yarn global add @bufbuild/protoc-gen-connect-es @bufbuild/protoc-gen-es

.PHONY: init
init: init-protoc init-protoc-tools ## Install commands
	go mod download
	go install github.com/k0kubun/sqldef/cmd/mysqldef@latest
	go install github.com/ktr0731/evans@latest

.PHONY: gogen
gogen: ## Generate go sources
	go generate ./...

.PHONY: protoc
protoc: ## Generate proto sources
	protoc $(PROTOC_OPTS) $(PROTOC_SOURCES)
	protoc $(PROTOC_OPTS_CLIENT) $(PROTOC_SOURCES_CLIENT)

.PHONY: db-gen-docs
db-gen-docs: ## Generate db schema docs
	@$(TBLS_CMD) doc --rm-dist

.PHONY: db-diff-docs
db-diff-docs: ## Calculate diff with current db schema docs
	@$(TBLS_CMD) diff

.PHONY: db-lint
db-lint: ## Lint current db schema docs
	@$(TBLS_CMD) lint

.PHONY: golangci-lint
golangci-lint: ## Lint go sources
	@golangci-lint run

.PHONY: build
build: ## Build containers
	@docker build -t ghcr.io/traptitech/ns-dashboard:main dashboard
	@docker compose build

.PHONY: ensure-network
ensure-network: ## Ensure apps network
	@docker network create neoshowcase_apps || return 0

.PHONY: ensure-mounts
ensure-mounts: ## Ensure local dev mounts
	@mkdir -p .local-dev/grafana
	@sudo chown -R 472:472 .local-dev/grafana
	@mkdir -p .local-dev/loki
	@sudo chown -R 10001:10001 .local-dev/loki

.PHONY: up
up: ensure-network ensure-mounts ## Setup development environment
	@docker compose up -d --build

.PHONY: down
down: ## Tear down development environment
	@docker compose down -v

.PHONY: migrate
migrate: ## Apply migration to development environment
	@$(SQLDEF_CMD) < ./migrations/schema.sql

.PHONY: ns-evans
ns-evans: ## Connect to ns api server service
	@$(EVANS_CMD) --path ./api/proto --proto neoshowcase/protobuf/gateway.proto --host ns.local.trapti.tech -p 80 repl

.PHONY: ns-controller-evans
ns-controller-evans: ## Connect to ns controller service
	@$(EVANS_CMD) --path ./api/proto --proto neoshowcase/protobuf/controller.proto --host localhost -p 10000 repl

.PHONY: db-update
db-update: migrate gogen db-gen-docs ## Apply migration, generate sqlboiler sources, and generate db schema docs

.PHONY: dind-up
dind-up: ## Setup docker-in-docker container
	docker run -it -d --privileged --name ns-test-dind -p 5555:2376 -e DOCKER_TLS_CERTDIR=/certs -v $$PWD/.local-dev/dind:/certs docker:dind

.PHONY: dind-down
dind-down: ## Tear down docker-in-docker container
	docker rm -vf ns-test-dind

.PHONY: docker-test
docker-test: ## Run docker tests
	@docker container inspect ns-test-dind > /dev/null || make dind-up
	ENABLE_DOCKER_TESTS=true DOCKER_HOST=tcp://localhost:5555 DOCKER_CERT_PATH=$$PWD/.local-dev/dind/client DOCKER_TLS_VERIFY=true go test -v ./pkg/infrastructure/backend/dockerimpl

.PHONY: k3s-import
k3s-import: ## Import images to k3s environment
	docker save ghcr.io/traptitech/ns-dashboard:main | sudo k3s ctr images import -
	docker save ghcr.io/traptitech/ns-auth-dev:main | sudo k3s ctr images import -
	docker save ghcr.io/traptitech/ns-builder:main | sudo k3s ctr images import -
	docker save ghcr.io/traptitech/ns-controller:main | sudo k3s ctr images import -
	docker save ghcr.io/traptitech/ns-gateway:main | sudo k3s ctr images import -
	# Uncomment if testing gitea-integration
	# docker save ghcr.io/traptitech/ns-gitea-integration:main | sudo k3s ctr images import -
	docker save ghcr.io/traptitech/ns-migrate:main | sudo k3s ctr images import -
	docker save ghcr.io/traptitech/ns-ssgen:main | sudo k3s ctr images import -

.PHONY: k3d-up
k3d-up: ## Setup k3s environment
	k3d cluster create ns-test --no-lb --k3s-arg "--disable=traefik,servicelb,metrics-server"
	kubectl apply -f https://raw.githubusercontent.com/traefik/traefik/v2.10/docs/content/reference/dynamic-configuration/kubernetes-crd-definition-v1.yml

.PHONY: k3d-down
k3d-down: ## Tear down k3s environment
	k3d cluster delete ns-test

.PHONY: k8s-test
k8s-test: ## Run k8s tests
	ENABLE_K8S_TESTS=true K8S_TESTS_CLUSTER_CONTEXT=k3d-ns-test go test -v ./pkg/infrastructure/backend/k8simpl
