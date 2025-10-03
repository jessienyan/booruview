FROM caddy:2.10-builder-alpine AS builder
RUN xcaddy build \
    --with github.com/lum8rjack/caddy-ja4h \
    --with github.com/shift72/caddy-geo-ip

FROM caddy:2.10-alpine
COPY --from=builder /usr/bin/caddy /usr/bin/caddy
COPY client/dist/ /srv

# This is intentionally at the end to prevent caching
RUN wget -O /etc/caddy/GeoLite2-Country.mmdb https://git.io/GeoLite2-Country.mmdb
