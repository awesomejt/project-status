# Project Status

This project is a full-stack project containing three parts:

* api : A simple API to handle all CRUD operations
* cli : Provide a CLI client that calls the API
* web : A web client for easy visualization

## Root Contract

AI assistants and humans should start with the root files:

- `AGENTS.md` - coding-project agent operating instructions.
- `README.md` - human-facing overview and setup.
- `TODO.md` - task lanes for discovery, planning, implementation, validation, blockers, and done items.
- `MEMORY.md` - persistent project decisions, milestones, blockers, and run notes.
- `status.yaml` - current agent workflow state.
- `PROJECT_BRIEF.md` - project goals, audience/users, constraints, stack preferences, and source material.
- `AGENT_WORKFLOW.md` - recurring local-agent workflow.
- `QUALITY_CHECKLIST.md` - engineering quality checklist before review or release.

## Quick Start

Start services:

```bash
# Start PostgreSQL and services
docker compose up -d

# Run database migrations
docker compose up migrations
```

API is available at `http://localhost:5000`:

```bash
# List status records
curl http://localhost:5000/api/project/status

# Create a record
curl -X POST http://localhost:5000/api/project/status \
  -H "Content-Type: application/json" \
  -d '{"project_name": "My Project", "short_name": "my-proj", "status": "active"}'
```

Build and use the CLI:

```bash
# Build
cd cli && go build -o ../build/project-status .

# Use
./build/project-status list
```

Run validation checks with Docker Compose-first workflows:

```bash
# Fast smoke validation
./scripts/smoke-curl.sh

# Extended integration checks in container
docker compose run --rm integration-test
```

## Intended Workflow

1. Fill out `PROJECT_BRIEF.md`.
2. Fill in `docs/Requirements.md` and `docs/Tech-Stack.md`.
3. A local or cloud agent reads `AGENTS.md`, `PROJECT_BRIEF.md`, `MEMORY.md`, `TODO.md`, and `status.yaml`.
4. The agent works the highest-priority unblocked engineering task.
5. The agent updates code, tests, docs, `TODO.md`, and `MEMORY.md`.
6. If the workflow calls for it, the agent commits a small logical change.
7. Jason reviews product decisions, security-sensitive behavior, deployments, and release readiness.

## Supporting Folders

- `docs/` - requirements, architecture, tech stack, implementation plan, and diagrams.
- `src/` - project source when the generated project has source code.
- `tests/` - project tests when applicable.
- `chats/` - optional local transcript workspace. Transcript files are ignored by Git.
- `working/` - temporary scratch files ignored by Git.
- `build/` - build artifacts ignored by Git.

## Chat Logs And Agent Output

Full chat transcripts are useful context but should not be committed by default. Keep temporary local transcripts under `chats/` if helpful. Workflow managers should copy or mirror transcripts and runtime logs to external storage.

Hermes-compatible defaults:

- Runtime logs: `/var/log/hermes`
- Mirrored logs: `/mnt/hermes/logs`
- Project output and transcripts: `/mnt/hermes/output/<project-name>/`

For n8n, OpenClaw, or another orchestrator, use equivalent configured storage.

## Git Setup

After creating a new project from this template:

```bash
git config commit.template .gitmessage
```

Use Conventional Commits and keep AI attribution in commit messages.
