# Matching API

The **Matching API** is the public-facing micro-service that pairs riders
with the nearest available driver.  
It authenticates user JWTs, calls the **Driver Location API** to fetch
nearby drivers, applies a circuit breaker for resilience, and returns
the best match (or `404` if none).

---

## Features

| Capability | Detail |
|------------|--------|
| **/v1/drivers/search** | POST GeoJSON point + radius + limit → returns closest drivers |
| **JWT auth** | Requires `{"authenticated": true}` claim |
| **Circuit breaker** | Protects outbound calls to Driver Location API |
| **Upstream health proxy** | `/v1/driver-healthcheck` checks Driver Location API |
| **Service health** | `healthz` |
| **OpenAPI 3 / Swagger** | `/swagger/index.html` (dev only) |
| **Tests** | Unit (handlers, service) |

---

## Stack

* **Go 1.24**   • **Echo v4**   • **gobreaker** (Sony)  
* **swaggo/swag** for docs  

---

## Quick start

```bash
git clone https://github.com/mrtuuro/matching-api.git
cd matching-api
touch .env

# PORT=:<enter your port>
# DRIVER_LOCATION_BASE_URL=http://127.0.0.1:10001
# DRIVER_LOCATION_API_TOKEN=<enter your `authenticated: true` jwt auth token>
vim .env

# additional to see the Makefile commands
make help

make run                   # build, generate docs, start on :9000
