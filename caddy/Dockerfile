FROM caddy:2.10-builder-alpine AS builder

RUN xcaddy build --with github.com/lum8rjack/caddy-ja4h


FROM caddy:2.10-alpine

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
COPY caddy/Caddyfile /etc/caddy/Caddyfile
COPY client/dist/ /srv
