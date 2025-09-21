# booruview

[booruview.com](https://booruview.com) is a web app for browsing gelbooru posts. It's very fast, has a simple interface, and anonymous with no signups or search limits.

### Running locally

Start the dev server

```bash
./dev.sh

# or
docker compose up --build
```

Then visit http://localhost:8080

### Infrastructure

The site uses three components:

- API built with Go. Proxies requests to the gelbooru API, transforms the response, and caches it
- Frontend built with Vue.js and Typescript
- Valkey (Redis) as the backend cache

It also includes [Loki](https://grafana.com/oss/loki/) and [Grafana](https://grafana.com/oss/grafana/) for monitoring.

### Donate
Booruview costs $25/mo to host. Help keep it running with a donation. Any amount helps!

- Bitcoin (BTC): `35jCje24EGQLDzW2gkmhdT1D8pxep4nN67`
- Bitcoin Cash (BCH): `19VgrR32NZGhS4pJxnBDWZD293SaCrphbs`
- Ethereum Classic (ETC): `0xB4cF6eCD1152D68204875B00d8A2D250105f51c1`
- Ethereum (ETH): `0x0f9dF73967fD57f04decb9AF9c4aa04342c153F7`
- Litecoin (LTC): `MLGX8DiuWKuSKoLXdzRtw1QHxZ6Hu6Veqr`
- Tether (USDT): `0xc5614040d18B1fC8e8E77e6599553F5EEC946C69`
- XRP: `rw2ciyaNshpHe7bCHo4bRWq6pqqynnWKQg`
    - memo: `783456807`
