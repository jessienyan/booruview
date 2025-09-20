FROM caddy:2.10-builder-alpine AS builder

RUN xcaddy build \
    --with github.com/lum8rjack/caddy-ja4h \
    --with github.com/WeidiDeng/caddy-cloudflare-ip \
    --with github.com/caddy-dns/cloudflare


FROM caddy:2.10-alpine

ARG CADDYFILE=Caddyfile

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
COPY caddy/$CADDYFILE /etc/caddy/Caddyfile
COPY client/dist/ /srv
