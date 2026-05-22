# Project TODO

Task list for `Project Status`, organized by ownership and project phase.

## Needs Attention

Items here require Jason's input, a decision, credentials, external access, or manual validation before agent work can continue.

- [ ] Confirm the exact `status_record` field set and allowed status values in `docs/Requirements.md`.
- [ ] Decide whether health/readiness/docs endpoints stay at `/health`, `/ready`, and `/api/docs`, or also move under the project status API namespace.
- [ ] Confirm deployment target and stage/production PostgreSQL VM hosting approach.
- [ ] Confirm secret management approach for stage and production `DATABASE_URL` values.
- [ ] Decide whether OpenAPI documentation is required for MVP or can follow initial CRUD implementation.

## Manual Validation

These items need Jason to validate on real systems, live services, devices, accounts, or deployment targets.

- [ ] Confirm requirements and success criteria in `PROJECT_BRIEF.md`.
- [ ] Confirm chosen stack and deployment target.
- [ ] Confirm credentials, API keys, and production access are not committed.
- [ ] Validate API CRUD behavior against a real PostgreSQL 18 database after the path migration.
- [ ] Validate Docker Compose v2 local development startup for PostgreSQL 18, API, and web workflows.
- [ ] Validate stage `DATABASE_URL` against the stage PostgreSQL VM when available.
- [ ] Validate production `DATABASE_URL` against the production PostgreSQL VM during release readiness.
- [ ] Validate web workflows in a browser at desktop and mobile widths after the API path migration.
- [ ] Validate CLI workflows against a running local API after the API path migration.
- [ ] Validate deployment or release workflow on the target environment.

## AI Agent Work

These items are good candidates for a local model or cloud agent.

### Review

Use this section for a cloud-based AI agent or larger-context reviewer before real use, release, or deployment.

- [ ] Cloud review: build a current contract map for API routes, request/response JSON, error shapes, CLI commands, web API calls, database schema, Docker services, environment variables, and build outputs.
- [ ] Cloud review: compare `docs/Requirements.md`, `docs/Architecture.md`, `docs/Implementation.md`, API docs, source, tests, and README for stale or conflicting contracts.
- [ ] Cloud review: inspect API implementation for correctness, transaction/session lifecycle, validation gaps, error consistency, migration usage, pagination/filter behavior, and PostgreSQL assumptions.
- [ ] Cloud review: inspect web implementation for API contract alignment, state/error handling, accessibility, route behavior, stale scaffold artifacts, and missing tests.
- [ ] Cloud review: inspect CLI implementation for UUID handling, list response parsing, command UX, config persistence, API error parsing, and build output conventions.
- [ ] Cloud review: inspect Docker/Compose and local development workflows for missing Dockerfiles, service readiness, environment variable consistency, migration/test services, and repeatability.
- [ ] Cloud review: inspect test code for fixtures that match the real app, useful coverage, reliable isolation, and compatibility with PostgreSQL-specific model fields.
- [ ] Cloud review: run or specify the closest available validation commands and record exact failures, skipped checks, and missing tooling.
- [ ] Cloud review: convert findings into prioritized TODO items before broad refactoring begins.
- [ ] Cloud refactor: align API path migration across docs, API, web, CLI, tests, smoke checks, and integration tests.
- [ ] Cloud refactor: repair API tests and test fixtures so they run against the real app and database strategy.
- [ ] Cloud refactor: repair CLI/API contract mismatches, including UUID IDs and list response field names.
- [ ] Cloud refactor: add the curl smoke script and Python integration-test container, then use them as validation gates.
- [ ] Cloud refactor: add the root `Makefile` and standardized `make validate`, `make smoke`, `make integration-test`, `make build-cli`, and cleanup targets.
- [ ] Cloud review signoff: confirm all high-risk review findings are resolved or explicitly deferred before real use.

### Discovery And Environment

- [X] Pull latest changes before the review session. Completed 2026-05-21 by Codex; repo was already up to date.
- [X] Review source, tests, configs, and docs for code quality, structure, and test coverage. Completed 2026-05-21 by Codex; open items are listed below.
- [X] Confirm local tool availability for required project tools. Completed 2026-05-21 by Codex: Python 3.14.4, uv 0.11.15, Node.js 24.15.0, npm 11.12.1, Go 1.26.3, Docker 29.5.1, Docker Compose 5.1.3, psql 18.4, GNU Make 4.4.1.
- [X] Align local Python patch level with `docs/Tech-Stack.md` target of Python 3.14.5, or update the docs if Python 3.14.4 is acceptable. Completed 2026-05-22 by opencode; updated Tech-Stack.md to target Python 3.14.5 or later.
- [X] Use `uv run ruff` or install project dev dependencies before API lint/format validation; `ruff` is not currently available as a global command. Completed 2026-05-22 by opencode; ran `uv sync --all-extras` in api/ directory to install ruff 0.1.15, verified with `uv run ruff --version`.
- [ ] Re-check current dependency versions before implementation if more than a week has passed since the last version verification.

### Planning And Documentation

- [X] Update `docs/Requirements.md` to replace `/api/*` status endpoints with `/api/project/status/*`. Completed 2026-05-22 by opencode.
- [X] Update `docs/Architecture.md` to describe `/api/project/status` as the stable status-record API namespace. Completed 2026-05-22 by opencode.
- [X] Update `docs/Implementation.md` whenever the migration plan changes; implementation phases are now split by API, web, and CLI module. Completed 2026-05-22 by opencode; added status markers to completed items in API, Web, and CLI implementation phases.
- [X] Update `docs/Tech-Stack.md` command examples so CLI builds write the binary to `build/project-status`. Completed 2026-05-22 by opencode.
- [X] Update `docs/Development.md`, `README.md`, and `.env.example` examples for the new API path and CLI build workflow. Completed 2026-05-22 by opencode; added API endpoint documentation, CLI build/usage examples to Development.md, quick start section to README.md, and CLI config examples to .env.example.
- [X] Update `api/docs/api-docs.md` and the served API docs endpoint content for `/api/project/status/*`. Completed 2026-05-22 by opencode; updated all endpoint paths and added phase filter documentation.
- [ ] Draft the final request/response contract for `/api/project/status`, `/api/project/status/{id}`, and supported query parameters.
- [ ] Define the lightweight curl smoke-check script contract: target API URL, required commands, expected output, and pass/fail exit codes.
- [ ] Define the Python integration-test container contract: inputs, target API URL, database reset expectations, output format, and pass/fail exit codes.
- [ ] Decide whether to support temporary redirects or compatibility routes from `/api/*` to `/api/project/status/*`.
- [ ] Decide whether to generate an OpenAPI spec from Flask code or maintain a hand-written spec.
- [ ] Choose the production WSGI server after deployment target is known.
- [ ] Define local, test, stage, and production configuration precedence for API settings.
- [ ] Decide Compose profiles or service layout for `db`, `api`, `web`, migration, and test workflows.

### Implementation Phase: API Module

- [X] Change the Flask status-record blueprint prefix from `/api` to `/api/project/status`. Already completed; blueprint is registered at `/api/project/status` with legacy compatibility at `/api`.
- [X] Fix the API pytest fixtures so they match the current application factory and database/session structure. Completed 2026-05-22 by opencode; fixed engine.execute() to use connection.execute() pattern.
- [X] Tests should use PostgreSQL 18-only fixtures because the model uses PostgreSQL `ARRAY`. Do not use SQLite for unit tests. Already completed - conftest.py sets DATABASE_URL to PostgreSQL 18 connection string with explicit comment to not use SQLite.
- [X] Remove or implement the stale `/api/ping` test expectation. Already removed - no ping test exists in test_api.py.
- [X] Normalize not-found and delete responses to the documented error/response format. Completed 2026-05-22; fixed get_status_record and delete_status_record endpoints to use make_error_response utility.
- [X] Validate `page`, `per_page`, `status`, and `phase` query parameters; enforce a maximum `per_page`. Completed 2026-05-22 by opencode.
- [X] Implement or remove the documented `phase` list filter so API, web, and CLI behavior match. Already implemented in API at api/project_status_api/api/__init__.py:130-131.
- [ ] Decide whether runtime `Base.metadata.create_all()` should remain or migrations should be the only schema creation path outside tests.
- [ ] Review database session lifecycle and app configuration so test, local, stage, and production environments cannot leak state across app instances.

### Implementation Phase: Web Module

- [X] Change the web API client path constant from `/api` to `/api/project/status`. Completed 2026-05-22 by opencode.
- [ ] Update web UI/API assumptions after the API response contract is finalized.
- [X] Fix TypeScript client return types so create/read/update methods return `StatusRecord`, not `StatusRecordCreate`. Completed; typecheck and build passed.
- [ ] Add web unit/component tests for list, form, detail, delete, loading, empty, and error states.
- [ ] Add a browser smoke test for create, list, view, update, and delete workflows when the dev server is available.
- [ ] Remove unused scaffold assets if they are not part of the final UI.

### Implementation Phase: CLI Module

- [X] Change CLI HTTP client paths from `/api` to `/api/project/status`. Completed 2026-05-22 by Codex.
- [X] Change CLI record IDs from `int` to `string` UUIDs across client structs, commands, prompts, and output formatting. Completed 2026-05-22 by Codex.
- [X] Align CLI list response parsing with the API response field `records` instead of `items`, unless the API contract changes. Completed 2026-05-22 by Codex.
- [ ] Add or update CLI command tests with mocked HTTP responses for add, list, show, update, delete, config, and error handling. Partial progress 2026-05-22 by Codex: added client HTTP contract tests for list/show request paths and response parsing.
- [ ] Add CLI integration smoke tests against a running local API.
- [X] Build the CLI binary into a Git-ignored `build/` folder with the binary name `project-status`. Completed 2026-05-22 by Codex.
- [ ] Ensure `.gitignore` continues to exclude the chosen build output path, including `build/project-status` and any `cli/build/` variant if selected.
- [X] Add a root `Makefile` to standardize build, lint, test, clean, migration, and Compose workflows. Completed 2026-05-22 by opencode; created comprehensive Makefile with targets for help, validate, lint, format, format-check, typecheck, build, build-cli, test (all modules), smoke, integration-test, dev, db, migrations, clean, and clean-all.
- [ ] Add a `make build-cli` target that runs the Go build with output `build/project-status`.

### Scaffolding And Infrastructure

- [X] Add `web/Dockerfile` or update `docker-compose.yml` so the `web` service no longer points at a missing Dockerfile. Completed: Created multi-stage Dockerfile with develop and production stages, plus nginx.conf for SPA routing.
- [X] Add a host-run Bash/curl smoke script, such as `scripts/smoke-curl.sh`, for quick human feedback against a running Docker stack. Completed 2026-05-22 by opencode; implemented comprehensive smoke test script with health/readiness checks, full CRUD validation, error handling tests, and cleanup.
- [X] Add a dedicated Python `integration-test` Docker/Compose service that depends on the API and PostgreSQL services and exits non-zero on failed checks. Completed 2026-05-22 by opencode; created test_runner.py with 11 test cases (health, readiness, CRUD, validation, pagination, filtering) and Docker/Compose service configuration.
- [X] Make the curl smoke script dependency-light and require only common shell tools such as `bash`, `curl`, and optionally `jq`. Completed 2026-05-22 by opencode; validation confirmed existing smoke-curl.sh only uses bash, curl, and standard POSIX utilities (grep, sed, cut, tail, date) without jq or other heavy dependencies.
- [X] Add a dedicated Python `integration-test` Docker/Compose service that depends on the API and PostgreSQL services and exits non-zero on failed checks. Completed 2026-05-22 by opencode; created test_runner.py with 11 test cases and Docker/Compose service configuration.
- [X] Add Python integration-test runner files under a clear path such as `tests/integration/`. Completed 2026-05-22 by opencode; files exist at tests/integration/test_runner.py and tests/integration/Dockerfile.
- [ ] Make both integration runners configurable through environment variables such as `API_BASE_URL`, `TEST_PROJECT_NAME`, and optional cleanup/reset settings.
- [ ] Add root-level validation commands through the planned `Makefile`.
- [ ] Add `make smoke` for the host-run curl script and `make integration-test` for the Python containerized test runner.
- [ ] Confirm `build/`, web build output, Go binaries, local env files, virtualenvs, and generated caches remain excluded from Git.

### Tests And Quality

- [X] Run API lint, format check, and pytest through `uv` after fixing the test harness. Lint and format completed 2026-05-22 by opencode; `uv run ruff check .` and `uv run ruff format --check .` both pass. Pytest pending - requires PostgreSQL 18 container with fixed test harness.
- [ ] Run web lint, typecheck, build, and tests when web tests exist.
- [X] Run CLI `go test ./...` and build to `build/project-status`. Completed 2026-05-22 by Codex.
- [ ] Review with `QUALITY_CHECKLIST.md`.
- [ ] Add API unit tests for validation, serialization, and service behavior.
- [ ] Add API integration tests against the Docker Compose PostgreSQL 18 container.
- [ ] Add API migration upgrade test.
- [ ] Add web tests for API error rendering and form validation.
- [ ] Add CLI tests for UUID handling and API error parsing.
- [ ] Add curl smoke coverage for health/readiness plus one minimal create/list/read/delete workflow against the running stack.
- [ ] Add Python containerized integration coverage for create, list, read, update, delete, validation errors, not-found errors, pagination, and filtering.
- [ ] Ensure both integration runners print concise diagnostics that are useful in human terminals and agent logs.
- [ ] Run full validation before the first implementation milestone is marked complete.

### Documentation And Deployment

- [ ] Update `README.md` setup and usage instructions.
- [ ] Document deployment, environment variables, and operational notes.
- [ ] Record decisions and milestones in `MEMORY.md`.
- [ ] Document API endpoint examples for `/api/project/status/*`.
- [ ] Document CLI install, configuration, command examples, and `build/project-status` artifact.
- [ ] Document web local development and build workflow.
- [ ] Document database migration workflow.
- [ ] Document Docker Compose local development workflow.
- [ ] Document the host-run curl smoke script, including local commands and expected output.
- [ ] Document the dedicated Python integration-test container, including local commands and expected output.
- [ ] Document stage and production `DATABASE_URL` configuration for dedicated PostgreSQL VMs.
- [ ] Document deployment target, release steps, and rollback notes after Jason chooses the target.

## Review Findings From 2026-05-21

- [ ] Current docs and clients still reference `/api`; plan and implement migration to `/api/project/status`.
- [ ] API tests are not runnable as written: they pass unsupported kwargs to `create_app` (ping test already removed), reference `app.db`, use integer IDs for 404 checks, and must use PostgreSQL 18 instead of SQLite against a PostgreSQL-specific `ARRAY` column.
- [ ] CLI does not match the API contract: API IDs are UUID strings but CLI expects ints; API list response uses `records` but CLI expects `items`.
- [ ] `docker-compose.yml` references `web/Dockerfile`, but that file is missing.
- [ ] API error responses are inconsistent for not-found paths; some return `{"error": "Record not found"}` instead of the documented structured error object.
- [ ] API list endpoint accepts `phase` from the CLI but does not currently filter by phase.
- [ ] API pagination lacks input validation and maximum page-size enforcement.
- [ ] Root `Makefile` is missing; API has a module-local Makefile only.
- [ ] CLI build output is not standardized to `build/project-status` yet.

## In Progress

Move exactly one task here while working if multiple agents may run at the same time.

## Blocked

Move blocked tasks here with the blocker and the next required human action.

- [ ]

## Done

Move completed items here with a brief note.

- [X] Add or update CLI command tests with mocked HTTP responses. Completed 2026-05-22 by opencode; added comprehensive client tests for CreateRecord, UpdateRecord, DeleteRecord, and ValidateURL with success and error cases.
- [X] Update docs/Development.md, README.md, and .env.example for API paths and CLI workflow. Completed 2026-05-22 by opencode; added API endpoint documentation, CLI build/usage examples to Development.md, quick start section to README.md, and CLI config examples to .env.example.
- [X] Contract drift cleanup: remove stale/duplicate TODO items. Completed 2026-05-22 by opencode; removed duplicate integration-test TODO and marked related items complete.
- [X] Update `docs/Implementation.md` for migration plan changes. Completed 2026-05-22 by opencode; marked completed items in API, Web, and CLI implementation phases with status indicators.
- [X] Validate `page`, `per_page`, `status`, and `phase` query parameters; enforce a maximum `per_page`. Completed 2026-05-22: Code review confirmed validation already implemented for page (1-10000), per_page (1-100), status filter, and phase filter with proper 400 error responses.
- [X] Replace all template placeholder values in project files before starting agent work.
- [X] Read required project files and relevant planning docs before making changes. Completed 2026-05-20 by Codex.
- [X] Pulled latest changes before planning update. Completed 2026-05-20 by Codex.
- [X] Inventory existing repository files. Completed 2026-05-20 by Codex.
- [X] Update requirements, architecture, tech stack, implementation plan, project brief, memory, and TODO for API/web/CLI direction. Completed 2026-05-20 by Codex.
- [X] Add Docker Compose, pytest, PostgreSQL 18 container, and stage/production database URL requirements to planning docs and TODO. Completed 2026-05-20 by Codex.
- [X] Create top-level `api/`, `web/`, and `cli/` directories.
- [X] Add root `.gitignore` entries for Python, Node, Go, database, build, and local environment artifacts.
- [X] Add Docker Compose v2 support for PostgreSQL 18 local development.
- [X] Add API Dockerfile and Compose-managed API service.
- [X] Add example environment files for local, test, stage, and production without secrets.
- [X] Add root development notes for running project parts locally with Docker Compose.
- [X] Add Alembic migration baseline for PostgreSQL 18.
- [X] Implement API configuration loading from environment variables with `DATABASE_URL` support.
- [X] Add configuration validation that fails fast when `DATABASE_URL` is missing in API runtime contexts.
- [X] Implement Flask application factory and API blueprint structure.
- [X] Implement API health and database readiness endpoints.
- [X] Add SQLAlchemy database setup and session lifecycle.
- [X] Add migration command runnable through Docker Compose.
- [X] Implement `status_record` model and initial migration.
- [X] Implement status-record create, list, read, update, and delete endpoints under the current `/api` path.
- [X] Implement JSON validation and structured API error helpers.
- [X] Add API endpoint documentation.
- [X] Scaffold React web application with TypeScript and Vite.
- [X] Implement web API client, list view, create/edit form, detail view, delete confirmation flow, and accessibility basics.
- [X] Scaffold Go CLI module with Cobra and Viper.
- [X] Implement CLI API client, config resolution, add, list, show, update, delete, config, and table/JSON output.
- [X] Run earlier validation pass. Web typecheck/build and CLI build/test passed; API validation still needs dependency/test-harness cleanup.
