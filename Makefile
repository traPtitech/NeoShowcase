TBLS_VERSION := 1.43.1
SPECTRAL_VERSION := 5.6.0

.PHONY: init
init:
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install github.com/markbates/pkger/cmd/pkger
	go install github.com/volatiletech/sqlboiler/v4
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql
	go install github.com/rubenv/sql-migrate/sql-migrate

.PHONY: gogen
gogen:
	go generate ./...

.PHONY: db-gen-docs
db-gen-docs:
	@if [ -d "./docs/dbschema" ]; then \
		rm -r ./docs/dbschema; \
	fi
	@docker run --rm --net=host -v $$PWD:/work k1low/tbls:$(TBLS_VERSION) doc

.PHONY: db-diff-docs
db-diff-docs:
	@docker run --rm --net=host -v $$PWD:/work k1low/tbls:$(TBLS_VERSION) diff

.PHONY: db-lint
db-lint:
	@docker run --rm --net=host -v $$PWD:/work k1low/tbls:$(TBLS_VERSION) lint

.PHONY: swagger-lint
swagger-lint:
	@docker run --rm -it -v $$PWD:/tmp stoplight/spectral:$(SPECTRAL_VERSION) lint -r /tmp/.spectral.yml -q /tmp/api/http/swagger.yaml

.PHONY: golangci-lint
golangci-lint:
	@golangci-lint run

.PHONY: migrate-up
migrate-up:
	@sql-migrate up

.PHONY: migrate-down
migrate-down:
	@sql-migrate down
