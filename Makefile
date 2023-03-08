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

.PHONY: init
init:
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/volatiletech/sqlboiler/v4@latest
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
	go install github.com/rubenv/sql-migrate/sql-migrate@latest
	go install github.com/ktr0731/evans@latest

.PHONY: gogen
gogen:
	go generate ./...

.PHONY: protoc
protoc:
	protoc $(PROTOC_OPTS) $(PROTOC_SOURCES)
	protoc $(PROTOC_OPTS_CLIENT) $(PROTOC_SOURCES_CLIENT)

.PHONY: db-gen-docs
db-gen-docs:
	@if [ -d "./docs/dbschema" ]; then \
		rm -r ./docs/dbschema; \
	fi
	@$(TBLS_CMD) doc

.PHONY: db-diff-docs
db-diff-docs:
	@$(TBLS_CMD) diff

.PHONY: db-lint
db-lint:
	@$(TBLS_CMD) lint

.PHONY: swagger-lint
swagger-lint:
	@$(SPECTRAL_CMD) lint -r /tmp/.spectral.yml -q /tmp/api/http/swagger.yaml

.PHONY: golangci-lint
golangci-lint:
	@golangci-lint run

.PHONY: up
up:
	@docker compose up -d --build

.PHONY: up-ns
up-ns:
	@docker compose up -d --build ns

.PHONY: up-ns-builder
up-ns-builder:
	@docker compose up -d --build ns-builder

.PHONY: up-ns-ssgem
up-ns-ssgen:
	@docker compose up -d --build ns-ssgen

.PHONY: down
down:
	@docker compose down

.PHONY: migrate-up
migrate-up:
	@$(SQL_MIGRATE_CMD) up

.PHONY: migrate-down
migrate-down:
	@$(SQL_MIGRATE_CMD) down

.PHONY: ns-evans
ns-evans:
	@$(EVANS_CMD) --host localhost -p 5009 -r repl

.PHONY: ns-builder-evans
ns-builder-evans:
	@$(EVANS_CMD) --host localhost -p 5006 -r repl

.PHONY: ns-ssgen-evans
ns-ssgen-evans:
	@$(EVANS_CMD) --host localhost -p 5007 -r repl

.PHONY: db-update
db-update: migrate-up gogen db-gen-docs

.PHONY: dind-up
dind-up:
	docker run -it -d --privileged --name ns-test-dind -p 5555:2376 -e DOCKER_TLS_CERTDIR=/certs -v $$PWD/local-dev/dind:/certs docker:dind

.PHONY: dind-down
dind-down:
	docker rm -vf ns-test-dind

.PHONY: docker-test
docker-test:
	@docker container inspect ns-test-dind > /dev/null || make dind-up
	ENABLE_DOCKER_TESTS=true DOCKER_HOST=tcp://localhost:5555 DOCKER_CERT_PATH=$$PWD/local-dev/dind/client DOCKER_TLS_VERIFY=true go test -v ./pkg/infrastructure/backend/dockerimpl

.PHONY: k3d-up
k3d-up:
	k3d cluster create ns-test --kubeconfig-switch-context=false --no-lb --k3s-arg "--no-deploy=traefik,servicelb,metrics-server"

.PHONY: k3d-down
k3d-down:
	k3d cluster delete ns-test

.PHONY: k8s-test
k8s-test:
	ENABLE_K8S_TESTS=true K8S_TESTS_CLUSTER_CONTEXT=k3d-ns-test go test -v ./pkg/infrastructure/backend/k8simpl
