# Technical Stack

Document the approved tools, languages, frameworks, libraries, and versions for this project.

## Runtime And Language

- API language/runtime: Python 3.14, targeting current patch release Python 3.14.5 or later.
- Web language/runtime: TypeScript on Node.js 24 LTS.
- CLI language/runtime: Go 1.26, targeting current patch release Go 1.26.3.
- Database runtime: PostgreSQL 18, targeting current patch release PostgreSQL 18.4.
- Local orchestration: Docker Engine with Docker Compose v2.
- API package manager: `uv` preferred for Python dependency locking and virtual environment management.
- Web package manager: `npm` unless scaffolding proves this repo already standardizes on another JavaScript package manager.
- CLI package manager: Go modules.

## Frameworks And Libraries

- API framework: Flask 3.1.x.
- API database access: SQLAlchemy 2.0.x, Alembic migrations, psycopg 3.x driver.
- API server: Flask development server locally; production WSGI server to be selected during deployment planning.
- UI framework: React 19 stable with TypeScript and Vite.
- CLI framework: Cobra v1.10.2.
- CLI configuration: Viper v1.21.0.
- Database or storage: PostgreSQL 18.
- Local development services:
  - PostgreSQL 18 container managed by Docker Compose v2.
  - Optional Compose-managed API and web services for repeatable local startup.
- Testing:
  - API: pytest, Flask test client, database-backed integration tests using the PostgreSQL 18 container.
  - Web: Vitest, React Testing Library, and Playwright when browser workflow coverage is needed.
  - CLI: Go `testing`, command tests, HTTP client tests with `httptest`.
  - Curl smoke checks: host-run Bash script using `curl` and optionally `jq` against a running local Docker stack.
  - Containerized integration: dedicated Docker/Compose Python test runner for black-box API checks.
  - End-to-end: API plus PostgreSQL integration tests, then optional web/CLI smoke tests.

## Commands

```bash
# Install dependencies after scaffolding
cd api && uv sync
cd web && npm install
cd cli && go mod download

# Start development services after scaffolding
docker compose up -d db
cd api && uv run flask --app project_status_api run
cd web && npm run dev

# Start all Compose-managed development services when available
docker compose up --build

# Run tests after scaffolding
docker compose up -d db
cd api && DATABASE_URL="${DATABASE_URL}" uv run pytest
cd web && npm test
cd cli && go test ./...

# Run quick host smoke checks against the Docker stack after the script exists
./scripts/smoke-curl.sh

# Run containerized Python integration tests after the service exists
docker compose run --rm integration-test

# Build after scaffolding
cd web && npm run build
cd cli && go build -o ../build/project-status ./...
```

## Environment

- Required API environment variables:
  - `DATABASE_URL`: PostgreSQL connection URL. Local points at the Compose `db` service or published localhost port; stage and production point at their PostgreSQL VM endpoints.
  - `APP_ENV`: environment name such as `local`, `test`, `stage`, or `production`.
  - `FLASK_ENV`: local development mode when needed.
  - `API_HOST` and `API_PORT`: optional local bind settings.
- Required web environment variables:
  - `VITE_API_BASE_URL`: base URL for the Flask API.
- Required CLI configuration:
  - `PROJECT_STATUS_API_URL`: API base URL, also configurable through Viper config files and flags.
- Local services:
  - PostgreSQL 18 database managed by Docker Compose v2 for development and tests.
  - Flask API on localhost.
  - Vite web development server on localhost.
- Stage and production:
  - Stage and production may use dedicated PostgreSQL VMs.
  - `DATABASE_URL` must be injectable per environment through deployment configuration, not hard-coded in code or committed env files.
  - Database credentials and VM connection details must remain outside Git.
- Deployment target:
  - Application deployment target is not yet selected. Keep deployment-specific tasks in TODO until Jason confirms the target.

## Version Notes

- Verified on 2026-05-20 MDT:
  - Go latest stable: 1.26.3 from `go.dev/VERSION?m=text`.
  - Python 3.14 current patch: 3.14.5 from Python.org release listings.
  - PostgreSQL 18 current patch: 18.4 from PostgreSQL release notes.
  - Cobra current module version: v1.10.2 from pkg.go.dev.
  - Viper current module version: v1.21.0 from pkg.go.dev.
  - React 19.2 is the latest documented stable React feature release; pin the latest safe React 19.x patch during web scaffolding and avoid canary/experimental React builds.
