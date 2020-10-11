FROM golang:1.15.0-alpine AS build
RUN apk add --update --no-cache git
WORKDIR /go/src/github.com/traPtitech/NeoShowcase
COPY ./go.* ./
RUN go mod download
COPY . .

ARG APP_VERSION=dev
ARG APP_REVISION=local
ENV CGO_ENABLED 0
RUN go build -o /app/ns-ssgen -ldflags "-s -w -X main.version=$APP_VERSION -X main.revision=$APP_REVISION" ./cmd/ns-ssgen

FROM alpine:3.12.0
WORKDIR /app

RUN apk add --update ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

COPY --from=build /app/* ./

ENTRYPOINT ["/app/ns-ssgen"]
CMD ["run"]
