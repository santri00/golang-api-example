# golang-api-example

A example REST API written in Go that demonstrates a layered architecture for a vehicle inventory management system. It covers HTTP handling, service logic, Redis-backed storage, unit testing, linting, and Docker deployment, making it a practical reference for building production-style Go APIs.

---

## Features

- **Vehicle Inventory API** — add, list, and retrieve vehicle records over HTTP
- **Layered architecture** — handler → service → repository, each with a clear interface boundary
- **Redis storage** — vehicle records stored as Redis hashes and lists using atomic transaction pipelines. > **Security note:** this example connects to Redis without authentication. Redis security (AUTH password, TLS, ACLs) is strongly recommended and should be enabled before deploying to any production environment.
- **Computed fields** — `days_in_inventory` and `is_under_warranty` derived at read time, not stored
- **Pagination** — `GET /vehicles` supports `offset` and `limit` query parameters (defaults: `0` / `-1` for all)
- **Health check** — `GET /ping` verifies Redis connectivity
- **Unit tests** — handler, service, and repository layers fully covered; repository tests use miniredis (no live Redis required)
- **Quality gates** — `go fmt`, `go vet`, and `golangci-lint` enforced via Makefile before any run
- **Dockerized** — multi-stage Dockerfile produces a minimal distroless image; `docker compose` wires up the API and Redis together

---

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/inventory/v1/ping` | Health check — pings Redis |
| `POST` | `/inventory/v1/vehicles` | Add a vehicle record |
| `GET` | `/inventory/v1/vehicles` | List vehicles (`?offset=0&limit=-1`) |
| `GET` | `/inventory/v1/vehicles/{stock_id}` | Get a vehicle by stock ID |

---

## Prerequisites

| Tool | Purpose | Install |
|------|---------|---------|
| Go 1.19+ | Build and test | https://go.dev/dl |
| Docker + Docker Compose | Run the full stack | https://docs.docker.com/get-docker |
| GNU Make | Run quality gates | `winget install GnuWin32.Make` |
| golangci-lint | Linting | `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` |
| jq | JSON pretty-printing in scripts | `winget install jqlang.jq` |
| curl | HTTP requests in scripts | https://curl.se |

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/santri00/golang-api-example.git
cd golang-api-example
```

### 2. Install Go dependencies

```bash
go mod download
```

### 3. Run quality gates

The Makefile runs `go fmt`, `go vet`, `golangci-lint`, and `go test` in sequence. Each step must pass before the next runs.

```bash
make check
```

Available targets:

```
make fmt        # Format all packages
make vet        # Vet all packages
make lint       # Lint all packages
make test       # Run all tests
make test-race  # Run tests with race detector (requires CGO + C compiler)
make check      # Run all quality gates: fmt → vet → lint → test
make up         # Run quality gates then start the stack via docker compose
```

### 4. Start the stack

```bash
make up
```

This runs all quality gates first, then starts the API (port `3001`) and Redis (port `6380`) via `docker compose --build`.

To skip the quality gates and start directly:

```bash
docker compose up --build
```

---

## Running Tests

```bash
# All tests (unit + repository tests using miniredis — no live Redis needed)
make test

# With race detector (requires CGO_ENABLED=1 and a C compiler)
make test-race

# A specific package
go test -count=1 ./pkg/api/...
go test -count=1 ./pkg/services/vehicle/...
go test -count=1 ./pkg/storage/vehicle/...
```

---

## Seeding Data and Testing the API

With the server running (`make up` or `docker compose up`), use the seed script to POST all 10 sample vehicle records and query each one back:

```bash
./scripts/seed_and_query.sh
```

The script will:

1. Check server health via `GET /inventory/v1/ping`
2. For each JSON file in `scripts/Requests/`:
   - `POST /inventory/v1/vehicles` with the file body
   - Extract the returned `stock_id`
   - `GET /inventory/v1/vehicles/{stock_id}` — pretty-printed via `jq`
3. `GET /inventory/v1/vehicles` — list all records, pretty-printed
4. Print a summary of all stock IDs created

To target a different host:

```bash
BASE_URL=http://localhost:8080 ./scripts/seed_and_query.sh
```

---

## Project Structure

```
.
├── cmd/api/              # Entry point — wires config, Redis, and server
├── configs/              # config.yaml (name, port, Redis address)
├── pkg/
│   ├── api/              # HTTP handlers and response writer
│   ├── apiconfig/        # Config loader (koanf)
│   ├── models/           # API-facing request/response types
│   ├── router/           # Gorilla Mux route registration
│   ├── services/vehicle/ # Business logic and model↔entity mapper
│   ├── storage/
│   │   ├── cachekeys/    # Redis key constants
│   │   ├── entities/     # Storage-facing types
│   │   └── vehicle/      # Redis repository implementation
│   └── utilities/        # UUID generation
├── scripts/
│   ├── Requests/         # Sample POST request bodies (JSON)
│   └── seed_and_query.sh # Seed and query script
├── .golangci.yml         # golangci-lint configuration
├── .gitattributes        # LF line ending enforcement
├── Makefile              # Quality gate and docker compose targets
├── Dockerfile            # Multi-stage build → distroless runtime
└── docker-compose.yaml   # API + Redis services
```
