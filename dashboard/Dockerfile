FROM --platform=$BUILDPLATFORM node:24 AS base

WORKDIR /work

RUN corepack enable

COPY package.json yarn.lock .yarnrc.yml ./
COPY .yarn .yarn
RUN yarn --immutable

FROM --platform=$BUILDPLATFORM base AS dev

ENTRYPOINT ["yarn", "dev"]
CMD ["--host", "--port=80"]

FROM --platform=$BUILDPLATFORM base AS builder

COPY . .
RUN yarn build

FROM caddy:2 AS prod

COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /work/dist/ /usr/share/caddy/
