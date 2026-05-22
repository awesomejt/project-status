# Instructions for AI Coding Assistants

Start here before doing any work.

This project is an agentic AI-ready template for coding projects: frontend, backend, web, CLI, automation, library, and service work.

## Agent Priorities

1. Build the smallest correct change that advances the highest-priority unblocked task.
2. Preserve project requirements, architecture, and style decisions.
3. Prefer readable, maintainable code with focused tests over broad rewrites.
4. Keep `TODO.md`, `MEMORY.md`, docs, and `status.yaml` current.
5. Stop and mark blockers clearly when a human decision, credential, external service, or unsafe operation is required.

## Required First Reads

Before starting a task, read:

- `PROJECT_BRIEF.md`
- `MEMORY.md`
- `TODO.md`
- `status.yaml`
- `docs/Requirements.md`
- `docs/Tech-Stack.md`
- Any source, test, or config files directly relevant to the task

If the task affects structure or rollout, also read `docs/Architecture.md` and `docs/Implementation.md`.

## Root Contract

- `AGENTS.md` - agent operating instructions.
- `README.md` - human-facing overview and setup.
- `TODO.md` - task lanes, priorities, blockers, and completed work.
- `MEMORY.md` - persistent project memory, decisions, milestones, and run notes.
- `status.yaml` - shared workflow state for humans and automation.
- `PROJECT_BRIEF.md` - product goals, constraints, users, and source material.
- `AGENT_WORKFLOW.md` - recurring local-agent, cloud-agent, and review workflow.
- `QUALITY_CHECKLIST.md` - pre-review, pre-PR, and pre-release quality gate.
- `.gitmessage` - Conventional Commit template with AI attribution.

## Engineering Rules

- Check `git status` before editing.
- Pull latest changes before starting when network and permissions allow.
- Do not overwrite user changes.
- Keep edits focused and reviewable.
- Follow the project's existing stack, formatting, naming, and architecture.
- Add or update tests for meaningful behavior changes.
- Run the most relevant validation before finishing.
- Update docs when behavior, setup, deployment, or public interfaces change.
- Do not commit secrets, credentials, local env files, private keys, or generated transcripts.

## Contract And Validation Discipline

This repo is developed by multiple AI agents, including short-context local loops. Prevent drift between docs, API, web, CLI, and tests.

- Treat public contracts as shared source material: routes, request/response JSON, CLI arguments, config names, environment variables, database schema, and build outputs.
- Before changing API behavior, read the API docs, web client, CLI client, tests, and TODO items that reference that behavior.
- When changing a public contract, update all affected surfaces in the same task or leave explicit TODOs if the task is intentionally planning-only.
- Do not mark a task complete just because code was written. "Done" requires the relevant validation to pass, or a clearly documented blocker/test gap.
- Prefer small, vertical changes that keep API, web, CLI, tests, and docs aligned over broad partially validated rewrites.
- Keep one canonical path/contract in docs and TODO. Remove stale endpoint names, imagined response fields, and old task wording as soon as the contract changes.
- If validation cannot run, record why in `MEMORY.md` and keep the task open unless the user explicitly requested documentation-only work.

## Generated Code Guardrails

Local autonomous loops are useful for momentum, but they can create cross-task drift. Each agent run should actively check for it.

- Compare implementation against the current requirements before adding new functionality.
- Search for old names and paths after migrations, especially API prefixes, response fields, command names, env vars, and build output paths.
- Verify tests are testing the actual app shape, not an older or imagined interface.
- Avoid marking TODO items complete from superficial build success alone. Prefer contract tests, smoke checks, or integration checks for cross-module behavior.
- If a task touches more than one module, validate the shared behavior at the boundary, not just each module in isolation.
- Use `QUALITY_CHECKLIST.md` before any PR, real-use milestone, or handoff to a cloud reviewer.

## Status Workflow

Use `status.yaml` as the shared state file:

- `active` - work may proceed.
- `paused` - do not perform automated work.
- `blocked` - waiting on a human decision, credential, source file, or validation.
- `working` - a human or agent is actively changing the repo; other agents should skip.
- `error` - repo or automation state is unsafe; stop and request recovery.
- `stopped` - project is complete or intentionally shut down.

Automated agents should set `working` only while actively editing, and return to `active`, `blocked`, `error`, or `stopped` before ending a run.

## Task Selection

Prefer tasks in this order unless `TODO.md` says otherwise:

1. Blocker removal and requirements clarification.
2. Contract drift across API, web, CLI, docs, tests, and build outputs.
3. Failing tests, broken builds, and safety/security issues.
4. Architecture or scaffolding that unlocks later work.
5. Core implementation tasks.
6. Tests and validation gaps.
7. Documentation, deployment notes, and cleanup.

## Cloud Review Gate

Before real use, release, or deployment, schedule a cloud-based AI review/refactor pass from `TODO.md`.

- Cloud review should prioritize correctness, contract alignment, test reliability, maintainability, and production-readiness risks.
- Review findings should be added to `TODO.md` before broad refactors begin.
- Refactors should be split by module or contract boundary and validated with the root workflow once it exists.
- Do not treat local-loop generated code as production-ready until the cloud review/refactor lane and relevant validation have completed.

## Chat Logs And External Agent Logs

Full chat transcripts should not be committed. Use `chats/` only as a local transcript workspace; Markdown transcript files there are ignored by Git.

Agent workflow managers should copy or mirror transcripts and runtime logs to external storage. Hermes-compatible defaults are:

- Runtime logs: `/var/log/hermes`
- Mirrored logs: `/mnt/hermes/logs`
- Project output and transcripts: `/mnt/hermes/output/<project-name>/`

For n8n, OpenClaw, or another orchestrator, use equivalent configured storage.

## Stop Conditions

Stop and mark a task blocked if:

- Required source files or requirements are missing.
- A decision depends on Jason's preference.
- Credentials, paid services, account access, or production systems are required.
- The next action could be destructive or security-sensitive.
- External facts, product APIs, laws, pricing, or platform rules must be current and cannot be verified.
