FROM --platform=$BUILDPLATFORM node:24-slim@sha256:2fe369e969550cde8e867afc3fe370b260140cab4a23d467074295b42163d553 AS base

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

FROM caddy:2@sha256:ec18ee54aab3315c22e25f3b2babda73ff8007d39b13b3bd1bfffa2f0444c7d9 AS prod
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /work/dist/ /usr/share/caddy/
