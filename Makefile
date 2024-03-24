PROTOC_VERSION := 25.3
TBLS_VERSION := 1.73.2

GO_REPO_ROOT_PACKAGE := "github.com/traPtitech/neoshowcase"
PROTOC_OPTS := -I ./api/proto --go_out=. --go_opt=module=$(GO_REPO_ROOT_PACKAGE) --connect-go_out=. --connect-go_opt=module=$(GO_REPO_ROOT_PACKAGE)
PROTOC_OPTS_CLIENT := -I ./api/proto --es_out=./dashboard/src/api --es_opt=target=ts --connect-es_out=./dashboard/src/api --connect-es_opt=target=ts
PROTOC_SOURCES ?= $(shell find ./api/proto/neoshowcase -type f -name "*.proto" -print)
PROTOC_SOURCES_CLIENT := ./api/proto/neoshowcase/protobuf/gateway.proto ./api/proto/neoshowcase/protobuf/null.proto

TBLS_CMD := docker run --rm --net=host -v $$(pwd):/work --workdir /work -u $$(id -u):$$(id -g) ghcr.io/k1low/tbls:v$(TBLS_VERSION)
SQLDEF_CMD := APP_VERSION=local APP_REVISION=makefile mysqldef --port=5004 --user=root --password=password neoshowcase
EVANS_CMD := evans

APP_VERSION ?= dev
APP_REVISION ?= local

.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ---- Init commands ----

.PHONY: init-k3d
init-k3d:
	curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash

.PHONY: init-protoc
init-protoc:
	@PROTOC_VERSION=$(PROTOC_VERSION) ./.local-dev/install-protoc.sh

.PHONY: init-protoc-tools
init-protoc-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	npm i -g @connectrpc/protoc-gen-connect-es @bufbuild/protoc-gen-es

.PHONY: init
init: init-k3d init-protoc init-protoc-tools ## Install / update required tools
	go mod download
	go install github.com/sqldef/sqldef/cmd/mysqldef@latest
	go install github.com/ktr0731/evans@latest

# ---- Ensure helpers ----

.PHONY: ensure-network
ensure-network:
	@docker network create neoshowcase_apps || return 0

.PHONY: ensure-mounts
ensure-mounts:
	@mkdir -p .local-dev/grafana
	@sudo chown -R 472:472 .local-dev/grafana
	@mkdir -p .local-dev/loki
	@sudo chown -R 10001:10001 .local-dev/loki

.PHONY: ensure-db
ensure-db: ensure-network
	docker compose up -d --wait mysql mongo

# ---- Code gen commands ----

.PHONY: migrate
migrate: ensure-db
	@$(SQLDEF_CMD) < ./migrations/schema.sql

.PHONY: gen-go
gen-go: ensure-db
	go generate ./...

.PHONY: gen-proto
gen-proto:
	protoc $(PROTOC_OPTS) $(PROTOC_SOURCES)
	protoc $(PROTOC_OPTS_CLIENT) $(PROTOC_SOURCES_CLIENT)

.PHONY: gen-db-docs
gen-db-docs: ensure-db
	@$(TBLS_CMD) doc --rm-dist

.PHONY: gen
gen: migrate gen-go gen-proto gen-db-docs ## Regenerate wire, sqlboiler, protobuf, and db docs

# ---- Test commands ----

.PHONY: test-up-docker
test-up-docker:
	@docker container inspect ns-test-dind > /dev/null \
	|| docker run -it -d --privileged --name ns-test-dind -p 5555:2376 -e DOCKER_TLS_CERTDIR=/certs -v $$PWD/.local-dev/dind:/certs docker:dind

.PHONY: test-down-docker
test-down-docker:
	docker rm -vf ns-test-dind

.PHONY: test-up-k8s
test-up-k8s:
	@k3d cluster list ns-test > /dev/null \
	|| k3d cluster create ns-test --no-lb --k3s-arg "--disable=traefik,servicelb,metrics-server" \
	&& kubectl apply -f https://raw.githubusercontent.com/traefik/traefik/v2.10/docs/content/reference/dynamic-configuration/kubernetes-crd-definition-v1.yml

.PHONY: test-down-k8s
test-down-k8s:
	k3d cluster delete ns-test

.PHONY: test-run-docker
test-run-docker:
	ENABLE_DOCKER_TESTS=true DOCKER_HOST=tcp://localhost:5555 DOCKER_CERT_PATH=$$PWD/.local-dev/dind/client DOCKER_TLS_VERIFY=true go test -v ./pkg/infrastructure/backend/dockerimpl

.PHONY: test-run-k8s
test-run-k8s:
	ENABLE_K8S_TESTS=true K8S_TESTS_CLUSTER_CONTEXT=k3d-ns-test go test -v ./pkg/infrastructure/backend/k8simpl

.PHONY: test-run
test-run:
	ENABLE_DOCKER_TESTS=true \
	DOCKER_HOST=tcp://localhost:5555 \
	DOCKER_CERT_PATH=$$PWD/.local-dev/dind/client \
	DOCKER_TLS_VERIFY=true \
	ENABLE_K8S_TESTS=true \
	K8S_TESTS_CLUSTER_CONTEXT=k3d-ns-test \
	go test -shuffle=on -v ./...

.PHONY: test-docker
test-docker: ensure-db test-up-docker test-run-docker test-down-docker

.PHONY: test-k8s
test-k8s: ensure-db test-up-k8s test-run-k8s test-down-k8s

.PHONY: test
test: ensure-db test-up-docker test-up-k8s test-run test-k8s test-down-docker ## Run all tests

# ---- Debug commands ----

.PHONY: debug-gateway
debug-gateway: ## Connect to gateway service
	@$(EVANS_CMD) --path ./api/proto --proto neoshowcase/protobuf/gateway.proto --host ns.local.trapti.tech -p 80 repl

.PHONY: debug-controller
debug-controller: ## Connect to controller service
	@$(EVANS_CMD) --path ./api/proto --proto neoshowcase/protobuf/controller.proto --host localhost -p 10000 repl

# ---- All in one commands ----

.PHONY: build-dashboard
build-dashboard:
	@docker build -t ghcr.io/traptitech/ns-dashboard:main dashboard

.PHONY: build
build:
	@docker compose build --build-arg APP_VERSION=$(APP_VERSION) --build-arg APP_REVISION=$(APP_REVISION)

.PHONY: up
up: ensure-network ensure-mounts build ## Start development environment
	@docker compose up -d

.PHONY: down
down: ## Tear down development environment
	@docker compose down -v
