FROM golang:1.15.0-alpine AS build
RUN apk add --update --no-cache git
WORKDIR /go/src/github.com/traPtitech/NeoShowcase
COPY ./go.* ./
RUN go mod download
COPY . .

ARG APP_VERSION=dev
ARG APP_REVISION=local
ENV CGO_ENABLED 0
RUN go build -o /app/ns-mc -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd/ns-mc

FROM alpine:3.12.0
WORKDIR /app

EXPOSE 8081

COPY --from=build /app/* ./

ENTRYPOINT ./ns-mc
CMD serve
