FROM node:18 AS base

WORKDIR /work

COPY package.json yarn.lock ./
RUN yarn --frozen-lockfile

FROM base AS dev

ENTRYPOINT ["yarn", "dev"]
CMD ["--host", "--port=80"]

FROM base AS builder

COPY . .
RUN yarn build

FROM caddy:2 AS prod

COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /work/dist/ /usr/share/caddy/
