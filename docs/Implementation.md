# Implementation Plan

High-level order of implementation. Keep this aligned with `TODO.md`.

## Discovery

- Confirm the three-part product shape: Flask API, React web client, Go CLI.
- Confirm status record fields, status values, and API workflows.
- Inventory source, tests, configs, docs, and deployment assumptions.
- Identify blockers and manual validation needs, especially deployment target and production database expectations.

## Planning

- Finalize API, web, CLI, database, and test stack versions.
- Define the monorepo layout:
  - `api/` for Flask service, migrations, tests, and API docs.
  - `web/` for React application and browser tests.
  - `cli/` for Cobra/Viper client and Go tests.
  - `docker-compose.yml` for PostgreSQL 18 and Compose-managed local development support.
- Draft initial REST API contract under `/api/v1`.
- Define environment configuration for `local`, `test`, `stage`, and `production`, with `DATABASE_URL` injected per environment.
- Convert open questions into TODO tasks or blockers.

## MVP

- Scaffold `api/`, `web/`, and `cli/` projects with pinned dependency manifests.
- Add Docker Compose v2 local development support with a PostgreSQL 18 `db` service.
- Add Compose support for API and web services when useful for repeatable development startup.
- Implement API application factory, environment-based config loading, health/readiness endpoints, SQLAlchemy setup, and Alembic migration baseline.
- Ensure API config accepts `DATABASE_URL` for local, test, stage, and production without code changes.
- Implement `status_record` database model and REST CRUD endpoints:
  - `GET /api/v1/status-records`
  - `POST /api/v1/status-records`
  - `GET /api/v1/status-records/{id}`
  - `PATCH /api/v1/status-records/{id}`
  - `DELETE /api/v1/status-records/{id}`
- Add API validation, JSON error responses, filtering, sorting, and pagination.
- Implement web list, detail, create/edit form, delete confirmation, and API error states.
- Implement CLI commands:
  - `status add`
  - `status list`
  - `status show`
  - `status update`
  - `status delete`
  - `status config`
- Add focused API pytest, web, and CLI unit tests.
- Document local setup, commands, environment variables, and common workflows.

## Refinement

- Add database-backed API integration tests with the PostgreSQL 18 Compose container.
- Add CLI integration tests against a test API server where practical.
- Add web component and browser smoke coverage for critical workflows.
- Add OpenAPI documentation or equivalent endpoint reference.
- Improve validation messages, empty states, loading states, and delete safeguards.
- Add lint, format, type-check, and build commands for all three parts.
- Update README, architecture notes, and deployment notes as the implementation matures.

## Release

- Run full validation across API, web, CLI, and database migrations.
- Confirm deployment target, runtime environment, dedicated stage/production PostgreSQL VM details, and secret handling with Jason.
- Document release and rollback steps once deployment is known.
- Record release notes, remaining follow-up, and manual validation findings in `MEMORY.md`.
