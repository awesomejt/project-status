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

- Use REST endpoints under `/api/v1`.
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
