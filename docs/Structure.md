# Repository Structure

This document captures the intentional project layout for MVP and beyond.

## MVP-Critical Paths

- `api/` Flask API + DB access + migrations/tests
- `web/` React client consuming API
- `cli/` Go CLI consuming API
- `docker-compose.yml` Compose-first local orchestration
- `scripts/` smoke and helper scripts
- `tests/integration/` containerized API integration checks

## Future-Ready Scaffolding

- `contracts/` API/public contract artifacts (`openapi/` placeholder)
- `deploy/` stage/production deployment scaffolding
- `ops/` runbooks and observability placeholders
- `infra/compose/profiles/` compose profile layering placeholders
- `tests/e2e/` reserved for post-MVP e2e coverage
- `docs/adr/` architecture decision records

## Visual Layout

```mermaid
flowchart TD
    Root[project-status repo]

    Root --> MVP[MVP Runtime]
    MVP --> API[api/]
    MVP --> WEB[web/]
    MVP --> CLI[cli/]
    MVP --> COMPOSE[docker-compose.yml]
    MVP --> SCRIPTS[scripts/]
    MVP --> INTEGRATION[tests/integration/]

    Root --> FUTURE[Future-Ready Scaffolding]
    FUTURE --> CONTRACTS[contracts/]
    CONTRACTS --> OPENAPI[contracts/openapi/]
    FUTURE --> DEPLOY[deploy/]
    DEPLOY --> STAGE[deploy/stage/]
    DEPLOY --> PROD[deploy/production/]
    FUTURE --> OPS[ops/]
    OPS --> RUNBOOKS[ops/runbooks/]
    OPS --> OBS[ops/observability/]
    FUTURE --> INFRA[infra/compose/profiles/]
    FUTURE --> E2E[tests/e2e/]
    FUTURE --> ADR[docs/adr/]
```
