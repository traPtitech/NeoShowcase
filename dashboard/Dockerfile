FROM --platform=$BUILDPLATFORM node:24-slim@sha256:242549cd46785b480c832479a730f4f2a20865d61ea2e404fdb2a5c3d3b73ecf AS base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

WORKDIR /work

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN corepack enable
RUN --mount=type=cache,target=/pnpm/store \
    pnpm install --frozen-lockfile

FROM --platform=$BUILDPLATFORM base AS dev
COPY . .
ENTRYPOINT ["pnpm", "dev"]
CMD ["--host", "--port=80"]

FROM --platform=$BUILDPLATFORM base AS builder
COPY . .
RUN pnpm build

FROM caddy:2@sha256:13ba145cba2f3e28fa801994876e4c086d1b95d5aa2a520a734765ffb6b12017 AS prod
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /work/dist/ /usr/share/caddy/
