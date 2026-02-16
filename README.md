# booruview

[booruview.com](https://booruview.com) is a web app for browsing gelbooru posts. It's fast, simple, anonymous, and free.

You can help support the project with a small [donation](#donate).

## Tech stack
- Frontend SPA built with Vue + Typescript. It's pretty minimal and lightweight
- Simple backend API built with Go. Acts as a "smart" proxy for the Gelbooru API
- Valkey (Redis) for caching tags, searches, and autocomplete. Also keeps track of rate limit bans
- Caddy as the web server
- Docker and Docker Compose

The production site uses Loki/Grafana ([docker-compose.prod.yml](./docker-compose.prod.yml)), but by no means is this necessary to get the site up and running.

## Quickstart

This guide will get you setup with your own booruview instance.

It also assumes two things:

- You are on Linux
- You have Docker and Docker Compose installed

### Step 0. Create a Gelbooru account

I highly recommend creating a GB account to get an API key. Without it the site might not function correctly.

Go to the [accounts page](https://gelbooru.com/index.php?page=account&s=home) and signup or login.

Once you are signed in, go to the [options page](https://gelbooru.com/index.php?page=account&s=options) and scroll all the way down. Have the API key and user ID ready for later.

### Step 1. Initial setup

Start by grabbing the code and creating the `.env` file:

```sh
git clone https://codeberg.org/jessienyan/booruview.git
cd booruview
cp .env.example .env
```

Open `.env` in your editor. Update these three vars:

- `CADDY_HOST`
- `GELBOORU_USERID`
- `GELBOORU_APIKEY`

Use the API key and user ID from the previous step here.

### Step 2. 
TODO

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
