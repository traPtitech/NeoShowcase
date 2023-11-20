FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

WORKDIR /work
ENV CGO_ENABLED 0

RUN apk add --update --no-cache git

COPY ./go.* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG APP_VERSION=dev
ARG APP_REVISION=local
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

FROM --platform=$BUILDPLATFORM builder AS builder-ns-migrate
ARG SQLDEF_VERSION=v0.16.12
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go install -ldflags "-s -w -X main.version=$SQLDEF_VERSION" github.com/sqldef/sqldef/cmd/mysqldef@$SQLDEF_VERSION
# keep output directory the same between platforms; workaround for https://github.com/golang/go/issues/57485
RUN cp /go/bin/mysqldef /mysqldef || cp /go/bin/"$GOOS"_"$GOARCH"/mysqldef /mysqldef

FROM --platform=$BUILDPLATFORM builder AS builder-ns
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app/ns -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd

FROM alpine:3 as base
WORKDIR /app

ARG APP_VERSION=dev
ARG APP_REVISION=local
ENV APP_VERSION=$APP_VERSION
ENV APP_REVISION=$APP_REVISION

FROM base as ns-migrate

COPY ./migrations/entrypoint.sh ./
COPY ./migrations/schema.sql ./
COPY --from=builder-ns-migrate /mysqldef /usr/local/bin/

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/schema.sql"]

FROM base AS ns
COPY --from=builder-ns /app/ns ./
ENTRYPOINT ["/app/ns"]

FROM ns as ns-auth-dev
ENTRYPOINT ["/app/ns", "auth-dev"]

FROM ns as ns-builder
ENTRYPOINT ["/app/ns", "builder"]

FROM ns as ns-controller
ENTRYPOINT ["/app/ns", "controller"]

FROM ns as ns-gateway
ENTRYPOINT ["/app/ns", "gateway"]

FROM ns as ns-gitea-integration
ENTRYPOINT ["/app/ns", "gitea-integration"]

FROM ns as ns-ssgen
ENTRYPOINT ["/app/ns", "ssgen"]
