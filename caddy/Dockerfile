FROM caddy:2.10-builder-alpine AS builder

RUN xcaddy build \
    --with github.com/lum8rjack/caddy-ja4h \
    --with github.com/WeidiDeng/caddy-cloudflare-ip


FROM caddy:2.10-alpine
ARG caddyfile=Caddyfile

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
COPY caddy/${caddyfile} /etc/caddy/Caddyfile
COPY client/dist/ /srv
