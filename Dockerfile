# build stage
FROM golang:1.19-alpine AS builder
RUN apk add --update --no-cache git
WORKDIR /go/src/github.com/traPtitech/NeoShowcase
COPY ./go.* ./
RUN go mod download
COPY . .

ARG APP_VERSION=dev
ARG APP_REVISION=local
ENV CGO_ENABLED 0

FROM builder as builder-ns
RUN go build -o /app/ns -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd/ns

FROM builder as builder-ns-builder
RUN go build -o /app/ns-builder -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd/ns-builder

FROM builder as builder-ns-mc
RUN go build -o /app/ns-mc -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd/ns-mc

FROM builder as builder-ns-migrate
RUN go build -o /app/ns-migrate -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd/ns-migrate

FROM builder as builder-ns-ssgen
RUN go build -o /app/ns-ssgen -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd/ns-ssgen

FROM alpine:3 as base
WORKDIR /app

# artifact images
FROM base as ns
EXPOSE 8080
COPY --from=builder-ns /app/ns ./
ENTRYPOINT ["/app/ns"]
CMD ["run"]

FROM base as ns-builder
COPY --from=builder-ns-builder /app/ns-builder ./
ENTRYPOINT ["/app/ns-builder"]
CMD ["run"]

FROM base as ns-mc
EXPOSE 8081
COPY --from=builder-ns-mc /app/ns-mc ./
ENTRYPOINT ["/app/ns-mc"]
CMD ["serve"]

FROM base as ns-ssgen
COPY --from=builder-ns-ssgen /app/ns-ssgen ./
ENTRYPOINT ["/app/ns-ssgen"]
CMD ["run"]

FROM base as ns-migrate
COPY --from=builder-ns-migrate /app/ns-migrate ./
ENTRYPOINT ["/app/ns-migrate"]

FROM base as ns-all
EXPOSE 8080 8081
COPY --from=builder-ns /app/ns ./
COPY --from=builder-ns-builder /app/ns-builder ./
COPY --from=builder-ns-mc /app/ns-mc ./
COPY --from=builder-ns-ssgen /app/ns-ssgen ./
COPY --from=builder-ns-migrate /app/ns-migrate ./
