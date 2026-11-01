# The pnpm version is decided by packageManager in package.json; this image only bootstraps pnpm.
FROM --platform=$BUILDPLATFORM ghcr.io/pnpm/pnpm:12@sha256:172468925f0ca368313c09867081a43e298eaeb13dee8f7e4356bfad45dd8b15 AS pnpm

FROM --platform=$BUILDPLATFORM node:26-slim@sha256:e4db4c914bf8ad241bd24804fc25d272d89bad97e68cdc8a1414862790a02c39 AS base

COPY --from=pnpm /opt/pnpm /opt/pnpm
RUN for bin in pnpm pnpx; do ln -s "/opt/pnpm/$bin" "/usr/local/bin/$bin"; done

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

WORKDIR /work

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN --mount=type=cache,target=/pnpm/store \
    pnpm install --frozen-lockfile

FROM --platform=$BUILDPLATFORM base AS dev
COPY . .
ENTRYPOINT ["pnpm", "dev"]
CMD ["--host", "--port=80"]

FROM --platform=$BUILDPLATFORM base AS builder
COPY . .
RUN pnpm build

FROM caddy:2@sha256:36c24ad91a7ff7bc4c87fe22bfa8df3c4b18da0773d73fc12ee2c52b98c0bb7e AS prod
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /work/dist/ /usr/share/caddy/
