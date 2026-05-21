# Project Requirements

## Purpose

`Project Status` provides a small API-backed status tracking system for coding agents and humans. It records project status entries in PostgreSQL, exposes all CRUD behavior through a Flask API, and offers both web and CLI clients that call the API instead of touching the database directly.

## Users

- Primary user: coding agents that need to create, inspect, and update project status records.
- Secondary users: humans who want a quick web or terminal view of active, blocked, and completed project work.

## Key Features

- REST API for project status record CRUD operations:
  - Create a project status record.
  - List records with pagination and filters.
  - Read one record by ID.
  - Update full or partial record details.
  - Delete a record.
- PostgreSQL-backed persistence for status records, timestamps, status history, and queryable metadata.
- Web client built with React that calls the API for list, detail, create, edit, and delete workflows.
- CLI client written in Go with Cobra commands and Viper configuration, calling the API for the same core workflows.
- Docker Compose v2 local development workflow for starting PostgreSQL 18 and project services consistently.
- Health and readiness endpoints for local development, automated checks, and deployment probes.
- Focused unit and integration test coverage for API behavior, database persistence, web API integration points, and CLI command behavior.

## Core Data Model

The first implementation should support a `status_record` resource with these fields unless later requirements override it:

- `id`: server-generated UUID.
- `project_name`: human-readable project name.
- `short_name`: short project identifier.
- `status`: one of `active`, `paused`, `blocked`, `working`, `error`, `stopped`, or `completed`.
- `phase`: optional workflow phase such as `planning`, `implementation`, `validation`, or `release`.
- `summary`: short user-facing status summary.
- `reason`: optional explanation for paused, blocked, error, or stopped states.
- `details`: optional longer notes.
- `tags`: optional string labels for filtering.
- `source`: optional origin such as `api`, `web`, `cli`, or an agent name.
- `created_at` and `updated_at`: server-managed timestamps.

## Non-Functional Requirements

- Performance: list endpoints should support pagination and common filters from the first API release.
- Security: no authentication or authorization in MVP, but inputs must be validated, SQL injection must be prevented through parameterized database access, and secrets must stay in environment variables or ignored local config files.
- Accessibility: the web client must support keyboard operation, semantic HTML, visible focus states, and readable contrast for the main workflows.
- Reliability: API writes should be transactional; migrations must be repeatable; health checks should distinguish application health from database readiness.
- Local development: Docker and Docker Compose v2 should be the primary way to run PostgreSQL 18 and should also support running API/web services when useful.
- Configuration: the API must read `DATABASE_URL` from environment-specific configuration so local, stage, and production can use different PostgreSQL endpoints without code changes.
- Environments: local development uses a PostgreSQL 18 container; stage and production may use dedicated PostgreSQL VMs reachable through stage/production database URLs.
- Compatibility: API runs on Python 3.14 and Flask; database runs on PostgreSQL 18; CLI builds with the latest stable Go toolchain; web targets current evergreen desktop and mobile browsers.

## Out Of Scope

- Authentication and authorization.
- Advanced logging, tracing, metrics, or audit trails beyond basic structured application logs.
- Direct database access from web or CLI clients.
- Multi-tenant account management.
- Offline-first client behavior.

## Acceptance Criteria

- API exposes documented JSON endpoints for create, list, read, update, and delete status records.
- API persists records in PostgreSQL 18 and includes migrations for schema creation and future changes.
- Local development can start a PostgreSQL 18 container with Docker Compose v2.
- API configuration supports local, stage, and production `DATABASE_URL` values, including external PostgreSQL VM endpoints for stage and production.
- Web client can list, create, edit, view, and delete status records through the API.
- CLI can configure an API base URL and perform add, list, show, update, and delete operations through the API.
- Unit tests and integration tests cover successful paths and important validation/error paths.
- README and docs explain local setup, environment variables, development commands, and validation commands.
