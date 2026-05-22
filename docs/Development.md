# Local Development Guide

This guide explains how to run the Project Status services locally using Docker Compose v2.

## Prerequisites

- Docker Engine
- Docker Compose v2 (Compose Plugin)
- `uv` (optional, for running API outside of Docker)
- Node.js 24 LTS (optional, for running web outside of Docker)
- Go 1.26 (optional, for running CLI outside of Docker)

## Quick Start

### Start PostgreSQL 18 Container

```bash
docker compose up -d db
```

This starts the PostgreSQL 18 database on port `5432`.

### Run API Tests Against PostgreSQL 18

```bash
docker compose up -d db
cd api && DATABASE_URL="postgresql://project_status:project_status_dev@localhost:5432/project_status" uv run pytest
```

### View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f db
docker compose logs -f api
```

### Stop Services

```bash
# Stop all services
docker compose down

# Stop and remove volumes (resets database)
docker compose down -v
```

## Service Details

### PostgreSQL 18 (db)

| Configuration | Value |
|---------------|-------|
| Container Name | `project-status-db` |
| Port | `5432:5432` |
| Database | `project_status` |
| User | `project_status` |
| Password | `project_status_dev` |

Connection string:
```bash
postgresql://project_status:project_status_dev@localhost:5432/project_status
```

### Flask API (api)

| Configuration | Value |
|---------------|-------|
| Container Name | `project-status-api` |
| Port | `5000:5000` |
| Depends On | `db` (healthy) |

**Note:** The API service depends on the database being healthy before starting.

### React Web (web)

| Configuration | Value |
|---------------|-------|
| Container Name | `project-status-web` |
| Port | `3000:3000` |
| Depends On | `api` |

**Note:** Web scaffolding is in progress. See `TODO.md` for status.

### Migrations (migrations)

| Configuration | Value |
|---------------|-------|
| Container Name | `project-status-migrations` |
| Depends On | `db` (healthy) |

Run database migrations:

```bash
docker compose up migrations
```

## Environment Variables

### API

Copy the example environment file and configure for your environment:

```bash
# Local development
cp api/.env.example.local api/.env.local
```

Required environment variables:
- `DATABASE_URL`: PostgreSQL connection URL
- `APP_ENV`: Environment name (`local`, `test`, `stage`, `production`)
- `FLASK_ENV`: Flask environment (e.g., `development`)

### Web

Configured via `VITE_API_BASE_URL`:
```
VITE_API_BASE_URL=http://localhost:5000
```

## Port Allocation

| Service | Port |
|---------|------|
| PostgreSQL | 5432 |
| Flask API | 5000 |
| React Web | 3000 |

## Development Workflows

### Run API Outside Docker (Hot Reload)

```bash
# Start database
docker compose up -d db

# Set environment variables
export DATABASE_URL="postgresql://project_status:project_status_dev@localhost:5432/project_status"
export APP_ENV="local"
export FLASK_ENV="development"

# Run API
cd api && uv run flask --app project_status_api run
```

### Run API Tests

```bash
docker compose up -d db
export DATABASE_URL="postgresql://project_status:project_status_dev@localhost:5432/project_status"
cd api && uv run pytest
```

### Run Migrations

```bash
# Via Docker Compose
docker compose up migrations

# Directly (if API is running)
cd api && uv run alembic -c alembic.ini upgrade head
```

### Reset Database

```bash
# Stop services and remove volumes
docker compose down -v

# Restart (clean database)
docker compose up -d db
```

## Troubleshooting

### Database Connection Issues

```bash
# Check if database is running
docker compose ps

# Check database logs
docker compose logs db

# Verify connection
docker compose exec db psql -U project_status -d project_status -c "SELECT 1;"
```

### API Not Starting

```bash
# Check API logs
docker compose logs api

# Verify database is healthy
docker compose ps db
```

### Port Already in Use

```bash
# Find what's using the port
lsof -i :5432  # PostgreSQL
lsof -i :5000  # API
lsof -i :3000  # Web

# Or change ports in docker-compose.yml
```

## API Endpoints

All status record operations use the `/api/project/status/*` namespace:

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/project/status` | List records (supports pagination and filters) |
| `POST` | `/api/project/status` | Create a new record |
| `GET` | `/api/project/status/{id}` | Read a specific record |
| `PATCH` | `/api/project/status/{id}` | Update a record |
| `DELETE` | `/api/project/status/{id}` | Delete a record |

Query parameters for list endpoint:
- `page` - Page number (min: 1, max: 10000)
- `per_page` - Records per page (min: 1, max: 100)
- `status` - Filter by status value
- `phase` - Filter by phase value

## CLI Build And Usage

Build the CLI:

```bash
# Build using Go directly
cd cli && go build -o ../build/project-status .

# Or use the Makefile
make build-cli
```

The binary is created at `build/project-status`.

CLI commands:

```bash
# Show configuration
./build/project-status config

# Set API URL
./build/project-status config set api-url http://localhost:5000

# List records
./build/project-status list

# Add a new record
./build/project-status add --project-name "My Project" --short-name "my-proj" --status active

# Show a specific record
./build/project-status show <id>

# Update a record
./build/project-status update <id> --status blocked --reason "Waiting on X"

# Delete a record
./build/project-status delete <id>
```
