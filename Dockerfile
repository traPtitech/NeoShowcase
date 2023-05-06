# build stage
FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS builder

WORKDIR /go/src/github.com/traPtitech/NeoShowcase
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

FROM --platform=$BUILDPLATFORM builder as builder-ns-auth-dev
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app/ns-auth-dev -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd/ns-auth-dev

FROM --platform=$BUILDPLATFORM builder as builder-ns-migrate
ARG SQLDEF_VERSION=v0.15.22
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go install -ldflags "-s -w -X main.version=$SQLDEF_VERSION" github.com/k0kubun/sqldef/cmd/mysqldef@$SQLDEF_VERSION
# keep output directory the same between platforms; workaround for https://github.com/golang/go/issues/57485
RUN cp /go/bin/mysqldef /mysqldef || cp /go/bin/"$GOOS"_"$GOARCH"/mysqldef /mysqldef

FROM --platform=$BUILDPLATFORM builder as builder-ns-builder
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app/ns-builder -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd/ns-builder

FROM --platform=$BUILDPLATFORM builder as builder-ns-controller
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app/ns-controller -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd/ns-controller

FROM --platform=$BUILDPLATFORM builder as builder-ns-gateway
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app/ns-gateway -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd/ns-gateway

FROM --platform=$BUILDPLATFORM builder as builder-ns-mc
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app/ns-mc -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd/ns-mc

FROM --platform=$BUILDPLATFORM builder as builder-ns-ssgen
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app/ns-ssgen -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd/ns-ssgen

FROM alpine:3 as base
WORKDIR /app

ARG APP_VERSION=dev
ARG APP_REVISION=local
ENV APP_VERSION=$APP_VERSION
ENV APP_REVISION=$APP_REVISION

FROM base AS ns-auth-dev
EXPOSE 4181
COPY --from=builder-ns-auth-dev /app/ns-auth-dev ./
ENTRYPOINT ["/app/ns-auth-dev"]

FROM base as ns-migrate

COPY ./migrations/entrypoint.sh ./
COPY ./migrations/schema.sql ./
COPY --from=builder-ns-migrate /mysqldef /usr/local/bin/

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/schema.sql"]

FROM base as ns-builder
COPY --from=builder-ns-builder /app/ns-builder ./
ENTRYPOINT ["/app/ns-builder"]
CMD ["run"]

FROM base as ns-controller
EXPOSE 10000
COPY --from=builder-ns-controller /app/ns-controller ./
ENTRYPOINT ["/app/ns-controller"]
CMD ["run"]

FROM base as ns-gateway
EXPOSE 8080
COPY --from=builder-ns-gateway /app/ns-gateway ./
ENTRYPOINT ["/app/ns-gateway"]
CMD ["run"]

FROM base as ns-mc
EXPOSE 8080
COPY --from=builder-ns-mc /app/ns-mc ./
ENTRYPOINT ["/app/ns-mc"]
CMD ["serve"]

FROM base as ns-ssgen
EXPOSE 8080
COPY --from=builder-ns-ssgen /app/ns-ssgen ./
ENTRYPOINT ["/app/ns-ssgen"]
CMD ["run"]

FROM base as ns-all
EXPOSE 8080 10000
COPY --from=builder-ns-builder /app/ns-builder ./
COPY --from=builder-ns-controller /app/ns-controller ./
COPY --from=builder-ns-gateway /app/ns-gateway ./
COPY --from=builder-ns-mc /app/ns-mc ./
COPY --from=builder-ns-ssgen /app/ns-ssgen ./
