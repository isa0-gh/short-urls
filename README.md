# Short URLs
---
A simple URL shortener with Swagger and a web UI.

Stack:
- Go
- Postgres
- Docker
- Redis
- Swagger
- sqlc

## Prerequisites

- Go 1.27 or newer
- Docker and Docker Compose
- [just](https://github.com/casey/just) for the repository shortcuts

## Quick start

```bash
# Clone the repository and copy the environment file.
git clone https://github.com/isa0-gh/short-urls
cd short-urls
cp .env.example .env

# Build the Docker image. Requires just.
just docker-build v1.0.0

# Without just, run:
# docker build -t github.com/isa0-gh/short-urls:v1.0.0 .

docker compose -f docker-compose.yaml up
```

The web UI is available at <http://localhost:8080> and the Swagger UI is
available at <http://localhost:8080/swagger/index.html>.

The API base URL is <http://localhost:8080/api/v1>.

## Configuration

The application loads configuration from `.env` when running locally. The
Docker Compose file provides these values automatically:

| Variable | Description | Default |
| --- | --- | --- |
| `PORT` | HTTP server port | `8080` |
| `SESSION_EXP` | Session duration | `1h` |
| `POSTGRES_HOST` | PostgreSQL host | `postgres` in Compose |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |
| `POSTGRES_USER` | PostgreSQL user | `postgres` |
| `POSTGRES_PASSWORD` | PostgreSQL password | `postgres` |
| `POSTGRES_DB` | PostgreSQL database | `short_urls` |
| `REDIS_ADDR` | Redis address | `redis:6379` in Compose |
| `LOG_FILE` | Optional log file path | stdout |

## API example

Register a user:

```bash
curl -X POST http://localhost:8080/api/v1/register \
	-H "Content-Type: application/json" \
	-d '{"username":"jane","email":"jane@example.com","password":"password"}'
```

## Development

Run these commands from the repository root:

```bash
just test       # Run tests
just check      # Run tests and govulncheck
just sqlc       # Regenerate database code
just gen-docs   # Regenerate Swagger files
just build      # Build the API binary
```

## Project status

This project is under active development.