# Implementation Plan

High-level order of implementation. Keep this aligned with `TODO.md`.

## Discovery

- Confirm the three-part product shape: Flask API, React web client, Go CLI.
- Confirm status record fields, status values, and API workflows.
- Inventory source, tests, configs, docs, and deployment assumptions.
- Identify blockers and manual validation needs, especially deployment target and production database expectations.
- Confirm local tool availability before implementation work starts.

## Planning

- Finalize API, web, CLI, database, and test stack versions.
- Define the monorepo layout:
  - `api/` for Flask service, migrations, tests, and API docs.
  - `web/` for React application and browser tests.
  - `cli/` for Cobra/Viper client and Go tests.
  - `docker-compose.yml` for PostgreSQL 18 and Compose-managed local development support.
- Migrate the status-record REST contract from `/api/*` to `/api/project/status/*`. [DONE]
- Add temporary compatibility routes from `/api/*` via legacy blueprint for smooth client migration. [DONE]
- Keep health/readiness endpoints separate unless Jason decides they should also move.
- Define two local integration feedback layers:
  - A host-run Bash/curl smoke script for quick human feedback against a running Docker stack. [DONE]
  - A Python-based Docker/Compose integration-test container for richer agentic assertions and diagnostics.

## Implementation Phase: API Module

- Update Flask route registration so status-record CRUD/list endpoints are served under `/api/project/status`. [DONE]
- Add legacy compatibility blueprint at `/api` that forwards to the same endpoints for smooth CLI/web migration. [DONE]
- Keep the endpoint contract explicit:
  - `GET /api/project/status`
  - `POST /api/project/status`
  - `GET /api/project/status/{id}`
  - `PATCH /api/project/status/{id}`
  - `DELETE /api/project/status/{id}`
- Update API documentation and examples to use the new namespace. [DONE]
- Fix API test fixtures so they match the current application factory and database/session lifecycle. [DONE]
- Use PostgreSQL-backed test coverage where PostgreSQL-specific model fields are required. [PENDING]
- Add validation and tests for pagination, supported filters, UUID IDs, not-found responses, and error response shape. [DONE]
- Review whether runtime schema creation should be replaced by migration-only startup outside tests.

## Implementation Phase: Web Module

- Update the shared web API client path constant to `/api/project/status`. [DONE]
- Keep all web workflows API-backed: list, detail, create, edit, delete, loading, empty, validation, and error states.
- Align TypeScript response types with the API contract. [DONE]
- Add component tests for the primary workflows and browser smoke coverage when the dev server is available.
- Update local development docs and environment examples for the new API path when needed.

## Implementation Phase: CLI Module

- Update the CLI HTTP client paths to `/api/project/status`. [DONE]
- Treat record IDs as UUID strings across command parsing, API calls, output formatting, and tests. [DONE]
- Align list response parsing with the API contract (`records` vs `items`). [DONE]
- Build the CLI binary into the ignored `build/` folder as `project-status`. [DONE]
- Add a root `Makefile` target for CLI builds, plus standard test/lint/build orchestration for the repo.
- Add command tests with mocked HTTP responses and integration smoke tests against a running local API.

## Refinement

- Add database-backed API integration tests with the PostgreSQL 18 Compose container.
- Add a host-run curl smoke script that checks health/readiness and one minimal create/list/read/delete path against the running stack.
- Add a dedicated Python integration-test Compose service that waits for the API, runs black-box HTTP workflow checks, reports concise diagnostics, and exits non-zero on failure.
- Cover create, list, read, update, delete, validation errors, not-found errors, pagination, and filtering in the Python integration-test container.
- Add root `Makefile` targets such as `make smoke` for curl checks and `make integration-test` for the Python container.
- Add CLI integration tests against a test API server where practical.
- Add web component and browser smoke coverage for critical workflows.
- Add OpenAPI documentation or equivalent endpoint reference.
- Improve validation messages, empty states, loading states, and delete safeguards.
- Add lint, format, type-check, and build commands for all three parts through a root `Makefile`.
- Update README, architecture notes, and deployment notes as the implementation matures.

## Review

Before real use, a cloud-based AI agent should perform a larger-context review and refactor pass.

- Build a current contract map for API routes, request/response JSON, error shapes, CLI commands, web API calls, database schema, Docker services, environment variables, and build artifacts.
- Compare that contract map against docs, source, tests, Docker/Compose, README, and TODO.
- Prioritize findings by risk: correctness, data integrity, test reliability, integration behavior, maintainability, security-sensitive assumptions, and production-readiness.
- Convert review findings into focused TODO items before broad refactoring begins.
- Refactor by module or contract boundary, with validation after each change.
- Complete or explicitly defer high-risk findings before treating the project as ready for real use.

## Release

- Run full validation across API, web, CLI, and database migrations.
- Confirm deployment target, runtime environment, dedicated stage/production PostgreSQL VM details, and secret handling with Jason.
- Document release and rollback steps once deployment is known.
- Record release notes, remaining follow-up, and manual validation findings in `MEMORY.md`.
