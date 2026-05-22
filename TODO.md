# Project TODO

Task list for `Project Status`, organized by ownership and project phase.

## Needs Attention

Items here require Jason's input, a decision, credentials, external access, or manual validation before agent work can continue.

- [X] Replace all template placeholder values in project files before starting agent work.
- [ ] Confirm the exact `status_record` field set and allowed status values in `docs/Requirements.md`.
- [ ] Confirm deployment target and stage/production PostgreSQL VM hosting approach.
- [ ] Confirm secret management approach for stage and production `DATABASE_URL` values.
- [ ] Decide whether OpenAPI documentation is required for MVP or can follow initial CRUD implementation.

## Manual Validation

These items need Jason to validate on real systems, live services, devices, accounts, or deployment targets.

- [ ] Confirm requirements and success criteria in `PROJECT_BRIEF.md`.
- [ ] Confirm chosen stack and deployment target.
- [ ] Confirm credentials, API keys, and production access are not committed.
- [ ] Validate API CRUD behavior against a real PostgreSQL 18 database.
- [ ] Validate Docker Compose v2 local development startup for PostgreSQL 18, API, and web workflows.
- [ ] Validate stage `DATABASE_URL` against the stage PostgreSQL VM when available.
- [ ] Validate production `DATABASE_URL` against the production PostgreSQL VM during release readiness.
- [ ] Validate web workflows in a browser at desktop and mobile widths.
- [ ] Validate CLI workflows against a running local API.
- [ ] Validate deployment or release workflow on the target environment.

## AI Agent Work

These items are good candidates for a local model or cloud agent.

### Discovery

- [ ] Re-check current dependency versions before scaffolding if more than a week has passed since 2026-05-20.
- [ ] Confirm local tool availability for Python 3.14, Go 1.26, Node.js 24 LTS, Docker, and Docker Compose v2.
- [ ] Confirm the official PostgreSQL 18 container tag to use for local development.

### Planning

- [X] Migrate API paths from `/api/v1/status-records/*` to `/api/*`: Completed 2026-05-21 - Routes already at `/api/*` (was `/api/v1/status-records`). Renamed `api_v1` module to `api`.
- [ ] Draft the initial `/api/*` request and response contract (updated from `/api/v1/status-records`).
- [ ] Decide whether to generate an OpenAPI spec from Flask code or maintain a hand-written spec.
- [ ] Choose the production WSGI server after deployment target is known.
- [ ] Define the status record database indexes for list filters and sort order.
- [ ] Define local, test, stage, and production configuration precedence for API settings.
- [ ] Decide Compose profiles or service layout for `db`, `api`, `web`, migration, and test workflows.

### Scaffolding

- [X] Create top-level `api/`, `web/`, and `cli/` directories. Completed - directories exist.
- [X] Add root `.gitignore` entries for Python, Node, Go, database, build, and local environment artifacts. Completed - .gitignore present.
- [X] Add Docker Compose v2 support for PostgreSQL 18 local development. Completed - docker-compose.yml configured.
- [X] Add Dockerfiles for API and web if Compose-managed service containers are part of the local workflow. Completed - api/Dockerfile present.
- [X] Add Compose services or profiles for `db`, `api`, `web`, migrations, and API tests. Completed - db, api, web services configured.
- [X] Add example environment files for local, test, stage, and production without secrets. Completed - .env.example.local, .env.example.test, .env.example.stage, .env.example.production added to api/ directory.
- [X] Add root development notes for running all three parts locally with Docker Compose. Completed - docs/Development.md created with PostgreSQL 18, API, web service documentation and troubleshooting guide.
- [X] Add CI-ready validation commands once project manifests exist. Completed - API: ruff lint/format, pytest; Web: ESLint, TypeScript typecheck, Vite build.

### Implementation

- [X] Add Alembic migration baseline for PostgreSQL 18. Completed - alembic/env.py configured, 001_initial_status_records.py migration added.
- [X] Implement API configuration loading from environment variables with `DATABASE_URL` support for local, test, stage, and production. Completed - config.py with multi-environment support.
- [X] Add configuration validation that fails fast when `DATABASE_URL` is missing in API runtime contexts. Completed - Config class raises ValueError if DATABASE_URL missing.
- [X] Implement Flask application factory and versioned blueprint structure. Completed - create_app() with api_v1 blueprint.
- [X] Implement API health and database readiness endpoints. Completed - /health and /ready endpoints.
- [X] Add SQLAlchemy database setup and session lifecycle. Completed - scoped_session with create_engine in __init__.py.
- [ ] Add Alembic migration baseline for PostgreSQL 18.
- [X] Add migration command runnable through Docker Compose. Completed - migrations service added to docker-compose.yml with automatic upgrade on startup.
- [X] Implement `status_record` model and migration. Completed - StatusRecord model in models.py (auto-create via init_db).
- [X] Implement create status record endpoint. Completed - POST /api/v1/status-records (migrate to POST /api).
- [X] Implement list status records endpoint with pagination, sorting, and filters. Completed - GET /api/v1/status-records (migrate to GET /api).
- [X] Implement read status record by ID endpoint. Completed - GET /api/v1/status-records/<id> (migrate to GET /api/<id>).
- [X] Implement partial update status record endpoint. Completed - PATCH /api/v1/status-records/<id> (migrate to PATCH /api/<id>).
- [X] Implement delete status record endpoint. Completed - DELETE /api/v1/status-records/<id> (migrate to DELETE /api/<id>).
- [X] Implement JSON validation and consistent API error responses. Completed 2026-05-21 - utils.py with validate_json, field validators, make_error_response; updated create and update endpoints.
- [W] Add API endpoint documentation or OpenAPI output.
- [X] Scaffold React web application with TypeScript and Vite. Completed - web/ directory scaffolded with React 19.2.6, Vite 8.0.12, TypeScript, environment variable support.
- [X] Implement web API client and environment-based API base URL configuration. Completed - Created types/statusRecord.ts and api/client.ts with CRUD operations and environment-based VITE_API_BASE_URL configuration.
- [X] Implement web status record list view. Completed - Created StatusListView component with table display, status filtering, pagination, loading/error/empty states.
- [X] Implement web create/edit status record form. Completed - StatusForm component with create and edit modes, validation, tags management, peek mode, loading/error states, and routing.
- [X] Implement web status record detail view. Completed - StatusDetailView component with read-only display, edit button, delete confirmation flow, loading/error states, clickable rows in list view. Updated routes: /detail/:id for view, /edit/:id for form.
- [ ] Implement web delete confirmation flow.
- [ ] Implement web loading, empty, validation, and API error states.
- [ ] Check web accessibility basics for the primary workflows.
- [ ] Scaffold Go CLI module with Cobra and Viper.
- [ ] Implement CLI API client and config resolution.
- [ ] Implement `status config` command.
- [ ] Implement `status add` command.
- [ ] Implement `status list` command with filter flags.
- [ ] Implement `status show` command.
- [ ] Implement `status update` command.
- [ ] Implement `status delete` command with confirmation or force flag.
- [ ] Add CLI output formats such as table and JSON.

### Tests And Quality

- [ ] Add or update unit tests.
- [ ] Add or update integration/e2e tests where risk justifies it.
- [ ] Run lint, format check, type check, build, and test commands when available.
- [ ] Review with `QUALITY_CHECKLIST.md`.
- [ ] Add API unit tests for validation, serialization, and service behavior.
- [ ] Add API pytest configuration, fixtures, and markers.
- [ ] Add API integration tests against the Docker Compose PostgreSQL 18 container.
- [ ] Add API migration upgrade test.
- [ ] Add a Compose command or profile for running API pytest against PostgreSQL 18.
- [ ] Add web unit/component tests for list, form, detail, delete, loading, empty, and error states.
- [ ] Add web browser smoke tests for critical workflows when the dev server is available.
- [ ] Add CLI command tests with mocked HTTP responses.
- [ ] Add CLI integration smoke tests against a running API.
- [ ] Add formatting and linting configuration for API, web, and CLI.
- [ ] Run full validation before the first implementation milestone is marked complete.

### Documentation And Deployment

- [ ] Update `README.md` setup and usage instructions.
- [ ] Document deployment, environment variables, and operational notes.
- [ ] Record decisions and milestones in `MEMORY.md`.
- [ ] Document API endpoint examples.
- [ ] Document CLI install, configuration, and command examples.
- [ ] Document web local development and build workflow.
- [ ] Document database migration workflow.
- [ ] Document Docker Compose local development workflow.
- [ ] Document stage and production `DATABASE_URL` configuration for dedicated PostgreSQL VMs.
- [ ] Document deployment target, release steps, and rollback notes after Jason chooses the target.

## In Progress

Move exactly one task here while working if multiple agents may run at the same time.

- [ ]

## Blocked

Move blocked tasks here with the blocker and the next required human action.

- [ ]

## Done

Move completed items here with a brief note.

- [X] Read required project files and relevant planning docs before making changes. Completed 2026-05-20 by Codex.
- [X] Pulled latest changes before planning update. Completed 2026-05-20 by Codex.
- [X] Inventory existing repository files. Completed 2026-05-20 by Codex.
- [X] Update requirements, architecture, tech stack, implementation plan, project brief, memory, and TODO for API/web/CLI direction. Completed 2026-05-20 by Codex.
- [X] Add Docker Compose, pytest, PostgreSQL 18 container, and stage/production database URL requirements to planning docs and TODO. Completed 2026-05-20 by Codex.
