# Project Memory

Persistent project memory for `Project Status`.

Agents should update this file after meaningful decisions, milestones, blockers, research findings, or implementation runs.

Keep this file concise and durable. Do not paste full chat transcripts here; store temporary transcripts under `chats/` and mirror workflow-manager logs to external storage.

## Current Status

- Current phase: planning complete for API, web, and CLI MVP.
- Last major milestone: documented requirements, architecture, tech stack, implementation phases, Docker Compose local development, pytest API testing, and TODO for three-part system.
- Next recommended task: scaffold the monorepo structure and local PostgreSQL 18 development environment.
- Current blocker: none blocking agent work, but Jason should confirm status fields, deployment target, stage/production database secrets, and OpenAPI timing.

## Key Decisions

- JSON validation framework uses a centralized `validate_json` function with pluggable field validators.
- Error responses follow a consistent format: `{"error": {"code": <http_code>, "message": "<human_readable>", "details": <optional>}}`.
- Project has three top-level parts: `api/`, `web/`, and `cli/`.
- API is the system of record and the only component that talks to PostgreSQL.
- Web and CLI clients call the API for all behavior.
- API implementation stack: Python 3.14, Flask, SQLAlchemy, Alembic, psycopg, PostgreSQL 18.
- Web implementation stack: React 19 stable, TypeScript, Vite.
- CLI implementation stack: Go 1.26, Cobra, Viper.
- Docker and Docker Compose v2 are the default local development orchestration tools.
- API testing uses pytest.
- Local development uses a PostgreSQL 18 container.
- Stage and production may use dedicated PostgreSQL VMs, selected by environment-specific `DATABASE_URL` values.
- Authentication, authorization, and advanced logging remain out of scope for MVP.

## Architecture Notes

- Use REST endpoints under `/api` (migrated from `/api/v1`).
- First resource is `status_record`, supporting create, list, read, update, and delete.
- List endpoints should support pagination and common filters from the first API release.
- Database migrations should be introduced with the first schema commit.
- The API must own database access and read its target database from `DATABASE_URL`.
- Docker Compose should support at least the PostgreSQL 18 container, with API/web Compose services added if useful for repeatable local workflows.

## Technical Notes

- Verified on 2026-05-20 MDT: Go 1.26.3, Python 3.14.5, PostgreSQL 18.4, Cobra v1.10.2, and Viper v1.21.0 from official release/package sources.
- React 19.2 is the latest documented stable React feature release; pin latest safe React 19.x patch during scaffolding and avoid canary/experimental builds.
- Keep local, test, stage, and production configuration environment-driven. Do not commit actual stage or production database URLs.

## Manual Validation Findings

Record findings from real systems, live services, browser/device testing, deployment targets, or Jason's checks.

- None yet.

## Open Questions

- Does Jason approve the initial `status_record` fields and status value set in `docs/Requirements.md`?
- What is the deployment target and production PostgreSQL hosting approach?
- How will stage and production `DATABASE_URL` secrets be provided to the deployed API?
- Should OpenAPI documentation be required for MVP?

## Blockers

- None.

## Agent Run Log

Newest entries first.

### 2026-05-21 19:25 - opencode

- Task: Add API endpoint documentation or OpenAPI output.
- Files changed: `api/docs/api-docs.md` (new), `api/project_status_api/api/__init__.py`.
- Validation: Python syntax verified.
- Result: Created comprehensive API documentation at `api/docs/api-docs.md` covering all endpoints (health, ready, CRUD for status records) with request/response examples, error formats, and parameter descriptions. Added `/api/docs` endpoint to serve the documentation as markdown.
- Commit: pending.
- Blockers or follow-up: none.

### 2026-05-21 19:00 - opencode

- Task: Scaffold Go CLI module with Cobra and Viper.
- Files changed: `cli/go.mod`, `cli/go.sum`, `cli/main.go`, `cli/cmd/root.go`, `cli/cmd/config.go`, `cli/cmd/add.go`, `cli/cmd/list.go`, `cli/cmd/show.go`, `cli/cmd/update.go`, `cli/cmd/delete.go`, `cli/internal/client/client.go`.
- Validation: Go build successful, CLI help command works.
- Result: Created Go CLI with Cobra v1.10.2 and Viper v1.21.0. Implemented commands: `config` (show/set), `add`, `list`, `show`, `update`, `delete`. API client in `internal/client/` handles HTTP requests. Supports `--api-url` flag for API base URL, `--output` flag for table/json format. Config uses Viper with environment variable support (`PROJECT_STATUS_API_URL`, `PROJECT_STATUS_OUTPUT`).
- Commit: pending.
- Blockers or follow-up: none.

### 2026-05-21 18:45 - opencode

- Task: Implement web status record detail view.
- Files changed: `web/src/components/StatusDetailView.tsx` (new), `web/src/App.tsx`, `web/src/components/StatusListView.tsx`, `web/src/api/client.ts`.
- Validation: TypeScript typecheck passed, Vite build successful (258kb bundle).
- Result: Created StatusDetailView component with read-only display of all status record fields, formatted metadata grid (status, phase, source, timestamps), edit button navigating to `/edit/:id`, delete confirmation flow with modal-style prompt. Updated StatusListView to make rows clickable (navigate to detail view). Refactored routes: `/detail/:id` for read-only view, `/edit/:id` for edit form.
- Commit: pending.
- Blockers or follow-up: none.

### 2026-05-21 17:35 - opencode

- Task: Migrate API paths from `/api/v1/status-records/*` to `/api/*`: Rename `api_v1` module to `api`.
- Files changed: `api/project_status_api/__init__.py`, `api/project_status_api/api_v1/` → `api/project_status_api/api/`.
- Validation: Code review confirmed correct module rename and blueprint naming. Routes already at `/api/*` (no migration needed for paths).
- Result: Renamed `api_v1` module to `api`. Updated import in `__init__.py` from `from . import api_v1` to `from . import api`. Updated blueprint name from `api_v1` to `api`.
- Commit: pending.
- Blockers or follow-up: none.

### 2026-05-21 03:35 - opencode

- Task: Implement web status record list view.
- Files changed: `web/src/components/StatusListView.tsx` (new), `web/src/App.tsx`.
- Validation: TypeScript typecheck passed.
- Result: Created StatusListView component with table display, status filtering (all 7 status values), pagination controls, loading spinner, error state display, empty state message, and status badges with color coding. Replaced default Vite/React welcome page with the status list view.
- Commit: pending.
- Blockers or follow-up: none.

### 2026-05-21 03:25 - opencode

- Task: Implement web API client and environment-based API base URL configuration.
- Files changed: `web/src/types/statusRecord.ts` (new), `web/src/api/client.ts` (new).
- Validation: TypeScript typecheck passed.
- Result: Created TypeScript types for StatusRecord, StatusRecordCreate, StatusRecordUpdate, StatusRecordListResponse, and ApiError. Implemented apiClient with getRecords, getRecord, createRecord, updateRecord, and deleteRecord methods. Uses VITE_API_BASE_URL from environment (defaults to http://localhost:5000).
- Commit: pending.
- Blockers or follow-up: none.

### 2026-05-21 03:15 - opencode

- Task: Add CI-ready validation commands for API and web.
- Files changed: `api/pyproject.toml`, `api/Makefile`, `web/package.json`.
- Validation: API lint (ruff check), format (ruff format), web lint (eslint), typecheck (tsc --noEmit), build (vite build) all pass.
- Result: Added ruff for linting/formatting, pytest configuration, Makefile with targets. Web already had ESLint; added typecheck script. All validation commands verified working.
- Commit: pending.
- Blockers or follow-up: none.

### 2026-05-21 03:05 - opencode

- Task: Scaffold React web application with TypeScript and Vite.
- Files changed: `web/` directory (new), `.gitignore`, `.env.example`, `TODO.md`, `status.yaml`.
- Validation: Build verified successful - TypeScript compiles, Vite builds production bundle.
- Result: Scaffolded web/ directory with React 19.2.6, Vite 8.0.12, TypeScript. Configured environment variable support for VITE_API_BASE_URL. Added example environment file with API placeholder.
- Commit: `e5f17d2` - feat(web): scaffold React web application with TypeScript and Vite.
- Blockers or follow-up: none.

### 2026-05-21 02:50 - opencode

- Task: Add root development notes for running all three parts locally with Docker Compose.
- Files changed: `docs/Development.md` (new), `TODO.md`, `status.yaml`.
- Validation: File created with comprehensive documentation covering prerequisites, quick start, service details, environment variables, port allocation, development workflows, and troubleshooting.
- Result: Added `docs/Development.md` with local development guide for PostgreSQL 18, Flask API, and React web services via Docker Compose v2. Includes commands for starting services, running tests, migrations, and common development workflows.
- Commit: pending.
- Blockers or follow-up: none.

### 2026-05-21 02:35 - opencode

- Task: Add example environment files for local, test, stage, and production without secrets.
- Files changed: `.gitignore`, `api/.env.example.local`, `api/.env.example.test`, `api/.env.example.stage`, `api/.env.example.production`, `TODO.md`, `status.yaml`.
- Validation: Files are not ignored by git, contain placeholder values only (no secrets).
- Result: Added four example environment files with placeholder DATABASE_URL values for each environment (local, test, stage, production). Updated .gitignore to track example files while ignoring actual .env files with secrets.
- Commit: pending.
- Blockers or follow-up: none.

### 2026-05-21 02:25 - opencode

### 2026-05-21 02:15 - opencode

- Task: Add Alembic migration baseline for PostgreSQL 18.
- Files changed: `api/alembic/env.py` (configured), `api/alembic/versions/001_initial_status_records.py` (new migration).
- Validation: Python syntax checks passed, Alembic history shows migration correctly.
- Result: Configured Alembic to use project models and DATABASE_URL. Added initial migration creating `status_records` table with indexes for `short_name`, `status`, `phase`, and `created_at`.
- Commit: `bad3526` - feat(api): add Alembic migration baseline for PostgreSQL 18.
- Blockers or follow-up: none.

### 2026-05-21 02:00 - opencode

- Task: Implement JSON validation and consistent API error responses.
- Files changed: `api/project_status_api/utils.py` (new), `api/project_status_api/api_v1/__init__.py`.
- Validation: Python syntax check passed for both files.
- Result: Added centralized validation utilities (`validate_json`, `validate_status`, `validate_string`, `validate_optional_string`, `validate_tags`, `make_error_response`). Updated `create_status_record` and `update_status_record` endpoints to use the new validation framework with consistent error responses.
- Commit: `b5472a2` - feat(api): add JSON validation and consistent error responses.
- Blockers or follow-up: none.

### 2026-05-20 23:57 - Codex

- Task: Add Docker Compose, pytest, PostgreSQL 18 container, and stage/production database URL requirements.
- Files changed: `PROJECT_BRIEF.md`, `MEMORY.md`, `TODO.md`, `docs/Requirements.md`, `docs/Architecture.md`, `docs/Tech-Stack.md`, `docs/Implementation.md`, `status.yaml`.
- Validation: pending final diff and whitespace checks for this docs-only update.
- Result: planning docs now require Docker Compose v2 local development, PostgreSQL 18 container use, pytest for API tests, and environment-driven `DATABASE_URL` for local/test/stage/production.
- Blockers or follow-up: confirm deployment target and secret management for dedicated PostgreSQL VMs.

### 2026-05-20 23:49 - Codex

- Task: Update project requirements, architecture, tech stack, implementation plan, and TODO for API, web, and CLI.
- Files changed: `PROJECT_BRIEF.md`, `MEMORY.md`, `TODO.md`, `docs/Requirements.md`, `docs/Architecture.md`, `docs/Tech-Stack.md`, `docs/Implementation.md`, `status.yaml`.
- Validation: reviewed docs and repository inventory; no code tests exist yet.
- Result: planning docs now describe a Flask/PostgreSQL API, React web client, and Go Cobra/Viper CLI.
- Blockers or follow-up: confirm status fields, deployment target, and OpenAPI timing.
