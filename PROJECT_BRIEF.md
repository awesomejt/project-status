# Project Brief

Fill this in before asking an agent to plan or implement the project.

## Project Identity

- Project name: `Project Status`
- Short name: `status`
- Repository: "https://github.com/awesomejt/project-status"
- Project type: `full-stack`
- Primary users: agents

## Purpose

What problem does this project solve?

- Provide an API-backed project status service for agents and humans to create, inspect, update, list, and delete status records from consistent web and CLI workflows.

## Success Criteria

The project is successful when:

- MVP priority: API layer and infrastructure are production-shaped and reliable before CLI/web polish.
- A Flask API persists status records in PostgreSQL 18 and exposes documented CRUD/list endpoints.
- A React web client can perform the core status workflows through the API.
- A Go Cobra/Viper CLI can perform the same core status workflows through the API.
- Unit and integration tests cover the important API, web, CLI, and database behavior.

## Users And Workflows

- Primary user: agents
- Secondary users: humans
- Most important workflow: api/add
- Repeated or high-frequency workflow: api/list
- Admin or maintenance workflow: api/readiness + integration test

## Must Include

- api crud operations
- api listing
- cli operations (calls api)
- web client (calls api)
- Docker and Docker Compose v2 for local development
- PostgreSQL 18 container for local development
- configurable database URLs for local, stage, and production
- unit tests
- integration tests

## Out Of Scope

- authentication
- authorization
- advanced logging

## Technical Preferences

- Preferred language/runtime: Python 3.14 for API, TypeScript/Node.js 24 LTS for web, Go 1.26 for CLI.
- Preferred framework: Flask for API, React for web, Cobra plus Viper for CLI.
- Preferred package manager: `uv` for Python, `npm` for web, Go modules for CLI.
- Preferred database/storage: PostgreSQL 18.
- Local development: Docker and Docker Compose v2 should be used extensively for running PostgreSQL 18 and orchestrating API/web development services.
- Deployment target: not yet selected.
- Authentication requirements: authentication and authorization are out of scope for MVP.
- Accessibility or browser/device support: current evergreen browsers with keyboard-accessible, responsive React UI.

## Source Material

Add local paths, links, screenshots, legacy systems, API docs, tickets, research notes, or examples.

| Source | Path or URL | How to use it |
| --- | --- | --- |
| AGENTS.md | `AGENTS.md` | Follow project workflow, status, and documentation rules. |
| Planning docs | `docs/` | Keep requirements, architecture, tech stack, and implementation plan aligned. |

## Validation Needed

Add items that need Jason, a browser, a real device, account access, credentials, production data, or current external confirmation.

- Confirm the status record fields and status value set.
- Confirm deployment target and production database hosting approach.
- Confirm stage and production PostgreSQL VM connection details and secret management approach.
- Confirm whether an OpenAPI spec is required for MVP or can follow initial CRUD implementation.
- Confirm when to start local-model dogfooding after cloud-AI implementation pass reaches MVP readiness.
