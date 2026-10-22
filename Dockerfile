FROM --platform=$BUILDPLATFORM golang:1.27-alpine@sha256:d530cf60723abf1ffae0cd16bb546b379ec0550d2f3815e5b16ed69dbc00fba4 AS builder

WORKDIR /work
ENV CGO_ENABLED=0

RUN apk add --update --no-cache git

COPY ./go.* ./
COPY ./third_party ./third_party
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
ARG SQLDEF_VERSION=v3.11.16
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go install -ldflags "-s -w -X main.version=$SQLDEF_VERSION" github.com/sqldef/sqldef/v3/cmd/mysqldef@$SQLDEF_VERSION
# keep output directory the same between platforms; workaround for https://github.com/golang/go/issues/57485
RUN cp /go/bin/mysqldef /mysqldef || cp /go/bin/"$GOOS"_"$GOARCH"/mysqldef /mysqldef

FROM --platform=$BUILDPLATFORM builder AS builder-ns
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app/ns -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd

FROM alpine:3@sha256:71e9b7bae160ba167fa17414401b94cbcdaf3c32feeab9298bd794a2a0614c35 AS base
WORKDIR /app

ARG APP_VERSION=dev
ARG APP_REVISION=local
ENV APP_VERSION=$APP_VERSION
ENV APP_REVISION=$APP_REVISION

FROM base AS ns-migrate

COPY ./migrations/entrypoint.sh ./
COPY ./migrations/schema.sql ./
COPY --from=builder-ns-migrate /mysqldef /usr/local/bin/

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/schema.sql"]

FROM base AS ns
COPY --from=builder-ns /app/ns ./
ENTRYPOINT ["/app/ns"]

FROM ns AS ns-auth-dev
ENTRYPOINT ["/app/ns", "auth-dev"]

FROM ns AS ns-builder
ENTRYPOINT ["/app/ns", "builder"]

FROM ns AS ns-controller
ENTRYPOINT ["/app/ns", "controller"]

FROM ns AS ns-gateway
ENTRYPOINT ["/app/ns", "gateway"]

FROM ns AS ns-gitea-integration
ENTRYPOINT ["/app/ns", "gitea-integration"]

FROM ns AS ns-ssgen
ENTRYPOINT ["/app/ns", "ssgen"]
