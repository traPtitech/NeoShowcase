TBLS_VERSION := 1.50.0
SPECTRAL_VERSION := 5.9.1

GO_REPO_ROOT_PACKAGE := "github.com/traPtitech/neoshowcase"
PROTOC_OPTS := -I ./api/proto --go_out=. --go_opt=module=$(GO_REPO_ROOT_PACKAGE) --go-grpc_out=. --go-grpc_opt=module=$(GO_REPO_ROOT_PACKAGE)
PROTOC_SOURCES ?= $(shell find ./api/proto/neoshowcase -type f -name "*.proto" -print)

SQL_MIGRATE_CMD := go run github.com/rubenv/sql-migrate/sql-migrate@latest
EVANS_CMD := go run github.com/ktr0731/evans@latest

.PHONY: init
init:
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/volatiletech/sqlboiler/v4@latest
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest

.PHONY: gogen
gogen:
	go generate ./...

.PHONY: protoc
protoc:
	protoc $(PROTOC_OPTS) $(PROTOC_SOURCES)

.PHONY: db-gen-docs
db-gen-docs:
	@if [ -d "./docs/dbschema" ]; then \
		rm -r ./docs/dbschema; \
	fi
	@docker run --rm --net=host -v $$PWD:/work ghcr.io/k1low/tbls:$(TBLS_VERSION) doc

.PHONY: db-diff-docs
db-diff-docs:
	@docker run --rm --net=host -v $$PWD:/work ghcr.io/k1low/tbls:$(TBLS_VERSION) diff

.PHONY: db-lint
db-lint:
	@docker run --rm --net=host -v $$PWD:/work ghcr.io/k1low/tbls:$(TBLS_VERSION) lint

.PHONY: swagger-lint
swagger-lint:
	@docker run --rm -it -v $$PWD:/tmp stoplight/spectral:$(SPECTRAL_VERSION) lint -r /tmp/.spectral.yml -q /tmp/api/http/swagger.yaml

.PHONY: golangci-lint
golangci-lint:
	@golangci-lint run

.PHONY: migrate-up
migrate-up:
	@$(SQL_MIGRATE_CMD) up

.PHONY: migrate-down
migrate-down:
	@$(SQL_MIGRATE_CMD) down

.PHONY: ns-builder-evans
ns-builder-evans:
	@$(EVANS_CMD) --host localhost -p 5006 -r repl

.PHONY: ns-builder-rebuild
ns-builder-rebuild:
	@docker compose up -d --build ns-builder

.PHONY: ns-ssgen-evans
ns-ssgen-evans:
	@$(EVANS_CMD) --host localhost -p 5007 -r repl

.PHONY: ns-ssgen-rebuild
ns-ssgen-rebuild:
	@docker compose up -d --build ns-ssgen

.PHONY: ns-rebuild
ns-rebuild:
	@docker compose up -d --build ns

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
	k3d cluster create ns-test --kubeconfig-switch-context=false --no-lb --k3s-server-arg "--no-deploy=traefik,servicelb,metrics-server"

.PHONY: k3d-down
k3d-down:
	k3d cluster delete ns-test

.PHONY: k8s-test
k8s-test:
	ENABLE_K8S_TESTS=true K8S_TESTS_CLUSTER_CONTEXT=k3d-ns-test go test -v ./pkg/infrastructure/backend/k8simpl
