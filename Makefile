TBLS_VERSION := 1.62.1
SPECTRAL_VERSION := 6.4.0

GO_REPO_ROOT_PACKAGE := "github.com/traPtitech/neoshowcase"
PROTOC_OPTS := -I ./api/proto --go_out=. --go_opt=module=$(GO_REPO_ROOT_PACKAGE) --go-grpc_out=. --go-grpc_opt=module=$(GO_REPO_ROOT_PACKAGE)
PROTOC_OPTS_CLIENT := -I ./api/proto --grpc-web_out=import_style=typescript,mode=grpcwebtext:./dashboard/src/api
PROTOC_SOURCES ?= $(shell find ./api/proto/neoshowcase -type f -name "*.proto" -print)
PROTOC_SOURCES_CLIENT := ./api/proto/neoshowcase/protobuf/apiserver.proto

TBLS_CMD := docker run --rm --net=host -v $$(pwd):/work --workdir /work -u $$(id -u):$$(id -g) ghcr.io/k1low/tbls:v$(TBLS_VERSION)
SPECTRAL_CMD := docker run --rm -it -v $$(pwd):/tmp stoplight/spectral:$(SPECTRAL_VERSION)
SQL_MIGRATE_CMD := sql-migrate
EVANS_CMD := evans

.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: init
init: ## Install commands
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/volatiletech/sqlboiler/v4@latest
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
	go install github.com/rubenv/sql-migrate/sql-migrate@latest
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
	@if [ -d "./docs/dbschema" ]; then \
		rm -r ./docs/dbschema; \
	fi
	@$(TBLS_CMD) doc

.PHONY: db-diff-docs
db-diff-docs: ## Calculate diff with current db schema docs
	@$(TBLS_CMD) diff

.PHONY: db-lint
db-lint: ## Lint current db schema docs
	@$(TBLS_CMD) lint

.PHONY: golangci-lint
golangci-lint: ## Lint go sources
	@golangci-lint run

.PHONY: up
up: ## Setup development environment
	@docker compose up -d --build

.PHONY: up-ns
up-ns: ## Rebuild ns api server
	@docker compose up -d --build ns

.PHONY: up-ns-builder
up-ns-builder: ## Rebuild ns builder
	@docker compose up -d --build ns-builder

.PHONY: up-ns-ssgem
up-ns-ssgen: ## Rebuild ns static site gen
	@docker compose up -d --build ns-ssgen

.PHONY: down
down: ## Tear down development environment
	@docker compose down

.PHONY: migrate-up
migrate-up: ## Apply migration to development environment
	@$(SQL_MIGRATE_CMD) up

.PHONY: migrate-down
migrate-down: ## Rollback migration of development environment
	@$(SQL_MIGRATE_CMD) down

.PHONY: ns-evans
ns-evans: ## Connect to ns api server service
	@$(EVANS_CMD) --host localhost -p 5009 -r repl

.PHONY: ns-builder-evans
ns-builder-evans: ## Connect to ns builder service
	@$(EVANS_CMD) --host localhost -p 5006 -r repl

.PHONY: ns-ssgen-evans
ns-ssgen-evans: ## Connect to ns static site gen service
	@$(EVANS_CMD) --host localhost -p 5007 -r repl

.PHONY: db-update
db-update: migrate-up gogen db-gen-docs ## Apply migration, generate sqlboiler sources, and generate db schema docs

.PHONY: dind-up
dind-up: ## Setup docker-in-docker container
	docker run -it -d --privileged --name ns-test-dind -p 5555:2376 -e DOCKER_TLS_CERTDIR=/certs -v $$PWD/local-dev/dind:/certs docker:dind

.PHONY: dind-down
dind-down: ## Tear down docker-in-docker container
	docker rm -vf ns-test-dind

.PHONY: docker-test
docker-test: ## Run docker tests
	@docker container inspect ns-test-dind > /dev/null || make dind-up
	ENABLE_DOCKER_TESTS=true DOCKER_HOST=tcp://localhost:5555 DOCKER_CERT_PATH=$$PWD/local-dev/dind/client DOCKER_TLS_VERIFY=true go test -v ./pkg/infrastructure/backend/dockerimpl

.PHONY: k3d-up
k3d-up: ## Setup k3s environment
	k3d cluster create ns-test --kubeconfig-switch-context=false --no-lb --k3s-arg "--disable=traefik,servicelb,metrics-server"

.PHONY: k3d-down
k3d-down: ## Tear down k3s environment
	k3d cluster delete ns-test

.PHONY: k8s-test
k8s-test: ## Run k8s tests
	ENABLE_K8S_TESTS=true K8S_TESTS_CLUSTER_CONTEXT=k3d-ns-test go test -v ./pkg/infrastructure/backend/k8simpl
