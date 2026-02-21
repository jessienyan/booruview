FROM caddy:2-builder-alpine AS builder
RUN xcaddy build \
    --with github.com/shift72/caddy-geo-ip

FROM caddy:2-alpine
COPY --from=builder /usr/bin/caddy /usr/bin/caddy

ENV USE_GRAFANA=0
ENV USE_MEDIA_PROXY=0

COPY caddy/generate_caddyfile.sh caddy/start.sh /etc/caddy
RUN mkdir -p /etc/caddy/components/ && chmod +x /etc/caddy/start.sh
COPY caddy/Caddyfile* /etc/caddy/components/

COPY client/dist/ /srv

# This is intentionally at the end to prevent caching
RUN wget -O /etc/caddy/GeoLite2-Country.mmdb https://git.io/GeoLite2-Country.mmdb

CMD ["/etc/caddy/start.sh"]
